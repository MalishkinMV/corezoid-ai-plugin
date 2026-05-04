# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Repository Is

A documentation and resource library for **Corezoid** — a visual process engine where business logic is defined as JSON graphs of connected nodes. It serves as the knowledge base for AI-assisted process generation.

**Key directories:**
- `knowledge/` — single source of truth: auth, process schema, node types, gotchas, quick-reference, validation checklist
- `agents/` — tool-agnostic workflow definitions (10 agents covering all workflows)
- `mcp-server/` — MCP server wrapping the Python client as 35+ AI-callable tools
- `playbooks/` — step-by-step guides for specific tasks
- `samples/` / `templates/` — working process examples and reusable skeletons
- `docs/` — detailed per-node docs, error codes, API ops map, APIGW and dashboard guides

## Authentication — Read This First

Corezoid uses **path-based SHA-1 auth**, NOT HTTP header auth:

```
URL:  {BASE_URL}/{API_LOGIN}/{timestamp}/{signature}
sig:  SHA1(timestamp + SECRET + request_body_compact_json + SECRET)
```

**Required credentials (from Account Settings → API in Corezoid UI):**

| Variable | Description |
|----------|-------------|
| `API_LOGIN` | Numeric user ID |
| `SECRET` | API secret key — a long static string. **NOT the JWT session token from the browser.** |
| `BASE_URL` | Must end with `/api/2/json` — e.g. `https://admin.corezoid.com/api/2/json` |
| `COMPANY_ID` | Workspace ID with `i` prefix — e.g. `i469b4ea8...` |

Copy `.env.example` → `.env` and fill it in before running any scripts.

## Python Client: corezoid_client.py

`corezoid_client.py` in the repo root is the **primary tool** for all Corezoid API operations. Use it instead of raw HTTP calls — it handles signing automatically.

```python
from corezoid_client import CorezoidClient
client = CorezoidClient(api_login='YOUR_API_LOGIN', secret='...', base_url='...', company_id='i469b4...')

# List folder contents
client.list_folder_contents(YOUR_FOLDER_ID)

# Get process scheme
result = client.get_process_scheme(12345)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']  # scheme is a LIST — index [0]

# Create and check a task
r = client.create_task(12345, {"city": "Kyiv"}, ref="test-1")
task_id = r['ops'][0]['obj_id']
client.show_task(12345, task_id=task_id)
```

**Always batch ops** — the API accepts multiple operations in one `ops` array. Never make two calls when one batched call works.

## MCP Server

`mcp-server/` wraps the Python client as MCP tools for AI agents. Setup:

```bash
cd mcp-server && uv sync
```

Configure in your AI tool (Cursor, Claude Desktop, Windsurf, Kiro) — see `docs/mcp/README.md`.
Full tool catalog: 35+ tools covering processes, nodes, tasks, aliases, analysis, and raw API.

## Corezoid Process Architecture

A **process** is a JSON object containing a `scheme.nodes` array. Each node has:
- `id` — 24-character hex string, unique within the process
- `uuid` — UUID v4, globally unique
- `obj_type` — `1` (Start), `0` (Logic/intermediate), `2` (End/Final/Error), `3` (Escalation/Reply)
- `condition.logics` — array of logic entries defining behavior and routing
- `condition.semaphors` — timeout/rate-limit controls
- `position` — `[x, y]` array (NOT separate `x`/`y` fields in API operations)

**Routing is exclusively through `go` logic entries:**
```json
{ "type": "go", "to_node_id": "<target_node_id>" }
```

### Node Types and Their Logic Types

| Node | `obj_type` | Key logic `type` |
|------|-----------|-----------------|
| Start | 1 | — |
| Code | 0 | `api_code` (JS/Erlang) |
| API Call | 0 | `api` |
| Call a Process (RPC) | 0 | `api_rpc` |
| Call a Process (async) | 0 | `api_copy` |
| Set Parameters | 0 | `set_param` |
| Condition | 0 | `go_if_const` or `condition` |
| Reply to Process | 3 | `api_rpc_reply` |
| Final / Error End | 2 | — |

### Standard Connector Pattern

```
Start → API Call → [Set Param/Code] → Reply Success → Final
              ↘ Reply Error (API) ↘ Escalation timeout → Error Final
```

Every functional node needs `err_node_id` pointing to an escalation (obj_type 3) node.
Every API call node needs a `semaphors` timeout pointing to a timeout escalation node.

### API Call Node Requirements

```json
{
  "type": "api",
  "is_migrate": true,
  "rfc_format": true,
  "format": "",
  "content_type": "application/json",
  "method": "GET",
  "url": "https://example.com/endpoint",
  "extra": {},
  "extra_type": {},
  "extra_headers": {},
  "cert_pem": "",
  "max_threads": 5,
  "send_sys": true,
  "debug_info": false,
  "err_node_id": "<escalation-node-id>",
  "customize_response": true,
  "response": { "header": "{{header}}", "response": "{{body}}" },
  "response_type": { "header": "object", "response": "object" },
  "version": 2
}
```

**Response variable access:** body is stored under the **second key** name (`response` in this example) — access as `{{response.field}}`, NOT `{{body.field}}`.

### Process JSON Format (for import/export)

For UI import or JSON-based deployment, the root must be a **JSON array**:

```json
[
  {
    "obj_type": 1,
    "obj_id": 1234567,
    "parent_id": 0,
    "title": "Process Name",
    "description": "",
    "status": "active",
    "params": [],
    "ref_mask": true,
    "conv_type": "process",
    "uuid": "<uuid-v4>",
    "scheme": {
      "nodes": [...],
      "web_settings": [[], []]
    }
  }
]
```

Every node needs both `id` (24-char hex) and `uuid` (UUID v4). Bare objects (not wrapped in array) are rejected on import.

### API-Based Process Creation (4 steps)

1. Create process (get `conv_id`)
2. Create nodes — batch Start + Final together (get node `obj_id`s)
3. Modify nodes — set positions, logics, connections (batch all nodes)
4. Commit

```python
# Step 1
result = client.make_request({'ops': [{'company_id': COMPANY_ID, 'conv_type': 'process',
    'create_mode': 'without_nodes', 'folder_id': FOLDER_ID, 'obj': 'conv', 'obj_type': 0,
    'status': 'active', 'title': 'My Process', 'type': 'create'}]})
conv_id = result['ops'][0]['obj_id']

# Step 2 — get version first
version = client.get_process_details(conv_id)['ops'][0]['change_time']
import uuid
result = client.make_request({'ops': [
    {'company_id': COMPANY_ID, 'conv_id': conv_id, 'id': str(uuid.uuid4()),
     'obj': 'node', 'obj_type': 1, 'title': 'Start', 'type': 'create', 'version': version},
    {'company_id': COMPANY_ID, 'conv_id': conv_id, 'id': str(uuid.uuid4()),
     'obj': 'node', 'obj_type': 2, 'title': 'Final', 'type': 'create', 'version': version}
]})
start_id = result['ops'][0]['obj_id']
final_id = result['ops'][1]['obj_id']

# Step 3 — modify with full payloads
client.make_request({'ops': [
    {'company_id': COMPANY_ID, 'conv_id': conv_id, 'obj': 'node', 'obj_id': start_id,
     'obj_type': 1, 'type': 'modify', 'version': version, 'title': 'Start', 'description': '',
     'extra': {'icon': '', 'modeForm': 'collapse'}, 'logics': [{'type': 'go', 'to_node_id': final_id}],
     'semaphors': [], 'position': [1000, 100]},
    {'company_id': COMPANY_ID, 'conv_id': conv_id, 'obj': 'node', 'obj_id': final_id,
     'obj_type': 2, 'type': 'modify', 'version': version, 'title': 'Final', 'description': '',
     'extra': {'icon': 'success', 'modeForm': 'collapse'}, 'logics': [],
     'semaphors': [], 'options': {'save_task': True}, 'position': [1000, 400]}
]})

# Step 4
client.make_request({'ops': [{'company_id': COMPANY_ID, 'conv_id': conv_id,
    'obj': 'commit', 'type': 'confirm', 'version': version}]})
```

## Knowledge Files (read these first)

| Purpose | Path |
|---------|------|
| Authentication & credentials | `knowledge/auth.md` |
| Process & node JSON structure | `knowledge/process-schema.md` |
| All logic types with examples | `knowledge/node-types.md` |
| 15 silent failures to avoid | `knowledge/gotchas.md` |
| Copy-paste snippets + API response structure | `knowledge/quick-reference.md` |
| Pre-deployment validation checklist | `knowledge/validation-checklist.md` |

## Agent Workflow Definitions

All agents live in `agents/`. See `agents/README.md` for the full list.

**Process design and building** (for process designers):

| Agent | File | When to use |
|-------|------|------------|
| architect | `agents/architect.md` | Design multi-process systems and choose patterns |
| process-builder | `agents/process-builder.md` | Build new connector or logic process |
| process-reviewer | `agents/process-reviewer.md` | Audit before deployment |
| process-updater | `agents/process-updater.md` | Modify or debug existing process |
| process-tech-writer | `agents/process-tech-writer.md` | Generate docs and enrich JSON |
| apigw-manager | `agents/apigw-manager.md` | Expose as HTTP endpoint via API Gateway |
| api-wrapper-manager | `agents/api-wrapper-manager.md` | Alternative HTTP wrapper via convctl |
| dashboard-manager | `agents/dashboard-manager.md` | Create monitoring dashboards |


## Detailed Documentation

| Purpose | Path |
|---------|------|
| Error code reference | `docs/error_explanations.md` |
| MCP server setup | `docs/mcp/README.md` |
| Node-specific docs | `docs/nodes/<type>-node.md` |
| API ops type map | `docs/api-ops-map.md` |
| Dashboard guide | `docs/dashboards/README.md` |
| APIGW guide | `docs/apigw/README.md` |

## Playbooks (step-by-step for specific scenarios)

`playbooks/create-connector.md`, `create-logic.md`, `create-openapi-spec.md`, `review-process.md`, `fix-connector.md`, `update-logic.md`

## Workflow for any process task

1. Read relevant `knowledge/` file for the task type
2. Check `agents/<task>.md` for the workflow
3. Fetch existing processes: `client.list_folder_contents(folder_id)`
4. Deploy via 4-step API create pattern (see `knowledge/quick-reference.md`)
5. Test: `client.create_task(proc_id, data)` → `client.show_task(proc_id, task_id=...)`

## Critical Rules (common sources of silent failures)

### Code node logic type is `api_code`, not `code`

The correct logic type name is `"type": "api_code"`. Using `"type": "code"` returns `Unexpected logic type` error immediately.

```json
{"type": "api_code", "lang": "js", "src": "try { ... } catch(e) { throw new Error(e.message); }", "err_node_id": "..."}
```

### `set_param` vs `api_code` nodes
`set_param` runs natively with zero overhead. `api_code` nodes spin up V8/BEAM. Always prefer `set_param` for: variable copy/rename, string concatenation (`"{{a}}_{{b}}"`), arithmetic (`$.math({{a}}+{{b}})`).

**`set_param` cannot extract array elements from API responses.** `{{geo.results.0.latitude}}` returns an empty string when `results` is a JSON array. Use `api_code` (JS) for any array access:
```js
data.lat = String(data.geo.results[0].latitude);
data.lon = String(data.geo.results[0].longitude);
```

Only use `api_code` for: array access/manipulation, regex, JSON parsing, complex conditionals.

### Commit version — always use creation `change_time`

When committing a process draft, always pass the `change_time` from when the process was **first created** (or the last successful commit), NOT the `change_time` after node modifications. Using the post-modify `change_time` returns `"This version hasn't any commit"`.

```python
# Get version ONCE after process creation — save it
version = client.get_process_details(conv_id)['ops'][0]['change_time']
# Use this same version for all node creates, modifies, AND the commit
client.make_request({'ops': [{'company_id': COMPANY_ID, 'conv_id': conv_id,
    'obj': 'commit', 'type': 'confirm', 'version': version}]})
```

### Identity and format
- **`company_id`** format: UUID without `i` prefix — e.g. `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` (from your workspace URL)
- **`BASE_URL`**: `https://www.corezoid.com/api/2/json` (NOT `admin.corezoid.com`)
- **`REF`** must be unique within a process for non-terminating task flows — duplicates silently drop tasks
- **Node `id`**: 24-char hex string. **Node `uuid`**: UUID v4. Both required and unique.
- **Process JSON for import**: root must be a JSON array `[{...}]`, never a bare object

### Copy Task vs Call a Process
| Pattern | Logic type | Behavior |
|---------|-----------|---------|
| **async fan-out** | `api_copy` | Fire-and-forget — sends task and continues immediately |
| **wait for result** | `api_rpc` | Synchronous — blocks until target replies |

### Node modify — always send full payload
When modifying any node via API, always include: `title`, `description`, `logics`, `semaphors`, `extra`, `options`, `position`. Partial updates silently reset fields to defaults.

### `extra` field format differs by context
- In **API create/modify ops**: `extra` is a **dict** `{"modeForm": "collapse", "icon": ""}`
- In **process JSON for import**: `extra` is an **escaped JSON string** `"{\"modeForm\":\"collapse\"}"`

### Dashboard chart types
Valid chart types: `column`, `pie`, `funnel`, `table`. The type `bar` does NOT exist.
Grid dimensions: use `width`/`height` (NOT `w`/`h`).
Chart `modify` requires `obj_type` + full `series` — partial modify is not supported.

### Process layout discipline
- Top-down layout; vertical grid step: 160px between nodes
- Main flow: x ≈ 1000 | Timeout handlers: x ≈ 740–800 | Error handlers: x ≈ 1150–1264
- Collapse Start and End nodes (`modeForm: "collapse"`); expand logic nodes (`"expand"`)
- No overlapping nodes

### Semaphors on API call nodes
Every `api` logic node MUST have a semaphor:
```json
{"type": "time", "value": 30, "dimension": "sec", "to_node_id": "<timeout-escalation-id>"}
```
Without it, tasks hang forever if the external API is unresponsive.

### Error reply nodes
Must include `exception_reason` when `throw_exception: true`:
```json
{"type": "api_rpc_reply", "throw_exception": true, "exception_reason": "API call failed",
 "res_data": {"description": "..."}, "res_data_type": {"description": "string"}, "mode": "key_value"}
```

## Onboarding and AI Tool Integration

**New to this repo?** Read `GETTING-STARTED.md` — credentials, first commands, workflow selection.

**Tool configs** are auto-generated from `knowledge/` and `agents/` by `scripts/generate-tool-configs.py`:

| Tool | Config file |
|------|------------|
| Claude Code | `CLAUDE.md` (this file — manually maintained) |
| Cursor | `.cursor/rules/corezoid.mdc` (generated) |
| Windsurf | `.windsurf/rules/corezoid.md` (generated) |
| GitHub Copilot | `.github/copilot-instructions.md` (generated) |
| Kiro | `.kiro/steering/corezoid.md` (generated) |
| Devin | `.devin.md` (generated) |

To regenerate after updating `knowledge/` or `agents/`: `python scripts/generate-tool-configs.py`

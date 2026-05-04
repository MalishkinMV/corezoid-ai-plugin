---
name: corezoid-process-builder
description: >
  Builds Corezoid process JSON from scratch or from a template — both connectors (wrapping
  an external API) and logic processes (orchestrating other processes via RPC calls).
  Activate whenever a user asks to create a connector, build a process, make a process that
  calls an API, implement business logic in Corezoid, or provides credentials alongside a
  process description. This skill knows all mandatory node fields and the exact JSON structure
  needed — use it before generating any process JSON to avoid structural errors that block deployment.
---

# Corezoid Process Builder


## Bundled References

This Codex plugin includes the upstream Corezoid AI docs repository at `../../assets/source/`. When this skill mentions paths from the original repo, resolve them under that directory unless the user has provided project-local versions. Common examples include `knowledge/`, `docs/`, `templates/`, `playbooks/`, `json-schema/`, `samples/`, `mcp-server/`, and `convctl.sh`.

## CRITICAL: Always use MCP tools — never write Python scripts

The MCP server (`mcp__corezoid__*`) is the correct interface for all Corezoid API operations.
**Never write Python scripts, never use `corezoid_client.py` directly.**

| Operation | Use this MCP tool |
|-----------|------------------|
| List folder contents | `mcp__corezoid__list_folder_contents` |
| Get process structure | `mcp__corezoid__get_process_scheme` |
| Any create/modify/commit API op | `mcp__corezoid__make_request` |
| Create and send a task | `mcp__corezoid__create_task` |
| Check task result | `mcp__corezoid__show_task` |
| Modify a single node | `mcp__corezoid__modify_node` |
| Batch modify nodes | `mcp__corezoid__batch_modify_nodes` |
| Commit changes | `mcp__corezoid__commit` |

---

## Required parameters (ask if missing)

| Parameter | Purpose |
|-----------|---------|
| `COMPANY_ID` | Workspace ID (from `.env`) |
| `FOLDER_ID` | Folder where the process will be deployed |

## Step 1 — Identify the pattern

**Connector** — wraps an external HTTP API:
```
Start → API Call → [Code: extract fields] → Reply Success → Final
              ↘ Reply Error (API) → Error
              ↘ Reply Timeout    → Error Timeout
```

**Logic** — orchestrates other Corezoid processes:
```
Start → Call Process (api_rpc) → Reply Success → Final
                  ↘ Reply Error (RPC) → Error
```

## Step 2 — Check existing folder context

```
mcp__corezoid__list_folder_contents(folder_id=FOLDER_ID)
```

Understand what's already there — reuse existing process IDs, discover aliases.

## Step 3 — Build process via MCP (4 steps)

### 3.1 Create empty process

```
mcp__corezoid__make_request(ops=[{
    "company_id": COMPANY_ID,
    "conv_type": "process",
    "create_mode": "without_nodes",
    "folder_id": FOLDER_ID,
    "obj": "conv",
    "obj_type": 0,
    "status": "active",
    "title": "Process Name",
    "description": "What this process does",
    "type": "create"
}])
```

Save `conv_id = result['ops'][0]['obj_id']`.

### 3.2 Get version + create all nodes in one batch

Get `version` (= `change_time`) immediately after creation — use this same value for ALL subsequent operations and the commit:

```
mcp__corezoid__get_process_details(process_id=conv_id)
→ save version = result['ops'][0]['change_time']
```

Create all nodes in a single batched call. Generate unique 24-char hex IDs for each node:

```
mcp__corezoid__make_request(ops=[
    {"company_id": COMPANY_ID, "conv_id": conv_id, "id": "<uuid>",
     "obj": "node", "obj_type": 1, "title": "Start", "type": "create", "version": version},
    {"company_id": COMPANY_ID, "conv_id": conv_id, "id": "<uuid>",
     "obj": "node", "obj_type": 0, "title": "API Call", "type": "create", "version": version},
    ... (all nodes in one call)
])
```

Save the hex `obj_id` returned for each node — these are the IDs used for connections.

### 3.3 Modify all nodes (positions + logics)

Use `mcp__corezoid__batch_modify_nodes` or `mcp__corezoid__make_request` with all modify ops batched. Send the **complete payload** for every node — partial updates silently reset fields:

```
mcp__corezoid__make_request(ops=[
    # Start node
    {"company_id": COMPANY_ID, "conv_id": conv_id, "obj": "node", "type": "modify",
     "obj_id": "<start_hex_id>", "obj_type": 1, "version": version,
     "title": "Start", "description": "Process entry point",
     "extra": {"icon": "", "modeForm": "collapse"},
     "logics": [{"type": "go", "to_node_id": "<api_call_hex_id>"}],
     "semaphors": [], "position": [1000, 100]},

    # API Call node
    {"company_id": COMPANY_ID, "conv_id": conv_id, "obj": "node", "type": "modify",
     "obj_id": "<api_call_hex_id>", "obj_type": 0, "version": version,
     "title": "API Call", "description": "Calls external API",
     "extra": {"icon": "", "modeForm": "expand"},
     "logics": [
         {"type": "api", "is_migrate": true, "rfc_format": true, "format": "",
          "content_type": "application/json", "method": "GET",
          "url": "https://example.com/endpoint",
          "extra": {}, "extra_type": {}, "extra_headers": {},
          "cert_pem": "", "max_threads": 5, "send_sys": true, "debug_info": false,
          "err_node_id": "<reply_error_hex_id>", "customize_response": true,
          "response": {"header": "{{header}}", "response": "{{body}}"},
          "response_type": {"header": "object", "response": "object"}, "version": 2},
         {"type": "go", "to_node_id": "<next_hex_id>"}
     ],
     "semaphors": [{"type": "time", "value": 30, "dimension": "sec",
                    "to_node_id": "<reply_timeout_hex_id>"}],
     "position": [1000, 280]},

    ... (all remaining nodes)
])
```

### 3.4 Commit

```
mcp__corezoid__commit(process_id=conv_id, version=version)
```

Or via make_request:
```
mcp__corezoid__make_request(ops=[{
    "company_id": COMPANY_ID, "conv_id": conv_id,
    "obj": "commit", "type": "confirm", "version": version
}])
```

## Step 4 — Test

```
mcp__corezoid__create_task(process_id=conv_id, data={"input_param": "value"}, ref="test-001")
→ save task_id = result['ops'][0]['obj_id']

mcp__corezoid__show_task(process_id=conv_id, task_id=task_id)
→ check result['ops'][0]['data'] for output fields
```

The test must pass before the work is done.

---

## Node skeletons

### Start node (`obj_type: 1`)
```json
{
  "obj_type": 1,
  "title": "Start",
  "description": "Process entry point",
  "extra": {"icon": "", "modeForm": "collapse"},
  "logics": [{"type": "go", "to_node_id": "<next_id>"}],
  "semaphors": [],
  "position": [1000, 100]
}
```

### API Call node (`obj_type: 0`)
```json
{
  "obj_type": 0,
  "title": "Call External API",
  "description": "Calls external API and stores response",
  "extra": {"icon": "", "modeForm": "expand"},
  "logics": [
    {
      "type": "api", "is_migrate": true, "rfc_format": true, "format": "",
      "content_type": "application/json", "method": "GET",
      "url": "https://api.example.com/endpoint?param={{input}}",
      "extra": {}, "extra_type": {}, "extra_headers": {},
      "cert_pem": "", "max_threads": 5, "send_sys": true, "debug_info": false,
      "err_node_id": "<reply_error_id>", "customize_response": true,
      "response": {"header": "{{header}}", "response": "{{body}}"},
      "response_type": {"header": "object", "response": "object"}, "version": 2
    },
    {"type": "go", "to_node_id": "<next_id>"}
  ],
  "semaphors": [{"type": "time", "value": 30, "dimension": "sec", "to_node_id": "<reply_timeout_id>"}],
  "position": [1000, 280]
}
```

### Code node (`obj_type: 0`, logic type `api_code`)
```json
{
  "obj_type": 0,
  "title": "Extract Fields",
  "description": "Extracts fields from API response",
  "extra": {"icon": "", "modeForm": "expand"},
  "logics": [
    {
      "type": "api_code", "lang": "js",
      "src": "try {\n  data.field = String(data.response.field);\n} catch(e) { throw new Error(e.message); }",
      "err_node_id": "<reply_error_id>"
    },
    {"type": "go", "to_node_id": "<next_id>"}
  ],
  "semaphors": [],
  "position": [1000, 440]
}
```

### Call Process node (`obj_type: 0`, logic type `api_rpc`)
```json
{
  "obj_type": 0,
  "title": "Call Process Name",
  "description": "Calls the X process and waits for reply",
  "extra": {"icon": "", "modeForm": "expand"},
  "logics": [
    {
      "type": "api_rpc", "conv_id": 1234567,
      "err_node_id": "<reply_error_id>",
      "extra": {"param1": "{{value1}}"}, "extra_type": {"param1": "string"},
      "group": "all"
    },
    {"type": "go", "to_node_id": "<next_id>"}
  ],
  "semaphors": [{"type": "time", "value": 30, "dimension": "sec", "to_node_id": "<reply_timeout_id>"}],
  "position": [1000, 280]
}
```

### Reply Success node (`obj_type: 3`)
```json
{
  "obj_type": 3,
  "title": "Reply Success",
  "description": "Returns result to caller",
  "extra": {"icon": "", "modeForm": "expand"},
  "logics": [
    {
      "type": "api_rpc_reply", "mode": "key_value",
      "res_data": {"field": "{{field}}"},
      "res_data_type": {"field": "string"},
      "throw_exception": false
    },
    {"type": "go", "to_node_id": "<final_id>"}
  ],
  "semaphors": [],
  "position": [1000, 600]
}
```

### Reply Error node (`obj_type: 3`)
```json
{
  "obj_type": 3,
  "title": "Reply Error — API",
  "description": "Returns error when API call fails",
  "extra": {"icon": "", "modeForm": "expand"},
  "logics": [
    {
      "type": "api_rpc_reply", "mode": "key_value",
      "res_data": {"description": "API call failed"},
      "res_data_type": {"description": "string"},
      "throw_exception": true, "exception_reason": "API call failed"
    },
    {"type": "go", "to_node_id": "<error_id>"}
  ],
  "semaphors": [],
  "position": [1264, 280]
}
```

### Final node (`obj_type: 2`)
```json
{
  "obj_type": 2,
  "title": "Final",
  "description": "Task completed successfully",
  "extra": {"icon": "success", "modeForm": "collapse"},
  "logics": [], "semaphors": [],
  "options": {"save_task": true},
  "position": [1000, 760]
}
```

### Error end node (`obj_type: 2`)
```json
{
  "obj_type": 2,
  "title": "Error",
  "description": "Error end node",
  "extra": {"icon": "error", "modeForm": "collapse"},
  "logics": [], "semaphors": [],
  "position": [1264, 380]
}
```

---

## Key rules

### `api_code` not `code`
`{"type": "code", ...}` → `Unexpected logic type`. Always use `"type": "api_code"`.

### `set_param` vs `api_code`
Prefer `set_param` for copy/rename/concat/arithmetic — it runs natively with zero overhead.
Use `api_code` only when you need array access, regex, or JSON parsing:
```js
// set_param cannot do this — use api_code:
data.lat = String(data.geo.results[0].latitude);
```

### Commit version
Save `version = change_time` once after process creation. Use the **same** version for all node creates, modifies, and the commit. Never re-fetch `change_time` after modifying nodes.

### Response field access
The response body is stored under the **second key name** in `response`:
```json
"response": {"header": "{{header}}", "response": "{{body}}"}
```
Access as `{{response.field}}` — not `{{body.field}}`.

### api_rpc vs api_copy
| Need | Logic type | Behavior |
|------|-----------|---------|
| Wait for result | `api_rpc` | Blocks until target replies |
| Fire and forget | `api_copy` | Continues immediately, no result |

### Layout discipline
- Main flow: x ≈ 1000 | Timeout handlers: x ≈ 740 | Error handlers: x ≈ 1264
- Vertical step: 160px between nodes
- Collapse Start/End nodes (`modeForm: "collapse"`); expand logic nodes (`"expand"`)

---

## Common mistakes checklist

- [ ] Using MCP tools — not Python scripts
- [ ] All nodes have full payload in modify (title, description, logics, semaphors, extra, position)
- [ ] API call nodes have `err_node_id` + timeout semaphor
- [ ] Error replies have `throw_exception: true` + `exception_reason`
- [ ] Success replies have `throw_exception: false`
- [ ] Final node has `options: {"save_task": true}`
- [ ] `api_rpc` (not `api_copy`) when waiting for result
- [ ] `set_param` preferred over `api_code` for simple operations
- [ ] No hardcoded URLs/secrets — use `{{env_var[@name]}}`
- [ ] REF is unique for non-terminating task flows
- [ ] Commit uses same version as node creation

## Reference

- Node & logic reference: `knowledge/node-types.md`
- Validation checklist: `knowledge/validation-checklist.md`
- Agent workflow definition: `agents/process-builder.md`
- Working example: `samples/stripe-checkout.json`

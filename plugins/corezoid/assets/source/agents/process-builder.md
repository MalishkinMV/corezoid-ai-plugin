# Agent: process-builder

## Purpose

Build a new Corezoid process from scratch — either a **connector** (wraps an external API) or a **logic** (orchestrates other Corezoid processes). Always produces working, validated process JSON deployed to the target workspace.

## When to invoke

- User describes a process that calls an external API, webhook, or data source
- User wants to compose existing Corezoid processes into a new logic
- User provides API credentials/URL and asks to "create a connector"
- User describes business logic ("route this", "fetch that", "when X happens do Y")

## Required inputs

| Parameter | How to get |
|-----------|-----------|
| `API_LOGIN` | Corezoid → Account Settings → API |
| `SECRET` | Corezoid → Account Settings → API (NOT the browser JWT token) |
| `BASE_URL` | e.g. `https://admin.corezoid.com/api/2/json` |
| `COMPANY_ID` | `i` + workspace ID from URL (e.g. `i469b4ea8...`) |
| `FOLDER_ID` | From project URL: `.../stage/<FOLDER_ID>` |
| Process description | What the process does, what inputs/outputs it has |

## MCP tools used

```
get_process_details      — inspect existing processes before building
list_folder_contents     — discover what exists in the target folder
create_process           — step 1: create the process shell
create_node              — step 2: add each node
modify_node              — step 3: configure each node with full payload
commit_process_draft     — step 4: publish draft
create_task              — test: send a task to the deployed process
show_task                — test: verify task result
list_task_history        — debug: trace task path through nodes
analyze_process_structure — review: validate the built process
find_nodes_missing_semaphors — review: check for missing timeouts
```

## Workflow

### Step 1 — Choose pattern

**Connector** (calls external HTTP API):
```
Start → Code (prepare params) → API Call → Reply Success → Final (save_task)
                                          ↘ Reply Error (API) → Error
              ↘ Reply Error (Code) → Error
```

**Logic** (calls other Corezoid processes):
```
Start → Code (optional transform) → Call Process → Reply Success → Final (save_task)
                                                  ↘ Reply Error (Call) → Error
              ↘ Reply Error (Code) → Error
```

### Step 2 — Generate unique IDs

Every node needs two IDs:
```python
import uuid, secrets
node_id = secrets.token_hex(12)   # 24-char hex → used in to_node_id, err_node_id
node_uuid = str(uuid.uuid4())     # UUID v4 → used in uuid field
```

### Step 3 — Build using the 4-step API pattern

```python
from corezoid_client import CorezoidClient

client = CorezoidClient(api_login=..., secret=..., base_url=..., company_id=...)

# Step 1: Create process shell
proc = client.create_process(title="My Process", folder_id=FOLDER_ID)
process_id = proc['ops'][0]['obj_id']

# Step 2: Create nodes (one call per node)
node = client.create_node(process_id=process_id, obj_type=0)
node_id = node['ops'][0]['obj_id']

# Step 3: Modify each node with FULL payload (partial resets fields to null)
client.modify_node(process_id=process_id, node_id=node_id, payload={
    "title": "Prepare Parameters",
    "description": "Builds API request body",
    "logics": [...],
    "semaphors": [],
    "extra": {"modeForm": "expand", "icon": ""},
    "options": {},
    "position": [500, 280]
})

# Step 4: Publish
client.commit_process_draft(process_id=process_id)
```

### Step 4 — Test before done

```python
r = client.create_task(process_id, data={"city": "Kyiv"}, ref="test-001")
task_id = r['ops'][0]['obj_id']
result = client.show_task(process_id, task_id=task_id)
# Check result['ops'][0]['task']['data'] for output
```

The test must return the expected output before the work is considered complete.

## Knowledge references

- `knowledge/auth.md` — correct auth pattern and credentials
- `knowledge/process-schema.md` — complete node/process JSON structure
- `knowledge/node-types.md` — all logic type schemas with examples
- `knowledge/gotchas.md` — 15 silent failures to avoid

## Critical rules (from knowledge/gotchas.md)

1. `err_node_id` is mandatory on every logic node — missing it causes unhandled crashes
2. Every API call node needs a semaphor — without it tasks hang forever
3. `api_rpc` blocks (use when you need the result); `api_copy` fires and forgets
4. Error reply nodes must have `throw_exception: true`
5. Final success node must have `options: {"save_task": true}`
6. Response body is under the second key name in `response` mapping — not under `body`
7. Node modify always sends the FULL payload — partial updates silently reset fields
8. Use `set_param` not `code` for variable copy, concatenation, arithmetic

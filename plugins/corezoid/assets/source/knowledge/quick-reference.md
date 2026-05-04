# Quick Reference — Common Corezoid Operations

Copy-paste snippets for the most frequent tasks. For full context on any item, see the linked knowledge file.

---

## Auth setup

```python
from corezoid_client import CorezoidClient
from dotenv import load_dotenv
import os

load_dotenv()
client = CorezoidClient(
    api_login=os.environ['API_LOGIN'],
    secret=os.environ['SECRET'],
    base_url=os.environ['BASE_URL'],
    company_id=os.environ['COMPANY_ID']
)
```

→ Details: `knowledge/auth.md`

---

## Generate node IDs

```python
import uuid, secrets

node_id = secrets.token_hex(12)   # "67f2a77682ba966c7f688d95" — 24-char hex
node_uuid = str(uuid.uuid4())     # "9e3d314c-29ef-4b0d-96d1-21e81e355797" — UUID v4
```

→ Details: `knowledge/gotchas.md` #5

---

## 4-step API process creation

```python
# 1. Create process shell
proc = client.create_process(title="My Process", folder_id=FOLDER_ID)
process_id = proc['ops'][0]['obj_id']

# 2. Create a node
node_resp = client.create_node(process_id=process_id, obj_type=0)
node_api_id = node_resp['ops'][0]['obj_id']

# 3. Modify with FULL payload (partial resets fields to null)
client.modify_node(process_id=process_id, node_id=node_api_id, payload={
    "title": "API Call",
    "description": "Calls external service",
    "logics": [
        {
            "type": "api",
            "is_migrate": True,
            "rfc_format": True,
            "format": "raw",
            "raw_body": "{{requestBody}}",
            "content_type": "application/json",
            "method": "POST",
            "url": "{{env_var[@service-url]}}",
            "extra": {}, "extra_type": {},
            "extra_headers": {"authorization": "Bearer {{token}}"},
            "cert_pem": "", "max_threads": 5,
            "send_sys": False, "debug_info": False,
            "err_node_id": reply_error_id,
            "customize_response": True,
            "response": {"body": "{{body}}", "header": "{{header}}"},
            "response_type": {"body": "object", "header": "object"},
            "version": 2
        },
        {"type": "go", "to_node_id": reply_success_id}
    ],
    "semaphors": [
        {"type": "time", "value": 30, "dimension": "sec", "to_node_id": timeout_error_id}
    ],
    "extra": {"modeForm": "expand", "icon": ""},
    "options": {},
    "position": [500, 460]
})

# 4. Commit
client.commit_process_draft(process_id)
```

→ Details: `knowledge/node-types.md`

---

## Fetch and navigate process nodes

```python
result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']  # LIST — always index [0]

# Find a node by type
start_node = next(n for n in nodes if n['obj_type'] == 1)
end_nodes = [n for n in nodes if n['obj_type'] == 2]
```

→ Details: `knowledge/gotchas.md` #10

---

## Test a process

```python
r = client.create_task(process_id, data={"input_field": "value"}, ref="test-001")
task_id = r['ops'][0]['obj_id']

result = client.show_task(process_id, task_id=task_id)
task_data = result['ops'][0]['task']['data']
print(task_data)

# If task is stuck, trace its path:
history = client.list_task_history(process_id, task_id)
```

---

## Analyze a process

```python
client.analyze_process_structure(process_id)
client.find_nodes_missing_semaphors(process_id)
client.find_hardcoded_values(process_id)
client.find_orphaned_nodes(process_id)
```

---

## Raw batched API call

```python
result = client.make_request({'ops': [
    {'type': 'show', 'obj': 'conv', 'obj_id': 123, 'company_id': COMPANY_ID},
    {'type': 'show', 'obj': 'conv', 'obj_id': 456, 'company_id': COMPANY_ID}
]})
# Always check: result['request_proc'] == 'ok'
# Per op: result['ops'][N]['proc'] == 'ok'
```

→ Details: `knowledge/auth.md`

---

## set_param logic (prefer over code nodes)

```python
{
    "type": "set_param",
    "extra": {
        "full_name": "{{first}}_{{last}}",
        "total": "$.math({{price}}*{{qty}})",
        "timestamp": "$.date()"
    },
    "extra_type": {
        "full_name": "string",
        "total": "number",
        "timestamp": "string"
    },
    "err_node_id": reply_error_id
}
```

→ Details: `knowledge/node-types.md`, `knowledge/gotchas.md` #12

---

## go_if_const conditions (NOT on logic object)

```python
# ❌ Wrong
logic.get('param')  # None

# ✅ Correct
for cond in logic['conditions']:
    param = cond['param']   # e.g. "status"
    const = cond['const']   # e.g. "active"
    fun   = cond['fun']     # e.g. "eq"
```

→ Details: `knowledge/gotchas.md` #13

---

## Error reply with exception

```python
{
    "type": "api_rpc_reply",
    "mode": "key_value",
    "res_data": {"description": "API call failed: {{error_detail}}"},
    "res_data_type": {"description": "string"},
    "throw_exception": True,
    "exception_reason": "API call failed"  # required when throw_exception=true
}
```

→ Details: `knowledge/gotchas.md` #14

---

## API response structure reference

### `get_process_scheme` path to nodes

```python
result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']
#                               ^--- LIST, must index [0]
#                                              ^--- inner 'scheme' is a DICT
```

### `get_process_details` — what's where

| Field | `get_process_details` | `get_process_scheme` |
|---|---|---|
| `conv_type` | ❌ absent | ✅ `scheme[0]['conv_type']` |
| `description`, `params` | ❌ absent | ✅ in `scheme[0]` |
| `change_time` (version) | ✅ `ops[0]['change_time']` | ❌ absent |
| `project_id`, `stage_id` | ✅ `ops[0]` | ❌ absent |

**Rule:** use `get_process_details` only for `change_time`/`project_id`/`stage_id`. Use `get_process_scheme` for everything else.

### Node `condition` field

```python
condition = node['condition']       # dict, NOT a list
logics    = condition['logics']     # list of logic entries
semaphors = condition['semaphors']  # list of semaphor entries
node_extra = node['extra']          # escaped JSON STRING, not a dict
```

---

## Large API response handling

**Never embed API responses as inline strings in Python or shell.** Schemes can be 10K–200K+ chars.

```python
# ✅ Always: fetch → save to file → load from file
scheme = client.get_process_scheme(process_id)
with open(f'process_{process_id}.json', 'w') as f:
    json.dump(scheme, f)

with open(f'process_{process_id}.json') as f:
    scheme = json.load(f)
nodes = scheme['ops'][0]['scheme'][0]['scheme']['nodes']
```

---

## Orphan and no-op detection

```python
# Find nodes unreachable from Start
result = client.find_orphaned_nodes(process_id)
# → reachable_count, orphaned_count, orphaned_nodes: [{id, title, obj_type}]

# Find functionally useless nodes
result = client.find_noop_nodes(process_id)
# → noop_conditions (all branches route to same node)
# → unused_set_params (variable never referenced downstream)

# Check if a variable is actually used
result = client.check_variable_usage(process_id, ["user_id", "channel"])
# → {variables: {name: {referenced, reference_count, references}}}
```

MCP tools: `find_orphaned_nodes`, `find_noop_nodes`, `check_variable_usage`

---

## Deleting a process — requires explicit confirmation

⚠️ Irreversible. Always confirm with the user before executing:

1. Fetch process titles via `get_process_details`
2. Show the list (title + ID) to the user
3. Ask user to type the exact title(s) to confirm
4. Only proceed if typed titles match exactly

```python
client.make_request({'ops': [{'type': 'delete', 'obj': 'conv', 'obj_id': process_id}]})
```

---

## Task operations reference

```python
# Create task
r = client.create_task(process_id, data={"key": "val"}, ref="unique-ref")
task_id = r['ops'][0]['obj_id']

# Show task (by task_id or ref)
result = client.show_task(process_id, task_id=task_id)
# or: client.show_task(process_id, ref="unique-ref")
task_data = result['ops'][0]['task']['data']

# Trace execution path (task_id only — ref not supported for history)
history = client.list_task_history(process_id, task_id)
nodes_visited = history['ops'][0]['list']  # [{node_id, node_prev_id, create_time_ms}]

# Modify task data and continue flow
client.modify_task(process_id, data={"field": "new_val"}, task_id=task_id)

# Delete task (auto-resolves node_id via show_task)
client.delete_task(process_id, task_id=task_id)
```

History notes:
- Per-node data snapshots are NOT available in history — use `show_task` for current data
- `ref` must be unique per process — duplicate ref submissions are rejected

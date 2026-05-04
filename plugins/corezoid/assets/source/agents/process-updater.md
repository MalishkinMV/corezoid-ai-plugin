# Agent: process-updater

## Purpose

Modify an existing Corezoid process — add or change nodes, fix broken logic, update parameters, debug failed tasks. Covers both planned updates and emergency fixes.

## When to invoke

- User says "update", "modify", "add", "change", "fix" an existing process
- User provides a failed task ID or `ResultNodeObjID` + `ResultTaskData`
- User wants to add a new API call or process call to existing logic
- User wants to change routing conditions or variable mappings

## Required inputs

| Input | When needed |
|-------|------------|
| `API_LOGIN`, `SECRET`, `BASE_URL`, `COMPANY_ID` | Always |
| `PROCESS_ID` | Always — to fetch and modify |
| `ResultNodeObjID` | For debugging — which node the task stopped at |
| `ResultTaskData` | For debugging — task data at failure point |
| Description of change | What to add/change/fix |

## MCP tools used

```
get_process_scheme       — fetch current process state
get_process_details      — get version/metadata
analyze_process_structure — understand current structure before changing
list_task_history        — trace failed task path through nodes
show_task                — inspect task data at any node
modify_node              — update node configuration (FULL payload required)
create_node              — add a new node
commit_process_draft     — publish changes
create_task              — re-test after fix
```

## Update workflow

### Step 1 — Fetch current state

```python
result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']
# Inspect nodes to understand current structure
```

### Step 2 — Plan changes

Before modifying, identify:
- Which nodes are affected
- Whether new nodes are needed (and generate their IDs)
- Which `to_node_id` / `err_node_id` references will change

### Step 3 — Modify nodes

**CRITICAL: Always send the complete payload.** Partial payload silently resets omitted fields to null.

```python
# Fetch the current node payload first
current = client.get_process_scheme(process_id)
node = next(n for n in nodes if n['id'] == target_node_id)

# Modify what you need, keep everything else
client.modify_node(process_id=process_id, node_id=node_id, payload={
    "title": node['title'],           # keep existing
    "description": "Updated description",  # change
    "logics": [new_logic, go_logic],  # full logics array
    "semaphors": node.get('semaphors', []),
    "extra": {"modeForm": "expand", "icon": ""},
    "options": node.get('options', {}),
    "position": [node['x'], node['y']]
})
```

### Step 4 — Commit and test

```python
client.commit_process_draft(process_id)

r = client.create_task(process_id, data={"test_field": "value"}, ref="fix-test-001")
task_id = r['ops'][0]['obj_id']
result = client.show_task(process_id, task_id=task_id)
```

## Debug workflow (for failed tasks)

### Step 1 — Understand where it failed

```python
# Get the task data at the failure node
history = client.list_task_history(process_id, task_id)
# Find the last node the task visited
```

### Step 2 — Read error tags in ResultTaskData

| Error tag | Meaning | Fix |
|-----------|---------|-----|
| `__conveyor_api_return_type_error__` | API returned unexpected format | Check `customize_response`, `response_type` |
| `__conveyor_api_return_code__` | API returned non-2xx status | Check URL, auth headers, request body |
| `__conveyor_code_return_type_error__` | Code node JS error | Check `try/catch`, variable access |
| `api_connection_error` | Network unreachable | Check URL, firewall, env vars |
| `api_bad_answer` | Unparseable API response | Check `content_type`, `format` |
| `api_timeout` | Exceeded semaphor time limit | Increase semaphor value or fix slow API |
| `api_validation_error` | Required field missing or wrong type | Check `extra_type` definitions |
| `api_code_syntax_error` | JS syntax error | Fix code node syntax |
| `api_code_runtime_error` | JS runtime exception | Add `try/catch`, check null references |

### Step 3 — Fix the node

After identifying the issue, apply the fix using `modify_node` with full payload, commit, and re-test.

## Adding a Call Process node

```python
call_process_logic = {
    "type": "api_rpc",           # NOT "call_process"
    "conv_id": target_process_id,
    "err_node_id": reply_error_node_id,
    "extra": {
        "param1": "{{value1}}",
        "param2": "{{value2}}"
    },
    "extra_type": {
        "param1": "string",
        "param2": "number"
    },
    "group": "all"               # wait for reply
}
```

Use `api_copy` instead of `api_rpc` only when you don't need the result back.

## Knowledge references

- `knowledge/gotchas.md` — especially #7 (missing semaphor), #8 (api_rpc vs api_copy), #11 (partial modify resets fields)
- `knowledge/node-types.md` — correct logic type schemas
- `knowledge/auth.md` — credentials and auth pattern

# Corezoid Process JSON Schema

## Root structure (for UI import / JSON-based deployment)

The root **must be a JSON array** — bare objects are rejected:

```json
[
  {
    "obj_type": 1,
    "obj_id": 1234567,
    "parent_id": 0,
    "title": "Process Name",
    "description": "What this process does",
    "status": "active",
    "params": [],
    "ref_mask": true,
    "conv_type": "process",
    "uuid": "25b6543f-99b9-4658-95e3-9ef44df05151",
    "scheme": {
      "nodes": [],
      "web_settings": [[], []]
    }
  }
]
```

## Required top-level fields

| Field | Type | Rules |
|-------|------|-------|
| `obj_type` | integer | Always `1` for a process |
| `obj_id` | integer | Unique positive integer within the file |
| `parent_id` | integer | `0` for root-level import |
| `title` | string | Max ~255 chars |
| `description` | string | Can be empty string |
| `status` | string | Use `"active"` |
| `params` | array | Input/output parameter declarations |
| `ref_mask` | boolean | Always `true` |
| `conv_type` | string | Always `"process"` |
| `uuid` | string | UUID v4, globally unique |
| `scheme` | object | Contains `nodes` and `web_settings` |

`web_settings` is always `[[], []]`.

## Node structure

```json
{
  "id": "67f2a77682ba966c7f688d95",
  "obj_type": 0,
  "condition": {
    "logics": [],
    "semaphors": []
  },
  "title": "Node Title",
  "description": "What this node does",
  "x": 900,
  "y": 400,
  "uuid": "9e3d314c-29ef-4b0d-96d1-21e81e355797",
  "extra": "{\"modeForm\":\"expand\",\"icon\":\"\"}",
  "options": null
}
```

## Node ID types

| Field | Format | Example |
|-------|--------|---------|
| `id` | 24-character hex string | `"67f2a77682ba966c7f688d95"` |
| `uuid` | UUID v4 | `"9e3d314c-29ef-4b0d-96d1-21e81e355797"` |

Both are required and must be **globally unique** within the process file.

## `extra` field format

In **import JSON**: `extra` is an **escaped JSON string**:
```json
"extra": "{\"modeForm\":\"expand\",\"icon\":\"\"}"
```

In **API modify operations**: `extra` is a **plain dict**:
```json
"extra": {"modeForm": "expand", "icon": ""}
```

This difference is a common source of errors when switching between import and API approaches.

## `params` array — declaring inputs and outputs

```json
{
  "name": "city",
  "type": "string",
  "descr": "City name to query weather for",
  "flags": ["required", "input"],
  "regex": "",
  "regex_error_text": ""
}
```

- `type`: `"string"`, `"number"`, `"object"`, `"boolean"`
- `flags`: `["required", "input"]` (mandatory input), `["input"]` (optional input), `["output"]` (output field)

## Node types

| `obj_type` | Name | Purpose |
|-----------|------|---------|
| `1` | Start | Entry point — exactly one per process |
| `0` | Standard | Logic nodes: API calls, set_param, conditions, replies |
| `2` | Final (End) | Terminal — tasks accumulate here |
| `3` | Escalation | Error/timeout handler — must route to a Final node |

## Validation checklist

### Structure
- [ ] Root is a JSON array `[...]`, not a bare object
- [ ] `parent_id` is `0`
- [ ] `obj_type` is `1` at process level
- [ ] `conv_type` is `"process"`
- [ ] `web_settings` is `[[], []]`

### Nodes
- [ ] Exactly one `obj_type: 1` (Start) node
- [ ] At least one `obj_type: 2` (Final) node
- [ ] Every node has unique `id` (24-char hex) and `uuid` (UUID v4)
- [ ] Every node has `condition` with `logics` and `semaphors` arrays
- [ ] Every non-final node has a `"go"` logic entry as its last logic
- [ ] Final nodes (`obj_type: 2`) have empty `logics` and `semaphors`
- [ ] `extra` is a valid escaped JSON string (not a dict)

### References
- [ ] Every `to_node_id` references an existing node `id`
- [ ] Every `err_node_id` references an existing escalation node (`obj_type: 3`)
- [ ] Every semaphor `to_node_id` references an existing escalation node
- [ ] No orphaned nodes (every node reachable from Start)

### JSON validity
- [ ] Valid JSON syntax — verify with `jq . file.json`
- [ ] Single backslash in strings: `"{\"key\":\"val\"}"` NOT `"{\\\"key\\\":\\\"val\\\"}"`
- [ ] No trailing commas
- [ ] No comments

### Logic rules
- [ ] API nodes have all required fields (see `knowledge/node-types.md`)
- [ ] Every API node has a timeout semaphor
- [ ] `set_param` keys match in `extra` and `extra_type`
- [ ] `api_rpc_reply` keys match in `res_data` and `res_data_type`
- [ ] Error replies: `throw_exception: true` + `exception_reason` present
- [ ] Success replies: `throw_exception: false`

## Common mistakes table

| Mistake | Consequence | Fix |
|---------|-------------|-----|
| Bare object (not array) | Import rejected | Wrap in `[...]` |
| `\\\"` double-escaping | Corezoid parses incorrectly | Use single `\"` |
| Duplicate `id` or `uuid` | Unpredictable behavior | Generate unique values |
| `to_node_id` → non-existent node | Task gets stuck | Verify all references |
| Missing `err_node_id` | Unhandled errors crash process | Always wire error handlers |
| Missing semaphor on API node | Tasks hang on timeout forever | Add `"time"` semaphor |
| `throw_exception: false` on error reply | Error not propagated | Use `true` for errors |
| `throw_exception: true` on success reply | Success treated as error | Use `false` |
| `obj_type: 2` node with logics | Invalid node | Keep `logics: []` |
| `extra` as dict in import JSON | UI rendering issues | Must be escaped string |

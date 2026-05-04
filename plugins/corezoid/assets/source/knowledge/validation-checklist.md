# Process JSON Validation Checklist

Run through this before finalizing any Corezoid process JSON.

---

## 1. Structural Validation

- [ ] The root is a JSON array `[...]`, not a bare object
- [ ] `parent_id` is `0`
- [ ] `obj_type` is `1` at the process level
- [ ] `conv_type` is `"process"`
- [ ] `scheme` contains both `nodes` and `web_settings`
- [ ] `web_settings` is `[ [], [] ]`

## 2. Node Validation

- [ ] Exactly one node with `obj_type: 1` (Start node) exists
- [ ] At least one node with `obj_type: 2` (Final node) exists
- [ ] Every node has a unique `id` (24-char hex)
- [ ] Every node has a unique `uuid` (UUID v4)
- [ ] Every node has `condition` with `logics` and `semaphors` arrays
- [ ] Every non-final, non-escalation node has a `"go"` logic entry as its last logic
- [ ] Final nodes (`obj_type: 2`) have empty `logics` and `semaphors`
- [ ] `extra` field is a valid escaped JSON string

## 3. Reference Integrity

- [ ] Every `to_node_id` references an existing node `id`
- [ ] Every `err_node_id` references an existing escalation node (`obj_type: 3`)
- [ ] Every `to_node_id` in semaphors references an existing escalation node (`obj_type: 3`)
- [ ] No circular references in the main flow (Start → ... → Final)
- [ ] No orphan nodes (every node reachable from Start, or is an error/timeout handler)

## 4. JSON Validity

- [ ] Valid JSON syntax (use `jq . filename.json` to verify)
- [ ] Single backslash escaping in strings: `"{\"key\":\"value\"}"`
- [ ] No double-backslash escaping: NOT `"{\\\"key\\\":\\\"value\\\"}"`
- [ ] No trailing commas
- [ ] No comments in JSON

## 5. Logic Validation

- [ ] API call logics have all required fields (`method`, `url`, `response`, `version`, etc.)
- [ ] `set_param` logics have matching keys in `extra` and `extra_type`
- [ ] `api_rpc_reply` logics have matching keys in `res_data` and `res_data_type`
- [ ] Success replies use `"throw_exception": false`
- [ ] Error replies use `"throw_exception": true` + `"exception_reason"` field
- [ ] `{{variable}}` references match defined parameter names and response paths
- [ ] Every API call node has a semaphor (`"type": "time"`, `"dimension": "sec"`)

## 6. Common Mistakes Quick Reference

| Mistake | Consequence | Fix |
|---------|------------|-----|
| Bare object instead of array | Import rejected | Wrap in `[...]` |
| Double-backslash `\\\"` in `extra` | Corezoid parses incorrectly | Use single `\"` |
| Duplicate node `id` or `uuid` | Unpredictable behavior | Generate unique values |
| `to_node_id` → non-existent node | Tasks get stuck | Verify all references |
| Missing `err_node_id` | Unhandled errors crash process | Always wire error handlers |
| Missing semaphor on API nodes | Tasks hang forever on timeout | Add `"time"` semaphor |
| `throw_exception: false` on error replies | Error not propagated to caller | Use `true` for errors |
| `throw_exception: true` on success replies | Success treated as error | Use `false` for success |
| Final node with logics | Invalid node definition | Keep `logics: []` |
| `extra` as object instead of string | UI rendering issues | Must be escaped JSON string |

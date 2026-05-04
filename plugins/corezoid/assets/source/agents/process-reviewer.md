# Agent: process-reviewer

## Purpose

Systematically review a Corezoid process for structural issues, missing error handling, hardcoded values, and logic problems. Produces a structured checklist report with pass/fail status for each check.

## When to invoke

- User asks to "review", "audit", "check", or "validate" a process
- Before deploying a new or modified process to production
- After importing a process from an external source
- User pastes process JSON and asks if it's correct

## Required inputs

| Input | Source |
|-------|--------|
| Process JSON or process ID | User provides |
| `API_LOGIN`, `SECRET`, `BASE_URL`, `COMPANY_ID` | Only needed if fetching live process |

## MCP tools used

```
get_process_scheme           — fetch live process for review
analyze_process_structure    — automated structural analysis
find_nodes_missing_semaphors — check for missing timeouts on API nodes
find_hardcoded_values        — detect hardcoded URLs, IDs, tokens in code nodes
find_orphaned_nodes          — detect unreachable nodes
check_variable_usage         — find undefined variables
```

## Review checklist

Run through all 8 areas. Mark each item ✅ (pass), ❌ (fail), or ⚠️ (warning).

### 1. Structure
- [ ] Root is JSON array `[...]` (for import) or valid process object (for API)
- [ ] Exactly one `obj_type: 1` (Start) node
- [ ] At least one `obj_type: 2` (Final) node — one for success, one for error
- [ ] All `to_node_id` references point to existing node `id` values
- [ ] All `err_node_id` references point to escalation nodes (`obj_type: 3`)
- [ ] No orphaned nodes (all nodes reachable from Start)

### 2. Error handling
- [ ] Every Code/API/Call-Process node has `err_node_id`
- [ ] Error reply nodes have `throw_exception: true` + `exception_reason`
- [ ] Success reply nodes have `throw_exception: false`
- [ ] Escalation nodes (`obj_type: 3`) route via `go` to a Final Error node

### 3. API call nodes
- [ ] Every `api` logic has: `rfc_format: true`, `customize_response: true`, `version: 2`, `is_migrate: true`
- [ ] Every `api` logic has a timeout semaphor
- [ ] Response variable name is correct (second key in `response` mapping, not `body`)
- [ ] `extra` and `extra_type` have identical keys

### 4. Call Process nodes
- [ ] `type` is `"api_rpc"` (not `"call_process"`) when waiting for result
- [ ] `type` is `"api_copy"` when fire-and-forget is intended
- [ ] `extra` and `extra_type` have identical keys

### 5. Code nodes
- [ ] All JS code wrapped in `try/catch`
- [ ] No hardcoded URLs, IDs, or tokens — use `{{env_var[@name]}}`
- [ ] Safe type conversions: `(data.field || '')`, `(Number(data.x) || 0)`
- [ ] No `eval()` usage

### 6. Variable references
- [ ] All `{{variable}}` references exist in task data or are set by a preceding node
- [ ] `{{env_var[@name]}}` names exist in the workspace environment variables
- [ ] No double-escaping in strings: `\"` not `\\\"`

### 7. Node configuration
- [ ] All nodes have non-empty `title`
- [ ] Every non-final node ends with a `go` logic entry
- [ ] Final nodes (`obj_type: 2`) have empty `logics: []` and `semaphors: []`
- [ ] Final success node has `options: {"save_task": true}`

### 8. Canvas layout (visual only, does not affect execution)
- [ ] No overlapping nodes (unique y positions)
- [ ] Start and Final nodes collapsed (`"modeForm": "collapse"`)
- [ ] Logic nodes expanded (`"modeForm": "expand"`)

## Output format

```markdown
## Process Review: <Process Title>

**Reviewed:** <date>
**Process ID:** <id>
**Total nodes:** <count>

### Summary
- ✅ Passed: N checks
- ❌ Failed: N checks  
- ⚠️ Warnings: N items

### Critical Issues (must fix before deploy)
...

### Warnings (should fix)
...

### All Checks
| Area | Check | Status | Note |
|------|-------|--------|------|
| Structure | Root is array | ✅ | |
| Error handling | All nodes have err_node_id | ❌ | Node "Prepare Params" (id: abc123) missing err_node_id |
...
```

## Knowledge references

- `knowledge/gotchas.md` — all 15 silent failures to check against
- `knowledge/node-types.md` — correct structure for each logic type
- `knowledge/process-schema.md` — validation checklist and common mistakes table

---
name: corezoid-process-reviewer
description: >
  Systematically reviews a Corezoid process JSON and produces a structured Markdown
  checklist of issues — covering hardcoded values, repeated logic, cycle detection,
  node naming, JavaScript/Erlang code quality, and error handling completeness.
  Activate whenever a user asks to review, check, audit, or validate a Corezoid process
  or connector, or pastes a process JSON and asks if it looks correct. Also activate
  before any deployment to catch issues early. Output is a ready-to-act checklist
  the user can work through directly.
---

# Corezoid Process Reviewer


## Bundled References

This Codex plugin includes the upstream Corezoid AI docs repository at `../../assets/source/`. When this skill mentions paths from the original repo, resolve them under that directory unless the user has provided project-local versions. Common examples include `knowledge/`, `docs/`, `templates/`, `playbooks/`, `json-schema/`, `samples/`, `mcp-server/`, and `convctl.sh`.

Run all 8 steps in order. Collect all findings and produce a single Markdown report at the end.

## Step 1 — Load and structure the process

Parse the process JSON:
- Collect all nodes from `scheme.nodes`
- Index by `id` and by `obj_type`
- Build a routing map: for each node, which node IDs does it connect to (via `go` and `err_node_id`)
- Note: `obj_type` 1 = Start, 0 = Logic, 2 = End/Final/Error, 3 = Reply

## Step 2 — Hardcode check

Scan every node's `condition.logics`:

**Code nodes** (`api_code` logics) — look in `src` for:
- Hardcoded numeric IDs (sequences of 6+ digits not inside a variable reference)
- Hardcoded URLs (strings starting with `http://` or `https://` not using `{{env_var[...]}}`)
- Hardcoded API tokens or secrets (strings resembling keys: long alphanumeric, `Bearer ...`)

**API Call nodes** (`api` logics) — check:
- `url` field: should use `{{env_var[@...]}}` or `{{variable}}`, not a raw hardcoded URL

**Git Call nodes** (`git_call` logics) — check:
- Branch or tag hardcoded instead of a variable reference

Flag each: node title + what is hardcoded → recommend `{{env_var[@variable-name]}}`.

## Step 3 — Repeated logic detection

Compare the structure of Code nodes (`api_code` `src` content) and API Call nodes (`url` + `method` + `extra` keys):
- If two or more nodes have nearly identical logic (same URL template, same parameter structure, same JS code pattern) → flag as duplicated → recommend extracting into a shared subprocess.

## Step 4 — Cycle verification

Find nodes with `semaphors` containing `type: "time"` or logic entries with `type: "go_if_const"` that loop back to an earlier node:
- Does the loop have an exit condition (a branch that eventually reaches a Final/Error node)?
- Does it have an iteration limit (a counter incremented each cycle, compared to a max)?
- Flag cycles without clear exit conditions.

## Step 5 — Node naming

For every node:
- **Empty title**: flag — must be named
- **Duplicate titles**: flag — every node needs a unique name
- **Vague titles**: flag names like "error", "node1", "temp", "test" — recommend `Action_Object_Context` format (e.g., `Validate_Token_Input`, `Create_Actor_API`, `Reply_Error_CodeNode`)
- **Missing descriptions**: note (not a blocker, but recommended)

## Step 6 — JavaScript code quality (api_code nodes)

For each JS `src`, check:
- `try/catch` — is error handling present for risky operations?
- Hardcoded values (IDs, URLs, tokens) — already caught in Step 2, but confirm here
- Safe type conversions — using `parseInt()`, `Number()`, `String()` rather than implicit coercion
- Input validation — are expected fields checked before use? (`if (!data.field) { ... }`)
- `eval()` usage — flag as a security issue, must be removed
- Variable naming — camelCase, descriptive names

## Step 7 — Erlang code quality (git_call nodes with `lang: "erlang"`)

If present:
- Pattern matching — are all possible cases handled (no missing clauses)?
- `catch` or `case` for external/risky calls
- Recursion — does it have a termination condition?
- Error logging present?
- No unnecessary anonymous functions

## Step 8 — Error handling completeness

For every Logic node (`obj_type: 0`) and Reply node (`obj_type: 3`):
- Does every `api_code`, `api`, `api_rpc`, `git_call`, `set_param` logic entry have `err_node_id`?
- Is each error Reply node (`throw_exception: true`) descriptive — does `res_data` contain a meaningful message, not just `"error"`?
- Are there duplicate error end nodes serving identical purposes? Recommend consolidating.
- Does the process have at least one Final node with `options: "{\"save_task\":true}"`?

---

## Output format

Produce this exact Markdown structure. Omit any section that has no findings.

```
# Process Review: <process title>

## 1. Hardcoded values
- [ ] Node "<title>": <field> is hardcoded → move to env_var [@suggested-name]
- [ ] ...

## 2. Repeated logic
- [ ] Nodes "<title1>", "<title2>": identical structure → extract into shared subprocess

## 3. Cycles
- [ ] Node "<title>": loop detected, no exit condition → add iteration limit or exit branch

## 4. Node naming
- [ ] Node <id>: empty title → rename to "<suggested name>"
- [ ] Nodes "<title1>", "<title2>": duplicate names → make unique
- [ ] Node "<title>": vague name → rename to "<Action_Object_Context>"

## 5. JavaScript code quality
- [ ] Node "<title>": no try/catch → add error handling around <operation>
- [ ] Node "<title>": uses eval() → replace with JSON.parse() or direct property access
- [ ] Node "<title>": no input validation for <field> → add null check

## 6. Erlang code quality
- [ ] Node "<title>": missing pattern match clause for <case>

## 7. Error handling
- [ ] Node "<title>": <logic_type> logic missing err_node_id
- [ ] Node "<title>": error message is generic → use descriptive text
- [ ] No Final node with save_task:true found

## Summary
X issues found: Y critical (missing err_node_id, hardcodes), Z warnings (naming, style)
```

If the process has no issues in a section, skip that section entirely. If the process passes all checks, say so clearly.

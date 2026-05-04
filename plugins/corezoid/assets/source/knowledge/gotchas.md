# Corezoid Gotchas — Silent Failures and Non-Obvious Rules

Things that cause incorrect behavior without obvious error messages.
Every item here has burned someone at least once.

---

## 1. Response body variable name

**Wrong:** `{{body.field}}` — never works.

**Why:** The response body is stored under the **second key** you define in the `response` mapping, not under a hardcoded `body` key.

```json
"response": {"header": "{{header}}", "response": "{{body}}"}
```
→ Access as `{{response.field}}` (not `{{body.field}}`)

```json
"response": {"header": "{{header}}", "data": "{{body}}"}
```
→ Access as `{{data.field}}`

The internal name `body` is Corezoid's variable for the raw response — your task key is whatever you name the second mapping entry.

---

## 2. JWT session token ≠ API SECRET

The `eyJhbGci...` token in your browser is a session cookie. It expires. It cannot sign API calls.

The `SECRET` is a separate static key in **Account Settings → API** tab. It looks like a long random string, not a JWT.

Using the JWT as SECRET: `403 cookie or headers are not valid`.

---

## 3. Process JSON must be an array

For UI import or JSON-based create, the root must be `[{...}]`:

```json
[{"obj_type": 1, "title": "My Process", ...}]
```

A bare object `{"obj_type": 1, ...}` → import rejected with a parse error.

---

## 4. `extra` format changes by context

In **import JSON** (for UI): `extra` must be an **escaped JSON string**:
```json
"extra": "{\"modeForm\":\"expand\",\"icon\":\"\"}"
```

In **API modify operations** (via `client.make_request`): `extra` is a **plain dict**:
```json
"extra": {"modeForm": "expand", "icon": ""}
```

Sending a string in an API modify call → field ignored or error. Sending a dict in import JSON → parsing error.

---

## 5. Node `id` vs `uuid` — both required and different

Every node needs two distinct identifiers:
- `id` — 24-character hex string (`"67f2a77682ba966c7f688d95"`) — used for routing (`to_node_id`, `err_node_id`)
- `uuid` — UUID v4 (`"9e3d314c-29ef-4b0d-96d1-21e81e355797"`) — used for global tracking

Generate them separately. A UUID v4 is NOT the same as a 24-char hex.

```python
import uuid, secrets
node_id = secrets.token_hex(12)  # 24-char hex
node_uuid = str(uuid.uuid4())    # UUID v4
```

---

## 6. REF uniqueness — duplicate REFs silently drop tasks

If a process receives two tasks with the same `ref` and the process uses `ref_mask: true`, the second task silently overwrites or drops. No error is raised.

For fire-and-forget patterns, generate a unique REF:
```json
"ref": "$.random()"
```
Or combine input fields: `"{{user_id}}_{{timestamp}}"`.

---

## 7. Missing semaphor = tasks hang forever

An `api` logic node without a timeout semaphor will hold tasks indefinitely if the external API is slow, unreachable, or returns no response.

Always add:
```json
"semaphors": [{"type": "time", "value": 30, "dimension": "sec", "to_node_id": "<timeout-node>"}]
```

---

## 8. `api_rpc` vs `api_copy` confusion

| Use | Logic type | Behavior |
|-----|-----------|---------|
| Need the result back | `api_rpc` | Blocks. Waits for target to reply. |
| Fire and forget | `api_copy` | Continues immediately. No result. |

Using `api_copy` when you need the result → your process gets no data back and continues with empty fields. Silent failure.

---

## 9. `company_id` format

Must be `i` + numeric workspace ID: `i469b4ea8...`

Plain numeric ID (without `i`) → most API operations fail or return wrong results.

Find it: Corezoid workspace URL contains the workspace ID segment after `admin.corezoid.com/`.

---

## 10. `get_process_scheme` response path

The scheme is nested inside a list:

```python
result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']
#                               ^--- LIST, must index [0]
```

`result['ops'][0]['scheme']` is a **list**, not a dict. Forgetting `[0]` → `TypeError: list indices must be integers`.

---

## 11. Node modify is never partial

When modifying a node via API, you must send the **complete payload**: `title`, `description`, `logics`, `semaphors`, `extra`, `options`, `position`. Any field you omit is silently reset to its default (usually empty/null).

Partial modify → logics disappear, connections break, semaphors vanish.

---

## 12. Code logic type is `api_code`, not `code`

**Wrong:** `{"type": "code", "lang": "js", ...}` → `Unexpected logic type` error immediately.

**Correct:** `{"type": "api_code", "lang": "js", "src": "..."}` — this is the actual type name in the API.

---

## 12a. `set_param` not `api_code` for simple operations

Code nodes (`api_code` logic) spin up a V8 (JS) or BEAM (Erlang) runtime for every task. `set_param` runs natively. For string concatenation, arithmetic, variable renaming — always `set_param`.

Using `api_code` where `set_param` works → measurable latency increase under load.

---

## 12b. `set_param` cannot extract array elements from API responses

`{{geo.results.0.latitude}}` returns an **empty string** when `results` is a JSON array from an API response body. Corezoid's template substitution does not traverse arrays via dot notation.

**Wrong:**
```json
"extra": {"lat": "{{geo.results.0.latitude}}"}
```
→ `lat` = `""` (silently empty, no error)

**Correct:** Use `api_code` with JS:
```js
data.lat = String(data.geo.results[0].latitude);
data.lon = String(data.geo.results[0].longitude);
```

Applies to any API response where the body contains a JSON array.

---

## 12c. Commit version — use creation `change_time`, not post-modify

After creating a process and modifying nodes, the `change_time` updates. But commit requires the **original** `change_time` from process creation (or last successful commit).

**Wrong:** using `change_time` fetched after node modifications → `"This version hasn't any commit"`

**Correct:** save `version` immediately after `create`, reuse it for all node ops AND the commit:
```python
version = client.get_process_details(conv_id)['ops'][0]['change_time']  # save once
# ... create nodes with version ...
# ... modify nodes with version ...
client.make_request({'ops': [{'obj': 'commit', 'type': 'confirm', 'version': version, ...}]})

---

## 13. `go_if_const` conditions are on the inner array, not the logic entry

```python
# ❌ Wrong — these keys don't exist at logic entry level
logic.get('param')   # None
logic.get('const')   # None

# ✅ Correct
for cond in logic['conditions']:
    param = cond['param']
    const = cond['const']
    fun   = cond['fun']
```

---

## 14. `api_rpc_reply` error replies require `exception_reason`

```json
{"throw_exception": true, "exception_reason": "API call failed"}
```

Missing `exception_reason` when `throw_exception: true` → the caller gets a generic or empty error.

---

## 15. Dashboard chart API — correct field names

Valid chart types: `column`, `pie`, `funnel`, `table`. `bar` does NOT exist.

**Chart create requires these fields** (discovered from live API — docs are wrong):
```json
{
  "obj": "chart",
  "type": "create",
  "name": "unique_snake_case_name",
  "sort": 1,
  "obj_type": "column",
  "dashboard_id": 123456,
  "title": "Human readable title",
  "series": [
    {"conv_id": 234567, "node_id": "hex-node-id", "title": "Metric Label"}
  ]
}
```

Key gotchas:
- `type: "create"` (NOT `obj_type: "create"`) for chart creation
- `obj_type` = the chart type string (`"column"`, `"pie"`, `"funnel"`, `"table"`)
- `series` (NOT `metrics`) — `metrics` returns `Key 'series' is required`
- `name` — required, unique snake_case identifier
- `sort` — required integer for ordering

Grid sizing uses `width`/`height`, NOT `w`/`h`.

Chart `modify` requires `obj_type` + full `series` — partial modify returns validation errors.

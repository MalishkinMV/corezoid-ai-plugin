# Corezoid Node Types and Logic Reference

## Node type map

| `obj_type` | Name | Logic types used |
|-----------|------|-----------------|
| `1` | Start | `go` only |
| `0` | Standard (API Call) | `api` |
| `0` | Standard (Set Param) | `set_param` |
| `0` | Standard (Code) | `api_code` |
| `0` | Standard (Call Process RPC) | `api_rpc` |
| `0` | Standard (Call Process async) | `api_copy` |
| `0` | Standard (Condition) | `go_if_const`, `condition` |
| `3` | Escalation / Reply | `api_rpc_reply` |
| `2` | Final / End | none (`logics: []`) |

## Logic types in full

### `go` — unconditional routing
```json
{"type": "go", "to_node_id": "<target-node-id>"}
```
Every non-final, non-escalation node must end with a `go` as its last logic entry.

---

### `api` — HTTP call to external endpoint

**All fields are required:**

```json
{
  "type": "api",
  "is_migrate": true,
  "rfc_format": true,
  "format": "",
  "content_type": "application/json",
  "method": "GET",
  "url": "https://api.example.com/endpoint",
  "extra": {},
  "extra_type": {},
  "extra_headers": {},
  "cert_pem": "",
  "max_threads": 5,
  "send_sys": true,
  "debug_info": false,
  "err_node_id": "<escalation-node-id>",
  "customize_response": true,
  "response": {"header": "{{header}}", "response": "{{body}}"},
  "response_type": {"header": "object", "response": "object"},
  "version": 2
}
```

**Response variable access gotcha:**
The body is stored under the **second key name** in `response`, not under `body`.
- With `"response": {"header": "{{header}}", "response": "{{body}}"}` → access as `{{response.field}}`
- With `"response": {"header": "{{header}}", "data": "{{body}}"}` → access as `{{data.field}}`
- `{{body.field}}` will NEVER work — `body` is the internal name, not the task key.

**POST body fields** go in `extra` / `extra_type`:
```json
"extra": {"param1": "{{value1}}", "param2": "{{value2}}"},
"extra_type": {"param1": "string", "param2": "string"}
```

**Timeout semaphor** — MANDATORY on every `api` node:
```json
"semaphors": [{"type": "time", "value": 30, "dimension": "sec", "to_node_id": "<timeout-escalation-id>"}]
```
Without this, tasks hang forever if the external API is unresponsive.

---

### `set_param` — native variable assignment (prefer over `code`)

```json
{
  "type": "set_param",
  "extra": {"output_var": "{{response.field}}"},
  "extra_type": {"output_var": "string"},
  "err_node_id": "<escalation-node-id>"
}
```

**Use `set_param` instead of `code` for:**
- Variable copy/rename: `"new_key": "{{old_key}}"`
- String concatenation: `"full": "{{first}}_{{last}}"`
- Arithmetic: `"total": "$.math({{price}}*{{qty}})"`
- Dynamic key access: `"value": "{{response.{{currency_code}}}}"`

`set_param` is native — zero overhead. `code` nodes spin up V8/BEAM — measurably slower.

**Built-in functions in `set_param`:**
- `$.math({{a}}+{{b}})` — arithmetic (exactly 2 operands; nest for more)
- `$.date()` — current timestamp
- `$.random()` — random value
- `$.md5_hex()`, `$.sha256_hex()` — hashing

---

### `api_code` — JS or Erlang code (use only when `set_param` can't)

**The logic type name is `api_code`, NOT `code`.** Using `"type": "code"` returns `Unexpected logic type` error.

```json
{
  "type": "api_code",
  "lang": "js",
  "src": "try {\n  data.result = data.input.split(',').length;\n} catch(e) {\n  throw new Error('Processing failed: ' + e.message);\n}",
  "err_node_id": "<escalation-node-id>"
}
```

- `lang`: `"js"` or `"erl"`
- Always wrap JS in `try/catch`
- Use defensive defaults: `(data.field || '')`, `(Number(data.x) || 0)`
- Never hardcode URLs, IDs, secrets — use `{{env_var[@name]}}` references
- **Array access from API responses:** use JS `data.geo.results[0].latitude` — `set_param` dot notation returns empty string for array elements

**Use `api_code` only for:** array access/manipulation, string `.length`, regex, JSON parsing, complex conditionals.

---

### `api_rpc` — synchronous call to another Corezoid process

```json
{
  "type": "api_rpc",
  "conv_id": "@process-alias",
  "extra": {"param1": "{{value1}}", "param2": "{{value2}}"},
  "extra_type": {"param1": "string", "param2": "number"},
  "err_node_id": "<escalation-node-id>"
}
```

- Blocks until the target process calls `api_rpc_reply`
- `conv_id`: use `@alias` (preferred) or numeric process ID
- `extra` and `extra_type` must have identical keys
- Use when you **need the result back**

---

### `api_copy` — fire-and-forget async send

```json
{
  "type": "api_copy",
  "conv_id": "@process-alias",
  "ref": "{{unique_ref}}",
  "mode": "create",
  "is_sync": false,
  "data": {"field": "{{value}}"},
  "data_type": {"field": "string"},
  "err_node_id": "<escalation-node-id>"
}
```

- Does NOT wait for the target to finish — continues immediately
- Use for async fan-out, logging, notifications
- Use `api_rpc` when you need the output from the called process

---

### `api_rpc_reply` — return response to calling process

```json
{
  "type": "api_rpc_reply",
  "mode": "key_value",
  "res_data": {"result": "{{computed_value}}", "city": "{{resolved_city}}"},
  "res_data_type": {"result": "object", "city": "string"},
  "throw_exception": false
}
```

**Error reply** — must include `exception_reason`:
```json
{
  "type": "api_rpc_reply",
  "mode": "key_value",
  "res_data": {"description": "API call failed: {{error_detail}}"},
  "res_data_type": {"description": "string"},
  "throw_exception": true,
  "exception_reason": "API call failed"
}
```

---

### `go_if_const` — conditional routing (most common branch type)

```json
{
  "type": "go_if_const",
  "to_node_id": "<target-node-id>",
  "conditions": [
    {"param": "status", "const": "active", "fun": "eq", "cast": "string"},
    {"param": "age",    "const": "18",     "fun": "ge", "cast": "number"}
  ]
}
```

- Multiple `conditions` items = AND logic (all must match)
- Multiple `go_if_const` entries in a node = OR / if-else-if (first match wins)
- `fun` values: `eq`, `ne`, `not_eq`, `gt`, `lt`, `ge`, `le`, `regex`
- **Gotcha:** `param`, `const`, `fun` are on `conditions` items, NOT on the `go_if_const` object

---

## Error handling pattern

Every functional node needs an error escalation path:

```
[Functional Node] --err_node_id--> [Escalation obj_type:3] --go--> [Final Error obj_type:2]
[API Node]        --semaphor-----> [Timeout Escalation]    --go--> [Timeout Final]
```

Escalation nodes (`obj_type: 3`) must:
- Have `api_rpc_reply` with `throw_exception: true` and `exception_reason`
- Route via `go` to a Final Error node (`obj_type: 2`)

---

## Canvas layout discipline

| Lane | x position | Contents |
|------|-----------|---------|
| Main flow | x ≈ 1000 | Start, logic nodes, Reply Success, Final |
| Timeout handlers | x ≈ 740–800 | Reply Timeout, Error Timeout |
| Error handlers | x ≈ 1150–1264 | Reply Error, Error Final |

- Vertical step: 160–200px between nodes
- Collapse Start and Final nodes (`modeForm: "collapse"`)
- Expand logic nodes (`modeForm: "expand"`)
- No overlapping nodes — each needs its own y position

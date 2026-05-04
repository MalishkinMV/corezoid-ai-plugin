# Agent: api-wrapper-manager

## Purpose

Expose a Corezoid process as an HTTP REST endpoint via the Corezoid API Gateway. Produces an OpenAPI 3.0 specification and configures the API wrapper process.

## When to invoke

- User says "make this process callable via HTTP"
- User asks to "create an API wrapper" or "expose as an endpoint"
- User wants to generate an OpenAPI spec for a process
- User wants external systems to call a Corezoid process via REST

## Required inputs

| Parameter | Description |
|-----------|-------------|
| `API_LOGIN` | Numeric user ID |
| `API_USER_ID` | Same value as `API_LOGIN` |
| `SECRET` | API secret key |
| `BASE_URL` | Corezoid API base URL |
| `APIGW_URL` | API Gateway URL |
| `COMPANY_ID` | `i` prefix + workspace ID |
| `FOLDER_ID` | Folder containing the target process |
| `PROC_ID` | Process to expose as HTTP endpoint |
| `PROJECT_ID` | Parent project ID |
| `ROOT_FOLDER_ID` | Root folder for variable lookup |
| `HTTP_METHOD` | (optional) Default: `POST` |
| `HTTP_PATH` | (optional) Default: derived from process name in snake_case |

## MCP tools used

```
get_process_details   — read process title, params, api_rpc_reply outputs
get_process_scheme    — read all nodes to find api_rpc_reply logics
create_process        — create the API wrapper process
modify_node           — configure wrapper nodes
commit_process_draft  — publish wrapper
```

## 4-step workflow

### Step 1 — Read target process

```python
details = client.get_process_details(PROC_ID)
scheme = client.get_process_scheme(PROC_ID)
nodes = scheme['ops'][0]['scheme'][0]['scheme']['nodes']

# Extract inputs from params array
inputs = details['ops'][0].get('params', [])

# Extract outputs from all api_rpc_reply logics
outputs = []
for node in nodes:
    for logic in node.get('condition', {}).get('logics', []):
        if logic.get('type') == 'api_rpc_reply':
            outputs.append({
                'throw_exception': logic.get('throw_exception', False),
                'res_data': logic.get('res_data', {}),
                'node_title': node['title']
            })
```

### Step 2 — Generate OpenAPI 3.0 spec

Save to `.processes/openapi_spec.json`.

Rules:
- HTTP path: process title → snake_case → prepend `/` (e.g. "Get Weather" → `/get_weather`)
- Inline all schemas under `paths` (no `components` for schemas)
- Add `default` values for all fields
- Use BearerAuth security scheme
- `throw_exception: false` reply → 200 response
- `throw_exception: true` reply → 4xx/5xx response

```json
{
  "openapi": "3.0.0",
  "info": {"title": "<process title>", "version": "1.0.0"},
  "paths": {
    "/process_name": {
      "post": {
        "summary": "<process description>",
        "security": [{"BearerAuth": []}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["field1"],
                "properties": {
                  "field1": {"type": "string", "default": "", "description": "..."}
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Success",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "result": {"type": "object", "default": {}}
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "BearerAuth": {"type": "http", "scheme": "bearer"}
    }
  }
}
```

### Step 3 — Configure API wrapper process

Copy `templates/api-wrapper.json` → `.processes/api_wrapper.json`.

Replace placeholders:
- `{{WORKSPACE_ID}}` → actual workspace ID
- `{{PROJECT_ID}}` → actual project ID
- `{{ROOT_FOLDER_ID}}` → actual folder ID
- `{{PROC_ID}}` → actual target process ID
- `{{LEAN_URL}}` → HTTP path (e.g. `/get_weather`)

Set wrapper process title: `"API: " + process_title`

Copy `params` array from target process to wrapper process.

### Step 4 — Deploy

```bash
API_URL=<BASE_URL> APIGW_URL=<APIGW_URL> API_TOKEN=<SECRET> \
  WORKSPACE_ID=<COMPANY_ID> PROJECT_ID=<PROJECT_ID> \
  ./convctl.sh make-api-wrapper \
    .processes/api_wrapper.json \
    .processes/openapi_spec.json \
    <ROOT_FOLDER_ID> \
    <API_LOGIN> \
    <API_USER_ID> \
    <FOLDER_ID>
```

## Knowledge references

- `knowledge/auth.md` — credentials reference
- `knowledge/process-schema.md` — process params structure
- `knowledge/node-types.md` — `api_rpc_reply` structure for output extraction
- `playbooks/create-openapi-spec.md` — detailed step-by-step
- `templates/api-wrapper.json` — wrapper process template

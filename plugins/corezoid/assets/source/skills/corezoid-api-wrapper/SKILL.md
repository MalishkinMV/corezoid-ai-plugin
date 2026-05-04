---
name: corezoid-api-wrapper
description: >
  Creates an HTTP API wrapper for an existing Corezoid process — generates an OpenAPI 3.0.0
  specification and a wrapper process JSON, then deploys both via the convctl make-api-wrapper
  command. Activate whenever a user asks to make a process callable via HTTP, create an API
  wrapper, expose a process as a REST endpoint, generate an OpenAPI spec for a process,
  or mentions "make-api-wrapper". This is the correct skill for turning any Corezoid process
  into an externally-callable HTTP API with proper authentication and API Gateway integration.
---

# Corezoid API Wrapper

## Required parameters (ask for all before starting)

| Parameter | Notes |
|-----------|-------|
| `API_URL` | Corezoid admin URL, e.g. `https://admin.corezoid.com` |
| `APIGW_URL` | API Gateway URL, e.g. `https://apigw.corezoid.com` |
| `API_LOGIN` | Login (email or username) |
| `API_USER_ID` | Numeric user ID |
| `API_TOKEN` | Auth token |
| `WORKSPACE_ID` | Workspace identifier |
| `PROC_ID` | ID of the target process to wrap |
| `PROJECT_ID` | Project identifier |
| `ROOT_FOLDER_ID` | Root folder ID |
| `FOLDER_ID` | Folder where the wrapper will be deployed |
| `HTTP_METHOD` | *(optional)* HTTP method — default `POST` |
| `HTTP_PATH` | *(optional)* HTTP path — default: auto-generated from process name |
| `LEAN_URL` | *(optional)* Lean URL — default `""` |

## Step 1 — Read the target process

The target process must already be exported at `.processes/target_process.json`.
If it isn't, export it first:

```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh fetch-process <PROC_ID> .processes/target_process.json
```

Extract from the JSON:
- **Inputs**: `params` array — each entry has `name`, `type`, `descr`, `flags`
- **Outputs**: find all nodes with `api_rpc_reply` in their `condition.logics`:
  - `throw_exception: false` → success response, note `res_data` keys and types
  - `throw_exception: true` → error response, note `res_data` keys, types, and `exception_reason` if present

Also note:
- Process `title` (used for HTTP path generation and wrapper title)
- Any `{{env_var[@...]}}` references (document as dependencies in the spec)

## Step 2 — Generate OpenAPI 3.0.0 spec

**HTTP method**: use `HTTP_METHOD` if provided, otherwise `POST`.

**HTTP path**: use `HTTP_PATH` if provided, otherwise convert process title to snake_case and prepend `/`:
- "Get User By Id" → `/get_user_by_id`
- "CreateOrder" → `/create_order`
- "Update Actor Status" → `/update_actor_status`

**Rules for the spec:**
- Version: `"openapi": "3.0.0"`
- All schemas go inline under `paths` — do NOT use `components/schemas` for request/response models
- Add `default` values for all fields based on type and description (use meaningful examples)
- For POST/PUT/PATCH: parameters go in `requestBody` as `application/json`
- For GET/DELETE: parameters go as `query` parameters
- Include BearerAuth security scheme and apply it globally

**Template:**

```json
{
  "info": {
    "description": "<process description or purpose>",
    "title": "<process title>",
    "version": "1.0.0"
  },
  "openapi": "3.0.0",
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  },
  "security": [{"BearerAuth": []}],
  "paths": {
    "<HTTP_PATH>": {
      "<http_method_lowercase>": {
        "summary": "<process title>",
        "description": "<process description>",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["<required_field_names>"],
                "properties": {
                  "<field_name>": {
                    "type": "<string|number|boolean|object>",
                    "description": "<descr from params>",
                    "default": "<meaningful example value>"
                  }
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
                    "<output_field>": {
                      "type": "<type>",
                      "default": "<example>"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Error",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {"type": "string", "default": "error description"}
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "servers": []
}
```

Save to `.processes/openapi_spec.json`.

## Step 3 — Create the wrapper process JSON

Copy the content of `templates/api-wrapper.json` to `.processes/api_wrapper.json`.

Then apply these changes:
1. Set `title` to `"API: " + <original process title>`
2. Copy the `params` array from `target_process.json` into the wrapper (same input interface)
3. Replace all placeholder tokens throughout the JSON:

| Placeholder | Replace with |
|-------------|-------------|
| `{{WORKSPACE_ID}}` | actual WORKSPACE_ID value |
| `{{PROJECT_ID}}` | actual PROJECT_ID value |
| `{{ROOT_FOLDER_ID}}` | actual ROOT_FOLDER_ID value |
| `{{PROC_ID}}` | actual PROC_ID value |
| `{{LEAN_URL}}` | actual LEAN_URL value (or `""` if not provided) |

Save the modified JSON to `.processes/api_wrapper.json`.

## Step 4 — Deploy with convctl

```bash
API_URL=<API_URL> APIGW_URL=<APIGW_URL> API_TOKEN=<API_TOKEN> \
  WORKSPACE_ID=<WORKSPACE_ID> PROJECT_ID=<PROJECT_ID> \
  ./convctl.sh make-api-wrapper \
    .processes/api_wrapper.json \
    .processes/openapi_spec.json \
    <ROOT_FOLDER_ID> <API_LOGIN> <API_USER_ID> <FOLDER_ID>
```

If the command fails, read the error output and fix the relevant file before retrying. Common issues are documented in `docs/error_explanations.md`.

## Output summary

After successful deployment, report:
- HTTP endpoint: `<APIGW_URL><HTTP_PATH>`
- Method: `<HTTP_METHOD>`
- Authentication: Bearer token required
- Input parameters: list from `params`
- Success response fields: from `res_data` of non-exception reply nodes
- Error response fields: from `res_data` of exception reply nodes

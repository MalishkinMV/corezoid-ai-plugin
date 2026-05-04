---
name: corezoid-apigw-manager
description: >
  Creates and manages Corezoid API Gateway (APIGW) endpoints — generates OpenAPI specs,
  builds wrapper processes, configures authentication, sets HTTP paths and methods, and
  deploys via convctl make-api-wrapper. Activate whenever a user asks to expose a process
  as an HTTP API, create an API Gateway endpoint, set up a REST endpoint for a Corezoid
  process, configure external access to a process, create a webhook, or asks about APIGW,
  API paths, or making a process callable from the outside. This skill covers the full
  workflow from zero to a live HTTP endpoint.
---

# Corezoid API Gateway Manager

## What API Gateway does

APIGW exposes any Corezoid process as an HTTP REST endpoint. External clients call a standard
HTTP API; the Gateway forwards requests as tasks into the process, waits for completion,
and returns the result — hiding all Corezoid internals from the outside world.

**The core pattern (synchronous mode):**
```
Client → APIGW → Wrapper Process → Target Process → API Call {{__callback_url}} → APIGW → Client
```

## Required parameters (collect all before starting)

| Parameter | Description |
|-----------|-------------|
| `API_URL` | Corezoid admin URL |
| `APIGW_URL` | API Gateway URL (e.g. `https://apigw.corezoid.com`) |
| `API_TOKEN` | Auth token |
| `API_LOGIN` | Login email |
| `API_USER_ID` | Numeric user ID |
| `WORKSPACE_ID` | Workspace ID |
| `PROC_ID` | ID of the process to expose |
| `PROJECT_ID` | Project ID |
| `ROOT_FOLDER_ID` | Root folder ID |
| `FOLDER_ID` | Folder where wrapper process will be deployed |
| `HTTP_METHOD` | HTTP verb — default `POST` |
| `HTTP_PATH` | URL path — default: auto-generated from process name |
| `LEAN_URL` | Optional lean URL — default `""` |

---

## Full workflow (4 steps)

### Step 1 — Read the target process

Export if not already done:
```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh fetch-process <PROC_ID> .processes/target_process.json
```

Extract from the JSON:
- **Inputs**: `params` array — `name`, `type`, `descr`, `flags` (required vs optional)
- **Success outputs**: `api_rpc_reply` nodes where `throw_exception: false` → `res_data` keys
- **Error outputs**: `api_rpc_reply` nodes where `throw_exception: true` → what triggers each error
- **Process title**: used for HTTP path auto-generation

### Step 2 — Determine HTTP method and path

**HTTP method** — ask the user or use these defaults:
- `POST` — creating resources, submitting data (default)
- `GET` — fetching data (params go as query string, not body)
- `PUT` / `PATCH` — updating
- `DELETE` — removing

**HTTP path** — if `HTTP_PATH` not provided:
- Convert process title to snake_case and prepend `/`
- "Get User By Id" → `/get_user_by_id`
- "Create Payment Order" → `/create_payment_order`

### Step 3 — Generate OpenAPI 3.0.0 spec

Save to `.processes/openapi_spec.json`.

Rules:
- All schemas inline under `paths` — no `components/schemas`
- For POST/PUT/PATCH: params go in `requestBody` as `application/json`
- For GET/DELETE: params go as `query` parameters
- Add `default` values to every field (meaningful examples)
- Required fields from `params` where `flags` includes `"required"`
- BearerAuth security applied globally

```json
{
  "info": {
    "title": "<process title>",
    "description": "<what this API does>",
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
      "<method>": {
        "summary": "<process title>",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["<required_fields>"],
                "properties": {
                  "<field>": {"type": "<type>", "description": "<descr>", "default": "<example>"}
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
                    "<output_field>": {"type": "<type>", "default": "<example>"}
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

### Step 4 — Build the wrapper process

Copy `templates/api-wrapper.json` → `.processes/api_wrapper.json`.

Apply these changes:
1. Set `title` → `"API: " + <original process title>`
2. Copy `params` array from `target_process.json` (same input interface)
3. Replace all placeholders:

| Placeholder | Replace with |
|-------------|-------------|
| `{{WORKSPACE_ID}}` | actual WORKSPACE_ID |
| `{{PROJECT_ID}}` | actual PROJECT_ID |
| `{{ROOT_FOLDER_ID}}` | actual ROOT_FOLDER_ID |
| `{{PROC_ID}}` | actual PROC_ID |
| `{{LEAN_URL}}` | actual LEAN_URL or `""` |

### Step 5 — Deploy

```bash
API_URL=<API_URL> APIGW_URL=<APIGW_URL> API_TOKEN=<API_TOKEN> \
  WORKSPACE_ID=<WORKSPACE_ID> PROJECT_ID=<PROJECT_ID> \
  ./convctl.sh make-api-wrapper \
    .processes/api_wrapper.json \
    .processes/openapi_spec.json \
    <ROOT_FOLDER_ID> <API_LOGIN> <API_USER_ID> <FOLDER_ID>
```

If it fails — check `docs/error_explanations.md` and fix before retrying.

---

## Authentication options

| Method | Configuration |
|--------|--------------|
| **Bearer token** (default) | `Authorization: Bearer <token>` — configured in OpenAPI spec BearerAuth |
| **API Key** | API key must have "Task management" rights on target process |
| **OAuth 2.0** | `client_id` + `secret_id` via Corezoid Single-Account |

Default wrapper uses Bearer token. Specify if the user needs a different auth method.

---

## Synchronous vs asynchronous

**Synchronous (default):**
- Client waits for response
- Process must have an API Call node calling `{{__callback_url}}`
- Timeout: 60 seconds (HTTP 500 if exceeded)
- Good for: real-time responses, small data operations

**Asynchronous:**
- Client gets immediate task ID acknowledgment
- Process runs in background
- Good for: long-running operations, batch processing

---

## Request data in the process

Inside the target process, access request data via:
- `{{api_gateway.request.body}}` — full request body as object
- `{{api_gateway.request.path}}` — the URL path
- `{{api_gateway.request.method}}` — HTTP method
- `{{api_gateway.request.headers.Header-Name}}` — specific header

---

## After deployment — report to user

```
✅ API endpoint deployed

URL: <APIGW_URL><HTTP_PATH>
Method: <HTTP_METHOD>
Auth: Bearer token required

Input parameters:
  - <field>: <type> (<required/optional>) — <description>

Success response (200):
  - <field>: <type>

Error response (400):
  - error: string
```

---

## Reference

Full documentation: `docs/apigw/README.md`
Full playbook: `playbooks/create-openapi-spec.md`
Wrapper template: `templates/api-wrapper.json`

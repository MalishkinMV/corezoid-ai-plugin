# Agent: apigw-manager

## Purpose

Expose an existing Corezoid process as an HTTP REST endpoint via the Corezoid API Gateway (APIGW). Produces an OpenAPI 3.0.0 spec and wrapper process JSON, then deploys both via `convctl make-api-wrapper`.

## When to invoke

- User wants to expose a process as an HTTP API
- User asks to create a REST endpoint, webhook, or external-facing API for a Corezoid process
- User asks about APIGW, HTTP path configuration, or making a process callable from outside

## Required inputs

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
| `FOLDER_ID` | Folder where wrapper process deploys |
| `HTTP_METHOD` | HTTP verb — default `POST` |
| `HTTP_PATH` | URL path — default: auto-generated from process name in snake_case |

## Workflow (4 steps)

### Step 1 — Read the target process

Extract from the process JSON:
- **Inputs**: `params` array — `name`, `type`, `descr`, `flags`
- **Success outputs**: `api_rpc_reply` nodes where `throw_exception: false` → `res_data` keys/types
- **Error cases**: `api_rpc_reply` nodes where `throw_exception: true` → what triggers each
- **Process title** → auto-generate HTTP path: `"Get User By Id"` → `/get_user_by_id`

### Step 2 — Generate OpenAPI spec

Save to `.processes/openapi_spec.json`:
- Use OpenAPI 3.0.0 format
- Inline schemas under `paths` — no `$ref` or `components`
- Add `default` values to all fields
- Use BearerAuth security scheme
- Include all error cases as separate response codes (400/500)

### Step 3 — Configure wrapper process

Copy `templates/api-wrapper.json` → `.processes/api_wrapper.json`:
- Set title to `"API: " + process_title`
- Copy `params` from target process
- Replace placeholders: `{{WORKSPACE_ID}}`, `{{PROJECT_ID}}`, `{{ROOT_FOLDER_ID}}`, `{{PROC_ID}}`, `{{LEAN_URL}}`

### Step 4 — Deploy

```bash
API_URL=<API_URL> APIGW_URL=<APIGW_URL> API_TOKEN=<TOKEN> WORKSPACE_ID=<WS_ID> \
  PROJECT_ID=<PROJ_ID> ./convctl.sh make-api-wrapper \
  .processes/api_wrapper.json .processes/openapi_spec.json \
  <ROOT_FOLDER_ID> <API_LOGIN> <API_USER_ID> <FOLDER_ID>
```

## The synchronous APIGW pattern

```
Client → APIGW → Wrapper Process → Target Process → API Call {{__callback_url}} → APIGW → Client
```

The wrapper process calls the target via `api_rpc`, which blocks until the target calls `api_rpc_reply`. The reply data is then returned to the client via the `{{__callback_url}}` mechanism.

## MCP tools used

```
get_process_scheme       — read target process inputs/outputs
list_folder_contents     — discover existing processes and folder structure
```

## Knowledge references

- `knowledge/quick-reference.md` — API patterns
- `knowledge/node-types.md` — api_rpc, api_rpc_reply logic schemas
- `templates/api-wrapper.json` — wrapper process template
- `docs/apigw/README.md` — APIGW configuration details

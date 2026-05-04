# Corezoid API Gateway (APIGW)

The Corezoid API Gateway is a reverse proxy that exposes Corezoid processes as HTTP REST endpoints.
External clients call a clean HTTP API; the Gateway forwards requests as tasks into Corezoid,
waits for the process to complete, and returns the result.

## How it works

```
Client → APIGW → Corezoid Process → API Call node → APIGW → Client
                      ↑
               task with __callback_url injected
```

1. Client sends HTTP request to the API Gateway endpoint
2. Gateway injects `__callback_url` into the task data
3. Task is sent to the target Corezoid process
4. Process runs, reaches an API Call node that POSTs to `{{__callback_url}}`
5. Gateway receives the callback and returns the response to the client

This means the **process must contain an API Call node** configured to call `{{__callback_url}}`
for synchronous (request-response) mode to work.

## Creating an API path

An **API path** = one HTTP endpoint backed by one Corezoid process.

### Required parameters

| Parameter | Description |
|-----------|-------------|
| `API_URL` | Corezoid admin URL |
| `APIGW_URL` | API Gateway URL (e.g. `https://apigw.corezoid.com`) |
| `API_TOKEN` | Auth token |
| `API_LOGIN` | Login email |
| `API_USER_ID` | Numeric user ID |
| `WORKSPACE_ID` | Workspace ID |
| `PROJECT_ID` | Project ID |
| `ROOT_FOLDER_ID` | Root folder ID |
| `PROC_ID` | ID of the process to expose |
| `FOLDER_ID` | Folder where the wrapper process will be deployed |
| `HTTP_METHOD` | HTTP verb: GET, POST, PUT, DELETE, PATCH |
| `HTTP_PATH` | URL path: e.g. `/create_user`, `/orders/{id}` |

### Via convctl (recommended)

The `make-api-wrapper` command creates both the wrapper process and registers the API path:

```bash
API_URL=<API_URL> APIGW_URL=<APIGW_URL> API_TOKEN=<API_TOKEN> \
  WORKSPACE_ID=<WORKSPACE_ID> PROJECT_ID=<PROJECT_ID> \
  ./convctl.sh make-api-wrapper \
    .processes/api_wrapper.json \
    .processes/openapi_spec.json \
    <ROOT_FOLDER_ID> <API_LOGIN> <API_USER_ID> <FOLDER_ID>
```

Follow `playbooks/create-openapi-spec.md` for the full step-by-step workflow.

### What the wrapper process does

The `templates/api-wrapper.json` template creates a process that:
1. Receives the incoming HTTP request
2. Calls the target process (`{{PROC_ID}}`) via `api_rpc`
3. Collects the response
4. POSTs back to `{{__callback_url}}` — completing the synchronous response
5. Handles errors with proper HTTP status codes via `Cz-Ag-Status-Code` header

## HTTP methods and routing

| Method | Use case |
|--------|---------|
| `POST` | Create resources, submit data (default) |
| `GET` | Fetch data — parameters sent as query string |
| `PUT` / `PATCH` | Update resources |
| `DELETE` | Remove resources |

**Important:** The HTTP method in your OpenAPI spec must match what the API Call node in the wrapper process uses.

## Synchronous vs asynchronous mode

### Synchronous (default)
- Client waits for response
- Process must have an API Call node that calls `{{__callback_url}}`
- Timeout: default 60 seconds (returns HTTP 500 if exceeded)
- Response: `{"request_proc": "ok", ...response_data}`

### Asynchronous
- Client gets immediate acknowledgment
- Process runs in background
- Response: immediate task ID
- Use when processing takes too long for a client to wait

## Authentication

API Gateway supports three authentication methods:

| Method | How |
|--------|-----|
| **Bearer token** | `Authorization: Bearer <token>` header — standard OAuth 2.0 |
| **API Key** | API key must have "Task management" rights on the target process |
| **OAuth 2.0** | `client_id` + `secret_id` generated via Corezoid Single-Account |

The `templates/api-wrapper.json` uses Bearer token by default (configured in the OpenAPI spec).

## Request data formats

| Content-Type | How body is passed to process |
|-------------|------------------------------|
| `application/json` | Body parsed as JSON object |
| `application/xml` | Body parsed as XML |
| `application/x-www-form-urlencoded` | Body parsed as form fields |

Access raw request body in the process via `{{api_gateway.request.body}}`.
Access request path via `{{api_gateway.request.path}}`, method via `{{api_gateway.request.method}}`.

## Path parameters

Define dynamic path segments in the OpenAPI spec:

```json
"/users/{id}": {
  "get": { ... }
}
```

The `{id}` value is available in the process as a task data field.

## Error responses

The wrapper process returns errors via the `Cz-Ag-Status-Code` header:
- `200` — success
- `4xx` / `5xx` — error, set via `Cz-Ag-Status-Code` header in the API Call node that calls `{{__callback_url}}`

## Domains and SSL

A production APIGW deployment requires three domains with valid SSL certificates:
- `api-apigw.yourdomain.com` — incoming API calls
- `cb-apigw.yourdomain.com` — callback responses
- `*.apigw.yourdomain.com` — user-defined API paths

## Related files

| File | Purpose |
|------|---------|
| `playbooks/create-openapi-spec.md` | Full step-by-step API wrapper creation |
| `templates/api-wrapper.json` | Wrapper process template |
| `docs/nodes/api-call-node.md` | API Call node used in wrapper for callback |
| `docs/nodes/waiting-for-callback-node.md` | Alternative for async patterns |

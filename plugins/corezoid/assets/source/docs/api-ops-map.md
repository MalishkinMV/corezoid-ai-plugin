# Corezoid API Operations Map

All Corezoid management API calls go to `POST <API_URL>/api/2/json` with HMAC-SHA1 auth.
The request body is always `{"ops": [...]}` where each operation specifies `obj` + `obj_type`.

## Authentication

Every request requires three headers:

```
X-API-Login: <your-login>
X-API-Timestamp: <GMT-unix-timestamp>
X-API-Sign: HMAC-SHA1(timestamp + secret + request_body + secret)
```

## Object families and operations

### `conv` — Processes

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Create a process | `title`, `folder_id` |
| `get` | Get process details | `id` |
| `modify` | Update a process | `id` + fields to change |
| `delete` | Soft-delete | `id` |
| `destroy` | Permanent delete | `id` |
| `clone` | Clone a process | `id`, `folder_id` |
| `list` | List processes in folder | `folder_id` |

```json
{
  "ops": [{
    "obj": "conv",
    "obj_type": "create",
    "title": "My Process",
    "folder_id": 12345
  }]
}
```

### `node` — Nodes within a process

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Add a node | `conv_id`, `obj_type` (node type), `condition` |
| `get` | Get node details | `id`, `conv_id` |
| `modify` | Update a node — **full payload required** | `id`, `conv_id`, all fields |
| `delete` | Remove a node | `id`, `conv_id` |
| `list` | List all nodes in a process | `conv_id` |

> **Warning:** `node modify` is NOT partial. Always include `title`, `description`, `logics`, `semaphors`, `extra`, `options`, and `position` — missing fields reset to defaults.

### `task` — Tasks in a process

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Create a task | `conv_id`, `node_id` (start node), `data` |
| `get` | Get task details | `id` |
| `modify` | Update task data | `id`, `data` |
| `delete` | Delete a task | `id` |
| `list` | List tasks in a node | `conv_id`, `node_id` |

### `folder` — Folders

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Create a folder | `title`, `parent_id` |
| `get` | Get folder contents | `id` |
| `modify` | Rename / move | `id` + fields to change |
| `delete` | Soft-delete | `id` |
| `destroy` | Permanent delete | `id` |
| `list` | List subfolders | `parent_id` |

### `dashboard` — Dashboards

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Create a dashboard | `title` |
| `show` | Get dashboard + charts | `id` |
| `modify` | Update dashboard | `id` + fields to change |
| `delete` | Soft-delete (trash) | `id` |
| `restore` | Restore from trash | `id` |
| `destroy` | Permanent delete | `id` |
| `link` | Share with users/groups | `id`, `users`/`groups`/`api_keys` |
| `favorite` | Add to favorites | `id` |

### `chart` — Charts within a dashboard

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Add a chart | `dashboard_id`, `title`, `chart_type`, `metrics` |
| `get` | Get chart details | `id` |
| `modify` | Update chart — **full payload required** | `id`, `obj_type`, full `series` |
| `delete` | Remove a chart | `id` |

Valid `chart_type` values: `column`, `pie`, `funnel`, `table`  
Grid sizing: use `width`/`height` (not `w`/`h`)

> **Warning:** `chart modify` requires `obj_type` + full `series`. Partial modify returns validation errors.

### `company` — Workspace / company operations

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `get` | Get company details | `id` (`i123456789` format) |
| `users` | List workspace members | `id` |

> **Note:** `company_id` must be in `i123456789` format (letter `i` prefix + numeric ID).

### `env_var` — Environment variables

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Create a variable | `title`, `description`, `value`, `folder_id` |
| `get` | Get variable | `id` |
| `modify` | Update variable | `id`, `value` |
| `delete` | Delete variable | `id` |
| `list` | List variables in folder | `folder_id` |

Access in process JSON: `{{env_var[@variable-name]}}`  
Naming: lowercase letters and hyphens only, minimum 3 characters.

### `apigw` — API Gateway paths

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Register API path | `path`, `method`, `conv_id` |
| `get` | Get path details | `id` |
| `modify` | Update path config | `id` + fields to change |
| `delete` | Remove path | `id` |
| `list` | List all API paths | — |

### `commit` — Deployment commits

| `obj_type` | Description | Required fields |
|-----------|-------------|-----------------|
| `create` | Deploy a process version | `conv_id` |
| `list` | List commits for a process | `conv_id` |
| `get` | Get commit details | `id` |

---

## Response format

All API responses follow this structure:

```json
{
  "request_proc": "ok",
  "ops": [
    {
      "id": "optional-client-id",
      "obj": "conv",
      "obj_type": "create",
      "proc": "ok",
      "obj_id": 12345
    }
  ]
}
```

- `request_proc: "ok"` — request was accepted and processed
- `ops[].proc: "ok"` — individual operation succeeded
- `ops[].proc: "error"` — individual operation failed; check `ops[].message` for details
- `ops[].obj_id` — ID of the created/modified object

## Common error messages

| Error | Meaning |
|-------|---------|
| `Key 'obj_type' is required` | `modify` call missing `obj_type` field |
| `Key 'series' is required` | Chart `modify` missing full `series` array |
| `invalid company_id format` | `company_id` not in `i123456789` format |
| `Duplicate REF` | Task REF already exists — use unique REF values |
| `Unknown chart_type` | Chart type not in `column`/`pie`/`funnel`/`table` |

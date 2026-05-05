---
name: corezoid-process-updater
description: >
  Modifies existing Corezoid processes and fixes broken connectors — covers exporting,
  editing process JSON, managing environment variables, and diagnosing failed tasks.
  Activate whenever a user asks to update a process, fix a connector, add a node,
  change how a process works, debug a failed task, or provides ResultNodeObjID and
  ResultTaskData showing where execution stopped. Also activate when the user says
  "add a call to process X", "the task failed with error Y", or "modify the logic".
  This skill knows the correct export→edit→variable→test workflow and error taxonomy.
---

# Corezoid Process Updater


## Bundled References

This Corezoid plugin package includes the upstream Corezoid AI docs repository at `../../assets/source/`. When this skill mentions paths from the original repo, resolve them under that directory unless the user has provided project-local versions. Common examples include `knowledge/`, `docs/`, `templates/`, `playbooks/`, `json-schema/`, `samples/`, `mcp-server/`, and `convctl.sh`.

## CRITICAL: Use MCP tools for all API operations

For reading process state and modifying nodes, always use MCP tools — never write Python scripts:

| Operation | Use this MCP tool |
|-----------|------------------|
| Read process structure | `mcp__corezoid__get_process_scheme` |
| Get process metadata | `mcp__corezoid__get_process_details` |
| Modify a single node | `mcp__corezoid__modify_node` |
| Modify multiple nodes | `mcp__corezoid__batch_modify_nodes` |
| Commit changes | `mcp__corezoid__commit` |
| Create a test task | `mcp__corezoid__create_task` |
| Check task result | `mcp__corezoid__show_task` |
| Trace task execution | `mcp__corezoid__list_task_history` |
| Any raw API op | `mcp__corezoid__make_request` |

`convctl.sh` is still used for fetch/deploy operations (it's the CLI deployment tool, not a Python replacement).

---

## Required parameters

| Parameter | Purpose |
|-----------|---------|
| `API_URL` | Corezoid instance URL |
| `API_TOKEN` | Auth token |
| `WORKSPACE_ID` | Workspace identifier |
| `PROJECT_ID` | Project identifier (needed for variable creation) |
| `ROOT_FOLDER_ID` | Root folder for export |
| `PROC_ID` | Process ID being modified |

---

## Path A — Updating existing logic

### 1. Export the project context

```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh fetch-folder <ROOT_FOLDER_ID> .processes/
```

Analyze the exported processes to understand the current structure and any process IDs referenced.

### 2. Export the specific process to edit

```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh fetch-process <PROC_ID> .processes/target_process.json
```

Work on `.processes/target_process.json` — this is the file to edit.

### 3. Make the changes

Apply the requested modifications directly to `.processes/target_process.json`.

**If adding a Call Process node**, use this structure:
```json
{
  "id": "<new-node-id>",
  "obj_type": 0,
  "condition": {
    "logics": [
      {
        "type": "api_rpc",
        "conv_id": <target_process_id>,
        "err_node_id": "<reply-error-node-id>",
        "extra": {
          "param1": "{{value1}}",
          "param2": "{{value2}}"
        },
        "extra_type": {
          "param1": "string",
          "param2": "number"
        },
        "group": "all"
      },
      {"type": "go", "to_node_id": "<next-node-id>"}
    ],
    "semaphors": []
  },
  "title": "Call <ProcessName>",
  "description": "Calls the X process with parameters",
  "x": 500, "y": 460,
  "extra": "{\"modeForm\":\"expand\",\"icon\":\"\"}",
  "options": null
}
```

**CRITICAL:** `type: "api_rpc"` — never use `"call_process"`. `extra` and `extra_type` must have identical keys with matching types.

### 4. Manage variables (for new constants)

All URLs, tokens, API endpoints, and configuration values must be stored as variables — never hardcoded.

**Check first**: look in `.processes/variables.json` — if the variable already exists, use it.

**Create if missing:**
```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  PROJECT_ID=<PROJECT_ID> \
  ./convctl.sh create-variable <ROOT_FOLDER_ID> <NAME> <DESCRIPTION> <VALUE>
```

Variable naming rules: lowercase letters, numbers, hyphens only — min 3 chars. Good: `api-endpoint`, `auth-token-prod`, `db-host`. Bad: `apiEndpoint`, `token`, `T`.

**Reference in process JSON:** `{{env_var[@variable-name]}}`

### 5. Test

```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh run-process <PROC_ID> .processes/target_process.json '{"key":"value"}'
```

Must pass ("Task completed") before the work is done.

---

## Path B — Fixing a broken connector (error diagnosis)

### 1. Collect error information from the user

Ask for:
- `ResultNodeObjID` — the node ID where execution stopped
- `ResultTaskData` — the full task data object at the time of failure

### 2. Export the current process

```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh fetch-process <PROC_ID> .processes/target_process.json
```

Find the node with `"id": "<ResultNodeObjID>"` in the JSON.

### 3. Diagnose the error

Read `ResultTaskData` and look for these system parameters:

**For API Call failures** (`api` logic):
| Parameter | Meaning |
|-----------|---------|
| `__conveyor_api_return_type_error__` | `"hardware"` = network/infra issue (retryable), `"software"` = logical error |
| `__conveyor_api_return_code__` | HTTP status code returned |
| `__conveyor_api_return_type_tag__` | Error tag: `api_connection_error`, `api_bad_answer`, `api_timeout`, `api_validation_error` |
| `__conveyor_api_return_description__` | Error description string |

**For Code Node failures** (`api_code` logic):
| Parameter | Meaning |
|-----------|---------|
| `__conveyor_code_return_type_error__` | `"hardware"` or `"software"` |
| `__conveyor_code_return_description__` | JS error message (syntax error, runtime exception) |
| `__conveyor_code_return_type_tag__` | `api_code_syntax_error` or `api_code_runtime_error` |

**Common error patterns and fixes:**

| Error tag | Likely cause | Fix |
|-----------|-------------|-----|
| `api_connection_error` | URL unreachable or DNS failure | Check URL variable, verify endpoint availability |
| `api_bad_answer` + `no_scheme` | URL missing `https://` prefix | Fix URL value in env_var or code node |
| `api_bad_answer` (4xx/5xx) | API returned an error | Check request body, headers, auth token |
| `api_timeout` | External API too slow | Increase semaphore timeout or check API |
| `api_code_syntax_error` | JS syntax error in Code node | Fix the `src` field — check for unclosed brackets, typos |
| `api_code_runtime_error` | JS runtime exception | Add null checks, verify data fields exist before use |
| `api_validation_error` | Request didn't match API schema | Check `raw_body` structure against API spec |

### 4. Fix the failing node

Apply the fix to `.processes/target_process.json` and re-test:

```bash
API_URL=<API_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> \
  ./convctl.sh run-process <PROC_ID> .processes/target_process.json '<input_task_data>'
```

Repeat until "Task completed".

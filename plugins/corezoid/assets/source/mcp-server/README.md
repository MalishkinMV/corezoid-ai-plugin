# Corezoid MCP Server

MCP server that exposes the Corezoid API as structured tools for AI agents.

## Prerequisites

- Python 3.14+
- [uv](https://docs.astral.sh/uv/) package manager
- `.env` file in the workspace root with Corezoid credentials

## Setup

```bash
cd corezoid-mcp-server
uv sync
```

## Environment Variables

The server reads from `../.env` (workspace root):

```
API_LOGIN=your_api_login
SECRET=your_secret_key
BASE_URL=https://your-corezoid-instance/api/2/json
COMPANY_ID=your_company_id
```

## Available Tools

### Process Management

| Tool | Description |
|------|-------------|
| `get_process_details` | Get process metadata by ID |
| `get_process_scheme` | Get full process scheme with nodes |
| `list_folder_contents` | List folder contents (0 = root) |
| `modify_node` | Modify node title/description/logics |
| `modify_node_code` | Update JS/Erlang code in a node |
| `batch_modify_nodes` | Modify multiple nodes + commit |
| `commit` | Commit pending changes |

### Task Management

| Tool | Description |
|------|-------------|
| `create_task` | Send a new task to a process (with optional ref) |
| `show_task` | Get task state by task_id and/or ref |
| `modify_task` | Update task data and continue the flow |
| `delete_task` | Delete a task by task_id or ref |
| `list_task_history` | Get execution path (node transitions) for a task |

### Process Copy

| Tool | Description |
|------|-------------|
| `copy_process` | Clone a process (uses /api/2/copy endpoint) |

### Aliases

| Tool | Description |
|------|-------------|
| `list_aliases` | List aliases in a project stage |
| `resolve_aliases` | Resolve alias names to process IDs |
| `create_alias` | Create a new alias |
| `link_alias` | Link alias to a process |
| `create_and_link_alias` | Create + link in one step |

### Dependencies

| Tool | Description |
|------|-------------|
| `crawl_dependencies` | BFS-crawl process dependencies |

### Process Analysis

| Tool | Description |
|------|-------------|
| `analyze_process_structure` | Node counts by type, logic distribution, basic stats |
| `find_nodes_missing_semaphors` | Find nodes missing timeout semaphors (by severity) |
| `extract_code_nodes` | Extract JS/Erlang code with quality flags |
| `list_external_dependencies` | List outbound process references (RPC, Copy, State Store) |
| `find_hardcoded_values` | Scan for hardcoded URLs and numeric conv_ids |
| `find_duplicate_patterns` | Detect duplicate set_param and @send-message patterns |
| `find_untitled_nodes` | Untitled nodes with flow context for naming audit |
| `validate_escalation_chains` | Verify escalation routing correctness |

### Raw API

| Tool | Description |
|------|-------------|
| `make_request` | Send raw authenticated API payload |

## Running

Standalone:
```bash
cd corezoid-mcp-server
uv run python main.py
```

Via Kiro MCP config (`.kiro/settings/mcp.json`):
```json
{
  "mcpServers": {
    "corezoid": {
      "command": "uv",
      "args": ["run", "--directory", "${workspaceFolder}/corezoid-mcp-server", "python", "main.py"],
      "disabled": false
    }
  }
}
```

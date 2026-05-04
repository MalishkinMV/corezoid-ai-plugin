# Corezoid MCP Server

The MCP server exposes the full Corezoid API as named tools that AI assistants call directly.
No manual JSON construction, no auth headers — the server handles everything.

## Authentication (critical — read this first)

The Corezoid API uses **path-based SHA-1 signing**, NOT HTTP headers:

```
URL: {BASE_URL}/{API_LOGIN}/{timestamp}/{signature}
signature = SHA1(timestamp + SECRET + request_body + SECRET)
```

Required credentials:
| Variable | Description | Example |
|----------|-------------|---------|
| `API_LOGIN` | Numeric user ID from Account Settings → API | `65281` |
| `SECRET` | API secret key from Account Settings → API | `abc123...` (long string) |
| `BASE_URL` | Instance API URL — must end with `/api/2/json` | `https://admin.corezoid.com/api/2/json` |
| `COMPANY_ID` | Workspace ID with `i` prefix | `i469b4ea8...` |

> **The JWT token from the browser session is NOT the SECRET.** The secret is a separate static key in Account Settings → API tab.

## Setup

```bash
cd mcp-server
uv sync        # installs mcp[cli] and python-dotenv
```

Or with pip:
```bash
pip install "mcp[cli]>=1.27.0" python-dotenv
```

Create `.env` in the repo root (copy from `.env.example`):
```
API_LOGIN=65281
SECRET=your_api_secret_here
BASE_URL=https://admin.corezoid.com/api/2/json
COMPANY_ID=i469b4ea8...
```

## AI Tool Configuration

### Claude Code / Claude Desktop

`~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "corezoid": {
      "command": "uv",
      "args": ["run", "--directory", "/path/to/corezoid-ai-doc/mcp-server", "python", "main.py"],
      "env": {}
    }
  }
}
```

Or with Python directly (if uv not installed):
```json
{
  "mcpServers": {
    "corezoid": {
      "command": "python",
      "args": ["/path/to/corezoid-ai-doc/mcp-server/main.py"],
      "env": {}
    }
  }
}
```

The `.env` file in the repo root is loaded automatically by the server.

### Cursor

`.cursor/mcp.json` (project-level) or `~/.cursor/mcp.json` (global):

```json
{
  "mcpServers": {
    "corezoid": {
      "command": "uv",
      "args": ["run", "--directory", "${workspaceFolder}/mcp-server", "python", "main.py"],
      "disabled": false
    }
  }
}
```

### Windsurf

`~/.windsurf/mcp.json`:

```json
{
  "mcpServers": {
    "corezoid": {
      "command": "uv",
      "args": ["run", "--directory", "/path/to/corezoid-ai-doc/mcp-server", "python", "main.py"]
    }
  }
}
```

### Kiro

`.kiro/settings/mcp.json`:

```json
{
  "mcpServers": {
    "corezoid": {
      "command": "uv",
      "args": ["run", "--directory", "${workspaceFolder}/mcp-server", "python", "main.py"],
      "disabled": false
    }
  }
}
```

### Multiple Corezoid instances

Run the server with different `.env` files by passing environment variables directly in the MCP config:

```json
{
  "mcpServers": {
    "corezoid-prod": {
      "command": "uv",
      "args": ["run", "--directory", "/path/to/mcp-server", "python", "main.py"],
      "env": {
        "API_LOGIN": "65281",
        "SECRET": "prod_secret",
        "BASE_URL": "https://admin.corezoid.com/api/2/json",
        "COMPANY_ID": "i469b4..."
      }
    },
    "corezoid-dev": {
      "command": "uv",
      "args": ["run", "--directory", "/path/to/mcp-server", "python", "main.py"],
      "env": {
        "API_LOGIN": "65281",
        "SECRET": "dev_secret",
        "BASE_URL": "https://dev.corezoid.com/api/2/json",
        "COMPANY_ID": "i111111..."
      }
    }
  }
}
```

## Available Tools

### Process management
| Tool | Description |
|------|-------------|
| `get_process_details` | Metadata: title, status, change_time (version), project_id, stage_id |
| `get_process_scheme` | Full process JSON: all nodes, logics, semaphors, params |
| `list_folder_contents` | List contents of a folder (folder_id=0 for root) |
| `copy_process` | Clone a process to same or different folder |

### Node modification
| Tool | Description |
|------|-------------|
| `modify_node` | Update title, description, or logics on a single node |
| `modify_node_code` | Replace JS/Erlang source in a code node |
| `batch_modify_nodes` | Apply multiple node modifications then commit in one call |
| `commit` | Commit pending changes to make them live |

### Task operations
| Tool | Description |
|------|-------------|
| `create_task` | Send a new task to a process (with optional unique ref) |
| `show_task` | Get current task state, data, and node by task_id or ref |
| `list_task_history` | Trace the node-by-node execution path of a task |
| `modify_task` | Update task data to continue the flow |
| `delete_task` | Remove a task from a process |

### Alias management
| Tool | Description |
|------|-------------|
| `list_aliases` | List all aliases in a project stage |
| `resolve_aliases` | Resolve @alias-names to process IDs |
| `create_alias` | Create a new alias |
| `link_alias` | Link an alias to a process |
| `create_and_link_alias` | Create + link in one step (auto-generates name from title) |

### Dependency analysis
| Tool | Description |
|------|-------------|
| `crawl_dependencies` | BFS-crawl all api_rpc/api_copy dependencies from a root process |

### Process analysis tools
| Tool | Description |
|------|-------------|
| `analyze_process_structure` | Node counts by type, logic distribution, basic stats |
| `find_nodes_missing_semaphors` | Find API/RPC nodes without timeout semaphors (by severity) |
| `extract_code_nodes` | Extract all JS/Erlang code with quality flags (try_catch, hardcodes) |
| `list_external_dependencies` | All outbound references: RPC, Copy, state store reads |
| `find_hardcoded_values` | Scan for hardcoded URLs and numeric conv_ids |
| `find_duplicate_patterns` | Detect duplicate set_param and call patterns |
| `find_untitled_nodes` | Untitled nodes with flow context for naming audit |
| `validate_escalation_chains` | Verify escalation routing correctness |
| `find_orphaned_nodes` | BFS from Start — detect unreachable (dead) nodes |
| `find_noop_nodes` | Detect no-op conditions and unused set_param nodes |
| `check_variable_usage` | Check if variables are referenced anywhere in a process |

### Raw API
| Tool | Description |
|------|-------------|
| `make_request` | Send any raw authenticated `ops` payload to the Corezoid API |

## Critical gotchas

**Always batch operations.** The API accepts multiple ops in one request. Never make two separate calls when one batched call works.

**`get_process_scheme` response path:**
```python
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']
#                               ^-- this is a LIST, index [0] first
```

**`modify_node` requires the full node payload** — title, description, logics, semaphors, extra, options, position. Partial updates silently reset missing fields.

**Process JSON for import must be a JSON array** `[{...}]`, not a bare object. Every node needs `uuid` (UUID v4). Use `position: [x, y]` not separate `x`/`y`.

**`company_id`** must appear in every op payload in `make_request` calls.

**After `create_task`**: use the returned `obj_id` as `task_id` for all subsequent lookups. If you pass a `ref`, you can look up by ref — but ref must be unique within the process.

## Python client (standalone use)

`corezoid_client.py` in the repo root can be used directly without the MCP server:

```python
from corezoid_client import CorezoidClient

client = CorezoidClient(
    api_login='65281',
    secret='your_secret',
    base_url='https://admin.corezoid.com/api/2/json',
    company_id='i469b4ea8...'
)

# Test connection
result = client.list_folder_contents(YOUR_FOLDER_ID)
print(result)

# Create a task
result = client.create_task(12345, data={"city": "Kyiv"}, ref="test-001")
task_id = result['ops'][0]['obj_id']

# Check result
result = client.show_task(12345, task_id=task_id)
print(result['ops'][0]['data'])
```

## Related files

| File | Purpose |
|------|---------|
| `corezoid_client.py` | Python client — use directly or via MCP server |
| `.env.example` | Credential template — copy to `.env` |
| `docs/api-ops-map.md` | Complete API ops reference |
| `playbooks/create-connector.md` | Building a connector process |
| `agents/` | Tool-agnostic workflow definitions for all agents |

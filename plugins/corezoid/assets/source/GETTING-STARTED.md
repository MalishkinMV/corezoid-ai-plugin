# Getting Started with Corezoid AI Console

This guide is for any AI coding assistant (or human) beginning to work with this repository.
Read this before touching any playbook or generating any process JSON.

---

## What this repository is

The **knowledge base and tooling** for the Corezoid AI Console — an AI assistant that helps users create, modify, validate, and document Corezoid processes. It contains:

- **`corezoid_client.py`** — Python API client (primary tool — handles all auth automatically)
- **`mcp-server/`** — MCP server wrapping the client as AI-callable tools
- **Playbooks** — step-by-step workflows for common tasks
- **Documentation** — node types, process structure, error handling
- **`knowledge/`** — single source of truth: auth, node types, gotchas, quick-reference, validation checklist
- **JSON schemas** — for validating process JSON
- **Templates** — reusable process skeletons
- **Samples** — real working processes to reference

---

## Step 1 — Obtain your credentials

You need these from the Corezoid UI → **Account Settings → API** tab:

| Credential | Description | Example |
|-----------|-------------|---------|
| `API_LOGIN` | Numeric user ID | `65281` |
| `SECRET` | API secret key (long static string) | `abc123xyz...` |
| `BASE_URL` | API endpoint URL | `https://admin.corezoid.com/api/2/json` |
| `COMPANY_ID` | Workspace ID with `i` prefix | `i469b4ea8...` |

> **Important:** The `SECRET` is NOT the JWT token visible in your browser. It's a separate static key in Account Settings → API.

Additional IDs you'll need for specific tasks:

| Parameter | When needed | How to find |
|-----------|-------------|-------------|
| `FOLDER_ID` | Deploying processes | From project URL: `.../stage/<FOLDER_ID>` |
| `PROC_ID` | Testing/running processes | From process URL or `list_folder_contents()` |
| `APIGW_URL` | API wrapper creation | From your Corezoid admin |
| `API_USER_ID` | API wrapper creation | Same as `API_LOGIN` numeric value |

---

## Step 2 — Configure credentials

```bash
cp .env.example .env
# Edit .env and fill in your API_LOGIN, SECRET, BASE_URL, COMPANY_ID
```

---

## Step 3 — Install Python dependencies

The Python client (`corezoid_client.py`) has no external dependencies — it uses only the standard library.

For the MCP server:
```bash
cd mcp-server && uv sync
```

Or: `pip install "mcp[cli]>=1.27.0" python-dotenv`

---

## Step 4 — Verify access

```python
from corezoid_client import CorezoidClient
from dotenv import load_dotenv
import os

load_dotenv()
client = CorezoidClient(
    api_login=os.environ['API_LOGIN'],
    secret=os.environ['SECRET'],
    base_url=os.environ['BASE_URL'],
    company_id=os.environ['COMPANY_ID']
)

# Test: list your project folder
result = client.list_folder_contents(int(os.environ.get('FOLDER_ID', 0)))
print(result)
```

---

## Step 5 — Install Claude Code skills (Claude Code users)

If you use Claude Code, install the skills so Claude knows the full Corezoid workflow without you having to explain it:

```bash
cp -r skills/corezoid-* ~/.claude/skills/
```

Then restart Claude Code. You can trigger any skill by describing what you want — e.g. "build a process that calls the Stripe API" or "review this process JSON".

See [`skills/README.md`](skills/README.md) for the full list and trigger phrases.

---

## Step 6 — Connect your AI tool to the MCP server (optional but recommended)

The MCP server gives your AI assistant direct access to all Corezoid operations as callable tools. Configure it once and the AI can create, modify, test, and review processes without manual API calls.

See `docs/mcp/README.md` for configuration for Claude Code, Cursor, Windsurf, Kiro.

---

## Step 6 — Understand the project structure

```
corezoid_client.py      ← Python API client (primary tool)
mcp-server/             ← MCP server (AI-callable tools)
  server.py             ← Tool definitions
  main.py               ← Entry point
  pyproject.toml        ← Dependencies
.env.example            ← Credential template
knowledge/              ← SSOT: auth, node types, gotchas, quick-reference, validation
agents/                 ← Tool-agnostic workflow definitions (8 agents)
skills/                 ← Claude Code skills (install with cp -r skills/corezoid-* ~/.claude/skills/)
playbooks/              ← Step-by-step guides for specific scenarios
docs/
  nodes/                ← Reference for all node types
  process/              ← Process structure, validation, error handling
  mcp/README.md         ← MCP server setup guide
  api-ops-map.md        ← Complete API operations reference
  corezoid-api.json     ← Full OpenAPI spec for Corezoid API
json-schema/            ← JSON Schema files for validation
templates/              ← Ready-to-copy process templates
samples/                ← Real working processes (use as reference)
.processes/             ← Your working directory (git-ignored, auto-created)
```

---

## Step 7 — Choose your workflow

| What you want to do | Use this |
|--------------------|----------|
| Design a multi-process system | `agents/architect.md` |
| Build a new connector to an external API | `agents/process-builder.md` + `playbooks/create-connector.md` |
| Build business logic using existing processes | `agents/process-builder.md` + `playbooks/create-logic.md` |
| Create from template | `playbooks/create-connector-from-template.md` |
| Expose a process as an HTTP endpoint | `agents/apigw-manager.md` + `playbooks/create-openapi-spec.md` |
| Modify an existing process | `agents/process-updater.md` + `playbooks/update-logic.md` |
| Fix a broken connector | `agents/process-updater.md` + `playbooks/fix-connector.md` |
| Review a process for issues | `agents/process-reviewer.md` + `playbooks/review-process.md` |
| Validate process JSON before import | `knowledge/validation-checklist.md` |
| Create monitoring dashboards | `agents/dashboard-manager.md` |
| Generate documentation for a process | `agents/process-tech-writer.md` |

---

## Key rules (memorize these)

1. **Auth is path-based** — `SECRET` + SHA1, not HTTP headers. Use `corezoid_client.py`, not raw curl.
2. **Process JSON for import must be an array** `[{...}]` — bare objects are rejected
3. **Every node needs `err_node_id`** — the #1 cause of broken processes
4. **Every API call node needs a semaphor** — without it, tasks hang forever on timeout
5. **Use `set_param` not code nodes** for variable copy, string concat, arithmetic
6. **Node modify requires full payload** — partial updates silently reset fields
7. **`company_id`** must appear in every `ops` entry for most API calls
8. **Test before done** — `create_task()` → `show_task()` — check the `data` field

---

## Common operations quick reference

```python
from corezoid_client import CorezoidClient
client = CorezoidClient(api_login=..., secret=..., base_url=..., company_id=...)

# List folder contents
client.list_folder_contents(FOLDER_ID)

# Get process details (metadata + change_time = version)
client.get_process_details(PROC_ID)

# Get full process scheme
result = client.get_process_scheme(PROC_ID)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']  # scheme is a LIST

# Create and test a task
r = client.create_task(PROC_ID, data={"city": "Kyiv"}, ref="test-001")
task_id = r['ops'][0]['obj_id']
client.show_task(PROC_ID, task_id=task_id)

# Trace task execution path
client.list_task_history(PROC_ID, task_id)

# Analyze process
client.analyze_process_structure(PROC_ID)
client.find_nodes_missing_semaphors(PROC_ID)
client.find_hardcoded_values(PROC_ID)
client.find_orphaned_nodes(PROC_ID)

# Raw API call (batch multiple ops)
client.make_request({'ops': [{...}, {...}]})
```

---

## If something goes wrong

- Check `docs/error_explanations.md` for known error codes and solutions
- Most common errors: `403 cookie or headers are not valid` → wrong auth pattern or expired SESSION token (use SECRET, not JWT)
- `import rejected` → process JSON is a bare object, not an array; or missing `uuid` on nodes
- Tasks stuck in a node → missing semaphor on API call node
- Silent task drops → duplicate REF values

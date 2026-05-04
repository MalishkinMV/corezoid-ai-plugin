# Corezoid AI Doc

Knowledge base and tooling for AI-assisted Corezoid process development. Provides documentation, agent workflow definitions, an MCP server, and ready-to-use configurations for Claude Code, Cursor, Windsurf, Kiro, GitHub Copilot, and Devin.

## What's in this repo

```
knowledge/      Single source of truth — auth, node types, gotchas, quick-reference
agents/         Tool-agnostic workflow definitions for 8 AI agents
skills/         Claude Code skills — install these for AI-assisted process development
mcp-server/     Python MCP server exposing 35+ Corezoid API tools
docs/           Full reference — all node types, process structure, APIGW, dashboards
playbooks/      Step-by-step guides for common tasks
samples/        Real working process JSON examples
templates/      Reusable process skeletons
json-schema/    JSON Schema files for process and node validation
scripts/        generate-tool-configs.py — regenerates all tool configs from knowledge/
```

## Quick start

**1. Get credentials** from Corezoid UI → Account Settings → API:

```
API_LOGIN   numeric user ID
SECRET      long static string (NOT the JWT session token)
BASE_URL    https://www.corezoid.com/api/2/json
COMPANY_ID  your workspace UUID
FOLDER_ID   target folder ID
```

**2. Configure credentials:**

```bash
cp .env.example .env
# fill in your values
```

**3. Connect your AI tool to the MCP server** — see [`docs/mcp/README.md`](docs/mcp/README.md)

**4. Read [`GETTING-STARTED.md`](GETTING-STARTED.md)** for the full onboarding guide.

## Claude Code Skills

Install the skills for AI-assisted Corezoid development in Claude Code:

```bash
cp -r skills/corezoid-* ~/.claude/skills/
```

| Skill | Purpose |
|-------|---------|
| [corezoid-architect](skills/corezoid-architect/SKILL.md) | Design multi-process systems |
| [corezoid-process-builder](skills/corezoid-process-builder/SKILL.md) | Build new connector or logic process |
| [corezoid-process-reviewer](skills/corezoid-process-reviewer/SKILL.md) | Audit process before deployment |
| [corezoid-process-updater](skills/corezoid-process-updater/SKILL.md) | Modify process or debug failed tasks |
| [corezoid-process-tech-writer](skills/corezoid-process-tech-writer/SKILL.md) | Generate docs and enrich process JSON |
| [corezoid-apigw-manager](skills/corezoid-apigw-manager/SKILL.md) | Expose process as HTTP endpoint |
| [corezoid-api-wrapper](skills/corezoid-api-wrapper/SKILL.md) | Alternative HTTP wrapper via convctl |
| [corezoid-dashboard-manager](skills/corezoid-dashboard-manager/SKILL.md) | Create monitoring dashboards |

Full setup: [`skills/README.md`](skills/README.md)

## AI tool configs

Auto-generated from `knowledge/` by `scripts/generate-tool-configs.py`:

| Tool | Config |
|------|--------|
| Claude Code | `CLAUDE.md` |
| Cursor | `.cursor/rules/corezoid.mdc` |
| Windsurf | `.windsurf/rules/corezoid.md` |
| GitHub Copilot | `.github/copilot-instructions.md` |
| Kiro | `.kiro/steering/corezoid.md` |
| Devin | `.devin.md` |
| MCP (any tool) | `.mcp.json.example` → copy to `.mcp.json` and fill credentials |

To regenerate after updating `knowledge/` or `agents/`:
```bash
python scripts/generate-tool-configs.py
```

## Agents

Eight workflow agents — each defines when to activate, what inputs are needed, which MCP tools to call, and which knowledge files to consult:

| Agent | When to use |
|-------|------------|
| [architect](agents/architect.md) | Design a multi-process system, choose patterns, plan decomposition |
| [process-builder](agents/process-builder.md) | Build a new connector or logic process from scratch |
| [process-reviewer](agents/process-reviewer.md) | Audit a process for structural errors before deployment |
| [process-updater](agents/process-updater.md) | Modify an existing process or debug a failed task |
| [process-tech-writer](agents/process-tech-writer.md) | Generate documentation and enrich process JSON with descriptions |
| [apigw-manager](agents/apigw-manager.md) | Expose a process as an HTTP REST endpoint via API Gateway |
| [api-wrapper-manager](agents/api-wrapper-manager.md) | Alternative HTTP wrapper workflow using convctl |
| [dashboard-manager](agents/dashboard-manager.md) | Create dashboards to monitor process task counters |

## MCP server

The MCP server wraps the Python client as callable tools for any AI assistant that supports MCP:

```bash
cd mcp-server && uv sync
```

Then configure your AI tool using `.mcp.json.example` as a template. Full setup: [`docs/mcp/README.md`](docs/mcp/README.md).

## Python client

`corezoid_client.py` handles all API operations — auth signing, request formatting, response parsing:

```python
from corezoid_client import CorezoidClient

client = CorezoidClient(
    api_login='YOUR_API_LOGIN',
    secret='YOUR_SECRET',
    base_url='https://www.corezoid.com/api/2/json',
    company_id='YOUR_COMPANY_ID'
)

result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']
```

## Knowledge files

| File | Contents |
|------|---------|
| [`knowledge/auth.md`](knowledge/auth.md) | SHA-1 auth pattern, credentials, batching |
| [`knowledge/node-types.md`](knowledge/node-types.md) | All logic types with full JSON examples |
| [`knowledge/process-schema.md`](knowledge/process-schema.md) | Process and node JSON structure |
| [`knowledge/gotchas.md`](knowledge/gotchas.md) | 15 silent failures — things that break without obvious errors |
| [`knowledge/quick-reference.md`](knowledge/quick-reference.md) | Copy-paste snippets for common operations |
| [`knowledge/validation-checklist.md`](knowledge/validation-checklist.md) | Pre-deployment checklist for process JSON |

## License

MIT — see [LICENSE](LICENSE).

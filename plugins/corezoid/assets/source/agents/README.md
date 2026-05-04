# Corezoid Agents

Tool-agnostic workflow definitions for all Corezoid AI workflows. Each agent definition specifies what the agent does, when to use it, what inputs it needs, which MCP tools it calls, and which knowledge files it references.

These definitions work with any AI framework that supports the MCP server (`mcp-server/`). They are also the source for AI-tool-specific skill files (Claude Code skills in `~/.claude/skills/corezoid-*/`).

## Agents

| Agent | File | Purpose |
|-------|------|---------|
| Architect | `architect.md` | Design multi-process systems — patterns, decomposition, project structure |
| Process Builder | `process-builder.md` | Build a new connector or logic process from scratch |
| Process Reviewer | `process-reviewer.md` | Audit process for errors before deployment |
| Process Updater | `process-updater.md` | Modify existing process or debug failed task |
| Process Tech Writer | `process-tech-writer.md` | Generate documentation and enrich process JSON |
| APIGW Manager | `apigw-manager.md` | Expose process as HTTP REST endpoint via Corezoid API Gateway |
| API Wrapper Manager | `api-wrapper-manager.md` | Alternative HTTP wrapper workflow using convctl |
| Dashboard Manager | `dashboard-manager.md` | Create dashboards to visualize process metrics |

All agents are for **building and managing Corezoid processes** — connectors, logic flows, API endpoints, dashboards, and documentation.

## How to use these definitions

### With the MCP server (any AI tool)

1. Configure the MCP server (see `docs/mcp/README.md`)
2. Reference the relevant agent definition when instructing your AI assistant
3. The AI will use the MCP tools listed in each agent definition

### With Claude Code

Skills in `~/.claude/skills/corezoid-*/` are Claude Code-specific implementations of these agent definitions. They activate automatically based on trigger phrases.

### Direct Python (no AI)

```python
from corezoid_client import CorezoidClient
client = CorezoidClient(api_login=..., secret=..., base_url=..., company_id=...)
# Use any method from the client directly
```

## Adding a new agent

1. Create `agents/<name>.md` with sections: Purpose, When to invoke, Required inputs, MCP tools used, Workflow, Knowledge references
2. Add it to the tables above
3. Run `scripts/generate-tool-configs.py` to regenerate AI-tool-specific configs
4. Create the corresponding skill in `~/.claude/skills/corezoid-<name>/SKILL.md`

# Claude Code Skills

These are Claude Code skills for working with Corezoid. Each skill encodes a complete workflow — when to activate, what to ask for, and exactly how to execute.

## Installation

Copy the skills you want into your Claude Code global skills directory:

```bash
# macOS / Linux
cp -r skills/corezoid-* ~/.claude/skills/

# Then restart Claude Code (or open a new session)
```

Or install all at once:

```bash
for skill in skills/corezoid-*/; do
  cp -r "$skill" ~/.claude/skills/
done
```

## Available Skills

| Skill | Trigger phrases | Purpose |
|-------|----------------|---------|
| `corezoid-architect` | "design a system", "plan the processes", "architect this" | Design multi-process systems before building |
| `corezoid-process-builder` | "build a process", "create a connector", "make a process that calls API X" | Build new connector or logic process from scratch |
| `corezoid-process-reviewer` | "review this process", "check the process", "audit the connector" | Audit process JSON for issues before deployment |
| `corezoid-process-updater` | "update the process", "fix the connector", "add a node", "the task failed" | Modify existing process or debug a failed task |
| `corezoid-process-tech-writer` | "document this process", "write docs for the connector" | Generate Markdown docs + enrich process JSON with descriptions |
| `corezoid-apigw-manager` | "expose as HTTP", "create API endpoint", "make this callable via REST" | Expose a process as HTTP endpoint via API Gateway |
| `corezoid-api-wrapper` | "create API wrapper", "make-api-wrapper", "generate openapi spec" | Alternative HTTP wrapper workflow using convctl |
| `corezoid-dashboard-manager` | "create a dashboard", "monitor this process", "add charts" | Create monitoring dashboards with charts |

## Recommended order

For a new system:
1. **architect** — design first, agree on the process map
2. **process-builder** — implement each process
3. **process-reviewer** — quality gate before deployment
4. **process-tech-writer** — generate documentation
5. **apigw-manager** or **api-wrapper** — expose as HTTP if needed
6. **dashboard-manager** — add monitoring

For an existing process:
- **process-updater** — to modify or debug
- **process-reviewer** — to audit
- **process-tech-writer** — to document

## How skills work in Claude Code

Skills are Markdown files that Claude Code loads when a matching trigger phrase is detected. The skill provides the full workflow — what parameters to collect, what steps to follow, what gotchas to avoid. You don't need to explain Corezoid to Claude; the skill handles that.

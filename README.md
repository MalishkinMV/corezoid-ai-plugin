# Corezoid AI Plugin Marketplace

This repository publishes the **Corezoid** plugin for both **Claude Code** and **Codex**. It packages Corezoid product knowledge, public docs, workflow skills, schemas, templates, playbooks, samples, and MCP reference material so AI coding tools can help design, build, review, document, and operate Corezoid processes.

## Marketplace Layout

```text
.agents/plugins/marketplace.json          # Codex marketplace catalog
.claude-plugin/marketplace.json           # Claude Code marketplace catalog
plugins/corezoid/.codex-plugin/           # Codex plugin manifest
plugins/corezoid/.claude-plugin/          # Claude Code plugin manifest
plugins/corezoid/skills/                  # Corezoid workflow skills shared by both platforms
plugins/corezoid/assets/                  # Icons, screenshots, public profile, bundled docs
plugins/corezoid/.mcp.json                # Empty by default; credentials are user-specific
```

## Install In Claude Code

```bash
claude plugin marketplace add corezoid/corezoid-ai-plugin
claude plugin install corezoid@corezoid
```

For local testing from this repository root:

```bash
claude plugin marketplace add ./
claude plugin install corezoid@corezoid
```

## Install In Codex

For the stable release:

```bash
codex plugin marketplace add corezoid/corezoid-ai-plugin --ref v1.1.0
codex plugin marketplace upgrade corezoid
```

For development tracking from `main`:

```bash
codex plugin marketplace add corezoid/corezoid-ai-plugin --ref main
codex plugin marketplace upgrade corezoid
```

For local testing from this repository root:

```bash
codex plugin marketplace add ./
codex plugin marketplace upgrade corezoid
```

Restart Codex, open the Plugin Directory, choose the **Corezoid** marketplace, and install **Corezoid**.

## Shared Skills

Both platforms expose the same Corezoid workflow skills:

| Skill | Purpose |
|-------|---------|
| `corezoid` | Product overview, reference map, and skill router |
| `corezoid-architect` | Design multi-process Corezoid systems |
| `corezoid-process-builder` | Build connector or logic process JSON |
| `corezoid-process-reviewer` | Audit process JSON before deployment |
| `corezoid-process-updater` | Modify existing processes or debug failed tasks |
| `corezoid-process-tech-writer` | Create Markdown docs and enrich process JSON |
| `corezoid-apigw-manager` | Expose Corezoid processes through API Gateway |
| `corezoid-api-wrapper` | Generate OpenAPI specs and wrapper processes |
| `corezoid-dashboard-manager` | Plan and create monitoring dashboards |

## Upstream Sources

The bundled technical corpus comes from the public repository `corezoid/corezoid-ai-doc`, recorded in `plugins/corezoid/assets/source-metadata.json` and indexed in `plugins/corezoid/assets/source-file-index.txt`.

Public product positioning and source links are captured in `plugins/corezoid/assets/public/`.

## Security Model

The plugin does not bundle credentials. `plugins/corezoid/.mcp.json` is intentionally empty by default, and any Corezoid API login, secret, workspace ID, or environment-specific MCP configuration must be added by the user locally.

## Publishing Status

| Platform | Status |
|----------|--------|
| Claude Code | Available via `claude plugin marketplace add corezoid/corezoid-ai-plugin` |
| Codex | Available via `codex plugin marketplace add corezoid/corezoid-ai-plugin --ref v1.1.0` |

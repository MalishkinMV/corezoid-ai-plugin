# Corezoid Plugin Marketplace

This repository publishes the **Corezoid** plugin for both **Codex** and **Claude Code**. The plugin packages Corezoid product knowledge, public docs, workflow skills, schemas, templates, playbooks, samples, and MCP reference material so AI coding tools can help design, build, review, document, and operate Corezoid processes.

## Marketplace Layout

```text
.agents/plugins/marketplace.json          # Codex marketplace catalog
.claude-plugin/marketplace.json           # Claude Code marketplace catalog
plugins/corezoid/.codex-plugin/           # Codex plugin manifest
plugins/corezoid/.claude-plugin/          # Claude Code plugin manifest
plugins/corezoid/skills/                  # Corezoid workflow skills (shared)
plugins/corezoid/assets/                  # Icons, screenshots, public profile, bundled docs
plugins/corezoid/.mcp.json                # Empty by default; credentials are user-specific
```

---

## Claude Code

### Install From GitHub

```bash
claude plugin marketplace add corezoid/corezoid-codex-plugin
claude plugin install corezoid@corezoid
```

### Local Test

```bash
claude plugin marketplace add ./
claude plugin install corezoid@corezoid
```

### Slash Commands

After installation, the following skills are available in Claude Code:

- `/corezoid` — product overview, reference map, and skill router
- `/corezoid-architect` — design multi-process Corezoid systems
- `/corezoid-process-builder` — build connector or logic process JSON
- `/corezoid-process-reviewer` — audit process JSON before deployment
- `/corezoid-process-updater` — modify existing processes or debug failed tasks
- `/corezoid-process-tech-writer` — create Markdown docs and enrich process JSON
- `/corezoid-apigw-manager` — expose Corezoid processes through API Gateway
- `/corezoid-api-wrapper` — generate OpenAPI specs and wrapper processes
- `/corezoid-dashboard-manager` — plan and create monitoring dashboards

---

## Codex

### Install From GitHub

```bash
codex plugin marketplace add corezoid/corezoid-codex-plugin --ref main
codex plugin marketplace upgrade corezoid
```

Then open the Codex Plugin Directory, choose the **Corezoid** marketplace, and install **Corezoid**.

For a pinned release, use a tag:

```bash
codex plugin marketplace add corezoid/corezoid-codex-plugin --ref v1.0.2
```

### Local Test

```bash
codex plugin marketplace add ./
codex plugin marketplace upgrade corezoid
```

Restart Codex, open the Plugin Directory, choose **Corezoid**, and verify that the plugin card renders with its icon, logo, screenshots, prompts, and skills.

---

## Plugin Contents

Both plugins include the same set of workflow skills:

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

The bundled technical corpus comes from the public repository `corezoid/corezoid-ai-doc`, recorded in `plugins/corezoid/assets/source-metadata.json`.

Public product positioning and source links are captured in `plugins/corezoid/assets/public/`.

## Publishing Status

| Platform | Status |
|----------|--------|
| Claude Code | Available — install via `claude plugin marketplace add corezoid/corezoid-codex-plugin` |
| Codex | Available — install via `codex plugin marketplace add corezoid/corezoid-codex-plugin` |

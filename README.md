# Corezoid Codex Marketplace

This repository publishes the **Corezoid** plugin for Codex. The plugin packages Corezoid product knowledge, public docs, workflow skills, schemas, templates, playbooks, samples, and MCP reference material so Codex can help design, build, review, document, and operate Corezoid processes.

## Marketplace Layout

```text
.agents/plugins/marketplace.json      # Codex marketplace catalog
plugins/corezoid/.codex-plugin/       # Required plugin manifest
plugins/corezoid/skills/              # Corezoid workflow skills
plugins/corezoid/assets/              # Icons, screenshots, public profile, bundled docs
plugins/corezoid/.mcp.json            # Empty by default; credentials are user-specific
```

## Install From GitHub

After this repository is pushed to GitHub, users can add the marketplace with:

```bash
codex plugin marketplace add corezoid/corezoid-codex-plugin --ref main
codex plugin marketplace upgrade corezoid
```

Then open Codex Plugin Directory, choose the **Corezoid** marketplace, and install **Corezoid**.

For a pinned release, use a tag:

```bash
codex plugin marketplace add corezoid/corezoid-codex-plugin --ref v1.0.2
```

## Local Test

From this repository root:

```bash
codex plugin marketplace add ./
codex plugin marketplace upgrade corezoid
```

Restart Codex, open the Plugin Directory, choose **Corezoid**, and verify that the plugin card renders with its icon, logo, screenshots, prompts, and skills.

## Plugin Contents

The Corezoid plugin includes:

- `corezoid`: product overview, reference map, and skill router.
- `corezoid-architect`: design multi-process Corezoid systems.
- `corezoid-process-builder`: build connector or logic process JSON.
- `corezoid-process-reviewer`: audit process JSON before deployment.
- `corezoid-process-updater`: modify existing processes or debug failed tasks.
- `corezoid-process-tech-writer`: create Markdown docs and enrich process JSON.
- `corezoid-apigw-manager`: expose Corezoid processes through API Gateway.
- `corezoid-api-wrapper`: generate OpenAPI specs and wrapper processes.
- `corezoid-dashboard-manager`: plan and create monitoring dashboards.

## Upstream Sources

The bundled technical corpus comes from the public repository `corezoid/corezoid-ai-doc`, recorded in `plugins/corezoid/assets/source-metadata.json`.

Public product positioning and source links are captured in `plugins/corezoid/assets/public/`.

## Publishing Status

Codex supports custom marketplaces from local paths and Git repositories. Official OpenAI Plugin Directory self-serve publishing is not yet generally available, so this repo is the distributable Corezoid marketplace source until that opens.

# Corezoid Codex Plugin

Corezoid is a product-focused Codex plugin for working with Corezoid Actor Engine: process architecture, process JSON, API connectors, API Gateway wrappers, dashboards, reviews, updates, documentation, and public reference lookup.

## What This Plugin Contains

- `skills/` exposes focused Corezoid workflow skills for Codex.
- `assets/source/` contains the full public `corezoid/corezoid-ai-doc` repository, excluding only `.git` metadata.
- `assets/public/` contains the product profile and public source links used for the marketplace-ready plugin card.
- `assets/source-file-index.txt` lists every bundled upstream source file for quick lookup.
- `.mcp.json` is intentionally empty so the plugin does not start a Corezoid MCP server without user credentials.

## Skills

- `corezoid`: product overview, reference map, and skill router.
- `corezoid-architect`: design multi-process Corezoid systems.
- `corezoid-process-builder`: build connector or logic process JSON.
- `corezoid-process-reviewer`: audit a process JSON before deployment.
- `corezoid-process-updater`: modify existing processes or debug failed tasks.
- `corezoid-process-tech-writer`: create Markdown docs and enrich process JSON.
- `corezoid-apigw-manager`: expose Corezoid processes through API Gateway.
- `corezoid-api-wrapper`: generate OpenAPI specs and wrapper processes.
- `corezoid-dashboard-manager`: plan and create monitoring dashboards.

## Main Bundled References

- `assets/source/knowledge/` for auth, gotchas, node types, process schema, quick reference, and validation.
- `assets/source/docs/` for node, process, task, APIGW, dashboard, MCP, and API reference documentation.
- `assets/source/json-schema/` for validation schemas.
- `assets/source/templates/` for reusable process, connector, and wrapper templates.
- `assets/source/playbooks/` for step-by-step implementation guides.
- `assets/source/samples/` for working sample process JSON.
- `assets/public/corezoid-profile.md` for public product positioning and marketplace metadata notes.

## MCP Note

The upstream MCP server and example config are bundled at `assets/source/mcp-server/` and `assets/source/.mcp.json.example`. Configure credentials before enabling any Corezoid MCP server locally.

## Marketplace Status

This plugin is published by the `corezoid` marketplace as `corezoid@corezoid` with installation policy `AVAILABLE`.

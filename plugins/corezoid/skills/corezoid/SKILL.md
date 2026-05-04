---
name: corezoid
description: >
  Universal Corezoid product assistant for Codex. Use when a user needs Corezoid Actor
  Engine guidance, process architecture, process JSON, API orchestration, API Gateway,
  dashboards, validation, debugging, public documentation, schemas, templates, samples,
  playbooks, or help choosing one of the bundled Corezoid workflow skills.
---

# Corezoid

This Codex plugin packages public Corezoid product knowledge and the full `corezoid/corezoid-ai-doc` technical corpus.

## Product Scope

Use this skill as the entry point for Corezoid product work:

- Process architecture and decomposition.
- Connector and logic process creation.
- External API orchestration, API Gateway, webhooks, and API wrappers.
- Process JSON validation, review, update, and debugging.
- Process documentation and enrichment.
- Dashboards, activity monitoring, task counters, and operational visibility.
- Navigation across public Corezoid docs, tutorials, schemas, templates, playbooks, and samples.

## Bundled Sources

The upstream technical repository is copied into this plugin at:

`../../assets/source/`

Public product metadata is stored at:

`../../assets/public/corezoid-profile.md`

Use these paths when this skill or any bundled Corezoid skill refers to source-repo paths such as:

- `knowledge/`
- `docs/`
- `json-schema/`
- `templates/`
- `playbooks/`
- `samples/`
- `agents/`
- `mcp-server/`
- `convctl.sh`

A complete source file list is available at `../../assets/source-file-index.txt`.

## Skill Map

Use the focused bundled skills when the user intent matches them:

- `corezoid-architect`: design a multi-process Corezoid system.
- `corezoid-process-builder`: build connector or logic process JSON.
- `corezoid-process-reviewer`: audit a process JSON before deployment.
- `corezoid-process-updater`: modify an existing process or debug failed tasks.
- `corezoid-process-tech-writer`: create Markdown docs and enrich process JSON.
- `corezoid-apigw-manager`: expose a Corezoid process through API Gateway.
- `corezoid-api-wrapper`: generate an OpenAPI spec and wrapper process.
- `corezoid-dashboard-manager`: create monitoring dashboards and charts.

## Reference Order

When working from the bundled docs, prefer these sources first:

- `../../assets/public/corezoid-profile.md` for product framing.
- `../../assets/source/knowledge/quick-reference.md` for common API and process snippets.
- `../../assets/source/knowledge/gotchas.md` for silent failure modes.
- `../../assets/source/knowledge/node-types.md` for logic schemas and examples.
- `../../assets/source/knowledge/process-schema.md` for process structure.
- `../../assets/source/knowledge/validation-checklist.md` before deployment.
- `../../assets/source/docs/` for full reference material.
- `../../assets/source/templates/` and `../../assets/source/samples/` for starting JSON.

## MCP Note

The plugin manifest points to an empty `.mcp.json` so it does not auto-start a Corezoid MCP server without user credentials. The upstream server and example config are bundled at `../../assets/source/mcp-server/` and `../../assets/source/.mcp.json.example`.

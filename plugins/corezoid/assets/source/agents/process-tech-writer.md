# Agent: process-tech-writer

## Purpose

Generate human-readable documentation for a Corezoid process, and enrich the process JSON itself with `description` fields on every node. Always produces both outputs together.

## When to invoke

- User asks to "document", "write docs for", "describe", "add descriptions to" a process
- User shares a process JSON and asks what it does
- User wants team-readable documentation for a deployed process
- Before publishing a process to a shared workspace

## Required inputs

| Input | When needed |
|-------|------------|
| Process JSON or process ID | Always |
| `API_LOGIN`, `SECRET`, `BASE_URL`, `COMPANY_ID` | Only if fetching from live system |

## MCP tools used

```
get_process_scheme    — fetch process to document
get_process_details   — get process title, description, params
```

## Two outputs (always produced together)

### Output 1: Markdown documentation

Save to `.processes/<process-name-in-snake-case>-docs.md`.

```markdown
# <Process Title>

## Overview
<1-2 sentences: what this process does and when to call it>

## Input Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| field | string | Yes | What this field is used for |

## Output
### Success response
| Field | Type | Description |
|-------|------|-------------|
| result | object | The response from the external API |

### Error cases
| Error | Trigger condition |
|-------|------------------|
| "Code node error" | JS preparation step failed |
| "API timeout" | External service did not respond within 30s |

## How to Call
\`\`\`json
{"field": "example_value"}
\`\`\`

## Process Flow
1. **Start** — Entry point, receives the task
2. **<Node title>** — Plain English: what this node does
3. **<API Call title>** — What is called, what is sent, what is returned
4. **Reply Success** — What is returned on success
5. **Final** — Task stored, process complete

**Error path (Code failure):** ...
**Error path (API failure):** ...

## External Dependencies
| Dependency | Type | Variable / URL |
|-----------|------|----------------|
| Weather API | HTTP GET | `{{env_var[@openmeteo-url]}}` |

## Notes
- API call has a 30-second timeout
- Rate limited to 5 concurrent tasks (`max_threads: 5`)
```

### Output 2: Enriched process JSON

Read the original process JSON, fill `description` on every node, write result to `.processes/<name>-enriched.json`.

**What to fill:**
- `description` on every node — one active-voice sentence, specific to this process
- `params[].descr` — if empty, infer from field name and context

**What NOT to change:**
- `id`, `obj_type`, `condition`, `logics`, `semaphors`, `x`, `y`, `extra`, `options`
- `title` — only fill if completely empty

**Description quality bar:**

| Node type | Bad | Good |
|-----------|-----|------|
| Code | "Prepares data" | "Builds request body with city name and units for Open-Meteo geocoder API" |
| API Call | "Makes API call" | "Sends GET request to Open-Meteo geocoder to resolve city name to latitude/longitude" |
| Reply Success | "Returns response" | "Returns coordinates and weather data back to the calling process" |
| Reply Error | "Returns error" | "Returns error reply when the geocoder API is unreachable or returns no results" |
| Final | "Final" | "Stores completed task with weather data — temperature, humidity, wind speed" |

## Extraction rules

**Inputs:** Read `process.params` — `flags: ["required", "input"]` = mandatory, `flags: ["input"]` = optional

**Outputs:** Find all nodes with `api_rpc_reply` in `condition.logics`:
- `throw_exception: false` → success response fields (`res_data` keys)
- `throw_exception: true` → error cases (use node `title` + `exception_reason`)

**Flow:** Walk `scheme.nodes` from `obj_type: 1` (Start), follow `go.to_node_id`, trace the happy path first, then error paths via `err_node_id` references.

**External deps:** API Call `url` field, `{{env_var[@name]}}` references in any field, code node references to external services.

## Knowledge references

- `knowledge/node-types.md` — understand what each logic type does
- `knowledge/process-schema.md` — process and node structure

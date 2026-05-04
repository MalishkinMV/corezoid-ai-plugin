---
name: corezoid-architect
description: >
  Designs multi-process Corezoid systems — chooses the right patterns (connector vs logic vs async fan-out),
  plans process decomposition and composition, reviews architectural decisions, and produces architecture
  documents with process maps. Activate whenever a user asks to design a system, choose between patterns,
  plan a multi-process flow, review an architecture, decompose a business problem into Corezoid processes,
  or asks "how should I structure this in Corezoid?" or "what processes do I need for X?".
  Also activate when the user describes a business capability (payments, onboarding, notifications, etc.)
  and wants to implement it in Corezoid — even if they don't use the word "architecture".
---

# Corezoid Architect


## Bundled References

This Codex plugin includes the upstream Corezoid AI docs repository at `../../assets/source/`. When this skill mentions paths from the original repo, resolve them under that directory unless the user has provided project-local versions. Common examples include `knowledge/`, `docs/`, `templates/`, `playbooks/`, `json-schema/`, `samples/`, `mcp-server/`, and `convctl.sh`.

## What this skill does

Helps design Corezoid process systems before building them. The output is an architecture document
(process map + decision rationale) that the process-builder skill can then implement.

Always produce a concrete architecture — not a list of questions. Make reasonable assumptions,
note them explicitly, and let the user correct them.

---

## Core patterns — choose before designing

### Pattern 1: Connector (wraps one external API)

```
Start → [Code: prepare params] → API Call → [Code/set_param: extract] → Reply Success → Final
                                          ↘ Reply Error → Error Final
              ↘ Reply Error Code → Error Final
```

Use when: a business action maps 1:1 to a single external HTTP call.
Examples: get weather, send SMS, charge payment, look up user in CRM.

### Pattern 2: Logic / Orchestrator (calls multiple Corezoid processes)

```
Start → [Code: validate/prepare] → Call Process A (api_rpc) → Call Process B (api_rpc) → Reply Success → Final
                                                            ↘ Reply Error A → Error Final
              ↘ Reply Error Code → Error Final
```

Use when: a business action requires coordinating multiple sub-processes or connectors.
Examples: register user (validate + create account + send welcome email), process order (validate + charge + notify).

### Pattern 3: Async Fan-out (fire-and-forget to multiple targets)

```
Start → [Code: prepare] → Copy Task to A (api_copy) → Copy Task to B (api_copy) → Final
```

Use when: one event triggers multiple independent downstream processes and you don't need results back.
Examples: event logging, sending notifications to multiple channels, publishing to audit trails.

### Pattern 4: State Machine (tasks wait for external callbacks)

```
Start → Set State "pending" → [Wait for Callback] → Condition → Reply Success → Final
                                                              ↘ Reply Error → Error Final
```

Use when: a step requires waiting for an external party (payment webhook, approval, SMS reply).
The task stays live in Corezoid until the callback arrives.

### Pattern 5: Pipeline (sequential enrichment)

```
Start → Connector A → Connector B → Connector C → Reply / Final
```

Use when: data must pass through multiple enrichment or transformation steps in sequence, each
depending on the previous result. Implement as a Logic process calling connectors in order via `api_rpc`.

---

## Decision framework

Ask these questions about each capability:

1. **Does it call an external HTTP API directly?**
   → Yes: Connector pattern. One connector per external endpoint.
   
2. **Does it coordinate multiple steps or other processes?**
   → Yes: Logic/Orchestrator pattern. The orchestrator calls connectors via `api_rpc`.

3. **Does it need the result back from each step?**
   → Yes: `api_rpc` (synchronous, blocks until reply).
   → No: `api_copy` (fire-and-forget, continues immediately).

4. **Does a step involve waiting for an external event (webhook, approval)?**
   → Yes: State Machine pattern with Waiting for Callback node.

5. **Do multiple independent things need to happen after one trigger?**
   → Yes: Async Fan-out. Each downstream flow is independent.

---

## System decomposition — how to break a business capability into processes

### Rule: one process = one responsibility

Good decomposition:
- `get-weather` — one connector, one API
- `resolve-location` — wraps geocoding
- `get-forecast` — wraps forecast API
- `get-current-weather` — orchestrator that calls resolve-location + get-forecast

Bad decomposition:
- A single 30-node process that geocodes, fetches weather, sends a Slack message, logs to a DB, and notifies a user — all in one flow. This is a "God process" — hard to reuse, hard to debug, impossible to change one part without risking another.

### Decomposition heuristics

| Signal | Action |
|--------|--------|
| "It calls API X and then API Y" | Split into two connectors, one orchestrator |
| "We also need to notify X in another context" | Extract the notification into a standalone connector |
| "This logic appears in multiple flows" | Extract it into a Logic process, call it from others |
| "We need to wait for a webhook before continuing" | Split at the wait point: first process sets state, second process handles callback |
| "We send the same message to 5 systems" | One fan-out process with 5 `api_copy` calls |

### Layered architecture (recommended for complex systems)

```
Layer 3: Business Logic
  └── "Register User" orchestrator
        └── calls Layer 2 processes

Layer 2: Domain Logic  
  ├── "Create Account" (calls connectors + validates)
  ├── "Send Welcome Email" (calls email connector)
  └── "Log Registration Event" (async, calls audit connector)

Layer 1: Connectors (raw API wrappers)
  ├── "Create User in Auth Service" connector
  ├── "Send Email via SendGrid" connector
  └── "Write to Audit DB" connector
```

Benefits: connectors can be reused across domains. Domain logic can be tested independently.
Business logic orchestrates without knowing implementation details.

---

## Architecture document output format

Always produce this when designing a system:

```markdown
# Architecture: <System Name>

## Overview
<1-3 sentences on what this system does and why it exists.>

## Process Map

| Process Name | Pattern | Calls | Called By | Notes |
|-------------|---------|-------|-----------|-------|
| get-weather | Logic | resolve-location, get-forecast | External callers | Entry point |
| resolve-location | Connector | Open-Meteo Geocoding API | get-weather | No auth required |
| get-forecast | Connector | Open-Meteo Forecast API | get-weather | Returns current metrics |

## Data Flow

<Describe how data flows from entry to exit. Focus on what goes in and what comes out of each step.>

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Async vs sync | `api_rpc` throughout | Caller needs weather data before continuing |
| Single connector vs split | Split geocoder + forecast | Two distinct APIs with independent error handling |
| Error isolation | Per-connector error nodes | Distinct error types: geo failure vs weather failure |

## Error Handling Plan

<For each failure mode, describe what the system does:>
- Geocoding fails → return "city not found" error, log to audit
- Forecast timeout → return "weather unavailable", caller falls back to cached value
- etc.

## Environment Variables Needed

| Variable | Purpose |
|----------|---------|
| `api-key-sendgrid` | SendGrid auth |
| `slack-webhook-url` | Slack notification target |

## Assumptions

<List any assumptions made that the user should confirm:>
- Assumed synchronous flow — user's caller can wait for response
- Assumed no caching required — all calls hit live APIs
```

---

## Common architectural mistakes to flag

| Mistake | Why it's a problem | Better approach |
|---------|-------------------|-----------------|
| One 30+ node "God process" | Hard to debug, reuse, or change | Decompose into connector + orchestrator layers |
| `api_copy` when result is needed | Silent failure — caller gets no data | Use `api_rpc` |
| `api_rpc` for fire-and-forget logging | Caller blocks unnecessarily | Use `api_copy` |
| No error escalation on API nodes | Tasks hang on timeout | Always add `err_node_id` + semaphor |
| Hardcoded URLs/tokens in nodes | Impossible to rotate keys, env-specific | Use `{{env_var[@name]}}` |
| State duplication across processes | Divergence and bugs | Extract shared state into a dedicated State process |
| Deep nesting of `api_rpc` chains | Hard to trace, latency compounds | Consider flattening or using async fan-out where reply isn't needed |

---

## After designing — hand off to the right skill

Once the architecture is agreed upon, route to:
- **corezoid-process-builder** — implement each process from the architecture
- **corezoid-process-reviewer** — review generated processes before deployment
- **corezoid-api-wrapper** — expose the top-level orchestrator as an HTTP endpoint
- **corezoid-dashboard-manager** — create monitoring dashboards for the deployed system
- **corezoid-process-tech-writer** — generate documentation after deployment

---

## Reference

- Node & logic types: `knowledge/node-types.md`
- Silent failure patterns: `knowledge/gotchas.md`
- Working example (2-API connector with full error handling): process 1831885, `.processes/get_current_weather-enriched.json`

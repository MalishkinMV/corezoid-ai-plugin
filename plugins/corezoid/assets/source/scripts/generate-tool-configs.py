#!/usr/bin/env python3
"""
Generate AI tool-specific configuration files from the knowledge/ SSOT.

Reads:
  knowledge/auth.md
  knowledge/process-schema.md
  knowledge/node-types.md
  knowledge/gotchas.md
  agents/*.md

Writes:
  .cursor/rules/corezoid.mdc
  .kiro/steering/corezoid.md
  .windsurf/rules/corezoid.md
  .github/copilot-instructions.md
  .devin.md

Run: python scripts/generate-tool-configs.py
"""

from __future__ import annotations

import os
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent


def read(rel: str) -> str:
    return (REPO_ROOT / rel).read_text()


def write(rel: str, content: str) -> None:
    path = REPO_ROOT / rel
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)
    print(f"  wrote {rel}")


def list_agents() -> list[str]:
    agents_dir = REPO_ROOT / "agents"
    return sorted(
        f.stem
        for f in agents_dir.glob("*.md")
        if f.stem != "README"
    )


# ---------------------------------------------------------------------------
# Shared core content (tool-agnostic)
# ---------------------------------------------------------------------------

CORE_INTRO = """\
# Corezoid AI Console

This repository is the knowledge base for AI-assisted Corezoid process development.
Corezoid is a visual process engine where business logic is defined as JSON graphs of connected nodes.

Knowledge base: `knowledge/` directory — auth, process schema, node types, gotchas.
Agent workflows: `agents/` directory — tool-agnostic workflow definitions.
MCP server: `mcp-server/` — 35+ tools for live Corezoid interaction.
"""

CORE_AUTH = """\
## Authentication (path-based SHA-1)

NEVER use HTTP header auth. Use `corezoid_client.py` which handles signing automatically.

```
URL:  {BASE_URL}/{API_LOGIN}/{timestamp}/{SHA1(ts+secret+body+secret)}
```

Credentials from Corezoid UI → Account Settings → API:
- `API_LOGIN` — numeric user ID
- `SECRET` — long static string. NOT the JWT browser token (`eyJhbGci...`).
- `BASE_URL` — must end with `/api/2/json`
- `COMPANY_ID` — `i` prefix + workspace ID (e.g. `i469b4ea8...`)

Copy `.env.example` → `.env` before running anything.
"""

CORE_CLIENT = """\
## Primary tool: corezoid_client.py

```python
from corezoid_client import CorezoidClient
client = CorezoidClient(api_login='65281', secret='...', base_url='...', company_id='i469b...')

result = client.get_process_scheme(process_id)
nodes = result['ops'][0]['scheme'][0]['scheme']['nodes']  # scheme is a LIST — always [0]

r = client.create_task(process_id, {"city": "Kyiv"}, ref="test-1")
task_id = r['ops'][0]['obj_id']
client.show_task(process_id, task_id=task_id)
```

Always batch ops — one request can contain multiple operations in the `ops` array.
"""

CORE_WORKFLOW = """\
## Workflow for any process task

1. Read `knowledge/` files for the relevant context
2. Fetch existing process: `client.get_process_scheme(process_id)`
3. Build or modify process JSON
4. Deploy: 4-step API pattern (create process → create nodes → modify nodes → commit)
5. Test: `client.create_task()` → `client.show_task()` — verify output before done

Agent definitions in `agents/` describe the detailed workflow for each task type.
"""

CORE_CRITICAL = """\
## Critical rules (silent failures)

- `err_node_id` is mandatory on every logic node — missing = unhandled crash
- Every API call node needs a timeout semaphor — missing = tasks hang forever
- `api_rpc` blocks (need result); `api_copy` fires and forgets — confusing them = silent data loss
- Error reply nodes must have `throw_exception: true`
- Final success node must have `options: {"save_task": true}`
- Node modify always sends the FULL payload — partial updates silently reset fields
- Response body is under the second key name in `response` mapping (not under `body`)
- `set_param` not `code` for copy/concat/arithmetic — code nodes spin up a JS runtime
- Dashboard chart type is `column` (not `bar`); grid uses `width`/`height` (not `w`/`h`)

Full list: `knowledge/gotchas.md`
"""

CORE_AGENTS = """\
## Available agents (workflows)

| Agent | File | Purpose |
|-------|------|---------|
| architect | `agents/architect.md` | Design multi-process systems, choose patterns |
| process-builder | `agents/process-builder.md` | Build new connector or logic process |
| process-reviewer | `agents/process-reviewer.md` | Audit process before deployment |
| process-updater | `agents/process-updater.md` | Modify process or debug failed task |
| process-tech-writer | `agents/process-tech-writer.md` | Generate docs and enrich JSON |
| apigw-manager | `agents/apigw-manager.md` | Expose process as HTTP endpoint via API Gateway |
| api-wrapper-manager | `agents/api-wrapper-manager.md` | Alternative HTTP wrapper via convctl |
| dashboard-manager | `agents/dashboard-manager.md` | Create monitoring dashboards |

"""


def build_core() -> str:
    return "\n".join([
        CORE_INTRO,
        CORE_AUTH,
        CORE_CLIENT,
        CORE_WORKFLOW,
        CORE_CRITICAL,
        CORE_AGENTS,
    ])


# ---------------------------------------------------------------------------
# Tool-specific formats
# ---------------------------------------------------------------------------

def generate_cursor() -> None:
    content = f"""---
description: Corezoid AI Console — rules for building, reviewing, and deploying Corezoid processes
alwaysApply: true
---

{build_core()}

## Knowledge files

Always consult before generating process JSON:
- `knowledge/auth.md` — SHA-1 auth, credentials, batching
- `knowledge/process-schema.md` — complete process and node JSON structure
- `knowledge/node-types.md` — all logic types with examples
- `knowledge/gotchas.md` — 15 documented silent failures

<!-- Auto-generated by scripts/generate-tool-configs.py — do not edit manually -->
"""
    write(".cursor/rules/corezoid.mdc", content)


def generate_kiro() -> None:
    content = f"""# Corezoid AI Console

{build_core()}

## Steering files

Load these when working on specific tasks:
- `knowledge/auth.md` — authentication reference
- `knowledge/process-schema.md` — process JSON structure and validation
- `knowledge/node-types.md` — node logic type reference
- `knowledge/gotchas.md` — known silent failures
- `agents/<task>.md` — workflow for the specific task at hand
- `playbooks/` — step-by-step guides for specific scenarios

<!-- Auto-generated by scripts/generate-tool-configs.py — do not edit manually -->
"""
    write(".kiro/steering/corezoid.md", content)


def generate_windsurf() -> None:
    content = f"""# Corezoid AI Console Rules

{build_core()}

## When in doubt

Read `knowledge/gotchas.md` — 15 documented silent failures that cause incorrect behavior without obvious error messages.

<!-- Auto-generated by scripts/generate-tool-configs.py — do not edit manually -->
"""
    write(".windsurf/rules/corezoid.md", content)


def generate_copilot() -> None:
    content = f"""# GitHub Copilot Instructions — Corezoid AI Console

{build_core()}

## Reference files for process generation

- `knowledge/auth.md`
- `knowledge/process-schema.md`
- `knowledge/node-types.md`
- `knowledge/gotchas.md`
- `samples/` — complete working process examples
- `templates/` — reusable process skeletons

<!-- Auto-generated by scripts/generate-tool-configs.py — do not edit manually -->
"""
    write(".github/copilot-instructions.md", content)


def generate_devin() -> None:
    content = f"""# Devin Instructions — Corezoid AI Console

{build_core()}

## Setup

```bash
cp .env.example .env
# Fill in API_LOGIN, SECRET, BASE_URL, COMPANY_ID
pip install python-dotenv  # only external dep for corezoid_client.py
```

## MCP server (optional)

```bash
cd mcp-server && uv sync
uv run python main.py
```

<!-- Auto-generated by scripts/generate-tool-configs.py — do not edit manually -->
"""
    write(".devin.md", content)


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> None:
    print("Generating AI tool configs from knowledge/ SSOT...")

    # Verify knowledge files exist
    required = ["knowledge/auth.md", "knowledge/process-schema.md",
                "knowledge/node-types.md", "knowledge/gotchas.md"]
    missing = [f for f in required if not (REPO_ROOT / f).exists()]
    if missing:
        print(f"ERROR: Missing knowledge files: {missing}", file=sys.stderr)
        sys.exit(1)

    generate_cursor()
    generate_kiro()
    generate_windsurf()
    generate_copilot()
    generate_devin()

    print("Done. All tool configs updated.")
    print("\nNote: These files are auto-generated. Edit knowledge/ or agents/ to change them.")


if __name__ == "__main__":
    main()

#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PLUGIN = ROOT / "plugins" / "corezoid"
CODEX_MARKETPLACE = ROOT / ".agents" / "plugins" / "marketplace.json"
CLAUDE_MARKETPLACE = ROOT / ".claude-plugin" / "marketplace.json"
REPOSITORY = "https://github.com/corezoid/corezoid-ai-plugin"


def load_json(path: Path) -> dict:
    with path.open() as handle:
        data = json.load(handle)
    if not isinstance(data, dict):
        raise SystemExit(f"{path} must contain a JSON object")
    return data


def assert_path(label: str, relative_path: str) -> None:
    if not relative_path.startswith("./"):
        raise SystemExit(f"{label} must start with ./, got {relative_path!r}")
    path = PLUGIN / relative_path[2:]
    if not path.exists():
        raise SystemExit(f"{label} points to missing path: {relative_path}")


def assert_plugin_path(label: str, path: Path, relative_path: str) -> None:
    if not relative_path.startswith("./"):
        raise SystemExit(f"{label} must start with ./, got {relative_path!r}")
    target = path / relative_path[2:]
    if not target.exists():
        raise SystemExit(f"{label} points to missing path: {relative_path}")


def main() -> None:
    marketplace = load_json(CODEX_MARKETPLACE)
    claude_marketplace = load_json(CLAUDE_MARKETPLACE)
    manifest = load_json(PLUGIN / ".codex-plugin" / "plugin.json")
    claude_manifest = load_json(PLUGIN / ".claude-plugin" / "plugin.json")
    load_json(PLUGIN / ".mcp.json")
    load_json(PLUGIN / "assets" / "public" / "source-links.json")
    load_json(PLUGIN / "assets" / "source-metadata.json")

    if marketplace.get("name") != "corezoid":
        raise SystemExit("marketplace name must be corezoid")
    if marketplace.get("interface", {}).get("displayName") != "Corezoid":
        raise SystemExit("marketplace interface.displayName must be Corezoid")

    entries = marketplace.get("plugins")
    if not isinstance(entries, list) or len(entries) != 1:
        raise SystemExit("marketplace must contain exactly one plugin entry")
    entry = entries[0]
    if entry.get("name") != "corezoid":
        raise SystemExit("marketplace plugin name must be corezoid")
    if entry.get("source", {}).get("path") != "./plugins/corezoid":
        raise SystemExit("marketplace source.path must be ./plugins/corezoid")
    policy = entry.get("policy", {})
    if policy.get("installation") not in {"AVAILABLE", "INSTALLED_BY_DEFAULT", "NOT_AVAILABLE"}:
        raise SystemExit("marketplace policy.installation is invalid")
    if policy.get("authentication") not in {"ON_INSTALL", "ON_USE"}:
        raise SystemExit("marketplace policy.authentication is invalid")
    if not entry.get("category"):
        raise SystemExit("marketplace category is required")

    if claude_marketplace.get("name") != "corezoid":
        raise SystemExit("Claude marketplace name must be corezoid")
    claude_entries = claude_marketplace.get("plugins")
    if not isinstance(claude_entries, list) or len(claude_entries) != 1:
        raise SystemExit("Claude marketplace must contain exactly one plugin entry")
    claude_entry = claude_entries[0]
    if claude_entry.get("name") != "corezoid":
        raise SystemExit("Claude marketplace plugin name must be corezoid")
    if claude_entry.get("source") != "./plugins/corezoid":
        raise SystemExit("Claude marketplace source must be ./plugins/corezoid")
    if not claude_entry.get("category"):
        raise SystemExit("Claude marketplace category is required")

    for label, plugin_manifest in (
        ("Codex plugin manifest", manifest),
        ("Claude plugin manifest", claude_manifest),
    ):
        if plugin_manifest.get("name") != "corezoid":
            raise SystemExit(f"{label} name must be corezoid")
        if plugin_manifest.get("repository") != REPOSITORY:
            raise SystemExit(f"{label} repository must be {REPOSITORY}")
        if "[TODO" in json.dumps(plugin_manifest):
            raise SystemExit(f"{label} contains TODO placeholders")

    if manifest.get("version") != claude_manifest.get("version"):
        raise SystemExit("Codex and Claude manifest versions must match")
    if claude_entry.get("version") != manifest.get("version"):
        raise SystemExit("Claude marketplace plugin version must match manifest version")

    for key in ("skills", "mcpServers"):
        value = manifest.get(key)
        if value:
            assert_path(key, value)
        claude_value = claude_manifest.get(key)
        if claude_value:
            assert_path(f"Claude {key}", claude_value)

    assert_plugin_path("Claude marketplace source", ROOT, claude_entry["source"])

    interface = manifest.get("interface", {})
    for key in ("composerIcon", "logo"):
        assert_path(f"interface.{key}", interface[key])
    for screenshot in interface.get("screenshots", []):
        assert_path("interface.screenshots[]", screenshot)

    skill_files = sorted((PLUGIN / "skills").glob("*/SKILL.md"))
    if len(skill_files) != 9:
        raise SystemExit(f"expected 9 skill files, found {len(skill_files)}")

    source_files = sorted(
        str(path.relative_to(PLUGIN / "assets" / "source"))
        for path in (PLUGIN / "assets" / "source").rglob("*")
        if path.is_file()
    )
    index_files = (PLUGIN / "assets" / "source-file-index.txt").read_text().splitlines()
    if source_files != index_files:
        raise SystemExit("assets/source-file-index.txt does not match bundled source files")

    print("Corezoid plugin marketplace validation passed")


if __name__ == "__main__":
    main()

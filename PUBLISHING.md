# Publishing Corezoid For Claude Code And Codex

This document describes how to publish the shared Corezoid plugin through Claude Code and Codex marketplace catalogs.

## 1. Validate Locally

Run these checks from the repository root:

```bash
python3 -m json.tool .agents/plugins/marketplace.json >/dev/null
python3 -m json.tool .claude-plugin/marketplace.json >/dev/null
python3 -m json.tool plugins/corezoid/.codex-plugin/plugin.json >/dev/null
python3 -m json.tool plugins/corezoid/.claude-plugin/plugin.json >/dev/null
python3 -m json.tool plugins/corezoid/.mcp.json >/dev/null
python3 -m json.tool plugins/corezoid/assets/public/source-links.json >/dev/null
```

Verify manifest paths, version sync, shared skills, and bundled source index:

```bash
python3 scripts/validate-plugin.py
```

## 2. Test In Claude Code

```bash
claude plugin marketplace add ./
claude plugin install corezoid@corezoid
```

Verify that the shared Corezoid skills are available as Claude Code plugin commands.

## 3. Test In Codex

```bash
codex plugin marketplace add ./
codex plugin marketplace upgrade corezoid
```

Restart Codex, open Plugin Directory, select **Corezoid**, and install the **Corezoid** plugin.

## 4. Push To GitHub

The public repository is:

```text
https://github.com/corezoid/corezoid-ai-plugin
```

Use this remote:

```bash
git remote set-url origin git@github.com:corezoid/corezoid-ai-plugin.git
git push -u origin main
```

## 5. Create A Release Tag

Use the shared plugin manifest version as the release tag:

```bash
git tag -a v1.1.0 -m "Release Corezoid AI plugin v1.1.0"
git push origin v1.1.0
```

## 6. Install From GitHub

Claude Code:

```bash
claude plugin marketplace add corezoid/corezoid-ai-plugin
claude plugin install corezoid@corezoid
```

Codex stable release:

```bash
codex plugin marketplace add corezoid/corezoid-ai-plugin --ref v1.1.0
codex plugin marketplace upgrade corezoid
```

Codex development tracking:

```bash
codex plugin marketplace add corezoid/corezoid-ai-plugin --ref main
codex plugin marketplace upgrade corezoid
```

## 7. Update Policy

When changing the plugin:

1. Update both manifest versions: `plugins/corezoid/.codex-plugin/plugin.json` and `plugins/corezoid/.claude-plugin/plugin.json`.
2. Update `.claude-plugin/marketplace.json` plugin version.
3. Update `CHANGELOG.md`.
4. Re-run `python3 scripts/validate-plugin.py`.
5. Test local install in Claude Code and Codex when those CLIs are available.
6. Commit, push, tag a new release, and publish release notes.
7. Ask users to upgrade their local marketplace/plugin.

## Official Directories

Until each platform provides a fully self-serve official plugin directory submission flow, this public Git-backed marketplace repository is the canonical distribution path for Corezoid.

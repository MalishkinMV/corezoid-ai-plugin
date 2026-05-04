# Publishing Corezoid For Codex

This document describes how to publish the Corezoid plugin through a Codex marketplace repository.

## 1. Validate Locally

Run these checks from the repository root:

```bash
python3 -m json.tool .agents/plugins/marketplace.json >/dev/null
python3 -m json.tool plugins/corezoid/.codex-plugin/plugin.json >/dev/null
python3 -m json.tool plugins/corezoid/.mcp.json >/dev/null
python3 -m json.tool plugins/corezoid/assets/public/source-links.json >/dev/null
```

Verify manifest paths:

```bash
python3 scripts/validate-plugin.py
```

## 2. Test In Codex

```bash
codex plugin marketplace add ./
codex plugin marketplace upgrade corezoid
```

Restart Codex, open Plugin Directory, select **Corezoid**, and install the **Corezoid** plugin.

## 3. Push To GitHub

Create a GitHub repository under the Corezoid organization, for example:

```text
https://github.com/corezoid/corezoid-codex-plugin
```

Push this repository:

```bash
git remote add origin git@github.com:corezoid/corezoid-codex-plugin.git
git push -u origin main
```

## 4. Create A Release Tag

Use the plugin manifest version as the release tag:

```bash
git tag v1.0.2
git push origin v1.0.2
```

## 5. Install From GitHub

For development tracking:

```bash
codex plugin marketplace add corezoid/corezoid-codex-plugin --ref main
```

For stable installs:

```bash
codex plugin marketplace add corezoid/corezoid-codex-plugin --ref v1.0.2
```

Then refresh:

```bash
codex plugin marketplace upgrade corezoid
```

## 6. Update Policy

When changing the plugin:

1. Update `plugins/corezoid/.codex-plugin/plugin.json` version.
2. Update `CHANGELOG.md`.
3. Re-run validation.
4. Test install locally.
5. Commit and tag a new release.
6. Ask users to run `codex plugin marketplace upgrade corezoid`.

## Official Plugin Directory

As of May 4, 2026, OpenAI documentation says adding plugins to the official Plugin Directory and self-serve plugin publishing are coming soon. Until then, Git-backed marketplace distribution is the supported path for sharing this plugin.

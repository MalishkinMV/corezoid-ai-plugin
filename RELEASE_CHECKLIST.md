# Corezoid AI Plugin Release Checklist

Use this before tagging a public release.

- [ ] `.agents/plugins/marketplace.json` has `name: corezoid` and `interface.displayName: Corezoid`.
- [ ] `.claude-plugin/marketplace.json` has `name: corezoid` and points to `./plugins/corezoid`.
- [ ] Codex marketplace entry includes `policy.installation`, `policy.authentication`, and `category`.
- [ ] Codex manifest version matches Claude manifest version.
- [ ] `.claude-plugin/marketplace.json` plugin version matches both manifests.
- [ ] Manifests have no TODO placeholders.
- [ ] Manifest asset paths resolve and files are under `plugins/corezoid/assets/`.
- [ ] `plugins/corezoid/.mcp.json` contains no credentials.
- [ ] `plugins/corezoid/assets/source-metadata.json` records the bundled upstream commit.
- [ ] `plugins/corezoid/assets/source-file-index.txt` matches bundled source files.
- [ ] `processes/` and other local test artifacts are not tracked.
- [ ] `python3 scripts/validate-plugin.py` passes.
- [ ] Claude Code can install the plugin from the local marketplace when the CLI is available.
- [ ] Codex can install the plugin from the local marketplace when the CLI is available.
- [ ] GitHub README install commands use `corezoid/corezoid-ai-plugin`.
- [ ] Release tag matches plugin version, for example `v1.1.0`.

# Corezoid Plugin Release Checklist

Use this before tagging a public release.

- [ ] `.agents/plugins/marketplace.json` has `name: corezoid` and `interface.displayName: Corezoid`.
- [ ] Marketplace entry points to `./plugins/corezoid`.
- [ ] Marketplace entry includes `policy.installation`, `policy.authentication`, and `category`.
- [ ] `plugins/corezoid/.codex-plugin/plugin.json` has the intended semantic version.
- [ ] Manifest has no TODO placeholders.
- [ ] Manifest asset paths resolve and files are under `plugins/corezoid/assets/`.
- [ ] `plugins/corezoid/.mcp.json` contains no credentials.
- [ ] `plugins/corezoid/assets/source-metadata.json` records the bundled upstream commit.
- [ ] `plugins/corezoid/assets/source-file-index.txt` matches bundled source files.
- [ ] `processes/` and other local test artifacts are not tracked.
- [ ] `python3 scripts/validate-plugin.py` passes.
- [ ] Codex can install the plugin from the local marketplace.
- [ ] GitHub README install command uses the actual repository name.
- [ ] Release tag matches plugin version, for example `v1.0.2`.

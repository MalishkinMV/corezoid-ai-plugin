# Security

The Corezoid AI plugin intentionally ships with an empty `plugins/corezoid/.mcp.json` so API credentials are never bundled in the marketplace package.

Do not commit Corezoid credentials, `.env` files, API secrets, workspace IDs tied to private environments, exported customer process data, or generated process exports containing private business logic.

The bundled `plugins/corezoid/assets/source/` corpus is copied from public Corezoid documentation and examples. If you discover a secret or private artifact in this repository, open a private security issue or contact Corezoid support before creating a public report.

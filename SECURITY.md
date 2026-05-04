# Security

The Corezoid Codex plugin intentionally ships with an empty `.mcp.json` so API credentials are never bundled in the marketplace package.

## Credentials

Do not commit Corezoid credentials, `.env` files, API secrets, workspace IDs tied to private environments, or exported customer process data.

Local experiments and process exports should stay under ignored paths such as `processes/` or `.processes/`.

## Reporting Issues

For issues related to this plugin package, open an issue in the GitHub repository that hosts this marketplace.

For Corezoid product or account security issues, use Corezoid's official support channels from https://corezoid.com.

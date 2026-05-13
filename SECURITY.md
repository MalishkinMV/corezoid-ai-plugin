# Security

## Credentials

`plugins/corezoid/.mcp.json` ships without API credentials so that tokens are never bundled in the marketplace package. Credentials are stored locally in `.env` in the user's working directory after OAuth2 login and are never committed to this repository.

Do not commit:

- Corezoid API tokens or OAuth2 refresh tokens
- `.env` files of any kind
- Workspace IDs or stage IDs tied to private environments
- Exported process files (`.conv.json`) containing private business logic or customer data
- Any other secrets or credentials

## Reporting a Vulnerability

If you discover a security issue — including a secret accidentally committed to this repository — please **do not open a public GitHub issue**. Instead:

1. Open a private security advisory at `https://github.com/corezoid/corezoid-ai-plugin/security/advisories/new`.
2. Or email `support@corezoid.com` with the subject line `[SECURITY] corezoid-ai-plugin`.

We will respond within 5 business days and coordinate a fix before any public disclosure.

# Corezoid API Authentication

## Method: Path-based SHA-1 signing

Every API call uses this URL pattern — no special HTTP headers needed:

```
POST {BASE_URL}/{API_LOGIN}/{timestamp}/{signature}
Content-Type: application/json
Body: compact JSON payload
```

Signature formula:
```
signature = SHA1( timestamp + SECRET + compact_json_body + SECRET )
```

## Required credentials

| Variable | Description | Where to get |
|----------|-------------|--------------|
| `API_LOGIN` | Numeric user ID | Corezoid UI → Account Settings → API |
| `SECRET` | API secret key (long static string) | Corezoid UI → Account Settings → API |
| `BASE_URL` | API endpoint — must end with `/api/2/json` | e.g. `https://admin.corezoid.com/api/2/json` |
| `COMPANY_ID` | Workspace ID with `i` prefix | e.g. `i469b4ea8...` (from workspace URL) |

## Critical: What the SECRET is NOT

The JWT token visible in your browser (`eyJhbGci...`) is a **session cookie**. It expires and cannot be used to sign API calls. The SECRET is a separate static key in Account Settings → API.

## Python implementation (use this, don't write your own)

`corezoid_client.py` in the repo root handles all signing:

```python
from corezoid_client import CorezoidClient

client = CorezoidClient(
    api_login='65281',
    secret='your_static_secret_here',
    base_url='https://admin.corezoid.com/api/2/json',
    company_id='i469b4ea8...'
)
```

## company_id format

Always `i` prefix + numeric workspace ID. Example: `i469b4ea8f1894c46b613d7be004aeac1`.

Find it in the Corezoid workspace URL: `admin.corezoid.com/<workspace-id>/workspace/...`

## Batching: always send multiple ops in one request

The API accepts any number of operations in a single `ops` array. Never make two separate calls when one batched call works — it wastes round-trips and token budget.

```python
# ✅ One request, two operations
client.make_request({'ops': [
    {'type': 'show', 'obj': 'conv', 'obj_id': 123, 'company_id': COMPANY_ID},
    {'type': 'show', 'obj': 'conv', 'obj_id': 456, 'company_id': COMPANY_ID}
]})

# ❌ Two separate requests for data that fits in one
client.make_request({'ops': [{'type': 'show', 'obj': 'conv', 'obj_id': 123, ...}]})
client.make_request({'ops': [{'type': 'show', 'obj': 'conv', 'obj_id': 456, ...}]})
```

## Response envelope

All API responses follow this structure:

```json
{
  "request_proc": "ok",
  "ops": [
    {
      "proc": "ok",
      "obj": "conv",
      "obj_id": 12345
    }
  ]
}
```

Always check `result['request_proc'] == 'ok'` and `result['ops'][N]['proc'] == 'ok'` before using the response.

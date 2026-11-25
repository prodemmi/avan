# API reliability and message status

## Idempotent sends

`POST /api/send` and `POST /api/quick-send` require an `Idempotency-Key` header
(1–255 bytes). Keys are unique within the API key's workspace. Reuse a key
only for a retry of the same logical request; reusing it with a different
provider, recipient, content, schedule, retry option, template, or payload
returns `409 Conflict`.

The first request creates one message row. Replays do not enqueue or invoke the
provider again. Responses identify the existing `message_id`, current `status`,
and `duplicate` flag: pending/scheduled requests return `202`, sent requests
return `200`, and failed requests return `502`. A replay of a failed request
does not retry delivery; submit a new logical request with a new key if a new
send is intended. Scheduled sends retain the same key and message through
dispatch.

## Workspace API security policy

In Dashboard → Settings, set the `api-security-policy` value to a JSON object:

```json
{"ip_allowlist":["192.0.2.0/24"],"api_key_allowlist":[12,15]}
```

Both fields are optional. IP entries accept IP addresses or CIDR prefixes;
`api_key_allowlist` contains IDs from the workspace's `api_keys` table. An
empty or absent policy preserves existing access. Invalid JSON, unknown fields,
invalid networks, or non-positive key IDs are rejected when configured and
fail closed on API requests.

The current server has no configured trusted-proxy list and explicitly ignores
forwarded client-IP headers for this policy. Allowlisting therefore uses the
direct network peer address. Deployments behind a proxy must establish an
explicit trusted-proxy policy before using forwarded client IPs.

## Message inquiry and provider callback limitation

`GET /api/messages/{id}` returns status only when the message belongs to the
API key's workspace. Current states are `pending`, `scheduled`, `sent`, and
`failed`; `sent` means the configured notifier accepted the send operation, not
that the recipient received or read it.

The current SMS, email, Telegram, Slack, and webhook integrations do not expose
an inbound delivery/read callback route or persist provider delivery receipts.
Consequently `delivered` and `read` are intentionally not reported. Those
states can be added when a provider callback with verifiable message
correlation is integrated.

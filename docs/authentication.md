# Authentication

Webull OpenAPI requests use two independent mechanisms: a per-request signature
and an access token. The SDK handles both. `client.New` builds the signing
headers and `Client.Do` signs every outgoing request; `Client.EnsureToken`
creates and activates an access token and installs it as the `x-access-token`
header on later requests.

## Request signing

Every request is signed with **HMAC-SHA1** over a percent-encoded canonical
string. The HMAC key is the app secret followed by `&`, and the signature is
returned as standard Base64.

The canonical string is assembled as follows:

1. Build a list of `name=value` entries from the query parameters and the
   participating headers (`x-app-key`, `x-signature-algorithm`,
   `x-signature-version`, `x-signature-nonce`, `x-timestamp`, `host`).
2. Sort the entries by name. For a repeated name, sort its values and join them
   with `&` inside the single entry.
3. Join the entries with `&` and append them to the request path with `&`:
   `path&entry1&entry2...`.
4. When there is a body, append `&` and the **uppercase hexadecimal MD5** of the
   exact bytes transmitted. A bodyless request has no digest.
5. Percent-encode the whole string using the RFC 3986 unreserved set (`A-Z`,
   `a-z`, `0-9`, `-`, `_`, `.`, `~`), with uppercase hex digits. This matches
   Python's `urllib.parse.quote(s, safe="")` and encodes `&`, `=`, `/`, `:`,
   `?`, and space.

The HTTP method is not part of the canonical string. `host` is the request host
as `hostname[:port]`, including the port only when it is not the scheme default.

Two normalization rules matter in practice:

- Query values are sorted so the bytes signed are the bytes sent.
- The body hash is computed over compact JSON with HTML escaping disabled (Go's
  `json.Marshal` would rewrite `<`, `>`, and `&`, changing the digest). The SDK
  signs the exact transmitted bytes.

A golden signature test pins this behavior so it cannot drift silently.

## Headers

Signed requests carry the following headers:

| Header | Description |
|--------|-------------|
| `x-app-key` | Your Webull OpenAPI app key |
| `x-timestamp` | Current UTC time in ISO-8601 form, `2006-01-02T15:04:05Z` |
| `x-signature-version` | Signature version, `1.0` |
| `x-signature-algorithm` | Signature algorithm, `HMAC-SHA1` |
| `x-signature-nonce` | Random 128-bit nonce, 32 lowercase hex characters |
| `host` | Request host; participates in the signature |
| `x-version` | API version, `v2` by default (or `v3` where configured) |
| `x-signature` | Base64 HMAC-SHA1 signature; not itself signed |

Authenticated endpoints additionally require:

| Header | Description |
|--------|-------------|
| `x-access-token` | Access token obtained from the token endpoint |

The app secret is **not** a header. It is used only to compute `x-signature` and
must never be transmitted.

## Token lifecycle

Tokens are created with `POST /openapi/auth/token/create` and checked with
`POST /openapi/auth/token/check`. Tokens are time-sensitive (15 days by
default). The public API is:

| Method | Behavior |
|--------|----------|
| `Client.CreateToken(ctx)` | Creates a token; does not touch the cache |
| `Client.CheckToken(ctx, token)` | Returns the current status of a token |
| `Client.EnsureToken(ctx)` | Reuses a cached valid token or creates one, waiting for `NORMAL` when needed, then caches it |
| `Client.CurrentToken()` | Returns the cached token, or `nil` |
| `Client.AccessToken()` | Returns the cached token value, or `""` |
| `Client.SetToken(token)` | Sets or clears the cached token |

| Status | Meaning |
|--------|---------|
| `PENDING` | Created, awaiting verification (Webull App code / 2FA) |
| `NORMAL` | Valid and usable |
| `INVALID` | Revoked, never used, or unused for 15 consecutive days |
| `EXPIRED` | Verification was not completed within five minutes; create a new token |

Sandbox tokens are issued as `NORMAL` automatically, with no 2FA step.
Production tokens start `PENDING` and become `NORMAL` after verification.
`EnsureToken` polls `CheckToken` until the token becomes `NORMAL`, becomes
terminal, or the poll timeout elapses. The default interval is 5 seconds and the
default timeout is 5 minutes, changeable with `Client.SetTokenPollInterval` and
`Client.SetTokenPollTimeout`.

Once a token is cached, `Client.Do` attaches it as the `x-access-token` header
after signing. The header is deliberately omitted from the token endpoints
themselves.

For sandbox automation, `client.WithAutoToken(true)` makes the first non-token
request obtain a token automatically when the environment is the sandbox. In
production the SDK never starts the 2FA flow implicitly; it returns
`client.ErrAccessTokenRequired` so you can drive `EnsureToken` explicitly.

## Sandbox

Point the client at the sandbox host while developing:

```
https://api.sandbox.webull.hk
```

```go
cl, err := client.New(
	client.WithCredentials(key, secret),
	client.WithSandbox(),
)
```

Never commit app keys, app secrets, or access tokens. Supply them through
environment variables (`WEBULL_APP_KEY`, `WEBULL_APP_SECRET`) or your own secret
manager. The token endpoint allows 10 requests per 30 seconds.

# Authentication

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Webull uses a dual layer: an HMAC-SHA1 request signature plus an access token. Server-to-server endpoints sign every request; Display Solution (Client-to-Server) uses an OAuth-style client token.

[<- Webull API Reference](../webull-api.md)

## Create Token

`POST /auth/tokens/create`

> Creates an access token. This interface is used to generate a new Token.

| | |
|---|---|
| **SDK** | `client.CreateToken` |
| **Reference** | [create-token.md](https://developer.webull.hk/apis/docs/reference/create-token.md) |
| **Note** | Returns a token in `PENDING`; completes 2FA in the Webull App. Sandbox tokens are `NORMAL` immediately. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `token` | string | yes | Access token string, used for authentication of subsequent API calls. Token is a 32-digit hexadecimal string that is unique and time-sensitive. |
| `expires_at` | integer | yes | Token expiration timestamp, a Unix timestamp in milliseconds. After this time, the token will become invalid and need to be recreated. |
| `status` | string | yes | Token validity status code, indicating the result of the token operation. PENDING indicates pending verification, NORMAL indicates valid, INVALID indicates the token is invalid, and EXPIRED indicates it has expired. — one of: `PENDING`, `NORMAL`, `INVALID`, `EXPIRED` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Check Token

`POST /auth/tokens/check`

> Retrieves Token status.

| | |
|---|---|
| **SDK** | `client.CheckToken` |
| **Reference** | [check-token.md](https://developer.webull.hk/apis/docs/reference/check-token.md) |
| **Note** | Polled by `client.EnsureToken` until the token is `NORMAL`. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `token` | string | yes | Access token, used for identity authentication and permission verification. This field is a unique identifier for token checking, refreshing, etc. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `token` | string | yes | Access token string, used for authentication of subsequent API calls. Token is a 32-digit hexadecimal string that is unique and time-sensitive. |
| `expires_at` | integer | yes | Token expiration timestamp, a Unix timestamp in milliseconds. After this time, the token will become invalid and need to be recreated. |
| `status` | string | yes | Token validity status code, indicating the result of the token operation. PENDING indicates pending verification, NORMAL indicates valid, INVALID indicates the token is invalid, and EXPIRED indicates it has expired. — one of: `PENDING`, `NORMAL`, `INVALID`, `EXPIRED` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Client Token (Display)

`POST /auth/client-tokens/create`

> Generates an access token and refresh token pair for end-user applications to directly access Webull services. This endpoint must be called from your institution's backend service with valid AK/SK credentials. The access token authorizes API calls, while the refresh token enables token renewal. Creating a new token pair for an existing customer immediately invalidates any previously issued tokens for that customer. • Access token expires in approximately 2 hours (may vary slightly) • Refresh token expires in 15 days • **Note**: This access_token is used to access Market Data API

| | |
|---|---|
| **SDK** | `display.Service.EnsureToken` |
| **Reference** | [create-client-token.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/create-client-token.md) |
| **Note** | Internal to `display.Service`; `client_user_id` is fixed to `openapi_client`. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_user_id` | string | yes | The unique identifier for the customer in your system. This identifier is used to associate the token pair with the specific end-user. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `access_token` | string | yes | Short-lived token for API authorization. |
| `expires_at` | integer | yes | Access Token Expiration Time in milliseconds since epoch. |
| `refresh_token` | string | yes | Long-lived token for obtaining new access tokens. |
| `refresh_expires_at` | integer | yes | Refresh Token Expiration Time in milliseconds since epoch. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Refresh Client Token (Display)

`POST /auth/client-tokens/refresh`

> Obtains a new access token and refresh token pair using a valid refresh token. This endpoint requires AK/SK signature authentication from your institution's backend service. Both the refresh token and request signature are validated before issuing new credentials. The old tokens will be invalidated after successful refresh. • Access token expires in approximately 2 hours (may vary slightly) • Refresh token expires in 15 days • **Note**: This access_token is used to access Market Data API

| | |
|---|---|
| **SDK** | `display.Service.RefreshClientToken` |
| **Reference** | [refresh-client-token.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/refresh-client-token.md) |
| **Note** | The SDK re-creates the client token instead of refreshing it. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `refresh_token` | string | yes | The refresh token obtained from the token creation or previous refresh operation. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `access_token` | string | yes | Short-lived token for API authorization. |
| `expires_at` | integer | yes | Access Token Expiration Time in milliseconds since epoch. |
| `refresh_token` | string | yes | Long-lived token for obtaining new access tokens. |
| `refresh_expires_at` | integer | yes | Refresh Token Expiration Time in milliseconds since epoch. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).


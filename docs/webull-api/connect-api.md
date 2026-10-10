# Connect API (OAuth)

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

OAuth 2.0 authorization-code flow for third-party applications (US only).

[<- Webull API Reference](../webull-api.md)

## Authorization Code

`GET /oauth2/auth-codes/get`

> This is the first step of the OAuth2 process. An authorization code is created when the user authorizes your application to access their account. If the user grants permission to your application, the callback URL registered in your application will be invoked. The interface for obtaining the authorization code is completed in the browser. <b>'SEND API REQUEST' function for this endpoint does not work in UAT environment</b>.

| | |
|---|---|
| **SDK** | `connect.AuthorizationURL` |
| **Reference** | [get-authorization-code.md](https://developer.webull.com/apis/docs/reference/connect-api/get-authorization-code.md) |
| **Note** | Browser redirect URL builder. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `response_type` | query | String | yes | Must be code to request an authorization code. |
| `client_id` | query | String | yes | Webull provides the client id |
| `scope` | query | String | yes | The application requests access to the list of scopes. user：user trade：trade wr：write read. |
| `state` | query | String | yes | An unguessable random string, used to protect against request forgery attacks. |
| `redirect_uri` | query | String | yes | The URL to which the user will be redirected after authorization. It must match the redirect URIs in whitelist. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Token

`POST /oauth2/tokens/create`

> This is the second step of the OAuth process. An access token is created using the authorization code from the first step's response. The access token is a key used for API access. These tokens should be protected like passwords.

| | |
|---|---|
| **SDK** | `connect.CreateToken` |
| **Reference** | [create-and-refresh-token.md](https://developer.webull.com/apis/docs/reference/connect-api/create-and-refresh-token.md) |
| **Note** | Takes `TokenRequest` with fields: `GrantType` (`"authorization_code"` or `"refresh_token"`), `Code`, `RedirectURI`, `RefreshToken`. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `access_token` | string | yes | Access token |
| `token_type` | string | yes | Access token type. Currently, only 'Bearer' is supported |
| `expires_in` | string | yes | Access token expiration time. Unit: seconds. |
| `rt_expires_in` | string | yes | Refresh token expiration time. Unit: seconds. |
| `refresh_token` | string | yes | Refresh token |
| `created_at` | string | yes | Token creation time. |
| `identity_id` | string | yes | User unique identifier. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).


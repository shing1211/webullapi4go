# Market Data — News

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

News summaries. The Non-Display endpoint streams Server-Sent Events.

[<- Webull API Reference](../webull-api.md)

## News Summary

`POST /market-data/news/summaries/get`

> Retrieves news summary with SSE stream response.

| | |
|---|---|
| **SDK** | `data.GetNewsSummary` |
| **Reference** | [news-summary.md](https://developer.webull.hk/apis/docs/reference/news-summary.md) |
| **Note** | SSE stream via `client.DoStream`; HK sandbox upstream returns 504. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `category_symbols` | array<object> | yes | List of security symbols by category. |
| `lang` | string |  | Support language, enum: [en]. |

*Nested — `category_symbols`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).


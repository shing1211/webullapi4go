# Streaming (MQTT)

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Real-time market data over MQTT. Subscribe/unsubscribe are HTTP calls that register the session; pushes arrive over MQTT. At most 5 concurrent connections per App Key.

[<- Webull API Reference](../webull-api.md)

## Subscribe

`POST /market-data/streaming/subscribe`

> Subscribes to real-time market data streaming.

| | |
|---|---|
| **SDK** | `stream.Subscribe` |
| **Reference** | [subscribe.md](https://developer.webull.hk/apis/docs/reference/subscribe.md) |
| **Note** | HTTP; registers the MQTT session. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `session_id` | string | yes | The session_id used to create the connection, and the connection must be successfully established. |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |
| `category` | string | yes | Security type, enum, refer to: Category — one of: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK` |
| `sub_types` | array<string> | yes | Subscription data type(s), multiple types separated by commas ",", enum, refer to: SubType, e.g.: ["SNAPSHOT"] |
| `grab` | string | yes | Whether to grab snapshot data, true/false |
| `depth` | string |  | LV2 subscription depth, default 10 levels, US stocks max 50 levels. |
| `overnight_required` | boolean |  | Whether to include overnight session, true/false. For US stock subscriptions, includes overnight session, only effective for US stocks, default is not included. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Unsubscribe

`POST /market-data/streaming/unsubscribe`

> Unsubscribes from real-time market data streaming.

| | |
|---|---|
| **SDK** | `stream.Unsubscribe` |
| **Reference** | [unsubscribe.md](https://developer.webull.hk/apis/docs/reference/unsubscribe.md) |
| **Note** | HTTP; `UnsubscribeAll` clears everything. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `session_id` | string |  | The session_id used to create the connection, and the connection must be successfully established. |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |
| `category` | string | yes | Security type, enum, refer to: Category — one of: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK` |
| `sub_types` | array<string> | yes | Subscription data type(s), multiple types separated by commas ",", enum, refer to: SubType, e.g.: ["SNAPSHOT"] |
| `unsubscribe_all` | boolean |  | Whether to unsubscribe all, true/false. When set to true, all subscriptions will be cancelled. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).


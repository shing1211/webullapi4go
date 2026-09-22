# Market Data — Watchlist

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Watchlist CRUD and instrument membership.

[<- Webull API Reference](../webull-api.md)

## Get Watchlists

`GET /market-data/watchlists/list`

> Retrieves user watchlists.

| | |
|---|---|
| **SDK** | `data.GetWatchlists` |
| **Reference** | [get-watchlist.md](https://developer.webull.hk/apis/docs/reference/get-watchlist.md) |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string |  | Watchlist unique identifier |
| `name` | string |  | Watchlist name |
| `sort` | integer |  | Sort order number for display ordering |
| `create_time` | string |  | Watchlist creation time |
| `update_time` | string |  | Watchlist last update time |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Watchlist

`POST /market-data/watchlists/create`

> Creates a new watchlist.

| | |
|---|---|
| **SDK** | `data.CreateWatchlist` |
| **Reference** | [create-watchlist.md](https://developer.webull.hk/apis/docs/reference/create-watchlist.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Watchlist name |
| `sort` | integer |  | Sort order number (optional) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string |  | Newly created watchlist unique identifier |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Update Watchlist

`POST /market-data/watchlists/update`

> Updates watchlist information.

| | |
|---|---|
| **SDK** | `data.UpdateWatchlist` |
| **Reference** | [update-watchlist.md](https://developer.webull.hk/apis/docs/reference/update-watchlist.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string | yes | Watchlist unique identifier |
| `name` | string |  | New watchlist name (optional) |
| `sort` | integer |  | New sort order number (optional) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `success` | boolean |  | Whether the operation was successful |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Delete Watchlist

`POST /market-data/watchlists/delete`

> Deletes a watchlist.

| | |
|---|---|
| **SDK** | `data.DeleteWatchlist` |
| **Reference** | [delete-watchlist.md](https://developer.webull.hk/apis/docs/reference/delete-watchlist.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string | yes | Watchlist unique identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `success` | boolean |  | Whether the operation was successful |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Get Watchlist Instruments

`GET /market-data/watchlists/instruments/list`

> Retrieves instruments in a watchlist.

| | |
|---|---|
| **SDK** | `data.GetWatchlistInstruments` |
| **Reference** | [get-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/get-watchlist-instruments.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `watchlist_id` | query | string | yes | Watchlist unique identifier. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string |  | Watchlist unique identifier |
| `instruments` | array<object> |  | List of instruments in the watchlist |

*Nested — `instruments`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Unique identifier of the instrument (stock ID, event contract market ID, etc. unified as instrument_id) |
| `symbol` | string |  | Trading symbol of the instrument |
| `name` | string |  | Full name of the instrument |
| `exchange_code` | string |  | Standardized exchange code |
| `sort` | integer |  | Sort order within the watchlist |
| `added_time` | string |  | Time when the instrument was added to the watchlist (ISO 8601) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Add Watchlist Instruments

`POST /market-data/watchlists/instruments/add`

> Adds instruments to a watchlist.

| | |
|---|---|
| **SDK** | `data.AddWatchlistInstruments` |
| **Reference** | [add-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/add-watchlist-instruments.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string | yes | Watchlist unique identifier |
| `instruments` | array<object> | yes | List of instruments to add |

*Nested — `instruments`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `category` | string | yes | Security category — one of: `US_STOCK`, `HK_STOCK` |
| `sort` | integer |  | Sort order within the watchlist (optional, used for update) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `success` | boolean |  | Whether the operation was successful |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Remove Watchlist Instruments

`POST /market-data/watchlists/instruments/remove`

> Removes instruments from a watchlist.

| | |
|---|---|
| **SDK** | `data.RemoveWatchlistInstruments` |
| **Reference** | [remove-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/remove-watchlist-instruments.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string | yes | Watchlist unique identifier |
| `instruments` | array<object> | yes | List of instruments to remove |

*Nested — `instruments`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `category` | string | yes | Security category — one of: `US_STOCK`, `HK_STOCK` |
| `sort` | integer |  | Sort order within the watchlist (optional, used for update) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `success` | boolean |  | Whether the operation was successful |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Update Watchlist Instruments

`POST /market-data/watchlists/instruments/update`

> Updates instruments sort order in a watchlist.

| | |
|---|---|
| **SDK** | `data.UpdateWatchlistInstruments` |
| **Reference** | [update-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/update-watchlist-instruments.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `watchlist_id` | string | yes | Watchlist unique identifier |
| `instruments` | array<object> | yes | List of instruments to update sort order |

*Nested — `instruments`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `category` | string | yes | Security category — one of: `US_STOCK`, `HK_STOCK` |
| `sort` | integer |  | Sort order within the watchlist (optional, used for update) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `success` | boolean |  | Whether the operation was successful |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).


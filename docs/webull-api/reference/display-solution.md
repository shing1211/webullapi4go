# Display Solution — Verbatim Reference

> Hosted Display Solution: a separate entitlement and host with Client-to-Server (Bearer) authentication. The SDK routes these through `display.Service`.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Stock Top Gainers/Losers

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/top-gainers-using-get-new.md>

### List Top Gainers/Losers

Retrieves a ranked list of top gaining or losing stocks for a specified time period. To get top gainers, pass direction=DESC; for top losers, pass direction=ASC. The rank_type parameter controls the time window (e.g., D1=today, W52=52-week)

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/screeners/gainers-losers/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "Retrieves a ranked list of top gaining or losing stocks for a specified time period. To get top gainers, pass direction=DESC; for top losers, pass direction=ASC. The rank_type parameter controls the time window (e.g., D1=today, W52=52-week)",
  "operationId": "topGainersUsingGETNew",
  "parameters": [
    {
      "name": "rank_type",
      "in": "query",
      "description": "Ranking type",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "PRE_MARKET",
          "AFTER_MARKET",
          "M3",
          "M5",
          "D1",
          "D5",
          "MO1",
          "MO3",
          "W52"
        ]
      },
      "example": "D1"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security category",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "sort_by",
      "in": "query",
      "description": "Sort Field",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "CHANGE_RATIO",
          "RELATIVE_VOLUME_10D",
          "MARKET_VALUE",
          "CLOSE",
          "PRICE",
          "PE_TTM",
          "HIGH",
          "LOW",
          "AMPLITUDE",
          "TURNOVER",
          "VOLUME"
        ]
      },
      "example": "CHANGE_RATIO"
    },
    {
      "name": "direction",
      "in": "query",
      "description": "Sorting direction. Ascending order: ASC, descending order: DESC",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "DESC"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key returned from previous page response. Pass null or omit for first page.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "data": {
                "type": "array",
                "description": "Data list",
                "items": {
                  "type": "object",
                  "properties": {
                    "symbol": {
                      "type": "string",
                      "description": "Security symbol",
                      "example": "TSLA"
                    },
                    "name": {
                      "type": "string",
                      "description": "Security name",
                      "example": "Tesla Inc"
                    },
                    "exchange_code": {
                      "type": "string",
                      "description": "Exchange code",
                      "example": "NSQ"
                    },
                    "currency": {
                      "type": "string",
                      "description": "Currency code",
                      "example": "USD"
                    },
                    "pre_close": {
                      "type": "string",
                      "description": "The closing price of the previous trading day",
                      "example": "11.01"
                    },
                    "open": {
                      "type": "string",
                      "description": "Open price for the current trading day",
                      "example": "11.01"
                    },
                    "high": {
                      "type": "string",
                      "description": "Today’s high",
                      "example": "12.11"
                    },
                    "low": {
                      "type": "string",
                      "description": "Today’s low",
                      "example": "11.01"
                    },
                    "close": {
                      "type": "string",
                      "description": "Latest intraday prices for the current trading day",
                      "example": "12.01"
                    },
                    "price": {
                      "type": "string",
                      "description": "The latest price for the current trading day",
                      "example": "11.11"
                    },
                    "change": {
                      "type": "string",
                      "description": "Trade change for the current trading day",
                      "example": "1.01"
                    },
                    "change_ratio": {
                      "type": "string",
                      "description": "Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%)",
                      "example": "0.0111"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Trade volume",
                      "format": "0",
                      "example": "1111111"
                    },
                    "turnover": {
                      "type": "string",
                      "description": "Transaction amount, US market stocks and ETFs do not return this data under Nb authorization",
                      "format": "0.00##",
                      "example": "1234567890"
                    },
                    "turnover_rate": {
                      "type": "string",
                      "description": "Turnover Rate",
                      "example": "0.1111"
                    },
                    "market_value": {
                      "type": "string",
                      "description": "Market Value",
                      "format": "0.00##",
                      "example": "12345678901"
                    },
                    "amplitude": {
                      "type": "string",
                      "description": "Amplitude Ratio",
                      "example": "0.1111"
                    }
                  },
                  "description": "Data list",
                  "title": "StandardGainersRankItemVO"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "GainersRankPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Top Gainers/Losers",
    "description": {
      "content": "Retrieves a ranked list of top gaining or losing stocks for a specified time period. To get top gainers, pass direction=DESC; for top losers, pass direction=ASC. The rank_type parameter controls the time window (e.g., D1=today, W52=52-week)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "screeners",
        "gainers-losers",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "Ranking type",
            "type": "text/plain"
          },
          "key": "rank_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security category",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort Field",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sorting direction. Ascending order: ASC, descending order: DESC",
            "type": "text/plain"
          },
          "key": "direction",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key returned from previous page response. Pass null or omit for first page.",
            "type": "text/plain"
          },
          "key": "pagination_key",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Top Active

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/top-active-using-get-new.md>

### List Top Actives

Retrieves stock top active rank list

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/screeners/top-actives/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "Retrieves stock top active rank list",
  "operationId": "topActiveUsingGETNew",
  "parameters": [
    {
      "name": "rank_type",
      "in": "query",
      "description": "Rank list type",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Rank Types",
        "enum": [
          "VOLUME",
          "RELATIVE_VOLUME_10D",
          "TURNOVER",
          "TURNOVER_RATE",
          "AMPLITUDE"
        ]
      },
      "example": "VOLUME"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security category",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "sort_by",
      "in": "query",
      "description": "Sort Field",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "CHANGE_RATIO",
          "RELATIVE_VOLUME_10D",
          "MARKET_VALUE",
          "CLOSE",
          "PRICE",
          "PE_TTM",
          "HIGH",
          "LOW",
          "AMPLITUDE",
          "TURNOVER",
          "VOLUME"
        ]
      },
      "example": "VOLUME"
    },
    {
      "name": "direction",
      "in": "query",
      "description": "Sorting direction. Ascending order: ASC, descending order: DESC",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "DESC"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key returned from previous page response. Pass null or omit for first page.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "data": {
                "type": "array",
                "description": "Data list",
                "items": {
                  "type": "object",
                  "properties": {
                    "symbol": {
                      "type": "string",
                      "description": "Security symbol",
                      "example": "TSLA"
                    },
                    "name": {
                      "type": "string",
                      "description": "Security name",
                      "example": "Tesla Inc"
                    },
                    "exchange_code": {
                      "type": "string",
                      "description": "Exchange code",
                      "example": "NSQ"
                    },
                    "currency": {
                      "type": "string",
                      "description": "Currency code",
                      "example": "USD"
                    },
                    "pre_close": {
                      "type": "string",
                      "description": "The closing price of the previous trading day",
                      "example": "11.01"
                    },
                    "open": {
                      "type": "string",
                      "description": "Open price for the current trading day",
                      "example": "11.01"
                    },
                    "high": {
                      "type": "string",
                      "description": "Today’s high",
                      "example": "12.11"
                    },
                    "low": {
                      "type": "string",
                      "description": "Today’s low",
                      "example": "11.01"
                    },
                    "close": {
                      "type": "string",
                      "description": "Latest intraday prices for the current trading day",
                      "example": "12.01"
                    },
                    "price": {
                      "type": "string",
                      "description": "The latest price for the current trading day",
                      "example": "11.11"
                    },
                    "change": {
                      "type": "string",
                      "description": "Trade change for the current trading day",
                      "example": "1.01"
                    },
                    "change_ratio": {
                      "type": "string",
                      "description": "Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%)",
                      "example": "0.0111"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Trade volume",
                      "format": "0",
                      "example": "1111111"
                    },
                    "turnover": {
                      "type": "string",
                      "description": "Transaction amount, US market stocks and ETFs do not return this data under Nb authorization",
                      "format": "0.00##",
                      "example": "1234567890"
                    },
                    "turnover_rate": {
                      "type": "string",
                      "description": "Turnover Rate",
                      "example": "0.1111"
                    },
                    "market_value": {
                      "type": "string",
                      "description": "Market Value",
                      "format": "0.00##",
                      "example": "12345678901"
                    },
                    "amplitude": {
                      "type": "string",
                      "description": "Amplitude Ratio",
                      "example": "0.1111"
                    },
                    "relative_volume_10d": {
                      "type": "string",
                      "description": "10 day average trading volume ratio",
                      "format": "0.00##",
                      "example": "11.91"
                    }
                  },
                  "description": "Data list",
                  "title": "StandardActiveRankItemVO"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "ActiveRankPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Top Actives",
    "description": {
      "content": "Retrieves stock top active rank list",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "screeners",
        "top-actives",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "Rank list type",
            "type": "text/plain"
          },
          "key": "rank_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security category",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort Field",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sorting direction. Ascending order: ASC, descending order: DESC",
            "type": "text/plain"
          },
          "key": "direction",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key returned from previous page response. Pass null or omit for first page.",
            "type": "text/plain"
          },
          "key": "pagination_key",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Snapshot

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/snapshot-using-get.md>

### List Stock Snapshots

Retrieves real-time market snapshot data for a security. Returns key market indicators such as latest price, price change, volume, turnover rate, etc. Supports querying various security types including US stocks, etc., with optional inclusion of pre-market, after-hours, and overnight trading data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/stocks/snapshots/list",
  "method": "post",
  "tags": [
    "Market Data"
  ],
  "description": "Retrieves real-time market snapshot data for a security. Returns key market indicators such as latest price, price change, volume, turnover rate, etc. Supports querying various security types including US stocks, etc., with optional inclusion of pre-market, after-hours, and overnight trading data.",
  "operationId": "snapshotUsingGET",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "category_symbols"
          ],
          "type": "object",
          "properties": {
            "category_symbols": {
              "type": "array",
              "description": "List of security symbols by category; maximum 100 symbols per query.",
              "items": {
                "required": [
                  "category",
                  "symbols"
                ],
                "type": "object",
                "properties": {
                  "category": {
                    "type": "string",
                    "description": "Security type. Category values are as shown in the enum.",
                    "example": "US_STOCK",
                    "enum": [
                      "US_STOCK"
                    ]
                  },
                  "symbols": {
                    "type": "array",
                    "description": "List of security symbols, supports JSON array format.",
                    "example": [
                      "AAPL",
                      "GOOG"
                    ],
                    "items": {
                      "type": "string",
                      "description": "List of security symbols, supports JSON array format.",
                      "example": "[\"AAPL\",\"GOOG\"]"
                    }
                  }
                },
                "description": "List of security symbols by security type",
                "title": "MultiCategory"
              }
            },
            "extend_hour_required": {
              "type": "string",
              "description": "Whether to include extend hour trading data.",
              "example": "false",
              "default": "false"
            },
            "overnight_required": {
              "type": "string",
              "description": "Whether to include overnight trading data.",
              "example": "false",
              "default": "false"
            }
          },
          "title": "SnapshotQuery"
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "symbol"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "AAPL"
              },
              "pre_close": {
                "type": "string",
                "description": "Previous close price",
                "example": "101"
              },
              "change_ratio": {
                "type": "string",
                "description": "Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%)",
                "example": "0.05"
              },
              "last_trade_time": {
                "type": "integer",
                "description": "Last trade time",
                "format": "int64",
                "example": 1640688000000
              },
              "price": {
                "type": "string",
                "description": "Current price",
                "example": "100"
              },
              "open": {
                "type": "string",
                "description": "Open price, for US stocks it's intraday open price, excluding pre/post market data. No return value if no trading occurred on the day",
                "example": "100"
              },
              "high": {
                "type": "string",
                "description": "Today's high price, for US stocks it's intraday high, excluding pre/post market data. No return value if no trading occurred on the day",
                "example": "105"
              },
              "low": {
                "type": "string",
                "description": "Today's low price, for US stocks it's intraday low, excluding pre/post market data. No return value if no trading occurred on the day",
                "example": "99"
              },
              "volume": {
                "type": "string",
                "description": "Volume. No return value if no trading occurred on the day",
                "example": "1000"
              },
              "change": {
                "type": "string",
                "description": "Change amount. No return value if no trading occurred on the day",
                "example": "1.0"
              },
              "close": {
                "type": "string",
                "description": "Intraday close price",
                "example": "101"
              },
              "ask": {
                "type": "string",
                "description": "Ask Price",
                "example": "13.99"
              },
              "ask_size": {
                "type": "string",
                "description": "Ask Size (Quantity)",
                "example": "50"
              },
              "bid": {
                "type": "string",
                "description": "Bid Price",
                "example": "13.91"
              },
              "bid_size": {
                "type": "string",
                "description": "Bid Size (Quantity)",
                "example": "100"
              },
              "extend_hour_last_price": {
                "type": "string",
                "description": "Pre/post market latest price",
                "example": "100.5"
              },
              "extend_hour_change": {
                "type": "string",
                "description": "Pre/post market change amount",
                "example": "0.5"
              },
              "extend_hour_change_ratio": {
                "type": "string",
                "description": "Pre/post market change ratio",
                "example": "0.005"
              },
              "extend_hour_volume": {
                "type": "string",
                "description": "Pre/post market volume",
                "example": "200"
              },
              "extend_hour_last_trade_time": {
                "type": "integer",
                "description": "Current pre/post market trade time",
                "format": "int64",
                "example": 1640688000000
              },
              "extend_hour_high": {
                "type": "string",
                "description": "Pre/post market high price",
                "example": "188.1"
              },
              "extend_hour_low": {
                "type": "string",
                "description": "Pre/post market low price",
                "example": "180.1"
              },
              "ovn_price": {
                "type": "string",
                "description": "Overnight price",
                "example": "100.25"
              },
              "ovn_high": {
                "type": "string",
                "description": "Overnight high price",
                "example": "101.0"
              },
              "ovn_low": {
                "type": "string",
                "description": "Overnight low price",
                "example": "99.5"
              },
              "ovn_volume": {
                "type": "string",
                "description": "Overnight volume",
                "example": "500"
              },
              "ovn_change": {
                "type": "string",
                "description": "Overnight change amount",
                "example": "0.25"
              },
              "ovn_change_ratio": {
                "type": "string",
                "description": "Overnight change ratio",
                "example": "0.0025"
              },
              "ovn_last_trade_time": {
                "type": "integer",
                "description": "Overnight trade time",
                "format": "int64",
                "example": 1640688000000
              },
              "ovn_ask": {
                "type": "string",
                "description": "Overnight Ask Price",
                "example": "13.91"
              },
              "ovn_ask_size": {
                "type": "string",
                "description": "Overnight Ask Size (Quantity)",
                "example": "50"
              },
              "ovn_bid": {
                "type": "string",
                "description": "Overnight Bid Price",
                "example": "13.9"
              },
              "ovn_bid_size": {
                "type": "string",
                "description": "Overnight Bid Size (Quantity)",
                "example": "100"
              },
              "pb_ratio": {
                "type": "string",
                "description": "Price to Book Ratio",
                "example": "54.47"
              },
              "ps_ratio": {
                "type": "string",
                "description": "Price to Sales Ratio",
                "example": "9.65"
              },
              "pe_ratio": {
                "type": "string",
                "description": "Price to Earnings Ratio",
                "example": "36.42"
              },
              "market_value": {
                "type": "string",
                "description": "Total market value",
                "example": "40200000000000"
              },
              "neg_market_value": {
                "type": "string",
                "description": "Free-Float market value",
                "example": "39500000000000"
              },
              "yield": {
                "type": "string",
                "description": "Dividend yield",
                "example": "1.04"
              },
              "total_shares": {
                "type": "string",
                "description": "Total shares outstanding",
                "example": "147760000000"
              },
              "out_standing_shares": {
                "type": "string",
                "description": "Free-Float shares outstanding,",
                "example": "145260000000"
              },
              "fifty_two_wk_high": {
                "type": "string",
                "description": "52 weeks high",
                "example": "288.62"
              },
              "fifty_two_wk_low": {
                "type": "string",
                "description": "52 weeks low",
                "example": "168.63"
              }
            },
            "description": "Market snapshot data response object",
            "title": "SnapshotVo"
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": {
    "category_symbols": [
      {
        "category": "US_STOCK",
        "symbols": [
          "AAPL",
          "GOOG"
        ]
      }
    ],
    "extend_hour_required": "false",
    "overnight_required": "false"
  },
  "postman": {
    "name": "List Stock Snapshots",
    "description": {
      "content": "Retrieves real-time market snapshot data for a security. Returns key market indicators such as latest price, price change, volume, turnover rate, etc. Supports querying various security types including US stocks, etc., with optional inclusion of pre-market, after-hours, and overnight trading data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "snapshots",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```

## Historical Bars (Batch)

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/query-batch-bars-using-post.md>

### List Stock Historical Bars

Retrieves the recent N bars of data based on stock symbols, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/stocks/bars/list",
  "method": "post",
  "tags": [
    "Market Data"
  ],
  "description": "Retrieves the recent N bars of data based on stock symbols, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.",
  "operationId": "queryBatchBarsUsingPOST",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "category_symbols"
          ],
          "type": "object",
          "properties": {
            "category_symbols": {
              "type": "array",
              "description": "List of security symbols by category; maximum 100 symbols per query.",
              "items": {
                "required": [
                  "category",
                  "symbols"
                ],
                "type": "object",
                "properties": {
                  "category": {
                    "type": "string",
                    "description": "Security type. Category values are as shown in the enum.",
                    "example": "US_STOCK",
                    "enum": [
                      "US_STOCK"
                    ]
                  },
                  "symbols": {
                    "type": "array",
                    "description": "List of security symbols, supports JSON array format.",
                    "example": [
                      "AAPL",
                      "GOOG"
                    ],
                    "items": {
                      "type": "string",
                      "description": "List of security symbols, supports JSON array format.",
                      "example": "[\"AAPL\",\"GOOG\"]"
                    }
                  }
                },
                "description": "List of security symbols by security type",
                "title": "MultiCategory"
              }
            },
            "count": {
              "type": "string",
              "description": "1-1200（M1：1-1650）",
              "example": "200"
            },
            "trading_sessions": {
              "type": "string",
              "description": "Specify trading session(s). Multiple sessions separated by \",\".",
              "example": "PRE,RTH,ATH,OVN"
            },
            "interval": {
              "type": "string",
              "description": "Bar time granularity:M1, M5, M15, M30, M60, M120, M240, D, W, M, Y",
              "example": "M1"
            },
            "last_time": {
              "type": "string",
              "description": "Last time of pre page. Example: 1763555670"
            },
            "real_time_required": {
              "type": "boolean",
              "description": "Return the latest trading data, default is true;\\n false: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request.\\n true: The returned data includes the latest market data.",
              "example": false
            }
          },
          "title": "StockChartsRequests"
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "instrument_id",
              "result",
              "special_times",
              "times"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "symbol",
                "example": "AAPL"
              },
              "result": {
                "type": "array",
                "description": "k-line data",
                "items": {
                  "required": [
                    "close",
                    "high",
                    "low",
                    "open",
                    "time",
                    "trading_sessions",
                    "volume"
                  ],
                  "type": "object",
                  "properties": {
                    "time": {
                      "type": "string",
                      "description": "Bar timestamp, Unix timestamp format.",
                      "example": "2021-12-28T09:00:09.945+0000"
                    },
                    "open": {
                      "type": "string",
                      "description": "Open price",
                      "example": "150.25"
                    },
                    "close": {
                      "type": "string",
                      "description": "Close price",
                      "example": "152.3"
                    },
                    "high": {
                      "type": "string",
                      "description": "High price",
                      "example": "153.15"
                    },
                    "low": {
                      "type": "string",
                      "description": "Low price",
                      "example": "149.8"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Volume",
                      "example": "1250000"
                    },
                    "trading_sessions": {
                      "type": "string",
                      "description": "Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market).",
                      "example": "RTH"
                    }
                  },
                  "description": "Single K-line price record",
                  "title": "PriceRecord"
                }
              },
              "times": {
                "type": "array",
                "description": "exchange trading hours",
                "items": {
                  "required": [
                    "end",
                    "start",
                    "trading_session"
                  ],
                  "type": "object",
                  "properties": {
                    "start": {
                      "type": "string",
                      "description": "exchange start time",
                      "example": "09:30:00"
                    },
                    "end": {
                      "type": "string",
                      "description": "exchange end time",
                      "example": "16:00:00"
                    },
                    "trading_session": {
                      "type": "string",
                      "description": "Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market).",
                      "example": "RTH"
                    }
                  },
                  "description": "exchange trading hours",
                  "title": "ExchangeTimesVo"
                }
              },
              "instrument_id": {
                "type": "string",
                "description": "instrument_id",
                "example": "913256135"
              },
              "special_times": {
                "type": "array",
                "description": "Special Trading Hours at the Exchange",
                "items": {
                  "required": [
                    "end",
                    "start",
                    "trading_session"
                  ],
                  "type": "object",
                  "properties": {
                    "start": {
                      "type": "integer",
                      "description": "exchange start time",
                      "format": "int64"
                    },
                    "end": {
                      "type": "integer",
                      "description": "exchange end time",
                      "format": "int64"
                    },
                    "trading_session": {
                      "type": "string",
                      "description": "Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market).",
                      "example": "RTH"
                    }
                  },
                  "description": "Special Trading Hours at the Exchange",
                  "title": "SpecialExchangeTimesVo"
                }
              }
            },
            "title": "StockChartsResponseVo"
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": {
    "category_symbols": [
      {
        "category": "US_STOCK",
        "symbols": [
          "AAPL",
          "GOOG"
        ]
      }
    ],
    "count": "200",
    "trading_sessions": "PRE,RTH,ATH,OVN",
    "interval": "M1",
    "last_time": "string",
    "real_time_required": false
  },
  "postman": {
    "name": "List Stock Historical Bars",
    "description": {
      "content": "Retrieves the recent N bars of data based on stock symbols, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "bars",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```

## Historical Bars (Single)

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/bars-using-get.md>

### Get Stock Historical Bars (single symbol)

Retrieves the recent N bars of data based on stock symbol, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/stocks/bars/get",
  "method": "get",
  "tags": [
    "Market Data"
  ],
  "description": "Retrieves the recent N bars of data based on stock symbol, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.",
  "operationId": "barsUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "interval",
      "in": "query",
      "description": "Bar time granularity:M1, M5, M15, M30, M60, M120, M240, D, W, M, Y",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "M1",
          "M5",
          "M15",
          "M30",
          "M60",
          "M120",
          "M240",
          "D",
          "W",
          "M",
          "Y"
        ]
      },
      "example": "D"
    },
    {
      "name": "last_time",
      "in": "query",
      "description": "Last time of pre page.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Empty represents querying the latest data. Example: 1763555670"
      }
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of bars, default 200, maximum limit 1200 (M1 supports up to 1650).",
      "required": false,
      "schema": {
        "type": "string",
        "description": "1-1200（M1：1-1650）",
        "default": "200"
      },
      "example": 500
    },
    {
      "name": "real_time_required",
      "in": "query",
      "description": "Return the latest trading data, default is true;\\n false: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request.\\n true: The returned data includes the latest market data.",
      "required": false,
      "schema": {
        "type": "string",
        "default": "true"
      },
      "example": true
    },
    {
      "name": "trading_sessions",
      "in": "query",
      "description": "Specify trading hours. Multiple selections are allowed. Separate multiple items with \",\".",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "PRE",
          "RTH",
          "ATH",
          "OVN"
        ]
      },
      "example": "PRE,RTH,ATH,OVN"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "instrument_id",
              "result",
              "special_times",
              "times"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "symbol",
                "example": "AAPL"
              },
              "result": {
                "type": "array",
                "description": "k-line data",
                "items": {
                  "required": [
                    "close",
                    "high",
                    "low",
                    "open",
                    "time",
                    "trading_sessions",
                    "volume"
                  ],
                  "type": "object",
                  "properties": {
                    "time": {
                      "type": "string",
                      "description": "Bar timestamp, Unix timestamp format.",
                      "example": "2021-12-28T09:00:09.945+0000"
                    },
                    "open": {
                      "type": "string",
                      "description": "Open price",
                      "example": "150.25"
                    },
                    "close": {
                      "type": "string",
                      "description": "Close price",
                      "example": "152.3"
                    },
                    "high": {
                      "type": "string",
                      "description": "High price",
                      "example": "153.15"
                    },
                    "low": {
                      "type": "string",
                      "description": "Low price",
                      "example": "149.8"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Volume",
                      "example": "1250000"
                    },
                    "trading_sessions": {
                      "type": "string",
                      "description": "Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market).",
                      "example": "RTH"
                    }
                  },
                  "description": "Single K-line price record",
                  "title": "PriceRecord"
                }
              },
              "times": {
                "type": "array",
                "description": "exchange trading hours",
                "items": {
                  "required": [
                    "end",
                    "start",
                    "trading_session"
                  ],
                  "type": "object",
                  "properties": {
                    "start": {
                      "type": "string",
                      "description": "exchange start time",
                      "example": "09:30:00"
                    },
                    "end": {
                      "type": "string",
                      "description": "exchange end time",
                      "example": "16:00:00"
                    },
                    "trading_session": {
                      "type": "string",
                      "description": "Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market).",
                      "example": "RTH"
                    }
                  },
                  "description": "exchange trading hours",
                  "title": "ExchangeTimesVo"
                }
              },
              "instrument_id": {
                "type": "string",
                "description": "instrument_id",
                "example": "913256135"
              },
              "special_times": {
                "type": "array",
                "description": "Special Trading Hours at the Exchange",
                "items": {
                  "required": [
                    "end",
                    "start",
                    "trading_session"
                  ],
                  "type": "object",
                  "properties": {
                    "start": {
                      "type": "integer",
                      "description": "exchange start time",
                      "format": "int64"
                    },
                    "end": {
                      "type": "integer",
                      "description": "exchange end time",
                      "format": "int64"
                    },
                    "trading_session": {
                      "type": "string",
                      "description": "Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market).",
                      "example": "RTH"
                    }
                  },
                  "description": "Special Trading Hours at the Exchange",
                  "title": "SpecialExchangeTimesVo"
                }
              }
            },
            "title": "StockChartsResponseVo"
          }
        }
      }
    }
  },
  "postman": {
    "name": "Get Stock Historical Bars (single symbol)",
    "description": {
      "content": "Retrieves the recent N bars of data based on stock symbol, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "bars",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Bar time granularity:M1, M5, M15, M30, M60, M120, M240, D, W, M, Y",
            "type": "text/plain"
          },
          "key": "interval",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Last time of pre page.",
            "type": "text/plain"
          },
          "key": "last_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of bars, default 200, maximum limit 1200 (M1 supports up to 1650).",
            "type": "text/plain"
          },
          "key": "count",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Return the latest trading data, default is true;\\n false: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request.\\n true: The returned data includes the latest market data.",
            "type": "text/plain"
          },
          "key": "real_time_required",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Specify trading hours. Multiple selections are allowed. Separate multiple items with \",\".",
            "type": "text/plain"
          },
          "key": "trading_sessions",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Tick

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/tick-using-get.md>

### List Stock Ticks

Retrieves tick-by-tick trade data for a security. Returns detailed tick trade records within a specified time range for a given security, including trade time, price, volume, direction, and other details. Data is sorted in reverse chronological order (latest first).

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/stocks/ticks/list",
  "method": "get",
  "tags": [
    "Market Data"
  ],
  "description": "Retrieves tick-by-tick trade data for a security. Returns detailed tick trade records within a specified time range for a given security, including trade time, price, volume, direction, and other details. Data is sorted in reverse chronological order (latest first).",
  "operationId": "tickUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "last_time",
      "in": "query",
      "description": "Last time of pre page.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Empty represents querying the latest data. Example: 1763555670"
      }
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of ticks, default 100, maximum limit 1000.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "30"
      },
      "example": 10
    },
    {
      "name": "trading_sessions",
      "in": "query",
      "description": "Specify trading hours. Multiple selections are allowed. Separate multiple items with \",\".",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "PRE",
          "RTH",
          "ATH",
          "OVN"
        ]
      },
      "example": "PRE,RTH"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "result",
              "symbol"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "AAPL"
              },
              "result": {
                "type": "array",
                "description": "Tick details",
                "items": {
                  "required": [
                    "price",
                    "side",
                    "time",
                    "trading_session",
                    "volume"
                  ],
                  "type": "object",
                  "properties": {
                    "time": {
                      "type": "string",
                      "description": "Trade time",
                      "example": "1669882089000"
                    },
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "294.4"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Volume",
                      "example": "300"
                    },
                    "side": {
                      "type": "string",
                      "description": "Trade direction, see trade direction for details",
                      "example": "N"
                    },
                    "trading_session": {
                      "type": "string",
                      "description": "Trading session",
                      "example": "RTH"
                    }
                  },
                  "description": "Tick detail",
                  "title": "Tick"
                }
              }
            },
            "description": "TickVo",
            "title": "TickVo"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock Ticks",
    "description": {
      "content": "Retrieves tick-by-tick trade data for a security. Returns detailed tick trade records within a specified time range for a given security, including trade time, price, volume, direction, and other details. Data is sorted in reverse chronological order (latest first).",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "ticks",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Last time of pre page.",
            "type": "text/plain"
          },
          "key": "last_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Number of ticks, default 100, maximum limit 1000.",
            "type": "text/plain"
          },
          "key": "count",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Specify trading hours. Multiple selections are allowed. Separate multiple items with \",\".",
            "type": "text/plain"
          },
          "key": "trading_sessions",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Quotes Depth

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/quotes-using-get.md>

### List Stock Depthes

Retrieves the latest bid/ask data for a security. Returns bid/ask information for a specified depth, including price, quantity, order details, etc.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/stocks/depths/list",
  "method": "get",
  "tags": [
    "Market Data"
  ],
  "description": "Retrieves the latest bid/ask data for a security. Returns bid/ask information for a specified depth, including price, quantity, order details, etc.",
  "operationId": "quotesUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "GOOG"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "depth",
      "in": "query",
      "description": "Market depth, L2-default 10 levels, etc.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": 10
    },
    {
      "name": "overnight_required",
      "in": "query",
      "description": "Whether to include overnight trading data.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "false"
      },
      "example": false
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "asks",
              "bids",
              "quote_time",
              "symbol"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "F"
              },
              "quote_time": {
                "type": "string",
                "description": "Quote time",
                "example": "1640688000000"
              },
              "asks": {
                "type": "array",
                "description": "Array of ask orders",
                "items": {
                  "required": [
                    "order",
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "13.9"
                    },
                    "size": {
                      "type": "string",
                      "description": "Size (Quantity)",
                      "example": "5"
                    },
                    "order": {
                      "type": "array",
                      "description": "Array of order details",
                      "items": {
                        "required": [
                          "mpid",
                          "size"
                        ],
                        "type": "object",
                        "properties": {
                          "mpid": {
                            "type": "string",
                            "description": "Market participant ID",
                            "example": "NSDQ"
                          },
                          "size": {
                            "type": "string",
                            "description": "Size (Quantity)",
                            "example": "5"
                          }
                        },
                        "description": "OrderItem",
                        "title": "OrderItem"
                      }
                    }
                  },
                  "description": "AskItem",
                  "title": "AskItem"
                }
              },
              "bids": {
                "type": "array",
                "description": "Array of bid orders",
                "items": {
                  "required": [
                    "order",
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "13.9"
                    },
                    "size": {
                      "type": "string",
                      "description": "Size (Quantity)",
                      "example": "5"
                    },
                    "order": {
                      "type": "array",
                      "description": "Array of order details",
                      "items": {
                        "required": [
                          "mpid",
                          "size"
                        ],
                        "type": "object",
                        "properties": {
                          "mpid": {
                            "type": "string",
                            "description": "Market participant ID",
                            "example": "NSDQ"
                          },
                          "size": {
                            "type": "string",
                            "description": "Size (Quantity)",
                            "example": "5"
                          }
                        },
                        "description": "OrderItem",
                        "title": "OrderItem"
                      }
                    }
                  },
                  "description": "AskItem",
                  "title": "AskItem"
                }
              }
            },
            "description": "DepthVo",
            "title": "DepthVo"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock Depthes",
    "description": {
      "content": "Retrieves the latest bid/ask data for a security. Returns bid/ask information for a specified depth, including price, quantity, order details, etc.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "depths",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Market depth, L2-default 10 levels, etc.",
            "type": "text/plain"
          },
          "key": "depth",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Whether to include overnight trading data.",
            "type": "text/plain"
          },
          "key": "overnight_required",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## News Summary

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/watchlist-summary-using-post.md>

### Get News Summary

Invokes LLM to generate news summaries for watchlist.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/news/summaries/get",
  "method": "post",
  "tags": [
    "News"
  ],
  "description": "Invokes LLM to generate news summaries for watchlist.",
  "operationId": "watchlistSummaryUsingPOST",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "type": "object",
          "properties": {
            "category_symbols": {
              "type": "array",
              "description": "List of security symbols by category.",
              "items": {
                "required": [
                  "category",
                  "symbols"
                ],
                "type": "object",
                "properties": {
                  "category": {
                    "type": "string",
                    "description": "Security type. Category values are as shown in the enum.",
                    "example": "US_STOCK",
                    "enum": [
                      "US_STOCK"
                    ]
                  },
                  "symbols": {
                    "type": "array",
                    "description": "List of security symbols, supports JSON array format.",
                    "example": [
                      "AAPL",
                      "GOOG"
                    ],
                    "items": {
                      "type": "string",
                      "description": "List of security symbols, supports JSON array format.",
                      "example": "[\"AAPL\",\"GOOG\"]"
                    }
                  }
                },
                "description": "List of security symbols by security type",
                "title": "MultiCategory"
              }
            },
            "lang": {
              "type": "string",
              "description": "Support language, enum: [en].",
              "example": "en"
            }
          },
          "description": "Chat request parameters",
          "title": "ChatRequest"
        },
        "example": {
          "category_symbols": {
            "category": "US_STOCK",
            "symbols": [
              "AAPL",
              "GOOG"
            ]
          }
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "Server-Sent Events stream. Each message contains JSON data.\n\nExample stream:\n```\nevent:message\ndata:{\"type\":\"meta\",\"args\":{\"sessionId\":\"1\",\"convId\":451107711450219}}\n\nevent:message\ndata:{\"type\":\"text\",\"message\":\"Hi\"}\n\nevent:message\ndata:{\"type\":\"text\",\"message\":\", I'm Wally an\"}\n\nevent:message\ndata:{\"type\":\"text\",\"message\":\"d I'll help you with this question\"}\n```\n\n**Response Types:**\n- **text**: Markdown text\n- **table**: `{ headers: {text: \"<Title>\"}[], rows: {text: \"Row Text\"}[][] }`",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "description": "Message type",
                "example": "text",
                "enum": [
                  "meta",
                  "text",
                  "table"
                ]
              },
              "message": {
                "type": "string",
                "description": "Message content (for text type)",
                "example": "Hi, I'm Wally"
              },
              "args": {
                "type": "object",
                "description": "Additional arguments (for meta type)"
              },
              "headers": {
                "type": "object",
                "description": "Table headers (for table type)"
              },
              "rows": {
                "type": "object",
                "description": "Table rows (for table type)"
              }
            },
            "description": "Chat stream response message",
            "title": "ChatStreamResponse"
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": {
    "category_symbols": [
      {
        "category": "US_STOCK",
        "symbols": [
          "AAPL",
          "GOOG"
        ]
      }
    ],
    "lang": "en"
  },
  "postman": {
    "name": "Get News Summary",
    "description": {
      "content": "Invokes LLM to generate news summaries for watchlist.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "news",
        "summaries",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```

## Market News

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-news-by-market-using-get.md>

### List Market News

Retrieves news from the market within the past 3 days.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/news/market-news/list",
  "method": "get",
  "tags": [
    "News"
  ],
  "description": "Retrieves news from the market within the past 3 days.",
  "operationId": "listNewsByMarketUsingGet",
  "parameters": [
    {
      "name": "market",
      "in": "query",
      "description": "Region，eg: United States：US，Thailand：TH，default US.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "US"
    },
    {
      "name": "language",
      "in": "query",
      "description": "News language.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "en"
    },
    {
      "name": "last_news_id",
      "in": "query",
      "description": "The ID of the last data item on the previous page,default 0.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64",
        "default": 0
      },
      "example": 0
    },
    {
      "name": "page_size",
      "in": "query",
      "description": "Number of data per page, default 15.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int32",
        "default": 10
      },
      "example": 10
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "description": "News id",
                "format": "int64",
                "example": 10739820929238016
              },
              "title": {
                "type": "string",
                "description": "News title",
                "example": "Apple Store Workers In Maryland Vote In Favor Of Strike Over Working Conditions"
              },
              "source_name": {
                "type": "string",
                "description": "News publish source",
                "example": "Benzinga"
              },
              "news_time": {
                "type": "string",
                "description": "News publish time",
                "format": "date-time"
              },
              "news_url": {
                "type": "string",
                "description": "News link",
                "example": "https://pre-news.webullbroker.com/us/news-html/20240513/10739820929238016.html"
              },
              "thumbnail": {
                "type": "string",
                "description": "News thumbnail",
                "example": "https://news-static.webullfintech.com/us/news-pic/20260107/14163175251215360.jpg"
              }
            },
            "description": "News data response object",
            "title": "NewsVO"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Market News",
    "description": {
      "content": "Retrieves news from the market within the past 3 days.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "news",
        "market-news",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Region，eg: United States：US，Thailand：TH，default US.",
            "type": "text/plain"
          },
          "key": "market",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) News language.",
            "type": "text/plain"
          },
          "key": "language",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The ID of the last data item on the previous page,default 0.",
            "type": "text/plain"
          },
          "key": "last_news_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of data per page, default 15.",
            "type": "text/plain"
          },
          "key": "page_size",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Symbol News

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-news-by-ticker-using-get.md>

### List Symbol News

Get news on stocks within the past 3 days.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/news/symbol-news/list",
  "method": "get",
  "tags": [
    "News"
  ],
  "description": "Get news on stocks within the past 3 days.",
  "operationId": "listNewsByTickerUsingGet",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "GOOG"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security category.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "US_STOCK"
    },
    {
      "name": "language",
      "in": "query",
      "description": "News language.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "en"
    },
    {
      "name": "last_news_id",
      "in": "query",
      "description": "The ID of the last data item on the previous page,default 0.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64",
        "default": 0
      },
      "example": 0
    },
    {
      "name": "page_size",
      "in": "query",
      "description": "Number of data per page, default 15.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int32",
        "default": 10
      },
      "example": 10
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "description": "News id",
                "format": "int64",
                "example": 10739820929238016
              },
              "title": {
                "type": "string",
                "description": "News title",
                "example": "Apple Store Workers In Maryland Vote In Favor Of Strike Over Working Conditions"
              },
              "source_name": {
                "type": "string",
                "description": "News publish source",
                "example": "Benzinga"
              },
              "news_time": {
                "type": "string",
                "description": "News publish time",
                "format": "date-time"
              },
              "news_url": {
                "type": "string",
                "description": "News link",
                "example": "https://pre-news.webullbroker.com/us/news-html/20240513/10739820929238016.html"
              },
              "thumbnail": {
                "type": "string",
                "description": "News thumbnail",
                "example": "https://news-static.webullfintech.com/us/news-pic/20260107/14163175251215360.jpg"
              }
            },
            "description": "News data response object",
            "title": "NewsVO"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Symbol News",
    "description": {
      "content": "Get news on stocks within the past 3 days.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "news",
        "symbol-news",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security category.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) News language.",
            "type": "text/plain"
          },
          "key": "language",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The ID of the last data item on the previous page,default 0.",
            "type": "text/plain"
          },
          "key": "last_news_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of data per page, default 15.",
            "type": "text/plain"
          },
          "key": "page_size",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Latest News

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-latest-news-using-get.md>

### List Latest News

Retrieves latest news within the past 3 days.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/news/latest-news/list",
  "method": "get",
  "tags": [
    "News"
  ],
  "description": "Retrieves latest news within the past 3 days.",
  "operationId": "listLatestNewsUsingGet",
  "parameters": [
    {
      "name": "language",
      "in": "query",
      "description": "News language.",
      "required": false,
      "schema": {
        "type": "string",
        "default": "en"
      },
      "example": "en"
    },
    {
      "name": "last_news_id",
      "in": "query",
      "description": "The ID of the last data item on the previous page,default 0.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64",
        "default": 0
      },
      "example": 0
    },
    {
      "name": "page_size",
      "in": "query",
      "description": "Number of data per page, default 15.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int32",
        "default": 10
      },
      "example": 10
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "description": "News id",
                "format": "int64",
                "example": 10739820929238016
              },
              "title": {
                "type": "string",
                "description": "News title",
                "example": "Apple Store Workers In Maryland Vote In Favor Of Strike Over Working Conditions"
              },
              "source_name": {
                "type": "string",
                "description": "News publish source",
                "example": "Benzinga"
              },
              "news_time": {
                "type": "string",
                "description": "News publish time",
                "format": "date-time"
              },
              "news_url": {
                "type": "string",
                "description": "News link",
                "example": "https://pre-news.webullbroker.com/us/news-html/20240513/10739820929238016.html"
              },
              "thumbnail": {
                "type": "string",
                "description": "News thumbnail",
                "example": "https://news-static.webullfintech.com/us/news-pic/20260107/14163175251215360.jpg"
              }
            },
            "description": "News data response object",
            "title": "NewsVO"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Latest News",
    "description": {
      "content": "Retrieves latest news within the past 3 days.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "news",
        "latest-news",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "News language.",
            "type": "text/plain"
          },
          "key": "language",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The ID of the last data item on the previous page,default 0.",
            "type": "text/plain"
          },
          "key": "last_news_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of data per page, default 15.",
            "type": "text/plain"
          },
          "key": "page_size",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Corporate Actions

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get.md>

### List Corporate Actions

Supports the query of the corporate events for stock splits and reverse stock split, including past and upcoming events.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/instruments/stocks/corporate-actions/list",
  "method": "get",
  "tags": [
    "Corporate Actions"
  ],
  "description": "Supports the query of the corporate events for stock splits and reverse stock split, including past and upcoming events.",
  "operationId": "corpActionUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "TSLA"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "US_STOCK"
    },
    {
      "name": "start_date",
      "in": "query",
      "description": "Event start date, UTC time. Format: yyyy-MM-dd",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2025-11-01"
    },
    {
      "name": "end_date",
      "in": "query",
      "description": "Event end date, UTC time. Format: yyyy-MM-dd",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2025-12-31"
    },
    {
      "name": "event_types",
      "in": "query",
      "description": "Event type collection. Multiple event_types should be separated by ,",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "REVERSE_SPLIT,FORWARD_SPLIT"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key returned from previous page response. Pass null or omit for first page.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "data": {
                "type": "array",
                "description": "Data list",
                "items": {
                  "type": "object",
                  "properties": {
                    "instrument_id": {
                      "type": "integer",
                      "description": "Security ID",
                      "format": "int32",
                      "example": 950175018
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Security symbol, e.g., AAPL, GOOG.",
                      "example": "AAPL"
                    },
                    "exchange_code": {
                      "type": "string",
                      "description": "Exchange code, e.g., NAS, OTC.",
                      "example": "NAS"
                    },
                    "event_type": {
                      "type": "string",
                      "description": "Event type",
                      "example": "REVERSE_SPLIT",
                      "enum": [
                        "NAME_CHANGE",
                        "CASH_DIVIDEND",
                        "STOCK_DIVIDEND",
                        "REVERSE_SPLIT",
                        "FORWARD_SPLIT",
                        "SPIN_OFF",
                        "UNIT_SPLIT",
                        "MERGER",
                        "REDEMPTION"
                      ]
                    },
                    "event_action": {
                      "type": "string",
                      "description": "Event status, e.g., I(Insert, Valid)/U(Update, Valid)/C(Cancellation, invalid)/D(Deletion, invalid).",
                      "example": "I"
                    },
                    "event_id": {
                      "type": "integer",
                      "description": "Event id",
                      "format": "int32",
                      "example": 500051270
                    },
                    "source": {
                      "type": "string",
                      "description": "Event source, e.g., WEBULL_ARTIFICIAL(Webull artificial)",
                      "example": "WEBULL_ARTIFICIAL"
                    },
                    "ratio_old": {
                      "type": "string",
                      "description": "Old ratio, Before the change",
                      "example": "10"
                    },
                    "ratio_new": {
                      "type": "string",
                      "description": "New ratio, After the change, During the REVERSE_SPLIT process, when ratio_old is 10 and ratio_new is 5, it means that 10 shares are combined into 5.",
                      "example": "5"
                    },
                    "event_date": {
                      "type": "string",
                      "description": "Event date, UTC time, e.g: 2021-12-28",
                      "example": "2021-12-28"
                    },
                    "update_time": {
                      "type": "string",
                      "description": "Update time, UTC time, e.g: 2021-12-28T09:00:09.945+0000",
                      "example": "2021-12-28T09:00:09.945+0000"
                    },
                    "create_time": {
                      "type": "string",
                      "description": "Create time, UTC time, e.g: 2021-12-28T09:00:09.945+0000",
                      "example": "2021-12-28T09:00:09.945+0000"
                    }
                  },
                  "description": "Corporate Actions data response object",
                  "title": "CorpVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "CorpVoPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Corporate Actions",
    "description": {
      "content": "Supports the query of the corporate events for stock splits and reverse stock split, including past and upcoming events.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "stocks",
        "corporate-actions",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Event start date, UTC time. Format: yyyy-MM-dd",
            "type": "text/plain"
          },
          "key": "start_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Event end date, UTC time. Format: yyyy-MM-dd",
            "type": "text/plain"
          },
          "key": "end_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Event type collection. Multiple event_types should be separated by ,",
            "type": "text/plain"
          },
          "key": "event_types",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key returned from previous page response. Pass null or omit for first page.",
            "type": "text/plain"
          },
          "key": "pagination_key",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Corporate Actions By Market

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get.md>

### List Corporate Actions by Market

Retrieves corporate action events for all securities in a specified market within a date range. Use this endpoint to bulk-fetch events across the entire US market, rather than querying by individual symbol.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/instruments/stocks/corporate-actions/list-by-market",
  "method": "get",
  "tags": [
    "Corporate Actions"
  ],
  "description": "Retrieves corporate action events for all securities in a specified market within a date range. Use this endpoint to bulk-fetch events across the entire US market, rather than querying by individual symbol.",
  "operationId": "corpMarketUsingGET",
  "parameters": [
    {
      "name": "market",
      "in": "query",
      "description": "Currently only `US` (US Stock Market) is supported.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "US"
    },
    {
      "name": "start_date",
      "in": "query",
      "description": "Event start date, UTC time. Format: yyyy-MM-dd",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2025-11-01"
    },
    {
      "name": "end_date",
      "in": "query",
      "description": "Event end date, UTC time. Format: yyyy-MM-dd",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2025-12-31"
    },
    {
      "name": "event_types",
      "in": "query",
      "description": "Event type collection. Multiple event_types should be separated by ,",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "REVERSE_SPLIT,FORWARD_SPLIT"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key returned from previous page response. Pass null or omit for first page.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "data": {
                "type": "array",
                "description": "Data list",
                "items": {
                  "type": "object",
                  "properties": {
                    "instrument_id": {
                      "type": "integer",
                      "description": "Security ID",
                      "format": "int32",
                      "example": 950175018
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Security symbol, e.g., AAPL, GOOG.",
                      "example": "AAPL"
                    },
                    "exchange_code": {
                      "type": "string",
                      "description": "Exchange code, e.g., NAS, OTC.",
                      "example": "NAS"
                    },
                    "event_type": {
                      "type": "string",
                      "description": "Event type",
                      "example": "REVERSE_SPLIT",
                      "enum": [
                        "NAME_CHANGE",
                        "CASH_DIVIDEND",
                        "STOCK_DIVIDEND",
                        "REVERSE_SPLIT",
                        "FORWARD_SPLIT",
                        "SPIN_OFF",
                        "UNIT_SPLIT",
                        "MERGER",
                        "REDEMPTION"
                      ]
                    },
                    "event_action": {
                      "type": "string",
                      "description": "Event status, e.g., I(Insert, Valid)/U(Update, Valid)/C(Cancellation, invalid)/D(Deletion, invalid).",
                      "example": "I"
                    },
                    "event_id": {
                      "type": "integer",
                      "description": "Event id",
                      "format": "int32",
                      "example": 500051270
                    },
                    "source": {
                      "type": "string",
                      "description": "Event source, e.g., WEBULL_ARTIFICIAL(Webull artificial)",
                      "example": "WEBULL_ARTIFICIAL"
                    },
                    "ratio_old": {
                      "type": "string",
                      "description": "Old ratio, Before the change",
                      "example": "10"
                    },
                    "ratio_new": {
                      "type": "string",
                      "description": "New ratio, After the change, During the REVERSE_SPLIT process, when ratio_old is 10 and ratio_new is 5, it means that 10 shares are combined into 5.",
                      "example": "5"
                    },
                    "event_date": {
                      "type": "string",
                      "description": "Event date, UTC time, e.g: 2021-12-28",
                      "example": "2021-12-28"
                    },
                    "update_time": {
                      "type": "string",
                      "description": "Update time, UTC time, e.g: 2021-12-28T09:00:09.945+0000",
                      "example": "2021-12-28T09:00:09.945+0000"
                    },
                    "create_time": {
                      "type": "string",
                      "description": "Create time, UTC time, e.g: 2021-12-28T09:00:09.945+0000",
                      "example": "2021-12-28T09:00:09.945+0000"
                    }
                  },
                  "description": "Corporate Actions data response object",
                  "title": "CorpVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "CorpVoPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Corporate Actions by Market",
    "description": {
      "content": "Retrieves corporate action events for all securities in a specified market within a date range. Use this endpoint to bulk-fetch events across the entire US market, rather than querying by individual symbol.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "stocks",
        "corporate-actions",
        "list-by-market"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Currently only `US` (US Stock Market) is supported.",
            "type": "text/plain"
          },
          "key": "market",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Event start date, UTC time. Format: yyyy-MM-dd",
            "type": "text/plain"
          },
          "key": "start_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Event end date, UTC time. Format: yyyy-MM-dd",
            "type": "text/plain"
          },
          "key": "end_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Event type collection. Multiple event_types should be separated by ,",
            "type": "text/plain"
          },
          "key": "event_types",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key returned from previous page response. Pass null or omit for first page.",
            "type": "text/plain"
          },
          "key": "pagination_key",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Get Instruments

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-using-get.md>

### List Stock Instruments

Retrieves security information for one or more instruments.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/instruments/stocks/profiles/list",
  "method": "post",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves security information for one or more instruments.",
  "operationId": "listUsingGET",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "type": "array",
          "items": {
            "required": [
              "category",
              "symbols"
            ],
            "type": "object",
            "properties": {
              "category": {
                "type": "string",
                "description": "Security type. Category values are as shown in the enum.",
                "example": "US_STOCK",
                "enum": [
                  "US_STOCK"
                ]
              },
              "symbols": {
                "type": "array",
                "description": "List of security symbols, supports JSON array format.",
                "example": [
                  "AAPL",
                  "GOOG"
                ],
                "items": {
                  "type": "string",
                  "description": "List of security symbols, supports JSON array format.",
                  "example": "[\"AAPL\",\"GOOG\"]"
                }
              }
            },
            "description": "List of security symbols by security type",
            "title": "MultiCategory"
          }
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "name": {
                  "type": "string",
                  "description": "Security name",
                  "example": "Apple"
                },
                "symbol": {
                  "type": "string",
                  "description": "Security symbol.",
                  "example": "AAPL"
                },
                "category": {
                  "type": "string",
                  "description": "Security category.",
                  "example": "US_STOCK"
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Exchange code",
                  "example": "NAQ"
                },
                "currency": {
                  "type": "string",
                  "description": "Currency",
                  "example": "USD"
                },
                "subtype": {
                  "type": "string",
                  "description": "Security Subtypes",
                  "example": "COMMON_STOCK",
                  "enum": [
                    "COMMON_STOCK",
                    "ETF",
                    "INDEX",
                    "PREFERRED_STOCK",
                    "WARRANT",
                    "UNITS",
                    "RIGHT"
                  ]
                },
                "is_adr": {
                  "type": "string",
                  "description": "Is ADR, true or false",
                  "example": "false"
                }
              },
              "description": "Securities Information",
              "title": "InstrumentVo"
            }
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": [
    {
      "category": "US_STOCK",
      "symbols": [
        "AAPL",
        "GOOG"
      ]
    }
  ],
  "postman": {
    "name": "List Stock Instruments",
    "description": {
      "content": "Retrieves security information for one or more instruments.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "stocks",
        "profiles",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```

## Batch Logos

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/batch-logo-using-post.md>

### List Logos

Retrieves logo image URLs for the specified securities. URLs are hosted on Webull's CDN.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/fundamentals/logos/list",
  "method": "post",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves logo image URLs for the specified securities. URLs are hosted on Webull's CDN.",
  "operationId": "batchLogoUsingPOST",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "type": "array",
          "items": {
            "required": [
              "category",
              "symbols"
            ],
            "type": "object",
            "properties": {
              "category": {
                "type": "string",
                "description": "Security type. Category values are as shown in the enum.",
                "example": "US_STOCK",
                "enum": [
                  "US_STOCK"
                ]
              },
              "symbols": {
                "type": "array",
                "description": "List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
                "example": [
                  "AAPL",
                  "GOOG"
                ],
                "items": {
                  "type": "string",
                  "description": "List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
                  "example": "[\"AAPL\",\"GOOG\"]"
                }
              }
            },
            "title": "TickerLogoParam"
          }
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "AAPL"
              },
              "category": {
                "type": "string",
                "description": "Security type. Category values are as shown in the enum",
                "example": "US_STOCK",
                "enum": [
                  "US_STOCK"
                ]
              },
              "logo_url": {
                "type": "string",
                "description": "Ticker's icon url. Returns null if no logo is available for the symbol",
                "example": "https://quotes-static.webullfintech.com/ticker-icon/950982786.png"
              }
            },
            "description": "Ticker Logos response object",
            "title": "TickerLogoVo"
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": [
    {
      "category": "US_STOCK",
      "symbols": [
        "AAPL",
        "GOOG"
      ]
    }
  ],
  "postman": {
    "name": "List Logos",
    "description": {
      "content": "Retrieves logo image URLs for the specified securities. URLs are hosted on Webull's CDN.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "fundamentals",
        "logos",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```

## Company Profile

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-company-profile-using-get.md>

### Get Company Profile

Retrieves company profile for one instrument.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/fundamentals/company-profiles/get",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves company profile for one instrument.",
  "operationId": "listCompanyProfileUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum. default is US_STOCK",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "NVDA"
              },
              "category": {
                "type": "string",
                "description": "Security type",
                "example": "US_STOCK"
              },
              "company_name": {
                "type": "string",
                "description": "Company name",
                "example": "NVIDIA Corp"
              },
              "establish_date": {
                "type": "string",
                "description": "Date of incorporation",
                "example": "1998-02-24"
              },
              "exhibition_code": {
                "type": "string",
                "description": "The exchange or market where the security is listed (e.g., NASDAQ, NYSE)",
                "example": "NASDAQ"
              },
              "profile": {
                "type": "string",
                "description": "Company profile",
                "example": "NVIDIA Corporation is a full-stack computing infrastructure company."
              },
              "employees": {
                "type": "string",
                "description": "Number of employees",
                "example": "36000"
              },
              "address": {
                "type": "string",
                "description": "Headquarters address",
                "example": "2788 San Tomas Expressway,SANTA CLARA,CA,United States"
              },
              "ceo": {
                "type": "string",
                "description": "Company CEO",
                "example": "Jen-Hsun Huang"
              },
              "industries": {
                "type": "array",
                "description": "Company industries",
                "example": [
                  "Semiconductors",
                  "Semiconductors & Semiconductor Equipment"
                ],
                "items": {
                  "type": "string",
                  "description": "Company industries",
                  "example": "[\"Semiconductors\",\"Semiconductors & Semiconductor Equipment\"]"
                }
              }
            },
            "title": "CompanyProfile"
          }
        }
      }
    }
  },
  "postman": {
    "name": "Get Company Profile",
    "description": {
      "content": "Retrieves company profile for one instrument.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "fundamentals",
        "company-profiles",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum. default is US_STOCK",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Analyst Target Price

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-analyst-target-price-using-get.md>

### Get Analyst Target Price

Retrieves analyst target price for one instrument.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/fundamentals/analysis/target-prices/get",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves analyst target price for one instrument.",
  "operationId": "listAnalystTargetPriceUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum. default is US_STOCK",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "NVDA"
              },
              "category": {
                "type": "string",
                "description": "Security type",
                "example": "US_STOCK"
              },
              "mean": {
                "type": "string",
                "description": "Average target price",
                "example": "85.0"
              },
              "low": {
                "type": "string",
                "description": "Lowest target price",
                "example": "85.0"
              },
              "high": {
                "type": "string",
                "description": "Highest target price",
                "example": "200.0"
              },
              "median": {
                "type": "string",
                "description": "Median target price",
                "example": "139.72"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD"
              },
              "effective_start_date": {
                "type": "string",
                "description": "The date from which the current consensus rating is effective, in ISO 8601 format (UTC).",
                "example": "2021-12-29T06:24:56.038+0000"
              }
            },
            "title": "TargetPrice"
          }
        }
      }
    }
  },
  "postman": {
    "name": "Get Analyst Target Price",
    "description": {
      "content": "Retrieves analyst target price for one instrument.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "fundamentals",
        "analysis",
        "target-prices",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum. default is US_STOCK",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Analyst Rating

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-analyst-rating-using-get.md>

### Get Analyst Rating

Retrieves analyst rating for one instrument.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://hk-co-branding-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/fundamentals/analysis/ratings/get",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves analyst rating for one instrument.",
  "operationId": "listAnalystRatingUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum. default is US_STOCK",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "NVDA"
              },
              "category": {
                "type": "string",
                "description": "Security type",
                "example": "US_STOCK"
              },
              "number": {
                "type": "string",
                "description": "Total number of analysts",
                "example": "58"
              },
              "under_perform": {
                "type": "string",
                "description": "Under perform count",
                "example": "0"
              },
              "buy": {
                "type": "string",
                "description": "Buy count",
                "example": "11"
              },
              "sell": {
                "type": "string",
                "description": "Sell count",
                "example": "0"
              },
              "strong_buy": {
                "type": "string",
                "description": "Strong buy count",
                "example": "43"
              },
              "hold": {
                "type": "string",
                "description": "Hold (neutral) count",
                "example": "4"
              },
              "effective_start_date": {
                "type": "string",
                "description": "The date from which the current consensus rating is effective, in ISO 8601 format (UTC).",
                "example": "2021-12-29T06:24:56.038+0000"
              }
            },
            "title": "AnalystRating"
          }
        }
      }
    }
  },
  "postman": {
    "name": "Get Analyst Rating",
    "description": {
      "content": "Retrieves analyst rating for one instrument.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "fundamentals",
        "analysis",
        "ratings",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum. default is US_STOCK",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Streaming Subscribe

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/subscribe-using-post.md>

### Subscribe

Subscribe to real-time market data streaming. This interface allows you to subscribe to various types of market data including quotes, snapshots, and tick data for specified securities.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://us-global-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/streaming/subscribe",
  "method": "post",
  "tags": [
    "Market Data/Streaming"
  ],
  "description": "Subscribe to real-time market data streaming. This interface allows you to subscribe to various types of market data including quotes, snapshots, and tick data for specified securities.",
  "operationId": "subscribeUsingPOST",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "category_symbols",
            "session_id",
            "sub_types"
          ],
          "type": "object",
          "properties": {
            "session_id": {
              "type": "string",
              "description": "The session_id used to create the connection, and the connection must be successfully established.",
              "example": "2d29ea02-8a35-11ec-8356-020017000b7b"
            },
            "category_symbols": {
              "type": "array",
              "description": "List of security symbols by category; maximum 100 symbols per query.",
              "items": {
                "required": [
                  "category",
                  "symbols"
                ],
                "type": "object",
                "properties": {
                  "category": {
                    "type": "string",
                    "description": "Security type. Category values are as shown in the enum.",
                    "example": "US_STOCK",
                    "enum": [
                      "US_STOCK"
                    ]
                  },
                  "symbols": {
                    "type": "array",
                    "description": "List of security symbols, supports JSON array format.",
                    "example": [
                      "AAPL",
                      "GOOG"
                    ],
                    "items": {
                      "type": "string",
                      "description": "List of security symbols, supports JSON array format.",
                      "example": "[\"AAPL\",\"GOOG\"]"
                    }
                  }
                },
                "description": "List of security symbols by security type",
                "title": "MultiCategory"
              }
            },
            "sub_types": {
              "type": "string",
              "description": "Subscription data type(s), multiple types separated by commas \",\", enum, refer to: SubType, e.g.: [SNAPSHOT]",
              "example": "[SNAPSHOT]",
              "enum": [
                "QUOTE",
                "SNAPSHOT",
                "TICK"
              ]
            },
            "depth": {
              "type": "string",
              "description": "LV2 subscription depth, default 10 levels, US stocks max 50 levels.",
              "example": "10"
            },
            "overnight_required": {
              "type": "boolean",
              "description": "Whether to include overnight session, true/false. For US stock subscriptions, includes overnight session, only effective for US stocks, default is not included.",
              "example": false
            }
          },
          "title": "SubscribeParams"
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "OK"
    }
  },
  "jsonRequestBodyExample": {
    "session_id": "2d29ea02-8a35-11ec-8356-020017000b7b",
    "category_symbols": [
      {
        "category": "US_STOCK",
        "symbols": [
          "AAPL",
          "GOOG"
        ]
      }
    ],
    "sub_types": "[SNAPSHOT]",
    "depth": "10",
    "overnight_required": false
  },
  "postman": {
    "name": "Subscribe",
    "description": {
      "content": "Subscribe to real-time market data streaming. This interface allows you to subscribe to various types of market data including quotes, snapshots, and tick data for specified securities.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "streaming",
        "subscribe"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```

## Streaming Unsubscribe

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/unsubscribe-using-post.md>

### Unsubscribe

After successfully establishing the market data streaming MQTT connection, call this interface to unsubscribe from real-time market data push. Successful call returns no value; failures return an Error. Unsubscribing will release the topic quota. Frequency limit: 1 call per second per App Key.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull Open API Reference",
    "description": "application.yml\\ncom\\ni18n\\nMETA-INF\\nstatic\\n\\r\\n",
    "contact": {
      "name": "",
      "url": "",
      "email": ""
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://us-global-openapi.uat.webullbroker.com"
    }
  ],
  "path": "/market-data/streaming/unsubscribe",
  "method": "post",
  "tags": [
    "Market Data/Streaming"
  ],
  "description": "After successfully establishing the market data streaming MQTT connection, call this interface to unsubscribe from real-time market data push. Successful call returns no value; failures return an Error. Unsubscribing will release the topic quota. Frequency limit: 1 call per second per App Key.",
  "operationId": "unsubscribeUsingPOST",
  "parameters": [
    {
      "name": "access_token",
      "in": "header",
      "description": "User's authenticated token.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "reqid",
      "in": "header",
      "description": "The unique ID for this request. Suggest using UUID.",
      "schema": {
        "type": "string"
      },
      "example": "2e46d5a4-bef9-4507-8cda-98f85d2f770c"
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "category_symbols",
            "sub_types"
          ],
          "type": "object",
          "properties": {
            "session_id": {
              "type": "string",
              "description": "The session_id used to create the connection, and the connection must be successfully established.",
              "example": "2d29ea02-8a35-11ec-8356-020017000b7b"
            },
            "category_symbols": {
              "type": "array",
              "description": "List of security symbols by category; maximum 100 symbols per query.",
              "items": {
                "required": [
                  "category",
                  "symbols"
                ],
                "type": "object",
                "properties": {
                  "category": {
                    "type": "string",
                    "description": "Security type. Category values are as shown in the enum.",
                    "example": "US_STOCK",
                    "enum": [
                      "US_STOCK"
                    ]
                  },
                  "symbols": {
                    "type": "array",
                    "description": "List of security symbols, supports JSON array format.",
                    "example": [
                      "AAPL",
                      "GOOG"
                    ],
                    "items": {
                      "type": "string",
                      "description": "List of security symbols, supports JSON array format.",
                      "example": "[\"AAPL\",\"GOOG\"]"
                    }
                  }
                },
                "description": "List of security symbols by security type",
                "title": "MultiCategory"
              }
            },
            "sub_types": {
              "type": "string",
              "description": "Subscription data type(s), multiple types separated by commas \",\", enum, refer to: SubType, e.g.: [SNAPSHOT]",
              "example": "[SNAPSHOT]",
              "enum": [
                "QUOTE",
                "SNAPSHOT",
                "TICK"
              ]
            },
            "unsubscribe_all": {
              "type": "boolean",
              "description": "Whether to unsubscribe all, true/false. When set to true, all subscriptions will be cancelled.",
              "example": false
            }
          },
          "title": "UnSubscribeParams"
        }
      }
    },
    "required": true
  },
  "responses": {
    "200": {
      "description": "OK"
    }
  },
  "jsonRequestBodyExample": {
    "session_id": "2d29ea02-8a35-11ec-8356-020017000b7b",
    "category_symbols": [
      {
        "category": "US_STOCK",
        "symbols": [
          "AAPL",
          "GOOG"
        ]
      }
    ],
    "sub_types": "[SNAPSHOT]",
    "unsubscribe_all": false
  },
  "postman": {
    "name": "Unsubscribe",
    "description": {
      "content": "After successfully establishing the market data streaming MQTT connection, call this interface to unsubscribe from real-time market data push. Successful call returns no value; failures return an Error. Unsubscribing will release the topic quota. Frequency limit: 1 call per second per App Key.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "streaming",
        "unsubscribe"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) User's authenticated token.",
          "type": "text/plain"
        },
        "key": "access_token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "The unique ID for this request. Suggest using UUID.",
          "type": "text/plain"
        },
        "key": "reqid",
        "value": ""
      },
      {
        "key": "Content-Type",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "raw",
      "raw": "",
      "options": {
        "raw": {
          "language": "json"
        }
      }
    }
  }
}
```


# Market Data — Stock — Verbatim Reference

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> HTTP on-demand stock/ETF market data (Non-Display Solution).

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Stock Snapshot

> Source: <https://developer.webull.hk/apis/docs/reference/snapshot.md>

### List Stock Snapshots

Retrieves stock real-time snapshot data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/snapshots/list",
  "method": "get",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves stock real-time snapshot data.",
  "operationId": "snapshot",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK",
          "US_ETF",
          "HK_STOCK",
          "CN_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "extend_hour_required",
      "in": "query",
      "description": "Whether to include pre-market and after-hours trading data.",
      "required": false,
      "schema": {
        "type": "string",
        "default": "false"
      },
      "example": false
    },
    {
      "name": "overnight_required",
      "in": "query",
      "description": "Whether to include overnight trading data.",
      "required": false,
      "schema": {
        "type": "string",
        "default": "false"
      },
      "example": false
    },
    {
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
    }
  ],
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "array",
            "items": {
              "required": [
                "symbol"
              ],
              "type": "object",
              "properties": {
                "instrument_id": {
                  "type": "string",
                  "description": "Instrument ID",
                  "example": "913256135"
                },
                "pre_close": {
                  "type": "string",
                  "description": "Previous close price",
                  "example": "101"
                },
                "change_ratio": {
                  "type": "string",
                  "description": "Change ratio",
                  "example": "0.05"
                },
                "symbol": {
                  "type": "string",
                  "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                  "example": "AAPL"
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
                "close": {
                  "type": "string",
                  "description": "Intraday close price",
                  "example": "101"
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
                "ask": {
                  "type": "string",
                  "description": "Ask",
                  "example": "13.9"
                },
                "ask_size": {
                  "type": "string",
                  "description": "Ask size (Quantity)",
                  "example": "5"
                },
                "bid": {
                  "type": "string",
                  "description": "Bid",
                  "example": "13.9"
                },
                "bid_size": {
                  "type": "string",
                  "description": "Bid size (Quantity)",
                  "example": "5"
                },
                "turnover": {
                  "type": "string",
                  "description": "Turnover rate.",
                  "example": "0.01"
                },
                "eps": {
                  "type": "string",
                  "description": "Earnings Per Share.",
                  "example": "7.465"
                },
                "eps_ttm": {
                  "type": "string",
                  "description": "Earnings Per Share (TTM).",
                  "example": "7.465"
                },
                "lot_size": {
                  "type": "string",
                  "description": "Shares per lot.",
                  "example": "1"
                },
                "bps": {
                  "type": "string",
                  "description": "Book Value Per Share.",
                  "example": "4.991"
                },
                "extend_hour_last_price": {
                  "type": "string",
                  "description": "Pre/post market latest price",
                  "example": "100.5"
                },
                "extend_hour_high": {
                  "type": "string",
                  "description": "Pre/post market high price",
                  "example": "101.0"
                },
                "extend_hour_low": {
                  "type": "string",
                  "description": "Pre/post market low price",
                  "example": "99.5"
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
                  "description": "Overnight ask",
                  "example": "13.9"
                },
                "ovn_ask_size": {
                  "type": "string",
                  "description": "Overnight ask size (Quantity)",
                  "example": "5"
                },
                "ovn_bid": {
                  "type": "string",
                  "description": "Overnight bid",
                  "example": "13.9"
                },
                "ovn_bid_size": {
                  "type": "string",
                  "description": "Overnight bid size (Quantity)",
                  "example": "5"
                }
              },
              "description": "Market snapshot data response object",
              "title": "SnapshotVo"
            }
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock Snapshots",
    "description": {
      "content": "Retrieves stock real-time snapshot data.",
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
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Whether to include pre-market and after-hours trading data.",
            "type": "text/plain"
          },
          "key": "extend_hour_required",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Whether to include overnight trading data.",
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
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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

## Stock Tick

> Source: <https://developer.webull.hk/apis/docs/reference/tick.md>

### List Stock Ticks

Retrieves stock tick-by-tick trade data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/ticks/list",
  "method": "get",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves stock tick-by-tick trade data.",
  "operationId": "tick",
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
      "description": "Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK",
          "US_ETF",
          "HK_STOCK",
          "CN_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of ticks, default 30, maximum limit 1000.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "30"
      },
      "example": 30
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
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
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
              "symbol"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "AAPL"
              },
              "instrument_id": {
                "type": "string",
                "description": "Instrument ID",
                "example": "913256409"
              },
              "result": {
                "type": "array",
                "description": "Tick details",
                "items": {
                  "required": [
                    "price",
                    "side",
                    "time",
                    "volume"
                  ],
                  "type": "object",
                  "properties": {
                    "time": {
                      "type": "string",
                      "description": "Trade time of this tick, expressed as Unix epoch timestamp in milliseconds",
                      "example": "1761182953043"
                    },
                    "price": {
                      "type": "string",
                      "description": "Executed trade price for this futures contract at this tick",
                      "example": "48.07"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Executed trade volume at this tick, expressed in number of futures contracts",
                      "example": "1"
                    },
                    "side": {
                      "type": "string",
                      "description": "Such as: B S G L N",
                      "example": "S"
                    }
                  },
                  "description": "Single futures trade tick detail, representing one executed trade.",
                  "title": "Tick"
                }
              }
            },
            "description": "TickVo",
            "title": "TickVo"
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock Ticks",
    "description": {
      "content": "Retrieves stock tick-by-tick trade data.",
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
            "content": "(Required) Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Number of ticks, default 30, maximum limit 1000.",
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
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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

## Stock Quotes (Depth)

> Source: <https://developer.webull.hk/apis/docs/reference/quotes.md>

### List Stock Depths

Retrieves stock quotes data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/depths/list",
  "method": "get",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves stock quotes data.",
  "operationId": "quotes",
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
      "description": "Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK",
          "US_ETF",
          "HK_STOCK",
          "CN_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "depth",
      "in": "query",
      "description": "Market depth, L1-1 level, L2-default 10 levels, etc.",
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
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
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
              "instrument_id",
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
              "instrument_id": {
                "type": "string",
                "description": "Instrument ID",
                "example": "913255275"
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
                    },
                    "broker": {
                      "type": "array",
                      "items": {
                        "required": [
                          "bid",
                          "name"
                        ],
                        "type": "object",
                        "properties": {
                          "bid": {
                            "type": "string",
                            "description": "Broker ID",
                            "example": "1"
                          },
                          "name": {
                            "type": "string",
                            "description": "Broker Name",
                            "example": "2"
                          }
                        },
                        "description": "BrokerItem",
                        "title": "BrokerItem"
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
                    },
                    "broker": {
                      "type": "array",
                      "items": {
                        "required": [
                          "bid",
                          "name"
                        ],
                        "type": "object",
                        "properties": {
                          "bid": {
                            "type": "string",
                            "description": "Broker ID",
                            "example": "1"
                          },
                          "name": {
                            "type": "string",
                            "description": "Broker Name",
                            "example": "2"
                          }
                        },
                        "description": "BrokerItem",
                        "title": "BrokerItem"
                      }
                    }
                  },
                  "description": "BidItem",
                  "title": "BidItem"
                }
              }
            },
            "description": "QuoteVo",
            "title": "QuoteVo"
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock Depths",
    "description": {
      "content": "Retrieves stock quotes data.",
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
            "content": "(Required) Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Market depth, L1-1 level, L2-default 10 levels, etc.",
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
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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

## Stock Historical Bars (Batch)

> Source: <https://developer.webull.hk/apis/docs/reference/historical-bars.md>

### List Stock Historical Bars

Retrieves historical bars data for multiple stock symbols in batch.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/bars/list",
  "method": "post",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves historical bars data for multiple stock symbols in batch.",
  "operationId": "historicalBars",
  "parameters": [
    {
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "category",
            "symbols",
            "timespan"
          ],
          "type": "object",
          "properties": {
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
            },
            "category": {
              "type": "string",
              "description": "Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported.",
              "example": "US_STOCK",
              "enum": [
                "US_STOCK",
                "US_ETF",
                "HK_STOCK",
                "CN_STOCK"
              ]
            },
            "timespan": {
              "type": "string",
              "description": "Bar time granularity.",
              "example": "D",
              "enum": [
                "S5",
                "S15",
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
            "count": {
              "type": "integer",
              "description": "Number of bars, default 200, maximum limit 1200 (M1 supports up to 1650).",
              "format": "int32",
              "example": 500
            },
            "real_time_required": {
              "type": "boolean",
              "description": "Return the latest trading data, default is true;<br/> true: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request.<br/> false: The returned data includes the latest market data.",
              "example": false
            },
            "trading_sessions": {
              "type": "string",
              "description": "Specify trading session(s). Multiple sessions separated by \",\".",
              "example": "PRE,RTH",
              "enum": [
                "PRE",
                "RTH",
                "ATH",
                "OVN"
              ]
            },
            "start_time": {
              "type": "integer",
              "description": "Start time as timestamp in milliseconds. Used to specify the beginning of the time range for bar data.",
              "format": "int64",
              "example": 1711262998500
            },
            "end_time": {
              "type": "integer",
              "description": "End time as timestamp in milliseconds. Used to specify the end of the time range for bar data. Delayed permission will automatically offset the time.",
              "format": "int64",
              "example": 1711349398500
            }
          },
          "title": "BatchBarRequest"
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
              "result"
            ],
            "type": "object",
            "properties": {
              "result": {
                "type": "array",
                "description": "List of batch bar data results, each element contains historical bar data for one stock.",
                "items": {
                  "required": [
                    "instrument_id",
                    "result",
                    "symbol"
                  ],
                  "type": "object",
                  "properties": {
                    "symbol": {
                      "type": "string",
                      "description": "Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange.",
                      "example": "SILZ5"
                    },
                    "instrument_id": {
                      "type": "string",
                      "description": "Unique instrument identifier for this futures contract in the Webull system or exchange.",
                      "example": "470059643"
                    },
                    "result": {
                      "type": "array",
                      "description": "List of historical bar data for this futures.",
                      "items": {
                        "required": [
                          "close",
                          "high",
                          "low",
                          "open",
                          "time",
                          "volume"
                        ],
                        "type": "object",
                        "properties": {
                          "time": {
                            "type": "string",
                            "description": "Bar UTC time",
                            "example": "2021-12-28T09:00:09.945+0000"
                          },
                          "open": {
                            "type": "string",
                            "description": "Open price",
                            "example": "1.3362"
                          },
                          "close": {
                            "type": "string",
                            "description": "Close price",
                            "example": "1.3362"
                          },
                          "high": {
                            "type": "string",
                            "description": "High price",
                            "example": "1.3362"
                          },
                          "low": {
                            "type": "string",
                            "description": "Low price",
                            "example": "1.3362"
                          },
                          "volume": {
                            "type": "string",
                            "description": "Volume",
                            "example": "10"
                          }
                        },
                        "description": "Single K-line price record",
                        "title": "PriceRecord"
                      }
                    }
                  },
                  "description": "Historical K-line data for a single futures",
                  "title": "SymbolData"
                }
              }
            },
            "description": "Batch historical K-line data response object",
            "title": "BatchBarsVo"
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": {
    "symbols": [
      "AAPL",
      "GOOG"
    ],
    "category": "US_STOCK",
    "timespan": "D",
    "count": 500,
    "real_time_required": false,
    "trading_sessions": "PRE,RTH",
    "start_time": 1711262998500,
    "end_time": 1711349398500
  },
  "postman": {
    "name": "List Stock Historical Bars",
    "description": {
      "content": "Retrieves historical bars data for multiple stock symbols in batch.",
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
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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

## Stock Footprint

> Source: <https://developer.webull.hk/apis/docs/reference/footprint.md>

### List Stock Footprints

Retrieves stock footprint data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/footprints/list",
  "method": "get",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves stock footprint data.",
  "operationId": "footprint",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 20 symbols per query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum; Only US_STOCK type queries are supported.",
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
      "name": "timespan",
      "in": "query",
      "description": "Supports granularities such as S5, S15, M1, M5, and M30.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "S5",
          "S15",
          "M1",
          "M5",
          "M30"
        ]
      },
      "example": "S5"
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of bars, default 200, maximum limit 1200.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "1-1200",
        "default": "200"
      },
      "example": 500
    },
    {
      "name": "real_time_required",
      "in": "query",
      "description": "Does it include the latest data? For candlesticks that are not yet finalized, the default is false (does not include). Only minute timespan is used.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "false"
      },
      "example": true
    },
    {
      "name": "trading_sessions",
      "in": "query",
      "description": "Specify trading hours. OVN type not supported.",
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
      "example": "RTH"
    },
    {
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
    }
  ],
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
                "symbol": {
                  "type": "string",
                  "description": "Security symbol",
                  "example": "AAPL"
                },
                "instrument_id": {
                  "type": "string",
                  "description": "Unique Identifier for Securities",
                  "example": "913256135"
                },
                "result": {
                  "type": "array",
                  "description": "Footprint chart candlestick chart",
                  "items": {
                    "type": "object",
                    "properties": {
                      "time": {
                        "type": "string",
                        "description": "Transaction date",
                        "example": "2025-09-30T05:47:00.000+0000"
                      },
                      "trading_session": {
                        "type": "string",
                        "description": "Trading Hours",
                        "example": "RTH"
                      },
                      "total": {
                        "type": "string",
                        "description": "The sum of the main buy and sell volumes",
                        "example": "1000"
                      },
                      "delta": {
                        "type": "string",
                        "description": "The difference in trading volume (primary buyers - primary sellers)",
                        "example": "200"
                      },
                      "buy_total": {
                        "type": "string",
                        "description": "Buy-initiated volume",
                        "example": "600"
                      },
                      "sell_total": {
                        "type": "string",
                        "description": "Sell-initiated volume",
                        "example": "400"
                      },
                      "buy_detail": {
                        "type": "object",
                        "additionalProperties": {
                          "type": "string",
                          "description": "The main purchase footprint details (quantity combined for items with the same price).",
                          "example": "{\"24.20\":\"100\",\"24.21\":\"60\"}"
                        },
                        "description": "The main purchase footprint details (quantity combined for items with the same price).",
                        "example": {
                          "24.20": "100",
                          "24.21": "60"
                        }
                      },
                      "sell_detail": {
                        "type": "object",
                        "additionalProperties": {
                          "type": "string",
                          "description": "The main seller's footprint shows details (quantities combined for the same price).",
                          "example": "{\"24.20\":\"50\",\"24.21\":\"50\"}"
                        },
                        "description": "The main seller's footprint shows details (quantities combined for the same price).",
                        "example": {
                          "24.20": "50",
                          "24.21": "50"
                        }
                      }
                    },
                    "description": "Footprint",
                    "title": "Footprint"
                  }
                }
              },
              "description": "Stock Footprint",
              "title": "FootprintVo"
            }
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock Footprints",
    "description": {
      "content": "Retrieves stock footprint data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "footprints",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 20 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Category values are as shown in the enum; Only US_STOCK type queries are supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Supports granularities such as S5, S15, M1, M5, and M30.",
            "type": "text/plain"
          },
          "key": "timespan",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of bars, default 200, maximum limit 1200.",
            "type": "text/plain"
          },
          "key": "count",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Does it include the latest data? For candlesticks that are not yet finalized, the default is false (does not include). Only minute timespan is used.",
            "type": "text/plain"
          },
          "key": "real_time_required",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Specify trading hours. OVN type not supported.",
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
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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

## NOII Bars

> Source: <https://developer.webull.hk/apis/docs/reference/get-noii-bars.md>

### List Stock NOII Bars

Retrieves NOII bars data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/noii-bars/list",
  "method": "get",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves NOII bars data.",
  "operationId": "getNoiiBars",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol. Currently only supports single symbol query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only supports US_STOCK.",
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
      "name": "imbalance_action_type",
      "in": "query",
      "description": "Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance).",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Imbalance action type for NOII data.",
        "enum": [
          "PRE_OPEN",
          "PRE_CLOSE"
        ]
      },
      "example": "PRE_OPEN"
    },
    {
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
    }
  ],
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
                "instrument_id": {
                  "type": "string",
                  "description": "Instrument unique identifier",
                  "example": "913256135"
                },
                "symbol": {
                  "type": "string",
                  "description": "Security symbol",
                  "example": "AAPL"
                },
                "imbalance_time": {
                  "type": "integer",
                  "description": "Timestamp of the imbalance data in milliseconds (data publish time)",
                  "format": "int64",
                  "example": 1711262998500
                },
                "imbalance_ref_price": {
                  "type": "string",
                  "description": "Reference price",
                  "example": "172.35"
                },
                "imbalance_near_price": {
                  "type": "string",
                  "description": "Indicative Match Price - the most likely execution price",
                  "example": "173.1"
                },
                "imbalance_far_price": {
                  "type": "string",
                  "description": "Far Price - the price at which orders could execute in extreme scenarios",
                  "example": "175.5"
                },
                "imbalance_action_type": {
                  "type": "string",
                  "description": "Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance)",
                  "example": "PRE_OPEN"
                }
              },
              "description": "NOII (Net Order Imbalance Indicator) bar data.",
              "title": "NoiiBarVo"
            }
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock NOII Bars",
    "description": {
      "content": "Retrieves NOII bars data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "noii-bars",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol. Currently only supports single symbol query.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only supports US_STOCK.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance).",
            "type": "text/plain"
          },
          "key": "imbalance_action_type",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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

## NOII Snapshot

> Source: <https://developer.webull.hk/apis/docs/reference/get-noii-snapshot.md>

### List Stock NOII Snapshots

Retrieves NOII snapshot data.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.hk/en/help",
      "email": "webull-api-support@webull.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.hk"
    }
  ],
  "path": "/market-data/stocks/noii-snapshots/list",
  "method": "get",
  "tags": [
    "Stock Market Data"
  ],
  "description": "Retrieves NOII snapshot data.",
  "operationId": "getNoiiSnapshot",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Security symbol. Currently only supports single symbol query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only supports US_STOCK.",
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
      "name": "imbalance_action_type",
      "in": "query",
      "description": "Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance).",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Imbalance action type for NOII data.",
        "enum": [
          "PRE_OPEN",
          "PRE_CLOSE"
        ]
      },
      "example": "PRE_OPEN"
    },
    {
      "name": "x-app-key",
      "in": "header",
      "description": "A unique identifier issued to a developer for accessing an application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-app-secret",
      "in": "header",
      "description": "A unique key issued to developers to access the application's API.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-timestamp",
      "in": "header",
      "description": "Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-signature-version",
      "in": "header",
      "description": "Signature algorithm version, default is 1.0.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "1.0"
      },
      "examples": {
        "1.0": {
          "value": "1.0"
        }
      }
    },
    {
      "name": "x-signature-algorithm",
      "in": "header",
      "description": "Signature algorithm, default is HMAC-SHA1.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "HMAC-SHA1"
      },
      "examples": {
        "HMAC-SHA1": {
          "value": "HMAC-SHA1"
        }
      }
    },
    {
      "name": "x-signature-nonce",
      "in": "header",
      "description": "Signature unique random number.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-access-token",
      "in": "header",
      "description": "An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
      "required": true,
      "schema": {
        "type": "string"
      }
    },
    {
      "name": "x-version",
      "in": "header",
      "description": "API interface version. Supported values: `v2`, `v3`.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v3"
      },
      "examples": {
        "v3": {
          "value": "v3"
        }
      }
    },
    {
      "name": "x-signature",
      "in": "header",
      "description": "A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
      "required": true,
      "schema": {
        "type": "string"
      }
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
              "instrument_id": {
                "type": "string",
                "description": "Instrument unique identifier",
                "example": "913256135"
              },
              "symbol": {
                "type": "string",
                "description": "Security symbol",
                "example": "AAPL"
              },
              "paired_shares": {
                "type": "string",
                "description": "Paired shares - the number of shares that can be matched under current conditions",
                "example": "701859"
              },
              "imbalance_shares": {
                "type": "string",
                "description": "Imbalance shares - the number of unmatched buy/sell shares",
                "example": "5715"
              },
              "imbalance_side": {
                "type": "string",
                "description": "Imbalance side (direction of imbalance)",
                "example": "2"
              },
              "imbalance_ref_price": {
                "type": "string",
                "description": "Reference price",
                "example": "253.83"
              },
              "imbalance_near_price": {
                "type": "string",
                "description": "Indicative Match Price - the most likely execution price",
                "example": "253.93"
              },
              "imbalance_far_price": {
                "type": "string",
                "description": "Far Price - the price at which orders could execute in extreme scenarios",
                "example": "253.98"
              },
              "imbalance_action_type": {
                "type": "string",
                "description": "Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance)",
                "example": "PRE_OPEN"
              },
              "imbalance_time": {
                "type": "integer",
                "description": "Timestamp in milliseconds",
                "format": "int64",
                "example": 1774272599000
              },
              "imbalance_var_indicator": {
                "type": "string",
                "description": "Volatility/imbalance status indicator",
                "example": "10"
              }
            },
            "description": "NOII (Net Order Imbalance Indicator) snapshot data.",
            "title": "NoiiSnapshotVo"
          }
        }
      }
    },
    "401": {
      "description": "Unauthorized: Authentication required",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "UNAUTHORIZED"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Insufficient permission"
              }
            }
          }
        }
      }
    },
    "417": {
      "description": "A business logic error triggered when the request cannot be processed due to domain-specific constraints.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "INVALID_PARAMETER"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Parameter error, phone"
              }
            }
          }
        }
      }
    },
    "500": {
      "description": "Internal Server Error.",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "Internal logic error code",
                "example": "SYSTEM_ERROR"
              },
              "message": {
                "type": "string",
                "description": "Error message",
                "example": "Internal Server Error"
              }
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Stock NOII Snapshots",
    "description": {
      "content": "Retrieves NOII snapshot data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "stocks",
        "noii-snapshots",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security symbol. Currently only supports single symbol query.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only supports US_STOCK.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance).",
            "type": "text/plain"
          },
          "key": "imbalance_action_type",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique identifier issued to a developer for accessing an application's API.",
          "type": "text/plain"
        },
        "key": "x-app-key",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A unique key issued to developers to access the application's API.",
          "type": "text/plain"
        },
        "key": "x-app-secret",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Timestamp of the request, follows ISO8601 format: YYYY-MM-DDThh:mm:ssZ, e.g. 2023-07-16T19:23:51Z, only supports UTC time zone.",
          "type": "text/plain"
        },
        "key": "x-timestamp",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm version, default is 1.0.",
          "type": "text/plain"
        },
        "key": "x-signature-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature algorithm, default is HMAC-SHA1.",
          "type": "text/plain"
        },
        "key": "x-signature-algorithm",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) Signature unique random number.",
          "type": "text/plain"
        },
        "key": "x-signature-nonce",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) An access token is a credential that represents the authorization granted to a client (e.g., a user or an application) to access specific protected resources on behalf of a user, without needing to share their password.",
          "type": "text/plain"
        },
        "key": "x-access-token",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) API interface version. Supported values: `v2`, `v3`.",
          "type": "text/plain"
        },
        "key": "x-version",
        "value": ""
      },
      {
        "disabled": false,
        "description": {
          "content": "(Required) A signature is a unique digital fingerprint, typically encrypted, that verifies the authenticity and integrity of a message or transaction, ensuring it has not been tampered with during transmission.",
          "type": "text/plain"
        },
        "key": "x-signature",
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


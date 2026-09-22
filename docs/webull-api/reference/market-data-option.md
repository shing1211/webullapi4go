# Market Data — Option — Verbatim Reference

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> Option tick, snapshot and historical bars (Non-Display Solution), plus the option contract list used to build an option chain.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Option Tick

> Source: <https://developer.webull.hk/apis/docs/reference/option-tick.md>

### List Option Ticks

Retrieves option tick-by-tick trade data.

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
  "path": "/market-data/options/ticks/list",
  "method": "get",
  "tags": [
    "Option Market Data"
  ],
  "description": "Retrieves option tick-by-tick trade data.",
  "operationId": "optionTick",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Option symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL260522C00300000"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_OPTION` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_OPTION"
        ]
      },
      "example": "US_OPTION"
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of ticks, maximum limit 1200 .",
      "required": true,
      "schema": {
        "type": "string",
        "default": "30"
      },
      "example": 30
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
                "description": "Option contract symbol used in trading and market data, e.g. front-month code provided by the exchange.",
                "example": "AAPL260522C00300000"
              },
              "instrument_id": {
                "type": "string",
                "description": "Unique instrument identifier for this option contract in the Webull system or exchange.",
                "example": "470059643"
              },
              "result": {
                "type": "array",
                "description": "List of tick details (trade prints) for this option contract.",
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
            "title": "OptionTickVo"
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
    "name": "List Option Ticks",
    "description": {
      "content": "Retrieves option tick-by-tick trade data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "options",
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
            "content": "(Required) Option symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_OPTION` is supported for this interface.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Number of ticks, maximum limit 1200 .",
            "type": "text/plain"
          },
          "key": "count",
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

## Option Snapshot

> Source: <https://developer.webull.hk/apis/docs/reference/option-snapshot.md>

### List Option Snapshots

Retrieves option real-time snapshot data.

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
  "path": "/market-data/options/snapshots/list",
  "method": "get",
  "tags": [
    "Option Market Data"
  ],
  "description": "Retrieves option real-time snapshot data.",
  "operationId": "optionSnapshot",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of option symbols, separated by commas; maximum 20 symbols per query. Example: AAPL260522C00300000.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL260522C00300000"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_OPTION` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_OPTION"
        ]
      },
      "example": "US_OPTION"
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
                  "description": "Unique instrument identifier for this option contract in the Webull system or exchange.",
                  "example": "470059643"
                },
                "symbol": {
                  "type": "string",
                  "description": "Option contract symbol used in trading and market data, e.g. front-month code provided by the exchange.",
                  "example": "AAPL260522C00300000"
                },
                "price": {
                  "type": "string",
                  "description": "Last traded price (last done) of the option contract, quoted in the contract's trading currency (e.g. USD).",
                  "example": "47.35"
                },
                "open": {
                  "type": "string",
                  "description": "Session open price for this option contract. Represents the first traded price of the current regular trading session. If no trade has occurred in the session, this field may be empty.",
                  "example": "48.17"
                },
                "high": {
                  "type": "string",
                  "description": "Session high price for this option contract during the current regular trading session. If no trade has occurred in the session, this field may be empty.",
                  "example": "48.655"
                },
                "low": {
                  "type": "string",
                  "description": "Session low price for this option contract during the current regular trading session. If no trade has occurred in the session, this field may be empty.",
                  "example": "46.855"
                },
                "pre_close": {
                  "type": "string",
                  "description": "Previous settlement/close price of the option contract (typically the official settlement price of the previous trading day), quoted in the contract's trading currency.",
                  "example": "47.704"
                },
                "volume": {
                  "type": "string",
                  "description": "Accumulated traded volume for the current session, expressed in number of option contracts. If no trading occurred in the session, this field may be empty.",
                  "example": "48906"
                },
                "change": {
                  "type": "string",
                  "description": "Absolute price change of the last traded price relative to the previous settlement/close price. If no valid reference price or last trade exists, this field may be empty.",
                  "example": "-0.354"
                },
                "change_ratio": {
                  "type": "string",
                  "description": "Price change ratio of the last traded price relative to the previous settlement/close price, Expressed as a decimal (e.g., -0.0074 represents -0.74%)",
                  "example": "-0.0074"
                },
                "last_trade_time": {
                  "type": "integer",
                  "description": "Timestamp of the last executed trade for this option contract, expressed as Unix epoch time in milliseconds",
                  "format": "int64",
                  "example": 1761131406558
                },
                "close": {
                  "type": "string",
                  "description": "Close price of the option contract, quoted in the contract's trading currency.",
                  "example": "0.05"
                },
                "strike_price": {
                  "type": "string",
                  "description": "Strike price of the option contract, the price at which the option can be exercised.",
                  "example": "135.0"
                },
                "gamma": {
                  "type": "string",
                  "description": "Gamma, the rate of change of delta with respect to the underlying asset's price. Measures the convexity of the option's value.",
                  "example": "1.0E-4"
                },
                "delta": {
                  "type": "string",
                  "description": "Delta, the rate of change of the option price with respect to the underlying asset's price. Ranges from -1 to 1 for options.",
                  "example": "-0.0023"
                },
                "rho": {
                  "type": "string",
                  "description": "Rho, the rate of change of the option price with respect to the risk-free interest rate.",
                  "example": "-0.001"
                },
                "theta": {
                  "type": "string",
                  "description": "Theta, the rate of change of the option price with respect to time (time decay). Usually expressed as the change in option price per day.",
                  "example": "-0.0037"
                },
                "vega": {
                  "type": "string",
                  "description": "Vega, the rate of change of the option price with respect to volatility. Measures sensitivity to implied volatility changes.",
                  "example": "0.0077"
                },
                "imp_vol": {
                  "type": "string",
                  "description": "Implied volatility, the market's forecast of the underlying asset's likely movement. Expressed as a decimal (e.g., 0.609 represents 60.9%).",
                  "example": "0.609"
                },
                "open_interest": {
                  "type": "string",
                  "description": "Open interest, representing the total number of outstanding and unsettled option contracts for this instrument, expressed in number of contracts.",
                  "example": "14331"
                },
                "quote_time": {
                  "type": "integer",
                  "description": "Quote timestamp of this snapshot, expressed as Unix epoch time in milliseconds (UTC). Represents the time when this snapshot data was generated.",
                  "format": "int64",
                  "example": 1761131409276
                },
                "bid": {
                  "type": "string",
                  "description": "Best bid price (top of book), i.e. the highest price currently offered by buyers, quoted in the contract's trading currency.",
                  "example": "47.345"
                },
                "ask": {
                  "type": "string",
                  "description": "Best ask price (top of book), i.e. the lowest price currently offered by sellers, quoted in the contract's trading currency.",
                  "example": "47.355"
                },
                "ask_size": {
                  "type": "string",
                  "description": "Best ask size, i.e. the total quantity available at the best ask price, expressed in number of option contracts (whole contract units).",
                  "example": "2"
                },
                "bid_size": {
                  "type": "string",
                  "description": "Best bid size, i.e. the total quantity available at the best bid price, expressed in number of option contracts (whole contract units).",
                  "example": "1"
                },
                "deal_amount": {
                  "type": "string",
                  "description": "Total deal amount (traded value) for the option contract, quoted in the contract's trading currency.",
                  "example": "70267.5"
                }
              },
              "description": "Market snapshot data response object",
              "title": "OptionSnapshotVo"
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
    "name": "List Option Snapshots",
    "description": {
      "content": "Retrieves option real-time snapshot data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "options",
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
            "content": "(Required) List of option symbols, separated by commas; maximum 20 symbols per query. Example: AAPL260522C00300000.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_OPTION` is supported for this interface.",
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

## Option Historical Bars

> Source: <https://developer.webull.hk/apis/docs/reference/option-historical-bars.md>

### List Option Historical Bars

Retrieves option historical bars data.

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
  "path": "/market-data/options/bars/list",
  "method": "get",
  "tags": [
    "Option Market Data"
  ],
  "description": "Retrieves option historical bars data.",
  "operationId": "optionHistoricalBars",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of option symbols, separated by commas; maximum 20 symbols per query. Example: AAPL260522C00300000",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "AAPL260522C00300000"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_OPTION` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_OPTION"
        ]
      },
      "example": "US_OPTION"
    },
    {
      "name": "timespan",
      "in": "query",
      "description": "Bar time granularity. eg: M1, M5, M15, M30, M60, M120, M240, D, W, M, Y",
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
      "example": "M1"
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of bars, maximum limit 1200 .",
      "required": false,
      "schema": {
        "type": "string",
        "default": "200"
      },
      "example": 200
    },
    {
      "name": "real_time_required",
      "in": "query",
      "description": "Include the latest data",
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
            "required": [
              "result"
            ],
            "type": "object",
            "properties": {
              "result": {
                "type": "array",
                "description": "List of batch bar data results, each element contains historical bar data for one option.",
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
                      "description": "Option contract symbol used in trading and market data, e.g. front-month code provided by the exchange.",
                      "example": "AAPL260522C00300000"
                    },
                    "instrument_id": {
                      "type": "string",
                      "description": "Unique instrument identifier for this option contract in the Webull system or exchange.",
                      "example": "470059643"
                    },
                    "result": {
                      "type": "array",
                      "description": "List of historical bar data for this option.",
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
                  "description": "Historical K-line data for a single option",
                  "title": "OptionSymbolData"
                }
              }
            },
            "description": "Batch historical K-line data response object",
            "title": "OptionBatchBarsVo"
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
    "name": "List Option Historical Bars",
    "description": {
      "content": "Retrieves option historical bars data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "options",
        "bars",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) List of option symbols, separated by commas; maximum 20 symbols per query. Example: AAPL260522C00300000",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_OPTION` is supported for this interface.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Bar time granularity. eg: M1, M5, M15, M30, M60, M120, M240, D, W, M, Y",
            "type": "text/plain"
          },
          "key": "timespan",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of bars, maximum limit 1200 .",
            "type": "text/plain"
          },
          "key": "count",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Include the latest data",
            "type": "text/plain"
          },
          "key": "real_time_required",
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

## Option Contract List (Chain)

> Source: <https://developer.webull.hk/apis/docs/reference/option-contract-list.md>

### List Option Contracts

Retrieves option contracts filtered by underlying symbol, status and other attributes. Contains static contract information.

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
  "path": "/trading/instruments/options/contracts/list",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves option contracts filtered by underlying symbol, status and other attributes. Contains static contract information.",
  "operationId": "optionContractList",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Option category. Currently only US_OPTION is supported.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Option Category<br/>US_OPTION - US option<br/>",
        "enum": [
          "US_OPTION"
        ]
      },
      "example": "US_OPTION"
    },
    {
      "name": "option_symbols",
      "in": "query",
      "description": "Option symbols, multiple separated by commas.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "AAPL250620C00150000,SPX241220P04200000"
    },
    {
      "name": "underlying_symbols",
      "in": "query",
      "description": "Underlying symbols, multiple separated by commas.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "AAPL,MSFT"
    },
    {
      "name": "status",
      "in": "query",
      "description": "Contract status, default LISTING. Enum: LISTING, DELISTING.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "LISTING"
    },
    {
      "name": "start_date",
      "in": "query",
      "description": "Exact expiration date, format: YYYY-MM-DD.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2025-06-20"
    },
    {
      "name": "end_date",
      "in": "query",
      "description": "Expiration date lower bound (inclusive), format: YYYY-MM-DD.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2025-06-20"
    },
    {
      "name": "root_symbol",
      "in": "query",
      "description": "Root symbol filter (series symbol, e.g. SPXW). Mainly used for index options and post-CA non-standard contracts.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "SPXW"
    },
    {
      "name": "option_type",
      "in": "query",
      "description": "Contract type: CALL / PUT.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "CALL"
    },
    {
      "name": "style",
      "in": "query",
      "description": "Exercise style: AMERICAN / EUROPEAN.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "AMERICAN"
    },
    {
      "name": "strike_price_gte",
      "in": "query",
      "description": "Strike price lower bound (inclusive).",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": 150
    },
    {
      "name": "strike_price_lte",
      "in": "query",
      "description": "Strike price upper bound (inclusive).",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": 200
    },
    {
      "name": "ppind",
      "in": "query",
      "description": "Penny Program Indicator: true = Penny Pilot contract, false = non-Penny Pilot.",
      "required": false,
      "schema": {
        "type": "boolean"
      },
      "example": true
    },
    {
      "name": "show_deliverables",
      "in": "query",
      "description": "Whether to return deliverables array in response: TRUE / FALSE, default FALSE.",
      "required": false,
      "schema": {
        "type": "boolean"
      },
      "example": "FALSE"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page. Not required for first request.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiIxMDMyNTExMTY0IiwicGFnZUluZGV4IjowLCJwYWdlU2l6ZSI6MTAwMCwiY29uZGl0aW9uIjoiVVNfT1BUSU9OO251bGw7bnVsbDtMSVNUSU5HO251bGw7bnVsbDtudWxsO251bGw7bnVsbDtudWxsO251bGw7bnVsbDtmYWxzZTsifQ=="
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
              "data": {
                "type": "array",
                "description": "Result data list",
                "items": {
                  "type": "object",
                  "properties": {
                    "instrument_id": {
                      "type": "string",
                      "description": "Contract unique identifier",
                      "example": "1038392011"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "OCC contract symbol",
                      "example": "AAPL250620C00150000"
                    },
                    "status": {
                      "type": "string",
                      "description": "Contract status: LISTING, DELISTING",
                      "example": "LISTING"
                    },
                    "tradable_status": {
                      "type": "string",
                      "description": "Trading restriction: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
                      "example": "OC"
                    },
                    "expiration_date": {
                      "type": "string",
                      "description": "Expiration date (effective expiration date after CA events), format: YYYY-MM-DD",
                      "example": "2025-06-20"
                    },
                    "root_symbol": {
                      "type": "string",
                      "description": "Root symbol (series symbol). For most equity options, same as underlying_symbol; may differ for index and post-CA contracts (e.g., SPXW)",
                      "example": "AAPL"
                    },
                    "underlying_symbol": {
                      "type": "string",
                      "description": "Underlying symbol",
                      "example": "AAPL"
                    },
                    "underlying_instrument_id": {
                      "type": "string",
                      "description": "Underlying instrument id",
                      "example": "913256135"
                    },
                    "underlying_type": {
                      "type": "string",
                      "description": "Underlying type: EQUITY_PUT_OPTION / EQUITY_CALL_OPTION / INDEX_CALL_OPTION / ETF_CALL_OPTION / ETF_PUT_OPTION",
                      "example": "EQUITY_CALL_OPTION"
                    },
                    "option_type": {
                      "type": "string",
                      "description": "Contract type: CALL / PUT",
                      "example": "CALL"
                    },
                    "style": {
                      "type": "string",
                      "description": "Exercise style: AMERICAN / EUROPEAN",
                      "example": "AMERICAN"
                    },
                    "strike_price": {
                      "type": "string",
                      "description": "Strike price",
                      "example": "150.0"
                    },
                    "multiplier": {
                      "type": "string",
                      "description": "Contract multiplier, typically 100 for US equity options",
                      "example": "100"
                    },
                    "settlement_method": {
                      "type": "string",
                      "description": "Settlement method: PHYSICAL / CASH",
                      "example": "PHYSICAL"
                    },
                    "expired_cycle": {
                      "type": "string",
                      "description": "Expiration cycle: DAILY / WEEKLY / MONTHLY / QUARTERLY / EOM",
                      "example": "MONTHLY"
                    },
                    "ppind": {
                      "type": "boolean",
                      "description": "Penny Program Indicator: true = Penny Pilot, false = non-Penny Pilot",
                      "example": true
                    },
                    "currency": {
                      "type": "string",
                      "description": "Pricing currency",
                      "example": "USD"
                    },
                    "def_type": {
                      "type": "string",
                      "description": "Definition type: STANDARD / BINARY / FLEX",
                      "example": "STANDARD"
                    },
                    "listed_exchanges": {
                      "type": "array",
                      "description": "Listed exchanges",
                      "items": {
                        "type": "string",
                        "description": "Listed exchanges"
                      }
                    },
                    "deliverables": {
                      "type": "array",
                      "description": "Deliverables configuration; after CA a single contract may correspond to multiple underlyings + cash",
                      "items": {
                        "type": "object",
                        "properties": {
                          "asset_type": {
                            "type": "string",
                            "description": "Asset type. settlement_method=CAFX/CFR/CADF returns CASH; settlement_method=BTOB/CCC/POST/PHYS returns EQUITY",
                            "example": "EQUITY"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Deliverable underlying symbol",
                            "example": "AAPL"
                          },
                          "instrument_id": {
                            "type": "string",
                            "description": "Deliverable underlying instrument id",
                            "example": "913256135"
                          },
                          "amount": {
                            "type": "string",
                            "description": "Delivery amount: number of shares when type=EQUITY, cash amount when type=CASH",
                            "example": "0"
                          },
                          "allocation_percentage": {
                            "type": "string",
                            "description": "Allocation percentage (0-100)",
                            "example": "100"
                          },
                          "settlement_type": {
                            "type": "string",
                            "description": "Settlement cycle: T_0 / T_1 / T_2 / T_3 / T_4",
                            "example": "T_1"
                          },
                          "settlement_method": {
                            "type": "string",
                            "description": "Clearing method: PHYSICAL (BTOB) / CASH_DIFF (CADF) / CASH_FIXED (CAFX) / CCC / CASH_FIXED_RETURN (CFR) / POSITIONAL (POST)",
                            "example": "PHYSICAL"
                          },
                          "settlement_status": {
                            "type": "string",
                            "description": "Settlement status: DELAYED / REGULAR",
                            "example": "REGULAR"
                          }
                        },
                        "description": "Option Contract Deliverable",
                        "title": "OptionContractDeliverable"
                      }
                    }
                  },
                  "description": "Option Contract Information",
                  "title": "OptionContractResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoOptionContractResult"
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
    "name": "List Option Contracts",
    "description": {
      "content": "Retrieves option contracts filtered by underlying symbol, status and other attributes. Contains static contract information.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "options",
        "contracts",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Option category. Currently only US_OPTION is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Option symbols, multiple separated by commas.",
            "type": "text/plain"
          },
          "key": "option_symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Underlying symbols, multiple separated by commas.",
            "type": "text/plain"
          },
          "key": "underlying_symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Contract status, default LISTING. Enum: LISTING, DELISTING.",
            "type": "text/plain"
          },
          "key": "status",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Exact expiration date, format: YYYY-MM-DD.",
            "type": "text/plain"
          },
          "key": "start_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Expiration date lower bound (inclusive), format: YYYY-MM-DD.",
            "type": "text/plain"
          },
          "key": "end_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Root symbol filter (series symbol, e.g. SPXW). Mainly used for index options and post-CA non-standard contracts.",
            "type": "text/plain"
          },
          "key": "root_symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Contract type: CALL / PUT.",
            "type": "text/plain"
          },
          "key": "option_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Exercise style: AMERICAN / EUROPEAN.",
            "type": "text/plain"
          },
          "key": "style",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Strike price lower bound (inclusive).",
            "type": "text/plain"
          },
          "key": "strike_price_gte",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Strike price upper bound (inclusive).",
            "type": "text/plain"
          },
          "key": "strike_price_lte",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Penny Program Indicator: true = Penny Pilot contract, false = non-Penny Pilot.",
            "type": "text/plain"
          },
          "key": "ppind",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Whether to return deliverables array in response: TRUE / FALSE, default FALSE.",
            "type": "text/plain"
          },
          "key": "show_deliverables",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key from previous response for next page. Not required for first request.",
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


# Market Data — Futures — Verbatim Reference

> Futures market data and instruments.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Futures Tick

> Source: <https://developer.webull.hk/apis/docs/reference/futures-tick.md>

### List Futures Ticks

Retrieves futures tick-by-tick trade data.

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
  "path": "/market-data/futures/ticks/list",
  "method": "get",
  "tags": [
    "Futures Market Data"
  ],
  "description": "Retrieves futures tick-by-tick trade data.",
  "operationId": "futuresTick",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Futures symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "SILZ5"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_FUTURES` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of ticks, maximum limit 1200 .",
      "required": true,
      "schema": {
        "type": "string",
        "default": "200"
      },
      "example": 200
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
                "description": "List of tick details (trade prints) for this futures contract.",
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
            "title": "FuturesTickVo"
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
    "name": "List Futures Ticks",
    "description": {
      "content": "Retrieves futures tick-by-tick trade data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "futures",
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
            "content": "(Required) Futures symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_FUTURES` is supported for this interface.",
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

## Futures Snapshot

> Source: <https://developer.webull.hk/apis/docs/reference/futures-snapshot.md>

### List Futures Snapshots

Retrieves futures real-time snapshot data.

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
  "path": "/market-data/futures/snapshots/list",
  "method": "get",
  "tags": [
    "Futures Market Data"
  ],
  "description": "Retrieves futures real-time snapshot data.",
  "operationId": "futuresSnapshot",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of futures symbols, separated by commas; maximum 20 symbols per query. Example: SILZ5,6BM6.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "SILZ5,6BM6"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_FUTURES` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
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
                "price": {
                  "type": "string",
                  "description": "Last traded price (last done) of the futures contract, quoted in the contract's trading currency (e.g. USD).",
                  "example": "47.35"
                },
                "open": {
                  "type": "string",
                  "description": "Session open price for this futures contract. Represents the first traded price of the current regular trading session. If no trade has occurred in the session, this field may be empty.",
                  "example": "48.17"
                },
                "high": {
                  "type": "string",
                  "description": "Session high price for this futures contract during the current regular trading session. If no trade has occurred in the session, this field may be empty.",
                  "example": "48.655"
                },
                "low": {
                  "type": "string",
                  "description": "Session low price for this futures contract during the current regular trading session. If no trade has occurred in the session, this field may be empty.",
                  "example": "46.855"
                },
                "pre_close": {
                  "type": "string",
                  "description": "Previous settlement/close price of the futures contract (typically the official settlement price of the previous trading day), quoted in the contract's trading currency.",
                  "example": "47.704"
                },
                "volume": {
                  "type": "string",
                  "description": "Accumulated traded volume for the current session, expressed in number of futures contracts. If no trading occurred in the session, this field may be empty.",
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
                  "description": "Timestamp of the last executed trade for this futures contract, expressed as Unix epoch time in milliseconds",
                  "format": "int64",
                  "example": 1761131406558
                },
                "open_interest": {
                  "type": "string",
                  "description": "Open interest, representing the total number of outstanding and unsettled futures contracts for this instrument, expressed in number of contracts.",
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
                "bid_size": {
                  "type": "string",
                  "description": "Best bid size, i.e. the total quantity available at the best bid price, expressed in number of futures contracts (whole contract units).",
                  "example": "1"
                },
                "ask_size": {
                  "type": "string",
                  "description": "Best ask size, i.e. the total quantity available at the best ask price, expressed in number of futures contracts (whole contract units).",
                  "example": "2"
                },
                "settle_date": {
                  "type": "string",
                  "description": "Settlement date of the latest official daily settlement, typically in ISO-8601 datetime format with timezone information.",
                  "example": "2025-10-21T12:00:00.000+0000"
                },
                "settle_price": {
                  "type": "string",
                  "description": "Settlement price of the latest official daily settlement, used for marking the contract to market and determining margin requirements, quoted in the contract's trading currency.",
                  "example": "47.704"
                }
              },
              "description": "Market snapshot data response object",
              "title": "FuturesSnapshotVo"
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
    "name": "List Futures Snapshots",
    "description": {
      "content": "Retrieves futures real-time snapshot data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "futures",
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
            "content": "(Required) List of futures symbols, separated by commas; maximum 20 symbols per query. Example: SILZ5,6BM6.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_FUTURES` is supported for this interface.",
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

## Futures Footprint

> Source: <https://developer.webull.hk/apis/docs/reference/futures-footprint.md>

### List Futures Footprints

Retrieves futures footprint data.

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
  "path": "/market-data/futures/footprints/list",
  "method": "get",
  "tags": [
    "Futures Market Data"
  ],
  "description": "Retrieves futures footprint data.",
  "operationId": "futuresFootprint",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 20 symbols per query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "SILZ5,6BM6"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Category values are as shown in the enum; Only US_FUTURES type queries are supported.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
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
      "description": "Whether to include the latest unfinalized bar. Default: false. Applies to minute level timespans only",
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
      "description": "Specify trading hours. Only supported RTH",
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
              "required": [
                "symbol"
              ],
              "type": "object",
              "properties": {
                "symbol": {
                  "type": "string",
                  "description": "Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange.",
                  "example": "ESmain"
                },
                "instrument_id": {
                  "type": "string",
                  "description": "Unique instrument identifier for this futures contract in the Webull system or exchange.",
                  "example": "470004426"
                },
                "result": {
                  "type": "array",
                  "description": "Footprint result",
                  "items": {
                    "type": "object",
                    "properties": {
                      "time": {
                        "type": "string",
                        "description": "Bar timestamp in ISO-8601 format",
                        "example": "2025-11-13T14:51:00.000+0000"
                      },
                      "trading_session": {
                        "type": "string",
                        "description": "Trading session identifier only supported RTH (e.g.,RTH)",
                        "example": "RTH"
                      },
                      "total": {
                        "type": "string",
                        "description": "The sum of the main buy and sell volumes",
                        "example": "6053"
                      },
                      "delta": {
                        "type": "string",
                        "description": "The difference in trading volume (primary buyers - primary sellers)",
                        "example": "389"
                      },
                      "buy_total": {
                        "type": "string",
                        "description": " Buy-initiated volume",
                        "example": "3221"
                      },
                      "sell_total": {
                        "type": "string",
                        "description": "Sell-initiated volume",
                        "example": "2832"
                      },
                      "buy_detail": {
                        "type": "object",
                        "additionalProperties": {
                          "type": "string",
                          "description": "The main purchase footprint details (quantity combined for items with the same price).",
                          "example": "{\"6823\":\"7\",\"6823.25\":\"34\"}"
                        },
                        "description": "The main purchase footprint details (quantity combined for items with the same price).",
                        "example": {
                          "6823": "7",
                          "6823.25": "34"
                        }
                      },
                      "sell_detail": {
                        "type": "object",
                        "additionalProperties": {
                          "type": "string",
                          "description": "The main seller's footprint shows details (quantities combined for the same price).",
                          "example": "{\"6822.75\":\"1\",\"6823\":\"45\"}"
                        },
                        "description": "The main seller's footprint shows details (quantities combined for the same price).",
                        "example": {
                          "6823": "45",
                          "6822.75": "1"
                        }
                      }
                    },
                    "description": "FuturesFootprint",
                    "title": "FuturesFootprint"
                  }
                }
              },
              "title": "FuturesFootprintVo"
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
    "name": "List Futures Footprints",
    "description": {
      "content": "Retrieves futures footprint data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "futures",
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
            "content": "(Required) Security type. Category values are as shown in the enum; Only US_FUTURES type queries are supported.",
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
            "content": "(Required) Whether to include the latest unfinalized bar. Default: false. Applies to minute level timespans only",
            "type": "text/plain"
          },
          "key": "real_time_required",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Specify trading hours. Only supported RTH",
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

## Futures Quotes (Depth)

> Source: <https://developer.webull.hk/apis/docs/reference/futures-depth-of-book.md>

### List Futures Depths

Retrieves futures depth of book data.

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
  "path": "/market-data/futures/depths/list",
  "method": "get",
  "tags": [
    "Futures Market Data"
  ],
  "description": "Retrieves futures depth of book data.",
  "operationId": "futuresDepthOfBook",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Futures symbol.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "SILZ5"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_FUTURES` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
    },
    {
      "name": "depth",
      "in": "query",
      "description": "User-defined number of bid/ask levels (per side) to return in the Level-2 order book. <br/>Valid range: 1 – 10.<br/>(Note: Level-1 data must be requested separately through the Snapshot API.)",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": 10
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
                "description": "Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange.",
                "example": "SILZ5"
              },
              "instrument_id": {
                "type": "string",
                "description": "Unique instrument identifier for this futures contract in the Webull system or exchange.",
                "example": "470059643"
              },
              "quote_time": {
                "type": "integer",
                "description": "Quote timestamp of this snapshot, expressed as Unix epoch time in milliseconds.",
                "format": "int64",
                "example": 1761131409276
              },
              "asks": {
                "type": "array",
                "description": "Array of ask orders",
                "items": {
                  "required": [
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Ask price for this level in the futures order book.",
                      "example": "6770.75"
                    },
                    "size": {
                      "type": "string",
                      "description": "Ask size (quantity) at this price level, expressed in number of contracts.",
                      "example": "10"
                    }
                  },
                  "description": "AskItem. Represents a single level on the futures order book ask side.",
                  "title": "FuturesAskItem"
                }
              },
              "bids": {
                "type": "array",
                "description": "Array of bid orders",
                "items": {
                  "required": [
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Bid price for this level in the futures order book.",
                      "example": "6770.75"
                    },
                    "size": {
                      "type": "string",
                      "description": "Bid size (quantity) at this price level, expressed in number of contracts.",
                      "example": "10"
                    }
                  },
                  "description": "BidItem. Represents a single level on the futures order book bid side.",
                  "title": "FuturesBidItem"
                }
              }
            },
            "description": "DepthOfBookVo",
            "title": "FuturesDepthOfBookVo"
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
    "name": "List Futures Depths",
    "description": {
      "content": "Retrieves futures depth of book data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "futures",
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
            "content": "(Required) Futures symbol.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_FUTURES` is supported for this interface.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) User-defined number of bid/ask levels (per side) to return in the Level-2 order book. <br/>Valid range: 1 – 10.<br/>(Note: Level-1 data must be requested separately through the Snapshot API.)",
            "type": "text/plain"
          },
          "key": "depth",
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

## Futures Historical Bars

> Source: <https://developer.webull.hk/apis/docs/reference/futures-historical-bars.md>

### List Futures Historical Bars

Retrieves futures historical bars data.

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
  "path": "/market-data/futures/bars/list",
  "method": "get",
  "tags": [
    "Futures Market Data"
  ],
  "description": "Retrieves futures historical bars data.",
  "operationId": "futuresHistoricalBars",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of futures symbols, separated by commas; maximum 20 symbols per query. Example: SILZ5,6BM6.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "SILZ5,6BM6"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Currently only `US_FUTURES` is supported for this interface.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
    },
    {
      "name": "timespan",
      "in": "query",
      "description": "Bar time granularity. eg: M1, M5, M15, M30, M60, M120, M240, D, W, M, Y",
      "required": true,
      "schema": {
        "type": "string"
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
                "description": "List of batch bar data results, each element contains historical bar data for one futures.",
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
            "title": "FuturesBatchBarsVo"
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
    "name": "List Futures Historical Bars",
    "description": {
      "content": "Retrieves futures historical bars data.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "futures",
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
            "content": "(Required) List of futures symbols, separated by commas; maximum 20 symbols per query. Example: SILZ5,6BM6.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Currently only `US_FUTURES` is supported for this interface.",
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

## Futures Instrument List

> Source: <https://developer.webull.hk/apis/docs/reference/futures-instrument-list.md>

### List Futures Contracts

Retrieves detail information of futures instruments.

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
  "path": "/trading/instruments/futures/contracts/list",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves detail information of futures instruments.",
  "operationId": "futuresInstrumentList",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Supported values: US_FUTURES, HK_FUTURES",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
    },
    {
      "name": "symbols",
      "in": "query",
      "description": "List of futures trading symbols. Accepts JSON array format or comma-separated strings. Maximum of 100 symbols per request. Note: Either symbols or code must be provided. ",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "ESZ5,NQZ5"
    },
    {
      "name": "code",
      "in": "query",
      "description": "List of futures trading code, remark:Either 'symbols' or 'code' must be present.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "ES"
    },
    {
      "name": "status",
      "in": "query",
      "description": "Tradable Status<br/>OC - Tradable: Security is available for trading<br/>CO - Liquidate only: Security can only be sold, no purchases allowed<br/>NT - Non-Tradable: Security cannot be traded<br/>Default: OC",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "OC"
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
                  "description": "Futures contract symbol used in trading and market data, e.g. front-month or continuous contract code such as ESZ5, ESmain, etc.",
                  "example": "ESZ5"
                },
                "instrument_id": {
                  "type": "string",
                  "description": "Unique identifier for this futures instrument in the Webull system or exchange. If the symbol represents a main/continuous contract, this ID is for the main contract itself. For order placement, it needs to be mapped to the actual month contract ID (see `contractId`).",
                  "example": "470059643"
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Exchange code, for example: CBOE, GLOBEX, XNYM, XCEC, XCME, XCBT, CDE.",
                  "example": "XCME"
                },
                "code": {
                  "type": "string",
                  "description": "Code for this futures contract, for example: ES.",
                  "example": "ES"
                },
                "name": {
                  "type": "string",
                  "description": "Display name of the futures contract.",
                  "example": "E-mini S&P 500 Futures Dec 2025"
                },
                "product_class_id": {
                  "type": "integer",
                  "description": "Futures product class id, For example: 2",
                  "format": "int32",
                  "example": 2
                },
                "product_class_name": {
                  "type": "string",
                  "description": "Futures product class name, For example: 2",
                  "example": "Equities"
                },
                "status": {
                  "type": "string",
                  "description": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)  OC, CO, NT",
                  "example": "OC"
                },
                "currency": {
                  "type": "string",
                  "description": "Trading currency of this futures contract, for example: USD.",
                  "example": "USD"
                },
                "contract_month": {
                  "type": "string",
                  "description": "Contract delivery month in the format yyyyMM, for example: 202512 means Dec 2025 (year + month).",
                  "example": "202512"
                },
                "settlement_date": {
                  "type": "string",
                  "description": "Final settlement (delivery) date of the contract in the format yyyy-MM-dd, for example: 2025-12-29.",
                  "example": "2025-12-29"
                },
                "size": {
                  "type": "string",
                  "description": "Contract size (multiplier). The notional value of one contract equals futures price multiplied by this size. ",
                  "example": "50.0"
                },
                "unit": {
                  "type": "string",
                  "description": "Contract unit, describing the pricing unit and quantity (for example: index points x USD).",
                  "example": "1-index points",
                  "enum": [
                    "1 - Index points",
                    "2 - Hong Kong dollars",
                    "3 - US dollars",
                    "4 - Bushels",
                    "5 - Bushels 2",
                    "6 - Futures contract",
                    "7 - Short tons, 2000 pounds",
                    "8 - Pounds",
                    "9 - Gallons",
                    "10 - Metric tons, 2204.6 pounds",
                    "11 - Brazilian real",
                    "12 - Troy ounces",
                    "13 - British pounds",
                    "14 - Euros",
                    "15 - Mexican peso",
                    "16 - Czech koruna",
                    "17 - Polish zloty",
                    "18 - Israeli shekel",
                    "19 - Barrels",
                    "20 - Metric ton",
                    "21 - Australian dollar",
                    "22 - New Zealand dollar",
                    "23 - Canadian dollar",
                    "24 - Swiss franc",
                    "25 - Japanese yen",
                    "26 - South African rand",
                    "27 - Hungarian forint",
                    "28 - Korean won",
                    "29 - Million British thermal units",
                    "30 - Chinese renminbi",
                    "31 - Megawatt hours",
                    "41 - Megawatt",
                    "42 - Therms",
                    "51 - Environmental offset",
                    "52 - Basis points",
                    "53 - Metric tons (thousands)",
                    "54 - Gross tons",
                    "55 - Tons (thousands)",
                    "56 - Ton",
                    "57 - Bitcoin",
                    "58 - Russian ruble",
                    "59 - Indian rupee",
                    "60 - 1 day of time charter",
                    "61 - Cubic meter",
                    "62 - Kiloliters",
                    "63 - Kilos",
                    "64 - Chilean peso",
                    "65 - Regional Greenhouse Gas Initiative allowances (RGGI)",
                    "66 - Hundredweight, 100 pounds",
                    "67 - Norwegian krone",
                    "68 - Allowance (emission)",
                    "69 - Board feet",
                    "70 - Grams",
                    "71 - Swedish krona",
                    "72 - Environmental credit",
                    "73 - Dry metric tons",
                    "74 - Shares",
                    "75 - Metric ton",
                    "76 - Malaysian ringgit",
                    "77 - Ether",
                    "78 - Pounds net weight",
                    "79 - Renewable Identification Number (RIN)",
                    "80 - Barrels (thousands)",
                    "81 - Troy ounce (millions)"
                  ]
                },
                "min_tick": {
                  "type": "string",
                  "description": "Minimum price increment (tick size) for the futures price. ",
                  "example": "0.25"
                },
                "first_notice_date": {
                  "type": "string",
                  "description": "First notice date. For physically delivered contracts, this is the first date on which physical delivery can be assigned. After this date, new long positions cannot be opened, and existing long positions are typically forced to close a few trading days before this date. For cash-settled or index futures, this field is usually empty.",
                  "example": "2025-11-25"
                },
                "last_notice_date": {
                  "type": "string",
                  "description": "Last notice date, i.e. the last date on which the buyer can be notified to take physical delivery.",
                  "example": "2025-11-28"
                },
                "first_trading_date": {
                  "type": "string",
                  "description": "First trading date on which this futures contract becomes tradable.",
                  "example": "2024-12-01"
                },
                "last_trading_date": {
                  "type": "string",
                  "description": "Last trading date, i.e. the final trading day in the delivery month. After this date, any outstanding futures positions must be closed out through physical delivery or cash settlement. For cash-settled contracts, trading is allowed normally before the last trading deadline. For non-cash-settled contracts, opening new positions is usually restricted from three trading days before the earlier of the last trading date or first notice date.",
                  "example": "2025-12-19"
                },
                "contract_type": {
                  "type": "string",
                  "description": "Contract type. MONTHLY means regular month contract; MAIN means main/continuous contract.",
                  "example": "MONTHLY",
                  "enum": [
                    "MONTHLY",
                    "MAIN"
                  ]
                },
                "settlement": {
                  "type": "string",
                  "description": "Settlement method of the contract. Cash means cash settlement; Physical means physical delivery. ",
                  "example": "Cash",
                  "enum": [
                    "Cash",
                    "Physical"
                  ]
                }
              },
              "description": "Futures instrument profile information",
              "title": "FuturesInstrumentVo"
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
    "name": "List Futures Contracts",
    "description": {
      "content": "Retrieves detail information of futures instruments.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "futures",
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
            "content": "(Required) Security type. Supported values: US_FUTURES, HK_FUTURES",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "List of futures trading symbols. Accepts JSON array format or comma-separated strings. Maximum of 100 symbols per request. Note: Either symbols or code must be provided. ",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "List of futures trading code, remark:Either 'symbols' or 'code' must be present.",
            "type": "text/plain"
          },
          "key": "code",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Tradable Status<br/>OC - Tradable: Security is available for trading<br/>CO - Liquidate only: Security can only be sold, no purchases allowed<br/>NT - Non-Tradable: Security cannot be traded<br/>Default: OC",
            "type": "text/plain"
          },
          "key": "status",
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

## Futures Product Codes

> Source: <https://developer.webull.hk/apis/docs/reference/futures-products.md>

### List Futures Product Codes

Retrieves futures product codes list.

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
  "path": "/trading/instruments/futures/product-codes/list",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves futures product codes list.",
  "operationId": "futuresProducts",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Supported values: US_FUTURES, HK_FUTURES",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
    },
    {
      "name": "product_class_id",
      "in": "query",
      "description": "Product class id.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int32"
      },
      "example": 1
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
                "name": {
                  "type": "string",
                  "description": "Display name of the futures product, e.g., E-Mini S&P 500",
                  "example": "E-Mini S&P 500"
                },
                "code": {
                  "type": "string",
                  "description": "Futures product code: often one to three letter codes identifying the asset that is attached to a specific contract. For example: ES.",
                  "example": "ES"
                },
                "product_class_id": {
                  "type": "integer",
                  "description": "Futures product class id, For example: 2",
                  "format": "int32",
                  "example": 2
                },
                "product_class_name": {
                  "type": "string",
                  "description": "Futures product class name, For example: 2",
                  "example": "Equities"
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Futures product for exchange code, For example: XCME",
                  "example": "XCME"
                }
              },
              "description": "Futures product information",
              "title": "FuturesProductsVo"
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
    "name": "List Futures Product Codes",
    "description": {
      "content": "Retrieves futures product codes list.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "futures",
        "product-codes",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Supported values: US_FUTURES, HK_FUTURES",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Product class id.",
            "type": "text/plain"
          },
          "key": "product_class_id",
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

## Futures Product Classes

> Source: <https://developer.webull.hk/apis/docs/reference/futures-products-class.md>

### List Futures Product Classes

Retrieves futures product classes.

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
  "path": "/trading/instruments/futures/product-classes/list",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves futures product classes.",
  "operationId": "futuresProductsClass",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security type. Supported values: US_FUTURES, HK_FUTURES",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_FUTURES",
          "HK_FUTURES"
        ]
      },
      "example": "US_FUTURES"
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
                "product_class_id": {
                  "type": "integer",
                  "description": "Futures product class id",
                  "format": "int32",
                  "example": 2
                },
                "product_class_name": {
                  "type": "string",
                  "description": "Futures product class name",
                  "example": "Equities"
                }
              },
              "title": "FuturesProductClass"
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
    "name": "List Futures Product Classes",
    "description": {
      "content": "Retrieves futures product classes.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "futures",
        "product-classes",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security type. Supported values: US_FUTURES, HK_FUTURES",
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


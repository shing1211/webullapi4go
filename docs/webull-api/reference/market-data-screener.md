# Market Data — Screener — Verbatim Reference

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> Ranked lists and sector data. Some rankers require the Display Solution entitlement (marked below).

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Top Gainers/Losers

> Source: <https://developer.webull.hk/apis/docs/reference/get-gainers-losers.md>

### List Top Gainers/Losers

• Function description: Top Gainers/Losers. The only difference between Top Gainers and Top Losers is that when order=CHANGE_RATIO, direction is passed as ASC (Losers) and DESC (Gainers). Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

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
  "path": "/market-data/screeners/gainers-losers/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "• Function description: Top Gainers/Losers. The only difference between Top Gainers and Top Losers is that when order=CHANGE_RATIO, direction is passed as ASC (Losers) and DESC (Gainers). Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
  "operationId": "getGainersLosers",
  "parameters": [
    {
      "name": "rank_type",
      "in": "query",
      "description": "Ranking time dimension that determines the calculation period for price change. Default: DAY_1.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Rank type for gainers/losers screener. Time period for ranking by price change percentage.",
        "enum": [
          "PRE_MARKET",
          "AFTER_MARKET",
          "MIN_3",
          "MIN_5",
          "DAY_1",
          "DAY_5",
          "MONTH_1",
          "MONTH_3",
          "WEEK_52"
        ]
      },
      "example": "DAY_1"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security market category. Default: US_STOCK.",
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
      "description": "Secondary sort field for further ordering within the ranking. Default: CHANGE_RATIO.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Secondary sort field for screener API.",
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
      "description": "Sort direction. Default: DESC for Gainers, use ASC for Losers.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort direction for screener API.",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "DESC"
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
                  "description": "Unique identifier of the tradable instrument",
                  "example": "913256135"
                },
                "symbol": {
                  "type": "string",
                  "description": "Trading symbol of the financial instrument",
                  "example": "00700"
                },
                "name": {
                  "type": "string",
                  "description": "Full name of the instrument",
                  "example": "Tencent Holdings Ltd."
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Standardized exchange code",
                  "example": "HKEX"
                },
                "currency_code": {
                  "type": "string",
                  "description": "Denomination currency of the instrument (ISO 4217)",
                  "example": "HKD"
                },
                "pre_close": {
                  "type": "string",
                  "description": "Previous trading day's closing price",
                  "example": "380.2"
                },
                "open": {
                  "type": "string",
                  "description": "Opening price for the current trading day",
                  "example": "382.0"
                },
                "high": {
                  "type": "string",
                  "description": "Intraday high price for the current trading day",
                  "example": "388.6"
                },
                "low": {
                  "type": "string",
                  "description": "Intraday low price for the current trading day",
                  "example": "380.0"
                },
                "close": {
                  "type": "string",
                  "description": "Latest traded price for the current trading day",
                  "example": "385.6"
                },
                "price": {
                  "type": "string",
                  "description": "Most recent quoted price within the selected time interval (pre market / post market / intraday)",
                  "example": "385.6"
                },
                "change": {
                  "type": "string",
                  "description": "Absolute price change within the selected time interval. Returns pre/post market change when in extended hours",
                  "example": "5.4"
                },
                "change_ratio": {
                  "type": "string",
                  "description": "Price change percentage within the selected time interval (decimal ratio). Returns pre/post market change percentage when in extended hours",
                  "example": "0.0142"
                },
                "volume": {
                  "type": "string",
                  "description": "Cumulative traded volume for the current day (in shares)",
                  "example": "12345678"
                },
                "turnover": {
                  "type": "string",
                  "description": "Cumulative turnover amount in denomination currency",
                  "example": "4756789012"
                },
                "turnover_rate": {
                  "type": "string",
                  "description": "Turnover rate as a decimal ratio (e.g. 0.05 represents 5%)",
                  "example": "0.0013"
                },
                "market_value": {
                  "type": "string",
                  "description": "Total market capitalization in denomination currency",
                  "example": "3650000000000"
                },
                "amplitude": {
                  "type": "string",
                  "description": "Price amplitude ((high - low) / pre_close) as a decimal ratio",
                  "example": "0.0226"
                },
                "relative_volume_10d": {
                  "type": "string",
                  "description": "Relative volume (current day volume / 10 day average volume)",
                  "example": "11.91"
                }
              },
              "description": "Stock information from gainers/losers screener API.",
              "title": "ScreenerStockVo"
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
    "name": "List Top Gainers/Losers",
    "description": {
      "content": "• Function description: Top Gainers/Losers. The only difference between Top Gainers and Top Losers is that when order=CHANGE_RATIO, direction is passed as ASC (Losers) and DESC (Gainers). Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
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
            "content": "(Required) Ranking time dimension that determines the calculation period for price change. Default: DAY_1.",
            "type": "text/plain"
          },
          "key": "rank_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security market category. Default: US_STOCK.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Secondary sort field for further ordering within the ranking. Default: CHANGE_RATIO.",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort direction. Default: DESC for Gainers, use ASC for Losers.",
            "type": "text/plain"
          },
          "key": "direction",
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

## Top Actives

> Source: <https://developer.webull.hk/apis/docs/reference/get-top-active.md>

### List Top Actives

• Function description: Stock Top Active Rank. Most actively traded stocks ranked by volume, relative volume, turnover, turnover rate, or amplitude. The relative_volume_10d field is unique to the Top Active response compared to the Gainers/Losers endpoint. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

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
  "path": "/market-data/screeners/top-actives/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "• Function description: Stock Top Active Rank. Most actively traded stocks ranked by volume, relative volume, turnover, turnover rate, or amplitude. The relative_volume_10d field is unique to the Top Active response compared to the Gainers/Losers endpoint. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
  "operationId": "getTopActive",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security market category. Default: US_STOCK.",
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
      "name": "rank_type",
      "in": "query",
      "description": "Ranking dimension that determines which activity metric is used for filtering. Default: VOLUME.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Rank type for most active screener. Activity metric for ranking.",
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
      "name": "sort_by",
      "in": "query",
      "description": "Secondary sort field for further ordering within the ranking. Default: VOLUME.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Secondary sort field for screener API.",
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
      "description": "Sort direction. Default: DESC.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort direction for screener API.",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "DESC"
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
                  "description": "Unique identifier of the tradable instrument",
                  "example": "913256135"
                },
                "symbol": {
                  "type": "string",
                  "description": "Trading symbol of the financial instrument",
                  "example": "00700"
                },
                "name": {
                  "type": "string",
                  "description": "Full name of the instrument",
                  "example": "Tencent Holdings Ltd."
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Standardized exchange code",
                  "example": "HKEX"
                },
                "currency_code": {
                  "type": "string",
                  "description": "Denomination currency of the instrument (ISO 4217)",
                  "example": "HKD"
                },
                "pre_close": {
                  "type": "string",
                  "description": "Previous trading day's closing price",
                  "example": "380.2"
                },
                "open": {
                  "type": "string",
                  "description": "Opening price for the current trading day",
                  "example": "382.0"
                },
                "high": {
                  "type": "string",
                  "description": "Intraday high price for the current trading day",
                  "example": "388.6"
                },
                "low": {
                  "type": "string",
                  "description": "Intraday low price for the current trading day",
                  "example": "380.0"
                },
                "close": {
                  "type": "string",
                  "description": "Latest traded price for the current trading day",
                  "example": "385.6"
                },
                "price": {
                  "type": "string",
                  "description": "Most recent quoted price within the selected time interval (pre market / post market / intraday)",
                  "example": "385.6"
                },
                "change": {
                  "type": "string",
                  "description": "Absolute price change within the selected time interval. Returns pre/post market change when in extended hours",
                  "example": "5.4"
                },
                "change_ratio": {
                  "type": "string",
                  "description": "Price change percentage within the selected time interval (decimal ratio). Returns pre/post market change percentage when in extended hours",
                  "example": "0.0142"
                },
                "volume": {
                  "type": "string",
                  "description": "Cumulative traded volume for the current day (in shares)",
                  "example": "12345678"
                },
                "turnover": {
                  "type": "string",
                  "description": "Cumulative turnover amount in denomination currency",
                  "example": "4756789012"
                },
                "turnover_rate": {
                  "type": "string",
                  "description": "Turnover rate as a decimal ratio (e.g. 0.05 represents 5%)",
                  "example": "0.0013"
                },
                "market_value": {
                  "type": "string",
                  "description": "Total market capitalization in denomination currency",
                  "example": "3650000000000"
                },
                "amplitude": {
                  "type": "string",
                  "description": "Price amplitude ((high - low) / pre_close) as a decimal ratio",
                  "example": "0.0226"
                },
                "relative_volume_10d": {
                  "type": "string",
                  "description": "Relative volume (current day volume / 10 day average volume)",
                  "example": "11.91"
                }
              },
              "description": "Stock information from most active screener API. Contains additional relative_volume_10d field compared to gainers/losers.",
              "title": "MostActiveStockVo"
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
    "name": "List Top Actives",
    "description": {
      "content": "• Function description: Stock Top Active Rank. Most actively traded stocks ranked by volume, relative volume, turnover, turnover rate, or amplitude. The relative_volume_10d field is unique to the Top Active response compared to the Gainers/Losers endpoint. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
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
            "content": "(Required) Security market category. Default: US_STOCK.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Ranking dimension that determines which activity metric is used for filtering. Default: VOLUME.",
            "type": "text/plain"
          },
          "key": "rank_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Secondary sort field for further ordering within the ranking. Default: VOLUME.",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort direction. Default: DESC.",
            "type": "text/plain"
          },
          "key": "direction",
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

## Market Sectors

> Source: <https://developer.webull.hk/apis/docs/reference/get-market-sectors.md>

### List Market Sectors

• Function description: Get all sector overview data including sector name, change ratio, volume, market value, and leading stocks.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

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
  "path": "/market-data/screeners/market-sectors/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "• Function description: Get all sector overview data including sector name, change ratio, volume, market value, and leading stocks.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
  "operationId": "getMarketSectors",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security market category",
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
      "name": "agg_type",
      "in": "query",
      "description": "Statistics type, default is MARKET_VALUE.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Statistics type for market sectors.",
        "enum": [
          "MARKET_VALUE",
          "VOLUME"
        ]
      },
      "example": "MARKET_VALUE"
    },
    {
      "name": "period",
      "in": "query",
      "description": "Statistics period, default is D1.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Statistics period for market sectors.",
        "enum": [
          "D1",
          "D5",
          "MO1",
          "MO3"
        ]
      },
      "example": "D1"
    },
    {
      "name": "direction",
      "in": "query",
      "description": "Sort direction. Default: ASC.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort direction for screener API.",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "ASC"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZUluZGV4IjoxLCJwYWd"
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
                "description": "List of market sectors",
                "items": {
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "string",
                      "description": "Sector ID",
                      "example": "6391"
                    },
                    "name": {
                      "type": "string",
                      "description": "Sector Name",
                      "example": "Energy - Fossil Fuels"
                    },
                    "change_ratio": {
                      "type": "string",
                      "description": "Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%)",
                      "example": "0.0111"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Trading volume within the current statistical period",
                      "example": "1111111"
                    },
                    "market_value": {
                      "type": "string",
                      "description": "Market value within the current statistical period",
                      "example": "12345678901"
                    },
                    "declined": {
                      "type": "string",
                      "description": "Number of stocks that have fallen",
                      "example": "39"
                    },
                    "advanced": {
                      "type": "string",
                      "description": "Number of stocks that have risen",
                      "example": "100"
                    },
                    "flat": {
                      "type": "string",
                      "description": "Number of stocks with a price change of 0",
                      "example": "220"
                    },
                    "data": {
                      "type": "array",
                      "description": "Leading stocks in the sector",
                      "items": {
                        "type": "object",
                        "properties": {
                          "instrument_id": {
                            "type": "string",
                            "description": "Security ID",
                            "example": "913256135"
                          },
                          "name": {
                            "type": "string",
                            "description": "Security name",
                            "example": "Tesla Inc"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Security symbol",
                            "example": "TSLA"
                          }
                        },
                        "description": "Leading stock in a market sector",
                        "title": "LeadingStockVo"
                      }
                    }
                  },
                  "description": "Market Sector item",
                  "title": "MarketSectorsVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZUluZGV4IjoxLCJwYWd"
              }
            },
            "description": "Market Sectors paginated response",
            "title": "MarketSectorsResponseVo"
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
    "name": "List Market Sectors",
    "description": {
      "content": "• Function description: Get all sector overview data including sector name, change ratio, volume, market value, and leading stocks.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "screeners",
        "market-sectors",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security market category",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Statistics type, default is MARKET_VALUE.",
            "type": "text/plain"
          },
          "key": "agg_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Statistics period, default is D1.",
            "type": "text/plain"
          },
          "key": "period",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort direction. Default: ASC.",
            "type": "text/plain"
          },
          "key": "direction",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key from previous response for next page",
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

## Market Sector Detail

> Source: <https://developer.webull.hk/apis/docs/reference/get-market-sectors-detail.md>

### Get Market Sector Detail

• Function description: Get stock list and statistics for a specific sector.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

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
  "path": "/market-data/screeners/market-sectors/get",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "• Function description: Get stock list and statistics for a specific sector.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
  "operationId": "getMarketSectorsDetail",
  "parameters": [
    {
      "name": "sector_id",
      "in": "query",
      "description": "Sector ID",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": 6391
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security market category",
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
      "name": "period",
      "in": "query",
      "description": "Statistics period, default is D1.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Statistics period for market sectors.",
        "enum": [
          "D1",
          "D5",
          "MO1",
          "MO3"
        ]
      },
      "example": "D1"
    },
    {
      "name": "sort_by",
      "in": "query",
      "description": "Sort field, default is CHANGE_RATIO.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort fields for High Dividend, Market Sectors Detail, and 52 Week High/Low.",
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
          "VOLUME",
          "YIELD",
          "DIVIDEND"
        ]
      },
      "example": "CHANGE_RATIO"
    },
    {
      "name": "direction",
      "in": "query",
      "description": "Sort direction. Default: ASC.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort direction for screener API.",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "ASC"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZUluZGV4IjoxLCJwYWd"
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
              "id": {
                "type": "string",
                "description": "Sector ID",
                "example": "6391"
              },
              "name": {
                "type": "string",
                "description": "Sector Name",
                "example": "Energy - Fossil Fuels"
              },
              "change_ratio": {
                "type": "string",
                "description": "Price change ratio",
                "example": "0.013"
              },
              "declined": {
                "type": "string",
                "description": "Number of stocks that have fallen",
                "example": "39"
              },
              "advanced": {
                "type": "string",
                "description": "Number of stocks that have risen",
                "example": "110"
              },
              "flat": {
                "type": "string",
                "description": "Number of stocks with a price change of 0",
                "example": "107"
              },
              "data": {
                "type": "array",
                "description": "List of stocks in the sector",
                "items": {
                  "type": "object",
                  "properties": {
                    "instrument_id": {
                      "type": "string",
                      "description": "Security ID",
                      "example": "913256135"
                    },
                    "category": {
                      "type": "string",
                      "description": "Security category",
                      "example": "US_STOCK"
                    },
                    "currency": {
                      "type": "string",
                      "description": "Currency code",
                      "example": "USD"
                    },
                    "name": {
                      "type": "string",
                      "description": "Security name",
                      "example": "Tesla Inc"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Security symbol",
                      "example": "TSLA"
                    },
                    "exchange_code": {
                      "type": "string",
                      "description": "Exchange code",
                      "example": "NSQ"
                    },
                    "close": {
                      "type": "string",
                      "description": "Latest intraday price",
                      "example": "12.01"
                    },
                    "change_ratio": {
                      "type": "string",
                      "description": "Price change ratio",
                      "example": "0.0111"
                    },
                    "price": {
                      "type": "string",
                      "description": "Latest price",
                      "example": "11.01"
                    },
                    "volume": {
                      "type": "string",
                      "description": "Trade volume",
                      "example": "1111111"
                    },
                    "market_value": {
                      "type": "string",
                      "description": "Market Value",
                      "example": "12345678901"
                    },
                    "turnover_rate": {
                      "type": "string",
                      "description": "Turnover Rate",
                      "example": "0.1111"
                    },
                    "amplitude": {
                      "type": "string",
                      "description": "Amplitude Ratio",
                      "example": "0.1111"
                    },
                    "high": {
                      "type": "string",
                      "description": "Today's high",
                      "example": "12.11"
                    },
                    "low": {
                      "type": "string",
                      "description": "Today's low",
                      "example": "11.01"
                    }
                  },
                  "description": "Stock in a market sector detail",
                  "title": "SectorDetailStockVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZUluZGV4IjoxLCJwYWd"
              }
            },
            "description": "Market Sectors Detail result",
            "title": "MarketSectorsDetailVo"
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
    "name": "Get Market Sector Detail",
    "description": {
      "content": "• Function description: Get stock list and statistics for a specific sector.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "screeners",
        "market-sectors",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Sector ID",
            "type": "text/plain"
          },
          "key": "sector_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security market category",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Statistics period, default is D1.",
            "type": "text/plain"
          },
          "key": "period",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort field, default is CHANGE_RATIO.",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort direction. Default: ASC.",
            "type": "text/plain"
          },
          "key": "direction",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key from previous response for next page",
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

## High Dividend Rank

> Source: <https://developer.webull.hk/apis/docs/reference/get-high-dividend.md>

### List High Dividend Rank

• Function description: Get high dividend rank list. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

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
  "path": "/market-data/screeners/high-dividend-ranks/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "• Function description: Get high dividend rank list. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
  "operationId": "getHighDividend",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security market category",
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
      "description": "Sort field, default is YIELD.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort fields for High Dividend, Market Sectors Detail, and 52 Week High/Low.",
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
          "VOLUME",
          "YIELD",
          "DIVIDEND"
        ]
      },
      "example": "YIELD"
    },
    {
      "name": "direction",
      "in": "query",
      "description": "Sort direction. Default: DESC.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort direction for screener API.",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "DESC"
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
                  "description": "Security ID",
                  "example": "913256135"
                },
                "category": {
                  "type": "string",
                  "description": "Security category",
                  "example": "US_STOCK"
                },
                "currency": {
                  "type": "string",
                  "description": "Currency code",
                  "example": "USD"
                },
                "name": {
                  "type": "string",
                  "description": "Security name",
                  "example": "Tesla Inc"
                },
                "symbol": {
                  "type": "string",
                  "description": "Security symbol",
                  "example": "TSLA"
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Exchange code",
                  "example": "NSQ"
                },
                "close": {
                  "type": "string",
                  "description": "Latest intraday price",
                  "example": "12.01"
                },
                "change": {
                  "type": "string",
                  "description": "Trade change for the current trading day",
                  "example": "1.01"
                },
                "change_ratio": {
                  "type": "string",
                  "description": "Price change ratio",
                  "example": "0.0111"
                },
                "price": {
                  "type": "string",
                  "description": "Latest price",
                  "example": "11.01"
                },
                "volume": {
                  "type": "string",
                  "description": "Trade volume",
                  "example": "1111111"
                },
                "market_value": {
                  "type": "string",
                  "description": "Market Value",
                  "example": "12345678901"
                },
                "turnover_rate": {
                  "type": "string",
                  "description": "Turnover Rate",
                  "example": "0.1111"
                },
                "amplitude": {
                  "type": "string",
                  "description": "Amplitude Ratio",
                  "example": "0.1111"
                },
                "high": {
                  "type": "string",
                  "description": "Today's high",
                  "example": "12.11"
                },
                "low": {
                  "type": "string",
                  "description": "Today's low",
                  "example": "11.01"
                },
                "turnover": {
                  "type": "string",
                  "description": "Turnover amount",
                  "example": "1234567890"
                },
                "yield": {
                  "type": "string",
                  "description": "Dividend Yield",
                  "example": "0.1111"
                },
                "dividend": {
                  "type": "string",
                  "description": "Dividend",
                  "example": "0.1111"
                },
                "ex_date": {
                  "type": "string",
                  "description": "Ex-Date",
                  "example": "2026-01-08"
                },
                "pe_ttm": {
                  "type": "string",
                  "description": "PE TTM",
                  "example": "0.1111"
                }
              },
              "description": "High Dividend stock item",
              "title": "HighDividendStockVo"
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
    "name": "List High Dividend Rank",
    "description": {
      "content": "• Function description: Get high dividend rank list. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "screeners",
        "high-dividend-ranks",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security market category",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort field, default is YIELD.",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort direction. Default: DESC.",
            "type": "text/plain"
          },
          "key": "direction",
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

## 52-Week High/Low

> Source: <https://developer.webull.hk/apis/docs/reference/get-week-52-high-low.md>

### List 52-Week High/Low

• Function description: Get 52 week high/low rank list. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

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
  "path": "/market-data/screeners/week52-high-low/list",
  "method": "get",
  "tags": [
    "Screeners"
  ],
  "description": "• Function description: Get 52 week high/low rank list. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
  "operationId": "getWeek52HighLow",
  "parameters": [
    {
      "name": "rank_type",
      "in": "query",
      "description": "52 week rank type.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "52 Week High/Low rank type.",
        "enum": [
          "NEW_HIGH",
          "NEAR_HIGH",
          "NEW_LOW",
          "NEAR_LOW"
        ]
      },
      "example": "NEW_HIGH"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Security market category",
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
      "description": "Sort field, default is CHANGE_RATIO_52W.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort fields for High Dividend, Market Sectors Detail, and 52 Week High/Low.",
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
          "VOLUME",
          "YIELD",
          "DIVIDEND"
        ]
      },
      "example": "CHANGE_RATIO_52W"
    },
    {
      "name": "direction",
      "in": "query",
      "description": "Sort direction. Default: ASC.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sort direction for screener API.",
        "enum": [
          "ASC",
          "DESC"
        ]
      },
      "example": "ASC"
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
                  "description": "Security ID",
                  "example": "913256135"
                },
                "category": {
                  "type": "string",
                  "description": "Security category",
                  "example": "US_STOCK"
                },
                "currency": {
                  "type": "string",
                  "description": "Currency code",
                  "example": "USD"
                },
                "name": {
                  "type": "string",
                  "description": "Security name",
                  "example": "Tesla Inc"
                },
                "symbol": {
                  "type": "string",
                  "description": "Security symbol",
                  "example": "TSLA"
                },
                "exchange_code": {
                  "type": "string",
                  "description": "Exchange code",
                  "example": "NSQ"
                },
                "close": {
                  "type": "string",
                  "description": "Latest intraday price",
                  "example": "12.01"
                },
                "change": {
                  "type": "string",
                  "description": "Trade change for the current trading day",
                  "example": "1.01"
                },
                "change_ratio": {
                  "type": "string",
                  "description": "Price change ratio",
                  "example": "0.0111"
                },
                "price": {
                  "type": "string",
                  "description": "Latest price",
                  "example": "11.01"
                },
                "volume": {
                  "type": "string",
                  "description": "Trade Volume",
                  "example": "1111111"
                },
                "market_value": {
                  "type": "string",
                  "description": "Market Value",
                  "example": "12345678901"
                },
                "turnover_rate": {
                  "type": "string",
                  "description": "Turnover Rate",
                  "example": "0.1111"
                },
                "high": {
                  "type": "string",
                  "description": "Today's high",
                  "example": "12.11"
                },
                "low": {
                  "type": "string",
                  "description": "Today's low",
                  "example": "11.01"
                },
                "turnover": {
                  "type": "string",
                  "description": "Turnover amount",
                  "example": "1234567890"
                },
                "amplitude": {
                  "type": "string",
                  "description": "Amplitude Ratio",
                  "example": "0.1111"
                },
                "price_1w": {
                  "type": "string",
                  "description": "This Week High/Low",
                  "example": "1.01"
                },
                "price_52w": {
                  "type": "string",
                  "description": "52W Last High/Low",
                  "example": "1.01"
                },
                "change_ratio_52w": {
                  "type": "string",
                  "description": "Change ratio since previous highest/lowest",
                  "example": "0.0111"
                },
                "pe_ttm": {
                  "type": "string",
                  "description": "PE TTM",
                  "example": "0.1111"
                }
              },
              "description": "52 Week High/Low stock item",
              "title": "Week52StockVo"
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
    "name": "List 52-Week High/Low",
    "description": {
      "content": "• Function description: Get 52 week high/low rank list. Returns top 200 results without pagination.<br/>• Frequency limit: Market-data interfaces rate limit is 600 requests per minute.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "screeners",
        "week52-high-low",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "52 week rank type.",
            "type": "text/plain"
          },
          "key": "rank_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Security market category",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort field, default is CHANGE_RATIO_52W.",
            "type": "text/plain"
          },
          "key": "sort_by",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sort direction. Default: ASC.",
            "type": "text/plain"
          },
          "key": "direction",
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


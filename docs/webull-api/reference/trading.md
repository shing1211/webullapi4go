# Trading API — Verbatim Reference

> Accounts, assets, the order lifecycle, order queries and instruments. Requests require an access token and default to API version `v3`.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Get Instruments

> Source: <https://developer.webull.hk/apis/docs/reference/instrument-list.md>

### List Stock Instruments

Retrieves profile information for one or more stock instruments.

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
  "path": "/trading/instruments/stocks/profiles/list",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves profile information for one or more stock instruments.",
  "operationId": "instrumentList",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security type.",
      "required": true,
      "schema": {
        "type": "string",
        "enum": [
          "US_STOCK",
          "HK_STOCK",
          "CN_STOCK"
        ]
      },
      "example": "US_STOCK"
    },
    {
      "name": "symbols",
      "in": "query",
      "description": "List of security symbols, maximum 100 symbols per query.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "AAPL,TSLA"
    },
    {
      "name": "status",
      "in": "query",
      "description": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
        "enum": [
          "OC",
          "CO",
          "NT"
        ]
      },
      "example": "CO"
    },
    {
      "name": "sub_category",
      "in": "query",
      "description": "Sub-category of the instrument. Only effective when symbols is not specified. When category = US_STOCK, supported values: COMMON_STOCK, ETF, PREFERRED_STOCK, WARRANT, UNITS, RIGHT. When category = HK_STOCK, CN_STOCK, supported values: COMMON_STOCK, ETF. If not specified, returns all sub-categories.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Sub-category of the instrument.",
        "enum": [
          "COMMON_STOCK",
          "ETF",
          "PREFERRED_STOCK",
          "WARRANT",
          "UNITS",
          "RIGHT"
        ]
      },
      "example": "ETF"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
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
                    "name": {
                      "type": "string",
                      "description": "Symbol name, e.g. Apple",
                      "example": "APPLE INC"
                    },
                    "instrument_id": {
                      "type": "string",
                      "description": "Unique identifier of the security",
                      "example": "10152734329"
                    },
                    "exchange_code": {
                      "type": "string",
                      "description": "Exchange code, e.g. CCC",
                      "example": "NSQ"
                    },
                    "category": {
                      "type": "string",
                      "description": "Instrument type, e.g. US_STOCK",
                      "example": "US_STOCK",
                      "enum": [
                        "US_STOCK",
                        "HK_STOCK",
                        "CN_STOCK"
                      ]
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Symbol of the instrument",
                      "example": "AAPL"
                    },
                    "status": {
                      "type": "string",
                      "description": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
                      "example": "OC",
                      "enum": [
                        "OC",
                        "CO",
                        "NT"
                      ]
                    },
                    "shortable": {
                      "type": "boolean",
                      "description": "Instrument is shortable or not",
                      "example": false
                    },
                    "fractionable": {
                      "type": "boolean",
                      "description": "Instrument is fractionable or not ",
                      "example": false
                    },
                    "marginable": {
                      "type": "boolean",
                      "description": "Instrument is marginable or not",
                      "example": false
                    },
                    "overnight_trading_supported": {
                      "type": "boolean",
                      "description": "Instrument support overnight trading or not ",
                      "example": false
                    },
                    "margin_requirement_long": {
                      "type": "string",
                      "description": "Margin requirement ratio for long position ",
                      "example": "0.5"
                    },
                    "margin_requirement_short": {
                      "type": "string",
                      "description": "Margin requirement ratio for short position ",
                      "example": "0.5"
                    },
                    "intraday_margin_long": {
                      "type": "string",
                      "description": "Intraday margin requirement ratio for long position ",
                      "example": "0.5"
                    },
                    "intraday_margin_short": {
                      "type": "string",
                      "description": "Intraday margin requirement ratio for short position ",
                      "example": "0.5"
                    },
                    "maintenance_margin_long": {
                      "type": "string",
                      "description": "Maintenance margin requirement ratio for long position ",
                      "example": "0.5"
                    },
                    "maintenance_margin_short": {
                      "type": "string",
                      "description": "Maintenance margin requirement ratio for short position ",
                      "example": "0.5"
                    },
                    "easy_to_borrow": {
                      "type": "boolean",
                      "description": "Instrument is easy to borrow or not",
                      "example": false
                    },
                    "lot_size": {
                      "type": "string",
                      "description": "Lot size",
                      "example": "1.0"
                    },
                    "currency": {
                      "type": "string",
                      "description": "currency",
                      "example": "USD"
                    },
                    "sub_category": {
                      "type": "string",
                      "description": "Sub-category of the instrument.",
                      "example": "COMMON_STOCK",
                      "enum": [
                        "COMMON_STOCK",
                        "ETF",
                        "PREFERRED_STOCK",
                        "WARRANT",
                        "UNITS",
                        "RIGHT"
                      ]
                    }
                  },
                  "description": "Instrument Information",
                  "title": "InstrumentStockDetailVO"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoInstrumentStockDetailVO"
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
    "name": "List Stock Instruments",
    "description": {
      "content": "Retrieves profile information for one or more stock instruments.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "stocks",
        "profiles",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
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
            "content": "List of security symbols, maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
            "type": "text/plain"
          },
          "key": "status",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Sub-category of the instrument. Only effective when symbols is not specified. When category = US_STOCK, supported values: COMMON_STOCK, ETF, PREFERRED_STOCK, WARRANT, UNITS, RIGHT. When category = HK_STOCK, CN_STOCK, supported values: COMMON_STOCK, ETF. If not specified, returns all sub-categories.",
            "type": "text/plain"
          },
          "key": "sub_category",
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

## Account List

> Source: <https://developer.webull.hk/apis/docs/reference/account-list.md>

### List Accounts

Retrieves the account list and returns account information.

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
  "path": "/trading/accounts/list",
  "method": "get",
  "tags": [
    "Account"
  ],
  "description": "Retrieves the account list and returns account information.",
  "operationId": "accountList",
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
                "account_id": {
                  "type": "string",
                  "description": "Account identifier",
                  "example": "LOJOQITOD49R6G9BPQM489CISA"
                },
                "account_number": {
                  "type": "string",
                  "description": "Brokerage account",
                  "example": "10010048"
                },
                "account_type": {
                  "type": "string",
                  "description": "Account type",
                  "example": "CASH",
                  "enum": [
                    "MARGIN",
                    "CASH"
                  ]
                },
                "account_class": {
                  "type": "string",
                  "description": "Account Class",
                  "example": "INDIVIDUAL_CASH",
                  "enum": [
                    "INDIVIDUAL_CASH",
                    "INDIVIDUAL_MRGN",
                    "FUTURES_MRGN",
                    "INSTITUTIONAL_CASH",
                    "INSTITUTIONAL_MRGN",
                    "INSTITUTIONAL_FUTURES_MRGN"
                  ]
                }
              },
              "title": "AccountListResult"
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
    "name": "List Accounts",
    "description": {
      "content": "Retrieves the account list and returns account information.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "accounts",
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
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Account Balance

> Source: <https://developer.webull.hk/apis/docs/reference/query-account-balance.md>

### Get Account Balance

Retrieves account details by account ID.

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
  "path": "/trading/assets/balances/get",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves account details by account ID.",
  "operationId": "queryAccountBalance",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LOJOQITOD49R6G9BPQM489CISA"
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
              "account_currency_assets",
              "total_asset_currency",
              "total_cash_balance",
              "total_market_value",
              "total_unrealized_profit_loss"
            ],
            "type": "object",
            "properties": {
              "total_asset_currency": {
                "type": "string",
                "description": "Currency",
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "total_cash_balance": {
                "type": "string",
                "description": "Cash Balance",
                "example": "485705.0"
              },
              "total_market_value": {
                "type": "string",
                "description": "Total holding market value",
                "example": "995705.0"
              },
              "total_unrealized_profit_loss": {
                "type": "string",
                "description": "Open P&L",
                "example": "227689.0"
              },
              "init_margin": {
                "type": "string",
                "description": "Initial margin",
                "example": "18000.0"
              },
              "account_currency_assets": {
                "type": "array",
                "description": "Currency assets Details",
                "items": {
                  "required": [
                    "available_withdrawal",
                    "buying_power",
                    "cash_balance",
                    "currency",
                    "interests_unpaid",
                    "market_value",
                    "settled_cash",
                    "unrealized_profit_loss",
                    "unsettled_cash"
                  ],
                  "type": "object",
                  "properties": {
                    "currency": {
                      "type": "string",
                      "description": "Currency",
                      "example": "HKD",
                      "enum": [
                        "CNH",
                        "HKD",
                        "USD"
                      ]
                    },
                    "cash_balance": {
                      "type": "string",
                      "description": "Cash Balance",
                      "example": "485705.95"
                    },
                    "settled_cash": {
                      "type": "string",
                      "description": "Settled Cash",
                      "example": "485705.95"
                    },
                    "unsettled_cash": {
                      "type": "string",
                      "description": "Unsettled Cash",
                      "example": "0.0"
                    },
                    "market_value": {
                      "type": "string",
                      "description": "holding market value",
                      "example": "0.0"
                    },
                    "held_amount": {
                      "type": "string",
                      "description": "In-transit funds",
                      "example": "0.0"
                    },
                    "frozen_amount": {
                      "type": "string",
                      "description": "Frozen funds",
                      "example": "485705"
                    },
                    "buying_power": {
                      "type": "string",
                      "description": "Buying Power",
                      "example": "484551"
                    },
                    "unrealized_profit_loss": {
                      "type": "string",
                      "description": "Open P&L",
                      "example": "227689"
                    },
                    "available_withdrawal": {
                      "type": "string",
                      "description": "The withdrawable amount",
                      "example": "3.0558743194E8"
                    },
                    "interests_unpaid": {
                      "type": "string",
                      "description": "Interest to be paid",
                      "example": "0.0"
                    },
                    "init_margin": {
                      "type": "string",
                      "description": "Init margin",
                      "example": "18000.0"
                    }
                  },
                  "description": "Currency assets Details",
                  "title": "AssetsCurrencyAssets"
                }
              }
            },
            "title": "AssetsBalanceResult"
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
    "name": "Get Account Balance",
    "description": {
      "content": "Retrieves account details by account ID.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "assets",
        "balances",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Account identifier",
            "type": "text/plain"
          },
          "key": "account_id",
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

## Account Positions

> Source: <https://developer.webull.hk/apis/docs/reference/query-account-position.md>

### List Account Positions

Retrieves positions according to the account ID

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
  "path": "/trading/assets/positions/list",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves positions according to the account ID",
  "operationId": "queryAccountPosition",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LOJOQITOD49R6G9BPQM489CISA"
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
                "cost_price",
                "currency",
                "instrument_type",
                "last_price",
                "option_strategy",
                "position_id",
                "quantity",
                "symbol",
                "unrealized_profit_loss"
              ],
              "type": "object",
              "properties": {
                "position_id": {
                  "type": "string",
                  "description": "Position ID",
                  "example": "N4I4SIM8TJF38KN2TAA0QVVNE9"
                },
                "currency": {
                  "type": "string",
                  "description": "Currency",
                  "example": "USD",
                  "enum": [
                    "CNH",
                    "HKD",
                    "USD"
                  ]
                },
                "quantity": {
                  "type": "string",
                  "description": "Quantity of the order. Specifies the number of shares or units to transact. <br/> For US stocks, fractional quantities are allowed and can include decimals.",
                  "example": "1"
                },
                "symbol": {
                  "type": "string",
                  "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                  "example": "AAPL"
                },
                "option_strategy": {
                  "type": "string",
                  "description": "Type of options strategy <br/> &bull; SINGLE: Indicates a single-leg options order",
                  "example": "SINGLE",
                  "enum": [
                    "SINGLE"
                  ]
                },
                "instrument_type": {
                  "type": "string",
                  "description": "Type of financial instrument associated with the request.",
                  "example": "EQUITY",
                  "enum": [
                    "EQUITY",
                    "OPTION",
                    "FUTURES"
                  ]
                },
                "last_price": {
                  "type": "string",
                  "description": "Last Price",
                  "example": "10.0"
                },
                "cost_price": {
                  "type": "string",
                  "description": "Cost Basis",
                  "example": "11.12"
                },
                "unrealized_profit_loss": {
                  "type": "string",
                  "description": "Open P&L",
                  "example": "0.08"
                },
                "legs": {
                  "type": "array",
                  "description": "legs",
                  "items": {
                    "required": [
                      "symbol"
                    ],
                    "type": "object",
                    "properties": {
                      "symbol": {
                        "type": "string",
                        "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                        "example": "AAPL"
                      },
                      "quantity": {
                        "type": "string",
                        "description": "Quantity of the order. Specifies the number of shares or units to transact.",
                        "example": "4"
                      },
                      "option_type": {
                        "type": "string",
                        "description": "Type of the option. <br/> &bull; CALL: Right to buy the underlying asset. <br/> &bull; PUT: Right to sell the underlying asset.",
                        "example": "CALL",
                        "enum": [
                          "CALL",
                          "PUT"
                        ]
                      },
                      "option_expire_date": {
                        "type": "string",
                        "description": "Option expiration date. Format: yyyy-MM-dd",
                        "example": "2019-09-20"
                      },
                      "option_exercise_price": {
                        "type": "string",
                        "description": "Exercise Price",
                        "example": "11.0"
                      },
                      "option_contract_multiplier": {
                        "type": "string",
                        "description": "The number of shares corresponding to each option contract",
                        "example": "100"
                      },
                      "option_contract_deliverable": {
                        "type": "string",
                        "description": "The number of shares required to exercise each contract",
                        "example": "100"
                      },
                      "expiration_type": {
                        "type": "string",
                        "description": "Option expiration types",
                        "example": "AM"
                      }
                    },
                    "description": "legs",
                    "title": "PositionItem"
                  }
                }
              },
              "title": "AssetsPositionResult"
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
    "name": "List Account Positions",
    "description": {
      "content": "Retrieves positions according to the account ID",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "assets",
        "positions",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Account identifier",
            "type": "text/plain"
          },
          "key": "account_id",
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

## Order Preview

> Source: <https://developer.webull.hk/apis/docs/reference/common-order-preview.md>

### Preview Order

Calculates the estimated amount and cost based on the provided information. Supports simple orders.

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
  "path": "/trading/orders/preview",
  "method": "post",
  "tags": [
    "Trading"
  ],
  "description": "Calculates the estimated amount and cost based on the provided information. Supports simple orders.",
  "operationId": "Common Order Preview",
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
    "description": "Order Preview",
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "account_id",
            "new_orders"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "client_combo_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the combined order<br/> If combo_type = NORMAL and client_combo_order_id not need to set<br/> If combo_type != NORMAL and client_combo_order_id not provided <br/> the server will automatically generate one.<br/> To sell and close an existing position with take-profit/stop-loss, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) grouped under the same client_combo_order_id; no MASTER order is required in this scenario.",
              "example": "0KGOHL4PR2SLC0DKIND4TI0001"
            },
            "new_orders": {
              "type": "array",
              "description": "Order Details",
              "items": {
                "required": [
                  "client_order_id",
                  "combo_type",
                  "entrust_type",
                  "instrument_type",
                  "market",
                  "order_type",
                  "side",
                  "support_trading_session",
                  "symbol",
                  "time_in_force"
                ],
                "type": "object",
                "properties": {
                  "combo_type": {
                    "type": "string",
                    "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> Futures Trading currently support only the NORMAL type.<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) <br/> Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened.<br/> Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO.",
                    "example": "NORMAL"
                  },
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY",
                      "OPTION",
                      "FUTURES"
                    ]
                  },
                  "market": {
                    "type": "string",
                    "description": "Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                    "example": "US",
                    "enum": [
                      "US",
                      "HK",
                      "CN"
                    ]
                  },
                  "symbol": {
                    "type": "string",
                    "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market.",
                    "example": "AAPL"
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order <br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order <br/>China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
                    "example": "MARKET",
                    "enum": [
                      "LIMIT",
                      "MARKET",
                      "STOP_LOSS",
                      "STOP_LOSS_LIMIT",
                      "ENHANCED_LIMIT",
                      "AT_AUCTION",
                      "AT_AUCTION_LIMIT",
                      "MARKET_ON_OPEN",
                      "TRAILING_STOP_LOSS",
                      "TRAILING_STOP_LOSS_LIMIT",
                      "TOUCH_MKT",
                      "TOUCH_LMT",
                      "ODD_LOT_LIMIT"
                    ]
                  },
                  "entrust_type": {
                    "type": "string",
                    "description": "Specifies the method for placing the order. <br/> &bull; QTY: Order specified by quantity of shares or units. <br/> &bull; AMOUNT: Order specified by total cash amount, applicable for fractional share trading of US stocks.",
                    "example": "QTY",
                    "enum": [
                      "QTY",
                      "AMOUNT"
                    ]
                  },
                  "support_trading_session": {
                    "type": "string",
                    "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only. <br/> Deprecated values: <br/> &bull; Y: [Deprecated]Include extended trading hours. <br/> &bull; N: [Deprecated]Only support regular trading hours. <br/> Active values: <br/> &bull; NIGHT: Only supports night trading. <br/> &bull; ALL: Include extended trading hours. <br/> &bull; CORE: Only support regular trading hours. <br/> &bull; ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day)",
                    "example": "CORE",
                    "enum": [
                      "Y",
                      "N",
                      "NIGHT",
                      "ALL",
                      "CORE",
                      "ALL_DAY"
                    ]
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTD",
                      "GTC"
                    ]
                  },
                  "side": {
                    "type": "string",
                    "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                    "example": "BUY",
                    "enum": [
                      "BUY",
                      "SELL",
                      "SHORT"
                    ]
                  },
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                    "example": "11.0"
                  },
                  "stop_price": {
                    "type": "string",
                    "description": "Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT.<br/> Specifies the trigger price at which the stop order becomes active.",
                    "example": "11.0"
                  },
                  "option_strategy": {
                    "type": "string",
                    "description": "Type of options strategy <br/> &bull; SINGLE: Indicates a single-leg options order",
                    "example": "SINGLE",
                    "enum": [
                      "SINGLE"
                    ]
                  },
                  "legs": {
                    "type": "array",
                    "description": "Option leg detail. Only required when previewing option orders.",
                    "items": {
                      "required": [
                        "instrument_type",
                        "market",
                        "side",
                        "symbol"
                      ],
                      "type": "object",
                      "properties": {
                        "instrument_type": {
                          "type": "string",
                          "description": "Type of financial instrument associated with the request.",
                          "example": "OPTION",
                          "enum": [
                            "EQUITY",
                            "OPTION",
                            "FUTURES"
                          ]
                        },
                        "market": {
                          "type": "string",
                          "description": "Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                          "example": "US",
                          "enum": [
                            "US",
                            "HK",
                            "CN"
                          ]
                        },
                        "symbol": {
                          "type": "string",
                          "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                          "example": "AAPL"
                        },
                        "side": {
                          "type": "string",
                          "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                          "example": "BUY",
                          "enum": [
                            "BUY",
                            "SELL",
                            "SHORT"
                          ]
                        },
                        "strike_price": {
                          "type": "string",
                          "description": "Exercise price (strike price) of the option. <br/> Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise.",
                          "example": "11.0"
                        },
                        "option_expire_date": {
                          "type": "string",
                          "description": "Expiration date. Format: yyyy-MM-dd",
                          "example": "2025-08-01"
                        },
                        "option_type": {
                          "type": "string",
                          "description": "Type of the option. <br/> &bull; CALL: Right to buy the underlying asset. <br/> &bull; PUT: Right to sell the underlying asset.",
                          "example": "CALL",
                          "enum": [
                            "CALL",
                            "PUT"
                          ]
                        },
                        "quantity": {
                          "type": "string",
                          "description": "Quantity of the order or strategy leg. <br/>For stock legs, specifies the number of shares to transact <br/>For option legs, specifies the number of option contracts to transact for this leg and is expressed in whole contracts.",
                          "example": "1"
                        }
                      },
                      "description": "Option leg detail. Only required when placing option orders.",
                      "title": "OptionCommonPlaceLegParam"
                    }
                  }
                },
                "description": "Order Details",
                "title": "OrderCommonPreviewItemParam"
              }
            }
          },
          "title": "OrderCommonPreviewParam"
        },
        "examples": {
          "Equity": {
            "summary": "Stock Limit Order",
            "description": "Equity",
            "value": {
              "account_id": "<your_account_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id>",
                  "combo_type": "NORMAL",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "180.00",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "support_trading_session": "CORE",
                  "entrust_type": "QTY"
                }
              ]
            }
          },
          "Single Option": {
            "summary": "Single Leg Order",
            "description": "Single Option",
            "value": {
              "account_id": "<your_account_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id>",
                  "combo_type": "NORMAL",
                  "order_type": "LIMIT",
                  "limit_price": "11.25",
                  "quantity": "1",
                  "option_strategy": "SINGLE",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "instrument_type": "OPTION",
                  "market": "US",
                  "symbol": "AAPL",
                  "legs": [
                    {
                      "side": "BUY",
                      "quantity": "1",
                      "symbol": "AAPL",
                      "strike_price": "220.00",
                      "option_expire_date": "2026-12-18",
                      "instrument_type": "OPTION",
                      "option_type": "CALL",
                      "market": "US"
                    }
                  ]
                }
              ]
            }
          },
          "Buy-to-Open TP/SL(Equity)": {
            "summary": "Open Position with Take-Profit and Stop-Loss (MASTER + STOP_PROFIT + STOP_LOSS)",
            "description": "Buy-to-Open TP/SL(Equity)",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "MASTER",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "limit_price": "309.53"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "STOP_PROFIT",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "limit_price": "378.33"
                },
                {
                  "client_order_id": "<unique_id_3>",
                  "combo_type": "STOP_LOSS",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "STOP_LOSS",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "stop_price": "292.34"
                }
              ]
            }
          },
          "Sell-to-Close TP/SL(Equity)": {
            "summary": "Sell to Close an Existing Position with Take-Profit and Stop-Loss (no MASTER order)",
            "description": "Sell-to-Close TP/SL(Equity)",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "STOP_PROFIT",
                  "symbol": "TSLA",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "limit_price": "326.59"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "STOP_LOSS",
                  "symbol": "TSLA",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "STOP_LOSS",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "stop_price": "252.36"
                }
              ]
            }
          },
          "OTO": {
            "summary": "One-Triggers-the-Other: Master Order Triggers a Single Follow-up Order",
            "description": "OTO",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "MASTER",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "180.00",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "support_trading_session": "CORE",
                  "entrust_type": "QTY"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "OTO",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "200.00",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE"
                }
              ]
            }
          },
          "OCO": {
            "summary": "One-Cancels-the-Other: Two Sub-orders on an Existing Position, Filling One Cancels the Other",
            "description": "OCO",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "market": "US",
                  "instrument_type": "EQUITY",
                  "symbol": "AAPL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "combo_type": "OCO",
                  "order_type": "LIMIT",
                  "side": "BUY",
                  "quantity": "1",
                  "limit_price": "360.00",
                  "client_order_id": "<unique_id_1>"
                },
                {
                  "market": "US",
                  "instrument_type": "EQUITY",
                  "symbol": "AAPL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "combo_type": "OCO",
                  "order_type": "LIMIT",
                  "side": "BUY",
                  "quantity": "1",
                  "limit_price": "370.80",
                  "client_order_id": "<unique_id_2>"
                }
              ]
            }
          },
          "OTOCO": {
            "summary": "One-Triggers-a-One-Cancels-the-Other: Master Order Triggers an OCO Order Set",
            "description": "OTOCO",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "MASTER",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "180.00",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "support_trading_session": "CORE",
                  "entrust_type": "QTY"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "OTOCO",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "200.00",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE"
                },
                {
                  "client_order_id": "<unique_id_3>",
                  "combo_type": "OTOCO",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "STOP_LOSS",
                  "stop_price": "170.00",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE"
                }
              ]
            }
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
            "required": [
              "estimated_cost",
              "estimated_transaction_fee"
            ],
            "type": "object",
            "properties": {
              "estimated_cost": {
                "type": "string",
                "description": "Estimated capital required for the order. The meaning varies by product type:<br/>The actual fee may differ based on final execution.",
                "example": "100"
              },
              "estimated_transaction_fee": {
                "type": "string",
                "description": "Estimated transaction fee for placing the order, including exchange, clearing, and commission fees. <br/> The actual fee may differ based on final execution.",
                "example": "1"
              },
              "estimated_transaction_fee_detail": {
                "type": "object",
                "properties": {
                  "commission": {
                    "type": "object",
                    "properties": {
                      "actual_commission": {
                        "type": "string",
                        "description": "Actual commission collected",
                        "example": "1.0"
                      },
                      "receivable_commission": {
                        "type": "string",
                        "description": "Receivable commission",
                        "example": "1.0"
                      }
                    },
                    "description": "Commission breakdown",
                    "title": "CommonCommissionResultVO"
                  },
                  "fees": {
                    "type": "array",
                    "description": "Fee breakdown",
                    "items": {
                      "type": "object",
                      "properties": {
                        "type": {
                          "type": "string",
                          "description": "Fee type",
                          "example": "FINRA_CAT_REGULATORY_FEE"
                        },
                        "actual_value": {
                          "type": "string",
                          "description": "Actual fee collected",
                          "example": "1.0"
                        },
                        "receivable_value": {
                          "type": "string",
                          "description": "Receivable fee",
                          "example": "1.0"
                        }
                      },
                      "description": "Fee breakdown",
                      "title": "CommonFeeResultVO"
                    }
                  }
                },
                "description": "Breakdown of the estimated transaction fee, including commission and itemized fees.",
                "title": "CommonFeeDetailResultVO"
              }
            },
            "title": "OrderCommonPreviewResult"
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
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "error code",
                "example": "OPENAPI_NO_NIGHT_TRADING_TIME"
              },
              "message": {
                "type": "string",
                "description": "error message",
                "example": "The current period does not support placing night orders"
              }
            },
            "description": "Business Response",
            "title": "BizResponse"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "client_combo_order_id": "0KGOHL4PR2SLC0DKIND4TI0001",
    "new_orders": [
      {
        "combo_type": "NORMAL",
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "instrument_type": "EQUITY",
        "market": "US",
        "symbol": "AAPL",
        "order_type": "MARKET",
        "entrust_type": "QTY",
        "support_trading_session": "CORE",
        "time_in_force": "DAY",
        "side": "BUY",
        "quantity": "1",
        "limit_price": "11.0",
        "stop_price": "11.0",
        "option_strategy": "SINGLE",
        "legs": [
          {
            "instrument_type": "OPTION",
            "market": "US",
            "symbol": "AAPL",
            "side": "BUY",
            "strike_price": "11.0",
            "option_expire_date": "2025-08-01",
            "option_type": "CALL",
            "quantity": "1"
          }
        ]
      }
    ]
  },
  "postman": {
    "name": "Preview Order",
    "description": {
      "content": "Calculates the estimated amount and cost based on the provided information. Supports simple orders.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "preview"
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

## Order Place

> Source: <https://developer.webull.hk/apis/docs/reference/common-order-place.md>

### Place Order

Places equity and options orders. The A-Share trading function is disabled by default.

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
  "path": "/trading/orders/place",
  "method": "post",
  "tags": [
    "Trading"
  ],
  "description": "Places equity and options orders. The A-Share trading function is disabled by default.",
  "operationId": "Common Order Place",
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
    "description": "Order Place",
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "account_id",
            "new_orders"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "client_combo_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the combined order<br/> If combo_type = NORMAL and client_combo_order_id not need to set<br/> If combo_type != NORMAL and client_combo_order_id not provided <br/> the server will automatically generate one.<br/> To sell and close an existing position with take-profit/stop-loss, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) grouped under the same client_combo_order_id; no MASTER order is required in this scenario.",
              "example": "0KGOHL4PR2SLC0DKIND4TI0001"
            },
            "new_orders": {
              "type": "array",
              "description": "Order Details",
              "items": {
                "required": [
                  "client_order_id",
                  "combo_type",
                  "entrust_type",
                  "instrument_type",
                  "market",
                  "order_type",
                  "side",
                  "symbol",
                  "time_in_force"
                ],
                "type": "object",
                "properties": {
                  "combo_type": {
                    "type": "string",
                    "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> Futures Trading currently support only the NORMAL type.<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) <br/> Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened.<br/> Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO.<br/> Sub-order quantity limits by combo_type:\n\n| Scenario | combo_type | Supported order_type | Quantity | Description |\n|---|---|---|---|---|\n| Take-Profit/Stop-Loss | MASTER | MARKET, LIMIT | 1 | Master Order |\n| Take-Profit/Stop-Loss | STOP_PROFIT | LIMIT | 0-1 | Take Profit Order |\n| Take-Profit/Stop-Loss | STOP_LOSS | STOP_LOSS | 0-1 | Stop Loss Order |\n| OTO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTO | OTO | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | Triggered Order(s) |\n| OCO | OCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 2-6 | Mutually Cancelling Orders |\n| OTOCO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTOCO | OTOCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | OCO Order Set Triggered by MASTER |",
                    "example": "NORMAL"
                  },
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY",
                      "OPTION",
                      "FUTURES"
                    ]
                  },
                  "market": {
                    "type": "string",
                    "description": "Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                    "example": "US",
                    "enum": [
                      "US",
                      "HK",
                      "CN"
                    ]
                  },
                  "symbol": {
                    "type": "string",
                    "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                    "example": "BULL"
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order <br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order <br/>China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
                    "example": "MARKET",
                    "enum": [
                      "LIMIT",
                      "MARKET",
                      "STOP_LOSS",
                      "STOP_LOSS_LIMIT",
                      "ENHANCED_LIMIT",
                      "AT_AUCTION",
                      "AT_AUCTION_LIMIT",
                      "MARKET_ON_OPEN",
                      "TRAILING_STOP_LOSS",
                      "TRAILING_STOP_LOSS_LIMIT",
                      "TOUCH_MKT",
                      "TOUCH_LMT",
                      "ODD_LOT_LIMIT"
                    ]
                  },
                  "entrust_type": {
                    "type": "string",
                    "description": "Specifies the method for placing the order. <br/> &bull; QTY: Order specified by quantity of shares or units. <br/> &bull; AMOUNT: Order specified by total cash amount, applicable for fractional share trading of US stocks.",
                    "example": "QTY",
                    "enum": [
                      "QTY",
                      "AMOUNT"
                    ]
                  },
                  "support_trading_session": {
                    "type": "string",
                    "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only. <br/> Deprecated values: <br/> &bull; Y: [Deprecated]Include extended trading hours. <br/> &bull; N: [Deprecated]Only support regular trading hours. <br/> Active values: <br/> &bull; NIGHT: Only supports night trading. <br/> &bull; ALL: Include extended trading hours. <br/> &bull; CORE: Only support regular trading hours. <br/> &bull; ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day)",
                    "example": "CORE",
                    "enum": [
                      "Y",
                      "N",
                      "NIGHT",
                      "ALL",
                      "CORE",
                      "ALL_DAY"
                    ]
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTD",
                      "GTC"
                    ]
                  },
                  "side": {
                    "type": "string",
                    "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                    "example": "BUY",
                    "enum": [
                      "BUY",
                      "SELL",
                      "SHORT"
                    ]
                  },
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "total_cash_amount": {
                    "type": "string",
                    "description": "The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount.",
                    "example": "100.4"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                    "example": "11.0"
                  },
                  "stop_price": {
                    "type": "string",
                    "description": "Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT.<br/> Specifies the trigger price at which the stop order becomes active.",
                    "example": "11.0"
                  },
                  "trailing_type": {
                    "type": "string",
                    "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.<br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                    "example": "AMOUNT",
                    "enum": [
                      "AMOUNT",
                      "PERCENTAGE"
                    ]
                  },
                  "trailing_stop_step": {
                    "type": "string",
                    "description": "Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.",
                    "example": "1"
                  },
                  "trailing_limit_price_offset": {
                    "type": "string",
                    "description": "The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT.<br/>When triggered, the limit order price is calculated as:<br/>&bull; Buy: limit price = stop price + trailing_limit_price_offset<br/>&bull; Sell: limit price = stop price &minus; trailing_limit_price_offset<br/>If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders.",
                    "example": "11.0"
                  },
                  "trigger_price_type": {
                    "type": "string",
                    "description": "Trigger price type of the order<br/> &bull; PRICE: Latest transaction price. <br/> &bull; PRICE_BID: Buy at one price. <br/> &bull; PRICE_ASK: Sell at one price. ",
                    "example": "PRICE",
                    "enum": [
                      "PRICE",
                      "PRICE_BID",
                      "PRICE_ASK"
                    ]
                  },
                  "sender_sub_id": {
                    "type": "string",
                    "description": "Identifier for the firm or sub-account in third-party transactions.<br/> For brokers, this field should contain the UUID of the broker user. Used to distinguish different entities or users within the same firm."
                  },
                  "no_party_ids": {
                    "type": "array",
                    "description": "List of party identifiers. Applicable only for Hong Kong stock orders. Required for Relevant Regulated Intermediaries; should be omitted otherwise.",
                    "items": {
                      "required": [
                        "party_id",
                        "party_id_source",
                        "party_role"
                      ],
                      "type": "object",
                      "properties": {
                        "party_id": {
                          "type": "string",
                          "description": "ID of the broker client submitting the order. Format examples:<br/> &bull; CE Number Format: ABC123<br/> &bull; BCAN Format: 2568<br/> Combined format example: ABC123.2568<br/> Must be certified by the BCAN system of the Hong Kong Stock Exchange.",
                          "example": "ABC123.2568"
                        },
                        "party_id_source": {
                          "type": "string",
                          "description": "Source of the party ID. Value must be \"D\" (Proprietary/Custom Code).",
                          "example": "D"
                        },
                        "party_role": {
                          "type": "string",
                          "description": "Role of the party. Value must be \"3\" (Client ID, BCAN Field).",
                          "example": "3"
                        }
                      },
                      "description": "List of party identifiers. Applicable only for Hong Kong stock orders. Required for Relevant Regulated Intermediaries; should be omitted otherwise.",
                      "title": "PartyId"
                    }
                  },
                  "expire_date": {
                    "type": "string",
                    "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                    "example": "2026-01-01"
                  },
                  "option_strategy": {
                    "type": "string",
                    "description": "Type of options strategy <br/> &bull; SINGLE: Indicates a single-leg options order",
                    "example": "SINGLE",
                    "enum": [
                      "SINGLE"
                    ]
                  },
                  "legs": {
                    "type": "array",
                    "description": "Option leg detail. Only required when placing option orders.",
                    "items": {
                      "required": [
                        "instrument_type",
                        "market",
                        "side",
                        "symbol"
                      ],
                      "type": "object",
                      "properties": {
                        "instrument_type": {
                          "type": "string",
                          "description": "Type of financial instrument associated with the request.",
                          "example": "OPTION",
                          "enum": [
                            "EQUITY",
                            "OPTION",
                            "FUTURES"
                          ]
                        },
                        "market": {
                          "type": "string",
                          "description": "Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                          "example": "US",
                          "enum": [
                            "US",
                            "HK",
                            "CN"
                          ]
                        },
                        "symbol": {
                          "type": "string",
                          "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                          "example": "AAPL"
                        },
                        "side": {
                          "type": "string",
                          "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                          "example": "BUY",
                          "enum": [
                            "BUY",
                            "SELL",
                            "SHORT"
                          ]
                        },
                        "strike_price": {
                          "type": "string",
                          "description": "Exercise price (strike price) of the option. <br/> Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise.",
                          "example": "11.0"
                        },
                        "option_expire_date": {
                          "type": "string",
                          "description": "Expiration date. Format: yyyy-MM-dd",
                          "example": "2025-08-01"
                        },
                        "option_type": {
                          "type": "string",
                          "description": "Type of the option. <br/> &bull; CALL: Right to buy the underlying asset. <br/> &bull; PUT: Right to sell the underlying asset.",
                          "example": "CALL",
                          "enum": [
                            "CALL",
                            "PUT"
                          ]
                        },
                        "quantity": {
                          "type": "string",
                          "description": "Quantity of the order or strategy leg. <br/>For stock legs, specifies the number of shares to transact <br/>For option legs, specifies the number of option contracts to transact for this leg and is expressed in whole contracts.",
                          "example": "1"
                        }
                      },
                      "description": "Option leg detail. Only required when placing option orders.",
                      "title": "OptionCommonPlaceLegParam"
                    }
                  }
                },
                "description": "Order Details",
                "title": "OrderCommonPlaceItemParam"
              }
            }
          },
          "title": "OrderCommonPlaceParam"
        },
        "examples": {
          "Equity": {
            "summary": "Stock Limit Order",
            "description": "Equity",
            "value": {
              "account_id": "<your_account_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id>",
                  "combo_type": "NORMAL",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "180.00",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "support_trading_session": "CORE",
                  "entrust_type": "QTY"
                }
              ]
            }
          },
          "Single Option": {
            "summary": "Single Leg Order",
            "description": "Single Option",
            "value": {
              "account_id": "<your_account_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id>",
                  "combo_type": "NORMAL",
                  "order_type": "LIMIT",
                  "limit_price": "11.25",
                  "quantity": "1",
                  "option_strategy": "SINGLE",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "instrument_type": "OPTION",
                  "market": "US",
                  "symbol": "AAPL",
                  "legs": [
                    {
                      "side": "BUY",
                      "quantity": "1",
                      "symbol": "AAPL",
                      "strike_price": "220.00",
                      "option_expire_date": "2026-12-18",
                      "instrument_type": "OPTION",
                      "option_type": "CALL",
                      "market": "US"
                    }
                  ]
                }
              ]
            }
          },
          "Buy-to-Open TP/SL(Equity)": {
            "summary": "Open Position with Take-Profit and Stop-Loss (MASTER + STOP_PROFIT + STOP_LOSS)",
            "description": "Buy-to-Open TP/SL(Equity)",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "MASTER",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "limit_price": "309.53"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "STOP_PROFIT",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "limit_price": "378.33"
                },
                {
                  "client_order_id": "<unique_id_3>",
                  "combo_type": "STOP_LOSS",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "STOP_LOSS",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "stop_price": "292.34"
                }
              ]
            }
          },
          "Sell-to-Close TP/SL(Equity)": {
            "summary": "Sell to Close an Existing Position with Take-Profit and Stop-Loss (no MASTER order)",
            "description": "Sell-to-Close TP/SL(Equity)",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "STOP_PROFIT",
                  "symbol": "TSLA",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "limit_price": "326.59"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "STOP_LOSS",
                  "symbol": "TSLA",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "STOP_LOSS",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "stop_price": "252.36"
                }
              ]
            }
          },
          "OTO": {
            "summary": "One-Triggers-the-Other: Master Order Triggers a Single Follow-up Order",
            "description": "OTO",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "MASTER",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "180.00",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "support_trading_session": "CORE",
                  "entrust_type": "QTY"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "OTO",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "200.00",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE"
                }
              ]
            }
          },
          "OCO": {
            "summary": "One-Cancels-the-Other: Two Sub-orders on an Existing Position, Filling One Cancels the Other",
            "description": "OCO",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "market": "US",
                  "instrument_type": "EQUITY",
                  "symbol": "AAPL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "combo_type": "OCO",
                  "order_type": "LIMIT",
                  "side": "BUY",
                  "quantity": "1",
                  "limit_price": "360.00",
                  "client_order_id": "<unique_id_1>"
                },
                {
                  "market": "US",
                  "instrument_type": "EQUITY",
                  "symbol": "AAPL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE",
                  "combo_type": "OCO",
                  "order_type": "LIMIT",
                  "side": "BUY",
                  "quantity": "1",
                  "limit_price": "370.80",
                  "client_order_id": "<unique_id_2>"
                }
              ]
            }
          },
          "OTOCO": {
            "summary": "One-Triggers-a-One-Cancels-the-Other: Master Order Triggers an OCO Order Set",
            "description": "OTOCO",
            "value": {
              "account_id": "<your_account_id>",
              "client_combo_order_id": "<unique_combo_id>",
              "new_orders": [
                {
                  "client_order_id": "<unique_id_1>",
                  "combo_type": "MASTER",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "180.00",
                  "quantity": "10",
                  "side": "BUY",
                  "time_in_force": "DAY",
                  "support_trading_session": "CORE",
                  "entrust_type": "QTY"
                },
                {
                  "client_order_id": "<unique_id_2>",
                  "combo_type": "OTOCO",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "LIMIT",
                  "limit_price": "200.00",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE"
                },
                {
                  "client_order_id": "<unique_id_3>",
                  "combo_type": "OTOCO",
                  "symbol": "AAPL",
                  "instrument_type": "EQUITY",
                  "market": "US",
                  "order_type": "STOP_LOSS",
                  "stop_price": "170.00",
                  "quantity": "10",
                  "side": "SELL",
                  "time_in_force": "DAY",
                  "entrust_type": "QTY",
                  "support_trading_session": "CORE"
                }
              ]
            }
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
              "client_order_id": {
                "type": "string",
                "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order for NORMAL order.",
                "example": "0KGOHL4PR2SLC0DKIND4TI0002"
              },
              "order_id": {
                "type": "string",
                "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order.",
                "example": "80HG7CPSFDPCAL3TP66LKBAS69"
              }
            },
            "title": "OrderCommonWriteResult"
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
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "error code",
                "example": "OPENAPI_NO_NIGHT_TRADING_TIME"
              },
              "message": {
                "type": "string",
                "description": "error message",
                "example": "The current period does not support placing night orders"
              }
            },
            "description": "Business Response",
            "title": "BizResponse"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "client_combo_order_id": "0KGOHL4PR2SLC0DKIND4TI0001",
    "new_orders": [
      {
        "combo_type": "NORMAL",
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "instrument_type": "EQUITY",
        "market": "US",
        "symbol": "BULL",
        "order_type": "MARKET",
        "entrust_type": "QTY",
        "support_trading_session": "CORE",
        "time_in_force": "DAY",
        "side": "BUY",
        "quantity": "1",
        "total_cash_amount": "100.4",
        "limit_price": "11.0",
        "stop_price": "11.0",
        "trailing_type": "AMOUNT",
        "trailing_stop_step": "1",
        "trailing_limit_price_offset": "11.0",
        "trigger_price_type": "PRICE",
        "sender_sub_id": "string",
        "no_party_ids": [
          {
            "party_id": "ABC123.2568",
            "party_id_source": "D",
            "party_role": "3"
          }
        ],
        "expire_date": "2026-01-01",
        "option_strategy": "SINGLE",
        "legs": [
          {
            "instrument_type": "OPTION",
            "market": "US",
            "symbol": "AAPL",
            "side": "BUY",
            "strike_price": "11.0",
            "option_expire_date": "2025-08-01",
            "option_type": "CALL",
            "quantity": "1"
          }
        ]
      }
    ]
  },
  "postman": {
    "name": "Place Order",
    "description": {
      "content": "Places equity and options orders. The A-Share trading function is disabled by default.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "place"
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

## Batch Place Orders

> Source: <https://developer.webull.com/apis/docs/reference/order-batch-place.md>

### Batch Place Orders

Places multiple orders in a single request.<br/>A maximum of 50 orders can be submitted once, Currently only stocks are supported. This service is not currently available to all clients. Please contact Webull if you require assistance.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data, and OAuth authentication for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.com/help",
      "email": "api-support@webull-us.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.com"
    }
  ],
  "path": "/trading/orders/batch-place",
  "method": "post",
  "tags": [
    "Trading"
  ],
  "description": "Places multiple orders in a single request.<br/>A maximum of 50 orders can be submitted once, Currently only stocks are supported. This service is not currently available to all clients. Please contact Webull if you require assistance.",
  "operationId": "Order Batch Place",
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
            "account_id",
            "batch_orders"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "batch_orders": {
              "type": "array",
              "description": "Batch Orders",
              "items": {
                "required": [
                  "client_order_id",
                  "combo_type",
                  "entrust_type",
                  "instrument_type",
                  "market",
                  "order_type",
                  "quantity",
                  "side",
                  "support_trading_session",
                  "symbol",
                  "time_in_force"
                ],
                "type": "object",
                "properties": {
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "combo_type": {
                    "type": "string",
                    "description": "Type of order combination. Currently only NORMAL is supported. It may be expanded to support other types in the future<br/> &bull; NORMAL: Indicates a standard single order",
                    "example": "NORMAL"
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY"
                    ]
                  },
                  "entrust_type": {
                    "type": "string",
                    "description": "Specifies the method for placing the order.<br/> &bull; QTY: Order specified by quantity of shares or units.",
                    "example": "QTY",
                    "enum": [
                      "QTY"
                    ]
                  },
                  "support_trading_session": {
                    "type": "string",
                    "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only.<br/> Algorithmic trading order currently supports only regular trading hours.<br/> &bull; NIGHT: Only supports night trading.<br/> &bull; ALL: Include extended trading hours.<br/> &bull; CORE: Only support regular trading hours.",
                    "example": "CORE",
                    "enum": [
                      "ALL",
                      "CORE",
                      "NIGHT"
                    ]
                  },
                  "symbol": {
                    "type": "string",
                    "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities).",
                    "example": "BULL"
                  },
                  "market": {
                    "type": "string",
                    "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                    "example": "US",
                    "enum": [
                      "US"
                    ]
                  },
                  "side": {
                    "type": "string",
                    "description": "The order side indicating the intended trading direction of the transaction. <br/> The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                    "example": "BUY",
                    "enum": [
                      "BUY",
                      "SELL"
                    ]
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>",
                    "example": "MARKET",
                    "enum": [
                      "MARKET",
                      "LIMIT"
                    ]
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.",
                    "example": "DAY",
                    "enum": [
                      "DAY"
                    ]
                  },
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.<br/> When event_trade_mode is set for an event contract trade, limit_price is not required.",
                    "example": "11.0"
                  }
                },
                "description": "Batch Orders",
                "title": "OrderCommonBatchPlaceItemParam"
              }
            }
          },
          "description": "Order Place Request Body",
          "title": "OrderCommonBatchPlaceParam"
        }
      }
    }
  },
  "responses": {
    "200": {
      "description": "Request successful",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "batch_orders",
              "failed",
              "success",
              "total"
            ],
            "type": "object",
            "properties": {
              "total": {
                "type": "integer",
                "description": "The total number of orders submitted each time",
                "format": "int32"
              },
              "success": {
                "type": "integer",
                "description": "The number of orders successfully submitted to the webull system",
                "format": "int32"
              },
              "failed": {
                "type": "integer",
                "description": "The number of failed order submitted to the webull system",
                "format": "int32"
              },
              "batch_orders": {
                "type": "array",
                "description": "Batch Order place result",
                "items": {
                  "type": "object",
                  "properties": {
                    "client_order_id": {
                      "type": "string",
                      "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order for NORMAL order.",
                      "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                    },
                    "order_id": {
                      "type": "string",
                      "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order.",
                      "example": "80HG7CPSFDPCAL3TP66LKBAS69"
                    },
                    "error_code": {
                      "type": "string",
                      "description": "Order place failed code",
                      "example": "OPENAPI_NO_TRADING_TIME"
                    },
                    "message": {
                      "type": "string",
                      "description": "Order place failed and detail failed reason.",
                      "example": "Non-trading time."
                    }
                  },
                  "description": "Batch Order place result",
                  "title": "OrderResp"
                }
              }
            },
            "title": "BatchOrderResp"
          },
          "examples": {
            " Special Instructions": {
              "description": " Special Instructions",
              "value": {
                "total": 2,
                "success": 1,
                "failed": 1,
                "batch_orders": [
                  {
                    "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0001",
                    "order_id": "80HG7CPSFDPCAL3TP66LKBAS69"
                  },
                  {
                    "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
                    "error_code": "OAUTH_OPENAPI_NO_TRADING_TIME",
                    "message": "Non-trading time."
                  }
                ]
              }
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
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "error code",
                "example": "OPENAPI_NO_NIGHT_TRADING_TIME"
              },
              "message": {
                "type": "string",
                "description": "error message",
                "example": "The current period does not support placing night orders"
              }
            },
            "description": "Business Response",
            "title": "BizResponse"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "batch_orders": [
      {
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "combo_type": "NORMAL",
        "instrument_type": "EQUITY",
        "entrust_type": "QTY",
        "support_trading_session": "CORE",
        "symbol": "BULL",
        "market": "US",
        "side": "BUY",
        "order_type": "MARKET",
        "time_in_force": "DAY",
        "quantity": "1",
        "limit_price": "11.0"
      }
    ]
  },
  "postman": {
    "name": "Batch Place Orders",
    "description": {
      "content": "Places multiple orders in a single request.<br/>A maximum of 50 orders can be submitted once, Currently only stocks are supported. This service is not currently available to all clients. Please contact Webull if you require assistance.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "batch-place"
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

## Order Replace

> Source: <https://developer.webull.hk/apis/docs/reference/common-order-replace.md>

### Replace Order

Modifies equity and options orders.

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
  "path": "/trading/orders/replace",
  "method": "post",
  "tags": [
    "Trading"
  ],
  "description": "Modifies equity and options orders.",
  "operationId": "Common Order Replace",
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
            "account_id",
            "modify_orders"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "modify_orders": {
              "type": "array",
              "description": "Order Details",
              "items": {
                "required": [
                  "client_order_id"
                ],
                "type": "object",
                "properties": {
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTD",
                      "GTC"
                    ]
                  },
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "expire_date": {
                    "type": "string",
                    "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                    "example": "2026-01-01"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                    "example": "11.0"
                  },
                  "stop_price": {
                    "type": "string",
                    "description": "Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT.<br/> Specifies the trigger price at which the stop order becomes active.",
                    "example": "11.0"
                  },
                  "trailing_type": {
                    "type": "string",
                    "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.<br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                    "example": "AMOUNT",
                    "enum": [
                      "AMOUNT",
                      "PERCENTAGE"
                    ]
                  },
                  "trailing_stop_step": {
                    "type": "string",
                    "description": "Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.",
                    "example": "1"
                  },
                  "trailing_limit_price_offset": {
                    "type": "string",
                    "description": "The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT.<br/>When triggered, the limit order price is calculated as:<br/>&bull; Buy: limit price = stop price + trailing_limit_price_offset<br/>&bull; Sell: limit price = stop price &minus; trailing_limit_price_offset<br/>If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders.",
                    "example": "11.0"
                  },
                  "trigger_price_type": {
                    "type": "string",
                    "description": "Trigger price type of the order<br/> &bull; PRICE: Latest transaction price. <br/> &bull; PRICE_BID: Buy at one price. <br/> &bull; PRICE_ASK: Sell at one price. ",
                    "example": "PRICE",
                    "enum": [
                      "PRICE",
                      "PRICE_BID",
                      "PRICE_ASK"
                    ]
                  }
                },
                "description": "Order Details",
                "title": "OrderCommonReplaceItemParam"
              }
            }
          },
          "title": "OrderCommonReplaceParam"
        }
      }
    }
  },
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "client_order_id": {
                "type": "string",
                "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order for NORMAL order.",
                "example": "0KGOHL4PR2SLC0DKIND4TI0002"
              },
              "order_id": {
                "type": "string",
                "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order.",
                "example": "80HG7CPSFDPCAL3TP66LKBAS69"
              }
            },
            "title": "OrderCommonWriteResult"
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
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "error code",
                "example": "OPENAPI_NO_NIGHT_TRADING_TIME"
              },
              "message": {
                "type": "string",
                "description": "error message",
                "example": "The current period does not support placing night orders"
              }
            },
            "description": "Business Response",
            "title": "BizResponse"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "modify_orders": [
      {
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "time_in_force": "DAY",
        "quantity": "1",
        "expire_date": "2026-01-01",
        "limit_price": "11.0",
        "stop_price": "11.0",
        "trailing_type": "AMOUNT",
        "trailing_stop_step": "1",
        "trailing_limit_price_offset": "11.0",
        "trigger_price_type": "PRICE"
      }
    ]
  },
  "postman": {
    "name": "Replace Order",
    "description": {
      "content": "Modifies equity and options orders.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "replace"
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

## Order Cancel

> Source: <https://developer.webull.hk/apis/docs/reference/common-order-cancel.md>

### Cancel Order

Cancels orders for equities and options.

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
  "path": "/trading/orders/cancel",
  "method": "post",
  "tags": [
    "Trading"
  ],
  "description": "Cancels orders for equities and options.",
  "operationId": "Common Order Cancel",
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
            "account_id",
            "client_order_id"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "client_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).<br/> Used to track or reference the order when interacting with the system.",
              "example": "0KGOHL4PR2SLC0DKIND4TI0002"
            }
          },
          "title": "OrderCommonCancelParam"
        }
      }
    }
  },
  "responses": {
    "200": {
      "description": "OK",
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "client_order_id": {
                "type": "string",
                "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order for NORMAL order.",
                "example": "0KGOHL4PR2SLC0DKIND4TI0002"
              },
              "order_id": {
                "type": "string",
                "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order.",
                "example": "80HG7CPSFDPCAL3TP66LKBAS69"
              }
            },
            "title": "OrderCommonWriteResult"
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
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "error_code": {
                "type": "string",
                "description": "error code",
                "example": "OPENAPI_NO_NIGHT_TRADING_TIME"
              },
              "message": {
                "type": "string",
                "description": "error message",
                "example": "The current period does not support placing night orders"
              }
            },
            "description": "Business Response",
            "title": "BizResponse"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002"
  },
  "postman": {
    "name": "Cancel Order",
    "description": {
      "content": "Cancels orders for equities and options.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "cancel"
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

## Open Orders

> Source: <https://developer.webull.hk/apis/docs/reference/order-open.md>

### List Open Orders

Retrieves pending orders by page. Orders can be modified or cancelled based on client_order_id.

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
  "path": "/trading/orders/open-orders/list",
  "method": "get",
  "tags": [
    "Order Query"
  ],
  "description": "Retrieves pending orders by page. Orders can be modified or cancelled based on client_order_id.",
  "operationId": "orderOpen",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": 20150320010101000
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
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
                  "required": [
                    "client_order_id",
                    "combo_type",
                    "orders"
                  ],
                  "type": "object",
                  "properties": {
                    "client_order_id": {
                      "type": "string",
                      "description": "Client-defined order identifier. Returned in the response for simple orders. <br/> Represents the unique order ID assigned by the user when placing the order.",
                      "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                    },
                    "combo_type": {
                      "type": "string",
                      "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> Futures Trading currently support only the NORMAL type.<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) <br/> Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened.<br/> Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO.<br/> Sub-order quantity limits by combo_type:\n\n| Scenario | combo_type | Supported order_type | Quantity | Description |\n|---|---|---|---|---|\n| Take-Profit/Stop-Loss | MASTER | MARKET, LIMIT | 1 | Master Order |\n| Take-Profit/Stop-Loss | STOP_PROFIT | LIMIT | 0-1 | Take Profit Order |\n| Take-Profit/Stop-Loss | STOP_LOSS | STOP_LOSS | 0-1 | Stop Loss Order |\n| OTO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTO | OTO | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | Triggered Order(s) |\n| OCO | OCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 2-6 | Mutually Cancelling Orders |\n| OTOCO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTOCO | OTOCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | OCO Order Set Triggered by MASTER |",
                      "example": "NORMAL",
                      "enum": [
                        "NORMAL",
                        "MASTER",
                        "STOP_PROFIT",
                        "STOP_LOSS",
                        "OTO",
                        "OCO",
                        "OTOCO"
                      ]
                    },
                    "orders": {
                      "type": "array",
                      "description": "Order Details",
                      "items": {
                        "required": [
                          "client_order_id",
                          "order_id",
                          "order_type",
                          "place_time_at",
                          "side",
                          "status",
                          "symbol",
                          "time_in_force",
                          "total_quantity"
                        ],
                        "type": "object",
                        "properties": {
                          "client_order_id": {
                            "type": "string",
                            "description": "Client-defined order identifier. Returned in the response for simple orders. <br/> Represents the unique order ID assigned by the user when placing the order.",
                            "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                          },
                          "order_id": {
                            "type": "string",
                            "description": "System-generated order identifier. Returned in the response for simple orders. <br/> Represents the unique Webull order ID assigned by the system.",
                            "example": "0352U72LQI6DT0KF41GK000000"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                            "example": "AAPL"
                          },
                          "side": {
                            "type": "string",
                            "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                            "example": "BUY",
                            "enum": [
                              "BUY",
                              "SELL",
                              "SHORT"
                            ]
                          },
                          "status": {
                            "type": "string",
                            "description": "&bull; PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion <br/> &bull; SUBMITTED: Indicates that the order has been submitted to the exchange and is awaiting completion<br/> <br/> &bull; CANCELLED: Indicates that the order has been successfully cancelled <br/> &bull; FILLED: Indicates that the order has been fully executed <br/> &bull; FAILED: Indicates a failed order, such as REJECTED <br/> &bull; PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed",
                            "example": "SUBMITTED",
                            "enum": [
                              "PENDING",
                              "SUBMITTED",
                              "CANCELLED",
                              "FILLED",
                              "FAILED",
                              "PARTIAL_FILLED"
                            ]
                          },
                          "order_type": {
                            "type": "string",
                            "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order <br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order <br/>China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
                            "example": "MARKET",
                            "enum": [
                              "LIMIT",
                              "MARKET",
                              "STOP_LOSS",
                              "STOP_LOSS_LIMIT",
                              "ENHANCED_LIMIT",
                              "AT_AUCTION",
                              "AT_AUCTION_LIMIT",
                              "MARKET_ON_OPEN",
                              "TRAILING_STOP_LOSS",
                              "TRAILING_STOP_LOSS_LIMIT",
                              "TOUCH_MKT",
                              "TOUCH_LMT",
                              "ODD_LOT_LIMIT"
                            ]
                          },
                          "instrument_type": {
                            "type": "string",
                            "description": "Type of financial instrument associated with the request.",
                            "example": "STOCK",
                            "enum": [
                              "EQUITY",
                              "OPTION",
                              "FUTURES"
                            ]
                          },
                          "support_trading_session": {
                            "type": "string",
                            "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only. <br/> Deprecated values: <br/> &bull; Y: [Deprecated]Include extended trading hours. <br/> &bull; N: [Deprecated]Only support regular trading hours. <br/> Active values: <br/> &bull; NIGHT: Only supports night trading. <br/> &bull; ALL: Include extended trading hours. <br/> &bull; CORE: Only support regular trading hours. <br/> &bull; ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day)",
                            "example": "CORE",
                            "enum": [
                              "Y",
                              "N",
                              "NIGHT",
                              "ALL",
                              "CORE",
                              "ALL_DAY"
                            ]
                          },
                          "time_in_force": {
                            "type": "string",
                            "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                            "example": "DAY",
                            "enum": [
                              "DAY",
                              "GTD",
                              "GTC"
                            ]
                          },
                          "total_quantity": {
                            "type": "string",
                            "description": "Total order quantity. Represents the total number of units submitted for this order.",
                            "example": "1"
                          },
                          "filled_quantity": {
                            "type": "string",
                            "description": "Quantity that has been executed. Represents the number of units that have been filled so far.",
                            "example": "1"
                          },
                          "filled_price": {
                            "type": "string",
                            "description": "Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null.",
                            "example": "11.0"
                          },
                          "limit_price": {
                            "type": "string",
                            "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                            "example": "11.0"
                          },
                          "stop_price": {
                            "type": "string",
                            "description": "Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT.<br/> Specifies the trigger price at which the stop order becomes active.",
                            "example": "11.0"
                          },
                          "trailing_type": {
                            "type": "string",
                            "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.<br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                            "example": "AMOUNT",
                            "enum": [
                              "AMOUNT",
                              "PERCENTAGE"
                            ]
                          },
                          "trailing_stop_step": {
                            "type": "string",
                            "description": "Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.",
                            "example": "1"
                          },
                          "trailing_limit_price_offset": {
                            "type": "string",
                            "description": "The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT.<br/>When triggered, the limit order price is calculated as:<br/>&bull; Buy: limit price = stop price + trailing_limit_price_offset<br/>&bull; Sell: limit price = stop price &minus; trailing_limit_price_offset<br/>If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders.",
                            "example": "11.0"
                          },
                          "trigger_price_type": {
                            "type": "string",
                            "description": "Trigger price type of the order<br/> &bull; PRICE: Latest transaction price. <br/> &bull; PRICE_BID: Buy at one price. <br/> &bull; PRICE_ASK: Sell at one price. ",
                            "example": "PRICE",
                            "enum": [
                              "PRICE",
                              "PRICE_BID",
                              "PRICE_ASK"
                            ]
                          },
                          "place_time": {
                            "type": "string",
                            "description": "Order placement time in milliseconds since Unix epoch.",
                            "example": "1726745361658",
                            "deprecated": true
                          },
                          "place_time_at": {
                            "type": "string",
                            "description": "Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                            "example": "2025-11-11T05:44:35.385Z"
                          },
                          "filled_time": {
                            "type": "string",
                            "description": "Time of the last executed trade in milliseconds since Unix epoch.",
                            "example": "1726745361871",
                            "deprecated": true
                          },
                          "filled_time_at": {
                            "type": "string",
                            "description": "Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                            "example": "2025-11-11T05:44:35.385Z"
                          },
                          "legs": {
                            "type": "array",
                            "description": "Leg detail",
                            "items": {
                              "required": [
                                "option_category",
                                "option_expire_date",
                                "option_type",
                                "quantity",
                                "side",
                                "strike_price",
                                "symbol"
                              ],
                              "type": "object",
                              "properties": {
                                "symbol": {
                                  "type": "string",
                                  "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                                  "example": "AAPL"
                                },
                                "side": {
                                  "type": "string",
                                  "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                                  "example": "BUY",
                                  "enum": [
                                    "BUY",
                                    "SELL",
                                    "SHORT"
                                  ]
                                },
                                "quantity": {
                                  "type": "string",
                                  "description": "Quantity of the order. Specifies the number of shares or units to transact.<br/>For US stocks, fractional quantities are allowed and can include decimals.",
                                  "example": "1"
                                },
                                "option_type": {
                                  "type": "string",
                                  "description": "Type of the option. <br/> &bull; CALL: Right to buy the underlying asset. <br/> &bull; PUT: Right to sell the underlying asset.",
                                  "example": "CALL",
                                  "enum": [
                                    "CALL",
                                    "PUT"
                                  ]
                                },
                                "option_category": {
                                  "type": "string",
                                  "description": "Category of the option, indicating its exercise style. <br/> Possible values: <br/> &bull; AMERICAN: Can be exercised any time before expiration. <br/> &bull; EUROPEAN: Can only be exercised at expiration. ",
                                  "example": "AMERICAN",
                                  "enum": [
                                    "AMERICAN",
                                    "EUROPEAN"
                                  ]
                                },
                                "option_strategy": {
                                  "type": "string",
                                  "description": "Type of options strategy <br/> &bull; SINGLE: Indicates a single-leg options order",
                                  "example": "SINGLE",
                                  "enum": [
                                    "SINGLE"
                                  ]
                                },
                                "strike_price": {
                                  "type": "string",
                                  "description": "Exercise price (strike price) of the option.<br/> Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise.",
                                  "example": "190.0"
                                },
                                "option_contract_multiplier": {
                                  "type": "string",
                                  "description": "The number of shares corresponding to each option contract",
                                  "example": "100"
                                },
                                "option_contract_deliverable": {
                                  "type": "string",
                                  "description": "The number of shares required to exercise each contract",
                                  "example": "100"
                                },
                                "option_expire_date": {
                                  "type": "string",
                                  "description": "Expiration date of the option.<br/> Format: yyyy-MM-dd. After this date, the option will no longer be valid.",
                                  "example": "2025-11-21"
                                }
                              },
                              "description": "Leg detail",
                              "title": "OrderListLeg"
                            }
                          }
                        },
                        "description": "Order Details",
                        "title": "OrderListItem"
                      }
                    }
                  },
                  "description": "Result data list",
                  "title": "OrderListResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoOrderListResult"
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
    "name": "List Open Orders",
    "description": {
      "content": "Retrieves pending orders by page. Orders can be modified or cancelled based on client_order_id.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "open-orders",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Account identifier.",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key from previous response for next page.",
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

## Order History

> Source: <https://developer.webull.hk/apis/docs/reference/order-history.md>

### List Order History

Retrieves historical orders for the past 7 days. If orders are group orders, they will be returned together.

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
  "path": "/trading/orders/historical-orders/list",
  "method": "get",
  "tags": [
    "Order Query"
  ],
  "description": "Retrieves historical orders for the past 7 days. If orders are group orders, they will be returned together.",
  "operationId": "orderHistory",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": 20150320010101000
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "The start date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 6 months.<br/> Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:59.012Z"
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "The end time of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:59.012Z"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
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
                  "required": [
                    "client_order_id",
                    "combo_type",
                    "orders"
                  ],
                  "type": "object",
                  "properties": {
                    "client_order_id": {
                      "type": "string",
                      "description": "Client-defined order identifier. Returned in the response for simple orders. <br/> Represents the unique order ID assigned by the user when placing the order.",
                      "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                    },
                    "combo_type": {
                      "type": "string",
                      "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> Futures Trading currently support only the NORMAL type.<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) <br/> Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened.<br/> Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO.<br/> Sub-order quantity limits by combo_type:\n\n| Scenario | combo_type | Supported order_type | Quantity | Description |\n|---|---|---|---|---|\n| Take-Profit/Stop-Loss | MASTER | MARKET, LIMIT | 1 | Master Order |\n| Take-Profit/Stop-Loss | STOP_PROFIT | LIMIT | 0-1 | Take Profit Order |\n| Take-Profit/Stop-Loss | STOP_LOSS | STOP_LOSS | 0-1 | Stop Loss Order |\n| OTO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTO | OTO | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | Triggered Order(s) |\n| OCO | OCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 2-6 | Mutually Cancelling Orders |\n| OTOCO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTOCO | OTOCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | OCO Order Set Triggered by MASTER |",
                      "example": "NORMAL",
                      "enum": [
                        "NORMAL",
                        "MASTER",
                        "STOP_PROFIT",
                        "STOP_LOSS",
                        "OTO",
                        "OCO",
                        "OTOCO"
                      ]
                    },
                    "orders": {
                      "type": "array",
                      "description": "Order Details",
                      "items": {
                        "required": [
                          "client_order_id",
                          "order_id",
                          "order_type",
                          "place_time_at",
                          "side",
                          "status",
                          "symbol",
                          "time_in_force",
                          "total_quantity"
                        ],
                        "type": "object",
                        "properties": {
                          "client_order_id": {
                            "type": "string",
                            "description": "Client-defined order identifier. Returned in the response for simple orders. <br/> Represents the unique order ID assigned by the user when placing the order.",
                            "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                          },
                          "order_id": {
                            "type": "string",
                            "description": "System-generated order identifier. Returned in the response for simple orders. <br/> Represents the unique Webull order ID assigned by the system.",
                            "example": "0352U72LQI6DT0KF41GK000000"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                            "example": "AAPL"
                          },
                          "side": {
                            "type": "string",
                            "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                            "example": "BUY",
                            "enum": [
                              "BUY",
                              "SELL",
                              "SHORT"
                            ]
                          },
                          "status": {
                            "type": "string",
                            "description": "&bull; PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion <br/> &bull; SUBMITTED: Indicates that the order has been submitted to the exchange and is awaiting completion<br/> <br/> &bull; CANCELLED: Indicates that the order has been successfully cancelled <br/> &bull; FILLED: Indicates that the order has been fully executed <br/> &bull; FAILED: Indicates a failed order, such as REJECTED <br/> &bull; PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed",
                            "example": "SUBMITTED",
                            "enum": [
                              "PENDING",
                              "SUBMITTED",
                              "CANCELLED",
                              "FILLED",
                              "FAILED",
                              "PARTIAL_FILLED"
                            ]
                          },
                          "order_type": {
                            "type": "string",
                            "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order <br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order <br/>China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
                            "example": "MARKET",
                            "enum": [
                              "LIMIT",
                              "MARKET",
                              "STOP_LOSS",
                              "STOP_LOSS_LIMIT",
                              "ENHANCED_LIMIT",
                              "AT_AUCTION",
                              "AT_AUCTION_LIMIT",
                              "MARKET_ON_OPEN",
                              "TRAILING_STOP_LOSS",
                              "TRAILING_STOP_LOSS_LIMIT",
                              "TOUCH_MKT",
                              "TOUCH_LMT",
                              "ODD_LOT_LIMIT"
                            ]
                          },
                          "instrument_type": {
                            "type": "string",
                            "description": "Type of financial instrument associated with the request.",
                            "example": "STOCK",
                            "enum": [
                              "EQUITY",
                              "OPTION",
                              "FUTURES"
                            ]
                          },
                          "support_trading_session": {
                            "type": "string",
                            "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only. <br/> Deprecated values: <br/> &bull; Y: [Deprecated]Include extended trading hours. <br/> &bull; N: [Deprecated]Only support regular trading hours. <br/> Active values: <br/> &bull; NIGHT: Only supports night trading. <br/> &bull; ALL: Include extended trading hours. <br/> &bull; CORE: Only support regular trading hours. <br/> &bull; ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day)",
                            "example": "CORE",
                            "enum": [
                              "Y",
                              "N",
                              "NIGHT",
                              "ALL",
                              "CORE",
                              "ALL_DAY"
                            ]
                          },
                          "time_in_force": {
                            "type": "string",
                            "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                            "example": "DAY",
                            "enum": [
                              "DAY",
                              "GTD",
                              "GTC"
                            ]
                          },
                          "total_quantity": {
                            "type": "string",
                            "description": "Total order quantity. Represents the total number of units submitted for this order.",
                            "example": "1"
                          },
                          "filled_quantity": {
                            "type": "string",
                            "description": "Quantity that has been executed. Represents the number of units that have been filled so far.",
                            "example": "1"
                          },
                          "filled_price": {
                            "type": "string",
                            "description": "Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null.",
                            "example": "11.0"
                          },
                          "limit_price": {
                            "type": "string",
                            "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                            "example": "11.0"
                          },
                          "stop_price": {
                            "type": "string",
                            "description": "Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT.<br/> Specifies the trigger price at which the stop order becomes active.",
                            "example": "11.0"
                          },
                          "trailing_type": {
                            "type": "string",
                            "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.<br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                            "example": "AMOUNT",
                            "enum": [
                              "AMOUNT",
                              "PERCENTAGE"
                            ]
                          },
                          "trailing_stop_step": {
                            "type": "string",
                            "description": "Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT.",
                            "example": "1"
                          },
                          "trailing_limit_price_offset": {
                            "type": "string",
                            "description": "The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT.<br/>When triggered, the limit order price is calculated as:<br/>&bull; Buy: limit price = stop price + trailing_limit_price_offset<br/>&bull; Sell: limit price = stop price &minus; trailing_limit_price_offset<br/>If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders.",
                            "example": "11.0"
                          },
                          "trigger_price_type": {
                            "type": "string",
                            "description": "Trigger price type of the order<br/> &bull; PRICE: Latest transaction price. <br/> &bull; PRICE_BID: Buy at one price. <br/> &bull; PRICE_ASK: Sell at one price. ",
                            "example": "PRICE",
                            "enum": [
                              "PRICE",
                              "PRICE_BID",
                              "PRICE_ASK"
                            ]
                          },
                          "place_time": {
                            "type": "string",
                            "description": "Order placement time in milliseconds since Unix epoch.",
                            "example": "1726745361658",
                            "deprecated": true
                          },
                          "place_time_at": {
                            "type": "string",
                            "description": "Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                            "example": "2025-11-11T05:44:35.385Z"
                          },
                          "filled_time": {
                            "type": "string",
                            "description": "Time of the last executed trade in milliseconds since Unix epoch.",
                            "example": "1726745361871",
                            "deprecated": true
                          },
                          "filled_time_at": {
                            "type": "string",
                            "description": "Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                            "example": "2025-11-11T05:44:35.385Z"
                          },
                          "legs": {
                            "type": "array",
                            "description": "Leg detail",
                            "items": {
                              "required": [
                                "option_category",
                                "option_expire_date",
                                "option_type",
                                "quantity",
                                "side",
                                "strike_price",
                                "symbol"
                              ],
                              "type": "object",
                              "properties": {
                                "symbol": {
                                  "type": "string",
                                  "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                                  "example": "AAPL"
                                },
                                "side": {
                                  "type": "string",
                                  "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                                  "example": "BUY",
                                  "enum": [
                                    "BUY",
                                    "SELL",
                                    "SHORT"
                                  ]
                                },
                                "quantity": {
                                  "type": "string",
                                  "description": "Quantity of the order. Specifies the number of shares or units to transact.<br/>For US stocks, fractional quantities are allowed and can include decimals.",
                                  "example": "1"
                                },
                                "option_type": {
                                  "type": "string",
                                  "description": "Type of the option. <br/> &bull; CALL: Right to buy the underlying asset. <br/> &bull; PUT: Right to sell the underlying asset.",
                                  "example": "CALL",
                                  "enum": [
                                    "CALL",
                                    "PUT"
                                  ]
                                },
                                "option_category": {
                                  "type": "string",
                                  "description": "Category of the option, indicating its exercise style. <br/> Possible values: <br/> &bull; AMERICAN: Can be exercised any time before expiration. <br/> &bull; EUROPEAN: Can only be exercised at expiration. ",
                                  "example": "AMERICAN",
                                  "enum": [
                                    "AMERICAN",
                                    "EUROPEAN"
                                  ]
                                },
                                "option_strategy": {
                                  "type": "string",
                                  "description": "Type of options strategy <br/> &bull; SINGLE: Indicates a single-leg options order",
                                  "example": "SINGLE",
                                  "enum": [
                                    "SINGLE"
                                  ]
                                },
                                "strike_price": {
                                  "type": "string",
                                  "description": "Exercise price (strike price) of the option.<br/> Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise.",
                                  "example": "190.0"
                                },
                                "option_contract_multiplier": {
                                  "type": "string",
                                  "description": "The number of shares corresponding to each option contract",
                                  "example": "100"
                                },
                                "option_contract_deliverable": {
                                  "type": "string",
                                  "description": "The number of shares required to exercise each contract",
                                  "example": "100"
                                },
                                "option_expire_date": {
                                  "type": "string",
                                  "description": "Expiration date of the option.<br/> Format: yyyy-MM-dd. After this date, the option will no longer be valid.",
                                  "example": "2025-11-21"
                                }
                              },
                              "description": "Leg detail",
                              "title": "OrderListLeg"
                            }
                          }
                        },
                        "description": "Order Details",
                        "title": "OrderListItem"
                      }
                    }
                  },
                  "description": "Result data list",
                  "title": "OrderListResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoOrderListResult"
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
    "name": "List Order History",
    "description": {
      "content": "Retrieves historical orders for the past 7 days. If orders are group orders, they will be returned together.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "historical-orders",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Account identifier.",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The start date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 6 months.<br/> Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'.",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The end time of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'.",
            "type": "text/plain"
          },
          "key": "end_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key from previous response for next page.",
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

## Order Detail

> Source: <https://developer.webull.hk/apis/docs/reference/order-detail.md>

### Get Order Detail

Retrieves the specified order details through the order ID.

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
  "path": "/trading/orders/get",
  "method": "get",
  "tags": [
    "Order Query"
  ],
  "description": "Retrieves the specified order details through the order ID.",
  "operationId": "orderDetail",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LOJOQITOD49R6G9BPQM489CISA"
    },
    {
      "name": "client_order_id",
      "in": "query",
      "description": "Unique client-defined identifier for the order. <br/>Maximum length is 32 characters and must be unique per account. <br/>Used to track or reference the order when interacting with the system.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "0KGOHL4PR2SLC0DKIND4TI0002"
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
              "client_order_id",
              "combo_type",
              "orders"
            ],
            "type": "object",
            "properties": {
              "client_order_id": {
                "type": "string",
                "description": "Client-defined order identifier. Returned in the response for simple orders. <br/> Represents the unique order ID assigned by the user when placing the order.",
                "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
              },
              "combo_type": {
                "type": "string",
                "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> Futures Trading currently support only the NORMAL type.<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) <br/> Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened.<br/> Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO.<br/> Sub-order quantity limits by combo_type:\n\n| Scenario | combo_type | Supported order_type | Quantity | Description |\n|---|---|---|---|---|\n| Take-Profit/Stop-Loss | MASTER | MARKET, LIMIT | 1 | Master Order |\n| Take-Profit/Stop-Loss | STOP_PROFIT | LIMIT | 0-1 | Take Profit Order |\n| Take-Profit/Stop-Loss | STOP_LOSS | STOP_LOSS | 0-1 | Stop Loss Order |\n| OTO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTO | OTO | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | Triggered Order(s) |\n| OCO | OCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 2-6 | Mutually Cancelling Orders |\n| OTOCO | MASTER | MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1 | Master Order |\n| OTOCO | OTOCO | LIMIT, STOP_LOSS, STOP_LOSS_LIMIT | 1-6 | OCO Order Set Triggered by MASTER |",
                "example": "NORMAL",
                "enum": [
                  "NORMAL",
                  "MASTER",
                  "STOP_PROFIT",
                  "STOP_LOSS",
                  "OTO",
                  "OCO",
                  "OTOCO"
                ]
              },
              "orders": {
                "type": "array",
                "description": "Order Details",
                "items": {
                  "required": [
                    "client_order_id",
                    "order_id",
                    "order_type",
                    "place_time_at",
                    "side",
                    "status",
                    "symbol",
                    "time_in_force",
                    "total_quantity"
                  ],
                  "type": "object",
                  "properties": {
                    "client_order_id": {
                      "type": "string",
                      "description": "Client-defined order identifier. Returned in the response for simple orders. <br/> Represents the unique order ID assigned by the user when placing the order.",
                      "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                    },
                    "order_id": {
                      "type": "string",
                      "description": "System-generated order identifier. Returned in the response for simple orders. <br/> Represents the unique Webull order ID assigned by the system.",
                      "example": "0352U72LQI6DT0KF41GK000000"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                      "example": "AAPL"
                    },
                    "side": {
                      "type": "string",
                      "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                      "example": "BUY",
                      "enum": [
                        "BUY",
                        "SELL",
                        "SHORT"
                      ]
                    },
                    "status": {
                      "type": "string",
                      "description": "&bull; PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion <br/> &bull; SUBMITTED: Indicates that the order has been submitted to the exchange and is awaiting completion<br/> <br/> &bull; CANCELLED: Indicates that the order has been successfully cancelled <br/> &bull; FILLED: Indicates that the order has been fully executed <br/> &bull; FAILED: Indicates a failed order, such as REJECTED <br/> &bull; PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed",
                      "example": "SUBMITTED",
                      "enum": [
                        "PENDING",
                        "SUBMITTED",
                        "CANCELLED",
                        "FILLED",
                        "FAILED",
                        "PARTIAL_FILLED"
                      ]
                    },
                    "order_type": {
                      "type": "string",
                      "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order <br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order <br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order <br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order <br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order <br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order <br/>China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
                      "example": "MARKET",
                      "enum": [
                        "LIMIT",
                        "MARKET",
                        "STOP_LOSS",
                        "STOP_LOSS_LIMIT",
                        "ENHANCED_LIMIT",
                        "AT_AUCTION",
                        "AT_AUCTION_LIMIT",
                        "MARKET_ON_OPEN",
                        "TRAILING_STOP_LOSS",
                        "TRAILING_STOP_LOSS_LIMIT",
                        "TOUCH_MKT",
                        "TOUCH_LMT",
                        "ODD_LOT_LIMIT"
                      ]
                    },
                    "instrument_type": {
                      "type": "string",
                      "description": "Type of financial instrument associated with the request.",
                      "example": "STOCK",
                      "enum": [
                        "EQUITY",
                        "OPTION",
                        "FUTURES"
                      ]
                    },
                    "support_trading_session": {
                      "type": "string",
                      "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only. <br/> Deprecated values: <br/> &bull; Y: [Deprecated]Include extended trading hours. <br/> &bull; N: [Deprecated]Only support regular trading hours. <br/> Active values: <br/> &bull; NIGHT: Only supports night trading. <br/> &bull; ALL: Include extended trading hours. <br/> &bull; CORE: Only support regular trading hours. <br/> &bull; ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day)",
                      "example": "CORE",
                      "enum": [
                        "Y",
                        "N",
                        "NIGHT",
                        "ALL",
                        "CORE",
                        "ALL_DAY"
                      ]
                    },
                    "time_in_force": {
                      "type": "string",
                      "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                      "example": "DAY",
                      "enum": [
                        "DAY",
                        "GTD",
                        "GTC"
                      ]
                    },
                    "total_quantity": {
                      "type": "string",
                      "description": "Total order quantity. Represents the total number of units submitted for this order.",
                      "example": "1"
                    },
                    "filled_quantity": {
                      "type": "string",
                      "description": "Quantity that has been executed. Represents the number of units that have been filled so far.",
                      "example": "1"
                    },
                    "filled_price": {
                      "type": "string",
                      "description": "Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null.",
                      "example": "11.0"
                    },
                    "limit_price": {
                      "type": "string",
                      "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                      "example": "11.0"
                    },
                    "stop_price": {
                      "type": "string",
                      "description": "Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT.<br/> Specifies the trigger price at which the stop order becomes active.",
                      "example": "11.0"
                    },
                    "place_time": {
                      "type": "string",
                      "description": "Order placement time in milliseconds since Unix epoch.",
                      "example": "1726745361658",
                      "deprecated": true
                    },
                    "place_time_at": {
                      "type": "string",
                      "description": "Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                      "example": "2025-11-11T05:44:35.385Z"
                    },
                    "filled_time": {
                      "type": "string",
                      "description": "Time of the last executed trade in milliseconds since Unix epoch.",
                      "example": "1726745361871",
                      "deprecated": true
                    },
                    "filled_time_at": {
                      "type": "string",
                      "description": "Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                      "example": "2025-11-11T05:44:35.385Z"
                    },
                    "legs": {
                      "type": "array",
                      "description": "Leg detail",
                      "items": {
                        "required": [
                          "option_category",
                          "option_expire_date",
                          "option_type",
                          "quantity",
                          "side",
                          "strike_price",
                          "symbol"
                        ],
                        "type": "object",
                        "properties": {
                          "symbol": {
                            "type": "string",
                            "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                            "example": "AAPL"
                          },
                          "side": {
                            "type": "string",
                            "description": "The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).",
                            "example": "BUY",
                            "enum": [
                              "BUY",
                              "SELL",
                              "SHORT"
                            ]
                          },
                          "quantity": {
                            "type": "string",
                            "description": "Quantity of the order. Specifies the number of shares or units to transact.<br/>For US stocks, fractional quantities are allowed and can include decimals.",
                            "example": "1"
                          },
                          "option_type": {
                            "type": "string",
                            "description": "Type of the option. <br/> &bull; CALL: Right to buy the underlying asset. <br/> &bull; PUT: Right to sell the underlying asset.",
                            "example": "CALL",
                            "enum": [
                              "CALL",
                              "PUT"
                            ]
                          },
                          "option_category": {
                            "type": "string",
                            "description": "Category of the option, indicating its exercise style. <br/> Possible values: <br/> &bull; AMERICAN: Can be exercised any time before expiration. <br/> &bull; EUROPEAN: Can only be exercised at expiration. ",
                            "example": "AMERICAN",
                            "enum": [
                              "AMERICAN",
                              "EUROPEAN"
                            ]
                          },
                          "option_strategy": {
                            "type": "string",
                            "description": "Type of options strategy <br/> &bull; SINGLE: Indicates a single-leg options order",
                            "example": "SINGLE",
                            "enum": [
                              "SINGLE"
                            ]
                          },
                          "strike_price": {
                            "type": "string",
                            "description": "Exercise price (strike price) of the option.<br/> Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise.",
                            "example": "190.0"
                          },
                          "option_contract_multiplier": {
                            "type": "string",
                            "description": "The number of shares corresponding to each option contract",
                            "example": "100"
                          },
                          "option_contract_deliverable": {
                            "type": "string",
                            "description": "The number of shares required to exercise each contract",
                            "example": "100"
                          },
                          "option_expire_date": {
                            "type": "string",
                            "description": "Expiration date of the option.<br/> Format: yyyy-MM-dd. After this date, the option will no longer be valid.",
                            "example": "2025-11-21"
                          }
                        },
                        "description": "Leg detail",
                        "title": "OrderListLeg"
                      }
                    },
                    "commission": {
                      "type": "object",
                      "properties": {
                        "actual_commission": {
                          "type": "string",
                          "description": "Actual commission collected",
                          "example": "1.0"
                        },
                        "receivable_commission": {
                          "type": "string",
                          "description": "Receivable commission",
                          "example": "1.0"
                        }
                      },
                      "description": "Commission breakdown",
                      "title": "CommonCommissionResultVO"
                    },
                    "fees": {
                      "type": "array",
                      "description": "Fee breakdown",
                      "items": {
                        "type": "object",
                        "properties": {
                          "type": {
                            "type": "string",
                            "description": "Fee type",
                            "example": "FINRA_CAT_REGULATORY_FEE"
                          },
                          "actual_value": {
                            "type": "string",
                            "description": "Actual fee collected",
                            "example": "1.0"
                          },
                          "receivable_value": {
                            "type": "string",
                            "description": "Receivable fee",
                            "example": "1.0"
                          }
                        },
                        "description": "Fee breakdown",
                        "title": "CommonFeeResultVO"
                      }
                    }
                  },
                  "description": "Order Details",
                  "title": "OrderDetailItem"
                }
              }
            },
            "title": "OrderDetailResult"
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
    "name": "Get Order Detail",
    "description": {
      "content": "Retrieves the specified order details through the order ID.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "orders",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Account identifier.",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique client-defined identifier for the order. <br/>Maximum length is 32 characters and must be unique per account. <br/>Used to track or reference the order when interacting with the system.",
            "type": "text/plain"
          },
          "key": "client_order_id",
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

## Cash Activities

> Source: <https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md>

### List Cash Activities

Lists an account's cash activities, filterable by type and time range. Defaults to the last 7 days if no date is provided.

### OpenAPI definition

```json
{
  "info": {
    "title": "Webull OpenAPI Documentation",
    "description": "The Webull OpenAPI enables integration of trading APIs, market data, and OAuth authentication for building trading applications and brokerage solutions. It supports HTTP-based historical and real-time market data and MQTT streaming via WebSocket/TCP, along with SDKs, secure authentication, and APIs for orders, accounts, and event contract trading.",
    "contact": {
      "name": "Webull Developer Support",
      "url": "https://www.webull.com/help",
      "email": "api-support@webull-us.com"
    },
    "version": "2.0",
    "x-logo": {
      "url": "static/png/logo.png"
    }
  },
  "servers": [
    {
      "url": "https://api.sandbox.webull.com"
    }
  ],
  "path": "/trading/activities/cash-activities/list",
  "method": "get",
  "tags": [
    "Activities"
  ],
  "description": "Lists an account's cash activities, filterable by type and time range. Defaults to the last 7 days if no date is provided.",
  "operationId": "tradeCashActivityByType",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Provide the target account id",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "943a9802f6c14983b3b4755c69c01717"
    },
    {
      "name": "activity_types",
      "in": "query",
      "description": "Account activity types.<br/>Note: EC_STATEMENT is deprecated; use EC_SETTLEMENT instead (EC_STATEMENT behaves the same as EC_SETTLEMENT).<br/>Note: Crypto accounts only support: TRADE, DEPOSIT, WITHDRAW, and FEES.\n",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "TRADE",
          "DEPOSIT",
          "WITHDRAW",
          "FEES",
          "TRANSFER",
          "DIVIDENDS",
          "TAX",
          "INTERESTS",
          "CORPORATE_ACTION",
          "OPTION_EA",
          "EC_SETTLEMENT",
          "JOURNAL"
        ]
      },
      "example": "DEPOSIT,TRADE"
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "Activity query start time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. Cross-year queries are not supported; start_time and end_time must be within the same year.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:59.012Z"
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "If not provided, the default query is the last 7 days.<br/>Activity query end time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. Cross-year queries are not supported; start_time and end_time must be within the same year.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-06T22:59:59.012Z"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
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
                  "required": [
                    "account_id",
                    "activity_sub_type",
                    "activity_type",
                    "biz_time",
                    "currency",
                    "id",
                    "net_amount",
                    "trade_date"
                  ],
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "string",
                      "description": "Unique ID",
                      "example": "a1b2c3d4e5f6g7h8i9j0"
                    },
                    "account_id": {
                      "type": "string",
                      "description": "Account ID",
                      "example": "943a9802f6c14983b3b4755c69c01717"
                    },
                    "activity_type": {
                      "type": "string",
                      "description": "Activity Type\n| Code             | Description                          |\n|------------------|--------------------------------------|\n| TRADE            | Trade Activity Type                  |\n| DEPOSIT          | Deposit Activity Type                |\n| WITHDRAW         | Withdraw Activity Type               |\n| FEES             | Fees Activity Type                   |\n| TRANSFER         | Transfer Activity Type               |\n| DIVIDENDS        | Dividends Activity Type              |\n| TAX              | Tax Activity Type                    |\n| INTERESTS        | Interests Activity Type              |\n| CORPORATE_ACTION | Corporate Action Activity Type       |\n| OPTION_EA        | Option Exercise and Assignment Type  |\n| JOURNAL          | Journal Activity Type                |\n| EC_SETTLEMENT    | EC Settlement Activity Type          |\n| OTHER            | Other Activity Type                  |",
                      "example": "TRADE",
                      "enum": [
                        "TRADE",
                        "DEPOSIT",
                        "WITHDRAW",
                        "FEES",
                        "TRANSFER",
                        "DIVIDENDS",
                        "TAX",
                        "INTERESTS",
                        "CORPORATE_ACTION",
                        "OPTION_EA",
                        "JOURNAL",
                        "EC_SETTLEMENT",
                        "OTHER"
                      ]
                    },
                    "activity_sub_type": {
                      "type": "string",
                      "description": "Activity Type and Sub Type Mapping.\n| ActivityType     | ActivitySubType            |\n|------------------|---------------------------|\n| TRADE            | BUY                       |\n| TRADE            | SELL                      |\n| TRADE            | BUY_CANCELLED             |\n| TRADE            | SELL_CANCELLED            |\n| TRADE            | PENNY_FOR_LOT             |\n| TRADE            | FX_EXCHANGE               |\n| DEPOSIT          | WIRE                      |\n| DEPOSIT          | ACH                       |\n| DEPOSIT          | REVERSAL                  |\n| DEPOSIT          | CHECK                     |\n| DEPOSIT          | ACH_REVERSE                     |\n| DEPOSIT          | INTERNAL_TRANSFER                     |\n| WITHDRAW         | WIRE                      |\n| WITHDRAW         | ACH                       |\n| WITHDRAW         | REVERSAL                  |\n| WITHDRAW         | CHECK                     |\n| WITHDRAW         | INTERNAL_TRANSFER                     |\n| FEES             | WIRE_FEE                  |\n| FEES             | REVERSAL_FEE              |\n| FEES             | TRANSFER_ACATS            |\n| FEES             | CA_HANDLING_FEE           |\n| FEES             | ADR                       |\n| FEES             | PAPER_STATEMENT_FEE       |\n| FEES             | PAPER_CONFIRM_FEE         |\n| FEES             | WRITE_OFF                 |\n| FEES             | CHECK_FEE                 |\n| FEES             | ADVISORY_FEE              |\n| FEES             | SUBSCRIPTION_FEE          |\n| FEES             | SERVICE_FEE               |\n| FEES             | ACH_REVERSE_FEE               |\n| FEES             | OTHER                     |\n| TRANSFER         | ACATS_IN                  |\n| TRANSFER         | ACATS_OUT                 |\n| TRANSFER         | INTERNAL_TRANSFER         |\n| TRANSFER         | BANK_SWEEP                |\n| DIVIDENDS        | INCOME                    |\n| DIVIDENDS        | PAYMENT_IN_LIEU           |\n| DIVIDENDS        | CASH_IN_LIEU              |\n| DIVIDENDS        | LP_DISTRIBUTION           |\n| TAX              | FOREIGN_TAX_WITHHELD      |\n| TAX              | US_TAX_WITHHOLDING        |\n| TAX              | IRA_FED_WITHHOLDING       |\n| TAX              | STATE_WITHHOLDING         |\n| TAX              | WITHHOLDING_TAX           |\n| TAX              | GP_TAX_WITHHELD           |\n| TAX              | OTHER                     |\n| INTERESTS        | CREDIT                    |\n| INTERESTS        | DEBIT                     |\n| INTERESTS        | STOCK_BORROW_INTEREST     |\n| INTERESTS        | SECURITIES_LENDING_INCOME |\n| INTERESTS        | ADJUSTMENT                |\n| INTERESTS        | PAYMENT                   |\n| INTERESTS        | INTEREST_REBATE           |\n| CORPORATE_ACTION | CASH_IN_LIEU              |\n| CORPORATE_ACTION | REDEMPTION                |\n| CORPORATE_ACTION | MISC_ADJUSTMENT           |\n| CORPORATE_ACTION | MERGER                    |\n| CORPORATE_ACTION | RIGHTS_OFFERING           |\n| CORPORATE_ACTION | IDENTIFIER_CHANGE         |\n| CORPORATE_ACTION | REVERSE_SPLIT             |\n| CORPORATE_ACTION | FORWARD_SPLIT             |\n| CORPORATE_ACTION | SPIN_OFF                  |\n| CORPORATE_ACTION | CONVERSION                |\n| CORPORATE_ACTION | LIQUIDATION               |\n| OPTION_EA        | CONTRACT_CLOSE            |\n| OPTION_EA        | OPTION_ASSIGNMENT         |\n| OPTION_EA        | OPTION_EXPIRATION         |\n| OPTION_EA        | OPTION_EXERCISE           |\n| JOURNAL          | CASH_JOURNAL              |\n| EC_SETTLEMENT    | EC_EXPIRATION             |\n| EC_SETTLEMENT    | EC_PAYOUT                 |\n| OTHER            | INCOME                    |\n| OTHER            | LENDING_REBATE            |\n| OTHER            | DVP                       |\n| OTHER            | GRID_TRANSFER                       |\n| OTHER            | OTHER                     |",
                      "example": "BUY",
                      "enum": [
                        "BUY",
                        "SELL",
                        "BUY_CANCELLED",
                        "SELL_CANCELLED",
                        "PENNY_FOR_LOT",
                        "TRADE",
                        "FX_EXCHANGE",
                        "OPTION_EXPIRATION",
                        "WIRE",
                        "ACH",
                        "REVERSAL",
                        "CHECK",
                        "ACH_REVERSE",
                        "WIRE_FEE",
                        "REVERSAL_FEE",
                        "TRANSFER_ACATS",
                        "CA_HANDLING_FEE",
                        "ADR",
                        "PAPER_STATEMENT_FEE",
                        "PAPER_CONFIRM_FEE",
                        "WRITE_OFF",
                        "CHECK_FEE",
                        "ADVISORY_FEE",
                        "SUBSCRIPTION_FEE",
                        "SERVICE_FEE",
                        "ACH_REVERSE_FEE",
                        "ACATS_IN",
                        "ACATS_OUT",
                        "INTERNAL_TRANSFER",
                        "BANK_SWEEP",
                        "INCOME",
                        "PAYMENT_IN_LIEU",
                        "CASH_IN_LIEU",
                        "LP_DISTRIBUTION",
                        "FOREIGN_TAX_WITHHELD",
                        "US_TAX_WITHHOLDING",
                        "IRA_FED_WITHHOLDING",
                        "STATE_WITHHOLDING",
                        "WITHHOLDING_TAX",
                        "GP_TAX_WITHHELD",
                        "CREDIT",
                        "DEBIT",
                        "STOCK_BORROW_INTEREST",
                        "SECURITIES_LENDING_INCOME",
                        "ADJUSTMENT",
                        "PAYMENT",
                        "INTEREST_REBATE",
                        "REDEMPTION",
                        "MISC_ADJUSTMENT",
                        "MERGER",
                        "RIGHTS_OFFERING",
                        "IDENTIFIER_CHANGE",
                        "REVERSE_SPLIT",
                        "FORWARD_SPLIT",
                        "SPIN_OFF",
                        "CONVERSION",
                        "LIQUIDATION",
                        "CONTRACT_CLOSE",
                        "OPTION_ASSIGNMENT",
                        "OPTION_EXERCISE",
                        "CASH_JOURNAL",
                        "EC_EXPIRATION",
                        "EC_PAYOUT",
                        "LENDING_REBATE",
                        "DVP",
                        "GRID_TRANSFER",
                        "OTHER"
                      ]
                    },
                    "currency": {
                      "type": "string",
                      "description": "Currency",
                      "example": "USD",
                      "enum": [
                        "USD"
                      ]
                    },
                    "market": {
                      "type": "string",
                      "description": "Market Code<br/>US - US Market<br/>",
                      "example": "US",
                      "enum": [
                        "US"
                      ]
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Activity Symbol",
                      "example": "AAPL"
                    },
                    "trade_date": {
                      "type": "string",
                      "description": "Accounting date of the transaction (trade date), format: yyyy-MM-dd",
                      "example": "2024-05-01"
                    },
                    "net_amount": {
                      "type": "string",
                      "description": "Net change amount of the transaction (positive for credit, negative for debit)",
                      "example": "1500.0"
                    },
                    "biz_time": {
                      "type": "string",
                      "description": "Business event time when the transaction occurred",
                      "example": "2024-05-01T10:15:30.691Z"
                    }
                  },
                  "description": "Activity Cash Result",
                  "title": "ActivityCashResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoActivityCashResult"
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
    "name": "List Cash Activities",
    "description": {
      "content": "Lists an account's cash activities, filterable by type and time range. Defaults to the last 7 days if no date is provided.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "activities",
        "cash-activities",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Provide the target account id",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Account activity types.<br/>Note: EC_STATEMENT is deprecated; use EC_SETTLEMENT instead (EC_STATEMENT behaves the same as EC_SETTLEMENT).<br/>Note: Crypto accounts only support: TRADE, DEPOSIT, WITHDRAW, and FEES.\n",
            "type": "text/plain"
          },
          "key": "activity_types",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Activity query start time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. Cross-year queries are not supported; start_time and end_time must be within the same year.",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "If not provided, the default query is the last 7 days.<br/>Activity query end time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. Cross-year queries are not supported; start_time and end_time must be within the same year.",
            "type": "text/plain"
          },
          "key": "end_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Pagination key from previous response for next page.",
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


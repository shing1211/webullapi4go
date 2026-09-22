# Event Contracts — Verbatim Reference

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> Event-contract instruments and market data (US documentation).

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Event Contract Categories

> Source: <https://developer.webull.com/apis/docs/reference/event-categories-list.md>

### List Event Contract Categories

Retrieves all categories under the Event Contract.

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
  "path": "/trading/instruments/event-contracts/categories/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves all categories under the Event Contract.",
  "operationId": "eventCategoriesList",
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
              "required": [
                "category_code",
                "category_id",
                "category_name"
              ],
              "type": "object",
              "properties": {
                "category_id": {
                  "type": "integer",
                  "description": "Category ID.",
                  "format": "int32",
                  "example": 1
                },
                "category_code": {
                  "type": "string",
                  "description": "Unique identifier code for category.",
                  "example": "ECONOMICS"
                },
                "category_name": {
                  "type": "string",
                  "description": "Category corresponding name.",
                  "example": "Economics"
                }
              },
              "description": "Event Categories",
              "title": "EventCategoriesVo"
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
    "name": "List Event Contract Categories",
    "description": {
      "content": "Retrieves all categories under the Event Contract.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "event-contracts",
        "categories",
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

## Event Contract Series

> Source: <https://developer.webull.com/apis/docs/reference/event-series-list.md>

### List Event Contract Series

Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., "Monthly Jobs Report").

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
  "path": "/trading/instruments/event-contracts/series/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., \"Monthly Jobs Report\").",
  "operationId": "eventSeriesList",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "The category which this series belongs to.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "ECONOMICS",
          "FINANCIALS",
          "POLITICS",
          "ENTERTAINMENT",
          "SCIENCE_TECHNOLOGY",
          "CLIMATE_WEATHER",
          "TRANSPORTATION",
          "CRYPTO",
          "SPORTS"
        ]
      },
      "example": "ECONOMICS"
    },
    {
      "name": "symbols",
      "in": "query",
      "description": "List of series symbols, maximum 100 symbols per query.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXRATECUTCOUNT,KXFEDDECISION"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiIxNTMxMTgiLCJwYWdlSW5kZXgiOjAsInBhZ2VTaXplIjo1MDAsImNvbmRpdGlvbiI6Im51bGw7bnVsbDsifQ=="
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
                    "category",
                    "frequency",
                    "name",
                    "series_id",
                    "symbol"
                  ],
                  "type": "object",
                  "properties": {
                    "category": {
                      "type": "string",
                      "description": "The category which this series belongs to.",
                      "example": "ECONOMICS"
                    },
                    "series_id": {
                      "type": "string",
                      "description": "ID that identifies this series.",
                      "example": "151750"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Symbol that identifies this series.",
                      "example": "KXRATECUTCOUNT"
                    },
                    "name": {
                      "type": "string",
                      "description": "Name that describes the series.",
                      "example": "Number of Rate Cuts"
                    },
                    "frequency": {
                      "type": "string",
                      "description": "Description of the frequency of the series.",
                      "example": "ANNUAL",
                      "enum": [
                        "HOURLY",
                        "DAILY",
                        "WEEKLY",
                        "MONTHLY",
                        "ANNUAL",
                        "ONE_OFF",
                        "CUSTOM"
                      ]
                    }
                  },
                  "description": "Event series",
                  "title": "EventSeriesVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoEventSeriesVo"
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
    "name": "List Event Contract Series",
    "description": {
      "content": "Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., \"Monthly Jobs Report\").",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "event-contracts",
        "series",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "The category which this series belongs to.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "List of series symbols, maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
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

## Event Contract Events

> Source: <https://developer.webull.com/apis/docs/reference/event-events-list.md>

### List Event Contract Events

Retrieves events under the Event Contract matching the query.

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
  "path": "/trading/instruments/event-contracts/events/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves events under the Event Contract matching the query.",
  "operationId": "eventEventsList",
  "parameters": [
    {
      "name": "series_symbol",
      "in": "query",
      "description": "Symbol that identifies this series.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "ECONOMICS"
    },
    {
      "name": "symbols",
      "in": "query",
      "description": "List of events symbols, maximum 100 symbols per query.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXRATECUTCOUNT,KXFEDDECISION"
    },
    {
      "name": "status",
      "in": "query",
      "description": "The status of the event.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "ACTIVE",
          "INACTIVE"
        ]
      },
      "example": "ACTIVE"
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
                "mutually_exclusive",
                "name",
                "series_id",
                "short_name",
                "status",
                "symbol"
              ],
              "type": "object",
              "properties": {
                "series_id": {
                  "type": "string",
                  "description": "ID identification of series.",
                  "example": "151950"
                },
                "symbol": {
                  "type": "string",
                  "description": "The symbol of the event",
                  "example": "KXGDP-27JAN30"
                },
                "name": {
                  "type": "string",
                  "description": "The name of the event.",
                  "example": "US GDP growth in Q4 2026?"
                },
                "status": {
                  "type": "string",
                  "description": "The status of the event.",
                  "example": "inactive",
                  "enum": [
                    "ACTIVE",
                    "INACTIVE"
                  ]
                },
                "short_name": {
                  "type": "string",
                  "description": "The abbreviation for event.",
                  "example": "In Q4 2026"
                },
                "strike_date": {
                  "type": "string",
                  "description": "Exercise Date.",
                  "example": "2023-11-02"
                },
                "strike_period": {
                  "type": "string",
                  "description": "Exercise period.",
                  "example": "Q4 2026"
                },
                "mutually_exclusive": {
                  "type": "boolean",
                  "description": "Whether mutually exclusive.",
                  "example": false
                }
              },
              "title": "EventEventsVo"
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
    "name": "List Event Contract Events",
    "description": {
      "content": "Retrieves events under the Event Contract matching the query.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "event-contracts",
        "events",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Symbol that identifies this series.",
            "type": "text/plain"
          },
          "key": "series_symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "List of events symbols, maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The status of the event.",
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

## Event Contract Instruments

> Source: <https://developer.webull.com/apis/docs/reference/event-market-list.md>

### List Event Contract Markets

Retrieves event contract market instruments for the given series symbol.

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
  "path": "/trading/instruments/event-contracts/markets/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves event contract market instruments for the given series symbol.",
  "operationId": "eventMarketList",
  "parameters": [
    {
      "name": "series_symbol",
      "in": "query",
      "description": "Symbol that identifies this series.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "ECONOMICS"
    },
    {
      "name": "event_symbol",
      "in": "query",
      "description": "Symbol of the event events.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXRATECUTCOUNT-25DEC31"
    },
    {
      "name": "symbols",
      "in": "query",
      "description": "List of security symbols, maximum 100 symbols per query.",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXRATECUTCOUNT-25DEC31-T8"
    },
    {
      "name": "expiration_date_after",
      "in": "query",
      "description": "Used to filter items whose expiration date is later than a specified date; the default selection is the current day (inclusive).",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "2026-12-31"
    },
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "eyJ2IjoxLCJsYXN0SWQiOiIxNTMxMTgiLCJwYWdlSW5kZXgiOjAsInBhZ2VTaXplIjo1MDAsImNvbmRpdGlvbiI6Im51bGw7bnVsbDsifQ=="
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
                    "can_close_early",
                    "event_name",
                    "event_symbol",
                    "expected_exp_date",
                    "fractionable",
                    "instrument_id",
                    "last_trading_date",
                    "latest_exp_date",
                    "name",
                    "payout_date",
                    "series_id",
                    "series_name",
                    "series_symbol",
                    "status",
                    "symbol",
                    "tradable_status",
                    "yes_condition"
                  ],
                  "type": "object",
                  "properties": {
                    "series_id": {
                      "type": "string",
                      "description": "ID that identifies this series.",
                      "example": "151750"
                    },
                    "series_symbol": {
                      "type": "string",
                      "description": "Symbol that identifies this series.",
                      "example": "KXRATECUTCOUNT"
                    },
                    "series_name": {
                      "type": "string",
                      "description": "Name that describes the series.",
                      "example": "Number of Rate Cuts"
                    },
                    "event_symbol": {
                      "type": "string",
                      "description": "Symbol that identifies this events.",
                      "example": "KXRATECUTCOUNT-26DEC31"
                    },
                    "event_name": {
                      "type": "string",
                      "description": "Name that describes the events.",
                      "example": "Number of rate cuts in 2026?"
                    },
                    "instrument_id": {
                      "type": "string",
                      "description": "Unique id of the event market.",
                      "example": "502257268"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Symbol of the event market.",
                      "example": "KXRATECUTCOUNT-25DEC31-T3"
                    },
                    "name": {
                      "type": "string",
                      "description": "Name of the event market.",
                      "example": "Will the Fed cut rates 3 times?"
                    },
                    "yes_condition": {
                      "type": "string",
                      "description": "Conditions for a 'Yes' outcome.",
                      "example": "Exactly 3 cuts"
                    },
                    "last_trading_date": {
                      "type": "string",
                      "description": "Last Notice Day.",
                      "example": "2025-12-31"
                    },
                    "status": {
                      "type": "string",
                      "description": "Listing status.",
                      "example": "LISTING",
                      "enum": [
                        "NOT_SET",
                        "LISTING",
                        "DELISTING",
                        "OTHER",
                        "UNRECOGNIZED"
                      ]
                    },
                    "tradable_status": {
                      "type": "string",
                      "description": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
                      "example": "NT",
                      "enum": [
                        "OC",
                        "CO",
                        "NT"
                      ]
                    },
                    "can_close_early": {
                      "type": "boolean",
                      "description": "Can the contract close early?",
                      "example": true
                    },
                    "expected_exp_date": {
                      "type": "string",
                      "description": "Expected expiration date of contract.",
                      "example": "2025-12-31"
                    },
                    "latest_exp_date": {
                      "type": "string",
                      "description": "Latest expiration date.",
                      "example": "2026-01-01"
                    },
                    "payout_date": {
                      "type": "string",
                      "description": "Settlement/Payment Date.",
                      "example": "2025-12-31"
                    },
                    "fractionable": {
                      "type": "boolean",
                      "description": "Support fragmented event contracts.",
                      "example": true
                    },
                    "price_ranges": {
                      "type": "array",
                      "description": "Price range.",
                      "items": {
                        "required": [
                          "end",
                          "start",
                          "step"
                        ],
                        "type": "object",
                        "properties": {
                          "start": {
                            "type": "string",
                            "description": "Start price.",
                            "example": "0.0"
                          },
                          "end": {
                            "type": "string",
                            "description": "End price.",
                            "example": "1.0"
                          },
                          "step": {
                            "type": "string",
                            "description": "Step length.",
                            "example": "0.01"
                          }
                        },
                        "description": "Price range.",
                        "title": "PriceRange"
                      }
                    }
                  },
                  "description": "Event market",
                  "title": "EventMarketVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoEventMarketVo"
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
    "name": "List Event Contract Markets",
    "description": {
      "content": "Retrieves event contract market instruments for the given series symbol.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "trading",
        "instruments",
        "event-contracts",
        "markets",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "Symbol that identifies this series.",
            "type": "text/plain"
          },
          "key": "series_symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Symbol of the event events.",
            "type": "text/plain"
          },
          "key": "event_symbol",
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
            "content": "Used to filter items whose expiration date is later than a specified date; the default selection is the current day (inclusive).",
            "type": "text/plain"
          },
          "key": "expiration_date_after",
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

## Event Snapshot

> Source: <https://developer.webull.com/apis/docs/reference/event-snapshot.md>

### List Event Snapshots

Retrieves a real-time snapshot for an event instrument.

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
  "path": "/market-data/event-contracts/snapshots/list",
  "method": "get",
  "tags": [
    "Event Market Data"
  ],
  "description": "Retrieves a real-time snapshot for an event instrument.",
  "operationId": "eventSnapshot",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "Symbol of the event market, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXU3-25OCT-T3.8"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ]
      },
      "example": "US_EVENT"
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
                "instrument_id",
                "last_trade_time",
                "name",
                "no_ask",
                "no_ask_size",
                "no_bid",
                "no_bid_size",
                "open_interest",
                "price",
                "symbol",
                "volume",
                "yes_ask",
                "yes_ask_size",
                "yes_bid",
                "yes_bid_size"
              ],
              "type": "object",
              "properties": {
                "instrument_id": {
                  "type": "string",
                  "description": "Unique id of the event market.",
                  "example": "504279491"
                },
                "symbol": {
                  "type": "string",
                  "description": "Symbol of the event market.",
                  "example": "KXCPI-26JAN-T0.3"
                },
                "name": {
                  "type": "string",
                  "description": "Name of the event market.",
                  "example": "Will CPI rise more than 0.3% in January 2026?"
                },
                "price": {
                  "type": "string",
                  "description": "The current market price for buying or selling an event contract",
                  "example": "0.13"
                },
                "volume": {
                  "type": "string",
                  "description": "Number of contracts bought on this event market.",
                  "example": "32"
                },
                "last_trade_time": {
                  "type": "integer",
                  "description": "Price for the last traded YES contract on this market.",
                  "format": "int64",
                  "example": 1768861137000
                },
                "open_interest": {
                  "type": "string",
                  "description": "Number of contracts bought on this event market disconsidering netting.",
                  "example": "14240"
                },
                "yes_bid": {
                  "type": "string",
                  "description": "Price for the highest YES buy offer on this event market.",
                  "example": "0.08"
                },
                "yes_bid_size": {
                  "type": "string",
                  "description": "Size for the highest YES buy offer on this event market.",
                  "example": "2115"
                },
                "yes_ask": {
                  "type": "string",
                  "description": "Price for the lowest YES sell offer on this event market.",
                  "example": "0.13"
                },
                "yes_ask_size": {
                  "type": "string",
                  "description": "Size for the lowest YES sell offer on this event market.",
                  "example": "543"
                },
                "no_bid": {
                  "type": "string",
                  "description": "Price for the highest NO buy offer on this event market.",
                  "example": "0.87"
                },
                "no_bid_size": {
                  "type": "string",
                  "description": "Size for the highest NO buy offer on this event market.",
                  "example": "543"
                },
                "no_ask": {
                  "type": "string",
                  "description": "Price for the lowest NO sell offer on this event market.",
                  "example": "0.92"
                },
                "no_ask_size": {
                  "type": "string",
                  "description": "Size for the lowest NO sell offer on this event market.",
                  "example": "2115"
                }
              },
              "description": "Event snapshot",
              "title": "EventSnapshotVo"
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
    "name": "List Event Snapshots",
    "description": {
      "content": "Retrieves a real-time snapshot for an event instrument.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
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
            "content": "(Required) Symbol of the event market, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
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

## Event Depth

> Source: <https://developer.webull.com/apis/docs/reference/event-depth.md>

### List Event Depths

Retrieves the order book for an event instrument. Only yes/no bids are returned (in binary markets a yes bid at X equals a no ask at 100-X).

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
  "path": "/market-data/event-contracts/depths/list",
  "method": "get",
  "tags": [
    "Event Market Data"
  ],
  "description": "Retrieves the order book for an event instrument. Only yes/no bids are returned (in binary markets a yes bid at X equals a no ask at 100-X).",
  "operationId": "eventDepth",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Symbol of the event market.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXU3-25OCT-T3.8"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ]
      },
      "example": "US_EVENT"
    },
    {
      "name": "depth",
      "in": "query",
      "description": "Depth of buying and selling orders, default 10 levels, etc.",
      "required": false,
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
            "type": "array",
            "items": {
              "required": [
                "instrument_id",
                "no_asks",
                "no_bids",
                "quote_time",
                "symbol",
                "yes_asks",
                "yes_bids"
              ],
              "type": "object",
              "properties": {
                "instrument_id": {
                  "type": "string",
                  "description": "Unique id of the event market.",
                  "example": "504279491"
                },
                "symbol": {
                  "type": "string",
                  "description": "Symbol of the event market.",
                  "example": "KXCPI-26JAN-T0.3"
                },
                "quote_time": {
                  "type": "integer",
                  "description": "Quotation Time.",
                  "format": "int64",
                  "example": 1768872168870
                },
                "yes_bids": {
                  "type": "array",
                  "description": "Yes, buy order array.",
                  "items": {
                    "required": [
                      "price",
                      "size"
                    ],
                    "type": "object",
                    "properties": {
                      "price": {
                        "type": "string",
                        "description": "Price",
                        "example": "0.13"
                      },
                      "size": {
                        "type": "string",
                        "description": "Trading volume.",
                        "example": "543"
                      }
                    },
                    "description": "No, sell order array.",
                    "title": "AskBid"
                  }
                },
                "yes_asks": {
                  "type": "array",
                  "description": "Yes, sell order array.",
                  "items": {
                    "required": [
                      "price",
                      "size"
                    ],
                    "type": "object",
                    "properties": {
                      "price": {
                        "type": "string",
                        "description": "Price",
                        "example": "0.13"
                      },
                      "size": {
                        "type": "string",
                        "description": "Trading volume.",
                        "example": "543"
                      }
                    },
                    "description": "No, sell order array.",
                    "title": "AskBid"
                  }
                },
                "no_bids": {
                  "type": "array",
                  "description": "No, buy order array.",
                  "items": {
                    "required": [
                      "price",
                      "size"
                    ],
                    "type": "object",
                    "properties": {
                      "price": {
                        "type": "string",
                        "description": "Price",
                        "example": "0.13"
                      },
                      "size": {
                        "type": "string",
                        "description": "Trading volume.",
                        "example": "543"
                      }
                    },
                    "description": "No, sell order array.",
                    "title": "AskBid"
                  }
                },
                "no_asks": {
                  "type": "array",
                  "description": "No, sell order array.",
                  "items": {
                    "required": [
                      "price",
                      "size"
                    ],
                    "type": "object",
                    "properties": {
                      "price": {
                        "type": "string",
                        "description": "Price",
                        "example": "0.13"
                      },
                      "size": {
                        "type": "string",
                        "description": "Trading volume.",
                        "example": "543"
                      }
                    },
                    "description": "No, sell order array.",
                    "title": "AskBid"
                  }
                }
              },
              "description": "Event depth",
              "title": "EventDepthVo"
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
    "name": "List Event Depths",
    "description": {
      "content": "Retrieves the order book for an event instrument. Only yes/no bids are returned (in binary markets a yes bid at X equals a no ask at 100-X).",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
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
            "content": "(Required) Symbol of the event market.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Depth of buying and selling orders, default 10 levels, etc.",
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

## Event Bars

> Source: <https://developer.webull.com/apis/docs/reference/event-bars.md>

### List Event Bars

Retrieves the most recent N bars for an event symbol.

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
  "path": "/market-data/event-contracts/bars/list",
  "method": "get",
  "tags": [
    "Event Market Data"
  ],
  "description": "Retrieves the most recent N bars for an event symbol.",
  "operationId": "eventBars",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "Symbol of the event market, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXU3-25OCT-T3.8"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ]
      },
      "example": "US_EVENT"
    },
    {
      "name": "timespan",
      "in": "query",
      "description": "Bar time granularity.",
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
          "D"
        ]
      },
      "example": "M1"
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
      "example": 200
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
                "instrument_id",
                "result",
                "symbol"
              ],
              "type": "object",
              "properties": {
                "instrument_id": {
                  "type": "string",
                  "description": "Unique id of the event market.",
                  "example": "504279491"
                },
                "symbol": {
                  "type": "string",
                  "description": "Symbol of the event market.",
                  "example": "KXCPI-26JAN-T0.3"
                },
                "result": {
                  "type": "array",
                  "description": "K-line data.",
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
                      "open": {
                        "type": "string",
                        "description": "Opening price.",
                        "example": "0.05"
                      },
                      "close": {
                        "type": "string",
                        "description": "Close price.",
                        "example": "0.05"
                      },
                      "high": {
                        "type": "string",
                        "description": "High price.",
                        "example": "0.05"
                      },
                      "low": {
                        "type": "string",
                        "description": "Low price.",
                        "example": "0.05"
                      },
                      "volume": {
                        "type": "string",
                        "description": "Turnover.",
                        "example": "1"
                      },
                      "time": {
                        "type": "string",
                        "description": "UTC time",
                        "example": "2021-12-28T09:00:09.945+0000"
                      }
                    },
                    "description": "K-line data.",
                    "title": "Kdata"
                  }
                }
              },
              "title": "EventBarsVo"
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
    "name": "List Event Bars",
    "description": {
      "content": "Retrieves the most recent N bars for an event symbol.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
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
            "content": "(Required) Symbol of the event market, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Bar time granularity.",
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

## Event Tick

> Source: <https://developer.webull.com/apis/docs/reference/event-tick.md>

### List Event Ticks

Retrieves tick-by-tick trades for an event contract, sorted latest first.

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
  "path": "/market-data/event-contracts/ticks/list",
  "method": "get",
  "tags": [
    "Event Market Data"
  ],
  "description": "Retrieves tick-by-tick trades for an event contract, sorted latest first.",
  "operationId": "eventTick",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Symbol of the event market.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXU3-25OCT-T3.8"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ]
      },
      "example": "US_EVENT"
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of tick, default 30, maximum limit 1200.",
      "required": false,
      "schema": {
        "type": "string",
        "description": "30-1200",
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
            "type": "array",
            "items": {
              "required": [
                "instrument_id",
                "result",
                "symbol"
              ],
              "type": "object",
              "properties": {
                "instrument_id": {
                  "type": "string",
                  "description": "Unique id of the event market.",
                  "example": "504279491"
                },
                "symbol": {
                  "type": "string",
                  "description": "Symbol of the event market.",
                  "example": "KXCPI-26JAN-T0.3"
                },
                "result": {
                  "type": "array",
                  "description": "Tick data.",
                  "items": {
                    "required": [
                      "no_price",
                      "side",
                      "time",
                      "trade_id",
                      "volume",
                      "yes_price"
                    ],
                    "type": "object",
                    "properties": {
                      "time": {
                        "type": "string",
                        "description": "Timestamp when this trade was executed.",
                        "example": "1772730554000"
                      },
                      "yes_price": {
                        "type": "string",
                        "description": "Yes price for this trade in dollars.",
                        "example": "0.05"
                      },
                      "no_price": {
                        "type": "string",
                        "description": "No price for this trade in dollars.",
                        "example": "0.95"
                      },
                      "volume": {
                        "type": "string",
                        "description": "Turnover.",
                        "example": "1.0"
                      },
                      "side": {
                        "type": "string",
                        "description": "Side for the taker of this trade(yes/no).",
                        "example": "no"
                      },
                      "trade_id": {
                        "type": "string",
                        "description": "Unique identifier for this trade.",
                        "example": "604d16d9-9000-4266-68f9-1ec898e2acf6"
                      }
                    },
                    "description": "Tick data.",
                    "title": "EventTick"
                  }
                }
              },
              "title": "EventTickVo"
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
    "name": "List Event Ticks",
    "description": {
      "content": "Retrieves tick-by-tick trades for an event contract, sorted latest first.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
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
            "content": "(Required) Symbol of the event market.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of tick, default 30, maximum limit 1200.",
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

## Event Contract Tags

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/all-tags-using-get.md>

### List All Tags

Retrieves all available category tags for event contract series. Use the returned tags to filter series via the tags parameter in the Series List endpoint.

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
  "path": "/market-data/instruments/event-contracts/categories/tags/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves all available category tags for event contract series. Use the returned tags to filter series via the tags parameter in the Series List endpoint.",
  "operationId": "allTagsUsingGET",
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
                "tags": {
                  "type": "array",
                  "description": "Category Tag list",
                  "example": "['Basketball', 'Baseball','Football']",
                  "items": {
                    "type": "string",
                    "description": "Category Tag list",
                    "example": "['Basketball', 'Baseball','Football']"
                  }
                },
                "category_id": {
                  "type": "integer",
                  "description": "Category ID",
                  "format": "int32",
                  "example": 1
                },
                "category_name": {
                  "type": "string",
                  "description": "Category Name",
                  "example": "Sports"
                },
                "category_code": {
                  "type": "string",
                  "description": "Category Code",
                  "example": "SPORTS"
                }
              },
              "description": "Serise Category Information",
              "title": "SeriseCategory"
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List All Tags",
    "description": {
      "content": "Retrieves all available category tags for event contract series. Use the returned tags to filter series via the tags parameter in the Series List endpoint.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "event-contracts",
        "categories",
        "tags",
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
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Event Contract Events List

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-list-using-get.md>

### List Events

Retrieves a list of tradable events under a series. Each event represents a specific question or market (e.g., '2026-27 College Football National Championship Winner'). Filter by series_symbol to get events for a specific series.

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
  "path": "/market-data/instruments/event-contracts/events/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves a list of tradable events under a series. Each event represents a specific question or market (e.g., '2026-27 College Football National Championship Winner'). Filter by series_symbol to get events for a specific series.",
  "operationId": "eventListUsingGET",
  "parameters": [
    {
      "name": "series_symbol",
      "in": "query",
      "description": "series symbol",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXNCAAF"
    },
    {
      "name": "status",
      "in": "query",
      "description": "status, ACTIVE=event is open for trading; INACTIVE=event is closed or settled, default:ACTIVE",
      "required": false,
      "schema": {
        "type": "string",
        "default": "ACTIVE"
      },
      "example": "ACTIVE"
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
                      "description": "Event symbol",
                      "example": "KXNCAAF-27"
                    },
                    "name": {
                      "type": "string",
                      "description": "Event name",
                      "example": "College Football National Championship Winner"
                    },
                    "status": {
                      "type": "string",
                      "description": "status, eg:ACTIVE=event is open for trading；INACTIVE=event is closed or settled, default:ACTIVE",
                      "example": "ACTIVE"
                    },
                    "series_id": {
                      "type": "integer",
                      "description": "Series ID",
                      "format": "int32",
                      "example": 152917
                    },
                    "event_id": {
                      "type": "integer",
                      "description": "Event ID",
                      "format": "int32",
                      "example": 152006
                    },
                    "short_name": {
                      "type": "string",
                      "description": "Event short name",
                      "example": "2026-27"
                    },
                    "strike_date": {
                      "type": "string",
                      "description": "strike date,Only present for sports category events. Represents the settlement date/period of the event.",
                      "example": "2026-10-28"
                    },
                    "strike_period": {
                      "type": "string",
                      "description": "strike period, Only present for sports category events. Represents the settlement date/period of the event."
                    },
                    "mutually_exclusive": {
                      "type": "boolean",
                      "description": "mutually exclusive, true or false",
                      "example": true
                    }
                  },
                  "description": "Data list",
                  "title": "EventVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "EventVoPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Events",
    "description": {
      "content": "Retrieves a list of tradable events under a series. Each event represents a specific question or market (e.g., '2026-27 College Football National Championship Winner'). Filter by series_symbol to get events for a specific series.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "event-contracts",
        "events",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "series symbol",
            "type": "text/plain"
          },
          "key": "series_symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "status, ACTIVE=event is open for trading; INACTIVE=event is closed or settled, default:ACTIVE",
            "type": "text/plain"
          },
          "key": "status",
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

## Event Contract Milestones

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/milestones-using-get.md>

### List Event Milestones

Retrieves a paginated list of milestones (individual game or economic release events). Each milestone contains match details, team info, and related event contract symbols.

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
  "path": "/market-data/instruments/event-contracts/milestones/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves a paginated list of milestones (individual game or economic release events). Each milestone contains match details, team info, and related event contract symbols.",
  "operationId": "milestonesUsingGET",
  "parameters": [
    {
      "name": "minimum_start_date",
      "in": "query",
      "description": "Query data with match start time (millisecond timestamp) that is later than that time.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64"
      },
      "example": 1715100000000
    },
    {
      "name": "category",
      "in": "query",
      "description": "category code",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "SPORTS"
    },
    {
      "name": "competition",
      "in": "query",
      "description": "competition name",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "College Football"
    },
    {
      "name": "related_event_symbol",
      "in": "query",
      "description": "related event symbol",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXNCAAFGAME-26AUG29UNCTCU"
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
                    "category": {
                      "type": "string",
                      "description": "Category Code",
                      "example": "SPORTS"
                    },
                    "type": {
                      "type": "string",
                      "description": "Milestone type",
                      "example": "football_game"
                    },
                    "title": {
                      "type": "string",
                      "description": "Milestone title",
                      "example": "North Carolina at TCU"
                    },
                    "details": {
                      "type": "object",
                      "additionalProperties": {
                        "type": "object",
                        "description": "Detail fields vary by type. Refer to SportsGameDetails when type=$SPORTS_GAME, EconomicReleaseDetails when type=ECONOMIC_RELEASE.",
                        "oneOf": [
                          {
                            "type": "object",
                            "properties": {
                              "venue": {
                                "type": "string",
                                "description": "Venue of the game",
                                "example": "Amon G. Carter Stadium"
                              },
                              "home_team_id": {
                                "type": "string",
                                "description": "Home team ID",
                                "example": "NCAAF:TCU"
                              },
                              "home_team_name": {
                                "type": "string",
                                "description": "Home team full name",
                                "example": "TCU Horned Frogs"
                              },
                              "home_team_short_name": {
                                "type": "string",
                                "description": "Home team short name",
                                "example": "TCU"
                              },
                              "away_team_id": {
                                "type": "string",
                                "description": "Away team ID",
                                "example": "NCAAF:UNC"
                              },
                              "away_team_name": {
                                "type": "string",
                                "description": "Away team full name",
                                "example": "North Carolina Tar Heels"
                              },
                              "away_team_short_name": {
                                "type": "string",
                                "description": "Away team short name",
                                "example": "UNC"
                              }
                            },
                            "description": "Details when type=SPORTS_GAME",
                            "title": "SportsGameDetails"
                          },
                          {
                            "type": "object",
                            "properties": {
                              "source": {
                                "type": "string",
                                "description": "Data source institution",
                                "example": "Bureau of Labor Statistics"
                              },
                              "indicator": {
                                "type": "string",
                                "description": "Economic indicator name",
                                "example": "Non-Farm Payrolls"
                              }
                            },
                            "description": "Details when type=ECONOMIC_RELEASE",
                            "title": "EconomicReleaseDetails"
                          }
                        ]
                      },
                      "description": "Detail fields vary by type. Refer to SportsGameDetails when type=$SPORTS_GAME, EconomicReleaseDetails when type=ECONOMIC_RELEASE."
                    },
                    "status": {
                      "type": "string",
                      "description": "status, eg:NOT_STARTED,INPROGRESS,CLOSED,CANCELLED,POSTPONED,DELAYED,SUSPENDED,UNKNOWN",
                      "example": "not_started"
                    },
                    "milestone_id": {
                      "type": "string",
                      "description": "Milestone Unique Identifier",
                      "example": "2526d351-9ffc-4182-8414-495cc84a0b64"
                    },
                    "start_date": {
                      "type": "string",
                      "description": "Start date",
                      "example": "2026-08-29T16:00:00Z"
                    },
                    "end_date": {
                      "type": "string",
                      "description": "End date",
                      "example": "2026-08-29T16:00:00Z"
                    },
                    "related_event_symbols": {
                      "type": "array",
                      "description": "related event symbol",
                      "example": "KXNCAAFGAME-26AUG29UNCTCU",
                      "items": {
                        "type": "string",
                        "description": "related event symbol",
                        "example": "KXNCAAFGAME-26AUG29UNCTCU"
                      }
                    },
                    "primary_event_symbols": {
                      "type": "array",
                      "description": "primary event symbol",
                      "example": "KXNCAAFGAME-26AUG29UNCTCU",
                      "items": {
                        "type": "string",
                        "description": "primary event symbol",
                        "example": "KXNCAAFGAME-26AUG29UNCTCU"
                      }
                    }
                  },
                  "description": "Data list",
                  "title": "MilestoneVo"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "MilestoneVoPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Event Milestones",
    "description": {
      "content": "Retrieves a paginated list of milestones (individual game or economic release events). Each milestone contains match details, team info, and related event contract symbols.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "event-contracts",
        "milestones",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "Query data with match start time (millisecond timestamp) that is later than that time.",
            "type": "text/plain"
          },
          "key": "minimum_start_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "category code",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "competition name",
            "type": "text/plain"
          },
          "key": "competition",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "related event symbol",
            "type": "text/plain"
          },
          "key": "related_event_symbol",
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

## Event Contract Series List

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/series-list-using-get.md>

### List Event Series

Retrieves a paginated list of event contract series. A series represents a recurring competition or event category (e.g., CFP National Champion). Use category and tags to filter by sport type.

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
  "path": "/market-data/instruments/event-contracts/series/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves a paginated list of event contract series. A series represents a recurring competition or event category (e.g., CFP National Champion). Use category and tags to filter by sport type.",
  "operationId": "seriesListUsingGET",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "category code",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "SPORTS"
    },
    {
      "name": "tags",
      "in": "query",
      "description": "tag name",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "Baseball,Football"
    },
    {
      "name": "symbols",
      "in": "query",
      "description": "series symbols",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "KXNCAAF, KXHEISMAN"
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
                  "required": [
                    "frequency"
                  ],
                  "type": "object",
                  "properties": {
                    "symbol": {
                      "type": "string",
                      "description": "Series Symbol",
                      "example": "KXNCAAF"
                    },
                    "name": {
                      "type": "string",
                      "description": "Series Name",
                      "example": "CFP National Champion"
                    },
                    "category": {
                      "type": "string",
                      "description": "Category Code",
                      "example": "SPORTS"
                    },
                    "frequency": {
                      "type": "string",
                      "description": "frequency, e.g., HOURLY,DAILY,WEEKLY,MONTHLY,ANNUAL,ONE_OFF,CUSTOM.",
                      "example": "ANNUAL"
                    },
                    "tags": {
                      "type": "array",
                      "description": "tags name",
                      "example": "['Football']",
                      "items": {
                        "type": "string",
                        "description": "tags name",
                        "example": "['Football']"
                      }
                    },
                    "series_id": {
                      "type": "integer",
                      "description": "Series ID",
                      "format": "int32",
                      "example": 152917
                    }
                  },
                  "description": "Series Information",
                  "title": "EventSeries"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. null means no more data.",
                "example": "eyJ2IjoxLCJwYWdlSW5kZXgiOjJ9"
              }
            },
            "title": "EventSeriesPageResponse"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Event Series",
    "description": {
      "content": "Retrieves a paginated list of event contract series. A series represents a recurring competition or event category (e.g., CFP National Champion). Use category and tags to filter by sport type.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "event-contracts",
        "series",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "category code",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "tag name",
            "type": "text/plain"
          },
          "key": "tags",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "series symbols",
            "type": "text/plain"
          },
          "key": "symbols",
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

## Event Contract Sports Filters

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/sports-filter-using-get.md>

### List Sports Filters

Retrieves available filter options for sports event contracts, including sport tags, competitions, and scopes. Use this to populate filter UI or discover available sports categories before querying series or events.

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
  "path": "/market-data/instruments/event-contracts/sports-filters/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves available filter options for sports event contracts, including sport tags, competitions, and scopes. Use this to populate filter UI or discover available sports categories before querying series or events.",
  "operationId": "sportsFilterUsingGET",
  "parameters": [
    {
      "name": "tag",
      "in": "query",
      "description": "tag name",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": "Football"
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
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "tag": {
                  "type": "string",
                  "description": "Tag name",
                  "example": "Basketball"
                },
                "competitions": {
                  "type": "array",
                  "description": "competition list",
                  "example": "College Football,Pro Football",
                  "items": {
                    "type": "object",
                    "properties": {
                      "competition": {
                        "type": "string",
                        "description": "Competition name (e.g., College Football, Pro Football).",
                        "example": "College Football"
                      },
                      "scopes": {
                        "type": "array",
                        "description": "scope list",
                        "example": "['Games', 'Futures','Awards']",
                        "items": {
                          "type": "string",
                          "description": "scope list",
                          "example": "['Games', 'Futures','Awards']"
                        }
                      }
                    },
                    "description": "Competition Information",
                    "example": "College Football,Pro Football",
                    "title": "Competition"
                  }
                },
                "scopes": {
                  "type": "array",
                  "description": "List of available scope categories under this tag (e.g., Games, Futures, Awards).",
                  "example": "['Games', 'Futures','Awards']",
                  "items": {
                    "type": "string",
                    "description": "List of available scope categories under this tag (e.g., Games, Futures, Awards).",
                    "example": "['Games', 'Futures','Awards']"
                  }
                }
              },
              "description": "Tag Information",
              "title": "EventTag"
            }
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Sports Filters",
    "description": {
      "content": "Retrieves available filter options for sports event contracts, including sport tags, competitions, and scopes. Use this to populate filter UI or discover available sports categories before querying series or events.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "instruments",
        "event-contracts",
        "sports-filters",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "tag name",
            "type": "text/plain"
          },
          "key": "tag",
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

## Event Game Stats

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-game-stats-using-get.md>

### Get Game Stats

Retrieves detailed play-by-play or drive-by-drive game statistics for a live sporting event. The response is a flat structure (wide table) containing fields for all sport types. Only fields relevant to the queried milestone's sport will be populated; others will be absent.

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
  "path": "/market-data/event-contracts/game-stats/get",
  "method": "get",
  "tags": [
    "Event Contract Market Data"
  ],
  "description": "Retrieves detailed play-by-play or drive-by-drive game statistics for a live sporting event. The response is a flat structure (wide table) containing fields for all sport types. Only fields relevant to the queried milestone's sport will be populated; others will be absent.",
  "operationId": "eventGameStatsUsingGET",
  "parameters": [
    {
      "name": "milestone_id",
      "in": "query",
      "description": "Milestone ID. A milestone represents a specific game or real-world occurrence tied to events.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "MS-NBA-20250514-LAL-BOS"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ],
        "default": "US_EVENT"
      },
      "example": "US_EVENT"
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
              "milestone_id"
            ],
            "type": "object",
            "properties": {
              "milestone_id": {
                "type": "string",
                "description": "Milestone ID",
                "example": "MS-NBA-20250514-LAL-BOS"
              },
              "periods": {
                "type": "array",
                "description": "Game periods/innings list",
                "items": {
                  "type": "object",
                  "properties": {
                    "period_number": {
                      "type": "integer",
                      "description": "Period/inning number",
                      "format": "int32",
                      "example": 1
                    },
                    "period_type": {
                      "type": "string",
                      "description": "Period type, e.g. quarter, top, bottom, period, half. Varies by sport.",
                      "example": "quarter"
                    },
                    "events": {
                      "type": "array",
                      "description": "Event list for this period",
                      "items": {
                        "type": "object",
                        "properties": {
                          "type": {
                            "type": "string",
                            "description": "Event type identifier, e.g. football_drive, basketball_play, baseball_play, hockey_play, soccer_event.",
                            "example": "basketball_play"
                          },
                          "description": {
                            "type": "string",
                            "description": "Event description.",
                            "example": "LeBron James makes 3-point jump shot"
                          },
                          "attribution": {
                            "type": "string",
                            "description": "Attribution team/player ID.",
                            "example": "LAL-001"
                          },
                          "away_points": {
                            "type": "integer",
                            "description": "Away team points at this event.",
                            "format": "int32",
                            "example": 3
                          },
                          "home_points": {
                            "type": "integer",
                            "description": "Home team points at this event.",
                            "format": "int32",
                            "example": 0
                          },
                          "clock": {
                            "type": "string",
                            "description": "Game clock.",
                            "example": "11:24"
                          },
                          "possession": {
                            "type": "string",
                            "description": "Possession team ID.",
                            "example": "BOS-001"
                          },
                          "event_type": {
                            "type": "string",
                            "description": "Event sub-type (e.g. three_point_made, field_goal_made, goal, yellow_card). May appear in basketball, soccer.",
                            "example": "three_point_made"
                          },
                          "wall_clock": {
                            "type": "integer",
                            "description": "Event wall clock UTC timestamp in seconds. May appear in basketball.",
                            "format": "int64",
                            "example": 1778640099
                          },
                          "half": {
                            "type": "string",
                            "description": "Half inning indicator. May appear in baseball.",
                            "example": "top"
                          },
                          "strength": {
                            "type": "string",
                            "description": "Strength status, e.g. even/powerplay/shorthanded. May appear in hockey.",
                            "example": "powerplay"
                          },
                          "competitor": {
                            "type": "string",
                            "description": "Competitor side: home/away. May appear in soccer.",
                            "example": "home"
                          },
                          "match_time": {
                            "type": "integer",
                            "description": "Match time in minutes. May appear in soccer.",
                            "format": "int32",
                            "example": 22
                          },
                          "player_name": {
                            "type": "string",
                            "description": "Player name. May appear in soccer.",
                            "example": "Vinicius Jr"
                          },
                          "stoppage_time": {
                            "type": "integer",
                            "description": "Stoppage time in minutes. May appear in soccer.",
                            "format": "int32",
                            "example": 0
                          },
                          "def_points": {
                            "type": "integer",
                            "description": "Defensive points. May appear in football.",
                            "format": "int32",
                            "example": 0
                          },
                          "end_reason": {
                            "type": "string",
                            "description": "Drive end reason. May appear in football.",
                            "example": "touchdown"
                          },
                          "gain": {
                            "type": "integer",
                            "description": "Drive gain in yards. May appear in football.",
                            "format": "int32",
                            "example": 75
                          },
                          "off_points": {
                            "type": "integer",
                            "description": "Offensive points. May appear in football.",
                            "format": "int32",
                            "example": 7
                          },
                          "play_count": {
                            "type": "integer",
                            "description": "Number of plays in drive. May appear in football.",
                            "format": "int32",
                            "example": 8
                          },
                          "plays": {
                            "type": "array",
                            "description": "Play details list. May appear in football.",
                            "items": {
                              "type": "object",
                              "properties": {
                                "type": {
                                  "type": "string",
                                  "description": "Play type identifier.",
                                  "example": "football_play"
                                },
                                "clock": {
                                  "type": "string",
                                  "description": "Game clock.",
                                  "example": "15:00"
                                },
                                "description": {
                                  "type": "string",
                                  "description": "Play description.",
                                  "example": "Jalen Hurts pass to A.J. Brown for 22 yards"
                                },
                                "down": {
                                  "type": "integer",
                                  "description": "Current down number.",
                                  "format": "int32",
                                  "example": 1
                                },
                                "yfd": {
                                  "type": "integer",
                                  "description": "Yards to first down.",
                                  "format": "int32",
                                  "example": 10
                                }
                              },
                              "description": "EventContractGameStatsPlayVo",
                              "title": "EventContractGameStatsPlayVo"
                            }
                          }
                        },
                        "description": "EventContractGameStatsEventVo",
                        "title": "EventContractGameStatsEventVo"
                      }
                    },
                    "period_name": {
                      "type": "string",
                      "description": "Period name. May appear in baseball, soccer.",
                      "example": "Top 1st"
                    },
                    "attribution": {
                      "type": "string",
                      "description": "Attribution team ID. May appear in baseball.",
                      "example": "NYY-001"
                    },
                    "half": {
                      "type": "string",
                      "description": "Half inning indicator. May appear in baseball.",
                      "example": "top"
                    },
                    "away_score": {
                      "type": "integer",
                      "description": "Away team current score. May appear in baseball, hockey, soccer.",
                      "format": "int32",
                      "example": 2
                    },
                    "away_team_id": {
                      "type": "string",
                      "description": "Away team ID. May appear in baseball, hockey, soccer.",
                      "example": "NYY-001"
                    },
                    "home_score": {
                      "type": "integer",
                      "description": "Home team current score. May appear in baseball, hockey, soccer.",
                      "format": "int32",
                      "example": 0
                    },
                    "home_team_id": {
                      "type": "string",
                      "description": "Home team ID. May appear in baseball, hockey, soccer.",
                      "example": "BOS-001"
                    }
                  },
                  "description": "EventContractGameStatsPeriodVo",
                  "title": "EventContractGameStatsPeriodVo"
                }
              }
            },
            "description": "EventContractGameStatsVo",
            "title": "EventContractGameStatsVo"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNAUTHORIZED",
            "message": "Insufficient permission"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNSUPPORTED_CATEGORY",
            "message": "Unsupported category:US_EVENTS"
          }
        }
      }
    }
  },
  "postman": {
    "name": "Get Game Stats",
    "description": {
      "content": "Retrieves detailed play-by-play or drive-by-drive game statistics for a live sporting event. The response is a flat structure (wide table) containing fields for all sport types. Only fields relevant to the queried milestone's sport will be populated; others will be absent.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
        "game-stats",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Milestone ID. A milestone represents a specific game or real-world occurrence tied to events.",
            "type": "text/plain"
          },
          "key": "milestone_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
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

## Event Live Data

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-live-data-using-get.md>

### Get Event Live Data

Retrieves real-time live game/event data for a specified milestone, including scores, game clock, period, and winner information. Use this to display live match status alongside event contract prices.

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
  "path": "/market-data/event-contracts/live-data/get",
  "method": "get",
  "tags": [
    "Event Contract Market Data"
  ],
  "description": "Retrieves real-time live game/event data for a specified milestone, including scores, game clock, period, and winner information. Use this to display live match status alongside event contract prices.",
  "operationId": "eventLiveDataUsingGET",
  "parameters": [
    {
      "name": "milestone_id",
      "in": "query",
      "description": "Milestone ID. A milestone represents a specific game or real-world occurrence tied to events.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "MS-NBA-20250514-LAL-BOS"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ],
        "default": "US_EVENT"
      },
      "example": "US_EVENT"
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
              "milestone_id",
              "status",
              "type"
            ],
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "description": "Sport type identifier, e.g. basketball_game, baseball_game, baseball_tournament, football_game, hockey_match, hockey_tournament, soccer_tournament_multi_leg, tennis_tournament_singles, golf_tournament, cricket_match.",
                "example": "basketball_game"
              },
              "milestone_id": {
                "type": "string",
                "description": "Milestone ID.",
                "example": "MS-NBA-20250514-LAL-BOS"
              },
              "status": {
                "type": "string",
                "description": "Game status. Values: NOT_STARTED (not begun), INPROGRESS (underway), CLOSED (concluded with result), CANCELLED (cancelled, contracts may be voided), POSTPONED (delayed to future date), DELAYED (temporarily paused, expected to resume), SUSPENDED (indefinitely paused, outcome pending).",
                "example": "INPROGRESS"
              },
              "winner": {
                "type": "string",
                "description": "Winner identifier. Empty string or absent if not yet determined.",
                "example": "LAL"
              },
              "last_play": {
                "type": "object",
                "properties": {
                  "description": {
                    "type": "string",
                    "description": "Last play description.",
                    "example": "LeBron James makes 3-point jump shot"
                  },
                  "occurence_ts": {
                    "type": "integer",
                    "description": "Occurrence timestamp in seconds.",
                    "format": "int64",
                    "example": 1778640200
                  }
                },
                "description": "EventContractLiveLastPlayVo",
                "title": "EventContractLiveLastPlayVo"
              },
              "last_updated_ts": {
                "type": "integer",
                "description": "Last updated timestamp in seconds.",
                "format": "int64",
                "example": 1778640200
              },
              "details": {
                "type": "object",
                "description": "Sport-specific details. Structure varies by type field.<br/> Type mapping:<br/>• basketball_game → EventContractLiveBasketballDetailsVo<br/>• baseball_game → EventContractLiveBaseballDetailsVo<br/>• baseball_tournament → EventContractLiveBaseballDetailsVo<br/>• football_game → EventContractLiveFootballDetailsVo<br/>• hockey_match → EventContractLiveHockeyDetailsVo<br/>• hockey_tournament → EventContractLiveHockeyDetailsVo<br/>• soccer_tournament_multi_leg → EventContractLiveSoccerDetailsVo<br/>• tennis_tournament_singles → EventContractLiveTennisDetailsVo<br/>• golf_tournament → EventContractLiveGolfDetailsVo<br/>• cricket_match → EventContractLiveCricketDetailsVo",
                "oneOf": [
                  {
                    "type": "object",
                    "properties": {
                      "away_points": {
                        "type": "integer",
                        "description": "Away team score.",
                        "format": "int32",
                        "example": 87
                      },
                      "home_points": {
                        "type": "integer",
                        "description": "Home team score.",
                        "format": "int32",
                        "example": 92
                      },
                      "last_event_created_ts": {
                        "type": "integer",
                        "description": "Last event created timestamp in seconds.",
                        "format": "int64",
                        "example": 1778640200
                      },
                      "last_event_is_timeout": {
                        "type": "boolean",
                        "description": "Whether last event is a timeout.",
                        "example": false
                      },
                      "period": {
                        "type": "integer",
                        "description": "Current period number.",
                        "format": "int32",
                        "example": 3
                      },
                      "period_remaining_time": {
                        "type": "string",
                        "description": "Remaining time in current period.",
                        "example": "04:32"
                      },
                      "period_type": {
                        "type": "string",
                        "description": "Period type, e.g. quarter, half, overtime. Varies by sport.",
                        "example": "quarter"
                      },
                      "possession": {
                        "type": "string",
                        "description": "Ball possession side, e.g. away, home.",
                        "example": "home"
                      },
                      "score_last_updated_ts": {
                        "type": "integer",
                        "description": "Score last updated timestamp in seconds.",
                        "format": "int64",
                        "example": 1778640200
                      }
                    },
                    "description": "Basketball game details (type=basketball_game)",
                    "title": "EventContractLiveBasketballDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "away_points": {
                        "type": "integer",
                        "description": "Away team score.",
                        "format": "int32",
                        "example": 3
                      },
                      "home_points": {
                        "type": "integer",
                        "description": "Home team score.",
                        "format": "int32",
                        "example": 5
                      },
                      "balls": {
                        "type": "integer",
                        "description": "Current ball count.",
                        "format": "int32",
                        "example": 2
                      },
                      "bases": {
                        "type": "array",
                        "description": "Base occupancy. Array of 4 booleans: [home, 1st, 2nd, 3rd].",
                        "example": [
                          false,
                          true,
                          false,
                          true
                        ],
                        "items": {
                          "type": "boolean",
                          "description": "Base occupancy. Array of 4 booleans: [home, 1st, 2nd, 3rd].",
                          "example": false
                        }
                      },
                      "inning": {
                        "type": "integer",
                        "description": "Current inning number.",
                        "format": "int32",
                        "example": 6
                      },
                      "inning_half": {
                        "type": "integer",
                        "description": "Current half inning. 0=top, 1=bottom.",
                        "format": "int32",
                        "example": 0
                      },
                      "outs": {
                        "type": "integer",
                        "description": "Out count.",
                        "format": "int32",
                        "example": 1
                      },
                      "strikes": {
                        "type": "integer",
                        "description": "Strike count.",
                        "format": "int32",
                        "example": 1
                      }
                    },
                    "description": "Baseball game details (type=baseball_game, baseball_tournament)",
                    "title": "EventContractLiveBaseballDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "away_points": {
                        "type": "integer",
                        "description": "Away team score.",
                        "format": "int32",
                        "example": 14
                      },
                      "home_points": {
                        "type": "integer",
                        "description": "Home team score.",
                        "format": "int32",
                        "example": 21
                      },
                      "clock": {
                        "type": "string",
                        "description": "Game clock.",
                        "example": "07:23"
                      },
                      "quarter": {
                        "type": "integer",
                        "description": "Current quarter number.",
                        "format": "int32",
                        "example": 3
                      },
                      "situation": {
                        "type": "object",
                        "properties": {
                          "down": {
                            "type": "integer",
                            "description": "Current down number.",
                            "format": "int32",
                            "example": 2
                          },
                          "goal_to_go": {
                            "type": "boolean",
                            "description": "Whether it is goal to go.",
                            "example": false
                          },
                          "possession_team_id": {
                            "type": "string",
                            "description": "Possession team ID.",
                            "example": "PHI-001"
                          },
                          "side_team_id": {
                            "type": "string",
                            "description": "Ball side team ID.",
                            "example": "KC-001"
                          },
                          "yardline": {
                            "type": "integer",
                            "description": "Yard line position.",
                            "format": "int32",
                            "example": 35
                          },
                          "yfd": {
                            "type": "integer",
                            "description": "Yards to first down.",
                            "format": "int32",
                            "example": 7
                          }
                        },
                        "description": "EventContractLiveFootballSituationVo",
                        "title": "EventContractLiveFootballSituationVo"
                      }
                    },
                    "description": "Football game details (type=football_game)",
                    "title": "EventContractLiveFootballDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "away_points": {
                        "type": "integer",
                        "description": "Away team score.",
                        "format": "int32",
                        "example": 2
                      },
                      "home_points": {
                        "type": "integer",
                        "description": "Home team score.",
                        "format": "int32",
                        "example": 3
                      },
                      "period": {
                        "type": "integer",
                        "description": "Current period number.",
                        "format": "int32",
                        "example": 2
                      },
                      "period_remaining_time": {
                        "type": "string",
                        "description": "Remaining time in current period.",
                        "example": "08:15"
                      }
                    },
                    "description": "Hockey match details (type=hockey_match, hockey_tournament)",
                    "title": "EventContractLiveHockeyDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "aggregate_text": {
                        "type": "string",
                        "description": "Aggregate score text.",
                        "example": "Agg: 3-2"
                      },
                      "away_aggregate_score": {
                        "type": "integer",
                        "description": "Away team aggregate score (multi-leg).",
                        "format": "int32",
                        "example": 2
                      },
                      "away_penalties": {
                        "type": "array",
                        "description": "Away team penalty results list.",
                        "items": {
                          "type": "string",
                          "description": "Away team penalty results list."
                        }
                      },
                      "away_same_game_score": {
                        "type": "integer",
                        "description": "Away team current game score.",
                        "format": "int32",
                        "example": 1
                      },
                      "away_significant_events": {
                        "type": "array",
                        "description": "Away team significant events list.",
                        "items": {
                          "type": "object",
                          "properties": {
                            "event_type": {
                              "type": "string",
                              "description": "Event type, e.g. goal, yellow_card, red_card.",
                              "example": "goal"
                            },
                            "player": {
                              "type": "string",
                              "description": "Player name.",
                              "example": "Vinicius Jr"
                            },
                            "time": {
                              "type": "string",
                              "description": "Event time.",
                              "example": "22'"
                            }
                          },
                          "description": "EventContractLiveSoccerSignificantEventVo",
                          "title": "EventContractLiveSoccerSignificantEventVo"
                        }
                      },
                      "half": {
                        "type": "string",
                        "description": "Current half, e.g. 1H, 2H, FT.",
                        "example": "2H"
                      },
                      "home_aggregate_score": {
                        "type": "integer",
                        "description": "Home team aggregate score (multi-leg).",
                        "format": "int32",
                        "example": 3
                      },
                      "home_penalties": {
                        "type": "array",
                        "description": "Home team penalty results list.",
                        "items": {
                          "type": "string",
                          "description": "Home team penalty results list."
                        }
                      },
                      "home_same_game_score": {
                        "type": "integer",
                        "description": "Home team current game score.",
                        "format": "int32",
                        "example": 2
                      },
                      "home_significant_events": {
                        "type": "array",
                        "description": "Home team significant events list.",
                        "items": {
                          "type": "object",
                          "properties": {
                            "event_type": {
                              "type": "string",
                              "description": "Event type, e.g. goal, yellow_card, red_card.",
                              "example": "goal"
                            },
                            "player": {
                              "type": "string",
                              "description": "Player name.",
                              "example": "Vinicius Jr"
                            },
                            "time": {
                              "type": "string",
                              "description": "Event time.",
                              "example": "22'"
                            }
                          },
                          "description": "EventContractLiveSoccerSignificantEventVo",
                          "title": "EventContractLiveSoccerSignificantEventVo"
                        }
                      },
                      "penalties_text": {
                        "type": "string",
                        "description": "Penalty shootout text."
                      },
                      "show_penalties": {
                        "type": "boolean",
                        "description": "Whether to show penalty shootout.",
                        "example": false
                      },
                      "status_text": {
                        "type": "string",
                        "description": "Status text.",
                        "example": "2nd Half"
                      },
                      "time": {
                        "type": "string",
                        "description": "Match time.",
                        "example": "62"
                      }
                    },
                    "description": "Soccer match details (type=soccer_tournament_multi_leg)",
                    "title": "EventContractLiveSoccerDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "advantage": {
                        "type": "string",
                        "description": "Advantage holder ID. Empty if deuce."
                      },
                      "away_overall_score": {
                        "type": "integer",
                        "description": "Away player overall score.",
                        "format": "int32",
                        "example": 1
                      },
                      "competitor1_current_round_score": {
                        "type": "integer",
                        "description": "Competitor 1 current round score.",
                        "format": "int32",
                        "example": 4
                      },
                      "competitor1_id": {
                        "type": "string",
                        "description": "Competitor 1 ID.",
                        "example": "SINNER-001"
                      },
                      "competitor1_is_home": {
                        "type": "boolean",
                        "description": "Whether competitor 1 is home.",
                        "example": true
                      },
                      "competitor1_overall_score": {
                        "type": "integer",
                        "description": "Competitor 1 overall score (sets won).",
                        "format": "int32",
                        "example": 2
                      },
                      "competitor1_round_scores": {
                        "type": "array",
                        "description": "Competitor 1 round scores.",
                        "items": {
                          "type": "object",
                          "properties": {
                            "outcome": {
                              "type": "string",
                              "description": "Round outcome. Values: won, lost, or empty (in progress).",
                              "example": "won"
                            },
                            "score": {
                              "type": "integer",
                              "description": "Games won in this set.",
                              "format": "int32",
                              "example": 6
                            },
                            "tiebreak_score": {
                              "type": "integer",
                              "description": "Tiebreak score. Null if no tiebreak.",
                              "format": "int32"
                            }
                          },
                          "description": "EventContractLiveTennisRoundScoreVo",
                          "title": "EventContractLiveTennisRoundScoreVo"
                        }
                      },
                      "competitor1_seed": {
                        "type": "integer",
                        "description": "Competitor 1 seed ranking.",
                        "format": "int32",
                        "example": 1
                      },
                      "competitor1_statistics": {
                        "type": "object",
                        "properties": {
                          "gamescore": {
                            "type": "string",
                            "description": "Current game score.",
                            "example": "30"
                          }
                        },
                        "description": "EventContractLiveTennisStatisticsVo",
                        "title": "EventContractLiveTennisStatisticsVo"
                      },
                      "competitor2_current_round_score": {
                        "type": "integer",
                        "description": "Competitor 2 current round score.",
                        "format": "int32",
                        "example": 3
                      },
                      "competitor2_id": {
                        "type": "string",
                        "description": "Competitor 2 ID.",
                        "example": "ALCARAZ-001"
                      },
                      "competitor2_overall_score": {
                        "type": "integer",
                        "description": "Competitor 2 overall score (sets won).",
                        "format": "int32",
                        "example": 1
                      },
                      "competitor2_round_scores": {
                        "type": "array",
                        "description": "Competitor 2 round scores.",
                        "items": {
                          "type": "object",
                          "properties": {
                            "outcome": {
                              "type": "string",
                              "description": "Round outcome. Values: won, lost, or empty (in progress).",
                              "example": "won"
                            },
                            "score": {
                              "type": "integer",
                              "description": "Games won in this set.",
                              "format": "int32",
                              "example": 6
                            },
                            "tiebreak_score": {
                              "type": "integer",
                              "description": "Tiebreak score. Null if no tiebreak.",
                              "format": "int32"
                            }
                          },
                          "description": "EventContractLiveTennisRoundScoreVo",
                          "title": "EventContractLiveTennisRoundScoreVo"
                        }
                      },
                      "competitor2_seed": {
                        "type": "integer",
                        "description": "Competitor 2 seed ranking.",
                        "format": "int32",
                        "example": 2
                      },
                      "competitor2_statistics": {
                        "type": "object",
                        "properties": {
                          "gamescore": {
                            "type": "string",
                            "description": "Current game score.",
                            "example": "30"
                          }
                        },
                        "description": "EventContractLiveTennisStatisticsVo",
                        "title": "EventContractLiveTennisStatisticsVo"
                      },
                      "completed_rounds": {
                        "type": "integer",
                        "description": "Number of completed rounds.",
                        "format": "int32",
                        "example": 2
                      },
                      "home_overall_score": {
                        "type": "integer",
                        "description": "Home player overall score.",
                        "format": "int32",
                        "example": 2
                      },
                      "round_winners": {
                        "type": "array",
                        "description": "Round winners ID list.",
                        "items": {
                          "type": "string",
                          "description": "Round winners ID list."
                        }
                      },
                      "server": {
                        "type": "string",
                        "description": "Current server ID.",
                        "example": "SINNER-001"
                      }
                    },
                    "description": "Tennis match details (type=tennis_tournament_singles)",
                    "title": "EventContractLiveTennisDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "current_round": {
                        "type": "integer",
                        "description": "Current round number.",
                        "format": "int32",
                        "example": 3
                      },
                      "leaderboard": {
                        "type": "array",
                        "description": "Leaderboard entries.",
                        "items": {
                          "type": "object",
                          "properties": {
                            "competitor_id": {
                              "type": "string",
                              "description": "Competitor ID.",
                              "example": "SCHEFFLER-001"
                            },
                            "current_round_score": {
                              "type": "integer",
                              "description": "Current round score (strokes relative to par).",
                              "format": "int32",
                              "example": -4
                            },
                            "current_round_thru": {
                              "type": "integer",
                              "description": "Holes completed in current round.",
                              "format": "int32",
                              "example": 14
                            },
                            "finished_current_round": {
                              "type": "boolean",
                              "description": "Whether current round is finished.",
                              "example": false
                            },
                            "position": {
                              "type": "integer",
                              "description": "Current leaderboard position.",
                              "format": "int32",
                              "example": 1
                            },
                            "started_current_round": {
                              "type": "boolean",
                              "description": "Whether current round has started.",
                              "example": true
                            },
                            "total_score": {
                              "type": "integer",
                              "description": "Total score (strokes relative to par).",
                              "format": "int32",
                              "example": -12
                            }
                          },
                          "description": "EventContractLiveGolfLeaderboardEntryVo",
                          "title": "EventContractLiveGolfLeaderboardEntryVo"
                        }
                      },
                      "round_label": {
                        "type": "string",
                        "description": "Round label.",
                        "example": "Round 3"
                      }
                    },
                    "description": "Golf tournament details (type=golf_tournament)",
                    "title": "EventContractLiveGolfDetailsVo"
                  },
                  {
                    "type": "object",
                    "properties": {
                      "away_overs": {
                        "type": "integer",
                        "description": "Away team overs.",
                        "format": "int32",
                        "example": 18
                      },
                      "away_score": {
                        "type": "integer",
                        "description": "Away team score.",
                        "format": "int32",
                        "example": 156
                      },
                      "away_wickets": {
                        "type": "integer",
                        "description": "Away team wickets.",
                        "format": "int32",
                        "example": 6
                      },
                      "batting": {
                        "type": "string",
                        "description": "Current batting side. Values: home, away.",
                        "example": "home"
                      },
                      "home_overs": {
                        "type": "integer",
                        "description": "Home team overs.",
                        "format": "int32",
                        "example": 12
                      },
                      "home_score": {
                        "type": "integer",
                        "description": "Home team score.",
                        "format": "int32",
                        "example": 98
                      },
                      "home_wickets": {
                        "type": "integer",
                        "description": "Home team wickets.",
                        "format": "int32",
                        "example": 3
                      }
                    },
                    "description": "Cricket match details (type=cricket_match)",
                    "title": "EventContractLiveCricketDetailsVo"
                  }
                ]
              }
            },
            "description": "EventContractLiveDataVo",
            "title": "EventContractLiveDataVo"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNAUTHORIZED",
            "message": "Insufficient permission"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNSUPPORTED_CATEGORY",
            "message": "Unsupported category:US_EVENTS"
          }
        }
      }
    }
  },
  "postman": {
    "name": "Get Event Live Data",
    "description": {
      "content": "Retrieves real-time live game/event data for a specified milestone, including scores, game clock, period, and winner information. Use this to display live match status alongside event contract prices.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
        "live-data",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Milestone ID. A milestone represents a specific game or real-world occurrence tied to events.",
            "type": "text/plain"
          },
          "key": "milestone_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
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

## Event Market Bars

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-bars-using-get.md>

### List Event Market Bars

Retrieves historical OHLCV bar data for one or more event contract markets, keyed by market symbol. Use this endpoint to build price charts for individual contract markets. Maximum 100 symbols per request.

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
  "path": "/market-data/event-contracts/markets/bars/list",
  "method": "get",
  "tags": [
    "Event Contract Market Data"
  ],
  "description": "Retrieves historical OHLCV bar data for one or more event contract markets, keyed by market symbol. Use this endpoint to build price charts for individual contract markets. Maximum 100 symbols per request.",
  "operationId": "eventMarketBarsUsingGET",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "Comma-separated market symbols. A market is a single binary contract within an event. Maximum 100 symbols.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXNBA-LAL-BOS-0415-ML,KXNBA-LAL-BOS-0415-SPREAD"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ],
        "default": "US_EVENT"
      },
      "example": "US_EVENT"
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "Start time (unix timestamp in milliseconds). Empty means no lower bound.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64"
      },
      "example": 1700000000000
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "End time (unix timestamp in milliseconds). Empty means no upper bound.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64"
      },
      "example": 1700100000000
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of bars. Range: 1-1200, default 200.",
      "required": false,
      "schema": {
        "type": "integer",
        "default": 200
      },
      "example": 200
    },
    {
      "name": "timespan",
      "in": "query",
      "description": "Bar time granularity. M1=1min, M5=5min, M15=15min, M30=30min, M60=1hour, M120=2hour, M240=4hour, D=Daily, W=Weekly, M=Monthly, Y=Yearly.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "M5"
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
      "description": "JSON object where key is market_symbol (string) and value is array of candlestick bar objects. See schema below for bar object fields.",
      "content": {
        "application/json": {
          "schema": {
            "type": "array",
            "items": {
              "required": [
                "end_period_time",
                "volume"
              ],
              "type": "object",
              "properties": {
                "end_period_time": {
                  "type": "string",
                  "description": "Bar end period time (UTC datetime)",
                  "example": "2021-12-28T09:00:09.945+0000"
                },
                "volume": {
                  "type": "string",
                  "description": "Volume during period",
                  "example": "450.0"
                },
                "open": {
                  "type": "string",
                  "description": "Open price (nullable)",
                  "example": "0.53"
                },
                "high": {
                  "type": "string",
                  "description": "High price (nullable)",
                  "example": "0.57"
                },
                "low": {
                  "type": "string",
                  "description": "Low price (nullable)",
                  "example": "0.51"
                },
                "close": {
                  "type": "string",
                  "description": "Close price (nullable)",
                  "example": "0.56"
                }
              },
              "description": "EventContractCandlestickVo",
              "title": "EventContractCandlestickVo"
            }
          },
          "example": {
            "KXNBAGAME-26MAY12MINSAS-MIN": [
              {
                "open": "0.23",
                "high": "0.24",
                "low": "0.18",
                "close": "0.19",
                "end_period_time": "2026-05-12T04:00:00.000+0000",
                "volume": "286472.18"
              }
            ],
            "KXNBAGAME-26MAY12MINSAS-SAS": [
              {
                "open": "0.78",
                "high": "0.95",
                "low": "0.71",
                "close": "0.94",
                "end_period_time": "2026-05-12T04:00:00.000+0000",
                "volume": "149897.77"
              }
            ]
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNAUTHORIZED",
            "message": "Insufficient permission"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNSUPPORTED_CATEGORY",
            "message": "Unsupported category:US_EVENTS"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Event Market Bars",
    "description": {
      "content": "Retrieves historical OHLCV bar data for one or more event contract markets, keyed by market symbol. Use this endpoint to build price charts for individual contract markets. Maximum 100 symbols per request.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
        "markets",
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
            "content": "(Required) Comma-separated market symbols. A market is a single binary contract within an event. Maximum 100 symbols.",
            "type": "text/plain"
          },
          "key": "symbols",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Start time (unix timestamp in milliseconds). Empty means no lower bound.",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "End time (unix timestamp in milliseconds). Empty means no upper bound.",
            "type": "text/plain"
          },
          "key": "end_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of bars. Range: 1-1200, default 200.",
            "type": "text/plain"
          },
          "key": "count",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Bar time granularity. M1=1min, M5=5min, M15=15min, M30=30min, M60=1hour, M120=2hour, M240=4hour, D=Daily, W=Weekly, M=Monthly, Y=Yearly.",
            "type": "text/plain"
          },
          "key": "timespan",
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

## Event Market Bars By Event

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-bars-by-event-using-get.md>

### List Event Bars by Event

Retrieves historical OHLCV bar data for all markets under a given event, keyed by market symbol. Use this endpoint to compare price movements across all contract outcomes within a single event.

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
  "path": "/market-data/event-contracts/markets/bars/list-by-event",
  "method": "get",
  "tags": [
    "Event Contract Market Data"
  ],
  "description": "Retrieves historical OHLCV bar data for all markets under a given event, keyed by market symbol. Use this endpoint to compare price movements across all contract outcomes within a single event.",
  "operationId": "eventMarketBarsByEventUsingGET",
  "parameters": [
    {
      "name": "event_symbol",
      "in": "query",
      "description": "Event unique identifier. An event contains multiple markets (binary contracts).",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXNBA-LAL-BOS-0415"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ],
        "default": "US_EVENT"
      },
      "example": "US_EVENT"
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "Start time (unix timestamp in milliseconds). Empty means no lower bound.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64"
      },
      "example": 1700000000000
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "End time (unix timestamp in milliseconds). Empty means no upper bound.",
      "required": false,
      "schema": {
        "type": "integer",
        "format": "int64"
      },
      "example": 1700100000000
    },
    {
      "name": "count",
      "in": "query",
      "description": "Number of bars. Range: 1-1200, default 200.",
      "required": false,
      "schema": {
        "type": "integer",
        "default": 200
      },
      "example": 200
    },
    {
      "name": "timespan",
      "in": "query",
      "description": "Bar time granularity. M1=1min, M5=5min, M15=15min, M30=30min, M60=1hour, M120=2hour, M240=4hour, D=Daily, W=Weekly, M=Monthly, Y=Yearly.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "M5"
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
      "description": "JSON object where key is market_symbol (string) and value is array of candlestick bar objects. See schema below for bar object fields.",
      "content": {
        "application/json": {
          "schema": {
            "type": "array",
            "items": {
              "required": [
                "end_period_time",
                "volume"
              ],
              "type": "object",
              "properties": {
                "end_period_time": {
                  "type": "string",
                  "description": "Bar end period time (UTC datetime)",
                  "example": "2021-12-28T09:00:09.945+0000"
                },
                "volume": {
                  "type": "string",
                  "description": "Volume during period",
                  "example": "450.0"
                },
                "open": {
                  "type": "string",
                  "description": "Open price (nullable)",
                  "example": "0.53"
                },
                "high": {
                  "type": "string",
                  "description": "High price (nullable)",
                  "example": "0.57"
                },
                "low": {
                  "type": "string",
                  "description": "Low price (nullable)",
                  "example": "0.51"
                },
                "close": {
                  "type": "string",
                  "description": "Close price (nullable)",
                  "example": "0.56"
                }
              },
              "description": "EventContractCandlestickVo",
              "title": "EventContractCandlestickVo"
            }
          },
          "example": {
            "KXNBAGAME-26MAY12MINSAS-MIN": [
              {
                "open": "0.23",
                "high": "0.24",
                "low": "0.18",
                "close": "0.19",
                "end_period_time": "2026-05-12T04:00:00.000+0000",
                "volume": "286472.18"
              }
            ],
            "KXNBAGAME-26MAY12MINSAS-SAS": [
              {
                "open": "0.78",
                "high": "0.95",
                "low": "0.71",
                "close": "0.94",
                "end_period_time": "2026-05-12T04:00:00.000+0000",
                "volume": "149897.77"
              }
            ]
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNAUTHORIZED",
            "message": "Insufficient permission"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNSUPPORTED_CATEGORY",
            "message": "Unsupported category:US_EVENTS"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Event Bars by Event",
    "description": {
      "content": "Retrieves historical OHLCV bar data for all markets under a given event, keyed by market symbol. Use this endpoint to compare price movements across all contract outcomes within a single event.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
        "markets",
        "bars",
        "list-by-event"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Event unique identifier. An event contains multiple markets (binary contracts).",
            "type": "text/plain"
          },
          "key": "event_symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Start time (unix timestamp in milliseconds). Empty means no lower bound.",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "End time (unix timestamp in milliseconds). Empty means no upper bound.",
            "type": "text/plain"
          },
          "key": "end_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Number of bars. Range: 1-1200, default 200.",
            "type": "text/plain"
          },
          "key": "count",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Bar time granularity. M1=1min, M5=5min, M15=15min, M30=30min, M60=1hour, M120=2hour, M240=4hour, D=Daily, W=Weekly, M=Monthly, Y=Yearly.",
            "type": "text/plain"
          },
          "key": "timespan",
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

## Event Market Depth

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-depth-using-get.md>

### List Event Market Depthes

Retrieves the order book (bid/ask depth) for a single event contract market. Each level shows the price and aggregate size (quantity of open orders) at that price point.

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
  "path": "/market-data/event-contracts/markets/depths/list",
  "method": "get",
  "tags": [
    "Event Contract Market Data"
  ],
  "description": "Retrieves the order book (bid/ask depth) for a single event contract market. Each level shows the price and aggregate size (quantity of open orders) at that price point.",
  "operationId": "eventMarketDepthUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Market symbol. A market is a single binary contract within an event.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXNBAGAME-26MAY11OKCLAL-LAL"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ],
        "default": "US_EVENT"
      },
      "example": "US_EVENT"
    },
    {
      "name": "depth",
      "in": "query",
      "description": "Market depth levels. Range: 0-100. Default 0 returns all available levels.",
      "required": false,
      "schema": {
        "type": "integer",
        "default": 0
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
            "required": [
              "instrument_id",
              "no_asks",
              "no_bids",
              "symbol",
              "yes_asks",
              "yes_bids"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Contract symbol",
                "example": "KXNBAGAME-26MAY11OKCLAL-LAL"
              },
              "instrument_id": {
                "type": "string",
                "description": "Contract instrument ID",
                "example": "505640444"
              },
              "yes_asks": {
                "type": "array",
                "description": "Yes side ask orders, sorted by price ascending",
                "items": {
                  "required": [
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "0.19"
                    },
                    "size": {
                      "type": "string",
                      "description": "Size (Quantity)",
                      "example": "3068826.79"
                    }
                  },
                  "description": "EventContractAskBidVo",
                  "title": "EventContractAskBidVo"
                }
              },
              "yes_bids": {
                "type": "array",
                "description": "Yes side bid orders, sorted by price descending",
                "items": {
                  "required": [
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "0.19"
                    },
                    "size": {
                      "type": "string",
                      "description": "Size (Quantity)",
                      "example": "3068826.79"
                    }
                  },
                  "description": "EventContractAskBidVo",
                  "title": "EventContractAskBidVo"
                }
              },
              "no_asks": {
                "type": "array",
                "description": "No side ask orders, sorted by price ascending",
                "items": {
                  "required": [
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "0.19"
                    },
                    "size": {
                      "type": "string",
                      "description": "Size (Quantity)",
                      "example": "3068826.79"
                    }
                  },
                  "description": "EventContractAskBidVo",
                  "title": "EventContractAskBidVo"
                }
              },
              "no_bids": {
                "type": "array",
                "description": "No side bid orders, sorted by price descending",
                "items": {
                  "required": [
                    "price",
                    "size"
                  ],
                  "type": "object",
                  "properties": {
                    "price": {
                      "type": "string",
                      "description": "Price",
                      "example": "0.19"
                    },
                    "size": {
                      "type": "string",
                      "description": "Size (Quantity)",
                      "example": "3068826.79"
                    }
                  },
                  "description": "EventContractAskBidVo",
                  "title": "EventContractAskBidVo"
                }
              }
            },
            "description": "EventContractDepthVo",
            "title": "EventContractDepthVo"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNAUTHORIZED",
            "message": "Insufficient permission"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNSUPPORTED_CATEGORY",
            "message": "Unsupported category:US_EVENTS"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Event Market Depthes",
    "description": {
      "content": "Retrieves the order book (bid/ask depth) for a single event contract market. Each level shows the price and aggregate size (quantity of open orders) at that price point.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
        "markets",
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
            "content": "(Required) Market symbol. A market is a single binary contract within an event.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
            "type": "text/plain"
          },
          "key": "category",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Market depth levels. Range: 0-100. Default 0 returns all available levels.",
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

## Event Market Snapshot

> Source: <https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-snapshot-using-get.md>

### List Event Market Snapshots

Retrieves the latest market snapshot for a single event contract market, including yes/no bid-ask prices, last trade price, volume, open interest, and market status.

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
  "path": "/market-data/event-contracts/markets/snapshots/list",
  "method": "get",
  "tags": [
    "Event Contract Market Data"
  ],
  "description": "Retrieves the latest market snapshot for a single event contract market, including yes/no bid-ask prices, last trade price, volume, open interest, and market status.",
  "operationId": "eventMarketSnapshotUsingGET",
  "parameters": [
    {
      "name": "symbol",
      "in": "query",
      "description": "Market symbol. A market is a single binary contract within an event.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXNBA-LAL-BOS-0415-ML"
    },
    {
      "name": "category",
      "in": "query",
      "description": "Category, default is US_EVENT, currently only US_EVENT is supported.",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "US_EVENT"
        ],
        "default": "US_EVENT"
      },
      "example": "US_EVENT"
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
              "symbol"
            ],
            "type": "object",
            "properties": {
              "symbol": {
                "type": "string",
                "description": "Contract symbol",
                "example": "KXNBA-LAL-BOS-0415-ML"
              },
              "instrument_id": {
                "type": "string",
                "description": "Contract instrument ID",
                "example": "9878278271"
              },
              "event_symbol": {
                "type": "string",
                "description": "Parent event symbol",
                "example": "KXNBA-LAL-BOS-0415"
              },
              "yes_sub_title": {
                "type": "string",
                "description": "Yes side subtitle (e.g. 'Lakers Win')",
                "example": "Lakers Win"
              },
              "no_sub_title": {
                "type": "string",
                "description": "No side subtitle (e.g. 'Celtics Win')",
                "example": "Celtics Win"
              },
              "status": {
                "type": "string",
                "description": "Market status. ACTIVE=Market is open for trading; INACTIVE=Market is closed, no new orders accepted.",
                "example": "ACTIVE"
              },
              "yes_bid": {
                "type": "string",
                "description": "Yes side best bid price",
                "example": "0.56"
              },
              "yes_ask": {
                "type": "string",
                "description": "Yes side best ask price",
                "example": "0.58"
              },
              "no_bid": {
                "type": "string",
                "description": "No side best bid price",
                "example": "0.42"
              },
              "no_ask": {
                "type": "string",
                "description": "No side best ask price",
                "example": "0.44"
              },
              "price": {
                "type": "string",
                "description": "Price for the last traded YES contract on this market in dollars.",
                "example": "0.57"
              },
              "volume": {
                "type": "string",
                "description": "String representation of the market volume in contracts.",
                "example": "12500.0"
              },
              "open_interest": {
                "type": "string",
                "description": "String representation of the number of contracts bought on this market disregarding netting.",
                "example": "8500.0"
              },
              "last_trade_time": {
                "type": "string",
                "description": "Timestamp of the most recent trade. Format: ISO 8601 with timezone offset.",
                "example": "2026-05-27T09:03:50.000+0000"
              }
            },
            "description": "EventContractSnapshotVo",
            "title": "EventContractSnapshotVo"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNAUTHORIZED",
            "message": "Insufficient permission"
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
                "description": "Internal logic error code"
              },
              "message": {
                "type": "string",
                "description": "Error message"
              }
            },
            "description": "ErrorResponseVo",
            "title": "ErrorResponseVo"
          },
          "example": {
            "error_code": "UNSUPPORTED_CATEGORY",
            "message": "Unsupported category:US_EVENTS"
          }
        }
      }
    }
  },
  "postman": {
    "name": "List Event Market Snapshots",
    "description": {
      "content": "Retrieves the latest market snapshot for a single event contract market, including yes/no bid-ask prices, last trade price, volume, open interest, and market status.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "market-data",
        "event-contracts",
        "markets",
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
            "content": "(Required) Market symbol. A market is a single binary contract within an event.",
            "type": "text/plain"
          },
          "key": "symbol",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Category, default is US_EVENT, currently only US_EVENT is supported.",
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


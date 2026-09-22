# Event Contracts — Verbatim Reference

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


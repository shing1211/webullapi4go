# Authentication — Verbatim Reference

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> Webull uses a dual layer: an HMAC-SHA1 request signature plus an access token. Server-to-server endpoints sign every request; Display Solution (Client-to-Server) uses an OAuth-style client token.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Create Token

> Source: <https://developer.webull.hk/apis/docs/reference/create-token.md>

### Create Token

Creates an access token. This interface is used to generate a new Token.

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
  "path": "/auth/tokens/create",
  "method": "post",
  "tags": [
    "Server-To-Server"
  ],
  "description": "Creates an access token. This interface is used to generate a new Token.",
  "operationId": "createToken",
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
              "expires_at",
              "status",
              "token"
            ],
            "type": "object",
            "properties": {
              "token": {
                "type": "string",
                "description": "Access token string, used for authentication of subsequent API calls. Token is a 32-digit hexadecimal string that is unique and time-sensitive.",
                "example": "ccb071f764864b65a1fb48484e940a56"
              },
              "expires_at": {
                "type": "integer",
                "description": "Token expiration timestamp, a Unix timestamp in milliseconds. After this time, the token will become invalid and need to be recreated.",
                "format": "int64",
                "example": 1755486723000
              },
              "status": {
                "type": "string",
                "description": "Token validity status code, indicating the result of the token operation. PENDING indicates pending verification, NORMAL indicates valid, INVALID indicates the token is invalid, and EXPIRED indicates it has expired.",
                "example": "PENDING",
                "enum": [
                  "PENDING",
                  "NORMAL",
                  "INVALID",
                  "EXPIRED"
                ]
              }
            },
            "description": "Token response information",
            "title": "TokenRespVo"
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
    "name": "Create Token",
    "description": {
      "content": "Creates an access token. This interface is used to generate a new Token.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "auth",
        "tokens",
        "create"
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
    "method": "POST"
  }
}
```

## Check Token

> Source: <https://developer.webull.hk/apis/docs/reference/check-token.md>

### Check Token

Retrieves Token status.

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
  "path": "/auth/tokens/check",
  "method": "post",
  "tags": [
    "Server-To-Server"
  ],
  "description": "Retrieves Token status.",
  "operationId": "checkToken",
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
            "token"
          ],
          "type": "object",
          "properties": {
            "token": {
              "type": "string",
              "description": "Access token, used for identity authentication and permission verification. This field is a unique identifier for token checking, refreshing, etc.",
              "example": "ccb071f764864b65a1fb48484e940a56"
            }
          },
          "description": "Token request parameters",
          "title": "TokenReq"
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
              "expires_at",
              "status",
              "token"
            ],
            "type": "object",
            "properties": {
              "token": {
                "type": "string",
                "description": "Access token string, used for authentication of subsequent API calls. Token is a 32-digit hexadecimal string that is unique and time-sensitive.",
                "example": "ccb071f764864b65a1fb48484e940a56"
              },
              "expires_at": {
                "type": "integer",
                "description": "Token expiration timestamp, a Unix timestamp in milliseconds. After this time, the token will become invalid and need to be recreated.",
                "format": "int64",
                "example": 1755486723000
              },
              "status": {
                "type": "string",
                "description": "Token validity status code, indicating the result of the token operation. PENDING indicates pending verification, NORMAL indicates valid, INVALID indicates the token is invalid, and EXPIRED indicates it has expired.",
                "example": "PENDING",
                "enum": [
                  "PENDING",
                  "NORMAL",
                  "INVALID",
                  "EXPIRED"
                ]
              }
            },
            "description": "Token response information",
            "title": "TokenRespVo"
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
    "token": "ccb071f764864b65a1fb48484e940a56"
  },
  "postman": {
    "name": "Check Token",
    "description": {
      "content": "Retrieves Token status.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "auth",
        "tokens",
        "check"
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

## Create Client Token (Display)

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/create-client-token.md>

### Create Client Token

Generates an access token and refresh token pair for end-user applications to directly access Webull services. This endpoint must be called from your institution's backend service with valid AK/SK credentials. The access token authorizes API calls, while the refresh token enables token renewal. Creating a new token pair for an existing customer immediately invalidates any previously issued tokens for that customer. <br/>• Access token expires in approximately 2 hours (may vary slightly)<br/>• Refresh token expires in 15 days<br/>• **Note**: This access_token is used to access Market Data API

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
  "path": "/auth/client-tokens/create",
  "method": "post",
  "tags": [
    "Client-To-Server"
  ],
  "description": "Generates an access token and refresh token pair for end-user applications to directly access Webull services. This endpoint must be called from your institution's backend service with valid AK/SK credentials. The access token authorizes API calls, while the refresh token enables token renewal. Creating a new token pair for an existing customer immediately invalidates any previously issued tokens for that customer. <br/>• Access token expires in approximately 2 hours (may vary slightly)<br/>• Refresh token expires in 15 days<br/>• **Note**: This access_token is used to access Market Data API",
  "operationId": "createClientToken",
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
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "client_user_id"
          ],
          "type": "object",
          "properties": {
            "client_user_id": {
              "type": "string",
              "description": "The unique identifier for the customer in your system. This identifier is used to associate the token pair with the specific end-user.",
              "example": "a9c3a8cbd6184a1292a3e36eff1e6f3f"
            }
          },
          "description": "Token Create Request",
          "title": "TokenCreateReq"
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
              "access_token",
              "expires_at",
              "refresh_expires_at",
              "refresh_token"
            ],
            "type": "object",
            "properties": {
              "access_token": {
                "type": "string",
                "description": "Short-lived token for API authorization.",
                "example": "TH.19b68a108ad-51627ef292044cc3b39c37800523d803"
              },
              "expires_at": {
                "type": "integer",
                "description": "Access Token Expiration Time in milliseconds since epoch.",
                "format": "int64",
                "example": 1766987799637
              },
              "refresh_token": {
                "type": "string",
                "description": "Long-lived token for obtaining new access tokens.",
                "example": "19b68a108ad-4013b3daec584bc2b078ae1902ce0986"
              },
              "refresh_expires_at": {
                "type": "integer",
                "description": "Refresh Token Expiration Time in milliseconds since epoch.",
                "format": "int64",
                "example": 1767004942637
              }
            },
            "description": "Token Create Result",
            "title": "TokenCreateResult"
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": {
    "client_user_id": "a9c3a8cbd6184a1292a3e36eff1e6f3f"
  },
  "postman": {
    "name": "Create Client Token",
    "description": {
      "content": "Generates an access token and refresh token pair for end-user applications to directly access Webull services. This endpoint must be called from your institution's backend service with valid AK/SK credentials. The access token authorizes API calls, while the refresh token enables token renewal. Creating a new token pair for an existing customer immediately invalidates any previously issued tokens for that customer. <br/>• Access token expires in approximately 2 hours (may vary slightly)<br/>• Refresh token expires in 15 days<br/>• **Note**: This access_token is used to access Market Data API",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "auth",
        "client-tokens",
        "create"
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

## Refresh Client Token (Display)

> Source: <https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/refresh-client-token.md>

### Refresh Client Token

Obtains a new access token and refresh token pair using a valid refresh token. This endpoint requires AK/SK signature authentication from your institution's backend service. Both the refresh token and request signature are validated before issuing new credentials. The old tokens will be invalidated after successful refresh. <br/>• Access token expires in approximately 2 hours (may vary slightly)<br/>• Refresh token expires in 15 days<br/>• **Note**: This access_token is used to access Market Data API

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
  "path": "/auth/client-tokens/refresh",
  "method": "post",
  "tags": [
    "Client-To-Server"
  ],
  "description": "Obtains a new access token and refresh token pair using a valid refresh token. This endpoint requires AK/SK signature authentication from your institution's backend service. Both the refresh token and request signature are validated before issuing new credentials. The old tokens will be invalidated after successful refresh. <br/>• Access token expires in approximately 2 hours (may vary slightly)<br/>• Refresh token expires in 15 days<br/>• **Note**: This access_token is used to access Market Data API",
  "operationId": "refreshClientToken",
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
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "refresh_token"
          ],
          "type": "object",
          "properties": {
            "refresh_token": {
              "type": "string",
              "description": "The refresh token obtained from the token creation or previous refresh operation.",
              "example": "19b68a108ad-4013b3daec584bc2b078ae1902ce0986"
            }
          },
          "description": "Token Refresh Request",
          "title": "TokenRefreshReq"
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
              "access_token",
              "expires_at",
              "refresh_expires_at",
              "refresh_token"
            ],
            "type": "object",
            "properties": {
              "access_token": {
                "type": "string",
                "description": "Short-lived token for API authorization.",
                "example": "US.19b68a108ad-51627ef292044cc3b39c37800523d803"
              },
              "expires_at": {
                "type": "integer",
                "description": "Access Token Expiration Time in milliseconds since epoch.",
                "format": "int64",
                "example": 1766987799637
              },
              "refresh_token": {
                "type": "string",
                "description": "Long-lived token for obtaining new access tokens.",
                "example": "19b68a108ad-4013b3daec584bc2b078ae1902ce0986"
              },
              "refresh_expires_at": {
                "type": "integer",
                "description": "Refresh Token Expiration Time in milliseconds since epoch.",
                "format": "int64",
                "example": 1767004942637
              }
            },
            "description": "Token Refresh Result",
            "title": "TokenRefreshResult"
          }
        }
      }
    }
  },
  "jsonRequestBodyExample": {
    "refresh_token": "19b68a108ad-4013b3daec584bc2b078ae1902ce0986"
  },
  "postman": {
    "name": "Refresh Client Token",
    "description": {
      "content": "Obtains a new access token and refresh token pair using a valid refresh token. This endpoint requires AK/SK signature authentication from your institution's backend service. Both the refresh token and request signature are validated before issuing new credentials. The old tokens will be invalidated after successful refresh. <br/>• Access token expires in approximately 2 hours (may vary slightly)<br/>• Refresh token expires in 15 days<br/>• **Note**: This access_token is used to access Market Data API",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "auth",
        "client-tokens",
        "refresh"
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


# Broker API — FD (US) — Verbatim Reference

> US Broker FD surface. Documented on the US site only and marked **provisional** in this SDK: all paths require live-probe confirmation and the HK sandbox returns 404.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## List Accounts

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/list-accounts.md>

### List Accounts

Retrieves a paginated list of accounts associated with the requesting institution.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/accounts/list",
  "method": "get",
  "tags": [
    "Accounts"
  ],
  "description": "Retrieves a paginated list of accounts associated with the requesting institution.",
  "operationId": "listAccounts",
  "parameters": [
    {
      "name": "pagination_key",
      "in": "query",
      "description": "Pagination key from previous response for next page.",
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
                    "account_number",
                    "account_type"
                  ],
                  "type": "object",
                  "properties": {
                    "account_id": {
                      "type": "string",
                      "description": "Unique identifier of the account to which restrictions are currently applied.",
                      "example": "f47ac10b58cc4372a5670e02b2c3d479"
                    },
                    "account_number": {
                      "type": "string",
                      "description": "Account number associated with the account.",
                      "example": "U1234567"
                    },
                    "account_type": {
                      "type": "string",
                      "description": "Account Type<br/>CASH - Cash Account Type<br/>MARGIN - Margin Account Type<br/>",
                      "enum": [
                        "CASH",
                        "MARGIN"
                      ]
                    },
                    "rep_code": {
                      "type": "string",
                      "description": "Rep Code",
                      "example": "REP001"
                    }
                  },
                  "description": "Account List Result",
                  "title": "AccountListResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoAccountListResult"
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
      "content": "Retrieves a paginated list of accounts associated with the requesting institution.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
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

## Account Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/get-account-detail.md>

### Get Account Detail

Retrieves the detailed information of a specific account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/accounts/get",
  "method": "get",
  "tags": [
    "Accounts"
  ],
  "description": "Retrieves the detailed information of a specific account.",
  "operationId": "getAccountDetail",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Unique identifier of the account to query restrictions for.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "f47ac10b58cc4372a5670e02b2c3d479"
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
              "account_id",
              "account_number",
              "account_type"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Unique identifier of the account to which restrictions are currently applied.",
                "example": "f47ac10b58cc4372a5670e02b2c3d479"
              },
              "account_number": {
                "type": "string",
                "description": "Account number associated with the account.",
                "example": "U1234567"
              },
              "account_type": {
                "type": "string",
                "description": "Account Type<br/>CASH - Cash Account Type<br/>MARGIN - Margin Account Type<br/>",
                "enum": [
                  "CASH",
                  "MARGIN"
                ]
              },
              "rep_code": {
                "type": "string",
                "description": "Rep Code",
                "example": "REP001"
              },
              "commission_code": {
                "type": "string",
                "description": "Commission Code",
                "example": "COM001"
              },
              "restrict_infos": {
                "type": "array",
                "description": "List of restrictions currently effective for this account. ",
                "items": {
                  "required": [
                    "restriction_code",
                    "restriction_start_time"
                  ],
                  "type": "object",
                  "properties": {
                    "restriction_code": {
                      "type": "string",
                      "description": "A code representing the type of restriction applied. <br/>RISK: Risk Restriction<br/>TRADING: Trading Restriction<br/>TRADING_STOCK_ETF: Stock and ETF Trading Restriction<br/>TRADING_FRACTIONAL: Fractional Share Trading Restriction<br/>TRADING_OPTIONS: Options Trading Restriction<br/>TRADING_EVENT: Event Contract Trading Restriction<br/>TRADING_OTC: OTC Trading Restriction<br/>TRADING_OVERNIGHT: Overnight Trading Restriction<br/>ACCOUNT_LINKING: Account Linking Restriction<br/>FUNDING_ACH: ACH Deposit Restriction<br/>FUNDING_DEBIT_CARD: Debit Card Deposit Restriction<br/>FUNDING_DIRECT_DEPOSIT: Direct Deposit Restriction<br/>WITHDRAWAL: Withdrawal Restriction<br/>TRANSFER_ACATS_IN: ACATS Incoming Transfer Restriction<br/>TRANSFER_ACATS_OUT: ACATS Outgoing Transfer Restriction",
                      "example": "RISK"
                    },
                    "restriction_start_time": {
                      "type": "string",
                      "description": "The date and time when the restriction takes effect. Time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",
                      "example": "2025-01-06T22:59:59.012Z"
                    },
                    "reason": {
                      "type": "string",
                      "description": "A description of the reason for the restriction.",
                      "example": "Unusual trading activity detected on 2025-01-05."
                    }
                  },
                  "description": "RestrictInfo",
                  "title": "RestrictInfo"
                }
              }
            },
            "description": "Account Detail Response",
            "title": "AccountDetailResult"
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
    "name": "Get Account Detail",
    "description": {
      "content": "Retrieves the detailed information of a specific account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique identifier of the account to query restrictions for.",
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

## Create Account

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/create-account-apply.md>

### Create Account

Submits a request to create an account for your end user. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Application Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/accounts/create",
  "method": "post",
  "tags": [
    "Accounts"
  ],
  "description": "Submits a request to create an account for your end user. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Application Events](../fd-events/application-events)",
  "operationId": "createAccountApply",
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
            "client_request_id",
            "forms"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client-supplied request ID, which is unique for each account submission request.",
              "example": "123e4567e89b12d3a456426614174000"
            },
            "application_id": {
              "type": "string",
              "description": "System-generated unique identifier assigned to the account application record. This field is omitted for the initial submission and must be provided for subsequent submissions related to the same account application.",
              "example": "4F5A13B5FDFD43D29C7A3A856B2458A2"
            },
            "related_account_id": {
              "type": "string",
              "description": "The unique identifier of the related account. This field is required when submitting an additional account application related to an existing account. For example, if you have an Event Contract account and want to add a Brokerage Margin account, you would provide the account_id of the Event Contract account in this field. If this field is not provided or is empty, it indicates that this submission is for an initial account application.",
              "example": "SVHC5L98E0D79UR0AB4QJ962JB"
            },
            "forms": {
              "type": "array",
              "description": "Array of one or more FormData objects representing the account form being submitted. Use form details API to retrieve the form schema and details.",
              "items": {
                "required": [
                  "form_info",
                  "json_data"
                ],
                "type": "object",
                "properties": {
                  "form_info": {
                    "required": [
                      "form_code",
                      "version"
                    ],
                    "type": "object",
                    "properties": {
                      "form_code": {
                        "type": "string",
                        "description": "Unique identifier of the form, typically used to distinguish different form definitions.",
                        "example": "NEW_ACCOUNT_BASIC_FORM"
                      },
                      "version": {
                        "type": "string",
                        "description": "Version number of the form definition, used for compatibility and evolutionary control.",
                        "example": "1.0"
                      }
                    },
                    "description": "Form Identifier.",
                    "title": "FormId"
                  },
                  "json_data": {
                    "type": "object",
                    "description": "Object containing the completed form data as defined by the form schema and will vary from form to form. ",
                    "example": {}
                  }
                },
                "description": "Account Form",
                "title": "AccountForm"
              }
            }
          },
          "description": "Account Submit Request",
          "title": "AccountSubmitRequest"
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
              "application_id",
              "client_request_id"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Unique client-generated identifier used for request tracking, idempotency and audit purposes.",
                "example": "4F5A13B5FDFD43D29C7A3A856B2458A1"
              },
              "application_id": {
                "type": "string",
                "description": "System-generated unique identifier assigned to the account application record.",
                "example": "4F5A13B5FDFD43D29C7A3A856B2458A2"
              }
            },
            "description": "Account Submit Response",
            "title": "AccountSubmitResult"
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
    "client_request_id": "123e4567e89b12d3a456426614174000",
    "application_id": "4F5A13B5FDFD43D29C7A3A856B2458A2",
    "related_account_id": "SVHC5L98E0D79UR0AB4QJ962JB",
    "forms": [
      {
        "form_info": {
          "form_code": "NEW_ACCOUNT_BASIC_FORM",
          "version": "1.0"
        },
        "json_data": {}
      }
    ]
  },
  "postman": {
    "name": "Create Account",
    "description": {
      "content": "Submits a request to create an account for your end user. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Application Events](../fd-events/application-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
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

## Update Account

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/update-account-apply.md>

### Update Account

Submit a request to update an existing account. This API only updates the fields submitted in the forms data. Any attributes not included in the request will remain unchanged. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Account Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/accounts/update",
  "method": "post",
  "tags": [
    "Accounts"
  ],
  "description": "Submit a request to update an existing account. This API only updates the fields submitted in the forms data. Any attributes not included in the request will remain unchanged. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Account Events](../fd-events/account-events)",
  "operationId": "updateAccountApply",
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
    "description": "Request payload containing the account details to update.",
    "content": {
      "application/json": {
        "schema": {
          "required": [
            "account_id",
            "client_request_id",
            "forms"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Unique identifier for the request, generated each time.",
              "example": "123e4567e89b12d3a456426614174000"
            },
            "account_id": {
              "type": "string",
              "description": "The account ID for which the update is being submitted.",
              "example": "943a9802f6c14983b3b4755c69c01717"
            },
            "forms": {
              "type": "array",
              "description": "Array of one or more FormData objects representing the account form being submitted. Use form details API to retrieve the form schema and details.",
              "items": {
                "required": [
                  "form_info",
                  "json_data"
                ],
                "type": "object",
                "properties": {
                  "form_info": {
                    "required": [
                      "form_code",
                      "version"
                    ],
                    "type": "object",
                    "properties": {
                      "form_code": {
                        "type": "string",
                        "description": "Unique identifier of the form, typically used to distinguish different form definitions.",
                        "example": "UPDATE_ACCOUNT_FORM"
                      },
                      "version": {
                        "type": "string",
                        "description": "Version number of the form definition, used for compatibility and evolutionary control.",
                        "example": "1.0"
                      }
                    },
                    "description": "Form Identifier.",
                    "title": "UpdateFormId"
                  },
                  "json_data": {
                    "type": "object",
                    "description": "Object containing the completed form data as defined by the form schema and will vary from form to form. ",
                    "example": {
                      "type": "object",
                      "description": "This form is used to submit post-onboarding account update requests for an existing account. It serves as a unified entry point for multiple account maintenance and compliance-related update scenarios after the account has been successfully opened. The specific update scenario is determined by the `type` field. Based on the selected type, the corresponding update information object must be provided. This form is commonly used for periodic reviews, regulatory compliance verification, risk control follow-ups, and other account data maintenance processes to ensure that account records remain accurate, valid, and up to date.",
                      "form_code": "ACCOUNT_UPDATE_FORM",
                      "properties": {
                        "type": {
                          "type": "string",
                          "description": "Specifies the category of account update to be performed after account opening. This field determines which account update scenario applies and which corresponding information object is required in the request. Each type represents a distinct account maintenance or compliance update operation."
                        },
                        "commission_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is COMMISSION_UPDATE.",
                          "properties": {
                            "commission_code": {
                              "type": "string",
                              "description": "Commission Code for account, determining the commission structure."
                            }
                          },
                          "required": [
                            "commission_code"
                          ]
                        },
                        "identity_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is IDENTITY_UPDATE. This form is used to collect updated identification information for an existing account. It is required during account maintenance, periodic review, or compliance verification to ensure the accuracy and validity of identity records.",
                          "properties": {
                            "id_doc_type": {
                              "type": "string",
                              "description": "Type of identification document provided by the applicant. The value must be retrieved from the Master Data API where data_type = ID_DOC_TYPE."
                            },
                            "id_number": {
                              "type": "string",
                              "description": "Identification document number exactly as shown on the provided ID."
                            },
                            "id_expiry_date": {
                              "type": "string",
                              "description": "Expiration date of the identification document.",
                              "format": "yyyy-MM-dd"
                            },
                            "picture_front_doc_id": {
                              "type": "string",
                              "description": "Document ID referencing the uploaded image of the front side of the identification document."
                            },
                            "picture_back_doc_id": {
                              "type": "string",
                              "description": "Document ID referencing the uploaded image of the back side of the identification document."
                            }
                          },
                          "required": [
                            "id_doc_type",
                            "id_number",
                            "id_expiry_date",
                            "picture_front_doc_id",
                            "picture_back_doc_id"
                          ]
                        },
                        "employment_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is EMPLOYMENT_UPDATE. Collects updated employment information required during account maintenance or compliance review.",
                          "properties": {
                            "employment_type": {
                              "type": "string",
                              "description": "Type of employment. Value must be obtained from the Master Data API with data_type = EMPLOYMENT_TYPE."
                            },
                            "company_name": {
                              "type": "string",
                              "description": "Name of the applicant current employer.",
                              "max_length": 128
                            },
                            "occupation": {
                              "type": "string",
                              "description": "Specific occupation under the employment type. Value must be obtained from the Master Data API with data_type = OCCUPATION."
                            },
                            "position": {
                              "type": "string",
                              "description": "Job position or title. Value must be obtained from the Master Data API with data_type = POSITION."
                            },
                            "years_employed": {
                              "type": "number",
                              "description": "Total duration (in years) the applicant has been employed with the current employer. Decimal values are allowed to represent partial years (e.g. 1.5)."
                            },
                            "employment_address": {
                              "type": "object",
                              "description": "Address of the applicant current employer.",
                              "properties": {
                                "country": {
                                  "type": "string",
                                  "description": "Employer address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                                },
                                "state": {
                                  "type": "string",
                                  "description": "Employer address state. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                                },
                                "city": {
                                  "type": "string",
                                  "description": "Employer address city.",
                                  "max_length": 256
                                },
                                "street_address": {
                                  "type": "string",
                                  "description": "Street address of the employer. Multiple lines are allowed.",
                                  "max_length": 512
                                },
                                "postal_code": {
                                  "type": "string",
                                  "description": "Postal code of the employer address.",
                                  "max_length": 80
                                },
                                "neighborhood_code": {
                                  "type": "string",
                                  "description": "Neighborhood or district code of the employer address, if applicable.",
                                  "max_length": 128
                                },
                                "apartment_number": {
                                  "type": "string",
                                  "description": "Apartment number of the employer address, if applicable.",
                                  "max_length": 128
                                },
                                "building_number": {
                                  "type": "string",
                                  "description": "Building number of the employer address, if applicable.",
                                  "max_length": 128
                                }
                              },
                              "required": [
                                "country",
                                "state",
                                "city",
                                "street_address",
                                "postal_code",
                                "neighborhood_code",
                                "apartment_number",
                                "building_number"
                              ]
                            }
                          },
                          "required": [
                            "employment_type",
                            "company_name",
                            "occupation",
                            "position",
                            "years_employed",
                            "employment_address"
                          ]
                        },
                        "trusted_contact_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is TRUSTED_UPDATE. Collects or updates trusted contact information for an account, used for account protection, compliance review, or situations where the account holder cannot be contacted.",
                          "properties": {
                            "first_name": {
                              "type": "string",
                              "description": "First name of the trusted contact person."
                            },
                            "middle_name": {
                              "type": "string",
                              "description": "Middle name of the trusted contact person, if applicable."
                            },
                            "last_name": {
                              "type": "string",
                              "description": "Last name of the trusted contact person."
                            },
                            "phone_country": {
                              "type": "string",
                              "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                            },
                            "phone_number": {
                              "type": "string",
                              "description": "Primary contact phone number of the trusted contact person.",
                              "example": "+1-2025550124",
                              "max_length": 20
                            },
                            "relationship": {
                              "type": "string",
                              "description": "Relationship between the account holder and the trusted contact person.",
                              "max_length": 32
                            },
                            "trusted_address": {
                              "type": "object",
                              "description": "Residential or mailing address of the trusted contact person.",
                              "properties": {
                                "country": {
                                  "type": "string",
                                  "description": "Country of the trusted contact person address. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                                },
                                "state": {
                                  "type": "string",
                                  "description": "State or province of the trusted contact person address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                                },
                                "city": {
                                  "type": "string",
                                  "description": "City of the trusted contact person address.",
                                  "max_length": 256
                                },
                                "street_address": {
                                  "type": "string",
                                  "description": "Street address of the trusted contact person. Multiple address lines are allowed.",
                                  "max_length": 512
                                },
                                "postal_code": {
                                  "type": "string",
                                  "description": "Postal or ZIP code of the trusted contact person address.",
                                  "max_length": 80
                                },
                                "neighborhood_code": {
                                  "type": "string",
                                  "description": "Neighborhood or district code of the trusted contact person address, if applicable.",
                                  "max_length": 128
                                },
                                "apartment_number": {
                                  "type": "string",
                                  "description": "Apartment or unit number of the trusted contact person address, if applicable.",
                                  "max_length": 128
                                },
                                "building_number": {
                                  "type": "string",
                                  "description": "Building number of the trusted contact person address, if applicable.",
                                  "max_length": 128
                                }
                              },
                              "required": [
                                "country",
                                "state",
                                "city",
                                "street_address",
                                "postal_code",
                                "neighborhood_code"
                              ]
                            }
                          },
                          "required": [
                            "first_name",
                            "last_name",
                            "phone_country",
                            "phone_number",
                            "relationship",
                            "trusted_address"
                          ]
                        },
                        "email_address_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is EMAIL_ADDRESS_UPDATE. Collects or updates the primary email address associated with an account. This information is used for account communication, security notifications, and compliance-related correspondence.",
                          "properties": {
                            "email_address": {
                              "type": "string",
                              "description": "Email address",
                              "max_length": 256
                            }
                          },
                          "required": [
                            "email_address"
                          ]
                        },
                        "address_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is ADDRESS_UPDATE. Collects or updates address information associated with an account, including residential and mailing addresses. This information is used for account administration, regulatory compliance, and official correspondence.",
                          "properties": {
                            "id_address": {
                              "type": "object",
                              "description": "Residential address of the account holder as recorded for identity verification and regulatory purposes.",
                              "properties": {
                                "country": {
                                  "type": "string",
                                  "description": "Country of residence. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                                },
                                "state": {
                                  "type": "string",
                                  "description": "State or province of residence. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                                },
                                "city": {
                                  "type": "string",
                                  "description": "City of residence.",
                                  "max_length": 256
                                },
                                "street_address": {
                                  "type": "string",
                                  "description": "Street address of the account holder residence. Multiple address lines are allowed.",
                                  "max_length": 512
                                },
                                "postal_code": {
                                  "type": "string",
                                  "description": "Postal or ZIP code of the residential address.",
                                  "max_length": 80
                                },
                                "neighborhood_code": {
                                  "type": "string",
                                  "description": "Neighborhood or district code of the id address, if applicable.",
                                  "max_length": 128
                                },
                                "apartment_number": {
                                  "type": "string",
                                  "description": "Apartment number of the id address, if applicable.",
                                  "max_length": 128
                                },
                                "building_number": {
                                  "type": "string",
                                  "description": "Building number of the id address, if applicable.",
                                  "max_length": 128
                                }
                              },
                              "required": [
                                "country",
                                "state",
                                "city",
                                "street_address",
                                "postal_code",
                                "neighborhood_code",
                                "apartment_number",
                                "building_number"
                              ]
                            },
                            "mail_address": {
                              "type": "object",
                              "description": "Mailing address used for receiving account-related correspondence, if different from the residential address.",
                              "properties": {
                                "country": {
                                  "type": "string",
                                  "description": "Mailing address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                                },
                                "state": {
                                  "type": "string",
                                  "description": "State or province of the mailing address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                                },
                                "city": {
                                  "type": "string",
                                  "description": "City of the mailing address.",
                                  "max_length": 256
                                },
                                "street_address": {
                                  "type": "string",
                                  "description": "Street address of the mailing address. Multiple address lines are allowed.",
                                  "max_length": 512
                                },
                                "postal_code": {
                                  "type": "string",
                                  "description": "Postal or ZIP code of the mailing address.",
                                  "max_length": 80
                                },
                                "neighborhood_code": {
                                  "type": "string",
                                  "description": "Neighborhood or district code of the mailing address.",
                                  "max_length": 128
                                },
                                "apartment_number": {
                                  "type": "string",
                                  "description": "Apartment number of the mailing address.",
                                  "max_length": 128
                                },
                                "building_number": {
                                  "type": "string",
                                  "description": "Building number of the mailing address.",
                                  "max_length": 128
                                }
                              },
                              "required": [
                                "country",
                                "state",
                                "city",
                                "street_address",
                                "postal_code",
                                "neighborhood_code",
                                "apartment_number",
                                "building_number"
                              ]
                            }
                          },
                          "required": [
                            "id_address",
                            "mail_address"
                          ]
                        },
                        "phone_number_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is PHONE_NUMBER_UPDATE. Collects or updates the primary phone number associated with an account, used for account communication, security verification, and regulatory notifications.",
                          "properties": {
                            "phone_country": {
                              "type": "string",
                              "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                            },
                            "phone_number": {
                              "type": "string",
                              "description": "Primary contact phone number associated with the account, in international format (E.164). Used for account notifications, verification, and security-related communications.",
                              "example": "+1-2025550123",
                              "max_length": 20
                            }
                          },
                          "required": [
                            "phone_country",
                            "phone_number"
                          ]
                        },
                        "politically_exposed_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is PEP_UPDATE. Collects or updates Politically Exposed Person (PEP) information for an account. This information is used for compliance review, enhanced due diligence (EDD), and regulatory reporting.",
                          "properties": {
                            "is_politically_exposed": {
                              "type": "boolean",
                              "description": "Indicates whether the account holder or any close associate or immediate family member is a Politically Exposed Person (PEP).  If true, the politically_exposed_info object must be provided with detailed disclosure information.  Please note: PEP accounts are currently not supported. Any request submitted with is_politically_exposed = true will be rejected and a failure response will be returned."
                            },
                            "politically_exposed_info": {
                              "type": "object",
                              "description": "Detail information about Politically Exposed Person (PEP) associated with the account holder.",
                              "properties": {
                                "pep_organization": {
                                  "type": "string",
                                  "description": "Name of the related political organization.",
                                  "max_length": 256
                                },
                                "immediate_family": {
                                  "type": "array",
                                  "description": "Name of family members, including former spouses.",
                                  "items": {
                                    "type": "string",
                                    "description": ""
                                  }
                                }
                              },
                              "required": [
                                "pep_organization",
                                "immediate_family"
                              ]
                            }
                          },
                          "required": [
                            "is_politically_exposed"
                          ]
                        },
                        "company_control_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is COMPANY_CONTROL_UPDATE. Collects or updates information about the applicant status as a control person of a public company. This information is used for compliance review, enhanced due diligence (EDD), and regulatory reporting.",
                          "properties": {
                            "is_company_control_person": {
                              "type": "boolean",
                              "description": "Indicates whether the applicant is a control person of a public company. If true, `public_company_infos` must be provided with details of the companies controlled."
                            },
                            "public_company_infos": {
                              "type": "array",
                              "description": "List of public companies controlled by the applicant. Each item must include the company stock symbol.",
                              "items": {
                                "type": "object",
                                "description": "",
                                "properties": {
                                  "symbol": {
                                    "type": "string",
                                    "description": "Stock symbol of the public company controlled by the applicant, as listed on the exchange."
                                  }
                                },
                                "required": [
                                  "symbol"
                                ]
                              }
                            }
                          },
                          "required": [
                            "is_company_control_person"
                          ]
                        },
                        "disclosure_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is DISCLOSURE_UPDATE. Collects or updates the applicant regulatory and compliance-related affiliations or disclosures. This form is used for compliance review, enhanced due diligence (EDD), and regulatory reporting, and may include multiple types of information beyond exchange or FINRA affiliations.",
                          "properties": {
                            "brokerage_disclosure_info": {
                              "type": "object",
                              "description": "Brokerage disclosure info",
                              "properties": {
                                "is_exchange_or_finra_affiliated": {
                                  "type": "boolean",
                                  "description": "Indicates whether the applicant is employed by, or directly associated with, an exchange or a FINRA member firm. If true, exchange_affiliations details must be provided."
                                },
                                "exchange_affiliations": {
                                  "type": "array",
                                  "description": "List of the applicant affiliations with exchanges or FINRA member firms. Required if is_exchange_or_finra_affiliated is true.",
                                  "items": {
                                    "type": "object",
                                    "description": "",
                                    "properties": {
                                      "company_name": {
                                        "type": "string",
                                        "description": "Name of the exchange or FINRA member firm the applicant is affiliated with.",
                                        "max_length": 128
                                      },
                                      "country": {
                                        "type": "string",
                                        "description": "Country of the affiliated exchange or FINRA firm. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                                      },
                                      "state": {
                                        "type": "string",
                                        "description": "State or province of the affiliated exchange or FINRA firm. Value must be obtained from the Master Data API with data_type = STATE."
                                      },
                                      "city": {
                                        "type": "string",
                                        "description": "City of the affiliated exchange or FINRA firm.",
                                        "max_length": 256
                                      },
                                      "street_address": {
                                        "type": "string",
                                        "description": "Street address of the affiliated exchange or FINRA firm. Multiple lines are allowed.",
                                        "max_length": 512
                                      },
                                      "postal_code": {
                                        "type": "string",
                                        "description": "Postal or ZIP code of the affiliated exchange or FINRA firm.",
                                        "max_length": 80
                                      }
                                    },
                                    "required": [
                                      "company_name",
                                      "country",
                                      "state",
                                      "city",
                                      "street_address",
                                      "postal_code"
                                    ]
                                  }
                                }
                              },
                              "required": [
                                "is_exchange_or_finra_affiliated"
                              ]
                            }
                          },
                          "required": [
                            "brokerage_disclosure_info"
                          ]
                        },
                        "dependent_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is DEPENDENT_UPDATE. Collects or updates the number of dependents associated with an account holder. This information is used for account maintenance, compliance review, and regulatory reporting purposes.",
                          "properties": {
                            "num_dependents": {
                              "type": "number",
                              "description": "Number of dependents financially or legally associated with the account holder."
                            }
                          },
                          "required": [
                            "num_dependents"
                          ]
                        },
                        "citizenship_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is CITIZENSHIP_UPDATE. Collects or updates the country of citizenship for an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
                          "properties": {
                            "country_of_citizenship": {
                              "type": "string",
                              "description": "Country of citizenship of the account holder. The value must be retrieved from the Master Data API with data_type = COUNTRY."
                            }
                          },
                          "required": [
                            "country_of_citizenship"
                          ]
                        },
                        "gender_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is GENDER_UPDATE. Collects or updates the gender of an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
                          "properties": {
                            "gender": {
                              "type": "string",
                              "description": "Gender of the account holder. Allowed values: F (Female), M (Male)."
                            }
                          },
                          "required": [
                            "gender"
                          ]
                        },
                        "marital_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is MARITAL_UPDATE. Collects or updates the marital status of an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
                          "properties": {
                            "marital_status": {
                              "type": "string",
                              "description": "Marital status of the account holder. Value must be obtained from the Master Data API with data_type = MARITAL_STATUS."
                            }
                          },
                          "required": [
                            "marital_status"
                          ]
                        },
                        "financial_investment_info": {
                          "type": "object",
                          "description": "This object is applicable when the update type is FINANCIAL_INVESTMENT_UPDATE. Collects or updates the investment profile and financial information of an account holder. This information is used for account maintenance, risk assessment, and regulatory compliance.",
                          "properties": {
                            "financial_info": {
                              "type": "object",
                              "description": "Financial information of the account holder.",
                              "properties": {
                                "liquidity_needs": {
                                  "type": "string",
                                  "description": "Indicates the client liquidity needs or preference for cash accessibility. Value must be obtained from the Master Data API with data_type = LIQUIDITY_NEEDS."
                                },
                                "annual_income": {
                                  "type": "string",
                                  "description": "Represents the client annual income. Value must be obtained from the Master Data API with data_type = ANNUAL_INCOME."
                                },
                                "total_net_worth": {
                                  "type": "string",
                                  "description": "Represents the client total net worth in USD. Value must be obtained from the Master Data API with data_type = TOTAL_NET_WORTH."
                                },
                                "liquid_net_worth": {
                                  "type": "string",
                                  "description": "Represents the client liquid net worth (USD), including cash and easily marketable assets. Value must be obtained from the Master Data API with data_type = LIQUID_NET_WORTH."
                                }
                              },
                              "required": [
                                "liquidity_needs",
                                "annual_income",
                                "total_net_worth",
                                "liquid_net_worth"
                              ]
                            },
                            "event_contract_investment_info": {
                              "type": "object",
                              "description": "Event contract account investment info",
                              "properties": {
                                "investment_experience": {
                                  "type": "string",
                                  "description": "Event contract investment experience. Value must be obtained from the Master Data API with data_type = INVESTMENT_EXPERIENCE."
                                },
                                "investment_knowledge": {
                                  "type": "string",
                                  "description": "Event contract investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE_V2."
                                },
                                "trade_per_year": {
                                  "type": "string",
                                  "description": "Event contract trade per year. Value must be obtained from the Master Data API with data_type = TRADE_PER_YEAR."
                                }
                              },
                              "required": [
                                "investment_experience",
                                "investment_knowledge",
                                "trade_per_year"
                              ]
                            },
                            "brokerage_investment_info": {
                              "type": "object",
                              "description": "Brokerage account investment info",
                              "properties": {
                                "investment_knowledge": {
                                  "type": "string",
                                  "description": "Investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE."
                                },
                                "investment_objective": {
                                  "type": "string",
                                  "description": "Investment objective. Value must be obtained from the Master Data API with data_type = INVESTMENT_OBJECTIVE."
                                },
                                "time_horizon": {
                                  "type": "string",
                                  "description": "Time horizon. Value must be obtained from the Master Data API with data_type = TIME_HORIZON."
                                }
                              },
                              "required": [
                                "investment_knowledge",
                                "investment_objective",
                                "time_horizon"
                              ]
                            }
                          }
                        }
                      },
                      "required": [
                        "type"
                      ]
                    }
                  }
                },
                "description": "Account Info Update Form",
                "title": "AccountUpdateForm"
              }
            }
          },
          "description": "Account Info Update Submit Request",
          "title": "AccountUpdateSubmitRequest"
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
              "account_id",
              "client_request_id"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Unique client-generated identifier used for request tracking, idempotency and audit purposes.",
                "example": "4F5A13B5FDFD43D29C7A3A856B2458A1"
              },
              "account_id": {
                "type": "string",
                "description": "System-generated unique identifier assigned to the account application record.",
                "example": "943a9802f6c14983b3b4755c69c01717"
              }
            },
            "description": "Account Update Response",
            "title": "AccountUpdateResult"
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
    "client_request_id": "123e4567e89b12d3a456426614174000",
    "account_id": "943a9802f6c14983b3b4755c69c01717",
    "forms": [
      {
        "form_info": {
          "form_code": "UPDATE_ACCOUNT_FORM",
          "version": "1.0"
        },
        "json_data": {
          "type": "object",
          "description": "This form is used to submit post-onboarding account update requests for an existing account. It serves as a unified entry point for multiple account maintenance and compliance-related update scenarios after the account has been successfully opened. The specific update scenario is determined by the `type` field. Based on the selected type, the corresponding update information object must be provided. This form is commonly used for periodic reviews, regulatory compliance verification, risk control follow-ups, and other account data maintenance processes to ensure that account records remain accurate, valid, and up to date.",
          "form_code": "ACCOUNT_UPDATE_FORM",
          "properties": {
            "type": {
              "type": "string",
              "description": "Specifies the category of account update to be performed after account opening. This field determines which account update scenario applies and which corresponding information object is required in the request. Each type represents a distinct account maintenance or compliance update operation."
            },
            "commission_info": {
              "type": "object",
              "description": "This object is applicable when the update type is COMMISSION_UPDATE.",
              "properties": {
                "commission_code": {
                  "type": "string",
                  "description": "Commission Code for account, determining the commission structure."
                }
              },
              "required": [
                "commission_code"
              ]
            },
            "identity_info": {
              "type": "object",
              "description": "This object is applicable when the update type is IDENTITY_UPDATE. This form is used to collect updated identification information for an existing account. It is required during account maintenance, periodic review, or compliance verification to ensure the accuracy and validity of identity records.",
              "properties": {
                "id_doc_type": {
                  "type": "string",
                  "description": "Type of identification document provided by the applicant. The value must be retrieved from the Master Data API where data_type = ID_DOC_TYPE."
                },
                "id_number": {
                  "type": "string",
                  "description": "Identification document number exactly as shown on the provided ID."
                },
                "id_expiry_date": {
                  "type": "string",
                  "description": "Expiration date of the identification document.",
                  "format": "yyyy-MM-dd"
                },
                "picture_front_doc_id": {
                  "type": "string",
                  "description": "Document ID referencing the uploaded image of the front side of the identification document."
                },
                "picture_back_doc_id": {
                  "type": "string",
                  "description": "Document ID referencing the uploaded image of the back side of the identification document."
                }
              },
              "required": [
                "id_doc_type",
                "id_number",
                "id_expiry_date",
                "picture_front_doc_id",
                "picture_back_doc_id"
              ]
            },
            "employment_info": {
              "type": "object",
              "description": "This object is applicable when the update type is EMPLOYMENT_UPDATE. Collects updated employment information required during account maintenance or compliance review.",
              "properties": {
                "employment_type": {
                  "type": "string",
                  "description": "Type of employment. Value must be obtained from the Master Data API with data_type = EMPLOYMENT_TYPE."
                },
                "company_name": {
                  "type": "string",
                  "description": "Name of the applicant current employer.",
                  "max_length": 128
                },
                "occupation": {
                  "type": "string",
                  "description": "Specific occupation under the employment type. Value must be obtained from the Master Data API with data_type = OCCUPATION."
                },
                "position": {
                  "type": "string",
                  "description": "Job position or title. Value must be obtained from the Master Data API with data_type = POSITION."
                },
                "years_employed": {
                  "type": "number",
                  "description": "Total duration (in years) the applicant has been employed with the current employer. Decimal values are allowed to represent partial years (e.g. 1.5)."
                },
                "employment_address": {
                  "type": "object",
                  "description": "Address of the applicant current employer.",
                  "properties": {
                    "country": {
                      "type": "string",
                      "description": "Employer address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                    },
                    "state": {
                      "type": "string",
                      "description": "Employer address state. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                    },
                    "city": {
                      "type": "string",
                      "description": "Employer address city.",
                      "max_length": 256
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street address of the employer. Multiple lines are allowed.",
                      "max_length": 512
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal code of the employer address.",
                      "max_length": 80
                    },
                    "neighborhood_code": {
                      "type": "string",
                      "description": "Neighborhood or district code of the employer address, if applicable.",
                      "max_length": 128
                    },
                    "apartment_number": {
                      "type": "string",
                      "description": "Apartment number of the employer address, if applicable.",
                      "max_length": 128
                    },
                    "building_number": {
                      "type": "string",
                      "description": "Building number of the employer address, if applicable.",
                      "max_length": 128
                    }
                  },
                  "required": [
                    "country",
                    "state",
                    "city",
                    "street_address",
                    "postal_code",
                    "neighborhood_code",
                    "apartment_number",
                    "building_number"
                  ]
                }
              },
              "required": [
                "employment_type",
                "company_name",
                "occupation",
                "position",
                "years_employed",
                "employment_address"
              ]
            },
            "trusted_contact_info": {
              "type": "object",
              "description": "This object is applicable when the update type is TRUSTED_UPDATE. Collects or updates trusted contact information for an account, used for account protection, compliance review, or situations where the account holder cannot be contacted.",
              "properties": {
                "first_name": {
                  "type": "string",
                  "description": "First name of the trusted contact person."
                },
                "middle_name": {
                  "type": "string",
                  "description": "Middle name of the trusted contact person, if applicable."
                },
                "last_name": {
                  "type": "string",
                  "description": "Last name of the trusted contact person."
                },
                "phone_country": {
                  "type": "string",
                  "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                },
                "phone_number": {
                  "type": "string",
                  "description": "Primary contact phone number of the trusted contact person.",
                  "example": "+1-2025550124",
                  "max_length": 20
                },
                "relationship": {
                  "type": "string",
                  "description": "Relationship between the account holder and the trusted contact person.",
                  "max_length": 32
                },
                "trusted_address": {
                  "type": "object",
                  "description": "Residential or mailing address of the trusted contact person.",
                  "properties": {
                    "country": {
                      "type": "string",
                      "description": "Country of the trusted contact person address. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                    },
                    "state": {
                      "type": "string",
                      "description": "State or province of the trusted contact person address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                    },
                    "city": {
                      "type": "string",
                      "description": "City of the trusted contact person address.",
                      "max_length": 256
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street address of the trusted contact person. Multiple address lines are allowed.",
                      "max_length": 512
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal or ZIP code of the trusted contact person address.",
                      "max_length": 80
                    },
                    "neighborhood_code": {
                      "type": "string",
                      "description": "Neighborhood or district code of the trusted contact person address, if applicable.",
                      "max_length": 128
                    },
                    "apartment_number": {
                      "type": "string",
                      "description": "Apartment or unit number of the trusted contact person address, if applicable.",
                      "max_length": 128
                    },
                    "building_number": {
                      "type": "string",
                      "description": "Building number of the trusted contact person address, if applicable.",
                      "max_length": 128
                    }
                  },
                  "required": [
                    "country",
                    "state",
                    "city",
                    "street_address",
                    "postal_code",
                    "neighborhood_code"
                  ]
                }
              },
              "required": [
                "first_name",
                "last_name",
                "phone_country",
                "phone_number",
                "relationship",
                "trusted_address"
              ]
            },
            "email_address_info": {
              "type": "object",
              "description": "This object is applicable when the update type is EMAIL_ADDRESS_UPDATE. Collects or updates the primary email address associated with an account. This information is used for account communication, security notifications, and compliance-related correspondence.",
              "properties": {
                "email_address": {
                  "type": "string",
                  "description": "Email address",
                  "max_length": 256
                }
              },
              "required": [
                "email_address"
              ]
            },
            "address_info": {
              "type": "object",
              "description": "This object is applicable when the update type is ADDRESS_UPDATE. Collects or updates address information associated with an account, including residential and mailing addresses. This information is used for account administration, regulatory compliance, and official correspondence.",
              "properties": {
                "id_address": {
                  "type": "object",
                  "description": "Residential address of the account holder as recorded for identity verification and regulatory purposes.",
                  "properties": {
                    "country": {
                      "type": "string",
                      "description": "Country of residence. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                    },
                    "state": {
                      "type": "string",
                      "description": "State or province of residence. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                    },
                    "city": {
                      "type": "string",
                      "description": "City of residence.",
                      "max_length": 256
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street address of the account holder residence. Multiple address lines are allowed.",
                      "max_length": 512
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal or ZIP code of the residential address.",
                      "max_length": 80
                    },
                    "neighborhood_code": {
                      "type": "string",
                      "description": "Neighborhood or district code of the id address, if applicable.",
                      "max_length": 128
                    },
                    "apartment_number": {
                      "type": "string",
                      "description": "Apartment number of the id address, if applicable.",
                      "max_length": 128
                    },
                    "building_number": {
                      "type": "string",
                      "description": "Building number of the id address, if applicable.",
                      "max_length": 128
                    }
                  },
                  "required": [
                    "country",
                    "state",
                    "city",
                    "street_address",
                    "postal_code",
                    "neighborhood_code",
                    "apartment_number",
                    "building_number"
                  ]
                },
                "mail_address": {
                  "type": "object",
                  "description": "Mailing address used for receiving account-related correspondence, if different from the residential address.",
                  "properties": {
                    "country": {
                      "type": "string",
                      "description": "Mailing address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                    },
                    "state": {
                      "type": "string",
                      "description": "State or province of the mailing address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                    },
                    "city": {
                      "type": "string",
                      "description": "City of the mailing address.",
                      "max_length": 256
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street address of the mailing address. Multiple address lines are allowed.",
                      "max_length": 512
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal or ZIP code of the mailing address.",
                      "max_length": 80
                    },
                    "neighborhood_code": {
                      "type": "string",
                      "description": "Neighborhood or district code of the mailing address.",
                      "max_length": 128
                    },
                    "apartment_number": {
                      "type": "string",
                      "description": "Apartment number of the mailing address.",
                      "max_length": 128
                    },
                    "building_number": {
                      "type": "string",
                      "description": "Building number of the mailing address.",
                      "max_length": 128
                    }
                  },
                  "required": [
                    "country",
                    "state",
                    "city",
                    "street_address",
                    "postal_code",
                    "neighborhood_code",
                    "apartment_number",
                    "building_number"
                  ]
                }
              },
              "required": [
                "id_address",
                "mail_address"
              ]
            },
            "phone_number_info": {
              "type": "object",
              "description": "This object is applicable when the update type is PHONE_NUMBER_UPDATE. Collects or updates the primary phone number associated with an account, used for account communication, security verification, and regulatory notifications.",
              "properties": {
                "phone_country": {
                  "type": "string",
                  "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                },
                "phone_number": {
                  "type": "string",
                  "description": "Primary contact phone number associated with the account, in international format (E.164). Used for account notifications, verification, and security-related communications.",
                  "example": "+1-2025550123",
                  "max_length": 20
                }
              },
              "required": [
                "phone_country",
                "phone_number"
              ]
            },
            "politically_exposed_info": {
              "type": "object",
              "description": "This object is applicable when the update type is PEP_UPDATE. Collects or updates Politically Exposed Person (PEP) information for an account. This information is used for compliance review, enhanced due diligence (EDD), and regulatory reporting.",
              "properties": {
                "is_politically_exposed": {
                  "type": "boolean",
                  "description": "Indicates whether the account holder or any close associate or immediate family member is a Politically Exposed Person (PEP).  If true, the politically_exposed_info object must be provided with detailed disclosure information.  Please note: PEP accounts are currently not supported. Any request submitted with is_politically_exposed = true will be rejected and a failure response will be returned."
                },
                "politically_exposed_info": {
                  "type": "object",
                  "description": "Detail information about Politically Exposed Person (PEP) associated with the account holder.",
                  "properties": {
                    "pep_organization": {
                      "type": "string",
                      "description": "Name of the related political organization.",
                      "max_length": 256
                    },
                    "immediate_family": {
                      "type": "array",
                      "description": "Name of family members, including former spouses.",
                      "items": {
                        "type": "string",
                        "description": ""
                      }
                    }
                  },
                  "required": [
                    "pep_organization",
                    "immediate_family"
                  ]
                }
              },
              "required": [
                "is_politically_exposed"
              ]
            },
            "company_control_info": {
              "type": "object",
              "description": "This object is applicable when the update type is COMPANY_CONTROL_UPDATE. Collects or updates information about the applicant status as a control person of a public company. This information is used for compliance review, enhanced due diligence (EDD), and regulatory reporting.",
              "properties": {
                "is_company_control_person": {
                  "type": "boolean",
                  "description": "Indicates whether the applicant is a control person of a public company. If true, `public_company_infos` must be provided with details of the companies controlled."
                },
                "public_company_infos": {
                  "type": "array",
                  "description": "List of public companies controlled by the applicant. Each item must include the company stock symbol.",
                  "items": {
                    "type": "object",
                    "description": "",
                    "properties": {
                      "symbol": {
                        "type": "string",
                        "description": "Stock symbol of the public company controlled by the applicant, as listed on the exchange."
                      }
                    },
                    "required": [
                      "symbol"
                    ]
                  }
                }
              },
              "required": [
                "is_company_control_person"
              ]
            },
            "disclosure_info": {
              "type": "object",
              "description": "This object is applicable when the update type is DISCLOSURE_UPDATE. Collects or updates the applicant regulatory and compliance-related affiliations or disclosures. This form is used for compliance review, enhanced due diligence (EDD), and regulatory reporting, and may include multiple types of information beyond exchange or FINRA affiliations.",
              "properties": {
                "brokerage_disclosure_info": {
                  "type": "object",
                  "description": "Brokerage disclosure info",
                  "properties": {
                    "is_exchange_or_finra_affiliated": {
                      "type": "boolean",
                      "description": "Indicates whether the applicant is employed by, or directly associated with, an exchange or a FINRA member firm. If true, exchange_affiliations details must be provided."
                    },
                    "exchange_affiliations": {
                      "type": "array",
                      "description": "List of the applicant affiliations with exchanges or FINRA member firms. Required if is_exchange_or_finra_affiliated is true.",
                      "items": {
                        "type": "object",
                        "description": "",
                        "properties": {
                          "company_name": {
                            "type": "string",
                            "description": "Name of the exchange or FINRA member firm the applicant is affiliated with.",
                            "max_length": 128
                          },
                          "country": {
                            "type": "string",
                            "description": "Country of the affiliated exchange or FINRA firm. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of the affiliated exchange or FINRA firm. Value must be obtained from the Master Data API with data_type = STATE."
                          },
                          "city": {
                            "type": "string",
                            "description": "City of the affiliated exchange or FINRA firm.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the affiliated exchange or FINRA firm. Multiple lines are allowed.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the affiliated exchange or FINRA firm.",
                            "max_length": 80
                          }
                        },
                        "required": [
                          "company_name",
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    }
                  },
                  "required": [
                    "is_exchange_or_finra_affiliated"
                  ]
                }
              },
              "required": [
                "brokerage_disclosure_info"
              ]
            },
            "dependent_info": {
              "type": "object",
              "description": "This object is applicable when the update type is DEPENDENT_UPDATE. Collects or updates the number of dependents associated with an account holder. This information is used for account maintenance, compliance review, and regulatory reporting purposes.",
              "properties": {
                "num_dependents": {
                  "type": "number",
                  "description": "Number of dependents financially or legally associated with the account holder."
                }
              },
              "required": [
                "num_dependents"
              ]
            },
            "citizenship_info": {
              "type": "object",
              "description": "This object is applicable when the update type is CITIZENSHIP_UPDATE. Collects or updates the country of citizenship for an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
              "properties": {
                "country_of_citizenship": {
                  "type": "string",
                  "description": "Country of citizenship of the account holder. The value must be retrieved from the Master Data API with data_type = COUNTRY."
                }
              },
              "required": [
                "country_of_citizenship"
              ]
            },
            "gender_info": {
              "type": "object",
              "description": "This object is applicable when the update type is GENDER_UPDATE. Collects or updates the gender of an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
              "properties": {
                "gender": {
                  "type": "string",
                  "description": "Gender of the account holder. Allowed values: F (Female), M (Male)."
                }
              },
              "required": [
                "gender"
              ]
            },
            "marital_info": {
              "type": "object",
              "description": "This object is applicable when the update type is MARITAL_UPDATE. Collects or updates the marital status of an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
              "properties": {
                "marital_status": {
                  "type": "string",
                  "description": "Marital status of the account holder. Value must be obtained from the Master Data API with data_type = MARITAL_STATUS."
                }
              },
              "required": [
                "marital_status"
              ]
            },
            "financial_investment_info": {
              "type": "object",
              "description": "This object is applicable when the update type is FINANCIAL_INVESTMENT_UPDATE. Collects or updates the investment profile and financial information of an account holder. This information is used for account maintenance, risk assessment, and regulatory compliance.",
              "properties": {
                "financial_info": {
                  "type": "object",
                  "description": "Financial information of the account holder.",
                  "properties": {
                    "liquidity_needs": {
                      "type": "string",
                      "description": "Indicates the client liquidity needs or preference for cash accessibility. Value must be obtained from the Master Data API with data_type = LIQUIDITY_NEEDS."
                    },
                    "annual_income": {
                      "type": "string",
                      "description": "Represents the client annual income. Value must be obtained from the Master Data API with data_type = ANNUAL_INCOME."
                    },
                    "total_net_worth": {
                      "type": "string",
                      "description": "Represents the client total net worth in USD. Value must be obtained from the Master Data API with data_type = TOTAL_NET_WORTH."
                    },
                    "liquid_net_worth": {
                      "type": "string",
                      "description": "Represents the client liquid net worth (USD), including cash and easily marketable assets. Value must be obtained from the Master Data API with data_type = LIQUID_NET_WORTH."
                    }
                  },
                  "required": [
                    "liquidity_needs",
                    "annual_income",
                    "total_net_worth",
                    "liquid_net_worth"
                  ]
                },
                "event_contract_investment_info": {
                  "type": "object",
                  "description": "Event contract account investment info",
                  "properties": {
                    "investment_experience": {
                      "type": "string",
                      "description": "Event contract investment experience. Value must be obtained from the Master Data API with data_type = INVESTMENT_EXPERIENCE."
                    },
                    "investment_knowledge": {
                      "type": "string",
                      "description": "Event contract investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE_V2."
                    },
                    "trade_per_year": {
                      "type": "string",
                      "description": "Event contract trade per year. Value must be obtained from the Master Data API with data_type = TRADE_PER_YEAR."
                    }
                  },
                  "required": [
                    "investment_experience",
                    "investment_knowledge",
                    "trade_per_year"
                  ]
                },
                "brokerage_investment_info": {
                  "type": "object",
                  "description": "Brokerage account investment info",
                  "properties": {
                    "investment_knowledge": {
                      "type": "string",
                      "description": "Investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE."
                    },
                    "investment_objective": {
                      "type": "string",
                      "description": "Investment objective. Value must be obtained from the Master Data API with data_type = INVESTMENT_OBJECTIVE."
                    },
                    "time_horizon": {
                      "type": "string",
                      "description": "Time horizon. Value must be obtained from the Master Data API with data_type = TIME_HORIZON."
                    }
                  },
                  "required": [
                    "investment_knowledge",
                    "investment_objective",
                    "time_horizon"
                  ]
                }
              }
            }
          },
          "required": [
            "type"
          ]
        }
      }
    ]
  },
  "postman": {
    "name": "Update Account",
    "description": {
      "content": "Submit a request to update an existing account. This API only updates the fields submitted in the forms data. Any attributes not included in the request will remain unchanged. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Account Events](../fd-events/account-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "update"
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

## Close Account

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/close-account.md>

### Close Account

Closes an existing account. The operation is irreversible and requires valid account information. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Account Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/accounts/close",
  "method": "post",
  "tags": [
    "Accounts"
  ],
  "description": "Closes an existing account. The operation is irreversible and requires valid account information. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Account Events](../fd-events/account-events)",
  "operationId": "closeAccount",
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
            "account_id",
            "client_request_id",
            "close_reason"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client-supplied request ID, which is unique for each account closure request.",
              "example": "ac0753e04f5a42099a02579225527856"
            },
            "account_id": {
              "type": "string",
              "description": "Unique identifier of the account to be closed.",
              "example": "ac0753e04f5a42099a02579225527856"
            },
            "close_reason": {
              "type": "string",
              "description": "Optional detailed textual explanation of the account closure reason. Limited to 512 characters.",
              "example": "Switching to another broker for better trading conditions."
            }
          },
          "description": "Request payload submitted when a client wants to close an account.",
          "title": "AccountCloseReq"
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
              "account_id",
              "client_request_id"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client-generated unique identifier for idempotency and audit purposes.",
                "example": "A9C3F2E84D0B4A7E9D1F8C6B12345678"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID to be closed.",
                "example": "ACC123456789"
              }
            },
            "description": "Account Closure Status",
            "title": "AccountCloseStatusResult"
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
    "client_request_id": "ac0753e04f5a42099a02579225527856",
    "account_id": "ac0753e04f5a42099a02579225527856",
    "close_reason": "Switching to another broker for better trading conditions."
  },
  "postman": {
    "name": "Close Account",
    "description": {
      "content": "Closes an existing account. The operation is irreversible and requires valid account information. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Account Events](../fd-events/account-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "close"
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

## Account Application Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/get-account-application-detail.md>

### Get Account Application Detail

Retrieves detailed information about a specific account application.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/accounts/applications/get",
  "method": "get",
  "tags": [
    "Accounts"
  ],
  "description": "Retrieves detailed information about a specific account application.",
  "operationId": "getAccountApplicationDetail",
  "parameters": [
    {
      "name": "application_id",
      "in": "query",
      "description": "Unique identifier of the account application to retrieve details for.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "a1b2c3d4e5f60718293a4b5c6d7e8f90"
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
              "additional_infos",
              "application_id",
              "contact_info",
              "customer_info",
              "disclosure_info",
              "employment_info",
              "financial_info",
              "tax_info"
            ],
            "type": "object",
            "properties": {
              "application_id": {
                "type": "string",
                "description": "Application ID",
                "example": "APP123456789"
              },
              "resubmission_codes": {
                "type": "array",
                "description": "Resubmission Code, only present for resubmission applications",
                "example": [
                  "OCR_SUSPENDED",
                  "CIP_SUSPENDED"
                ],
                "items": {
                  "type": "string",
                  "description": "Application Resubmission Code",
                  "example": "[\"OCR_SUSPENDED\",\"CIP_SUSPENDED\"]",
                  "enum": [
                    "OCR_SUSPENDED",
                    "CIP_SUSPENDED",
                    "NOT_A_GOVERNMENT_ISSUED",
                    "EXPIRED_ID",
                    "ID_CARD_ILLEGIBLE",
                    "ID_ON_SCREEN",
                    "UNMATCHED_ID_NUMBER_WITH_THE_ID_PICTURE",
                    "UNMATCHED_SSN_PICTURE",
                    "UNMATCHED_ID_EXPIRE_DATE",
                    "UNMATCHED_NATIONALITY_WITH_THE_ID_PICTURE",
                    "UNMATCHED_NAME_WITH_THE_ID_PICTURE",
                    "UNMATCHED_DOB_WITH_THE_ID_PICTURE",
                    "SELFIE_NOT_MEET_THE_REQUIREMENTS",
                    "NEED_407_LETTERS"
                  ]
                }
              },
              "customer_info": {
                "required": [
                  "country_of_citizenship",
                  "date_of_birth",
                  "first_name",
                  "gender",
                  "id_doc_type",
                  "id_expiry_date",
                  "id_number",
                  "ip_address",
                  "last_name",
                  "marital_status",
                  "num_dependents",
                  "permanent_resident"
                ],
                "type": "object",
                "properties": {
                  "ip_address": {
                    "type": "string",
                    "description": "The IP address when the user submits the application.",
                    "example": "0.0.0.0"
                  },
                  "first_name": {
                    "type": "string",
                    "description": "First name",
                    "example": "John"
                  },
                  "middle_name": {
                    "type": "string",
                    "description": "Middle name",
                    "example": "A"
                  },
                  "last_name": {
                    "type": "string",
                    "description": "Last name",
                    "example": "Doe"
                  },
                  "date_of_birth": {
                    "type": "string",
                    "description": "Date of birth",
                    "example": "1990-01-01"
                  },
                  "gender": {
                    "type": "string",
                    "description": "Gender. Allowed values: F (Female), M (Male).",
                    "example": "M"
                  },
                  "country_of_citizenship": {
                    "type": "string",
                    "description": "Country of citizenship. Value must be obtained from the Master Data API with data_type = COUNTRY",
                    "example": "US"
                  },
                  "id_doc_type": {
                    "type": "string",
                    "description": "Type of identification document. Value must be obtained from the Master Data API with data_type = ID_DOC_TYPE.",
                    "example": "PASSPORT"
                  },
                  "picture_front_doc_id": {
                    "type": "string",
                    "description": "Picture front document id",
                    "example": "DOC123456789"
                  },
                  "picture_back_doc_id": {
                    "type": "string",
                    "description": "Picture back document id",
                    "example": "DOC987654321"
                  },
                  "id_number": {
                    "type": "string",
                    "description": "Number of the identification document.",
                    "example": "X1234567"
                  },
                  "id_expiry_date": {
                    "type": "string",
                    "description": "Expiration date of the identification document.",
                    "example": "2030-12-31"
                  },
                  "selfie_img_key": {
                    "type": "string",
                    "description": "User’s selfie, face photo, or liveness document id",
                    "example": "SELFIE123456789"
                  },
                  "tax_id_doc_id": {
                    "type": "string",
                    "description": "Tax ID or SSN document id for tax/SSN verification",
                    "example": "TAXID123456789"
                  },
                  "permanent_resident": {
                    "type": "string",
                    "description": "Permanent resident. Allowed values: YES, NO.",
                    "example": "NO"
                  },
                  "marital_status": {
                    "type": "string",
                    "description": "The user's marital status. Value must be obtained from the Master Data API with data_type = MARITAL_STATUS.",
                    "example": "SINGLE"
                  },
                  "num_dependents": {
                    "type": "integer",
                    "description": "Num dependents",
                    "format": "int32",
                    "example": 1
                  },
                  "ext_attr_list": {
                    "type": "array",
                    "description": "Ext attr list",
                    "items": {
                      "required": [
                        "name",
                        "value"
                      ],
                      "type": "object",
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "Attribute Name",
                          "example": "custom_attribute_1"
                        },
                        "value": {
                          "type": "string",
                          "description": "Attribute Value",
                          "example": "custom_value_1"
                        }
                      },
                      "description": "Extended Attributes",
                      "title": "ExtAttr"
                    }
                  }
                },
                "description": "Customer Info",
                "title": "CustomerInfo"
              },
              "contact_info": {
                "required": [
                  "email_address",
                  "id_address",
                  "phone_country",
                  "phone_number"
                ],
                "type": "object",
                "properties": {
                  "phone_country": {
                    "type": "string",
                    "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                    "example": "USA"
                  },
                  "phone_number": {
                    "type": "string",
                    "description": "Phone number",
                    "example": "+1-2025550143"
                  },
                  "email_address": {
                    "type": "string",
                    "description": "Email address",
                    "example": "demo@example.com"
                  },
                  "id_address": {
                    "required": [
                      "city",
                      "country",
                      "postal_code",
                      "state",
                      "street_address"
                    ],
                    "type": "object",
                    "properties": {
                      "country": {
                        "type": "string",
                        "description": "Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                        "example": "USA"
                      },
                      "state": {
                        "type": "string",
                        "description": "State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA).",
                        "example": "USA-CA"
                      },
                      "city": {
                        "type": "string",
                        "description": "City.",
                        "example": "San Francisco"
                      },
                      "street_address": {
                        "type": "string",
                        "description": "Street address",
                        "example": "123 Main St"
                      },
                      "postal_code": {
                        "type": "string",
                        "description": "Postal code",
                        "example": "94105"
                      },
                      "neighborhood_code": {
                        "type": "string",
                        "description": "Neighborhood code.",
                        "example": "NBHD12345"
                      },
                      "apartment_number": {
                        "type": "string",
                        "description": "Apartment number.",
                        "example": "Apt 4B"
                      },
                      "building_number": {
                        "type": "string",
                        "description": "Building number.",
                        "example": "Building 5"
                      }
                    },
                    "description": "Address Info",
                    "title": "AddressInfo"
                  },
                  "mail_address": {
                    "required": [
                      "city",
                      "country",
                      "postal_code",
                      "state",
                      "street_address"
                    ],
                    "type": "object",
                    "properties": {
                      "country": {
                        "type": "string",
                        "description": "Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                        "example": "USA"
                      },
                      "state": {
                        "type": "string",
                        "description": "State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA).",
                        "example": "USA-CA"
                      },
                      "city": {
                        "type": "string",
                        "description": "City.",
                        "example": "San Francisco"
                      },
                      "street_address": {
                        "type": "string",
                        "description": "Street address",
                        "example": "123 Main St"
                      },
                      "postal_code": {
                        "type": "string",
                        "description": "Postal code",
                        "example": "94105"
                      },
                      "neighborhood_code": {
                        "type": "string",
                        "description": "Neighborhood code.",
                        "example": "NBHD12345"
                      },
                      "apartment_number": {
                        "type": "string",
                        "description": "Apartment number.",
                        "example": "Apt 4B"
                      },
                      "building_number": {
                        "type": "string",
                        "description": "Building number.",
                        "example": "Building 5"
                      }
                    },
                    "description": "Address Info",
                    "title": "AddressInfo"
                  }
                },
                "description": "Contact Information",
                "title": "ContactInfo"
              },
              "trusted_contact_info": {
                "required": [
                  "first_name",
                  "last_name",
                  "phone_country",
                  "phone_number",
                  "relationship",
                  "trusted_address"
                ],
                "type": "object",
                "properties": {
                  "first_name": {
                    "type": "string",
                    "description": "First name of the trusted contact.",
                    "example": "John"
                  },
                  "middle_name": {
                    "type": "string",
                    "description": "Middle name of the trusted contact.",
                    "example": "A."
                  },
                  "last_name": {
                    "type": "string",
                    "description": "Last name of the trusted contact.",
                    "example": "Doe"
                  },
                  "phone_country": {
                    "type": "string",
                    "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                    "example": "USA"
                  },
                  "phone_number": {
                    "type": "string",
                    "description": "Phone number",
                    "example": "+1-2025550143"
                  },
                  "relationship": {
                    "type": "string",
                    "description": "Relationship between the account holder and the trusted contact person. ",
                    "example": "Parent"
                  },
                  "trusted_address": {
                    "required": [
                      "city",
                      "country",
                      "postal_code",
                      "state",
                      "street_address"
                    ],
                    "type": "object",
                    "properties": {
                      "country": {
                        "type": "string",
                        "description": "Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                        "example": "USA"
                      },
                      "state": {
                        "type": "string",
                        "description": "State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA).",
                        "example": "USA-CA"
                      },
                      "city": {
                        "type": "string",
                        "description": "City.",
                        "example": "San Francisco"
                      },
                      "street_address": {
                        "type": "string",
                        "description": "Street address",
                        "example": "123 Main St"
                      },
                      "postal_code": {
                        "type": "string",
                        "description": "Postal code",
                        "example": "94105"
                      },
                      "neighborhood_code": {
                        "type": "string",
                        "description": "Neighborhood code.",
                        "example": "NBHD12345"
                      },
                      "apartment_number": {
                        "type": "string",
                        "description": "Apartment number.",
                        "example": "Apt 4B"
                      },
                      "building_number": {
                        "type": "string",
                        "description": "Building number.",
                        "example": "Building 5"
                      }
                    },
                    "description": "Trusted Contact Address",
                    "title": "TrustedContactAddress"
                  }
                },
                "description": "Trusted Contact Info",
                "title": "TrustedContactInfo"
              },
              "employment_info": {
                "required": [
                  "company_name",
                  "employment_address",
                  "employment_type",
                  "occupation",
                  "position",
                  "years_employed"
                ],
                "type": "object",
                "properties": {
                  "employment_type": {
                    "type": "string",
                    "description": "Type of employment. Value must be obtained from the Master Data API with data_type = EMPLOYMENT_TYPE.",
                    "example": "EMPLOYED"
                  },
                  "company_name": {
                    "type": "string",
                    "description": "Name of the applicant current employer.",
                    "example": "Tech Corp"
                  },
                  "occupation": {
                    "type": "string",
                    "description": "Specific occupation under the employment type. Value must be obtained from the Master Data API with data_type = OCCUPATION.",
                    "example": "ACCOUNTING_AUDITING"
                  },
                  "position": {
                    "type": "string",
                    "description": "Position. Value must be obtained from the Master Data API with data_type = POSITION",
                    "example": "SECRETARY"
                  },
                  "years_employed": {
                    "type": "string",
                    "description": "Total duration (in years) the applicant has been employed with the current employer. Decimal values are allowed to represent partial years (e.g. 1.5).",
                    "example": "1.5"
                  },
                  "employment_address": {
                    "required": [
                      "city",
                      "country",
                      "postal_code",
                      "state",
                      "street_address"
                    ],
                    "type": "object",
                    "properties": {
                      "country": {
                        "type": "string",
                        "description": "Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                        "example": "USA"
                      },
                      "state": {
                        "type": "string",
                        "description": "State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA).",
                        "example": "USA-CA"
                      },
                      "city": {
                        "type": "string",
                        "description": "City.",
                        "example": "San Francisco"
                      },
                      "street_address": {
                        "type": "string",
                        "description": "Street address",
                        "example": "123 Main St"
                      },
                      "postal_code": {
                        "type": "string",
                        "description": "Postal code",
                        "example": "94105"
                      },
                      "neighborhood_code": {
                        "type": "string",
                        "description": "Neighborhood code.",
                        "example": "NBHD12345"
                      },
                      "apartment_number": {
                        "type": "string",
                        "description": "Apartment number.",
                        "example": "Apt 4B"
                      },
                      "building_number": {
                        "type": "string",
                        "description": "Building number.",
                        "example": "Building 5"
                      }
                    },
                    "description": "Employment Address",
                    "title": "EmploymentAddress"
                  }
                },
                "description": "Employment Info",
                "title": "EmploymentInfo"
              },
              "financial_info": {
                "required": [
                  "annual_income",
                  "liquid_net_worth",
                  "liquidity_needs",
                  "total_net_worth"
                ],
                "type": "object",
                "properties": {
                  "liquidity_needs": {
                    "type": "string",
                    "description": "Indicates the client liquidity needs or preference for cash accessibility. Value must be obtained from the Master Data API with data_type = LIQUIDITY_NEEDS.",
                    "example": "VERY_IMPORTANT"
                  },
                  "annual_income": {
                    "type": "string",
                    "description": "Annual income. Value must be obtained from the Master Data API with data_type = ANNUAL_INCOME.",
                    "example": "LEVEL1"
                  },
                  "total_net_worth": {
                    "type": "string",
                    "description": "Client's total net worth (USD). Value must be obtained from the Master Data API with data_type = TOTAL_NET_WORTH.",
                    "example": "LEVEL1"
                  },
                  "liquid_net_worth": {
                    "type": "string",
                    "description": "Liquidity assets. Value must be obtained from the Master Data API with data_type = LIQUID_NET_WORTH.",
                    "example": "LEVEL1"
                  }
                },
                "description": "Financial Info",
                "title": "FinancialInfo"
              },
              "disclosure_info": {
                "required": [
                  "is_cftc_or_nfa_registered",
                  "is_company_control_person",
                  "is_politically_exposed"
                ],
                "type": "object",
                "properties": {
                  "is_politically_exposed": {
                    "type": "boolean",
                    "description": "Indicates whether the account holder or any close associate or immediate family member is a Politically Exposed Person (PEP).  If true, the politically_exposed_info object must be provided with detailed disclosure information.  Please note: PEP accounts are currently not supported. Any request submitted with is_politically_exposed = true will be rejected and a failure response will be returned.",
                    "example": false
                  },
                  "politically_exposed_info": {
                    "type": "object",
                    "properties": {
                      "pep_organization": {
                        "type": "string",
                        "description": "Name of the related political organization.",
                        "example": "United Nations"
                      },
                      "immediate_family": {
                        "type": "array",
                        "description": "Name of family members, including former spouses.",
                        "example": "Secretary-General",
                        "items": {
                          "type": "string",
                          "description": "Name of family members, including former spouses.",
                          "example": "Secretary-General"
                        }
                      }
                    },
                    "description": "Politically Exposed Info",
                    "title": "PoliticallyExposedInfo"
                  },
                  "is_cftc_or_nfa_registered": {
                    "type": "boolean",
                    "description": "Indicates whether the applicant is registered with the CFTC or is a member of the NFA.",
                    "example": false
                  },
                  "is_company_control_person": {
                    "type": "boolean",
                    "description": "Indicates whether the applicant is a control person of a public company. If true, `public_company_infos` must be provided with details of the companies controlled.",
                    "example": false
                  },
                  "public_company_infos": {
                    "type": "array",
                    "description": "List of public companies controlled by the applicant. Each item must include the company stock symbol.",
                    "items": {
                      "required": [
                        "symbol"
                      ],
                      "type": "object",
                      "properties": {
                        "symbol": {
                          "type": "string",
                          "description": "Stock symbol of the public company controlled by the applicant, as listed on the exchange.",
                          "example": "AAPL"
                        }
                      },
                      "description": "Public Company Info",
                      "title": "PublicCompanyInfo"
                    }
                  }
                },
                "description": "Disclosure Info",
                "title": "DisclosureInfo"
              },
              "tax_info": {
                "required": [
                  "tax_id",
                  "tax_type"
                ],
                "type": "object",
                "properties": {
                  "tax_type": {
                    "type": "string",
                    "description": "Tax Type. Value must be obtained from the Master Data API with data_type = TAX_TYPE.",
                    "example": "SSN"
                  },
                  "tax_id": {
                    "type": "string",
                    "description": "Tax ID",
                    "example": "123-45-6789"
                  }
                },
                "description": "Tax Info",
                "title": "TaxInfo"
              },
              "w8ben_info": {
                "required": [
                  "foreign_tax_id",
                  "sign_date",
                  "treaty_country"
                ],
                "type": "object",
                "properties": {
                  "foreign_tax_id": {
                    "type": "string",
                    "description": "Foreign tax ID (FTIN) issued by non-U.S. authority (W-8BEN Line 6)"
                  },
                  "treaty_country": {
                    "type": "string",
                    "description": "Tax treaty country for reduced withholding. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                  },
                  "sign_date": {
                    "type": "string",
                    "description": "Date when the W-8BEN form was signed by the customer. The date should be in ISO 8601 format (YYYY-MM-DD).",
                    "example": "2026-03-09"
                  }
                },
                "description": "W8ben Info",
                "title": "W8benInfo"
              },
              "additional_infos": {
                "type": "array",
                "description": "Additional Information for the application review process. ",
                "items": {
                  "required": [
                    "form_code",
                    "is_exchange_or_finra_affiliated",
                    "rep_code",
                    "status"
                  ],
                  "type": "object",
                  "properties": {
                    "status": {
                      "type": "string",
                      "description": "Application Status",
                      "enum": [
                        "SUBMITTED",
                        "PENDING_REVIEW",
                        "ACTION_REQUIRED",
                        "REJECTED",
                        "APPROVED"
                      ]
                    },
                    "form_code": {
                      "type": "string",
                      "description": "Form code",
                      "example": "BROKERAGE_ADDITIONAL_FORM"
                    },
                    "rep_code": {
                      "type": "string",
                      "description": "Rep code",
                      "example": "REP001"
                    },
                    "account_id": {
                      "type": "string",
                      "description": "Unique account identifier. Populated only when the status is APPROVED; otherwise, this field is empty.",
                      "example": "5RHU6AG2HQRSC7JC2D772Q534A"
                    },
                    "account_number": {
                      "type": "string",
                      "description": "Account number. Populated only when the status is APPROVED; otherwise, this field is empty.",
                      "example": "CUA0900026"
                    },
                    "is_exchange_or_finra_affiliated": {
                      "type": "boolean",
                      "description": "Indicates whether the applicant is employed by, or directly associated with, an exchange or a FINRA member firm. If true, exchange_affiliations details must be provided.",
                      "example": false
                    },
                    "exchange_or_finra_affiliations": {
                      "type": "array",
                      "description": "List of the applicant affiliations with exchanges or FINRA member firms. Required if is_exchange_or_finra_affiliated is true.",
                      "items": {
                        "required": [
                          "city",
                          "company_name",
                          "country",
                          "postal_code",
                          "state",
                          "street_address"
                        ],
                        "type": "object",
                        "properties": {
                          "company_name": {
                            "type": "string",
                            "description": "Company name",
                            "example": "Finance Corp"
                          },
                          "country": {
                            "type": "string",
                            "description": "Country of the affiliated exchange or FINRA firm. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code).",
                            "example": "USA"
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of the affiliated exchange or FINRA firm. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>",
                            "example": "USA-NY"
                          },
                          "city": {
                            "type": "string",
                            "description": "City of the affiliated exchange or FINRA firm.",
                            "example": "New York"
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address",
                            "example": "456 Finance St"
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal code",
                            "example": "10005"
                          }
                        },
                        "description": "Exchange Affiliations",
                        "title": "ExchangeAffiliations"
                      }
                    },
                    "investment_experience": {
                      "type": "string",
                      "description": "Investment experience. Value must be obtained from the Master Data API with data_type = INVESTMENT_EXPERIENCE."
                    },
                    "investment_knowledge": {
                      "type": "string",
                      "description": "Investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE. Note: For Event Contract , the value of this field should be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE_V2 instead."
                    },
                    "investment_objective": {
                      "type": "string",
                      "description": "Investment objective. Value must be obtained from the Master Data API with data_type = INVESTMENT_OBJECTIVE."
                    },
                    "trade_per_year": {
                      "type": "string",
                      "description": "Trade per year. Value must be obtained from the Master Data API with data_type = TRADE_PER_YEAR."
                    },
                    "time_horizon": {
                      "type": "string",
                      "description": "Time horizon. Value must be obtained from the Master Data API with data_type = TIME_HORIZON."
                    }
                  },
                  "description": "Additional Info",
                  "title": "AdditionalInfo"
                }
              },
              "document_attachments": {
                "type": "array",
                "description": "Document Attachments",
                "items": {
                  "required": [
                    "doc_type",
                    "document_id"
                  ],
                  "type": "object",
                  "properties": {
                    "doc_type": {
                      "type": "string",
                      "description": "Document Type"
                    },
                    "document_id": {
                      "type": "string",
                      "description": "Document ID"
                    }
                  },
                  "description": "Document Attachments",
                  "title": "DocumentAttachments"
                }
              }
            },
            "description": "Account Application Detail Result",
            "title": "AccountApplicationDetailResult"
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
    "name": "Get Account Application Detail",
    "description": {
      "content": "Retrieves detailed information about a specific account application.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "applications",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique identifier of the account application to retrieve details for.",
            "type": "text/plain"
          },
          "key": "application_id",
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

## List Forms

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-list.md>

### List Forms

Retrieves a list of available form codes.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/forms/list",
  "method": "get",
  "tags": [
    "Forms"
  ],
  "description": "Retrieves a list of available form codes.",
  "operationId": "getFormList",
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
            "type": "array",
            "items": {
              "type": "string",
              "example": "NEW_ACCOUNT_BASIC_FORM"
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
    "name": "List Forms",
    "description": {
      "content": "Retrieves a list of available form codes.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "forms",
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

## List Form Versions

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-version-list.md>

### List Form Versions

Retrieves a list of available versions for the specified form code.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/forms/versions/list",
  "method": "get",
  "tags": [
    "Forms"
  ],
  "description": "Retrieves a list of available versions for the specified form code.",
  "operationId": "getFormVersionList",
  "parameters": [
    {
      "name": "form_code",
      "in": "query",
      "description": "Unique identifier of the form, typically used to distinguish different form definitions.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "create_account_form"
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
              "type": "string",
              "example": "1.0"
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
    "name": "List Form Versions",
    "description": {
      "content": "Retrieves a list of available versions for the specified form code.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "forms",
        "versions",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique identifier of the form, typically used to distinguish different form definitions.",
            "type": "text/plain"
          },
          "key": "form_code",
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

## Form Content

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-content.md>

### Get Form Detail

Retrieves the JSON schema for the specified form code and version.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/forms/get",
  "method": "get",
  "tags": [
    "Forms"
  ],
  "description": "Retrieves the JSON schema for the specified form code and version.",
  "operationId": "getFormContent",
  "parameters": [
    {
      "name": "form_code",
      "in": "query",
      "description": "Unique identifier of the form, typically used to distinguish different form definitions.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "NEW_ACCOUNT_BASIC_FORM"
    },
    {
      "name": "version",
      "in": "query",
      "description": "Version number of the form definition",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "1.0"
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
      "description": "Request successful",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "form_code",
              "version"
            ],
            "type": "object",
            "properties": {
              "required": {
                "type": "array",
                "description": "List of required property names for the current schema object.",
                "example": [
                  "first_name",
                  "last_name"
                ],
                "items": {
                  "type": "string",
                  "description": "List of required property names for the current schema object.",
                  "example": "[\"first_name\",\"last_name\"]"
                }
              },
              "type": {
                "type": "string",
                "description": "Enumeration of supported JSON Schema property types used for form definition and validation.",
                "example": "object",
                "enum": [
                  "string",
                  "number",
                  "boolean",
                  "object",
                  "array"
                ]
              },
              "enum_values": {
                "type": "string",
                "description": "Enumerated values allowed for this property. Applies when the schema defines a fixed set of permitted values.",
                "example": "NEW_ACCOUNT_FORM"
              },
              "format": {
                "type": "string",
                "description": "Format keyword that further constrains the property type (e.g., date, date-time).",
                "example": "yyyy-MM-dd"
              },
              "description": {
                "type": "string",
                "description": "Description of the property, explaining its purpose or usage within the schema.",
                "example": "User age"
              },
              "example": {
                "type": "string",
                "description": "Representative example value to illustrate the expected input format.",
                "example": "2001-12-02"
              },
              "min_items": {
                "type": "integer",
                "description": "Minimum number of items allowed in the array. Applicable only when the property type is array.",
                "format": "int32",
                "example": 1
              },
              "max_items": {
                "type": "integer",
                "description": "Maximum number of items allowed in the array. Applicable only when the property type is array.",
                "format": "int32",
                "example": 5
              },
              "max_length": {
                "type": "integer",
                "description": "Maximum length allowed for string-type properties.",
                "format": "int32",
                "example": 64
              },
              "form_code": {
                "type": "string",
                "description": "Unique identifier of the form, typically used to distinguish different form definitions.",
                "example": "NEW_ACCOUNT_BASIC_FORM"
              },
              "version": {
                "type": "string",
                "description": "Version number of the form definition, used for compatibility and evolutionary control.",
                "example": "1.0"
              }
            },
            "description": "Form Content Response",
            "title": "FormContentResult"
          },
          "examples": {
            "NEW_ACCOUNT_BASIC_FORM": {
              "description": "NEW_ACCOUNT_BASIC_FORM",
              "value": {
                "type": "object",
                "description": "The New Account Basic Form collects the core information required for account opening.  This form cannot be submitted independently and must be used in combination with one or more applicable additional forms, depending on the account features or products being applied for.",
                "form_code": "NEW_ACCOUNT_BASIC_FORM",
                "properties": {
                  "customer_info": {
                    "type": "object",
                    "description": "Identity information",
                    "properties": {
                      "ip_address": {
                        "type": "string",
                        "description": "The IP address when the user submits the application."
                      },
                      "first_name": {
                        "type": "string",
                        "description": "First name of the applicant as it appears on the identification document.",
                        "max_length": 32
                      },
                      "middle_name": {
                        "type": "string",
                        "description": "Middle name of the applicant as it appears on the identification document, if applicable.",
                        "max_length": 32
                      },
                      "last_name": {
                        "type": "string",
                        "description": "Last name of the applicant as it appears on the identification document.",
                        "max_length": 32
                      },
                      "date_of_birth": {
                        "type": "string",
                        "description": "Date of birth of the applicant as shown on the identification document.",
                        "format": "yyyy-MM-dd"
                      },
                      "gender": {
                        "type": "string",
                        "description": "Gender of the account holder. Allowed values: F (Female), M (Male)."
                      },
                      "country_of_citizenship": {
                        "type": "string",
                        "description": "Country of citizenship of the account holder. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                      },
                      "id_doc_type": {
                        "type": "string",
                        "description": "Type of identification document provided by the applicant. The value must be retrieved from the Master Data API where data_type = ID_DOC_TYPE."
                      },
                      "id_number": {
                        "type": "string",
                        "description": "Identification document number exactly as shown on the provided ID."
                      },
                      "id_expiry_date": {
                        "type": "string",
                        "description": "Expiration date of the identification document.",
                        "format": "yyyy-MM-dd"
                      },
                      "permanent_resident": {
                        "type": "string",
                        "description": "Permanent resident. Allowed values: YES, NO."
                      },
                      "marital_status": {
                        "type": "string",
                        "description": "Marital status of the account holder. Value must be obtained from the Master Data API with data_type = MARITAL_STATUS."
                      },
                      "num_dependents": {
                        "type": "number",
                        "description": "Number of dependents financially or legally associated with the account holder.",
                        "example": "1"
                      },
                      "ext_attr_list": {
                        "type": "array",
                        "description": "Ext attr list",
                        "items": {
                          "type": "object",
                          "description": "Item",
                          "properties": {
                            "name": {
                              "type": "string",
                              "description": "The name of the extension."
                            },
                            "value": {
                              "type": "string",
                              "description": "Value"
                            }
                          }
                        }
                      }
                    },
                    "required": [
                      "ip_address",
                      "first_name",
                      "last_name",
                      "date_of_birth",
                      "gender",
                      "country_of_citizenship",
                      "id_doc_type",
                      "id_number",
                      "id_expiry_date",
                      "permanent_resident",
                      "marital_status",
                      "num_dependents"
                    ]
                  },
                  "contact_info": {
                    "type": "object",
                    "description": "Contact info",
                    "properties": {
                      "phone_country": {
                        "type": "string",
                        "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                      },
                      "phone_number": {
                        "type": "string",
                        "description": "Primary contact phone number associated with the account, in international format (E.164). Used for account notifications, verification, and security-related communications.",
                        "example": "+1-2025550123",
                        "max_length": 20
                      },
                      "email_address": {
                        "type": "string",
                        "description": "Email address",
                        "max_length": 256
                      },
                      "id_address": {
                        "type": "object",
                        "description": "Residential address of the account holder as recorded for identity verification and regulatory purposes.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Country of residence. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of residence. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>",
                            "example": "USA-CA"
                          },
                          "city": {
                            "type": "string",
                            "description": "City of residence.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the residence of the account holder.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the residential address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the residential address.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment number of the residential address.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the residential address.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      },
                      "mail_address": {
                        "type": "object",
                        "description": "Mailing address used for receiving account-related correspondence, if different from the residential address.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Mailing address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of the mailing address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                          },
                          "city": {
                            "type": "string",
                            "description": "City of the mailing address.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the mailing address.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the mailing address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the mailing address.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment number of the mailing address.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the mailing address.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    },
                    "required": [
                      "phone_country",
                      "phone_number",
                      "email_address",
                      "id_address",
                      "mail_address"
                    ]
                  },
                  "trusted_contact_info": {
                    "type": "object",
                    "description": "Trusted contact info",
                    "properties": {
                      "first_name": {
                        "type": "string",
                        "description": "First name of the trusted contact person."
                      },
                      "middle_name": {
                        "type": "string",
                        "description": "Middle name of the trusted contact person, if applicable."
                      },
                      "last_name": {
                        "type": "string",
                        "description": "Last name of the trusted contact person."
                      },
                      "phone_country": {
                        "type": "string",
                        "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                      },
                      "phone_number": {
                        "type": "string",
                        "description": "Primary contact phone number of the trusted contact person.",
                        "example": "+1-2025550124",
                        "max_length": 20
                      },
                      "relationship": {
                        "type": "string",
                        "description": "Relationship between the account holder and the trusted contact person. ",
                        "max_length": 32
                      },
                      "trusted_address": {
                        "type": "object",
                        "description": "Residential or mailing address of the trusted contact person.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Country of the trusted contact address. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of the address of the trusted contact person. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA)."
                          },
                          "city": {
                            "type": "string",
                            "description": "City of the trusted contact address.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the trusted contact person.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the trusted contact address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the trusted contact address.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment or unit number of the trusted contact address.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the trusted contact address.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    },
                    "required": [
                      "first_name",
                      "last_name",
                      "phone_country",
                      "phone_number",
                      "relationship",
                      "trusted_address"
                    ]
                  },
                  "employment_info": {
                    "type": "object",
                    "description": "Employment info",
                    "properties": {
                      "employment_type": {
                        "type": "string",
                        "description": "Type of employment. Value must be obtained from the Master Data API with data_type = EMPLOYMENT_TYPE."
                      },
                      "company_name": {
                        "type": "string",
                        "description": "Name of the applicant current employer.",
                        "max_length": 128
                      },
                      "occupation": {
                        "type": "string",
                        "description": "Specific occupation under the employment type. Value must be obtained from the Master Data API with data_type = OCCUPATION."
                      },
                      "position": {
                        "type": "string",
                        "description": "Job position or title. Value must be obtained from the Master Data API with data_type = POSITION."
                      },
                      "years_employed": {
                        "type": "number",
                        "description": "Total duration (in years) the applicant has been employed with the current employer. Decimal values are allowed to represent partial years (e.g. 1.5).",
                        "example": "1.5"
                      },
                      "employment_address": {
                        "type": "object",
                        "description": "Address of the applicant current employer.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Employer address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "Employer address state. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                          },
                          "city": {
                            "type": "string",
                            "description": "Employer address city.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the employer.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal code of the employer address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the employer address.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment number of the employer address.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the employer address.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    },
                    "required": [
                      "employment_type",
                      "company_name",
                      "occupation",
                      "position",
                      "years_employed",
                      "employment_address"
                    ]
                  },
                  "financial_info": {
                    "type": "object",
                    "description": "Financial information of the account holder.",
                    "properties": {
                      "liquidity_needs": {
                        "type": "string",
                        "description": "Indicates the client liquidity needs or preference for cash accessibility. Value must be obtained from the Master Data API with data_type = LIQUIDITY_NEEDS."
                      },
                      "annual_income": {
                        "type": "string",
                        "description": "Represents the client annual income. Value must be obtained from the Master Data API with data_type = ANNUAL_INCOME."
                      },
                      "total_net_worth": {
                        "type": "string",
                        "description": "Represents the client total net worth in USD. Value must be obtained from the Master Data API with data_type = TOTAL_NET_WORTH."
                      },
                      "liquid_net_worth": {
                        "type": "string",
                        "description": "Represents the client liquid net worth (USD), including cash and easily marketable assets. Value must be obtained from the Master Data API with data_type = LIQUID_NET_WORTH."
                      }
                    },
                    "required": [
                      "liquidity_needs",
                      "annual_income",
                      "total_net_worth",
                      "liquid_net_worth"
                    ]
                  },
                  "disclosure_info": {
                    "type": "object",
                    "description": "Disclosure",
                    "properties": {
                      "is_politically_exposed": {
                        "type": "boolean",
                        "description": "Indicates whether the account holder or any close associate or immediate family member is a Politically Exposed Person (PEP).  If true, the politically_exposed_info object must be provided with detailed disclosure information.  Please note: PEP accounts are currently not supported. Any request submitted with is_politically_exposed = true will be rejected and a failure response will be returned."
                      },
                      "politically_exposed_info": {
                        "type": "object",
                        "description": "Detail information about Politically Exposed Person (PEP) associated with the account holder.",
                        "properties": {
                          "pep_organization": {
                            "type": "string",
                            "description": "Name of the related political organization.",
                            "max_length": 256
                          },
                          "immediate_family": {
                            "type": "array",
                            "description": "Name of family members, including former spouses.",
                            "items": {
                              "type": "string",
                              "description": ""
                            }
                          }
                        },
                        "required": [
                          "pep_organization",
                          "immediate_family"
                        ]
                      },
                      "is_company_control_person": {
                        "type": "boolean",
                        "description": "Indicates whether the applicant is a control person of a public company. If true, `public_company_infos` must be provided with details of the companies controlled."
                      },
                      "public_company_infos": {
                        "type": "array",
                        "description": "List of public companies controlled by the applicant. Each item must include the company stock symbol.",
                        "items": {
                          "type": "object",
                          "description": "",
                          "properties": {
                            "symbol": {
                              "type": "string",
                              "description": "Stock symbol of the public company controlled by the applicant, as listed on the exchange."
                            }
                          },
                          "required": [
                            "symbol"
                          ]
                        }
                      }
                    },
                    "required": [
                      "is_politically_exposed",
                      "is_company_control_person"
                    ]
                  },
                  "tax_info": {
                    "type": "object",
                    "description": "Tax info",
                    "properties": {
                      "tax_type": {
                        "type": "string",
                        "description": "Tax type. Value must be obtained from the Master Data API with data_type = TAX_TYPE."
                      },
                      "tax_id": {
                        "type": "string",
                        "description": "Tax id. SSN format is (123-45-6789).",
                        "example": "123-45-6789"
                      }
                    },
                    "required": [
                      "tax_type",
                      "tax_id"
                    ]
                  },
                  "w8ben_info": {
                    "type": "object",
                    "description": "W8ben info",
                    "properties": {
                      "foreign_tax_id": {
                        "type": "string",
                        "description": "Foreign tax ID (FTIN) issued by non-U.S. authority (W-8BEN Line 6)"
                      },
                      "treaty_country": {
                        "type": "string",
                        "description": "Tax treaty country for reduced withholding. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                      },
                      "sign_date": {
                        "type": "string",
                        "description": "Certification date. Using UTC time.",
                        "format": "yyyy-MM-dd"
                      }
                    },
                    "required": [
                      "sign_date"
                    ]
                  },
                  "document_attachments": {
                    "type": "array",
                    "description": "Document attachments",
                    "min_items": 1,
                    "max_items": 10,
                    "items": {
                      "type": "object",
                      "description": "Item",
                      "properties": {
                        "doc_type": {
                          "type": "string",
                          "description": "Doc type"
                        },
                        "document_id": {
                          "type": "string",
                          "description": "Document id. Retrieved from the Upload Document API response."
                        }
                      },
                      "required": [
                        "doc_type",
                        "document_id"
                      ]
                    }
                  }
                },
                "required": [
                  "customer_info",
                  "contact_info",
                  "trusted_contact_info",
                  "employment_info",
                  "financial_info",
                  "disclosure_info",
                  "tax_info"
                ]
              }
            },
            "BROKERAGE_ADDITIONAL_FORM": {
              "description": "BROKERAGE_ADDITIONAL_FORM",
              "value": {
                "type": "object",
                "description": "The Brokerage Additional Form is used to collect supplemental information required to enable and configure brokerage trading services during the account opening process. This form captures brokerage-specific settings, investment profile details, and agreement acknowledgements necessary for regulatory compliance, risk assessment, and account configuration. It cannot be submitted independently and must be submitted together with the New Account Basic Form as part of the brokerage account opening workflow.",
                "form_code": "BROKERAGE_ADDITIONAL_FORM",
                "properties": {
                  "rep_code": {
                    "type": "string",
                    "description": "Specific code that Webull assigned to the Introducing broker, to distinguish between different IB customers, e.g. REP_DEMO."
                  },
                  "account_type": {
                    "type": "string",
                    "description": "Account type",
                    "example": "CASH"
                  },
                  "commission_code": {
                    "type": "string",
                    "description": "Commission Code for account, determining the commission structure"
                  },
                  "investment_knowledge": {
                    "type": "string",
                    "description": "Investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE."
                  },
                  "investment_objective": {
                    "type": "string",
                    "description": "Investment objective. Value must be obtained from the Master Data API with data_type = INVESTMENT_OBJECTIVE."
                  },
                  "time_horizon": {
                    "type": "string",
                    "description": "Time horizon. Value must be obtained from the Master Data API with data_type = TIME_HORIZON."
                  },
                  "is_exchange_or_finra_affiliated": {
                    "type": "boolean",
                    "description": "Indicates whether the applicant is employed by, or directly associated with, an exchange or a FINRA member firm. If true, exchange_affiliations details must be provided."
                  },
                  "exchange_affiliations": {
                    "type": "array",
                    "description": "List of the applicant affiliations with exchanges or FINRA member firms. Required if is_exchange_or_finra_affiliated is true.",
                    "items": {
                      "type": "object",
                      "description": "",
                      "properties": {
                        "company_name": {
                          "type": "string",
                          "description": "Name of the exchange or FINRA member firm the applicant is affiliated with.",
                          "max_length": 128
                        },
                        "country": {
                          "type": "string",
                          "description": "Country of the affiliated exchange or FINRA firm. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                        },
                        "state": {
                          "type": "string",
                          "description": "State or province of the affiliated exchange or FINRA firm. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                        },
                        "city": {
                          "type": "string",
                          "description": "City of the affiliated exchange or FINRA firm.",
                          "max_length": 256
                        },
                        "street_address": {
                          "type": "string",
                          "description": "Street address of the affiliated exchange or FINRA firm.",
                          "max_length": 512
                        },
                        "postal_code": {
                          "type": "string",
                          "description": "Postal or ZIP code of the affiliated exchange or FINRA firm.",
                          "max_length": 80
                        }
                      },
                      "required": [
                        "company_name",
                        "country",
                        "state",
                        "city",
                        "street_address",
                        "postal_code"
                      ]
                    }
                  },
                  "agreement_infos": {
                    "type": "object",
                    "description": "Agreement infos",
                    "properties": {
                      "accepted_at": {
                        "type": "string",
                        "description": "Accepted at. Agreement Confirmation Date (ET). If an agreement is signed, it must be filled in.",
                        "format": "yyyy-MM-dd"
                      },
                      "webull_financial_customer_agreement": {
                        "type": "object",
                        "description": "Webull Financial Customer Agreement (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_financial_annual_disclosure_statement": {
                        "type": "object",
                        "description": "Webull Financial Annual Disclosure Statement (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "statement_of_financial_condition": {
                        "type": "object",
                        "description": "Statement of Financial Condition",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "finra_customer_identification_notice": {
                        "type": "object",
                        "description": "FINRA - Customer Identification Notice (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "investor_education_and_protection": {
                        "type": "object",
                        "description": "Investor Education and Protection (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "electronic_delivery_of_trade_and_e_signature": {
                        "type": "object",
                        "description": "Electronic Delivery of Trade and E-Signature (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_financial_sipc_and_excess_sipc_information": {
                        "type": "object",
                        "description": "Webull Financial SIPC and Excess SIPC Information",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_financial_customer_relationship_summary": {
                        "type": "object",
                        "description": "Webull Financial Customer Relationship Summary",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "stop_and_advanced_orders_disclosure": {
                        "type": "object",
                        "description": "Stop & Advanced Orders Disclosure (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "order_handling_and_sec_rule_606_and_607_disclosures": {
                        "type": "object",
                        "description": "Order Handling and SEC Rule 606 and 607 Disclosures",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_business_continuity_plan_statement": {
                        "type": "object",
                        "description": "Webull Business Continuity Plan Statement (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Is accepted"
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_financial_llc_privacy_notice": {
                        "type": "object",
                        "description": "Webull Financial LLC Privacy Notice (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "day_trading_risk_disclosure": {
                        "type": "object",
                        "description": "Day Trading Risk Disclosure (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "extended_hours_trading_disclosure": {
                        "type": "object",
                        "description": "Extended Hours Trading Disclosure (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "low_priced_securities_alert": {
                        "type": "object",
                        "description": "Low-Priced Securities Alert (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_financial_customer_agreement_margin": {
                        "type": "object",
                        "description": "Webull Financial Customer Agreement - Margin(Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "webull_financial_margin_disclosure": {
                        "type": "object",
                        "description": "Webull Financial Margin Disclosure (Omnibus)",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      }
                    }
                  }
                },
                "required": [
                  "rep_code",
                  "account_type",
                  "commission_code",
                  "investment_knowledge",
                  "investment_objective",
                  "time_horizon",
                  "is_exchange_or_finra_affiliated",
                  "agreement_infos"
                ]
              }
            },
            "EVENT_CONTRACT_ADDITIONAL_FORM": {
              "description": "EVENT_CONTRACT_ADDITIONAL_FORM",
              "value": {
                "type": "object",
                "description": "The Event Contract Additional Form collects supplemental information required for enabling event contract trading.  This form cannot be submitted independently and must be submitted together with the New Account Basic Form as part of the account opening process.",
                "form_code": "EVENT_CONTRACT_ADDITIONAL_FORM",
                "properties": {
                  "rep_code": {
                    "type": "string",
                    "description": "Specific code that Webull assigned to the Introducing broker, to distinguish between different IB customers, e.g. REP_DEMO."
                  },
                  "commission_code": {
                    "type": "string",
                    "description": "Commission Code for account, determining the commission structure"
                  },
                  "investment_experience": {
                    "type": "string",
                    "description": "Event contract investment experience. Value must be obtained from the Master Data API with data_type = INVESTMENT_EXPERIENCE."
                  },
                  "investment_knowledge": {
                    "type": "string",
                    "description": "Event contract investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE_V2."
                  },
                  "trade_per_year": {
                    "type": "string",
                    "description": "Event contract trade per year. Value must be obtained from the Master Data API with data_type = TRADE_PER_YEAR."
                  },
                  "is_cftc_or_nfa_registered": {
                    "type": "boolean",
                    "description": "Indicates whether the applicant is registered with the CFTC or is a member of the NFA."
                  },
                  "is_exchange_affiliated": {
                    "type": "boolean",
                    "description": "Indicates whether the applicant is employed by, or directly associated with, an exchange member. If true, exchange_affiliations details must be provided."
                  },
                  "is_bp_shared": {
                    "type": "boolean",
                    "description": "Whether to enable buying power sharing.  true: Enable buying power sharing. false: Do not enable buying power sharing. By default, it will share with the CASH account; if no CASH account exists, it will share with the MARGIN account.  "
                  },
                  "exchange_affiliations": {
                    "type": "array",
                    "description": "List of the applicant affiliations with exchanges member. Required if is_exchange_affiliated is true.",
                    "items": {
                      "type": "object",
                      "description": "",
                      "properties": {
                        "company_name": {
                          "type": "string",
                          "description": "Name of the exchange member the applicant is affiliated with.",
                          "max_length": 128
                        },
                        "country": {
                          "type": "string",
                          "description": "Country of the affiliated exchange. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                        },
                        "state": {
                          "type": "string",
                          "description": "State or province of the affiliated exchange. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                        },
                        "city": {
                          "type": "string",
                          "description": "City of the affiliated exchange.",
                          "max_length": 256
                        },
                        "street_address": {
                          "type": "string",
                          "description": "Street address of the affiliated exchange.",
                          "max_length": 512
                        },
                        "postal_code": {
                          "type": "string",
                          "description": "Postal or ZIP code of the affiliated exchange.",
                          "max_length": 80
                        }
                      },
                      "required": [
                        "company_name",
                        "country",
                        "state",
                        "city",
                        "street_address",
                        "postal_code"
                      ]
                    }
                  },
                  "agreement_infos": {
                    "type": "object",
                    "description": "Agreement infos",
                    "properties": {
                      "accepted_at": {
                        "type": "string",
                        "description": "Accepted at. Agreement Confirmation Date (ET). If an agreement is signed, it must be filled in.",
                        "format": "yyyy-MM-dd"
                      },
                      "cleared_swap_arbitration_agreement": {
                        "type": "object",
                        "description": "Cleared Swap Arbitration Agreement",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "event_contract_heightened_risk_disclosure": {
                        "type": "object",
                        "description": "Event Contract Heightened Risk Disclosure",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "event_contract_sports_disclaimer": {
                        "type": "object",
                        "description": "Event Contract Sports Disclaimer",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "form_1_55k_firm_specific_disclosure_document": {
                        "type": "object",
                        "description": "1.55(k) Firm Specific Disclosure Document",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "event_contract_customer_agreement": {
                        "type": "object",
                        "description": "Event Contract Customer Agreement",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "event_contract_risk_disclosure": {
                        "type": "object",
                        "description": "Event Contract Risk Disclosure",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      },
                      "investment_and_trading_disclosures_booklet_futures": {
                        "type": "object",
                        "description": "Investment and Trading Disclosures Booklet - Futures",
                        "properties": {
                          "agreement_id": {
                            "type": "string",
                            "description": "Agreement id"
                          },
                          "sign_method": {
                            "type": "string",
                            "description": "Signature method. Value must be obtained from the Master Data API with data_type = AGREEMENT_SIGN_METHOD."
                          },
                          "is_accepted": {
                            "type": "boolean",
                            "description": "Whether the user accepted the agreement."
                          }
                        },
                        "required": [
                          "agreement_id",
                          "is_accepted"
                        ]
                      }
                    }
                  }
                },
                "required": [
                  "rep_code",
                  "commission_code",
                  "investment_experience",
                  "investment_knowledge",
                  "trade_per_year",
                  "is_cftc_or_nfa_registered",
                  "is_exchange_affiliated",
                  "agreement_infos"
                ]
              }
            },
            "APPLICATION_RESUBMISSION_FORM": {
              "description": "APPLICATION_RESUBMISSION_FORM",
              "value": {
                "type": "object",
                "description": "This form is used to collect additional information or supporting documents required during the account opening process. It is triggered by KYC/AML screening, Enhanced Due Diligence (EDD), or compliance review, and is conditionally required to complete the account application.",
                "form_code": "APPLICATION_RESUBMISSION_FORM",
                "properties": {
                  "codes": {
                    "type": "array",
                    "description": "List of resubmission codes and the corresponding fields that need to be provided:  -- OCR_SUSPENDED: Provide picture_front_doc_id and picture_back_doc_id for identity document images.   -- CIP_SUSPENDED, UNMATCHED_SSN_PICTURE: Provide tax_id_doc_id for tax/SSN verification.  -- NOT_A_GOVERNMENT_ISSUED, EXPIRED_ID, ID_CARD_ILLEGIBLE, ID_ON_SCREEN, UNMATCHED_ID_NUMBER_WITH_THE_ID_PICTURE: Provide id_doc_type, id_expiry_date, picture_front_doc_id, picture_back_doc_id.  -- UNMATCHED_ID_EXPIRE_DATE: Provide id_expiry_date.   -- UNMATCHED_NATIONALITY_WITH_THE_ID_PICTURE: Provide country_of_citizenship.  --  UNMATCHED_NAME_WITH_THE_ID_PICTURE: Provide first_name, middle_name, last_name.  -- UNMATCHED_DOB_WITH_THE_ID_PICTURE: Provide date_of_birth.  -- SELFIE_NOT_MEET_THE_REQUIREMENTS: Provide selfie_img_key for selfie verification. -- NEED_407_LETTERS: Provide rule_3210_doc_id when is_exchange_or_finra_affiliated is true.",
                    "items": {
                      "type": "string",
                      "description": "Resubmission code indicating which type of information is required."
                    }
                  },
                  "id_doc_type": {
                    "type": "string",
                    "description": "Type of identification document provided by the applicant. The value must be retrieved from the Master Data API where data_type = ID_DOC_TYPE."
                  },
                  "id_number": {
                    "type": "string",
                    "description": "Identification document number exactly as shown on the provided ID."
                  },
                  "id_expiry_date": {
                    "type": "string",
                    "description": "Expiration date of the identification document."
                  },
                  "picture_front_doc_id": {
                    "type": "string",
                    "description": "Document ID referencing the uploaded image of the front side of the identification document."
                  },
                  "picture_back_doc_id": {
                    "type": "string",
                    "description": "Document ID referencing the uploaded image of the back side of the identification document."
                  },
                  "first_name": {
                    "type": "string",
                    "description": "First name of the applicant as it appears on the identification document."
                  },
                  "middle_name": {
                    "type": "string",
                    "description": "Middle name of the applicant as it appears on the identification document, if applicable."
                  },
                  "last_name": {
                    "type": "string",
                    "description": "Last name of the applicant as it appears on the identification document."
                  },
                  "date_of_birth": {
                    "type": "string",
                    "description": "Date of birth of the applicant as shown on the identification document."
                  },
                  "country_of_citizenship": {
                    "type": "string",
                    "description": "Country of citizenship of the account holder. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                  },
                  "selfie_img_key": {
                    "type": "string",
                    "description": "Document ID referencing the uploaded selfie image used for identity verification."
                  },
                  "tax_id_doc_id": {
                    "type": "string",
                    "description": "Document ID referencing the uploaded supporting document for tax identification or SSN verification."
                  },
                  "rule_3210_doc_id": {
                    "type": "string",
                    "description": "Document ID referencing the uploaded Rule 3210 Approval Form. Required when is_exchange_or_finra_affiliated is true."
                  }
                },
                "required": [
                  "codes"
                ]
              }
            },
            "ACCOUNT_UPDATE_FORM": {
              "description": "ACCOUNT_UPDATE_FORM",
              "value": {
                "type": "object",
                "description": "This form is used to submit post-onboarding account update requests for an existing account. It serves as a unified entry point for multiple account maintenance and compliance-related update scenarios after the account has been successfully opened. The specific update scenario is determined by the `type` field. Based on the selected type, the corresponding update information object must be provided. This form is commonly used for periodic reviews, regulatory compliance verification, risk control follow-ups, and other account data maintenance processes to ensure that account records remain accurate, valid, and up to date.",
                "form_code": "ACCOUNT_UPDATE_FORM",
                "properties": {
                  "type": {
                    "type": "string",
                    "description": "Specifies the category of account update to be performed after account opening. This field determines which account update scenario applies and which corresponding information object is required in the request. Each type represents a distinct account maintenance or compliance update operation."
                  },
                  "commission_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is COMMISSION_UPDATE.",
                    "properties": {
                      "commission_code": {
                        "type": "string",
                        "description": "Commission Code for account, determining the commission structure."
                      }
                    },
                    "required": [
                      "commission_code"
                    ]
                  },
                  "identity_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is IDENTITY_UPDATE. This form is used to collect updated identification information for an existing account. It is required during account maintenance, periodic review, or compliance verification to ensure the accuracy and validity of identity records.",
                    "properties": {
                      "id_doc_type": {
                        "type": "string",
                        "description": "Type of identification document provided by the applicant. The value must be retrieved from the Master Data API where data_type = ID_DOC_TYPE."
                      },
                      "id_number": {
                        "type": "string",
                        "description": "Identification document number exactly as shown on the provided ID."
                      },
                      "id_expiry_date": {
                        "type": "string",
                        "description": "Expiration date of the identification document.",
                        "format": "yyyy-MM-dd"
                      },
                      "picture_front_doc_id": {
                        "type": "string",
                        "description": "Document ID referencing the uploaded image of the front side of the identification document."
                      },
                      "picture_back_doc_id": {
                        "type": "string",
                        "description": "Document ID referencing the uploaded image of the back side of the identification document."
                      }
                    },
                    "required": [
                      "id_doc_type",
                      "id_number",
                      "id_expiry_date",
                      "picture_front_doc_id",
                      "picture_back_doc_id"
                    ]
                  },
                  "bp_share_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is SHARE_ACCOUNT_UPDATE. Buypower share info",
                    "properties": {
                      "share_direction": {
                        "type": "string",
                        "description": "Buypower share direction"
                      }
                    },
                    "required": [
                      "share_direction"
                    ]
                  },
                  "employment_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is EMPLOYMENT_UPDATE. Collects updated employment information required during account maintenance or compliance review.",
                    "properties": {
                      "employment_type": {
                        "type": "string",
                        "description": "Type of employment. Value must be obtained from the Master Data API with data_type = EMPLOYMENT_TYPE."
                      },
                      "company_name": {
                        "type": "string",
                        "description": "Name of the applicant current employer.",
                        "max_length": 128
                      },
                      "occupation": {
                        "type": "string",
                        "description": "Specific occupation under the employment type. Value must be obtained from the Master Data API with data_type = OCCUPATION."
                      },
                      "position": {
                        "type": "string",
                        "description": "Job position or title. Value must be obtained from the Master Data API with data_type = POSITION."
                      },
                      "years_employed": {
                        "type": "number",
                        "description": "Total duration (in years) the applicant has been employed with the current employer. Decimal values are allowed to represent partial years (e.g. 1.5)."
                      },
                      "employment_address": {
                        "type": "object",
                        "description": "Address of the applicant current employer.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Employer address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "Employer address state. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                          },
                          "city": {
                            "type": "string",
                            "description": "Employer address city.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the employer. Multiple lines are allowed.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal code of the employer address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the employer address, if applicable.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment number of the employer address, if applicable.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the employer address, if applicable.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    },
                    "required": [
                      "employment_type",
                      "company_name",
                      "occupation",
                      "position",
                      "years_employed",
                      "employment_address"
                    ]
                  },
                  "trusted_contact_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is TRUSTED_UPDATE. Collects or updates trusted contact information for an account, used for account protection, compliance review, or situations where the account holder cannot be contacted.",
                    "properties": {
                      "first_name": {
                        "type": "string",
                        "description": "First name of the trusted contact person."
                      },
                      "middle_name": {
                        "type": "string",
                        "description": "Middle name of the trusted contact person, if applicable."
                      },
                      "last_name": {
                        "type": "string",
                        "description": "Last name of the trusted contact person."
                      },
                      "phone_country": {
                        "type": "string",
                        "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                      },
                      "phone_number": {
                        "type": "string",
                        "description": "Primary contact phone number of the trusted contact person.",
                        "example": "+1-2025550124",
                        "max_length": 20
                      },
                      "relationship": {
                        "type": "string",
                        "description": "Relationship between the account holder and the trusted contact person.",
                        "max_length": 32
                      },
                      "trusted_address": {
                        "type": "object",
                        "description": "Residential or mailing address of the trusted contact person.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Country of the trusted contact person address. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of the trusted contact person address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                          },
                          "city": {
                            "type": "string",
                            "description": "City of the trusted contact person address.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the trusted contact person. Multiple address lines are allowed.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the trusted contact person address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the trusted contact person address, if applicable.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment or unit number of the trusted contact person address, if applicable.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the trusted contact person address, if applicable.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    },
                    "required": [
                      "first_name",
                      "last_name",
                      "phone_country",
                      "phone_number",
                      "relationship",
                      "trusted_address"
                    ]
                  },
                  "email_address_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is EMAIL_ADDRESS_UPDATE. Collects or updates the primary email address associated with an account. This information is used for account communication, security notifications, and compliance-related correspondence.",
                    "properties": {
                      "email_address": {
                        "type": "string",
                        "description": "Email address",
                        "max_length": 256
                      }
                    },
                    "required": [
                      "email_address"
                    ]
                  },
                  "address_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is ADDRESS_UPDATE. Collects or updates address information associated with an account, including residential and mailing addresses. This information is used for account administration, regulatory compliance, and official correspondence.",
                    "properties": {
                      "id_address": {
                        "type": "object",
                        "description": "Residential address of the account holder as recorded for identity verification and regulatory purposes.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Country of residence. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of residence. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                          },
                          "city": {
                            "type": "string",
                            "description": "City of residence.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the account holder residence. Multiple address lines are allowed.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the residential address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the id address, if applicable.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment number of the id address, if applicable.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the id address, if applicable.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      },
                      "mail_address": {
                        "type": "object",
                        "description": "Mailing address used for receiving account-related correspondence, if different from the residential address.",
                        "properties": {
                          "country": {
                            "type": "string",
                            "description": "Mailing address country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                          },
                          "state": {
                            "type": "string",
                            "description": "State or province of the mailing address. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code>"
                          },
                          "city": {
                            "type": "string",
                            "description": "City of the mailing address.",
                            "max_length": 256
                          },
                          "street_address": {
                            "type": "string",
                            "description": "Street address of the mailing address. Multiple address lines are allowed.",
                            "max_length": 512
                          },
                          "postal_code": {
                            "type": "string",
                            "description": "Postal or ZIP code of the mailing address.",
                            "max_length": 80
                          },
                          "neighborhood_code": {
                            "type": "string",
                            "description": "Neighborhood or district code of the mailing address.",
                            "max_length": 128
                          },
                          "apartment_number": {
                            "type": "string",
                            "description": "Apartment number of the mailing address.",
                            "max_length": 128
                          },
                          "building_number": {
                            "type": "string",
                            "description": "Building number of the mailing address.",
                            "max_length": 128
                          }
                        },
                        "required": [
                          "country",
                          "state",
                          "city",
                          "street_address",
                          "postal_code"
                        ]
                      }
                    },
                    "required": [
                      "id_address",
                      "mail_address"
                    ]
                  },
                  "phone_number_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is PHONE_NUMBER_UPDATE. Collects or updates the primary phone number associated with an account, used for account communication, security verification, and regulatory notifications.",
                    "properties": {
                      "phone_country": {
                        "type": "string",
                        "description": "Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                      },
                      "phone_number": {
                        "type": "string",
                        "description": "Primary contact phone number associated with the account, in international format (E.164). Used for account notifications, verification, and security-related communications.",
                        "example": "+1-2025550123",
                        "max_length": 20
                      }
                    },
                    "required": [
                      "phone_country",
                      "phone_number"
                    ]
                  },
                  "politically_exposed_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is PEP_UPDATE. Collects or updates Politically Exposed Person (PEP) information for an account. This information is used for compliance review, enhanced due diligence (EDD), and regulatory reporting.",
                    "properties": {
                      "is_politically_exposed": {
                        "type": "boolean",
                        "description": "Indicates whether the account holder or any close associate or immediate family member is a Politically Exposed Person (PEP).  If true, the politically_exposed_info object must be provided with detailed disclosure information.  Please note: PEP accounts are currently not supported. Any request submitted with is_politically_exposed = true will be rejected and a failure response will be returned."
                      },
                      "politically_exposed_info": {
                        "type": "object",
                        "description": "Detail information about Politically Exposed Person (PEP) associated with the account holder.",
                        "properties": {
                          "pep_organization": {
                            "type": "string",
                            "description": "Name of the related political organization.",
                            "max_length": 256
                          },
                          "immediate_family": {
                            "type": "array",
                            "description": "Name of family members, including former spouses.",
                            "items": {
                              "type": "string",
                              "description": ""
                            }
                          }
                        },
                        "required": [
                          "pep_organization",
                          "immediate_family"
                        ]
                      }
                    },
                    "required": [
                      "is_politically_exposed"
                    ]
                  },
                  "company_control_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is COMPANY_CONTROL_UPDATE. Collects or updates information about the applicant status as a control person of a public company. This information is used for compliance review, enhanced due diligence (EDD), and regulatory reporting.",
                    "properties": {
                      "is_company_control_person": {
                        "type": "boolean",
                        "description": "Indicates whether the applicant is a control person of a public company. If true, `public_company_infos` must be provided with details of the companies controlled."
                      },
                      "public_company_infos": {
                        "type": "array",
                        "description": "List of public companies controlled by the applicant. Each item must include the company stock symbol.",
                        "items": {
                          "type": "object",
                          "description": "",
                          "properties": {
                            "symbol": {
                              "type": "string",
                              "description": "Stock symbol of the public company controlled by the applicant, as listed on the exchange."
                            }
                          },
                          "required": [
                            "symbol"
                          ]
                        }
                      }
                    },
                    "required": [
                      "is_company_control_person"
                    ]
                  },
                  "disclosure_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is DISCLOSURE_UPDATE. Collects or updates the applicant regulatory and compliance-related affiliations or disclosures. This form is used for compliance review, enhanced due diligence (EDD), and regulatory reporting, and may include multiple types of information beyond exchange or FINRA affiliations.",
                    "properties": {
                      "brokerage_disclosure_info": {
                        "type": "object",
                        "description": "Brokerage disclosure info",
                        "properties": {
                          "is_exchange_or_finra_affiliated": {
                            "type": "boolean",
                            "description": "Indicates whether the applicant is employed by, or directly associated with, an exchange or a FINRA member firm. If true, exchange_affiliations details must be provided."
                          },
                          "exchange_affiliations": {
                            "type": "array",
                            "description": "List of the applicant affiliations with exchanges or FINRA member firms. Required if is_exchange_or_finra_affiliated is true.",
                            "items": {
                              "type": "object",
                              "description": "",
                              "properties": {
                                "company_name": {
                                  "type": "string",
                                  "description": "Name of the exchange or FINRA member firm the applicant is affiliated with.",
                                  "max_length": 128
                                },
                                "country": {
                                  "type": "string",
                                  "description": "Country of the affiliated exchange or FINRA firm. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                                },
                                "state": {
                                  "type": "string",
                                  "description": "State or province of the affiliated exchange or FINRA firm. Value must be obtained from the Master Data API with data_type = STATE."
                                },
                                "city": {
                                  "type": "string",
                                  "description": "City of the affiliated exchange or FINRA firm.",
                                  "max_length": 256
                                },
                                "street_address": {
                                  "type": "string",
                                  "description": "Street address of the affiliated exchange or FINRA firm. Multiple lines are allowed.",
                                  "max_length": 512
                                },
                                "postal_code": {
                                  "type": "string",
                                  "description": "Postal or ZIP code of the affiliated exchange or FINRA firm.",
                                  "max_length": 80
                                }
                              },
                              "required": [
                                "company_name",
                                "country",
                                "state",
                                "city",
                                "street_address",
                                "postal_code"
                              ]
                            }
                          }
                        },
                        "required": [
                          "is_exchange_or_finra_affiliated"
                        ]
                      }
                    },
                    "required": [
                      "brokerage_disclosure_info"
                    ]
                  },
                  "dependent_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is DEPENDENT_UPDATE. Collects or updates the number of dependents associated with an account holder. This information is used for account maintenance, compliance review, and regulatory reporting purposes.",
                    "properties": {
                      "num_dependents": {
                        "type": "number",
                        "description": "Number of dependents financially or legally associated with the account holder."
                      }
                    },
                    "required": [
                      "num_dependents"
                    ]
                  },
                  "citizenship_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is CITIZENSHIP_UPDATE. Collects or updates the country of citizenship for an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
                    "properties": {
                      "country_of_citizenship": {
                        "type": "string",
                        "description": "Country of citizenship of the account holder. The value must be retrieved from the Master Data API with data_type = COUNTRY."
                      }
                    },
                    "required": [
                      "country_of_citizenship"
                    ]
                  },
                  "gender_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is GENDER_UPDATE. Collects or updates the gender of an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
                    "properties": {
                      "gender": {
                        "type": "string",
                        "description": "Gender of the account holder. Allowed values: F (Female), M (Male)."
                      }
                    },
                    "required": [
                      "gender"
                    ]
                  },
                  "marital_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is MARITAL_UPDATE. Collects or updates the marital status of an account holder. This information is used for account maintenance, compliance review, and regulatory reporting.",
                    "properties": {
                      "marital_status": {
                        "type": "string",
                        "description": "Marital status of the account holder. Value must be obtained from the Master Data API with data_type = MARITAL_STATUS."
                      }
                    },
                    "required": [
                      "marital_status"
                    ]
                  },
                  "financial_investment_info": {
                    "type": "object",
                    "description": "This object is applicable when the update type is FINANCIAL_INVESTMENT_UPDATE. Collects or updates the investment profile and financial information of an account holder. This information is used for account maintenance, risk assessment, and regulatory compliance.",
                    "properties": {
                      "financial_info": {
                        "type": "object",
                        "description": "Financial information of the account holder.",
                        "properties": {
                          "liquidity_needs": {
                            "type": "string",
                            "description": "Indicates the client liquidity needs or preference for cash accessibility. Value must be obtained from the Master Data API with data_type = LIQUIDITY_NEEDS."
                          },
                          "annual_income": {
                            "type": "string",
                            "description": "Represents the client annual income. Value must be obtained from the Master Data API with data_type = ANNUAL_INCOME."
                          },
                          "total_net_worth": {
                            "type": "string",
                            "description": "Represents the client total net worth in USD. Value must be obtained from the Master Data API with data_type = TOTAL_NET_WORTH."
                          },
                          "liquid_net_worth": {
                            "type": "string",
                            "description": "Represents the client liquid net worth (USD), including cash and easily marketable assets. Value must be obtained from the Master Data API with data_type = LIQUID_NET_WORTH."
                          }
                        },
                        "required": [
                          "liquidity_needs",
                          "annual_income",
                          "total_net_worth",
                          "liquid_net_worth"
                        ]
                      },
                      "event_contract_investment_info": {
                        "type": "object",
                        "description": "Event contract account investment info",
                        "properties": {
                          "investment_experience": {
                            "type": "string",
                            "description": "Event contract investment experience. Value must be obtained from the Master Data API with data_type = INVESTMENT_EXPERIENCE."
                          },
                          "investment_knowledge": {
                            "type": "string",
                            "description": "Event contract investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE_V2."
                          },
                          "trade_per_year": {
                            "type": "string",
                            "description": "Event contract trade per year. Value must be obtained from the Master Data API with data_type = TRADE_PER_YEAR."
                          }
                        },
                        "required": [
                          "investment_experience",
                          "investment_knowledge",
                          "trade_per_year"
                        ]
                      },
                      "brokerage_investment_info": {
                        "type": "object",
                        "description": "Brokerage account investment info",
                        "properties": {
                          "investment_knowledge": {
                            "type": "string",
                            "description": "Investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE."
                          },
                          "investment_objective": {
                            "type": "string",
                            "description": "Investment objective. Value must be obtained from the Master Data API with data_type = INVESTMENT_OBJECTIVE."
                          },
                          "time_horizon": {
                            "type": "string",
                            "description": "Time horizon. Value must be obtained from the Master Data API with data_type = TIME_HORIZON."
                          }
                        },
                        "required": [
                          "investment_knowledge",
                          "investment_objective",
                          "time_horizon"
                        ]
                      }
                    }
                  }
                },
                "required": [
                  "type"
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
    "name": "Get Form Detail",
    "description": {
      "content": "Retrieves the JSON schema for the specified form code and version.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "forms",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique identifier of the form, typically used to distinguish different form definitions.",
            "type": "text/plain"
          },
          "key": "form_code",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Version number of the form definition",
            "type": "text/plain"
          },
          "key": "version",
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

## Upload Document

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/document-upload.md>

### Upload Document

Uploads a document file. Supports multipart/form-data requests.Supported file types: image/jpeg, image/png, application/pdf. Maximum file size: 5MB.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/documents/upload",
  "method": "post",
  "tags": [
    "Document"
  ],
  "description": "Uploads a document file. Supports multipart/form-data requests.Supported file types: image/jpeg, image/png, application/pdf. Maximum file size: 5MB.",
  "operationId": "documentUpload",
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
    "description": "Multipart request containing the document file and optional metadata.",
    "content": {
      "multipart/form-data": {
        "schema": {
          "type": "object",
          "properties": {
            "doc_type": {
              "type": "string",
              "description": "Enumeration of supported document types used for identity verification, profile updates, and compliance checks.",
              "example": "FACE_IMG",
              "enum": [
                "FACE_IMG",
                "ID_DOCUMENT",
                "ADDRESS_PROOF",
                "OTHER"
              ]
            },
            "file": {
              "type": "string",
              "format": "binary"
            }
          },
          "title": "FileUploadRequest"
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
              "document_id": {
                "type": "string",
                "description": "Unique identifier of the uploaded document.",
                "example": "dadd1f4d5ed04c25a095830d8576cbe2"
              }
            },
            "description": "Result of a file upload operation",
            "title": "FileUploadResult"
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
    "name": "Upload Document",
    "description": {
      "content": "Uploads a document file. Supports multipart/form-data requests.Supported file types: image/jpeg, image/png, application/pdf. Maximum file size: 5MB.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "documents",
        "upload"
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
        "value": "multipart/form-data"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "formdata",
      "formdata": []
    }
  }
}
```

## Download Document

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/document-download.md>

### Download Document

Downloads a document by document_id.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/documents/download",
  "method": "get",
  "tags": [
    "Document"
  ],
  "description": "Downloads a document by document_id.",
  "operationId": "documentDownload",
  "parameters": [
    {
      "name": "document_id",
      "in": "query",
      "description": "Unique identifier of the uploaded document.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "dadd1f4d5ed04c25a095830d8576cbe2"
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
              "doc_name": {
                "type": "string",
                "description": "Original name of the document file.",
                "example": "passport_front.jpg"
              },
              "doc_type": {
                "type": "string",
                "description": "Enumeration of supported document types used for identity verification, profile updates, and compliance checks.",
                "example": "ID_FRONT_IMG",
                "enum": [
                  "FACE_IMG",
                  "ID_DOCUMENT",
                  "ADDRESS_PROOF",
                  "OTHER"
                ]
              },
              "mime_type": {
                "type": "string",
                "description": "Enumeration of supported file MIME types for document upload and download.",
                "example": "image/jpeg",
                "enum": [
                  "image/jpeg",
                  "image/png",
                  "application/pdf"
                ]
              },
              "doc_content_base64": {
                "type": "string",
                "description": "Base64-encoded content of the document file.",
                "example": "iVBORw0KGgoAAAANSUhEUgAA..."
              }
            },
            "description": "Response object containing document metadata and Base64-encoded content for document download.",
            "title": "DocumentDownloadResp"
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
    "name": "Download Document",
    "description": {
      "content": "Downloads a document by document_id.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "documents",
        "download"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique identifier of the uploaded document.",
            "type": "text/plain"
          },
          "key": "document_id",
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

## Assets Summary

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/summary.md>

### Get Account Summary

Retrieves summary-level asset information for a specified account, including available cash, position market value, and total asset value.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/assets/summaries/get",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves summary-level asset information for a specified account, including available cash, position market value, and total asset value.",
  "operationId": "summary",
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
              "balance",
              "positions"
            ],
            "type": "object",
            "properties": {
              "balance": {
                "required": [
                  "account_currency_assets",
                  "total_asset_currency",
                  "total_cash_balance"
                ],
                "type": "object",
                "properties": {
                  "total_asset_currency": {
                    "type": "string",
                    "description": "Currency",
                    "example": "USD",
                    "enum": [
                      "USD"
                    ]
                  },
                  "total_cash_balance": {
                    "type": "string",
                    "description": "Cash Balance, ",
                    "example": "485705.0"
                  },
                  "total_market_value": {
                    "type": "string",
                    "description": "Total holding market value",
                    "example": "995705.0"
                  },
                  "total_unrealized_profit_loss": {
                    "type": "string",
                    "description": "The unrealized profit or loss of open positions for the specified account.",
                    "example": "227689.0"
                  },
                  "total_net_liquidation_value": {
                    "type": "string",
                    "description": "Net Account Value, ",
                    "example": "727687.04"
                  },
                  "total_day_profit_loss": {
                    "type": "string",
                    "description": "Day's P&L<br/> Not returned at the Omnibus master account level.",
                    "example": "11798.6"
                  },
                  "maintenance_margin": {
                    "type": "string",
                    "description": "Maintenance Margin（Margin Account）<br/> Only returned for margin accounts.",
                    "example": "0.0"
                  },
                  "open_margin_calls": {
                    "type": "string",
                    "description": "Open Margin Calls:<br/> Only returned for margin accounts.<br/>RM: An RM call is triggered when an account's margin equity falls below the maintenance requirement, often due to a decline in the value of the account's positions or an increase in margin requirements for the holdings.\n<br/>RT: A Reg T call occurs when there is not enough equity in an account to cover the 50% initial margin requirement. This can happen due to open positions depreciating or from holding positions overnight that were opened with day trade buying power.\n<br/>EM: An EM call is triggered when a flagged Pattern Day Trader (PDT) account closes the prior business day below the $25,000 minimum Net Account Value (NAV) requirement. Only positions held in the margin account will count toward this requirement. Crypto, futures, event contracts, and any assets held outside the margin account are excluded from the calculation.\n<br/>DT: A Day Trade (DT) Call is issued when you exceed your available Day Trading Buying Power (DTBP) and then place a day trade. A DT call may also occur if you execute a day trade while an Equity Maintenance (EM) call is still active on your account. Exceeding day trade buying power most commonly results from trading on intraday profits. For more details on how buying power is replenished, please refer to our Margin Buying Power page.",
                    "example": "['EM']",
                    "enum": [
                      "EM",
                      "RM",
                      "RT",
                      "DT"
                    ]
                  },
                  "account_currency_assets": {
                    "type": "array",
                    "description": "Currency assets Details",
                    "items": {
                      "required": [
                        "cash_balance",
                        "currency"
                      ],
                      "type": "object",
                      "properties": {
                        "currency": {
                          "type": "string",
                          "description": "Currency",
                          "example": "USD",
                          "enum": [
                            "USD"
                          ]
                        },
                        "cash_balance": {
                          "type": "string",
                          "description": "Cash Balance.",
                          "example": "485705.95"
                        },
                        "settled_cash": {
                          "type": "string",
                          "description": "Settled Cash，<br/> Not returned for event contract accounts.<br/> Not returned for margin accounts.",
                          "example": "485705.95"
                        },
                        "unsettled_cash": {
                          "type": "string",
                          "description": "Unsettled Cash，<br/> Not returned for event contract accounts.<br/> Not returned for margin accounts.",
                          "example": "0.0"
                        },
                        "market_value": {
                          "type": "string",
                          "description": "holding market value",
                          "example": "0.0"
                        },
                        "held_amount": {
                          "type": "string",
                          "description": "In-transit funds<br/> Not returned for event contract accounts.",
                          "example": "0.0"
                        },
                        "buying_power": {
                          "type": "string",
                          "description": "Buying Power<br/> Not returned for margin accounts.",
                          "example": "484551"
                        },
                        "day_buying_power": {
                          "type": "string",
                          "description": "Day-Trade Buying Power（Margin Account）<br/> Only returned for margin accounts.",
                          "example": "0.0"
                        },
                        "overnight_buying_power": {
                          "type": "string",
                          "description": "Overnight BP（Margin Account）<br/> Only returned for margin accounts.",
                          "example": "0.0"
                        },
                        "night_trading_buying_power": {
                          "type": "string",
                          "description": "Night Trading Buying Power<br/> Not returned for event contract accounts.",
                          "example": "0.0"
                        },
                        "day_profit_loss": {
                          "type": "string",
                          "description": "Day's P&L<br/> Not returned at the Omnibus master account level.",
                          "example": "0.0"
                        },
                        "unrealized_profit_loss": {
                          "type": "string",
                          "description": "Open P&L",
                          "example": "227689"
                        },
                        "available_withdrawal": {
                          "type": "string",
                          "description": "The amount of funds currently available for withdrawal.",
                          "example": "3.0558743194E8"
                        },
                        "interests_unpaid": {
                          "type": "string",
                          "description": "Interest to be paid<br/> Not returned for event contract accounts.",
                          "example": "0.0"
                        },
                        "net_liquidation_value": {
                          "type": "string",
                          "description": "Net Account Value",
                          "example": "0.0"
                        }
                      },
                      "description": "Currency assets Details",
                      "title": "AssetsCurrencyAssets"
                    }
                  }
                },
                "description": "Account Balance",
                "title": "AssetsBalanceResult"
              },
              "positions": {
                "type": "array",
                "description": "Account Position",
                "items": {
                  "required": [
                    "cost_price",
                    "currency",
                    "instrument_type",
                    "last_price",
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
                        "USD"
                      ]
                    },
                    "quantity": {
                      "type": "string",
                      "description": "Quantity of the position",
                      "example": "1"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market.",
                      "example": "AAPL"
                    },
                    "instrument_type": {
                      "type": "string",
                      "description": "Type of financial instrument associated with the request.",
                      "example": "EQUITY",
                      "enum": [
                        "EQUITY",
                        "EVENT"
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
                    "event_outcome": {
                      "type": "string",
                      "description": "Event outcome decision, only applicable to event orders.",
                      "example": "yes",
                      "enum": [
                        "yes",
                        "no"
                      ]
                    }
                  },
                  "description": "Account Position",
                  "title": "AssetsPositionResult"
                }
              }
            },
            "title": "AssetsSummaryResult"
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
    "name": "Get Account Summary",
    "description": {
      "content": "Retrieves summary-level asset information for a specified account, including available cash, position market value, and total asset value.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "assets",
        "summaries",
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

## Assets Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/account-balance.md>

### Get Account Balance

Retrieves the balance information for a specified account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/assets/balances/get",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves the balance information for a specified account.",
  "operationId": "accountBalance",
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
              "account_currency_assets",
              "total_asset_currency",
              "total_cash_balance"
            ],
            "type": "object",
            "properties": {
              "total_asset_currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "total_cash_balance": {
                "type": "string",
                "description": "Cash Balance, ",
                "example": "485705.0"
              },
              "total_market_value": {
                "type": "string",
                "description": "Total holding market value",
                "example": "995705.0"
              },
              "total_unrealized_profit_loss": {
                "type": "string",
                "description": "The unrealized profit or loss of open positions for the specified account.",
                "example": "227689.0"
              },
              "total_net_liquidation_value": {
                "type": "string",
                "description": "Net Account Value, ",
                "example": "727687.04"
              },
              "total_day_profit_loss": {
                "type": "string",
                "description": "Day's P&L<br/> Not returned at the Omnibus master account level.",
                "example": "11798.6"
              },
              "maintenance_margin": {
                "type": "string",
                "description": "Maintenance Margin（Margin Account）<br/> Only returned for margin accounts.",
                "example": "0.0"
              },
              "open_margin_calls": {
                "type": "string",
                "description": "Open Margin Calls:<br/> Only returned for margin accounts.<br/>RM: An RM call is triggered when an account's margin equity falls below the maintenance requirement, often due to a decline in the value of the account's positions or an increase in margin requirements for the holdings.\n<br/>RT: A Reg T call occurs when there is not enough equity in an account to cover the 50% initial margin requirement. This can happen due to open positions depreciating or from holding positions overnight that were opened with day trade buying power.\n<br/>EM: An EM call is triggered when a flagged Pattern Day Trader (PDT) account closes the prior business day below the $25,000 minimum Net Account Value (NAV) requirement. Only positions held in the margin account will count toward this requirement. Crypto, futures, event contracts, and any assets held outside the margin account are excluded from the calculation.\n<br/>DT: A Day Trade (DT) Call is issued when you exceed your available Day Trading Buying Power (DTBP) and then place a day trade. A DT call may also occur if you execute a day trade while an Equity Maintenance (EM) call is still active on your account. Exceeding day trade buying power most commonly results from trading on intraday profits. For more details on how buying power is replenished, please refer to our Margin Buying Power page.",
                "example": "['EM']",
                "enum": [
                  "EM",
                  "RM",
                  "RT",
                  "DT"
                ]
              },
              "account_currency_assets": {
                "type": "array",
                "description": "Currency assets Details",
                "items": {
                  "required": [
                    "cash_balance",
                    "currency"
                  ],
                  "type": "object",
                  "properties": {
                    "currency": {
                      "type": "string",
                      "description": "Currency",
                      "example": "USD",
                      "enum": [
                        "USD"
                      ]
                    },
                    "cash_balance": {
                      "type": "string",
                      "description": "Cash Balance.",
                      "example": "485705.95"
                    },
                    "settled_cash": {
                      "type": "string",
                      "description": "Settled Cash，<br/> Not returned for event contract accounts.<br/> Not returned for margin accounts.",
                      "example": "485705.95"
                    },
                    "unsettled_cash": {
                      "type": "string",
                      "description": "Unsettled Cash，<br/> Not returned for event contract accounts.<br/> Not returned for margin accounts.",
                      "example": "0.0"
                    },
                    "market_value": {
                      "type": "string",
                      "description": "holding market value",
                      "example": "0.0"
                    },
                    "held_amount": {
                      "type": "string",
                      "description": "In-transit funds<br/> Not returned for event contract accounts.",
                      "example": "0.0"
                    },
                    "buying_power": {
                      "type": "string",
                      "description": "Buying Power<br/> Not returned for margin accounts.",
                      "example": "484551"
                    },
                    "day_buying_power": {
                      "type": "string",
                      "description": "Day-Trade Buying Power（Margin Account）<br/> Only returned for margin accounts.",
                      "example": "0.0"
                    },
                    "overnight_buying_power": {
                      "type": "string",
                      "description": "Overnight BP（Margin Account）<br/> Only returned for margin accounts.",
                      "example": "0.0"
                    },
                    "night_trading_buying_power": {
                      "type": "string",
                      "description": "Night Trading Buying Power<br/> Not returned for event contract accounts.",
                      "example": "0.0"
                    },
                    "day_profit_loss": {
                      "type": "string",
                      "description": "Day's P&L<br/> Not returned at the Omnibus master account level.",
                      "example": "0.0"
                    },
                    "unrealized_profit_loss": {
                      "type": "string",
                      "description": "Open P&L",
                      "example": "227689"
                    },
                    "available_withdrawal": {
                      "type": "string",
                      "description": "The amount of funds currently available for withdrawal.",
                      "example": "3.0558743194E8"
                    },
                    "interests_unpaid": {
                      "type": "string",
                      "description": "Interest to be paid<br/> Not returned for event contract accounts.",
                      "example": "0.0"
                    },
                    "net_liquidation_value": {
                      "type": "string",
                      "description": "Net Account Value",
                      "example": "0.0"
                    }
                  },
                  "description": "Currency assets Details",
                  "title": "AssetsCurrencyAssets"
                }
              }
            },
            "description": "Account Balance",
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
      "content": "Retrieves the balance information for a specified account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "(Required) Account identifier.",
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

## Positions

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/account-position.md>

### List Account Positions

Retrieves positions according to the account ID.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/assets/positions/list",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves positions according to the account ID.",
  "operationId": "accountPosition",
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
                    "USD"
                  ]
                },
                "quantity": {
                  "type": "string",
                  "description": "Quantity of the position",
                  "example": "1"
                },
                "symbol": {
                  "type": "string",
                  "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market.",
                  "example": "AAPL"
                },
                "instrument_type": {
                  "type": "string",
                  "description": "Type of financial instrument associated with the request.",
                  "example": "EQUITY",
                  "enum": [
                    "EQUITY",
                    "EVENT"
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
                "event_outcome": {
                  "type": "string",
                  "description": "Event outcome decision, only applicable to event orders.",
                  "example": "yes",
                  "enum": [
                    "yes",
                    "no"
                  ]
                }
              },
              "description": "Account Position",
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
      "content": "Retrieves positions according to the account ID.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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

## Cash Activities By Type

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-cash-activity-by-type.md>

### List Account Cash Activities

Retrieves account transaction activities records with details.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/activities/cash-activities/list",
  "method": "get",
  "tags": [
    "Activity"
  ],
  "description": "Retrieves account transaction activities records with details.",
  "operationId": "brokerCashActivityByType",
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
      "description": "Account activity types",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "TRADE",
          "DEPOSIT",
          "WITHDRAW",
          "FEES",
          "JOURNAL",
          "TRANSFER",
          "INTERESTS",
          "EC_STATEMENT"
        ]
      },
      "example": "DEPOSIT,TRADE"
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "Activity query start time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:59.012Z"
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "If not provided, the default query is the last 7 days.<br/>Activity query end time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",
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
                      "description": "Activity Type\n| Code          | Description                     |\n|---------------|---------------------------------|\n|TRADE \t     | Trade Activity Type             |\n|DEPOSIT \t\t | Deposit Activity Type           |\n|WITHDRAW \t\t | Withdraw Activity Type          |\n|FEES \t\t\t | Fees Activity Type              |\n|JOURNAL \t\t | Journal Activity Type           |\n|TRANSFER \t\t | Transfer Activity Type          |\n|INTERESTS \t | Interests Activity Type         |\n|EC_STATEMENT \t | EC Statement Activity Type      |\n|OTHER \t\t | Other Activity Type             |",
                      "example": "TRADE",
                      "enum": [
                        "TRADE",
                        "DEPOSIT",
                        "WITHDRAW",
                        "FEES",
                        "JOURNAL",
                        "TRANSFER",
                        "INTERESTS",
                        "EC_STATEMENT"
                      ]
                    },
                    "activity_sub_type": {
                      "type": "string",
                      "description": "Activity Type and Sub Type Mapping.\n| ActivityType   | ActivitySubType                     |\n|----------------|-------------------------------------|\n| TRADE          | BTO                      \t\t\t|\n| TRADE          | STC                      \t\t\t|\n| DEPOSIT        | ACH                                 |\n| DEPOSIT        | WIRE                                |\n| DEPOSIT        | WIRE_REVERSE                        |\n| DEPOSIT        | ACH_REVERSE                         |\n| WITHDRAW       | ACH                                 |\n| WITHDRAW       | WIRE                                |\n| WITHDRAW       | ACH_REVERSE                         |\n| WITHDRAW       | WIRE_REVERSE                        |\n| FEES           | WIRE_DEPOSIT                        |\n| FEES           | WIRE_WITHDRAW                       |\n| FEES           | ACH_REVERSAL                        |\n| FEES           | WIRE_REVERSAL                       |\n| FEES           | SERVICE_FEE                         |\n| FEES           | ADR                                 |\n| FEES           | TRANSFER_ACATS                      |\n| FEES           | PAPER_CONFIRM_FEE                   |\n| FEES           | PAPER_STATEMENT_FEE                 |\n| FEES           | WRITE_OFF                           |\n| FEES           | PROCESSING_FEE                      |\n| FEES           | COMMISSION                          |\n| FEES           | TRANSFER_FOP                        |\n| FEES           | MONTH_END_REPORT_CHARGES            |\n| FEES           | FINANCIAL_TRANSACTION_TAX           |\n| FEES           | SETTLEMENT_FEE                      |\n| FEES           | TRANSACTION_FEE                     |\n| JOURNAL        | CASH_JOURNAL                        |\n| JOURNAL        | CREDIT                              |\n| TRANSFER        | INTERNAL_THIRD_PARTY_TRANSFER      |\n| INTERESTS      | CREDIT_CASH                         |\n| EC_STATEMENT   | EC_EXPIRATION                       |\n| EC_STATEMENT   | EC_PAYOUT                           |\n| OTHER          | OTHER                               |",
                      "example": "BTO",
                      "enum": [
                        "BTO",
                        "STC",
                        "WIRE",
                        "ACH",
                        "WIRE_REVERSE",
                        "ACH_REVERSE",
                        "WIRE_DEPOSIT",
                        "WIRE_WITHDRAW",
                        "ACH_REVERSAL",
                        "WIRE_REVERSAL",
                        "CASH_JOURNAL",
                        "INTERNAL_THIRD_PARTY_TRANSFER",
                        "CREDIT",
                        "CREDIT_CASH",
                        "EC_EXPIRATION",
                        "EC_PAYOUT",
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
    "name": "List Account Cash Activities",
    "description": {
      "content": "Retrieves account transaction activities records with details.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "Account activity types",
            "type": "text/plain"
          },
          "key": "activity_types",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Activity query start time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "If not provided, the default query is the last 7 days.<br/>Activity query end time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",
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

## Create Bank Relationship

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/create-bank-relationship.md>

### Create Bank Relationship

Creates a Bank Relationship for an Account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/bank-relationships/create",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Creates a Bank Relationship for an Account.",
  "operationId": "createBankRelationship",
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
            "account_id",
            "bank_account_name",
            "bank_account_number",
            "bank_code",
            "bank_code_type",
            "client_request_id"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4S1VI5UU6L3K8NU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account identifier.",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            },
            "bank_name": {
              "type": "string",
              "description": "Required if bank_code_type = BIC, Name of recipient bank"
            },
            "bank_code": {
              "type": "string",
              "description": "Alphanumeric ABA RTN (Routing Number)"
            },
            "bank_code_type": {
              "type": "string",
              "description": "ABA (US based bank accounts)、BIC",
              "example": "ABA",
              "enum": [
                "ABA",
                "BIC"
              ]
            },
            "bank_account_name": {
              "type": "string",
              "description": "Name of the bank account holder, as registered with the bank",
              "example": "ABA"
            },
            "bank_account_number": {
              "type": "string",
              "description": "Bank account number."
            },
            "country": {
              "type": "string",
              "description": "Required if bank_code_type = BIC, The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
            },
            "state_province": {
              "type": "string",
              "description": "Required if bank_code_type = BIC"
            },
            "postal_code": {
              "type": "string",
              "description": "Required if bank_code_type = BIC"
            },
            "city": {
              "type": "string",
              "description": "Required if bank_code_type = BIC"
            },
            "intermediary_bank_details": {
              "type": "object",
              "properties": {
                "intermediary_bank_name": {
                  "type": "string",
                  "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nIntermediary Bank Name"
                },
                "intermediary_bank_code": {
                  "type": "string",
                  "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nAlphanumeric ABA RTN (Routing Number)"
                },
                "intermediary_bank_account_number": {
                  "type": "string",
                  "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nBank account number.",
                  "example": "****6789"
                },
                "intermediary_bank_address": {
                  "type": "string",
                  "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nBank Address.",
                  "example": "****6789"
                }
              },
              "description": "Intermediary bank details.",
              "title": "IntermediaryBank"
            }
          },
          "description": "create bank relationship",
          "title": "CreateBankRelationshipRequest"
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
              "bank_account_name",
              "bank_account_number",
              "bank_code",
              "bank_code_type",
              "bank_relationship_id",
              "client_request_id",
              "create_time",
              "status",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4S1VI5UU6L3K8NU10IA"
              },
              "bank_relationship_id": {
                "type": "string",
                "description": "Relationship ID",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "bank_code": {
                "type": "string",
                "description": "Alphanumeric ABA RTN (Routing Number)"
              },
              "bank_name": {
                "type": "string",
                "description": "Bank name"
              },
              "bank_account_name": {
                "type": "string",
                "description": "Name of the bank account holder, as registered with the bank",
                "example": "ABA"
              },
              "bank_account_number": {
                "type": "string",
                "description": "Bank account number"
              },
              "bank_code_type": {
                "type": "string",
                "description": "ABA (US based bank accounts)、BIC",
                "example": "ABA",
                "enum": [
                  "ABA",
                  "BIC"
                ]
              },
              "country": {
                "type": "string",
                "description": "Only for Non-US based bank accounts, The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
              },
              "state_province": {
                "type": "string",
                "description": "Only for Non-US based bank accounts"
              },
              "postal_code": {
                "type": "string",
                "description": "Only for Non-US based bank accounts"
              },
              "city": {
                "type": "string",
                "description": "Only for Non-US based bank accounts"
              },
              "status": {
                "type": "string",
                "description": "Status",
                "example": "APPROVED",
                "enum": [
                  "APPROVED"
                ]
              },
              "intermediary_bank_details": {
                "type": "object",
                "properties": {
                  "intermediary_bank_name": {
                    "type": "string",
                    "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nIntermediary Bank Name"
                  },
                  "intermediary_bank_code": {
                    "type": "string",
                    "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nAlphanumeric ABA RTN (Routing Number)"
                  },
                  "intermediary_bank_account_number": {
                    "type": "string",
                    "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nBank account number.",
                    "example": "****6789"
                  },
                  "intermediary_bank_address": {
                    "type": "string",
                    "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nBank Address.",
                    "example": "****6789"
                  }
                },
                "description": "Intermediary bank details.",
                "title": "IntermediaryBank"
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "BankRelationshipResult"
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
    "client_request_id": "89JGKD4S1VI5UU6L3K8NU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649",
    "bank_name": "string",
    "bank_code": "string",
    "bank_code_type": "ABA",
    "bank_account_name": "ABA",
    "bank_account_number": "string",
    "country": "string",
    "state_province": "string",
    "postal_code": "string",
    "city": "string",
    "intermediary_bank_details": {
      "intermediary_bank_name": "string",
      "intermediary_bank_code": "string",
      "intermediary_bank_account_number": "****6789",
      "intermediary_bank_address": "****6789"
    }
  },
  "postman": {
    "name": "Create Bank Relationship",
    "description": {
      "content": "Creates a Bank Relationship for an Account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "bank-relationships",
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

## Delete Bank Relationship

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/delete-bank-relationship.md>

### Delete Bank Relationship

Deletes a Bank Relationship for an Account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/bank-relationships/delete",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Deletes a Bank Relationship for an Account.",
  "operationId": "deleteBankRelationship",
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
            "account_id",
            "bank_relationship_id"
          ],
          "type": "object",
          "properties": {
            "bank_relationship_id": {
              "type": "string",
              "description": "Bank relationship ID. The bank_relationship_id created from the Create Bank Relationship interface",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account identifier.",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            }
          },
          "description": "delete bank relationship",
          "title": "DeleteBankRelationshipRequest"
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
              "account_id",
              "bank_relationship_id"
            ],
            "type": "object",
            "properties": {
              "bank_relationship_id": {
                "type": "string",
                "description": "Bank relationship ID. The bank_relationship_id created from the Create Bank Relationship interface",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account identifier.",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              }
            },
            "title": "DeleteBankRelationshipResult"
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
    "bank_relationship_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649"
  },
  "postman": {
    "name": "Delete Bank Relationship",
    "description": {
      "content": "Deletes a Bank Relationship for an Account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "bank-relationships",
        "delete"
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

## List Bank Accounts

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/list-linked-bank-accounts.md>

### List Bank Relationships

Retrieves Bank Relationships for an Account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/bank-relationships/list",
  "method": "get",
  "tags": [
    "funding"
  ],
  "description": "Retrieves Bank Relationships for an Account.",
  "operationId": "listLinkedBankAccounts",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
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
                "bank_account_name",
                "bank_account_number",
                "bank_code",
                "bank_code_type",
                "bank_relationship_id",
                "create_time",
                "status",
                "update_time"
              ],
              "type": "object",
              "properties": {
                "bank_relationship_id": {
                  "type": "string",
                  "description": "Relationship ID",
                  "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
                },
                "bank_code": {
                  "type": "string",
                  "description": "Alphanumeric ABA RTN (Routing Number)"
                },
                "bank_name": {
                  "type": "string",
                  "description": "Bank name"
                },
                "bank_account_name": {
                  "type": "string",
                  "description": "Name of the bank account holder, as registered with the bank",
                  "example": "ABA"
                },
                "bank_account_number": {
                  "type": "string",
                  "description": "Bank account number"
                },
                "bank_code_type": {
                  "type": "string",
                  "description": "ABA (US based bank accounts)、BIC",
                  "example": "ABA",
                  "enum": [
                    "ABA",
                    "BIC"
                  ]
                },
                "country": {
                  "type": "string",
                  "description": "Only for Non-US based bank accounts, The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code)."
                },
                "state_province": {
                  "type": "string",
                  "description": "Only for Non-US based bank accounts"
                },
                "postal_code": {
                  "type": "string",
                  "description": "Only for Non-US based bank accounts"
                },
                "city": {
                  "type": "string",
                  "description": "Only for Non-US based bank accounts"
                },
                "status": {
                  "type": "string",
                  "description": "Status",
                  "example": "APPROVED",
                  "enum": [
                    "APPROVED"
                  ]
                },
                "intermediary_bank_details": {
                  "type": "object",
                  "properties": {
                    "intermediary_bank_name": {
                      "type": "string",
                      "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nIntermediary Bank Name"
                    },
                    "intermediary_bank_code": {
                      "type": "string",
                      "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nAlphanumeric ABA RTN (Routing Number)"
                    },
                    "intermediary_bank_account_number": {
                      "type": "string",
                      "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nBank account number.",
                      "example": "****6789"
                    },
                    "intermediary_bank_address": {
                      "type": "string",
                      "description": "This is an optional field;it is only required if an intermediary bank is involved in the transfer.\n\nBank Address.",
                      "example": "****6789"
                    }
                  },
                  "description": "Intermediary bank details.",
                  "title": "IntermediaryBank"
                },
                "create_time": {
                  "type": "string",
                  "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                  "example": "2025-01-05T22:59:59.012Z"
                },
                "update_time": {
                  "type": "string",
                  "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                  "example": "2025-01-05T22:59:59.012Z"
                }
              },
              "title": "ListBankRelationshipResult"
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
    "name": "List Bank Relationships",
    "description": {
      "content": "Retrieves Bank Relationships for an Account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "bank-relationships",
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

## Create ACH Relationship

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/create-ach-relationship.md>

### Create ACH Relationship

Creates an ACH Relationship for an Account. Only one ACH relationship is allowed per account, and creating a new one automatically replaces the existing one.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/ach-relationships/create",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Creates an ACH Relationship for an Account. Only one ACH relationship is allowed per account, and creating a new one automatically replaces the existing one.",
  "operationId": "createAchRelationship",
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
            "account_id",
            "client_request_id",
            "processor_token",
            "token_type"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4S1VI5UU6L3K8NU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account identifier.",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            },
            "token_type": {
              "type": "string",
              "description": "Token Type only support 'Plaid' ",
              "example": "Plaid",
              "enum": [
                "Plaid"
              ]
            },
            "processor_token": {
              "type": "string",
              "description": "Using Plaid, you can specify a Plaid processor token here."
            }
          },
          "description": "create bank relationship",
          "title": "CreateAchRelationshipRequest"
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
              "account_id",
              "account_owner_name",
              "ach_relationship_id",
              "bank_account_number",
              "bank_account_type",
              "bank_routing_number",
              "client_request_id",
              "create_time",
              "status",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4S1VI5UU6L3K8NU10IA"
              },
              "ach_relationship_id": {
                "type": "string",
                "description": "Relationship ID. The ach_relationship_id created from the Create ACH Relationship interface.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account identifier.",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              },
              "account_owner_name": {
                "type": "string",
                "description": "Legal full name of the bank account owner."
              },
              "bank_account_type": {
                "type": "string",
                "description": "Must be CHECKING or SAVINGS",
                "example": "CHECKING",
                "enum": [
                  "CHECKING",
                  "SAVINGS"
                ]
              },
              "bank_account_number": {
                "type": "string",
                "description": "Bank account number"
              },
              "bank_routing_number": {
                "type": "string",
                "description": "ABA routing transit number (RTN), Alphanumeric identifier for the financial institution, used for ACH transaction routing.",
                "example": "021000021"
              },
              "status": {
                "type": "string",
                "description": "Status",
                "example": "APPROVED",
                "enum": [
                  "APPROVED"
                ]
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "AchRelationshipResult"
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
    "client_request_id": "89JGKD4S1VI5UU6L3K8NU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649",
    "token_type": "Plaid",
    "processor_token": "string"
  },
  "postman": {
    "name": "Create ACH Relationship",
    "description": {
      "content": "Creates an ACH Relationship for an Account. Only one ACH relationship is allowed per account, and creating a new one automatically replaces the existing one.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "ach-relationships",
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

## Delete ACH Relationship

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/delete-ach-relationship.md>

### Delete ACH Relationship

Deletes an ACH Relationship for an Account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/ach-relationships/delete",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Deletes an ACH Relationship for an Account.",
  "operationId": "deleteAchRelationship",
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
            "account_id",
            "ach_relationship_id"
          ],
          "type": "object",
          "properties": {
            "ach_relationship_id": {
              "type": "string",
              "description": "ACH relationship ID",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account identifier.",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            }
          },
          "description": "delete ach relationship",
          "title": "DeleteAchRelationshipRequest"
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
              "account_id",
              "ach_relationship_id"
            ],
            "type": "object",
            "properties": {
              "ach_relationship_id": {
                "type": "string",
                "description": "Relationship ID. The ach_relationship_id created from the Create ACH Relationship interface.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account identifier.",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              }
            },
            "title": "DeleteAchRelationshipResult"
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
    "ach_relationship_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649"
  },
  "postman": {
    "name": "Delete ACH Relationship",
    "description": {
      "content": "Deletes an ACH Relationship for an Account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "ach-relationships",
        "delete"
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

## List ACH Relationships

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/list-ach-relationships.md>

### List ACH Relationships

Retrieves ACH Relationships for an Account.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/ach-relationships/list",
  "method": "get",
  "tags": [
    "funding"
  ],
  "description": "Retrieves ACH Relationships for an Account.",
  "operationId": "listAchRelationships",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
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
                "account_id",
                "account_owner_name",
                "ach_relationship_id",
                "bank_account_number",
                "bank_account_type",
                "bank_routing_number",
                "create_time",
                "status",
                "update_time"
              ],
              "type": "object",
              "properties": {
                "ach_relationship_id": {
                  "type": "string",
                  "description": "Relationship ID. The ach_relationship_id created from the Create ACH Relationship interface.",
                  "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
                },
                "account_id": {
                  "type": "string",
                  "description": "Account identifier.",
                  "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
                },
                "account_owner_name": {
                  "type": "string",
                  "description": "Legal full name of the bank account owner."
                },
                "bank_account_type": {
                  "type": "string",
                  "description": "Must be CHECKING or SAVINGS",
                  "example": "CHECKING",
                  "enum": [
                    "CHECKING",
                    "SAVINGS"
                  ]
                },
                "bank_account_number": {
                  "type": "string",
                  "description": "Bank account number"
                },
                "bank_routing_number": {
                  "type": "string",
                  "description": "ABA routing transit number (RTN), Alphanumeric identifier for the financial institution, used for ACH transaction routing.",
                  "example": "021000021"
                },
                "status": {
                  "type": "string",
                  "description": "Status",
                  "example": "APPROVED",
                  "enum": [
                    "APPROVED"
                  ]
                },
                "create_time": {
                  "type": "string",
                  "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                  "example": "2025-01-05T22:59:59.012Z"
                },
                "update_time": {
                  "type": "string",
                  "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                  "example": "2025-01-05T22:59:59.012Z"
                }
              },
              "title": "ListAchRelationshipResult"
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
    "name": "List ACH Relationships",
    "description": {
      "content": "Retrieves ACH Relationships for an Account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "ach-relationships",
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

## Create Transfer

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/create-transfer.md>

### Create Transfer

Initiates a funding transfer to an account. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Funding Transfer Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/transfers/create",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Initiates a funding transfer to an account. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Funding Transfer Events](/apis/docs/reference/fd-events/funding-events#funding-transfer-events)",
  "operationId": "createTransfer",
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
            "account_id",
            "amount",
            "client_request_id",
            "currency",
            "direction",
            "transfer_type"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4S1VI5UU6L3K8NU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account ID",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            },
            "transfer_type": {
              "type": "string",
              "description": "Transfer type,current support WIRE or ACH",
              "example": "WIRE",
              "enum": [
                "ACH",
                "WIRE"
              ]
            },
            "ach_relationship_id": {
              "type": "string",
              "description": "Required if transfer_type = ACH ",
              "example": "89JGKD4S1VIX1U609K8NU10IA"
            },
            "bank_relationship_id": {
              "type": "string",
              "description": "Required if transfer_type = WIRE ",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "amount": {
              "type": "string",
              "description": "Transfer amount",
              "example": "1000.0"
            },
            "currency": {
              "type": "string",
              "description": "Currency code",
              "example": "USD"
            },
            "direction": {
              "type": "string",
              "description": "Transfer direction",
              "example": "DEPOSIT",
              "enum": [
                "DEPOSIT",
                "WITHDRAWAL"
              ]
            }
          },
          "description": "create bank relationship",
          "title": "CreateTransferRequest"
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
              "account_id",
              "amount",
              "client_request_id",
              "create_time",
              "currency",
              "direction",
              "status",
              "transfer_id",
              "transfer_type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              },
              "transfer_id": {
                "type": "string",
                "description": "Transfer id generated by the system.",
                "example": "0KGOHL4PR2SLC0DKIND4TI0002"
              },
              "transfer_type": {
                "type": "string",
                "description": "transfer type",
                "example": "ACH",
                "enum": [
                  "ACH",
                  "WIRE"
                ]
              },
              "ach_relationship_id": {
                "type": "string",
                "description": "Required if transfer_type is ACH ",
                "example": "89JGKD4S1VIX1U609K8NU10IA"
              },
              "bank_relationship_id": {
                "type": "string",
                "description": "Required if transfer_type is WIRE.",
                "example": "C0DKID4LKIVI5UU6L3KHNU10HL4"
              },
              "amount": {
                "type": "string",
                "description": "Transfer amount",
                "example": "1000.0"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "direction": {
                "type": "string",
                "description": "Transfer direction",
                "example": "DEPOSIT",
                "enum": [
                  "DEPOSIT",
                  "WITHDRAWAL"
                ]
              },
              "status": {
                "type": "string",
                "description": "Transfer status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "COMPLETED",
                  "RETURNED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "reason of the status"
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "fee": {
                "type": "string",
                "description": "Fee amount to be collected. Only applies when type is WIRE"
              }
            },
            "title": "TransferResult"
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
    "client_request_id": "89JGKD4S1VI5UU6L3K8NU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649",
    "transfer_type": "WIRE",
    "ach_relationship_id": "89JGKD4S1VIX1U609K8NU10IA",
    "bank_relationship_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "amount": "1000.0",
    "currency": "USD",
    "direction": "DEPOSIT"
  },
  "postman": {
    "name": "Create Transfer",
    "description": {
      "content": "Initiates a funding transfer to an account. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Funding Transfer Events](/apis/docs/reference/fd-events/funding-events#funding-transfer-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "transfers",
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

## Transfer List

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/transfer-list.md>

### List Transfer Records

Retrieves account fund transfers by specified criteria. Transfer records are sorted by client_request_id. If the client_request_id values in the request are provided in order, the response will preserve the same ordering. 

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/transfers/list",
  "method": "get",
  "tags": [
    "funding"
  ],
  "description": "Retrieves account fund transfers by specified criteria. Transfer records are sorted by client_request_id. If the client_request_id values in the request are provided in order, the response will preserve the same ordering. ",
  "operationId": "transferList",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "The start datetime to include in the date range, formatted as yyyy-MM-dd'T'HH:mm:ss.SSSZ.<br/> Both start_time and end_time should be provided, end_time should be greater than start_time.<br/> By default, records from the last 7 days will be queried if no time range is specified.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:59.012Z"
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "The end datetime to include in the date range, formatted as yyyy-MM-dd'T'HH:mm:ss.SSSZ.<br/> Both start_time and end_time should be provided, end_date should be greater than start_time.<br/> By default, records from the last 7 days will be queried if no time range is specified.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-06T22:59:59.012Z"
    },
    {
      "name": "transfer_types",
      "in": "query",
      "description": "Transfer types",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "ACH",
          "WIRE"
        ]
      },
      "example": "ACH,WIRE"
    },
    {
      "name": "directions",
      "in": "query",
      "description": "Transfer directions",
      "required": false,
      "schema": {
        "type": "string",
        "description": "Transfer direction",
        "enum": [
          "DEPOSIT",
          "WITHDRAWAL"
        ]
      },
      "example": "DEPOSIT,WITHDRAWAL"
    },
    {
      "name": "statuses",
      "in": "query",
      "description": "Transfer statuses",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "SUBMITTED",
          "CANCELED",
          "REJECTED",
          "COMPLETED",
          "RETURNED"
        ]
      },
      "example": "CANCELED,REJECTED"
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
                    "amount",
                    "client_request_id",
                    "create_time",
                    "currency",
                    "direction",
                    "status",
                    "transfer_id",
                    "transfer_type",
                    "update_time"
                  ],
                  "type": "object",
                  "properties": {
                    "client_request_id": {
                      "type": "string",
                      "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                      "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
                    },
                    "account_id": {
                      "type": "string",
                      "description": "Account ID",
                      "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
                    },
                    "transfer_id": {
                      "type": "string",
                      "description": "Transfer id generated by the system.",
                      "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                    },
                    "transfer_type": {
                      "type": "string",
                      "description": "transfer type",
                      "example": "ACH",
                      "enum": [
                        "ACH",
                        "WIRE"
                      ]
                    },
                    "ach_relationship_id": {
                      "type": "string",
                      "description": "Required if transfer_type is ACH ",
                      "example": "89JGKD4S1VIX1U609K8NU10IA"
                    },
                    "bank_relationship_id": {
                      "type": "string",
                      "description": "Required if transfer_type is WIRE.",
                      "example": "C0DKID4LKIVI5UU6L3KHNU10HL4"
                    },
                    "amount": {
                      "type": "string",
                      "description": "Transfer amount",
                      "example": "1000.0"
                    },
                    "currency": {
                      "type": "string",
                      "description": "Currency",
                      "example": "USD",
                      "enum": [
                        "USD"
                      ]
                    },
                    "direction": {
                      "type": "string",
                      "description": "Transfer direction",
                      "example": "DEPOSIT",
                      "enum": [
                        "DEPOSIT",
                        "WITHDRAWAL"
                      ]
                    },
                    "status": {
                      "type": "string",
                      "description": "Transfer status",
                      "example": "SUBMITTED",
                      "enum": [
                        "SUBMITTED",
                        "CANCELED",
                        "REJECTED",
                        "COMPLETED",
                        "RETURNED"
                      ]
                    },
                    "reason": {
                      "type": "string",
                      "description": "reason of the status"
                    },
                    "create_time": {
                      "type": "string",
                      "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                      "example": "2025-01-05T22:59:59.012Z"
                    },
                    "update_time": {
                      "type": "string",
                      "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                      "example": "2025-01-05T22:59:59.012Z"
                    },
                    "fee": {
                      "type": "string",
                      "description": "Fee amount to be collected. Only applies when type is WIRE"
                    }
                  },
                  "title": "TransferResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoTransferResult"
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
    "name": "List Transfer Records",
    "description": {
      "content": "Retrieves account fund transfers by specified criteria. Transfer records are sorted by client_request_id. If the client_request_id values in the request are provided in order, the response will preserve the same ordering. ",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "transfers",
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
        },
        {
          "disabled": false,
          "description": {
            "content": "The start datetime to include in the date range, formatted as yyyy-MM-dd'T'HH:mm:ss.SSSZ.<br/> Both start_time and end_time should be provided, end_time should be greater than start_time.<br/> By default, records from the last 7 days will be queried if no time range is specified.",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The end datetime to include in the date range, formatted as yyyy-MM-dd'T'HH:mm:ss.SSSZ.<br/> Both start_time and end_time should be provided, end_date should be greater than start_time.<br/> By default, records from the last 7 days will be queried if no time range is specified.",
            "type": "text/plain"
          },
          "key": "end_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Transfer types",
            "type": "text/plain"
          },
          "key": "transfer_types",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Transfer directions",
            "type": "text/plain"
          },
          "key": "directions",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Transfer statuses",
            "type": "text/plain"
          },
          "key": "statuses",
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

## Transfer Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/transfer-detail.md>

### Get Transfer Detail

Retrieves Account Fund Transfer Details.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/transfers/get",
  "method": "get",
  "tags": [
    "funding"
  ],
  "description": "Retrieves Account Fund Transfer Details.",
  "operationId": "transferDetail",
  "parameters": [
    {
      "name": "client_request_id",
      "in": "query",
      "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LJIS16BACHQG9LPP44L9IQHGAB"
    },
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
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
              "account_id",
              "amount",
              "client_request_id",
              "create_time",
              "currency",
              "direction",
              "status",
              "transfer_id",
              "transfer_type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              },
              "transfer_id": {
                "type": "string",
                "description": "Transfer id generated by the system.",
                "example": "0KGOHL4PR2SLC0DKIND4TI0002"
              },
              "transfer_type": {
                "type": "string",
                "description": "transfer type",
                "example": "ACH",
                "enum": [
                  "ACH",
                  "WIRE"
                ]
              },
              "ach_relationship_id": {
                "type": "string",
                "description": "Required if transfer_type is ACH ",
                "example": "89JGKD4S1VIX1U609K8NU10IA"
              },
              "bank_relationship_id": {
                "type": "string",
                "description": "Required if transfer_type is WIRE.",
                "example": "C0DKID4LKIVI5UU6L3KHNU10HL4"
              },
              "amount": {
                "type": "string",
                "description": "Transfer amount",
                "example": "1000.0"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "direction": {
                "type": "string",
                "description": "Transfer direction",
                "example": "DEPOSIT",
                "enum": [
                  "DEPOSIT",
                  "WITHDRAWAL"
                ]
              },
              "status": {
                "type": "string",
                "description": "Transfer status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "COMPLETED",
                  "RETURNED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "reason of the status"
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "fee": {
                "type": "string",
                "description": "Fee amount to be collected. Only applies when type is WIRE"
              }
            },
            "title": "TransferResult"
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
    "name": "Get Transfer Detail",
    "description": {
      "content": "Retrieves Account Fund Transfer Details.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "transfers",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
            "type": "text/plain"
          },
          "key": "client_request_id",
          "value": ""
        },
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

## Cancel Transfer

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/cancel-transfer.md>

### Cancel Transfer

Cancels a fund transfer.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/transfers/cancel",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Cancels a fund transfer.",
  "operationId": "cancelTransfer",
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
            "account_id",
            "client_request_id"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account ID",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            }
          },
          "description": "cancel transfer",
          "title": "CancelTransferRequest"
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
              "account_id",
              "client_request_id",
              "transfer_id"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "transfer_id": {
                "type": "string",
                "description": "tranfer_id",
                "example": "06I479T2KPTN1S7U89396NTLB9"
              },
              "account_id": {
                "type": "string",
                "description": "Account identifier.",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              }
            },
            "title": "CancelTransferResult"
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
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649"
  },
  "postman": {
    "name": "Cancel Transfer",
    "description": {
      "content": "Cancels a fund transfer.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "transfers",
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

## Create Instant Funding

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-instant-create.md>

### Create Instant Funding

Initiates an immediate deposit request. The system instantly processes and credits the funds to the target account, making them available for trading immediately. No settlement institution is involved in this instant funding process.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Instant Funding Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/instant-funding/create",
  "method": "post",
  "tags": [
    "funding"
  ],
  "description": "Initiates an immediate deposit request. The system instantly processes and credits the funds to the target account, making them available for trading immediately. No settlement institution is involved in this instant funding process.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Instant Funding Events](/apis/docs/reference/fd-events/funding-events#instant-funding-events)",
  "operationId": "brokerFundingInstantCreate",
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
            "account_id",
            "amount",
            "client_request_id",
            "currency",
            "type"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "account_id": {
              "type": "string",
              "description": "Account ID",
              "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
            },
            "type": {
              "type": "string",
              "description": "Transfer direction",
              "example": "DEPOSIT",
              "enum": [
                "DEPOSIT",
                "WITHDRAWAL"
              ]
            },
            "amount": {
              "type": "string",
              "description": "Amount.",
              "example": "100"
            },
            "currency": {
              "type": "string",
              "description": "Currency",
              "example": "USD",
              "enum": [
                "USD"
              ]
            }
          },
          "title": "CreateInstantFundingParam"
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
              "account_id",
              "amount",
              "client_request_id",
              "create_time",
              "currency",
              "instant_funding_id",
              "status",
              "type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              },
              "type": {
                "type": "string",
                "description": "Transfer direction",
                "example": "DEPOSIT",
                "enum": [
                  "DEPOSIT",
                  "WITHDRAWAL"
                ]
              },
              "instant_funding_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount.",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "InstantFundingResult"
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
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "account_id": "J6HA4EBQRQFJD2J6NQH0F7M649",
    "type": "DEPOSIT",
    "amount": "100",
    "currency": "USD"
  },
  "postman": {
    "name": "Create Instant Funding",
    "description": {
      "content": "Initiates an immediate deposit request. The system instantly processes and credits the funds to the target account, making them available for trading immediately. No settlement institution is involved in this instant funding process.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Instant Funding Events](/apis/docs/reference/fd-events/funding-events#instant-funding-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "instant-funding",
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

## Instant Funding Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-instant-query.md>

### Get Instant Funding Detail

Retrieves instant funding record details.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/funding/instant-funding/get",
  "method": "get",
  "tags": [
    "funding"
  ],
  "description": "Retrieves instant funding record details.",
  "operationId": "brokerFundingInstantQuery",
  "parameters": [
    {
      "name": "client_request_id",
      "in": "query",
      "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LJIS16BACHQG9LPP44L9IQHGAB"
    },
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
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
              "account_id",
              "amount",
              "client_request_id",
              "create_time",
              "currency",
              "instant_funding_id",
              "status",
              "type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "J6HA4EBQRQFJD2J6NQH0F7M649"
              },
              "type": {
                "type": "string",
                "description": "Transfer direction",
                "example": "DEPOSIT",
                "enum": [
                  "DEPOSIT",
                  "WITHDRAWAL"
                ]
              },
              "instant_funding_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount.",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "InstantFundingResult"
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
    "name": "Get Instant Funding Detail",
    "description": {
      "content": "Retrieves instant funding record details.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "instant-funding",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
            "type": "text/plain"
          },
          "key": "client_request_id",
          "value": ""
        },
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

## Create Fee

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-fee-create.md>

### Create Fee Deduction

Initial Fee Deduction debits a specified account and credits an equal amount to a contra account, without involving any settlement institutions.Use this API to deduct various fees (such as ADR, SERVICE_FEE, etc.) from a specified account.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Fee Deduction Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/fees/create",
  "method": "post",
  "tags": [
    "Fees and Credits"
  ],
  "description": "Initial Fee Deduction debits a specified account and credits an equal amount to a contra account, without involving any settlement institutions.Use this API to deduct various fees (such as ADR, SERVICE_FEE, etc.) from a specified account.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Fee Deduction Events](/apis/docs/reference/fd-events/fee-credit-events#fee-deduction-events)",
  "operationId": "brokerFundingFeeCreate",
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
            "account",
            "amount",
            "client_request_id",
            "contra_account",
            "currency",
            "type"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "account": {
              "type": "string",
              "description": "Account id.",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "contra_account": {
              "type": "string",
              "description": "Firm Account",
              "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
            },
            "type": {
              "type": "string",
              "description": "Fee Type",
              "example": "SERVICE_FEE",
              "enum": [
                "SERVICE_FEE",
                "ADR",
                "TRANSFER_ACATS",
                "PAPER_CONFIRM_FEE",
                "PAPER_STATEMENT_FEE",
                "WRITE_OFF",
                "PROCESSING_FEE",
                "COMMISSION",
                "TRANSFER_FOP",
                "MONTH_END_REPORT_CHARGES",
                "FTT",
                "SETTLEMENT_FEES",
                "TRANSACTION_FEES"
              ]
            },
            "amount": {
              "type": "string",
              "description": "Amount.",
              "example": "100"
            },
            "currency": {
              "type": "string",
              "description": "Currency",
              "example": "USD",
              "enum": [
                "USD"
              ]
            }
          },
          "title": "CreateFeeParam"
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
              "account",
              "amount",
              "client_request_id",
              "contra_account",
              "create_time",
              "currency",
              "fee_id",
              "status",
              "type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account": {
                "type": "string",
                "description": "Account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "contra_account": {
                "type": "string",
                "description": "Firm Account",
                "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
              },
              "type": {
                "type": "string",
                "description": "Fee Type",
                "example": "SERVICE_FEE",
                "enum": [
                  "SERVICE_FEE",
                  "ADR",
                  "TRANSFER_ACATS",
                  "PAPER_CONFIRM_FEE",
                  "PAPER_STATEMENT_FEE",
                  "WRITE_OFF",
                  "PROCESSING_FEE",
                  "COMMISSION",
                  "TRANSFER_FOP",
                  "MONTH_END_REPORT_CHARGES",
                  "FTT",
                  "SETTLEMENT_FEES",
                  "TRANSACTION_FEES"
                ]
              },
              "fee_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount.",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "FeeResult"
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
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "account": "93IUJ28O9VO2KBGHDHR4H9",
    "contra_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
    "type": "SERVICE_FEE",
    "amount": "100",
    "currency": "USD"
  },
  "postman": {
    "name": "Create Fee Deduction",
    "description": {
      "content": "Initial Fee Deduction debits a specified account and credits an equal amount to a contra account, without involving any settlement institutions.Use this API to deduct various fees (such as ADR, SERVICE_FEE, etc.) from a specified account.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Fee Deduction Events](/apis/docs/reference/fd-events/fee-credit-events#fee-deduction-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "fees",
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

## Get Fee

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-fee-query.md>

### Get Fee Deduction Detail

Retrieves fee record details..

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/fees/get",
  "method": "get",
  "tags": [
    "Fees and Credits"
  ],
  "description": "Retrieves fee record details..",
  "operationId": "brokerFundingFeeQuery",
  "parameters": [
    {
      "name": "client_request_id",
      "in": "query",
      "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LJIS16BACHQG9LPP44L9IQHGAB"
    },
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "93IUJ28O9VO2KBGHDHR4H9"
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
              "account",
              "amount",
              "client_request_id",
              "contra_account",
              "create_time",
              "currency",
              "fee_id",
              "status",
              "type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account": {
                "type": "string",
                "description": "Account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "contra_account": {
                "type": "string",
                "description": "Firm Account",
                "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
              },
              "type": {
                "type": "string",
                "description": "Fee Type",
                "example": "SERVICE_FEE",
                "enum": [
                  "SERVICE_FEE",
                  "ADR",
                  "TRANSFER_ACATS",
                  "PAPER_CONFIRM_FEE",
                  "PAPER_STATEMENT_FEE",
                  "WRITE_OFF",
                  "PROCESSING_FEE",
                  "COMMISSION",
                  "TRANSFER_FOP",
                  "MONTH_END_REPORT_CHARGES",
                  "FTT",
                  "SETTLEMENT_FEES",
                  "TRANSACTION_FEES"
                ]
              },
              "fee_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount.",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "FeeResult"
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
    "name": "Get Fee Deduction Detail",
    "description": {
      "content": "Retrieves fee record details..",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "fees",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
            "type": "text/plain"
          },
          "key": "client_request_id",
          "value": ""
        },
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

## Create Credit

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-credit-create.md>

### Create New Credit

Initial Fee Credit credits a specified account and deducts an equal amount from a contra account, without involving any settlement institutions. Use this API to credit funds to a specified account (such as GIFTING). <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Credit Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/credits/create",
  "method": "post",
  "tags": [
    "Fees and Credits"
  ],
  "description": "Initial Fee Credit credits a specified account and deducts an equal amount from a contra account, without involving any settlement institutions. Use this API to credit funds to a specified account (such as GIFTING). <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Credit Events](/apis/docs/reference/fd-events/fee-credit-events#credit-events)",
  "operationId": "brokerFundingCreditCreate",
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
            "account",
            "amount",
            "client_request_id",
            "contra_account",
            "currency",
            "type"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "account": {
              "type": "string",
              "description": "Account id.",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "contra_account": {
              "type": "string",
              "description": "Firm Account",
              "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
            },
            "type": {
              "type": "string",
              "description": "Credit Type",
              "example": "GIFTING",
              "enum": [
                "GIFTING"
              ]
            },
            "amount": {
              "type": "string",
              "description": "Amount.",
              "example": "100"
            },
            "currency": {
              "type": "string",
              "description": "Currency",
              "example": "USD",
              "enum": [
                "USD"
              ]
            }
          },
          "title": "CreateCreditParam"
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
              "account",
              "amount",
              "client_request_id",
              "contra_account",
              "create_time",
              "credit_id",
              "currency",
              "status",
              "type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account": {
                "type": "string",
                "description": "Account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "contra_account": {
                "type": "string",
                "description": "Firm Account",
                "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
              },
              "type": {
                "type": "string",
                "description": "Credit Type",
                "example": "GIFTING",
                "enum": [
                  "GIFTING"
                ]
              },
              "credit_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount.",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "CreditResult"
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
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "account": "93IUJ28O9VO2KBGHDHR4H9",
    "contra_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
    "type": "GIFTING",
    "amount": "100",
    "currency": "USD"
  },
  "postman": {
    "name": "Create New Credit",
    "description": {
      "content": "Initial Fee Credit credits a specified account and deducts an equal amount from a contra account, without involving any settlement institutions. Use this API to credit funds to a specified account (such as GIFTING). <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Credit Events](/apis/docs/reference/fd-events/fee-credit-events#credit-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "credits",
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

## Get Credit

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-credit-query.md>

### Get Credit Detail

Retrieves credit record details.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/credits/get",
  "method": "get",
  "tags": [
    "Fees and Credits"
  ],
  "description": "Retrieves credit record details.",
  "operationId": "brokerFundingCreditQuery",
  "parameters": [
    {
      "name": "client_request_id",
      "in": "query",
      "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LJIS16BACHQG9LPP44L9IQHGAB"
    },
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "93IUJ28O9VO2KBGHDHR4H9"
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
              "account",
              "amount",
              "client_request_id",
              "contra_account",
              "create_time",
              "credit_id",
              "currency",
              "status",
              "type",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "account": {
                "type": "string",
                "description": "Account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "contra_account": {
                "type": "string",
                "description": "Firm Account",
                "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
              },
              "type": {
                "type": "string",
                "description": "Credit Type",
                "example": "GIFTING",
                "enum": [
                  "GIFTING"
                ]
              },
              "credit_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount.",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "CreditResult"
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
    "name": "Get Credit Detail",
    "description": {
      "content": "Retrieves credit record details.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "credits",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
            "type": "text/plain"
          },
          "key": "client_request_id",
          "value": ""
        },
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

## List Stock Instruments

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/list-stock-instruments.md>

### List Stock Instruments

Retrieves profile information for one or more instruments.<br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Instrument Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/instruments/stocks/profiles/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves profile information for one or more instruments.<br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Instrument Events](../fd-events/instruments-events)",
  "operationId": "listStockInstruments",
  "parameters": [
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
      "name": "category",
      "in": "query",
      "description": "Security type.",
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
      "name": "sub_category",
      "in": "query",
      "description": "Sub-category of the instrument. Only effective when symbols is not specified. When category = US_STOCK, supported values: COMMON_STOCK, ETF, PREFERRED_STOCK, WARRANT, UNITS, RIGHT. If not specified, returns all sub-categories.",
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
      "name": "status",
      "in": "query",
      "description": "Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable)",
      "required": false,
      "schema": {
        "type": "string",
        "enum": [
          "OC",
          "CO",
          "NT"
        ]
      },
      "example": "CO"
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
                        "US_STOCK"
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
                      "description": "Margin requirement ratio for long position. 0.5 represents 50%",
                      "example": "0.5"
                    },
                    "margin_requirement_short": {
                      "type": "string",
                      "description": "Margin requirement ratio for short position. 0.5 represents 50%",
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
      "content": "Retrieves profile information for one or more instruments.<br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Instrument Events](../fd-events/instruments-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "List of security symbols, maximum 100 symbols per query.",
            "type": "text/plain"
          },
          "key": "symbols",
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
            "content": "Sub-category of the instrument. Only effective when symbols is not specified. When category = US_STOCK, supported values: COMMON_STOCK, ETF, PREFERRED_STOCK, WARRANT, UNITS, RIGHT. If not specified, returns all sub-categories.",
            "type": "text/plain"
          },
          "key": "sub_category",
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

## Event Contract Categories

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-categories-list.md>

### List Event Contract Categories

Retrieves all categories under the Event Contract.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/instruments/event-contracts/categories/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves all categories under the Event Contract.",
  "operationId": "brokerEventCategoriesList",
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
        "broker",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-series-list.md>

### List Event Contract Series

Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., “Monthly Jobs Report” ). This endpoint allows you to browse and discover available series templates by category.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/instruments/event-contracts/series/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., “Monthly Jobs Report” ). This endpoint allows you to browse and discover available series templates by category.",
  "operationId": "brokerEventSeriesList",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "The category which this series belongs to.",
      "required": true,
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
      "example": "eyJ2IjoxLCJsYXN0SWQiOiIxNTE3NTAiLCJwYWdlSW===="
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
                    "series_id": {
                      "type": "integer",
                      "description": "Series ID",
                      "format": "int32",
                      "example": 152917
                    }
                  },
                  "description": "Event series information",
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
      "content": "Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., “Monthly Jobs Report” ). This endpoint allows you to browse and discover available series templates by category.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "(Required) The category which this series belongs to.",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-events-list.md>

### List Event Contract Events

Retrieves events under the Event Contract.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/instruments/event-contracts/events/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves events under the Event Contract.",
  "operationId": "brokerEventEventsList",
  "parameters": [
    {
      "name": "series_symbol",
      "in": "query",
      "description": "Symbol that identifies this series.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXRATECUTCOUNT"
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
      "example": "KXRATECUTCOUNT,KXFEDDECISION"
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
                  "description": "status, eg:ACTIVE=event is open for trading; INACTIVE=event is closed or settled, default:ACTIVE",
                  "example": "ACTIVE"
                },
                "series_id": {
                  "type": "integer",
                  "description": "Series ID",
                  "format": "int32",
                  "example": 152917
                },
                "short_name": {
                  "type": "string",
                  "description": "Event short name",
                  "example": "2026-27"
                },
                "strike_date": {
                  "type": "string",
                  "description": "strike date, Only present for sports category events. Represents the settlement date/period of the event.",
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
              "description": "Event contract event",
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
      "content": "Retrieves events under the Event Contract.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-market-list.md>

### List Event Contract Instruments

Retrieves profile information for event contract markets based on the series symbol.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/instruments/event-contracts/markets/list",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves profile information for event contract markets based on the series symbol.",
  "operationId": "brokerEventMarketList",
  "parameters": [
    {
      "name": "series_symbol",
      "in": "query",
      "description": "Symbol that identifies this series.",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "KXRATECUTCOUNT"
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
      "example": "eyJ2IjoxLCJsYXN0SWQiOiI1MDIyNTcyNTciLCJwYWdlSW===="
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
    "name": "List Event Contract Instruments",
    "description": {
      "content": "Retrieves profile information for event contract markets based on the series symbol.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "(Required) Symbol that identifies this series.",
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

## Corporate Actions Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-corporate-actions-detail.md>

### Get Corporate Actions Detail

Retrieves corporate actions detail information,only support stocks. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Corporate Actions Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/instruments/stocks/corporate-actions/get",
  "method": "get",
  "tags": [
    "Instruments"
  ],
  "description": "Retrieves corporate actions detail information,only support stocks. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Corporate Actions Events](../fd-events/ca-events)",
  "operationId": "brokerCorporateActionsDetail",
  "parameters": [
    {
      "name": "event_id",
      "in": "query",
      "description": "Corporate Event Action ID",
      "required": true,
      "schema": {
        "type": "string"
      },
      "example": "CA123456789"
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
              "event_id": {
                "type": "string",
                "description": "Company Event ID",
                "example": "1234567890"
              },
              "event_type": {
                "type": "string",
                "description": "Corporate Event Type<br/>IDENTIFIER_CHANGE - Identifier-related changes to a security without economic impact, including symbol, exchange, ISIN, CUSIP, or instrument name changes<br/>DIVIDEND - Dividend and distribution events, including cash dividend, stock dividend, optional dividend, and return of capital<br/>REVERSE_SPLIT - Reverse stock split that consolidates shares and reduces the number of outstanding shares<br/>FORWARD_SPLIT - Forward stock split that increases the number of outstanding shares<br/>BONUS_ISSUE - Bonus issue of additional shares distributed to existing shareholders at no cost<br/>RIGHTS_OFFERING - Rights offering allowing shareholders to subscribe for additional shares<br/>DISTRIBUTION - Distribution of cash, securities, or other assets to shareholders<br/>SPIN_OFF - Spin-off event where shares of a subsidiary or new entity are distributed to existing shareholders<br/>UNIT_SPLIT - Unit split event affecting composite or unit-based securities<br/>MERGER - Merger or acquisition event involving the combination of two or more entities<br/>FULL_CALL - Full call redemption of the entire outstanding security issue<br/>PARTIAL_CALL - Partial call redemption affecting only a portion of the outstanding issue<br/>EXCHANGE - Exchange event where existing securities are exchanged for new securities or other consideration<br/>DTC_EXIT - Event indicating a security is no longer eligible for DTC settlement or custody<br/>LIQUIDATION - Liquidation event involving the winding up of an issuer and asset distribution<br/>WORTHLESS - Worthless security event indicating the security has no residual value<br/>ADR_GDR_TERMINATION - Termination of an ADR or GDR program<br/>MATURITY - Maturity event where a security reaches its contractual maturity date<br/>ADR_FEE - ADR fee charged to holders of American Depositary Receipts<br/>CONVERSION - Conversion event where securities are converted into another class or form<br/>OPEN_OFFER - Open offer allowing shareholders to subscribe for additional securities<br/>PREFERENTIAL_OFFER - Preferential offer made to selected shareholders under specific terms<br/>PERFORMANCE_COMPENSATION - Performance compensation event related to performance commitments, commonly in A-share markets<br/>DELISTING - Delisting event where a security is removed from exchange trading",
                "example": "DIVIDEND",
                "enum": [
                  "IDENTIFIER_CHANGE",
                  "DIVIDEND",
                  "REVERSE_SPLIT",
                  "FORWARD_SPLIT",
                  "BONUS_ISSUE",
                  "RIGHTS_OFFERING",
                  "DISTRIBUTION",
                  "SPIN_OFF",
                  "UNIT_SPLIT",
                  "MERGER",
                  "FULL_CALL",
                  "PARTIAL_CALL",
                  "EXCHANGE",
                  "DTC_EXIT",
                  "LIQUIDATION",
                  "WORTHLESS",
                  "ADR_GDR_TERMINATION",
                  "MATURITY",
                  "ADR_FEE",
                  "CONVERSION",
                  "OPEN_OFFER",
                  "PREFERENTIAL_OFFER",
                  "PERFORMANCE_COMPENSATION",
                  "DELISTING"
                ]
              },
              "event_version": {
                "type": "string",
                "description": "Event Version",
                "example": "1"
              },
              "instrument_id": {
                "type": "string",
                "description": "Instrument ID",
                "example": "943a9802f6c14983b3b4755c69c01717"
              },
              "category": {
                "type": "string",
                "description": "Instrument Category",
                "example": "US_STOCK",
                "enum": [
                  "US_STOCK"
                ]
              },
              "record_date": {
                "type": "string",
                "description": "Record Date (YYYY-MM-DD)",
                "example": "2024-12-31"
              },
              "ex_date": {
                "type": "string",
                "description": "Ex Date (YYYY-MM-DD)",
                "example": "2024-12-30"
              },
              "payment_date": {
                "type": "string",
                "description": "Payment Date (YYYY-MM-DD)",
                "example": "2025-01-15"
              },
              "final_pay_date": {
                "type": "string",
                "description": "Final Payment Date (YYYY-MM-DD)",
                "example": "2025-02-15"
              },
              "country_code": {
                "type": "string",
                "description": "Country Code, ISO 3166-1 alpha-2 format",
                "example": "US"
              },
              "listing_country_of_code": {
                "type": "string",
                "description": "Listing Country Code, ISO 3166-1 alpha-2 format",
                "example": "US"
              },
              "issuer_country_code": {
                "type": "string",
                "description": "Issuer Country Code, ISO 3166-1 alpha-2 format",
                "example": "US"
              },
              "from": {
                "type": "object",
                "properties": {
                  "symbol": {
                    "type": "string",
                    "description": "Symbol of the instrument",
                    "example": "NVDA"
                  },
                  "name": {
                    "type": "string",
                    "description": "Name of the instrument"
                  },
                  "exchange": {
                    "type": "string",
                    "description": "Exchange code",
                    "example": "CCC"
                  }
                },
                "description": "Event From Info. Position's instrument information",
                "title": "EventFromInfo"
              },
              "to": {
                "type": "array",
                "description": "Event To Info. Corporate action's target instrument information",
                "items": {
                  "type": "object",
                  "properties": {
                    "option_number": {
                      "type": "string",
                      "description": "Option Number. Identifier for the payout option",
                      "example": "1"
                    },
                    "description": {
                      "type": "string",
                      "description": "Description of the payout option",
                      "example": "Securities"
                    },
                    "default_option_flag": {
                      "type": "string",
                      "description": "Default Option Flag. Indicates if this is the default payout option",
                      "example": "true"
                    },
                    "payouts": {
                      "type": "array",
                      "description": "Payouts associated with this option",
                      "items": {
                        "type": "object",
                        "properties": {
                          "type": {
                            "type": "string",
                            "description": "Payout Nature Type<br/>DV - Dividend: Cash or stock distribution paid to shareholders<br/>FR - Franked Dividend: Dividend paid with franking credits attached<br/>IN - Interest: Interest income distribution<br/>L2 - Long Term Capital Gains: Gains from disposal of assets held longer than one year<br/>OT - Other: Other types of income or entitlement (see extended terms)<br/>C - Cash: Cash payment (only applicable to events created prior to release 2)<br/>PC - Cash/Principal/Return of Capital: Cash returned as principal or capital<br/>PM - Premium: Additional amount paid over base entitlement<br/>S - Securities: Distribution of securities instead of cash<br/>ST - Short Term Capital Gains: Gains from disposal of assets held less than one year<br/>SI - Sundry Income: Miscellaneous income distributions<br/>UF - Unfranked Dividend: Dividend paid without franking credits<br/>PI - Property Income Distribution: Income derived from property assets<br/>TD - Tax Deferred: Income deferred for tax purposes<br/>TE - Tax Exempted: Income exempted from taxation<br/>FI - Foreign Income: Income sourced from foreign jurisdictions<br/>CD - Capital gain on disposal of taxable property - Discounted<br/>CO - Capital gain on disposal of taxable property - Other<br/>CC - Capital gain on disposal of taxable property - Concessional<br/>CN - Capital gain on disposal of non-taxable property<br/>RT - Royalties: Payment received for intellectual property usage<br/>TX - Tax Credit: Credit applied against tax liability<br/>BP - Buy Permitted: Security is eligible for purchase under corporate action<br/>CL - Cash in Lieu of Fractional Share: Cash payment for fractional share entitlement<br/>DF - Drop Fraction: Fractional shares are dropped without compensation<br/>EX - Extend and Retain Fractions: Fractional shares are retained and adjusted<br/>NC - Round to Nearest Cent: Monetary amounts rounded to the nearest cent<br/>NW - Round to Nearest Whole Number if .5 or above: Rounding rule for fractional shares<br/>PR - Purchase Required: Mandatory purchase of securities as part of corporate action<br/>RD - Round Down to Nearest Whole Number: Fractional shares rounded down<br/>RU - Round Up to Nearest Whole Number: Fractional shares rounded up<br/>SR - Sale Required: Mandatory sale of securities as part of corporate action<br/>BU - Round up: Generic rounding up rule<br/>BC - Beneficial Owner Cash in Lieu: Cash payment to beneficial owner for fractional share<br/>BD - Beneficial Owner Round Down: Fractional shares of beneficial owner rounded down<br/>CT - Security Convert To Cash: Conversion of security into cash (used for tokenized assets)",
                            "example": "DV",
                            "enum": [
                              "DV",
                              "FR",
                              "IN",
                              "L2",
                              "OT",
                              "C",
                              "PC",
                              "PM",
                              "S",
                              "ST",
                              "SI",
                              "UF",
                              "PI",
                              "TD",
                              "TE",
                              "FI",
                              "CD",
                              "CO",
                              "CC",
                              "CN",
                              "RT",
                              "TX",
                              "BP",
                              "CL",
                              "DF",
                              "EX",
                              "NC",
                              "NW",
                              "PR",
                              "RD",
                              "RU",
                              "SR",
                              "BU",
                              "BC",
                              "BD",
                              "CT"
                            ]
                          },
                          "pay_type": {
                            "type": "string",
                            "description": "Payout Delivery Type<br/>CASH - Cash settlement<br/>SECURITY - Security settlement (stock, right, warrant, etc.)<br/>SCRIP - Scrip dividend (dividend paid in shares instead of cash)<br/>SECURITY_AND_CASH - Combination of security and cash (logical type, not for persistence)",
                            "example": "CASH",
                            "enum": [
                              "CASH",
                              "SECURITY",
                              "SCRIP",
                              "SECURITY_AND_CASH"
                            ]
                          },
                          "payout_number": {
                            "type": "integer",
                            "description": "Sequence number of the payout within the same option. Used for ordering and identification.",
                            "format": "int32",
                            "example": 1
                          },
                          "adr_fee_rate": {
                            "type": "string",
                            "description": "ADR fee rate applied to this payout, if applicable.",
                            "example": "0.02"
                          },
                          "fraction_share_rule": {
                            "type": "string",
                            "description": "Fraction Share Rule<br/>NONE - No special handling for fractional shares<br/>ROUND_DOWN - Round down fractional shares<br/>ROUND_UP - Round up fractional shares<br/>CASH_IN_LIEU - Cash in lieu for fractional shares<br/>DISTRIBUTION - Fractional share distribution<br/>STANDARD - Standard rounding for fractional shares",
                            "example": "ROUND_DOWN",
                            "enum": [
                              "NONE",
                              "ROUND_DOWN",
                              "ROUND_UP",
                              "CASH_IN_LIEU",
                              "DISTRIBUTION",
                              "STANDARD"
                            ]
                          },
                          "cancellation_fee": {
                            "type": "string",
                            "description": "Cancellation fee rate applied if the corporate action is cancelled.",
                            "example": "0"
                          },
                          "issuance_fee": {
                            "type": "string",
                            "description": "Issuance fee rate applied for newly issued securities.",
                            "example": "0"
                          },
                          "tax_status": {
                            "type": "string",
                            "description": "IRS income classification for tax reporting purposes.",
                            "example": "0001"
                          },
                          "currency": {
                            "type": "string",
                            "description": "Currency of the cash payout. Applicable only when payType involves CASH.",
                            "example": "USD"
                          },
                          "amount": {
                            "type": "string",
                            "description": "Cash amount paid per share held. Applicable only when payType involves CASH.",
                            "example": "0.85"
                          },
                          "withholding_tax_rate": {
                            "type": "string",
                            "description": "Withholding tax rate applied to the cash payout.",
                            "example": "0.15"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Trading symbol of the distributed security.",
                            "example": "AAPL"
                          },
                          "name": {
                            "type": "string",
                            "description": "Name of the distributed security.",
                            "example": "Apple Inc."
                          },
                          "exchange": {
                            "type": "string",
                            "description": "Exchange where the distributed security is listed.",
                            "example": "CCC"
                          },
                          "from_ratio": {
                            "type": "string",
                            "description": "Original holding quantity used as the base for ratio calculation.",
                            "example": "10"
                          },
                          "to_ratio": {
                            "type": "string",
                            "description": "Distributed quantity received for the given base holding.",
                            "example": "1"
                          },
                          "cash_in_lieu_price": {
                            "type": "string",
                            "description": "Cash-in-lieu price used to settle fractional shares.",
                            "example": "125.3"
                          },
                          "reinvest_price": {
                            "type": "string",
                            "description": "Reinvestment price used for scrip dividend calculation.",
                            "example": "132.5"
                          }
                        },
                        "description": "Corporate Action Payout Detail. Represents a single payout rule applied per unit holding under a specific corporate action option.",
                        "title": "EventPayout"
                      }
                    }
                  },
                  "description": "Event To Info. Corporate action's target instrument information",
                  "title": "EventToInfo"
                }
              }
            },
            "description": "Corporate Action Result",
            "title": "CorporateActionResult"
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
    "name": "Get Corporate Actions Detail",
    "description": {
      "content": "Retrieves corporate actions detail information,only support stocks. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Corporate Actions Events](../fd-events/ca-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "instruments",
        "stocks",
        "corporate-actions",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Corporate Event Action ID",
            "type": "text/plain"
          },
          "key": "event_id",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-preview.md>

### Preview Order

Calculates the estimated cost and fees for an order based on the provided parameters. Supports simple orders. 

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/preview",
  "method": "post",
  "tags": [
    "Orders"
  ],
  "description": "Calculates the estimated cost and fees for an order based on the provided parameters. Supports simple orders. ",
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
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
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
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "combo_type": {
                    "type": "string",
                    "description": "Type of order combination.<br/> &bull; NORMAL: Indicates a standard single order",
                    "example": "NORMAL",
                    "enum": [
                      "NORMAL"
                    ]
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY",
                      "EVENT"
                    ]
                  },
                  "entrust_type": {
                    "type": "string",
                    "description": "Specifies the method for placing the order.<br/> &bull; QTY: Order specified by quantity of shares or units.<br/> &bull; AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK.",
                    "example": "QTY",
                    "enum": [
                      "QTY",
                      "AMOUNT"
                    ]
                  },
                  "support_trading_session": {
                    "type": "string",
                    "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only.<br/> &bull; NIGHT: Only supports night trading.<br/> &bull; ALL: Include extended trading hours.<br/> &bull; CORE: Only support regular trading hours.",
                    "example": "CORE",
                    "enum": [
                      "ALL",
                      "CORE",
                      "NIGHT"
                    ]
                  },
                  "symbol": {
                    "type": "string",
                    "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market.",
                    "example": "AAPL"
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
                    "description": "The order side indicating the intended trading direction of the transaction. <br/> The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).<br/> Event trading supports BUY and SELL sides only.<br/> Equity trading supports BUY, SELL, and SHORT sides.",
                    "example": "BUY",
                    "enum": [
                      "BUY",
                      "SELL",
                      "SHORT"
                    ]
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>For equity trading support all order types listed below.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>&nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/>&nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/>&nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/>For event trading: <br/>&nbsp; &bull; Supported order types: LIMIT.<br/> ",
                    "example": "MARKET",
                    "enum": [
                      "MARKET",
                      "LIMIT",
                      "STOP_LOSS",
                      "STOP_LOSS_LIMIT",
                      "TRAILING_STOP_LOSS"
                    ]
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK.<br/> U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC.<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.<br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days).<br/> &bull; IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled.<br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time.<br/> &bull; FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled.",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTC",
                      "IOC",
                      "GTD",
                      "FOK"
                    ]
                  },
                  "expire_date": {
                    "type": "string",
                    "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                    "example": "2026-12-01"
                  },
                  "stop_price": {
                    "type": "string",
                    "description": "Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT.<br/> Specifies the trigger price at which the stop order becomes active.",
                    "example": "11.0"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.<br/> When event_trade_mode is set for an event contract trade, limit_price is not required.",
                    "example": "11.0"
                  },
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "trailing_type": {
                    "type": "string",
                    "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. <br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                    "example": "AMOUNT",
                    "enum": [
                      "PERCENTAGE",
                      "AMOUNT"
                    ]
                  },
                  "trailing_stop_step": {
                    "type": "string",
                    "description": "Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%).",
                    "example": "1"
                  },
                  "event_outcome": {
                    "type": "string",
                    "description": "Event outcome decision, only applicable to event orders.",
                    "example": "yes",
                    "enum": [
                      "yes",
                      "no"
                    ]
                  },
                  "event_trade_mode": {
                    "type": "string",
                    "description": "Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored.<br/> &bull; TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price.<br/> &bull; TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field.",
                    "example": "TRADE_IN_AMOUNT",
                    "enum": [
                      "TRADE_IN_AMOUNT",
                      "TRADE_IN_CONTRACT"
                    ]
                  }
                },
                "description": "Order Details",
                "title": "OrderCommonPreviewItemParam"
              }
            }
          },
          "description": "Order Preview Request Body",
          "title": "OrderCommonPreviewParam"
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
                "description": "Estimated capital required for the order. <br/>The final capital usage may differ depending on the actual execution and fee settlement.",
                "example": "100"
              },
              "estimated_transaction_fee": {
                "type": "string",
                "description": "Estimated transaction fee for placing the order, including exchange, clearing, and commission fees. <br/> The actual fee may differ based on final execution.",
                "example": "1"
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
    "new_orders": [
      {
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "combo_type": "NORMAL",
        "instrument_type": "EQUITY",
        "entrust_type": "QTY",
        "support_trading_session": "CORE",
        "symbol": "AAPL",
        "market": "US",
        "side": "BUY",
        "order_type": "MARKET",
        "time_in_force": "DAY",
        "expire_date": "2026-12-01",
        "stop_price": "11.0",
        "limit_price": "11.0",
        "quantity": "1",
        "trailing_type": "AMOUNT",
        "trailing_stop_step": "1",
        "event_outcome": "yes",
        "event_trade_mode": "TRADE_IN_AMOUNT"
      }
    ]
  },
  "postman": {
    "name": "Preview Order",
    "description": {
      "content": "Calculates the estimated cost and fees for an order based on the provided parameters. Supports simple orders. ",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-place.md>

### Place Order

Places order for a specific account. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Trade Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/place",
  "method": "post",
  "tags": [
    "Orders"
  ],
  "description": "Places order for a specific account. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Trade Events](/apis/docs/reference/fd-events/trading-events#trade-events)",
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
            "new_orders"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
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
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "combo_type": {
                    "type": "string",
                    "description": "Type of order combination.<br/> &bull; NORMAL: Indicates a standard single order",
                    "example": "NORMAL",
                    "enum": [
                      "NORMAL"
                    ]
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY",
                      "EVENT"
                    ]
                  },
                  "entrust_type": {
                    "type": "string",
                    "description": "Specifies the method for placing the order.<br/> &bull; QTY: Order specified by quantity of shares or units.<br/> &bull; AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK.",
                    "example": "QTY",
                    "enum": [
                      "QTY",
                      "AMOUNT"
                    ]
                  },
                  "support_trading_session": {
                    "type": "string",
                    "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only.<br/> &bull; NIGHT: Only supports night trading.<br/> &bull; ALL: Include extended trading hours.<br/> &bull; CORE: Only support regular trading hours.",
                    "example": "CORE",
                    "enum": [
                      "ALL",
                      "CORE",
                      "NIGHT"
                    ]
                  },
                  "symbol": {
                    "type": "string",
                    "description": "Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market.",
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
                    "description": "The order side indicating the intended trading direction of the transaction. <br/> The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).<br/> Event trading supports BUY and SELL sides only.<br/> Equity trading supports BUY, SELL, and SHORT sides.",
                    "example": "BUY",
                    "enum": [
                      "BUY",
                      "SELL",
                      "SHORT"
                    ]
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>For equity trading support all order types listed below.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>&nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/>&nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/>&nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/>For event trading: <br/>&nbsp; &bull; Supported order types: LIMIT.<br/> ",
                    "example": "MARKET",
                    "enum": [
                      "MARKET",
                      "LIMIT",
                      "STOP_LOSS",
                      "STOP_LOSS_LIMIT",
                      "TRAILING_STOP_LOSS"
                    ]
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK.<br/> U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC.<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.<br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days).<br/> &bull; IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled.<br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time.<br/> &bull; FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled.",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTC",
                      "IOC",
                      "GTD",
                      "FOK"
                    ]
                  },
                  "expire_date": {
                    "type": "string",
                    "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                    "example": "2026-12-01"
                  },
                  "stop_price": {
                    "type": "string",
                    "description": "Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT.<br/> Specifies the trigger price at which the stop order becomes active.",
                    "example": "11.0"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.<br/> When event_trade_mode is set for an event contract trade, limit_price is not required.",
                    "example": "11.0"
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
                  "trailing_type": {
                    "type": "string",
                    "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. <br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                    "example": "AMOUNT",
                    "enum": [
                      "PERCENTAGE",
                      "AMOUNT"
                    ]
                  },
                  "trailing_stop_step": {
                    "type": "string",
                    "description": "Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%).",
                    "example": "1"
                  },
                  "event_outcome": {
                    "type": "string",
                    "description": "Event outcome decision, only applicable to event orders.",
                    "example": "yes",
                    "enum": [
                      "yes",
                      "no"
                    ]
                  },
                  "event_trade_mode": {
                    "type": "string",
                    "description": "Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored.<br/> &bull; TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price.<br/> &bull; TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field.",
                    "example": "TRADE_IN_AMOUNT",
                    "enum": [
                      "TRADE_IN_AMOUNT",
                      "TRADE_IN_CONTRACT"
                    ]
                  }
                },
                "description": "Order Details",
                "title": "OrderCommonPlaceItemParam"
              }
            }
          },
          "description": "Order Place Request Body",
          "title": "OrderCommonPlaceParam"
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
    "new_orders": [
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
        "expire_date": "2026-12-01",
        "stop_price": "11.0",
        "limit_price": "11.0",
        "quantity": "1",
        "total_cash_amount": "100.4",
        "trailing_type": "AMOUNT",
        "trailing_stop_step": "1",
        "event_outcome": "yes",
        "event_trade_mode": "TRADE_IN_AMOUNT"
      }
    ]
  },
  "postman": {
    "name": "Place Order",
    "description": {
      "content": "Places order for a specific account. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Trade Events](/apis/docs/reference/fd-events/trading-events#trade-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-replace.md>

### Replace Order

Modifies an existing order. Only the provided fields are updated; all other attributes remain unchanged. 

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/replace",
  "method": "post",
  "tags": [
    "Orders"
  ],
  "description": "Modifies an existing order. Only the provided fields are updated; all other attributes remain unchanged. ",
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
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Used to track or reference the order when interacting with the system.",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK.<br/> U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC.<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.<br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days).<br/> &bull; IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled.<br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time.<br/> &bull; FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled.",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTC",
                      "IOC",
                      "GTD",
                      "FOK"
                    ]
                  },
                  "stop_price": {
                    "type": "string",
                    "description": "Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT.<br/> Specifies the trigger price at which the stop order becomes active.",
                    "example": "11.0"
                  },
                  "limit_price": {
                    "type": "string",
                    "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT.<br/> Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
                    "example": "11.0"
                  },
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>For equity trading support all order types listed below.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>&nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/>&nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/>&nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/>For event trading: <br/>&nbsp; &bull; Supported order types: LIMIT.<br/> ",
                    "example": "MARKET",
                    "enum": [
                      "MARKET",
                      "LIMIT",
                      "STOP_LOSS",
                      "STOP_LOSS_LIMIT",
                      "TRAILING_STOP_LOSS"
                    ]
                  },
                  "trailing_stop_step": {
                    "type": "string",
                    "description": "Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%).",
                    "example": "1"
                  }
                },
                "description": "Order Details",
                "title": "OrderCommonReplaceItemParam"
              }
            }
          },
          "description": "Order Replace Request Body",
          "title": "OrderCommonReplaceParam"
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
    "modify_orders": [
      {
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "time_in_force": "DAY",
        "stop_price": "11.0",
        "limit_price": "11.0",
        "quantity": "1",
        "order_type": "MARKET",
        "trailing_stop_step": "1"
      }
    ]
  },
  "postman": {
    "name": "Replace Order",
    "description": {
      "content": "Modifies an existing order. Only the provided fields are updated; all other attributes remain unchanged. ",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-cancel.md>

### Cancel Order

Cancels a previously submitted order that has not yet been fully filled. Only orders in open status can be cancelled.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/cancel",
  "method": "post",
  "tags": [
    "Orders"
  ],
  "description": "Cancels a previously submitted order that has not yet been fully filled. Only orders in open status can be cancelled.",
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
            "client_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters and must be unique per account.<br/> Used to track or reference the order when interacting with the system.",
              "example": "0KGOHL4PR2SLC0DKIND4TI0002"
            },
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            }
          },
          "description": "Order Cancel Request Body",
          "title": "OrderCommonCancelParam"
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
    "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
    "account_id": "93IUJ28O9VO2KBGHDHR4H9"
  },
  "postman": {
    "name": "Cancel Order",
    "description": {
      "content": "Cancels a previously submitted order that has not yet been fully filled. Only orders in open status can be cancelled.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/order-open.md>

### List Open Orders

Retrieves pending orders by page. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/open-orders/list",
  "method": "get",
  "tags": [
    "Orders"
  ],
  "description": "Retrieves pending orders by page. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.",
  "operationId": "orderOpen",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "93IUJ28O9VO2KBGHDHR4H9"
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
                    "combo_type",
                    "orders"
                  ],
                  "type": "object",
                  "properties": {
                    "client_order_id": {
                      "type": "string",
                      "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order.",
                      "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                    },
                    "combo_type": {
                      "type": "string",
                      "description": "Type of order combination.<br/> &bull; NORMAL: Indicates a standard single order",
                      "example": "NORMAL"
                    },
                    "orders": {
                      "type": "array",
                      "description": "Order Details",
                      "items": {
                        "required": [
                          "client_order_id",
                          "entrust_type",
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
                            "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order.",
                            "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                          },
                          "order_id": {
                            "type": "string",
                            "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system.",
                            "example": "0352U72LQI6DT0KF41GK000000"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market.",
                            "example": "AAPL"
                          },
                          "side": {
                            "type": "string",
                            "description": "The order side indicating the intended trading direction of the transaction. <br/> The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).<br/> Event trading supports BUY and SELL sides only.<br/> Equity trading supports BUY, SELL, and SHORT sides.",
                            "example": "BUY",
                            "enum": [
                              "BUY",
                              "SELL",
                              "SHORT"
                            ]
                          },
                          "status": {
                            "type": "string",
                            "description": "&bull; PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion<br/> &bull; SUBMITTED: Indicates that the order has been submitted to the exchange or webull<br/> &bull; CANCELLED: Indicates that the order has been successfully cancelled<br/> &bull; FILLED: Indicates that the order has been fully executed<br/> &bull; FAILED: Indicates a failed order, such as REJECTED<br/> &bull; PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed",
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
                            "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>For equity trading support all order types listed below.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>&nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/>&nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/>&nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/>For event trading: <br/>&nbsp; &bull; Supported order types: LIMIT.<br/> ",
                            "example": "MARKET",
                            "enum": [
                              "MARKET",
                              "LIMIT",
                              "STOP_LOSS",
                              "STOP_LOSS_LIMIT",
                              "TRAILING_STOP_LOSS"
                            ]
                          },
                          "instrument_type": {
                            "type": "string",
                            "description": "Type of financial instrument associated with the request.",
                            "example": "EQUITY",
                            "enum": [
                              "EQUITY",
                              "EVENT"
                            ]
                          },
                          "support_trading_session": {
                            "type": "string",
                            "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only.<br/> &bull; NIGHT: Only supports night trading.<br/> &bull; ALL: Include extended trading hours.<br/> &bull; CORE: Only support regular trading hours.",
                            "example": "CORE",
                            "enum": [
                              "ALL",
                              "CORE",
                              "NIGHT"
                            ]
                          },
                          "entrust_type": {
                            "type": "string",
                            "description": "Specifies the method for placing the order.<br/> &bull; QTY: Order specified by quantity of shares or units.<br/> &bull; AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK.",
                            "example": "QTY",
                            "enum": [
                              "QTY",
                              "AMOUNT"
                            ]
                          },
                          "time_in_force": {
                            "type": "string",
                            "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK.<br/> U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC.<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.<br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days).<br/> &bull; IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled.<br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time.<br/> &bull; FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled.",
                            "example": "DAY",
                            "enum": [
                              "DAY",
                              "GTC",
                              "IOC",
                              "GTD",
                              "FOK"
                            ]
                          },
                          "expire_date": {
                            "type": "string",
                            "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                            "example": "2026-12-01"
                          },
                          "total_cash_amount": {
                            "type": "string",
                            "description": "The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount.",
                            "example": "100.4"
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
                            "description": "Limit Price",
                            "example": "11.0"
                          },
                          "stop_price": {
                            "type": "string",
                            "description": "Stop Price",
                            "example": "11.0"
                          },
                          "trailing_type": {
                            "type": "string",
                            "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. <br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                            "example": "AMOUNT",
                            "enum": [
                              "PERCENTAGE",
                              "AMOUNT"
                            ]
                          },
                          "trailing_stop_step": {
                            "type": "string",
                            "description": "Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%).",
                            "example": "1"
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
                          "event_outcome": {
                            "type": "string",
                            "description": "Event outcome decision, only applicable to event orders.",
                            "example": "yes",
                            "enum": [
                              "yes",
                              "no"
                            ]
                          },
                          "event_trade_mode": {
                            "type": "string",
                            "description": "Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored.<br/> &bull; TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price.<br/> &bull; TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field.",
                            "example": "TRADE_IN_AMOUNT",
                            "enum": [
                              "TRADE_IN_AMOUNT",
                              "TRADE_IN_CONTRACT"
                            ]
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
      "content": "Retrieves pending orders by page. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "(Required) Account identifier",
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/order-detail.md>

### Get Order Detail

Retrieves the specified order details through the order ID or client order ID.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/get",
  "method": "get",
  "tags": [
    "Orders"
  ],
  "description": "Retrieves the specified order details through the order ID or client order ID.",
  "operationId": "orderDetail",
  "parameters": [
    {
      "name": "client_order_id",
      "in": "query",
      "description": "The last order ID returned from the previous response.<br/> Used for cursor-based pagination.<br/> Not required for the first page query.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "0KGOHL4PR2SLC0DKIND4TI0002"
    },
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "93IUJ28O9VO2KBGHDHR4H9"
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
              "combo_type",
              "orders"
            ],
            "type": "object",
            "properties": {
              "client_order_id": {
                "type": "string",
                "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order.",
                "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
              },
              "combo_type": {
                "type": "string",
                "description": "Type of order combination.<br/> &bull; NORMAL: Indicates a standard single order",
                "example": "NORMAL"
              },
              "orders": {
                "type": "array",
                "description": "Order Details",
                "items": {
                  "required": [
                    "client_order_id",
                    "entrust_type",
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
                      "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order.",
                      "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                    },
                    "order_id": {
                      "type": "string",
                      "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system.",
                      "example": "0352U72LQI6DT0KF41GK000000"
                    },
                    "symbol": {
                      "type": "string",
                      "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market.",
                      "example": "AAPL"
                    },
                    "side": {
                      "type": "string",
                      "description": "The order side indicating the intended trading direction of the transaction. <br/> The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).<br/> Event trading supports BUY and SELL sides only.<br/> Equity trading supports BUY, SELL, and SHORT sides.",
                      "example": "BUY",
                      "enum": [
                        "BUY",
                        "SELL",
                        "SHORT"
                      ]
                    },
                    "status": {
                      "type": "string",
                      "description": "&bull; PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion<br/> &bull; SUBMITTED: Indicates that the order has been submitted to the exchange or webull<br/> &bull; CANCELLED: Indicates that the order has been successfully cancelled<br/> &bull; FILLED: Indicates that the order has been fully executed<br/> &bull; FAILED: Indicates a failed order, such as REJECTED<br/> &bull; PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed",
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
                      "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>For equity trading support all order types listed below.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>&nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/>&nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/>&nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/>For event trading: <br/>&nbsp; &bull; Supported order types: LIMIT.<br/> ",
                      "example": "MARKET",
                      "enum": [
                        "MARKET",
                        "LIMIT",
                        "STOP_LOSS",
                        "STOP_LOSS_LIMIT",
                        "TRAILING_STOP_LOSS"
                      ]
                    },
                    "instrument_type": {
                      "type": "string",
                      "description": "Type of financial instrument associated with the request.",
                      "example": "EQUITY",
                      "enum": [
                        "EQUITY",
                        "EVENT"
                      ]
                    },
                    "support_trading_session": {
                      "type": "string",
                      "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only.<br/> &bull; NIGHT: Only supports night trading.<br/> &bull; ALL: Include extended trading hours.<br/> &bull; CORE: Only support regular trading hours.",
                      "example": "CORE",
                      "enum": [
                        "ALL",
                        "CORE",
                        "NIGHT"
                      ]
                    },
                    "entrust_type": {
                      "type": "string",
                      "description": "Specifies the method for placing the order.<br/> &bull; QTY: Order specified by quantity of shares or units.<br/> &bull; AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK.",
                      "example": "QTY",
                      "enum": [
                        "QTY",
                        "AMOUNT"
                      ]
                    },
                    "time_in_force": {
                      "type": "string",
                      "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK.<br/> U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC.<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.<br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days).<br/> &bull; IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled.<br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time.<br/> &bull; FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled.",
                      "example": "DAY",
                      "enum": [
                        "DAY",
                        "GTC",
                        "IOC",
                        "GTD",
                        "FOK"
                      ]
                    },
                    "expire_date": {
                      "type": "string",
                      "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                      "example": "2026-12-01"
                    },
                    "total_cash_amount": {
                      "type": "string",
                      "description": "The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount.",
                      "example": "100.4"
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
                      "description": "Limit Price",
                      "example": "11.0"
                    },
                    "stop_price": {
                      "type": "string",
                      "description": "Stop Price",
                      "example": "11.0"
                    },
                    "trailing_type": {
                      "type": "string",
                      "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. <br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                      "example": "AMOUNT",
                      "enum": [
                        "PERCENTAGE",
                        "AMOUNT"
                      ]
                    },
                    "trailing_stop_step": {
                      "type": "string",
                      "description": "Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%).",
                      "example": "1"
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
                    "event_outcome": {
                      "type": "string",
                      "description": "Event outcome decision, only applicable to event orders.",
                      "example": "yes",
                      "enum": [
                        "yes",
                        "no"
                      ]
                    },
                    "event_trade_mode": {
                      "type": "string",
                      "description": "Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored.<br/> &bull; TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price.<br/> &bull; TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field.",
                      "example": "TRADE_IN_AMOUNT",
                      "enum": [
                        "TRADE_IN_AMOUNT",
                        "TRADE_IN_CONTRACT"
                      ]
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
            "description": "Result data list",
            "title": "OrderListResult"
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
      "content": "Retrieves the specified order details through the order ID or client order ID.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "(Required) The last order ID returned from the previous response.<br/> Used for cursor-based pagination.<br/> Not required for the first page query.",
            "type": "text/plain"
          },
          "key": "client_order_id",
          "value": ""
        },
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

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/order-history.md>

### List Order History

Retrieves historical orders within a specified date range (default: last 7 days). If orders are group orders, they will be returned together, and the number of orders returned on one page may exceed the page_size. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/orders/historical-orders/list",
  "method": "get",
  "tags": [
    "Orders"
  ],
  "description": "Retrieves historical orders within a specified date range (default: last 7 days). If orders are group orders, they will be returned together, and the number of orders returned on one page may exceed the page_size. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.",
  "operationId": "orderHistory",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Account identifier",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "93IUJ28O9VO2KBGHDHR4H9"
    },
    {
      "name": "start_date",
      "in": "query",
      "description": "The start date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 2 years.<br/> Format: yyyy-MM-dd.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-09-25"
    },
    {
      "name": "end_date",
      "in": "query",
      "description": "The end date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 2 years.<br/> Format: yyyy-MM-dd.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-10-25"
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
                    "combo_type",
                    "orders"
                  ],
                  "type": "object",
                  "properties": {
                    "client_order_id": {
                      "type": "string",
                      "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order.",
                      "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                    },
                    "combo_type": {
                      "type": "string",
                      "description": "Type of order combination.<br/> &bull; NORMAL: Indicates a standard single order",
                      "example": "NORMAL"
                    },
                    "orders": {
                      "type": "array",
                      "description": "Order Details",
                      "items": {
                        "required": [
                          "client_order_id",
                          "entrust_type",
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
                            "description": "Client-defined order identifier. Returned in the response for simple orders.<br/> Represents the unique order ID assigned by the user when placing the order.",
                            "example": "THI82O5JB7MQ2K76LL5FSDS2CB"
                          },
                          "order_id": {
                            "type": "string",
                            "description": "System-generated order identifier. Returned in the response for simple orders.<br/> Represents the unique Webull order ID assigned by the system.",
                            "example": "0352U72LQI6DT0KF41GK000000"
                          },
                          "symbol": {
                            "type": "string",
                            "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market.",
                            "example": "AAPL"
                          },
                          "side": {
                            "type": "string",
                            "description": "The order side indicating the intended trading direction of the transaction. <br/> The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash).<br/> Event trading supports BUY and SELL sides only.<br/> Equity trading supports BUY, SELL, and SHORT sides.",
                            "example": "BUY",
                            "enum": [
                              "BUY",
                              "SELL",
                              "SHORT"
                            ]
                          },
                          "status": {
                            "type": "string",
                            "description": "&bull; PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion<br/> &bull; SUBMITTED: Indicates that the order has been submitted to the exchange or webull<br/> &bull; CANCELLED: Indicates that the order has been successfully cancelled<br/> &bull; FILLED: Indicates that the order has been fully executed<br/> &bull; FAILED: Indicates a failed order, such as REJECTED<br/> &bull; PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed",
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
                            "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/>Available order types depend on the market and instrument type.<br/>For equity trading support all order types listed below.<br/>&nbsp; &bull; <b>LIMIT:</b> Limit Order<br/>&nbsp; &bull; <b>MARKET:</b> Market Order<br/>&nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/>&nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/>&nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/>For event trading: <br/>&nbsp; &bull; Supported order types: LIMIT.<br/> ",
                            "example": "MARKET",
                            "enum": [
                              "MARKET",
                              "LIMIT",
                              "STOP_LOSS",
                              "STOP_LOSS_LIMIT",
                              "TRAILING_STOP_LOSS"
                            ]
                          },
                          "instrument_type": {
                            "type": "string",
                            "description": "Type of financial instrument associated with the request.",
                            "example": "EQUITY",
                            "enum": [
                              "EQUITY",
                              "EVENT"
                            ]
                          },
                          "support_trading_session": {
                            "type": "string",
                            "description": "Specifies the trading session for the order. Applicable to U.S. stock market orders only.<br/> &bull; NIGHT: Only supports night trading.<br/> &bull; ALL: Include extended trading hours.<br/> &bull; CORE: Only support regular trading hours.",
                            "example": "CORE",
                            "enum": [
                              "ALL",
                              "CORE",
                              "NIGHT"
                            ]
                          },
                          "entrust_type": {
                            "type": "string",
                            "description": "Specifies the method for placing the order.<br/> &bull; QTY: Order specified by quantity of shares or units.<br/> &bull; AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK.",
                            "example": "QTY",
                            "enum": [
                              "QTY",
                              "AMOUNT"
                            ]
                          },
                          "time_in_force": {
                            "type": "string",
                            "description": "Specifies the duration for which the order remains active in the market (Time-In-Force).<br/> Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK.<br/> U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC.<br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day.<br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days).<br/> &bull; IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled.<br/> &bull; GTD: order that will automatically expire and be cancelled at a specific future date and time.<br/> &bull; FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled.",
                            "example": "DAY",
                            "enum": [
                              "DAY",
                              "GTC",
                              "IOC",
                              "GTD",
                              "FOK"
                            ]
                          },
                          "expire_date": {
                            "type": "string",
                            "description": "GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format",
                            "example": "2026-12-01"
                          },
                          "total_cash_amount": {
                            "type": "string",
                            "description": "The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount.",
                            "example": "100.4"
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
                            "description": "Limit Price",
                            "example": "11.0"
                          },
                          "stop_price": {
                            "type": "string",
                            "description": "Stop Price",
                            "example": "11.0"
                          },
                          "trailing_type": {
                            "type": "string",
                            "description": "When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. <br/> &bull; AMOUNT: By amount. <br/> &bull; PERCENTAGE: By percentage.",
                            "example": "AMOUNT",
                            "enum": [
                              "PERCENTAGE",
                              "AMOUNT"
                            ]
                          },
                          "trailing_stop_step": {
                            "type": "string",
                            "description": "Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%).",
                            "example": "1"
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
                          "event_outcome": {
                            "type": "string",
                            "description": "Event outcome decision, only applicable to event orders.",
                            "example": "yes",
                            "enum": [
                              "yes",
                              "no"
                            ]
                          },
                          "event_trade_mode": {
                            "type": "string",
                            "description": "Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored.<br/> &bull; TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price.<br/> &bull; TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field.",
                            "example": "TRADE_IN_AMOUNT",
                            "enum": [
                              "TRADE_IN_AMOUNT",
                              "TRADE_IN_CONTRACT"
                            ]
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
      "content": "Retrieves historical orders within a specified date range (default: last 7 days). If orders are group orders, they will be returned together, and the number of orders returned on one page may exceed the page_size. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
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
            "content": "(Required) Account identifier",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The start date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 2 years.<br/> Format: yyyy-MM-dd.",
            "type": "text/plain"
          },
          "key": "start_date",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The end date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 2 years.<br/> Format: yyyy-MM-dd.",
            "type": "text/plain"
          },
          "key": "end_date",
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

## Create Cash Journal

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-journal-cash-create.md>

### Create Cash Journal

Creates a cash journal request. Cash journal supports both proprietary account (PAB) level and end client account level. Journals are supported only between accounts with the same owner.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Cash Journal Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/journals/cash-journals/create",
  "method": "post",
  "tags": [
    "Journals"
  ],
  "description": "Creates a cash journal request. Cash journal supports both proprietary account (PAB) level and end client account level. Journals are supported only between accounts with the same owner.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Cash Journal Events](/apis/docs/reference/fd-events/journal-events#cash-journal-events)",
  "operationId": "brokerJournalCashCreate",
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
            "amount",
            "client_request_id",
            "currency",
            "from_account",
            "to_account"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "from_account": {
              "type": "string",
              "description": "From account id.",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "to_account": {
              "type": "string",
              "description": "To account id.",
              "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
            },
            "amount": {
              "type": "string",
              "description": "Amount",
              "example": "100"
            },
            "currency": {
              "type": "string",
              "description": "Currency",
              "example": "USD",
              "enum": [
                "USD"
              ]
            }
          },
          "title": "BrokerCreateCashJournalParam"
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
              "amount",
              "client_request_id",
              "create_time",
              "currency",
              "from_account",
              "journal_id",
              "status",
              "to_account",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "from_account": {
                "type": "string",
                "description": "From account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "to_account": {
                "type": "string",
                "description": "To account id.",
                "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
              },
              "journal_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "BrokerCashJournalResult"
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
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "from_account": "93IUJ28O9VO2KBGHDHR4H9",
    "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
    "amount": "100",
    "currency": "USD"
  },
  "postman": {
    "name": "Create Cash Journal",
    "description": {
      "content": "Creates a cash journal request. Cash journal supports both proprietary account (PAB) level and end client account level. Journals are supported only between accounts with the same owner.  <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Cash Journal Events](/apis/docs/reference/fd-events/journal-events#cash-journal-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "journals",
        "cash-journals",
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

## Cash Journal Detail

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-journal-cash-query.md>

### Get Cash Journal Detail

Retrieves details information for cash journal request.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/journals/cash-journals/get",
  "method": "get",
  "tags": [
    "Journals"
  ],
  "description": "Retrieves details information for cash journal request.",
  "operationId": "brokerJournalCashQuery",
  "parameters": [
    {
      "name": "client_request_id",
      "in": "query",
      "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "LJIS16BACHQG9LPP44L9IQHGAB"
    },
    {
      "name": "account_id",
      "in": "query",
      "description": "From account identifier.",
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
              "amount",
              "client_request_id",
              "create_time",
              "currency",
              "from_account",
              "journal_id",
              "status",
              "to_account",
              "update_time"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "from_account": {
                "type": "string",
                "description": "From account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "to_account": {
                "type": "string",
                "description": "To account id.",
                "example": "41IO9QG4M5O65B0EA4LSJ4UJ99"
              },
              "journal_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "amount": {
                "type": "string",
                "description": "Amount",
                "example": "100"
              },
              "currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "USD"
                ]
              },
              "status": {
                "type": "string",
                "description": "Request status",
                "example": "SUBMITTED",
                "enum": [
                  "SUBMITTED",
                  "CANCELED",
                  "REJECTED",
                  "FAILED",
                  "COMPLETED"
                ]
              },
              "reason": {
                "type": "string",
                "description": "Reason. If the terminal state is not “COMPLETED”, return the reason.",
                "example": "Failed."
              },
              "create_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              },
              "update_time": {
                "type": "string",
                "description": "UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ.",
                "example": "2025-01-05T22:59:59.012Z"
              }
            },
            "title": "BrokerCashJournalResult"
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
    "name": "Get Cash Journal Detail",
    "description": {
      "content": "Retrieves details information for cash journal request.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "journals",
        "cash-journals",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
            "type": "text/plain"
          },
          "key": "client_request_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) From account identifier.",
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

## List Enums

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/list-enums.md>

### List Enums

Get master data by data type and optional parent code.
Master Data definition, hierarchy, and code description:

| data_type | description |
|-----------|-------------|
| ID_DOC_TYPE | Identification document type |
| MARITAL_STATUS | Marital status options |
| EMPLOYMENT_TYPE | Type of employment |
| OCCUPATION | Specific occupation under the employment type |
| POSITION | Job position or title |
| LIQUIDITY_NEEDS | Liquidity needs options |
| ANNUAL_INCOME | Annual income range |
| TOTAL_NET_WORTH | Estimated total net worth range |
| LIQUID_NET_WORTH | Liquid net worth range |
| INVESTMENT_KNOWLEDGE | Investment knowledge level (old) |
| INVESTMENT_OBJECTIVE | Investment goal or objective |
| TIME_HORIZON | Investment time horizon |
| INVESTMENT_EXPERIENCE | Experience level in investments |
| INVESTMENT_KNOWLEDGE_V2 | Investment knowledge level |
| TRADE_PER_YEAR | Estimated number of trades per year |
| TAX_TYPE | Tax classification type |
| AGREEMENT_TYPE | Agreement type |
| AGREEMENT_CONTENT_TYPE | Agreement content type |
| AGREEMENT_SIGN_METHOD | Agreement signing method |

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/master-data/enums/list",
  "method": "get",
  "tags": [
    "Master Data"
  ],
  "description": "Get master data by data type and optional parent code.\nMaster Data definition, hierarchy, and code description:\n\n| data_type | description |\n|-----------|-------------|\n| ID_DOC_TYPE | Identification document type |\n| MARITAL_STATUS | Marital status options |\n| EMPLOYMENT_TYPE | Type of employment |\n| OCCUPATION | Specific occupation under the employment type |\n| POSITION | Job position or title |\n| LIQUIDITY_NEEDS | Liquidity needs options |\n| ANNUAL_INCOME | Annual income range |\n| TOTAL_NET_WORTH | Estimated total net worth range |\n| LIQUID_NET_WORTH | Liquid net worth range |\n| INVESTMENT_KNOWLEDGE | Investment knowledge level (old) |\n| INVESTMENT_OBJECTIVE | Investment goal or objective |\n| TIME_HORIZON | Investment time horizon |\n| INVESTMENT_EXPERIENCE | Experience level in investments |\n| INVESTMENT_KNOWLEDGE_V2 | Investment knowledge level |\n| TRADE_PER_YEAR | Estimated number of trades per year |\n| TAX_TYPE | Tax classification type |\n| AGREEMENT_TYPE | Agreement type |\n| AGREEMENT_CONTENT_TYPE | Agreement content type |\n| AGREEMENT_SIGN_METHOD | Agreement signing method |",
  "operationId": "listEnums",
  "parameters": [
    {
      "name": "data_type",
      "in": "query",
      "description": "Type of master data",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Enumeration of master data types used for account opening, profile management, and questionnaire configuration.",
        "enum": [
          "ID_DOC_TYPE",
          "MARITAL_STATUS",
          "EMPLOYMENT_TYPE",
          "OCCUPATION",
          "POSITION",
          "LIQUIDITY_NEEDS",
          "ANNUAL_INCOME",
          "TOTAL_NET_WORTH",
          "LIQUID_NET_WORTH",
          "INVESTMENT_KNOWLEDGE_OLD",
          "INVESTMENT_OBJECTIVE",
          "TIME_HORIZON",
          "INVESTMENT_EXPERIENCE",
          "INVESTMENT_KNOWLEDGE",
          "TRADE_PER_YEAR",
          "TAX_TYPE",
          "AGREEMENT_CONTENT_TYPE",
          "AGREEMENT_SIGN_METHOD",
          "COUNTRY",
          "STATE",
          "CITY",
          "CONTACT_RELATIONSHIP"
        ]
      },
      "example": "AGREEMENT_TYPE"
    },
    {
      "name": "parent_code",
      "in": "query",
      "description": "Parent code of the master data, used for hierarchical queries",
      "required": false,
      "schema": {
        "type": "string"
      },
      "example": 121
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
                "code": {
                  "type": "string",
                  "description": "Business code of the master data item, typically used as a stable reference value.",
                  "example": "LIMITED"
                },
                "name": {
                  "type": "string",
                  "description": "Display name of the master data item in the default language.",
                  "example": "Limited"
                },
                "parent_code": {
                  "type": "string",
                  "description": "Parent code of the master data item, used to represent hierarchical relationships.",
                  "example": "121"
                }
              },
              "description": "Response object representing a master data item, commonly used for dictionaries, reference data, or configurable options.",
              "title": "MasterDataResp"
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
    "name": "List Enums",
    "description": {
      "content": "Get master data by data type and optional parent code.\nMaster Data definition, hierarchy, and code description:\n\n| data_type | description |\n|-----------|-------------|\n| ID_DOC_TYPE | Identification document type |\n| MARITAL_STATUS | Marital status options |\n| EMPLOYMENT_TYPE | Type of employment |\n| OCCUPATION | Specific occupation under the employment type |\n| POSITION | Job position or title |\n| LIQUIDITY_NEEDS | Liquidity needs options |\n| ANNUAL_INCOME | Annual income range |\n| TOTAL_NET_WORTH | Estimated total net worth range |\n| LIQUID_NET_WORTH | Liquid net worth range |\n| INVESTMENT_KNOWLEDGE | Investment knowledge level (old) |\n| INVESTMENT_OBJECTIVE | Investment goal or objective |\n| TIME_HORIZON | Investment time horizon |\n| INVESTMENT_EXPERIENCE | Experience level in investments |\n| INVESTMENT_KNOWLEDGE_V2 | Investment knowledge level |\n| TRADE_PER_YEAR | Estimated number of trades per year |\n| TAX_TYPE | Tax classification type |\n| AGREEMENT_TYPE | Agreement type |\n| AGREEMENT_CONTENT_TYPE | Agreement content type |\n| AGREEMENT_SIGN_METHOD | Agreement signing method |",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "master-data",
        "enums",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Type of master data",
            "type": "text/plain"
          },
          "key": "data_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Parent code of the master data, used for hierarchical queries",
            "type": "text/plain"
          },
          "key": "parent_code",
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

## Trade Calendar

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/list-trade-calendar.md>

### List Trade Calendar

Retrieves trading and settlement calendar for the specified country and market. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: Trade Calendar Events

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/master-data/trading-calendars/list",
  "method": "get",
  "tags": [
    "Master Data"
  ],
  "description": "Retrieves trading and settlement calendar for the specified country and market. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Trade Calendar Events](/apis/docs/reference/fd-events/broker-master-data-events#trade-calendar-events)",
  "operationId": "listTradeCalendar",
  "parameters": [
    {
      "name": "market",
      "in": "query",
      "description": "Country code, US etc.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
        "enum": [
          "US"
        ]
      }
    },
    {
      "name": "instrument_type",
      "in": "query",
      "description": "Type of financial instrument associated with the request. ",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Type of financial instrument associated with the request.",
        "enum": [
          "EQUITY",
          "EVENT"
        ]
      }
    },
    {
      "name": "year",
      "in": "query",
      "description": "Query Year in YYYY format",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": 2025
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
                "date",
                "is_settlement_day",
                "is_trading_day"
              ],
              "type": "object",
              "properties": {
                "date": {
                  "type": "string",
                  "description": "Date string in YYYY-MM-DD format.",
                  "example": "2025-01-06"
                },
                "settlement_date": {
                  "type": "string",
                  "description": "Settlement Date string in YYYY-MM-DD format.",
                  "example": "2025-01-08"
                },
                "is_trading_day": {
                  "type": "boolean",
                  "description": "Is Trading Day",
                  "example": true
                },
                "is_early_close": {
                  "type": "boolean",
                  "description": "Whether it is a half-day trading session",
                  "example": true
                },
                "is_settlement_day": {
                  "type": "boolean",
                  "description": "Is Settlement Day",
                  "example": true
                }
              },
              "title": "BrokerCalendarResult"
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
    "name": "List Trade Calendar",
    "description": {
      "content": "Retrieves trading and settlement calendar for the specified country and market. <br/><br/>Related Event Notifications: <br/>Please refer to the Event Push API documentation: [Trade Calendar Events](/apis/docs/reference/fd-events/broker-master-data-events#trade-calendar-events)",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "master-data",
        "trading-calendars",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Country code, US etc.",
            "type": "text/plain"
          },
          "key": "market",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Type of financial instrument associated with the request. ",
            "type": "text/plain"
          },
          "key": "instrument_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Query Year in YYYY format",
            "type": "text/plain"
          },
          "key": "year",
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

## List Agreements

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-list-agreements-by-type.md>

### List Agreements by Type

Helps retrieve all available agreement templates of a specific type. These templates can be presented to users for signing during account setup or feature activation. Note: This endpoint returns agreement templates, not user-specific agreement status. To check if a user has signed an agreement, use the account-specific endpoints..

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/agreements/list",
  "method": "get",
  "tags": [
    "Agreements"
  ],
  "description": "Helps retrieve all available agreement templates of a specific type. These templates can be presented to users for signing during account setup or feature activation. Note: This endpoint returns agreement templates, not user-specific agreement status. To check if a user has signed an agreement, use the account-specific endpoints..",
  "operationId": "brokerListAgreementsByType",
  "parameters": [
    {
      "name": "lang",
      "in": "query",
      "description": "Language parameter, default is EN <br/>EN: English <br/>",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "EN"
    },
    {
      "name": "agreement_type",
      "in": "query",
      "description": "Agreement type parameter. Value must be obtained from the Master Data API with data_type = AGREEMENT_TYPE.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "EVENT_CONTRACT_AGREEMENT"
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
                "id",
                "name",
                "version"
              ],
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "description": "Agreement ID",
                  "example": "JPGSRHN6QCAI68SEO24UL316TB_1"
                },
                "version": {
                  "type": "string",
                  "description": "Agreement Version",
                  "example": "1.0"
                },
                "name": {
                  "type": "string",
                  "description": "Agreement Name",
                  "example": "Customer Agreement"
                }
              },
              "description": "Agreement List Result",
              "title": "AgreementListResult"
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
    "name": "List Agreements by Type",
    "description": {
      "content": "Helps retrieve all available agreement templates of a specific type. These templates can be presented to users for signing during account setup or feature activation. Note: This endpoint returns agreement templates, not user-specific agreement status. To check if a user has signed an agreement, use the account-specific endpoints..",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "agreements",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "Language parameter, default is EN <br/>EN: English <br/>",
            "type": "text/plain"
          },
          "key": "lang",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Agreement type parameter. Value must be obtained from the Master Data API with data_type = AGREEMENT_TYPE.",
            "type": "text/plain"
          },
          "key": "agreement_type",
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

## Agreement Details

> Source: <https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-get-agreement-details.md>

### Get Agreement Detail

Retrieves detailed information about a specific agreement.

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
      "url": "https://broker-api.sandbox.webull.com"
    }
  ],
  "path": "/broker/agreements/get",
  "method": "get",
  "tags": [
    "Agreements"
  ],
  "description": "Retrieves detailed information about a specific agreement.",
  "operationId": "brokerGetAgreementDetails",
  "parameters": [
    {
      "name": "ids",
      "in": "query",
      "description": "The unique identifier of the agreement to retrieve details for.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "JPGSRHN6QCAI68SEO24UL316TB_1,JPGSRHN6QCAI68SEO24UL316TB_2"
    },
    {
      "name": "lang",
      "in": "query",
      "description": "Language parameter <br/>EN: English <br/>",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "EN"
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
                "content",
                "content_type",
                "id",
                "name"
              ],
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "description": "Agreement ID",
                  "example": "JPGSRHN6QCAI68SEO24UL316TB_1"
                },
                "name": {
                  "type": "string",
                  "description": "Agreement Name",
                  "example": "Customer Agreement"
                },
                "description": {
                  "type": "string",
                  "description": "Agreement Description",
                  "example": "This agreement outlines the terms and conditions for customers."
                },
                "content": {
                  "type": "string",
                  "description": "Agreement Content",
                  "example": "Full text of the agreement goes here..."
                },
                "content_type": {
                  "type": "string",
                  "description": "Agreement content type",
                  "example": "TEXT",
                  "enum": [
                    "TEXT",
                    "URL",
                    "PDF"
                  ]
                }
              },
              "description": "Agreement Details",
              "title": "AgreementDetails"
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
    "name": "Get Agreement Detail",
    "description": {
      "content": "Retrieves detailed information about a specific agreement.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "agreements",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) The unique identifier of the agreement to retrieve details for.",
            "type": "text/plain"
          },
          "key": "ids",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Language parameter <br/>EN: English <br/>",
            "type": "text/plain"
          },
          "key": "lang",
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


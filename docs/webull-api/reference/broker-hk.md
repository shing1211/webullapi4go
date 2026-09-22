# Broker API — HK — Verbatim Reference

> Institutional Broker API for Hong Kong. Uses the `DoBroker` transport. The HK sandbox returns `401 ROUTE_NOT_PERMITTED` (missing app scope).

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Create Virtual Account

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-create.md>

### Create Virtual Account

Creates a Virtual Account for a new customer.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/accounts/virtual-accounts/create",
  "method": "post",
  "tags": [
    "Account (ND)"
  ],
  "description": "Creates a Virtual Account for a new customer.",
  "operationId": "brokerAccountCreate",
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
            "account_type",
            "belong_account_id",
            "client_request_id",
            "trading_permissions"
          ],
          "type": "object",
          "properties": {
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "LJIS16BACHQG9LPP44L9IQHGAB"
            },
            "belong_account_id": {
              "type": "string",
              "description": "Belonging Account ID, the master account under which the VA account is created",
              "example": "TUHO91TNJ6B870VJ2SUL65U4LB"
            },
            "external_account_number": {
              "type": "string",
              "description": "The account number of the end customer in your own system, used as the client-defined segment of the virtual account number. The final account_number is composed as: V + Entity Code + this value, where the Entity Code is assigned by Webull during onboarding.<br/>Example: Entity Code ABCD + 10000001 → VABCD10000001.<br/>• Length: 1-10 characters.<br/>• Allowed characters: ASCII digits (0-9), English letters (a-z, A-Z), hyphen (-) and underscore (_).<br/>• Case-sensitive: ab123 and AB123 are treated as different accounts.<br/>• Leading zeros are significant: 00123 and 123 are different accounts.<br/>• No sequential or other format requirement beyond the above — any value is accepted as long as it is unique within your entity (which makes the composed account_number globally unique).<br/>• Immutable after account creation.<br/>• If omitted, the account number is generated by the system.<br/>The composed account_number is used consistently across order placement (FIX Tag 1), queries, and SOD reconciliation files.",
              "example": "0000000001"
            },
            "account_type": {
              "type": "string",
              "description": "Account Type<br/>CASH: Cash Account<br/>MARGIN: Margin Account",
              "example": "CASH",
              "enum": [
                "CASH",
                "MARGIN"
              ]
            },
            "trading_permissions": {
              "type": "array",
              "description": "Trading Permissions for the VA account, list of permission codes",
              "example": [
                "US_STOCK_NORMAL"
              ],
              "items": {
                "type": "string",
                "description": "Trading Permission<br/>US_STOCK_NORMAL - US Stock Normal Trading Permission<br/>HK_STOCK_NORMAL - HK Stock Normal Trading Permission<br/>US_OPTION_NORMAL - US Option Normal Trading Permission<br/>CN_STOCK_NORMAL - CN Stock Trading Permission",
                "example": "[\"US_STOCK_NORMAL\"]",
                "enum": [
                  "US_STOCK_NORMAL",
                  "HK_STOCK_NORMAL",
                  "CN_STOCK_NORMAL",
                  "US_OPTION_NORMAL"
                ]
              }
            },
            "option_level": {
              "type": "string",
              "description": "The option level must be less than or equal to the option level of the master account.When the trading_permissions contains US_OPTION_NORMAL, it needs to be filled in",
              "example": "LV1",
              "enum": [
                "LV1",
                "LV2",
                "LV3",
                "LV4"
              ]
            },
            "commission_code": {
              "type": "string",
              "description": "Commission Code for the VA account, determining the commission structure.The specific value of the Commission Code needs to be confirmed with the Webull business team, as it determines how commissions will be charged.\n\n",
              "example": "COMMISSION_GROUP_1",
              "enum": [
                "COMMISSION_GROUP_1",
                "COMMISSION_GROUP_2",
                "COMMISSION_GROUP_3",
                "COMMISSION_GROUP_4",
                "COMMISSION_GROUP_5"
              ]
            },
            "w8ben_info": {
              "required": [
                "first_name",
                "home_address",
                "last_name",
                "mail_address",
                "sign_date",
                "tax_id",
                "treaty_country"
              ],
              "type": "object",
              "properties": {
                "treaty_country": {
                  "type": "string",
                  "description": "Country of Treaty, using ISO 3166-1 alpha-2 format",
                  "example": "US"
                },
                "tax_id": {
                  "type": "string",
                  "description": "Tax Identification Number",
                  "example": "123-45-6789"
                },
                "sign_date": {
                  "type": "string",
                  "description": "W-8BEN Form Sign Date, format: YYYY-MM-DD",
                  "example": "2026-01-01"
                },
                "first_name": {
                  "type": "string",
                  "description": "First Name of the account holder",
                  "example": "John"
                },
                "last_name": {
                  "type": "string",
                  "description": "Last Name of the account holder",
                  "example": "Doe"
                },
                "middle_name": {
                  "type": "string",
                  "description": "Middle Name of the account holder",
                  "example": "Michael"
                },
                "home_address": {
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
                      "description": "Country, using ISO 3166-1 alpha-2 format",
                      "example": "US"
                    },
                    "state": {
                      "type": "string",
                      "description": "State or Province",
                      "example": "California"
                    },
                    "city": {
                      "type": "string",
                      "description": "City",
                      "example": "Central"
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street Address",
                      "example": "1 Queen's Road"
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal Code",
                      "example": "999077"
                    }
                  },
                  "description": "W-8BEN Address Information",
                  "title": "W8BenAddress"
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
                      "description": "Country, using ISO 3166-1 alpha-2 format",
                      "example": "US"
                    },
                    "state": {
                      "type": "string",
                      "description": "State or Province",
                      "example": "California"
                    },
                    "city": {
                      "type": "string",
                      "description": "City",
                      "example": "Central"
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street Address",
                      "example": "1 Queen's Road"
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal Code",
                      "example": "999077"
                    }
                  },
                  "description": "W-8BEN Address Information",
                  "title": "W8BenAddress"
                }
              },
              "description": "W-8BEN Information,When trading U.S. stocks in the account, this information needs to be supplemented.\n\n",
              "title": "W8BenInfoCreateRequest"
            },
            "china_connect_investor_info": {
              "required": [
                "country_of_issuance",
                "first_name",
                "id_number",
                "id_type",
                "last_name"
              ],
              "type": "object",
              "properties": {
                "first_name": {
                  "type": "string",
                  "description": "First Name",
                  "example": "John"
                },
                "last_name": {
                  "type": "string",
                  "description": "Last Name",
                  "example": "Doe"
                },
                "middle_name": {
                  "type": "string",
                  "description": "Middle Name",
                  "example": "Michael"
                },
                "id_type": {
                  "type": "string",
                  "description": "ID Type<br/>ID_CARD: Identification Card<br/>PASSPORT: Passport<br/>CERT_INCORP: Certificate of Incorporation<br/>LEI: Legal Entity Identifier<br/>OTHER_OFFICIAL_ID_DOC: Other Official ID Document",
                  "example": "ID_CARD",
                  "enum": [
                    "ID_CARD",
                    "PASSPORT",
                    "CERT_INCORP",
                    "LEI",
                    "OTHER_OFFICIAL_ID_DOC"
                  ]
                },
                "id_number": {
                  "type": "string",
                  "description": "ID Number",
                  "example": "X12345678"
                },
                "country_of_issuance": {
                  "type": "string",
                  "description": "Country of Issuance, using ISO 3166-1 alpha-2 format",
                  "example": "US"
                }
              },
              "description": "China Connect Investor Information,This information only needs to be supplemented when trading Chinese A-shares.\n\n",
              "title": "ChinaConnectInvestorInfo"
            }
          },
          "description": "Create VA Account Request",
          "title": "AccountVaCreateRequest"
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
              "account_number",
              "account_status",
              "client_request_id"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Create VA Account Request ID",
                "example": "LJIS16BACHQG9LPP44L9IQHGAB"
              },
              "account_number": {
                "type": "string",
                "description": "Account Number",
                "example": "VA000000001"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "943a9802f6c14983b3b4755c69c01717"
              },
              "account_status": {
                "type": "string",
                "description": "Account Status<br/>CREATED - Account is created but not yet active<br/>ACTIVE - Account is active and operational<br/>SUSPENDED - Account temporarily suspended<br/>RESTRICTED - Account restricted, limited operations<br/>FROZEN - Account frozen due to compliance or risk<br/>CLOSED - Account permanently closed",
                "example": "CREATED",
                "enum": [
                  "CREATED",
                  "ACTIVE",
                  "SUSPENDED",
                  "RESTRICTED",
                  "FROZEN",
                  "CLOSED"
                ]
              }
            },
            "description": "Create VA Account Result",
            "title": "AccountVaCreateResult"
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
    "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
    "belong_account_id": "TUHO91TNJ6B870VJ2SUL65U4LB",
    "external_account_number": "0000000001",
    "account_type": "CASH",
    "trading_permissions": [
      "US_STOCK_NORMAL"
    ],
    "option_level": "LV1",
    "commission_code": "COMMISSION_GROUP_1",
    "w8ben_info": {
      "treaty_country": "US",
      "tax_id": "123-45-6789",
      "sign_date": "2026-01-01",
      "first_name": "John",
      "last_name": "Doe",
      "middle_name": "Michael",
      "home_address": {
        "country": "US",
        "state": "California",
        "city": "Central",
        "street_address": "1 Queen's Road",
        "postal_code": "999077"
      },
      "mail_address": {
        "country": "US",
        "state": "California",
        "city": "Central",
        "street_address": "1 Queen's Road",
        "postal_code": "999077"
      }
    },
    "china_connect_investor_info": {
      "first_name": "John",
      "last_name": "Doe",
      "middle_name": "Michael",
      "id_type": "ID_CARD",
      "id_number": "X12345678",
      "country_of_issuance": "US"
    }
  },
  "postman": {
    "name": "Create Virtual Account",
    "description": {
      "content": "Creates a Virtual Account for a new customer.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "virtual-accounts",
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

## Update Virtual Account

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-update.md>

### Update Virtual Account

Updates Virtual Account information for an existing virtual account.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/accounts/virtual-accounts/update",
  "method": "post",
  "tags": [
    "Account (ND)"
  ],
  "description": "Updates Virtual Account information for an existing virtual account.",
  "operationId": "brokerAccountUpdate",
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
              "example": "LJIS16BACHQG9LPP44L9IQHGAB"
            },
            "account_id": {
              "type": "string",
              "description": "VA Account ID to be updated",
              "example": "943a9802f6c14983b3b4755c69c01717"
            },
            "trading_permissions": {
              "type": "array",
              "description": "Trading Permissions for the VA account, list of permission codes",
              "example": [
                "US_STOCK_NORMAL"
              ],
              "items": {
                "type": "string",
                "description": "Trading Permission<br/>US_STOCK_NORMAL - US Stock Normal Trading Permission<br/>HK_STOCK_NORMAL - HK Stock Normal Trading Permission<br/>US_OPTION_NORMAL - US Option Normal Trading Permission<br/>CN_STOCK_NORMAL - CN Stock Trading Permission",
                "example": "[\"US_STOCK_NORMAL\"]",
                "enum": [
                  "US_STOCK_NORMAL",
                  "HK_STOCK_NORMAL",
                  "CN_STOCK_NORMAL",
                  "US_OPTION_NORMAL"
                ]
              }
            },
            "option_level": {
              "type": "string",
              "description": "The option level must be less than or equal to the option level of the master account.When the trading_permissions contains US_OPTION_NORMAL, it needs to be filled in",
              "example": "LV1",
              "enum": [
                "LV1",
                "LV2",
                "LV3",
                "LV4"
              ]
            },
            "commission_code": {
              "type": "string",
              "description": "Commission Code for the VA account, determining the commission structure",
              "example": "STANDARD_COMMISSION"
            },
            "w8ben_info": {
              "required": [
                "first_name",
                "last_name",
                "sign_date",
                "tax_id",
                "treaty_country"
              ],
              "type": "object",
              "properties": {
                "treaty_country": {
                  "type": "string",
                  "description": "Country of Treaty, using ISO 3166-1 alpha-2 format",
                  "example": "US"
                },
                "tax_id": {
                  "type": "string",
                  "description": "Tax Identification Number",
                  "example": "123-45-6789"
                },
                "sign_date": {
                  "type": "string",
                  "description": "W-8BEN Form Sign Date, format: YYYY-MM-DD",
                  "example": "2026-01-01"
                },
                "first_name": {
                  "type": "string",
                  "description": "First Name of the account holder",
                  "example": "John"
                },
                "last_name": {
                  "type": "string",
                  "description": "Last Name of the account holder",
                  "example": "Doe"
                },
                "middle_name": {
                  "type": "string",
                  "description": "Middle Name of the account holder",
                  "example": "Michael"
                },
                "home_address": {
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
                      "description": "Country, using ISO 3166-1 alpha-2 format",
                      "example": "US"
                    },
                    "state": {
                      "type": "string",
                      "description": "State or Province",
                      "example": "California"
                    },
                    "city": {
                      "type": "string",
                      "description": "City",
                      "example": "Central"
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street Address",
                      "example": "1 Queen's Road"
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal Code",
                      "example": "999077"
                    }
                  },
                  "description": "W-8BEN Address Information",
                  "title": "W8BenAddress"
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
                      "description": "Country, using ISO 3166-1 alpha-2 format",
                      "example": "US"
                    },
                    "state": {
                      "type": "string",
                      "description": "State or Province",
                      "example": "California"
                    },
                    "city": {
                      "type": "string",
                      "description": "City",
                      "example": "Central"
                    },
                    "street_address": {
                      "type": "string",
                      "description": "Street Address",
                      "example": "1 Queen's Road"
                    },
                    "postal_code": {
                      "type": "string",
                      "description": "Postal Code",
                      "example": "999077"
                    }
                  },
                  "description": "W-8BEN Address Information",
                  "title": "W8BenAddress"
                }
              },
              "description": "W-8BEN Information,When trading U.S. stocks in the account, this information needs to be supplemented.\n\n",
              "title": "W8BenInfoUpdateRequest"
            },
            "china_connect_investor_info": {
              "required": [
                "country_of_issuance",
                "first_name",
                "id_number",
                "id_type",
                "last_name"
              ],
              "type": "object",
              "properties": {
                "first_name": {
                  "type": "string",
                  "description": "First Name",
                  "example": "John"
                },
                "last_name": {
                  "type": "string",
                  "description": "Last Name",
                  "example": "Doe"
                },
                "middle_name": {
                  "type": "string",
                  "description": "Middle Name",
                  "example": "Michael"
                },
                "id_type": {
                  "type": "string",
                  "description": "ID Type<br/>ID_CARD: Identification Card<br/>PASSPORT: Passport<br/>CERT_INCORP: Certificate of Incorporation<br/>LEI: Legal Entity Identifier<br/>OTHER_OFFICIAL_ID_DOC: Other Official ID Document",
                  "example": "ID_CARD",
                  "enum": [
                    "ID_CARD",
                    "PASSPORT",
                    "CERT_INCORP",
                    "LEI",
                    "OTHER_OFFICIAL_ID_DOC"
                  ]
                },
                "id_number": {
                  "type": "string",
                  "description": "ID Number",
                  "example": "X12345678"
                },
                "country_of_issuance": {
                  "type": "string",
                  "description": "Country of Issuance, using ISO 3166-1 alpha-2 format",
                  "example": "US"
                }
              },
              "description": "China Connect Investor Information,This information only needs to be supplemented when trading Chinese A-shares.\n\n",
              "title": "ChinaConnectInvestorInfo"
            }
          },
          "description": "Account VA Update Request, used to update VA account details, such as trading permissions and commission code",
          "title": "AccountVaUpdateRequest"
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
              "account_number",
              "account_status",
              "client_request_id"
            ],
            "type": "object",
            "properties": {
              "client_request_id": {
                "type": "string",
                "description": "Create VA Account Request ID",
                "example": "LJIS16BACHQG9LPP44L9IQHGAB"
              },
              "account_number": {
                "type": "string",
                "description": "Account Number",
                "example": "VA000000001"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "943a9802f6c14983b3b4755c69c01717"
              },
              "account_status": {
                "type": "string",
                "description": "Account Status<br/>CREATED - Account is created but not yet active<br/>ACTIVE - Account is active and operational<br/>SUSPENDED - Account temporarily suspended<br/>RESTRICTED - Account restricted, limited operations<br/>FROZEN - Account frozen due to compliance or risk<br/>CLOSED - Account permanently closed",
                "example": "CREATED",
                "enum": [
                  "CREATED",
                  "ACTIVE",
                  "SUSPENDED",
                  "RESTRICTED",
                  "FROZEN",
                  "CLOSED"
                ]
              }
            },
            "description": "Update VA Account Result",
            "title": "AccountVaUpdateResult"
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
    "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
    "account_id": "943a9802f6c14983b3b4755c69c01717",
    "trading_permissions": [
      "US_STOCK_NORMAL"
    ],
    "option_level": "LV1",
    "commission_code": "STANDARD_COMMISSION",
    "w8ben_info": {
      "treaty_country": "US",
      "tax_id": "123-45-6789",
      "sign_date": "2026-01-01",
      "first_name": "John",
      "last_name": "Doe",
      "middle_name": "Michael",
      "home_address": {
        "country": "US",
        "state": "California",
        "city": "Central",
        "street_address": "1 Queen's Road",
        "postal_code": "999077"
      },
      "mail_address": {
        "country": "US",
        "state": "California",
        "city": "Central",
        "street_address": "1 Queen's Road",
        "postal_code": "999077"
      }
    },
    "china_connect_investor_info": {
      "first_name": "John",
      "last_name": "Doe",
      "middle_name": "Michael",
      "id_type": "ID_CARD",
      "id_number": "X12345678",
      "country_of_issuance": "US"
    }
  },
  "postman": {
    "name": "Update Virtual Account",
    "description": {
      "content": "Updates Virtual Account information for an existing virtual account.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "virtual-accounts",
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

## Get Virtual Account Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-detail.md>

### Get Virtual Account

Retrieves detailed information for a specific virtual account. 

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/accounts/virtual-accounts/get",
  "method": "get",
  "tags": [
    "Account (ND)"
  ],
  "description": "Retrieves detailed information for a specific virtual account. ",
  "operationId": "brokerAccountDetail",
  "parameters": [
    {
      "name": "account_id",
      "in": "query",
      "description": "Virtual Account ID to retrieve details for.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "943a9802f6c14983b3b4755c69c01717"
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
              "account_status",
              "account_type",
              "belong_account_id",
              "belong_account_number",
              "trading_permissions"
            ],
            "type": "object",
            "properties": {
              "belong_account_number": {
                "type": "string",
                "description": "Belonging Account Number, the master account belong to",
                "example": "FV000123"
              },
              "belong_account_id": {
                "type": "string",
                "description": "Belonging Account ID, the master account belong to",
                "example": "4942451778e14177b9ab066b9662bda9"
              },
              "account_number": {
                "type": "string",
                "description": "Account Number",
                "example": "VA000000001"
              },
              "account_id": {
                "type": "string",
                "description": "Account ID",
                "example": "943a9802f6c14983b3b4755c69c01717"
              },
              "account_type": {
                "type": "string",
                "description": "Account Type<br/>CASH: Cash Account<br/>MARGIN: Margin Account",
                "example": "CASH",
                "enum": [
                  "CASH",
                  "MARGIN"
                ]
              },
              "account_status": {
                "type": "string",
                "description": "Account Status<br/>CREATED - Account is created but not yet active<br/>ACTIVE - Account is active and operational<br/>SUSPENDED - Account temporarily suspended<br/>RESTRICTED - Account restricted, limited operations<br/>FROZEN - Account frozen due to compliance or risk<br/>CLOSED - Account permanently closed",
                "example": "CREATED",
                "enum": [
                  "CREATED",
                  "ACTIVE",
                  "SUSPENDED",
                  "RESTRICTED",
                  "FROZEN",
                  "CLOSED"
                ]
              },
              "trading_permissions": {
                "type": "array",
                "description": "Trading Permissions",
                "example": [
                  "US_STOCK_NORMAL",
                  "US_STOCK_FRACTIONAL",
                  "US_STOCK_OVERNIGHT",
                  "HK_STOCK_NORMAL",
                  "HK_STOCK_FRACTIONAL"
                ],
                "items": {
                  "type": "string",
                  "description": "Trading Permission<br/>US_STOCK_NORMAL - US Stock Normal Trading Permission<br/>HK_STOCK_NORMAL - HK Stock Normal Trading Permission<br/>US_OPTION_NORMAL - US Option Normal Trading Permission<br/>CN_STOCK_NORMAL - CN Stock Trading Permission",
                  "example": "[\"US_STOCK_NORMAL\",\"US_STOCK_FRACTIONAL\",\"US_STOCK_OVERNIGHT\",\"HK_STOCK_NORMAL\",\"HK_STOCK_FRACTIONAL\"]",
                  "enum": [
                    "US_STOCK_NORMAL",
                    "HK_STOCK_NORMAL",
                    "CN_STOCK_NORMAL",
                    "US_OPTION_NORMAL"
                  ]
                }
              },
              "option_level": {
                "type": "string",
                "description": "The option level must be less than or equal to the option level of the master account.When the trading_permissions contains US_OPTION_NORMAL, it needs to be filled in",
                "example": "LV1",
                "enum": [
                  "LV1",
                  "LV2",
                  "LV3",
                  "LV4"
                ]
              },
              "commission_code": {
                "type": "string",
                "description": "Commission Code",
                "example": "STANDARD_COMMISSION"
              },
              "w8ben_info": {
                "required": [
                  "first_name",
                  "home_address",
                  "last_name",
                  "sign_date",
                  "tax_id",
                  "treaty_country"
                ],
                "type": "object",
                "properties": {
                  "treaty_country": {
                    "type": "string",
                    "description": "Country of Treaty, using ISO 3166-1 alpha-2 format",
                    "example": "US"
                  },
                  "tax_id": {
                    "type": "string",
                    "description": "Tax Identification Number",
                    "example": "123-45-6789"
                  },
                  "sign_date": {
                    "type": "string",
                    "description": "W-8BEN Form Sign Date, format: YYYY-MM-DD",
                    "example": "2026-01-01"
                  },
                  "first_name": {
                    "type": "string",
                    "description": "First Name of the account holder",
                    "example": "John"
                  },
                  "last_name": {
                    "type": "string",
                    "description": "Last Name of the account holder",
                    "example": "Doe"
                  },
                  "middle_name": {
                    "type": "string",
                    "description": "Middle Name of the account holder",
                    "example": "Michael"
                  },
                  "home_address": {
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
                        "description": "Country, using ISO 3166-1 alpha-2 format",
                        "example": "US"
                      },
                      "state": {
                        "type": "string",
                        "description": "State or Province",
                        "example": "California"
                      },
                      "city": {
                        "type": "string",
                        "description": "City",
                        "example": "Central"
                      },
                      "street_address": {
                        "type": "string",
                        "description": "Street Address",
                        "example": "1 Queen's Road"
                      },
                      "postal_code": {
                        "type": "string",
                        "description": "Postal Code",
                        "example": "999077"
                      }
                    },
                    "description": "W-8BEN Address Information",
                    "title": "W8BenAddress"
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
                        "description": "Country, using ISO 3166-1 alpha-2 format",
                        "example": "US"
                      },
                      "state": {
                        "type": "string",
                        "description": "State or Province",
                        "example": "California"
                      },
                      "city": {
                        "type": "string",
                        "description": "City",
                        "example": "Central"
                      },
                      "street_address": {
                        "type": "string",
                        "description": "Street Address",
                        "example": "1 Queen's Road"
                      },
                      "postal_code": {
                        "type": "string",
                        "description": "Postal Code",
                        "example": "999077"
                      }
                    },
                    "description": "W-8BEN Address Information",
                    "title": "W8BenAddress"
                  }
                },
                "description": "W-8BEN Information,When trading U.S. stocks in the account, this information needs to be supplemented.\n\n",
                "title": "W8BenInfoResult"
              },
              "china_connect_investor_info": {
                "required": [
                  "country_of_issuance",
                  "first_name",
                  "id_number",
                  "id_type",
                  "last_name"
                ],
                "type": "object",
                "properties": {
                  "first_name": {
                    "type": "string",
                    "description": "First Name",
                    "example": "John"
                  },
                  "last_name": {
                    "type": "string",
                    "description": "Last Name",
                    "example": "Doe"
                  },
                  "middle_name": {
                    "type": "string",
                    "description": "Middle Name",
                    "example": "Michael"
                  },
                  "id_type": {
                    "type": "string",
                    "description": "ID Type<br/>ID_CARD: Identification Card<br/>PASSPORT: Passport<br/>CERT_INCORP: Certificate of Incorporation<br/>LEI: Legal Entity Identifier<br/>OTHER_OFFICIAL_ID_DOC: Other Official ID Document",
                    "example": "ID_CARD",
                    "enum": [
                      "ID_CARD",
                      "PASSPORT",
                      "CERT_INCORP",
                      "LEI",
                      "OTHER_OFFICIAL_ID_DOC"
                    ]
                  },
                  "id_number": {
                    "type": "string",
                    "description": "ID Number",
                    "example": "X12345678"
                  },
                  "country_of_issuance": {
                    "type": "string",
                    "description": "Country of Issuance, using ISO 3166-1 alpha-2 format",
                    "example": "US"
                  }
                },
                "description": "China Connect Investor Information,This information only needs to be supplemented when trading Chinese A-shares.\n\n",
                "title": "ChinaConnectInvestorInfo"
              }
            },
            "description": "VA Account Detail Result",
            "title": "AccountVaDetailResult"
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
    "name": "Get Virtual Account",
    "description": {
      "content": "Retrieves detailed information for a specific virtual account. ",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "virtual-accounts",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Virtual Account ID to retrieve details for.",
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

## List Virtual Accounts

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-list.md>

### List Virtual Accounts

Retrieves a paginated list of Virtual Accounts.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/accounts/virtual-accounts/list",
  "method": "get",
  "tags": [
    "Account (ND)"
  ],
  "description": "Retrieves a paginated list of Virtual Accounts.",
  "operationId": "brokerAccountList",
  "parameters": [
    {
      "name": "start_time",
      "in": "query",
      "description": "Inclusive start time for the query. formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:11.591Z"
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "Inclusive end time for the query, formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-06T22:59:11.591Z"
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
                    "account_number",
                    "account_status",
                    "account_type",
                    "belong_account_id",
                    "belong_account_number",
                    "trading_permissions"
                  ],
                  "type": "object",
                  "properties": {
                    "belong_account_number": {
                      "type": "string",
                      "description": "Belonging Account Number, the master account belong to",
                      "example": "FV000123"
                    },
                    "belong_account_id": {
                      "type": "string",
                      "description": "Belonging Account ID, the master account belong to",
                      "example": "4942451778e14177b9ab066b9662bda9"
                    },
                    "account_number": {
                      "type": "string",
                      "description": "Account Number",
                      "example": "VA000000001"
                    },
                    "account_id": {
                      "type": "string",
                      "description": "Account ID",
                      "example": "943a9802f6c14983b3b4755c69c01717"
                    },
                    "account_type": {
                      "type": "string",
                      "description": "Account Type<br/>CASH: Cash Account<br/>MARGIN: Margin Account",
                      "example": "CASH",
                      "enum": [
                        "CASH",
                        "MARGIN"
                      ]
                    },
                    "account_status": {
                      "type": "string",
                      "description": "Account Status<br/>CREATED - Account is created but not yet active<br/>ACTIVE - Account is active and operational<br/>SUSPENDED - Account temporarily suspended<br/>RESTRICTED - Account restricted, limited operations<br/>FROZEN - Account frozen due to compliance or risk<br/>CLOSED - Account permanently closed",
                      "example": "CREATED",
                      "enum": [
                        "CREATED",
                        "ACTIVE",
                        "SUSPENDED",
                        "RESTRICTED",
                        "FROZEN",
                        "CLOSED"
                      ]
                    },
                    "trading_permissions": {
                      "type": "array",
                      "description": "Trading Permissions",
                      "example": [
                        "US_STOCK_NORMAL",
                        "US_STOCK_FRACTIONAL",
                        "US_STOCK_OVERNIGHT",
                        "HK_STOCK_NORMAL",
                        "HK_STOCK_FRACTIONAL"
                      ],
                      "items": {
                        "type": "string",
                        "description": "Trading Permission<br/>US_STOCK_NORMAL - US Stock Normal Trading Permission<br/>HK_STOCK_NORMAL - HK Stock Normal Trading Permission<br/>US_OPTION_NORMAL - US Option Normal Trading Permission<br/>CN_STOCK_NORMAL - CN Stock Trading Permission",
                        "example": "[\"US_STOCK_NORMAL\",\"US_STOCK_FRACTIONAL\",\"US_STOCK_OVERNIGHT\",\"HK_STOCK_NORMAL\",\"HK_STOCK_FRACTIONAL\"]",
                        "enum": [
                          "US_STOCK_NORMAL",
                          "HK_STOCK_NORMAL",
                          "CN_STOCK_NORMAL",
                          "US_OPTION_NORMAL"
                        ]
                      }
                    },
                    "option_level": {
                      "type": "string",
                      "description": "The option level must be less than or equal to the option level of the master account.When the trading_permissions contains US_OPTION_NORMAL, it needs to be filled in",
                      "example": "LV1",
                      "enum": [
                        "LV1",
                        "LV2",
                        "LV3",
                        "LV4"
                      ]
                    },
                    "commission_code": {
                      "type": "string",
                      "description": "Commission Code",
                      "example": "STANDARD_COMMISSION"
                    },
                    "w8ben_info": {
                      "required": [
                        "first_name",
                        "home_address",
                        "last_name",
                        "sign_date",
                        "tax_id",
                        "treaty_country"
                      ],
                      "type": "object",
                      "properties": {
                        "treaty_country": {
                          "type": "string",
                          "description": "Country of Treaty, using ISO 3166-1 alpha-2 format",
                          "example": "US"
                        },
                        "tax_id": {
                          "type": "string",
                          "description": "Tax Identification Number",
                          "example": "123-45-6789"
                        },
                        "sign_date": {
                          "type": "string",
                          "description": "W-8BEN Form Sign Date, format: YYYY-MM-DD",
                          "example": "2026-01-01"
                        },
                        "first_name": {
                          "type": "string",
                          "description": "First Name of the account holder",
                          "example": "John"
                        },
                        "last_name": {
                          "type": "string",
                          "description": "Last Name of the account holder",
                          "example": "Doe"
                        },
                        "middle_name": {
                          "type": "string",
                          "description": "Middle Name of the account holder",
                          "example": "Michael"
                        },
                        "home_address": {
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
                              "description": "Country, using ISO 3166-1 alpha-2 format",
                              "example": "US"
                            },
                            "state": {
                              "type": "string",
                              "description": "State or Province",
                              "example": "California"
                            },
                            "city": {
                              "type": "string",
                              "description": "City",
                              "example": "Central"
                            },
                            "street_address": {
                              "type": "string",
                              "description": "Street Address",
                              "example": "1 Queen's Road"
                            },
                            "postal_code": {
                              "type": "string",
                              "description": "Postal Code",
                              "example": "999077"
                            }
                          },
                          "description": "W-8BEN Address Information",
                          "title": "W8BenAddress"
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
                              "description": "Country, using ISO 3166-1 alpha-2 format",
                              "example": "US"
                            },
                            "state": {
                              "type": "string",
                              "description": "State or Province",
                              "example": "California"
                            },
                            "city": {
                              "type": "string",
                              "description": "City",
                              "example": "Central"
                            },
                            "street_address": {
                              "type": "string",
                              "description": "Street Address",
                              "example": "1 Queen's Road"
                            },
                            "postal_code": {
                              "type": "string",
                              "description": "Postal Code",
                              "example": "999077"
                            }
                          },
                          "description": "W-8BEN Address Information",
                          "title": "W8BenAddress"
                        }
                      },
                      "description": "W-8BEN Information,When trading U.S. stocks in the account, this information needs to be supplemented.\n\n",
                      "title": "W8BenInfoResult"
                    },
                    "china_connect_investor_info": {
                      "required": [
                        "country_of_issuance",
                        "first_name",
                        "id_number",
                        "id_type",
                        "last_name"
                      ],
                      "type": "object",
                      "properties": {
                        "first_name": {
                          "type": "string",
                          "description": "First Name",
                          "example": "John"
                        },
                        "last_name": {
                          "type": "string",
                          "description": "Last Name",
                          "example": "Doe"
                        },
                        "middle_name": {
                          "type": "string",
                          "description": "Middle Name",
                          "example": "Michael"
                        },
                        "id_type": {
                          "type": "string",
                          "description": "ID Type<br/>ID_CARD: Identification Card<br/>PASSPORT: Passport<br/>CERT_INCORP: Certificate of Incorporation<br/>LEI: Legal Entity Identifier<br/>OTHER_OFFICIAL_ID_DOC: Other Official ID Document",
                          "example": "ID_CARD",
                          "enum": [
                            "ID_CARD",
                            "PASSPORT",
                            "CERT_INCORP",
                            "LEI",
                            "OTHER_OFFICIAL_ID_DOC"
                          ]
                        },
                        "id_number": {
                          "type": "string",
                          "description": "ID Number",
                          "example": "X12345678"
                        },
                        "country_of_issuance": {
                          "type": "string",
                          "description": "Country of Issuance, using ISO 3166-1 alpha-2 format",
                          "example": "US"
                        }
                      },
                      "description": "China Connect Investor Information,This information only needs to be supplemented when trading Chinese A-shares.\n\n",
                      "title": "ChinaConnectInvestorInfo"
                    }
                  },
                  "description": "VA Account Detail Result",
                  "title": "AccountVaDetailResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoAccountVaDetailResult"
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
    "name": "List Virtual Accounts",
    "description": {
      "content": "Retrieves a paginated list of Virtual Accounts.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "accounts",
        "virtual-accounts",
        "list"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "Inclusive start time for the query. formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Inclusive end time for the query, formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
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

## Get Stock Instrument

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-instrument-list.md>

### List Stock Instruments

Retrieves detail information of instruments associated to customer's account.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/instruments/stocks/profiles/list",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves detail information of instruments associated to customer's account.",
  "operationId": "brokerInstrumentList",
  "parameters": [
    {
      "name": "category",
      "in": "query",
      "description": "Security type.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Instrument Stock Category<br/>US_STOCK - US stock, <br/>HK_STOCK - HK stock, <br/>CN_STOCK - China A share<br/>",
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
        "description": "Tradable Status<br/>OC - Tradable: Security is available for trading<br/>CO - Liquidate only: Security can only be sold, no purchases allowed<br/>NT - Non-Tradable: Security cannot be traded",
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
                      "description": "Instrument Stock Category<br/>US_STOCK - US stock, <br/>HK_STOCK - HK stock, <br/>CN_STOCK - China A share<br/>",
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
                      "description": "Tradable Status<br/>OC - Tradable: Security is available for trading<br/>CO - Liquidate only: Security can only be sold, no purchases allowed<br/>NT - Non-Tradable: Security cannot be traded",
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
                    },
                    "etf_leveraged_flag": {
                      "type": "string",
                      "description": "Is Leveraged ETF. Allowed values: YES / NO. Null if not an ETF",
                      "example": "YES"
                    },
                    "etf_leveraged_factor": {
                      "type": "string",
                      "description": "ETF Leveraged Factor. Numeric string, positive for bullish, negative for inverse. Null if not an ETF",
                      "example": "2"
                    },
                    "inverse_etf": {
                      "type": "string",
                      "description": "Is Inverse ETF. Allowed values: true / false. Null if not an ETF",
                      "example": "true"
                    },
                    "crypto_etf": {
                      "type": "string",
                      "description": "Is Crypto ETF. Allowed values: true / false. Null if not an ETF",
                      "example": "true"
                    },
                    "single_stock_etf": {
                      "type": "string",
                      "description": "Is Single Stock ETF. Allowed values: true / false. Null if not an ETF",
                      "example": "true"
                    }
                  },
                  "description": "Broker Securities Information",
                  "title": "BrokerInstrumentResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "BrokerInstrumentPageResult"
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
      "content": "Retrieves detail information of instruments associated to customer's account.",
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

## Get Stock Locate Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-stock-locate-detail.md>

### Get Stock Locate

Retrieves stock short trading related information, including stock borrow type (ETB/HTB), available short quantity, short interest rate, etc. (Reference use only).

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/instruments/stock-locates/get",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves stock short trading related information, including stock borrow type (ETB/HTB), available short quantity, short interest rate, etc. (Reference use only).",
  "operationId": "brokerStockLocateDetail",
  "parameters": [
    {
      "name": "symbols",
      "in": "query",
      "description": "List of security symbols, maximum 100 symbols per query.",
      "required": true,
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
        "description": "Locate Category<br/>US_STOCK - US stock<br/>",
        "enum": [
          "US_STOCK"
        ]
      },
      "example": "US_STOCK"
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
                  "description": "Symbol of the instrument",
                  "example": "AAPL"
                },
                "shortable": {
                  "type": "string",
                  "description": "Whether the security is shortable: Y = Yes, N = No",
                  "example": "Y"
                },
                "short_type": {
                  "type": "string",
                  "description": "Short availability type: ETB (Easy To Borrow), HTB (Hard To Borrow)",
                  "example": "ETB"
                },
                "available_num": {
                  "type": "integer",
                  "description": "remaining available shares for short selling (not returned for ETB)",
                  "format": "int32",
                  "example": 1000
                },
                "short_interest_rate": {
                  "type": "number",
                  "description": "Annualized short interest rate",
                  "format": "double",
                  "example": 5.5
                }
              },
              "description": "Stock Locate Detail",
              "title": "StockLocateResult"
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
    "name": "Get Stock Locate",
    "description": {
      "content": "Retrieves stock short trading related information, including stock borrow type (ETB/HTB), available short quantity, short interest rate, etc. (Reference use only).",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "instruments",
        "stock-locates",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) List of security symbols, maximum 100 symbols per query.",
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

## Get Corporate Actions Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-corporate-actions-detail.md>

### Get Corporate Actions Detail

Retrieves detailed information for corporate action events.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/instruments/stocks/corporate-actions/get",
  "method": "get",
  "tags": [
    "Instrument"
  ],
  "description": "Retrieves detailed information for corporate action events.",
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
                "description": "Instrument Stock Category<br/>US_STOCK - US stock, <br/>HK_STOCK - HK stock, <br/>CN_STOCK - China A share<br/>",
                "example": "HK_STOCK",
                "enum": [
                  "US_STOCK",
                  "HK_STOCK",
                  "CN_STOCK"
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
                "example": "HK"
              },
              "listing_country_of_code": {
                "type": "string",
                "description": "Listing Country Code, ISO 3166-1 alpha-2 format",
                "example": "HK"
              },
              "issuer_country_code": {
                "type": "string",
                "description": "Issuer Country Code, ISO 3166-1 alpha-2 format",
                "example": "HK"
              },
              "from": {
                "type": "object",
                "properties": {
                  "symbol": {
                    "type": "string",
                    "description": "Symbol of the instrument",
                    "example": "01029"
                  },
                  "name": {
                    "type": "string",
                    "description": "Name of the instrument",
                    "example": "IRC Limited"
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
                      "description": "Default Option Flag. Indicates if the default payout option is exposed",
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
      "content": "Retrieves detailed information for corporate action events.",
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

## Account Activities By Type

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-activity-by-type.md>

### List Cash Activities

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/activities/cash-activities/list",
  "method": "get",
  "tags": [
    "Activity"
  ],
  "description": "Retrieves account transaction activities records with details.",
  "operationId": "brokerActivityByType",
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
        "description": "Activity Type<br/>ADJUSTMENT - Adjustment Activity Type<br/>ALLOCATION - Allocation Activity Type<br/>DEPOSIT - Deposit Activity Type<br/>DIVIDENDS - Dividends Activity Type<br/>EXECUTION - Execution Activity Type<br/>FEES - Fees Activity Type<br/>INTERESTS - Interests Activity Type<br/>JOURNAL - Journal Activity Type<br/>OPTION_EA - Option EA Activity Type<br/>REORGANIZATION - Reorganization Activity Type<br/>SUBSCRIBE - Subscribe Activity Type<br/>TAX - Tax Activity Type<br/>TRADE - Trade Activity Type<br/>TRANSFER - Transfer Activity Type<br/>WITHDRAW - Withdraw Activity Type",
        "enum": [
          "ADJUSTMENT",
          "ALLOCATION",
          "DEPOSIT",
          "DIVIDENDS",
          "EXECUTION",
          "FEES",
          "INTERESTS",
          "JOURNAL",
          "OPTION_EA",
          "REORGANIZATION",
          "SUBSCRIBE",
          "TAX",
          "TRADE",
          "TRANSFER",
          "WITHDRAW"
        ]
      },
      "example": "TRADE,TRANSFER"
    },
    {
      "name": "start_time",
      "in": "query",
      "description": "Activity query start time, formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T22:59:11.591Z"
    },
    {
      "name": "end_time",
      "in": "query",
      "description": "Activity query end time, formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2025-01-05T23:59:11.591Z"
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
                    "account_number",
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
                    "account_number": {
                      "type": "string",
                      "description": "Account Number",
                      "example": "VA000000001"
                    },
                    "activity_type": {
                      "type": "string",
                      "description": "Activity Type<br/>ADJUSTMENT - Adjustment Activity Type<br/>ALLOCATION - Allocation Activity Type<br/>DEPOSIT - Deposit Activity Type<br/>DIVIDENDS - Dividends Activity Type<br/>EXECUTION - Execution Activity Type<br/>FEES - Fees Activity Type<br/>INTERESTS - Interests Activity Type<br/>JOURNAL - Journal Activity Type<br/>OPTION_EA - Option EA Activity Type<br/>REORGANIZATION - Reorganization Activity Type<br/>SUBSCRIBE - Subscribe Activity Type<br/>TAX - Tax Activity Type<br/>TRADE - Trade Activity Type<br/>TRANSFER - Transfer Activity Type<br/>WITHDRAW - Withdraw Activity Type<br/>OTHER - Other Activity Type",
                      "example": "TRADE",
                      "enum": [
                        "ADJUSTMENT",
                        "ALLOCATION",
                        "DEPOSIT",
                        "DIVIDENDS",
                        "EXECUTION",
                        "FEES",
                        "INTERESTS",
                        "JOURNAL",
                        "OPTION_EA",
                        "REORGANIZATION",
                        "SUBSCRIBE",
                        "TAX",
                        "TRADE",
                        "TRANSFER",
                        "WITHDRAW",
                        "OTHER"
                      ]
                    },
                    "activity_sub_type": {
                      "type": "string",
                      "description": "Activity Sub Type. Provides detailed classification under each ActivityType.<br/>Multiple sub-types may share the same enum value but belong to different parent types.<br/>For example, FOREIGN_TAX_WITHHELD can appear under DIVIDENDS, FEES, or TAX types.<br/>NOT_SET - No sub type<br/>CASH_ADJUSTMENT - Cash Adjustment Sub Type, super type is ADJUSTMENT<br/>WIRE - Wire Deposit Sub Type, super type is DEPOSIT or WITHDRAW<br/>TAX_WITHHOLD - Tax Withhold Dividend Sub Type, super type is DIVIDENDS<br/>COLLECTION_FEE - Collection Fee Dividend Sub Type, super type is DIVIDENDS or INTERESTS<br/>INCOME - Income Dividend Sub Type, super type is DIVIDENDS<br/>PAYMENT_IN_LIEU - Payment In Lieu Dividend Sub Type, super type is DIVIDENDS<br/>SCRIP_FEE - Scrip Fee Dividend Sub Type, super type is DIVIDENDS or REORGANIZATION<br/>US_TAX_WITHHOLDING - US Tax Withholding Dividend or Tax Sub Type, super type is DIVIDENDS or TAX<br/>CA_HANDLING_FEE - CA Handling Fee Sub Type, super type is DIVIDENDS or REORGANIZATION<br/>FOREIGN_TAX_WITHHELD - Foreign Tax Withheld Sub Type, super type is DIVIDENDS or FEES or TAX<br/>CASH_IN_LIEU - Cash In Lieu Sub Type, super type is DIVIDENDS or REORGANIZATION<br/>ADJUSTMENT - Adjustment Dividend Sub Type, super type is DIVIDENDS<br/>ADR - ADR Fees Sub Type, super type is FEES<br/>IPO_SUBSCRIPTION_FEE - IPO Subscription Fee Sub Type, super type is FEES<br/>TRANSACTION_FEES - Transaction Fees Sub Type, super type is FEES<br/>CCASS_SETTLEMENT_FEE - CCASS Settlement Fee Sub Type, super type is FEES<br/>JOURNAL_BETWEEN_ACCOUNTS - Journal Between Accounts Sub Type, super type is FEES or JOURNAL<br/>CUSTODIAN_FEE - Custodian Fee Sub Type, super type is FEES<br/>SETTLEMENT_FEES - Settlement Fees Sub Type, super type is FEES<br/>TRANSFER_FOP - Transfer FOP Fees Sub Type, super type is FEES<br/>FTT - FTT Fees Sub Type, super type is FEES<br/>COMMISSION - Commission Fees Sub Type, super type is FEES<br/>SYSTEM_FEE - System Fee Sub Type, super type is FEES<br/>CORRESPONDENT_BILLING - Correspondent Billing Fees Sub Type, super type is FEES<br/>CUSTODY_FEE - Custody Fee Sub Type, super type is FEES<br/>STOCK_BORROW_INTEREST - Stock Borrow Interest Sub Type, super type is INTERESTS<br/>DEBIT_CASH - Debit Cash Sub Type, super type is INTERESTS<br/>IPO_FINANCING - IPO Financing Interest Sub Type, super type is INTERESTS<br/>PAYMENT - Payment Interests Sub Type, super type is INTERESTS<br/>ACCOUNT_MIGRATION - Account Migration Sub Type, super type is JOURNAL<br/>IPO_REFUND - IPO Refund Sub Type, super type is JOURNAL<br/>JOURNAL_BETWEEN_TYPES - Journal Between Types Sub Type, super type is JOURNAL<br/>GIFTING - Gifting Sub Type, super type is JOURNAL<br/>BANK_TRANSFER - Bank Transfer Sub Type, super type is JOURNAL<br/>FX_EXCHANGE - FX Exchange Sub Type, super type is JOURNAL<br/>CA_BREAK_RESOLVING - CA Break Resolving Sub Type, super type is JOURNAL<br/>MARK_TO_MARKET - Mark To Market Sub Type, super type is JOURNAL<br/>INV_DEPOSIT - INV Deposit Sub Type, super type is JOURNAL<br/>IPO_SUBSCRIPTION - IPO Subscription Sub Type, super type is JOURNAL<br/>FIXED_INCOME_SPREAD - Fixed Income Spread Sub Type, super type is JOURNAL<br/>HK_CBBC_EXERCISE - HK CBBC Exercise Sub Type, super type is OPTION_EA<br/>HK_CBBC_EXPIRATION - HK CBBC Expiration Sub Type, super type is OPTION_EA<br/>HK_WARRANT_EXERCISE - HK Warrant Exercise Sub Type, super type is OPTION_EA<br/>OPTION_EXERCISE - Option Exercise Sub Type, super type is OPTION_EA or TRADE<br/>HK_INLINE_WARRANT_EXERCISE - HK Inline Warrant Exercise Sub Type, super type is OPTION_EA<br/>OPTION_ASSIGNMENT - Option Assignment Sub Type, super type is OPTION_EA or TRADE<br/>MERGER - Merger Sub Type, super type is REORGANIZATION<br/>RIGHT_SUBSCRIPTION - Right Subscription Sub Type, super type is REORGANIZATION<br/>CORP_ACTION_FEE - Corp Action Fee Sub Type, super type is REORGANIZATION<br/>LIQUIDATION - Liquidation Sub Type, super type is REORGANIZATION<br/>REDEMPTION - Redemption Sub Type, super type is REORGANIZATION<br/>REPURCHASE - Repurchase Sub Type, super type is REORGANIZATION<br/>FEE - Subscribe Fee Sub Type, super type is SUBSCRIBE<br/>FEE_RETURN - Subscribe Fee Return Sub Type, super type is SUBSCRIBE<br/>INTERNAL_TRANSFER - Internal Transfer Sub Type, super type is TRANSFER<br/>DTC_OUT - DTC Out Sub Type, super type is TRANSFER<br/>WIRE_FEE_LOCAL - Wire Fee Local Sub Type, super type is WITHDRAW<br/>OPTION_TRADE - Option Trade Sub Type, super type is TRADE<br/>FUND_TRADE - Fund Trade Sub Type, super type is TRADE<br/>OTHER - Other",
                      "example": "NOT_SET",
                      "enum": [
                        "NOT_SET",
                        "CASH_ADJUSTMENT",
                        "WIRE",
                        "TAX_WITHHOLD",
                        "COLLECTION_FEE",
                        "INCOME",
                        "PAYMENT_IN_LIEU",
                        "SCRIP_FEE",
                        "US_TAX_WITHHOLDING",
                        "CA_HANDLING_FEE",
                        "FOREIGN_TAX_WITHHELD",
                        "CASH_IN_LIEU",
                        "ADJUSTMENT",
                        "ADR",
                        "IPO_SUBSCRIPTION_FEE",
                        "TRANSACTION_FEES",
                        "CCASS_SETTLEMENT_FEE",
                        "JOURNAL_BETWEEN_ACCOUNTS",
                        "CUSTODIAN_FEE",
                        "SETTLEMENT_FEES",
                        "TRANSFER_FOP",
                        "FTT",
                        "COMMISSION",
                        "SYSTEM_FEE",
                        "CORRESPONDENT_BILLING",
                        "CUSTODY_FEE",
                        "STOCK_BORROW_INTEREST",
                        "DEBIT_CASH",
                        "IPO_FINANCING",
                        "PAYMENT",
                        "ACCOUNT_MIGRATION",
                        "IPO_REFUND",
                        "JOURNAL_BETWEEN_TYPES",
                        "GIFTING",
                        "BANK_TRANSFER",
                        "FX_EXCHANGE",
                        "CA_BREAK_RESOLVING",
                        "MARK_TO_MARKET",
                        "INV_DEPOSIT",
                        "IPO_SUBSCRIPTION",
                        "FIXED_INCOME_SPREAD",
                        "HK_CBBC_EXERCISE",
                        "HK_CBBC_EXPIRATION",
                        "HK_WARRANT_EXERCISE",
                        "OPTION_EXERCISE",
                        "HK_INLINE_WARRANT_EXERCISE",
                        "OPTION_ASSIGNMENT",
                        "MERGER",
                        "RIGHT_SUBSCRIPTION",
                        "CORP_ACTION_FEE",
                        "LIQUIDATION",
                        "REDEMPTION",
                        "REPURCHASE",
                        "FEE",
                        "FEE_RETURN",
                        "INTERNAL_TRANSFER",
                        "DTC_OUT",
                        "WIRE_FEE_LOCAL",
                        "OPTION_TRADE",
                        "FUND_TRADE",
                        "OTHER"
                      ]
                    },
                    "currency": {
                      "type": "string",
                      "description": "Currency<br/>CNH - RMB<br/>HKD - Hong Kong Dollar<br/>USD - US Dollar",
                      "example": "USD",
                      "enum": [
                        "CNH",
                        "HKD",
                        "USD"
                      ]
                    },
                    "market": {
                      "type": "string",
                      "description": "Market Code<br/>HK - Hong Kong Market<br/>US - US Market<br/>CN - China Market",
                      "example": "US",
                      "enum": [
                        "HK",
                        "US",
                        "CN"
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
                      "example": "2024-05-01T10:15:30Z"
                    }
                  },
                  "description": "Activity Result",
                  "title": "ActivityResult"
                }
              },
              "pagination_key": {
                "type": "string",
                "description": "Pagination key for next page. If absent, indicates this is the last page.",
                "example": "eyJ2IjoxLCJsYXN0SWQiOiI5MTMyNDQ3NjkiLCJwYWdlSW===="
              }
            },
            "description": "Paginated result with cursor-based pagination",
            "title": "PaginatedResultVoActivityResult"
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
            "content": "Activity query start time, formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
            "type": "text/plain"
          },
          "key": "start_time",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "Activity query end time, formats: yyyy-MM-dd'T'HH:mm:ss.SSSZ",
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

## Account Balance

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-assets-balance.md>

### Get Account Balance

Retrieves account asset information, including cash balance, position market value, profit & loss, buying power, interest payable, etc.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/assets/balances/get",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves account asset information, including cash balance, position market value, profit & loss, buying power, interest payable, etc.",
  "operationId": "brokerAssetsBalance",
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
                  "title": "BrokerAssetsCurrencyAssets"
                }
              }
            },
            "title": "BrokerAssetsBalanceResult"
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
      "content": "Retrieves account asset information, including cash balance, position market value, profit & loss, buying power, interest payable, etc.",
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

## Account Positions

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-assets-positions.md>

### List Account Positions

Retrieves position details for a specific account.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/assets/positions/list",
  "method": "get",
  "tags": [
    "Assets"
  ],
  "description": "Retrieves position details for a specific account.",
  "operationId": "brokerAssetsPositions",
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
      "description": "Request successful",
      "content": {
        "application/json": {
          "schema": {
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
                "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                "example": "SINGLE",
                "enum": [
                  "SINGLE",
                  "COVERED_STOCK",
                  "STRADDLE",
                  "STRANGLE",
                  "VERTICAL",
                  "CALENDAR",
                  "DIAGONAL",
                  "COLLAR_WITH_STOCK",
                  "BUTTERFLY",
                  "CONDOR",
                  "IRON_BUTTERFLY",
                  "IRON_CONDOR",
                  "RATIO"
                ]
              },
              "instrument_type": {
                "type": "string",
                "description": "Type of financial instrument associated with the request.",
                "example": "EQUITY",
                "enum": [
                  "EQUITY",
                  "OPTION"
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
                    "leg_id",
                    "symbol"
                  ],
                  "type": "object",
                  "properties": {
                    "symbol": {
                      "type": "string",
                      "description": "Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives).",
                      "example": "AAPL"
                    },
                    "leg_id": {
                      "type": "string",
                      "description": "Leg Id",
                      "example": "T8HO5CR1CI7MEF4RF9U98JM189"
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
                  "title": "BrokerPositionItem"
                }
              }
            },
            "title": "BrokerAssetsPositionResult"
          },
          "examples": {
            "Example": {
              "description": "Example",
              "value": [
                {
                  "currency": "USD",
                  "quantity": "10.0000000000",
                  "legs": [
                    {
                      "quantity": "1000.0000000000",
                      "symbol": "AAPL",
                      "leg_id": "037SCRJQ8G8F60KHJOS8000001"
                    },
                    {
                      "quantity": "10.0000000000",
                      "symbol": "AAPL",
                      "leg_id": "037SCRJQ9M8F60KHJOS8000001",
                      "option_type": "PUT",
                      "option_expire_date": "2026-11-20",
                      "option_exercise_price": "330.00",
                      "option_contract_multiplier": "100.00",
                      "option_contract_deliverable": "100.00"
                    }
                  ],
                  "position_id": "037SCRJQ8A8F60KHJOS8000000",
                  "symbol": "AAPL",
                  "option_strategy": "COVERED_STOCK",
                  "instrument_type": "OPTION",
                  "cost_price": "347.24",
                  "last_price": "347.37",
                  "unrealized_profit_loss": "119.51"
                }
              ]
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
      "content": "Retrieves position details for a specific account.",
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

## Order Preview

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-preview.md>

### Preview Order

Calculates the estimated trading cost and fees based on the provided trade information. 

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/preview",
  "method": "post",
  "tags": [
    "Order"
  ],
  "description": "Calculates the estimated trading cost and fees based on the provided trade information. ",
  "operationId": "brokerOrderPreview",
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
            "client_combo_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the combined order<br/> If combo_type = NORMAL and client_combo_order_id not need to set<br/> If combo_type != NORMAL and client_combo_order_id not provided <br/> the server will automatically generate one.",
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
                    "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) ",
                    "example": "NORMAL"
                  },
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY",
                      "OPTION"
                    ]
                  },
                  "market": {
                    "type": "string",
                    "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
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
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order<br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> HK Options<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
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
                      "MARKET_ON_CLOSE",
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
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                    "example": "DAY",
                    "enum": [
                      "DAY",
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
                    "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                    "example": "SINGLE",
                    "enum": [
                      "SINGLE",
                      "COVERED_STOCK",
                      "STRADDLE",
                      "STRANGLE",
                      "VERTICAL",
                      "CALENDAR",
                      "DIAGONAL",
                      "COLLAR_WITH_STOCK",
                      "BUTTERFLY",
                      "CONDOR",
                      "IRON_BUTTERFLY",
                      "IRON_CONDOR",
                      "RATIO"
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
                            "OPTION"
                          ]
                        },
                        "market": {
                          "type": "string",
                          "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
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
                      "description": "Option leg detail. Only required when previewing option orders.",
                      "title": "BrokerOptionPlaceLegParam"
                    }
                  }
                },
                "description": "Order Details",
                "title": "BrokerOrderPreviewItemParam"
              }
            }
          },
          "title": "BrokerOrderPreviewParam"
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
            "title": "BrokerOrderPreviewResult"
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
      "content": "Calculates the estimated trading cost and fees based on the provided trade information. ",
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

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-place.md>

### Place Order

Places an order for a specific account. For order execution of accounts with Omnibus with VA model, please submit virtual account number instead of omnibus master account number.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/place",
  "method": "post",
  "tags": [
    "Order"
  ],
  "description": "Places an order for a specific account. For order execution of accounts with Omnibus with VA model, please submit virtual account number instead of omnibus master account number.",
  "operationId": "brokerOrderPlace",
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
            "client_combo_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the combined order<br/> If combo_type = NORMAL and client_combo_order_id not need to set<br/> If combo_type != NORMAL and client_combo_order_id not provided <br/> the server will automatically generate one.",
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
                    "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) ",
                    "example": "NORMAL"
                  },
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "instrument_type": {
                    "type": "string",
                    "description": "Type of financial instrument associated with the request.",
                    "example": "EQUITY",
                    "enum": [
                      "EQUITY",
                      "OPTION"
                    ]
                  },
                  "market": {
                    "type": "string",
                    "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
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
                    "example": "BULL"
                  },
                  "order_type": {
                    "type": "string",
                    "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order<br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> HK Options<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
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
                      "MARKET_ON_CLOSE",
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
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                    "example": "DAY",
                    "enum": [
                      "DAY",
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
                    "description": "List of party identifiers. Applicable only for Hong Kong stock orders. Required for Relevant Regulated Intermediaries; should be omitted otherwise. ",
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
                      "description": "List of party identifiers. Applicable only for Hong Kong stock orders. Required for Relevant Regulated Intermediaries; should be omitted otherwise. ",
                      "title": "BrokerPartyId"
                    }
                  },
                  "option_strategy": {
                    "type": "string",
                    "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                    "example": "SINGLE",
                    "enum": [
                      "SINGLE",
                      "COVERED_STOCK",
                      "STRADDLE",
                      "STRANGLE",
                      "VERTICAL",
                      "CALENDAR",
                      "DIAGONAL",
                      "COLLAR_WITH_STOCK",
                      "BUTTERFLY",
                      "CONDOR",
                      "IRON_BUTTERFLY",
                      "IRON_CONDOR",
                      "RATIO"
                    ]
                  },
                  "position_intent": {
                    "type": "string",
                    "description": "Position intent for option orders, indicating whether the trade is opening or closing a position.<br/> &bull; BUY_TO_OPEN: Buy to open a new position<br/> &bull; BUY_TO_CLOSE: Buy to close an existing position<br/> &bull; SELL_TO_OPEN: Sell to open a new position<br/> &bull; SELL_TO_CLOSE: Sell to close an existing position",
                    "example": "BUY_TO_OPEN",
                    "enum": [
                      "BUY_TO_OPEN",
                      "BUY_TO_CLOSE",
                      "SELL_TO_OPEN",
                      "SELL_TO_CLOSE"
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
                            "OPTION"
                          ]
                        },
                        "market": {
                          "type": "string",
                          "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
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
                      "description": "Option leg detail. Only required when previewing option orders.",
                      "title": "BrokerOptionPlaceLegParam"
                    }
                  }
                },
                "description": "Order Details",
                "title": "BrokerOrderPlaceItemParam"
              }
            }
          },
          "title": "BrokerOrderPlaceParam"
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
            "title": "BrokerOrderWriteResult"
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
        "option_strategy": "SINGLE",
        "position_intent": "BUY_TO_OPEN",
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
      "content": "Places an order for a specific account. For order execution of accounts with Omnibus with VA model, please submit virtual account number instead of omnibus master account number.",
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

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-replace.md>

### Replace Order

Modifies a submitted order with updated parameters.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/replace",
  "method": "post",
  "tags": [
    "Order"
  ],
  "description": "Modifies a submitted order with updated parameters.",
  "operationId": "brokerOrderReplace",
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
                  "client_order_id",
                  "combo_type"
                ],
                "type": "object",
                "properties": {
                  "combo_type": {
                    "type": "string",
                    "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) ",
                    "example": "NORMAL"
                  },
                  "client_order_id": {
                    "type": "string",
                    "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
                    "example": "0KGOHL4PR2SLC0DKIND4TI0002"
                  },
                  "time_in_force": {
                    "type": "string",
                    "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                    "example": "DAY",
                    "enum": [
                      "DAY",
                      "GTC"
                    ]
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
                  "quantity": {
                    "type": "string",
                    "description": "Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks.",
                    "example": "1"
                  },
                  "legs": {
                    "type": "array",
                    "description": "Option leg detail for modifying option orders.",
                    "items": {
                      "required": [
                        "id",
                        "quantity"
                      ],
                      "type": "object",
                      "properties": {
                        "id": {
                          "type": "string",
                          "description": "Unique defined identifier for the leg.",
                          "example": "G2JAJPOR4KUA0F5I9LONH8J83A"
                        },
                        "quantity": {
                          "type": "string",
                          "description": "Quantity of the order or strategy leg. <br/>For stock legs, specifies the number of shares to transact <br/>For option legs, specifies the number of option contracts to transact for this leg and is expressed in whole contracts.",
                          "example": "1"
                        }
                      },
                      "description": "Option leg detail for modifying option orders.",
                      "title": "BrokerOptionReplaceLegParam"
                    }
                  }
                },
                "description": "Order Details",
                "title": "BrokerOrderReplaceItemParam"
              }
            }
          },
          "title": "BrokerOrderReplaceParam"
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
            "title": "BrokerOrderWriteResult"
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
        "combo_type": "NORMAL",
        "client_order_id": "0KGOHL4PR2SLC0DKIND4TI0002",
        "time_in_force": "DAY",
        "limit_price": "11.0",
        "stop_price": "11.0",
        "trailing_type": "AMOUNT",
        "trailing_stop_step": "1",
        "trailing_limit_price_offset": "11.0",
        "trigger_price_type": "PRICE",
        "quantity": "1",
        "legs": [
          {
            "id": "G2JAJPOR4KUA0F5I9LONH8J83A",
            "quantity": "1"
          }
        ]
      }
    ]
  },
  "postman": {
    "name": "Replace Order",
    "description": {
      "content": "Modifies a submitted order with updated parameters.",
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

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-cancel.md>

### Cancel Order

Cancels a submitted order that is still under open status.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/cancel",
  "method": "post",
  "tags": [
    "Order"
  ],
  "description": "Cancels a submitted order that is still under open status.",
  "operationId": "brokerOrderCancel",
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
            "account_id": {
              "type": "string",
              "description": "Account identifier",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "client_order_id": {
              "type": "string",
              "description": "Unique client-defined identifier for the order.<br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "0KGOHL4PR2SLC0DKIND4TI0002"
            }
          },
          "title": "BrokerOrderCancelParam"
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
            "title": "BrokerOrderWriteResult"
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
      "content": "Cancels a submitted order that is still under open status.",
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

## Order Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-detail.md>

### Get Order Detail

Retrieves the order details for a specific trade via the order ID.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/get",
  "method": "get",
  "tags": [
    "Order Query"
  ],
  "description": "Retrieves the order details for a specific trade via the order ID.",
  "operationId": "brokerOrderDetail",
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
      "description": "Unique client-defined identifier for the order. <br/>Maximum length is 32 characters and must be unique per user/account. <br/>Used to track or reference the order when interacting with the system.",
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
              "client_order_id",
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
                "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) ",
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
                      "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order<br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> HK Options<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
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
                        "MARKET_ON_CLOSE",
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
                      "example": "EQUITY",
                      "enum": [
                        "EQUITY",
                        "OPTION"
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
                      "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                      "example": "DAY",
                      "enum": [
                        "DAY",
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
                      "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT, ENHANCED_LIMIT, or AT_AUCTION_LIMIT or TOUCH_LMT.<br/>Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
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
                      "description": "Order placement time in ISO8601 format (UTC). Format: YYYY-MM-DDThh:mm:ssZ",
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
                      "description": "Time of the last executed trade in ISO8601 format (UTC). Format: YYYY-MM-DDThh:mm:ssZ",
                      "example": "2025-11-11T05:44:35.385Z"
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
                    },
                    "legs": {
                      "type": "array",
                      "description": "Leg detail",
                      "items": {
                        "required": [
                          "quantity",
                          "side",
                          "symbol"
                        ],
                        "type": "object",
                        "properties": {
                          "id": {
                            "type": "string",
                            "description": "Unique identifier for the leg.",
                            "example": "G2JAJPOR4KUA0F5I9LONH8J83A"
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
                            "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                            "example": "SINGLE",
                            "enum": [
                              "SINGLE",
                              "COVERED_STOCK",
                              "STRADDLE",
                              "STRANGLE",
                              "VERTICAL",
                              "CALENDAR",
                              "DIAGONAL",
                              "COLLAR_WITH_STOCK",
                              "BUTTERFLY",
                              "CONDOR",
                              "IRON_BUTTERFLY",
                              "IRON_CONDOR",
                              "RATIO"
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
                        "title": "BrokerOrderListLeg"
                      }
                    }
                  },
                  "description": "Order Details",
                  "title": "BrokerOrderDetailItem"
                }
              }
            },
            "title": "BrokerOrderDetailResult"
          },
          "examples": {
            "Example": {
              "description": "Example",
              "value": {
                "client_order_id": "41099381429646009294",
                "orders": [
                  {
                    "symbol": "AAPL",
                    "side": "BUY",
                    "status": "FILLED",
                    "legs": [
                      {
                        "id": "037RUBB7QO8F60KHJOIO000000",
                        "quantity": "100.000000",
                        "side": "BUY",
                        "symbol": "AAPL"
                      },
                      {
                        "id": "037RUBB7QO8F60KHJOIO000001",
                        "quantity": "1.000000",
                        "side": "SELL",
                        "symbol": "AAPL",
                        "option_type": "CALL",
                        "option_expire_date": "2026-09-18",
                        "strike_price": "315",
                        "option_category": "AMERICAN",
                        "option_contract_multiplier": "100",
                        "option_contract_deliverable": "100"
                      }
                    ],
                    "client_order_id": "41099381429646009294",
                    "order_type": "MARKET",
                    "instrument_type": "OPTION",
                    "order_id": "037RUBB7HG8F60KHJOIO000000",
                    "support_trading_session": "CORE",
                    "total_quantity": "1.000000",
                    "filled_quantity": "1.000000",
                    "place_time": "1784531308312",
                    "option_strategy": "COVERED_STOCK",
                    "place_time_at": "2026-07-20T07:08:28.312Z",
                    "filled_time": "1784531309477",
                    "filled_time_at": "2026-07-20T07:08:29.477Z",
                    "filled_price": "304.90",
                    "time_in_force": "DAY"
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
      "content": "Retrieves the order details for a specific trade via the order ID.",
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
            "content": "(Required) Account identifier.",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Unique client-defined identifier for the order. <br/>Maximum length is 32 characters and must be unique per user/account. <br/>Used to track or reference the order when interacting with the system.",
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

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-history.md>

### List Order History

Retrieves historical orders.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/historical-orders/list",
  "method": "get",
  "tags": [
    "Order Query"
  ],
  "description": "Retrieves historical orders.",
  "operationId": "brokerOrderHistory",
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
      "name": "start_date",
      "in": "query",
      "description": "The start date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 6 months.<br/> Format: yyyy-MM-dd.",
      "required": false,
      "schema": {
        "type": "String"
      },
      "example": "2024-09-25"
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
      "description": "Request successful",
      "content": {
        "application/json": {
          "schema": {
            "required": [
              "client_order_id",
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
                "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) ",
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
                      "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order<br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> HK Options<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
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
                        "MARKET_ON_CLOSE",
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
                      "example": "EQUITY",
                      "enum": [
                        "EQUITY",
                        "OPTION"
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
                      "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                      "example": "DAY",
                      "enum": [
                        "DAY",
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
                      "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT, ENHANCED_LIMIT, or AT_AUCTION_LIMIT or TOUCH_LMT.<br/>Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
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
                      "description": "Order placement time in ISO8601 format (UTC). Format: YYYY-MM-DDThh:mm:ssZ",
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
                      "description": "Time of the last executed trade in ISO8601 format (UTC). Format: YYYY-MM-DDThh:mm:ssZ",
                      "example": "2025-11-11T05:44:35.385Z"
                    },
                    "option_strategy": {
                      "type": "string",
                      "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                      "example": "SINGLE",
                      "enum": [
                        "SINGLE",
                        "COVERED_STOCK",
                        "STRADDLE",
                        "STRANGLE",
                        "VERTICAL",
                        "CALENDAR",
                        "DIAGONAL",
                        "COLLAR_WITH_STOCK",
                        "BUTTERFLY",
                        "CONDOR",
                        "IRON_BUTTERFLY",
                        "IRON_CONDOR",
                        "RATIO"
                      ]
                    },
                    "legs": {
                      "type": "array",
                      "description": "Leg detail",
                      "items": {
                        "required": [
                          "quantity",
                          "side",
                          "symbol"
                        ],
                        "type": "object",
                        "properties": {
                          "id": {
                            "type": "string",
                            "description": "Unique identifier for the leg.",
                            "example": "G2JAJPOR4KUA0F5I9LONH8J83A"
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
                            "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                            "example": "SINGLE",
                            "enum": [
                              "SINGLE",
                              "COVERED_STOCK",
                              "STRADDLE",
                              "STRANGLE",
                              "VERTICAL",
                              "CALENDAR",
                              "DIAGONAL",
                              "COLLAR_WITH_STOCK",
                              "BUTTERFLY",
                              "CONDOR",
                              "IRON_BUTTERFLY",
                              "IRON_CONDOR",
                              "RATIO"
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
                        "title": "BrokerOrderListLeg"
                      }
                    }
                  },
                  "description": "Order Details",
                  "title": "BrokerOrderListItem"
                }
              }
            },
            "title": "BrokerOrderListResult"
          },
          "examples": {
            "Example": {
              "description": "Example",
              "value": {
                "client_order_id": "41099381429646009294",
                "orders": [
                  {
                    "symbol": "AAPL",
                    "side": "BUY",
                    "status": "FILLED",
                    "legs": [
                      {
                        "id": "037RUBB7QO8F60KHJOIO000000",
                        "quantity": "100.000000",
                        "side": "BUY",
                        "symbol": "AAPL"
                      },
                      {
                        "id": "037RUBB7QO8F60KHJOIO000001",
                        "quantity": "1.000000",
                        "side": "SELL",
                        "symbol": "AAPL",
                        "option_type": "CALL",
                        "option_expire_date": "2026-09-18",
                        "strike_price": "315",
                        "option_category": "AMERICAN",
                        "option_contract_multiplier": "100",
                        "option_contract_deliverable": "100"
                      }
                    ],
                    "client_order_id": "41099381429646009294",
                    "order_type": "MARKET",
                    "instrument_type": "OPTION",
                    "order_id": "037RUBB7HG8F60KHJOIO000000",
                    "total_quantity": "1.000000",
                    "filled_quantity": "1.000000",
                    "place_time": "1784531308312",
                    "option_strategy": "COVERED_STOCK",
                    "place_time_at": "2026-07-20T07:08:28.312Z",
                    "filled_time": "1784531309477",
                    "filled_time_at": "2026-07-20T07:08:29.477Z",
                    "filled_price": "304.90",
                    "time_in_force": "DAY"
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
      "content": "Retrieves historical orders.",
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
            "content": "(Required) Account identifier.",
            "type": "text/plain"
          },
          "key": "account_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "The start date of the query period.<br/> If not provided, the default query period is the last 7 days.<br/> Users can specify an earlier date, but the maximum allowed look-back period is 6 months.<br/> Format: yyyy-MM-dd.",
            "type": "text/plain"
          },
          "key": "start_date",
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

## Open Orders

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-open.md>

### List Open Orders

Retrieves current pending orders list with pagination support.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/orders/open-orders/list",
  "method": "get",
  "tags": [
    "Order Query"
  ],
  "description": "Retrieves current pending orders list with pagination support.",
  "operationId": "brokerOrderOpen",
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
              "client_order_id",
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
                "description": "Specifies the type of order combination. For details, please refer to [Combo Order](/apis/docs/trade-api/stock/#combo-orders-us-only).<br/> &bull; NORMAL: A standard single order.<br/> &bull; MASTER: A primary order that triggers a take-profit or stop-loss order upon execution <br/> &bull; STOP_PROFIT: A take-profit order <br/> &bull; STOP_LOSS: A stop-loss order <br/> &bull; OTO: An order that triggers another order upon execution (One-Triggers-the-Other) <br/> &bull; OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) <br/> &bull; OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) ",
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
                      "description": "Specifies the type of order to be placed. Determines how the order will be executed in the market.<br/> Available order types depend on the market and instrument type.<br/> Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported.<br/> U.S. Stock<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> Hong Kong Stock<br/> &nbsp; &bull; <b>ENHANCED_LIMIT:</b> Enhanced Limit Order<br/> &nbsp; &bull; <b>AT_AUCTION:</b> At-auction order<br/> &nbsp; &bull; <b>AT_AUCTION_LIMIT:</b> At-auction limit order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order<br/> &nbsp; &bull; <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order<br/> &nbsp; &bull; <b>TOUCH_MKT:</b> Touch Market Order<br/> &nbsp; &bull; <b>TOUCH_LMT:</b> Touch Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order<br/> &nbsp; &bull; <b>MARKET_ON_OPEN:</b> Opening market order<br/> &nbsp; &bull; <b>MARKET_ON_CLOSE:</b> Closing market order<br/> HK Options<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> &nbsp; &bull; <b>MARKET:</b> Market Order<br/> &nbsp; &bull; <b>STOP_LOSS:</b> Stop Order<br/> &nbsp; &bull; <b>STOP_LOSS_LIMIT:</b> Stop Limit Order<br/> China Connect<br/> &nbsp; &bull; <b>LIMIT:</b> Limit Order<br/> ",
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
                        "MARKET_ON_CLOSE",
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
                      "example": "EQUITY",
                      "enum": [
                        "EQUITY",
                        "OPTION"
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
                      "description": "Specifies the duration for which the order remains active in the market (Time-In-Force). <br/> &bull; DAY: The order is valid only for the current trading day and expires at the end of the day. <br/> &bull; GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days).",
                      "example": "DAY",
                      "enum": [
                        "DAY",
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
                      "description": "Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT, ENHANCED_LIMIT, or AT_AUCTION_LIMIT or TOUCH_LMT.<br/>Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed.",
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
                      "description": "Order placement time in ISO8601 format (UTC). Format: YYYY-MM-DDThh:mm:ssZ",
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
                      "description": "Time of the last executed trade in ISO8601 format (UTC). Format: YYYY-MM-DDThh:mm:ssZ",
                      "example": "2025-11-11T05:44:35.385Z"
                    },
                    "option_strategy": {
                      "type": "string",
                      "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                      "example": "SINGLE",
                      "enum": [
                        "SINGLE",
                        "COVERED_STOCK",
                        "STRADDLE",
                        "STRANGLE",
                        "VERTICAL",
                        "CALENDAR",
                        "DIAGONAL",
                        "COLLAR_WITH_STOCK",
                        "BUTTERFLY",
                        "CONDOR",
                        "IRON_BUTTERFLY",
                        "IRON_CONDOR",
                        "RATIO"
                      ]
                    },
                    "legs": {
                      "type": "array",
                      "description": "Leg detail",
                      "items": {
                        "required": [
                          "quantity",
                          "side",
                          "symbol"
                        ],
                        "type": "object",
                        "properties": {
                          "id": {
                            "type": "string",
                            "description": "Unique identifier for the leg.",
                            "example": "G2JAJPOR4KUA0F5I9LONH8J83A"
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
                            "description": "Type of options strategy <br/> &bull; SINGLE: Single-leg options order<br/> &bull; COVERED_STOCK: Covered call (stock + sell CALL)<br/> &bull; STRADDLE: Same strike, same expiry CALL + PUT<br/> &bull; STRANGLE: Different strike, same expiry CALL + PUT<br/> &bull; VERTICAL: Same expiry, different strike, same type (Bull/Bear Spread)<br/> &bull; CALENDAR: Same strike, different expiry, same type<br/> &bull; DIAGONAL: Different strike, different expiry, same type<br/> &bull; COLLAR_WITH_STOCK: Stock + buy PUT + sell CALL<br/> &bull; BUTTERFLY: Buy low + sell 2 mid + buy high (same type)<br/> &bull; CONDOR: Four different strikes, same type<br/> &bull; IRON_BUTTERFLY: Buy low PUT + sell mid PUT + sell mid CALL + buy high CALL<br/> &bull; IRON_CONDOR: Buy low PUT + sell mid-low PUT + sell mid-high CALL + buy high CALL<br/> &bull; RATIO: Same expiry, same option type, unequal leg quantities (e.g. buy 1 low-strike CALL + sell 2 high-strike CALLs)",
                            "example": "SINGLE",
                            "enum": [
                              "SINGLE",
                              "COVERED_STOCK",
                              "STRADDLE",
                              "STRANGLE",
                              "VERTICAL",
                              "CALENDAR",
                              "DIAGONAL",
                              "COLLAR_WITH_STOCK",
                              "BUTTERFLY",
                              "CONDOR",
                              "IRON_BUTTERFLY",
                              "IRON_CONDOR",
                              "RATIO"
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
                        "title": "BrokerOrderListLeg"
                      }
                    }
                  },
                  "description": "Order Details",
                  "title": "BrokerOrderListItem"
                }
              }
            },
            "title": "BrokerOrderListResult"
          },
          "examples": {
            "Example": {
              "description": "Example",
              "value": [
                {
                  "client_order_id": "93424901665872540803",
                  "orders": [
                    {
                      "symbol": "AAPL",
                      "side": "BUY",
                      "status": "SUBMITTED",
                      "legs": [
                        {
                          "id": "037RUB9NBE8F60KHJOIO000000",
                          "quantity": "100.000000",
                          "side": "BUY",
                          "symbol": "AAPL"
                        },
                        {
                          "id": "037RUB9NBE8F60KHJOIO000001",
                          "quantity": "1.000000",
                          "side": "SELL",
                          "symbol": "AAPL",
                          "option_type": "CALL",
                          "option_expire_date": "2026-09-18",
                          "strike_price": "315",
                          "option_category": "AMERICAN",
                          "option_contract_multiplier": "100",
                          "option_contract_deliverable": "100"
                        }
                      ],
                      "client_order_id": "93424901665872540803",
                      "order_type": "LIMIT",
                      "instrument_type": "OPTION",
                      "order_id": "037RUB9N068F60KHJOIO000000",
                      "total_quantity": "1.000000",
                      "filled_quantity": "0.000000",
                      "place_time": "1784531283459",
                      "option_strategy": "COVERED_STOCK",
                      "place_time_at": "2026-07-20T07:08:03.459Z",
                      "filled_price": "0.00",
                      "time_in_force": "DAY",
                      "limit_price": "2.80"
                    }
                  ]
                }
              ]
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
    "name": "List Open Orders",
    "description": {
      "content": "Retrieves current pending orders list with pagination support.",
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

## Get FX Rate

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-query-rate.md>

### Get FX Rate

Retrieves the currency exchange rate based on currency pair.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/fx-rates/get",
  "method": "get",
  "tags": [
    "FX"
  ],
  "description": "Retrieves the currency exchange rate based on currency pair.",
  "operationId": "brokerFundingQueryRate",
  "parameters": [
    {
      "name": "from_currency",
      "in": "query",
      "description": "Which currency you use for exchange.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Currency",
        "enum": [
          "CNH",
          "HKD",
          "USD"
        ]
      }
    },
    {
      "name": "to_currency",
      "in": "query",
      "description": "Target currency you plan to get for exchange. ",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Currency",
        "enum": [
          "CNH",
          "HKD",
          "USD"
        ]
      }
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
              "from_currency",
              "fx_rate",
              "rate_effective_time",
              "rate_expire_time",
              "to_currency"
            ],
            "type": "object",
            "properties": {
              "from_currency": {
                "type": "string",
                "description": "Currency",
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "to_currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "fx_rate": {
                "type": "string",
                "description": "FX rate",
                "example": "8.01"
              },
              "rate_effective_time": {
                "type": "string",
                "description": "Effective time of the FX rate (ISO 8601 format, UTC time zone).",
                "example": "2025-11-11T05:44:35.385Z"
              },
              "rate_expire_time": {
                "type": "string",
                "description": "Expiration time of the FX rate (ISO 8601 format, UTC time zone).",
                "example": "2025-11-11T05:44:35.385Z"
              }
            },
            "title": "BrokerExchangeRateResult"
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
    "name": "Get FX Rate",
    "description": {
      "content": "Retrieves the currency exchange rate based on currency pair.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "fx-rates",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Which currency you use for exchange.",
            "type": "text/plain"
          },
          "key": "from_currency",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Target currency you plan to get for exchange. ",
            "type": "text/plain"
          },
          "key": "to_currency",
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

## Create FX Request

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-create-fx.md>

### Create FX Exchange

Creates a currency exchange request at the Omnibus master account level. Requires base currency, base amount and target currency as inputs. Uses the latest FX rate for currency exchange processing.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/fx-exchanges/create",
  "method": "post",
  "tags": [
    "FX"
  ],
  "description": "Creates a currency exchange request at the Omnibus master account level. Requires base currency, base amount and target currency as inputs. Uses the latest FX rate for currency exchange processing.",
  "operationId": "brokerFundingCreateFX",
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
            "from_amount",
            "from_currency",
            "to_currency"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier, Exchange on which account id.",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "from_amount": {
              "type": "string",
              "description": "Base amount for exchange.",
              "example": "100"
            },
            "from_currency": {
              "type": "string",
              "description": "Currency",
              "example": "USD",
              "enum": [
                "CNH",
                "HKD",
                "USD"
              ]
            },
            "to_currency": {
              "type": "string",
              "description": "Currency",
              "example": "HKD",
              "enum": [
                "CNH",
                "HKD",
                "USD"
              ]
            }
          },
          "title": "BrokerCreateExchangeParam"
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
              "from_amount",
              "from_currency",
              "fx_id",
              "status",
              "to_currency"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Account identifier, Exchange on which account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "LJIS16BACHQG9LPP44L9IQHGAB"
              },
              "fx_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "from_amount": {
                "type": "string",
                "description": "Base amount for exchange.",
                "example": "100"
              },
              "from_currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "to_currency": {
                "type": "string",
                "description": "Currency",
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
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
              }
            },
            "title": "BrokerCreateExchangeResult"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "from_amount": "100",
    "from_currency": "USD",
    "to_currency": "HKD"
  },
  "postman": {
    "name": "Create FX Exchange",
    "description": {
      "content": "Creates a currency exchange request at the Omnibus master account level. Requires base currency, base amount and target currency as inputs. Uses the latest FX rate for currency exchange processing.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "fx-exchanges",
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

## Get FX Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-query-fx.md>

### Get FX Exchange Detail

Retrieves currency exchange record details at the Omnibus master account level.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/fx-exchanges/get",
  "method": "get",
  "tags": [
    "FX"
  ],
  "description": "Retrieves currency exchange record details at the Omnibus master account level.",
  "operationId": "brokerFundingQueryFX",
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
              "account_id",
              "client_request_id",
              "from_currency",
              "fx_id",
              "status",
              "to_currency"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Account identifier, Exchange on which account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "LJIS16BACHQG9LPP44L9IQHGAB"
              },
              "fx_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "from_amount": {
                "type": "string",
                "description": "Base amount for exchange.",
                "example": "100"
              },
              "to_amount": {
                "type": "string",
                "description": "Amount after exchange.",
                "example": "700"
              },
              "from_currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "to_currency": {
                "type": "string",
                "description": "Currency",
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "fx_rate": {
                "type": "string",
                "description": "FX rate",
                "example": "8.01"
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
              }
            },
            "title": "BrokerQueryExchangeResult"
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
    "name": "Get FX Exchange Detail",
    "description": {
      "content": "Retrieves currency exchange record details at the Omnibus master account level.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "fx-exchanges",
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

## Create Instant Exchange

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-create-instant-fx.md>

### Create Instant Exchange

Creates an instant currency exchange request at the Virtual account level. Requires base currency, base amount, target currency, and target amount as inputs. This feature does not rely on the latest FX rate provided by Webull.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/instant-exchanges/create",
  "method": "post",
  "tags": [
    "FX"
  ],
  "description": "Creates an instant currency exchange request at the Virtual account level. Requires base currency, base amount, target currency, and target amount as inputs. This feature does not rely on the latest FX rate provided by Webull.",
  "operationId": "brokerFundingCreateInstantFX",
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
            "from_amount",
            "from_currency",
            "to_amount",
            "to_currency"
          ],
          "type": "object",
          "properties": {
            "account_id": {
              "type": "string",
              "description": "Account identifier, Exchange on which account id.",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "from_amount": {
              "type": "string",
              "description": "Base amount for exchange.",
              "example": "100"
            },
            "to_amount": {
              "type": "string",
              "description": "Amount after exchange.",
              "example": "700"
            },
            "from_currency": {
              "type": "string",
              "description": "Currency",
              "example": "USD",
              "enum": [
                "CNH",
                "HKD",
                "USD"
              ]
            },
            "to_currency": {
              "type": "string",
              "description": "Currency",
              "example": "HKD",
              "enum": [
                "CNH",
                "HKD",
                "USD"
              ]
            }
          },
          "title": "BrokerCreateInstantExchangeParam"
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
              "from_amount",
              "from_currency",
              "instant_exchange_id",
              "status",
              "to_amount",
              "to_currency"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Account identifier, Exchange on which account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "LJIS16BACHQG9LPP44L9IQHGAB"
              },
              "instant_exchange_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "from_amount": {
                "type": "string",
                "description": "Base amount for exchange.",
                "example": "100"
              },
              "to_amount": {
                "type": "string",
                "description": "Amount after exchange.",
                "example": "700"
              },
              "from_currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "to_currency": {
                "type": "string",
                "description": "Currency",
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
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
              }
            },
            "title": "BrokerInstantExchangeResult"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "from_amount": "100",
    "to_amount": "700",
    "from_currency": "USD",
    "to_currency": "HKD"
  },
  "postman": {
    "name": "Create Instant Exchange",
    "description": {
      "content": "Creates an instant currency exchange request at the Virtual account level. Requires base currency, base amount, target currency, and target amount as inputs. This feature does not rely on the latest FX rate provided by Webull.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "instant-exchanges",
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

## Get Instant Exchange Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-query-instant-fx.md>

### Get Instant Exchange Detail

Retrieves instant currency exchange record details at the Omnibus master account level.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/instant-exchanges/get",
  "method": "get",
  "tags": [
    "FX"
  ],
  "description": "Retrieves instant currency exchange record details at the Omnibus master account level.",
  "operationId": "brokerFundingQueryInstantFX",
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
              "account_id",
              "client_request_id",
              "from_amount",
              "from_currency",
              "instant_exchange_id",
              "status",
              "to_amount",
              "to_currency"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Account identifier, Exchange on which account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "LJIS16BACHQG9LPP44L9IQHGAB"
              },
              "instant_exchange_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "from_amount": {
                "type": "string",
                "description": "Base amount for exchange.",
                "example": "100"
              },
              "to_amount": {
                "type": "string",
                "description": "Amount after exchange.",
                "example": "700"
              },
              "from_currency": {
                "type": "string",
                "description": "Currency",
                "example": "USD",
                "enum": [
                  "CNH",
                  "HKD",
                  "USD"
                ]
              },
              "to_currency": {
                "type": "string",
                "description": "Currency",
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
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
              }
            },
            "title": "BrokerInstantExchangeResult"
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
    "name": "Get Instant Exchange Detail",
    "description": {
      "content": "Retrieves instant currency exchange record details at the Omnibus master account level.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "funding",
        "instant-exchanges",
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

## Create Instant Funding

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-instant-create.md>

### Create Instant Funding Request

Creates an instant funding request at the Virtual account level.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/instant-funding/create",
  "method": "post",
  "tags": [
    "Instant Funding"
  ],
  "description": "Creates an instant funding request at the Virtual account level.",
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
            "account_id": {
              "type": "string",
              "description": "Account identifier, Exchange on which account id.",
              "example": "93IUJ28O9VO2KBGHDHR4H9"
            },
            "type": {
              "type": "string",
              "description": "Instant Funding Request Type",
              "example": "DEPOSIT",
              "enum": [
                "DEPOSIT",
                "WITHDRAWAL"
              ]
            },
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
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
              "example": "HKD",
              "enum": [
                "CNH",
                "HKD",
                "USD"
              ]
            }
          },
          "title": "BrokerCreateInstantFundingParam"
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
              "currency",
              "instant_funding_id",
              "status",
              "type"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Account identifier, Exchange on which account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "type": {
                "type": "string",
                "description": "Instant Funding Request Type",
                "example": "DEPOSIT",
                "enum": [
                  "DEPOSIT",
                  "WITHDRAWAL"
                ]
              },
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
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
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
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
              }
            },
            "title": "BrokerInstantFundingResult"
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
    "account_id": "93IUJ28O9VO2KBGHDHR4H9",
    "type": "DEPOSIT",
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "amount": "100",
    "currency": "HKD"
  },
  "postman": {
    "name": "Create Instant Funding Request",
    "description": {
      "content": "Creates an instant funding request at the Virtual account level.",
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

## Get Instant Funding Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-instant-query.md>

### Get Instant Funding Detail

Retrieves instant funding record details at the Virtual account level.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/funding/instant-funding/get",
  "method": "get",
  "tags": [
    "Instant Funding"
  ],
  "description": "Retrieves instant funding record details at the Virtual account level.",
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
              "account_id",
              "amount",
              "client_request_id",
              "currency",
              "instant_funding_id",
              "status",
              "type"
            ],
            "type": "object",
            "properties": {
              "account_id": {
                "type": "string",
                "description": "Account identifier, Exchange on which account id.",
                "example": "93IUJ28O9VO2KBGHDHR4H9"
              },
              "type": {
                "type": "string",
                "description": "Instant Funding Request Type",
                "example": "DEPOSIT",
                "enum": [
                  "DEPOSIT",
                  "WITHDRAWAL"
                ]
              },
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
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
                "example": "HKD",
                "enum": [
                  "CNH",
                  "HKD",
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
              }
            },
            "title": "BrokerInstantFundingResult"
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
      "content": "Retrieves instant funding record details at the Virtual account level.",
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

## Create Cash Journal

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-cash-create.md>

### Create Cash Journal

Creates a cash journal request. Cash journal supports both Omnibus master account level and Virtual account level. 

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/journals/cash-journals/create",
  "method": "post",
  "tags": [
    "Journals"
  ],
  "description": "Creates a cash journal request. Cash journal supports both Omnibus master account level and Virtual account level. ",
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
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
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
                "CNH",
                "HKD",
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
              "currency",
              "from_account",
              "journal_id",
              "status",
              "to_account"
            ],
            "type": "object",
            "properties": {
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
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
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
                  "CNH",
                  "HKD",
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
    "from_account": "93IUJ28O9VO2KBGHDHR4H9",
    "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "amount": "100",
    "currency": "USD"
  },
  "postman": {
    "name": "Create Cash Journal",
    "description": {
      "content": "Creates a cash journal request. Cash journal supports both Omnibus master account level and Virtual account level. ",
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

## Get Cash Journal Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-cash-query.md>

### Get Cash Journal Detail

Retrieves detailed information for a cash journal request.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/journals/cash-journals/get",
  "method": "get",
  "tags": [
    "Journals"
  ],
  "description": "Retrieves detailed information for a cash journal request.",
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
              "currency",
              "from_account",
              "journal_id",
              "status",
              "to_account"
            ],
            "type": "object",
            "properties": {
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
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
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
                  "CNH",
                  "HKD",
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
      "content": "Retrieves detailed information for a cash journal request.",
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

## Create Position Journal

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-position-create.md>

### Create Position Journal

Creates a position journal request. Position journal only supports Virtual account level. 

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/journals/position-journals/create",
  "method": "post",
  "tags": [
    "Journals"
  ],
  "description": "Creates a position journal request. Position journal only supports Virtual account level. ",
  "operationId": "brokerJournalPositionCreate",
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
            "from_account",
            "instrument_type",
            "market",
            "quantity",
            "symbol",
            "to_account"
          ],
          "type": "object",
          "properties": {
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
            "client_request_id": {
              "type": "string",
              "description": "Client Request ID, unique for each request. <br/> Maximum length is 32 characters. <br/> Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_).",
              "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
            },
            "instrument_type": {
              "type": "string",
              "description": "Type of financial instrument associated with the request.",
              "example": "EQUITY",
              "enum": [
                "EQUITY",
                "OPTION"
              ]
            },
            "market": {
              "type": "string",
              "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
              "example": "US",
              "enum": [
                "US",
                "HK",
                "CN"
              ]
            },
            "quantity": {
              "type": "string",
              "description": "Quantity, positive number.",
              "example": "100"
            },
            "symbol": {
              "type": "string",
              "description": "Symbol name. Options use OCC format.",
              "example": "BULL"
            }
          },
          "title": "BrokerCreatePositionJournalParam"
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
              "client_request_id",
              "from_account",
              "instrument_type",
              "journal_id",
              "market",
              "quantity",
              "status",
              "symbol",
              "to_account"
            ],
            "type": "object",
            "properties": {
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
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "journal_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "instrument_type": {
                "type": "string",
                "description": "Type of financial instrument associated with the request.",
                "example": "EQUITY",
                "enum": [
                  "EQUITY",
                  "OPTION"
                ]
              },
              "market": {
                "type": "string",
                "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                "example": "US",
                "enum": [
                  "US",
                  "HK",
                  "CN"
                ]
              },
              "quantity": {
                "type": "string",
                "description": "Quantity, positive number.",
                "example": "100"
              },
              "symbol": {
                "type": "string",
                "description": "Symbol name. Options use OCC format.",
                "example": "BULL"
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
                "description": "Reason for failure. Returned when the terminal state is not COMPLETED.",
                "example": "Failed."
              }
            },
            "title": "BrokerPositionJournalResult"
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
    "from_account": "93IUJ28O9VO2KBGHDHR4H9",
    "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
    "client_request_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
    "instrument_type": "EQUITY",
    "market": "US",
    "quantity": "100",
    "symbol": "BULL"
  },
  "postman": {
    "name": "Create Position Journal",
    "description": {
      "content": "Creates a position journal request. Position journal only supports Virtual account level. ",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "journals",
        "position-journals",
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

## Get Position Journal Detail

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-position-query.md>

### Get Position Journal Detail

Retrieves detailed information for a position journal request.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/journals/position-journals/get",
  "method": "get",
  "tags": [
    "Journals"
  ],
  "description": "Retrieves detailed information for a position journal request.",
  "operationId": "brokerJournalPositionQuery",
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
              "client_request_id",
              "from_account",
              "instrument_type",
              "journal_id",
              "market",
              "quantity",
              "status",
              "symbol",
              "to_account"
            ],
            "type": "object",
            "properties": {
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
              "client_request_id": {
                "type": "string",
                "description": "Client Request ID, unique for each request. ",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "journal_id": {
                "type": "string",
                "description": "Request id generated by the system.",
                "example": "89JGKD4LKIVI5UU6L3KHNU10IA"
              },
              "instrument_type": {
                "type": "string",
                "description": "Type of financial instrument associated with the request.",
                "example": "EQUITY",
                "enum": [
                  "EQUITY",
                  "OPTION"
                ]
              },
              "market": {
                "type": "string",
                "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
                "example": "US",
                "enum": [
                  "US",
                  "HK",
                  "CN"
                ]
              },
              "quantity": {
                "type": "string",
                "description": "Quantity, positive number.",
                "example": "100"
              },
              "symbol": {
                "type": "string",
                "description": "Symbol name. Options use OCC format.",
                "example": "BULL"
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
                "description": "Reason for failure. Returned when the terminal state is not COMPLETED.",
                "example": "Failed."
              }
            },
            "title": "BrokerPositionJournalResult"
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
    "name": "Get Position Journal Detail",
    "description": {
      "content": "Retrieves detailed information for a position journal request.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "broker",
        "journals",
        "position-journals",
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

## Trade Calendar

> Source: <https://developer.webull.hk/apis/docs/reference/broker-api/broker-trade-calendar.md>

### List Trade Calendar

Retrieves trading and settlement calendar based on markets and product types.

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
      "url": "https://broker-api.sandbox.webull.hk"
    }
  ],
  "path": "/broker/master-data/trading-calendars/list",
  "method": "get",
  "tags": [
    "Master"
  ],
  "description": "Retrieves trading and settlement calendar based on markets and product types.",
  "operationId": "brokerTradeCalendar",
  "parameters": [
    {
      "name": "market",
      "in": "query",
      "description": "Country code, US,HK,etc.",
      "required": true,
      "schema": {
        "type": "string",
        "description": "Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument.",
        "enum": [
          "US",
          "HK",
          "CN"
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
          "OPTION"
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
      "content": "Retrieves trading and settlement calendar based on markets and product types.",
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
            "content": "(Required) Country code, US,HK,etc.",
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

## Account Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-account-events.md>

### Account Events

Receive notifications for account-related events, including:
- Account restrictions and status changes
- Virtual Account (VA) BCAN code approval status changes
- Account closure processing

#### Account Update Event

This event is triggered when an Omni-level account's restrictions or status changes.

##### Event Structure

```json
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "ACCOUNT",
    "position": "CJO1fxACGAAgADAB",
    "timestamp": "2025-03-29T07:02:33.200Z",
    "payload": {
        "account_id": "036UUF6TQA8CD0KHJPEC000000",
        "account_number": "5KT05001",
        "account_type": "CASH",
        "restrictions": [
            "HK_STOCK_NO_TRADE",
            "HK_STOCK_LIQUIDATE_ONLY",
            "US_STOCK_NO_TRADE",
            "US_STOCK_LIQUIDATE_ONLY",
            "CN_STOCK_LIQUIDATE_ONLY",
            "CN_STOCK_NO_TRADE"
        ],
        "update_time": "2025-03-29T07:02:33.200Z",
        "biz_type": "ACCOUNT_UPDATE"
    }
}
```

##### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique event identifier |
| `position` | string |  Cursor for event replay; re-pushes subsequent events within the current business type.|
| `event_type` | string | Event type, fixed as `ACCOUNT` |
| `timestamp` | string | Event timestamp in ISO 8601 format |
| `payload` | object | Event payload data |

##### Payload Fields

| Field            | Type | Description                                           |
|------------------|------|-------------------------------------------------------|
| `account_id`     | string | Account id                                            |
| `account_number` | string | Account number                                        |
| `account_type`   | string | Account type: `CASH` / `MARGIN`                       |
| `restrictions`   | array | Array of current account restrictions (full snapshot) |
| `update_time`    | string | Update timestamp in ISO 8601 format                   |
| `biz_type`       | string | Business type, fixed as `ACCOUNT_UPDATE`              |

##### Restriction Types

The `restrictions` array contains current account restrictions. Each push provides a complete snapshot of all active restrictions:

| Restriction Code | Description |
|------------------|-------------|
| `HK_STOCK_NO_TRADE` | Hong Kong stocks - No trading allowed |
| `HK_STOCK_LIQUIDATE_ONLY` | Hong Kong stocks - Liquidate only (can only close positions) |
| `US_STOCK_NO_TRADE` | US stocks - No trading allowed |
| `US_STOCK_LIQUIDATE_ONLY` | US stocks - Liquidate only (can only close positions) |
| `CN_STOCK_NO_TRADE` | China A-shares - No trading allowed |
| `CN_STOCK_LIQUIDATE_ONLY` | China A-shares - Liquidate only (can only close positions) |

- Each event push contains a **full snapshot** of current restrictions, not incremental changes
- An empty `restrictions` array means no restrictions are currently active on the account

#### Virtual Account BCAN Update Event

This event is triggered when a Virtual Account's (VA) BCAN code approval status changes at the exchange.

##### Event Structure - Approval Success

When the exchange approves the BCAN code, the corresponding category trading permission is granted for the VA account.

```json
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "ACCOUNT",
    "position": "CJO1fxACGAAgADAB",
    "timestamp": "2025-03-29T07:02:33.200Z",
    "payload": {
       "account_id": "O3I8NNB882S2BK7MF1ACJ5UTK9",
       "account_number": "VA0000001",
        "account_type": "MARGIN",
        "type": "CN",
        "status": "SUCCESS",
        "update_time": "2025-03-29T07:02:33.200Z",
        "biz_type": "VIRTUAL_ACCOUNT_BCAN_UPDATE"
    }
}
```

##### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique event identifier |
| `position` | string |  Cursor for event replay; re-pushes subsequent events within the current business type.|
| `event_type` | string | Event type, fixed as `ACCOUNT` |
| `timestamp` | string | Event timestamp in ISO 8601 format |
| `payload` | object | Event payload data |

##### Payload Fields

| Field            | Type | Description                                                 |
|------------------|------|-------------------------------------------------------------|
| `account_id`     | string | Virtual Account id                                          |
| `account_number` | string | Virtual Account number                                      |
| `account_type`   | string | Account type: `CASH` / `MARGIN`                             |
| `type`           | string | BCAN type, currently only `CN` (China A-share) is supported |
| `status`         | string | BCAN status: `SUCCESS` / `FAILED`                           |
| `update_time`    | string | Update timestamp in ISO 8601 format                         |
| `biz_type`       | string | Business type, fixed as `VIRTUAL_ACCOUNT_BCAN_UPDATE`       |

##### BCAN Status Values

| Status | Description |
|--------|-------------|
| `SUCCESS` | BCAN code approved by exchange - Trading permission granted |
| `FAILED` | BCAN code rejected by exchange - Trading permission denied |

##### Event Structure - Approval Failed

When the exchange rejects the BCAN code, the corresponding category trading permission is denied for the VA account.

```json
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "ACCOUNT",
    "position": "CJO1fxACGAAgADAB",
    "timestamp": "2025-03-29T07:02:33.200Z",
    "payload": {
        "account_id": "O3I8NNB882S2BK7MF1ACJ5UTK9",
        "account_number": "VA0000001",
        "account_type": "MARGIN",
        "type": "CN",
        "status": "FAILED",
        "update_time": "2025-03-29T07:02:33.200Z",
        "biz_type": "VIRTUAL_ACCOUNT_BCAN_UPDATE"
    }
}
```

The field structure is identical to the approval success case, with `status` set to `FAILED`.

##### Use Cases

Listening to these events can be used for:

1. **Account Status Monitoring**: Track real-time account restriction changes and update UI/access controls accordingly
2. **Trading Permission Management**: Automatically enable or disable trading features based on account restrictions
3. **BCAN Status Tracking**: Monitor Virtual Account BCAN approval status for China A-share trading
4. **Risk Management**: Adjust risk controls and position management when liquidate-only restrictions are applied
5. **Client Notification**: Alert clients about account status changes requiring attention or action
6. **Compliance Monitoring**: Track account restrictions for regulatory and compliance purposes

## Instrument Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-instrument-events.md>

### Instrument Events

Receive notifications for instrument property changes, including:
- Trading permission changes (Tradable/Liquidate only/Non-Tradable)
- Shortable status changes (Shortable/Non-Shortable)
- Marginable status changes
- ETF properties changes (crypto ETF, leveraged ETF, etc.)

#### Trading Property Change Event

This event is triggered when an instrument's trading-related properties are updated.

##### Event Structure

```json
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "INSTRUMENT",
    "position": "CJO1fxACGAAgADAB",
    "timestamp": "2025-03-29T07:02:33.200Z",
    "payload": {
        "instrument_id": "10152734329",
        "status": "NT",
        "shortable": "true",
        "marginable": "true",
        "biz_type": "PROPERTY_CHANGE"
    }
}
```

##### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique event identifier |
| `position` | string |  Cursor for event replay; re-pushes subsequent events within the current business type.|
| `event_type` | string | Event type, fixed as `INSTRUMENT` |
| `timestamp` | string | Event timestamp in ISO 8601 format |
| `payload` | object | Event payload data |

##### Payload Fields

| Field | Type | Description |
|-------|------|-------------|
| `instrument_id` | string | Unique instrument identifier |
| `status` | string | Trading status:<br/>- `OC`: Tradable<br/>- `CO`: Liquidate only<br/>- `NT`: Non-Tradable |
| `shortable` | string | Whether shortable: `true` / `false` |
| `marginable` | string | Whether margin trading supported: `true` / `false` |
| `biz_type` | string | Business type, fixed as `PROPERTY_CHANGE` |

##### Status Values

Trading status (`status`) values:

| Value | Description |
|-------|-------------|
| `OC` | Tradable - Instrument can be bought and sold normally |
| `CO` | Liquidate only - Can only sell existing positions, cannot open new positions |
| `NT` | Non-Tradable - No trading operations allowed |

##### Use Cases

Listening to this event can be used for:

1. **Trading Permission Monitoring**: Get real-time updates on instrument trading permission changes and update UI displays accordingly
2. **Short Status Tracking**: Monitor changes in shortable status and adjust trading strategies
3. **Risk Management**: Adjust positions and risk controls based on margin requirement changes

---

#### Basic Property Change Event

This event is triggered when an instrument's basic properties (such as ETF properties) are updated.

##### Event Structure

```json
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "INSTRUMENT",
    "position": "CJO1fxACGAAgADAB",
    "timestamp": "2025-03-29T07:02:33.200Z",
    "payload": {
        "instrument_id": "10152734329",
        "crypto_etf": "false",
        "etf_leveraged_flag": "YES",
        "etf_leveraged_factor": "-1.0",
        "single_stock_etf": "false",
        "inverse_etf": "false",
        "biz_type": "BASIC_PROPERTY_CHANGE"
    }
}
```

##### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique event identifier |
| `position` | string |  Cursor for event replay; re-pushes subsequent events within the current business type.|
| `event_type` | string | Event type, fixed as `INSTRUMENT` |
| `timestamp` | string | Event timestamp in ISO 8601 format |
| `payload` | object | Event payload data |

##### Payload Fields

| Field | Type | Description |
|-------|------|-------------|
| `instrument_id` | string | Unique instrument identifier |
| `crypto_etf` | string | Whether crypto ETF: `true` / `false`. Empty if not ETF |
| `etf_leveraged_flag` | string | Whether leveraged ETF: `YES` / `NO`. Empty if not ETF |
| `etf_leveraged_factor` | string | ETF leverage factor. Empty if not ETF. Negative value for inverse ETF, positive value for regular leveraged ETF |
| `single_stock_etf` | string | Whether single stock ETF: `true` / `false`. Empty if not ETF |
| `inverse_etf` | string | Whether inverse ETF: `true` / `false`. Empty if not ETF |
| `biz_type` | string | Business type, fixed as `BASIC_PROPERTY_CHANGE` |

##### Use Cases

Listening to this event can be used for:

1. **ETF Product Monitoring**: Track updates to ETF-specific properties
2. **Risk Classification**: Identify high-risk products like leveraged or inverse ETFs
3. **Product Filtering**: Filter instruments based on ETF characteristics (crypto ETF, single stock ETF, etc.)
4. **Investment Strategy**: Adjust strategies based on leverage factors and ETF types

## Corporate Actions Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-ca-events.md>

### Corporate Actions Events

Receive notifications for corporate action events, including:
- Identifier changes (symbol, CUSIP, ISIN changes)
- Dividends (cash dividends, stock dividends)
- Stock splits and reverse splits
- Mergers and acquisitions
- Rights offerings
- Other corporate action events

#### Corporate Action Event

This event is triggered when a corporate action is announced or updated.

##### Event Structure

```json
{
  "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
  "event_type": "INSTRUMENT",
  "position": "CJO1fxACGAAgADAB",
  "timestamp": "2025-03-29T07:02:33.200Z",
  "payload": {
    "event_id": "925299289_1569116747",
    "event_type": "DIVIDEND",
    "event_version": "1768542619139",
    "instrument_id": "925299289",
    "category": "US_STOCK",
    "record_date": "2026-01-14",
    "ex_date": "2026-01-13",
    "payment_date": "2026-01-15",
    "final_pay_date": "2026-01-15",
    "country_code": "US",
    "listing_country_of_code": "US",
    "issuer_country_code": "CA",
    "from": {
      "symbol": "AFRRF",
      "name": "AFR NUVENTURE RESOURCES INC",
      "exchange": "PSGM"
    },
    "to": [
      {
        "option_number": "2",
        "description": "Cash",
        "default_option_flag": "true",
        "payouts": [
          {
            "type": "DV",
            "pay_type": "CASH",
            "payout_number": 1,
            "adr_fee_rate": "0.003",
            "fraction_share_rule": "NONE",
            "cancellation_fee": "0",
            "issuance_fee": "0",
            "currency": "USD",
            "amount": "0.176885",
            "withholding_tax_rate": "0.011"
          }
        ]
      }
    ],
    "biz_type": "CA_EVENT"
  }
}
```

##### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique event identifier |
| `position` | string |  Cursor for event replay; re-pushes subsequent events within the current business type.|
| `event_type` | string | Event type, fixed as `INSTRUMENT` |
| `timestamp` | string | Event timestamp in ISO 8601 format |
| `payload` | object | Event payload data |

##### Payload Fields

| Field | Type | Description |
|-------|------|-------------|
| `eventId` | string | Unique corporate action event identifier |
| `event_type` | string | Corporate action type (e.g., `IDENTIFIER_CHANGE`, `DIVIDEND`, `SPLIT`) |
| `event_version` | string | Event version number, increments with updates. Higher numbers represent newer versions |
| `instrument_id` | string | Instrument identifier affected by the corporate action |
| `category` | string | Instrument category (e.g., `US_STOCK`) |
| `record_date` | string | Record date for the corporate action |
| `ex_date` | string | Ex-dividend or ex-date |
| `payment_date` | string | Payment date |
| `final_pay_date` | string | Final payment date |
| `country_code` | string | Country code |
| `listing_country_of_code` | string | Listing country code |
| `issuer_country_code` | string | Issuer country code |
| `from` | object | Original instrument information |
| `to` | array | Target payment options (multiple options may exist) |
| `biz_type` | string | Business type, fixed as `CA_EVENT` |

##### From Object Fields

| Field | Type | Description |
|-------|------|-------------|
| `symbol` | string | Original instrument symbol |
| `name` | string | Original instrument name |
| `exchange` | string | Original exchange code |

##### To Array - Option Fields

Each option represents a payment choice available to shareholders:

| Field | Type | Description |
|-------|------|-------------|
| `option_number` | string | Option number identifier |
| `description` | string | Option description |
| `default_option_flag` | string | Whether this is the default option: `1` (yes) / `0` (no) |
| `payouts` | array | Array of payout details for this option |

- Multiple payment options may be available, with one marked as the default option
- Currently, only the default option is provided in the event

##### Payout Fields

Each payout represents a specific payment type within an option. A single option may contain multiple payouts (e.g., cash + security).

###### Common Fields (All Payment Types)

These fields are available for all payment types (`CASH`, `SECURITY`, `SCRIP`):

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Payment category code (e.g., `DV`, `FR`, `IN`) |
| `pay_type` | string | Payment type: `CASH`, `SECURITY`, or `SCRIP` |
| `payout_number` | string | Payout number within the option |
| `adr_fee_rate` | string | ADR fee rate |
| `fraction_share_rule` | string | Fractional share handling: `ROUND_DOWN`, `ROUND_UP`, `CASH_IN_LIEU` |
| `cancellation_fee` | string | Cancellation fee rate |
| `issuance_fee` | string | Issuance fee rate |
| `tax_status` | string | IRS income type code |

###### Cash Payment Specific Fields (`pay_type`: `CASH`)

Additional fields for cash payments:

| Field | Type | Description |
|-------|------|-------------|
| `currency` | string | Payment currency |
| `amount` | string | Payment amount per share |
| `withholding_tax_rate` | string | Withholding tax rate |

###### Security/SCRIP Common Fields (`pay_type`: `SECURITY` or `SCRIP`)

These fields are shared by both security and SCRIP payments:

| Field | Type | Description |
|-------|------|-------------|
| `symbol` | string | Payout instrument symbol |
| `name` | string | Payout instrument name |
| `exchange` | string | Payout instrument exchange |
| `from_ratio` | string | Original share ratio |
| `to_ratio` | string | Target share ratio |
| `cash_in_lieu_price` | string | Cash in lieu price for fractional shares |
| `currency` | string | Payment currency |

###### Security Payment Specific Fields (`pay_type`: `SECURITY`)

No additional security-specific fields beyond the common Security/SCRIP fields.

###### SCRIP Payment Specific Fields (`pay_type`: `SCRIP`)

SCRIP dividends (dividend reinvestment) are primarily available for Hong Kong stocks.

| Field | Type | Description |
|-------|------|-------------|
| `reinvest_price` | string | Reinvestment price per share |
| `amount` | string | Payment amount per share |

##### Payout Type Examples

###### Example 1: Cash Dividend

```json
{
    "type": "DV",
    "pay_type": "CASH",
    "payout_number": "1",
    "currency": "USD",
    "amount": "0.25",
    "withholding_tax_rate": "0.15"
}
```

###### Example 2: Stock Dividend (10 shares → 1 share)

```json
{
    "type": "DV",
    "pay_type": "SECURITY",
    "payout_number": "1",
    "symbol": "AAPL",
    "from_ratio": "10",
    "to_ratio": "1",
    "fraction_share_rule": "CASH_IN_LIEU",
    "cash_in_lieu_price": "150.00"
}
```

###### Example 3: SCRIP Dividend (Reinvestment)

```json
{
    "type": "DV",
    "pay_type": "SCRIP",
    "payout_number": "1",
    "symbol": "0700.HK",
    "reinvest_price": "350.00",
    "from_ratio": "1",
    "to_ratio": "1"
}
```

###### Example 4: Mixed Payout (Cash + Security)

An option may contain multiple payouts:

```json
{
    "option_number": "1",
    "description": "Cash and Stock",
    "default_option_flag": "1",
    "payouts": [
        {
            "type": "DV",
            "pay_type": "CASH",
            "payout_number": "1",
            "currency": "USD",
            "amount": "0.10"
        },
        {
            "type": "DV",
            "pay_type": "SECURITY",
            "payout_number": "2",
            "symbol": "AAPL",
            "from_ratio": "20",
            "to_ratio": "1"
        }
    ]
}
```

##### Use Cases

Listening to this event can be used for:

1. **Corporate Action Tracking**: Monitor all corporate actions affecting held positions
2. **Dividend Processing**: Automatically process dividend payments and reinvestments
3. **Position Adjustment**: Update positions based on stock splits, mergers, or conversions
4. **Client Notification**: Alert clients about upcoming corporate actions requiring action
5. **Tax Reporting**: Track withholding taxes and taxable events for reporting purposes
6. **Symbol Changes**: Update instrument identifiers when symbols change

## Trade Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-trade-events.md>

### Trade Events

To enable third-party systems to promptly obtain order execution results and order status updates, Broker OpenAPI provides an asynchronous active push mechanism for order trading events.

Clients can subscribe to order-related events to receive notifications, allowing them to monitor order processing results and execution status changes in real time.

#####  Order
Order Event Notification

<Tabs groupId="programming-language">

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",    
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {  
   "request_id": "1045474398137483264",
    "account_id": "4MHSOMIJ88O7E80VBG0O4G6E9A",
    "client_order_id": "db74f19918054a7e9bb72067731c9ae4",
    "instrument_id": "913256135",
    "order_status": "PARTIAL_FILLED",
    "symbol": "AAPL",
    "qty": "10.00",
    "filled_price": "180.00",
    "filled_qty": "1.00",
    "filled_time": "2025-11-21T06:27:43.312+0000",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "FILLED",
    "biz_type":"TRADE"
    }
}
```

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",    
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {  
   "request_id": "1045474398137483264",
    "account_id": "4MHSOMIJ88O7E80VBG0O4G6E9A",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "db74f19918054a7e9bb72067731c9ae4",
    "instrument_id": "913256135",
    "order_status": "FILLED",
    "symbol": "AAPL",
    "qty": "10.00",
    "filled_price": "180.00",
    "filled_qty": "10.00",
    "filled_time": "2025-11-21T06:27:43.312+0000",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "FINAL_FILLED",
    "biz_type":"TRADE"
    }
}
```

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",    
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {
    "request_id": "1045474643156140032",
    "account_id": "4MHSOMIJ88O7E80VBG0O4G6E9A",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "de2868b71c154bcaafd2baca61127966",
    "instrument_id": "913256135",
    "order_status": "FAILED",
    "symbol": "AAPL",
    "qty": "10.00",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "PLACE_FAILED",
    "biz_type":"TRADE"
    }
}
```

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",      
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV5P4I8BV0KHKN60000000",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "04cda8db7ed940f6afeb26be6201ee53",
    "instrument_id": "913256135",
    "order_status": "SUBMITTED",
    "symbol": "AAPL",
    "qty": "4.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "MODIFY_SUCCESS",
    "biz_type":"TRADE"
}
}
```

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",    
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV5P4I8BV0KHKN60000000",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "04cda8db7ed940f6afeb26be6201ee53",
    "instrument_id": "913256135",
    "order_status": "SUBMITTED",
    "symbol": "AAPL",
    "qty": "4.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "MODIFY_FAILED",
    "biz_type":"TRADE"
}
}
```

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",      
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV5P4I8BV0KHKN60000000",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "04cda8db7ed940f6afeb26be6201ee53",
    "instrument_id": "913256135",
    "order_status": "CANCELLED",
    "symbol": "AAPL",
    "qty": "4.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "CANCEL_SUCCESS",
    "biz_type":"TRADE"
}
}
```

```python
{
"id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",    
"event_type": "TRADE",      
"timestamp": "2025-03-29T07:02:33.200962333Z",    
"payload": {
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV5P4I8BV0KHKN60000000",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "04cda8db7ed940f6afeb26be6201ee53",
    "instrument_id": "913256135",
    "order_status": "SUBMITTED",
    "symbol": "AAPL",
    "qty": "4.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "scene_type": "CANCEL_FAILED",
    "biz_type":"TRADE"
}
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `TRADE`       |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field           | Type   | Description                                                                                                                 |
|-----------------|--------|-----------------------------------------------------------------------------------------------------------------------------|
| account_id      | string | Account id                                                                                                                  |
| request_id      | string | Request Id                                                                                                                  |
| order_id        | string | System-generated order identifier.                                                                                          |
| client_order_id | array  | Client-defined order identifier.                                                                                            |
| instrument_id   | string | Instrument Id                                                                                                               |
| order_status    | string | Order Status, See the `status` field in the Order Detail API response.    |
| symbol          | string | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market.        |
| qty             | string | Total order quantity. Represents the total number of units submitted for this order.                                        |
| filled_qty      | string | Quantity that has been executed. Represents the number of units that have been filled so far.                               |
| filled_price    | string | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null.         |
| filled_time     | string | Time of the last executed trade. Example: 2025-11-21T06:27:43.312+0000                                                         |
| side            | string | Order Side, See the `side` field in the Order Detail API response.        |
| category        | string | Category, `HK_STOCK`  or `US_STOCK` .                                                                                       |
| order_type      | string | Order Type, See the `order_type`  field in the Order Detail API response. |
| scene_type      | string | Indicates the order event scenario or execution result.                                                                     |
| biz_type        | string | Business type, fixed as `TRADE`                                                                                             |

###### Scene Types

| scene_type     | Description               |
|----------------|---------------------------|
| FILLED         | Partially filled          |
| FINAL_FILLED   | All filled                |
| PLACE_FAILED   | Order failed              |
| MODIFY_SUCCESS | Change order successfully |
| MODIFY_FAILED  | Change order failed       |
| CANCEL_SUCCESS | Cancellation succeeded    |
| CANCEL_FAILED  | Cancellation failed       |

## Funding Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-funding-events.md>

### Funding Events

To enable third-party systems to promptly obtain the execution results of foreign exchange and deposit/withdrawal operations, Broker OpenAPI provides an asynchronous proactive event push mechanism for FX conversion and deposit/withdrawal transaction events.

Clients can subscribe to receive relevant events, allowing them to monitor transaction processing results and status changes in real time.

#####  Instant Funding
Instant Funding Event Notification

<Tabs groupId="programming-language">

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_funding_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "amount": "100",
        "currency": "HKD",
        "type": "DEPOSIT",
        "status": "CANCELED",
        "reason": "xxxxxx",
        "biz_type": "INSTANT_FUNDING"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_funding_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "amount": "100",
        "currency": "HKD",
        "type": "DEPOSIT",
        "status": "REJECTED",
        "reason": "xxxxxx",
        "biz_type": "INSTANT_FUNDING"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_funding_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "amount": "100",
        "currency": "HKD",
        "type": "DEPOSIT",
        "status": "FAILED",
        "reason": "xxxxxx",
        "biz_type": "INSTANT_FUNDING"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_funding_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "amount": "100",
        "currency": "HKD",
        "type": "DEPOSIT",
        "status": "COMPLETED",
        "biz_type": "INSTANT_FUNDING"
    }
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `FUNDING`     |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field              | Type   | Description                                                                                                                                                   |
|--------------------|--------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| account_id         | string | Account id                                                                                                                                                    |
| client_request_id  | string | Client Request ID, unique for each request.                                                                                                                   |
| instant_funding_id | string | Request id generated by the system.                                                                                                                           |
| amount             | string | Amount.                                                                                                                                                       |
| currency           | string | Currency, See the `currency` field in the Get Instant Funding Detail API response.                 |
| type               | string | Instant Funding Request Type, See the `type` field in the Get Instant Funding Detail API response. |
| status             | string | Request status, See the `status` field in the Get Instant Funding Detail API response.             |
| reason             | string | Reason. If the terminal state is not “COMPLETED”, return the reason.                                                                                          |
| biz_type           | string | Business type, fixed as `INSTANT_FUNDING`                                                                                                                     |

#####  Instant Exchange
Instant Exchange Event Notification

<Tabs groupId="programming-language">

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_exchange_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "to_amount": "700",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "CANCELED",
        "reason": "xxxxxx",
        "biz_type": "INSTANT_EXCHANGE"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_exchange_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "to_amount": "700",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "REJECTED",
        "reason": "xxxxxx",
        "biz_type": "INSTANT_EXCHANGE"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_exchange_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "to_amount": "700",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "FAILED",
        "reason": "xxxxxx",
        "biz_type": "INSTANT_EXCHANGE"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "instant_exchange_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "to_amount": "700",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "COMPLETED",
        "biz_type": "INSTANT_EXCHANGE"
    }
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `FUNDING`     |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field             | Type   | Description                                                                                                                                                               |
|-------------------|--------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| account_id        | string | Account id                                                                                                                                                                |
| client_request_id | string | Client Request ID, unique for each request.                                                                                                                               |
| instant_exchange_id     | string | Request id generated by the system.                                                                                                                                       |
| from_amount       | string | Base amount for exchange.                                                                                                                                                 |
| to_amount         | string | Amount after exchange.                                                                                                                                                    |
| from_currency     | string | Base currency for exchange., See the `from_currency` field in the Get Instant Exchange Detail API response. |
| to_currency       | string | Currency after exchange. , See the `to_currency` field in the Get Instant Exchange Detail API response.     |
| status            | string | Request status, See the `status` field in the Get Instant Exchange Detail API response.                     |
| reason            | string | Reason. If the terminal state is not “COMPLETED”, return the reason.                                                                                                      |
| biz_type          | string | Business type, fixed as `INSTANT_EXCHANGE`                                                                                                                                |

#####  Exchange
Foreign Exchange Conversion Event Notification

<Tabs groupId="programming-language">

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "fx_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "CANCELED",
        "reason": "xxxxxx",
        "biz_type": "EXCHANGE"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "fx_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "REJECTED",
        "reason": "xxxxxx",
        "biz_type": "EXCHANGE"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "fx_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "from_currency": "USD",
        "to_currency": "HKD",
        "status": "FAILED",
        "reason": "xxxxxx",
        "biz_type": "EXCHANGE"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "FUNDING",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "account_id": "93IUJ28O9VO2KBGHDHR4H9",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "fx_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "from_amount": "100",
        "to_amount": "801",
        "from_currency": "USD",
        "to_currency": "HKD",
        "fx_rate": "8.01",
        "status": "COMPLETED",
        "biz_type": "EXCHANGE"
    }
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `FUNDING`     |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field             | Type   | Description                                                                                                                                         |
|-------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| account_id        | string | Account id                                                                                                                                          |
| client_request_id | string | Client Request ID, unique for each request.                                                                                                         |
| fx_id             | string | Request id generated by the system.                                                                                                                 |
| from_amount       | string | Base amount for exchange.                                                                                                                           |
| to_amount         | string | Amount after exchange.                                                                                                                              |
| from_currency     | string | Base currency for exchange., See the `from_currency` field in the Get FX Detail API response. |
| to_currency       | string | Currency after exchange. , See the `to_currency` field in the Get FX Detail API response.     |
| fx_rate           | string | FX ratre.                                                                                                                                           |
| status            | string | Request status, See the `status` field in the Get FX Detail API response.                     |
| reason            | string | Reason. If the terminal state is not “COMPLETED”, return the reason.                                                                                |
| biz_type          | string | Business type, fixed as `EXCHANGE`                                                                                                                  |

## Journal Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-journal-events.md>

### Journal Events
To enable third-party systems to promptly obtain the results of journal transfers (cash / position), Broker OpenAPI supports an asynchronous event push mechanism for journal transfer–related operations (cash / position).

When the status of a journal transfer operation (cash / position) changes, the system will deliver event notifications, and subscribed clients will receive the final or intermediate execution results.

#####  Cash Journal
Cash Journal Event Notification

<Tabs groupId="programming-language">

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "CASH",
        "amount": "100",
        "currency": "USD",
        "status": "CANCELED",
        "reason": "xxxxxx",
        "biz_type": "CASH_JOURNAL"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "CASH",
        "amount": "100",
        "currency": "USD",
        "status": "REJECTED",
        "reason": "xxxxxx",
        "biz_type": "CASH_JOURNAL"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "CASH",
        "amount": "100",
        "currency": "USD",
        "status": "FAILED",
        "reason": "xxxxxx",
        "biz_type": "CASH_JOURNAL"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "CASH",
        "amount": "100",
        "currency": "USD",
        "status": "COMPLETED",
        "biz_type": "CASH_JOURNAL"
    }
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `JOURNAL`     |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field             | Type   | Description                                                                                                                                   |
|-------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| from_account      | string | From account id.                                                                                                                              |
| to_account        | string | To account id.                                                                                                                                |
| client_request_id | string | Client Request ID, unique for each request.                                                                                                   |
| journal_id        | string | Request id generated by the system.                                                                                                           |
| journal_type      | string | Journal type, fixed as `CASH`                                                                                                                 |
| amount            | string | Amount                                                                                                                                        |
| currency          | string | Currency, See the `currency` field in the Query Cash Journal Detail API response.     |
| status            | string | Request status, See the `status` field in the Query Cash Journal Detail API response. |
| reason            | string | Reason. If the terminal state is not “COMPLETED”, return the reason.                                                                          |
| biz_type          | string | Business type, fixed as `CASH_JOURNAL`                                                                                                        |

#####  Position Journal
Position Journal Event Notification

<Tabs groupId="programming-language">

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "POSITION",
        "instrument_type": "EQUITY",
        "market": "US",
        "symbol": "BULL",
        "quantity": "100",
        "status": "CANCELED",
        "reason": "xxxxxx",
        "biz_type": "POSITION_JOURNAL"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "POSITION",
        "instrument_type": "EQUITY",
        "market": "US",
        "symbol": "BULL",
        "quantity": "100",
        "status": "REJECTED",
        "reason": "xxxxxx",
        "biz_type": "POSITION_JOURNAL"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "POSITION",
        "instrument_type": "EQUITY",
        "market": "US",
        "symbol": "BULL",
        "quantity": "100",
        "status": "FAILED",
        "reason": "xxxxxx",
        "biz_type": "POSITION_JOURNAL"
    }
}
```

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "JOURNAL",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "from_account": "93IUJ28O9VO2KBGHDHR4H9",
        "to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
        "client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
        "journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
        "journal_type": "POSITION",
        "instrument_type": "EQUITY",
        "market": "US",
        "symbol": "BULL",
        "quantity": "100",
        "status": "COMPLETED",
        "biz_type": "POSITION_JOURNAL"
    }
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `JOURNAL`     |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field             | Type   | Description                                                                                                                                                                                                                      |
|-------------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| from_account      | string | From account id.                                                                                                                                                                                                                 |
| to_account        | string | To account id.                                                                                                                                                                                                                   |
| client_request_id | string | Client Request ID, unique for each request.                                                                                                                                                                                      |
| journal_id        | string | Request id generated by the system.                                                                                                                                                                                              |
| journal_type      | string | Journal type, fixed as `POSITION`                                                                                                                                                                                                |
| market            | string | Market code indicating the trading venue or regulatory region of the financial instrument. See the `market` field in the Query Position Journal Detail API response. |
| instrument_type   | string | Type of financial instrument associated with the request. See the `instrument_type` field in the Query Position Journal Detail API response.                         |
| symbol            | string | Symbol name.                                                                                                                                                                                                                     |
| quantity          | string | Quantity                                                                                                                                                                                                                         |
| status            | string | Request status, See the `status` field in the Query Position Journal Detail API response.                                                                            |
| reason            | string | Reason. If the terminal state is not “COMPLETED”, return the reason.                                                                                                                                                             |
| biz_type          | string | Business type, fixed as `POSITION_JOURNAL`                                                                                                                                                                                       |

## Master Data Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/broker-master-data-events.md>

### Master Data Events

#####  Trade Calendar
To enable third-party systems to promptly obtain updates on changes to the trading and settlement calendar, Broker OpenAPI provides an asynchronous proactive event push mechanism for trading and settlement calendar events.

Clients can subscribe to receive these events, allowing them to monitor adjustments to trading days and settlement dates in real time, ensuring accurate time calculations in trading and fund management systems.

<Tabs groupId="programming-language">

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "MASTER_DATA",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "request_id": "036LVV5P4I8BV0KHKN60000000",
        "market": "HK",
        "instrument_type": "EQUITY",
        "year": 2026,
        "biz_type": "CALENDAR_UPDATE"
    }
}
```

##### Response Fields

| Field      | Type   | Description                        |
|------------|--------|------------------------------------|
| id         | string | Unique event identifier            |
| event_type | string | Event type, fixed as `MASTER_DATA` |
| timestamp  | string | Event timestamp in ISO 8601 format |
| payload    | object | Event payload data                 |

##### Payload Fields

| Field           | Type   | Description                                                                                                                              |
|-----------------|--------|------------------------------------------------------------------------------------------------------------------------------------------|
| request_id      | string | Request Id.                                                                                                                              |
| market          | string | See the `market` field in the Query Trade Calendar API request for details.          |
| instrument_type | string | See the `instrument_type` field in the Query Trade Calendar API request for details. |
| year            | string | Year in YYYY format.                                                                                                                     |
| biz_type        | string | Business type, fixed as `CALENDAR_UPDATE`                                                                                                |


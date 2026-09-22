# Connect API (OAuth) — Verbatim Reference

> OAuth 2.0 authorization-code flow for third-party applications (US only).

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Authorization Code

> Source: <https://developer.webull.com/apis/docs/reference/connect-api/get-authorization-code.md>

### Get Authorization Code

This is the first step of the OAuth2 process. An authorization code is created when the user authorizes your application to access their account. If the user grants permission to your application, the callback URL registered in your application will be invoked. The interface for obtaining the authorization code is completed in the browser.<br/> <b>'SEND API REQUEST' function for this endpoint does not work in UAT environment</b>.

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
      "url": "https://oauth-open-api.sandbox.webull.com"
    }
  ],
  "path": "/oauth2/auth-codes/get",
  "method": "get",
  "tags": [
    "Connect"
  ],
  "description": "This is the first step of the OAuth2 process. An authorization code is created when the user authorizes your application to access their account. If the user grants permission to your application, the callback URL registered in your application will be invoked. The interface for obtaining the authorization code is completed in the browser.<br/> <b>'SEND API REQUEST' function for this endpoint does not work in UAT environment</b>.",
  "operationId": "getAuthorizationCode",
  "parameters": [
    {
      "name": "response_type",
      "in": "query",
      "description": "Must be code to request an authorization code.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "code"
    },
    {
      "name": "client_id",
      "in": "query",
      "description": "Webull provides the client id",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "CLINETTEST"
    },
    {
      "name": "scope",
      "in": "query",
      "description": "The application requests access to the list of scopes. user：user trade：trade wr：write read.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "user:trade:wr"
    },
    {
      "name": "state",
      "in": "query",
      "description": "An unguessable random string, used to protect against request forgery attacks.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "MiLCJjb25uZWN0aW9uX3R5cGUiOiJs"
    },
    {
      "name": "redirect_uri",
      "in": "query",
      "description": "The URL to which the user will be redirected after authorization. It must match the redirect URIs in whitelist.",
      "required": true,
      "schema": {
        "type": "String"
      },
      "example": "http://callbackurl.com"
    }
  ],
  "responses": {
    "302": {
      "description": "After successful authorization, it will call back to the redirect_uri in the request parameters, structured as follows:<br/> http://testcallbackurl.com?code=NjVhODIxODItYTAzMC00Y2IxLTkzNzQt&state=MiLCJjb25uZWN0aW9uX3R5cGUiOiJs"
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
    "name": "Get Authorization Code",
    "description": {
      "content": "This is the first step of the OAuth2 process. An authorization code is created when the user authorizes your application to access their account. If the user grants permission to your application, the callback URL registered in your application will be invoked. The interface for obtaining the authorization code is completed in the browser.<br/> <b>'SEND API REQUEST' function for this endpoint does not work in UAT environment</b>.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "oauth2",
        "auth-codes",
        "get"
      ],
      "host": [
        "{{baseUrl}}"
      ],
      "query": [
        {
          "disabled": false,
          "description": {
            "content": "(Required) Must be code to request an authorization code.",
            "type": "text/plain"
          },
          "key": "response_type",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) Webull provides the client id",
            "type": "text/plain"
          },
          "key": "client_id",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) The application requests access to the list of scopes. user：user trade：trade wr：write read.",
            "type": "text/plain"
          },
          "key": "scope",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) An unguessable random string, used to protect against request forgery attacks.",
            "type": "text/plain"
          },
          "key": "state",
          "value": ""
        },
        {
          "disabled": false,
          "description": {
            "content": "(Required) The URL to which the user will be redirected after authorization. It must match the redirect URIs in whitelist.",
            "type": "text/plain"
          },
          "key": "redirect_uri",
          "value": ""
        }
      ],
      "variable": []
    },
    "header": [
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "GET"
  }
}
```

## Create Token

> Source: <https://developer.webull.com/apis/docs/reference/connect-api/create-and-refresh-token.md>

### Create And Refresh Token

This is the second step of the OAuth process. An access token is created using the authorization code from the first step's response. The access token is a key used for API access. These tokens should be protected like passwords.

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
      "url": "https://oauth-open-api.sandbox.webull.com"
    }
  ],
  "path": "/oauth2/tokens/create",
  "method": "post",
  "tags": [
    "Connect"
  ],
  "description": "This is the second step of the OAuth process. An access token is created using the authorization code from the first step's response. The access token is a key used for API access. These tokens should be protected like passwords.",
  "operationId": "CreateAndRefreshToken",
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
      "description": "Interface version, Accepts only the value v2.",
      "required": true,
      "schema": {
        "type": "string",
        "default": "v2"
      },
      "examples": {
        "v2": {
          "value": "v2"
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
      "application/x-www-form-urlencoded": {
        "schema": {
          "required": [
            "client_id",
            "client_secret",
            "grant_type"
          ],
          "type": "object",
          "properties": {
            "client_id": {
              "type": "string",
              "description": "Webull provides the client_id",
              "example": "CLINETTEST"
            },
            "client_secret": {
              "type": "string",
              "description": "The client password provided by Webull.",
              "example": "YzExZTVkMDQtYzU210D0MzAxLWE1M2UtNGQ1YmQwYzU="
            },
            "grant_type": {
              "type": "string",
              "description": "The request to create an access token be set to authorized_code. refresh using refresh_token.",
              "example": "authorization_code",
              "enum": [
                "authorization_code",
                "refresh_token"
              ]
            },
            "code": {
              "type": "string",
              "description": "The authorization code received in Get An Authorization Code, Creating token Required",
              "example": "MDM2T0IyUFNRNDk4UzBLSEtCVDgwMDAwMDA="
            },
            "refresh_token": {
              "type": "string",
              "description": "Required for refresh, use the refresh_token from the return value.",
              "example": "MDM2VTFFUzlHSTk4UzBLSEs2RTgwMDAwMDE="
            }
          },
          "title": "TokenReq"
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
            "required": [
              "access_token",
              "created_at",
              "expires_in",
              "identity_id",
              "refresh_token",
              "rt_expires_in",
              "token_type"
            ],
            "type": "object",
            "properties": {
              "access_token": {
                "type": "string",
                "description": "Access token",
                "example": "MDM2VTFFUzlHSTk4UzBLSEs2RTgwMDAwMDA="
              },
              "token_type": {
                "type": "string",
                "description": "Access token type. Currently, only 'Bearer' is supported",
                "example": "Bearer"
              },
              "expires_in": {
                "type": "string",
                "description": "Access token expiration time. Unit: seconds.",
                "example": "1800"
              },
              "rt_expires_in": {
                "type": "string",
                "description": "Refresh token expiration time. Unit: seconds.",
                "example": "1296000"
              },
              "refresh_token": {
                "type": "string",
                "description": "Refresh token",
                "example": "MDM2VTFFUzlHSTk4UzBLSEs2RTgwMDAwMDE="
              },
              "created_at": {
                "type": "string",
                "description": "Token creation time.",
                "example": "2024-09-01T19:37:22.532+0000"
              },
              "identity_id": {
                "type": "string",
                "description": "User unique identifier.",
                "example": "a0920c03dcba35a56a33a7f4a0df6275"
              }
            },
            "title": "TokenResult"
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
    "name": "Create And Refresh Token",
    "description": {
      "content": "This is the second step of the OAuth process. An access token is created using the authorization code from the first step's response. The access token is a key used for API access. These tokens should be protected like passwords.",
      "type": "text/plain"
    },
    "url": {
      "path": [
        "oauth2",
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
          "content": "(Required) Interface version, Accepts only the value v2.",
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
        "value": "application/x-www-form-urlencoded"
      },
      {
        "key": "Accept",
        "value": "application/json"
      }
    ],
    "method": "POST",
    "body": {
      "mode": "urlencoded",
      "urlencoded": []
    }
  }
}
```


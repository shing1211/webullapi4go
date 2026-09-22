# Broker API — FD (US)

US Broker FD surface. Documented on the US site only.

[<- Webull API Reference](../webull-api.md)

## List Accounts

`GET /broker/accounts/list`

> Retrieves a paginated list of accounts associated with the requesting institution.

| | |
|---|---|
| **SDK** | `brokerfd.ListFDAccounts` |
| **Reference** | [list-accounts.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-accounts.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `pagination_key` | query | string |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Unique identifier of the account to which restrictions are currently applied. |
| `account_number` | string | yes | Account number associated with the account. |
| `account_type` | string | yes | Account Type CASH - Cash Account Type MARGIN - Margin Account Type — one of: `CASH`, `MARGIN` |
| `rep_code` | string |  | Rep Code |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Account Detail

`GET /broker/accounts/get`

> Retrieves the detailed information of a specific account.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDAccountDetail` |
| **Reference** | [get-account-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-account-detail.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | string | yes | Unique identifier of the account to query restrictions for. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Unique identifier of the account to which restrictions are currently applied. |
| `account_number` | string | yes | Account number associated with the account. |
| `account_type` | string | yes | Account Type CASH - Cash Account Type MARGIN - Margin Account Type — one of: `CASH`, `MARGIN` |
| `rep_code` | string |  | Rep Code |
| `commission_code` | string |  | Commission Code |
| `restrict_infos` | array<object> |  | List of restrictions currently effective for this account. |

*Nested — `restrict_infos`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `restriction_code` | string | yes | A code representing the type of restriction applied. RISK: Risk Restriction TRADING: Trading Restriction TRADING_STOCK_ETF: Stock and ETF Trading Restriction TRADING_FRACTIONAL: Fractional Share Trading Restriction TRADING_OPTIONS: Options Trading Restriction TRADING_EVENT: Event Contract Trading Restriction TRADING_OTC: OTC Trading Restriction TRADING_OVERNIGHT: Overnight Trading Restriction ACCOUNT_LINKING: Account Linking Restriction FUNDING_ACH: ACH Deposit Restriction FUNDING_DEBIT_CARD: Debit Card Deposit Restriction FUNDING_DIRECT_DEPOSIT: Direct Deposit Restriction WITHDRAWAL: Withdrawal Restriction TRANSFER_ACATS_IN: ACATS Incoming Transfer Restriction TRANSFER_ACATS_OUT: ACATS Outgoing Transfer Restriction |
| `restriction_start_time` | string | yes | The date and time when the restriction takes effect. Time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z' |
| `reason` | string |  | A description of the reason for the restriction. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Account

`POST /broker/accounts/create`

> Submits a request to create an account for your end user. Related Event Notifications: Please refer to the Event Push API documentation: Application Events

| | |
|---|---|
| **SDK** | `brokerfd.CreateFDAccount` |
| **Reference** | [create-account-apply.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-account-apply.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client-supplied request ID, which is unique for each account submission request. |
| `application_id` | string |  | System-generated unique identifier assigned to the account application record. This field is omitted for the initial submission and must be provided for subsequent submissions related to the same account application. |
| `related_account_id` | string |  | The unique identifier of the related account. This field is required when submitting an additional account application related to an existing account. For example, if you have an Event Contract account and want to add a Brokerage Margin account, you would provide the account_id of the Event Contract account in this field. If this field is not provided or is empty, it indicates that this submission is for an initial account application. |
| `forms` | array<object> | yes | Array of one or more FormData objects representing the account form being submitted. Use form details API to retrieve the form schema and details. |

*Nested — `forms`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `form_info` | object | yes | Form Identifier. |
| `json_data` | object | yes | Object containing the completed form data as defined by the form schema and will vary from form to form. |

*Nested — `form_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `form_code` | string | yes | Unique identifier of the form, typically used to distinguish different form definitions. |
| `version` | string | yes | Version number of the form definition, used for compatibility and evolutionary control. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Unique client-generated identifier used for request tracking, idempotency and audit purposes. |
| `application_id` | string | yes | System-generated unique identifier assigned to the account application record. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Update Account

`POST /broker/accounts/update`

> Submit a request to update an existing account. This API only updates the fields submitted in the forms data. Any attributes not included in the request will remain unchanged. Related Event Notifications: Please refer to the Event Push API documentation: Account Events

| | |
|---|---|
| **SDK** | `brokerfd.UpdateFDAccount` |
| **Reference** | [update-account-apply.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/update-account-apply.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Unique identifier for the request, generated each time. |
| `account_id` | string | yes | The account ID for which the update is being submitted. |
| `forms` | array<object> | yes | Array of one or more FormData objects representing the account form being submitted. Use form details API to retrieve the form schema and details. |

*Nested — `forms`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `form_info` | object | yes | Form Identifier. |
| `json_data` | object | yes | Object containing the completed form data as defined by the form schema and will vary from form to form. |

*Nested — `form_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `form_code` | string | yes | Unique identifier of the form, typically used to distinguish different form definitions. |
| `version` | string | yes | Version number of the form definition, used for compatibility and evolutionary control. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Unique client-generated identifier used for request tracking, idempotency and audit purposes. |
| `account_id` | string | yes | System-generated unique identifier assigned to the account application record. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Close Account

`POST /broker/accounts/close`

> Closes an existing account. The operation is irreversible and requires valid account information. Related Event Notifications: Please refer to the Event Push API documentation: Account Events

| | |
|---|---|
| **SDK** | `brokerfd.CloseFDAccount` |
| **Reference** | [close-account.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/close-account.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client-supplied request ID, which is unique for each account closure request. |
| `account_id` | string | yes | Unique identifier of the account to be closed. |
| `close_reason` | string | yes | Optional detailed textual explanation of the account closure reason. Limited to 512 characters. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client-generated unique identifier for idempotency and audit purposes. |
| `account_id` | string | yes | Account ID to be closed. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Account Application Detail

`GET /broker/accounts/applications/get`

> Retrieves detailed information about a specific account application.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [get-account-application-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-account-application-detail.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `application_id` | query | string | yes | Unique identifier of the account application to retrieve details for. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `application_id` | string | yes | Application ID |
| `resubmission_codes` | array<string> |  | Resubmission Code, only present for resubmission applications |
| `customer_info` | object | yes | Customer Info |
| `contact_info` | object | yes | Contact Information |
| `trusted_contact_info` | object |  | Trusted Contact Info |
| `employment_info` | object | yes | Employment Info |
| `financial_info` | object | yes | Financial Info |
| `disclosure_info` | object | yes | Disclosure Info |
| `tax_info` | object | yes | Tax Info |
| `w8ben_info` | object |  | W8ben Info |
| `additional_infos` | array<object> | yes | Additional Information for the application review process. |
| `document_attachments` | array<object> |  | Document Attachments |

*Nested — `customer_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `ip_address` | string | yes | The IP address when the user submits the application. |
| `first_name` | string | yes | First name |
| `middle_name` | string |  | Middle name |
| `last_name` | string | yes | Last name |
| `date_of_birth` | string | yes | Date of birth |
| `gender` | string | yes | Gender. Allowed values: F (Female), M (Male). |
| `country_of_citizenship` | string | yes | Country of citizenship. Value must be obtained from the Master Data API with data_type = COUNTRY |
| `id_doc_type` | string | yes | Type of identification document. Value must be obtained from the Master Data API with data_type = ID_DOC_TYPE. |
| `picture_front_doc_id` | string |  | Picture front document id |
| `picture_back_doc_id` | string |  | Picture back document id |
| `id_number` | string | yes | Number of the identification document. |
| `id_expiry_date` | string | yes | Expiration date of the identification document. |
| `selfie_img_key` | string |  | User’s selfie, face photo, or liveness document id |
| `tax_id_doc_id` | string |  | Tax ID or SSN document id for tax/SSN verification |
| `permanent_resident` | string | yes | Permanent resident. Allowed values: YES, NO. |
| `marital_status` | string | yes | The user's marital status. Value must be obtained from the Master Data API with data_type = MARITAL_STATUS. |
| `num_dependents` | integer | yes | Num dependents |
| `ext_attr_list` | array<object> |  | Ext attr list |

*Nested — `ext_attr_list`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Attribute Name |
| `value` | string | yes | Attribute Value |

*Nested — `contact_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `phone_country` | string | yes | Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `phone_number` | string | yes | Phone number |
| `email_address` | string | yes | Email address |
| `id_address` | object | yes | Address Info |
| `mail_address` | object |  | Address Info |

*Nested — `id_address`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `country` | string | yes | Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state` | string | yes | State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA). |
| `city` | string | yes | City. |
| `street_address` | string | yes | Street address |
| `postal_code` | string | yes | Postal code |
| `neighborhood_code` | string |  | Neighborhood code. |
| `apartment_number` | string |  | Apartment number. |
| `building_number` | string |  | Building number. |

*Nested — `mail_address`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `country` | string | yes | Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state` | string | yes | State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA). |
| `city` | string | yes | City. |
| `street_address` | string | yes | Street address |
| `postal_code` | string | yes | Postal code |
| `neighborhood_code` | string |  | Neighborhood code. |
| `apartment_number` | string |  | Apartment number. |
| `building_number` | string |  | Building number. |

*Nested — `trusted_contact_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `first_name` | string | yes | First name of the trusted contact. |
| `middle_name` | string |  | Middle name of the trusted contact. |
| `last_name` | string | yes | Last name of the trusted contact. |
| `phone_country` | string | yes | Phone country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `phone_number` | string | yes | Phone number |
| `relationship` | string | yes | Relationship between the account holder and the trusted contact person. |
| `trusted_address` | object | yes | Trusted Contact Address |

*Nested — `trusted_address`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `country` | string | yes | Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state` | string | yes | State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA). |
| `city` | string | yes | City. |
| `street_address` | string | yes | Street address |
| `postal_code` | string | yes | Postal code |
| `neighborhood_code` | string |  | Neighborhood code. |
| `apartment_number` | string |  | Apartment number. |
| `building_number` | string |  | Building number. |

*Nested — `employment_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `employment_type` | string | yes | Type of employment. Value must be obtained from the Master Data API with data_type = EMPLOYMENT_TYPE. |
| `company_name` | string | yes | Name of the applicant current employer. |
| `occupation` | string | yes | Specific occupation under the employment type. Value must be obtained from the Master Data API with data_type = OCCUPATION. |
| `position` | string | yes | Position. Value must be obtained from the Master Data API with data_type = POSITION |
| `years_employed` | string | yes | Total duration (in years) the applicant has been employed with the current employer. Decimal values are allowed to represent partial years (e.g. 1.5). |
| `employment_address` | object | yes | Employment Address |

*Nested — `employment_address`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `country` | string | yes | Country. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state` | string | yes | State. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: ISO 3166-1 alpha-3 country code followed by a hyphen and the subdivision code (for example: USA-CA). |
| `city` | string | yes | City. |
| `street_address` | string | yes | Street address |
| `postal_code` | string | yes | Postal code |
| `neighborhood_code` | string |  | Neighborhood code. |
| `apartment_number` | string |  | Apartment number. |
| `building_number` | string |  | Building number. |

*Nested — `financial_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `liquidity_needs` | string | yes | Indicates the client liquidity needs or preference for cash accessibility. Value must be obtained from the Master Data API with data_type = LIQUIDITY_NEEDS. |
| `annual_income` | string | yes | Annual income. Value must be obtained from the Master Data API with data_type = ANNUAL_INCOME. |
| `total_net_worth` | string | yes | Client's total net worth (USD). Value must be obtained from the Master Data API with data_type = TOTAL_NET_WORTH. |
| `liquid_net_worth` | string | yes | Liquidity assets. Value must be obtained from the Master Data API with data_type = LIQUID_NET_WORTH. |

*Nested — `disclosure_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `is_politically_exposed` | boolean | yes | Indicates whether the account holder or any close associate or immediate family member is a Politically Exposed Person (PEP). If true, the politically_exposed_info object must be provided with detailed disclosure information. Please note: PEP accounts are currently not supported. Any request submitted with is_politically_exposed = true will be rejected and a failure response will be returned. |
| `politically_exposed_info` | object |  | Politically Exposed Info |
| `is_cftc_or_nfa_registered` | boolean | yes | Indicates whether the applicant is registered with the CFTC or is a member of the NFA. |
| `is_company_control_person` | boolean | yes | Indicates whether the applicant is a control person of a public company. If true, `public_company_infos` must be provided with details of the companies controlled. |
| `public_company_infos` | array<object> |  | List of public companies controlled by the applicant. Each item must include the company stock symbol. |

*Nested — `politically_exposed_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `pep_organization` | string |  | Name of the related political organization. |
| `immediate_family` | array<string> |  | Name of family members, including former spouses. |

*Nested — `public_company_infos`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Stock symbol of the public company controlled by the applicant, as listed on the exchange. |

*Nested — `tax_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `tax_type` | string | yes | Tax Type. Value must be obtained from the Master Data API with data_type = TAX_TYPE. |
| `tax_id` | string | yes | Tax ID |

*Nested — `w8ben_info`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `foreign_tax_id` | string | yes | Foreign tax ID (FTIN) issued by non-U.S. authority (W-8BEN Line 6) |
| `treaty_country` | string | yes | Tax treaty country for reduced withholding. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `sign_date` | string | yes | Date when the W-8BEN form was signed by the customer. The date should be in ISO 8601 format (YYYY-MM-DD). |

*Nested — `additional_infos`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `status` | string | yes | Application Status — one of: `SUBMITTED`, `PENDING_REVIEW`, `ACTION_REQUIRED`, `REJECTED`, `APPROVED` |
| `form_code` | string | yes | Form code |
| `rep_code` | string | yes | Rep code |
| `account_id` | string |  | Unique account identifier. Populated only when the status is APPROVED; otherwise, this field is empty. |
| `account_number` | string |  | Account number. Populated only when the status is APPROVED; otherwise, this field is empty. |
| `is_exchange_or_finra_affiliated` | boolean | yes | Indicates whether the applicant is employed by, or directly associated with, an exchange or a FINRA member firm. If true, exchange_affiliations details must be provided. |
| `exchange_or_finra_affiliations` | array<object> |  | List of the applicant affiliations with exchanges or FINRA member firms. Required if is_exchange_or_finra_affiliated is true. |
| `investment_experience` | string |  | Investment experience. Value must be obtained from the Master Data API with data_type = INVESTMENT_EXPERIENCE. |
| `investment_knowledge` | string |  | Investment knowledge. Value must be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE. Note: For Event Contract , the value of this field should be obtained from the Master Data API with data_type = INVESTMENT_KNOWLEDGE_V2 instead. |
| `investment_objective` | string |  | Investment objective. Value must be obtained from the Master Data API with data_type = INVESTMENT_OBJECTIVE. |
| `trade_per_year` | string |  | Trade per year. Value must be obtained from the Master Data API with data_type = TRADE_PER_YEAR. |
| `time_horizon` | string |  | Time horizon. Value must be obtained from the Master Data API with data_type = TIME_HORIZON. |

*Nested — `exchange_or_finra_affiliations`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `company_name` | string | yes | Company name |
| `country` | string | yes | Country of the affiliated exchange or FINRA firm. The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state` | string | yes | State or province of the affiliated exchange or FINRA firm. The value must follow the ISO 3166-2 subdivision standard and include the country code prefix. Format: <ISO 3166-1 alpha-3 country code>-<Subdivision Code> |
| `city` | string | yes | City of the affiliated exchange or FINRA firm. |
| `street_address` | string | yes | Street address |
| `postal_code` | string | yes | Postal code |

*Nested — `document_attachments`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `doc_type` | string | yes | Document Type |
| `document_id` | string | yes | Document ID |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List Forms

`GET /broker/forms/list`

> Retrieves a list of available form codes.

| | |
|---|---|
| **SDK** | `brokerfd.ListAccountForms` |
| **Reference** | [get-form-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-list.md) |

**Response 200**

Array of `string`.

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List Form Versions

`GET /broker/forms/versions/list`

> Retrieves a list of available versions for the specified form code.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [get-form-version-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-version-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `form_code` | query | String | yes | Unique identifier of the form, typically used to distinguish different form definitions. |

**Response 200**

Array of `string`.

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Form Content

`GET /broker/forms/get`

> Retrieves the JSON schema for the specified form code and version.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [get-form-content.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-content.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `form_code` | query | String | yes | Unique identifier of the form, typically used to distinguish different form definitions. |
| `version` | query | String | yes | Version number of the form definition |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `required` | array<string> |  | List of required property names for the current schema object. |
| `type` | string |  | Enumeration of supported JSON Schema property types used for form definition and validation. — one of: `string`, `number`, `boolean`, `object`, `array` |
| `enum_values` | string |  | Enumerated values allowed for this property. Applies when the schema defines a fixed set of permitted values. |
| `format` | string |  | Format keyword that further constrains the property type (e.g., date, date-time). |
| `description` | string |  | Description of the property, explaining its purpose or usage within the schema. |
| `example` | string |  | Representative example value to illustrate the expected input format. |
| `min_items` | integer |  | Minimum number of items allowed in the array. Applicable only when the property type is array. |
| `max_items` | integer |  | Maximum number of items allowed in the array. Applicable only when the property type is array. |
| `max_length` | integer |  | Maximum length allowed for string-type properties. |
| `form_code` | string | yes | Unique identifier of the form, typically used to distinguish different form definitions. |
| `version` | string | yes | Version number of the form definition, used for compatibility and evolutionary control. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Upload Document

`POST /broker/documents/upload`

> Uploads a document file. Supports multipart/form-data requests.Supported file types: image/jpeg, image/png, application/pdf. Maximum file size: 5MB.

| | |
|---|---|
| **SDK** | `brokerfd.UploadDocument` |
| **Reference** | [document-upload.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/document-upload.md) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `document_id` | string |  | Unique identifier of the uploaded document. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Download Document

`GET /broker/documents/download`

> Downloads a document by document_id.

| | |
|---|---|
| **SDK** | `brokerfd.DownloadDocument` |
| **Reference** | [document-download.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/document-download.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `document_id` | query | string | yes | Unique identifier of the uploaded document. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `doc_name` | string |  | Original name of the document file. |
| `doc_type` | string |  | Enumeration of supported document types used for identity verification, profile updates, and compliance checks. — one of: `FACE_IMG`, `ID_DOCUMENT`, `ADDRESS_PROOF`, `OTHER` |
| `mime_type` | string |  | Enumeration of supported file MIME types for document upload and download. — one of: `image/jpeg`, `image/png`, `application/pdf` |
| `doc_content_base64` | string |  | Base64-encoded content of the document file. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Assets Summary

`GET /broker/assets/summaries/get`

> Retrieves summary-level asset information for a specified account, including available cash, position market value, and total asset value.

| | |
|---|---|
| **SDK** | `brokerfd.GetAccountsSummary / GetFDAssetsSummary` |
| **Reference** | [summary.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/summary.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `balance` | object | yes | Account Balance |
| `positions` | array<object> | yes | Account Position |

*Nested — `balance`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `total_asset_currency` | string | yes | Currency — one of: `USD` |
| `total_cash_balance` | string | yes | Cash Balance, |
| `total_market_value` | string |  | Total holding market value |
| `total_unrealized_profit_loss` | string |  | The unrealized profit or loss of open positions for the specified account. |
| `total_net_liquidation_value` | string |  | Net Account Value, |
| `total_day_profit_loss` | string |  | Day's P&L Not returned at the Omnibus master account level. |
| `maintenance_margin` | string |  | Maintenance Margin（Margin Account） Only returned for margin accounts. |
| `open_margin_calls` | string |  | Open Margin Calls: Only returned for margin accounts. RM: An RM call is triggered when an account's margin equity falls below the maintenance requirement, often due to a decline in the value of the account's positions or an increase in margin requirements for the holdings. RT: A Reg T call occurs when there is not enough equity in an account to cover the 50% initial margin requirement. This can happen due to open positions depreciating or from holding positions overnight that were opened with day trade buying power. EM: An EM call is triggered when a flagged Pattern Day Trader (PDT) account closes the prior business day below the $25,000 minimum Net Account Value (NAV) requirement. Only positions held in the margin account will count toward this requirement. Crypto, futures, event contracts, and any assets held outside the margin account are excluded from the calculation. DT: A Day Trade (DT) Call is issued when you exceed your available Day Trading Buying Power (DTBP) and then place a day trade. A DT call may also occur if you execute a day trade while an Equity Maintenance (EM) call is still active on your account. Exceeding day trade buying power most commonly results from trading on intraday profits. For more details on how buying power is replenished, please refer to our Margin Buying Power page. — one of: `EM`, `RM`, `RT`, `DT` |
| `account_currency_assets` | array<object> | yes | Currency assets Details |

*Nested — `account_currency_assets`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `currency` | string | yes | Currency — one of: `USD` |
| `cash_balance` | string | yes | Cash Balance. |
| `settled_cash` | string |  | Settled Cash， Not returned for event contract accounts. Not returned for margin accounts. |
| `unsettled_cash` | string |  | Unsettled Cash， Not returned for event contract accounts. Not returned for margin accounts. |
| `market_value` | string |  | holding market value |
| `held_amount` | string |  | In-transit funds Not returned for event contract accounts. |
| `buying_power` | string |  | Buying Power Not returned for margin accounts. |
| `day_buying_power` | string |  | Day-Trade Buying Power（Margin Account） Only returned for margin accounts. |
| `overnight_buying_power` | string |  | Overnight BP（Margin Account） Only returned for margin accounts. |
| `night_trading_buying_power` | string |  | Night Trading Buying Power Not returned for event contract accounts. |
| `day_profit_loss` | string |  | Day's P&L Not returned at the Omnibus master account level. |
| `unrealized_profit_loss` | string |  | Open P&L |
| `available_withdrawal` | string |  | The amount of funds currently available for withdrawal. |
| `interests_unpaid` | string |  | Interest to be paid Not returned for event contract accounts. |
| `net_liquidation_value` | string |  | Net Account Value |

*Nested — `positions`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `position_id` | string | yes | Position ID |
| `currency` | string | yes | Currency — one of: `USD` |
| `quantity` | string | yes | Quantity of the position |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market. |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `last_price` | string | yes | Last Price |
| `cost_price` | string | yes | Cost Basis |
| `unrealized_profit_loss` | string | yes | Open P&L |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Assets Detail

`GET /broker/assets/balances/get`

> Retrieves the balance information for a specified account.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDAssetsDetail` |
| **Reference** | [account-balance.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/account-balance.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `total_asset_currency` | string | yes | Currency — one of: `USD` |
| `total_cash_balance` | string | yes | Cash Balance, |
| `total_market_value` | string |  | Total holding market value |
| `total_unrealized_profit_loss` | string |  | The unrealized profit or loss of open positions for the specified account. |
| `total_net_liquidation_value` | string |  | Net Account Value, |
| `total_day_profit_loss` | string |  | Day's P&L Not returned at the Omnibus master account level. |
| `maintenance_margin` | string |  | Maintenance Margin（Margin Account） Only returned for margin accounts. |
| `open_margin_calls` | string |  | Open Margin Calls: Only returned for margin accounts. RM: An RM call is triggered when an account's margin equity falls below the maintenance requirement, often due to a decline in the value of the account's positions or an increase in margin requirements for the holdings. RT: A Reg T call occurs when there is not enough equity in an account to cover the 50% initial margin requirement. This can happen due to open positions depreciating or from holding positions overnight that were opened with day trade buying power. EM: An EM call is triggered when a flagged Pattern Day Trader (PDT) account closes the prior business day below the $25,000 minimum Net Account Value (NAV) requirement. Only positions held in the margin account will count toward this requirement. Crypto, futures, event contracts, and any assets held outside the margin account are excluded from the calculation. DT: A Day Trade (DT) Call is issued when you exceed your available Day Trading Buying Power (DTBP) and then place a day trade. A DT call may also occur if you execute a day trade while an Equity Maintenance (EM) call is still active on your account. Exceeding day trade buying power most commonly results from trading on intraday profits. For more details on how buying power is replenished, please refer to our Margin Buying Power page. — one of: `EM`, `RM`, `RT`, `DT` |
| `account_currency_assets` | array<object> | yes | Currency assets Details |

*Nested — `account_currency_assets`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `currency` | string | yes | Currency — one of: `USD` |
| `cash_balance` | string | yes | Cash Balance. |
| `settled_cash` | string |  | Settled Cash， Not returned for event contract accounts. Not returned for margin accounts. |
| `unsettled_cash` | string |  | Unsettled Cash， Not returned for event contract accounts. Not returned for margin accounts. |
| `market_value` | string |  | holding market value |
| `held_amount` | string |  | In-transit funds Not returned for event contract accounts. |
| `buying_power` | string |  | Buying Power Not returned for margin accounts. |
| `day_buying_power` | string |  | Day-Trade Buying Power（Margin Account） Only returned for margin accounts. |
| `overnight_buying_power` | string |  | Overnight BP（Margin Account） Only returned for margin accounts. |
| `night_trading_buying_power` | string |  | Night Trading Buying Power Not returned for event contract accounts. |
| `day_profit_loss` | string |  | Day's P&L Not returned at the Omnibus master account level. |
| `unrealized_profit_loss` | string |  | Open P&L |
| `available_withdrawal` | string |  | The amount of funds currently available for withdrawal. |
| `interests_unpaid` | string |  | Interest to be paid Not returned for event contract accounts. |
| `net_liquidation_value` | string |  | Net Account Value |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Positions

`GET /broker/assets/positions/list`

> Retrieves positions according to the account ID.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDPositions / GetPositions` |
| **Reference** | [account-position.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/account-position.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `position_id` | string | yes | Position ID |
| `currency` | string | yes | Currency — one of: `USD` |
| `quantity` | string | yes | Quantity of the position |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market. |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `last_price` | string | yes | Last Price |
| `cost_price` | string | yes | Cost Basis |
| `unrealized_profit_loss` | string | yes | Open P&L |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Cash Activities By Type

`GET /broker/activities/cash-activities/list`

> Retrieves account transaction activities records with details.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDActivities` |
| **Reference** | [broker-cash-activity-by-type.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-cash-activity-by-type.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Provide the target account id |
| `activity_types` | query | string |  | Account activity types — one of: `TRADE`, `DEPOSIT`, `WITHDRAW`, `FEES`, `JOURNAL`, `TRANSFER`, `INTERESTS`, `EC_STATEMENT` |
| `start_time` | query | String |  | Activity query start time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z' |
| `end_time` | query | String |  | If not provided, the default query is the last 7 days. Activity query end time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z' |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | string | yes | Unique ID |
| `account_id` | string | yes | Account ID |
| `activity_type` | string | yes | Activity Type \| Code \| Description \| \|---------------\|---------------------------------\| \|TRADE \| Trade Activity Type \| \|DEPOSIT \| Deposit Activity Type \| \|WITHDRAW \| Withdraw Activity Type \| \|FEES \| Fees Activity Type \| \|JOURNAL \| Journal Activity Type \| \|TRANSFER \| Transfer Activity Type \| \|INTERESTS \| Interests Activity Type \| \|EC_STATEMENT \| EC Statement Activity Type \| \|OTHER \| Other Activity Type \| — one of: `TRADE`, `DEPOSIT`, `WITHDRAW`, `FEES`, `JOURNAL`, `TRANSFER`, `INTERESTS`, `EC_STATEMENT` |
| `activity_sub_type` | string | yes | Activity Type and Sub Type Mapping. \| ActivityType \| ActivitySubType \| \|----------------\|-------------------------------------\| \| TRADE \| BTO \| \| TRADE \| STC \| \| DEPOSIT \| ACH \| \| DEPOSIT \| WIRE \| \| DEPOSIT \| WIRE_REVERSE \| \| DEPOSIT \| ACH_REVERSE \| \| WITHDRAW \| ACH \| \| WITHDRAW \| WIRE \| \| WITHDRAW \| ACH_REVERSE \| \| WITHDRAW \| WIRE_REVERSE \| \| FEES \| WIRE_DEPOSIT \| \| FEES \| WIRE_WITHDRAW \| \| FEES \| ACH_REVERSAL \| \| FEES \| WIRE_REVERSAL \| \| FEES \| SERVICE_FEE \| \| FEES \| ADR \| \| FEES \| TRANSFER_ACATS \| \| FEES \| PAPER_CONFIRM_FEE \| \| FEES \| PAPER_STATEMENT_FEE \| \| FEES \| WRITE_OFF \| \| FEES \| PROCESSING_FEE \| \| FEES \| COMMISSION \| \| FEES \| TRANSFER_FOP \| \| FEES \| MONTH_END_REPORT_CHARGES \| \| FEES \| FINANCIAL_TRANSACTION_TAX \| \| FEES \| SETTLEMENT_FEE \| \| FEES \| TRANSACTION_FEE \| \| JOURNAL \| CASH_JOURNAL \| \| JOURNAL \| CREDIT \| \| TRANSFER \| INTERNAL_THIRD_PARTY_TRANSFER \| \| INTERESTS \| CREDIT_CASH \| \| EC_STATEMENT \| EC_EXPIRATION \| \| EC_STATEMENT \| EC_PAYOUT \| \| OTHER \| OTHER \| — one of: `BTO`, `STC`, `WIRE`, `ACH`, `WIRE_REVERSE`, `ACH_REVERSE`, `WIRE_DEPOSIT`, `WIRE_WITHDRAW`, `ACH_REVERSAL`, `WIRE_REVERSAL`, `CASH_JOURNAL`, `INTERNAL_THIRD_PARTY_TRANSFER`, `CREDIT`, `CREDIT_CASH`, `EC_EXPIRATION`, `EC_PAYOUT`, `OTHER` |
| `currency` | string | yes | Currency — one of: `USD` |
| `market` | string |  | Market Code US - US Market — one of: `US` |
| `symbol` | string |  | Activity Symbol |
| `trade_date` | string | yes | Accounting date of the transaction (trade date), format: yyyy-MM-dd |
| `net_amount` | string | yes | Net change amount of the transaction (positive for credit, negative for debit) |
| `biz_time` | string | yes | Business event time when the transaction occurred |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Bank Relationship

`POST /broker/funding/bank-relationships/create`

> Creates a Bank Relationship for an Account.

| | |
|---|---|
| **SDK** | `brokerfd.AddFDBankAccount` |
| **Reference** | [create-bank-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-bank-relationship.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account identifier. |
| `bank_name` | string |  | Required if bank_code_type = BIC, Name of recipient bank |
| `bank_code` | string | yes | Alphanumeric ABA RTN (Routing Number) |
| `bank_code_type` | string | yes | ABA (US based bank accounts)、BIC — one of: `ABA`, `BIC` |
| `bank_account_name` | string | yes | Name of the bank account holder, as registered with the bank |
| `bank_account_number` | string | yes | Bank account number. |
| `country` | string |  | Required if bank_code_type = BIC, The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state_province` | string |  | Required if bank_code_type = BIC |
| `postal_code` | string |  | Required if bank_code_type = BIC |
| `city` | string |  | Required if bank_code_type = BIC |
| `intermediary_bank_details` | object |  | Intermediary bank details. |

*Nested — `intermediary_bank_details`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `intermediary_bank_name` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Intermediary Bank Name |
| `intermediary_bank_code` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Alphanumeric ABA RTN (Routing Number) |
| `intermediary_bank_account_number` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Bank account number. |
| `intermediary_bank_address` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Bank Address. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `bank_relationship_id` | string | yes | Relationship ID |
| `bank_code` | string | yes | Alphanumeric ABA RTN (Routing Number) |
| `bank_name` | string |  | Bank name |
| `bank_account_name` | string | yes | Name of the bank account holder, as registered with the bank |
| `bank_account_number` | string | yes | Bank account number |
| `bank_code_type` | string | yes | ABA (US based bank accounts)、BIC — one of: `ABA`, `BIC` |
| `country` | string |  | Only for Non-US based bank accounts, The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state_province` | string |  | Only for Non-US based bank accounts |
| `postal_code` | string |  | Only for Non-US based bank accounts |
| `city` | string |  | Only for Non-US based bank accounts |
| `status` | string | yes | Status — one of: `APPROVED` |
| `intermediary_bank_details` | object |  | Intermediary bank details. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

*Nested — `intermediary_bank_details`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `intermediary_bank_name` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Intermediary Bank Name |
| `intermediary_bank_code` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Alphanumeric ABA RTN (Routing Number) |
| `intermediary_bank_account_number` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Bank account number. |
| `intermediary_bank_address` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Bank Address. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Delete Bank Relationship

`POST /broker/funding/bank-relationships/delete`

> Deletes a Bank Relationship for an Account.

| | |
|---|---|
| **SDK** | `brokerfd.RemoveFDBankAccount` |
| **Reference** | [delete-bank-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/delete-bank-relationship.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `bank_relationship_id` | string | yes | Bank relationship ID. The bank_relationship_id created from the Create Bank Relationship interface |
| `account_id` | string | yes | Account identifier. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `bank_relationship_id` | string | yes | Bank relationship ID. The bank_relationship_id created from the Create Bank Relationship interface |
| `account_id` | string | yes | Account identifier. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List Bank Accounts

`GET /broker/funding/bank-relationships/list`

> Retrieves Bank Relationships for an Account.

| | |
|---|---|
| **SDK** | `brokerfd.ListFDBankAccounts` |
| **Reference** | [list-linked-bank-accounts.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-linked-bank-accounts.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `bank_relationship_id` | string | yes | Relationship ID |
| `bank_code` | string | yes | Alphanumeric ABA RTN (Routing Number) |
| `bank_name` | string |  | Bank name |
| `bank_account_name` | string | yes | Name of the bank account holder, as registered with the bank |
| `bank_account_number` | string | yes | Bank account number |
| `bank_code_type` | string | yes | ABA (US based bank accounts)、BIC — one of: `ABA`, `BIC` |
| `country` | string |  | Only for Non-US based bank accounts, The value must be a valid ISO 3166-1 alpha-3 country code in uppercase format (3-letter code). |
| `state_province` | string |  | Only for Non-US based bank accounts |
| `postal_code` | string |  | Only for Non-US based bank accounts |
| `city` | string |  | Only for Non-US based bank accounts |
| `status` | string | yes | Status — one of: `APPROVED` |
| `intermediary_bank_details` | object |  | Intermediary bank details. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

*Nested — `intermediary_bank_details`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `intermediary_bank_name` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Intermediary Bank Name |
| `intermediary_bank_code` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Alphanumeric ABA RTN (Routing Number) |
| `intermediary_bank_account_number` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Bank account number. |
| `intermediary_bank_address` | string |  | This is an optional field;it is only required if an intermediary bank is involved in the transfer. Bank Address. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create ACH Relationship

`POST /broker/funding/ach-relationships/create`

> Creates an ACH Relationship for an Account. Only one ACH relationship is allowed per account, and creating a new one automatically replaces the existing one.

| | |
|---|---|
| **SDK** | `brokerfd.AddFDAchAccount` |
| **Reference** | [create-ach-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-ach-relationship.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account identifier. |
| `token_type` | string | yes | Token Type only support 'Plaid' — one of: `Plaid` |
| `processor_token` | string | yes | Using Plaid, you can specify a Plaid processor token here. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `ach_relationship_id` | string | yes | Relationship ID. The ach_relationship_id created from the Create ACH Relationship interface. |
| `account_id` | string | yes | Account identifier. |
| `account_owner_name` | string | yes | Legal full name of the bank account owner. |
| `bank_account_type` | string | yes | Must be CHECKING or SAVINGS — one of: `CHECKING`, `SAVINGS` |
| `bank_account_number` | string | yes | Bank account number |
| `bank_routing_number` | string | yes | ABA routing transit number (RTN), Alphanumeric identifier for the financial institution, used for ACH transaction routing. |
| `status` | string | yes | Status — one of: `APPROVED` |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Delete ACH Relationship

`POST /broker/funding/ach-relationships/delete`

> Deletes an ACH Relationship for an Account.

| | |
|---|---|
| **SDK** | `brokerfd.RemoveFDAchAccount` |
| **Reference** | [delete-ach-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/delete-ach-relationship.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `ach_relationship_id` | string | yes | ACH relationship ID |
| `account_id` | string | yes | Account identifier. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `ach_relationship_id` | string | yes | Relationship ID. The ach_relationship_id created from the Create ACH Relationship interface. |
| `account_id` | string | yes | Account identifier. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List ACH Relationships

`GET /broker/funding/ach-relationships/list`

> Retrieves ACH Relationships for an Account.

| | |
|---|---|
| **SDK** | `brokerfd.ListFDAchAccounts` |
| **Reference** | [list-ach-relationships.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-ach-relationships.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `ach_relationship_id` | string | yes | Relationship ID. The ach_relationship_id created from the Create ACH Relationship interface. |
| `account_id` | string | yes | Account identifier. |
| `account_owner_name` | string | yes | Legal full name of the bank account owner. |
| `bank_account_type` | string | yes | Must be CHECKING or SAVINGS — one of: `CHECKING`, `SAVINGS` |
| `bank_account_number` | string | yes | Bank account number |
| `bank_routing_number` | string | yes | ABA routing transit number (RTN), Alphanumeric identifier for the financial institution, used for ACH transaction routing. |
| `status` | string | yes | Status — one of: `APPROVED` |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Transfer

`POST /broker/funding/transfers/create`

> Initiates a funding transfer to an account. Related Event Notifications: Please refer to the Event Push API documentation: Funding Transfer Events

| | |
|---|---|
| **SDK** | `brokerfd.InitiateFDTransfer` |
| **Reference** | [create-transfer.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-transfer.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account ID |
| `transfer_type` | string | yes | Transfer type,current support WIRE or ACH — one of: `ACH`, `WIRE` |
| `ach_relationship_id` | string |  | Required if transfer_type = ACH |
| `bank_relationship_id` | string |  | Required if transfer_type = WIRE |
| `amount` | string | yes | Transfer amount |
| `currency` | string | yes | Currency code |
| `direction` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account ID |
| `transfer_id` | string | yes | Transfer id generated by the system. |
| `transfer_type` | string | yes | transfer type — one of: `ACH`, `WIRE` |
| `ach_relationship_id` | string |  | Required if transfer_type is ACH |
| `bank_relationship_id` | string |  | Required if transfer_type is WIRE. |
| `amount` | string | yes | Transfer amount |
| `currency` | string | yes | Currency — one of: `USD` |
| `direction` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |
| `status` | string | yes | Transfer status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `COMPLETED`, `RETURNED` |
| `reason` | string |  | reason of the status |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `fee` | string |  | Fee amount to be collected. Only applies when type is WIRE |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Transfer List

`GET /broker/funding/transfers/list`

> Retrieves account fund transfers by specified criteria. Transfer records are sorted by client_request_id. If the client_request_id values in the request are provided in order, the response will preserve the same ordering.

| | |
|---|---|
| **SDK** | `brokerfd.ListFDTransfers` |
| **Reference** | [transfer-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/transfer-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |
| `start_time` | query | String |  | The start datetime to include in the date range, formatted as yyyy-MM-dd'T'HH:mm:ss.SSSZ. Both start_time and end_time should be provided, end_time should be greater than start_time. By default, records from the last 7 days will be queried if no time range is specified. |
| `end_time` | query | String |  | The end datetime to include in the date range, formatted as yyyy-MM-dd'T'HH:mm:ss.SSSZ. Both start_time and end_time should be provided, end_date should be greater than start_time. By default, records from the last 7 days will be queried if no time range is specified. |
| `transfer_types` | query | string |  | Transfer types — one of: `ACH`, `WIRE` |
| `directions` | query | string |  | Transfer directions — one of: `DEPOSIT`, `WITHDRAWAL` |
| `statuses` | query | string |  | Transfer statuses — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `COMPLETED`, `RETURNED` |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account ID |
| `transfer_id` | string | yes | Transfer id generated by the system. |
| `transfer_type` | string | yes | transfer type — one of: `ACH`, `WIRE` |
| `ach_relationship_id` | string |  | Required if transfer_type is ACH |
| `bank_relationship_id` | string |  | Required if transfer_type is WIRE. |
| `amount` | string | yes | Transfer amount |
| `currency` | string | yes | Currency — one of: `USD` |
| `direction` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |
| `status` | string | yes | Transfer status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `COMPLETED`, `RETURNED` |
| `reason` | string |  | reason of the status |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `fee` | string |  | Fee amount to be collected. Only applies when type is WIRE |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Transfer Detail

`GET /broker/funding/transfers/get`

> Retrieves Account Fund Transfer Details.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDTransferDetail` |
| **Reference** | [transfer-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/transfer-detail.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `client_request_id` | query | String | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | query | String | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account ID |
| `transfer_id` | string | yes | Transfer id generated by the system. |
| `transfer_type` | string | yes | transfer type — one of: `ACH`, `WIRE` |
| `ach_relationship_id` | string |  | Required if transfer_type is ACH |
| `bank_relationship_id` | string |  | Required if transfer_type is WIRE. |
| `amount` | string | yes | Transfer amount |
| `currency` | string | yes | Currency — one of: `USD` |
| `direction` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |
| `status` | string | yes | Transfer status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `COMPLETED`, `RETURNED` |
| `reason` | string |  | reason of the status |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `fee` | string |  | Fee amount to be collected. Only applies when type is WIRE |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Cancel Transfer

`POST /broker/funding/transfers/cancel`

> Cancels a fund transfer.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [cancel-transfer.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/cancel-transfer.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account ID |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `transfer_id` | string | yes | tranfer_id |
| `account_id` | string | yes | Account identifier. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Instant Funding

`POST /broker/funding/instant-funding/create`

> Initiates an immediate deposit request. The system instantly processes and credits the funds to the target account, making them available for trading immediately. No settlement institution is involved in this instant funding process. Related Event Notifications: Please refer to the Event Push API documentation: Instant Funding Events

| | |
|---|---|
| **SDK** | `brokerfd.CreateFDInstantFunding` |
| **Reference** | [broker-funding-instant-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-instant-create.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | string | yes | Account ID |
| `type` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. |
| `account_id` | string | yes | Account ID |
| `type` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |
| `instant_funding_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Instant Funding Detail

`GET /broker/funding/instant-funding/get`

> Retrieves instant funding record details.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDInstantFundingDetail` |
| **Reference** | [broker-funding-instant-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-instant-query.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `client_request_id` | query | String | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | query | String | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. |
| `account_id` | string | yes | Account ID |
| `type` | string | yes | Transfer direction — one of: `DEPOSIT`, `WITHDRAWAL` |
| `instant_funding_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Fee

`POST /broker/fees/create`

> Initial Fee Deduction debits a specified account and credits an equal amount to a contra account, without involving any settlement institutions.Use this API to deduct various fees (such as ADR, SERVICE_FEE, etc.) from a specified account. Related Event Notifications: Please refer to the Event Push API documentation: Fee Deduction Events

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [broker-funding-fee-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-fee-create.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account` | string | yes | Account id. |
| `contra_account` | string | yes | Firm Account |
| `type` | string | yes | Fee Type — one of: `SERVICE_FEE`, `ADR`, `TRANSFER_ACATS`, `PAPER_CONFIRM_FEE`, `PAPER_STATEMENT_FEE`, `WRITE_OFF`, `PROCESSING_FEE`, `COMMISSION`, `TRANSFER_FOP`, `MONTH_END_REPORT_CHARGES`, `FTT`, `SETTLEMENT_FEES`, `TRANSACTION_FEES` |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account` | string | yes | Account id. |
| `contra_account` | string | yes | Firm Account |
| `type` | string | yes | Fee Type — one of: `SERVICE_FEE`, `ADR`, `TRANSFER_ACATS`, `PAPER_CONFIRM_FEE`, `PAPER_STATEMENT_FEE`, `WRITE_OFF`, `PROCESSING_FEE`, `COMMISSION`, `TRANSFER_FOP`, `MONTH_END_REPORT_CHARGES`, `FTT`, `SETTLEMENT_FEES`, `TRANSACTION_FEES` |
| `fee_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Get Fee

`GET /broker/fees/get`

> Retrieves fee record details..

| | |
|---|---|
| **SDK** | `brokerfd.GetFDTransferFees` |
| **Reference** | [broker-funding-fee-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-fee-query.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `client_request_id` | query | String | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | query | String | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account` | string | yes | Account id. |
| `contra_account` | string | yes | Firm Account |
| `type` | string | yes | Fee Type — one of: `SERVICE_FEE`, `ADR`, `TRANSFER_ACATS`, `PAPER_CONFIRM_FEE`, `PAPER_STATEMENT_FEE`, `WRITE_OFF`, `PROCESSING_FEE`, `COMMISSION`, `TRANSFER_FOP`, `MONTH_END_REPORT_CHARGES`, `FTT`, `SETTLEMENT_FEES`, `TRANSACTION_FEES` |
| `fee_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Credit

`POST /broker/credits/create`

> Initial Fee Credit credits a specified account and deducts an equal amount from a contra account, without involving any settlement institutions. Use this API to credit funds to a specified account (such as GIFTING). Related Event Notifications: Please refer to the Event Push API documentation: Credit Events

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [broker-funding-credit-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-credit-create.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account` | string | yes | Account id. |
| `contra_account` | string | yes | Firm Account |
| `type` | string | yes | Credit Type — one of: `GIFTING` |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account` | string | yes | Account id. |
| `contra_account` | string | yes | Firm Account |
| `type` | string | yes | Credit Type — one of: `GIFTING` |
| `credit_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Get Credit

`GET /broker/credits/get`

> Retrieves credit record details.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDCreditInfo` |
| **Reference** | [broker-funding-credit-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-credit-query.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `client_request_id` | query | String | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | query | String | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account` | string | yes | Account id. |
| `contra_account` | string | yes | Firm Account |
| `type` | string | yes | Credit Type — one of: `GIFTING` |
| `credit_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount. |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List Stock Instruments

`GET /broker/instruments/stocks/profiles/list`

> Retrieves profile information for one or more instruments. Related Event Notifications: Please refer to the Event Push API documentation: Instrument Events

| | |
|---|---|
| **SDK** | `brokerfd.GetFDStockInstruments` |
| **Reference** | [list-stock-instruments.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-stock-instruments.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string |  | List of security symbols, maximum 100 symbols per query. |
| `category` | query | string | yes | Security type. — one of: `US_STOCK` |
| `sub_category` | query | string |  | Sub-category of the instrument. Only effective when symbols is not specified. When category = US_STOCK, supported values: COMMON_STOCK, ETF, PREFERRED_STOCK, WARRANT, UNITS, RIGHT. If not specified, returns all sub-categories. — one of: `COMMON_STOCK`, `ETF`, `PREFERRED_STOCK`, `WARRANT`, `UNITS`, `RIGHT` |
| `status` | query | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Symbol name, e.g. Apple |
| `instrument_id` | string |  | Unique identifier of the security |
| `exchange_code` | string |  | Exchange code, e.g. CCC |
| `category` | string |  | Instrument type, e.g. US_STOCK — one of: `US_STOCK` |
| `symbol` | string |  | Symbol of the instrument |
| `status` | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `shortable` | boolean |  | Instrument is shortable or not |
| `fractionable` | boolean |  | Instrument is fractionable or not |
| `marginable` | boolean |  | Instrument is marginable or not |
| `overnight_trading_supported` | boolean |  | Instrument support overnight trading or not |
| `margin_requirement_long` | string |  | Margin requirement ratio for long position. 0.5 represents 50% |
| `margin_requirement_short` | string |  | Margin requirement ratio for short position. 0.5 represents 50% |
| `intraday_margin_long` | string |  | Intraday margin requirement ratio for long position |
| `intraday_margin_short` | string |  | Intraday margin requirement ratio for short position |
| `maintenance_margin_long` | string |  | Maintenance margin requirement ratio for long position |
| `maintenance_margin_short` | string |  | Maintenance margin requirement ratio for short position |
| `easy_to_borrow` | boolean |  | Instrument is easy to borrow or not |
| `lot_size` | string |  | Lot size |
| `currency` | string |  | currency |
| `sub_category` | string |  | Sub-category of the instrument. — one of: `COMMON_STOCK`, `ETF`, `PREFERRED_STOCK`, `WARRANT`, `UNITS`, `RIGHT` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Categories

`GET /broker/instruments/event-contracts/categories/list`

> Retrieves all categories under the Event Contract.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [broker-event-categories-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-categories-list.md) |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `category_id` | integer | yes | Category ID. |
| `category_code` | string | yes | Unique identifier code for category. |
| `category_name` | string | yes | Category corresponding name. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Series

`GET /broker/instruments/event-contracts/series/list`

> Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., “Monthly Jobs Report” ). This endpoint allows you to browse and discover available series templates by category.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [broker-event-series-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-series-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | The category which this series belongs to. — one of: `ECONOMICS`, `FINANCIALS`, `POLITICS`, `ENTERTAINMENT`, `SCIENCE_TECHNOLOGY`, `CLIMATE_WEATHER`, `TRANSPORTATION`, `CRYPTO`, `SPORTS` |
| `symbols` | query | string |  | List of series symbols, maximum 100 symbols per query. |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Series Symbol |
| `name` | string |  | Series Name |
| `category` | string |  | Category Code |
| `frequency` | string | yes | frequency, e.g., HOURLY,DAILY,WEEKLY,MONTHLY,ANNUAL,ONE_OFF,CUSTOM. |
| `series_id` | integer |  | Series ID |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Events

`GET /broker/instruments/event-contracts/events/list`

> Retrieves events under the Event Contract.

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [broker-event-events-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-events-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `series_symbol` | query | string | yes | Symbol that identifies this series. |
| `symbols` | query | string |  | List of events symbols, maximum 100 symbols per query. |
| `status` | query | string |  | The status of the event. — one of: `ACTIVE`, `INACTIVE` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Event symbol |
| `name` | string |  | Event name |
| `status` | string |  | status, eg:ACTIVE=event is open for trading; INACTIVE=event is closed or settled, default:ACTIVE |
| `series_id` | integer |  | Series ID |
| `short_name` | string |  | Event short name |
| `strike_date` | string |  | strike date, Only present for sports category events. Represents the settlement date/period of the event. |
| `strike_period` | string |  | strike period, Only present for sports category events. Represents the settlement date/period of the event. |
| `mutually_exclusive` | boolean |  | mutually exclusive, true or false |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Instruments

`GET /broker/instruments/event-contracts/markets/list`

> Retrieves profile information for event contract markets based on the series symbol.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDECInstruments` |
| **Reference** | [broker-event-market-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-market-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `series_symbol` | query | string | yes | Symbol that identifies this series. |
| `event_symbol` | query | string |  | Symbol of the event events. |
| `symbols` | query | string |  | List of security symbols, maximum 100 symbols per query. |
| `expiration_date_after` | query | string |  | Used to filter items whose expiration date is later than a specified date; the default selection is the current day (inclusive). |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `series_id` | string | yes | ID that identifies this series. |
| `series_symbol` | string | yes | Symbol that identifies this series. |
| `series_name` | string | yes | Name that describes the series. |
| `event_symbol` | string | yes | Symbol that identifies this events. |
| `event_name` | string | yes | Name that describes the events. |
| `instrument_id` | string | yes | Unique id of the event market. |
| `symbol` | string | yes | Symbol of the event market. |
| `name` | string | yes | Name of the event market. |
| `yes_condition` | string | yes | Conditions for a 'Yes' outcome. |
| `last_trading_date` | string | yes | Last Notice Day. |
| `status` | string | yes | Listing status. — one of: `NOT_SET`, `LISTING`, `DELISTING`, `OTHER`, `UNRECOGNIZED` |
| `tradable_status` | string | yes |  — one of: `OC`, `CO`, `NT` |
| `can_close_early` | boolean | yes | Can the contract close early? |
| `expected_exp_date` | string | yes | Expected expiration date of contract. |
| `latest_exp_date` | string | yes | Latest expiration date. |
| `payout_date` | string | yes | Settlement/Payment Date. |
| `fractionable` | boolean | yes | Support fragmented event contracts. |
| `price_ranges` | array<object> |  | Price range. |

*Nested — `price_ranges`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `start` | string | yes | Start price. |
| `end` | string | yes | End price. |
| `step` | string | yes | Step length. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Corporate Actions Detail

`GET /broker/instruments/stocks/corporate-actions/get`

> Retrieves corporate actions detail information,only support stocks. Related Event Notifications: Please refer to the Event Push API documentation: Corporate Actions Events

| | |
|---|---|
| **SDK** | `brokerfd.GetFDCorporateActions` |
| **Reference** | [broker-corporate-actions-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-corporate-actions-detail.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `event_id` | query | string | yes | Corporate Event Action ID |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `event_id` | string |  | Company Event ID |
| `event_type` | string |  | Corporate Event Type IDENTIFIER_CHANGE - Identifier-related changes to a security without economic impact, including symbol, exchange, ISIN, CUSIP, or instrument name changes DIVIDEND - Dividend and distribution events, including cash dividend, stock dividend, optional dividend, and return of capital REVERSE_SPLIT - Reverse stock split that consolidates shares and reduces the number of outstanding shares FORWARD_SPLIT - Forward stock split that increases the number of outstanding shares BONUS_ISSUE - Bonus issue of additional shares distributed to existing shareholders at no cost RIGHTS_OFFERING - Rights offering allowing shareholders to subscribe for additional shares DISTRIBUTION - Distribution of cash, securities, or other assets to shareholders SPIN_OFF - Spin-off event where shares of a subsidiary or new entity are distributed to existing shareholders UNIT_SPLIT - Unit split event affecting composite or unit-based securities MERGER - Merger or acquisition event involving the combination of two or more entities FULL_CALL - Full call redemption of the entire outstanding security issue PARTIAL_CALL - Partial call redemption affecting only a portion of the outstanding issue EXCHANGE - Exchange event where existing securities are exchanged for new securities or other consideration DTC_EXIT - Event indicating a security is no longer eligible for DTC settlement or custody LIQUIDATION - Liquidation event involving the winding up of an issuer and asset distribution WORTHLESS - Worthless security event indicating the security has no residual value ADR_GDR_TERMINATION - Termination of an ADR or GDR program MATURITY - Maturity event where a security reaches its contractual maturity date ADR_FEE - ADR fee charged to holders of American Depositary Receipts CONVERSION - Conversion event where securities are converted into another class or form OPEN_OFFER - Open offer allowing shareholders to subscribe for additional securities PREFERENTIAL_OFFER - Preferential offer made to selected shareholders under specific terms PERFORMANCE_COMPENSATION - Performance compensation event related to performance commitments, commonly in A-share markets DELISTING - Delisting event where a security is removed from exchange trading — one of: `IDENTIFIER_CHANGE`, `DIVIDEND`, `REVERSE_SPLIT`, `FORWARD_SPLIT`, `BONUS_ISSUE`, `RIGHTS_OFFERING`, `DISTRIBUTION`, `SPIN_OFF`, `UNIT_SPLIT`, `MERGER`, `FULL_CALL`, `PARTIAL_CALL`, `EXCHANGE`, `DTC_EXIT`, `LIQUIDATION`, `WORTHLESS`, `ADR_GDR_TERMINATION`, `MATURITY`, `ADR_FEE`, `CONVERSION`, `OPEN_OFFER`, `PREFERENTIAL_OFFER`, `PERFORMANCE_COMPENSATION`, `DELISTING` |
| `event_version` | string |  | Event Version |
| `instrument_id` | string |  | Instrument ID |
| `category` | string |  | Instrument Category — one of: `US_STOCK` |
| `record_date` | string |  | Record Date (YYYY-MM-DD) |
| `ex_date` | string |  | Ex Date (YYYY-MM-DD) |
| `payment_date` | string |  | Payment Date (YYYY-MM-DD) |
| `final_pay_date` | string |  | Final Payment Date (YYYY-MM-DD) |
| `country_code` | string |  | Country Code, ISO 3166-1 alpha-2 format |
| `listing_country_of_code` | string |  | Listing Country Code, ISO 3166-1 alpha-2 format |
| `issuer_country_code` | string |  | Issuer Country Code, ISO 3166-1 alpha-2 format |
| `from` | object |  | Event From Info. Position's instrument information |
| `to` | array<object> |  | Event To Info. Corporate action's target instrument information |

*Nested — `from`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Symbol of the instrument |
| `name` | string |  | Name of the instrument |
| `exchange` | string |  | Exchange code |

*Nested — `to`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `option_number` | string |  | Option Number. Identifier for the payout option |
| `description` | string |  | Description of the payout option |
| `default_option_flag` | string |  | Default Option Flag. Indicates if this is the default payout option |
| `payouts` | array<object> |  | Payouts associated with this option |

*Nested — `payouts`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Payout Nature Type DV - Dividend: Cash or stock distribution paid to shareholders FR - Franked Dividend: Dividend paid with franking credits attached IN - Interest: Interest income distribution L2 - Long Term Capital Gains: Gains from disposal of assets held longer than one year OT - Other: Other types of income or entitlement (see extended terms) C - Cash: Cash payment (only applicable to events created prior to release 2) PC - Cash/Principal/Return of Capital: Cash returned as principal or capital PM - Premium: Additional amount paid over base entitlement S - Securities: Distribution of securities instead of cash ST - Short Term Capital Gains: Gains from disposal of assets held less than one year SI - Sundry Income: Miscellaneous income distributions UF - Unfranked Dividend: Dividend paid without franking credits PI - Property Income Distribution: Income derived from property assets TD - Tax Deferred: Income deferred for tax purposes TE - Tax Exempted: Income exempted from taxation FI - Foreign Income: Income sourced from foreign jurisdictions CD - Capital gain on disposal of taxable property - Discounted CO - Capital gain on disposal of taxable property - Other CC - Capital gain on disposal of taxable property - Concessional CN - Capital gain on disposal of non-taxable property RT - Royalties: Payment received for intellectual property usage TX - Tax Credit: Credit applied against tax liability BP - Buy Permitted: Security is eligible for purchase under corporate action CL - Cash in Lieu of Fractional Share: Cash payment for fractional share entitlement DF - Drop Fraction: Fractional shares are dropped without compensation EX - Extend and Retain Fractions: Fractional shares are retained and adjusted NC - Round to Nearest Cent: Monetary amounts rounded to the nearest cent NW - Round to Nearest Whole Number if .5 or above: Rounding rule for fractional shares PR - Purchase Required: Mandatory purchase of securities as part of corporate action RD - Round Down to Nearest Whole Number: Fractional shares rounded down RU - Round Up to Nearest Whole Number: Fractional shares rounded up SR - Sale Required: Mandatory sale of securities as part of corporate action BU - Round up: Generic rounding up rule BC - Beneficial Owner Cash in Lieu: Cash payment to beneficial owner for fractional share BD - Beneficial Owner Round Down: Fractional shares of beneficial owner rounded down CT - Security Convert To Cash: Conversion of security into cash (used for tokenized assets) — one of: `DV`, `FR`, `IN`, `L2`, `OT`, `C`, `PC`, `PM`, `S`, `ST`, `SI`, `UF`, `PI`, `TD`, `TE`, `FI`, `CD`, `CO`, `CC`, `CN`, `RT`, `TX`, `BP`, `CL`, `DF`, `EX`, `NC`, `NW`, `PR`, `RD`, `RU`, `SR`, `BU`, `BC`, `BD`, `CT` |
| `pay_type` | string |  | Payout Delivery Type CASH - Cash settlement SECURITY - Security settlement (stock, right, warrant, etc.) SCRIP - Scrip dividend (dividend paid in shares instead of cash) SECURITY_AND_CASH - Combination of security and cash (logical type, not for persistence) — one of: `CASH`, `SECURITY`, `SCRIP`, `SECURITY_AND_CASH` |
| `payout_number` | integer |  | Sequence number of the payout within the same option. Used for ordering and identification. |
| `adr_fee_rate` | string |  | ADR fee rate applied to this payout, if applicable. |
| `fraction_share_rule` | string |  | Fraction Share Rule NONE - No special handling for fractional shares ROUND_DOWN - Round down fractional shares ROUND_UP - Round up fractional shares CASH_IN_LIEU - Cash in lieu for fractional shares DISTRIBUTION - Fractional share distribution STANDARD - Standard rounding for fractional shares — one of: `NONE`, `ROUND_DOWN`, `ROUND_UP`, `CASH_IN_LIEU`, `DISTRIBUTION`, `STANDARD` |
| `cancellation_fee` | string |  | Cancellation fee rate applied if the corporate action is cancelled. |
| `issuance_fee` | string |  | Issuance fee rate applied for newly issued securities. |
| `tax_status` | string |  | IRS income classification for tax reporting purposes. |
| `currency` | string |  | Currency of the cash payout. Applicable only when payType involves CASH. |
| `amount` | string |  | Cash amount paid per share held. Applicable only when payType involves CASH. |
| `withholding_tax_rate` | string |  | Withholding tax rate applied to the cash payout. |
| `symbol` | string |  | Trading symbol of the distributed security. |
| `name` | string |  | Name of the distributed security. |
| `exchange` | string |  | Exchange where the distributed security is listed. |
| `from_ratio` | string |  | Original holding quantity used as the base for ratio calculation. |
| `to_ratio` | string |  | Distributed quantity received for the given base holding. |
| `cash_in_lieu_price` | string |  | Cash-in-lieu price used to settle fractional shares. |
| `reinvest_price` | string |  | Reinvestment price used for scrip dividend calculation. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Preview

`POST /broker/orders/preview`

> Calculates the estimated cost and fees for an order based on the provided parameters. Supports simple orders.

| | |
|---|---|
| **SDK** | `brokerfd.PreviewFDOrder` |
| **Reference** | [common-order-preview.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-preview.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string |  | Account identifier |
| `new_orders` | array<object> |  | Order Details |

*Nested — `new_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). Used to track or reference the order when interacting with the system. |
| `combo_type` | string | yes | Type of order combination. - NORMAL: Indicates a standard single order — one of: `NORMAL` |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK. — one of: `QTY`, `AMOUNT` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. — one of: `ALL`, `CORE`, `NIGHT` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market. |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US` |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). Event trading supports BUY and SELL sides only. Equity trading supports BUY, SELL, and SHORT sides. — one of: `BUY`, `SELL`, `SHORT` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. For equity trading support all order types listed below. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order For event trading: - Supported order types: LIMIT. — one of: `MARKET`, `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK. U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC. - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days). - IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled. - GTD: order that will automatically expire and be cancelled at a specific future date and time. - FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled. — one of: `DAY`, `GTC`, `IOC`, `GTD`, `FOK` |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT. Specifies the trigger price at which the stop order becomes active. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. When event_trade_mode is set for an event contract trade, limit_price is not required. |
| `quantity` | string |  | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `PERCENTAGE`, `AMOUNT` |
| `trailing_stop_step` | string |  | Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%). |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |
| `event_trade_mode` | string |  | Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored. - TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price. - TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field. — one of: `TRADE_IN_AMOUNT`, `TRADE_IN_CONTRACT` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `estimated_cost` | string | yes | Estimated capital required for the order. The final capital usage may differ depending on the actual execution and fee settlement. |
| `estimated_transaction_fee` | string | yes | Estimated transaction fee for placing the order, including exchange, clearing, and commission fees. The actual fee may differ based on final execution. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Place

`POST /broker/orders/place`

> Places order for a specific account. Related Event Notifications: Please refer to the Event Push API documentation: Trade Events

| | |
|---|---|
| **SDK** | `brokerfd.PlaceFDOrder` |
| **Reference** | [common-order-place.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-place.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `new_orders` | array<object> | yes | Order Details |

*Nested — `new_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Used to track or reference the order when interacting with the system. |
| `combo_type` | string | yes | Type of order combination. - NORMAL: Indicates a standard single order — one of: `NORMAL` |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK. — one of: `QTY`, `AMOUNT` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. — one of: `ALL`, `CORE`, `NIGHT` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market. |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US` |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). Event trading supports BUY and SELL sides only. Equity trading supports BUY, SELL, and SHORT sides. — one of: `BUY`, `SELL`, `SHORT` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. For equity trading support all order types listed below. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order For event trading: - Supported order types: LIMIT. — one of: `MARKET`, `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK. U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC. - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days). - IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled. - GTD: order that will automatically expire and be cancelled at a specific future date and time. - FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled. — one of: `DAY`, `GTC`, `IOC`, `GTD`, `FOK` |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT. Specifies the trigger price at which the stop order becomes active. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. When event_trade_mode is set for an event contract trade, limit_price is not required. |
| `quantity` | string |  | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `total_cash_amount` | string |  | The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount. |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `PERCENTAGE`, `AMOUNT` |
| `trailing_stop_step` | string |  | Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%). |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |
| `event_trade_mode` | string |  | Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored. - TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price. - TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field. — one of: `TRADE_IN_AMOUNT`, `TRADE_IN_CONTRACT` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Replace

`POST /broker/orders/replace`

> Modifies an existing order. Only the provided fields are updated; all other attributes remain unchanged.

| | |
|---|---|
| **SDK** | `brokerfd.ReplaceFDOrder` |
| **Reference** | [common-order-replace.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-replace.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `modify_orders` | array<object> | yes | Order Details |

*Nested — `modify_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Used to track or reference the order when interacting with the system. |
| `time_in_force` | string |  | Specifies the duration for which the order remains active in the market (Time-In-Force). Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK. U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC. - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days). - IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled. - GTD: order that will automatically expire and be cancelled at a specific future date and time. - FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled. — one of: `DAY`, `GTC`, `IOC`, `GTD`, `FOK` |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT. Specifies the trigger price at which the stop order becomes active. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `quantity` | string |  | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `order_type` | string |  | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. For equity trading support all order types listed below. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order For event trading: - Supported order types: LIMIT. — one of: `MARKET`, `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS` |
| `trailing_stop_step` | string |  | Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%). |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Cancel

`POST /broker/orders/cancel`

> Cancels a previously submitted order that has not yet been fully filled. Only orders in open status can be cancelled.

| | |
|---|---|
| **SDK** | `brokerfd.CancelFDOrder` |
| **Reference** | [common-order-cancel.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-cancel.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Used to track or reference the order when interacting with the system. |
| `account_id` | string | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Open Orders

`GET /broker/orders/open-orders/list`

> Retrieves pending orders by page. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDOpenOrders` |
| **Reference** | [order-open.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/order-open.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `combo_type` | string | yes | Type of order combination. - NORMAL: Indicates a standard single order |
| `orders` | array<object> | yes | Order Details |

*Nested — `orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `order_id` | string | yes | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market. |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). Event trading supports BUY and SELL sides only. Equity trading supports BUY, SELL, and SHORT sides. — one of: `BUY`, `SELL`, `SHORT` |
| `status` | string | yes | - PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion - SUBMITTED: Indicates that the order has been submitted to the exchange or webull - CANCELLED: Indicates that the order has been successfully cancelled - FILLED: Indicates that the order has been fully executed - FAILED: Indicates a failed order, such as REJECTED - PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed — one of: `PENDING`, `SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, `PARTIAL_FILLED` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. For equity trading support all order types listed below. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order For event trading: - Supported order types: LIMIT. — one of: `MARKET`, `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS` |
| `instrument_type` | string |  | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. — one of: `ALL`, `CORE`, `NIGHT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK. — one of: `QTY`, `AMOUNT` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK. U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC. - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days). - IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled. - GTD: order that will automatically expire and be cancelled at a specific future date and time. - FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled. — one of: `DAY`, `GTC`, `IOC`, `GTD`, `FOK` |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `total_cash_amount` | string |  | The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount. |
| `total_quantity` | string | yes | Total order quantity. Represents the total number of units submitted for this order. |
| `filled_quantity` | string |  | Quantity that has been executed. Represents the number of units that have been filled so far. |
| `filled_price` | string |  | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null. |
| `limit_price` | string |  | Limit Price |
| `stop_price` | string |  | Stop Price |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `PERCENTAGE`, `AMOUNT` |
| `trailing_stop_step` | string |  | Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%). |
| `place_time` | string |  | Order placement time in milliseconds since Unix epoch. |
| `place_time_at` | string | yes | Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `filled_time` | string |  | Time of the last executed trade in milliseconds since Unix epoch. |
| `filled_time_at` | string |  | Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |
| `event_trade_mode` | string |  | Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored. - TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price. - TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field. — one of: `TRADE_IN_AMOUNT`, `TRADE_IN_CONTRACT` |
| `commission` | object |  | Commission breakdown |
| `fees` | array<object> |  | Fee breakdown |

*Nested — `commission`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `actual_commission` | string |  | Actual commission collected |
| `receivable_commission` | string |  | Receivable commission |

*Nested — `fees`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Fee type |
| `actual_value` | string |  | Actual fee collected |
| `receivable_value` | string |  | Receivable fee |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Detail

`GET /broker/orders/get`

> Retrieves the specified order details through the order ID or client order ID.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDOrderDetail` |
| **Reference** | [order-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/order-detail.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `client_order_id` | query | String | yes | The last order ID returned from the previous response. Used for cursor-based pagination. Not required for the first page query. |
| `account_id` | query | String | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `combo_type` | string | yes | Type of order combination. - NORMAL: Indicates a standard single order |
| `orders` | array<object> | yes | Order Details |

*Nested — `orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `order_id` | string | yes | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market. |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). Event trading supports BUY and SELL sides only. Equity trading supports BUY, SELL, and SHORT sides. — one of: `BUY`, `SELL`, `SHORT` |
| `status` | string | yes | - PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion - SUBMITTED: Indicates that the order has been submitted to the exchange or webull - CANCELLED: Indicates that the order has been successfully cancelled - FILLED: Indicates that the order has been fully executed - FAILED: Indicates a failed order, such as REJECTED - PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed — one of: `PENDING`, `SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, `PARTIAL_FILLED` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. For equity trading support all order types listed below. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order For event trading: - Supported order types: LIMIT. — one of: `MARKET`, `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS` |
| `instrument_type` | string |  | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. — one of: `ALL`, `CORE`, `NIGHT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK. — one of: `QTY`, `AMOUNT` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK. U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC. - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days). - IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled. - GTD: order that will automatically expire and be cancelled at a specific future date and time. - FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled. — one of: `DAY`, `GTC`, `IOC`, `GTD`, `FOK` |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `total_cash_amount` | string |  | The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount. |
| `total_quantity` | string | yes | Total order quantity. Represents the total number of units submitted for this order. |
| `filled_quantity` | string |  | Quantity that has been executed. Represents the number of units that have been filled so far. |
| `filled_price` | string |  | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null. |
| `limit_price` | string |  | Limit Price |
| `stop_price` | string |  | Stop Price |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `PERCENTAGE`, `AMOUNT` |
| `trailing_stop_step` | string |  | Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%). |
| `place_time` | string |  | Order placement time in milliseconds since Unix epoch. |
| `place_time_at` | string | yes | Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `filled_time` | string |  | Time of the last executed trade in milliseconds since Unix epoch. |
| `filled_time_at` | string |  | Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |
| `event_trade_mode` | string |  | Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored. - TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price. - TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field. — one of: `TRADE_IN_AMOUNT`, `TRADE_IN_CONTRACT` |
| `commission` | object |  | Commission breakdown |
| `fees` | array<object> |  | Fee breakdown |

*Nested — `commission`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `actual_commission` | string |  | Actual commission collected |
| `receivable_commission` | string |  | Receivable commission |

*Nested — `fees`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Fee type |
| `actual_value` | string |  | Actual fee collected |
| `receivable_value` | string |  | Receivable fee |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order History

`GET /broker/orders/historical-orders/list`

> Retrieves historical orders within a specified date range (default: last 7 days). If orders are group orders, they will be returned together, and the number of orders returned on one page may exceed the page_size. This endpoint may not return the most recent order data in real time due to processing delays. To ensure you get the latest order status, please query the Order Detail endpoint by client_order_id.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDOrderHistory` |
| **Reference** | [order-history.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/order-history.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |
| `start_date` | query | String |  | The start date of the query period. If not provided, the default query period is the last 7 days. Users can specify an earlier date, but the maximum allowed look-back period is 2 years. Format: yyyy-MM-dd. |
| `end_date` | query | String |  | The end date of the query period. If not provided, the default query period is the last 7 days. Users can specify an earlier date, but the maximum allowed look-back period is 2 years. Format: yyyy-MM-dd. |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `combo_type` | string | yes | Type of order combination. - NORMAL: Indicates a standard single order |
| `orders` | array<object> | yes | Order Details |

*Nested — `orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `order_id` | string | yes | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market. |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). Event trading supports BUY and SELL sides only. Equity trading supports BUY, SELL, and SHORT sides. — one of: `BUY`, `SELL`, `SHORT` |
| `status` | string | yes | - PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion - SUBMITTED: Indicates that the order has been submitted to the exchange or webull - CANCELLED: Indicates that the order has been successfully cancelled - FILLED: Indicates that the order has been fully executed - FAILED: Indicates a failed order, such as REJECTED - PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed — one of: `PENDING`, `SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, `PARTIAL_FILLED` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. For equity trading support all order types listed below. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order For event trading: - Supported order types: LIMIT. — one of: `MARKET`, `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS` |
| `instrument_type` | string |  | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. — one of: `ALL`, `CORE`, `NIGHT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount. Supported for U.S. stock trading and Event Contract trading. When placing Event Contract orders using AMOUNT, only BUY orders are supported (side=BUY), and time_in_force must be FOK. — one of: `QTY`, `AMOUNT` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). Event trading supports the following Time in Force (TIF) values: DAY, GTC, IOC, GTD, and FOK. U.S. Equity trading support the following Time in Force (TIF) values: DAY and GTC. - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 60 days). - IOC: Immediate-Or-Cancel, the order attempts to execute immediately. Any portion that can be filled right away will be executed; any unfilled remainder is immediately cancelled. - GTD: order that will automatically expire and be cancelled at a specific future date and time. - FOK: Fill or Kill. The order must be filled in its entirety immediately; otherwise, the entire order will be canceled. — one of: `DAY`, `GTC`, `IOC`, `GTD`, `FOK` |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `total_cash_amount` | string |  | The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount. |
| `total_quantity` | string | yes | Total order quantity. Represents the total number of units submitted for this order. |
| `filled_quantity` | string |  | Quantity that has been executed. Represents the number of units that have been filled so far. |
| `filled_price` | string |  | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null. |
| `limit_price` | string |  | Limit Price |
| `stop_price` | string |  | Stop Price |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `PERCENTAGE`, `AMOUNT` |
| `trailing_stop_step` | string |  | Trailing spread. When trailing_type is PERCENTAGE, the value must be greater than or equal to 0.01 and cannot exceed 1.0, and it represents a percentage in decimal form (e.g., 1.0 = 100%, 0.1 = 10%, 0.01 = 1%). |
| `place_time` | string |  | Order placement time in milliseconds since Unix epoch. |
| `place_time_at` | string | yes | Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `filled_time` | string |  | Time of the last executed trade in milliseconds since Unix epoch. |
| `filled_time_at` | string |  | Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `event_outcome` | string |  | Event outcome decision, only applicable to event orders. — one of: `yes`, `no` |
| `event_trade_mode` | string |  | Specifies how the order quantity is expressed for event contract trading. Only applicable to event orders. When this field is set, the order executes at the best available market price; the limit_price field is ignored. - TRADE_IN_AMOUNT: Specifies how the order quantity is expressed for event contract trading. When this field is set, the order executes at the best available market price. - TRADE_IN_CONTRACT: The order is specified by the number of contracts the user wants to buy or sell at the best available market price. Requires quantity field. — one of: `TRADE_IN_AMOUNT`, `TRADE_IN_CONTRACT` |
| `commission` | object |  | Commission breakdown |
| `fees` | array<object> |  | Fee breakdown |

*Nested — `commission`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `actual_commission` | string |  | Actual commission collected |
| `receivable_commission` | string |  | Receivable commission |

*Nested — `fees`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Fee type |
| `actual_value` | string |  | Actual fee collected |
| `receivable_value` | string |  | Receivable fee |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Create Cash Journal

`POST /broker/journals/cash-journals/create`

> Creates a cash journal request. Cash journal supports both proprietary account (PAB) level and end client account level. Journals are supported only between accounts with the same owner. Related Event Notifications: Please refer to the Event Push API documentation: Cash Journal Events

| | |
|---|---|
| **SDK** | `—` |
| **Reference** | [broker-journal-cash-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-journal-cash-create.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `from_account` | string | yes | From account id. |
| `to_account` | string | yes | To account id. |
| `amount` | string | yes | Amount |
| `currency` | string | yes | Currency — one of: `USD` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. |
| `from_account` | string | yes | From account id. |
| `to_account` | string | yes | To account id. |
| `journal_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Cash Journal Detail

`GET /broker/journals/cash-journals/get`

> Retrieves details information for cash journal request.

| | |
|---|---|
| **SDK** | `brokerfd.GetFDCashJournalDetail` |
| **Reference** | [broker-journal-cash-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-journal-cash-query.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `client_request_id` | query | String | yes | Client Request ID, unique for each request. Maximum length is 32 characters. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). |
| `account_id` | query | String | yes | From account identifier. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_request_id` | string | yes | Client Request ID, unique for each request. |
| `from_account` | string | yes | From account id. |
| `to_account` | string | yes | To account id. |
| `journal_id` | string | yes | Request id generated by the system. |
| `amount` | string | yes | Amount |
| `currency` | string | yes | Currency — one of: `USD` |
| `status` | string | yes | Request status — one of: `SUBMITTED`, `CANCELED`, `REJECTED`, `FAILED`, `COMPLETED` |
| `reason` | string |  | Reason. If the terminal state is not “COMPLETED”, return the reason. |
| `create_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |
| `update_time` | string | yes | UTC time, format: yyyy-MM-dd'T'HH:mm:ss.SSSZ. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List Enums

`GET /broker/master-data/enums/list`

> Get master data by data type and optional parent code. Master Data definition, hierarchy, and code description: \| data_type \| description \| \|-----------\|-------------\| \| ID_DOC_TYPE \| Identification document type \| \| MARITAL_STATUS \| Marital status options \| \| EMPLOYMENT_TYPE \| Type of employment \| \| OCCUPATION \| Specific occupation under the employment type \| \| POSITION \| Job position or title \| \| LIQUIDITY_NEEDS \| Liquidity needs options \| \| ANNUAL_INCOME \| Annual income range \| \| TOTAL_NET_WORTH \| Estimated total net worth range \| \| LIQUID_NET_WORTH \| Liquid net worth range \| \| INVESTMENT_KNOWLEDGE \| Investment knowledge level (old) \| \| INVESTMENT_OBJECTIVE \| Investment goal or objective \| \| TIME_HORIZON \| Investment time horizon \| \| INVESTMENT_EXPERIENCE \| Experience level in investments \| \| INVESTMENT_KNOWLEDGE_V2 \| Investment knowledge level \| \| TRADE_PER_YEAR \| Estimated number of trades per year \| \| TAX_TYPE \| Tax classification type \| \| AGREEMENT_TYPE \| Agreement type \| \| AGREEMENT_CONTENT_TYPE \| Agreement content type \| \| AGREEMENT_SIGN_METHOD \| Agreement signing method \|

| | |
|---|---|
| **SDK** | `brokerfd.GetFDEnums` |
| **Reference** | [list-enums.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-enums.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `data_type` | query | string | yes | Type of master data — one of: `ID_DOC_TYPE`, `MARITAL_STATUS`, `EMPLOYMENT_TYPE`, `OCCUPATION`, `POSITION`, `LIQUIDITY_NEEDS`, `ANNUAL_INCOME`, `TOTAL_NET_WORTH`, `LIQUID_NET_WORTH`, `INVESTMENT_KNOWLEDGE_OLD`, `INVESTMENT_OBJECTIVE`, `TIME_HORIZON`, `INVESTMENT_EXPERIENCE`, `INVESTMENT_KNOWLEDGE`, `TRADE_PER_YEAR`, `TAX_TYPE`, `AGREEMENT_CONTENT_TYPE`, `AGREEMENT_SIGN_METHOD`, `COUNTRY`, `STATE`, `CITY`, `CONTACT_RELATIONSHIP` |
| `parent_code` | query | string |  | Parent code of the master data, used for hierarchical queries |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `code` | string |  | Business code of the master data item, typically used as a stable reference value. |
| `name` | string |  | Display name of the master data item in the default language. |
| `parent_code` | string |  | Parent code of the master data item, used to represent hierarchical relationships. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Trade Calendar

`GET /broker/master-data/trading-calendars/list`

> Retrieves trading and settlement calendar for the specified country and market. Related Event Notifications: Please refer to the Event Push API documentation: Trade Calendar Events

| | |
|---|---|
| **SDK** | `brokerfd.GetFDTradeCalendar` |
| **Reference** | [list-trade-calendar.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-trade-calendar.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `market` | query | string | yes | Country code, US etc. — one of: `US` |
| `instrument_type` | query | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `EVENT` |
| `year` | query | String | yes | Query Year in YYYY format |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `date` | string | yes | Date string in YYYY-MM-DD format. |
| `settlement_date` | string |  | Settlement Date string in YYYY-MM-DD format. |
| `is_trading_day` | boolean | yes | Is Trading Day |
| `is_early_close` | boolean |  | Whether it is a half-day trading session |
| `is_settlement_day` | boolean | yes | Is Settlement Day |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## List Agreements

`GET /broker/agreements/list`

> Helps retrieve all available agreement templates of a specific type. These templates can be presented to users for signing during account setup or feature activation. Note: This endpoint returns agreement templates, not user-specific agreement status. To check if a user has signed an agreement, use the account-specific endpoints..

| | |
|---|---|
| **SDK** | `brokerfd.ListAgreements` |
| **Reference** | [broker-list-agreements-by-type.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-list-agreements-by-type.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `lang` | query | String |  | Language parameter, default is EN EN: English |
| `agreement_type` | query | String | yes | Agreement type parameter. Value must be obtained from the Master Data API with data_type = AGREEMENT_TYPE. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | string | yes | Agreement ID |
| `version` | string | yes | Agreement Version |
| `name` | string | yes | Agreement Name |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Agreement Details

`GET /broker/agreements/get`

> Retrieves detailed information about a specific agreement.

| | |
|---|---|
| **SDK** | `brokerfd.GetAgreementDetail` |
| **Reference** | [broker-get-agreement-details.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-get-agreement-details.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `ids` | query | String | yes | The unique identifier of the agreement to retrieve details for. |
| `lang` | query | String | yes | Language parameter EN: English |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | string | yes | Agreement ID |
| `name` | string | yes | Agreement Name |
| `description` | string |  | Agreement Description |
| `content` | string | yes | Agreement Content |
| `content_type` | string | yes | Agreement content type — one of: `TEXT`, `URL`, `PDF` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).


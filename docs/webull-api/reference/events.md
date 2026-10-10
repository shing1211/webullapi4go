# Events (gRPC) — Verbatim Reference

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> Server-streaming gRPC subscriptions, signed with HMAC-SHA256 over the serialized request.

> Verbatim snapshot of Webull's published OpenAPI definitions. No SDK-specific content.

[<- Master Reference](../master-reference.md) · [<- Webull API Reference](../../webull-api.md)

## Subscribe Trade Events

> Source: <https://developer.webull.hk/apis/docs/reference/custom/subscribe-trade-events.md>

### Subscribe Trade Events

##### Interface Description

Trade events subscription is a `server streaming` persistent connection implemented based on `gRPC`, which is suitable for connecting Webull customers through the OpenAPI development platform. The trade events subscription fully follows the `gRPC` open source protocol, and you can refer to the [gRPC open source](https://grpc.io/docs/) library when using it.

Currently, the interface supports order status change message push, and the supported scenarios are as follows:

| scene_type     | Description               |
|----------------|---------------------------|
| FILLED         | Partially filled          |
| FINAL_FILLED   | All filled                |
| PLACE_FAILED   | Order failed              |
| MODIFY_SUCCESS | Change order successfully |
| MODIFY_FAILED  | Change order failed       |
| CANCEL_SUCCESS | Cancellation succeeded    |
| CANCEL_FAILED  | Cancellation failed       |

##### Trade events subscribe Proto protocol definition.

**Request Proto**

```protobuf
message SubscribeRequest {
	 uint32 subscribeType = 1; // Subscription type
	 int64 timestamp = 2; // Timestamp
	 string contentType = 3; // Content type
	 string payload = 4; // Content
	 repeated string accounts = 5; // Account ID
}
```

**Response Proto**

```protobuf
message SubscribeResponse {
	 EventType eventType = 1; // Event type
	 uint32 subscribeType = 2; // Subscription type
	 string contentType = 3; // Subscription type
	 string payload = 4; // Content
	 string requestId = 5; // Request id
	 int64  timestamp = 6; // Timestamp
}
```

**EventType enumeration**

```protobuf
enum EventType {
	SubscribeSuccess = 0; // Subscription succeeded
	Ping = 1; // Heartbeat information
	AuthError = 2; // Authentication error
	NumOfConnExceed = 3; // Connection limit exceeded
	SubscribeExpired = 4; // Subscription expired
}
```

##### Request Example

<Tabs groupId="programming-language">

When using sdk request, subscribeType, timestamp, contentType, and payload can be ignored. Just pass in the accounts. subscribeType currently only supports =1.
In the following case, the _on_log method is used to output the log. The my_on_events_message method is to receive order status change messages.

```python
import logging

from webull.trade.events.types import ORDER_STATUS_CHANGED, EVENT_TYPE_ORDER
from webull.trade.trade_events_client import TradeEventsClient

your_app_key = "<your_app_key>"
your_app_secret = "<your_app_secret>"
account_id = "<your_account_id>"
region_id = "hk"

# PRD env host: events-api.webull.hk
# Test env host: events-api.sandbox.webull.hk
optional_api_endpoint = "<event_api_endpoint>"

def _on_log(level, log_content):
    print(logging.getLevelName(level), log_content)

def my_on_events_message(event_type, subscribe_type, payload, raw_message):
    if EVENT_TYPE_ORDER == event_type and ORDER_STATUS_CHANGED == subscribe_type:
        print('%s' % payload)

if __name__ == '__main__':

    # Create EventsClient instance
    trade_events_client = TradeEventsClient(your_app_key, your_app_secret, region_id)
    # For non production environment, you need to set the domain name of the subscription service through eventsclient. For example, the domain name of the UAT environment is set here
    # trade_events_client = TradeEventsClient(your_app_key, your_app_secret, region_id, host=optional_api_endpoint)
    trade_events_client.on_log = _on_log

    # Set the callback function when the event data is received.
    # The data of order status change is printed here

    trade_events_client.on_events_message = my_on_events_message
    # Set the account ID to be subscribed and initiate the subscription. This method is synchronous
    trade_events_client.do_subscribe([account_id])
```

The handleEventMessage method is to receive order status change messages.

```java
import com.google.gson.reflect.TypeToken;
import com.webull.openapi.core.execption.ClientException;
import com.webull.openapi.core.execption.ServerException;
import com.webull.openapi.core.logger.Logger;
import com.webull.openapi.core.logger.LoggerFactory;
import com.webull.openapi.core.serialize.JsonSerializer;
import com.webull.openapi.samples.config.Env;
import com.webull.openapi.trade.events.subscribe.ISubscription;
import com.webull.openapi.trade.events.subscribe.ITradeEventClient;
import com.webull.openapi.trade.events.subscribe.message.EventType;
import com.webull.openapi.trade.events.subscribe.message.SubscribeRequest;
import com.webull.openapi.trade.events.subscribe.message.SubscribeResponse;

import java.util.Map;

public class TradeEventsClient {

    private static final Logger logger = LoggerFactory.getLogger(TradeEventsClient.class);

    public static void main(String[] args) {
        try (ITradeEventClient client = ITradeEventClient.builder()
                .appKey(Env.APP_KEY)
                .appSecret(Env.APP_SECRET)
                .regionId(Env.REGION_ID)
                // .host("<event_api_endpoint>")
                .onMessage(TradeEventsClient::handleEventMessage)
                .build()) {

            SubscribeRequest request = new SubscribeRequest("<your_account_id>");

            ISubscription subscription = client.subscribe(request);
            subscription.blockingAwait();

        } catch (ClientException ex) {
            logger.error("Client error", ex);
        } catch (ServerException ex) {
            logger.error("Sever error", ex);
        } catch (Exception ex) {
            logger.error("Unknown error", ex);
        }
    }

    private static void handleEventMessage(SubscribeResponse response) {
        if (SubscribeResponse.CONTENT_TYPE_JSON.equals(response.getContentType())) {
            Map<String, String> payload = JsonSerializer.fromJson(response.getPayload(),
                    new TypeToken<Map<String, String>>(){}.getType());
            if (EventType.Order.getCode() == response.getEventType() || EventType.Position.getCode() == response.getEventType()) {
                logger.info("{}", payload);
            }
        }
    }
}
```

##### Response Example
Transaction event scene type

<Tabs groupId="programming-language">

```python
{
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVUOVRA8BV0KHKN60000000",
    "order_id": "036LVUOVRA8BV0KHKN60000000",
    "client_order_id": "6adeba36fc174acd92538e990c06274e",
    "instrument_id": "913256135",
    "order_status": "PARTIAL_FILLED",
    "symbol": "AAPL",
    "short_name": "Apple Inc",
    "qty": "3.0000000000",
    "filled_price": "10.0",
    "filled_qty": "1.0",
    "filled_time": "2025-11-26T11:40:35.524+0000",
    "side": "BUY",
    "scene_type": "FILLED",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "actual_commission": "0.01",
    "receivable_commission": "0.01",
    "fees": [
        {
            "type": "FINRA_CAT_REGULATORY_FEE",
            "actual_value": "0.01",
            "receivable_value": "0.01"
        }
    ]
}
DEBUG response:eventType: Ping
subscribeType: 1
contentType: "text/plain"
requestId: "5ca8b9d7-6f57-4218-9fb3-d3951c2fa399"
timestamp: 1764157239726
```
  

```python
{
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVUAB7C8BV0KHKN60000000",
    "order_id": "036LVUAB7C8BV0KHKN60000000",
    "client_order_id": "9b0440781feb4c0cb8523224ce0926ac",
    "instrument_id": "913256135",
    "order_status": "FILLED",
    "symbol": "AAPL",
    "short_name": "Apple Inc",
    "qty": "2.0000000000",
    "filled_price": "277.98",
    "filled_qty": "2.0",
    "filled_time": "2025-11-26T11:35:38.513+0000",
    "side": "BUY",
    "scene_type": "FINAL_FILLED",
    "category": "US_STOCK",
    "order_type": "MARKET",
    "actual_commission": "0.01",
    "receivable_commission": "0.01",
    "fees": [
        {
            "type": "FINRA_CAT_REGULATORY_FEE",
            "actual_value": "0.01",
            "receivable_value": "0.01"
        }
    ]
}
DEBUG response:eventType: Ping
subscribeType: 1
contentType: "text/plain"
requestId: "5ca8b9d7-6f57-4218-9fb3-d3951c2fa399"
timestamp: 1764156939725
```

```python
{
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV3TME8BV0KHKN60000000",
    "order_id": "036LVV3TME8BV0KHKN60000000",
    "client_order_id": "794cabc6783640d7869a37f966ae4819",
    "instrument_id": "913256135",
    "order_status": "FAILED",
    "symbol": "AAPL",
    "short_name": "Apple Inc",
    "qty": "92.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "scene_type": "PLACE_FAILED",
    "category": "US_STOCK",
    "order_type": "LIMIT"
}
DEBUG response:eventType: Ping
subscribeType: 1
contentType: "text/plain"
requestId: "5ca8b9d7-6f57-4218-9fb3-d3951c2fa399"
timestamp: 1764157359725
```
  

```python
{
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV5P4I8BV0KHKN60000000",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "04cda8db7ed940f6afeb26be6201ee53",
    "instrument_id": "913256135",
    "order_status": "SUBMITTED",
    "symbol": "AAPL",
    "short_name": "Apple Inc",
    "qty": "4.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "scene_type": "MODIFY_SUCCESS",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "actual_commission": "0.01",
    "receivable_commission": "0.01",
    "fees": [
        {
            "type": "FINRA_CAT_REGULATORY_FEE",
            "actual_value": "0.01",
            "receivable_value": "0.01"
        }
    ]
}
DEBUG response:eventType: Ping
subscribeType: 1
contentType: "text/plain"
requestId: "5ca8b9d7-6f57-4218-9fb3-d3951c2fa399"
timestamp: 1764157419725
```
  

```python
{
    "account_id": "PHIUK08VAKH7EOVG85ULCAG3JB",
    "request_id": "036LVV5P4I8BV0KHKN60000000",
    "order_id": "036LVV5P4I8BV0KHKN60000000",
    "client_order_id": "04cda8db7ed940f6afeb26be6201ee53",
    "instrument_id": "913256135",
    "order_status": "CANCELLED",
    "symbol": "AAPL",
    "short_name": "Apple Inc",
    "qty": "4.0000000000",
    "filled_price": "0E-10",
    "filled_qty": "0E-10",
    "side": "BUY",
    "scene_type": "CANCEL_SUCCESS",
    "category": "US_STOCK",
    "order_type": "LIMIT",
    "actual_commission": "0.01",
    "receivable_commission": "0.01",
    "fees": [
        {
            "type": "FINRA_CAT_REGULATORY_FEE",
            "actual_value": "0.01",
            "receivable_value": "0.01"
        }
    ]
}
DEBUG response:eventType: Ping
subscribeType: 1
contentType: "text/plain"
requestId: "5ca8b9d7-6f57-4218-9fb3-d3951c2fa399"
timestamp: 1764157479725
```

## Subscribe Position Events

> Source: <https://developer.webull.com/apis/docs/reference/custom/subscribe-position-events.md>

### Subscribe Position Events

##### Interface Description

The interface supports event contract position settlement message push, subscribeType only supports = 2. Event settlement messages require the latest SDK version.

##### Position events subscribe Proto protocol definition.

**Request Proto**

```protobuf
message SubscribeRequest {
	 uint32 subscribeType = 1; // Subscription type
	 int64 timestamp = 2; // Timestamp
	 string contentType = 3; // Content type
	 string payload = 4; // Content
	 repeated string accounts = 5; // Account ID
}
```

**Response Proto**

```protobuf
message SubscribeResponse {
	 EventType eventType = 1; // Event type
	 uint32 subscribeType = 2; // Subscription type
	 string contentType = 3; // Subscription type
	 string payload = 4; // Content
	 string requestId = 5; // Request id
	 int64  timestamp = 6; // Timestamp
}
```

**EventType enumeration**

```protobuf
enum EventType {
	SubscribeSuccess = 0; // Subscription succeeded
	Ping = 1; // Heartbeat information
	AuthError = 2; // Authentication error
	NumOfConnExceed = 3; // Connection limit exceeded
	SubscribeExpired = 4; // Subscription expired
}
```

##### Request Example

<Tabs groupId="programming-language">

In the following case, the _on_log method is used to output the log. The my_on_events_message method is to receive order status change messages.

```python
import logging

from webull.trade.events.types import EVENT_TYPE_POSITION, POSITION_STATUS_CHANGED
from webull.trade.trade_events_client import TradeEventsClient

your_app_key = "<your_app_key>"
your_app_secret = "<your_app_secret>"
account_id = "<your_account_id>"
region_id = "us"

# PRD env host: events-api.webull.com
# Test env host: us-openapi-alb.uat.webullbroker.com
optional_api_endpoint = "<event_api_endpoint>"

def _on_log(level, log_content):
    print(logging.getLevelName(level), log_content)

def my_on_events_message(event_type, subscribe_type, payload, raw_message):
    if EVENT_TYPE_POSITION == event_type and POSITION_STATUS_CHANGED == subscribe_type:
        print('event payload:%s' % payload)

if __name__ == '__main__':

    # Create EventsClient instance
    trade_events_client = TradeEventsClient(your_app_key, your_app_secret, region_id)
    # For non production environment, you need to set the domain name of the subscription service through eventsclient. For example, the domain name of the UAT environment is set here
    # trade_events_client = TradeEventsClient(your_app_key, your_app_secret, region_id, host=optional_api_endpoint)
    trade_events_client.on_log = _on_log

    # Set the callback function when the event data is received.
    # The data of order status change is printed here

    trade_events_client.on_events_message = my_on_events_message
    # Set the account ID to be subscribed and initiate the subscription. This method is synchronous
    trade_events_client.do_subscribe([account_id])
```

The handleEventMessage method is to receive order status change messages.

```java
import com.google.gson.reflect.TypeToken;
import com.webull.openapi.core.execption.ClientException;
import com.webull.openapi.core.execption.ServerException;
import com.webull.openapi.core.logger.Logger;
import com.webull.openapi.core.logger.LoggerFactory;
import com.webull.openapi.core.serialize.JsonSerializer;
import com.webull.openapi.samples.config.Env;
import com.webull.openapi.trade.events.subscribe.ISubscription;
import com.webull.openapi.trade.events.subscribe.ITradeEventClient;
import com.webull.openapi.trade.events.subscribe.message.EventType;
import com.webull.openapi.trade.events.subscribe.message.SubscribeRequest;
import com.webull.openapi.trade.events.subscribe.message.SubscribeResponse;

import java.util.Map;

public class TradeEventsClient {

    private static final Logger logger = LoggerFactory.getLogger(TradeEventsClient.class);

    public static void main(String[] args) {
        try (ITradeEventClient client = ITradeEventClient.builder()
                .appKey(Env.APP_KEY)
                .appSecret(Env.APP_SECRET)
                .regionId(Env.REGION_ID)
                // .host("<event_api_endpoint>")
                .onMessage(TradeEventsClient::handleEventMessage)
                .build()) {

            SubscribeRequest request = new SubscribeRequest("<your_account_id>");

            ISubscription subscription = client.subscribe(request);
            subscription.blockingAwait();

        } catch (ClientException ex) {
            logger.error("Client error", ex);
        } catch (ServerException ex) {
            logger.error("Sever error", ex);
        } catch (Exception ex) {
            logger.error("Unknown error", ex);
        }
    }

    private static void handleEventMessage(SubscribeResponse response) {
        if (SubscribeResponse.CONTENT_TYPE_JSON.equals(response.getContentType())) {
            Map<String, String> payload = JsonSerializer.fromJson(response.getPayload(),
                    new TypeToken<Map<String, String>>(){}.getType());
            if (EventType.Position.getCode() == response.getEventType()) {
                logger.info("{}", payload);
            }
        }
    }
}
```

##### Response Example
Position event scene type

<Tabs groupId="programming-language">

```python
{
    "event_name": "Number of rate cuts in 2025",
    "yes_condition": "Exactly 8 cuts",
    "settle_result": "Yes",
    "settle_side": "Yes",
    "quantity": "40",
    "cost": "20.00",
    "settle_amount": "40.00"
}
DEBUG response:eventType: Ping
subscribeType: 2
contentType: "text/plain"
requestId: "ab39b532-3ad4-46c4-823d-3dff12e0f1b9"
timestamp: 1768994319673
```

## Subscribe Events (FD)

> Source: <https://developer.webull.com/apis/docs/reference/fd-events/subscribe-events.md>

### Event Subscription Guide

##### Interface Description

The Events Subscription service is implemented as a server-streaming persistent connection based on the `gRPC` framework. It is designed for Webull customers
integrating through the OpenAPI development platform to receive real-time
event notifications in a reliable and efficient manner.

This service fully complies with the official `gRPC` open-source protocol
specification. Clients can leverage standard `gRPC` libraries and SDKs
(https://grpc.io/docs/) to establish and maintain the streaming connection.

Through this subscription channel, the server continuously pushes structured
event messages to clients once the connection is established, eliminating
the need for repeated polling and ensuring low-latency data delivery.

##### Event Modules Overview

The message push service is organized into the following functional modules.
Each module can be subscribed to independently based on business needs.

| Module Name                   | Description                                                                 |
|-------------------------------|-----------------------------------------------------------------------------|
| Application Events            | Account application lifecycle events, including submission and status updates. |
| Account Events                | Account-level updates such as status changes and configuration updates.    |
| Funding Events                | Deposit and withdrawal related events and status notifications.            |
| Fees and Credits Events       | Fee charges, and credit postings.                             |
| Journal Events                | Internal transfers and journal entry notifications.                        |
| Master Data Events            |     Trading calendar change notifications.      |
| Instrument Events             | Instrument-level updates such as trading status or metadata changes.       |
| Trade Events                  | Order execution, trade confirmations, and related transaction updates.     |
| Corporate Actions Events      | Corporate action announcements and processing updates.                     |
| Event Contract Events         | Event contract lifecycle and status updates.                               |
| SOD Events                    | Start-of-Day processing results and daily initialization updates.          |

Clients may selectively subscribe to one or more modules depending on their
integration scope and operational requirements.

##### Events subscribe Proto protocol definition.

**Complete Description of Event gRPC/Protobuf Protocol**

```protobuf
syntax = "proto3";

package grpc.event;
option go_package = ".;grpc_event";

service EventService {
	rpc Subscribe(SubscribeRequest) returns (stream SubscribeResponse) {}
}

message SubscribeRequest {
	 repeated SubscribeEvent events = 1;
	 string timestamp = 2;
}

message SubscribeResponse {
	 Type type = 1;
	 string id = 2;
	 string eventType = 3;
	 string position = 4;
	 string  timestamp = 5; 
	 string version = 6;
	 string payload = 7;
}

enum Type {
	SubscribeSuccess = 0;
	Ping = 1;
	Event = 2;
	AuthError = 3;
	NumOfConnExceed = 4;
	SubscribeExpired = 5;
}

message SubscribeEvent {
	string eventType = 1;
	string position = 2;
}
```

**Request Proto**

```protobuf
message SubscribeRequest {
	 repeated SubscribeEvent events = 1;
	 string timestamp = 2;
}
message SubscribeEvent {
	string eventType = 1;
	string position = 2;
}
```
**Response Proto**

```protobuf
message SubscribeResponse {
	 Type type = 1;
	 string id = 2;
	 string eventType = 3;
	 string position = 4;
	 string  timestamp = 5; 
	 string version = 6;
	 string payload = 7;
}

```

**Type enumeration**

```protobuf
enum Type {
	SubscribeSuccess = 0;
	Ping = 1;
	Event = 2;
	AuthError = 3;
	NumOfConnExceed = 4;
	SubscribeExpired = 5;
}
```

##### Request Example

For native Python integration, please use the sample code included in the ZIP package: [python_event_subscription_sample.zip](https://uat-static.webullbroker.com/inst-bo/O3GILMTVCGQQ77CQDCRA8J9VO9.zip)

##### Response Payload Example
 Event scene type

<Tabs groupId="programming-language">
  

```python
{
    "id": "event_c4b2c210-ce32-41d4-a9a1-cfad4fdf191c",
    "event_type": "APPLICATION",
    "position": "CJO1fxACGAAgADAB",
    "timestamp": "2025-03-29T07:02:33.200962333Z",
    "payload": {
        "client_request_id": "ndjanww3004abe894b51244560kmdcd",
        "application_id": "d44caf9eb12d4abe894b529699059ec5",
        "biz_type": "ACCOUNT_CREATE",
        "status": "APPROVED",
        "reason": "",
        "timestamp": "2025-11-21T06:27:43.312Z"
    }
}
```


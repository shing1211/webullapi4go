// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trade

import (
	"encoding/json"
	"testing"
)

// BatchPlaceOrder's page documents total, success, failed and batch_orders, and marks
// all four required, while the SDK's type carried a single Results field that appears
// on no page. A caller could not learn how much of a batch was placed, and a rejected
// order inside a batch was indistinguishable from a success that returned no
// identifier, because the documented element's error_code and message reached no field.
//
// The addition is safe for the same reason every required-name addition is: a
// conforming server always sends these, so a caller reading them today reads zero and
// reads the real value after, and no call that works today can break.

func TestBatchPlaceOrderResponseReadsTheDocumentedCounts(t *testing.T) {
	const body = `{
		"total": 3,
		"success": 2,
		"failed": 1,
		"batch_orders": [
			{"client_order_id": "C1", "order_id": "O1"},
			{"client_order_id": "C2", "order_id": "O2"},
			{"client_order_id": "C3", "error_code": "OPENAPI_NO_TRADING_TIME", "message": "Non-trading time."}
		]
	}`
	var got BatchPlaceOrderResponse
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Total != 3 {
		t.Errorf("Total = %d, want 3: the documented name reached no field", got.Total)
	}
	if got.Success != 2 {
		t.Errorf("Success = %d, want 2", got.Success)
	}
	if got.Failed != 1 {
		t.Errorf("Failed = %d, want 1", got.Failed)
	}
	if len(got.BatchOrders) != 3 {
		t.Fatalf("BatchOrders has %d entries, want 3", len(got.BatchOrders))
	}
	// The rejected order is the point of the element's error fields: without them a
	// partial failure looked like a success.
	rejected := got.BatchOrders[2]
	if rejected.ErrorCode != "OPENAPI_NO_TRADING_TIME" {
		t.Errorf("ErrorCode = %q, want OPENAPI_NO_TRADING_TIME", rejected.ErrorCode)
	}
	if rejected.Message != "Non-trading time." {
		t.Errorf("Message = %q, want Non-trading time.", rejected.Message)
	}
	if got.BatchOrders[0].OrderID != "O1" {
		t.Errorf("BatchOrders[0].OrderID = %q, want O1", got.BatchOrders[0].OrderID)
	}
}

// TestBatchPlaceOrderResponseStillDecodesTheSDKSpelling keeps the change additive.
// Results appears on no documented page but removing it would break callers, so a
// body carrying it must still decode.
func TestBatchPlaceOrderResponseStillDecodesTheSDKSpelling(t *testing.T) {
	const body = `{"results": [{"client_order_id": "C1", "order_id": "O1"}]}`
	var got BatchPlaceOrderResponse
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Results) != 1 || got.Results[0].OrderID != "O1" {
		t.Errorf("Results = %+v, want the SDK's own spelling to still decode", got.Results)
	}
	// The documented counts are absent from this body, so they must stay zero rather
	// than borrowing the undocumented field.
	if got.Total != 0 || got.Success != 0 || got.Failed != 0 || len(got.BatchOrders) != 0 {
		t.Errorf("documented fields were populated from an undocumented body: %+v", got)
	}
}

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

package trade_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/shing1211/webullapi4go/pkg/domain/order"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

func TestTrackedOrderRegistryIsolatesAccounts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"shared-id","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	for _, accountID := range []string{"ACCOUNT-A", "ACCOUNT-B"} {
		req := validPlaceRequest()
		req.AccountID = accountID
		req.NewOrders[0].ClientOrderID = "shared-id"
		if _, err := c.PlaceOrder(context.Background(), req); err != nil {
			t.Fatalf("PlaceOrder(%q) error = %v", accountID, err)
		}
	}

	a, ok := c.GetTrackedOrder("ACCOUNT-A", "shared-id")
	if !ok {
		t.Fatal("GetTrackedOrder(ACCOUNT-A) = not found")
	}
	b, ok := c.GetTrackedOrder("ACCOUNT-B", "shared-id")
	if !ok {
		t.Fatal("GetTrackedOrder(ACCOUNT-B) = not found")
	}
	if a == b {
		t.Fatal("orders for different accounts share the same tracked object")
	}
	if a.AccountID != "ACCOUNT-A" || b.AccountID != "ACCOUNT-B" {
		t.Fatalf("account IDs = %q/%q", a.AccountID, b.AccountID)
	}
	if _, ok := c.GetTrackedOrder("ACCOUNT-C", "shared-id"); ok {
		t.Fatal("GetTrackedOrder(ACCOUNT-C) unexpectedly found an order")
	}
	if got, ok := c.TrackedOrderState("ACCOUNT-A", "shared-id"); !ok || got != order.StatePending {
		t.Fatalf("TrackedOrderState(ACCOUNT-A) = (%q, %v), want (PENDING, true)", got, ok)
	}
}

func TestOrderGuardrailKeepsCompatibilityCode(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("guardrail request reached server: %s", r.URL.Path)
	}))
	defer srv.Close()

	c := newTradeClient(t, srv.URL, trade.WithMaxOrderQuantity("1"))
	req := validPlaceRequest()
	req.NewOrders[0].Quantity = mp("2")
	_, err := c.PlaceOrder(context.Background(), req)
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PlaceOrder() error = %v, want invalid_config", err)
	}
	if !errors.Is(err, errs.ErrOrderGuardrail) {
		t.Fatalf("PlaceOrder() error = %v, want ErrOrderGuardrail in the chain", err)
	}
}

func TestPlaceAndBatchResultsAreRegistered(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/trading/orders/place", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"place-id","order_id":"OID-PLACE"}`))
	})
	mux.HandleFunc("/trading/orders/batch-place", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"client_order_id":"batch-1","order_id":"OID-B1"},{"client_order_id":"batch-2","order_id":"OID-B2"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	placeReq := validPlaceRequest()
	placeReq.NewOrders[0].ClientOrderID = "place-id"
	placed, err := c.PlaceOrder(context.Background(), placeReq)
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	trackedPlace, ok := c.GetTrackedOrder(placeReq.AccountID, "place-id")
	if !ok || trackedPlace != placed {
		t.Fatalf("placed order was not registered: got %p, want %p (found %v)", trackedPlace, placed, ok)
	}

	batchReq := placeReq
	batchReq.NewOrders = append([]trade.OrderRequest(nil), placeReq.NewOrders...)
	batchReq.NewOrders[0].ClientOrderID = "batch-1"
	second := validOrder()
	second.ClientOrderID = "batch-2"
	batchReq.NewOrders = append(batchReq.NewOrders, second)
	batchResponse, err := c.BatchPlaceOrder(context.Background(), batchReq)
	if err != nil {
		t.Fatalf("BatchPlaceOrder() error = %v", err)
	}
	if len(batchResponse.Results) != 2 {
		t.Fatalf("len(BatchPlaceOrder().Results) = %d, want 2", len(batchResponse.Results))
	}
	for _, clientOrderID := range []string{"batch-1", "batch-2"} {
		tracked, ok := c.GetTrackedOrder(batchReq.AccountID, clientOrderID)
		if !ok {
			t.Fatalf("batch order %q was not registered", clientOrderID)
		}
		if tracked.State() != order.StatePending {
			t.Fatalf("batch order %q state = %q, want PENDING", clientOrderID, tracked.State())
		}
	}
}

func TestApplyEventAndStatusReconciliation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"tracked-id","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	req := validPlaceRequest()
	req.NewOrders[0].ClientOrderID = "tracked-id"
	if _, err := c.PlaceOrder(context.Background(), req); err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}

	got, err := c.ApplyOrderEvent(req.AccountID, "tracked-id", order.EventAcknowledge)
	if err != nil || got != order.StatePreSubmitted {
		t.Fatalf("ApplyOrderEvent(acknowledge) = (%q, %v), want (PRE_SUBMITTED, nil)", got, err)
	}
	got, err = c.ReconcileOrderStatus(req.AccountID, "tracked-id", trade.OrderStatusFilled)
	if err != nil || got != order.StateFilled {
		t.Fatalf("ReconcileOrderStatus(FILLED) = (%q, %v), want (FILLED, nil)", got, err)
	}
	got, err = c.ReconcileOrderStatus(req.AccountID, "tracked-id", order.StateSubmitted)
	if err != nil || got != order.StateFilled {
		t.Fatalf("stale ReconcileOrderStatus(SUBMITTED) = (%q, %v), want (FILLED, nil)", got, err)
	}
	if _, err := c.ApplyOrderEvent("OTHER-ACCOUNT", "tracked-id", order.EventPlace); !errs.Is(err, errs.CodeNotInitialized) {
		t.Fatalf("ApplyOrderEvent(unknown order) error = %v, want not_initialized", err)
	}
}

func TestSuccessfulAndFailedActionsUpdateTrackedState(t *testing.T) {
	t.Parallel()

	var failCancel atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/trading/orders/place":
			_, _ = w.Write([]byte(`{"client_order_id":"server-id","order_id":"OID-1"}`))
		case "/trading/orders/replace":
			_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
		case "/trading/orders/cancel":
			if failCancel.Load() {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"message":"cancel rejected"}`))
				return
			}
			_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	placed, err := c.PlaceOrder(context.Background(), validPlaceRequest())
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	if _, err := c.ReplaceOrder(context.Background(), validReplaceRequest()); err != nil {
		t.Fatalf("ReplaceOrder() error = %v", err)
	}
	if got := placed.State(); got != order.StateSubmitted {
		t.Fatalf("state after successful replace = %q, want SUBMITTED", got)
	}

	failCancel.Store(true)
	if _, err := c.CancelOrder(context.Background(), validCancelRequest()); err == nil {
		t.Fatal("CancelOrder() error = nil, want API failure")
	}
	if got := placed.State(); got != order.StateSubmitted {
		t.Fatalf("state after failed cancel = %q, want SUBMITTED", got)
	}

	failCancel.Store(false)
	if _, err := c.CancelOrder(context.Background(), validCancelRequest()); err != nil {
		t.Fatalf("CancelOrder() after recovery error = %v", err)
	}
	if got := placed.State(); got != order.StateCancelled {
		t.Fatalf("state after successful cancel = %q, want CANCELLED", got)
	}
}

func TestTerminalPreflightUsesAccountAndClientOrderID(t *testing.T) {
	t.Parallel()

	var replaceHits atomic.Int32
	var cancelHits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/trading/orders/place":
			_, _ = w.Write([]byte(`{"client_order_id":"same-id","order_id":"OID-1"}`))
		case "/trading/orders/replace":
			replaceHits.Add(1)
			_, _ = w.Write([]byte(`{"client_order_id":"same-id","order_id":"OID-1"}`))
		case "/trading/orders/cancel":
			cancelHits.Add(1)
			_, _ = w.Write([]byte(`{"client_order_id":"same-id","order_id":"OID-1"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	for _, accountID := range []string{"ACCOUNT-A", "ACCOUNT-B"} {
		req := validPlaceRequest()
		req.AccountID = accountID
		req.NewOrders[0].ClientOrderID = "same-id"
		if _, err := c.PlaceOrder(context.Background(), req); err != nil {
			t.Fatalf("PlaceOrder(%q) error = %v", accountID, err)
		}
	}
	if _, err := c.ApplyOrderEvent("ACCOUNT-A", "same-id", order.EventCancel); err != nil {
		t.Fatalf("ApplyOrderEvent(ACCOUNT-A) error = %v", err)
	}

	replaceReq := validReplaceRequest()
	replaceReq.AccountID = "ACCOUNT-A"
	replaceReq.ModifyOrders[0].ClientOrderID = "same-id"
	if _, err := c.ReplaceOrder(context.Background(), replaceReq); !errs.Is(err, errs.CodeInvalidTransition) {
		t.Fatalf("ReplaceOrder(terminal ACCOUNT-A) error = %v, want invalid_transition", err)
	}
	cancelReq := validCancelRequest()
	cancelReq.AccountID = "ACCOUNT-A"
	cancelReq.ClientOrderID = "same-id"
	if _, err := c.CancelOrder(context.Background(), cancelReq); !errs.Is(err, errs.CodeInvalidTransition) {
		t.Fatalf("CancelOrder(terminal ACCOUNT-A) error = %v, want invalid_transition", err)
	}
	if got := replaceHits.Load(); got != 0 {
		t.Fatalf("replace endpoint received %d requests, want 0", got)
	}
	if got := cancelHits.Load(); got != 0 {
		t.Fatalf("cancel endpoint received %d requests for terminal account, want 0", got)
	}

	cancelReq.AccountID = "ACCOUNT-B"
	if _, err := c.CancelOrder(context.Background(), cancelReq); err != nil {
		t.Fatalf("CancelOrder(ACCOUNT-B) error = %v", err)
	}
	if got := cancelHits.Load(); got != 1 {
		t.Fatalf("cancel endpoint received %d requests, want 1", got)
	}
}

func TestAutoClientOrderIDStableAndRequestNotMutated(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var placeIDs []string
	var batchIDs []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req trade.PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ids := make([]string, len(req.NewOrders))
		for i := range req.NewOrders {
			ids[i] = req.NewOrders[i].ClientOrderID
		}
		mu.Lock()
		if r.URL.Path == "/trading/orders/place" {
			placeIDs = append(placeIDs, ids...)
		} else {
			batchIDs = append(batchIDs, ids...)
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/trading/orders/place" {
			_, _ = w.Write([]byte(`{"client_order_id":"` + ids[0] + `","order_id":"OID-PLACE"}`))
			return
		}
		response := trade.BatchPlaceOrderResponse{Results: make([]trade.BatchPlaceOrderResult, len(ids))}
		for i := range ids {
			response.Results[i] = trade.BatchPlaceOrderResult{ClientOrderID: ids[i], OrderID: "OID-BATCH"}
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	c := newTradeClient(t, srv.URL, trade.WithAutoClientOrderID(true))
	placeReq := validPlaceRequest()
	placeReq.NewOrders[0].ClientOrderID = ""
	first, err := c.PlaceOrder(context.Background(), placeReq)
	if err != nil {
		t.Fatalf("first PlaceOrder() error = %v", err)
	}
	second, err := c.PlaceOrder(context.Background(), placeReq)
	if err != nil {
		t.Fatalf("second PlaceOrder() error = %v", err)
	}
	if first != second {
		t.Fatal("retry did not return the existing tracked order")
	}
	if placeReq.NewOrders[0].ClientOrderID != "" {
		t.Fatalf("PlaceOrder mutated caller request client_order_id to %q", placeReq.NewOrders[0].ClientOrderID)
	}
	mu.Lock()
	gotPlaceIDs := append([]string(nil), placeIDs...)
	mu.Unlock()
	if len(gotPlaceIDs) != 2 || gotPlaceIDs[0] == "" || gotPlaceIDs[0] != gotPlaceIDs[1] {
		t.Fatalf("captured PlaceOrder IDs = %#v, want two equal non-empty IDs", gotPlaceIDs)
	}
	if !trade.ValidClientOrderID(gotPlaceIDs[0]) {
		t.Fatalf("generated ID %q is invalid", gotPlaceIDs[0])
	}

	batchReq := validPlaceRequest()
	batchReq.NewOrders[0].ClientOrderID = ""
	secondOrder := validOrder()
	secondOrder.ClientOrderID = ""
	batchReq.NewOrders = []trade.OrderRequest{batchReq.NewOrders[0], secondOrder}
	if _, err := c.BatchPlaceOrder(context.Background(), batchReq); err != nil {
		t.Fatalf("first BatchPlaceOrder() error = %v", err)
	}
	if _, err := c.BatchPlaceOrder(context.Background(), batchReq); err != nil {
		t.Fatalf("second BatchPlaceOrder() error = %v", err)
	}
	if batchReq.NewOrders[0].ClientOrderID != "" || batchReq.NewOrders[1].ClientOrderID != "" {
		t.Fatalf("BatchPlaceOrder mutated caller request IDs to %q/%q", batchReq.NewOrders[0].ClientOrderID, batchReq.NewOrders[1].ClientOrderID)
	}
	mu.Lock()
	gotBatchIDs := append([]string(nil), batchIDs...)
	mu.Unlock()
	if len(gotBatchIDs) != 4 {
		t.Fatalf("captured BatchPlaceOrder IDs = %#v, want four IDs", gotBatchIDs)
	}
	for i := 0; i < 2; i++ {
		if gotBatchIDs[i] == "" || gotBatchIDs[i] != gotBatchIDs[i+2] {
			t.Fatalf("batch retry IDs = %#v, want stable per-index IDs", gotBatchIDs)
		}
	}
	for _, clientOrderID := range gotBatchIDs[:2] {
		if _, ok := c.GetTrackedOrder(batchReq.AccountID, clientOrderID); !ok {
			t.Fatalf("batch generated ID %q was not registered", clientOrderID)
		}
	}
}

func TestConcurrentTrackedOrderAccess(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"concurrent-id","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	req := validPlaceRequest()
	req.NewOrders[0].ClientOrderID = "concurrent-id"
	if _, err := c.PlaceOrder(context.Background(), req); err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}

	const workers = 32
	const iterations = 100
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := range workers {
		go func(worker int) {
			defer wg.Done()
			<-start
			for iteration := range iterations {
				if _, ok := c.GetTrackedOrder("ACC1", "concurrent-id"); !ok {
					t.Errorf("worker %d iteration %d lost tracked order", worker, iteration)
					return
				}
				_, _ = c.TrackedOrderState("ACC1", "concurrent-id")
				_, _ = c.ApplyOrderEvent("ACC1", "concurrent-id", order.EventCancelRequest)
				_, _ = c.ReconcileOrderStatus("ACC1", "concurrent-id", order.StatePending)
			}
		}(worker)
	}
	close(start)
	wg.Wait()
}

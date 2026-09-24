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
	"fmt"
	"strings"

	"github.com/shing1211/webullapi4go/pkg/domain/order"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

type orderRegistryKey struct {
	accountID     string
	clientOrderID string
}

// GetTrackedOrder returns the locally tracked order for accountID and
// clientOrderID. The returned pointer is the live tracked order, including its
// concurrency-safe state machine; a false result means the order is not tracked
// by this client.
func (c *Client) GetTrackedOrder(accountID, clientOrderID string) (*order.Order, bool) {
	return c.getOrder(accountID, clientOrderID)
}

// TrackedOrder returns the locally tracked order for accountID and
// clientOrderID. It is an accessor alias for [Client.GetTrackedOrder].
func (c *Client) TrackedOrder(accountID, clientOrderID string) (*order.Order, bool) {
	return c.GetTrackedOrder(accountID, clientOrderID)
}

// TrackedOrderState returns the current state of a locally tracked order. The
// second result is false when the order is not tracked or has no state machine.
func (c *Client) TrackedOrderState(accountID, clientOrderID string) (order.State, bool) {
	o, ok := c.GetTrackedOrder(accountID, clientOrderID)
	if !ok || o == nil || o.Machine == nil {
		return order.StateUnknown, false
	}
	return o.Machine.State(), true
}

// TrackedOrders returns a snapshot of the tracked-order pointers. The registry
// lock is released before the result is returned; callers must use the order's
// state-machine methods for synchronized state access.
func (c *Client) TrackedOrders() []*order.Order {
	if c == nil {
		return nil
	}
	c.ordMu.RLock()
	orders := make([]*order.Order, 0, len(c.orderRegistry))
	for _, o := range c.orderRegistry {
		orders = append(orders, o)
	}
	c.ordMu.RUnlock()
	return orders
}

// ApplyOrderEvent applies a domain order event to the locally tracked order for
// accountID and clientOrderID. It returns the resulting state and leaves the
// state unchanged when the event is invalid. The method accepts an [order.Event]
// directly, so callers do not need an event-stream-specific adapter. It returns
// [errs.CodeNotInitialized] when the requested order is not tracked.
func (c *Client) ApplyOrderEvent(accountID, clientOrderID string, event order.Event) (order.State, error) {
	o, ok := c.GetTrackedOrder(accountID, clientOrderID)
	if !ok || o == nil {
		return order.StateUnknown, orderNotTrackedError(accountID, clientOrderID)
	}
	if o.Machine == nil {
		return order.StateUnknown, errs.New(errs.CodeNotInitialized,
			fmt.Sprintf("order %q for account %q has no state machine", clientOrderID, accountID))
	}
	state, err := o.Machine.ApplyEvent(event)
	if err != nil {
		return state, errs.Wrap(errs.CodeInvalidTransition,
			fmt.Sprintf("applying event %q to order %q for account %q", event, clientOrderID, accountID), err)
	}
	return state, nil
}

// ReconcileOrderStatus reconciles a locally tracked order with a Webull status
// snapshot. status may be a string, [OrderStatus], or [order.State]. The
// state machine never moves a terminal order back to a non-terminal state and
// ignores stale non-terminal snapshots.
func (c *Client) ReconcileOrderStatus(accountID, clientOrderID string, status any) (order.State, error) {
	state, err := normalizeOrderStatus(status)
	if err != nil {
		return order.StateUnknown, err
	}
	return c.reconcileOrderState(accountID, clientOrderID, state)
}

// ReconcileOrderState reconciles a locally tracked order with an already
// normalized [order.State]. It is the typed counterpart to
// [Client.ReconcileOrderStatus].
func (c *Client) ReconcileOrderState(accountID, clientOrderID string, state order.State) (order.State, error) {
	if !knownOrderState(state) {
		return order.StateUnknown, errs.New(errs.CodeValidation,
			fmt.Sprintf("order status for account %q and client_order_id %q is not recognized", accountID, clientOrderID))
	}
	return c.reconcileOrderState(accountID, clientOrderID, state)
}

func (c *Client) reconcileOrderState(accountID, clientOrderID string, state order.State) (order.State, error) {
	o, ok := c.GetTrackedOrder(accountID, clientOrderID)
	if !ok || o == nil {
		return order.StateUnknown, orderNotTrackedError(accountID, clientOrderID)
	}
	if o.Machine == nil {
		return order.StateUnknown, errs.New(errs.CodeNotInitialized,
			fmt.Sprintf("order %q for account %q has no state machine", clientOrderID, accountID))
	}
	got, err := o.Machine.ReconcileStatus(state)
	if err != nil {
		return got, errs.Wrap(errs.CodeInvalidTransition,
			fmt.Sprintf("reconciling order %q for account %q to %q", clientOrderID, accountID, state), err)
	}
	return got, nil
}

func (c *Client) applyTrackedEvent(accountID, clientOrderID string, event order.Event) {
	_, _ = c.ApplyOrderEvent(accountID, clientOrderID, event)
}

func (c *Client) registerBatchResults(req PlaceOrderRequest, results []BatchPlaceOrderResult) {
	for i := range results {
		clientOrderID := batchResultClientOrderID(req, results[i], i)
		if clientOrderID == "" {
			continue
		}
		c.registerOrder(&order.Order{
			PlaceOrderResult: order.PlaceOrderResult{
				ClientOrderID: clientOrderID,
				OrderID:       results[i].OrderID,
			},
			AccountID: req.AccountID,
			Machine:   order.New(order.StatePending),
		})
	}
}

func batchResultClientOrderID(req PlaceOrderRequest, result BatchPlaceOrderResult, index int) string {
	if result.ClientOrderID != "" {
		for _, requestOrder := range req.NewOrders {
			if requestOrder.ClientOrderID == result.ClientOrderID {
				return result.ClientOrderID
			}
		}
	}
	if index < len(req.NewOrders) {
		return req.NewOrders[index].ClientOrderID
	}
	return result.ClientOrderID
}

func (c *Client) preparePlaceRequest(req PlaceOrderRequest) (PlaceOrderRequest, error) {
	if !c.cfg.autoClientOrderID {
		return req, nil
	}
	req.NewOrders = append([]OrderRequest(nil), req.NewOrders...)
	for i := range req.NewOrders {
		if req.NewOrders[i].ClientOrderID != "" {
			continue
		}
		req.NewOrders[i].ClientOrderID = stableClientOrderID(req, i)
	}
	return req, nil
}

func stableClientOrderID(req PlaceOrderRequest, index int) string {
	orderRequest := req.NewOrders[index]
	orderRequest.ClientOrderID = ""
	seed := struct {
		AccountID          string
		ClientComboOrderID string
		Index              int
		Order              OrderRequest
	}{
		AccountID:          req.AccountID,
		ClientComboOrderID: req.ClientComboOrderID,
		Index:              index,
		Order:              orderRequest,
	}
	data, err := json.Marshal(seed)
	if err != nil {
		data = []byte(fmt.Sprintf("%s\x00%s\x00%d\x00%#v", seed.AccountID, seed.ClientComboOrderID, seed.Index, seed.Order))
	}
	return ClientOrderIDFrom(data)
}

func normalizeOrderStatus(status any) (order.State, error) {
	var value string
	switch v := status.(type) {
	case string:
		value = v
	case order.State:
		value = string(v)
	case OrderStatus:
		value = string(v)
	case []byte:
		value = string(v)
	default:
		return order.StateUnknown, errs.New(errs.CodeValidation,
			fmt.Sprintf("order status has unsupported type %T", status))
	}
	state := order.FromWebullStatus(strings.TrimSpace(value))
	if state == order.StateUnknown {
		return order.StateUnknown, errs.New(errs.CodeValidation,
			fmt.Sprintf("order status %q is not recognized", value))
	}
	return state, nil
}

func knownOrderState(state order.State) bool {
	switch state {
	case order.StatePending, order.StateSubmitted, order.StatePreSubmitted, order.StateConfirmed,
		order.StatePartialFilled, order.StateFilled, order.StateCancelled, order.StateFailed, order.StateExpired:
		return true
	default:
		return false
	}
}

func orderNotTrackedError(accountID, clientOrderID string) error {
	return errs.New(errs.CodeNotInitialized,
		fmt.Sprintf("order %q for account %q is not tracked", clientOrderID, accountID))
}

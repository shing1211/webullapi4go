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

package events

import (
	"context"
	"io"
	"time"

	"google.golang.org/grpc"

	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
	"github.com/shing1211/webullapi4go/internal/errs"
)

type EventType = eventsevents.EventType

const (
	EventTypeSubscribeSuccess = eventsevents.EventType_SubscribeSuccess
	EventTypePing             = eventsevents.EventType_Ping
	EventTypeAuthError        = eventsevents.EventType_AuthError
	EventTypeNumOfConnExceed  = eventsevents.EventType_NumOfConnExceed
	EventTypeSubscribeExpired = eventsevents.EventType_SubscribeExpired
)

type SubscribeResponse struct {
	*eventsevents.SubscribeResponse
}

func (r *SubscribeResponse) EventType() EventType { return r.GetEventType() }
func (r *SubscribeResponse) ContentType() string  { return r.GetContentType() }
func (r *SubscribeResponse) Payload() string      { return r.GetPayload() }
func (r *SubscribeResponse) RequestId() string    { return r.GetRequestId() }
func (r *SubscribeResponse) Timestamp() int64     { return r.GetTimestamp() }

type Client struct {
	svc     eventsevents.EventServiceClient
	timeout time.Duration
}

func NewClient(svc eventsevents.EventServiceClient) *Client {
	return &Client{svc: svc, timeout: 30 * time.Second}
}

type Option func(*Client)

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.timeout = d }
}

func (c *Client) Subscribe(ctx context.Context, req *SubscribeRequest, opts ...grpc.CallOption) (*streamClient, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	stream, err := c.svc.Subscribe(ctx, req, opts...)
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, "subscribe failed", err)
	}
	return &streamClient{stream: stream}, nil
}

type streamClient struct {
	stream eventsevents.EventService_SubscribeClient
}

func (s *streamClient) Recv() (*SubscribeResponse, error) {
	resp, err := s.stream.Recv()
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, "recv failed", err)
	}
	return &SubscribeResponse{SubscribeResponse: resp}, nil
}

type SubscribeRequest = eventsevents.SubscribeRequest

func NewSubscribeRequest(subscribeType uint32, accounts []string) *SubscribeRequest {
	return &SubscribeRequest{
		SubscribeType: subscribeType,
		Timestamp:     time.Now().UnixMilli(),
		ContentType:   "application/json",
		Accounts:      accounts,
	}
}

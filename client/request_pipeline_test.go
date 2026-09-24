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

package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/internal/auth"
)

type pipelineEvents struct {
	mu           sync.Mutex
	order        []string
	attempts     []int
	statuses     []int
	errors       int
	correlations []string
}

func (e *pipelineEvents) add(value string) {
	e.mu.Lock()
	e.order = append(e.order, value)
	e.mu.Unlock()
}

func (e *pipelineEvents) addAttempt(attempt int) {
	e.mu.Lock()
	e.attempts = append(e.attempts, attempt)
	e.mu.Unlock()
}

func (e *pipelineEvents) addStatus(status int) {
	e.mu.Lock()
	e.statuses = append(e.statuses, status)
	e.mu.Unlock()
}

func (e *pipelineEvents) addError() {
	e.mu.Lock()
	e.errors++
	e.mu.Unlock()
}

func (e *pipelineEvents) addCorrelation(id string) {
	e.mu.Lock()
	e.correlations = append(e.correlations, id)
	e.mu.Unlock()
}

func (e *pipelineEvents) snapshot() ([]string, []int, []int, int, []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	order := append([]string(nil), e.order...)
	attempts := append([]int(nil), e.attempts...)
	statuses := append([]int(nil), e.statuses...)
	correlations := append([]string(nil), e.correlations...)
	return order, attempts, statuses, e.errors, correlations
}

type pipelineRateLimiter struct {
	events *pipelineEvents
}

func (l *pipelineRateLimiter) Wait(context.Context, string) error {
	l.events.add("rate")
	return nil
}

type pipelineBreaker struct {
	events *pipelineEvents
}

func (b *pipelineBreaker) Allow() bool {
	b.events.add("breaker")
	return true
}

func (b *pipelineBreaker) RecordSuccess() {}
func (b *pipelineBreaker) RecordFailure() {}

func TestEntrypointsShareHooksInterceptorsAndLogger(t *testing.T) {
	cases := []struct {
		name string
		call func(*client.Client) error
	}{
		{
			name: "do",
			call: func(cl *client.Client) error {
				return cl.Do(context.Background(), http.MethodGet, "/test", nil, nil)
			},
		},
		{
			name: "broker",
			call: func(cl *client.Client) error {
				return cl.DoBroker(context.Background(), http.MethodGet, "/test", nil, nil)
			},
		},
		{
			name: "stream",
			call: func(cl *client.Client) error {
				resp, err := cl.DoStream(context.Background(), http.MethodGet, "/test", nil)
				if resp != nil {
					_ = resp.Body.Close()
				}
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			events := &pipelineEvents{}
			var logOutput bytes.Buffer
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				events.add("server")
				events.addCorrelation(r.Header.Get("x-correlation-id"))
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer srv.Close()

			cl, err := client.New(
				client.WithCredentials(testAppKey, testAppSecret),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}),
				client.WithRateLimiter(&pipelineRateLimiter{events: events}),
				client.WithBreaker(&pipelineBreaker{events: events}),
				client.WithHooks(client.Hooks{
					OnRequest:  func(string, string) { events.add("hook") },
					OnLatency:  func(attempt int, _ time.Duration) { events.addAttempt(attempt) },
					OnResponse: func(status int, _ time.Duration) { events.addStatus(status) },
					OnError:    func(error, time.Duration) { events.addError() },
				}),
				client.WithInterceptor(func(ctx context.Context, next func(context.Context) error) error {
					events.add("a-before")
					err := next(ctx)
					events.add("a-after")
					return err
				}),
				client.WithInterceptor(func(ctx context.Context, next func(context.Context) error) error {
					events.add("b-before")
					err := next(ctx)
					events.add("b-after")
					return err
				}),
				client.WithLogger(slog.New(slog.NewJSONHandler(&logOutput, nil))),
			)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			defer cl.Close()

			if err := tc.call(cl); err != nil {
				t.Fatalf("call error = %v", err)
			}

			order, attempts, statuses, errorCount, correlations := events.snapshot()
			wantOrder := []string{"rate", "breaker", "hook", "a-before", "b-before", "server", "b-after", "a-after"}
			if len(order) != len(wantOrder) {
				t.Fatalf("pipeline order = %v, want %v", order, wantOrder)
			}
			for i := range wantOrder {
				if order[i] != wantOrder[i] {
					t.Fatalf("pipeline order = %v, want %v", order, wantOrder)
				}
			}
			if len(attempts) != 1 || attempts[0] != 1 {
				t.Fatalf("attempts = %v, want [1]", attempts)
			}
			if len(statuses) != 1 || statuses[0] != http.StatusCreated {
				t.Fatalf("statuses = %v, want [%d]", statuses, http.StatusCreated)
			}
			if errorCount != 0 {
				t.Fatalf("error hook count = %d, want 0", errorCount)
			}
			if len(correlations) != 1 || correlations[0] == "" {
				t.Fatalf("correlation IDs = %v, want one non-empty ID", correlations)
			}

			logs := logOutput.String()
			for _, want := range []string{"webull request start", "webull request done", `"attempt":1`, `"status":201`, correlations[0]} {
				if !bytes.Contains([]byte(logs), []byte(want)) {
					t.Errorf("logs do not contain %q: %s", want, logs)
				}
			}
		})
	}
}

func TestCorrelationIDIsStableAcrossRetries(t *testing.T) {
	cases := []struct {
		name string
		call func(*client.Client) error
	}{
		{
			name: "do",
			call: func(cl *client.Client) error {
				return cl.Do(context.Background(), http.MethodGet, "/test", nil, nil)
			},
		},
		{
			name: "broker",
			call: func(cl *client.Client) error {
				return cl.DoBroker(context.Background(), http.MethodGet, "/test", nil, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var ids []string
			var attempts int
			var hookAttempts []int
			var statuses []int
			var errorCount int
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				ids = append(ids, r.Header.Get("x-correlation-id"))
				attempts++
				attempt := attempts
				mu.Unlock()
				if attempt == 1 {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"message":"retry"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer srv.Close()

			cl, err := client.New(
				client.WithCredentials(testAppKey, testAppSecret),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}),
				client.WithRetry(client.RetryConfig{
					MaxAttempts: 2,
					BaseDelay:   time.Nanosecond,
					MaxDelay:    time.Nanosecond,
				}),
				client.WithHooks(client.Hooks{
					OnLatency: func(attempt int, _ time.Duration) {
						mu.Lock()
						hookAttempts = append(hookAttempts, attempt)
						mu.Unlock()
					},
					OnResponse: func(status int, _ time.Duration) {
						mu.Lock()
						statuses = append(statuses, status)
						mu.Unlock()
					},
					OnError: func(error, time.Duration) {
						mu.Lock()
						errorCount++
						mu.Unlock()
					},
				}),
			)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			defer cl.Close()

			if err := tc.call(cl); err != nil {
				t.Fatalf("call error = %v", err)
			}
			mu.Lock()
			got := append([]string(nil), ids...)
			gotAttempts := append([]int(nil), hookAttempts...)
			gotStatuses := append([]int(nil), statuses...)
			gotErrors := errorCount
			mu.Unlock()
			if len(got) != 2 {
				t.Fatalf("request count = %d, want 2", len(got))
			}
			if got[0] == "" || got[0] != got[1] {
				t.Fatalf("correlation IDs = %v, want one stable non-empty ID", got)
			}
			if len(gotAttempts) != 2 || gotAttempts[0] != 1 || gotAttempts[1] != 2 {
				t.Fatalf("hook attempts = %v, want [1 2]", gotAttempts)
			}
			if len(gotStatuses) != 1 || gotStatuses[0] != http.StatusOK {
				t.Fatalf("hook statuses = %v, want [%d]", gotStatuses, http.StatusOK)
			}
			if gotErrors != 1 {
				t.Fatalf("hook error count = %d, want 1", gotErrors)
			}
		})
	}
}

func TestRawMessageBytesArePreservedAcrossRetries(t *testing.T) {
	wantBody := json.RawMessage(`{"raw":"<value>&more"}`)
	var mu sync.Mutex
	var bodies [][]byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies = append(bodies, body)
		attempt := len(bodies)
		mu.Unlock()
		verifySignature(t, r, body)
		if attempt == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"retry"}`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl, err := client.New(
		client.WithCredentials(testAppKey, testAppSecret),
		client.WithBaseURL(srv.URL),
		client.WithRetry(client.RetryConfig{
			MaxAttempts:        2,
			BaseDelay:          time.Nanosecond,
			MaxDelay:           time.Nanosecond,
			RetryNonIdempotent: true,
		}),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()

	if err := cl.Do(context.Background(), http.MethodPost, "/test", wantBody, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	mu.Lock()
	got := append([][]byte(nil), bodies...)
	mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("request count = %d, want 2", len(got))
	}
	for i, body := range got {
		if string(body) != string(wantBody) {
			t.Errorf("request %d body = %q, want %q", i+1, body, wantBody)
		}
	}
}

func TestClockOffsetIsLearnedFromBrokerAndStreamResponses(t *testing.T) {
	cases := []struct {
		name string
		call func(*client.Client) error
	}{
		{
			name: "broker",
			call: func(cl *client.Client) error {
				return cl.DoBroker(context.Background(), http.MethodGet, "/test", nil, nil)
			},
		},
		{
			name: "stream",
			call: func(cl *client.Client) error {
				resp, err := cl.DoStream(context.Background(), http.MethodGet, "/test", nil)
				if resp != nil {
					_ = resp.Body.Close()
				}
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			serverTime := time.Now().Add(2 * time.Minute).UTC().Truncate(time.Second)
			var mu sync.Mutex
			var timestamps []string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				timestamps = append(timestamps, r.Header.Get(auth.HeaderTimestamp))
				mu.Unlock()
				w.Header().Set("Date", serverTime.Format(http.TimeFormat))
				_, _ = w.Write([]byte(`{}`))
			}))
			defer srv.Close()

			cl, err := client.New(
				client.WithCredentials(testAppKey, testAppSecret),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}),
				client.WithClockDriftCorrection(true),
			)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			defer cl.Close()

			if err := tc.call(cl); err != nil {
				t.Fatalf("first call error = %v", err)
			}
			if err := tc.call(cl); err != nil {
				t.Fatalf("second call error = %v", err)
			}
			mu.Lock()
			got := append([]string(nil), timestamps...)
			mu.Unlock()
			if len(got) != 2 {
				t.Fatalf("request count = %d, want 2", len(got))
			}
			ts, err := time.Parse(auth.TimestampFormat, got[1])
			if err != nil {
				t.Fatalf("parse second timestamp %q: %v", got[1], err)
			}
			want := time.Now().Add(2 * time.Minute)
			if delta := ts.Sub(want); delta < -3*time.Second || delta > 3*time.Second {
				t.Fatalf("second timestamp = %v, want within 3s of %v", ts, want)
			}
		})
	}
}

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

package stream

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"

	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const defaultChannelBuffer = 100

type DropPolicy int

const (
	DropBlock DropPolicy = iota
	DropOldest
	DropSample
)

func (p DropPolicy) String() string {
	switch p {
	case DropBlock:
		return "block"
	case DropOldest:
		return "drop-oldest"
	case DropSample:
		return "sample"
	default:
		return "block"
	}
}

type channelConfig struct {
	policy      DropPolicy
	bufSize     int
	dropCnt     atomic.Int64
	otelCounter metric.Int64Counter
	topic       string
}

func (c *channelConfig) recordDrop() { c.dropCnt.Add(1) }

func (c *channelConfig) recordDropToMeter() {
	c.dropCnt.Add(1)
	// OTel metric recording is intentionally fire-and-forget: the metric SDK
	// records synchronously and cannot block, so context.Background() is safe.
	if c.otelCounter != nil {
		c.otelCounter.Add(context.Background(), 1,
			metric.WithAttributes(attribute.String("topic", c.topic)))
	}
}

func (c *channelConfig) DropCount() int64 { return c.dropCnt.Load() }

type chanRegistry struct {
	mu       sync.RWMutex
	quote    map[*chanQuote]*channelConfig
	snapshot map[*chanSnapshot]*channelConfig
	tick     map[*chanTick]*channelConfig
}

type chanQuote struct {
	ch   chan<- *marketdatav1.Quote
	stop func()
}

type chanSnapshot struct {
	ch   chan<- *marketdatav1.Snapshot
	stop func()
}

type chanTick struct {
	ch   chan<- *marketdatav1.Tick
	stop func()
}

func newChanRegistry() *chanRegistry {
	return &chanRegistry{
		quote:    make(map[*chanQuote]*channelConfig),
		snapshot: make(map[*chanSnapshot]*channelConfig),
		tick:     make(map[*chanTick]*channelConfig),
	}
}

func (r *chanRegistry) stopAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for entry := range r.quote {
		entry.stop()
	}
	for entry := range r.snapshot {
		entry.stop()
	}
	for entry := range r.tick {
		entry.stop()
	}
}

func (r *chanRegistry) dispatchQuote(msg *marketdatav1.Quote) {
	r.mu.RLock()
	for entry, cfg := range r.quote {
		select {
		case entry.ch <- msg:
		default:
			switch cfg.policy {
			case DropOldest:
				select {
				case entry.ch <- msg:
				default:
					cfg.recordDropToMeter()
				}
			case DropSample:
				if rand.Intn(2) == 0 {
					cfg.recordDropToMeter()
				} else {
					select {
					case entry.ch <- msg:
					default:
						cfg.recordDropToMeter()
					}
				}
			default:
				entry.ch <- msg
			}
		}
	}
	r.mu.RUnlock()
}

func (r *chanRegistry) dispatchSnapshot(msg *marketdatav1.Snapshot) {
	r.mu.RLock()
	for entry, cfg := range r.snapshot {
		select {
		case entry.ch <- msg:
		default:
			switch cfg.policy {
			case DropOldest:
				select {
				case entry.ch <- msg:
				default:
					cfg.recordDropToMeter()
				}
			case DropSample:
				if rand.Intn(2) == 0 {
					cfg.recordDropToMeter()
				} else {
					select {
					case entry.ch <- msg:
					default:
						cfg.recordDropToMeter()
					}
				}
			default:
				entry.ch <- msg
			}
		}
	}
	r.mu.RUnlock()
}

func (r *chanRegistry) dispatchTick(msg *marketdatav1.Tick) {
	r.mu.RLock()
	for entry, cfg := range r.tick {
		select {
		case entry.ch <- msg:
		default:
			switch cfg.policy {
			case DropOldest:
				select {
				case entry.ch <- msg:
				default:
					cfg.recordDropToMeter()
				}
			case DropSample:
				if rand.Intn(2) == 0 {
					cfg.recordDropToMeter()
				} else {
					select {
					case entry.ch <- msg:
					default:
						cfg.recordDropToMeter()
					}
				}
			default:
				entry.ch <- msg
			}
		}
	}
	r.mu.RUnlock()
}

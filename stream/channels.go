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

type channelLifecycle struct {
	initOnce sync.Once
	stopOnce sync.Once
	stopCh   chan struct{}
	sendMu   sync.Mutex
}

func (l *channelLifecycle) done() chan struct{} {
	l.initOnce.Do(func() {
		l.stopCh = make(chan struct{})
	})
	return l.stopCh
}

func (l *channelLifecycle) stop(closeFn func()) {
	l.stopOnce.Do(func() {
		close(l.done())
		l.sendMu.Lock()
		defer l.sendMu.Unlock()
		if closeFn != nil {
			closeFn()
		}
	})
}

func sendChannelMessage[T any](l *channelLifecycle, ch chan T, cfg *channelConfig, msg T, policy DropPolicy) {
	l.sendMu.Lock()
	defer l.sendMu.Unlock()

	stopCh := l.done()
	if channelStopped(stopCh) {
		return
	}

	select {
	case ch <- msg:
		return
	default:
	}

	switch policy {
	case DropOldest:
		select {
		case <-ch:
			recordChannelDrop(cfg)
		default:
		}
		if channelStopped(stopCh) {
			return
		}
		select {
		case ch <- msg:
		default:
			recordChannelDrop(cfg)
		}
	case DropSample:
		if rand.Intn(2) == 0 {
			recordChannelDrop(cfg)
			return
		}
		select {
		case ch <- msg:
		default:
			recordChannelDrop(cfg)
		}
	default:
		select {
		case ch <- msg:
		case <-stopCh:
		}
	}
}

func channelStopped(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func recordChannelDrop(cfg *channelConfig) {
	if cfg != nil {
		cfg.recordDropToMeter()
	}
}

type chanRegistry struct {
	mu       sync.RWMutex
	closed   bool
	quote    map[*chanQuote]*channelConfig
	snapshot map[*chanSnapshot]*channelConfig
	tick     map[*chanTick]*channelConfig
}

type chanQuote struct {
	ch        chan *marketdatav1.Quote
	stop      func()
	lifecycle channelLifecycle
}

func newChanQuote(ch chan *marketdatav1.Quote) *chanQuote {
	return &chanQuote{ch: ch, stop: func() { close(ch) }}
}

func (e *chanQuote) shutdown() {
	e.lifecycle.stop(func() {
		if e.stop != nil {
			e.stop()
			return
		}
		close(e.ch)
	})
}

func (e *chanQuote) send(msg *marketdatav1.Quote, cfg *channelConfig) {
	policy := DropBlock
	if cfg != nil {
		policy = cfg.policy
	}
	sendChannelMessage(&e.lifecycle, e.ch, cfg, msg, policy)
}

type chanSnapshot struct {
	ch        chan *marketdatav1.Snapshot
	stop      func()
	lifecycle channelLifecycle
}

func newChanSnapshot(ch chan *marketdatav1.Snapshot) *chanSnapshot {
	return &chanSnapshot{ch: ch, stop: func() { close(ch) }}
}

func (e *chanSnapshot) shutdown() {
	e.lifecycle.stop(func() {
		if e.stop != nil {
			e.stop()
			return
		}
		close(e.ch)
	})
}

func (e *chanSnapshot) send(msg *marketdatav1.Snapshot, cfg *channelConfig) {
	policy := DropBlock
	if cfg != nil {
		policy = cfg.policy
	}
	sendChannelMessage(&e.lifecycle, e.ch, cfg, msg, policy)
}

type chanTick struct {
	ch        chan *marketdatav1.Tick
	stop      func()
	lifecycle channelLifecycle
}

func newChanTick(ch chan *marketdatav1.Tick) *chanTick {
	return &chanTick{ch: ch, stop: func() { close(ch) }}
}

func (e *chanTick) shutdown() {
	e.lifecycle.stop(func() {
		if e.stop != nil {
			e.stop()
			return
		}
		close(e.ch)
	})
}

func (e *chanTick) send(msg *marketdatav1.Tick, cfg *channelConfig) {
	policy := DropBlock
	if cfg != nil {
		policy = cfg.policy
	}
	sendChannelMessage(&e.lifecycle, e.ch, cfg, msg, policy)
}

func newChanRegistry() *chanRegistry {
	return &chanRegistry{
		quote:    make(map[*chanQuote]*channelConfig),
		snapshot: make(map[*chanSnapshot]*channelConfig),
		tick:     make(map[*chanTick]*channelConfig),
	}
}

func (r *chanRegistry) addQuote(entry *chanQuote, cfg *channelConfig) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return false
	}
	if r.quote == nil {
		r.quote = make(map[*chanQuote]*channelConfig)
	}
	r.quote[entry] = cfg
	return true
}

func (r *chanRegistry) addSnapshot(entry *chanSnapshot, cfg *channelConfig) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return false
	}
	if r.snapshot == nil {
		r.snapshot = make(map[*chanSnapshot]*channelConfig)
	}
	r.snapshot[entry] = cfg
	return true
}

func (r *chanRegistry) addTick(entry *chanTick, cfg *channelConfig) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return false
	}
	if r.tick == nil {
		r.tick = make(map[*chanTick]*channelConfig)
	}
	r.tick[entry] = cfg
	return true
}

func (r *chanRegistry) removeQuote(entry *chanQuote) {
	r.mu.Lock()
	delete(r.quote, entry)
	r.mu.Unlock()
	entry.shutdown()
}

func (r *chanRegistry) removeSnapshot(entry *chanSnapshot) {
	r.mu.Lock()
	delete(r.snapshot, entry)
	r.mu.Unlock()
	entry.shutdown()
}

func (r *chanRegistry) removeTick(entry *chanTick) {
	r.mu.Lock()
	delete(r.tick, entry)
	r.mu.Unlock()
	entry.shutdown()
}

func (r *chanRegistry) stopAll() {
	r.mu.Lock()
	r.closed = true
	quote := make([]*chanQuote, 0, len(r.quote))
	for entry := range r.quote {
		quote = append(quote, entry)
	}
	snapshot := make([]*chanSnapshot, 0, len(r.snapshot))
	for entry := range r.snapshot {
		snapshot = append(snapshot, entry)
	}
	tick := make([]*chanTick, 0, len(r.tick))
	for entry := range r.tick {
		tick = append(tick, entry)
	}
	r.quote = make(map[*chanQuote]*channelConfig)
	r.snapshot = make(map[*chanSnapshot]*channelConfig)
	r.tick = make(map[*chanTick]*channelConfig)
	r.mu.Unlock()

	for _, entry := range quote {
		entry.shutdown()
	}
	for _, entry := range snapshot {
		entry.shutdown()
	}
	for _, entry := range tick {
		entry.shutdown()
	}
}

func (r *chanRegistry) dispatchQuote(msg *marketdatav1.Quote) {
	type dispatchEntry struct {
		entry *chanQuote
		cfg   *channelConfig
	}

	r.mu.RLock()
	entries := make([]dispatchEntry, 0, len(r.quote))
	for entry, cfg := range r.quote {
		entries = append(entries, dispatchEntry{entry: entry, cfg: cfg})
	}
	r.mu.RUnlock()
	for _, item := range entries {
		item.entry.send(msg, item.cfg)
	}
}

func (r *chanRegistry) dispatchSnapshot(msg *marketdatav1.Snapshot) {
	type dispatchEntry struct {
		entry *chanSnapshot
		cfg   *channelConfig
	}

	r.mu.RLock()
	entries := make([]dispatchEntry, 0, len(r.snapshot))
	for entry, cfg := range r.snapshot {
		entries = append(entries, dispatchEntry{entry: entry, cfg: cfg})
	}
	r.mu.RUnlock()
	for _, item := range entries {
		item.entry.send(msg, item.cfg)
	}
}

func (r *chanRegistry) dispatchTick(msg *marketdatav1.Tick) {
	type dispatchEntry struct {
		entry *chanTick
		cfg   *channelConfig
	}

	r.mu.RLock()
	entries := make([]dispatchEntry, 0, len(r.tick))
	for entry, cfg := range r.tick {
		entries = append(entries, dispatchEntry{entry: entry, cfg: cfg})
	}
	r.mu.RUnlock()
	for _, item := range entries {
		item.entry.send(msg, item.cfg)
	}
}

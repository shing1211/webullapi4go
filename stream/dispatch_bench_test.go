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
	"strconv"
	"sync"
	"testing"
	"time"

	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
)

// These benchmarks measure the synchronous dispatch path only. They exist to
// quantify the head-of-line behavior described in ARCHITECTURE.md, not to assert
// a latency objective: no target has been agreed, so the numbers are recorded as
// measurements and the retain-or-change decision is still open.
//
// Dispatch is synchronous by design. Handlers run in registration order on the
// receive path, and channel subscribers are sorted by subscription order and
// sent to one at a time, so a slow subscriber delays every later delivery. The
// DropBlock policy makes that explicit: a full buffer blocks the sender until
// the consumer reads or the subscription is cancelled.

// benchQuote builds a realistic quote payload. Price, size, and timestamp are
// decimal strings on the wire, so the benchmark allocates strings like the real
// decode path does rather than reusing one shared value.
func benchQuote(i int) *marketdatav1.Quote {
	s := strconv.Itoa(i)
	return &marketdatav1.Quote{
		Basic: &marketdatav1.Basic{
			Symbol:    "AAPL",
			Timestamp: s,
		},
		Asks: []*marketdatav1.AskBid{{Price: s, Size: "1"}},
		Bids: []*marketdatav1.AskBid{{Price: s, Size: "2"}},
	}
}

// drain reads from ch until it is closed, so a channel subscriber does not
// become the bottleneck being measured.
func drain[T any](ch <-chan T, done *sync.WaitGroup) {
	defer done.Done()
	for range ch {
	}
}

// BenchmarkCallbackDispatch measures the cost of routing one message to N
// registered OnQuote handlers, with no channel subscribers. This is the floor
// imposed by synchronous fan-out.
func BenchmarkCallbackDispatch(b *testing.B) {
	for _, n := range []int{0, 1, 4, 16, 64} {
		b.Run(subName("handlers", n), func(b *testing.B) {
			c := &Client{chanReg: newChanRegistry()}
			var sink int
			for i := 0; i < n; i++ {
				c.OnQuote(func(*marketdatav1.Quote) { sink++ })
			}
			msg := benchQuote(1)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.emitQuote(msg)
			}
			b.StopTimer()
			if sink < 0 {
				b.Fatal("unreachable")
			}
		})
	}
}

// BenchmarkChannelDispatchPolicy measures dispatch to a single channel
// subscriber that is drained concurrently, so the buffer never fills and the
// result reflects send cost rather than blocking.
func BenchmarkChannelDispatchPolicy(b *testing.B) {
	policies := []struct {
		name   string
		policy DropPolicy
	}{
		{"DropBlock", DropBlock},
		{"DropOldest", DropOldest},
		{"DropSample", DropSample},
	}

	for _, p := range policies {
		b.Run(p.name, func(b *testing.B) {
			c := &Client{chanReg: newChanRegistry()}
			ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: p.policy, BufferSize: 100})
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			go drain(ch, &wg)

			msg := benchQuote(1)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.emitQuote(msg)
			}
			b.StopTimer()

			cancel()
			wg.Wait()
		})
	}
}

// BenchmarkChannelDispatchSubscriberCount measures how dispatch cost scales with
// the number of channel subscribers on a single message. All subscribers are
// drained concurrently, so this isolates fan-out cost from blocking.
func BenchmarkChannelDispatchSubscriberCount(b *testing.B) {
	for _, n := range []int{1, 4, 16, 64} {
		b.Run(subName("subscribers", n), func(b *testing.B) {
			c := &Client{chanReg: newChanRegistry()}
			var wg sync.WaitGroup
			cancels := make([]func(), 0, n)
			for i := 0; i < n; i++ {
				ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 100})
				cancels = append(cancels, cancel)
				wg.Add(1)
				go drain(ch, &wg)
			}

			msg := benchQuote(1)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.emitQuote(msg)
			}
			b.StopTimer()

			for _, cancel := range cancels {
				cancel()
			}
			wg.Wait()
		})
	}
}

// releaseAfter is how long the blocked subscriber is held before it is
// cancelled. Windows timer granularity can overshoot this, so the reported
// blocked time is a lower bound on the intended delay, which only strengthens
// the contrast with the non-blocking policy.
const releaseAfter = 2 * time.Millisecond

// sink prevents the compiler from eliminating a callback whose result is unused.
var sink int

// BenchmarkHeadOfLineBlockedDelivery is the measurement that motivates the open
// design question.
//
// A subscriber registered first is never read, so its buffer fills. Because
// dispatch is synchronous and ordered by subscription, a second subscriber that
// is ready to receive waits behind it: the delivery of a message to a ready
// subscriber is delayed until the blocked one is released. The blocked case is
// released by cancellation, which is the only thing that can unblock a DropBlock
// send, and the benchmark reports how long a ready subscriber actually waited.
//
// The dropOldest case is the control: the same shape with a non-blocking policy
// never waits, so the ready subscriber is served immediately.
func BenchmarkHeadOfLineBlockedDelivery(b *testing.B) {
	b.Run("DropBlock", func(b *testing.B) {
		var total time.Duration
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c := &Client{chanReg: newChanRegistry()}

			_, cancelSlow := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 1})
			// Fill the slow subscriber's buffer; nothing reads it.
			c.emitQuote(benchQuote(0))

			// A ready subscriber registered after the blocked one.
			fast, cancelFast := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 1})
			var wg sync.WaitGroup
			wg.Add(1)
			go drain(fast, &wg)

			go func() {
				time.Sleep(releaseAfter)
				cancelSlow()
			}()

			start := time.Now()
			// This send must wait behind the blocked subscriber.
			c.emitQuote(benchQuote(1))
			total += time.Since(start)

			cancelFast()
			wg.Wait()
		}
		b.StopTimer()
		if b.N > 0 {
			b.ReportMetric(float64(total.Nanoseconds())/float64(b.N), "blocked-ns/op")
		}
	})

	b.Run("DropOldest", func(b *testing.B) {
		var total time.Duration
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c := &Client{chanReg: newChanRegistry()}

			// Same shape: a subscriber that is never read.
			_, cancelSlow := c.SubscribeQuoteChan(ChannelConfig{Policy: DropOldest, BufferSize: 1})
			c.emitQuote(benchQuote(0))

			fast, cancelFast := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 1})
			var wg sync.WaitGroup
			wg.Add(1)
			go drain(fast, &wg)

			start := time.Now()
			// The non-blocking policy discards instead of waiting.
			c.emitQuote(benchQuote(1))
			total += time.Since(start)

			cancelFast()
			cancelSlow()
			wg.Wait()
		}
		b.StopTimer()
		if b.N > 0 {
			b.ReportMetric(float64(total.Nanoseconds())/float64(b.N), "blocked-ns/op")
		}
	})
}

// BenchmarkCallbackCost measures the per-handler cost of a minimal callback so
// the callback-dispatch numbers above can be attributed.
func BenchmarkCallbackCost(b *testing.B) {
	fn := func(*marketdatav1.Quote) { sink++ }
	msg := benchQuote(1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(msg)
	}
	b.StopTimer()
}

// BenchmarkDispatchJitter measures dispatch with a realistic mix of registered
// handlers and channel subscribers together, which is closer to production
// fan-out than either path measured alone.
func BenchmarkDispatchJitter(b *testing.B) {
	c := &Client{chanReg: newChanRegistry()}
	for i := 0; i < 8; i++ {
		c.OnQuote(func(*marketdatav1.Quote) { sink++ })
	}
	var wg sync.WaitGroup
	cancels := make([]func(), 0, 32)
	for i := 0; i < 32; i++ {
		ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 100})
		cancels = append(cancels, cancel)
		wg.Add(1)
		go drain(ch, &wg)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// benchQuote varies the payload so the benchmark does not benefit from
		// cache reuse of an identical message.
		c.emitQuote(benchQuote(i))
	}
	b.StopTimer()

	for _, cancel := range cancels {
		cancel()
	}
	wg.Wait()
}

// subName builds a stable sub-benchmark name.
func subName(kind string, n int) string {
	return kind + "=" + itoa(n)
}

// itoa avoids importing strconv purely for benchmark labels.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// BenchmarkShutdownTime measures releasing N subscriptions, which is the cost a
// caller pays to stop a blocked subscriber.
func BenchmarkShutdownTime(b *testing.B) {
	for _, n := range []int{1, 16, 64} {
		b.Run(subName("subscriptions", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				c := &Client{chanReg: newChanRegistry()}
				var wg sync.WaitGroup
				for j := 0; j < n; j++ {
					ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 1})
					_ = cancel
					wg.Add(1)
					go drain(ch, &wg)
				}
				c.chanReg.stopAll()
				wg.Wait()
			}
		})
	}
}

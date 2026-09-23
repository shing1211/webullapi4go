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

// Package resilience is deprecated: use
// [github.com/shing1211/webullapi4go/pkg/resilience/breaker],
// [github.com/shing1211/webullapi4go/pkg/resilience/ratelimit],
// [github.com/shing1211/webullapi4go/pkg/resilience/retry], and
// [github.com/shing1211/webullapi4go/pkg/resilience/clock].
package resilience

import (
	breakerpkg "github.com/shing1211/webullapi4go/pkg/resilience/breaker"
	clockpkg "github.com/shing1211/webullapi4go/pkg/resilience/clock"
	ratelimitpkg "github.com/shing1211/webullapi4go/pkg/resilience/ratelimit"
	retrypkg "github.com/shing1211/webullapi4go/pkg/resilience/retry"
)

type (
	Breaker = breakerpkg.Breaker
	Clock   = clockpkg.Clock
	Config  = breakerpkg.Config
	Fake    = clockpkg.Fake
	Keyed   = ratelimitpkg.Keyed
	Limiter = ratelimitpkg.Limiter
	Option  = breakerpkg.Option
	Retrier = retrypkg.Retrier
	State   = breakerpkg.State
)

const (
	DefaultCooldown       = breakerpkg.DefaultCooldown
	DefaultHalfOpenMax    = breakerpkg.DefaultHalfOpenMax
	DefaultMaxAttempts    = retrypkg.DefaultMaxAttempts
	DefaultMaxDelay       = retrypkg.DefaultMaxDelay
	DefaultRetryBaseDelay = retrypkg.DefaultBaseDelay
	DefaultRetryMaxDelay  = retrypkg.DefaultMaxDelay
	DefaultRetryAttempts  = retrypkg.DefaultMaxAttempts
	DefaultThreshold      = breakerpkg.DefaultThreshold
	StateClosed           = breakerpkg.StateClosed
	StateHalfOpen         = breakerpkg.StateHalfOpen
	StateOpen             = breakerpkg.StateOpen
)

var (
	ErrOpen = breakerpkg.ErrOpen
)

func NewBreaker(opts ...breakerpkg.Option) *breakerpkg.Breaker {
	return breakerpkg.New(opts...)
}
func NewKeyed(factory func() *ratelimitpkg.Limiter) *ratelimitpkg.Keyed {
	return ratelimitpkg.NewKeyed(factory)
}
func NewLimiter(rate float64, burst int, opts ...ratelimitpkg.Option) *ratelimitpkg.Limiter {
	return ratelimitpkg.New(rate, burst, opts...)
}
func NewRetrier(opts ...retrypkg.Option) *retrypkg.Retrier {
	return retrypkg.New(opts...)
}
func Permanent(err error) error {
	return retrypkg.Permanent(err)
}
func IsPermanent(err error) bool {
	return retrypkg.IsPermanent(err)
}
func DefaultIsRetryable(err error) bool {
	return retrypkg.DefaultIsRetryable(err)
}
func NewWithConfig(cfg retrypkg.Config) *retrypkg.Retrier {
	return retrypkg.NewWithConfig(cfg)
}
func SystemClock() clockpkg.Clock {
	return clockpkg.System()
}

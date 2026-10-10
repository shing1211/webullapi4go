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

// Package resilience groups the SDK's resilience primitives:
//
//   - [github.com/shing1211/webullapi4go/pkg/resilience/retry] provides
//     context-aware retries with exponential backoff and jitter.
//   - [github.com/shing1211/webullapi4go/pkg/resilience/ratelimit]
//     provides token-bucket rate limiters, including a keyed variant for
//     per-endpoint budgets.
//   - [github.com/shing1211/webullapi4go/pkg/resilience/breaker] provides
//     a circuit breaker.
//   - [github.com/shing1211/webullapi4go/pkg/resilience/clock] provides
//     the shared time source used for deterministic tests.
package resilience

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

package broker

// config is the resolved broker client configuration. It is populated by
// [Option] functions during [New] and treated as immutable afterwards.
type config struct{}

// defaultConfig returns the broker client defaults.
func defaultConfig() config { return config{} }

// Option mutates the [config] resolved by [New]. Options are applied in order,
// so a later option overrides an earlier one.
type Option func(*config)

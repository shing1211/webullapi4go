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

// Package broker provides the Webull Broker API HTTP client for HK.
//
// A Client is constructed from the public client, so signing, token handling,
// retries, and rate limiting are shared with the rest of the SDK:
//
//	cl, err := client.New(
//		client.WithAppKey(key),
//		client.WithAppSecret(secret),
//		client.WithSandbox(),
//	)
//	if err != nil {
//		return err
//	}
//	broker := broker.New(cl)
package broker

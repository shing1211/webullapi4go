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

// Package webull provides the public entry point for the Webull OpenAPI SDK.
//
// The recommended import path is github.com/shing1211/webullapi4go/webull:
//
//	import "github.com/shing1211/webullapi4go/webull"
//
// Create a client with webull.New:
//
//	cl, err := webull.New(
//	    webull.WithCredentials(appKey, appSecret),
//	    webull.WithRegion(webull.RegionHK),
//	)
//
// Obtain typed service clients from the returned Client:
//
//	mdCl := data.New(cl)         // market data
//	trCl := trade.New(cl)        // trading
//	stCl := stream.New(cl)       // MQTT streaming
//	evCl := events.New(cl)       // gRPC events
//
// All clients are safe for concurrent use.
//
// # Package-level re-exports
//
// The following types and constructors are re-exported from webull for
// convenience so callers need only one import:
//
//	Option                    from client
//	Region                    from client  (webull.RegionHK, webull.RegionUS)
//	Environment               from client  (webull.EnvProduction, webull.EnvSandbox)
//	Endpoints                 from client
//	WithCredentials           from client
//	WithAppKey, WithAppSecret from client
//	WithRegion, WithEnvironment from client
//	WithSandbox               from client
//	WithEndpoints, WithBaseURL from client
//	WithHTTPClient, WithTimeout from client
//	WithRetry, WithoutRetry    from client
//	WithRateLimiter, WithBreaker from client
//	NewRateLimiter, NewBreaker from client
package webull

import (
	"github.com/shing1211/webullapi4go/client"
)

type Option = client.Option

var WithCredentials = client.WithCredentials
var WithAppKey = client.WithAppKey
var WithAppSecret = client.WithAppSecret
var WithRegion = client.WithRegion
var WithEnvironment = client.WithEnvironment
var WithSandbox = client.WithSandbox
var WithEndpoints = client.WithEndpoints
var WithBaseURL = client.WithBaseURL
var WithHTTPClient = client.WithHTTPClient
var WithTimeout = client.WithTimeout
var WithUserAgent = client.WithUserAgent
var WithRetry = client.WithRetry
var WithoutRetry = client.WithoutRetry
var WithRateLimiter = client.WithRateLimiter
var WithBreaker = client.WithBreaker
var NewRateLimiter = client.NewRateLimiter
var NewBreaker = client.NewBreaker

type Region = client.Region

const (
	RegionHK Region = client.HK
	RegionUS Region = client.US
)

type Environment = client.Environment

const (
	EnvProduction Environment = client.Production
	EnvSandbox    Environment = client.Sandbox
)

type Endpoints = client.Endpoints

type Client = client.Client

func New(opts ...Option) (*Client, error) {
	return client.New(opts...)
}

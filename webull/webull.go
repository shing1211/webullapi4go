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

// Package webull provides optional aliases for selected core-client types and
// options. It is a convenience layer, not an aggregate service facade.
//
// Import the service packages directly when typed service clients are needed:
//
//	cl, err := webull.New(
//	    webull.WithCredentials(appKey, appSecret),
//	    webull.WithRegion(webull.RegionHK),
//	)
//	if err != nil {
//	    return err
//	}
//
//	market := data.New(cl)
//	trading := trade.New(cl)
//	streaming, err := stream.New(cl)
//	if err != nil {
//	    return err
//	}
//	ev, err := events.New(cl)
//	if err != nil {
//	    return err
//	}
//
// webull.New delegates to client.New and returns a *client.Client. The
// returned client does not construct or return data, trade, stream, or events
// service clients. The root service packages remain canonical.
//
// # Package-level aliases
//
// The following identifiers are aliases or thin delegations for the
// corresponding client constructors, types, and options:
//
//	Option                    from client
//	Region                    from client  (webull.RegionHK, webull.RegionUS)
//	Environment               from client  (webull.EnvProduction, webull.EnvSandbox)
//	Endpoints                 from client
//	Client                    from client
//	New                       from client
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

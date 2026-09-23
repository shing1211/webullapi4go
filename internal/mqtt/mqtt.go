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

// Package mqtt is a deprecated alias for [github.com/shing1211/webullapi4go/pkg/transport/mqtt].
// New code should import that package directly. This package will be removed in v3.
package mqtt

import (
	pkgmqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

type (
	ConnackError = pkgmqtt.ConnackError
	Message      = pkgmqtt.Message
	Handler      = pkgmqtt.Handler
	Config       = pkgmqtt.Config
	Client       = pkgmqtt.Client
)

const (
	DefaultKeepAlive            = pkgmqtt.DefaultKeepAlive
	DefaultConnectTimeout       = pkgmqtt.DefaultConnectTimeout
	DefaultWriteTimeout         = pkgmqtt.DefaultWriteTimeout
	DefaultPingTimeout          = pkgmqtt.DefaultPingTimeout
	DefaultMessageChannelDepth  = pkgmqtt.DefaultMessageChannelDepth
	DefaultMaxReconnectInterval = pkgmqtt.DefaultMaxReconnectInterval
	ConnackConnectionLimit      = pkgmqtt.ConnackConnectionLimit
)

var (
	ErrConnectionRefused = pkgmqtt.ErrConnectionRefused
	ErrConnectionLimit   = pkgmqtt.ErrConnectionLimit
)

func New(cfg Config) (*Client, error) {
	return pkgmqtt.New(cfg)
}

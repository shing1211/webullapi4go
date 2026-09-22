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

import (
	"testing"

	"github.com/shing1211/webullapi4go/client"
)

func TestNewWithNoOptions(t *testing.T) {
	t.Parallel()

	cl, err := client.New(
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}

	got := New(cl)
	if got == nil {
		t.Fatal("New(cl) returned nil")
	}
	if got.cfg != (config{}) {
		t.Errorf("config = %+v, want empty config{}", got.cfg)
	}
}

func TestNewWithNilOptionsIgnored(t *testing.T) {
	t.Parallel()

	cl, err := client.New(
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}

	got := New(cl, nil, nil)
	if got == nil {
		t.Fatal("New(cl, nil, nil) returned nil")
	}
	if got.cfg != (config{}) {
		t.Errorf("config = %+v, want empty config{}", got.cfg)
	}
}

func TestNewWithMultipleNilOptions(t *testing.T) {
	t.Parallel()

	cl, err := client.New(
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}

	got := New(cl, nil, nil, nil)
	if got == nil {
		t.Fatal("New(cl, nil, nil, nil) returned nil")
	}
	if got.cfg != (config{}) {
		t.Errorf("config = %+v, want empty config{}", got.cfg)
	}
}

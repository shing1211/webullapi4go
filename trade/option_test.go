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

package trade

import "testing"

func TestGuardrailOptionsResolve(t *testing.T) {
	t.Parallel()

	c := New(nil, WithMaxOrderNotional("2500.00"), WithMaxOrderQuantity("10"))
	if got, want := c.cfg.maxOrderNotional, "2500.00"; got != want {
		t.Errorf("maxOrderNotional = %q, want %q", got, want)
	}
	if got, want := c.cfg.maxOrderQuantity, "10"; got != want {
		t.Errorf("maxOrderQuantity = %q, want %q", got, want)
	}
}

func TestGuardrailOptionsDefaultOff(t *testing.T) {
	t.Parallel()

	c := New(nil)
	if c.cfg.maxOrderNotional != "" || c.cfg.maxOrderQuantity != "" {
		t.Errorf("defaults = %q/%q, want both empty", c.cfg.maxOrderNotional, c.cfg.maxOrderQuantity)
	}
}

func TestGuardrailOptionsEmptyDisables(t *testing.T) {
	t.Parallel()

	c := New(nil, WithMaxOrderNotional(""), WithMaxOrderQuantity(""))
	if c.cfg.maxOrderNotional != "" || c.cfg.maxOrderQuantity != "" {
		t.Errorf("empty values = %q/%q, want both empty", c.cfg.maxOrderNotional, c.cfg.maxOrderQuantity)
	}
}

func TestGuardrailOptionsAcceptZero(t *testing.T) {
	t.Parallel()

	c := New(nil, WithMaxOrderNotional("0"), WithMaxOrderQuantity("0.0"))
	if c.cfg.maxOrderNotional != "0" || c.cfg.maxOrderQuantity != "0.0" {
		t.Errorf("zero values = %q/%q, want 0/0.0", c.cfg.maxOrderNotional, c.cfg.maxOrderQuantity)
	}
}

func TestGuardrailOptionsPanicOnInvalid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opt  func() Option
	}{
		{"notional non-numeric", func() Option { return WithMaxOrderNotional("abc") }},
		{"notional negative", func() Option { return WithMaxOrderNotional("-1") }},
		{"notional NaN", func() Option { return WithMaxOrderNotional("NaN") }},
		{"notional infinity", func() Option { return WithMaxOrderNotional("Inf") }},
		{"quantity non-numeric", func() Option { return WithMaxOrderQuantity("x") }},
		{"quantity negative", func() Option { return WithMaxOrderQuantity("-0.5") }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", tc.name)
				}
			}()
			_ = tc.opt()
		})
	}
}

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

package money

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestMoneyJSONRoundTrip(t *testing.T) {
	var got Money
	if err := json.Unmarshal([]byte(`"12.3400"`), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.String() != "12.34" {
		t.Fatalf("Money.String() = %q, want %q", got.String(), "12.34")
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != `"12.34"` {
		t.Fatalf("Marshal() = %s, want %s", data, `"12.34"`)
	}
}

func TestMoneyUnmarshalAcceptsExistingNumericWireForms(t *testing.T) {
	cases := []struct {
		name string
		data string
		want string
	}{
		{name: "number", data: `1.25`, want: "1.25"},
		{name: "empty string", data: `""`, want: "0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got Money
			if err := json.Unmarshal([]byte(tc.data), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if got.String() != tc.want {
				t.Fatalf("Money.String() = %q, want %q", got.String(), tc.want)
			}
		})
	}
}

func TestMoneyRejectsInvalidValues(t *testing.T) {
	if _, err := NewFromString("not-a-number"); !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("NewFromString() error = %v, want ErrInvalidFormat", err)
	}
	var got Money
	if err := json.Unmarshal([]byte(`"not-a-number"`), &got); !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("Unmarshal() error = %v, want ErrInvalidFormat", err)
	}
}

func TestMoneyArithmetic(t *testing.T) {
	a := MustNew("10.25")
	b := MustNew("2.5")
	if got := a.Add(b).String(); got != "12.75" {
		t.Fatalf("Add() = %q, want %q", got, "12.75")
	}
	if got := a.Sub(b).String(); got != "7.75" {
		t.Fatalf("Sub() = %q, want %q", got, "7.75")
	}
	if got := a.Mul(b).String(); got != "25.625" {
		t.Fatalf("Mul() = %q, want %q", got, "25.625")
	}
	if got := a.Div(b).String(); got != "4.1" {
		t.Fatalf("Div() = %q, want %q", got, "4.1")
	}
	rounded, err := a.DivRound(b, 2)
	if err != nil {
		t.Fatalf("DivRound() error = %v", err)
	}
	if got := rounded.String(); got != "4.1" {
		t.Fatalf("DivRound() = %q, want %q", got, "4.1")
	}
	if got := a.Div(Zero()); !got.IsZero() {
		t.Fatalf("Div() by zero = %q, want zero", got.String())
	}
	if _, err := a.DivRound(Zero(), 2); err == nil {
		t.Fatal("DivRound() by zero error = nil")
	}
}

func TestMoneyPredicatesAndRat(t *testing.T) {
	zero := Zero()
	positive := MustNew("0.1")
	negative := MustNew("-0.1")
	if !zero.IsZero() || zero.IsPositive() || zero.IsNegative() {
		t.Fatalf("zero predicates = zero:%t positive:%t negative:%t", zero.IsZero(), zero.IsPositive(), zero.IsNegative())
	}
	if !positive.IsPositive() || positive.IsNegative() {
		t.Fatal("positive predicates are incorrect")
	}
	if !negative.IsNegative() || negative.IsPositive() {
		t.Fatal("negative predicates are incorrect")
	}
	want := positive.Rat().RatString()
	if got := positive.Rat().RatString(); got != want {
		t.Fatalf("Rat().RatString() = %q, want %q", got, want)
	}
}

func TestMoneyConstructorsAndParseDecimal(t *testing.T) {
	if got := NewFromInt64(42).String(); got != "42" {
		t.Fatalf("NewFromInt64() = %q, want %q", got, "42")
	}
	if got, ok := ParseMoney("3.5"); !ok || got.String() != "3.5" {
		t.Fatalf("ParseMoney() = %q, %t", got.String(), ok)
	}
	if _, ok := ParseMoney("bad"); ok {
		t.Fatal("ParseMoney() accepted invalid input")
	}
	if _, ok := ParseDecimal("3/4"); !ok {
		t.Fatal("ParseDecimal() rejected a valid rational")
	}
	if _, ok := ParseDecimal("not-a-number"); ok {
		t.Fatal("ParseDecimal() accepted invalid input")
	}
}

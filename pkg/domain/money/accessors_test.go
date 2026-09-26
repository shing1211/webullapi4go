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
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewWrapsDecimal(t *testing.T) {
	t.Parallel()

	m := New(decimal.RequireFromString("12.34"))
	if got := m.String(); got != "12.34" {
		t.Fatalf("New().String() = %q, want %q", got, "12.34")
	}
	if got := m.Decimal().String(); got != "12.34" {
		t.Fatalf("Decimal() = %q, want %q", got, "12.34")
	}
	if !m.Decimal().Equal(decimal.RequireFromString("12.34")) {
		t.Fatal("Decimal() did not round-trip the wrapped value")
	}
}

func TestMustNewPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	if got := MustNew("1.5").String(); got != "1.5" {
		t.Fatalf("MustNew(1.5) = %q, want %q", got, "1.5")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustNew(\"not-a-number\") did not panic")
		}
		if !errors.Is(r.(error), ErrInvalidFormat) {
			t.Fatalf("MustNew panic = %v, want ErrInvalidFormat", r)
		}
	}()
	MustNew("not-a-number")
}

func TestNewFromStringEmptyIsZeroValue(t *testing.T) {
	t.Parallel()

	got, err := NewFromString("")
	if err != nil {
		t.Fatalf("NewFromString(\"\") error = %v, want nil", err)
	}
	if !got.IsZero() {
		t.Fatalf("NewFromString(\"\") = %q, want zero", got)
	}

	if _, err := NewFromString("abc"); !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("NewFromString(abc) error = %v, want ErrInvalidFormat", err)
	}
}

func TestMustPanicsOnError(t *testing.T) {
	t.Parallel()

	if got := Must(MustNew("2.5"), nil).String(); got != "2.5" {
		t.Fatalf("Must(value, nil) = %q, want %q", got, "2.5")
	}

	// A nil error must not panic even for the zero value.
	if !Must(Money{}, nil).IsZero() {
		t.Fatal("Must(zero, nil) is not zero")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Must(value, err) did not panic")
		}
	}()
	Must(Money{}, errors.New("boom"))
}

func TestZeroIsZero(t *testing.T) {
	t.Parallel()

	if got := Zero().String(); got != "0" {
		t.Fatalf("Zero() = %q, want %q", got, "0")
	}
	if !Zero().IsZero() {
		t.Fatal("Zero().IsZero() = false, want true")
	}
}

func TestComparisonOperators(t *testing.T) {
	t.Parallel()

	small := MustNew("1.00")
	big := MustNew("2.00")
	same := MustNew("1.00")

	cases := []struct {
		name string
		got  bool
		want bool
	}{
		{"small Cmp big is negative", small.Cmp(big) < 0, true},
		{"big Cmp small is positive", big.Cmp(small) > 0, true},
		{"equal Cmp is zero", small.Cmp(same) == 0, true},
		{"small GreaterThan big", small.GreaterThan(big), false},
		{"big GreaterThan small", big.GreaterThan(small), true},
		{"equal GreaterThan", same.GreaterThan(small), false},
		{"big GreaterThanOrEqual small", big.GreaterThanOrEqual(small), true},
		{"equal GreaterThanOrEqual", same.GreaterThanOrEqual(small), true},
		{"small LessThan big", small.LessThan(big), true},
		{"equal LessThan", same.LessThan(small), false},
		{"small LessThanOrEqual big", small.LessThanOrEqual(big), true},
		{"equal LessThanOrEqual", same.LessThanOrEqual(small), true},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
}

func TestAbsAndNeg(t *testing.T) {
	t.Parallel()

	if got := MustNew("-4.25").Abs().String(); got != "4.25" {
		t.Errorf("Abs() = %q, want %q", got, "4.25")
	}
	if got := MustNew("4.25").Abs().String(); got != "4.25" {
		t.Errorf("Abs() of a positive = %q, want %q", got, "4.25")
	}
	if got := MustNew("4.25").Neg().String(); got != "-4.25" {
		t.Errorf("Neg() = %q, want %q", got, "-4.25")
	}
	if got := MustNew("-4.25").Neg().String(); got != "4.25" {
		t.Errorf("Neg() of a negative = %q, want %q", got, "4.25")
	}
}

func TestRoundAndShift(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		got   string
		want  string
		shift int32
		round int32
		doRnd bool
	}{
		{name: "round half up to 2", got: MustNew("1.005").Round(2).String(), want: "1.01", round: 2, doRnd: true},
		{name: "round down to 1", got: MustNew("1.004").Round(1).String(), want: "1", round: 1, doRnd: true},
		{name: "round to 0 places", got: MustNew("2.5").Round(0).String(), want: "3", round: 0, doRnd: true},
		{name: "shift up", got: MustNew("1.5").Shift(2).String(), want: "150", shift: 2},
		{name: "shift down", got: MustNew("150").Shift(-2).String(), want: "1.5", shift: -2},
		{name: "shift zero", got: MustNew("1.5").Shift(0).String(), want: "1.5", shift: 0},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

func TestBigIntTruncatesTowardZero(t *testing.T) {
	t.Parallel()

	got := MustNew("12.99").BigInt()
	if got.Cmp(big.NewInt(12)) != 0 {
		t.Fatalf("BigInt() = %v, want 12", got)
	}

	if got := MustNew("-12.99").BigInt(); got.Cmp(big.NewInt(-12)) != 0 {
		t.Fatalf("BigInt() of a negative = %v, want -12", got)
	}
}

func TestRatIsExact(t *testing.T) {
	t.Parallel()

	got := MustNew("0.1").Rat()
	want, ok := new(big.Rat).SetString("1/10")
	if !ok {
		t.Fatal("could not build the expected rational")
	}
	if got.Cmp(want) != 0 {
		t.Fatalf("Rat() = %v, want %v", got, want)
	}
}

func TestNewFromInt64(t *testing.T) {
	t.Parallel()

	if got := NewFromInt64(42).String(); got != "42" {
		t.Fatalf("NewFromInt64(42) = %q, want %q", got, "42")
	}
	if got := NewFromInt64(-7).String(); got != "-7" {
		t.Fatalf("NewFromInt64(-7) = %q, want %q", got, "-7")
	}
	if got := NewFromInt64(0); !got.IsZero() {
		t.Fatalf("NewFromInt64(0) = %q, want zero", got)
	}
}

func TestParseMoney(t *testing.T) {
	t.Parallel()

	if m, ok := ParseMoney("3.5"); !ok || m.String() != "3.5" {
		t.Fatalf("ParseMoney(3.5) = (%q, %v), want (3.5, true)", m, ok)
	}
	if m, ok := ParseMoney(""); !ok || !m.IsZero() {
		t.Fatalf("ParseMoney(\"\") = (%q, %v), want (zero, true)", m, ok)
	}
	if m, ok := ParseMoney("nope"); ok {
		t.Fatalf("ParseMoney(nope) = (%q, %v), want ok=false", m, ok)
	}
}

func TestUnmarshalJSONRejectsNonFinite(t *testing.T) {
	t.Parallel()

	// A JSON string that is not a decimal must surface ErrInvalidFormat rather
	// than a partially populated value.
	var m Money
	if err := m.UnmarshalJSON([]byte(`"not-a-number"`)); !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("UnmarshalJSON(not-a-number) error = %v, want ErrInvalidFormat", err)
	}
}

func TestMoneyJSONRoundTripThroughStruct(t *testing.T) {
	t.Parallel()

	type payload struct {
		Price Money `json:"price"`
		Size  Money `json:"size"`
	}

	in := payload{Price: MustNew("187.25"), Size: NewFromInt64(3)}
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if want := `{"price":"187.25","size":"3"}`; string(raw) != want {
		t.Fatalf("Marshal() = %s, want %s", raw, want)
	}

	var out payload
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if out.Price.String() != "187.25" || out.Size.String() != "3" {
		t.Fatalf("round trip = %+v, want price 187.25 and size 3", out)
	}
}

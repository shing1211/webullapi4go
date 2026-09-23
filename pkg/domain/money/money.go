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
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidFormat = errors.New("money: invalid decimal format")
)

type Money struct {
	d decimal.Decimal
}

func New(d decimal.Decimal) Money {
	return Money{d: d}
}

func MustNew(s string) Money {
	m, err := NewFromString(s)
	if err != nil {
		panic(err)
	}
	return m
}

func NewFromString(s string) (Money, error) {
	if s == "" {
		return Money{}, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, ErrInvalidFormat
	}
	return Money{d: d}, nil
}

func NewFromInt64(i int64) Money {
	return Money{d: decimal.NewFromInt(i)}
}

func ParseMoney(s string) (Money, bool) {
	m, err := NewFromString(s)
	return m, err == nil
}

func Must(m Money, err error) Money {
	if err != nil {
		panic(err)
	}
	return m
}

func Zero() Money {
	return Money{d: decimal.Zero}
}

func (m Money) Decimal() decimal.Decimal {
	return m.d
}

func (m Money) String() string {
	return m.d.String()
}

func (m Money) IsZero() bool {
	return m.d.IsZero()
}

func (m Money) IsPositive() bool {
	return m.d.IsPositive()
}

func (m Money) IsNegative() bool {
	return m.d.IsNegative()
}

func (m Money) Add(o Money) Money {
	return Money{d: m.d.Add(o.d)}
}

func (m Money) Sub(o Money) Money {
	return Money{d: m.d.Sub(o.d)}
}

func (m Money) Mul(o Money) Money {
	return Money{d: m.d.Mul(o.d)}
}

func (m Money) Div(o Money) Money {
	if o.d.IsZero() {
		return Money{d: decimal.Zero}
	}
	return Money{d: m.d.Div(o.d)}
}

func (m Money) DivRound(o Money, precision int32) (Money, error) {
	if o.d.IsZero() {
		return Money{d: decimal.Zero}, errors.New("money: division by zero")
	}
	rounded := m.d.DivRound(o.d, precision)
	return Money{d: rounded}, nil
}

func (m Money) Cmp(o Money) int {
	return m.d.Cmp(o.d)
}

func (m Money) GreaterThan(o Money) bool {
	return m.d.GreaterThan(o.d)
}

func (m Money) GreaterThanOrEqual(o Money) bool {
	return m.d.GreaterThanOrEqual(o.d)
}

func (m Money) LessThan(o Money) bool {
	return m.d.LessThan(o.d)
}

func (m Money) LessThanOrEqual(o Money) bool {
	return m.d.LessThanOrEqual(o.d)
}

func (m Money) Abs() Money {
	return Money{d: m.d.Abs()}
}

func (m Money) Neg() Money {
	return Money{d: m.d.Neg()}
}

func (m Money) Round(precision int32) Money {
	return Money{d: m.d.Round(precision)}
}

func (m Money) Shift(shift int32) Money {
	return Money{d: m.d.Shift(shift)}
}

func (m Money) BigInt() *big.Int {
	return m.d.BigInt()
}

func (m Money) Rat() *big.Rat {
	return m.d.Rat()
}

func ParseDecimal(s string) (*big.Rat, bool) {
	r, ok := new(big.Rat).SetString(strings.TrimSpace(s))
	return r, ok
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.d.String())
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var f float64
		if err2 := json.Unmarshal(data, &f); err2 == nil {
			m.d = decimal.NewFromFloat(f)
			return nil
		}
		return err
	}
	if s == "" {
		m.d = decimal.Zero
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return ErrInvalidFormat
	}
	m.d = d
	return nil
}

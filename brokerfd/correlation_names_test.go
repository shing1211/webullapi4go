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

package brokerfd

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// This file covers the three names added to the Broker FD response types in v2.1.30:
// client_request_id, create_time and update_time.
//
// The reason they were added is not that a page describes them. 28 recorded
// divergences were, all of the same shape: the page marks the name required, the type
// carries no field for it, and a caller reading it got the empty string with nothing
// in the return path to distinguish that from a genuinely empty value. It is the
// defect class item 20 recorded, and the same additive fix applies because a
// *required* name is one a conforming server always sends - so adding a field cannot
// break a call that works today, and no probe is needed to justify it.
//
// The names were chosen because they recur: client_request_id on 11 symbols,
// update_time on 11, create_time on 6, for 28 of the 119 rows in that class from 19
// field additions. client_request_id already appeared in the sibling broker/ module
// in this exact role, so the spelling follows it rather than being invented here.

// TestRequiredCorrelationNamesDecode pins the new fields against literal bodies, so
// a tag typo or a wrong spelling is caught here rather than by the conformance gate
// noticing a name went missing again.
func TestRequiredCorrelationNamesDecode(t *testing.T) {
	const body = `{
		"client_request_id": "LJIS16BACHQG9LPP44L9IQHGAB",
		"create_time": "1761131409276",
		"update_time": "1761131499276"
	}`

	tests := []struct {
		name   string
		target any
		fields []string
	}{
		{"ACHAccount", &ACHAccount{}, []string{"ClientRequestID", "CreateTime", "UpdateTime"}},
		{"BankAccount", &BankAccount{}, []string{"ClientRequestID", "CreateTime", "UpdateTime"}},
		{"CreditInfo", &CreditInfo{}, []string{"ClientRequestID", "CreateTime", "UpdateTime"}},
		{"TransferFee", &TransferFee{}, []string{"ClientRequestID", "CreateTime", "UpdateTime"}},
		{"InstantFunding", &InstantFunding{}, []string{"ClientRequestID", "UpdateTime"}},
		{"Transfer", &Transfer{}, []string{"ClientRequestID", "UpdateTime"}},
		{"FDCashJournal", &FDCashJournal{}, []string{"ClientRequestID", "UpdateTime"}},
		{"FDAccount", &FDAccount{}, []string{"ClientRequestID"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(body), tc.target); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			v := reflect.ValueOf(tc.target).Elem()
			want := map[string]string{
				"ClientRequestID": "LJIS16BACHQG9LPP44L9IQHGAB",
				"CreateTime":      "1761131409276",
				"UpdateTime":      "1761131499276",
			}
			for _, name := range tc.fields {
				f := v.FieldByName(name)
				if !f.IsValid() {
					t.Errorf("%s has no field %s", tc.name, name)
					continue
				}
				if got := f.String(); got != want[name] {
					t.Errorf("%s.%s = %q, want %q: the required name reached no field",
						tc.name, name, got, want[name])
				}
			}
		})
	}
}

// TestCorrelationNamesAreDistinctWireNames guards the tag spelling against the class
// v2.1.28 recorded: two fields of one struct claiming one json name makes
// encoding/json drop both, silently. These types already carry create_time on
// FDAccount, FDCashJournal, InstantFunding and Transfer, so the check is that the
// struct still has one field per wire name.
//
// The required set is per type, because the pages differ: FDAccount requires only
// client_request_id, and asserting the union of all three names on every type would
// fail on a type that is correct.
func TestCorrelationNamesAreDistinctWireNames(t *testing.T) {
	tests := []struct {
		target   any
		required []string
	}{
		{&ACHAccount{}, []string{"client_request_id", "create_time", "update_time"}},
		{&BankAccount{}, []string{"client_request_id", "create_time", "update_time"}},
		{&CreditInfo{}, []string{"client_request_id", "create_time", "update_time"}},
		{&TransferFee{}, []string{"client_request_id", "create_time", "update_time"}},
		{&InstantFunding{}, []string{"client_request_id", "update_time"}},
		{&Transfer{}, []string{"client_request_id", "update_time"}},
		{&FDCashJournal{}, []string{"client_request_id", "update_time"}},
		{&FDAccount{}, []string{"client_request_id"}},
	}
	for _, tc := range tests {
		rt := reflect.TypeOf(tc.target).Elem()
		name := rt.Name()
		seen := map[string]string{}
		for i := range rt.NumField() {
			f := rt.Field(i)
			raw, ok := f.Tag.Lookup("json")
			if !ok {
				continue
			}
			if k := strings.IndexByte(raw, ','); k >= 0 {
				raw = raw[:k]
			}
			if raw == "" || raw == "-" {
				continue
			}
			if prev, dup := seen[raw]; dup {
				t.Errorf("%s: fields %s and %s both encode under %q, so encoding/json "+
					"drops both", name, prev, f.Name, raw)
				continue
			}
			seen[raw] = f.Name
		}
		for _, required := range tc.required {
			if _, ok := seen[required]; !ok {
				t.Errorf("%s carries no field for %s", name, required)
			}
		}
	}
}

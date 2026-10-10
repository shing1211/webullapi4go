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

package region

import "testing"

func TestRegionString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		region Region
		want   string
	}{
		{HK, "hk"},
		{US, "us"},
		{JP, "jp"},
		{SG, "sg"},
		{TH, "th"},
		{AU, "au"},
		{MY, "my"},
		{UK, "uk"},
		{BR, "br"},
		{MX, "mx"},
		{ZA, "za"},
		{EU, "eu"},
		{Region(""), ""},
		{Region("nowhere"), "nowhere"},
	}

	for _, tc := range cases {
		if got := tc.region.String(); got != tc.want {
			t.Errorf("Region(%q).String() = %q, want %q", tc.region, got, tc.want)
		}
	}
}

func TestRegionIsValid(t *testing.T) {
	t.Parallel()

	for _, r := range []Region{HK, US, JP, SG, TH, AU, MY, UK, BR, MX, ZA, EU} {
		if !r.IsValid() {
			t.Errorf("Region(%q).IsValid() = false, want true", r)
		}
	}

	for _, r := range []Region{"", "nowhere", "HK"} {
		if r.IsValid() {
			t.Errorf("Region(%q).IsValid() = true, want false", r)
		}
	}
}

func TestParseRegion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in     string
		want   Region
		wantOK bool
	}{
		{"hk", HK, true},
		{"HK", HK, true},
		{"  US  ", US, true},
		{"Jp", JP, true},
		{"EU", EU, true},
		{"", Region(""), false},
		{"   ", Region(""), false},
		{"nowhere", Region("nowhere"), false},
	}

	for _, tc := range cases {
		got, ok := ParseRegion(tc.in)
		if ok != tc.wantOK {
			t.Errorf("ParseRegion(%q) ok = %v, want %v", tc.in, ok, tc.wantOK)
		}
		if got != tc.want {
			t.Errorf("ParseRegion(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestEnvironmentString(t *testing.T) {
	t.Parallel()

	if got := Production.String(); got != "production" {
		t.Errorf("Production.String() = %q, want %q", got, "production")
	}
	if got := Sandbox.String(); got != "sandbox" {
		t.Errorf("Sandbox.String() = %q, want %q", got, "sandbox")
	}
	if got := Environment("").String(); got != "" {
		t.Errorf("empty Environment.String() = %q, want empty", got)
	}
}

func TestEnvironmentIsValid(t *testing.T) {
	t.Parallel()

	if !Production.IsValid() {
		t.Error("Production.IsValid() = false, want true")
	}
	if !Sandbox.IsValid() {
		t.Error("Sandbox.IsValid() = false, want true")
	}
	for _, e := range []Environment{"", "uat", "Production", "SANDBOX", "staging"} {
		if e.IsValid() {
			t.Errorf("Environment(%q).IsValid() = true, want false", e)
		}
	}
}

// TestDomainForEveryKnownRegion guards that each documented region still has a
// root domain, because endpoints are derived from it.
func TestDomainForEveryKnownRegion(t *testing.T) {
	t.Parallel()

	want := map[Region]string{
		HK: "webull.hk",
		US: "webull.com",
		JP: "webull.co.jp",
		SG: "webull.com.sg",
		TH: "webull.co.th",
		AU: "webull.com.au",
		MY: "webull.com.my",
		UK: "webull-uk.com",
		BR: "webull.com.br",
		MX: "webull.com.mx",
		ZA: "webull.co.za",
		EU: "webull.eu",
	}

	for region, wantDomain := range want {
		got, ok := Domain(region)
		if !ok {
			t.Errorf("Domain(%q) ok = false, want true", region)
			continue
		}
		if got != wantDomain {
			t.Errorf("Domain(%q) = %q, want %q", region, got, wantDomain)
		}
	}

	if got, ok := Domain("nowhere"); ok || got != "" {
		t.Errorf("Domain(nowhere) = (%q, %v), want (empty, false)", got, ok)
	}
}

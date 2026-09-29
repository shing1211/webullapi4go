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

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// --------------------------------------------------------------------------
// The guard: a capture may only write under live/
// --------------------------------------------------------------------------

// A skeleton key that resolves into the documentation fixture tree is refused,
// because overwriting a committed fixture would silently change the existing
// gate rather than add evidence to it. The key below is a real documentation
// fixture id, so the case is the concrete one rather than a hypothetical path:
// conformance/testdata/manifest.json names it, gen_fixtures.py --check regenerates
// it from Webull's published schema, and conformance.CompareBody reads it.
func TestWriteCaptureRefusesToTouchDocumentationFixtures(t *testing.T) {
	dir := t.TempDir()
	// A fixture id that resolves into the documentation fixture tree must be rejected,
	// because overwriting a committed fixture would silently change the existing gate.
	err := WriteCapture(dir, map[string][]byte{
		"broker-fd-us/GET-broker-accounts-get.json": []byte(`{}`),
	}, nil)
	if err == nil {
		t.Fatal("WriteCapture accepted a documentation-fixture path")
	}
	if _, statErr := os.Stat(filepath.Join(dir, "broker-fd-us")); statErr == nil {
		t.Fatal("the refused key created a directory anyway")
	}
}

// The other ways out of live/ are refused too, and the corpus is the set of
// shapes a key can take: a path separator from another platform, an absolute
// path, a key that has to be cleaned, and the live directory named as a file.
func TestWriteCaptureRefusesEveryKeyOutsideTheLiveTree(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{"documentation fixture", "trading/GET-trading-orders-get.json"},
		{"the manifest itself", "manifest.json"},
		{"the documentation manifest", "manifest.json.bak"},
		{"top level", "trading.json"},
		{"backslash separator", `live\trading\x.json`},
		{"absolute path", "/etc/passwd"},
		{"absolute windows path", `C:\conformance\testdata\trading\x.json`},
		{"parent escape", "live/../trading/GET-trading-orders-get.json"},
		{"deep parent escape", "live/a/../../trading/x.json"},
		{"leading current directory", "./live/trading/x.json"},
		{"doubled separator", "live//trading/x.json"},
		{"trailing separator", "live/trading/"},
		{"the live directory itself", "live/"},
		{"the live directory as a bare name", "live"},
		{"empty", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			err := WriteCapture(dir, map[string][]byte{tc.key: []byte(`{}`)}, nil)
			if err == nil {
				t.Fatalf("WriteCapture accepted %q", tc.key)
			}
			// A refusal must leave the tree exactly as it found it, so the check
			// is that nothing at all was created: a run that half-writes before
			// refusing has already done the damage the guard exists to prevent.
			assertTreeEmpty(t, dir)
		})
	}
}

// A refusal is a refusal of the whole run, not of one file. The corpus has one
// good key and one bad one, and the good one must not be written either: a
// capture that wrote what it could and then failed leaves the tree in a state no
// manifest describes.
func TestWriteCaptureRefusesTheRunNotTheFile(t *testing.T) {
	dir := t.TempDir()
	err := WriteCapture(dir, map[string][]byte{
		"live/trading/GET-trading-orders-get.json": []byte("{}\n"),
		"trading/GET-trading-orders-get.json":      []byte("{}\n"),
	}, nil)
	if err == nil {
		t.Fatal("WriteCapture accepted a run holding a documentation-fixture path")
	}
	assertTreeEmpty(t, dir)
}

// assertTreeEmpty fails if dir holds any file, at any depth.
func assertTreeEmpty(t *testing.T, dir string) {
	t.Helper()
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			t.Errorf("a refused write left %s behind", p)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", dir, err)
	}
}

// --------------------------------------------------------------------------
// What WriteCapture writes
// --------------------------------------------------------------------------

// The exact bytes a skeleton reaches the repository with, pinned as a literal.
//
// This is the test that matters most in the file, because a skeleton is a
// committed artefact and its content is the whole claim. Three things are fixed
// at once: every value is a placeholder, the members are in sorted order, and
// the file is one-space indented with a trailing newline. The first is the
// invariant the reducer exists to provide; the second is what makes a re-run
// byte-identical; the third is what makes the file readable next to
// conformance/testdata/manifest.json. A change to any of the three is a change to
// what a committed file means.
func TestWriteCaptureWritesTheSkeletonVerbatim(t *testing.T) {
	dir := t.TempDir()
	key := "live/trading/GET-trading-orders-get.json"
	entry := ManifestEntry{
		Symbol:         "trade.GetOrderDetail",
		Fixture:        "trading/GET-trading-orders-get.json",
		Skeleton:       key,
		Status:         200,
		DecodedCleanly: true,
		Host:           "api.sandbox.webull.hk",
		ProbedAt:       "2026-09-29T04:05:06Z",
	}
	skeleton := []byte("{\n \"order_id\": \"1\",\n \"quantity\": -1,\n \"symbol\": \"1\",\n" +
		" \"filled_quantity\": null\n}\n")

	if err := WriteCapture(dir, map[string][]byte{key: skeleton}, []ManifestEntry{entry}); err != nil {
		t.Fatalf("WriteCapture: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(key))) //nolint:gosec // G304: the key is the test's own, and dir is t.TempDir().
	if err != nil {
		t.Fatalf("reading the skeleton: %v", err)
	}
	if !bytes.Equal(got, skeleton) {
		t.Errorf("skeleton on disk:\n%s\nwant:\n%s", got, skeleton)
	}
}

// A real capture writes the bytes Skeletonify and encodeSkeleton produce, so the
// literal above is checked against the pipeline rather than trusted.
func TestCaptureWritesTheReducedTreeAndNothingElse(t *testing.T) {
	// A body carrying every kind at depth, a name the SDK does not have, and the
	// figures a leak would look like.
	const body = `{"order_id":"0352U72LQI6DT0KF41GK000000","total_quantity":"12.5",` +
		`"filled_quantity":100,"close":385.6,"place_time_at":1756000000000,` +
		`"ratio":null,"active":false,"lots":[{"symbol":"AAPL","price":"1.0"}],` +
		`"a_name_the_sdk_does_not_carry":"sentinel-value","empty":[]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	dir := t.TempDir()
	skeleton, entry, err := captureThroughServer(t, srv, orderEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if entry.Skeleton == "" {
		t.Fatalf("nothing captured: %s", entry.NotCaptured)
	}
	if err := WriteCapture(dir, map[string][]byte{entry.Skeleton: skeleton}, []ManifestEntry{entry}); err != nil {
		t.Fatalf("WriteCapture: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(entry.Skeleton))) //nolint:gosec // G304: the key came from the capture under test and dir is t.TempDir().
	if err != nil {
		t.Fatalf("reading the skeleton: %v", err)
	}

	const want = "{\n \"a_name_the_sdk_does_not_carry\": \"1\",\n \"active\": true,\n" +
		" \"close\": -1,\n \"empty\": [],\n \"filled_quantity\": -1,\n" +
		" \"lots\": [\n  {\n   \"price\": \"1\",\n   \"symbol\": \"1\"\n  }\n ],\n" +
		" \"order_id\": \"1\",\n \"place_time_at\": -1,\n \"ratio\": null,\n" +
		" \"total_quantity\": \"1\"\n}\n"
	if string(raw) != want {
		t.Errorf("committed skeleton:\n%s\nwant:\n%s", raw, want)
	}

	// The other half of the same claim, on the bytes rather than on the literal:
	// every leaf in the file on disk is a placeholder, and none of the server's
	// figures is anywhere in it.
	assertNoLiveValueInFile(t, raw)
	for _, leaked := range []string{
		"0352U72LQI6DT0KF41GK000000", "12.5", "100", "385.6", "1756000000000",
		"AAPL", "sentinel-value",
	} {
		if bytes.Contains(raw, []byte(leaked)) {
			t.Errorf("the committed skeleton contains the server literal %q", leaked)
		}
	}
	// The member name is evidence and must survive: stripping it would empty
	// the evidence set as silently as keeping a value would pollute it.
	if !bytes.Contains(raw, []byte("a_name_the_sdk_does_not_carry")) {
		t.Error("the committed skeleton lost a member name")
	}
}

// The manifest a reader relies on has to name the file, the SDK symbol, the
// status, the host and the time, and has to say in its own text that it holds no
// values. Those five facts are the difference between a committed file a reader
// can judge and one they can only hope at.
func TestLiveManifestStatesWhatItHoldsAndFromWhere(t *testing.T) {
	dir := t.TempDir()
	key := "live/trading/GET-trading-orders-get.json"
	entries := []ManifestEntry{
		{
			Symbol: "trade.GetOrderDetail", Fixture: "trading/GET-trading-orders-get.json",
			Skeleton: key, Status: 200, DecodedCleanly: true,
			Host: "api.sandbox.webull.hk", ProbedAt: "2026-09-29T04:05:06Z",
		},
		{
			Symbol: "trade.GetPositions", Fixture: "trading/GET-trading-assets-positions-list.json",
			Status: 200, NotCaptured: "the 200 response declares Content-Type text/html",
			Host: "api.sandbox.webull.hk", ProbedAt: "2026-09-29T04:05:06Z",
		},
	}
	if err := WriteCapture(dir, map[string][]byte{key: []byte("{}\n")}, entries); err != nil {
		t.Fatalf("WriteCapture: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, liveManifestName)) //nolint:gosec // G304: a fixed name under t.TempDir().
	if err != nil {
		t.Fatalf("reading %s: %v", liveManifestName, err)
	}

	var doc LiveManifest
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parsing %s: %v", liveManifestName, err)
	}
	if doc.Kind != liveKind {
		t.Errorf("kind = %q, want %q", doc.Kind, liveKind)
	}
	if doc.Generator.ContainsValues {
		t.Error("containsValues is true, so the manifest claims to hold readings")
	}
	if !strings.Contains(doc.Warning, "no value it sent") {
		t.Errorf("the warning does not say the tree holds no values: %q", doc.Warning)
	}
	if !strings.Contains(doc.Scope, "HTTP 200") {
		t.Errorf("the scope does not name the status that was captured: %q", doc.Scope)
	}
	if doc.Host != "api.sandbox.webull.hk" {
		t.Errorf("host = %q, want the host every row recorded", doc.Host)
	}
	// Entries are sorted by fixture, so the file is a function of the run rather
	// than of the order the walk happened to produce.
	if len(doc.Entries) != 2 || doc.Entries[0].Fixture != "trading/GET-trading-assets-positions-list.json" {
		t.Errorf("entries are not sorted by fixture: %+v", doc.Entries)
	}
	// The counts are derived from the rows, and the two partitions must add up.
	if got, want := doc.Totals.Considered, 2; got != want {
		t.Errorf("totals.considered = %d, want %d", got, want)
	}
	if got, want := doc.Totals.Captured, 1; got != want {
		t.Errorf("totals.captured = %d, want %d", got, want)
	}
	if got, want := doc.Totals.NotCaptured, 1; got != want {
		t.Errorf("totals.notCaptured = %d, want %d", got, want)
	}
	if got, want := doc.Totals.DecodedCleanly+doc.Totals.DecodeRejected, doc.Totals.Considered; got != want {
		t.Errorf("decoded %d + rejected %d != considered %d", doc.Totals.DecodedCleanly,
			doc.Totals.DecodeRejected, doc.Totals.Considered)
	}
	if len(doc.Totals.NotCapturedReasons) != 1 {
		t.Errorf("totals.notCapturedReasons = %v, want the one reason", doc.Totals.NotCapturedReasons)
	}
}

// The manifest must be a function of the run and not of Go's map iteration, so
// two writes of the same set produce the same bytes. Without this a re-run
// produces a diff nobody can read, which is the same failure as a diff nobody
// can trust.
func TestWriteCaptureIsByteIdenticalAcrossRuns(t *testing.T) {
	build := func() map[string][]byte {
		// Nine members: more than the eight a hash-ordered map iteration is
		// unlikely to align on, so a map that leaked its order into the output
		// would show up here rather than once in a hundred runs.
		reduced := map[string]any{}
		for i := 0; i < 9; i++ {
			reduced["m"+strconv.Itoa(i)] = map[string]any{"n": "1", "v": -1}
		}
		encoded, err := encodeSkeleton(reduced)
		if err != nil {
			t.Fatalf("encodeSkeleton: %v", err)
		}
		return map[string][]byte{"live/fundamentals/x.json": encoded}
	}
	entries := []ManifestEntry{{
		Symbol: "data.GetX", Fixture: "fundamentals/x.json", Skeleton: "live/fundamentals/x.json",
		Status: 200, DecodedCleanly: true, Host: "api.sandbox.webull.hk",
		ProbedAt: "2026-09-29T04:05:06Z",
	}}

	read := func() (string, string) {
		dir := t.TempDir()
		if err := WriteCapture(dir, build(), entries); err != nil {
			t.Fatalf("WriteCapture: %v", err)
		}
		skeleton, err := os.ReadFile(filepath.Join(dir, "live", "fundamentals", "x.json")) //nolint:gosec // G304: a fixed name under t.TempDir().
		if err != nil {
			t.Fatalf("reading the skeleton: %v", err)
		}
		manifest, err := os.ReadFile(filepath.Join(dir, liveManifestName)) //nolint:gosec // G304: a fixed name under t.TempDir().
		if err != nil {
			t.Fatalf("reading the manifest: %v", err)
		}
		return string(skeleton), string(manifest)
	}
	firstSkeleton, firstManifest := read()
	secondSkeleton, secondManifest := read()
	if firstSkeleton != secondSkeleton {
		t.Errorf("the skeleton differs between runs:\n%q\n%q", firstSkeleton, secondSkeleton)
	}
	if firstManifest != secondManifest {
		t.Errorf("the manifest differs between runs:\n%q\n%q", firstManifest, secondManifest)
	}
	// And the members really are in sorted order rather than merely coincidentally
	// equal twice.
	if !strings.Contains(firstSkeleton, "\"m0\"") ||
		strings.Index(firstSkeleton, "\"m0\"") > strings.Index(firstSkeleton, "\"m1\"") {
		t.Errorf("members are not written in sorted order:\n%s", firstSkeleton)
	}
}

// A file under live/ that the manifest does not name is evidence nobody can
// trace, so the writer reports it rather than leaving it to be found later. The
// case is the one a real re-run produces: an endpoint that answered 200 last
// month and does not answer it this month.
func TestWriteCaptureReportsAFileItDidNotWrite(t *testing.T) {
	dir := t.TempDir()
	entries := []ManifestEntry{{
		Symbol: "data.GetX", Fixture: "fundamentals/x.json", Skeleton: "live/fundamentals/x.json",
		Status: 200, DecodedCleanly: true, Host: "h", ProbedAt: "2026-09-29T04:05:06Z",
	}}
	if err := WriteCapture(dir, map[string][]byte{"live/fundamentals/x.json": []byte("{}\n")}, entries); err != nil {
		t.Fatalf("WriteCapture: %v", err)
	}
	// Now simulate a previous, longer run: a skeleton this one did not write.
	stale := filepath.Join(dir, "live", "fundamentals", "gone.json")
	if err := os.WriteFile(stale, []byte("{}\n"), 0o600); err != nil {
		t.Fatalf("writing the stale file: %v", err)
	}
	err := WriteCapture(dir, map[string][]byte{"live/fundamentals/x.json": []byte("{}\n")}, entries)
	if err == nil {
		t.Fatal("WriteCapture accepted a file under live/ that its own manifest does not name")
	}
	if !strings.Contains(err.Error(), "gone.json") {
		t.Errorf("the error does not name the stale file: %v", err)
	}
	// The stale file is reported, never deleted: it is either evidence this run
	// could not reproduce or a file somebody placed by hand, and a writer that
	// removes it has made that decision without looking.
	if _, statErr := os.Stat(stale); statErr != nil {
		t.Errorf("the stale file was removed rather than reported: %v", statErr)
	}
}

// The tree's README is prose, not evidence, so the writer neither writes it nor
// holds it stale. A capture that treated it as an unnamed fixture would fail on
// every run after the first.
func TestWriteCaptureExemptsTheTreeReadmeFromTheStaleCheck(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "live", liveReadmeName)
	if err := os.MkdirAll(filepath.Dir(readme), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(readme, []byte("# live\n"), 0o600); err != nil {
		t.Fatalf("writing the readme: %v", err)
	}
	entries := []ManifestEntry{{
		Symbol: "data.GetX", Fixture: "fundamentals/x.json", Skeleton: "live/fundamentals/x.json",
		Status: 200, DecodedCleanly: true, Host: "h", ProbedAt: "2026-09-29T04:05:06Z",
	}}
	if err := WriteCapture(dir, map[string][]byte{"live/fundamentals/x.json": []byte("{}\n")}, entries); err != nil {
		t.Fatalf("WriteCapture rejected the tree's own README: %v", err)
	}
}

// A skeleton the manifest does not name, and a manifest row naming a skeleton
// that was not produced, are the same failure in two directions: a file with no
// provenance, and a claim about a file that is not there.
func TestWriteCaptureRefusesAnUnindexedSkeleton(t *testing.T) {
	cases := []struct {
		name     string
		skeleton map[string][]byte
		entries  []ManifestEntry
	}{
		{
			name:     "produced but not named",
			skeleton: map[string][]byte{"live/fundamentals/x.json": []byte("{}\n")},
			entries:  nil,
		},
		{
			name:     "named but not produced",
			skeleton: nil,
			entries: []ManifestEntry{{
				Fixture: "fundamentals/x.json", Skeleton: "live/fundamentals/x.json", Status: 200,
			}},
		},
		{
			name:     "named twice",
			skeleton: map[string][]byte{"live/fundamentals/x.json": []byte("{}\n")},
			entries: []ManifestEntry{
				{Fixture: "fundamentals/x.json", Skeleton: "live/fundamentals/x.json", Status: 200},
				{Fixture: "fundamentals/y.json", Skeleton: "live/fundamentals/x.json", Status: 200},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := WriteCapture(dir, tc.skeleton, tc.entries); err == nil {
				t.Fatal("WriteCapture accepted a tree and a manifest that disagree")
			}
			assertTreeEmpty(t, dir)
		})
	}
}

// --------------------------------------------------------------------------
// The capture call
// --------------------------------------------------------------------------

// tokenEndpoint is the endpoint whose SDK response type the decode tests run
// against. client.CreateToken decodes auth.Token, whose token member is a Go
// string, so a number there is a mismatch the Go decoder must reject - which is
// what makes it usable here. An unknown member would not do: encoding/json
// ignores one, so a body carrying only {"a":1} decodes cleanly into almost every
// DTO and would have made a decode test pass for the wrong reason.
var tokenEndpoint = Endpoint{
	Symbol:  "client.CreateToken",
	Fixture: "authentication/POST-auth-tokens-create.json",
	Method:  "POST",
	Path:    "/auth/tokens/create",
}

// A body the SDK type cannot hold at all: the skeleton is still valid, and the
// decode failure is a finding rather than a reason to drop the capture. This is
// the case the whole phase exists for - conformance can show a DTO matches the
// documentation, and only a live body can show it matches the server - so it
// must be recorded on the row and must not stop the capture.
func TestCaptureRecordsDecodeFailureWithoutFailing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// token is a string in auth.Token and a number here.
		_, _ = w.Write([]byte(`{"token":1}`))
	}))
	defer srv.Close()

	skeleton, entry, err := captureThroughServer(t, srv, tokenEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture returned an error for a decode failure: %v", err)
	}
	if skeleton == nil {
		t.Fatalf("the capture was dropped: %s", entry.NotCaptured)
	}
	if entry.Skeleton == "" {
		t.Error("no skeleton was recorded")
	}
	if entry.DecodedCleanly {
		t.Error("decodedCleanly is true for a body the SDK type cannot hold")
	}
	if entry.DecodeErr == "" {
		t.Error("decodeError is empty, so the row cannot say why")
	}
	// The rejection is recorded as a kind and a field path, never as a value.
	if !strings.Contains(entry.DecodeErr, "token") {
		t.Errorf("decodeError does not name the field that failed: %q", entry.DecodeErr)
	}
	if strings.Contains(entry.DecodeErr, `{"token"`) {
		t.Errorf("decodeError quotes the body: %q", entry.DecodeErr)
	}
}

// The converse: a body the SDK type accepts is recorded as such, because
// "decodedCleanly" is only evidence if false means something.
func TestCaptureRecordsACleanDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// The three members auth.Token declares, with the kinds it declares.
		_, _ = w.Write([]byte(`{"token":"a-token","expires_at":1756000000000,"status":"ok"}`))
	}))
	defer srv.Close()

	_, entry, err := captureThroughServer(t, srv, tokenEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if !entry.DecodedCleanly {
		t.Errorf("decodedCleanly is false for a body the SDK type accepts: %s", entry.DecodeErr)
	}
	if entry.DecodeErr != "" {
		t.Errorf("decodeError = %q on a clean decode", entry.DecodeErr)
	}
	// The token value is in the response and is nowhere in the row or the file.
	for _, leaked := range []string{"a-token", "1756000000000"} {
		if strings.Contains(entry.DecodeErr, leaked) || strings.Contains(entry.NotCaptured, leaked) {
			t.Errorf("the row carries the live value %q", leaked)
		}
	}
}

// orderEndpoint is a plain read endpoint in an area the mutation gate does not
// cover, used by the tests that are about the shape rather than about the SDK's
// tolerance of it. Its DTO is trade.OrderGroup, which ignores an unknown member,
// so it is not the endpoint a decode test should use; see tokenEndpoint.
var orderEndpoint = Endpoint{
	Symbol:  "trade.GetOrderDetail",
	Fixture: "trading/GET-trading-orders-get.json",
	Method:  "GET",
	Path:    "/trading/orders/1",
}

// A 200 whose body is not JSON is recorded and skipped. It is not written as an
// empty skeleton, because an empty file reads as "the endpoint answered with
// nothing" - a claim about the server that the server did not make.
func TestCaptureSkipsANonJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>502 Bad Gateway</body></html>"))
	}))
	defer srv.Close()

	skeleton, entry, err := captureThroughServer(t, srv, orderEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if skeleton != nil {
		t.Errorf("a skeleton was written for a non-JSON body: %s", skeleton)
	}
	if entry.NotCaptured == "" {
		t.Fatal("nothing was written and nothing says why")
	}
	if !strings.Contains(entry.NotCaptured, "text/html") {
		t.Errorf("the reason does not name the media type: %q", entry.NotCaptured)
	}
}

// A 200 whose body claims to be JSON and is not is a different statement, and it
// is recorded as a different one.
func TestCaptureSkipsABodyThatDoesNotReduce(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"a": `))
	}))
	defer srv.Close()

	skeleton, entry, err := captureThroughServer(t, srv, orderEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if skeleton != nil {
		t.Errorf("a skeleton was written for a body that does not reduce: %s", skeleton)
	}
	if !strings.Contains(entry.NotCaptured, "did not reduce") {
		t.Errorf("the reason does not say the body did not reduce: %q", entry.NotCaptured)
	}
}

// A body past the bound is refused rather than truncated. A prefix of a JSON
// document either fails to reduce - losing the endpoint - or closes inside the
// prefix and describes a shape the endpoint never sent, which is the one failure
// the harness cannot detect about itself.
func TestCaptureRefusesAnOversizedBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// A body that closes well inside the bound and is then padded past it, so
		// a truncating reader would have had a complete document to reduce.
		body := `{"a":[` + strings.Repeat(`1,`, maxCaptureBytes) + `1]}`
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	skeleton, entry, err := captureThroughServer(t, srv, orderEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if skeleton != nil {
		t.Error("a skeleton was written for a body past the bound")
	}
	if !strings.Contains(entry.NotCaptured, "exceeds") {
		t.Errorf("the reason does not mention the bound: %q", entry.NotCaptured)
	}
}

// A body exactly at the bound is accepted. The reader reads one byte past the
// limit precisely so that "at the limit" and "one past it" are distinguishable.
func TestCaptureAcceptsABodyExactlyAtTheBound(t *testing.T) {
	// {"a":"<padding>"} is 10 bytes of syntax around the string content.
	const overhead = len(`{"a":""}`)
	inner := `{"a":"` + strings.Repeat("x", maxCaptureBytes-overhead) + `"}`
	if len(inner) != maxCaptureBytes {
		t.Fatalf("the body is %d bytes, not %d", len(inner), maxCaptureBytes)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(inner))
	}))
	defer srv.Close()

	skeleton, entry, err := captureThroughServer(t, srv, orderEndpoint, 200)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if skeleton == nil {
		t.Fatalf("a body exactly at the bound was refused: %s", entry.NotCaptured)
	}
	if !bytes.Contains(skeleton, []byte(`"a": "1"`)) {
		t.Errorf("the reduced tree lost the member: %s", skeleton)
	}
}

// The decode error is reconstructed rather than quoted, because
// encoding/json returns a custom UnmarshalJSON error verbatim and
// money.Money.UnmarshalJSON can name the price it was handed. This is the one
// path by which a server's bytes could reach the manifest, and it is the reason
// the reconstruction exists.
func TestDecodeFailureQuotesNothingFromTheBody(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
		// never is a substring the rendered failure must not contain, chosen
		// per case so the assertion is about this input rather than about a
		// word list that would match the message's own vocabulary.
		never []string
	}{
		{
			name:  "a number where the SDK declares a string",
			raw:   `{"quantity":1756000000.5}`,
			want:  "cannot unmarshal a JSON number into Go field",
			never: []string{`{"quantity"`, "1756000000"},
		},
		{
			name:  "an object where the SDK declares a string",
			raw:   `{"quantity":{"nested":1756000000}}`,
			want:  "cannot unmarshal a JSON object into Go field",
			never: []string{"nested", "1756000000"},
		},
		{
			name:  "a body that is not JSON at all",
			raw:   `<html>Access denied for account 9110101000000000001</html>`,
			want:  "the body is not well-formed JSON",
			never: []string{"html", "9110101000000000001", "Access denied"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := decodeFailure(json.Unmarshal([]byte(tc.raw), newOrderProbe()))
			if !strings.Contains(got, tc.want) {
				t.Errorf("decodeFailure = %q, want it to contain %q", got, tc.want)
			}
			for _, leaked := range tc.never {
				if strings.Contains(got, leaked) {
					t.Errorf("decodeFailure = %q, which quotes %q from the body", got, leaked)
				}
			}
		})
	}

	// A DTO whose UnmarshalJSON fails on the value it was handed: the one shape
	// of error whose message genuinely can carry a reading, and the reason the
	// reconstruction exists.
	t.Run("a custom UnmarshalJSON error is classified, not quoted", func(t *testing.T) {
		secret := "1784500000.123456"
		got := decodeFailure(&echoingError{value: secret})
		if strings.Contains(got, secret) {
			t.Errorf("decodeFailure = %q, which quotes the value the DTO was handed", got)
		}
		if !strings.Contains(got, "not recorded") {
			t.Errorf("decodeFailure = %q, want it to say the message was withheld", got)
		}
	})
}

// orderProbe is the DTO decodeFailure is exercised against: a struct whose
// quantity is a string, so a number there is the type mismatch the
// reconstruction is written for.
type orderProbe struct {
	AccountNumber string `json:"account_number"`
	Quantity      string `json:"quantity"`
}

func newOrderProbe() *orderProbe { return &orderProbe{} }

// echoingError stands in for a DTO's own UnmarshalJSON error, which is the only
// error shape that can carry a value into a message.
type echoingError struct {
	value string
}

func (e *echoingError) Error() string {
	return "money: " + e.value + " is not a decimal"
}

// The two refusals happen before a request is sent, and each exists because the
// alternative is a committed file that is not evidence. The mutating one is the
// second barrier behind the census gate rather than a replacement for it.
func TestCaptureRefusesBeforeItSends(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl := newProbeTestClient(t, srv.URL)
	cases := []struct {
		name string
		ep   Endpoint
		prev Outcome
		want error
	}{
		{
			name: "the census recorded a 404",
			ep: Endpoint{Symbol: "trade.GetOrderDetail", Fixture: "trading/x.json",
				Method: "GET", Path: "/trading/orders/123"},
			prev: Outcome{Status: 404},
			want: errNotOK,
		},
		{
			name: "the census recorded a 403",
			ep: Endpoint{Symbol: "trade.GetOrderDetail", Fixture: "trading/x.json",
				Method: "GET", Path: "/trading/orders/123"},
			prev: Outcome{Status: 403},
			want: errNotOK,
		},
		{
			name: "the endpoint can mutate",
			ep: Endpoint{Symbol: "trade.PlaceOrder", Fixture: "trading/y.json",
				Method: "POST", Path: "/trading/orders/place"},
			prev: Outcome{Status: 200},
			want: errMutating,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			skeleton, entry, err := Capture(t.Context(), cl, tc.ep, tc.prev)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if skeleton != nil {
				t.Errorf("a skeleton was returned: %s", skeleton)
			}
			if entry.Skeleton != "" {
				t.Errorf("the row names a skeleton: %s", entry.Skeleton)
			}
		})
	}
	if called {
		t.Error("a refused capture sent a request")
	}
}

// A capture that cannot build the request phase 1 built is a probe failure, and
// it is reported as one rather than recorded as a shape.
func TestCaptureRefusesARequestItCannotBuild(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl := newProbeTestClient(t, srv.URL)
	// An account-scoped endpoint with no account id threaded in is the block the
	// census itself can produce, and a capture must not paper over it.
	SetAccountID("")
	t.Cleanup(func() { SetAccountID("") })

	ep := Endpoint{Symbol: "trade.GetPositions", Fixture: "trading/positions.json",
		Method: "GET", Path: "/trading/accounts/1/assets/positions",
		RequiredParams: []NamedParam{{Name: "account_id", Spec: ParamSpec{Type: "string"}}}}
	skeleton, _, err := Capture(t.Context(), cl, ep, Outcome{Status: 200})
	if err == nil {
		t.Fatal("Capture accepted a request it could not build")
	}
	if skeleton != nil {
		t.Errorf("a skeleton was returned: %s", skeleton)
	}
}

// A capture run is refused outright when the mutation gate is open, because a
// capture re-calls the endpoints that answered 200 and an open gate would put
// the order-placing ones among them. The refusal is up front, so the operator is
// told once rather than through 34 rows that silently went missing.
func TestRunCaptureRefusesAnOpenMutationGate(t *testing.T) {
	t.Setenv(MutateOptInEnv, MutateOptInValue)
	t.Setenv(NonSandboxOverrideEnv, NonSandboxOverrideValue)

	cl := newProbeTestClient(t, "https://not-a-sandbox.example")
	err := runCapture(t.Context(), cl, nil, nil, CensusSummary{}, t.TempDir(),
		"not-a-sandbox.example", "hk", "live")
	if err == nil {
		t.Fatal("runCapture ran with the mutation gate open")
	}
	if !strings.Contains(err.Error(), "refusing to capture") {
		t.Errorf("the refusal does not say why: %v", err)
	}
}

// The one clock reading per run is what keeps a re-run's manifest diffable, so
// it is pinned here: SetProbedAt overrides it, and every row carries the same
// value rather than its own.
func TestProbedAtIsOneReadingPerRun(t *testing.T) {
	at := time.Date(2026, 9, 29, 4, 5, 6, 0, time.UTC)
	SetProbedAt(at)
	t.Cleanup(func() { SetProbedAt(time.Time{}) })

	if got := probedAt(); got != "2026-09-29T04:05:06Z" {
		t.Errorf("probedAt() = %q, want the recorded reading", got)
	}
	// With no reading recorded, a caller that does not care still gets a usable
	// timestamp rather than a zero time.
	SetProbedAt(time.Time{})
	if got := probedAt(); got == "" {
		t.Error("probedAt() is empty with no reading recorded")
	}
}

// liveKey is the join between the two trees, and it has to be the fixture name
// unchanged under the live prefix, or a comparison between them is a lookup.
func TestLiveKeyMirrorsTheFixtureName(t *testing.T) {
	const fixture = "trading/GET-trading-orders-get.json"
	if got, want := liveKey(fixture), "live/"+fixture; got != want {
		t.Errorf("liveKey(%q) = %q, want %q", fixture, got, want)
	}
}

// isJSONMediaType decides whether a 200's body is worth reducing, and the
// tolerance set is the interesting part: a missing header must not lose an
// endpoint, and a vendor media type must not either.
func TestIsJSONMediaType(t *testing.T) {
	cases := []struct {
		header string
		want   bool
	}{
		{"application/json", true},
		{"application/json; charset=utf-8", true},
		{"APPLICATION/JSON", true},
		{"text/json", true},
		{"application/vnd.webull+json", true},
		{"text/plain", true},
		{"", true},
		{"text/html", false},
		{"text/html; charset=utf-8", false},
		{"application/octet-stream", false},
	}
	for _, tc := range cases {
		t.Run(tc.header, func(t *testing.T) {
			if got := isJSONMediaType(tc.header); got != tc.want {
				t.Errorf("isJSONMediaType(%q) = %v, want %v", tc.header, got, tc.want)
			}
		})
	}
}

// --------------------------------------------------------------------------
// Shared helpers
// --------------------------------------------------------------------------

// captureThroughServer runs one Capture against an httptest server, with the
// census row that admits it to phase 2.
func captureThroughServer(t *testing.T, srv *httptest.Server, ep Endpoint, status int) (
	[]byte, ManifestEntry, error) {
	t.Helper()
	cl := newProbeTestClient(t, srv.URL)
	prev := Outcome{Symbol: ep.Symbol, Fixture: ep.Fixture, Status: status}
	return Capture(t.Context(), cl, ep, prev)
}

// assertNoLiveValueInFile walks the bytes a writer produced and fails on any leaf
// that could carry something the server sent. It is the same walk
// skeleton_test.go applies to a tree in memory, run here on the committed form:
// a reduction can be correct in memory and a serialisation can still be wrong,
// and only the second one reaches git.
func assertNoLiveValueInFile(t *testing.T, raw []byte) {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var decoded any
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("the committed skeleton does not decode: %v", err)
	}
	assertNoLiveValue(t, decoded, "$")
}

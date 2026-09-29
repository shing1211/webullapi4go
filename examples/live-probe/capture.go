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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/conformance"
)

// --------------------------------------------------------------------------
// What phase 2 commits
// --------------------------------------------------------------------------
//
// Phase 1 asked 193 documented endpoints whether the sandbox answers them and
// recorded 55 that answered HTTP 200. Those 55 are the only bodies worth
// reducing, and the restriction is not politeness. A 404, a 403, a 417 and a 504
// are live answers about reachability, but their bodies are Webull's error
// shape; reducing one and comparing it with a documented 200 response would
// manufacture a divergence out of an endpoint the sandbox simply does not serve.
// So the scope is 200, and every other status stays where phase 1 put it: in
// the census.
//
// What is committed is a shape, not a reading. Skeletonify reduces a live body
// to a value-free tree - every member name kept, every value replaced by a typed
// placeholder - and this file writes that tree out. Nothing else from the
// response is recorded: not a header, not a length of the original body, not the
// server's message, and not the account id, symbol, price or timestamp the
// request carried. A skeleton in the repository is a file of "1", -1, true and
// null, so committing 55 of them puts 55 shapes in git and zero readings.
//
// Two trees now sit under conformance/testdata. The documentation fixtures are
// derived from Webull's published OpenAPI JSON by tools/conformance, and
// conformance.CompareBody reads them. The live skeletons are derived from the
// sandbox, and nothing in the conformance package reads them yet. They are kept
// in separate subtrees with separate manifests so a reader scanning the
// directory can tell which is which without opening a file, and so a change to
// one can never be mistaken for drift in the other.

// liveTreeName is the single directory this probe may write into under
// conformance/testdata.
const liveTreeName = "live"

// liveTree is liveTreeName as a path prefix, which is the form every skeleton
// key is written in and the form the guard in resolveLivePath tests.
const liveTree = liveTreeName + "/"

// liveManifestName is the manifest that indexes the live tree, beside it and
// not inside it: a directory that holds evidence should not also hold the index
// of that evidence, and conformance/testdata already holds manifest.json for the
// documentation tree.
const liveManifestName = "live-manifest.json"

// liveReadmeName is the note inside the live tree that tells a reader what the
// tree is before they open a file in it. WriteLiveCapture leaves it alone: it is
// prose, not evidence, so it is neither written nor held stale.
const liveReadmeName = "README.md"

// maxCaptureBytes bounds the response this probe will read into memory.
//
// The bound is a correctness rule and not a politeness one. A body past it is
// recorded as too large and skipped, where truncating it would hand Skeletonify
// a prefix of a JSON document: the reduction would either fail, which loses the
// endpoint from the evidence set, or - far worse - succeed on a document that
// happened to close inside the prefix and describe a shape the endpoint never
// sent. A missing skeleton is recoverable; a wrong one is evidence the harness
// does not know is wrong. 1 MiB is far above every response in the current
// corpus and far below anything that would strain the process.
const maxCaptureBytes = 1 << 20

// errBodyTooLarge is returned by readBounded for a response past
// [maxCaptureBytes]. It is a sentinel because the right response is to record
// the endpoint as not captured and carry on, not to fail a run over one large
// body.
var errBodyTooLarge = errors.New("live-probe: the response body exceeds " +
	"1048576 bytes, so it was not read rather than truncated")

// ManifestEntry is one row of conformance/testdata/live-manifest.json: the
// provenance of one capture, or the reason there is no capture.
//
// The row answers four questions and carries nothing else. Which SDK symbol
// produced it is Symbol. From which HTTP status is Status, and it is the status
// of the capture call itself rather than the census call, because the two are
// separate requests to a server that can answer differently: an endpoint that
// was a 200 in phase 1 and a 500 here is recorded as the latter, with the
// disagreement left visible rather than reconciled. Against which host is Host,
// read off the client that sent the request. When is ProbedAt.
//
// Fixture is the documentation fixture the row belongs to - "trading/GET-trading-orders-get.json" -
// so the row joins to conformance/testdata/manifest.json by name. Skeleton is
// where the shape was written, and it is empty exactly when NotCaptured is not.
//
// No field carries a value the server sent. NotCaptured is this probe's own
// sentence about what it saw, DecodeErr is a reconstruction from the decoder's
// own type rather than its message (see decodeFailure), and the four fields that
// could plausibly have held a reading - a body, a length, a header value, an
// account id - do not exist on this struct.
type ManifestEntry struct {
	// Symbol is the SDK method the documentation manifest maps this endpoint to.
	Symbol string `json:"symbol"`
	// Fixture is the documentation fixture file name, "<area>/<METHOD>-<slug>.json",
	// which is the join key into conformance/testdata/manifest.json.
	Fixture string `json:"fixture"`
	// Skeleton is the live skeleton this entry produced, relative to
	// conformance/testdata and always under live/. Empty when nothing was written.
	Skeleton string `json:"skeleton,omitempty"`
	// Status is the HTTP status the capture call was answered with. A row is
	// only written for an endpoint phase 1 recorded as a 200, so a different
	// status here is the server changing its mind between two calls, and it is
	// recorded rather than reconciled with the census.
	Status int `json:"status"`
	// DecodedCleanly reports that the live body unmarshalled into the response
	// type conformance.SDKTypes records for Symbol. It is the one finding the
	// manifest carries, because a body the SDK's own type rejects is a defect
	// the documentation comparison cannot see: the documentation says the name
	// is a string, the fixture holds a string, and only the live body says
	// otherwise.
	DecodedCleanly bool `json:"decodedCleanly"`
	// DecodeErr explains a false DecodedCleanly. It is a reconstruction from the
	// decoder's error type, never the error text, and the reason is in
	// decodeFailure: a custom UnmarshalJSON may name the value it was handed.
	DecodeErr string `json:"decodeError,omitempty"`
	// NotCaptured says why no skeleton was written, and is empty on every row
	// that has one. A 200 whose body is not JSON, whose body is larger than
	// [maxCaptureBytes], or whose body does not reduce to a type skeleton is
	// recorded here rather than as an empty file, because an empty file reads
	// as "the endpoint answered with nothing" and is a claim about the server.
	NotCaptured string `json:"notCaptured,omitempty"`
	// Host is the host that answered: the client's own base URL host, with the
	// scheme, path, query and any userinfo dropped. A base URL is
	// operator-supplied and a URL may carry a credential; the host is the one
	// fact that says which deployment this evidence came from.
	Host string `json:"host"`
	// ProbedAt is when the capture call was sent, in UTC RFC 3339 to the second.
	// It is the one field in this tree that changes on every re-run, which is
	// why every row in one run carries the same reading rather than a per-row
	// one: a re-run's manifest then differs from the previous one in exactly one
	// value per row, so a diff shows what the server sent differently.
	ProbedAt string `json:"probedAt"`
}

// LiveGenerator records how the live tree was produced, in the shape
// conformance/testdata/manifest.json uses for its own generator block.
type LiveGenerator struct {
	// Tool is the program that produced the tree.
	Tool string `json:"tool"`
	// Command is the invocation that reproduces it.
	Command string `json:"command"`
	// Source is where the input came from, which is a running server rather than
	// a committed file, and the credentials for it come from the environment.
	Source string `json:"source"`
	// DerivedFrom is the transformation applied to each response body.
	DerivedFrom string `json:"derivedFrom"`
	// ContainsValues is false, and it is stated as a field rather than left to
	// be inferred because it is the property a reader of a committed fixture has
	// no other way to check.
	ContainsValues bool `json:"containsValues"`
	// LeafForms spells out what a placeholder is, so a reader who finds a `1` in
	// a live file knows it is the string placeholder and not a value.
	LeafForms string `json:"leafForms"`
	// Determinism states what a re-run can and cannot reproduce byte for byte.
	Determinism string `json:"determinism"`
}

// LiveTotals summarises one capture run.
type LiveTotals struct {
	// Considered is how many endpoints phase 1 recorded as a 200, which is the
	// whole scope of this manifest.
	Considered int `json:"considered"`
	// Captured is how many of those produced a written skeleton.
	Captured int `json:"captured"`
	// NotCaptured is how many produced none, and NotCapturedReasons counts the
	// distinct sentences so the reason is a list a reader can act on.
	NotCaptured        int            `json:"notCaptured"`
	NotCapturedReasons map[string]int `json:"notCapturedReasons,omitempty"`
	// DecodedCleanly and DecodeRejected partition Considered by whether the
	// live body unmarshalled into the SDK response type.
	DecodedCleanly int `json:"decodedCleanly"`
	DecodeRejected int `json:"decodeRejected"`
	// SkeletonBytes is the combined length of the written skeletons, so the
	// tree's size is on the record rather than only in a diff.
	SkeletonBytes int `json:"skeletonBytes"`
}

// LiveManifest is the whole of conformance/testdata/live-manifest.json: the
// index of the live tree, and the only place a reader has to look to learn what
// the tree is.
type LiveManifest struct {
	// Kind names the artifact, so a tool - or a reader - that opens the file
	// without its path in view can still tell what it holds. It differs from
	// conformance/testdata/manifest.json's own shape deliberately: the two
	// files have different authors, different inputs and different consumers,
	// and the one thing they must never be is interchangeable.
	Kind string `json:"kind"`
	// Warning is the single sentence a reader needs before the entries.
	Warning string `json:"warning"`
	// Generator records how the tree was produced.
	Generator LiveGenerator `json:"generator"`
	// Scope states what was and was not captured, so a reader never has to infer
	// the boundary from the row count.
	Scope string `json:"scope"`
	// Host is the host that answered this run.
	Host string `json:"host,omitempty"`
	// Region and Environment are the client's resolved deployment selection.
	Region      string `json:"region,omitempty"`
	Environment string `json:"environment,omitempty"`
	// Totals summarises the run.
	Totals LiveTotals `json:"totals"`
	// Census is the reachability census this run is a subset of, so the 55 rows
	// below can be read against the 193 they were drawn from and against the
	// status mix of the ones that were left out.
	Census *CensusSummary `json:"census,omitempty"`
	// Entries is one row per endpoint phase 1 recorded as a 200, sorted by
	// Fixture so the file is a function of the run rather than of the walk
	// order.
	Entries []ManifestEntry `json:"entries"`
}

// liveKind identifies the artifact to anything reading the file.
const liveKind = "live-response-skeletons"

// liveWarning is what a reader has to know before the entries.
const liveWarning = "These files hold the SHAPE of a live Webull response - every member " +
	"name the sandbox sent, and no value it sent. Every leaf is a placeholder: a string " +
	"is \"1\", a number is -1, a boolean is true, a null is null. Nothing in this tree is a " +
	"price, a size, a timestamp, an account id or a symbol, and there is no reading in it to " +
	"recover."

// liveScope states the boundary of the capture.
const liveScope = "Only endpoints the reachability census recorded as HTTP 200 are captured. " +
	"A 404, 401, 403, 417, 500 or 504 is a live answer about whether an endpoint is served, " +
	"but its body is Webull's error shape, and reducing one would compare an error against a " +
	"documented success response and manufacture a divergence. The census block below records " +
	"every other status this run saw and why it was not captured."

// liveGeneratorBlock records the provenance of the tree.
func liveGeneratorBlock() LiveGenerator {
	return LiveGenerator{
		Tool:    "examples/live-probe (phase 2, -capture)",
		Command: "go run . -capture -out ../../conformance/testdata",
		Source: "the sandbox named by host below, reached through client.WithEnv; the " +
			"app key and app secret come from the process environment and reach no file here",
		DerivedFrom: "the live response body of each endpoint, reduced by Skeletonify, which " +
			"keeps every member name and replaces every value",
		ContainsValues: false,
		LeafForms: `"1" is a string, -1 is a number, true is a boolean, null is a null; ` +
			"an array reduces to an array and an object to an object, so the shape is preserved " +
			"and the reading is not",
		Determinism: "A skeleton is a function of the response body alone: object members are " +
			"written in sorted order, so the same body always produces the same bytes and a diff " +
			"means the server sent something different. This manifest is the one file that " +
			"changes on every re-run, because probedAt is a clock reading.",
	}
}

// NewLiveManifest assembles the manifest for one capture run from its rows and
// the skeletons they produced.
//
// The counts are derived rather than passed in, because a total a caller
// supplies and a total the rows do not agree on is a manifest that misreports
// its own tree. entries is copied and sorted, so a caller's slice is neither
// reordered nor aliased.
func NewLiveManifest(entries []ManifestEntry, skeletons map[string][]byte) LiveManifest {
	doc := LiveManifest{
		Kind:      liveKind,
		Warning:   liveWarning,
		Generator: liveGeneratorBlock(),
		Scope:     liveScope,
		Entries:   append([]ManifestEntry(nil), entries...),
	}
	sort.Slice(doc.Entries, func(i, j int) bool { return doc.Entries[i].Fixture < doc.Entries[j].Fixture })

	for _, e := range doc.Entries {
		doc.Totals.Considered++
		if e.Skeleton != "" {
			doc.Totals.Captured++
			doc.Totals.SkeletonBytes += len(skeletons[e.Skeleton])
		}
		if e.NotCaptured != "" {
			doc.Totals.NotCaptured++
			if doc.Totals.NotCapturedReasons == nil {
				doc.Totals.NotCapturedReasons = map[string]int{}
			}
			doc.Totals.NotCapturedReasons[e.NotCaptured]++
		}
		if e.DecodedCleanly {
			doc.Totals.DecodedCleanly++
		} else {
			doc.Totals.DecodeRejected++
		}
		if e.Host != "" && doc.Host == "" {
			doc.Host = e.Host
		}
	}
	return doc
}

// WriteCapture writes the skeletons and the live manifest for one run under
// dir, which is a conformance/testdata directory.
//
// It is the narrow form of [WriteLiveCapture]: it builds the manifest from the
// rows alone, so a caller with no census to attach can write a tree a reader can
// check. The census is a convenience, not a requirement, and requiring it here
// would make the function untestable without a cache.
func WriteCapture(dir string, skeletons map[string][]byte, entries []ManifestEntry) error {
	return WriteLiveCapture(dir, skeletons, NewLiveManifest(entries, skeletons))
}

// WriteLiveCapture writes skeletons and doc under dir, which is a
// conformance/testdata directory.
//
// The guard is the reason this function exists as a function and not as a
// `os.WriteFile` at each call site. Everything is validated before anything is
// written, so a rejected run leaves the tree exactly as it found it rather than
// half-updated; then the skeletons are written; then the manifest; then the tree
// is walked to prove it holds nothing but what this run wrote. A file under live/
// that the manifest does not name is an error and not a deletion: it is either
// evidence from an earlier run that this one failed to reproduce, which a
// reader must look at, or a file somebody put there by hand.
//
// Every key must name a file under live/. That is the whole of the refusal, and
// it is a refusal with teeth: conformance/testdata holds 193 committed
// documentation fixtures that conformance.CompareBody reads, and
// gen_fixtures.py --check regenerates from Webull's published schemas. Writing
// over one of them does not add evidence, it silently changes the gate that
// compares the SDK against the documentation - which is the last thing a probe
// whose purpose is to be evidence should be able to do.
func WriteLiveCapture(dir string, skeletons map[string][]byte, doc LiveManifest) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New("live-probe: WriteLiveCapture needs the conformance/testdata " +
			"directory to write into, and was given none")
	}
	root, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("live-probe: resolving %s: %w", dir, err)
	}

	// Every key is resolved before a single file is written, and the manifest is
	// checked against the set, so a run that would write an unindexed file or
	// index a file it did not write fails having written nothing.
	resolved := make(map[string]string, len(skeletons))
	for key := range skeletons {
		target, err := resolveLivePath(root, key)
		if err != nil {
			return err
		}
		resolved[key] = target
	}
	if err := checkIndexed(skeletons, doc.Entries); err != nil {
		return err
	}

	keys := make([]string, 0, len(skeletons))
	for key := range skeletons {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err := os.MkdirAll(filepath.Dir(resolved[key]), 0o750); err != nil {
			return fmt.Errorf("live-probe: creating the directory for %s: %w", key, err)
		}
		//nolint:gosec // G306: a skeleton is a few hundred bytes of placeholders, and it is committed content a reviewer reads.
		if err := os.WriteFile(resolved[key], skeletons[key], 0o644); err != nil {
			return fmt.Errorf("live-probe: writing %s: %w", key, err)
		}
	}
	manifestPath := filepath.Join(root, liveManifestName)
	encoded, err := json.MarshalIndent(doc, "", " ")
	if err != nil {
		return fmt.Errorf("live-probe: encoding %s: %w", liveManifestName, err)
	}
	encoded = append(encoded, '\n')
	//nolint:gosec // G306: the manifest is committed evidence, not a secret, and 0644 is what the tracked tree already uses.
	if err := os.WriteFile(manifestPath, encoded, 0o644); err != nil {
		return fmt.Errorf("live-probe: writing %s: %w", liveManifestName, err)
	}
	return checkNoStale(root, keys)
}

// checkIndexed holds the manifest and the file set to each other in both
// directions, for the same reason conformance holds its own two directions: a
// skeleton the manifest does not name has no provenance, and a manifest row
// naming a file that was not written is a claim about a file that is not there.
func checkIndexed(skeletons map[string][]byte, entries []ManifestEntry) error {
	named := make(map[string]int, len(entries))
	for _, e := range entries {
		if e.Skeleton == "" {
			continue
		}
		named[e.Skeleton]++
		if _, found := skeletons[e.Skeleton]; !found {
			return fmt.Errorf("live-probe: %s names skeleton %s, which this run did not produce",
				e.Fixture, e.Skeleton)
		}
	}
	for key := range skeletons {
		if named[key] == 0 {
			return fmt.Errorf("live-probe: skeleton %s was produced but no manifest entry names it, "+
				"so it would be committed with no provenance", key)
		}
	}
	for key, count := range named {
		if count > 1 {
			return fmt.Errorf("live-probe: skeleton %s is named by %d manifest entries", key, count)
		}
	}
	return nil
}

// checkNoStale reports every .json file under live/ that this run did not write.
//
// It is a walk rather than a list because the failure it catches is a file left
// by a previous, longer run: the set of endpoints that answer 200 shrinks, and a
// run that only adds and overwrites would leave the dropped ones behind, where
// the next comparison would find a skeleton for an endpoint this run could not
// reach. The failure mode it prevents is a manifest that is quietly out of date.
//
// Non-JSON files are exempt so the tree's README is not treated as evidence.
// The exemption is narrow on purpose: it is a filename suffix, not a content
// sniff.
func checkNoStale(root string, written []string) error {
	liveRoot := filepath.Join(root, filepath.FromSlash(liveTree))
	known := make(map[string]bool, len(written))
	for _, key := range written {
		known[key] = true
	}
	var stale []string
	err := filepath.WalkDir(liveRoot, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return relErr
		}
		key := filepath.ToSlash(rel)
		if !known[key] {
			stale = append(stale, key)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("live-probe: walking %s: %w", liveRoot, err)
	}
	if len(stale) == 0 {
		return nil
	}
	sort.Strings(stale)
	return fmt.Errorf("live-probe: %d file(s) under %s are not named by this run's manifest: %s. "+
		"A capture run writes every skeleton it captured, so a file it did not write is either "+
		"evidence from an earlier run that this one could not reproduce, or a file somebody "+
		"placed by hand. Read them before deleting them",
		len(stale), liveTree, strings.Join(stale, ", "))
}

// resolveLivePath turns a skeleton key into an absolute path under root/live, or
// refuses it.
//
// The rules are four and each rejects a different way out of the live tree: a
// backslash, because a key is forward-slashed the way the documentation manifest
// writes one and a Windows separator would make a key mean different files on
// different platforms; a leading slash, because an absolute key would name a
// file outside the directory the operator passed; a key that is not already in
// its simplest form, because a key that has to be cleaned is a key that can be
// cleaned into something else; and a first segment that is not live, which is
// the rule the last check below restates and the one that matters.
func resolveLivePath(root, key string) (string, error) {
	switch {
	case key == "":
		return "", errors.New("live-probe: a skeleton key is empty")
	case strings.ContainsRune(key, '\\'):
		return "", fmt.Errorf("live-probe: skeleton key %q holds a backslash; keys are "+
			"forward-slashed and relative, so one key names one file on every platform", key)
	case strings.HasPrefix(key, "/"):
		return "", fmt.Errorf("live-probe: skeleton key %q is absolute; a key names a file "+
			"under %s and never a path of its own", key, liveTree)
	case path.Clean(key) != key:
		return "", fmt.Errorf("live-probe: skeleton key %q is not in its simplest form; %q is "+
			"the same file, and a key that has to be cleaned is a key that can be cleaned "+
			"somewhere else", key, path.Clean(key))
	}

	first := key
	if index := strings.Index(key, "/"); index >= 0 {
		first = key[:index]
	}
	if first != liveTreeName {
		return "", fmt.Errorf("live-probe: skeleton key %q is not under %s. Only %s may be "+
			"written: a key resolving anywhere else in this directory would overwrite a "+
			"committed documentation fixture, and the conformance gate reads those files, so "+
			"overwriting one changes a gate rather than adding evidence to one",
			key, liveTree, liveTree)
	}
	if len(key) <= len(liveTree) {
		return "", fmt.Errorf("live-probe: skeleton key %q names the live directory rather than "+
			"a file inside it", key)
	}

	// The prefix test above is the readable statement of the rule; this is the
	// same rule in the form that cannot be wrong, checked on the resolved path
	// rather than on the string, so a future change to the first branch cannot
	// quietly open a way out of the tree.
	liveRoot := filepath.Join(root, filepath.FromSlash(liveTree))
	resolved := filepath.Join(root, filepath.FromSlash(key))
	if !strings.HasPrefix(resolved, liveRoot+string(filepath.Separator)) {
		return "", fmt.Errorf("live-probe: skeleton key %q resolves to %s, which is outside %s",
			key, resolved, liveRoot)
	}
	return resolved, nil
}

// encodeSkeleton renders a reduced tree as the bytes that get committed.
//
// The form is json.MarshalIndent with one space of indent and a trailing
// newline, which is what conformance/testdata/manifest.json uses, so the two
// files in the same directory read the same way. encoding/json writes the
// members of a map in sorted order, so the bytes are a function of the tree
// alone: the same response produces the same file, and a diff after a re-run
// means the server sent something different rather than that Go randomised a
// map.
//
// Only the reduced tree is marshalled. Nothing else about the response reaches
// this function, so there is no path by which a live value can reach the file it
// writes.
func encodeSkeleton(reduced any) ([]byte, error) {
	encoded, err := json.MarshalIndent(reduced, "", " ")
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}

// liveKey is where an endpoint's skeleton is written: the documentation fixture
// name under the live tree, so the two trees line up row for row and a
// comparison between them is a comparison of two names.
func liveKey(fixture string) string {
	return liveTree + fixture
}

// readBounded reads r whole, or refuses it.
//
// It reads one byte past [maxCaptureBytes] so that a body exactly at the limit
// is accepted and one byte beyond it is detectable. The alternative - reading to
// the limit and reducing the prefix - is the failure this refuses: a prefix of
// a JSON document is either a syntax error, which drops the endpoint from the
// evidence set, or a document that happened to close inside the prefix, which
// describes a shape the endpoint never sent.
func readBounded(r io.Reader) ([]byte, error) {
	raw, err := io.ReadAll(io.LimitReader(r, maxCaptureBytes+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > maxCaptureBytes {
		return nil, errBodyTooLarge
	}
	return raw, nil
}

// isJSONMediaType reports whether a Content-Type header names a JSON media type.
//
// It is a check on the declaration, not on the bytes, and it is what separates
// "the sandbox answered 200 with an HTML error page" from "the sandbox answered
// 200 with a body this probe could not parse". Both are recorded, but they are
// different statements about the endpoint, and a 200 carrying a text/html error
// page is a fact about the endpoint's configuration rather than about JSON.
//
// A missing header is accepted: HTTP permits it, several gateways drop it, and
// Skeletonify is the real arbiter of whether the body is JSON. Refusing here on a
// missing header would lose an endpoint for a reason the body does not bear out.
func isJSONMediaType(header string) bool {
	mediaType := header
	if index := strings.IndexByte(mediaType, ';'); index >= 0 {
		mediaType = mediaType[:index]
	}
	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "", "application/json", "text/json", "application/x-json", "text/plain":
		return true
	}
	// The vendor and structured-suffix forms: anything/json,
	// application/vnd.webull+json.
	return strings.HasSuffix(strings.ToLower(strings.TrimSpace(mediaType)), "json")
}

// decodeFailure renders why a live body did not unmarshal into an SDK response
// type, without quoting anything the body said.
//
// This is the one place in the probe where a server's bytes could reach a
// committed file through an error message, and it is closed. json.Unmarshal
// returns three error types of its own, and two of them are safe to render from
// their fields: an UnmarshalTypeError names the JSON kind, the field path and
// the Go type it wanted, all three of which are properties of the SDK type and
// the wire names rather than of any reading, and a SyntaxError names a byte
// offset. The third is everything else - and everything else is the problem.
// encoding/json returns an UnmarshalJSON error verbatim, so a DTO that fails to
// parse a price can return an error containing the price, and
// money.Money.UnmarshalJSON is exactly such a DTO. Rendering that message would
// put a live reading in a committed manifest, so it is classified by type and
// not quoted.
//
// The UnmarshalTypeError branch also carries more signal than the message did:
// Field is the wire path the decoder stopped at, which is a member name, and
// Value is the kind it found there, which is the finding itself - a number where
// the SDK declares a string is the change the documentation cannot show.
func decodeFailure(err error) string {
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		found := typeErr.Value
		if found == "" {
			found = "value"
		}
		if typeErr.Field != "" {
			return fmt.Sprintf("json: cannot unmarshal a JSON %s into Go field %s of type %s",
				found, typeErr.Field, typeErr.Type)
		}
		return fmt.Sprintf("json: cannot unmarshal a JSON %s into a value of Go type %s",
			found, typeErr.Type)
	}
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return fmt.Sprintf("the body is not well-formed JSON, at byte offset %d", syntaxErr.Offset)
	}
	var invalidErr *json.InvalidUnmarshalError
	if errors.As(err, &invalidErr) {
		return fmt.Sprintf("the body cannot be unmarshalled into %s, which is not a pointer",
			invalidErr.Type)
	}
	// Not quoted, and the reason is in this function's documentation: an error a
	// custom UnmarshalJSON returned may name the value it was handed.
	return fmt.Sprintf("the SDK response type rejected the body (%T); the decoder's own "+
		"message is not recorded because a custom UnmarshalJSON may name the value it was given",
		err)
}

// recordedAt holds the clock reading every row of one run carries.
//
// It is a package variable behind a mutex rather than a call to time.Now in
// Capture for the reason the field comment gives: 55 rows each reading their own
// clock would put 55 distinct timestamps in a committed file, and a re-run
// would then differ from the previous one in 55 places even when the server sent
// something identical. One reading per run makes a re-run's manifest differ in
// one value per row, so a diff is still readable - and the shape it is
// authoritative about, the skeleton, is byte-identical either way.
var recordedAt struct {
	mu sync.RWMutex
	at time.Time
}

// SetProbedAt records the clock reading every row of a capture run carries. A
// capture run calls it once, before the first endpoint. A nil or unset reading
// makes Capture fall back to the wall clock, so a caller that does not care -
// every test - gets a plausible timestamp without having to arrange one.
func SetProbedAt(at time.Time) {
	recordedAt.mu.Lock()
	defer recordedAt.mu.Unlock()
	recordedAt.at = at
}

// probedAt returns the recorded clock reading in the form the manifest writes,
// or the wall clock when none was recorded.
func probedAt() string {
	recordedAt.mu.RLock()
	at := recordedAt.at
	recordedAt.mu.RUnlock()
	if at.IsZero() {
		at = time.Now()
	}
	return at.UTC().Format(time.RFC3339)
}

// --------------------------------------------------------------------------
// The capture call
// --------------------------------------------------------------------------

// errNotOK and errMutating are the two refusals Capture makes before it sends
// anything. They are sentinels so a caller can recognise them without matching
// on text, and so a run can tell a refusal - which is the system working - from
// a transport failure, which is not.
var (
	errNotOK    = errors.New("live-probe: only an endpoint the census recorded as HTTP 200 is captured")
	errMutating = errors.New("live-probe: an endpoint in the mutating class is not captured")
)

// Capture calls one endpoint again and reduces its response to a value-free
// skeleton.
//
// It returns the bytes to write, the manifest row describing them, and an error.
// The three are independent on purpose. A nil skeleton with a nil error is a
// row that records why nothing was written - a 200 whose body is not JSON, whose
// body is larger than [maxCaptureBytes], or which does not reduce - because those
// are findings about the endpoint and dropping them would lose exactly the
// endpoints a reader most needs to know about. An error is reserved for a
// failure of the probe: a client it cannot use, a request it cannot build, or a
// response it could not read.
//
// Three refusals happen before a byte is sent, and each exists because the
// alternative is a committed file that is not evidence:
//
//   - prev.Status must be 200. A body captured from a 403 is Webull's error shape
//     and would compare against a documented success.
//   - the endpoint must not be in the mutating class. The census gate is closed
//     in a capture run, and this is the second barrier rather than the first: a
//     caller that walked the gate wrongly gets a refusal here instead of an
//     order.
//   - the request must build. A capture that filled a parameter phase 1 left
//     unresolved would record the server's rejection of the fabrication as the
//     endpoint's shape.
//
// The SDK decode is attempted on every 200 body, and a rejection is recorded
// rather than raised. That is the finding this whole phase exists to collect:
// conformance can prove a DTO matches the documentation, and only a live body
// can prove it matches the server.
//
// Nothing in the returned row, and nothing in the returned bytes, is a value the
// server sent. The body is reduced to placeholders, the decode error is
// reconstructed rather than quoted, and the row carries a symbol, a fixture
// name, a status, a host and a clock reading.
func Capture(ctx context.Context, cl *client.Client, ep Endpoint, prev Outcome) ([]byte, ManifestEntry, error) {
	if cl == nil {
		return nil, ManifestEntry{}, errors.New("live-probe: Capture needs a client")
	}
	entry := ManifestEntry{
		Symbol:   ep.Symbol,
		Fixture:  ep.Fixture,
		Host:     hostOf(cl),
		ProbedAt: probedAt(),
	}
	if prev.Status != http.StatusOK {
		return nil, entry, fmt.Errorf("live-probe: %s: the census recorded HTTP %d: %w",
			ep.Fixture, prev.Status, errNotOK)
	}
	if IsMutating(ep) {
		return nil, entry, fmt.Errorf("live-probe: %s: %w", ep.Fixture, errMutating)
	}
	request := prepare(ep)
	if request.Blocked != "" {
		return nil, entry, fmt.Errorf("live-probe: %s: the census built this request and this "+
			"call cannot: %s", ep.Fixture, request.Blocked)
	}

	target := request.Path
	if encoded := request.Query.Encode(); encoded != "" {
		target += "?" + encoded
	}
	var body any
	if len(request.Body) > 0 {
		body = json.RawMessage(request.Body)
	}
	resp, err := cl.DoStream(ctx, ep.Method, target, body)
	if resp == nil {
		return nil, entry, fmt.Errorf("live-probe: %s: no response: %w", ep.Fixture, err)
	}
	defer func() { _ = resp.Body.Close() }()
	entry.Status = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		entry.NotCaptured = fmt.Sprintf("the capture call answered HTTP %d where the census "+
			"recorded 200, so the two calls disagree and neither body is evidence of a "+
			"documented success response", resp.StatusCode)
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxCaptureBytes))
		return nil, entry, nil
	}

	raw, readErr := readBounded(resp.Body)
	if readErr != nil {
		if errors.Is(readErr, errBodyTooLarge) {
			entry.NotCaptured = errBodyTooLarge.Error() + ", and it was not read past that " +
				"bound because a truncated prefix of a JSON document describes a shape the " +
				"endpoint never sent"
			return nil, entry, nil
		}
		return nil, entry, fmt.Errorf("live-probe: %s: reading the response: %w", ep.Fixture, readErr)
	}

	recordDecode(&entry, ep.Symbol, raw)

	contentType := resp.Header.Get("Content-Type")
	if !isJSONMediaType(contentType) {
		declared := strings.TrimSpace(contentType)
		if declared == "" {
			declared = "none"
		}
		entry.NotCaptured = "the 200 response declares Content-Type " + declared + ", which is " +
			"not a JSON media type, so there is no response shape to reduce and the bytes were " +
			"not committed"
		return nil, entry, nil
	}
	reduced, reduceErr := Skeletonify(raw)
	if reduceErr != nil {
		entry.NotCaptured = "the 200 response body did not reduce to a type skeleton: " +
			reduceErr.Error() + "; the bytes were not committed rather than committed as an " +
			"empty shape, which would read as an endpoint that answered with nothing"
		return nil, entry, nil
	}
	skeleton, encodeErr := encodeSkeleton(reduced)
	if encodeErr != nil {
		return nil, entry, fmt.Errorf("live-probe: %s: a reduced tree did not marshal, which no "+
			"reduced tree should: %w", ep.Fixture, encodeErr)
	}
	entry.Skeleton = liveKey(ep.Fixture)
	return skeleton, entry, nil
}

// recordDecode unmarshals raw into the response type conformance.SDKTypes records
// for symbol, and sets entry.DecodedCleanly and entry.DecodeErr from the result.
//
// The type comes from the same table the conformance gate compares against, for
// two reasons: it is the type the SDK method actually hands its transport, so
// the answer is about the SDK rather than about a type the probe chose; and
// using the same table means a live capture and a documentation comparison
// cannot disagree about what the SDK decodes into.
//
// Every failure is recorded on the entry rather than returned. A body the SDK's
// own type rejects is a finding, and a run that aborted on the first one would
// capture none of the endpoints after it. A symbol with no type in the table is
// recorded too, with a message that says so plainly, because "no response type
// is recorded for this symbol" is a limit of the probe and must not be readable
// as a defect in the endpoint.
func recordDecode(entry *ManifestEntry, symbol string, raw []byte) {
	subject, _, named := conformance.Subject(symbol)
	if !named {
		entry.DecodeErr = "the documentation manifest maps " + symbol + " to no single SDK " +
			"symbol, so no response type could be named; this is a limit of the probe and " +
			"not a statement about the endpoint"
		return
	}
	table, found := conformance.SDKTypes[subject]
	if !found {
		entry.DecodeErr = "the conformance symbol table records no response type for " + subject +
			", so this run could not attempt the decode; this is a limit of the probe and not " +
			"a statement about the endpoint"
		return
	}
	target, err := table.DecodeTarget()
	if err != nil {
		entry.DecodeErr = "the symbol table entry for " + subject + " does not resolve to a " +
			"type: " + err.Error()
		return
	}
	if target == nil {
		entry.DecodeErr = "the conformance symbol table records that " + subject + " decodes no " +
			"response body, so there was nothing to unmarshal into"
		return
	}
	// reflect.New gives a *T for a type T and a **T for a pointer type, and
	// encoding/json walks either, so this covers a table entry that names a
	// pointer as well as one that names a value.
	if err := json.Unmarshal(raw, reflect.New(target).Interface()); err != nil {
		entry.DecodeErr = decodeFailure(err)
		return
	}
	entry.DecodedCleanly = true
}

// --------------------------------------------------------------------------
// The capture run
// --------------------------------------------------------------------------

// runCapture is phase 2: it captures the response of every endpoint phase 1
// recorded as a 200 and writes the live tree under root.
//
// The whole run is refused when the mutation gate is open. A capture re-calls
// each 200 endpoint, so an open gate would put the endpoints that can place an
// order and move money into the set of endpoints this phase calls - and
// [Capture] would then refuse each of them individually, producing a manifest
// that silently omits them. Refusing the run says so once, before any call,
// which is the difference between an operator who was told and an operator who
// has to notice.
//
// The census is a precondition rather than an input: runCapture receives the
// outcomes phase 1 produced and calls each 200 endpoint a second time. Two calls
// to one endpoint is a deliberate cost, and it buys the thing that matters -
// each captured body was fetched against a request the probe built from the
// documentation, with the same account, the same query and the same body
// construction the census used, so a divergence between the two statuses is a
// fact about the server and not a fact about two different requests.
//
// The walk is sequential, for the reason Census is: the rate limits are per App
// Key, and a concurrent capture would turn a run into a self-inflicted 429.
func runCapture(ctx context.Context, cl *client.Client, endpoints []Endpoint, outcomes []Outcome,
	summary CensusSummary, root, host, region, environment string) error {
	if gate := ResolveMutationGate(host); gate.Open() {
		return fmt.Errorf("live-probe: refusing to capture: the mutation gate is open (%s). A "+
			"capture run re-calls every endpoint the census recorded as a 200, and with the "+
			"gate open that set includes the endpoints that can place an order, cancel one or "+
			"move money. The gate is for a census; this run does not need it", mutationNote("live", gate))
	}

	SetProbedAt(time.Now())
	defer SetProbedAt(time.Time{})

	answered := 0
	for _, o := range outcomes {
		if o.Status == http.StatusOK {
			answered++
		}
	}
	// The census is printed before the capture so the number an operator is about
	// to believe - 55 - can be read against the 193 it came out of, and so the
	// four classes stay apart: a 200 is an answer, a blocked row is a request
	// this probe could not build, a refused row is one it declined to send, and
	// none of the three is the others.
	fmt.Printf("capture: host %s  region %s  environment %s\n", host, region, environment)
	fmt.Printf("census: %d endpoints, %d answered, %d blocked, %d refused by the mutation "+
		"gate, %d unanswered\n",
		summary.Total, summary.Reachable, summary.Blocked, summary.Skipped, summary.Unanswered)
	fmt.Printf("capture: only the %d endpoints that answered HTTP 200 are re-called; every "+
		"other status is left in the census, because an error body compared with a "+
		"documented success response manufactures a finding\n\n", answered)

	skeletons := make(map[string][]byte)
	entries := make([]ManifestEntry, 0, len(outcomes))
	var failures []string
	for i, o := range outcomes {
		if o.Status != http.StatusOK {
			continue
		}
		if i >= len(endpoints) || endpoints[i].Fixture != o.Fixture {
			failures = append(failures, fmt.Sprintf("the census and the endpoint inventory "+
				"disagree at position %d, so this 200 was not attributed to an endpoint", i))
			continue
		}
		skeleton, entry, err := Capture(ctx, cl, endpoints[i], o)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if skeleton != nil {
			skeletons[entry.Skeleton] = skeleton
		}
		entries = append(entries, entry)
	}

	doc := NewLiveManifest(entries, skeletons)
	doc.Host = host
	doc.Region = region
	doc.Environment = environment
	doc.Census = &summary

	if err := WriteLiveCapture(root, skeletons, doc); err != nil {
		return err
	}
	printCapture(root, doc, failures)
	if len(failures) > 0 {
		return fmt.Errorf("live-probe: %d endpoint(s) the census recorded as a 200 could not be "+
			"captured, so %s is short of the %d rows the census implies; the tree was still "+
			"written with the %d it did capture: %s",
			len(failures), liveManifestName, doc.Totals.Considered, doc.Totals.Captured,
			strings.Join(failures, "; "))
	}
	return nil
}

// printCapture reports what the run wrote, in the order an operator needs it:
// what it did, what it did not, and which rows the SDK could not decode.
func printCapture(root string, doc LiveManifest, failures []string) {
	t := doc.Totals
	fmt.Printf("capture: wrote %d skeleton(s) and %s under %s\n", t.Captured, liveManifestName, root)
	fmt.Printf("  considered %d   captured %d   not captured %d   decoded cleanly %d   "+
		"decode rejected %d   skeleton bytes %d\n",
		t.Considered, t.Captured, t.NotCaptured, t.DecodedCleanly, t.DecodeRejected, t.SkeletonBytes)
	if len(t.NotCapturedReasons) > 0 {
		fmt.Println("\nnot captured, by reason:")
		for _, reason := range sortedKeysOfCounts(t.NotCapturedReasons) {
			fmt.Printf("  %4d  %s\n", t.NotCapturedReasons[reason], reason)
		}
	}
	var rejected []ManifestEntry
	for _, e := range doc.Entries {
		if !e.DecodedCleanly {
			rejected = append(rejected, e)
		}
	}
	if len(rejected) > 0 {
		fmt.Printf("\nthe SDK response type did not accept the live body for %d endpoint(s). "+
			"This is the finding this phase exists to collect: conformance can show a DTO "+
			"matches the documentation, and only a live body can show it matches the server.\n",
			len(rejected))
		for _, e := range rejected {
			fmt.Printf("  %-64s %-28s %s\n", e.Fixture, e.Symbol, e.DecodeErr)
		}
	}
	if len(failures) > 0 {
		fmt.Printf("\ncould not be captured: %d\n", len(failures))
		for _, f := range failures {
			fmt.Printf("  %s\n", f)
		}
	}
}

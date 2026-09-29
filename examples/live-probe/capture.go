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
	"strconv"
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
// response reaches a file: not a length of the original body, not the server's
// message, not the account id, symbol, price or timestamp the request carried,
// and not the server's Content-Type, which is recorded as a classification this
// probe recognises rather than as the string the server sent (see
// describeContentType). A skeleton in the repository is a file of "1", -1, true
// and null, so committing 55 of them puts 55 shapes in git and zero readings
// beyond the placeholders themselves.
//
// The one qualification is a member name, and it is stated here because this
// file is where a reader looks first. The reduction keeps every name verbatim, so
// a name is the one place a reading can survive it, and the leak gate checks
// names as well as leaves: a key that is a JSON number, or one of the four
// placeholders, is refused. A key that is a string reading - a ticker, an order
// id - is not separable from a member name by any rule over the document, since
// in JSON it is the same token; keyCarriesAReading states that limit rather than
// implying the tree excludes one.
//
// Two trees now sit under conformance/testdata. The documentation fixtures are
// derived from Webull's published OpenAPI JSON by tools/conformance, and
// conformance.CompareBody reads them. The live skeletons are derived from the
// sandbox, and conformance.CompareLive reads those -- as a second body for the
// same five checks, with the findings recorded in conformance/live-divergences.json
// rather than in the SDK's own backlog. They are kept in separate subtrees with
// separate manifests so a reader scanning the directory can tell which is which
// without opening a file, and so a change to one can never be mistaken for drift
// in the other.
//
// The dependency runs this way and not the other: this program already imports
// the conformance package for the symbol table, so the writer and the reader share
// one encoder and one survey rather than holding a copy of each. A reimplementation
// here could disagree with the reader with no test failing at the boundary, and the
// committed file it produced would be a tree the harness could not read.

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
// own type rather than its message (see decodeFailure), and the three fields
// that could plausibly have held a reading - a body, a length, an account id -
// do not exist on this struct. NotCaptured quotes no part of the response
// either: the one place a server string would have reached it, a
// Content-Type, is recorded as a classification this probe recognises (see
// describeContentType). What does reach NotCaptured verbatim is a member name
// and an SDK symbol, and both are code, not readings: member names are the
// evidence the tree exists to carry.
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
	// WithoutElementShape names the captured skeletons that record no array
	// element shape: every array in them is empty, so the file carries the name
	// of an envelope or of a list and nothing about what a list holds.
	// trading/GET-trading-orders-get.json is `{"orders": []}` and is the clearest
	// case. The list is here because Captured counts files and not shapes, so a
	// reader who divides the row count by the files has read would otherwise take
	// all 55 for endpoints that answered with data. Sorted, so the field is a
	// function of the run.
	WithoutElementShape []string `json:"withoutElementShape,omitempty"`
}

// LiveSize records what the tree's bytes are spent on, in the shape
// conformance/testdata/manifest.json uses for its own sizeTripwire block.
//
// The documentation fixtures and this tree are not the same artifact and the
// reason is recorded rather than borrowed: a fixture is a minimal instance of a
// published schema, so its size is a property of the generator, while a skeleton
// is whatever the sandbox answered, so its size is a property of the responses.
// No bound is imposed here and none is proposed; the block exists so a reader who
// finds the tree large has the composition of the bytes instead of a guess.
type LiveSize struct {
	// ArrayElements is every element of every array in the written skeletons.
	ArrayElements int `json:"arrayElements"`
	// DistinctElementShapes is how many distinct reduced shapes those elements
	// take, counted over the whole tree. Two elements are the same shape when
	// their compact encodings are equal, and encoding/json writes a map's members
	// in sorted order, so the comparison does not depend on member order.
	DistinctElementShapes int `json:"distinctElementShapes"`
	// Reason states what the two numbers mean, and is the one place the tree's
	// size is characterised rather than merely reported.
	Reason string `json:"reason"`
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
	// Size records what the tree's bytes are, so a reader who finds it large is
	// not left with a number to guess at. It is beside Totals rather than inside
	// it because SkeletonBytes, which is the size itself, is a total.
	Size LiveSize `json:"size"`
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
//
// It says "no reading other than a placeholder" rather than "no value", and the
// weaker phrase is the accurate one. A placeholder is a fixed literal, so a
// server that sent -1 as a price, or "1" as a symbol, has produced bytes
// identical to the placeholder and no test can separate the two: the literal
// carries one bit, so it carries no information about whether the server sent it.
// What is true, and what the tree rests on, is that every leaf is a constant the
// reduction chose rather than a literal the body carried.
const liveWarning = "These files hold the SHAPE of a live Webull response - every member " +
	"name the sandbox sent, and no value it sent. Every leaf is a placeholder: a string " +
	"is \"1\", a number is -1, a boolean is true, a null is null. No reading in this tree " +
	"is recoverable other than a placeholder, so there is no price, size, timestamp, " +
	"account id or symbol here to recover. A placeholder is a fixed literal and a server " +
	"can send the same one, so a placeholder is not proof of what the body held - the " +
	"claim is that every leaf is the constant the reduction chose."

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
			"keeps every member name and replaces every value, and written by " +
			"conformance.EncodeSkeleton, which stores each distinct array element shape once " +
			"beside a count of how many elements took it",
		ContainsValues: false,
		LeafForms: `"1" is a string, -1 is a number, true is a boolean, null is a null; ` +
			"an object reduces to an object and an array to a list of its distinct element " +
			"shapes with a count beside each, so the shape is preserved and the reading is not",
		Determinism: "A skeleton is a function of the response body alone: object members are " +
			"written in sorted order and the array records are sorted by the encoding of the " +
			"shape they hold, so the same body always produces the same bytes and a diff " +
			"means the server sent something different. This manifest is the one file that " +
			"changes on every re-run, because probedAt is a clock reading.",
	}
}

// liveSizeReason states what the size block's two counts mean, in the terms the
// written tree uses. It is a function of the counts rather than a constant because
// it names them, and a sentence with the numbers transcribed into it goes stale
// the moment a capture changes the tree -- which is the drift this package exists
// to end, applied to its own prose.
//
// It corrects a claim this tree previously implied and did not make. Bounding an
// array and moving the tree are not the only two things that could be done to the
// size, and treating them as the whole set was wrong: the numbers below show where
// the bytes are, and what they rule in as well as out. Nothing here proposes an
// action, and no bound is imposed.
func liveSizeReason(elements, shapes int) string {
	return "arrayElements is every element of every array the skeletons expand to, and " +
		"distinctElementShapes is how many distinct reduced shapes those elements take. The " +
		"written tree holds each distinct shape once beside a count of how many elements " +
		"took it, so the two numbers are the tree's evidence rather than its size: " +
		groupThousands(elements) + " elements are written as " + strconv.Itoa(shapes) +
		" shapes. The sandbox answers with hundreds of rows of a handful of shapes where one " +
		"row would carry the same names and kinds, and the tree stores the shapes and the " +
		"counts rather than the repetition, which is lossless for member names and JSON " +
		"kinds because nothing was dropped and the counts still sum to the original length. " +
		"These files are live data rather than minimal instances, so the documentation " +
		"fixtures' sizeTripwire does not apply to them; the tree's size is a fact about the " +
		"responses that produced it and no bound is imposed on it here."
}

// groupThousands renders n with a comma every three digits, so the sentence
// liveSizeReason builds reads the way a reader of 8,182 elements expects. It is a
// helper rather than a format string because the manifest's own fields carry raw
// integers and this prose is the only place a grouped figure belongs.
func groupThousands(n int) string {
	digits := strconv.Itoa(n)
	if len(digits) <= 3 {
		return digits
	}
	var b strings.Builder
	for i, r := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// NewLiveManifest assembles the manifest for one capture run from its rows and
// the skeletons they produced.
//
// The counts are derived rather than passed in, because a total a caller
// supplies and a total the rows do not agree on is a manifest that misreports
// its own tree. entries is copied and sorted, so a caller's slice is neither
// reordered nor aliased.
//
// The size block and the WithoutElementShape list are derived from the skeleton
// bytes for the same reason: they describe the tree, so reading them off the tree
// is the only way they cannot disagree with it. The reason block is rendered from
// the two counts it names, so a sentence in it cannot survive a capture that
// changed them.
func NewLiveManifest(entries []ManifestEntry, skeletons map[string][]byte) LiveManifest {
	doc := LiveManifest{
		Kind:      liveKind,
		Warning:   liveWarning,
		Generator: liveGeneratorBlock(),
		Scope:     liveScope,
		Entries:   append([]ManifestEntry(nil), entries...),
	}
	sort.Slice(doc.Entries, func(i, j int) bool { return doc.Entries[i].Fixture < doc.Entries[j].Fixture })
	// One set across the whole tree, so DistinctElementShapes is the number of
	// shapes the tree holds rather than the sum of the number each file holds.
	shapes := map[string]bool{}

	for _, e := range doc.Entries {
		doc.Totals.Considered++
		if e.Skeleton != "" {
			doc.Totals.Captured++
			raw, present := skeletons[e.Skeleton]
			doc.Totals.SkeletonBytes += len(raw)
			if !present {
				// checkIndexed refuses this before any write; returning a manifest
				// that silently omits the accounting is worse than recording what
				// was actually there.
				continue
			}
			found, elements, carries := surveySkeleton(raw)
			for shape := range found {
				shapes[shape] = true
			}
			doc.Size.ArrayElements += elements
			if !carries {
				doc.Totals.WithoutElementShape = append(doc.Totals.WithoutElementShape, e.Skeleton)
			}
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
	sort.Strings(doc.Totals.WithoutElementShape)
	doc.Size.DistinctElementShapes = len(shapes)
	doc.Size.Reason = liveSizeReason(doc.Size.ArrayElements, doc.Size.DistinctElementShapes)
	return doc
}

// surveySkeleton reads one committed skeleton and reports the three facts the
// manifest records about it: the distinct array element shapes it holds, how many
// array elements it holds in all, and whether it holds any element at all.
//
// It delegates to conformance.SurveySkeleton, and the delegation is the point
// rather than a convenience. The manifest describes the tree the reader sees, and
// the reader sees the deduplicated form, so the accounts have to be computed over
// the *expanded* tree: read the committed bytes directly and a record holding a
// count of 997 counts as one element, so every multi-element array in the tree
// would be undercounted and totals.arrayElements would describe something the
// file does not hold.
//
// That is the same reason the encoder is conformance.EncodeSkeleton. Two
// implementations of one definition, one in this program and one in the reader,
// can disagree with no test failing at the boundary: a probe that walked the
// committed form while the reader walked the expanded one would each be internally
// consistent and jointly wrong. The dependency already runs this way -- this file
// imports the conformance package for Subject and SDKTypes -- so sharing the
// reader's definitions costs nothing in structure and removes the class.
//
// The shapes are keyed by the compact encoding of the element, so two elements
// with the same shape have the same key and the count does not depend on the order
// the members happen to appear in. A skeleton that cannot be read is reported as
// carrying an element, and a skeleton with no array at all is reported as carrying
// one too: both are the conservative direction, because the field this feeds is a
// list of skeletons whose every array is empty, and adding a skeleton that holds a
// member-name-and-kind shape to that list would tell a reader it holds nothing when
// it holds something.
func surveySkeleton(raw []byte) (shapes map[string]bool, elements int, carriesElement bool) {
	survey := conformance.SurveySkeleton(raw)
	return survey.DistinctElementShapes, survey.ArrayElements, survey.CarriesElementShape
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
// `os.WriteFile` at each call site. Every key is resolved and the manifest is
// held against the file set before a byte is written, so a run refused by either
// of those checks leaves the tree exactly as it found it rather than half-updated.
// Then the skeletons are written; then the manifest; then the tree is walked to
// prove it holds nothing but what this run wrote. That last check is the one that
// runs after the writes, so a run failing it has already updated the tree and left
// the file it objects to in place. That is deliberate and it is not the same
// guarantee: a file under live/ that the manifest does not name is an error and
// not a deletion, because it is either evidence from an earlier run that this one
// failed to reproduce, which a reader must look at, or a file somebody put there
// by hand. So the two failure classes differ - a bad key or an unindexed skeleton
// writes nothing, a stale file leaves the tree written and asks a question - and
// the error says which happened by naming the file.
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

// mediaTypeOf strips a Content-Type's parameters and normalises it, so a
// classification below is about the media type rather than about its spelling.
func mediaTypeOf(header string) string {
	mediaType := header
	if index := strings.IndexByte(mediaType, ';'); index >= 0 {
		mediaType = mediaType[:index]
	}
	return strings.ToLower(strings.TrimSpace(mediaType))
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
	mediaType := mediaTypeOf(header)
	switch mediaType {
	case "", "application/json", "text/json", "application/x-json", "text/plain":
		return true
	}
	// The vendor and structured-suffix forms: anything/json,
	// application/vnd.webull+json.
	return strings.HasSuffix(mediaType, "json")
}

// contentTypeClassifications is the fixed set of media types describeContentType
// is willing to name in a committed manifest.
//
// The set is the media types a 200 has actually been answered with in the error
// cases, not a registry: this is a classification for a human reading a manifest,
// and inventing a token for a form nobody has seen would be a name with no
// evidence behind it. An entry's key and its value are the same string, which is
// deliberate - the table records the media type's own spelling, so a reader can
// match it against what the endpoint declares without this file restating it.
var contentTypeClassifications = map[string]string{
	"text/html":                         "text/html",
	"application/xhtml+xml":             "application/xhtml+xml",
	"text/xml":                          "text/xml",
	"application/xml":                   "application/xml",
	"application/octet-stream":          "application/octet-stream",
	"application/x-www-form-urlencoded": "application/x-www-form-urlencoded",
	"text/csv":                          "text/csv",
	"text/event-stream":                 "text/event-stream",
	"application/pdf":                   "application/pdf",
	"application/zip":                   "application/zip",
	"image/png":                         "image/png",
	"image/jpeg":                        "image/jpeg",
	"image/gif":                         "image/gif",
}

// describeContentType returns the token a manifest records in place of a
// Content-Type, or a sentence saying the media type is outside the set above.
//
// The header is server-supplied and the manifest is committed, so quoting it would
// put an unvalidated server-supplied string in the repository - the same hazard as
// a value in a skeleton, one layer up, and the reason the probe already
// reconstructs a decode error rather than quoting it. The finding does not need
// the raw string: what a reader has to know is which media type the endpoint
// answered with, and that is a value from a fixed set. So a recognised form is
// named from this file's own table and everything else is reported as unnamed.
//
// Nothing derived from an unrecognised header reaches the result, which is the
// stricter half of the rule: the server's own spelling is left out entirely,
// because quoting an unrecognised token would reintroduce the server-supplied
// string this function exists to remove. The cost is that a reader learns an
// endpoint answered with a media type the probe does not recognise rather than
// which one it was, and the answer to that is the gateway, not the manifest.
//
// A missing header is its own token rather than a member of the set: an endpoint
// that answered 200 with no declaration is a different fact from one that answered
// with a recognised type, and lumping it into "unnamed" would hide which happened.
func describeContentType(header string) string {
	mediaType := mediaTypeOf(header)
	if mediaType == "" {
		return "no media type at all"
	}
	if name, known := contentTypeClassifications[mediaType]; known {
		return name
	}
	return "a media type this probe does not name"
}

// decodeFailure renders why a live body did not unmarshal into an SDK response
// type, without quoting anything the body said.
//
// This is the first of two places in the probe where a server's bytes could reach
// a committed file through an error message. json.Unmarshal returns three error
// types of its own, and two of them are safe to render from their fields: an
// UnmarshalTypeError names the JSON kind, the field path and the Go type it
// wanted, all three of which are properties of the SDK type and the wire names
// rather than of any reading, and a SyntaxError names a byte offset. The third is
// everything else - and everything else is the problem. encoding/json returns an
// UnmarshalJSON error verbatim, so a DTO that fails to parse a price can return an
// error containing the price, and money.Money.UnmarshalJSON is exactly such a DTO.
// Rendering that message would put a live reading in a committed manifest, so it is
// classified by type and not quoted.
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
	if rendered, ok := renderSyntaxError(err); ok {
		return rendered
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

// renderSyntaxError renders the offset of a *json.SyntaxError and nothing else.
//
// It is shared with reduceFailure because both paths can reach the same error type
// and because the reason a syntax error is safe is a property of the type: a
// SyntaxError carries an Offset and its message quotes the offending byte, so the
// offset is the whole of it that is safe and the message is the whole of it that
// is not.
func renderSyntaxError(err error) (string, bool) {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return fmt.Sprintf("the body is not well-formed JSON, at byte offset %d", syntaxErr.Offset), true
	}
	return "", false
}

// reduceFailure renders why a live body did not reduce to a type skeleton, without
// quoting anything the body said.
//
// It exists because decodeFailure closed the leak on one path only, and this is
// the other one. It is reachable: isJSONMediaType accepts a missing header and
// text/plain, so a 200 carrying an HTML error page or an unlabelled body lands
// here rather than being refused by the media type, and then a json.SyntaxError
// carrying the message "invalid character '<' looking for beginning of value"
// puts one byte of the live body into a committed manifest. Nothing in the
// committed tree exercises it - every captured body reduced - which is exactly why
// a latent path is the one worth closing rather than the one a commit already
// shows.
//
// The classification is the same shape as decodeFailure's and for the same reason:
// a SyntaxError's offset is a property of the document's position and not of its
// content, an UnmarshalTypeError cannot arise from decoding into `any`, and
// anything else is named by type and never quoted.
func reduceFailure(err error) string {
	if rendered, ok := renderSyntaxError(err); ok {
		return rendered
	}
	// A truncated document is not a SyntaxError: the decoder runs out of input
	// before the value closes and returns io.ErrUnexpectedEOF, whose own message
	// is value-free but which says nothing a reader could act on. Naming it here
	// is what tells "the body stopped mid-document" apart from "the body is not
	// JSON at all", which are different statements about an endpoint.
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "the body ended before the JSON document closed"
	}
	// The probe's own refusal, rendered by identity rather than by text: it is
	// this file's prose and carries nothing from the body, so quoting it is safe
	// and it names a real finding - a body with two documents in it.
	if errors.Is(err, errTrailingData) {
		return "the body holds more than one JSON document"
	}
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		found := typeErr.Value
		if found == "" {
			found = "value"
		}
		return fmt.Sprintf("the reduction read a JSON %s it cannot represent at %s",
			found, typeErr.Field)
	}
	return fmt.Sprintf("the body did not reduce to a type skeleton (%T); the decoder's own "+
		"message is not recorded because a JSON syntax error quotes a byte of the body",
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
		entry.NotCaptured = "the 200 response declares " + describeContentType(contentType) +
			", which is not a JSON media type, so there is no response shape to reduce and " +
			"the bytes were not committed"
		return nil, entry, nil
	}
	reduced, reduceErr := Skeletonify(raw)
	if reduceErr != nil {
		entry.NotCaptured = "the 200 response body did not reduce to a type skeleton: " +
			reduceFailure(reduceErr) + "; the bytes were not committed rather than committed as an " +
			"empty shape, which would read as an endpoint that answered with nothing"
		return nil, entry, nil
	}
	// conformance.EncodeSkeleton is the encoder, not a second one beside it. The
	// committed form is deduplicated -- an array is written as its distinct element
	// shapes with a count beside each -- and that form is a shared contract: the
	// conformance package reads it, its leak gate walks it, and its canonicality
	// check re-encodes it. A local marshaller would write the expanded form over a
	// tree the reader expects to be deduped, and the only signal would be a
	// 1.9 MB diff that looks like evidence and is not one.
	skeleton, encodeErr := conformance.EncodeSkeleton(reduced)
	if encodeErr != nil {
		return nil, entry, fmt.Errorf("live-probe: %s: a reduced tree did not encode, which no "+
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
// This is the second of two refusals, and the first is in main's run, before the
// token exchange and before the walk. A check that only lived here would be
// reached with the census already taken, so it would protect the capture and
// leave the 34 mutating requests the walk made behind it. Two checks, one
// message, and neither of them claims the other is unnecessary.
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
		return captureRefusal(gate)
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

// captureRefusal is the error a -capture run exits with when the mutation gate
// is open. Both refusal sites build it, so the sentence an operator reads is the
// same one whichever of the two caught the run.
//
// It states the conditions that opened the gate rather than reusing
// mutationNote, and that is not a preference. mutationNote is the line a census
// prints after the walk, so it has to report what the walk did and it reads
// "every endpoint in the mutating class was called". A refusal has to say the
// opposite: the two call sites are at different points in the run, and the
// message that says nothing was sent would be false of the census and the
// message that says everything was sent would be false of the refusal. So the
// two report the same gate in their own tense.
//
// It also says what the run does not have, because a refused -capture produces
// no census and no skeleton and an operator who expected a file needs to know
// that its absence is the refusal rather than a crash.
func captureRefusal(gate MutationGate) error {
	return fmt.Errorf("live-probe: refusing to capture: the mutation gate is open, because %s. "+
		"A capture re-calls every endpoint a census recorded as a 200, and with the gate open "+
		"that set includes the endpoints that can place an order, cancel one or move money. "+
		"This refusal is made before the token exchange, before the census walk and before the "+
		"capture, so no request was sent and this run produced neither a census nor a skeleton. "+
		"Close the gate - unset %s, or point the run at a sandbox host - and run the capture again",
		openGateConditions(gate), MutateOptInEnv)
}

// openGateConditions states which of the gate's conditions held, and which
// variable decides each, for a gate that is already known to be open.
//
// Only the open case is spelled, because a closed gate cannot be the reason for
// this sentence: Open() requires the opt-in, so the two branches below are the
// only two ways to get here. The two failures the gate does not open on are
// reported by mutationNote, which is the line a census prints.
func openGateConditions(gate MutationGate) string {
	if gate.SandboxHost {
		return MutateOptInEnv + "=" + MutateOptInValue + " is set and " + gate.Host +
			" is a sandbox host the SDK derives for a region it knows"
	}
	return MutateOptInEnv + "=" + MutateOptInValue + " and " + NonSandboxOverrideEnv + "=" +
		NonSandboxOverrideValue + " are set and " + gate.Host + " is NOT a sandbox host"
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

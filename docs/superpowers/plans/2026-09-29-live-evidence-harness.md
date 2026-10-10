# Live Evidence Harness Implementation Plan

> ### Completed — all 32 steps done
>
> The live-evidence harness was built entirely in the v2.1.36 release cycle.
> Implementation ran outside the OpenSpec workflow and closed all checkboxes manually.
> See the v2.1.36 and v2.1.37 changelog entries for the full commit list.

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A probe that calls every documented endpoint with a synthesised minimum valid request, records what answered, and — for those that answered — compares a value-free type skeleton of the live response against both the SDK type and the documentation-derived fixture.

**Architecture:** A nested Go module under `examples/live-probe/` calls the SDK through `client.WithEnv()` and writes two value-free artefacts: a per-endpoint type skeleton and a live manifest. A new `conformance/live.go` reads those and runs the *existing* `conformance.CompareBody` in two directions, emitting a third divergence set kept out of `known-divergences.json`. Phases 0, 1 and 4 ship first; phases 2 and 3 are gated on phase 4's SDK-disagreement count.

**Tech Stack:** Go 1.26, `encoding/json`, `reflect`, `net/http`; reuses `conformance.CompareBody`, `conformance.Describe`, `client.WithEnv`, `data`/`trade`/`brokerfd`/`broker` clients.

**Spec:** `docs/superpowers/specs/2026-09-29-live-evidence-harness-design.md`

## Global Constraints

- **No value from a live response may reach a commit.** The committed artefact is a type skeleton: names, JSON kinds and nesting only. A skeleton contains no values, not even placeholders.
- **Credentials are read from the environment only**, via `client.WithEnv()`. Never commit an app key, secret or account ID. `AGENTS.md` forbids inlining them, and a rotated shared credential would otherwise surface as conformance findings.
- **`string` and `number` are distinct skeleton kinds and must never be conflated.** `money.Money.UnmarshalJSON` and `data.QuoteTime` both accept either, so SDK permissiveness cannot be allowed to mask what the server actually sent.
- **Existing conformance code is not modified.** The live path reuses `CompareBody` as-is. `known-divergences.json` is not written by any task here.
- **Apache-2.0 header** on every new hand-written `.go` file, copied from a neighbouring file. `Copyright 2026 shing1211`.
- **gofmt clean** for the whole tree, including the new nested module, on every commit.
- **`goimports` local-prefixes `github.com/shing1211/webullapi4go`**: stdlib first, blank line, then local module imports last.
- **Licence headers and GoDoc on every exported identifier**, GoDoc starting with the identifier name and explaining behaviour rather than restating the signature.
- **Non-goal:** fixing any of the 25 remaining baseline rows. Credentials cannot settle them.
- **A type skeleton is already a valid fixture input, and that is why Task 5 needs no new comparison code.** `conformance.jsonKind` (`shapes.go:361`) classifies decoded values by kind, and `leafSamples` (`shapes.go:1068`) already reasons with kind-level literals like `"1"` and `1`. So a skeleton fed to `CompareBody` produces the same names and leaf kinds a real body would. An implementer who writes a second comparison routine is duplicating five checks that already work.
- **Request bodies must be synthesised, not omitted.** 50 of the 193 endpoints require a request body and 127 require a query parameter. An endpoint called with no body is one the probe failed to ask properly, not one that did not answer.
- **Unexported response types are already bridged and need no new work.** `conformance/envelopes.go` exists because six response envelopes are unexported, and it already exposes them to the symbol table. The probe never names a method's return type, so capture is unaffected and `live.go` inherits the bridges for free.
- **Reachability must not be pre-filtered by a path heuristic.** `/market-data/screeners/*` returns 404 in the HK sandbox while `/market-data/stocks/*` works, so any prefix-based guess is wrong in both directions. The census *is* the measurement; treating the probe as needing a prior reachability list is the mistake.
- **Windows:** set `GOTMPDIR` and `GOCACHE` under `C:\Users\Tchan\AppData\Local\Temp\opencode` before Go commands, or a build may fail on `a.out.exe` file locks.

## Review Focus

Input classes and conditions the spec implies that a naive implementation would get wrong. Each is pinned by a test in the owning task.

1. **A numeric field arrives as a JSON number and the skeleton must record `number`, not `string`.** A caller comparing the skeleton against `money.Money` expects the real wire kind, because `Money` accepts both and would hide the difference.
2. **A field present in the live body with a `null` value.** The skeleton must record the name with kind `null` rather than omitting it, or a field the SDK lacks looks absent rather than unreadable.
3. **A credential that has been rotated.** The probe must report it as an auth failure at the top of the report, never as a set of missing-field findings.
4. **A live body that is a JSON array where the page documents an object** (or the reverse) — the top-level kind disagreement that the remaining 25 rows are made of.
5. **An endpoint that answers `200` with an empty array.** This is an answer, not a failure, and must be recorded as such rather than treated as unreachable.

---
---

## File Structure

**Created — probe (nested module, `examples/live-probe/`):**

| Path | Responsibility |
|---|---|
| `go.mod` | Module `github.com/shing1211/webullapi4go/examples/live-probe`, with a `replace` to the repo root, mirroring `examples/broker-probe/go.mod`. |
| `skeleton.go` | `Skeletonify(raw []byte) (any, error)` — reduces a decoded body to names, kinds and nesting. The only file that knows how a kind is named. |
| `skeleton_test.go` | Table tests for `Skeletonify`, including the number/string and null cases. |
| `params.go` | `Param(name string, spec ParamSpec) any` — the §5.5 resolution order. |
| `params_test.go` | Table tests per parameter class. |
| `census.go` | Per-endpoint call loop, account discovery, outcome recording. |
| `main.go` | Flag parsing, client construction, phase sequencing, report output. |

**Created — harness:**

| Path | Responsibility |
|---|---|
| `conformance/live.go` | Loads skeletons, runs `CompareBody` in two directions, writes `live-divergences.json`. |
| `conformance/live_test.go` | Tests the loader and both comparison directions. |
| `conformance/testdata/live/` | Committed skeletons, one per reachable endpoint. |
| `conformance/testdata/live-manifest.json` | Provenance: host, timestamp, per-endpoint status. |

**Read-only inputs:** `conformance/testdata/manifest.json` (endpoint list, `Documented.Method`, `Documented.Path`), `conformance/testdata/*.json` (the documentation-derived fixture each skeleton is compared against).

**Not created:** anything that writes `known-divergences.json`.

---

## Task 1: `Skeletonify` — the value-free reduction

The foundation everything else depends on, and the only component that decides what a "kind" is. It must be right before any network call happens, so it comes first and is tested exhaustively with no credential required.

**Files:**
- Create: `examples/live-probe/go.mod`
- Create: `examples/live-probe/skeleton.go`
- Test: `examples/live-probe/skeleton_test.go`

**Interfaces:**
- Consumes: nothing. Standard library only.
- Produces:
  - `func Skeletonify(raw []byte) (any, error)` — decodes `raw` with `json.Decoder` using `UseNumber()`, returns the reduced tree as `map[string]any`, `[]any`, or a **typed placeholder** leaf, and returns an error only when `raw` is not valid JSON. **The leaf representation was revised during Task 1**; the tree is no longer spelled with kind strings. The binding form is in the Task 1 report at `.superpowers/sdd/2026-09-29-live-evidence-harness/task-1-report.md`, and it is the form a later task must code against. Briefly: a JSON string becomes the Go string `"1"`, a JSON number becomes a fixed synthetic `json.Number`, a JSON boolean becomes `true`, and a JSON null becomes a nil interface; an empty array is `[]any{}` and never nil. Two properties drive it and both matter downstream. `UseNumber()` is still required at decode time, or a body carrying a value outside `float64` range (`1e400`, a 40-digit integer) is a decode error and the whole body is silently dropped from the evidence set. And the number leaf is synthetic rather than the server's literal text, which is what makes a reduced tree safe to serialise into the committed files Task 4 writes — `conformance.jsonKind` classifies by Go type, not by a number's digits, so the kind is readable straight off the placeholder and a committed fixture can contain no value the server sent.
  - `func SkeletonKind(v any) string` — the kind name for a reduced value; used by `census.go` to print a one-line summary.

- [x] ** Create the module and write the failing test**

`go.mod` (mirror `examples/broker-probe/go.mod`):

```
module github.com/shing1211/webullapi4go/examples/live-probe

go 1.26

require github.com/shing1211/webullapi4go v0.0.0

replace github.com/shing1211/webullapi4go => ../..
```

`skeleton_test.go`, table-driven, with these cases — **the first two are the Review Focus cases 1 and 2 and must be explicit**:

```go
func TestSkeletonify(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want any
	}{
		{"object", `{"a":"1"}`, map[string]any{"a": "string"}},
		// A number must stay a number: Money and QuoteTime accept both forms, so
		// conflating them is the one mistake this function exists to prevent.
		{"number is not string", `{"close":385.6}`, map[string]any{"close": "number"}},
		// A present-but-null field is a name the SDK may not carry, so it must
		// survive as "null" rather than vanish.
		{"null is recorded", `{"outstanding":null}`, map[string]any{"outstanding": "null"}},
		{"array of objects", `[{"a":1}]`, []any{map[string]any{"a": "number"}}},
		{"empty array is not a wrong kind", `[]`, []any{}},
		{"bool", `{"ok":true}`, map[string]any{"ok": "boolean"}},
		{"nested", `{"d":{"e":[]}}`, map[string]any{"d": map[string]any{"e": []any{}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Skeletonify([]byte(tc.in))
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestSkeletonifyRejectsInvalidJSON(t *testing.T) {
	if _, err := Skeletonify([]byte(`{"a":`)); err == nil {
		t.Error("invalid JSON returned no error")
	}
}
```

- [x] ** Run the test to verify it fails**

```
$env:GOTMPDIR="C:\Users\Tchan\AppData\Local\Temp\opencode\gotmp"
$env:GOCACHE="C:\Users\Tchan\AppData\Local\Temp\opencode\gocache"
cd examples/live-probe
go test ./... -run TestSkeletonify -v
```
Expected: FAIL — `undefined: Skeletonify`.

- [x] ** Implement `Skeletonify` in `examples/live-probe/skeleton.go`**

Apache-2.0 header, then `package main`.

The kind names are exactly the strings in the test: `object`, `array`, `string`, `number`, `boolean`, `null`. Use `json.NewDecoder(bytes.NewReader(raw))` with `d.UseNumber()` so an integer is not turned into a float and then indistinguishable from a decimal; `json.Number` maps to `number`.

Arrays reduce to `[]any` of the reduced element. An empty array reduces to `[]any{}`, never `nil`, so an empty body is distinguishable from a wrong top-level kind.

`SkeletonKind` returns `"object"` for a map, `"array"` for a slice, `"string"` for the string placeholder, `"number"` for a `json.Number`, `"boolean"` for a bool, and `"null"` for a nil, and `""` for any other value. It classifies by Go type, so a kind is never an arbitrary string echoed back from the input. **This superseded the brief's original wording** ("the string itself for the kind names"), which failed open on an unreduced string; see the Task 1 report.

- [x] ** Run the test to verify it passes**

```
go test ./... -v
```
Expected: PASS, all subtests.

- [x] ** Commit**

```
cd examples/live-probe
gofmt -l .
git add examples/live-probe
git commit -m "feat(live-probe): reduce a live body to a value-free type skeleton"
```

---

## Task 2: parameter resolution per spec §5.5

The second-largest source of a false negative: a request the probe cannot build is reported as an endpoint that did not answer, when in fact the probe never reached it.

**Files:**
- Create: `examples/live-probe/params.go`
- Test: `examples/live-probe/params_test.go`

**Interfaces:**
- Consumes: nothing from Task 1.
- Produces:
  - `type ParamSpec struct { Type string; Enum []any; Example any; HasExample bool; Minimum, Maximum *float64 }`
  - `func Param(name string, spec ParamSpec) (any, error)` — returns the value to send, and a non-nil error when the class cannot be resolved; `census.go` records that as `blocked: unresolvable parameter` rather than attempting the call.
  - `const DefaultSymbol = "AAPL"` — the sandbox carries AAPL only, per `AGENTS.md`.

- [x] ** Write the failing test**

One case per class in spec §5.5, plus the four that matter most:

```go
func TestParam(t *testing.T) {
	// Class: caller-supplied correlation key. Any unique value is accepted.
	got, err := Param("client_request_id", ParamSpec{Type: "string"})
	if err != nil || got == "" {
		t.Errorf("client_request_id = %v, %v; want a synthesised non-empty value", got, err)
	}
	// Two calls must differ, or a second endpoint's correlation key would collide.
	a, _ := Param("client_request_id", ParamSpec{Type: "string"})
	b, _ := Param("client_request_id", ParamSpec{Type: "string"})
	if a == b {
		t.Error("two correlation keys are identical")
	}
	// Class: enum. First value wins, and the page's own example is preferred.
	got, _ = Param("category", ParamSpec{Type: "string", Enum: []any{"US_STOCK", "HK_STOCK"}})
	if got != "US_STOCK" {
		t.Errorf("enum first = %v, want US_STOCK", got)
	}
	got, _ = Param("category", ParamSpec{Type: "string", Enum: []any{"US_STOCK"}, Example: "HK_STOCK", HasExample: true})
	if got != "HK_STOCK" {
		t.Errorf("documented example = %v, want HK_STOCK: the page is authoritative", got)
	}
	// Class: declared scalar, respecting the published range.
	got, _ = Param("count", ParamSpec{Type: "integer", Minimum: ptr(1.0), Maximum: ptr(20.0)})
	if n, ok := got.(int); !ok || n < 1 || n > 20 {
		t.Errorf("count = %v, want an int within [1,20]", got)
	}
	// Class: domain constant.
	got, _ = Param("symbol", ParamSpec{Type: "string"})
	if got != DefaultSymbol {
		t.Errorf("symbol = %v, want %q", got, DefaultSymbol)
	}
	// Class: caller identity. Not synthesisable — the server issues it.
	if _, err := Param("account_id", ParamSpec{Type: "string"}); err == nil {
		t.Error("account_id resolved without discovery; it must be threaded in, not guessed")
	}
}

func TestParamRejectsUnresolvable(t *testing.T) {
	if _, err := Param("mystery", ParamSpec{}); err == nil {
		t.Error("a spec with no type, enum or example resolved; it must not")
	}
}
```

- [x] ** Run the test to verify it fails**

```
cd examples/live-probe
go test ./... -run TestParam -v
```
Expected: FAIL — `undefined: Param`.

- [x] ** Implement `Param` in `examples/live-probe/params.go`**

Resolution order, each step a short branch:

1. `HasExample` → return `Example`.
2. `len(Enum) > 0` → return `Enum[0]`.
3. `name == "account_id"` → return an error naming account discovery. **This branch must come before the type-based defaults**, or the caller-identity case silently becomes a fabricated string and 34 endpoints are reported against an account that does not exist.
4. `name` in `{symbol, symbols, series_symbol}` → `DefaultSymbol` (or `[]string{DefaultSymbol}` for the plural form the page declares).
5. `Type` in `{integer, number}` → clamp `1` into `[Minimum, Maximum]` when both are present, else `1`.
6. `Type == "boolean"` → `true`.
7. `Type == "string"` and `name` matches `client_request_id` or `client_order_id` → a unique value: a monotonic counter combined with the current nanosecond, so two calls in one run never collide.
8. Otherwise → error.

- [x] ** Run the test to verify it passes**

```
go test ./... -v
```
Expected: PASS.

- [x] ** Commit**

```
gofmt -l .
git add examples/live-probe/params.go examples/live-probe/params_test.go
git commit -m "feat(live-probe): resolve required parameters from the page schema"
```

---

## Task 3: phase 0 and phase 1 — account discovery and the reachability census

The first deliverable, and the one worth having even if nothing after it happens. Ships as a runnable command that writes `census.json`; it commits nothing on its own beyond the report the operator inspects.

**Files:**
- Create: `examples/live-probe/census.go`
- Create: `examples/live-probe/main.go`
- Test: `examples/live-probe/census_test.go`

**Interfaces:**
- Consumes: `Param`, `ParamSpec`, `DefaultSymbol` (Task 2).
- Produces:
  - `type Outcome struct { Symbol, Fixture, Method, Path string; Status int; Err string; Blocked string }`
  - `func Census(ctx context.Context, cl *client.Client, m []Endpoint) ([]Outcome, error)` — the per-endpoint loop.
  - `type Endpoint struct { Symbol, Fixture, Method, Path string; RequiredParams []ParamSpec; HasBody bool }` — one row per manifest fixture.
  - `func DiscoverAccountID(ctx context.Context, cl *client.Client) (string, error)` — phase 0.
  - `main.go` writes `census.json` to the path given by `-out`.

- [x] ** Write the failing test**

The census loop is testable without a credential by pointing it at an `httptest` server, and that is the point: the loop's job is to record outcomes, and a test that needs the sandbox cannot prove it records a 404 correctly.

```go
func TestCensusRecordsNonSuccessStatuses(t *testing.T) {
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.URL.Path)
		switch r.URL.Path {
		case "/ok":
			w.Write([]byte(`{"a":1}`))
		case "/gone":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer srv.Close()

	cl := testClient(t, srv.URL)
	got, err := Census(context.Background(), cl, []Endpoint{
		{Symbol: "ok", Method: "GET", Path: "/ok"},
		{Symbol: "gone", Method: "GET", Path: "/gone"},
		{Symbol: "blocked", Method: "GET", Path: "/nope", RequiredParams: []ParamSpec{{Name: "account_id"}}},
	})
	if err != nil {
		t.Fatalf("Census: %v", err)
	}
	by := map[string]Outcome{}
	for _, o := range got {
		by[o.Symbol] = o
	}
	if by["ok"].Status != http.StatusOK {
		t.Errorf("ok status = %d", by["ok"].Status)
	}
	if by["gone"].Status != http.StatusNotFound {
		t.Errorf("gone status = %d", by["gone"].Status)
	}
	// The unresolvable-parameter case must never reach the network: it is not an
	// endpoint that failed to answer, it is an endpoint the probe could not ask.
	if by["blocked"].Blocked == "" {
		t.Error("blocked endpoint has no Blocked reason")
	}
	if by["blocked"].Status != 0 {
		t.Errorf("blocked endpoint has status %d; it should not have been called", by["blocked"].Status)
	}
	for _, p := range seen {
		if p == "/nope" {
			t.Error("the probe called an endpoint whose parameters it could not resolve")
		}
	}
}
```

- [x] ** Run the test to verify it fails**

```
cd examples/live-probe
go test ./... -run TestCensusRecords -v
```
Expected: FAIL — `undefined: Census`.

- [x] ** Implement the census in `census.go` and `main.go`**

`main.go`: flags `-out` (default `census.json`), `-base` (optional override, defaulting to `client.WithEnv()`), and `-accounts` (optional `account_id` override, which skips phase 0). Construct with `client.New(client.WithEnv())` and call `cl.EnsureToken(ctx)`. **A failure of `EnsureToken` must be reported as a top-level auth failure and exit non-zero, before any endpoint is called** — that is Review Focus case 3, and it is why the credential problem can never be mistaken for conformance findings.

`DiscoverAccountID`: call `trade.New(cl).ListAccounts(ctx)`, return the first `AccountID`. On error, return a wrapped error and let `main` record every `account_id`-dependent endpoint as `blocked: account discovery failed`, which is a more useful statement than a generic 404.

`Census`: for each `Endpoint`, first resolve every `RequiredParams` entry through `Param`, and synthesise the request body when `HasBody` is set — **50 of the 193 endpoints require one, and calling them with an empty body records the probe's failure as the server's.** If any parameter resolution returns an error, append an `Outcome` with `Blocked` set and **Status 0, without making a request**. Otherwise build the query and body, call `cl.Do(ctx, method, path+"?"+query, body, &raw)` with the response captured into `[]byte`, and append an `Outcome` carrying the HTTP status or the error text.

Body synthesis reuses the same leaf rule as `Param`: a required string property takes its documented `example`, then its first `enum` value, then a type default. A body that cannot be built sets `Blocked` rather than sending a wrong one.

Reading the manifest: iterate `Manifest.Fixtures`, take `Documented.Method` and `Documented.Path`, and read required parameters from the cached page via the same `select_block`/`schema_200` logic `gen_fixtures.py` uses. **If the cache is absent, phase 1 exits 0 with a message, matching the existing `conformance-fixtures` behaviour** — a fresh checkout must not fail on a missing cache.

- [x] ** Run the test to verify it passes**

```
go test ./... -v
```
Expected: PASS.

- [x] ** Commit**

```
gofmt -l .
git add examples/live-probe
git commit -m "feat(live-probe): census endpoint reachability and discover the account id"
```

- [x] ** Run it for real and record the result**

```
$env:WEBULL_APP_KEY="<read from the Webull test-accounts page>"
$env:WEBULL_APP_SECRET="<read from the Webull test-accounts page>"
$env:WEBULL_REGION="hk"
$env:WEBULL_ENVIRONMENT="sandbox"
go run . -out ..\..\census.json
```
Expected: a printed table of outcomes and a `census.json`. **Report the reachable count to the human partner and stop for the phase-4 decision before continuing** — if the reachable count is near zero, phases 2 and 3 should not be built at all.

---

## Task 4: phase 2 — skeleton capture

Writes the committed evidence. Only runs against endpoints phase 1 proved reachable.

**Files:**
- Create: `examples/live-probe/capture.go`
- Test: `examples/live-probe/capture_test.go`
- Create at run time: `conformance/testdata/live/<fixture-id>.json`, `conformance/testdata/live-manifest.json`

**Interfaces:**
- Consumes: `Skeletonify` (Task 1), `Param` (Task 2), `Census`/`Outcome` (Task 3).
- Produces:
  - `func Capture(ctx context.Context, cl *client.Client, ep Endpoint, prev Outcome) (skeleton []byte, entry ManifestEntry, err error)`
  - `type ManifestEntry struct { Symbol, Fixture, Host, ProbedAt string; Status int; DecodedCleanly bool; DecodeErr string; Skeleton string }`
  - `func WriteCapture(dir string, skeletons map[string][]byte, entries []ManifestEntry) error` — refuses to write into `conformance/testdata/` if any file would change the documentation fixtures; writes only under `live/`.

- [x] ** Write the failing test**

```go
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
}

func TestCaptureRecordsDecodeFailureWithoutFailing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A body the SDK type cannot hold at all: the skeleton is still valid, and the
		// decode failure is a finding rather than a reason to drop the capture.
		w.Write([]byte(`{"a":1}`))
	}))
	defer srv.Close()
	// ... call Capture with a type known to reject "a", assert Skeleton is non-nil,
	// DecodedCleanly is false, DecodeErr is non-empty, and err is nil.
}
```

- [x] ** Run the test to verify it fails**

```
go test ./... -run 'TestWriteCapture|TestCaptureRecords' -v
```
Expected: FAIL — `undefined: WriteCapture`.

- [x] ** Implement `capture.go`**

`Capture` performs the call, runs `Skeletonify` on the raw bytes, then attempts `json.Unmarshal(raw, into)` into the SDK type obtained from the symbol table, recording success or error. **A decode failure must not abort the capture** — it is precisely one of the findings the harness exists to report.

`WriteCapture` refuses any path that does not begin with `live/` and writes `live-manifest.json` alongside. Keys are sorted so a re-run produces a byte-identical file when nothing changed; a diff then means the server changed.

- [x] ** Run the test to verify it passes**

```
go test ./... -v
```
Expected: PASS.

- [x] ** Commit**

```
gofmt -l .
git add examples/live-probe
git commit -m "feat(live-probe): capture live bodies as value-free type skeletons"
```

- [x] ** Run the capture, then review before committing the evidence**

```
go run . -capture -out ..\..\conformance\testdata
git diff --stat conformance/testdata/live conformance/testdata/live-manifest.json
grep -rE '[A-Za-z0-9+/]{24,}' conformance/testdata/live/ conformance/testdata/live-manifest.json
```
Expected: the `grep` finds **no** value-shaped token. If it does, stop and fix the redactor before committing — this is the last automated check before real data reaches git.

---

## Task 5: phase 3 — the third divergence set

Reads the skeletons and runs the existing checks. This is where the "of N probed, Y disagree with the SDK" number comes from.

**Files:**
- Create: `conformance/live.go`
- Test: `conformance/live_test.go`

**Interfaces:**
- Consumes: `conformance.CompareBody(f Fixture, symbol string, t reflect.Type, body []byte) Outcome` and `conformance.Describe(t reflect.Type) WireShape` — **both already exist in `conformance/shapes.go` and are used unmodified.** `CompareBody` is the only entry point; a new comparison routine would be a second implementation of the same checks.
- Produces:
  - `type LiveDivergenceKind string` with `LiveSDKKind = "live-sdk"` and `LiveDocsKind = "live-docs"`.
  - `type LiveDivergence struct { Symbol, Fixture, Direction, Kind, Name, Detail string }`
  - `func CompareLive(f Fixture, symbol string, t reflect.Type, skeleton []byte) []LiveDivergence` — `live-docs` when the skeleton disagrees with the committed fixture, `live-sdk` when it disagrees with the SDK type.
  - `func LoadLive(dir string) (map[string]Skeleton, error)` where `Skeleton` is `{ Symbol, Fixture string; Body []byte }`.
  - `const LiveBaseline = "live-divergences.json"` — deliberately a different filename from `known-divergences.json`.

- [x] ** Write the failing test**

Two directions, and the one that matters is the SDK direction. A skeleton whose names the SDK cannot reach must produce exactly one `live-sdk` row and **no** `live-docs` row.

```go
func TestCompareLiveReportsTheSDKDirection(t *testing.T) {
	sk := []byte(`{"a":"string","b":"string"}`)
	type onlyA struct {
		A string `json:"a"`
	}
	// Build a Fixture from the committed documentation for a symbol, or construct a
	// minimal one, then assert exactly one live-sdk row naming "b".
	got := CompareLive(f, "test", mustTypeFor(t, reflect.TypeOf(onlyA{})), sk)
	// assert: one live-sdk divergence, name "b", and zero live-docs rows
}

// Review Focus case 4: the top-level kind disagreement the remaining 25 rows are
// made of. A live array against a documented object must be reported, not absorbed.
func TestCompareLiveReportsTheTopLevelKindDisagreement(t *testing.T) {
	// Documented as an object, live body is an array.
	// Assert one live-docs row whose detail names the container kind.
}
```

- [x] ** Run the test to verify it fails**

```
cd <repo root>
$env:GOTMPDIR="C:\Users\Tchan\AppData\Local\Temp\opencode\gotmp"
$env:GOCACHE="C:\Users\Tchan\AppData\Local\Temp\opencode\gocache"
go test ./conformance/ -run TestCompareLive -v
```
Expected: FAIL — `undefined: CompareLive`.

- [x] ** Implement `conformance/live.go`**

`CompareLive` calls `CompareBody(f, symbol, t, skeleton)` and maps each returned `Divergence` to a `LiveDivergence` with `Direction: "sdk"`. A second call with the documentation fixture's own body and the SDK type gives the `live-docs` direction, where a name in the documentation that the skeleton lacks means the server omits a documented name, and a name in the skeleton the documentation lacks means the server sends an undocumented one.

**Do not reuse `known-divergences.json`'s loader, its `DivergenceKind` values, or its exact-set gate.** A live finding is a claim about Webull's server, not about the SDK, and writing it into the SDK's baseline would make 25 documented rows plus N live rows read as one backlog of 25+N SDK defects. `live.go` gets its own file, its own kinds and its own gate.

- [x] ** Run the test to verify it passes**

```
go test ./conformance/ -v
```
Expected: PASS, and every pre-existing conformance test still passing.

- [x] ** Commit**

```
gofmt -l .
git add conformance/live.go conformance/live_test.go
git commit -m "feat(conformance): compare live skeletons against the SDK and the docs"
```

---

## Task 6: phase 4 — the report and the write-up

The deliverable that decides whether phases 2 and 3 were worth having.

**Files:**
- Create: `conformance/live_test.go` — the report test (extend the file from Task 5)
- Modify: `IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md`
- Modify: `CHANGELOG.md`, `docs/runs/index.md`

**Interfaces:**
- Consumes: `CompareLive`, `LoadLive` (Task 5).
- Produces: `conformance/live-divergences.json` and a section in `IMPLEMENTATION_STATUS.md` reporting, in one line, **of N probed, X disagree with their documentation and Y disagree with the SDK**.

- [x] ** Write the failing report test**

```go
func TestLiveReportCountsBothDirections(t *testing.T) {
	// Assert the report distinguishes the two directions and that the numbers it
	// prints are the counts it found, not a constant.
}
```

- [x] ** Run it to verify it fails**

```
go test ./conformance/ -run TestLiveReport -v
```
Expected: FAIL.

- [x] ** Implement the report and write it up**

Print and record both counts. Then write the section, and state whichever of these is true:

- **Y = 0** — no SDK disagreed with a live server. Say that plainly, and say that phases 2–3 were gated on this number and the gate closed the question.
- **Y > 0** — list each disagreeing endpoint with the name and the detail, and open a defect item for each.

Either way, record the **reachable count** and name the surfaces that stayed unreachable, so the 14 envelope methods and the ~159 endpoints needing other credentials are on the record as bounds rather than as surprises.

- [x] ** Verify the whole gate, and that the existing baseline is untouched**

```
gofmt -l .
go build ./...
go vet ./...
go test ./... -count=1
go test ./... -race -count=1
cd examples/live-probe && go test ./... && cd ../..
cd broker && go test ./... && cd ..
golangci-lint run ./...
python tools/conformance/gen_fixtures.py --check --self-test
python tools/citations/check.py
mkdocs build --strict
go test . -run TestStatus -count=1
git diff --stat conformance/known-divergences.json
```
Expected: all green, and the last command **empty** — the documentation baseline must be byte-identical.

- [x] ** Commit**

```
git add conformance/ IMPLEMENTATION_STATUS.md docs/implementation-status.md CHANGELOG.md docs/runs/index.md
git commit -m "docs: report the live evidence census and its two disagreement counts"
```

---

## Self-Review

**Spec coverage.** §5.1 skeleton → Task 1. §5.2 separate live set → Task 5, with the separation rationale restated in Task 5 step 3. §5.3 manifest → Task 4. §5.4 env-only credentials → Global Constraints, Task 3 step 3. §5.5 parameter resolution → Task 2, including the `account_id` branch and its ordering. §6 phases 0–4 → Tasks 3–6, with the gate on phase 4 → Task 3 step 6 and Task 6 step 3. §7 bounds → Task 6 step 3. §8 risks → Global Constraints and Review Focus. §9 verification → Task 4 step 6 (the grep), Task 6 step 4 (the empty baseline diff). §10 open questions → all resolved, none carried.

**Step scan.** No step decides nothing. The narrowest judgement left to the implementer is the correlation-key format in Task 2 step 3, which the test pins to "unique and non-empty" rather than a specific shape — correct, since the format is irrelevant and the uniqueness is the requirement.

**Type consistency.** `ParamSpec` is produced in Task 2 and consumed in Tasks 3 and 4 with the same name and field set. `Skeletonify`/`SkeletonKind` from Task 1 are consumed in Task 4. `Outcome`/`Endpoint`/`Census` from Task 3 are consumed in Task 4. `CompareLive`/`LoadLive` from Task 5 are consumed in Task 6. `CompareBody` and `Describe` are existing names used unchanged.

**Review Focus.** All five map to a test: number-vs-string to Task 1 step 1; the null field to Task 1 step 1; a rotated credential to Task 3 step 3; the top-level kind disagreement to Task 5 step 3; a 200 with an empty array to Task 1 step 1 (`{"empty array is not a wrong kind"}`).

**Proportion.** Five code sketches, all in test blocks, all for assertions the spec fixes. No function body is written for the implementer. The census loop, which is the largest piece, is described by its interface and its behavioural requirements rather than transcribed.

## Deferred, and why

- **The `docs/runs/index.md` row** is written in Task 6, not sooner, because its outcome column is the phase-4 number and writing it earlier would be a guess.
- **No `AGENTS.md` change.** The reachability matrix supersedes several recorded HTTP constraints, but rewriting them is a separate decision for the human partner, not a side effect of building a probe.

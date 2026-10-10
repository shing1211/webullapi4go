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

// Command live-probe asks every documented Webull endpoint whether the sandbox
// answers it, and records the answer as a status code.
//
// It exists because the conformance harness can only compare the SDK against
// Webull's published schemas. That proves the SDK matches the documentation; it
// has never called an endpoint, so it cannot say whether the SDK works. This
// command is the first thing in the repository that does.
//
// It writes census.json and prints a per-area table. What it records is whether
// an endpoint answered and how: never a response body, never a request value,
// and never a credential. A census row is evidence about the endpoint, and the
// three numbers it produces - endpoints the sandbox answered, endpoints the
// probe could not ask, and endpoints the probe refused to ask - are different
// numbers and are never summed.
//
// 34 of the 193 documented endpoints can place an order, cancel one, move money
// or modify an account. The probe does not call any of them unless BOTH of the
// gate's conditions hold: WEBULL_TRADE_MUTATE=1, the same opt-in the mutating
// sandbox tests use, AND a resolved base URL on a host the SDK derives for a
// sandbox. Two conditions because one is not enough: a shell set up to run the
// mutating tests already exports the first variable, and the base URL comes from
// the environment too, so the opt-in alone would let a run pointed at production
// place real orders. A run that genuinely needs a non-sandbox base URL sets
// WEBULL_LIVE_PROBE_MUTATE_NON_SANDBOX=1 to say so by name. A refused endpoint is
// reported as skipped rather than as any status: a run that did not ask must not
// be recorded as one that was told no.
//
// Usage:
//
//	live-probe -out census.json                          # walk the sandbox, mutating endpoints refused
//	live-probe -dry-run -out census-construct.json       # build every request, send nothing
//	live-probe -capture -out ..\..\conformance\testdata  # reduce the 200s and write live/
//
// A run needs WEBULL_APP_KEY and WEBULL_APP_SECRET, and it reads them from the
// environment through client.WithEnv. A -dry-run needs neither and contacts
// nothing. There is no flag that opens the mutating gate: the two environment
// variables are the only door, so the closed default cannot be widened by a
// mistyped argument.
//
// -capture is phase 2, and it runs on the reachability census this command
// produces: it re-calls every endpoint the census recorded as HTTP 200, reduces
// each response to a value-free type skeleton, and writes one file per endpoint
// plus live-manifest.json under the -out directory. It is refused outright if
// the mutation gate is open, because a capture re-calls the endpoints that
// answered 200 and an open gate would put the order-placing ones among them.
// The refusal is reached before the token exchange, before account discovery and
// before the walk, so a refused run sends no request of any kind and produces no
// census at all. runCapture holds the same check a second time, so a caller that
// reached the capture phase another way is refused there too; the gate in run is
// the one that has to come first, because a check the walk has already passed is
// not a refusal.
// With -capture, -out names a conformance/testdata directory rather than a
// census file, and the census is printed rather than written: the per-endpoint
// rows are a run-local artefact, and a run-local artefact has no business inside
// a tracked evidence tree.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// report is what census.json holds.
//
// It carries no timestamp, no hostname, no path outside the repository and no
// value the server sent, so two runs of the same corpus are byte-comparable and
// a reader can commit one without committing a reading. The endpoint's host is
// the single environment fact it records, and only the host: a base URL may
// carry userinfo, and the key and secret never reach this struct at all.
//
// Mutation is the one field that is always written, never omitted. It says
// whether the run was allowed to call the endpoints that can change state and,
// when it was not, WHICH of the gate's two conditions was unmet - so a census
// that called them can never be read as one that did not, and a census that was
// refused for pointing at a non-sandbox host can never be read as one that was
// merely un-authorised. Both of those are the ambiguities a shared opt-in
// variable creates, and the field is the reason they are recorded rather than
// derived from the counts after the fact.
type report struct {
	Mode        string                   `json:"mode"`
	Host        string                   `json:"host,omitempty"`
	Region      string                   `json:"region,omitempty"`
	Environment string                   `json:"environment,omitempty"`
	Account     string                   `json:"accountDiscovery"`
	Mutation    string                   `json:"mutationGate"`
	Endpoints   int                      `json:"endpointCount"`
	Summary     CensusSummary            `json:"summary"`
	Outcomes    []Outcome                `json:"outcomes"`
	Diagnostics map[string]DiagnosticSet `json:"diagnostics,omitempty"`
}

// DiagnosticSet is one class of probe limitation counted across the corpus. It
// exists so a reader can tell a real answer from a probe bug without reading
// every row: a row naming an SDK-signed header is a filter bug, and a row blocked
// for a parameter the page published an example for is a schema-reading bug, and
// both would otherwise sit in the same column as a genuine 404.
type DiagnosticSet struct {
	// Endpoints is how many endpoints the class touches.
	Endpoints int `json:"endpoints"`
	// Detail names the parameter, property or page each row is about, so a
	// class is a list a reader can act on rather than a count.
	Detail []string `json:"detail,omitempty"`
}

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "live-probe: "+err.Error())
		os.Exit(1)
	}
}

// run is main's body, separated so every exit path is a return rather than an
// os.Exit buried in a branch, and so the argument parsing is testable.
func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("live-probe", flag.ContinueOnError)
	out := flags.String("out", "census.json",
		"where to write the census as JSON; with -capture it is instead the "+
			"conformance/testdata directory the live tree is written into")
	base := flags.String("base", "",
		"override the base URL; empty uses the region and environment from the environment")
	account := flags.String("accounts", "",
		"account id to use instead of discovering one; empty runs phase 0")
	manifestPath := flags.String("manifest", "",
		"conformance/testdata/manifest.json; empty finds it by walking up from the working directory")
	cacheDir := flags.String("cache", "",
		"docgen cache directory; empty uses WEBULL_DOCGEN_CACHE, then the cache beside the manifest")
	dryRun := flags.Bool("dry-run", false,
		"build every documented request and send nothing; needs no credential")
	capture := flags.Bool("capture", false,
		"phase 2: re-call every endpoint the census recorded as HTTP 200, reduce each "+
			"response to a value-free type skeleton, and write the live tree and "+
			"live-manifest.json under the -out directory")
	assumeDiscovered := flags.Bool("assume-discovered", false,
		"dry-run only: resolve account_id and access_token from placeholder values, to "+
			"measure constructibility as it would be once phase 0 and the token exchange have both run")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *capture && *dryRun {
		return errors.New("-capture re-calls the endpoints that answered 200, so it cannot be " +
			"combined with -dry-run, which sends nothing")
	}
	if *assumeDiscovered && !*dryRun {
		return errors.New("-assume-discovered measures request construction and sends nothing, " +
			"so it is only meaningful with -dry-run; it is refused here rather than " +
			"silently substituting a placeholder for a real account or token")
	}

	manifest, err := resolveManifest(*manifestPath)
	if err != nil {
		return err
	}
	cache := resolveCache(*cacheDir, manifest)

	// The endpoint inventory is read before any credential is touched, so a
	// missing cache is a skip that says what to do about it rather than an
	// authentication error the operator has to decode.
	endpoints, err := LoadEndpoints(manifest, cache)
	if err != nil {
		if errors.Is(err, errCacheAbsent) {
			return skip(cache, manifest)
		}
		return err
	}

	if *dryRun {
		return runDryRun(endpoints, *out, *assumeDiscovered)
	}

	cl, err := newClient(*base)
	if err != nil {
		return fmt.Errorf("AUTH FAILURE: no endpoint was called, so nothing in this run is "+
			"evidence about an endpoint: %w", err)
	}
	defer func() { _ = cl.Close() }()

	// The host is resolved once, here, and the gate with it, because the gate is
	// judged on the host this run resolved rather than on one read again later.
	host := hostOf(cl)

	// The refusal comes before the token exchange, before phase 0 and before the
	// walk, and that position is the whole point. A capture re-calls every
	// endpoint a census recorded as a 200, so with the gate open that set includes
	// the endpoints that place an order, cancel one and move money - and the walk
	// that precedes runCapture would already have called them. The check in
	// runCapture is the second barrier, not the only one: by the time it is
	// reached the census has been taken, so a refusal there protects the capture
	// and nothing before it.
	//
	// Refusing here also settles what a refused run has and has not got. It has
	// no census: not one request was sent, so there is no per-endpoint evidence to
	// write and no artefact for an operator to mistake for one. That is the
	// intended behaviour rather than a side effect, and the refusal says so.
	if *capture {
		if gate := ResolveMutationGate(host); gate.Open() {
			return captureRefusal(gate)
		}
	}

	// Authentication is resolved before a single endpoint is called, and a
	// failure here is a top-level auth failure with a non-zero exit. It cannot be
	// mistaken for a conformance finding, which is the whole reason it is
	// separated: 193 endpoints answering 401 is a credential problem, and
	// recording it as 193 findings would bury it.
	//
	// The token is read back through Client.AccessToken rather than off the
	// returned *Token, so no local variable in this function ever holds a
	// credential. An empty token after a successful exchange is reported as the
	// auth failure it is, because it would otherwise surface as 26 endpoints
	// mysteriously blocked on access_token.
	if _, err := cl.EnsureToken(ctx); err != nil {
		return fmt.Errorf("AUTH FAILURE: could not obtain an access token, so no endpoint was "+
			"called and no row below is evidence about an endpoint: %w", err)
	}
	SetSessionToken(cl.AccessToken())
	if SessionToken() == "" {
		return errors.New("AUTH FAILURE: the token exchange reported success but yielded no " +
			"access token, so no endpoint was called")
	}

	// Phase 0. A failure here is not fatal: the endpoints that need an account
	// are recorded as blocked and the other 168 are still asked, which is a more
	// useful census than refusing to start.
	accountDiscovery := "not run: an account id was supplied"
	if strings.TrimSpace(*account) != "" {
		SetAccountID(strings.TrimSpace(*account))
	} else {
		discovered, discoverErr := DiscoverAccountID(ctx, cl)
		switch {
		case discoverErr != nil:
			SetAccountID("")
			accountDiscovery = discoveryFailure(discoverErr)
			fmt.Fprintln(os.Stderr, "live-probe: account discovery failed: "+accountDiscovery)
			fmt.Fprintln(os.Stderr, "live-probe: every account-scoped endpoint is reported blocked, "+
				"which is a statement about the probe and not about those endpoints")
		default:
			SetAccountID(discovered)
			accountDiscovery = "succeeded: one account id is threaded into every account-scoped request"
		}
	}

	outcomes, err := Census(ctx, cl, endpoints)
	if err != nil {
		return err
	}
	// The gate is not read here. It is resolved once, above, before anything was
	// called, and reading the environment again at this point would be a second
	// source of truth that could disagree with the rows below - and a check
	// reached after the walk is not a check the walk obeyed. Census resolves the
	// same gate from the same client host, so the mutationGate line above the
	// table and the rows beneath it are decided by one fact.
	summary := Summarise(outcomes)
	if *capture {
		return runCapture(ctx, cl, endpoints, outcomes, summary, *out, host,
			cl.Region().String(), cl.Environment().String())
	}
	built := report{
		Mode:        "live",
		Host:        host,
		Region:      cl.Region().String(),
		Environment: cl.Environment().String(),
		Account:     accountDiscovery,
		Mutation:    mutationNote("live", ResolveMutationGate(host)),
		Endpoints:   len(endpoints),
		Summary:     summary,
		Outcomes:    outcomes,
	}
	return finish(built, *out)
}

// runDryRun builds every documented request and sends nothing.
//
// It exists because the number that decides whether a live run is worth
// attempting is how many endpoints the probe can construct a request for at all,
// and that number is knowable without a credential. A census reporting blocked
// for nearly everything is not a finding about the endpoints; it is the probe
// failing to run, and this mode is where that shows up.
//
// assumeDiscovered supplies placeholder values for the two inputs phase 0 and
// the token exchange provide, so the count can be read as constructibility once
// both have run. The placeholders never leave the process: no request is sent,
// and neither value is written to the report.
//
// The mutating class is not gated here and that is the point, not an oversight.
// The gate is on sending; a dry run sends nothing, so gating it would subtract
// 34 endpoints from a constructibility measurement that is still true of them and
// would be the only place it can be made. Instead every mutating endpoint carries
// Mutating, and the summary prints the size of the class, so a reader knows what a
// live run would not call before choosing to call it.
func runDryRun(endpoints []Endpoint, out string, assumeDiscovered bool) error {
	accountNote := "not run: a dry run resolves no account, so every account-scoped request is blocked"
	if assumeDiscovered {
		SetAccountID("live-probe-assumed-account")
		SetSessionToken("live-probe-assumed-session-token")
		accountNote = "assumed: a placeholder account id and session token stand in for phase 0 " +
			"and the token exchange, so this count is constructibility once both have run"
		defer func() {
			SetAccountID("")
			SetSessionToken("")
		}()
	}
	outcomes := make([]Outcome, 0, len(endpoints))
	for _, ep := range endpoints {
		outcome := Outcome{
			Symbol:   ep.Symbol,
			Fixture:  ep.Fixture,
			Area:     areaOf(ep.Fixture),
			Method:   ep.Method,
			Path:     ep.Path,
			Mutating: IsMutating(ep),
			DryRun:   true,
		}
		request := prepare(ep)
		if request.Blocked != "" {
			outcome.Blocked = request.Blocked
		} else {
			outcome.QueryKeys = sortedKeys(request.Query)
			outcome.BodyBytes = len(request.Body)
			outcome.BodyKeys = request.BodyKeys
		}
		// Read from the document rather than from the built body, so the count is
		// the same whether or not the account id is known. See undeterminedArrays.
		outcome.SynthesisedEmptyArrays = undeterminedArrays(ep)
		outcome.DroppedSDKSignedHeaders = ep.DroppedSDKSignedHeaders
		outcomes = append(outcomes, outcome)
	}
	built := report{
		Mode:        "dry-run: no request was sent and no credential was read",
		Account:     accountNote,
		Mutation:    mutationNote("dry-run", MutationGate{}),
		Endpoints:   len(endpoints),
		Summary:     Summarise(outcomes),
		Outcomes:    outcomes,
		Diagnostics: diagnose(outcomes),
	}
	return finish(built, out)
}

// mutationNote states, in one line, whether this run was allowed to call the
// endpoints that can change state, which of the two conditions were satisfied,
// and which variable decides. It is the line an operator reads before believing
// any number below it, because the numbers mean different things depending on the
// answer: 158 reachable with the gate closed is a complete census of the read
// surface, and the same 158 with it open is a census that placed orders. (158 is
// the run recorded in conformance/testdata/live-manifest.json: the host answered
// 158 of 193, and 55 of those answered HTTP 200 - three different numbers that
// are never added.)
//
// The line names WHICH condition failed rather than only that the gate is closed,
// because the two failures call for different actions. A run that was never
// authorised needs one variable; a run that was authorised against a non-sandbox
// host needs its base URL changed, or the override set on purpose. Reporting both
// as "closed" would leave a run pointed at production looking merely un-approved.
func mutationNote(mode string, gate MutationGate) string {
	if isDryRun(mode) {
		return "not applicable: a dry run sends no request, so no endpoint that can mutate was " +
			"called; a live run needs " + MutateOptInEnv + "=" + MutateOptInValue + " and a sandbox " +
			"base URL to call any"
	}
	switch {
	case !gate.OptedIn:
		return "closed: " + MutateOptInEnv + "=" + MutateOptInValue + " is not set, so no endpoint in " +
			"the mutating class was called"
	case gate.SandboxHost:
		return "OPEN: " + MutateOptInEnv + "=" + MutateOptInValue + " was set and " + gate.Host +
			" is a sandbox host, so every endpoint in the mutating class was called and may have " +
			"changed sandbox account state"
	case gate.Overridden:
		return "OPEN BY OVERRIDE: " + NonSandboxOverrideEnv + "=" + NonSandboxOverrideValue +
			" was set and " + MutateOptInEnv + "=" + MutateOptInValue + " was set, so every endpoint " +
			"in the mutating class was called against " + gate.Host + ", which is NOT a sandbox host, " +
			"and may have changed real account state"
	default:
		return "closed: " + MutateOptInEnv + "=" + MutateOptInValue + " was set but " + gate.Host +
			" is not a sandbox host, so no endpoint in the mutating class was called; point the run " +
			"at a sandbox, or set " + NonSandboxOverrideEnv + "=" + NonSandboxOverrideValue +
			" to call them anyway"
	}
}

// isDryRun reports whether a mode string names a dry run. Both report builders
// and the table depend on it, so it is one predicate rather than three copies of
// a string prefix.
func isDryRun(mode string) bool {
	return strings.HasPrefix(mode, "dry-run")
}

// discoveryFailure renders why phase 0 could not resolve an account id, without
// quoting what the server said.
//
// The report's contract is that no value the server sent appears in it, and every
// other server string in the artefact is classified rather than printed: through
// classify for a request that produced a status, describeContentType for a media
// type, decodeFailure and reduceFailure for a body. Phase 0 was the one place
// that did not, and it is the one that could have written a whole HTML error page
// into census.json, because errs.FromHTTPStatus folds the response body into its
// message and falls back to the raw trimmed body when the body is not a
// recognised Webull error envelope. The file is gitignored and written 0600, so
// the blast radius was local rather than committed; the discipline is the same
// either way, and this artefact is the evidence the write-up quotes.
//
// So the line names the status and the class and stops. A status and an SDK error
// code are facts about the request, and they are what a reader acts on - an
// operator who is told phase 0 was refused with FORBIDDEN looks at the
// entitlement, and one handed a sentence of server prose looks at nothing.
func discoveryFailure(err error) string {
	if errors.Is(err, errNoAccountReported) {
		return "failed: the account-list request was answered with an empty list, so no account id " +
			"could be threaded into an account-scoped request and every account-scoped endpoint " +
			"is reported blocked"
	}
	var typed *errs.Error
	if !errors.As(err, &typed) {
		return "failed: the account-list request produced no HTTP answer, which the SDK " +
			"classifies as " + classify(err) + "; the error's own message is not recorded"
	}
	if typed.Status == 0 {
		return "failed: the account-list request produced no HTTP answer, which the SDK " +
			"classifies as " + string(typed.Code) + "; the error's own message is not recorded"
	}
	return "failed: the account-list request was answered HTTP " + strconv.Itoa(typed.Status) +
		", which the SDK classifies as " + string(typed.Code) + "; the server's own message is " +
		"not recorded"
}

// finish renders the table and writes the JSON, and is where the four numbers -
// answered, unbuildable, refused and unanswered - are printed far enough apart
// that a reader cannot add any two of them.
func finish(built report, out string) error {
	built.Diagnostics = diagnose(built.Outcomes)
	printSummary(built)
	if err := writeReport(built, out); err != nil {
		return err
	}
	fmt.Printf("\nwrote %s\n", out)
	return nil
}

// printSummary prints the per-area table and the diagnostics.
//
// The columns that matter are Reachable, Blocked and Skipped, and the gate line
// above them is why: a reachable endpoint may have answered 404, a blocked one
// was never asked, and a skipped one was deliberately not asked, and the third
// says nothing whatever about the server. The Reachable column is the summary's
// own Reachable and not Total minus Blocked, because that subtraction would count
// every refused endpoint as an answer - which is the one number a safety gate
// must not produce.
func printSummary(built report) {
	s := built.Summary
	// A dry run sent nothing, so "reachable" would be zero by construction and
	// the number a reader compares against an expectation is the complement of
	// blocked and skipped. The two modes are named for what they measure so they
	// cannot be read as the same table.
	primary, primaryLabel := s.Reachable, "reachable"
	secondLabel := "no response"
	if isDryRun(built.Mode) {
		primaryLabel, secondLabel = "constructible", "not called"
		primary = s.Total - s.Blocked - s.Skipped
	}
	fmt.Printf("%s\n", built.Mode)
	if built.Host != "" {
		fmt.Printf("host %s  region %s  environment %s\n", built.Host, built.Region, built.Environment)
	}
	fmt.Printf("account discovery: %s\n", built.Account)
	fmt.Printf("mutating gate: %s\n", built.Mutation)
	fmt.Printf("endpoints %d   %s %d   blocked %d   skipped %d   %s %d   mutating class %d\n\n",
		s.Total, primaryLabel, primary, s.Blocked, s.Skipped, secondLabel, s.Unanswered,
		s.Mutating)

	areas := make([]string, 0, len(s.Area))
	for name := range s.Area {
		areas = append(areas, name)
	}
	sort.Strings(areas)
	fmt.Printf("%-24s %5s %12s %7s %8s %9s %12s  %s\n", "area", "total", primaryLabel, "blocked",
		"skipped", "mutating", secondLabel, "statuses")
	for _, name := range areas {
		a := s.Area[name]
		shown := a.Reachable
		if isDryRun(built.Mode) {
			shown = a.Total - a.Blocked - a.Skipped
		}
		fmt.Printf("%-24s %5d %12d %7d %8d %9d %12d  %s\n", name, a.Total, shown, a.Blocked,
			a.Skipped, a.Mutating, a.Unanswered, renderCounts(a.Statuses))
	}
	fmt.Printf("%-24s %5d %12d %7d %8d %9d %12d  %s\n", "TOTAL", s.Total, primary, s.Blocked,
		s.Skipped, s.Mutating, s.Unanswered, renderCounts(s.Statuses))

	if len(s.BlockedReasons) > 0 {
		fmt.Println("\nblocked, by reason:")
		for _, reason := range sortedKeysOfCounts(s.BlockedReasons) {
			fmt.Printf("  %4d  %s\n", s.BlockedReasons[reason], reason)
		}
	}
	if s.Skipped > 0 {
		fmt.Printf("\nskipped: %d endpoint(s) were not called. A skipped row is a statement about\n", s.Skipped)
		fmt.Println("this probe, not about the endpoint: it is not a 404, a 403 or any other status,")
		fmt.Println("because no request was sent. Calling them needs two things, both deliberate:")
		fmt.Printf("%s=%s, and a base URL on a sandbox host. To call them against a\n",
			MutateOptInEnv, MutateOptInValue)
		fmt.Printf("non-sandbox host as well, set %s=%s.\n", NonSandboxOverrideEnv, NonSandboxOverrideValue)
	}
	if len(s.Errs) > 0 {
		fmt.Println("\nrequests that produced no status, by SDK error code:")
		for _, code := range sortedKeysOfCounts(s.Errs) {
			fmt.Printf("  %4d  %s\n", s.Errs[code], code)
		}
	}

	// One line per endpoint, so a reader can see which row says what rather than
	// inferring it from a total. The status column is the answer; the reason
	// column is the reason the probe has it or does not.
	fmt.Println("\nper-endpoint outcomes:")
	for _, o := range built.Outcomes {
		fmt.Printf("  %-6s %-58s %-18s %s\n", o.Method, o.Path, statusText(o), rowReason(o))
	}

	// A reachable endpoint is not a working endpoint, and the reason is the
	// status. Saying it here means nobody has to infer it from the table.
	fmt.Println("\nread this next: reachable counts every endpoint that answered, whatever it")
	fmt.Println("answered. A 404 from broker-fd-us, a 403 from display-solution and a 401 from")
	fmt.Println("broker-hk are live answers about an endpoint's reachability, and they are not")
	fmt.Println("evidence that the SDK's decoding is right. blocked means the probe could not")
	fmt.Println("build the request and skipped means it declined to send one; neither is an")
	fmt.Println("answer from any server. Request-constructibility is the separate number in a")
	fmt.Println("dry run, and the two diverge by area: constructible is everything in market")
	fmt.Println("data and nothing in display-solution.")

	if len(built.Diagnostics) > 0 {
		fmt.Println("\nprobe diagnostics - a non-zero count here is a statement about this probe:")
		for _, name := range diagnosticNames(built.Diagnostics) {
			d := built.Diagnostics[name]
			fmt.Printf("  %-34s %4d endpoint(s)\n", name, d.Endpoints)
			for _, detail := range d.Detail {
				fmt.Printf("      %s\n", detail)
			}
		}
	}
}

// diagnose classifies the rows that describe a limitation of the probe rather
// than an answer from an endpoint.
//
// The three classes that matter are the three a reader would otherwise have to
// catch by eye. A row blocked on a parameter named like an SDK-signed header is a
// filter bug; a row blocked on a parameter the page published an example for is a
// schema-reading bug; a row in the mutating class is a gate the operator did not
// open. Each is counted here so a real 404 is never read as one of them.
//
// The mutating class is keyed off Outcome.Mutating and not off Skipped, so it
// reports the class in a dry run too. A dry run sends nothing, so every row in it
// is uncalled; the class is still the number a reader needs before a live run, and
// the gate line above the table says which of them a live run would skip.
func diagnose(outcomes []Outcome) map[string]DiagnosticSet {
	header := DiagnosticSet{}
	exampleMissed := DiagnosticSet{}
	bodyString := DiagnosticSet{}
	noBody := DiagnosticSet{}
	account := DiagnosticSet{}
	sessionToken := DiagnosticSet{}
	emptyArray := DiagnosticSet{}
	mutating := DiagnosticSet{}

	for _, o := range outcomes {
		if len(o.SynthesisedEmptyArrays) > 0 {
			emptyArray.Detail = append(emptyArray.Detail, o.Fixture+": "+
				strings.Join(o.SynthesisedEmptyArrays, ", "))
		}
		if o.Mutating {
			mutating.Detail = append(mutating.Detail, o.Fixture+": "+o.Method+" in "+o.Area)
		}
		if o.Blocked == "" {
			continue
		}
		reason := o.Blocked
		// The order is the contract. A signed header can only come from the
		// filter failing, so it is checked first: it is the one class that is
		// always a bug. The session token is checked next because its reason
		// text is Param's own, which names a parameter and a type, and it would
		// otherwise fall through to the example class and make that class - the
		// one that catches a schema-reading bug - look broken.
		switch {
		case mentionsSignedHeader(reason):
			header.Detail = append(header.Detail, o.Fixture+": "+reason)
		case strings.Contains(reason, blockedAccountReason):
			account.Detail = append(account.Detail, o.Fixture)
		case strings.Contains(reason, "no JSON request body"):
			noBody.Detail = append(noBody.Detail, o.Fixture+": "+reason)
		case strings.Contains(reason, "body property"):
			bodyString.Detail = append(bodyString.Detail, o.Fixture+": "+reason)
		case strings.Contains(reason, sessionTokenParam):
			sessionToken.Detail = append(sessionToken.Detail, o.Fixture+": "+reason)
		default:
			exampleMissed.Detail = append(exampleMissed.Detail, o.Fixture+": "+reason)
		}
	}

	set := map[string]DiagnosticSet{}
	for name, d := range map[string]DiagnosticSet{
		"sdk-signed header not filtered":                  header,
		"example published but not read":                  exampleMissed,
		"undocumented required body string":               bodyString,
		"request body is not JSON":                        noBody,
		"account id unavailable":                          account,
		"session token unavailable":                       sessionToken,
		"body rests on a synthesised empty array":         emptyArray,
		"endpoint can mutate, gated by " + MutateOptInEnv: mutating,
	} {
		d.Endpoints = len(d.Detail)
		sort.Strings(d.Detail)
		if d.Endpoints > 0 {
			set[name] = d
		}
	}
	return set
}

// mentionsSignedHeader reports whether a blocked reason names one of the nine
// headers the SDK signs. Those are the only names a filter bug can produce, so
// this is the check that separates a probe bug from a finding.
func mentionsSignedHeader(reason string) bool {
	for name := range sdkSignedHeaders {
		if strings.Contains(reason, name) {
			return true
		}
	}
	return false
}

// diagnosticNames returns the diagnostic class names in a stable order.
func diagnosticNames(set map[string]DiagnosticSet) []string {
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// renderCounts renders a count map in a stable order.
func renderCounts(counts map[string]int) string {
	if len(counts) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(counts))
	for _, key := range sortedKeysOfCounts(counts) {
		parts = append(parts, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return strings.Join(parts, " ")
}

// sortedKeysOfCounts returns the keys of a count map in ascending order.
func sortedKeysOfCounts(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// newClient builds the SDK client from the environment, with an optional base
// URL override. The credentials come from client.WithEnv and from nowhere else:
// no flag takes a key, no flag takes a secret, and nothing this program prints
// or writes carries one.
func newClient(base string) (*client.Client, error) {
	opts := []client.Option{client.WithEnv()}
	if trimmed := strings.TrimSpace(base); trimmed != "" {
		opts = append(opts, client.WithBaseURL(trimmed))
	}
	cl, err := client.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("building the client (are WEBULL_APP_KEY and WEBULL_APP_SECRET set?): %w", err)
	}
	return cl, nil
}

// hostOf returns the host of the client's base URL, and nothing else. The path,
// the query and any userinfo are dropped: a base URL is operator-supplied and a
// URL may carry a credential, while the host is the one fact that says which
// deployment answered.
func hostOf(cl *client.Client) string {
	endpoints := cl.Endpoints()
	for _, raw := range []string{endpoints.HTTP, endpoints.BrokerHTTP} {
		if raw == "" {
			continue
		}
		if parsed, err := url.Parse(raw); err == nil && parsed.Host != "" {
			return parsed.Host
		}
	}
	return ""
}

// writeReport writes the census as indented JSON.
func writeReport(built report, out string) error {
	encoded, err := json.MarshalIndent(built, "", " ")
	if err != nil {
		return fmt.Errorf("encoding the census: %w", err)
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(out, encoded, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", out, err)
	}
	return nil
}

// skip reports an absent cache and asks for it to be populated, then returns nil.
//
// The exit status is 0 on purpose. The cache is gitignored, so a fresh checkout
// cannot judge the corpus in either direction, and a non-zero exit would report a
// defect in the SDK that is not there. The conformance-fixtures target skips for
// the same reason, and this message names the directory and the two ways to fill
// it so the skip is actionable rather than a shrug.
func skip(cache, manifest string) error {
	fmt.Printf("SKIP: no docgen cache at %s\n", cache)
	fmt.Println("      The cache is gitignored, so a checkout without it cannot read a single")
	fmt.Println("      documented request schema and the census has nothing to ask.")
	fmt.Println("      Populate it with tools/webull-docgen (needs network) or point")
	fmt.Println("      WEBULL_DOCGEN_CACHE, or -cache, at a populated directory.")
	fmt.Printf("      The endpoint inventory is %s; it names the pages, not their inputs.\n", manifest)
	return nil
}

// resolveManifest finds conformance/testdata/manifest.json, walking up from the
// working directory when -manifest is not given. The census has to be runnable
// from the module directory and from the repository root, and a relative default
// would work in one and not the other.
func resolveManifest(flagValue string) (string, error) {
	if trimmed := strings.TrimSpace(flagValue); trimmed != "" {
		if _, err := os.Stat(trimmed); err != nil {
			return "", fmt.Errorf("reading the manifest at %s: %w", trimmed, err)
		}
		return trimmed, nil
	}
	if fromEnv := strings.TrimSpace(os.Getenv("WEBULL_CONFORMANCE_MANIFEST")); fromEnv != "" {
		if _, err := os.Stat(fromEnv); err != nil {
			return "", fmt.Errorf("reading the manifest at %s: %w", fromEnv, err)
		}
		return fromEnv, nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolving the working directory: %w", err)
	}
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "conformance", "testdata", "manifest.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("live-probe: could not find conformance/testdata/manifest.json by " +
		"walking up from the working directory; pass -manifest or set " +
		"WEBULL_CONFORMANCE_MANIFEST")
}

// resolveCache picks the docgen cache directory, in the order the repo's own
// tools use: an explicit -cache, then WEBULL_DOCGEN_CACHE, then the cache beside
// the manifest. The last is what makes the flag unnecessary in the normal case,
// and it is derived from the manifest rather than from a relative path so it
// resolves identically from any working directory.
func resolveCache(flagValue, manifest string) string {
	if trimmed := strings.TrimSpace(flagValue); trimmed != "" {
		return trimmed
	}
	if fromEnv := strings.TrimSpace(os.Getenv("WEBULL_DOCGEN_CACHE")); fromEnv != "" {
		return fromEnv
	}
	// <repo>/conformance/testdata/manifest.json -> <repo>/tools/webull-docgen/.cache
	root := filepath.Dir(filepath.Dir(filepath.Dir(manifest)))
	return filepath.Join(root, "tools", "webull-docgen", ".cache")
}

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
// two numbers it produces - endpoints the sandbox answered and endpoints the
// probe was able to ask - are different numbers and are never summed.
//
// Usage:
//
//	live-probe -out census.json                     # walk the sandbox
//	live-probe -dry-run -out census-construct.json # build every request, send nothing
//
// A run needs WEBULL_APP_KEY and WEBULL_APP_SECRET, and it reads them from the
// environment through client.WithEnv. A -dry-run needs neither and contacts
// nothing.
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
	"strings"

	"github.com/shing1211/webullapi4go/client"
)

// report is what census.json holds.
//
// It carries no timestamp, no hostname, no path outside the repository and no
// value the server sent, so two runs of the same corpus are byte-comparable and
// a reader can commit one without committing a reading. The endpoint's host is
// the single environment fact it records, and only the host: a base URL may
// carry userinfo, and the key and secret never reach this struct at all.
type report struct {
	Mode        string                   `json:"mode"`
	Host        string                   `json:"host,omitempty"`
	Region      string                   `json:"region,omitempty"`
	Environment string                   `json:"environment,omitempty"`
	Account     string                   `json:"accountDiscovery"`
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
		"where to write the census as JSON")
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
	if err := flags.Parse(args); err != nil {
		return err
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
		return runDryRun(endpoints, *out, manifest, cache)
	}

	cl, err := newClient(*base)
	if err != nil {
		return fmt.Errorf("AUTH FAILURE: no endpoint was called, so nothing in this run is "+
			"evidence about an endpoint: %w", err)
	}
	defer func() { _ = cl.Close() }()

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
			accountDiscovery = "failed: " + discoverErr.Error()
			fmt.Fprintln(os.Stderr, "live-probe: account discovery failed: "+discoverErr.Error())
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
	built := report{
		Mode:        "live",
		Host:        hostOf(cl),
		Region:      cl.Region().String(),
		Environment: cl.Environment().String(),
		Account:     accountDiscovery,
		Endpoints:   len(endpoints),
		Summary:     Summarise(outcomes),
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
func runDryRun(endpoints []Endpoint, out, manifest, cache string) error {
	outcomes := make([]Outcome, 0, len(endpoints))
	for _, ep := range endpoints {
		outcome := Outcome{
			Symbol:  ep.Symbol,
			Fixture: ep.Fixture,
			Area:    areaOf(ep.Fixture),
			Method:  ep.Method,
			Path:    ep.Path,
			DryRun:  true,
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
		Account:     "not run: a dry run resolves no account",
		Endpoints:   len(endpoints),
		Summary:     Summarise(outcomes),
		Outcomes:    outcomes,
		Diagnostics: diagnose(outcomes),
	}
	return finish(built, out)
}

// finish renders the table and writes the JSON, and is where the two numbers
// are printed far enough apart that a reader cannot add them.
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
// The column that matters is Reachable, and the note under it is the reason the
// table exists: a reachable endpoint may have answered 404, and constructibility
// is a different number again, and by area the two diverge sharply.
func printSummary(built report) {
	s := built.Summary
	// A dry run sent nothing, so "reachable" would be zero by construction and
	// the number a reader compares against an expectation is the complement of
	// blocked. The two columns are named for what they measure so the two modes
	// cannot be read as the same table.
	buildable := "reachable"
	second := "no response"
	if strings.HasPrefix(built.Mode, "dry-run") {
		buildable, second = "constructible", "not called"
	}
	fmt.Printf("%s\n", built.Mode)
	if built.Host != "" {
		fmt.Printf("host %s  region %s  environment %s\n", built.Host, built.Region, built.Environment)
	}
	fmt.Printf("account discovery: %s\n", built.Account)
	fmt.Printf("endpoints %d   %s %d   blocked %d   %s %d\n\n",
		s.Total, buildable, s.Total-s.Blocked, s.Blocked, second, s.Unanswered)

	areas := make([]string, 0, len(s.Area))
	for name := range s.Area {
		areas = append(areas, name)
	}
	sort.Strings(areas)
	fmt.Printf("%-24s %5s %9s %7s %12s  %s\n", "area", "total", buildable, "blocked",
		second, "statuses")
	for _, name := range areas {
		a := s.Area[name]
		shown := a.Unanswered
		if strings.HasPrefix(built.Mode, "dry-run") {
			shown = a.Total - a.Blocked
		}
		fmt.Printf("%-24s %5d %9d %7d %12d  %s\n", name, a.Total, shown, a.Blocked,
			a.Unanswered, renderCounts(a.Statuses))
	}
	fmt.Printf("%-24s %5d %9d %7d %12d  %s\n", "TOTAL", s.Total, s.Total-s.Blocked, s.Blocked,
		s.Unanswered, renderCounts(s.Statuses))

	if len(s.BlockedReasons) > 0 {
		fmt.Println("\nblocked, by reason:")
		for _, reason := range sortedKeysOfCounts(s.BlockedReasons) {
			fmt.Printf("  %4d  %s\n", s.BlockedReasons[reason], reason)
		}
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
		fmt.Printf("  %-6s %-58s %-12s %s\n", o.Method, o.Path, statusText(o), blockedOrErr(o))
	}

	// A reachable endpoint is not a working endpoint, and the reason is the
	// status. Saying it here means nobody has to infer it from the table.
	fmt.Println("\nread this next: reachable counts every endpoint that answered, whatever it")
	fmt.Println("answered. A 404 from broker-fd-us, a 403 from display-solution and a 401 from")
	fmt.Println("broker-hk are live answers about an endpoint's reachability, and they are not")
	fmt.Println("evidence that the SDK's decoding is right. Request-constructibility is the")
	fmt.Println("separate number in a dry run, and the two diverge by area: constructible is")
	fmt.Println("everything in market data and nothing in display-solution.")

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
// The two classes that matter are the two a reader would otherwise have to catch
// by eye. A row blocked on a parameter named like an SDK-signed header is a
// filter bug, and a row blocked on a parameter the page published an example for
// is a schema-reading bug; both are counted here so a real 404 is never read as
// one of them.
func diagnose(outcomes []Outcome) map[string]DiagnosticSet {
	header := DiagnosticSet{}
	exampleMissed := DiagnosticSet{}
	bodyString := DiagnosticSet{}
	noBody := DiagnosticSet{}
	account := DiagnosticSet{}
	sessionToken := DiagnosticSet{}
	emptyArray := DiagnosticSet{}

	for _, o := range outcomes {
		if len(o.SynthesisedEmptyArrays) > 0 {
			emptyArray.Detail = append(emptyArray.Detail, o.Fixture+": "+
				strings.Join(o.SynthesisedEmptyArrays, ", "))
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
		"sdk-signed header not filtered":          header,
		"example published but not read":          exampleMissed,
		"undocumented required body string":       bodyString,
		"request body is not JSON":                noBody,
		"account id unavailable":                  account,
		"session token unavailable":               sessionToken,
		"body rests on a synthesised empty array": emptyArray,
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

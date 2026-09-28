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

package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// This file adds the one contract comparison the harness was missing: the HTTP
// verb.
//
// Every other check here compares a documented *response* against what the SDK
// decodes. None of them looks at the *request* the SDK sends, so a method that
// sends the right path with the wrong verb reconciles as a clean match.
// broker.UpdateVirtualAccount did exactly that for as long as the harness has
// existed: it issued PUT where the page documents POST, put account_id in the
// query where the page requires it in the body, and sent an account_name field
// the page does not declare at all. The path matched, so every row was green.
//
// The verb is read from the SDK source rather than from a hand-maintained table,
// for the same reason envelopes.go parses unexported types instead of asking a
// caller to register them: a table would be a second place to forget to update,
// and a stale table would report agreement about code that no longer exists. The
// generator already records documented.method for all 193 fixtures, so the
// documented side needed nothing.

// verbHelpers maps the unexported per-verb request helpers the service packages
// define. Every package has get and post; the ones added later have put, patch
// and delete.
var verbHelpers = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"put":    "PUT",
	"patch":  "PATCH",
	"delete": "DELETE",
}

// exportedVerbHelpers are the same verbs on a sub-service, reached as
// c.DisplayService().Get(...). The data package routes its Display Solution
// endpoints through display.Service this way rather than through its own helper.
var exportedVerbHelpers = map[string]string{
	"Get":     "GET",
	"Post":    "POST",
	"Put":     "PUT",
	"Patch":   "PATCH",
	"Delete":  "DELETE",
	"Head":    "HEAD",
	"Options": "OPTIONS",
}

// httpMethodConstants maps net/http's method constants. A method reaches its verb
// through one as a bare argument - buildSignedRequest(ctx, http.MethodPost, path, ...)
// - which is why these are collected from any argument position rather than from
// a call's function.
var httpMethodConstants = map[string]string{
	"MethodGet":     "GET",
	"MethodPost":    "POST",
	"MethodPut":     "PUT",
	"MethodPatch":   "PATCH",
	"MethodDelete":  "DELETE",
	"MethodHead":    "HEAD",
	"MethodOptions": "OPTIONS",
}

// doHelpers are the names a shared request helper is reached under. A verb written
// as a string literal - connect.CreateToken's c.core.Do(ctx, "POST", path, ...)
// - is only read at one of these calls, so that a method which merely compares a
// string against "GET" is not read as sending one.
var doHelpers = map[string]bool{
	"do":       true,
	"Do":       true,
	"DoBroker": true,
	"DoStream": true,
}

// maxDelegationDepth bounds the walk through same-package helpers.
//
// Delegation is what makes the extractor usable rather than merely working on the
// easy methods: trade.GetOpenOrders calls GetOpenOrdersPage and returns a field of
// its result, and display.Service.EnsureToken calls fetchToken, which is the
// function that actually names a verb. Both are common shapes here, so refusing to
// follow them would have left a third of the surface unreadable. The bound and the
// visited set exist because a cycle would otherwise hang the test run.
const maxDelegationDepth = 6

// SDKVerb is the HTTP verb an SDK method sends, as read from its source.
type SDKVerb struct {
	// Verb is the single verb the method sends, or "" when there is not exactly
	// one.
	Verb string
	// Verbs is every verb found, in order, including the single one. A method
	// that sends two different verbs is not unreadable and is not a disagreement
	// either: it has no single counterpart verb to compare against a page, so it
	// is reported as its own state rather than forced into either answer.
	Verbs []string
	// Reason explains a blank Verb when Verbs is empty.
	Reason string
}

// Known reports whether any verb was read, which is the question the coverage
// test asks.
func (v SDKVerb) Known() bool { return len(v.Verbs) > 0 }

// Single reports whether exactly one verb was read, which is the question the
// comparison asks.
func (v SDKVerb) Single() bool { return len(v.Verbs) == 1 }

var (
	verbCacheMu sync.Mutex
	verbCache   = map[string]SDKVerb{}
	pkgIndexMu  sync.Mutex
	pkgIndex    = map[string]map[string]*ast.FuncDecl{}
	pkgIndexErr = map[string]string{}
)

// SDKVerbOf reports the HTTP verb the named symbol's method sends.
//
// The symbol is the manifest's spelling. It is usually "package.Method" and is
// sometimes "package.Type.Method" - display.Service.EnsureToken is the latter - so
// the method is the last dotted segment.
func SDKVerbOf(symbol string) SDKVerb {
	verbCacheMu.Lock()
	if v, ok := verbCache[symbol]; ok {
		verbCacheMu.Unlock()
		return v
	}
	verbCacheMu.Unlock()

	v := readSDKVerb(symbol)

	verbCacheMu.Lock()
	verbCache[symbol] = v
	verbCacheMu.Unlock()
	return v
}

func readSDKVerb(symbol string) SDKVerb {
	pkg, method, ok := splitSymbol(symbol)
	if !ok {
		return SDKVerb{Reason: fmt.Sprintf("the manifest symbol %q is not in package.Method form", symbol)}
	}
	funcs, err := packageFuncs(pkg)
	if err != "" {
		return SDKVerb{Reason: fmt.Sprintf("cannot read package %s: %s", pkg, err)}
	}
	if _, ok := funcs[method]; !ok {
		return SDKVerb{Reason: fmt.Sprintf("no function or method named %s is declared in package %s", method, pkg)}
	}

	// Breadth-first through same-package delegation, collecting every verb any
	// reachable body can be shown to send. Agreement everywhere is agreement;
	// disagreement is reported rather than resolved by picking one.
	verbs := map[string]bool{}
	visited := map[string]bool{}
	queue := []string{method}
	for depth := 0; len(queue) > 0 && depth <= maxDelegationDepth; depth++ {
		var next []string
		for _, name := range queue {
			if visited[name] {
				continue
			}
			visited[name] = true
			fn, ok := funcs[name]
			if !ok || fn.Body == nil {
				continue
			}
			for _, v := range verbsIn(fn, receiverName(fn)) {
				verbs[v] = true
			}
			next = append(next, calleesIn(fn)...)
		}
		queue = next
	}

	uniq := sortedVerbs(verbs)
	switch len(uniq) {
	case 0:
		return SDKVerb{Reason: fmt.Sprintf("the body of %s sends no recognisable HTTP verb, not even "+
			"through %d levels of same-package delegation, so it is reached by some route this "+
			"extractor does not model; read the method rather than trusting a blank row",
			symbol, maxDelegationDepth)}
	case 1:
		return SDKVerb{Verb: uniq[0], Verbs: uniq}
	default:
		return SDKVerb{Verbs: uniq, Reason: fmt.Sprintf("the body of %s sends %d verbs (%s), so it "+
			"has no single counterpart verb to compare against a page", symbol, len(uniq),
			strings.Join(uniq, ", "))}
	}
}

// splitSymbol separates a manifest symbol into its package and the method to look
// for. The method is the last dotted segment, so display.Service.EnsureToken
// resolves to package display and method EnsureToken.
func splitSymbol(symbol string) (pkg, method string, ok bool) {
	i := strings.LastIndex(symbol, ".")
	if i <= 0 || i == len(symbol)-1 {
		return "", "", false
	}
	pkg, method = symbol[:i], symbol[i+1:]
	if strings.Contains(pkg, ".") {
		// A three-part symbol's package is its first segment.
		pkg = strings.SplitN(pkg, ".", 2)[0]
	}
	return pkg, method, true
}

// packageFuncs parses a package's non-test sources once and indexes every
// top-level function and method by name.
func packageFuncs(pkg string) (map[string]*ast.FuncDecl, string) {
	pkgIndexMu.Lock()
	if f, ok := pkgIndex[pkg]; ok {
		pkgIndexMu.Unlock()
		return f, ""
	}
	if msg, ok := pkgIndexErr[pkg]; ok {
		pkgIndexMu.Unlock()
		return nil, msg
	}
	pkgIndexMu.Unlock()

	dir, err := packageDir(pkg)
	var funcs map[string]*ast.FuncDecl
	var errMsg string
	if err != nil {
		errMsg = fmt.Sprintf("cannot locate package %s: %v", pkg, err)
	} else {
		fset := token.NewFileSet()
		pkgs, perr := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
			return !strings.HasSuffix(fi.Name(), "_test.go")
		}, 0)
		if perr != nil {
			errMsg = fmt.Sprintf("cannot parse package %s: %v", pkg, perr)
		} else {
			funcs = map[string]*ast.FuncDecl{}
			for _, p := range pkgs {
				for _, file := range p.Files {
					for _, d := range file.Decls {
						fn, ok := d.(*ast.FuncDecl)
						if !ok {
							continue
						}
						// A name declared more than once at package scope cannot be
						// indexed unambiguously, so it is left out and the symbol
						// reads as undeclared rather than picking one.
						if _, dup := funcs[fn.Name.Name]; !dup {
							funcs[fn.Name.Name] = fn
						} else {
							delete(funcs, fn.Name.Name)
						}
					}
				}
			}
		}
	}

	pkgIndexMu.Lock()
	if errMsg != "" {
		pkgIndexErr[pkg] = errMsg
	} else {
		pkgIndex[pkg] = funcs
	}
	pkgIndexMu.Unlock()
	return funcs, errMsg
}

// receiverName is the method's receiver variable, or "" for a package-level
// function. connect.CreateToken is a method and display.Service.EnsureToken is a
// method on Service, so neither shape is skipped, but a package-level function is
// accepted too rather than reported unreadable for a reason unrelated to its verb.
func receiverName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 || len(fn.Recv.List[0].Names) == 0 {
		return ""
	}
	return fn.Recv.List[0].Names[0].Name
}

// verbsIn collects every verb the body can be shown to send.
//
// recv is the function's own receiver variable, and the distinction matters: an
// unexported get/post on the receiver is this SDK's transport, while a Get on some
// other identifier is usually url.Values.Get reading a query parameter. Counting
// that second one as a GET verb would make every method that reads a query
// parameter look like a GET, including the POST ones.
func verbsIn(fn *ast.FuncDecl, recv string) []string {
	var out []string
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			// http.MethodPost and friends, wherever they appear. A method that names
			// one is sending that verb; nothing else in this SDK refers to them.
			if verb, ok := httpMethodConstants[node.Sel.Name]; ok {
				if pkgIdent, isIdent := node.X.(*ast.Ident); isIdent && pkgIdent.Name == "http" {
					out = append(out, verb)
				}
			}
		case *ast.CallExpr:
			sel, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				// A transport named as a plain function, c(path, ...).
				if ident, isIdent := node.Fun.(*ast.Ident); isIdent {
					out = append(out, verbLiterals(node, ident.Name)...)
				}
				return true
			}
			name := sel.Sel.Name
			// The transport reached through this SDK's own unexported helper.
			if verb, ok := verbHelpers[name]; ok {
				if x, isIdent := sel.X.(*ast.Ident); isIdent && x.Name == recv && recv != "" {
					out = append(out, verb)
				}
				return true
			}
			// The transport reached through a sub-service, as
			// c.DisplayService().Get(...). The receiver is a call or a field
			// selector, never a bare variable, which is what keeps query.Get out.
			if verb, ok := exportedVerbHelpers[name]; ok {
				switch x := sel.X.(type) {
				case *ast.CallExpr, *ast.SelectorExpr, *ast.IndexExpr:
					out = append(out, verb)
				case *ast.Ident:
					if x.Name == recv && recv != "" {
						out = append(out, verb)
					}
				}
				return true
			}
			// A verb written as a string literal, as connect.CreateToken's
			// c.core.Do(ctx, "POST", pathTokenCreate, body, &out) does.
			out = append(out, verbLiterals(node, name)...)
		}
		return true
	})
	return out
}

// verbLiterals returns the HTTP verbs named by string-literal arguments of a call
// to the named function.
//
// It is restricted to a do-family transport so that a method which merely compares
// a string against "GET" is not read as sending one. The literal has to be a
// whole verb, and no path in this SDK is spelled like one.
func verbLiterals(call *ast.CallExpr, fnName string) []string {
	if !doHelpers[fnName] {
		return nil
	}
	var out []string
	for _, arg := range call.Args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}
		if verb := strings.ToUpper(strings.Trim(lit.Value, `"`)); verbHelpers[verb] != "" ||
			verb == "GET" || verb == "POST" || verb == "PUT" || verb == "PATCH" ||
			verb == "DELETE" || verb == "HEAD" {
			out = append(out, verb)
		}
	}
	return out
}

// calleesIn names the same-package functions and methods the body calls, so the
// walk can follow delegation. A call through a package qualifier is skipped
// because the function is not in this package and the index would miss it.
func calleesIn(fn *ast.FuncDecl) []string {
	var out []string
	add := func(name string) {
		if name != "" && name != fn.Name.Name {
			out = append(out, name)
		}
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			add(fun.Name)
		case *ast.SelectorExpr:
			// c.helper(...) is same-package; pkg.helper(...) is not.
			if _, isIdent := fun.X.(*ast.Ident); isIdent {
				add(fun.Sel.Name)
			}
		}
		return true
	})
	return out
}

func sortedVerbs(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// verbMismatch is the divergence CheckVerb produces.
//
// The name is the documented verb, so the finding reads "POST disagrees" rather
// than repeating the symbol, and the detail carries both sides so a report can
// print the disagreement without a lookup.
func verbMismatch(documented, sent string) Divergence {
	return Divergence{
		Kind: VerbMismatch,
		Name: documented,
		Detail: fmt.Sprintf("the page documents %s but the SDK method sends %s, so the request "+
			"cannot succeed; a path comparison cannot see this, because the path matches",
			documented, sent),
	}
}

// CompareVerb is the verb check: the verb the page declares against the verb the
// SDK method sends.
//
// It is a separate entry point from Compare rather than a sixth arm of it, because
// the two answer different questions. Compare reads the documented response and
// asks what the SDK can decode; this reads the documented request and asks what the
// SDK sends. Folding them together would have made a response-shape row carry a
// request finding, and would have made the not-comparable reason, which is about
// paths, read as though it also voided the verb.
//
// An unreadable verb, and a method that sends several verbs, both yield no
// divergence: neither is a statement about what the page documents.
//
// The unreadable case is a defect in this extractor rather than a fact about the
// SDK, and it is deliberately kept out of the baseline. A row recording it would
// put a harness defect in the one file whose purpose is to record SDK defects
// against the documentation, under a reason that would have to be "the tool could
// not look" - which is not a reason about Webull. It is instead required to be zero
// by TestVerbExtractorCoversEverySymbol, so the extractor cannot rot silently.
//
// The several-verbs case is a property of the method: client.CreateToken creates a
// token, polls its status and revokes it, so it has no single counterpart verb to
// compare against one page. That test lists them so the skip stays visible.
func CompareVerb(f Fixture, symbol string) []Divergence {
	documented := strings.ToUpper(strings.TrimSpace(f.Documented.Method))
	if documented == "" {
		return nil
	}
	// A row whose path does not match is not this page's call, so its verb is not
	// comparable either; the path divergence already reports the row as not
	// comparable and this must not also claim a verb disagreement about it.
	if f.SDKPathMatch != PathMatchSame {
		return nil
	}
	sent := SDKVerbOf(symbol)
	if !sent.Single() {
		return nil
	}
	if sent.Verb != documented {
		return []Divergence{verbMismatch(documented, sent.Verb)}
	}
	return nil
}

// moduleGoFiles lists the non-test Go files of the module rooted at dir.
//
// A subdirectory holding its own go.mod is a separate module and is skipped, so
// that scanning the root and then scanning that module does not count its structs
// twice. Each nested module is scanned as its own entry instead; see
// ScannedModuleDirs.
//
// gen/ is skipped because it is generated protobuf whose tags are not this
// repository's to maintain, and which .golangci.yml already excludes.
func moduleGoFiles(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name != "." && (strings.HasPrefix(name, ".") || name == "testdata" ||
				name == "vendor" || name == "site" || name == "gen") {
				return filepath.SkipDir
			}
			if p != dir && hasGoMod(p) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".go") && !strings.HasSuffix(d.Name(), "_test.go") {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out, err
}

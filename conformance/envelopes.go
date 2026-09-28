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
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/shing1211/webullapi4go/data"
)

// This file exists because six of the SDK's response envelopes are unexported,
// and an unexported type cannot be named from another package -- not in a
// reflect.TypeOf call, not in a type switch, not in a composite literal. The
// symbol table's own rule is that the type it names is the one the method's
// `var out` declares, and for those six that type is data.corpActionResponse,
// data.optionBarsResponse, data.optionContractsResponse,
// data.cryptoInstrumentsResponse, data.instrumentProfilesV3Resp, and
// instrument_v3.go's and corporate_actions.go's method bodies are the only place
// that fact is written down.
//
// The two obvious workarounds are both wrong. Naming the public projection the
// method returns (data.StockProfilesV3Result, data.CryptoInstrumentsResult,
// data.OptionContractsResult) compares a tagless Go struct against a wire
// contract, which reports a missing name for every documented property -- three
// or four phantom findings per row, all of them artifacts of the harness's own
// choice. Dropping the rows as "not comparable" would be honest about the gap and
// silent about the cost: six documented endpoints would stop being checked at
// all, which is the one outcome the package documentation says this instrument
// exists to prevent.
//
// So the type is read from the source instead. The declarations are parsed with
// go/parser, and the fields are rebuilt into a struct type with reflect.StructOf.
// Three properties make that a reading rather than a guess:
//
//   - It fails closed. A field whose type is not a builtin or a name this file
//     registers is an error, not a silently dropped field, and an embedded field
//     is an error, so a declaration this parser cannot represent makes the gate
//     red instead of quietly comparing against a type that is not the SDK's.
//   - It cannot disagree with the source. The tags and field names in the rebuilt
//     type are the ones the parser read; there is no second copy of the fact
//     anywhere for the two to drift apart. TestEnvelopeTypesMatchTheSource
//     re-parses every envelope and compares the rebuild field by field.
//   - It is checked on every run. The parse happens inside the comparison, so a
//     renamed or reshaped envelope changes the observed set and the baseline
//     gate fails, which is the same signal a renamed exported type would give.
//
// The cost is that the comparison now reads files under ../<package>. That is
// read-only, it is why the fixtures are embedded and the baseline is not
// regenerable, and it is confined to the five declarations listed above.

// envelopeNamedTypes registers the SDK types a source-declared envelope's fields
// refer to by name. Only the element types of those five envelopes are needed;
// a field naming anything else is an error, so this map cannot become a silent
// gap. A registered name resolves to the real exported type, so everything below
// it is compared by ordinary reflection.
var envelopeNamedTypes = map[string]reflect.Type{
	"data.CorporateAction":  reflect.TypeOf(data.CorporateAction{}),
	"data.CryptoInstrument": reflect.TypeOf(data.CryptoInstrument{}),
	"data.OptionContract":   reflect.TypeOf(data.OptionContract{}),
	"data.OptionSymbolBars": reflect.TypeOf(data.OptionSymbolBars{}),
	"data.StockInstrument":  reflect.TypeOf(data.StockInstrument{}),
}

// envelopeBuiltins are the predeclared types a field may name. It is the same set
// the harness would otherwise reach through reflect, listed so that a field
// naming a type the parser does not understand produces one clear error.
var envelopeBuiltins = map[string]reflect.Type{
	"bool":    reflect.TypeOf(false),
	"float64": reflect.TypeOf(float64(0)),
	"int":     reflect.TypeOf(int(0)),
	"int32":   reflect.TypeOf(int32(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"string":  reflect.TypeOf(""),
}

// envelopeField is one field of a source-declared envelope, as parsed.
type envelopeField struct {
	// Name is the Go field name.
	Name string
	// Type is the field's type as written in the source, without any package
	// qualifier, because a same-package type is written unqualified.
	Type string
	// Tag is the struct tag with its source quoting removed, so it is directly
	// usable as a reflect.StructTag, or "" when the field carries none.
	Tag string
	// Pos is the file and line the field is declared on, for an error message
	// and for a drift report.
	Pos string
}

// envelopeCache memoises resolutions, keyed by the "pkg.Type" spelling. The
// comparison resolves the same handful of envelopes once per manifest row, and a
// resolved type is immutable, so caching is safe and keeps the parse off the hot
// path of a repeated report run.
//
// The same map is how a rebuilt type keeps its name. reflect.StructOf builds an
// anonymous struct, and t.String() on one prints the whole field list with escaped
// tags, which is unreadable in a divergence detail and would change if the
// declaration did. The declared spelling is what belongs there.
var (
	envelopeCacheMu sync.RWMutex
	envelopeCache   = map[string]reflect.Type{}
	envelopeNames   = map[reflect.Type]string{}
)

// declaredName is the "pkg.Type" spelling a rebuilt type was built from, if it is
// a rebuilt type at all. Anything else is named by the Go type system already.
func declaredName(t reflect.Type) (string, bool) {
	envelopeCacheMu.RLock()
	defer envelopeCacheMu.RUnlock()
	name, ok := envelopeNames[t]
	return name, ok
}

// envelopeType resolves a "pkg.Type" spelling to the struct type the SDK package
// declares under that name.
//
// The returned type carries the source's field names and json tags and nothing
// else: it has no methods, so it cannot stand in for the SDK type at run time.
// That is the right trade for this use, because the five checks read field tags,
// field types and the top-level kind, and the decode check exercises exactly the
// same encoding/json behaviour a tagless mirror of the SDK type would.
func envelopeType(spelling string) (reflect.Type, error) {
	envelopeCacheMu.RLock()
	if t, ok := envelopeCache[spelling]; ok {
		envelopeCacheMu.RUnlock()
		return t, nil
	}
	envelopeCacheMu.RUnlock()

	t, err := resolveEnvelope(spelling)

	envelopeCacheMu.Lock()
	defer envelopeCacheMu.Unlock()
	// A failure is not cached, so a declaration this parser cannot represent is
	// re-reported on every run instead of being remembered as absent.
	if err == nil {
		envelopeCache[spelling] = t
		envelopeNames[t] = spelling
	}
	return t, err
}

// resolveEnvelope does the work behind envelopeType.
func resolveEnvelope(spelling string) (reflect.Type, error) {
	pkg, name, ok := strings.Cut(spelling, ".")
	if !ok || pkg == "" || name == "" {
		return nil, fmt.Errorf("envelope %q is not spelled \"package.TypeName\"", spelling)
	}
	dir, err := packageDir(pkg)
	if err != nil {
		return nil, err
	}
	fields, _, err := parseEnvelopeFields(dir, pkg, name)
	if err != nil {
		return nil, err
	}
	return buildEnvelopeType(pkg, name, fields)
}

// parseEnvelopeFields reads the fields of `type <name> struct` out of the
// non-test Go files of dir, and returns them together with the position of the
// declaration.
func parseEnvelopeFields(dir, pkg, name string) (fields []envelopeField, declPos string, err error) {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil, "", fmt.Errorf("read package %s: %w", pkg, readErr)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		names = append(names, name)
	}
	// Sorted so a declaration that moves between two files is reported from the
	// same one on every run.
	sort.Strings(names)

	fset := token.NewFileSet()
	for _, file := range names {
		path := filepath.Join(dir, file)
		parsed, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			return nil, "", fmt.Errorf("parse %s: %w", path, perr)
		}
		for _, decl := range parsed.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != name {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return nil, "", fmt.Errorf("%s declares %s.%s as %T, not a struct; "+
						"this resolver handles only the response envelopes",
						path, pkg, name, ts.Type)
				}
				if st.Fields == nil {
					return []envelopeField{}, fset.Position(ts.Pos()).String(), nil
				}
				for _, f := range st.Fields.List {
					if len(f.Names) == 0 {
						return nil, "", fmt.Errorf("%s: %s.%s has an embedded field at %s; "+
							"an embedded struct is inlined by encoding/json and this "+
							"resolver does not rebuild the inlined shape, so the row is "+
							"refused rather than compared against the wrong type",
							path, pkg, name, fset.Position(f.Pos()))
					}
					if len(f.Names) > 1 {
						return nil, "", fmt.Errorf("%s: %s.%s declares %d names on one "+
							"field at %s, which no envelope in this table does",
							path, pkg, name, len(f.Names), fset.Position(f.Pos()))
					}
					tag := ""
					if f.Tag != nil {
						// go/ast reports a tag as the source literal, so a raw
						// string arrives wrapped in backticks. reflect.StructTag
						// then fails every lookup on it and wireNameOf falls back
						// to the Go field name, which agrees only for camelCase
						// and silently loses a snake_case name such as
						// pagination_key. Unquoting here is what makes such a name
						// comparable at all.
						//
						// A literal that will not unquote is refused rather than
						// defaulted to the empty tag, because defaulting is the
						// exact quiet failure this line exists to remove.
						unquoted, uerr := strconv.Unquote(f.Tag.Value)
						if uerr != nil {
							return nil, "", fmt.Errorf("%s: %s.%s field %s has a tag at %s that is not a quoted string: %w",
								path, pkg, name, f.Names[0].Name, fset.Position(f.Tag.Pos()), uerr)
						}
						tag = unquoted
					}
					fields = append(fields, envelopeField{
						Name: f.Names[0].Name,
						Type: types.ExprString(f.Type),
						Tag:  tag,
						Pos:  fset.Position(f.Pos()).String(),
					})
				}
				return fields, fset.Position(ts.Pos()).String(), nil
			}
		}
	}
	return nil, "", fmt.Errorf("no non-test file in %s declares `type %s struct`", dir, name)
}

// buildEnvelopeType turns parsed fields into a struct type, resolving each field's
// declared type. It recovers a reflect.StructOf panic into an error so a field
// combination the reflector rejects is reported as a message rather than a stack.
func buildEnvelopeType(pkg, name string, fields []envelopeField) (typ reflect.Type, err error) {
	defer func() {
		if r := recover(); r != nil {
			typ, err = nil, fmt.Errorf("rebuild %s.%s: %v", pkg, name, r)
		}
	}()

	// The type is named so a divergence message reads "data.corpActionResponse
	// decodes it as object" rather than naming an anonymous struct.
	reflected := make([]reflect.StructField, 0, len(fields))
	seen := map[string]bool{}
	for _, f := range fields {
		if !isExportedIdent(f.Name) {
			return nil, fmt.Errorf("%s.%s field %s is unexported; encoding/json "+
				"ignores it, and reflect.StructOf cannot rebuild an unexported "+
				"field, so the row is refused rather than compared against a type "+
				"that carries a name the wire never has", pkg, name, f.Name)
		}
		if seen[f.Name] {
			return nil, fmt.Errorf("%s.%s declares %s twice", pkg, name, f.Name)
		}
		seen[f.Name] = true
		ft, ferr := resolveFieldType(pkg, f)
		if ferr != nil {
			return nil, fmt.Errorf("%s.%s field %s at %s: %w", pkg, name, f.Name, f.Pos, ferr)
		}
		reflected = append(reflected, reflect.StructField{
			Name: f.Name,
			Type: ft,
			Tag:  reflect.StructTag(f.Tag),
		})
	}
	return reflect.StructOf(reflected), nil
}

// resolveFieldType resolves one declared field type.
//
// The grammar is deliberately tiny: a builtin, a slice or pointer of one, or a
// name registered above. Anything else is an error. A wider grammar would be
// code that can only ever be exercised by a future envelope, and an unexercised
// branch in the type resolver is exactly the kind of quiet gap this package is
// built to not have.
func resolveFieldType(pkg string, f envelopeField) (reflect.Type, error) {
	if t, ok := envelopeBuiltins[f.Type]; ok {
		return t, nil
	}
	if registered, ok := envelopeNamedTypes[pkg+"."+f.Type]; ok {
		return registered, nil
	}
	for prefix, wrap := range map[string]func(reflect.Type) reflect.Type{
		"[]": reflect.SliceOf,
		"*":  reflect.PointerTo,
	} {
		elem, ok := strings.CutPrefix(f.Type, prefix)
		if !ok {
			continue
		}
		if t, ok := envelopeBuiltins[elem]; ok {
			return wrap(t), nil
		}
		if registered, ok := envelopeNamedTypes[pkg+"."+elem]; ok {
			return wrap(registered), nil
		}
	}
	if strings.HasPrefix(f.Type, "[]") || strings.HasPrefix(f.Type, "*") {
		return nil, fmt.Errorf("%q is not a builtin or a registered SDK type, so the "+
			"element type cannot be resolved", f.Type)
	}
	return nil, fmt.Errorf("%q is not a builtin and not one of the registered SDK "+
		"types this resolver knows, so the field type cannot be resolved", f.Type)
}

// isExportedIdent reports whether name can begin an exported Go identifier, so a
// field reflect.StructOf would refuse is refused with a message naming the field
// rather than a panic from the reflector.
func isExportedIdent(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	return c >= 'A' && c <= 'Z'
}

// modulePath is the root module's path, and goModName the file that declares it.
// Together they are what tells the SDK's source tree apart from a directory that
// merely happens to contain a data/ package.
const (
	modulePath = "github.com/shing1211/webullapi4go"
	goModName  = "go.mod"
)

// moduleRoot returns the directory holding the root module's go.mod.
//
// runtime.Caller gives the compile-time path of this file, which is the direct
// route and is what a plain `go test` uses. It is not always a usable path: a
// -trimpath build rewrites it to an import-path-shaped string. So the walk up
// from the working directory is kept as a second route, and the go.mod check is
// what tells the two apart -- without it, a directory that merely happens to
// contain data/ would be read as the SDK.
func moduleRoot() (string, error) {
	var tried []string
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(filepath.Dir(file))
		if isModuleRoot(dir) {
			return dir, nil
		}
		tried = append(tried, dir)
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("locate the module root: %w", err)
	}
	for dir := wd; ; {
		if isModuleRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no %s declaring module %s found at or above %s "+
		"(also tried %s); the symbol table needs the SDK sources to read an "+
		"unexported response envelope", goModName, modulePath, wd, strings.Join(tried, ", "))
}

// isModuleRoot reports whether dir holds the root module's go.mod.
//
// It never returns an error, because a directory on the way up that has no
// readable go.mod is the ordinary case rather than a failure, and the caller
// wants a yes or no per directory. The path is built from a directory this file
// walked itself, the leaf is a constant, and the content is compared against a
// constant module path, so nothing outside the repository reaches the read.
func isModuleRoot(dir string) bool {
	//nolint:gosec // the path is this function's own walk, the leaf is a constant, and the content is matched against modulePath
	raw, err := os.ReadFile(filepath.Join(dir, goModName))
	if err != nil {
		return false
	}
	rest, ok := strings.CutPrefix(strings.TrimSpace(string(raw)), "module ")
	if !ok {
		return false
	}
	return strings.TrimSpace(strings.SplitN(rest, "\n", 2)[0]) == modulePath
}

// packageDir returns the directory a root-module package's sources live in. The
// harness reads only root-module packages; broker/ is a separate module and its
// types are out of scope by a boundary no edit here can cross.
func packageDir(pkg string) (string, error) {
	root, err := moduleRoot()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(root, filepath.FromSlash(pkg))
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("package %s has no source directory at %s; the "+
			"symbol table can only read a package of the root module", pkg, dir)
	}
	return dir, nil
}

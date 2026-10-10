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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoStructHasADuplicateJSONTag is the whole-module version of the two
// per-type tests. It is here because the class is invisible to every other check
// in this harness: a duplicated tag makes no name *missing*, and this harness looks
// for missing names. See tags.go for how FDPosition hid that way.
func TestNoStructHasADuplicateJSONTag(t *testing.T) {
	root, err := moduleRoot()
	if err != nil {
		t.Fatalf("moduleRoot: %v", err)
	}
	dirs, err := ScannedModuleDirs(root)
	if err != nil {
		t.Fatalf("ScannedModuleDirs: %v", err)
	}

	structs := 0
	var found []DuplicateTag
	for _, dir := range dirs {
		dupes, derr := ScanDuplicateTags(dir)
		if derr != nil {
			t.Fatalf("ScanDuplicateTags(%s): %v", dir, derr)
		}
		rel, _ := filepath.Rel(root, dir)
		if rel == "." {
			rel = "(root module)"
		}
		n, nerr := countStructs(dir)
		if nerr != nil {
			t.Fatalf("countStructs(%s): %v", dir, nerr)
		}
		structs += n
		t.Logf("%-46s %3d structs, %d collision(s)", rel, n, len(dupes))
		found = append(found, dupes...)
	}
	t.Logf("scanned %d structs across %d module directories", structs, len(dirs))
	if structs == 0 {
		t.Fatal("no structs were scanned, so this checked nothing")
	}
	if len(found) > 0 {
		t.Errorf("%d struct field pair(s) claim the same json name, so encoding/json drops "+
			"every one of them and no error is reported:", len(found))
		for _, d := range found {
			t.Errorf("  %s", d)
		}
	}
}

// TestDuplicateTagScannerBites is how the scan is known to work. The per-type tests
// cover the two structs that shipped the defect, but a scanner that always returns
// nothing would pass them too.
func TestDuplicateTagScannerBites(t *testing.T) {
	dir := t.TempDir()
	// One struct with a real collision, and two that must not be reported: a field
	// differing only after the comma, which encoding/json does not treat as
	// distinct, and a field with no tag at all.
	src := `package sample

type Collides struct {
	Open  string ` + "`" + `json:"pl"` + "`" + `
	Close string ` + "`" + `json:"pl,omitempty"` + "`" + `
}

type WithOmitOnly struct {
	A string ` + "`" + `json:"a"` + "`" + `
	B string ` + "`" + `json:"b,omitempty"` + "`" + `
}

type Untagged struct {
	A string
	B string
	skip string ` + "`" + `json:"-"` + "`" + `
}
`
	path := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}

	dupes, err := ScanDuplicateTags(dir)
	if err != nil {
		t.Fatalf("ScanDuplicateTags: %v", err)
	}
	if len(dupes) != 1 {
		t.Fatalf("scan found %d collision(s), want exactly 1: %v", len(dupes), dupes)
	}
	d := dupes[0]
	if d.Type != "sample.Collides" {
		t.Errorf("type = %q, want sample.Collides", d.Type)
	}
	if d.First != "Open" || d.Second != "Close" {
		t.Errorf("fields = %q and %q, want Open and Close", d.First, d.Second)
	}
	if d.Tag != "pl" {
		t.Errorf("tag = %q, want pl", d.Tag)
	}
	if !strings.Contains(d.String(), "encoding/json") {
		t.Errorf("String() must say why this is silent, got %q", d.String())
	}
}

// countStructs is what makes the coverage claim in the test name checkable: a scan
// that walked no files would otherwise report zero collisions and look clean.
func countStructs(dir string) (int, error) {
	files, err := moduleGoFiles(dir)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, path := range files {
		parsed, perr := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if perr != nil {
			return 0, perr
		}
		for _, decl := range parsed.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if _, isStruct := ts.Type.(*ast.StructType); isStruct {
						n++
					}
				}
			}
		}
	}
	return n, nil
}

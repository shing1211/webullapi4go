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
	"strings"
)

// This file generalises one finding.
//
// brokerfd.FDPosition tagged both UnrealizedPL and RealizedPL
// json:"unrealized_pl". encoding/json drops both fields when two fields of one
// struct claim the same name - verified directly, and with no error reported - so
// open P&L had been silently zero for every caller and realized_pl was never
// decodable at all.
//
// It survived 285 rows of conformance work because the harness looks for names
// that are *missing*. Nothing was missing: both tags were present, and tagsOf
// collects wire names into a map and keeps one of the two, so the name read as
// covered while the decoder filled neither field. A name-lookup harness has no
// blind spot here; a tag-lookup harness is blind to it completely.
//
// Two per-type tests now pin FDPosition and broker.Position, but a per-type test
// only catches the type it names. This is the whole-module version.

// DuplicateTag is one struct in which two fields encode under the same wire name.
type DuplicateTag struct {
	// File is the source file, relative to the scanned directory.
	File string
	// Line is the line of the second field, which is the one a reader finds.
	Line int
	// Type is the struct's name, qualified with its package.
	Type string
	// First and Second are the two field names.
	First, Second string
	// Tag is the wire name both claim.
	Tag string
}

func (d DuplicateTag) String() string {
	return fmt.Sprintf("%s:%d: %s fields %s and %s both encode under %q; encoding/json "+
		"drops both when two fields of one struct claim the same name, so neither is "+
		"ever filled and neither reads as missing", d.File, d.Line, d.Type, d.First, d.Second, d.Tag)
}

// ScanDuplicateTags walks the Go files under dir and reports every struct in which
// two fields carry the same json wire name.
//
// A field is skipped when it has no json tag, when the tag is "-", or when the name
// before the first comma is empty, which is how an embedded or untagged field is
// spelled. The name is taken before the comma because that is what encoding/json
// matches on: `json:"name,omitempty"` and `json:"name"` collide.
func ScanDuplicateTags(dir string) ([]DuplicateTag, error) {
	files, err := moduleGoFiles(dir)
	if err != nil {
		return nil, err
	}
	var out []DuplicateTag
	for _, path := range files {
		fset := token.NewFileSet()
		parsed, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			// A file that does not parse is a build failure, not a finding here, and
			// swallowing it would make the scan quietly cover less than it claims.
			return nil, fmt.Errorf("parse %s: %w", path, perr)
		}
		out = append(out, duplicateTagsIn(fset, parsed)...)
	}
	return out, nil
}

func duplicateTagsIn(fset *token.FileSet, file *ast.File) []DuplicateTag {
	var out []DuplicateTag
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			seen := map[string]*ast.Field{}
			for _, f := range st.Fields.List {
				name := jsonTagName(f)
				if name == "" {
					continue
				}
				if prev, dup := seen[name]; dup {
					// Report once per collision, naming the first field that claimed it.
					out = append(out, DuplicateTag{
						File:  fset.Position(f.Pos()).Filename,
						Line:  fset.Position(f.Pos()).Line,
						Type:  file.Name.Name + "." + ts.Name.Name,
						First: fieldName(prev), Second: fieldName(f), Tag: name,
					})
					continue
				}
				seen[name] = f
			}
		}
	}
	return out
}

// jsonTagName is the wire name a field claims, or "" when it claims none.
func jsonTagName(f *ast.Field) string {
	if f.Tag == nil {
		return ""
	}
	raw := strings.Trim(f.Tag.Value, "`")
	// The tag is `json:"name,omitempty"`.
	i := strings.Index(raw, `json:"`)
	if i < 0 {
		return ""
	}
	rest := raw[i+len(`json:"`):]
	j := strings.IndexByte(rest, '"')
	if j < 0 {
		return ""
	}
	name := rest[:j]
	if k := strings.IndexByte(name, ','); k >= 0 {
		name = name[:k]
	}
	if name == "" || name == "-" {
		return ""
	}
	return name
}

func fieldName(f *ast.Field) string {
	if len(f.Names) == 0 {
		return "(embedded)"
	}
	names := make([]string, 0, len(f.Names))
	for _, n := range f.Names {
		names = append(names, n.Name)
	}
	return strings.Join(names, ", ")
}

// ScannedModuleDirs are the directories the duplicate-tag scan treats as module
// roots: the repository root, and each nested module, identified by its own go.mod.
//
// The root walk already descends into broker/, so that module's structs are covered
// by the root entry even though it is separately linted. The nested example modules
// are not reached that way, because a nested module's go.mod makes the Go toolchain
// skip it during a ./... walk, and this scan mirrors that by looking for go.mod
// rather than by walking blindly. That is also what keeps the coverage claim honest:
// every directory here is a module the build actually compiles.
func ScannedModuleDirs(root string) ([]string, error) {
	dirs := []string{root}

	// A nested module sits at examples/<name>, matching the Makefile's MODULES and
	// the CI lint matrix.
	examples := filepath.Join(root, "examples")
	entries, err := os.ReadDir(examples)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			dir := filepath.Join(examples, e.Name())
			if hasGoMod(dir) {
				dirs = append(dirs, dir)
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// A nested module can also sit directly under the root, as broker/ does.
	rootEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, e := range rootEntries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		if dir == examples || !hasGoMod(dir) {
			continue
		}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}

func hasGoMod(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "go.mod"))
	return err == nil && !info.IsDir()
}

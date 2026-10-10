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
	"encoding/json"
	"sort"
	"testing"
)

// readLiveManifest reads the live manifest as the tests need it: the entries and
// the size accounts, typed rather than pulled apart field by field.
//
// The type is local to the test file and mirrors only what these tests read, for
// the reason loadLiveManifest in fixtures_test.go gives: the file is written by an
// example program, and a strict parse here would fail on a field this package
// does not use and is entitled to ignore.
type testLiveManifest struct {
	Totals struct {
		SkeletonBytes       int      `json:"skeletonBytes"`
		DecodedCleanly      int      `json:"decodedCleanly"`
		DecodeRejected      int      `json:"decodeRejected"`
		WithoutElementShape []string `json:"withoutElementShape"`
	} `json:"totals"`
	Size struct {
		ArrayElements         int `json:"arrayElements"`
		DistinctElementShapes int `json:"distinctElementShapes"`
	} `json:"size"`
	Entries []liveEntry `json:"entries"`
}

// readLiveManifest reads and parses the committed live manifest.
func readLiveManifest(t *testing.T) testLiveManifest {
	t.Helper()
	raw, err := fixturesFS.ReadFile(testdataDir + "/" + liveManifestFile)
	if err != nil {
		t.Fatalf("read %s: %v", liveManifestFile, err)
	}
	var m testLiveManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse %s: %v", liveManifestFile, err)
	}
	if len(m.Entries) == 0 {
		t.Fatalf("%s records no entries, so the live tree is unindexed", liveManifestFile)
	}
	return m
}

// decodeRejected is the live manifest's record that the SDK's own type could not
// unmarshal a live body. Four rows carry it, and each is a container-kind
// disagreement between the server's answer and the SDK's type, which is the class
// the documented comparison is structurally blind to. Three of the four are also
// invisible to the shape check, because for them the page and the SDK agree and
// only the server differs; data.GetStockInstruments is the fourth and carries a
// top-level-shape-mismatch as well, since there the page disagrees with the type.
func (m testLiveManifest) decodeRejected() []string {
	var out []string
	for _, e := range m.Entries {
		if !e.DecodedCleanly {
			out = append(out, e.Symbol)
		}
	}
	sort.Strings(out)
	return out
}

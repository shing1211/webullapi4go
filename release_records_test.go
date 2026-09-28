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

package webullapi4go

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// This gate exists because the release record fell behind three releases running.
//
// v2.1.19 and v2.1.20 were tagged with no CHANGELOG entry at all. The runs index
// had skipped v2.1.16 and v2.1.18 as well as both of those. Both status documents
// went on naming v2.1.17 as the latest tag after v2.1.20 existed. Nothing in the
// build could see any of it: a path comparison cannot observe how a number in a
// document was arrived at, and a fixture byte check cannot observe that a release
// has no record. Each gap was found by a person checking after the fact, which is
// the same way the wrong trade denominator in v2.1.18 was found.
//
// So this asserts the records exist. It is a narrow gate and it is deliberately
// bounded rather than comprehensive, because the alternative is a gate that cannot
// pass: 24 earlier tags have no runs-index row and v0.8.0 has no changelog entry,
// and writing those rows now would be inventing records for releases whose run
// artifacts were never kept. A gate that cannot pass is a gate nobody runs, which
// is worth less than a narrower one that does.

// recordsCompleteFrom is the tag from which this gate applies.
//
// It is not an arbitrary starting point. From v2.1.10 onward every tagged release
// carries both a CHANGELOG entry and a runs-index row, so the window is the whole
// period for which the record has actually been kept per release. Moving it earlier
// would require backfilling; TestReleaseRecordCutoffNamesATag exists so the constant
// cannot be quietly raised past the newest tag to make the gate pass vacuously.
const recordsCompleteFrom = "v2.1.10"

// statusDocHeaders are the two documents whose "Latest repository tag" line is
// asserted. They are kept in step by hand and this is what checks that they are.
var statusDocHeaders = []string{
	"IMPLEMENTATION_STATUS.md",
	"docs/implementation-status.md",
}

// changelogEntryRE matches a Keep a Changelog heading for a version, e.g.
// "## [2.1.21] - 2026-09-28".
var changelogEntryRE = regexp.MustCompile(`(?m)^## \[(\d+\.\d+\.\d+)\]`)

// headerTagRE matches the tag named by a status document's header line, e.g.
// "- Latest repository tag: **`v2.1.21`** (2026-09-28)".
var headerTagRE = regexp.MustCompile(`Latest repository tag: \*\*` + "`v" + `([0-9]+\.[0-9]+\.[0-9]+)` + "`" + `\*\*`)

// TestReleaseRecordsAreComplete is the gate: every tagged release from
// recordsCompleteFrom onward carries a changelog entry, a runs-index row, and is
// named by both status-document headers.
//
// It skips rather than passes when there are no tags to check. CI checks out with
// fetch-depth 1 and no tags, so a gate that treated "no tags" as "all recorded"
// would report green in CI while never having looked at anything, which is the one
// failure mode this file exists to remove.
func TestReleaseRecordsAreComplete(t *testing.T) {
	tags := annotatedReleaseTags(t)
	if len(tags) == 0 {
		t.Skip("no annotated v* tags in this clone, so there is no release record to check; " +
			"CI checks out with fetch-depth 1 and fetches no tags")
	}
	inWindow := tagsFrom(tags, recordsCompleteFrom)
	if len(inWindow) == 0 {
		t.Fatalf("recordsCompleteFrom is %q but no tag is at or after it; the gate would "+
			"check nothing and report green", recordsCompleteFrom)
	}
	t.Logf("checking %d tag(s) from %s onward, of %d annotated", len(inWindow), recordsCompleteFrom, len(tags))

	changelog := readRepoFile(t, "CHANGELOG.md")
	index := readRepoFile(t, "docs/runs/index.md")

	// The header is read from the root document and the two documents are asserted
	// separately below, because each carries its own copy of the line.
	headerTag := headerTagOf(readRepoFile(t, statusDocHeaders[0]))
	for _, gap := range releaseRecordGaps(inWindow, changelog, index, headerTag) {
		t.Errorf("%s", gap)
	}

	// Per document: the header must not be stale. It is checked as staleness rather
	// than equality on purpose. A release writes the header naming the tag it is
	// about to create, so the header legitimately runs one ahead of the newest
	// existing tag; what must never happen is a tag existing that the header does
	// not name. This is asserted for both documents, not just the one the gap
	// function read, because they are maintained by hand and can drift apart.
	newest := strings.TrimPrefix(inWindow[len(inWindow)-1], "v")
	for _, path := range statusDocHeaders {
		doc := readRepoFile(t, path)
		m := headerTagRE.FindStringSubmatch(doc)
		if m == nil {
			t.Errorf("%s has no 'Latest repository tag' line, so the header cannot be "+
				"checked for staleness", path)
			continue
		}
		if compareVersions(m[1], newest) < 0 {
			t.Errorf("%s names v%s but v%s exists and is newer; the header is stale",
				path, m[1], newest)
		} else {
			t.Logf("%s names v%s, at or ahead of the newest tag v%s", path, m[1], newest)
		}
	}
}

// TestReleaseRecordCutoffNamesATag keeps the cutoff honest.
//
// A cutoff is a way for this gate to check nothing, so the constant must name a
// real annotated tag. Without this, raising recordsCompleteFrom to something past
// HEAD would make TestReleaseRecordsAreComplete pass by having nothing to do.
func TestReleaseRecordCutoffNamesATag(t *testing.T) {
	tags := annotatedReleaseTags(t)
	if len(tags) == 0 {
		t.Skip("no annotated tags in this clone")
	}
	for _, tag := range tags {
		if tag == recordsCompleteFrom {
			return
		}
	}
	t.Fatalf("recordsCompleteFrom is %q, which is not an annotated tag; the gate would "+
		"skip every tag and report green. Either point it at a real tag or record why "+
		"it is a lower bound", recordsCompleteFrom)
}

// TestReleaseRecordGapsAreDetected is how this gate is known to bite.
//
// A check that has never been observed failing is not evidence of anything, and
// the defect this file addresses was invisible precisely because nothing asserted
// it. So each arm is fed a repository that is missing exactly that record and is
// required to notice. Without this table, "the gate passes" and "the gate works"
// would be the same claim.
func TestReleaseRecordGapsAreDetected(t *testing.T) {
	const changelog = "## [2.1.9] - 2026-09-27\n## [2.1.11] - 2026-09-27\n"
	const index = "| 2026-09-27 | some-run | BUILD | released v2.1.9 | v2.1.9 |\n" +
		"| 2026-09-27 | other-run | BUILD | released v2.1.11 | v2.1.11 |\n"
	const header = "- Latest repository tag: **`v2.1.11`** (2026-09-28)\n"

	tests := []struct {
		name      string
		from      string
		changelog string
		index     string
		header    string
		wantGap   bool
		why       string
	}{
		{
			name: "complete", from: "v2.1.9", changelog: changelog, index: index,
			header: header, wantGap: false,
			why: "every record present",
		},
		{
			name: "missing changelog entry", from: "v2.1.9",
			changelog: "## [2.1.9] - 2026-09-27\n", index: index, header: header,
			wantGap: true, why: "v2.1.11 has a row but no heading",
		},
		{
			name: "missing index row", from: "v2.1.9", changelog: changelog,
			index:   "| 2026-09-27 | some-run | BUILD | released v2.1.9 | v2.1.9 |\n",
			header:  header,
			wantGap: true,
			why:     "v2.1.11 has a heading but no row",
		},
		{
			name: "stale header", from: "v2.1.9", changelog: changelog, index: index,
			header:  "- Latest repository tag: **`v2.1.9`** (2026-09-28)\n",
			wantGap: true, why: "v2.1.11 exists and the header names an older tag",
		},
		{
			name: "header one ahead is not stale", from: "v2.1.9", changelog: changelog,
			index:   index,
			header:  "- Latest repository tag: **`v2.1.12`** (2026-09-28)\n",
			wantGap: false,
			why:     "a release writes the tag it is about to cut, so the header may lead",
		},
		{
			name: "cutoff above every tag checks nothing", from: "v9.9.9",
			changelog: changelog, index: index, header: header,
			wantGap: true, why: "no tag is at or after the cutoff, which is a silent pass",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tags := []string{"v2.1.9", "v2.1.11"}
			got := releaseRecordGaps(tagsFrom(tags, tc.from), tc.changelog, tc.index, headerTagOf(tc.header))
			if (len(got) > 0) != tc.wantGap {
				t.Errorf("gaps = %v, wantGap = %v (%s)", got, tc.wantGap, tc.why)
			}
		})
	}
}

// releaseRecordGaps is the whole check as a pure function, so the table above can
// drive it and the repository test can drive the same code.
func releaseRecordGaps(tags []string, changelog, index, headerTag string) []string {
	var gaps []string
	if len(tags) == 0 {
		return []string{fmt.Sprintf("no tag at or after the cutoff, so nothing was checked")}
	}
	for _, tag := range tags {
		version := strings.TrimPrefix(tag, "v")
		if !hasChangelogEntry(changelog, version) {
			gaps = append(gaps, tag+": no CHANGELOG.md entry")
		}
		if !strings.Contains(index, tag) {
			gaps = append(gaps, tag+": no docs/runs/index.md row")
		}
	}
	if newest := tags[len(tags)-1]; headerTag != "" &&
		compareVersions(headerTag, strings.TrimPrefix(newest, "v")) < 0 {
		gaps = append(gaps, fmt.Sprintf("status header names v%s but %s exists", headerTag, newest))
	}
	return gaps
}

// hasChangelogEntry reports whether the file carries a heading for a version.
func hasChangelogEntry(changelog, version string) bool {
	for _, m := range changelogEntryRE.FindAllStringSubmatch(changelog, -1) {
		if m[1] == version {
			return true
		}
	}
	return false
}

// headerTagOf pulls the version out of a status-document header line.
func headerTagOf(header string) string {
	if m := headerTagRE.FindStringSubmatch(header); m != nil {
		return m[1]
	}
	return ""
}

// cloneIsShallow reports whether the working tree is a shallow clone.
//
// This is the signal that decides whether the gate can be trusted at all, and it
// was added because the gate's first real run failed. The first tag-triggered CI
// run checked out with the default fetch-depth of 1, which fetched the tag being
// built and none of its predecessors, so the tag set held exactly one entry. That
// is a partial set that looks complete: it is non-empty, so a "no tags" skip did
// not fire, and the cutoff check then concluded v2.1.10 was not a tag and failed
// six build jobs and a coverage job over a misconfigured clone. Asking git is
// authoritative and works from a worktree too, unlike looking for .git/shallow.
func cloneIsShallow() bool {
	out, err := exec.Command("git", "rev-parse", "--is-shallow-repository").Output()
	if err != nil {
		// Not being able to tell is not being shallow, and claiming otherwise
		// would turn an unanswerable question into a silent skip.
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// annotatedReleaseTags returns every annotated v* tag, ordered by version.
//
// Order is by parsed version, not by string or by creation date, because "v2.1.10"
// sorts before "v2.1.5" as a string and a newest-tag comparison that gets this
// wrong silently stops comparing.
func annotatedReleaseTags(t *testing.T) []string {
	t.Helper()
	if cloneIsShallow() {
		t.Skip("clone is shallow, so its tag set is necessarily partial and a tag list " +
			"that looks complete is not; the release-records CI job checks out full " +
			"history and runs this gate for real")
	}
	out, err := exec.Command("git", "for-each-ref", "--format=%(refname:short) %(objecttype)", "refs/tags").Output()
	if err != nil {
		t.Skipf("git is unavailable in this environment (%v), so there is no release record "+
			"to check; this is a skip, not a pass", err)
	}
	var tags []string
	for _, line := range strings.Split(string(out), "\n") {
		name, kind, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok || kind != "tag" || !strings.HasPrefix(name, "v") {
			continue
		}
		if _, err := parseVersion(strings.TrimPrefix(name, "v")); err != nil {
			continue
		}
		tags = append(tags, name)
	}
	sort.Slice(tags, func(i, j int) bool {
		a, _ := parseVersion(strings.TrimPrefix(tags[i], "v"))
		b, _ := parseVersion(strings.TrimPrefix(tags[j], "v"))
		return a.less(b)
	})
	return tags
}

// tagsFrom returns the tags at or after a cutoff, preserving order.
func tagsFrom(tags []string, from string) []string {
	fv, err := parseVersion(strings.TrimPrefix(from, "v"))
	if err != nil {
		return nil
	}
	var out []string
	for _, tag := range tags {
		if v, err := parseVersion(strings.TrimPrefix(tag, "v")); err == nil && !v.less(fv) {
			out = append(out, tag)
		}
	}
	return out
}

// version is a three-part tag version, compared numerically so 2.1.10 sorts above
// 2.1.9.
type version struct{ major, minor, patch int }

func parseVersion(s string) (version, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return version{}, fmt.Errorf("%q is not major.minor.patch", s)
	}
	var out [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return version{}, fmt.Errorf("%q has a non-numeric part %q", s, p)
		}
		out[i] = n
	}
	return version{major: out[0], minor: out[1], patch: out[2]}, nil
}

func (v version) less(o version) bool {
	switch {
	case v.major != o.major:
		return v.major < o.major
	case v.minor != o.minor:
		return v.minor < o.minor
	default:
		return v.patch < o.patch
	}
}

func compareVersions(a, b string) int {
	av, aerr := parseVersion(a)
	bv, berr := parseVersion(b)
	switch {
	case aerr != nil || berr != nil:
		return 0
	case av.less(bv):
		return -1
	case bv.less(av):
		return 1
	default:
		return 0
	}
}

// readRepoFile reads a file from the repository root. Tests run in the package
// directory, which is the root for this package.
func readRepoFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

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
	"os"
	"testing"
)

// TestLiveBaselineReasonsAreWritten guards the reasons the committed baseline is
// regenerated from.
//
// live-divergences.json is written by hand and checked by TestLiveDivergenceBaseline,
// so regenerating it is an act of judgement: every observed finding has to be
// matched with a reason and an unblock requirement, and a finding nobody has
// agreed to look at is a suppressed one wearing a comment. This is the test that
// makes the act repeatable -- it is driven from the same committed tree the gate
// reads, so it goes stale in exactly the way the gate would.
//
// It skips by default, because regenerating is a decision and not a routine. Run
// it with GEN_LIVE=1 after reading TestLiveDivergenceBaseline's output, and only
// adopt what you have read.
//
// The way to break this test is to add a finding without a reason: the
// regeneration fails naming the finding, and the file on disk is left untouched
// rather than rewritten with a gap in it.
func TestLiveBaselineReasonsAreWritten(t *testing.T) {
	if os.Getenv("GEN_LIVE") == "" {
		t.Skip("set GEN_LIVE=1 to regenerate " + LiveBaselineName)
	}
	run := CompareAllLive()
	for _, e := range run.Errs {
		t.Errorf("the run could not compare everything: %v", e)
	}
	rows := make([]LiveEntry, 0, len(run.Classifications))
	for _, c := range run.Classifications {
		reason, unblock, ok := liveReasonFor(c)
		if !ok {
			t.Errorf("no reason is written for [%s] %s / %s / %q -- %s\n"+
				"Every entry needs a reason and an unblock requirement, because a "+
				"reasonless entry is a suppressed finding rather than a recorded one. "+
				"The file on disk has not been touched.",
				c.Bucket, c.Symbol, c.Kind, c.Name, c.Detail)
			continue
		}
		rows = append(rows, LiveEntry{
			LiveDivergence: LiveDivergence{
				Symbol:    c.Symbol,
				Fixture:   c.Fixture,
				Direction: directionOf(c.Bucket),
				Kind:      DivergenceKind(c.Kind),
				Name:      c.Name,
				Detail:    c.Detail,
			},
			Reason:  reason,
			Unblock: unblock,
			Bucket:  c.Bucket,
		})
	}
	if t.Failed() {
		return
	}
	b := &LiveBaseline{
		Policy:   liveBaselinePolicy,
		Notes:    liveBaselineNotes(run),
		Coverage: run.Coverage,
		Entries:  rows,
	}
	out, err := b.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	// 0600 rather than 0644: the file is written to be read by git and by a
	// reviewer, not by another user on this machine, and a baseline that only its
	// owner can read is the safer default for a file whose whole purpose is to be
	// committed. A checkout materialises it with whatever the repository records.
	if err := os.WriteFile(LiveBaselineName, out, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d entries to %s", len(rows), LiveBaselineName)
}

// directionOf is the direction a bucket's rows are recorded on.
//
// Every classification stands for one finding observed on one or both sides, and
// the direction records the side carrying the evidence. A both or
// both-detail-changed classification is recorded on the live side, because that
// is the side saying something the documentation does not; a docs-only one is
// recorded on the documented side, because that is the only side that said it.
func directionOf(b LiveBucket) LiveDirection {
	if b == LiveDocsOnly {
		return LiveDocsDirection
	}
	return LiveSDKDirection
}

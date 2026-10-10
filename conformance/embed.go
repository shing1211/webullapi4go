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

import "embed"

// fixturesFS carries the committed tree into the test binary.
//
// Embedding rather than reading testdata off disk is deliberate on two counts.
// It removes a filesystem path built from manifest data, which is what a
// reviewer should want removed from a package that reads 194 files. And it makes
// the tree immutable at build time: nothing a test does can regenerate, move, or
// overwrite the input it is checking, which is the property the whole design
// rests on. A test that could rewrite its own fixture could not fail honestly.
//
// The embed pattern is a directory rather than individual files, so a fixture
// added without touching this source is still compiled in. A missing file then
// fails at build time instead of at run time.
//
//go:embed testdata
var fixturesFS embed.FS

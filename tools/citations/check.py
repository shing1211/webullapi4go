# Copyright 2026 shing1211
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Validate the `file:line` citations in this repository's status documents.

Why this exists
---------------
`AGENTS.md` requires live-blocked defects to be recorded with `file:line`
citations, and names `IMPLEMENTATION_STATUS.md` and `docs/implementation-status.md`
as authoritative for them. Across five consecutive releases those citations went
stale and shipped: a comment inserted above one const block invalidated nine
references across three documents, one manifest citation was transposed at birth,
another pointed at a blank line, and a human reviewer caught some of them. Nothing
mechanical checked any of them. This is that check.

Usage
-----
    python tools/citations/check.py                   # the default documents
    python tools/citations/check.py --verbose         # also list what resolved
    python tools/citations/check.py path/to/doc.md    # a different document

Documents
---------
`DEFAULT_DOCUMENTS` are the three status documents: `IMPLEMENTATION_STATUS.md`,
`docs/implementation-status.md` and `AGENTS.md`. `WAIVED_DOCUMENTS` is a
deliberate, greppable exclusion, not an oversight. The document list may be
overridden on the command line for testing; the defaults stay those three.

Citation syntax
---------------
Only inline code spans (backtick-delimited) are read, and only tokens matching
one of these two shapes, where `N` is a decimal line number:

    `path/to/file.go:123`      a single line
    `path/to/file.go:45-49`    a range; both endpoints must be in range
    `path/to/file.go:26,30`    a comma list of singles and/or ranges
    `:63`                      a bare continuation of the previous citation

A bare `:NN` continuation is resolved against the most recent *full* `file:NNN`
citation in the same document and is reported as `unresolved-continuation` when
the document has not cited a file yet. This form is used heavily by the status
documents, so a checker that ignored it would check far less than it appears to.

A citation whose path carries no directory (`assets.go:25`) is resolved in this
order: the most recently cited file in the document with that basename, then a
file of that name sitting beside the most recent citation (`brokerfd/assets.go`),
then a file of that name at the repository root, then the most recent citation
itself. `--verbose` prints the resolution of every citation, inferred or not, so
a reader can see which file each line was read from.

That order is a lookup by memory wearing the clothes of a lookup by name, and it
can bind to the wrong package. 41 Go basenames here name more than one file
(`client.go` seven times, `accounts.go` three), and 70 citable basenames do so
across all file types, so a bare `accounts.go:23` resolves to whichever of them
the document cited most recently: `broker/` or `brokerfd/`, two files that share
a name and nothing else. A resolved basename is therefore authoritative only
when the repository holds exactly one file of that name. When it holds two or
more the citation is reported as `ambiguous-basename`, a warning naming every
candidate and the one actually checked. That is a visibility control, not a
correctness control, and the honest claim is the narrow one: no citation is ever
checked against a file this tool has not named, and a green run never means an
ambiguity went unnoticed. The exit status stays zero, because the alternative is
a gate red on every run of documents that are mostly correct, and a gate nobody
reads catches nothing either.

The last rule in that order has no name to be right about at all: the basename
matches nothing. A guess is still checked, but its findings are downgraded to
`assumed-` warnings, because an assumed file must never fail the gate.

Two shapes of citation cannot be checked at all and are reported rather than
silently skipped: a bare `:NN` continuation with no earlier full citation
(`unresolved-continuation`, an error, since no file can be inferred), and a bare
basename that matches no file at all (`unresolved-basename`, a warning).

Failure conditions (each exits non-zero)
----------------------------------------
1. `missing-file`      the cited file does not exist. Only a path carrying its own
                       directory can produce this: a mistyped bare basename names
                       no file at all, so it is a warning, not a failure
2. `line-out-of-range` a cited line, or a range end, is outside the file
3. `blank-line`        a single cited line is blank or whitespace-only
4. `comment-line`      a single cited line is a pure comment: `//` in Go, `#` in
                       Python, `<!--` in Markdown, per `COMMENT_PREFIXES`
5. `empty-waiver`      a `# citation-waiver:` marker carrying no reason. An
                       unexplained marker is not a waiver: the citation is
                       checked normally and the marker is reported

Plus `unresolved-continuation`, the unresolvable `:NN` case above.

A range is held to a weaker rule on purpose, and that difference is a real
limitation rather than a convenience. A range fails only when *no* line in it is
code, so a range that starts on a GoDoc comment, or that spans a struct whose
interior is field comments, still passes: `brokerfd/assets.go:92-100` starts on
`// GetFDPositions retrieves ...`, and `data/snapshot.go:35-47` is a struct whose
interior is all field comments, and both citations in the current documents are
correct. The number of *passing* ranges with an interior comment or blank line
is reported in the summary so the weakening stays visible. The weakening is not
theoretical: `data/display_quotes.go:41-49` is a range that drifted onto a const
and this tool does not report it, and that is one of the historical misses it
cannot see.

Warning conditions (exit zero, still reported)
----------------------------------------------
`duplicate-label`: the same document cites the same `file:line` twice with
different descriptive labels, where a label is the adjacent inline-code symbol
immediately before the citation, as in `` `GetFDPositions` (`brokerfd/assets.go:92-100`) ``.
This is a *partial* mitigation for the label/number transposition and it does
not prevent that class of error; see the limits below.

`ambiguous-basename`: a bare filename the repository holds more than one copy of,
reported with every candidate and the file this tool read, as described above.

`guessed-basename` and `unresolved-basename` report a bare filename this tool
could not locate from the document, plus any finding produced against the file it
assumed, so an unverifiable citation is never silently dropped.

Waiver
------
A citation on a line carrying `# citation-waiver: <reason>` is skipped and
counted. The marker must sit outside every code span, because a marker inside
one documents the syntax rather than waiving anything — including in this
docstring and in `AGENTS.md`. The reason is not decoration: a marker with nothing
after the colon is not a waiver at all, so the citation is checked and
`empty-waiver` is reported. That is what makes a waiver unable to be added
silently, and it is why a green run means every skipped citation said why.

What this tool will not do
--------------------------
It never rewrites, repairs, or auto-fixes a citation, and it never writes to any
file, including the documents it reads. An auto-rewritten citation is a wrong
citation that looks verified, which is worse than a red gate. There is no
`--fix` and there will not be one; the only output is this report.

Honest limits
-------------
This tool does *not* verify that the cited line contains what the surrounding
prose claims. Doing that needs a symbol-naming convention or semantic analysis,
and neither belongs in a Makefile gate. It therefore does *not* catch a
label/number transposition, where the number is a valid line of the right file
but the wrong line; the duplicate-label warning is a partial mitigation only and
fires far too rarely to be trusted. Nor can it see a stale citation that drifted
onto a line which still reads as code. Measured on 2026-09-27 against a copy of
`CHANGELOG.md`, the only document that still holds the historical misses, it
reports the four that drifted onto a blank or comment line —
`tools/webull-docgen/_common.py:313`, `data/display_quotes.go:28`, `:40`, and
`examples/options-multi-leg/main.go:193` — as 5 `comment-line` and 2 `blank-line`
findings over 4 distinct targets. It cannot see the other three:
`tools/webull-docgen/_common.py:298`, which is the `[` that opens a list;
`data/display_quotes.go:52`; and the `data/display_quotes.go:41-49` range, which
contains a const. A green run is evidence that the citations resolve to existing,
non-blank, non-comment lines inside the cited file, and nothing more.
"""

import argparse
import os
import posixpath
import re
import sys

# Documents checked unless the command line overrides the list.
DEFAULT_DOCUMENTS = (
    "IMPLEMENTATION_STATUS.md",
    "docs/implementation-status.md",
    "AGENTS.md",
)

# Documents deliberately NOT checked. `CHANGELOG.md` is waived because it is
# released history, not a live status document. It narrates the known-stale
# citations of past releases — `tools/webull-docgen/_common.py:313` and `:298`,
# `examples/options-multi-leg/main.go:193`, and `data/display_quotes.go:28`, `:40`,
# `:52`, and `:41-49` — as the record of what was wrong and what was later
# corrected. Those citations quote the defects deliberately: rewriting them to
# match today's line numbers would falsify the release record, so they must stay
# byte-unchanged. That is why the exclusion is a named constant here rather than
# silence in the default list. Measured on 2026-09-27 by passing a copy of the
# file on the command line, the checker reports 7 findings on 4 distinct stale
# targets, which is the waiver working as intended and the reason it exists; the
# other 3 stale targets land on lines that still read as code and are invisible to
# it. The default run never opens the file: `main` skips a waived document before
# reading it, so the measurement has to be taken from outside the default list.
WAIVED_DOCUMENTS = ("CHANGELOG.md",)

# A citation on a line carrying this marker is skipped.
WAIVER_MARKER = "# citation-waiver:"

# First non-whitespace characters of a pure comment line, by file extension. An
# extension with no entry here still gets the range and blank-line checks, and
# the fact that it has no comment rule is counted in the summary.
COMMENT_PREFIXES = {
    ".cfg": "#",
    ".go": "//",
    ".ini": "#",
    ".markdown": "<!--",
    ".md": "<!--",
    ".mod": "//",
    ".py": "#",
    ".pyi": "#",
    ".sh": "#",
    ".toml": "#",
    ".yaml": "#",
    ".yml": "#",
}

# Extensionless basenames that may still be cited.
EXTENSIONLESS = ("Dockerfile", "Makefile")

_HELP = """\
exit status
  0  no failure condition found; warnings do not change the exit status
  1  at least one failure condition found
  2  the tool could not run: a document is missing or unreadable

citation syntax
  `path/to/file.go:123`    a single line
  `path/to/file.go:45-49`  a range, both endpoints in range
  `path/to/file.go:26,30`  a comma list of singles and/or ranges
  `:63`                    bare continuation of the most recent full citation in
                           the same document; unresolvable when there is none

failure conditions (exit 1)
  missing-file             the cited file does not exist. Only a path carrying its
                           own directory can produce this; a mistyped bare
                           basename names no file at all and is a warning
  line-out-of-range        a cited line, or a range end, is outside the file
  blank-line               a single cited line is blank or whitespace-only
  comment-line             a single cited line is a pure comment: // in Go, # in
                           Python, <!-- in Markdown
  empty-waiver             a `# citation-waiver:` marker with no reason after it.
                           An unexplained marker is not a waiver: the citation is
                           checked and this is reported
  unresolved-continuation  a bare `:NN` with no earlier full citation

warnings (exit 0, still reported)
  duplicate-label          the same file:line cited twice in one document with
                           different adjacent labels. This is a partial
                           mitigation for the label/number transposition and it
                           does not prevent one.
  ambiguous-basename       a bare filename the repository holds more than one copy
                           of, reported with every candidate and the file actually
                           checked
  guessed-basename         a bare filename that matches no cited file, no
                           sibling and no root file, checked against an assumed
                           file
  unresolved-basename      a bare filename matching no file at all, so the line
                           it names went unchecked
  assumed-<kind>           a finding produced against an assumed file, so it is
                           not authoritative and cannot fail the gate

a bare filename is resolved against, in order: the most recently cited file with
that basename, a sibling of the most recent citation, a file at the repository
root, and only then a guess. That is authoritative only when the repository holds
exactly one file of that name; otherwise the citation is reported as
`ambiguous-basename` with the candidates and the file checked. --verbose prints
the resolution of every citation, one line per target line rather than per span,
so a span carrying a comma list appears once per item.

a range fails only when it contains no code line at all, so a range may start on
a doc comment or span field comments.

waiver
  A citation on a line carrying `{marker} <reason>` is skipped. The marker must
  sit outside every code span, since one inside a span documents the syntax
  rather than waiving anything. The reason is required: a marker with nothing
  after the colon is not a waiver, the citation is checked, and `empty-waiver` is
  reported. Waived citations are counted in the summary.

this tool never writes to any file
  It reports and exits. There is no --fix: an auto-rewritten citation is a wrong
  citation that looks verified, which is worse than a failing gate.
""".format(marker=WAIVER_MARKER)

_LIMITS = """\
limits: this checks only that a cited line exists, is not blank, and is not a
pure comment. It does not check that the line contains what the prose claims, so
it cannot detect a label/number transposition; the duplicate-label warning is a
partial mitigation only. A stale citation that drifted onto a line which still
reads as code is invisible to it: of the historical misses the release record
names, it reports the four now reading as a blank or comment line and cannot see
the three now reading as code."""

# A path is a sequence of slash-separated name segments whose final segment
# carries an extension starting with a letter. That last rule is what keeps
# `wss://...:8883/mqtt` and a bare host:port pair out of the citation grammar.
_NAME = r"[A-Za-z0-9_][A-Za-z0-9_.\-]*"
_PATH = r"(?:%s/)*%s" % (_NAME, _NAME)
_ITEM = r"[0-9]+(?:-[0-9]+)?"
_SPEC = r"%s(?:,%s)*" % (_ITEM, _ITEM)
_EXTENSION = re.compile(r"\.([A-Za-z][A-Za-z0-9]*)$")

CODE_SPAN_RE = re.compile(r"`([^`\n]+)`")
CITATION_RE = re.compile(r"^(?P<path>%s):(?P<spec>%s)$" % (_PATH, _SPEC))
CONTINUATION_RE = re.compile(r"^:(?P<spec>%s)$" % _SPEC)
IDENTIFIER_RE = re.compile(r"^[A-Za-z_][A-Za-z0-9_.\-]*$")
FENCE_RE = re.compile(r"^\s{0,3}(?:```|~~~)")
WAIVER_RE = re.compile(re.escape(WAIVER_MARKER))

# A label must be adjacent to its citation: at most this many separator
# characters, none of them alphanumeric. `GetFDPositions` (`x.go:9`) qualifies; a
# symbol in a sentence twenty words earlier does not.
_LABEL_WINDOW = 4

EXIT_OK = 0
EXIT_FINDINGS = 1
EXIT_TOOL_ERROR = 2

# The failure conditions above decide the exit status. A citation whose file this
# tool could only infer is reported as a warning instead, because the uncertain
# party there is the tool's resolution heuristic, not the document. That covers a
# guessed basename and every finding derived from one, which is why the check on
# an assumed file is prefixed `assumed-`, and a basename the repository holds
# more than once, whose resolution the document chose and this tool only confirmed.
WARNING_KINDS = ("duplicate-label", "ambiguous-basename",
                 "guessed-basename", "unresolved-basename")

# Directories a basename index skips: build output, dependencies, and every
# dot-directory, which is what keeps `.git` out of the walk.
INDEX_SKIP_DIRS = frozenset(("site", "node_modules", "__pycache__"))

# Resolution mode -> the stats counter it increments.
STATS_MODES = {
    "recalled": "recalled",
    "sibling": "siblings",
    "rooted": "rooted",
    "guessed": "guessed",
}


def severity(kind):
    return "warning" if kind in WARNING_KINDS or kind.startswith("assumed-") else "error"


class ToolError(Exception):
    """A condition that stops the tool itself rather than the documents."""


def read_lines(path):
    """Return the lines of path, tolerating a byte-order mark and CRLF endings."""
    try:
        with open(path, encoding="utf-8-sig") as handle:
            return handle.read().splitlines()
    except OSError as exc:
        raise ToolError("cannot read %s: %s" % (path, exc))


def comment_prefix(target):
    """Return the pure-comment prefix for target, or None when it has no rule."""
    extension = _EXTENSION.search(posixpath.basename(target))
    if extension is None:
        return None
    return COMMENT_PREFIXES.get("." + extension.group(1).lower())


def is_blank(text):
    return not text.strip()


def is_comment(text, prefix):
    return prefix is not None and text.lstrip().startswith(prefix)


def is_citable(path):
    """True when path's final segment names a file this tool may cite."""
    name = posixpath.basename(path)
    return name in EXTENSIONLESS or _EXTENSION.search(name) is not None


def is_citation(token):
    """True when token is a full `file:line` citation this grammar accepts."""
    match = CITATION_RE.match(token)
    return match is not None and is_citable(match.group("path"))


def describe(text):
    """Render a cited line for a one-line report."""
    stripped = text.strip()
    if len(stripped) > 72:
        stripped = stripped[:69] + "..."
    return "'%s'" % (stripped if stripped else "<blank>")


def label_before(previous, span, line):
    """Return the inline-code label adjacent to span, or '<none>'."""
    if previous is None:
        return "<none>"
    separator = line[previous.end():span.start()]
    if len(separator) > _LABEL_WINDOW:
        return "<none>"
    for char in separator:
        if char.isalnum():
            return "<none>"
    token = previous.group(1)
    if not IDENTIFIER_RE.match(token):
        return "<none>"
    if is_citation(token) or CONTINUATION_RE.match(token):
        return "<none>"
    return token


class FileCache:
    """Read each cited file once, keyed by its repository-relative path."""

    def __init__(self, root):
        self.root = root
        self._lines = {}

    def get(self, target):
        """Return (lines, absolute_path) for target, or ([], None) if absent."""
        if target not in self._lines:
            absolute = os.path.join(self.root, *target.split("/"))
            if os.path.isfile(absolute):
                self._lines[target] = (read_lines(absolute), absolute)
            else:
                self._lines[target] = ([], None)
        return self._lines[target]

    def exists(self, target):
        return self.get(target)[1] is not None


class BasenameIndex:
    """Map each citable basename in the repository to every path carrying it.

    This is what makes an inferred resolution reportable. The document's own
    citation history decides which file a bare `accounts.go:23` means, and in a
    repository where 41 basenames name more than one file that decision is a
    guess; comparing the name against the tree is the only way to know whether
    the guess had anywhere to go. Built on first use, so a document of fully
    qualified citations never walks the tree.
    """

    def __init__(self, root):
        self.root = root
        self._by_name = None

    def candidates(self, name):
        """Return every repository-relative path named name, sorted."""
        if self._by_name is None:
            self._build()
        return self._by_name.get(name, ())

    def _build(self):
        found = {}
        for directory, subdirectories, names in os.walk(self.root):
            subdirectories[:] = [entry for entry in subdirectories
                                 if not entry.startswith(".")
                                 and entry not in INDEX_SKIP_DIRS]
            for name in names:
                if not is_citable(name):
                    continue
                absolute = os.path.join(directory, name)
                relative = os.path.relpath(absolute, self.root)
                found.setdefault(name, []).append(
                    relative.replace(os.sep, "/"))
        self._by_name = dict((name, tuple(sorted(paths)))
                             for name, paths in found.items())


def resolve(path_part, history, cache):
    """Resolve a citation path against the earlier citations of the document.

    Returns (target, mode). mode is "literal" for a path carrying its own
    directory, "recalled" for a bare basename matching an earlier citation,
    "sibling" for a bare basename sitting beside the most recent citation,
    "rooted" for a bare basename naming a file at the repository root,
    "guessed" when only the most recent citation is left to assume, and
    "unknown" when the document has cited nothing yet. target is None only in
    that last case.

    Only the first four modes are authoritative, and only where the repository
    holds a single file of that name; `ambiguous_basename` is the caller that
    says so. A "guessed" target is checked anyway, but its findings are
    downgraded to warnings, because an assumed file must never fail the gate.
    """
    if "/" in path_part:
        return path_part, "literal"
    for target in reversed(history):
        if posixpath.basename(target) == path_part:
            return target, "recalled"
    if not history:
        if cache.exists(path_part):
            return path_part, "rooted"
        return None, "unknown"
    sibling = posixpath.join(posixpath.dirname(history[-1]), path_part)
    if cache.exists(sibling):
        return sibling, "sibling"
    if cache.exists(path_part):
        return path_part, "rooted"
    return history[-1], "guessed"


def ambiguous_basename(index, name, token, target):
    """Return an `ambiguous-basename` finding, or None when the name is unique.

    A bare basename the repository holds exactly one copy of resolves the same
    way whichever rule fired, so there is nothing to disclose. Two or more copies
    means the file this tool read is a choice, and a choice made from document
    history rather than from the tree is the one the reader cannot reconstruct.
    """
    candidates = index.candidates(name)
    if len(candidates) < 2:
        return None
    return ("ambiguous-basename",
            "bare `%s` was checked against %s, but %d files in this repository "
            "are named %s: %s. The document's own citation history chose which "
            "one to read, so the file checked is not necessarily the file the "
            "prose means"
            % (token, target, len(candidates), name, ", ".join(candidates)))


def check_item(cache, target, start, end, prefix):
    """Check one cited item. Returns a list of (kind, message) findings."""
    if not cache.exists(target):
        return [("missing-file", "no such file: %s" % target)]
    lines, _ = cache.get(target)
    total = len(lines)
    if start < 1 or end < 1 or end > total or start > end:
        if start > end:
            return [("line-out-of-range",
                     "range %d-%d in %s is inverted" % (start, end, target))]
        offender = start if not 1 <= start <= total else end
        return [("line-out-of-range",
                 "line %d of %s is out of range; the file has %d line%s"
                 % (offender, target, total, "" if total == 1 else "s"))]
    if start == end:
        text = lines[start - 1]
        if is_blank(text):
            return [("blank-line",
                     "line %d of %s is blank or whitespace-only"
                     % (start, target))]
        if is_comment(text, prefix):
            return [("comment-line", "line %d of %s is a comment line %s"
                     % (start, target, describe(text)))]
        return []
    interior = lines[start - 1:end]
    findings = []
    if all(is_blank(text) for text in interior):
        findings.append(("blank-line", "range %d-%d of %s contains no non-blank "
                         "line" % (start, end, target)))
    if all(is_blank(text) or is_comment(text, prefix) for text in interior):
        findings.append(("comment-line", "range %d-%d of %s contains no code "
                         "line" % (start, end, target)))
    return findings


def waiver_reason(line, spans):
    """Return the reason a line offers as a waiver, or None if it offers none.

    Only a marker outside every code span counts, because a marker inside one
    documents the syntax — as this very tool's help and `AGENTS.md` do — and must
    not quietly waive whatever else the sentence cites. The reason is what
    separates a recorded exception from a suppressed gate, so an empty one is not
    a reason: it returns "" and the caller reports `empty-waiver` and checks the
    citation anyway.
    """
    for match in WAIVER_RE.finditer(line):
        if any(span.start() <= match.start() < span.end() for span in spans):
            continue
        return line[match.end():].strip()
    return None


def check_document(document, root, cache, index, verbose=False, stream=None):
    """Check one document. Returns (findings, stats) for that document."""
    findings = []
    stats = {"citations": 0, "full": 0, "continuations": 0, "targets": 0,
             "recalled": 0, "siblings": 0, "rooted": 0, "guessed": 0,
             "ambiguous": 0, "waived": 0, "comment_ranges": 0, "blank_ranges": 0,
             "uncommented": set()}
    lines = read_lines(document)
    history = []
    labels = {}
    in_fence = False
    for number, line in enumerate(lines, 1):
        if FENCE_RE.match(line):
            in_fence = not in_fence
            continue
        if in_fence:
            continue
        spans = list(CODE_SPAN_RE.finditer(line))
        reason = waiver_reason(line, spans)
        waived = reason is not None
        if reason:
            stats["waived"] += sum(1 for span in spans
                                   if is_citation(span.group(1))
                                   or CONTINUATION_RE.match(span.group(1)))
        elif reason is not None:
            # Counted as a checked citation and reported, never as a waiver: the
            # marker is present but unexplained, and the footer must not present
            # it as a recorded exception.
            findings.append(
                (number, "empty-waiver",
                 "`%s` carries no reason, so it waives nothing: every citation "
                 "on this line is checked. An exception that cannot say why it "
                 "is an exception is a suppressed gate, not a recorded one"
                 % WAIVER_MARKER))
            waived = False
        previous = None
        for span in spans:
            token = span.group(1)
            label = label_before(previous, span, line)
            previous = span
            citation = CITATION_RE.match(token)
            continuation = None if citation else CONTINUATION_RE.match(token)
            if citation is not None and not is_citable(citation.group("path")):
                citation = None
                continuation = None
            if citation is None and continuation is None:
                continue
            if citation is not None:
                spec = citation.group("spec")
                target, mode = resolve(citation.group("path"), history, cache)
                if target is None:
                    if not waived:
                        findings.append(
                            (number, "unresolved-basename",
                             "bare `%s` cannot be resolved: no file of that name "
                             "sits at the repository root and the document "
                             "cites no file yet, so the line it names went "
                             "unchecked" % token))
                    continue
                # Appended unconditionally: a bare `:NN` resolves against the
                # most recently *cited* file, so re-citing a file has to move it
                # back to the end. Deduplicating here would silently resolve
                # continuations against a file the prose moved on from.
                history.append(target)
                stats["full"] += 1
                if mode in STATS_MODES:
                    stats[STATS_MODES[mode]] += 1
                if not waived and "/" not in citation.group("path"):
                    ambiguous = ambiguous_basename(
                        index, citation.group("path"), token, target)
                    if ambiguous is not None:
                        stats["ambiguous"] += 1
                        findings.append((number,) + ambiguous)
                if mode == "guessed" and not waived:
                    findings.append(
                        (number, "guessed-basename",
                         "bare `%s` names no file this document has cited, no "
                         "such file beside the latest citation, and none at the "
                         "repository root, so it was checked against %s on the "
                         "assumption that it means that file; that check is not "
                         "authoritative" % (token, target)))
            else:
                spec = continuation.group("spec")
                target, mode = history[-1] if history else None, "continuation"
                stats["continuations"] += 1
                if target is None:
                    if not waived:
                        findings.append(
                            (number, "unresolved-continuation",
                             "bare `:%s` has no preceding full citation in this "
                             "document to resolve against" % spec))
                    continue
            stats["citations"] += 1
            prefix = None
            if not waived:
                prefix = comment_prefix(target)
                if prefix is None:
                    stats["uncommented"].add(target)
            assumed = mode == "guessed"
            for item in spec.split(","):
                low, _, high = item.partition("-")
                start, end = int(low), int(high or low)
                # Printed before the waiver short-circuit, because a CI reader
                # asking "which file did this line come from" wants the waived
                # citations most of all: they are the ones nothing else reports.
                if verbose:
                    stream.write("  %s:%d -> %s:%s [%s] label %s%s\n"
                                 % (document, number, target, item, mode, label,
                                    " waived" if waived else ""))
                if waived:
                    continue
                stats["targets"] += 1
                if not assumed:
                    labels.setdefault((target, start, end), {}).setdefault(
                        label, []).append(number)
                item_findings = check_item(cache, target, start, end, prefix)
                body, _ = cache.get(target)
                # Counted only for an item that passed. The footer prints these
                # under "not failures", and a range that failed is a finding,
                # not an instance of the weaker range rule being exercised.
                if (not item_findings and body and start < end
                        and end <= len(body)):
                    interior = body[start - 1:end]
                    if any(is_comment(text, prefix) for text in interior):
                        stats["comment_ranges"] += 1
                    if any(is_blank(text) for text in interior):
                        stats["blank_ranges"] += 1
                for kind, message in item_findings:
                    if assumed:
                        findings.append(
                            (number, "assumed-" + kind,
                             "%s, in the assumed file %s rather than in a file "
                             "this document cited (cited as `%s`)"
                             % (message, target, token)))
                    else:
                        findings.append((number, kind, "%s (cited as `%s`)"
                                         % (message, token)))
    for (target, start, end), by_label in sorted(labels.items()):
        if len(by_label) < 2:
            continue
        where = target + (":%d" % start if start == end
                          else ":%d-%d" % (start, end))
        findings.append((min(number for group in by_label.values()
                             for number in group), "duplicate-label",
                         "%s is cited %d times in this document with different "
                         "adjacent labels (%s). This is a partial mitigation for "
                         "a label/number transposition and does not prevent one."
                         % (where, len(by_label),
                            "; ".join("label %s on line%s %s"
                                      % (label, "" if len(numbers) == 1 else "s",
                                         ", ".join(str(n) for n in sorted(numbers)))
                                      for label, numbers in sorted(by_label.items())))))
    return findings, stats


def report(document, findings, stream):
    for number, kind, message in sorted(findings, key=lambda item: (item[0], item[1])):
        stream.write("%s:%d: %s: %s: %s\n"
                     % (document, number, severity(kind), kind, message))


def default_root():
    """Return the repository root, derived from this file's location."""
    return os.path.dirname(os.path.dirname(
        os.path.dirname(os.path.abspath(__file__))))


def main(argv=None):
    parser = argparse.ArgumentParser(
        prog="tools/citations/check.py",
        description="Validate the file:line citations in the status documents. "
                    "Reports only; it never writes to any file.",
        epilog=_HELP + "\n" + _LIMITS + "\n",
        formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("document", nargs="*",
                        help="documents to check (default: %s)"
                             % ", ".join(DEFAULT_DOCUMENTS))
    parser.add_argument("--root", default=None,
                        help="repository root that cited paths resolve against "
                             "(default: the root of this checkout)")
    parser.add_argument("-v", "--verbose", action="store_true",
                        help="print every citation as it resolves: the document "
                             "line, the file and line actually checked, which "
                             "rule resolved the file, and the adjacent label")
    args = parser.parse_args(argv)

    root = args.root or default_root()
    if not os.path.isdir(root):
        sys.stderr.write("error: --root %s is not a directory\n" % root)
        return EXIT_TOOL_ERROR
    documents = args.document or list(DEFAULT_DOCUMENTS)

    cache = FileCache(root)
    index = BasenameIndex(root)
    errors = warnings = 0
    totals = {"citations": 0, "full": 0, "continuations": 0, "targets": 0,
              "recalled": 0, "siblings": 0, "rooted": 0, "guessed": 0,
              "ambiguous": 0, "waived": 0, "comment_ranges": 0, "blank_ranges": 0}
    uncommented = set()
    try:
        for document in documents:
            if not os.path.isfile(document):
                # Tolerate being invoked by absolute path from another
                # directory: a default document is named relative to the
                # repository root, not to the current directory.
                candidate = os.path.join(root, *document.split("/"))
                if os.path.isfile(candidate):
                    document = candidate
                else:
                    sys.stderr.write("error: no such document: %s\n" % document)
                    return EXIT_TOOL_ERROR
            if posixpath.normpath(document).replace("\\", "/") in WAIVED_DOCUMENTS:
                sys.stderr.write("note: %s is waived (see WAIVED_DOCUMENTS in "
                                 "tools/citations/check.py); not checked\n"
                                 % document)
                continue
            findings, stats = check_document(document, root, cache, index,
                                            args.verbose, sys.stdout)
            report(document, findings, sys.stdout)
            errors += sum(1 for _, kind, _ in findings if severity(kind) == "error")
            warnings += sum(1 for _, kind, _ in findings
                            if severity(kind) == "warning")
            for key in totals:
                totals[key] += stats[key]
            uncommented |= stats["uncommented"]
    except ToolError as exc:
        sys.stderr.write("error: %s\n" % exc)
        return EXIT_TOOL_ERROR

    inferred = (totals["recalled"] + totals["siblings"] + totals["rooted"]
                + totals["guessed"])
    named = totals["full"] - inferred
    sys.stdout.write(
        "checked %d document%s: %d citation span%s (%d full, %d bare "
        "continuation%s), %d target line%s\n"
        % (len(documents), "" if len(documents) == 1 else "s",
           totals["citations"], "" if totals["citations"] == 1 else "s",
           totals["full"], totals["continuations"],
           "" if totals["continuations"] == 1 else "s",
           totals["targets"], "" if totals["targets"] == 1 else "s"))
    sys.stdout.write(
        "of the %d full citations, %d name their own path and %d name a bare "
        "filename, so this tool chose the file for them; run --verbose to see "
        "which file each one was checked against\n"
        % (totals["full"], named, inferred))
    sys.stdout.write(
        "bare filenames resolved: %d recalled from an earlier citation, %d "
        "beside the latest citation, %d at the repository root, %d guessed "
        "(a guess is a warning and cannot fail the gate); %d of them name a "
        "file the repository holds more than once, each reported as "
        "ambiguous-basename\n"
        % (totals["recalled"], totals["siblings"], totals["rooted"],
           totals["guessed"], totals["ambiguous"]))
    sys.stdout.write(
        "not failures: %d range citation%s that passed with an interior comment "
        "line, %d with an interior blank line, %d waived citation%s, %d target "
        "file%s with no comment rule\n"
        % (totals["comment_ranges"],
           "" if totals["comment_ranges"] == 1 else "s", totals["blank_ranges"],
           totals["waived"], "" if totals["waived"] == 1 else "s",
           len(uncommented), "" if len(uncommented) == 1 else "s"))
    sys.stdout.write("%d error%s, %d warning%s\n"
                     % (errors, "" if errors == 1 else "s",
                        warnings, "" if warnings == 1 else "s"))
    sys.stdout.write(_LIMITS + "\n")
    return EXIT_FINDINGS if errors else EXIT_OK


if __name__ == "__main__":
    if hasattr(sys.stdout, "reconfigure"):
        try:
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
        except (ValueError, OSError):
            pass
    sys.exit(main())

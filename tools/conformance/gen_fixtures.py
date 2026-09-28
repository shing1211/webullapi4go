#!/usr/bin/env python3
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

"""Emit the committed wire-conformance fixtures under ``conformance/testdata``.

Why this tool exists
--------------------
Most of the SDK's local-server tests serve a response body produced by
marshalling the SDK's own Go type and then assert the SDK decodes it, e.g.
``brokerfd/assets_test.go`` encodes ``[]FDPosition`` and asserts
``TestGetFDPositions`` reads it back. That proves the type is self-consistent.
It cannot prove the type matches the endpoint: the documented response for
``GET /broker/assets/positions/list`` requires ``cost_price``, ``last_price``
and ``unrealized_profit_loss`` while ``FDPosition`` tags ``average_cost``,
``market_value`` and ``unrealized_pl``. The self-round-trip passes and the live
endpoint returns three zero-valued ``money.Money`` with no error.

The cure is a response body that came from somewhere other than the SDK. This
tool builds one per documented endpoint out of Webull's own published OpenAPI
schema, so the fixture is a snapshot of the documentation and the SDK struct is
what has to agree with it.

Design rule: this tool emits data, never verdicts
-------------------------------------------------
No output of this tool is compared against a Go type, and no Go file is read.
A fixture derived from the SDK reproduces the bug it is meant to catch, so the
derivation must terminate in the documentation. The comparison lives in the Go
consumer under ``conformance/``, where a reviewer can read it against the SDK
types. See ``conformance/doc.go``.

What a fixture is
-----------------
A *minimal conforming instance* of the page's documented ``200`` schema: for
each object node, each required property; for each leaf, its inline
``example``; for each array, one element. Optional properties are omitted, which
is what makes the instance minimal. JSON-Schema-wise the result satisfies the
schema, so it is a legal body for the endpoint and a fixed target for a decoder.

The rule's cost is honest and recorded, not hidden: a page that declares no
``required`` list anywhere yields a structurally correct but nearly empty
instance. 101 of the 193 in-scope pages are in that position. The manifest marks
each one ``checks.declaresRequired = false`` and separately records the full
declared name inventory, so a consumer can still check SDK tag coverage for
names the minimal instance does not carry, and nobody can mistake "no names to
check" for "names verified".

Inputs and invariants
---------------------
The only input is the gitignored docgen cache, read with ``encoding="utf-8"``
because these pages carry emoji that the Windows cp1252 default codec cannot
decode. The cache is strictly read-only here: this tool never writes into it,
never re-fetches, and never runs ``docgen.py``. It reaches the cache format
through ``tools/webull-docgen/_common.py`` by import -- ``extract_json``,
``resolve_method_path``, ``parse_consts`` and the ``AREAS`` manifest -- rather
than re-implementing any of it.

One duplication is deliberate: ``_common.fetch`` derives a cache filename from a
URL, then creates the cache directory and re-fetches on a miss. Both effects are
forbidden here, so the three-line name derivation is repeated in
``cache_filename`` instead. A page whose name no longer resolves is reported as
a miss and fails ``--check``, so the duplication cannot drift unnoticed.

Usage
-----
    python tools/conformance/gen_fixtures.py --check     # gate: fail on drift
    python tools/conformance/gen_fixtures.py --write     # adopt new output
    python tools/conformance/gen_fixtures.py --self-test # emitter unit test only

``--self-test`` composes with either mode rather than replacing it, so the
Makefile targets can run the emitter's unit tests and the drift gate in one
invocation.

``--check`` and ``--write`` both exit 0 with a SKIP message when the cache is
absent, because the cache is gitignored and a checkout without it cannot judge
drift either way.
"""

import argparse
import hashlib
import json
import os
import re
import sys

_HERE = os.path.dirname(os.path.abspath(__file__))
_ROOT = os.path.abspath(os.path.join(_HERE, "..", ".."))
sys.path.insert(0, os.path.join(_ROOT, "tools", "webull-docgen"))

import _common as common  # noqa: E402  (needs the path above)

OUT_DIR = os.path.join(_ROOT, "conformance", "testdata")
MANIFEST_NAME = "manifest.json"

# A fixture this large is still readable in a diff; a fixture much larger than
# this is not a schema-derived instance any more. A future page that inlines a
# base64 blob as an `example` would blow past both, and the committed tree would
# become unreviewable, so the ceiling fails the run instead of warning.
LOG_ABOVE_BYTES = 8 * 1024
FAIL_ABOVE_BYTES = 32 * 1024

# Key used when a required free-form object has to be given a shape. Plainly
# synthetic so it cannot be mistaken for a documented name.
SYNTHETIC_KEY = "example_key"

# Recursion guard. The deepest schema in the cache nests a handful of levels;
# this only exists so a pathological page cannot exhaust the stack.
MAX_DEPTH = 16

# A JSON Pointer segment escapes these two characters; ``application/json``
# therefore becomes ``application~1json``.
_PT_ESCAPE = re.compile(r"~(?![01])|/")


def pointer(*segments):
    """Join ``segments`` into a JSON Pointer into the schema node."""
    out = []
    for seg in segments:
        if isinstance(seg, int):
            out.append(str(seg))
        else:
            out.append(_PT_ESCAPE.sub(lambda m: "~1" if m.group() == "/" else "~0", str(seg)))
    return "/" + "/".join(out)


# --------------------------------------------------------------------------
# Reading the cache (read-only)
# --------------------------------------------------------------------------
def cache_filename(url):
    """Return the cache filename ``_common.fetch`` would use for ``url``.

    Deliberately a copy of three lines of ``_common.fetch`` rather than a call
    to it. ``fetch`` calls ``os.makedirs`` on the cache and, on a miss, performs
    a network request; this tool may do neither. See the module docstring.
    """
    return re.sub(r"[^A-Za-z0-9]+", "_", url)[-120:] + ".md"


def cache_is_populated():
    """True when the cache directory holds at least one ``.md`` page."""
    cache = common.cache_dir()
    if not os.path.isdir(cache):
        return False
    return any(f.endswith(".md") for f in os.listdir(cache))


def read_page(url):
    """Return the cached markdown for ``url``, or ``None`` when not cached.

    ``encoding="utf-8"`` is mandatory, not a default: several pages contain
    emoji and the Windows cp1252 default codec raises on them.
    """
    path = os.path.join(common.cache_dir(), cache_filename(url))
    if not os.path.exists(path) or os.path.getsize(path) == 0:
        return None
    with open(path, "r", encoding="utf-8") as fh:
        return fh.read()


def _json_blocks(md):
    """Yield ``(index, parsed)`` for every ```json block in ``md``."""
    for i, raw in enumerate(re.findall(r"```json\s*\n(.*?)\n```", md, re.S)):
        try:
            yield i, json.loads(raw)
        except ValueError:
            continue


def select_block(md):
    """Return ``(index, document)`` for the OpenAPI document in ``md``.

    ``_common.extract_json`` is the shared definition of "the JSON block in this
    page" and is tried first, so a page whose first block is the endpoint
    document is read exactly as the doc generator reads it. Only when that block
    lacks both ``path`` and ``responses`` do we look further, preferring a block
    that has both. Selecting costs nothing and keeps a multi-block guide page
    from supplying the wrong document.
    """
    first = common.extract_json(md)
    if isinstance(first, dict) and "path" in first and "responses" in first:
        return 0, first
    for i, doc in _json_blocks(md):
        if isinstance(doc, dict) and "path" in doc and "responses" in doc:
            return i, doc
    if isinstance(first, dict):
        return 0, first
    return None, None


def schema_200(document):
    """Return the documented ``200`` ``application/json`` schema, or ``None``."""
    responses = document.get("responses") or {}
    ok = responses.get("200") or responses.get(200) or {}
    for mime, body in (ok.get("content") or {}).items():
        if "json" in mime and isinstance(body, dict):
            return body.get("schema")
    return None


SCHEMA_POINTER = pointer("responses", "200", "content", "application/json", "schema")


# --------------------------------------------------------------------------
# Emitting a minimal conforming instance
# --------------------------------------------------------------------------
def ordered_required(properties, required):
    """Required names in documented property order.

    Webull lists ``properties`` in page order and ``required`` alphabetically.
    Emitting in ``properties`` order makes a fixture line up with the page's own
    property table, so a reviewer can read the two side by side. The order is a
    property of the page, not of this tool, so it is stable across runs.
    """
    wanted = set(required or ())
    return [name for name in (properties or {}) if name in wanted]


def leaf_value(node):
    """Return ``(value, synthesized)`` for a leaf ``node``.

    ``synthesized`` is True when the value did not come from an inline
    ``example`` and had to be invented from the declared type. The caller records
    every such leaf in the manifest, so a reader can tell a documented example
    from a stand-in.
    """
    if "example" in node:
        return node["example"], False
    kind = node.get("type")
    if kind == "string":
        return "", True
    if kind == "integer":
        return 0, True
    if kind == "number":
        return 0, True
    if kind == "boolean":
        return False, True
    if kind == "null":
        return None, False
    return None, True


def discriminator_of(properties):
    """Return ``(name, value)`` for the discriminator a sibling ``type`` implies.

    The two ``oneOf`` pages in the cache both document that a sibling ``type``
    selects the branch. Reading that example lets the branch choice be derived
    from the page instead of asserted, which is why it is captured here and
    handed down the recursion.
    """
    node = (properties or {}).get("type")
    if not isinstance(node, dict):
        return None, None
    if "example" in node:
        return "type", node["example"]
    enum = node.get("enum")
    if isinstance(enum, list) and enum:
        return "type", enum[0]
    return "type", None


def choose_one_of(branch_nodes, disc):
    """Choose one ``oneOf`` branch. Returns ``(index, title, how)``.

    The choice is explicit and recorded, never implicit. When the containing
    object carries a discriminating ``type`` example, the branch whose ``title``
    names that variant wins, and the record says the choice was made that way.
    Otherwise the first branch is taken and the record says *that*, so a reader
    can see no branch was selected on evidence.

    Only the head token of the example is matched, e.g. ``basketball`` out of
    ``basketball_game``. That is what the pages themselves do: the
    ``event-contracts/live-data/get`` description tabulates
    ``basketball_game -> EventContractLiveBasketballDetailsVo`` and each branch
    title carries the same head word. Matching every token would instead pair
    ``football_game`` with ``SportsGameDetails`` on the shared word ``game``,
    which is a coincidence and would read as evidence.
    """
    if isinstance(disc, (str, int, float, bool)):
        head = re.split(r"[^A-Za-z0-9]+", str(disc).lower())[0]
        if len(head) > 2:
            for i, branch in enumerate(branch_nodes):
                title = (branch.get("title") or "").lower() if isinstance(branch, dict) else ""
                if title and head in title:
                    return i, branch.get("title"), "title-matches-discriminator"
    first = branch_nodes[0] if branch_nodes else {}
    return 0, first.get("title") if isinstance(first, dict) else None, "first-branch"


class Emitter(object):
    """Builds one minimal conforming instance and records how it got there."""

    def __init__(self):
        self.notes = []          # oneOf / additionalProperties / synthesized leaves
        self.emitted = set()     # pointers the instance actually reached
        self.unknown_required = []  # required names with no declared property

    # -- notes ---------------------------------------------------------
    def _note(self, kind, ptr, emitted, **fields):
        """Append one record. ``emitted`` is passed in, never inferred.

        Inferring it from a set of visited pointers looks tempting and is wrong:
        a ``oneOf`` note names the ``oneOf`` keyword while the traversal marks
        the node that holds it, so the two pointers never coincide and every
        record reads ``emitted: false``.
        """
        record = {"kind": kind, "pointer": ptr, "emitted": emitted}
        record.update(fields)
        self.notes.append(record)

    # -- traversal -----------------------------------------------------
    def node(self, schema, ptr, disc, depth=0):
        """Return the minimal instance for ``schema``, or ``None`` to omit it."""
        if not isinstance(schema, dict) or depth > MAX_DEPTH:
            return None
        self.emitted.add(ptr)
        if "oneOf" in schema and isinstance(schema["oneOf"], list) and schema["oneOf"]:
            idx, title, how = choose_one_of(schema["oneOf"], disc[1] if disc else None)
            branch = schema["oneOf"][idx] or {}
            self._note("oneOf", ptr + pointer("oneOf"), True,
                       branchIndex=idx, title=title, selection=how,
                       branchCount=len(schema["oneOf"]),
                       branchRequiredNameCount=len(branch.get("required") or []),
                       discriminator=None if not disc else disc[0],
                       discriminatorValue=None if not disc else disc[1])
            return self.node(branch, ptr + pointer("oneOf", idx), disc, depth + 1)
        kind = schema.get("type")
        if kind == "array" or (kind is None and "items" in schema):
            element = self.node(schema.get("items") or {}, ptr + pointer("items"),
                                disc, depth + 1)
            return [element]
        if kind == "object" or (kind is None and "properties" in schema):
            return self.obj(schema, ptr, depth)
        value, synthesized = leaf_value(schema)
        if synthesized:
            self._note("synthesizedLeaf", ptr, True, type=kind or "unspecified")
        return value

    def obj(self, schema, ptr, depth):
        """Emit an object node: its required properties, nothing else."""
        properties = schema.get("properties") or {}
        required = schema.get("required") or []
        inner_disc = discriminator_of(properties)
        out = {}
        for name in ordered_required(properties, required):
            out[name] = self.node(properties[name], ptr + pointer("properties", name),
                                  inner_disc, depth + 1)
        for name in required:
            if name not in properties:
                self.unknown_required.append(ptr + pointer("required", name))
        extra = schema.get("additionalProperties")
        if not out and isinstance(extra, dict):
            # A free-form object that the minimal rule would leave empty. One
            # synthetic key keeps the declared value type exercised instead of
            # collapsing every map-shaped response to "{}".
            out[SYNTHETIC_KEY] = self.node(
                extra, ptr + pointer("additionalProperties", SYNTHETIC_KEY),
                inner_disc, depth + 1)
            self._note("additionalProperties", ptr + pointer("additionalProperties"),
                       True, valueType=extra.get("type"), syntheticKey=SYNTHETIC_KEY)
        return out

    def walk_all(self, schema, ptr, disc=None, depth=0):
        """Record every ``oneOf`` and typed ``additionalProperties`` in the tree.

        A separate pass from :meth:`node` on purpose. The minimal rule omits
        optional properties, so a ``oneOf`` or ``additionalProperties`` living
        under one is never reached by the emitted instance -- but the decision
        to pick a branch still has to be made and written down rather than left
        to the first consumer that happens to want the branch. The ``emitted``
        flag on each record says which of the two happened, so a reader never
        has to guess whether a recorded choice was exercised.
        """
        if not isinstance(schema, dict) or depth > MAX_DEPTH:
            return
        if isinstance(schema.get("oneOf"), list) and schema["oneOf"]:
            idx, title, how = choose_one_of(schema["oneOf"], disc[1] if disc else None)
            branch = schema["oneOf"][idx] or {}
            self._note("oneOf", ptr + pointer("oneOf"), ptr in self.emitted,
                       branchIndex=idx, title=title, selection=how,
                       branchCount=len(schema["oneOf"]),
                       branchRequiredNameCount=len(branch.get("required") or []),
                       discriminator=None if not disc else disc[0],
                       discriminatorValue=None if not disc else disc[1])
        extra = schema.get("additionalProperties")
        if isinstance(extra, dict):
            self._note("additionalProperties", ptr + pointer("additionalProperties"),
                       ptr in self.emitted, valueType=extra.get("type"),
                       syntheticKey=SYNTHETIC_KEY)
        properties = schema.get("properties") or {}
        inner = discriminator_of(properties)
        for name, sub in properties.items():
            self.walk_all(sub, ptr + pointer("properties", name), inner, depth + 1)
        if "items" in schema:
            self.walk_all(schema.get("items") or {}, ptr + pointer("items"), disc, depth + 1)
        for combinator in ("oneOf", "allOf", "anyOf"):
            for i, sub in enumerate(schema.get(combinator) or []):
                self.walk_all(sub, ptr + pointer(combinator, i), disc, depth + 1)
        if isinstance(extra, dict):
            self.walk_all(extra, ptr + pointer("additionalProperties"), disc, depth + 1)

    # -- description ---------------------------------------------------
    def declared_names(self, schema):
        """Names declared at the two levels a Go consumer can compare directly.

        The top-level object, or the array element when the top level is an
        array. Deeper names are reachable from the fixture itself, so the Go
        side reads them from the committed bytes; only these two levels need to
        be restated for a page whose minimal instance is empty.
        """
        if not isinstance(schema, dict):
            return [], []
        if schema.get("type") == "array" or (schema.get("type") is None and "items" in schema):
            items = schema.get("items") or {}
            return [], list((items.get("properties") or {}).keys())
        return list((schema.get("properties") or {}).keys()), []

    def all_property_names(self, schema, acc=None):
        """Every property name the schema declares at any depth, deduplicated."""
        if acc is None:
            acc = []
        if not isinstance(schema, dict):
            return acc
        for name, sub in (schema.get("properties") or {}).items():
            if name not in acc:
                acc.append(name)
            self.all_property_names(sub, acc)
        items = schema.get("items")
        if isinstance(items, dict):
            self.all_property_names(items, acc)
        for combinator in ("oneOf", "allOf", "anyOf"):
            for sub in schema.get(combinator) or []:
                self.all_property_names(sub, acc)
        extra = schema.get("additionalProperties")
        if isinstance(extra, dict):
            self.all_property_names(extra, acc)
        return acc


# --------------------------------------------------------------------------
# Building one record per in-scope endpoint
# --------------------------------------------------------------------------
def slugify(path):
    return re.sub(r"[^a-z0-9]+", "-", str(path).lower()).strip("-")


def property_names_in(instance, acc=None):
    """Every object key present in the emitted instance, at any depth."""
    if acc is None:
        acc = set()
    if isinstance(instance, dict):
        for key, val in instance.items():
            acc.add(key)
            property_names_in(val, acc)
    elif isinstance(instance, list):
        for item in instance:
            property_names_in(item, acc)
    return acc


def latest_notes(notes):
    """Collapse the notes of a two-pass walk, keeping the last of each node.

    ``Emitter.node`` and ``Emitter.walk_all`` both record a ``oneOf`` or
    ``additionalProperties`` decision, in that order, and only the second pass
    knows whether the minimal instance actually reached the node. Keeping the
    last record per ``(kind, pointer)`` therefore keeps the accurate one.
    """
    latest = {}
    for note in notes:
        latest[(note["kind"], note["pointer"])] = note
    return list(latest.values())


def build_records(consts):
    """Return ``(records, skipped)`` for every manifest entry with a 200 schema.

    An entry is in scope when its cached page carries a documented ``200``
    ``application/json`` schema. Pages without one (guides, gRPC subscription
    references) are reported as skipped rather than silently dropped, so a
    reviewer can see what was not covered.
    """
    records = []
    skipped = []
    for area, (title, _blurb, entries) in common.AREAS.items():
        for label, url, sdk, note in entries:
            page = cache_filename(url)
            md = read_page(url)
            if md is None:
                skipped.append({"area": area, "label": label, "page": page,
                                "reason": "page not in cache"})
                continue
            block, document = select_block(md)
            if not isinstance(document, dict):
                skipped.append({"area": area, "label": label, "page": page,
                                "reason": "no JSON block parses"})
                continue
            schema = schema_200(document)
            if schema is None:
                skipped.append({"area": area, "label": label, "page": page,
                                "reason": "no 200 application/json schema"})
                continue
            sdk_path, path_const = common.resolve_method_path(sdk, consts)
            method = str(document.get("method") or "get").upper()
            doc_path = document.get("path") or ""
            fid = "%s/%s-%s" % (area, method, slugify(doc_path))

            emitter = Emitter()
            instance = emitter.node(schema, "", None)
            emitter.walk_all(schema, "")
            top_names, element_names = emitter.declared_names(schema)
            top_type = schema.get("type") or ("array" if "items" in schema else "object")
            element_type = None
            if top_type == "array":
                items = schema.get("items") or {}
                element_type = items.get("type") or ("array" if "items" in items else "object")
            required = []
            required_source = None
            if top_type == "array":
                required = (schema.get("items") or {}).get("required") or []
                required_source = "items" if required else None
            else:
                required = schema.get("required") or []
                required_source = "topLevel" if required else None
            notes = [n for n in latest_notes(emitter.notes)
                     if n["kind"] in ("oneOf", "additionalProperties")]
            fixture = render_fixture(instance)
            record = {
                "id": fid,
                "fixture": fid + ".json",
                "area": area,
                "pageTitle": label,
                "sdkSymbol": sdk,
                "source": {
                    "url": url,
                    "cacheFile": page,
                    "jsonBlockIndex": block,
                    "schemaPointer": SCHEMA_POINTER,
                },
                "documented": {"method": method, "path": doc_path},
                "sdkPath": sdk_path,
                "sdkPathConst": path_const,
                "sdkPathMatch": ("noSdkPath" if not sdk_path
                                 else "same" if sdk_path == doc_path else "differs"),
                "checks": {
                    "topLevel": top_type,
                    "elementType": element_type,
                    "declaresRequired": bool(required),
                    "requiredNamesSource": required_source,
                    "requiredNames": sorted(required),
                    "declaredTopLevelNames": top_names,
                    "declaredElementNames": element_names,
                    "declaredPropertyNameCount": len(emitter.all_property_names(schema)),
                    "propertyNameCountInFixture": len(property_names_in(instance)),
                    "synthesizedLeafCount": sum(1 for n in emitter.notes
                                                if n["kind"] == "synthesizedLeaf"),
                    "requiredWithoutExample": sorted(emitter.unknown_required),
                },
                "oneOfChoices": [n for n in notes if n["kind"] == "oneOf"],
                "additionalProperties": [n for n in notes
                                        if n["kind"] == "additionalProperties"],
                "bytes": len(fixture.encode("utf-8")),
                "sha256": hashlib.sha256(
                    fixture.replace("\r\n", "\n").encode("utf-8")).hexdigest(),
                # Private build input, popped before the manifest is rendered.
                "_fixture": fixture,
            }
            records.append(record)
    return records, skipped


def render_fixture(instance):
    """Serialise ``instance`` as compact JSON with a trailing newline."""
    return json.dumps(instance, separators=(",", ":"), ensure_ascii=False) + "\n"


def render_manifest(records, skipped, oversized):
    """Assemble ``manifest.json``.

    Deliberately free of timestamps, absolute paths and hostnames: the ``--check``
    gate compares these bytes against the committed file, so any run-varying
    value would report drift on every run and destroy the signal.
    """
    total = sum(r["bytes"] for r in records)
    no_required = sorted(r["id"] for r in records
                         if not r["checks"]["declaresRequired"])
    return {
        "generator": {
            "tool": "tools/conformance/gen_fixtures.py",
            "command": "make conformance-fixtures",
            "source": "tools/webull-docgen/.cache (gitignored, read-only)",
            "derivedFrom": "the 200 application/json schema of each page",
            "fixtureBytesComeFromGo": False,
            "goAccess": ("resolve_method_path reads Go path constants, so the run "
                         "does open .go files. The audit in this tool reports every "
                         "one it opens and fails on any outside that path-const "
                         "scan. Those constants annotate sdkPath and sdkPathMatch "
                         "below and never reach a fixture: every fixture byte is a "
                         "function of the documented schema alone."),
            "minimalInstanceRule": ("for each object node emit its required "
                                    "properties; for each leaf emit its inline "
                                    "example; for each array emit one element; "
                                    "omit optional properties"),
        },
        "sizeTripwire": {
            "logAboveBytes": LOG_ABOVE_BYTES,
            "failAboveBytes": FAIL_ABOVE_BYTES,
            "reason": ("a fixture far larger than the mean is not a "
                       "schema-derived instance any more; a page inlining a "
                       "base64 blob as an example would put an unreviewable "
                       "blob in the committed tree"),
        },
        "totals": {
            "fixtures": len(records),
            "fixtureBytes": total,
            "pagesWithoutRequired": len(no_required),
            "pagesWithOneOf": sum(1 for r in records if r["oneOfChoices"]),
            "pagesWithAdditionalProperties": sum(1 for r in records
                                                 if r["additionalProperties"]),
            "oversizedFixtures": [r["id"] for r in oversized],
        },
        "pagesWithoutRequiredNames": no_required,
        "notes": {
            "declaresRequired": ("false means the page declares no required list "
                                 "at the top level or on the array element, so "
                                 "the minimal instance carries no names from "
                                 "Webull. It is NOT a statement that the names "
                                 "were verified; use declaredTopLevelNames / "
                                 "declaredElementNames for that."),
            "emitted": ("on a oneOfChoices or additionalProperties record, "
                        "whether the minimal instance actually reached that "
                        "node. false means the node sits under an optional "
                        "property, which the minimal rule omits."),
            "branchRequiredNameCount": ("required names the chosen oneOf branch "
                                        "declares. 0 means the recorded branch "
                                        "choice cannot be observed in the "
                                        "instance, because the minimal rule "
                                        "emits only required names."),
            "sdkPath": ("resolve_method_path's answer for the manifest's SDK "
                        "symbol. sdkPathMatch compares it against "
                        "documented.path so a path defect is visible here; "
                        "see IMPLEMENTATION_STATUS.md."),
            "sdkPathMatch": ("same when the SDK sends the documented path, "
                             "differs when it sends another, noSdkPath when "
                             "the entry is deliberately mapped to no SDK "
                             "symbol (the AREAS manifest's em dash)."),
        },
        "skippedPages": skipped,
        "fixtures": records,
    }


# --------------------------------------------------------------------------
# Output: write, or gate on drift
# --------------------------------------------------------------------------
def out_path(name):
    """Resolve an ``OUT_DIR``-relative, forward-slash name to a native path."""
    return os.path.normpath(os.path.join(OUT_DIR, *name.split("/")))


def committed_files():
    """Map every file under ``conformance/testdata`` to its text.

    Keys are forward-slashed and relative to ``OUT_DIR``, the same key space
    ``record["fixture"]`` uses. Keying on the raw path instead made ``--check``
    compare 193 planned paths against 0 committed ones on Windows and report
    every file as added, which is indistinguishable from real drift.
    """
    found = {}
    for root, _dirs, files in os.walk(OUT_DIR):
        for name in files:
            path = os.path.join(root, name)
            with open(path, "r", encoding="utf-8") as fh:
                found[os.path.relpath(path, OUT_DIR).replace(os.sep, "/")] = fh.read()
    return found


def relative(path):
    return os.path.relpath(path, _ROOT).replace(os.sep, "/")


# --------------------------------------------------------------------------
# Go-access audit
# --------------------------------------------------------------------------
# `resolve_method_path` -- which the manifest is required to reuse -- works by
# reading Go path constants: `_common.parse_consts` scans every non-test .go
# file under the SDK package directories, and `_common.find_method` reads the
# method bodies. So this tool's run *does* open Go files, and any claim
# otherwise would be false.
#
# What matters is narrower and is what the audit below enforces: no Go file
# outside that path-constant scan is read, and no Go-derived value reaches a
# fixture. Every fixture byte is a function of the documented schema alone;
# `sdkPath` is the only Go-derived value in the output and it appears solely as
# manifest annotation. `parse_consts` reads path strings, never struct fields,
# so it cannot inform a fixture's property names even in principle.
_GO_ROOTS = tuple(os.path.normpath(os.path.join(common.ROOT, d))
                  for d in common.GO_DIRS)


class GoAccessAudit(object):
    """Records every ``.go`` file the run opens and flags anything out of scope.

    Installed through :func:`sys.addaudithook`, which cannot be uninstalled, so
    it stays disarmed until :meth:`arm` is called. A permanent hook that fired
    during import would tax every other use of the interpreter.
    """

    def __init__(self):
        self.armed = False
        self.opened = []
        self.violations = []

    def arm(self):
        self.armed = True

    def disarm(self):
        self.armed = False

    def hook(self, event, args):
        if not self.armed or event != "open" or not args:
            return
        target = args[0]
        if isinstance(target, bytes):
            target = target.decode("utf-8", "replace")
        if not isinstance(target, str) or not target.endswith(".go"):
            return
        path = os.path.normpath(target)
        self.opened.append(path)
        if path.endswith("_test.go") or not path.startswith(_GO_ROOTS):
            self.violations.append(path)

    def report(self):
        """Print what the run read, and fail if it read out of scope."""
        # relative() already normalises to forward slashes, so split on "/" and
        # not on os.sep: on Windows the latter finds no separator and prints one
        # 2,000-element "directory". Only the first segment is wanted -- the SDK
        # package -- not the file within it.
        packages = sorted({relative(p).split("/")[0] for p in self.opened})
        print("go files read: %d opens across %d distinct files in: %s"
              % (len(self.opened), len(set(self.opened)),
                 ", ".join(packages) if packages else "-"))
        print("  path constants only, via _common.parse_consts; no fixture byte "
              "comes from Go")
        if self.violations:
            print("FAIL: read %d .go file(s) outside the path-const scan: %s"
                  % (len(self.violations), ", ".join(sorted(set(self.violations)))))
            return 1
        return 0


GO_AUDIT = GoAccessAudit()
sys.addaudithook(GO_AUDIT.hook)


def write(planned):
    for name, data in planned.items():
        path = out_path(name)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w", encoding="utf-8", newline="\n") as fh:
            fh.write(data)
    for name in committed_files():
        if name not in planned:
            os.remove(out_path(name))


def check(planned):
    """Return a list of drift descriptions; empty means the tree is current."""
    current = committed_files()
    root = relative(OUT_DIR)
    drift = []
    for name, data in sorted(planned.items()):
        if name not in current:
            drift.append("added    %s" % "/".join((root, name)))
        elif current[name] != data:
            drift.append("changed  %s" % "/".join((root, name)))
    for name in sorted(current):
        if name not in planned:
            drift.append("removed  %s" % "/".join((root, name)))
    return drift


def report_sizes(records, oversized):
    """Print the size distribution against the two thresholds."""
    sizes = sorted(r["bytes"] for r in records)
    n = len(sizes)
    print("fixtures: %d   total bytes: %d   mean %d   median %d   max %d"
          % (n, sum(sizes), sum(sizes) // n, sizes[n // 2], sizes[-1]))
    buckets = [("<= 1 KiB", lambda s: s <= 1024),
               ("1-8 KiB", lambda s: 1024 < s <= LOG_ABOVE_BYTES),
               ("8-32 KiB (logged)", lambda s: LOG_ABOVE_BYTES < s <= FAIL_ABOVE_BYTES),
               ("> 32 KiB (fatal)", lambda s: s > FAIL_ABOVE_BYTES)]
    for label, pred in buckets:
        print("  %-22s %3d" % (label, sum(1 for s in sizes if pred(s))))
    for record in records:
        if record["bytes"] > LOG_ABOVE_BYTES:
            print("LOG: %s is %d bytes, above the %d byte review threshold"
                  % (record["id"], record["bytes"], LOG_ABOVE_BYTES))
    if oversized:
        print("FAIL: %d fixture(s) exceed the %d byte ceiling: %s"
              % (len(oversized), FAIL_ABOVE_BYTES,
                 ", ".join("%s (%d B)" % (r["id"], r["bytes"]) for r in oversized)))
        print("      %s" % manifest_reason())


def manifest_reason():
    return ("a fixture above this size is not a schema-derived minimal instance "
            "any more; a page inlining a base64 blob as an `example` would put "
            "an unreviewable blob into the committed tree")


# --------------------------------------------------------------------------
# Emitter self-test
# --------------------------------------------------------------------------
SELFTEST_SCHEMA = {
    "type": "object",
    "required": ["plain", "no_example", "flag", "rows", "variant", "freeform"],
    "properties": {
        "plain": {"type": "string", "example": "AAPL"},
        "no_example": {"type": "integer"},
        "flag": {"type": "boolean", "example": True},
        "rows": {"type": "array", "items": {"type": "string", "example": "row"}},
        "variant": {
            "type": "object",
            "oneOf": [
                {"type": "object", "title": "BasketballDetails",
                 "properties": {"home_points": {"type": "integer", "example": 92}}},
                {"type": "object", "title": "HockeyDetails",
                 "properties": {"goals": {"type": "integer", "example": 3}}},
            ],
        },
        "freeform": {
            "type": "object",
            "additionalProperties": {
                "type": "array",
                "items": {"type": "object", "required": ["value"],
                          "properties": {"value": {"type": "string", "example": "0.1"}}},
            },
        },
    },
}

# Neither branch of `variant` declares a required list, which is the case for
# both `oneOf` pages in the real cache: the branch choice is recorded, and the
# instance is "{}" whichever branch was picked.
SELFTEST_EXPECTED = (
    '{"plain":"AAPL","no_example":0,"flag":true,"rows":["row"],'
    '"variant":{},"freeform":{"example_key":[{"value":"0.1"}]}}'
)

VARIANT = SELFTEST_SCHEMA["properties"]["variant"]


def discriminated(discriminator_example, branches):
    """A schema whose only required property is a ``oneOf`` chosen by ``type``."""
    return {
        "type": "object",
        "required": ["variant"],
        "properties": {
            "type": {"type": "string", "example": discriminator_example},
            "variant": {"type": "object", "oneOf": branches},
        },
    }


REQUIRED_BRANCHES = [
    {"type": "object", "title": "BasketballDetails", "required": ["home_points"],
     "properties": {"home_points": {"type": "integer", "example": 92}}},
    {"type": "object", "title": "HockeyDetails", "required": ["goals"],
     "properties": {"goals": {"type": "integer", "example": 3}}},
]


def self_test():
    """Exercise the emitter on a synthetic schema covering every branch.

    The real cache never reaches the ``additionalProperties`` path -- all six of
    its occurrences sit under an optional property, which the minimal rule
    omits -- so this is the only place that code is actually executed. It also
    pins the ``oneOf`` rules the cache cannot exercise, since it never visits a
    branch: that a discriminator picks the *named* branch rather than always
    branch 0, and that a branch is only observable when it declares required
    names.
    """
    failures = []
    ran = []

    def check(label, got, want):
        ran.append(label)
        if got != want:
            failures.append("%s\n    got  %r\n    want %r" % (label, got, want))

    def emit(schema):
        emitter = Emitter()
        instance = emitter.node(schema, "", None)
        emitter.walk_all(schema, "")
        return emitter, instance

    # -- the shapes the cache does contain ----------------------------
    emitter, instance = emit(SELFTEST_SCHEMA)
    check("minimal instance", json.dumps(instance, separators=(",", ":")),
          SELFTEST_EXPECTED)
    check("synthesized leaves", [n["pointer"] for n in emitter.notes
                                  if n["kind"] == "synthesizedLeaf"],
          ["/properties/no_example"])
    variant_note = [n for n in latest_notes(emitter.notes)
                    if n["kind"] == "oneOf"][0]
    check("oneOf without a sibling type", (variant_note["branchIndex"],
                                            variant_note["selection"]),
          (0, "first-branch"))
    check("oneOf branch declares no required name",
          variant_note["branchRequiredNameCount"], 0)
    check("oneOf emitted", variant_note["emitted"], True)
    free = [n for n in latest_notes(emitter.notes)
            if n["kind"] == "additionalProperties"][0]
    check("additionalProperties record", (free["valueType"], free["syntheticKey"],
                                          free["emitted"]),
          ("array", SYNTHETIC_KEY, True))

    # -- the shapes the cache does not reach --------------------------
    emitter, instance = emit(discriminated("basketball_game", REQUIRED_BRANCHES))
    check("discriminator picks branch 0", json.dumps(instance, separators=(",", ":")),
          '{"variant":{"home_points":92}}')
    note = [n for n in latest_notes(emitter.notes) if n["kind"] == "oneOf"][0]
    check("branch 0 chosen on evidence", (note["branchIndex"], note["title"],
                                          note["selection"],
                                          note["discriminator"],
                                          note["discriminatorValue"]),
          (0, "BasketballDetails", "title-matches-discriminator", "type",
           "basketball_game"))

    # The decisive case: a discriminator that names the *second* branch. If the
    # rule silently degraded to "always branch 0" this would still pass the case
    # above, so it is asserted separately.
    emitter, instance = emit(discriminated("hockey_match", REQUIRED_BRANCHES))
    check("discriminator picks branch 1", json.dumps(instance, separators=(",", ":")),
          '{"variant":{"goals":3}}')
    note = [n for n in latest_notes(emitter.notes) if n["kind"] == "oneOf"][0]
    check("branch 1 chosen on evidence", (note["branchIndex"], note["title"]),
          (1, "HockeyDetails"))

    # A discriminator no branch title mentions must fall back honestly rather
    # than pair on a shared word, the way "football_game" would against
    # "SportsGameDetails".
    emitter, _ = emit(discriminated("football_game", [
        {"type": "object", "title": "SportsGameDetails",
         "properties": {"venue": {"type": "string", "example": "Amon G. Carter"}}},
        {"type": "object", "title": "EconomicReleaseDetails",
         "properties": {"source": {"type": "string", "example": "BLS"}}},
    ]))
    note = [n for n in latest_notes(emitter.notes) if n["kind"] == "oneOf"][0]
    check("unmatched discriminator falls back", (note["branchIndex"],
                                                 note["selection"]),
          (0, "first-branch"))

    # A typed additionalProperties on a *required* free-form object is the path
    # the cache never reaches; assert the synthetic key and the value shape.
    emitter, instance = emit({
        "type": "object", "required": ["values"],
        "properties": {"values": {
            "type": "object",
            "additionalProperties": {
                "type": "array",
                "items": {"type": "object", "required": ["fiscal_year"],
                          "properties": {"fiscal_year": {"type": "integer",
                                                         "example": 2025}}}}}},
    })
    check("required free-form object", json.dumps(instance, separators=(",", ":")),
          '{"values":{"example_key":[{"fiscal_year":2025}]}}')

    # A pointer escapes a "/" inside a name so it stays a valid JSON Pointer.
    check("pointer escaping",
          pointer("responses", "200", "content", "application/json"),
          "/responses/200/content/application~1json")

    if failures:
        print("self-test FAILED (%d of %d checks):" % (len(failures), len(ran)))
        for f in failures:
            print("  " + f)
        return 1
    # Counted from what actually ran, not written out by hand. A hardcoded total
    # is a claim about this function that no run verifies, and it was wrong once:
    # it said 12 while self_test held 13 checks, so the tool under-reported its
    # own coverage in the one line a reader looks at.
    print("self-test: all %d emitter checks passed" % len(ran))
    return 0


# --------------------------------------------------------------------------
def run(mode):
    GO_AUDIT.arm()
    try:
        consts = common.parse_consts()
        records, skipped = build_records(consts)
    finally:
        GO_AUDIT.disarm()
    rc = GO_AUDIT.report()
    if rc != 0:
        return rc

    oversized = [r for r in records if r["bytes"] > FAIL_ABOVE_BYTES]
    report_sizes(records, oversized)

    # Pop the fixture text before the manifest is rendered. Leaving it in place
    # would serialise every fixture a second time inside its own record, and the
    # manifest would then be fifteen times the size of the tree it describes.
    planned = {}
    for record in records:
        planned[record["fixture"]] = record.pop("_fixture")
    manifest = render_manifest(records, skipped, oversized)
    planned[MANIFEST_NAME] = json.dumps(manifest, indent=1, ensure_ascii=False) + "\n"

    print("manifest bytes: %d" % len(planned[MANIFEST_NAME].encode("utf-8")))
    print("pages skipped (no 200 JSON schema or not cached): %d" % len(skipped))

    if oversized:
        return 1

    if mode == "write":
        write(planned)
        print("wrote %d file(s) under %s" % (len(planned), relative(OUT_DIR)))
        return 0
    drift = check(planned)
    if drift:
        print("DRIFT: %d file(s) differ from the committed tree:" % len(drift))
        for line in drift:
            print("  " + line)
        print("Regenerate and review with: make conformance-fixtures-update")
        return 1
    print("conformance fixtures match the committed tree (%d file(s))" % len(planned))
    return 0


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("--check", action="store_true",
                        help="fail when the regenerated tree differs (default)")
    parser.add_argument("--write", action="store_true",
                        help="overwrite the committed tree")
    parser.add_argument("--self-test", action="store_true",
                        help="also unit-test the emitter")
    args = parser.parse_args(argv)

    if args.check and args.write:
        parser.error("--check and --write are mutually exclusive")

    # The self-test runs in addition to the mode, never instead of it. Both
    # Makefile targets pass it, and returning early here would leave the drift
    # gate silently unrun while still exiting 0.
    if args.self_test:
        rc = self_test()
        if rc != 0:
            return rc

    if not cache_is_populated():
        print("SKIP: no docgen cache at %s" % common.cache_dir())
        print("      The cache is gitignored, so a checkout without it cannot "
              "regenerate or judge the fixtures.")
        print("      Populate it with tools/webull-docgen (needs network) or "
              "point WEBULL_DOCGEN_CACHE at a populated directory.")
        return 0

    return run("write" if args.write else "check")


if __name__ == "__main__":
    sys.exit(main())

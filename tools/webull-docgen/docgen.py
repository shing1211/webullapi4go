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

"""Generate the Webull API documentation artifacts.

Targets:
  reference       SDK-mapped endpoint pages (docs/webull-api/<area>.md)
  master          verbatim Webull guides and reference (docs/webull-api/master-*.md,
                  docs/webull-api/reference/*.md)
  reconciliation  SDK <-> official API coverage table (docs/reconciliation.md)

Usage:
  python docgen.py reference|master|reconciliation|all
"""
import argparse
import datetime
import json
import os
import re
import sys

import _common as c

HK_LLMS = "https://developer.webull.hk/apis/llms.txt"
US_LLMS = "https://developer.webull.com/apis/llms.txt"

BANNER = ("> ⚠️ **Generated file — do not edit.** Regenerate with "
          "`python tools/webull-docgen/docgen.py <target>` "
          "(`reference`, `master`, `reconciliation` or `all`).")


# --------------------------------------------------------------------------
# Target: SDK-mapped reference pages
# --------------------------------------------------------------------------
def generate_reference():
    os.makedirs(c.REFERENCE_OUT, exist_ok=True)
    for area, (title, blurb, eps) in c.AREAS.items():
        parts = ["# %s" % title, "", BANNER, "", blurb, "",
                 "[<- Webull API Reference](../webull-api.md)", ""]
        for (label, url, sdk, note) in eps:
            try:
                parts.append(c.render_endpoint(label, url, sdk, note))
            except Exception as exc:  # noqa: BLE001
                parts.append("### %s\n\n*Error generating: %s* — [official page](%s)\n"
                             % (label, exc, url))
            sys.stderr.write(".")
            sys.stderr.flush()
        with open(os.path.join(c.REFERENCE_OUT, area + ".md"), "w", encoding="utf-8") as fh:
            fh.write("\n".join(parts) + "\n")
    print("\nwrote reference pages (%d areas)" % len(c.AREAS))


# --------------------------------------------------------------------------
# Target: verbatim master guides + reference
# --------------------------------------------------------------------------
def generate_master():
    now = datetime.date.today().isoformat()
    os.makedirs(c.MASTER_OUT, exist_ok=True)

    guides = [
        "# Webull OpenAPI — Master Guides (verbatim)",
        "",
        BANNER,
        "",
        "> Verbatim snapshot of Webull's published OpenAPI **guides**. No "
        "SDK-specific content. Prices/sizes are strings on the wire; see the "
        "endpoint fields in [Master Reference](master-reference.md).",
        "",
        "| | |",
        "|---|---|",
        "| **Snapshot** | %s |" % now,
        "| **Sources** | [developer.webull.hk](%s) and [developer.webull.com](%s) |" % (HK_LLMS, US_LLMS),
        "| **Contents** | Getting started, authentication, market data, trading, broker, connect, errors, FAQ, changelog, AI tools |",
        "| **Refresh** | `python docgen.py master` (re-fetches the `.md` variant of each official page). |",
        "",
        "[<- Webull API Reference](../webull-api.md)",
        "",
    ]
    for section, pages in c.GUIDE_SECTIONS:
        guides.append("## %s" % section)
        guides.append("")
        for label, url in pages:
            try:
                body = c.render_verbatim_page(url)
            except Exception as exc:  # noqa: BLE001
                body = "*Unavailable: %s*" % exc
            guides.append("### %s" % label)
            guides.append("")
            guides.append("> Source: <%s>" % url)
            guides.append("")
            guides.append(body)
            guides.append("")
    with open(os.path.join(c.MASTER_OUT, "master-guides.md"), "w", encoding="utf-8") as fh:
        fh.write("\n".join(guides) + "\n")
    print("wrote master-guides.md")

    ref = [
        "# Webull OpenAPI — Master Reference (verbatim)",
        "",
        BANNER,
        "",
        "> Verbatim snapshot of every Webull-published **endpoint definition** "
        "(OpenAPI schema), split by area. No SDK-specific content. See "
        "[Master Guides](master-guides.md) for authentication, streaming "
        "protocol, trading rules and error codes.",
        "",
        "| | |",
        "|---|---|",
        "| **Snapshot** | %s |" % now,
        "| **Sources** | [developer.webull.hk](%s) and [developer.webull.com](%s) |" % (HK_LLMS, US_LLMS),
        "| **Scope** | All documented endpoints, HTTP and gRPC |",
        "| **Refresh** | `python docgen.py master` |",
        "",
        "[<- Webull API Reference](../webull-api.md)",
        "",
        "## Areas",
        "",
    ]
    for area, (title, blurb, eps) in c.AREAS.items():
        ref.append("- [%s](reference/%s.md) — %d endpoints" % (title, area, len(eps)))
    ref.append("")
    with open(os.path.join(c.MASTER_OUT, "master-reference.md"), "w", encoding="utf-8") as fh:
        fh.write("\n".join(ref) + "\n")
    print("wrote master-reference.md (index)")

    refdir = os.path.join(c.MASTER_OUT, "reference")
    os.makedirs(refdir, exist_ok=True)
    for area, (title, blurb, eps) in c.AREAS.items():
        parts = [
            "# %s — Verbatim Reference" % title,
            "",
            BANNER,
            "",
            "> %s" % blurb,
            "",
            "> Verbatim snapshot of Webull's published OpenAPI definitions. No "
            "SDK-specific content.",
            "",
            "[<- Master Reference](../master-reference.md) · "
            "[<- Webull API Reference](../../webull-api.md)",
            "",
        ]
        for (label, url, sdk, note) in eps:
            try:
                body = c.render_verbatim_page(url)
            except Exception as exc:  # noqa: BLE001
                body = "*Unavailable: %s*" % exc
            parts.append("## %s" % label)
            parts.append("")
            parts.append("> Source: <%s>" % url)
            parts.append("")
            parts.append(body)
            parts.append("")
        with open(os.path.join(refdir, area + ".md"), "w", encoding="utf-8") as fh:
            fh.write("\n".join(parts) + "\n")
    print("wrote reference/*.md (%d files)" % len(c.AREAS))


# --------------------------------------------------------------------------
# Target: reconciliation
# --------------------------------------------------------------------------
def _reference_urls_from_llms(text):
    urls = []
    for line in text.splitlines():
        m = re.search(r"\((https?://[^)]+)\)", line)
        if not m:
            continue
        u = c.norm_url(m.group(1))
        if "/reference/" not in u:
            continue
        if not u.endswith(".md"):
            u = u.rstrip("/") + ".md"
        urls.append(u)
    return list(dict.fromkeys(urls))


def _summary_paths_from_llms(text):
    out = {}
    for line in text.splitlines():
        m = re.search(r"\((https?://[^)]+)\)\s*:\s*(GET|POST|PUT|DELETE)\s+(/\S+)", line)
        if m:
            out[c.norm_url(m.group(1))] = m.group(3)
    return out


def _endpoint_info(url):
    try:
        spec = c.extract_json(c.fetch(url))
    except Exception:  # noqa: BLE001
        return None
    if not spec:
        return None
    return {"method": (spec.get("method") or "").upper(), "path": spec.get("path", "")}


_CATEGORY_RULES = [
    ("Connect API (OAuth)", lambda p: "/oauth2" in p),
    ("Crypto", lambda p: "/crypto" in p),
    ("Display Event Contracts", lambda p: "event-contracts" in p),
    ("Fund Data", lambda p: re.search(r"/fundamentals/fund-", p) is not None),
    ("Authentication (Display client token)", lambda p: "/auth/client-tokens" in p),
    ("Trading", lambda p: p.startswith("/trading")),
    ("Broker FD (US)", lambda p: p.startswith("/broker-fd")),
    ("Broker HK", lambda p: p.startswith("/broker")),
    ("Market Data", lambda p: p.startswith("/market-data")),
]


def _category_for(path):
    for name, rule in _CATEGORY_RULES:
        if rule(path):
            return name
    return "Other"


def _status_row(status):
    return {
        "match": "✅ match",
        "summary": "🟡 SDK matches docs summary, not OpenAPI JSON",
        "differs": "⚠️ path differs from both",
        "no-sdk-path": "❓ SDK path unresolved",
        "no-openapi": "❓ no OpenAPI schema on page",
        "intentional": "ℹ️ intentionally not implemented",
    }.get(status, status)


def _reconcile_data():
    consts = c.parse_consts()

    hk_llms = c.fetch(HK_LLMS)
    us_llms = c.fetch(US_LLMS)
    summary = {}
    summary.update(_summary_paths_from_llms(hk_llms))
    summary.update(_summary_paths_from_llms(us_llms))

    areas_url_to_entry = {}
    for area, (title, blurb, eps) in c.AREAS.items():
        for (label, url, sdk, note) in eps:
            areas_url_to_entry[url] = (area, label, sdk, note)

    llms_urls = set(_reference_urls_from_llms(hk_llms)) | set(_reference_urls_from_llms(us_llms))
    official_all = {}
    for url in sorted(llms_urls | set(areas_url_to_entry)):
        info = _endpoint_info(url)
        if info:
            official_all[url] = info

    implemented_keys = set()
    for url in areas_url_to_entry:
        info = official_all.get(url)
        if info:
            implemented_keys.add((info["method"], c.normalise_path(info["path"])))

    rows = []
    for url, (area, label, sdk, note) in areas_url_to_entry.items():
        info = official_all.get(url)
        sdk_path, sdk_const = c.resolve_method_path(sdk, consts)
        summ = summary.get(url)
        if sdk in ("not implemented", "not exposed"):
            status = "intentional"
        elif info is None:
            status = "no-openapi"
        elif not sdk_path:
            status = "no-sdk-path"
        else:
            ns, nj = c.normalise_path(sdk_path), c.normalise_path(info["path"])
            if ns == nj:
                status = "match"
            elif summ and ns == c.normalise_path(summ):
                status = "summary"
            else:
                status = "differs"
        rows.append((area, label, url, sdk, sdk_path, sdk_const,
                     info["method"] if info else "", info["path"] if info else "",
                     summ or "", status, note))

    uniq = {}
    for url, info in official_all.items():
        if not info["path"]:
            continue
        uniq.setdefault((info["method"], c.normalise_path(info["path"])), []).append((url, info))

    gaps_by_cat = {}
    for key, urls in uniq.items():
        if key in implemented_keys:
            continue
        hk = [u for u, _ in urls if "developer.webull.hk" in u]
        ref = (hk or [u for u, _ in urls])[0]
        canon = next(info["path"] for u, info in urls if u == ref)
        gaps_by_cat.setdefault(_category_for(canon), []).append((key[0], canon, ref))

    counts = {}
    for r in rows:
        counts[r[9]] = counts.get(r[9], 0) + 1
    return rows, gaps_by_cat, counts


def generate_changes():
    """Emit a machine-readable change list for the differing paths."""
    rows, gaps_by_cat, counts = _reconcile_data()
    changes = []
    for (area, label, url, sdk, sdk_path, sdk_const, omethod, opath, summ, status, note) in rows:
        if status != "differs":
            continue
        changes.append({
            "area": area,
            "label": label,
            "sdk_func": sdk,
            "sdk_const": sdk_const,
            "sdk_path": sdk_path,
            "official_method": omethod,
            "official_path": opath,
            "reference": url,
        })
    os.makedirs(c.CACHE, exist_ok=True)
    out = os.path.join(c.CACHE, "changes.json")
    with open(out, "w", encoding="utf-8") as fh:
        json.dump({"count": len(changes), "changes": changes}, fh, indent=2)
    print("wrote %s (%d changes)" % (out, len(changes)))


def generate_reconciliation():
    rows, gaps_by_cat, counts = _reconcile_data()
    gaps_total = sum(len(v) for v in gaps_by_cat.values())

    def render(label, url, sdk, sdk_path, sdk_const, omethod, opath, summ, status, note):
        lines = ["### %s" % label, "", "| | |", "|---|---|"]
        lines.append("| **SDK** | `%s` |" % sdk)
        if omethod and opath:
            lines.append("| **Official (OpenAPI JSON)** | `%s %s` |" % (omethod, opath))
        else:
            lines.append("| **Official** | _no OpenAPI schema_ |")
        if summ:
            lines.append("| **Official (llms.txt summary)** | `%s` |" % summ)
        if sdk_path:
            lines.append("| **SDK path** | `%s` (%s) |" % (sdk_path, sdk_const or "?"))
        lines.append("| **Status** | %s |" % _status_row(status))
        if note:
            lines.append("| **Note** | %s |" % note)
        lines.append("")
        lines.append("Reference: [%s](%s)" % (url.rsplit("/", 1)[-1], url))
        lines.append("")
        return "\n".join(lines)

    out = [
        "# SDK ↔ Webull API Reconciliation",
        "",
        BANNER,
        "",
        "> Reconciles every implemented `webullapi4go` function against the "
        "official Webull OpenAPI. **Official (OpenAPI JSON)** is the canonical "
        "path embedded in the docs; **Official (llms.txt summary)** is the path "
        "in Webull's machine-readable index (they disagree for some endpoints). "
        "**SDK path** is what the code actually calls.",
        "",
        "| | |",
        "|---|---|",
        "| **Snapshot** | %s |" % datetime.date.today().isoformat(),
        "| **Sources** | [HK llms.txt](%s), [US llms.txt](%s) |" % (HK_LLMS, US_LLMS),
        "| **Implemented endpoints** | %d |" % len(rows),
        "| **Documented-only endpoints (gaps)** | %d |" % gaps_total,
        "| ✅ Path matches OpenAPI JSON | %d |" % counts.get("match", 0),
        "| 🟡 Matches docs summary only | %d |" % counts.get("summary", 0),
        "| ⚠️ Path differs from both | %d |" % counts.get("differs", 0),
        "| ❓ Unresolved | %d |" % (counts.get("no-sdk-path", 0) + counts.get("no-openapi", 0)),
        "| ℹ️ Intentionally not implemented | %d |" % counts.get("intentional", 0),
        "",
        "## Implemented endpoints",
        "",
    ]
    current = None
    for area, label, url, sdk, sdk_path, sdk_const, omethod, opath, summ, status, note in rows:
        if area != current:
            current = area
            out.append("## %s" % c.AREAS[area][0])
            out.append("")
        out.append(render(label, url, sdk, sdk_path, sdk_const, omethod, opath, summ, status, note))

    out.append("## Documented but not implemented")
    out.append("")
    out.append("Unique official endpoints (deduplicated by method and path) that "
               "`webullapi4go` does not implement. Duplicate HK/US references to "
               "the same endpoint are collapsed into one row.")
    out.append("")
    for cat in sorted(gaps_by_cat):
        out.append("### %s" % cat)
        out.append("")
        out.append("| Method | Path | Reference |")
        out.append("|---|---|---|")
        for method, path, ref in sorted(gaps_by_cat[cat]):
            out.append("| %s | `%s` | [%s](%s) |" % (method, path, ref.rsplit("/", 1)[-1], ref))
        out.append("")

    with open(c.RECON_OUT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(out) + "\n")
    print("wrote %s (implemented=%d gaps=%d match=%d summary=%d differs=%d unresolved=%d intentional=%d)"
          % (c.RECON_OUT, len(rows), gaps_total, counts.get("match", 0), counts.get("summary", 0),
             counts.get("differs", 0), counts.get("no-sdk-path", 0) + counts.get("no-openapi", 0),
             counts.get("intentional", 0)))


def main():
    ap = argparse.ArgumentParser(description="Webull API documentation generator")
    ap.add_argument("target", choices=["reference", "master", "reconciliation", "changes", "all"])
    args = ap.parse_args()
    if args.target in ("reference", "all"):
        generate_reference()
    if args.target in ("master", "all"):
        generate_master()
    if args.target in ("reconciliation", "all"):
        generate_reconciliation()
    if args.target in ("changes", "all"):
        generate_changes()


if __name__ == "__main__":
    main()

# Webull doc generator

Generates the Webull API documentation under `docs/` from Webull's own
machine-readable documentation. Every official page has a `.md` variant that
embeds the OpenAPI definition JSON; this tool fetches and caches those pages and
renders them.

No SDK-specific knowledge lives outside the area manifest in `_common.py`; the
tool only fetches public Webull documentation (no credentials, no secrets).

## Usage

```sh
python docgen.py reference       # SDK-mapped endpoint pages  -> docs/webull-api/<area>.md
python docgen.py master          # verbatim guides + reference -> docs/webull-api/master-*.md, docs/webull-api/reference/*.md
python docgen.py reconciliation  # SDK <-> official coverage   -> docs/reconciliation.md
python docgen.py all             # all of the above
```

`docs/webull-api.md` (the reference hub) is hand-written and is not generated.

Every generated file starts with a "Generated file — do not edit" banner; if you
find yourself editing one, change the template in `docgen.py` (or the manifest in
`_common.py`) and regenerate instead.

## Sources

- HK: <https://developer.webull.hk/apis/llms.txt>
- US: <https://developer.webull.com/apis/llms.txt>

Fetched pages are cached in `.cache/` (git-ignored). Delete it to force a
refresh; the first run needs network access, later runs are offline.

## Outputs

| Target | Files |
|--------|-------|
| `reference` | `docs/webull-api/{authentication,market-data-*,fundamentals,event-contracts,trading,broker-hk,broker-fd-us,display-solution,streaming,events}.md` |
| `master` | `docs/webull-api/master-guides.md`, `docs/webull-api/master-reference.md`, `docs/webull-api/reference/*.md` |
| `reconciliation` | `docs/reconciliation.md` |

The rendered pages apply a minimal, documented normalization: Docusaurus
directives are dropped, inner headings demoted, and unresolvable internal links
flattened. Fenced code blocks (the OpenAPI JSON) are preserved byte-for-byte.

## Requirements

Python 3.8+ with the standard library only.

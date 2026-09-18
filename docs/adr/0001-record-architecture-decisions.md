---
Status: Accepted
Date: 2026-09-18
Deciders: maintainer
---

# 0001. Record architecture decisions

## Context

`webullapi4go` is a public Go SDK that wraps the Webull OpenAPI. It ships in
stages (v0.1 Foundation + Market Data, v0.2 Trading HTTP, v0.3 gRPC events, and
so on). Several decisions have non-obvious constraints: licensing of vendored
protobuf definitions, signature/normalization quirks, and which Webull API
surfaces are stable enough to expose.

Without a written record, the reasoning behind these choices is lost between
phases, and the same trade-offs get re-argued in issues and pull requests.

## Decision

We will use Architecture Decision Records (ADRs).

- ADRs live in `docs/adr/`.
- Each ADR is a numbered Markdown file, `NNNN-short-title.md`.
- Each ADR records status, date, deciders, context, the options considered, the
  decision, the rationale, and the consequences.
- Accepted ADRs are immutable; a changed decision is recorded as a new ADR that
  supersedes the old one.
- `docs/adr/README.md` documents the format and indexes all ADRs.

## Rationale

- Cheap to write, reviewable in a pull request, and versioned with the code.
- Makes the boundary between verified facts and open questions explicit, which
  matters for surfaces Webull does not fully document.
- Standard practice (see Michael Nygard's original ADR write-up); contributors
  are likely to already recognize the format.

## Consequences

- Contributors are expected to check the index and add or supersede an ADR when
  they make a significant, hard-to-reverse decision.
- Small, mechanical changes do not need an ADR.
- The ADR set becomes part of the project's documentation, not a separate wiki.

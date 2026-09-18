# Architecture Decision Records (ADRs)

This directory holds the Architecture Decision Records for `webullapi4go`.
An ADR captures a single significant architectural decision: the context that
forced it, the options considered, the choice made, and the consequences.

## Why ADRs

- Decisions are recorded once, in one place, instead of being re-litigated in
  issues and pull requests.
- The rationale survives contributor turnover.
- Superseded decisions stay on disk, so the history of the design is auditable.

## Format

Each ADR is a Markdown file named `NNNN-short-title.md`, where `NNNN` is a
zero-padded, monotonically increasing number. Numbers are never reused.

Every ADR starts with a metadata block:

```markdown
---
Status: Proposed | Accepted | Superseded by ADR-NNNN | Deprecated
Date: YYYY-MM-DD
Deciders: <who decided>
Supersedes: ADR-NNNN   # optional
---
```

Status meanings:

- **Proposed** — written up but not yet ratified by the maintainer.
- **Accepted** — the decision is in force.
- **Superseded by ADR-NNNN** — replaced; the new ADR is authoritative.
- **Deprecated** — no longer recommended, but not yet replaced.

Recommended body sections: `Context`, `Problem`, `Options Considered`
(with pros/cons), `Decision`, `Rationale`, `Consequences`, `Risks and
Unknowns`, `Spike Plan` (when follow-up verification is required).

## Conventions

- One decision per ADR. Do not bundle unrelated choices.
- State facts and inferences separately. Mark anything inferred or unverified
  explicitly with **Inferred** or **To verify**.
- Do not fabricate details (message names, field numbers, endpoints). When a
  detail is unknown, write "to be determined".
- ADRs are immutable once `Accepted`. To change a decision, add a new ADR that
  supersedes the old one; do not rewrite history.

## Index

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [0001](0001-record-architecture-decisions.md) | Record architecture decisions | Accepted | 2026-09-18 |
| [0002](0002-grpc-event-protobuf-strategy.md) | gRPC event protobuf strategy | Proposed | 2026-09-18 |

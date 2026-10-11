---
Status: Accepted
Date: 2026-10-11
Deciders: maintainer
---

# 0005. Module path and import versioning

## Context

Go's semantic import versioning requires major-version bumps in the module path when
breaking changes are introduced (`module/path/v2`, `module/path/v3`, etc.).
`webullapi4go` is not yet at v1.0. Decisions about when (or whether) to adopt `/v2`
need to be recorded now, before v1.0, so that consumers understand the upgrade
contract.

## Decision

The module path is and remains:

```
github.com/shing1211/webullapi4go
```

No `/v2`, `/v3`, etc. suffix is used, now or in any planned future version.

The module is versioned with Git tags on the `main` branch using semver tags
(`v2.1.39`, etc.). However, **the module proxy does not serve these tags**;
it serves only the `v1.x` line (the last tag is `v1.1.1`). Tags `v2.x`
are present in the git tree but are not installable via `go get`.

To pin a specific commit, consumers use:

```
go get github.com/shing1211/webullapi4go@<commit-sha>
```

## Rationale

1. **No v2 planned**: The maintainer does not anticipate a breaking redesign that
   would require a major version bump. The API surface (typed clients, error codes,
   resilience primitives) is considered stable enough to not need a clean break.

2. **Migration cost**: Moving to `/v2` requires updating every import statement in
   every consumer. For an in-house SDK used across multiple repositories, this is
   a high-friction change with no immediate benefit.

3. **v2 tags are for documentation**: Tags `v2.x` record the progression of the
   documentation/hardening effort, not a Go module major version. They are not
   installable and should not be used as version selectors in `go.mod`.

4. **Upgrade path preserved**: If a `/v2` is eventually required, it can be created
   from the same `main` branch at that time. The absence of a `/v2` today does
   not prevent one in the future.

## Consequences

- `go get github.com/shing1211/webullapi4go` resolves to `v1.1.1`.
- Consumers who want newer work must pin a commit SHA or use the GitHub raw URL.
- Import paths never need to change, even if breaking changes are introduced
  (the maintainer commits to ADR if that ever happens).
- The `v2.x` tags in the git tree are an internal versioning aid, not a Go
  module versioning signal.

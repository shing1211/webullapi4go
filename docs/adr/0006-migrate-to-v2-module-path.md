---
Status: Accepted
Date: 2026-10-28
Deciders: maintainer
Supersedes: ADR-0005
---

# 0006. Migrate to /v2 module path

## Context

ADR 0005 decided to keep the module on the v1 import path
(`github.com/shing1211/webullapi4go`) and treat `v2.x` Git tags as
repository-only, non-installable markers. The Go module proxy therefore served
only `v1.1.1`, and consumers pinned commits to reach newer work.

The maintainer has now authorized a `/v2` migration and a published Go-semver
v2 module release.

## Decision

The module path is now:

```
github.com/shing1211/webullapi4go/v2
```

Every import statement, `go.mod` declaration, `replace` directive, `buf.gen.yaml`
M mapping, and `golangci-lint` `local-prefixes` carries the `/v2` suffix. The
first installable Go-semver v2 tag is `v2.2.0`. The 66 prior `v2.1.x` repository
tags declared the pre-migration module path and were not installable; they have
been deleted from both remotes (`origin` and `gitee`).

## Rationale

1. **Installability**: `go get github.com/shing1211/webullapi4go/v2@latest`
   resolves to `v2.2.0` instead of failing or falling back to `v1.1.1`.
2. **Semantic import versioning**: Go requires the major version in the module
   path for v2+. Carrying `/v2` is the idiomatic mechanism.
3. **Clean release**: Deleting the non-installable `v2.1.x` tags prevents them
   from polluting `@latest` resolution — a `v2.1.40` tag whose `go.mod` declares
   the old path would make `@latest` error with "module declares its path as
   ... not .../v2".

## Consequences

- `go get github.com/shing1211/webullapi4go/v2@latest` installs `v2.2.0`.
- All consumers must update import paths from `github.com/shing1211/webullapi4go`
  to `github.com/shing1211/webullapi4go/v2`.
- The `v1.x` line (`v1.0.0`–`v1.1.1`) remains installable at the old path
  `github.com/shing1211/webullapi4go` for consumers that have not migrated.
- The `broker/` nested module path is
  `github.com/shing1211/webullapi4go/v2/broker`.
- ADR 0005 is superseded.

# Contributing to webullapi4go

Thanks for your interest in `github.com/shing1211/webullapi4go`, a public Go SDK
for the Webull OpenAPI. This document explains how to set up a development
environment, the standards we expect, and where to get help.

By participating you agree to follow our [Code of Conduct](CODE_OF_CONDUCT.md).
For security issues, do not open a public issue: follow [SECURITY.md](SECURITY.md)
instead.

- Questions, ideas, and show-and-tell: [Discussions](https://github.com/shing1211/webullapi4go/discussions)
- Bugs and feature requests: [Issues](https://github.com/shing1211/webullapi4go/issues)

## Prerequisites

- Go 1.26 or newer. The module declares `go 1.26` in [`go.mod`](go.mod).
  Install from https://go.dev/dl/.
- Git.
- Recommended: [`golangci-lint`](https://golangci-lint.run/) matching the version
  used by CI (see `.github/workflows/`).
- Optional: [`buf`](https://buf.build/) for protobuf regeneration, and a Webull
  sandbox account for integration tests.

Confirm your toolchain:

```sh
go version   # go1.26 or newer
```

## Getting started

```sh
git clone https://github.com/shing1211/webullapi4go.git
cd webullapi4go
go mod download
```

Build and check the tree:

```sh
go build ./...
go vet ./...
gofmt -l .        # must print nothing
```

## Testing

Unit tests need no credentials and must pass offline:

```sh
go test ./...
go test -race ./...
```

Lint with:

```sh
golangci-lint run
```

### Documentation

The MkDocs Material site lives in `docs/`. Build it in strict mode before
shipping doc changes:

```sh
python -m venv .venv
.venv/Scripts/pip install -r requirements-docs.txt   # Windows
pip install -r requirements-docs.txt                 # macOS / Linux
mkdocs build --strict
```

Internal run artifacts under `docs/runs/` are excluded from the published site.

### Sandbox integration tests

Integration tests hit the Webull sandbox and are skipped unless explicitly
enabled. Webull publishes shared public test accounts for the HK sandbox —
see [Sandbox > Test credentials](sandbox.md#test-credentials) for the
values.

```sh
# macOS / Linux
WEBULL_SANDBOX=1 \
WEBULL_APP_KEY=your-sandbox-app-key \
WEBULL_APP_SECRET=your-sandbox-app-secret \
go test ./... -run Integration
```

```powershell
# Windows PowerShell
$env:WEBULL_SANDBOX = "1"
$env:WEBULL_APP_KEY = "your-sandbox-app-key"
$env:WEBULL_APP_SECRET = "your-sandbox-app-secret"
go test ./... -run Integration
```

Private credentials must never be committed. The shared public test accounts
published by Webull are the exception — they are public by design.

## Branch naming

Create short-lived branches from `main`:

- `feat/<short-description>` for new features
- `fix/<short-description>` for bug fixes
- `docs/<short-description>` for documentation
- `refactor/<short-description>` for internal changes
- `test/<short-description>` for test-only changes
- `chore/<short-description>` for tooling and maintenance

## Commit convention

Commits follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<optional scope>): <description>

[optional body]

[optional footer(s)]
```

Common types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `build`, `ci`,
`chore`. Keep each commit focused, write the subject in the imperative mood, and
keep it to 72 characters or fewer.

### Sign-off

All commits must carry a `Signed-off-by` line certifying the
[Developer Certificate of Origin](https://developercertificate.org/):

```sh
git commit -s -m "feat(marketdata): add snapshot endpoint"
```

Cryptographically signed commits (GPG or SSH) are welcome but not required.
Contributions are accepted under the project's [Apache-2.0 license](LICENSE).

## Pull requests

- Open one logical change per pull request and keep the diff small.
- Add or update tests for behavior changes.
- Update the README, docs, and ADRs where behavior or decisions change. See
  `docs/adr/index.md` for the ADR format and conventions.
- Ensure `go build ./...`, `go vet ./...`, `go test ./...`, and
  `golangci-lint run` pass before requesting review.
- Fill in the [pull request template](.github/pull_request_template.md).
- Pull requests are squash-merged into `main`.

## Community

Discussions is enabled for this repository. Use it for questions, design
discussions, and announcements:

- https://github.com/shing1211/webullapi4go/discussions

If Discussions is ever unavailable, a maintainer can re-enable it in the
repository settings: **Settings → Features → check "Discussions" → Set up
discussions**.

## Reporting security issues

Never report vulnerabilities in public issues, pull requests, or Discussions.
Use GitHub Security Advisories as described in [SECURITY.md](SECURITY.md).

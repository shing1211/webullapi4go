# Todos — v03-grpc-events

Single source of truth. Statuses: `todo`, `doing`, `blocked`, `review`, `done`.

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|------------|------------|
| T15.1 | Vendor `events.proto` + attribution + generate Go; ADR-0002 → Accepted | data | done | — | byte-identical (`7ebc4728…`); grpc stub generated |
| T15.2 | Algorithm-parameterised signer (HMAC-SHA256 + SHA256 body) | security | done | — | SHA1 golden unchanged; SHA256 vector passes |
| T15.3 | `events` package: connect, signed metadata, Subscribe, dispatch | backend | done | T15.1, T15.2 | sandbox SubscribeSuccess (lowercase SHA256; account must match app key) |
| T15.4 | Typed JSON payloads + reconnect/resubscribe | backend | done | T15.3 | live OrderEvent decoded; reconnect tests pass |
| T15.5 | Docs/events.md + example + README/CHANGELOG | docs | done | T15.4 | events docs + example; mkdocs --strict OK |
| T15.6 | Nightly read-only live CI job | devops | done | T15.3 | nightly-live.yml valid; mutation off |
| T15.7 | Release v0.3.0 | release | doing | T15.5, T15.6 | tag on both remotes |

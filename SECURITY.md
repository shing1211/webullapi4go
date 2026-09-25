# Security Policy

## Supported versions

The latest repository tag is `v2.1.1` (2026-09-25). Current request, OMS,
streaming, telemetry, and documentation hardening is tagged in repository
`v2.1.1`, but this is a repository Git patch release, not a
Go-semver-compatible v2 module. The root module path remains
`github.com/shing1211/webullapi4go`, and publication of a v2 module remains
deferred. The tagged hardening was not newly live-verified. Security fixes for
that work are prepared on `main`; the repository tag is not a published v2
module security-fix line. The pinned `v1.1.1` install command is historical v1
guidance.

| Version or state | Supported |
|---|---|
| `v2.1.1` | Repository tag only; no published v2 module; fixes are tracked in the repository |
| `v2.1.0` and earlier `v2.x` tags | No — historical Git tags only; v2 module publication is deferred |
| `v2.0.x` and earlier | No |

Older lines may receive a fix when the correction is low risk and backporting
it does not create a disproportionate maintenance burden. A backport is not
guaranteed; use a separately published module release when one is authorized.

## Reporting a vulnerability

Do not report security vulnerabilities in public issues, pull requests, or
Discussions.

Report privately through GitHub Security Advisories:

1. Open `https://github.com/shing1211/webullapi4go/security/advisories/new`.
2. Describe the issue, affected versions, impact, and reproduction steps.
3. If the advisory form is unavailable, open a minimal public issue asking for
   a private channel without including vulnerability details.

Include:

- A description and impact assessment.
- Steps to reproduce or a proof of concept.
- Affected versions and environment (OS and Go version).
- Any suggested remediation.

Response expectations:

- Acknowledgement within 3 business days.
- Initial assessment and severity within 10 business days.
- Coordinate updates and the disclosure timeline.
- Credit reporters in the advisory unless they prefer anonymity.

## Credentials are secrets

This SDK authenticates with a Webull App Key, App Secret, access tokens, and
account-scoped data. Treat all of them as secrets.

- Never commit App Keys, App Secrets, access tokens, account credentials, or
  `.env` files. `.env` is ignored by Git.
- Never paste real credentials into issues, pull requests, Discussions, logs,
  screenshots, examples, or test fixtures.
- Supply credentials at runtime through `WEBULL_APP_KEY` and
  `WEBULL_APP_SECRET`, direct `client.WithCredentials` configuration, or a
  secret manager.
- Use throwaway sandbox credentials for development and tests.
- Treat an exposed credential as compromised and rotate it immediately in the
  Webull developer portal.
- Do not place secrets in correlation IDs, trace attributes, metric labels, or
  application log fields.

Shared public test accounts published by Webull in its official documentation
are not application credentials and may be referenced in documentation and
tests.

## Telemetry and logging

The SDK's REST, MQTT, and gRPC telemetry does not include App Keys, App
Secrets, access tokens, or signing values. Production telemetry uses
`observability.SafeErrorText`: typed errors are reduced to their category,
context errors retain only the standard cancellation/deadline text, and other
errors become `operation failed`. Response bodies, gRPC status messages, and
wrapped error causes are not copied into SDK spans or logs.

Telemetry still includes request paths without query strings, HTTP status,
timing, attempt numbers, topic/category labels, RPC metadata names, and
correlation IDs. Those fields can be operationally sensitive even though they
are not authentication material.

Applications remain responsible for:

- configuring `slog` redaction and backend access controls;
- avoiding raw account/order payloads in sensitive log levels;
- keeping secrets out of trace attributes, baggage, metric labels, and
  correlation IDs;
- treating exported telemetry as sensitive operational data.

See [Observability](docs/observability.md) for the emitted fields and setup.

## Reporting an exposed credential

If a credential was exposed in the repository, issues, logs, or telemetry,
report it through the private channel above so it can be rotated and repository
history can be cleaned where appropriate.

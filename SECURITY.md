# Security Policy

## Supported versions

The latest tagged release is `v2.1.0`. Security fixes are applied to the
current `v2.1.x` release line on `main`.

| Version | Supported |
|---|---|
| `v2.1.x` | Yes |
| `v2.0.x` and earlier | No |

Older lines may receive a fix when the correction is low risk and backporting
it does not create a disproportionate maintenance burden. A backport is not
guaranteed; upgrade to the current release line for security fixes.

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

The SDK's request and event telemetry does not include App Keys, App Secrets,
access tokens, or signing values. It does include request paths, status,
timing, attempt numbers, RPC metadata names, and correlation IDs.

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

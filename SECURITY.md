# Security Policy

## Supported versions

`webullapi4go` is pre-1.0. Security fixes are applied to the latest release on
the `main` branch.

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | Yes                |
| < 0.1   | No                 |

We do not backport security fixes to older pre-1.0 minor versions.

## Reporting a vulnerability

Please do not report security vulnerabilities in public issues, pull requests,
or Discussions.

Report privately through GitHub Security Advisories:

1. Open https://github.com/shing1211/webullapi4go/security/advisories/new
2. Describe the issue, including affected versions, impact, and reproduction
   steps.
3. If you cannot use the advisory form, open a minimal public issue asking for a
   private channel without including any vulnerability details.

Please include:

- A description of the vulnerability and its impact.
- Steps to reproduce, or a proof of concept.
- Affected versions and environment (OS, Go version).
- Any suggested remediation.

## Response expectations

- Acknowledgement within 3 business days.
- Initial assessment and severity within 10 business days.
- We will keep you updated and coordinate a disclosure timeline with you.
- Reporters are credited in the advisory unless they prefer to remain
  anonymous.

## Never commit credentials

This SDK authenticates with a Webull App Key, App Secret, and access tokens.
These are secrets.

- Never commit App Keys, App Secrets, tokens, or `.env` files. `.env` is already
  covered by [`.gitignore`](.gitignore).
- Never paste real credentials into issues, pull requests, Discussions, logs,
  screenshots, or test fixtures.
- Pass credentials at runtime through environment variables
  (`WEBULL_APP_KEY`, `WEBULL_APP_SECRET`) or a secret manager.
- Use the sandbox environment and throwaway sandbox credentials for development
  and tests.
- Treat any exposed credential as compromised and rotate it in the Webull
  developer portal immediately.

If you believe a credential was exposed in this repository, report it through
the private channel above so it can be rotated and the repository history can be
cleaned.

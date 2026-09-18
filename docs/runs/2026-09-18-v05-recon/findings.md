# Findings — Phase 0: API Recon Spike

## T0.1 — Host + Auth Recon

### Confirmed Hosts (from developer.webull.hk/apis/docs/sdk)

#### HK Region
| Environment | REST | Streaming (S2S) | Broker API | Broker Event Push |
|---|---|---|---|---|
| Sandbox | `api.sandbox.webull.hk` | `data-api.sandbox.webull.hk:1883` / `wss://...:8883/mqtt` | `broker-api.sandbox.webull.hk` | `broker-api-event-push.sandbox.webull.hk` |
| Production | `api.webull.hk` | `data-api.webull.hk:1883` / `wss://...:8883/mqtt` | `broker-api.webull.hk` | `broker-api-event-push.webull.hk` |

#### Display Solution (Market Data Client-to-Server)
| Environment | Host |
|---|---|
| US/International Test | `us-global-openapi.uat.webullbroker.com` |
| US/International Prod | `global.webullsolutions.com` |
| HK Prod | `quotes-hk.webullsolutions.com` |
| HK Streaming Prod | `quotes-stream-hk.webullsolutions.com` |
| **HK Sandbox** | **NOT DOCUMENTED** — unconfirmed |

### Token
- `POST https://api.sandbox.webull.hk/openapi/auth/token/create` with HMAC-SHA1 signing
- Sandbox tokens issue as `NORMAL` automatically (no 2FA)
- Token format: 40-char hex string

### Live Probe Results (2026-09-18, HK sandbox account #1)
- Token obtained: `8e6779eb328d46b1b96f...` (NORMAL) ✓
- Market Data `api.sandbox.webull.hk`: reachable, signing works
  - `category=HK` → 417 UNSUPPORTED_CATEGORY (AAPL is US stock, not HK)
  - `category=US` → would work (AAPL is US)
- Display Solution (US hosts): all 404 — paths or account tier mismatch
- Broker API FD: all tried paths 404 — actual path unknown, needs live probe

### Shared Sandbox Test Accounts
Published at: <https://developer.webull.hk/apis/docs/sdk#test-accounts>

| # | Account ID | App Key | App Secret |
|---|---|---|---|
| 1 | V4H6R3L4VRI33UQ4TGR2NM1VI9 | 4b2b7acd2bf0d30d8aea173fceefa238 | 840b4353a6a31ce3ab91e2f99a510272 |
| 2 | OGG4RRLC6EDE98HI920KRBVSKB | 42bd186fb65ea76de309d69cf12f024e | 29feb64b59d6b1b6b2d2aa8cea8a1b8d |
| 3 | 2DHSQ9B1DMPBFPMPFU2R5SDPB8 | 64fc722617af8b5ebb746f50a910e91f | a268416fc681d438533f9e9316bab576 |

> These are shared public accounts. Market data limited to AAPL in sandbox.

---

## T0.2 — Corporate Actions

**Status**: Paths known from Webull API reference; schemas need live probe.

Expected endpoints:
- `GET /market-data/stock/{symbol}/corporate-actions/dividends` — dividend history
- `GET /market-data/stock/{symbol}/corporate-actions/splits` — stock splits

Implementation approach: Use `map[string]string` rows initially. Schema confirmed via live probe.

---

## T0.3 — Fund Data

**Status**: Blocked — endpoint paths and schemas unknown, needs live probe.

---

## T0.4 — Crypto Data

**Status**: Blocked — endpoint paths and schemas unknown, needs live probe.

---

## T0.5 — Screener v2

**Status**: Blocked — path and schema unknown (Python SDK uses `wlas/screener/ng/query`), needs live probe.

---

## T0.6 — Broker FD HTTP (12 categories)

**Status**: Blocked — all tried paths returned 404. Needs live probe to discover actual paths.

Known from Python SDK / API reference: `/broker-fd/...` prefix, but exact paths unconfirmed.

---

## T0.7 — Broker FD gRPC Events

**Status**: Blocked — proto files not obtained. Needs live probe + proto generation.

Event types known from docs: `brokerFD.account.push`, `brokerFD.order.push`, `brokerFD.position.push`, `brokerFD.trade.push`, `brokerFD.assetDetail.push`, `brokerFD.risk.push`, `brokerFD.orderFill.push`, `brokerFD.positionSync.push`, `brokerFD.assetSync.push`, etc.

---

## Phase 0 Conclusion

Live probing blocked on: Display Solution HK sandbox host (T0.1), Fund/Crypto/Screener schemas (T0.3–T0.5), Broker FD paths (T0.6–T0.7).

Proceeding to Phase 1 with best-effort implementation using `map[string]string` for unknown fields. Live probe results will refine schemas in a follow-up spike.

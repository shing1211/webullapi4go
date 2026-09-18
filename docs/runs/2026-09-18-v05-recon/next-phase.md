# Next Phase — webullapi4go Phase 0: API Recon Spike

## Status
Phase 0 (T0.1–T0.7) is blocked on `WEBULL_APP_KEY` and `WEBULL_APP_SECRET` credentials.
Hosts and auth flows are confirmed from official docs; live schema probing cannot proceed.

## Candidate Next Phase Items

### 1. Phase 1: v0.5 Implementation (Corporate Actions, Instrument v3, Logos, Fund/Crypto/Screener v2)
**Priority: HIGH** — Adds substantial new coverage with well-understood signing.

Known work (no live probing required):
- Instrument v3: new path `/market-data/instruments/stocks/list` (new, not overlap)
- Corporate Actions: 2 endpoints from API reference (paths known; schemas need live probe)
- Logos batch: new endpoint, confirmed from docs

Unknown work (blocked on T0.2–T0.5 live probing):
- Fund data schemas
- Crypto tick/depth/orderbook schemas
- Screener v2 path + schema

### 2. Phase 2: Broker FD (HTTP + gRPC)
**Priority: MEDIUM** — Full live probe needed (T0.6–T0.7). Blocked until credentials provided.

### 3. Phase 3: v1.0 API Lock
**Priority: LOW** — Ship when v0.6 is done; no rush.

## Recommended First Task

**Provide sandbox credentials** (`WEBULL_APP_KEY` + `WEBULL_APP_SECRET`) to unblock T0.2–T0.7.

Once credentials are available:
1. Run live probing (T0.2–T0.7) to fill in all unknown schemas
2. Proceed to Phase 1 implementation with confirmed schemas

## What Can Be Done Without Credentials

While waiting for credentials, the following can proceed:
- Write `data/corporate_actions.go` for the 2 confirmed Corporate Actions endpoints using `map[string]string` rows
- Write `data/instrument_v3.go` for Instrument v3
- Write `data/logos.go` for Logos batch
- Write the protobuf files for Broker FD gRPC (event types known from docs: `brokerFD.account.push`, `brokerFD.order.push`, `brokerFD.position.push`, `brokerFD.trade.push`, `brokerFD.assetDetail.push`, `brokerFD.risk.push`, etc.)

These can all use `map[string]string` for unknown fields and be refined with live data later.

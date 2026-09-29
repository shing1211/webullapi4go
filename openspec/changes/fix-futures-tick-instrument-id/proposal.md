# Proposal

## Why

The 2026-09-29 live-evidence walk found one demonstrated SDK defect: the sandbox response for `data.GetFuturesTick` carries `instrumentId`, while the documented contract and `data.StockTicks` use `instrument_id`. `encoding/json` does not treat the underscore difference as a case variation, so callers receive a successful response with `InstrumentID` silently set to `""`. The SDK should tolerate both published and observed spellings without changing its public field or method signatures.

## What Changes

- Add dual-name decoding for `data.StockTicks` so `instrument_id` and `instrumentId` both populate the existing `InstrumentID` field.
- Keep marshaling and the documented `json:"instrument_id"` tag unchanged; the change is additive and non-breaking.
- Define deterministic precedence when both spellings are present, with the documented spelling winning.
- Add regression tests covering both wire spellings, precedence, empty/unknown fields, and unchanged output serialization.
- Re-run the live conformance comparison and update the recorded live divergence set and status documentation so the finding is gated rather than silently removed.
- Do not claim that one sandbox host establishes the server-wide contract; a second host remains a future confirmation step.

## Capabilities

### New Capabilities

- `market-data/futures-tick-wire-compatibility`: Decode futures tick responses regardless of whether Webull sends the documented `instrument_id` name or the observed `instrumentId` name, while preserving the public `StockTicks` API and documented serialization.

### Modified Capabilities

None. The repository has no existing OpenSpec capabilities yet.

## Impact

- Affected code: `data/tick.go`, futures/tick regression tests, and the live-conformance divergence record.
- Affected public API: none at the type or method-signature level; decoding behavior becomes more tolerant.
- Dependencies: none.
- Evidence boundary: the change is justified by one HK sandbox run and the documentation disagreement; it does not claim production verification.
- Non-goals: changing the public tag to camelCase, adding a second exported alias field, changing unrelated tick endpoints, or fabricating a second-host result.

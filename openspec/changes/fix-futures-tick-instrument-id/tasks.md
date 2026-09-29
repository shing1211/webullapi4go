# Tasks

## 1. Specification and design

- [ ] 1.1 Write the spec delta at `specs/market-data/futures-tick-wire-compatibility/spec.md` and verify `openspec validate fix-futures-tick-instrument-id --type change --strict` reports no error

## 2. Decoder

- [ ] 2.1 Add `UnmarshalJSON` to `*data.StockTicks` in `data/tick.go` using a type alias plus a shadow camelCase field, and verify a response carrying `instrument_id` alone still decodes to the same value
- [ ] 2.2 Verify a response carrying `instrumentId` alone populates `InstrumentID`, and that a response carrying both leaves the documented spelling in place
- [ ] 2.3 Verify marshalling still emits `instrument_id` and not `instrumentId`, so no existing consumer sees a renamed member

## 3. Tests

- [ ] 3.1 Add table-driven cases to `data/tick_test.go` covering snake_case, camelCase, both-present precedence, neither-present, and the marshal round-trip; verify `go test ./data/...` passes
- [ ] 3.2 Add a `GetFuturesTick` case whose body carries `instrumentId` and verify the returned `InstrumentID` is populated from it
- [ ] 3.3 Run `gofmt -l .` and verify it prints nothing

## 4. Conformance prose

- [ ] 4.1 Rewrite the `data.GetFuturesTick` reason and unblock in `conformance/livereasons_test.go` so neither claims a caller receives an empty identifier, and verify the wording matches the JSON baseline
- [ ] 4.2 Mirror that prose in `conformance/live-divergences.json`, keeping the entry's key, kind, name, and bucket unchanged
- [ ] 4.3 Update the stale assertions in `conformance/live_test.go` and `conformance/doc.go` that describe the SDK as currently returning a zero value
- [ ] 4.4 Run `make conformance-live-gate` and verify it passes with the row still recorded and the live-only count still 13

## 5. Status documentation

- [ ] 5.1 Rewrite item 28 and next step 10 in `IMPLEMENTATION_STATUS.md` to record the decoder as applied and the second-host probe as the only remaining question, keeping the 13-finding count
- [ ] 5.2 Apply the same rewrites to the corresponding sections of `docs/implementation-status.md`
- [ ] 5.3 Run `make citations` and verify it exits zero
- [ ] 5.4 Add a `CHANGELOG.md` entry describing the decoder, the preserved public tag, and the retained live row

## 6. Integration checks

- [ ] 6.1 Run `make build`, `make vet`, and `make test`, and verify all three pass without credentials
- [ ] 6.2 Run `make conformance-gate` and verify the documented baseline is unchanged at its recorded entry count

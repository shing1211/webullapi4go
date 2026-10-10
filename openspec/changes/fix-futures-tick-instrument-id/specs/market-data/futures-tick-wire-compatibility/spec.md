# Spec Delta

## Purpose

Decodes futures and stock tick responses regardless of whether Webull sends the documented `instrument_id` name or the observed `instrumentId` name, so a caller reading `StockTicks.InstrumentID` receives the identifier the server sent.

## ADDED Requirements

### Requirement: Both instrument identifier spellings decode

The SDK SHALL populate `StockTicks.InstrumentID` from either the `instrument_id` member or the `instrumentId` member, and SHALL NOT return an error for either spelling.

#### Scenario: Documented snake_case spelling

- **WHEN** a tick response carries `instrument_id`
- **THEN** `StockTicks.InstrumentID` holds that value

#### Scenario: Observed camelCase spelling

- **WHEN** a tick response carries `instrumentId` and no `instrument_id`
- **THEN** `StockTicks.InstrumentID` holds that value

#### Scenario: Neither spelling present

- **WHEN** a tick response carries neither `instrument_id` nor `instrumentId`
- **THEN** `StockTicks.InstrumentID` is the empty string and decoding succeeds without error

### Requirement: Precedence is deterministic

When a tick response carries both spellings, the SDK SHALL use the documented `instrument_id` member, because the documented contract is the one the SDK publishes.

#### Scenario: Both spellings present with different values

- **WHEN** a tick response carries `instrument_id` "a" and `instrumentId` "b"
- **THEN** `StockTicks.InstrumentID` is "a"

### Requirement: Serialization is unchanged

The SDK SHALL encode `StockTicks` with the member name `instrument_id`, so the wire form a caller observes is the one already documented and no existing consumer sees a renamed member.

#### Scenario: Marshal after decoding the camelCase spelling

- **WHEN** a `StockTicks` decoded from `instrumentId` is marshalled
- **THEN** the encoded object carries `instrument_id` and not `instrumentId`

### Requirement: The public API does not change

The change SHALL NOT alter the `StockTicks` field set, field types, method signatures, or exported identifiers, and SHALL NOT require a caller to opt in.

#### Scenario: Existing callers compile unchanged

- **WHEN** existing code reads `StockTicks.InstrumentID` and calls `GetTick`, `GetFuturesTick`, or `GetDisplayTick`
- **THEN** it compiles and runs without modification

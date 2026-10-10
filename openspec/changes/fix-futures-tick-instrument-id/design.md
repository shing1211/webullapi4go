# Design

## Context

See `proposal.md` for motivation. The shape that matters: `data.StockTicks` (`data/tick.go:61`) carries `InstrumentID string` with `json:"instrument_id"`, and `encoding/json` matches a member name exactly, then case-insensitively. An underscore is not a case, so `instrumentId` matches neither and the field stays at its zero value with no error. Three methods return this type — `GetTick`, `GetFuturesTick`, and `GetDisplayTick`.

## Goals / Non-Goals

**Goals:**
- Populate `InstrumentID` from either spelling, with documented-name precedence.
- Leave marshaling, the exported surface, and every other field's behavior untouched.
- Keep the live conformance row and make its prose describe what it can actually see.

**Non-Goals:**
- Retagging the field. `StockTicks` is public and the retag would break anyone reading it.
- Resolving whether Webull's server-wide convention is camelCase. That needs a second host and is out of scope here.
- Touching the conformance harness's checks. The name check compares tag inventory, which a decoder does not change.

## Decisions

**A custom `UnmarshalJSON` on `*StockTicks`, using a type alias plus one shadow field.**

The method decodes into an alias of `StockTicks` — a distinct type with no methods, so it cannot recurse — while a second struct field carries the camelCase spelling. The alias keeps `Tick` and `money.Money` decoding on the stdlib path, so this change cannot alter how a price or a tick decodes. After decoding, the field is copied only when the documented spelling is absent.

Alternatives considered:

- *Retag the field to `instrumentId`.* Rejected: `StockTicks` is a public type, so this is a breaking change for anyone who has read the member, and it would abandon responses that send the documented name.
- *Add a second exported field such as `InstrumentIDCamel`.* Rejected: it grows the public API permanently to carry information a caller never needs, and two exported fields that must be kept in agreement is a worse contract than one.
- *Wrap decoding only in `GetFuturesTick`.* Rejected: the type is shared by three methods, so a per-endpoint wrapper leaves the same silent zero reachable through `GetTick` and `GetDisplayTick` and adds three call sites that can drift apart.

**The live conformance row stays, with rewritten prose.**

The `missing-required-name` check compares the SDK's tag inventory against the names a body carries. A custom decoder changes neither the tag nor the fixture, so the gate will keep reporting `data.GetFuturesTick|missing-required-name|instrument_id` after the fix. Removing the row would turn the gate red; the alternative, teaching the harness that a type accepts extra names, is a larger change to a safety check and is not warranted by one endpoint. What changes is the row's claim: it currently says a caller receives `""`, which becomes false, so its reason and unblock are rewritten to say the SDK now tolerates both spellings and the remaining question is only whether the camelCase spelling is endpoint-specific.

## Risks / Trade-offs

- **The prose in five committed files asserts the SDK is broken and will become false** (`conformance/livereasons_test.go`, `conformance/live-divergences.json`, `conformance/live_test.go`, `conformance/doc.go`, and the two status documents) → each is updated in the same change, so no file states a defect that no longer exists.
- **The three methods sharing `StockTicks` all gain the tolerance** → deliberate and additive; accepting a valid server spelling on the other two cannot break a caller that sends or reads the documented one.
- **A hand-written decoder is code the stdlib no longer owns** → `data/fuzz_test.go:159` already fuzzes `data.StockTicks` decoding, which is the coverage this risks.
- **The live count stays 13, so a reader could take a fixed defect for a live one** → the item-28 prose states plainly that the SDK now tolerates both spellings and the row records a server observation, not an SDK defect.

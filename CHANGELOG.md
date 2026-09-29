# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
version labels follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
The `v2.x` entries below preserve repository Git-tag facts. The module stays on
the v1 import path by decision, so the module proxy serves only the `v1.x` line
and these tags are not published Go-semver v2 modules; `v1.1.1` remains the
newest installable version.

## [Unreleased]

## [2.1.38] - 2026-09-30

**Release records, and the OpenSpec assistant scaffolding.** No SDK code changed
and nothing is live-verified. This tag exists because `v2.1.37` was cut before its
own release records were committed, which left the tag's own tree failing
`TestReleaseRecordsAreComplete`. Rather than rewrite a published tag, the records
land here.

### Added

- **The OpenSpec assistant scaffolding for all six tools.** `openspec init` wrote a
  skills-and-commands layer into `.opencode/`, `.claude/`, `.cursor/`, `.agents/`,
  `.github/prompts/` and `.github/skills/`, and it was untracked, so a fresh clone
  received the change workflow in `openspec/` with no way to drive it. All 55 files
  are CLI-generated (`generatedBy: 1.13.2`) and regenerate with
  `openspec init --tools ...`. The six tools differ only in invocation, which each
  directory records: `/opsx-propose` for OpenCode, Cursor and GitHub Copilot,
  `/opsx:propose` for Claude Code, and the `openspec-propose` skill for Codex, which
  gets no commands because the CLI resolves it from the skills directory.

### Fixed

- **The `v2.1.37` release records.** `TestReleaseRecordsAreComplete` reports that a
  tag existing which the status-document headers do not name is exactly the state
  the gate exists to prevent, and that state is invisible to build, vet, lint, the
  docs build and the citations check — none of which read release records. Both
  headers now name `v2.1.38` and `docs/runs/index.md` carries a row for `v2.1.37`,
  which is the record its own tag was missing.

### Note for anyone reading this tag

`v2.1.37` itself still fails `TestReleaseRecordsAreComplete` when checked out. That
is left in place deliberately: the tag is a published object and moving it would
rewrite it on two remotes. The failure is a record gap in a historical snapshot,
not a defect in the SDK, and it is fixed for every tree from this tag onward.

## [2.1.37] - 2026-09-30

**The demonstrated defect from v2.1.36 is fixed, and two release gates were
repaired for Windows.** Patch release: no exported symbol changed signature, no
existing behaviour was removed, and nothing is live-verified.

### Fixed

- **`data.StockTicks` now decodes either instrument-identifier spelling.** The 2026-09-29
  live-evidence run (v2.1.36) found the one demonstrated SDK defect in the tree: the
  sandbox sends `instrumentId` where the page documents `instrument_id`, and
  `encoding/json` matches a member name exactly and then case-insensitively — an
  underscore is not a case — so `GetFuturesTick` returned a `StockTicks` whose
  `InstrumentID` was `""`, with no error and nothing to indicate a value was missed.
  `StockTicks` gained an `UnmarshalJSON` that accepts either name and prefers the
  documented one when a body carries both. `data/tick.go`.
- **This is additive, and the public wire name is unchanged.** The `InstrumentID`
  field, its `json:"instrument_id"` tag, the marshalled form, and the signatures of
  all three methods that return the type (`GetTick`, `GetFuturesTick`,
  `GetDisplayTick`) are exactly as before, so the fix is not a breaking change and
  requires no action from a caller. The tolerance is deliberately wider than the
  single endpoint it was found on: a decoder that read one name for futures and
  another for stocks would leave the same silent zero reachable through two other
  methods.
- **The conformance gates no longer fail for a reason unrelated to the SDK.**
  `text=auto` normalized these files on commit but still wrote CRLF into the working
  tree wherever `core.autocrlf` is true, so re-encoding a committed fixture produced
  different bytes than the file held — 55 live fixtures plus the two baseline files,
  failing `TestEncodeSkeletonIsIdempotent`, `TestExpandSkeletonIsTheInverseOfEncode`,
  `TestLiveBaselineIsCanonical` and `TestLiveDivergenceBaseline`. The fixture trees
  and the two divergence baselines are now marked `-text`, so the bytes on disk are
  the bytes in the blob on every platform.
- **`make lint` passes on Windows.** gofmt rejects a CRLF working-tree copy, and the
  same `core.autocrlf` behaviour was putting CRLF into 157 Go files. golangci-lint
  reports at most three issues per run by default, so this presented as three files
  rather than 157. One `*.go text eol=lf` rule replaces the per-file listing that
  would otherwise have had to be grown one checkout at a time. No source content
  changed: the blobs already held LF.

### Changed

- **The live-conformance row for `data.GetFuturesTick` is retained, with rewritten
  prose.** `conformance/live-divergences.json` still records
  `data.GetFuturesTick|missing-required-name|instrument_id` and the live-only count
  is still 13. That is not an oversight: the name check compares the SDK's *tag
  inventory* against the names a body carries, and a custom decoder changes neither
  the tag nor the fixture, so the row still reproduces. What the row now records is
  that the server sends a name the documentation does not describe — a claim about
  Webull rather than about this SDK — and its `reason`, `unblock`, and the
  surrounding prose in `conformance/doc.go`, `conformance/live.go`, and the status
  documents were rewritten to say so. Removing the row would have turned the gate
  red for a reason that is not a defect. A second host remains the only open
  question: it would establish whether the camelCase spelling is endpoint-specific
  or this environment's convention, which is a question about the server and not
  something this SDK can answer.
- **Item 19's status text no longer reads as pending work.** The document recommended
  restoring an explicit unverified marker for `data.GetDisplaySnapshot`; commit
  `3478a6f` had already done it two days earlier, and the sentence survived the
  commit that implemented it. The marker is in place and is now recorded as done.
- **OpenSpec is now the project's change workflow.** `openspec/` holds the first
  change, `fix-futures-tick-instrument-id`, with its proposal, spec delta, design and
  tasks, and `openspec validate --strict` passes.

## [2.1.36] - 2026-09-29

**First live evidence in the repository.** Until this release the conformance
harness could only compare the SDK against fixtures generated from Webull's
published OpenAPI JSON, so it could say whether the SDK matched the
*documentation* and never whether it matched the *server*. This release calls
the sandbox, reduces each response to a value-free shape, and compares those
shapes against the SDK type. No SDK code changed and no existing exported symbol
was removed, so this is a patch release; the one demonstrated defect is recorded
rather than fixed, because fixing it changes a public DTO's wire name.

### Added

- **`examples/live-probe`**: a reachability census and a value-free response-shape
  capture. One authorised walk of all 193 documented endpoints against the HK
  sandbox recorded one outcome per endpoint, and the 55 that answered HTTP 200 were
  reduced to member-name skeletons -- every leaf a fixed placeholder the reduction
  chose -- and committed under `conformance/testdata/live/`. The census records
  73 `404`, 13 `417`, 9 `500`, 7 `403` and 1 `504` across the rest, 34 mutating
  endpoints deliberately not called, and 1 blocked because no JSON request body is
  documented for it.
- **`conformance.CompareAllLive`** and the **third divergence set**,
  `conformance/live-divergences.json`. Each captured shape is compared against the
  SDK type *and* against the documentation, which is the class a documented
  comparison is structurally blind to: a page and a fixture made to agree with each
  other report an endpoint as clean whether or not either matches the server. The
  set is separate from `known-divergences.json` on purpose, because a live finding
  is a claim about Webull's server rather than about this SDK. **Of 55 probed, 2
  disagree with their documentation and 12 disagree with the SDK** -- 14 directional
  observations over 13 distinct findings. The documented baseline is unchanged.
- **`TestLiveReportCountsBothDirections`**: the phase-4 report gate. It renders
  "of N probed, X disagree with their documentation and Y disagree with the SDK"
  from the run, reads the numbers back out of the rendered text, and cross-checks
  each direction against the raw rows, so a report and the set it describes cannot
  drift apart.

### The findings, and what they are not

- **One demonstrated SDK defect: `data.GetFuturesTick` decodes `instrument_id` to
  `""` on every response.** The page requires `instrument_id` and `data.StockTicks`
  tags the field `instrument_id`, so the documented comparison is green. The sandbox
  sends `instrumentId`, and `encoding/json` matches neither exactly nor
  case-insensitively because the underscore is not a case. **Not fixed here**: a
  public DTO's wire name is a breaking change, and one host cannot establish
  whether the camelCase spelling is endpoint-specific. Unblock is one probe of the
  same endpoint against a second host.
- **The other 12 are observations, not defects.** 4 decode rejections whose
  direction is undecided, 3 rows of one container disagreement seen from three
  angles, 3 absence-of-evidence rows against an empty account and an equity
  position, and the 2 rows of the one documentation divergence the run settled
  (`data.GetCapitalFlow`, where the server and the SDK agree and the page is the
  outlier). Written up as item 28 in `IMPLEMENTATION_STATUS.md` and
  `docs/implementation-status.md`.
- **2 of the 25 open documentation rows could be settled; 23 could not**, because
  21 sit on endpoints the HK sandbox does not serve and 2 on the one endpoint it
  answered `504`. Until this run that was an assertion; it is now measured. Every
  **per-area** figure behind that split is re-derivable from the committed
  `conformance/testdata/live-manifest.json`, which carries a `census.area` block of
  per-area statuses; the **per-endpoint** attribution is not, because the statuses
  of the endpoints that were never captured live only in the gitignored
  `examples/live-probe/census.json`. A reader with only the committed tree can
  check the areas and not the rows.
- **Nothing here is production evidence.** One run, one host, one day, and the
  census is not reproducible from a clean checkout: request schemas exist only in
  the gitignored docgen cache.

No `.go` file outside `conformance/` and `examples/live-probe/` changed, and no SDK
behaviour changed at all.

## [2.1.35] - 2026-09-28

**Breaking release. Fourteen methods could not read a conforming response at all, and
none of them could page.** Each of the fourteen documented its 200 body as
`{"data": [...], "pagination_key": "..."}` and decoded a bare slice, so
`encoding/json` reported `cannot unmarshal object into Go value of type []T` and the
call failed. That is a loud failure rather than a silent zero, so it is not the class
of defect the response-contract gate was built to catch quietly -- but it means these
methods had never succeeded against a conforming server, and `pagination_key` was
unreachable while twelve of the fourteen could not even send one. Root baseline 81 to
**25 rows over 12 symbols**, and the `missing-declared-name` class is now **empty** for
the first time: it held 104 rows when opened, 28 after v2.1.34, and its last 28 were
these same fourteen wrapper names. No fixture changed. Nothing is live-verified.

### Migration

Every one of the fourteen follows the same shape, so the change is one line each:

| | before | after |
|---|---|---|
| rows | `out[i]` | `out.Data[i]` |
| rows length | `len(out)` | `len(out.Data)` |
| next page | impossible | pass `out.PaginationKey` |

- **`data.GetMarketSectors`** gains a trailing `paginationKey string` parameter:
  `GetMarketSectors(ctx)` becomes `GetMarketSectors(ctx, "")`.
- **The other six `data` methods** take a cursor on their existing query struct, which
  gained a `PaginationKey` field: `GetDisplayGainersLosers`, `GetDisplayTopActive`,
  `GetMarketSectorDetail`, `GetFundDividends`, `GetEventContractMarkets`,
  `GetEventContractSeries`. No call site changes; the result does.
- **The seven `brokerfd` methods** take a trailing `paginationKey string`, so both the
  call and the result change: `ListFDAccounts(ctx)` becomes `ListFDAccounts(ctx, "")`,
  and `ListFDTransfers`, `GetFDActivities`, `GetFDECInstruments`, `GetFDOpenOrders`,
  `GetFDOrderHistory` and `GetFDStockInstruments` gain the parameter alongside the one
  they already take.
- **`types.Page[T]` is the new shared return type**, in `pkg/types`. One generic
  covers all eleven distinct element types instead of fourteen near-identical result
  structs, and it is the shape `GetStockProfilesV3` and `GetCorporateActions` already
  modelled locally. A Page is one response, not the whole collection: send
  `PaginationKey` back to get the next page and stop when it comes back empty. There
  is no total count and no page number on the wire.

### Fixed - tests

- **Fourteen tests encoded the SDK's own bare slice and decoded it.** That is why
  fourteen broken methods shipped with a green suite: a round-trip is green whether the
  shape is right or wrong. Each now serves the documented envelope, and each asserts
  the cursor arrives, because reading a cursor is only half of pagination and sending
  it is the half that fails silently. That second half is not hypothetical -- the first
  pass of this change built the query in `GetMarketSectors`, wrote the cursor into it,
  and then passed `nil` to `c.get`, and the only thing that caught it was a compile
  error about an undefined variable. `TestEnvelopeCursorIsSent`,
  `TestFDCursorIsSent`, `TestEnvelopeOmitCursorOnFirstPage` and
  `TestEnvelopeRoundTripsAcrossTwoPages` now cover all fourteen, in both directions.

### Noted

- **The `DeclaredNameCoverage` bite test has been rebuilt rather than repointed.** It
  had outlived four subjects across four releases -- `GetFDCorporateActions`,
  `GetFundInfo`, `GetMarketSectors`, `GetMarketSectorDetail` -- each time because the
  release drained the class it was pinned to. A bite coupled to the class under repair
  expires exactly when the repair succeeds, which is the wrong moment. It now
  manufactures its condition the way `LeafType` already did: the missing-name state is
  created by stripping a tag from a type that otherwise matches its page, so the bite
  holds whatever the baseline says. It also asserts a repaired row reports nothing,
  which is a stronger claim than the one it replaced.
- **One new test helper, `withPageElem`.** The decode target is now a `types.Page`
  struct wrapping a slice, and `withField` only descends into a slice, so a tag inside
  `Data` could not be reached for the bite to strip.
- **11 container-kind rows and 11 decode rows remain**, over 10 symbols, and none of
  them is the envelope family. Five sit on a page whose name ends `-list` and four of
  those document a single item's own fields at the top level of a list path, which is
  implausible and is the strongest documentation-error candidate in the set. They stay
  recorded: the harness has no basis to prefer the page over the type, and a wrapper
  there would be a decision rather than a reading.
- `brokerfd.GetFDOrderDetail` keeps its 2 required-name rows, and
  `brokerfd.ListAccountForms` keeps its free-form row, both recorded deliberately.

## [2.1.34] - 2026-09-28

**The evidence here is weaker than in v2.1.33, and the release leads with that.** Every
page behind these 76 names declares the property but publishes **no `required` list**, so
a name may be optional, conditionally sent, or absent from the example the page chose.
Each field added here exists because the page publishes it, not because a server must
send it. That is the whole difference from the required-name work, and each field's GoDoc
says so where a reader of the type will actually see it.

76 declared names across 17 response types now reach a field. Each one decoded to the
zero value with no error reported, so a caller reading a name the page publishes got an
empty string, a zero or an empty slice. Root baseline 157 to 81 rows over 26 symbols.
Required-name is unchanged at 2 rows over 1 symbol. No fixture changed. Nothing is
live-verified.

### Fixed - not breaking

- **`brokerfd.FDCorporateAction` gains 12**, the largest single batch, including `from`
  and `to` as new types: the two sides of the event, where `to` is an array because a
  split can have several targets. Without them a caller cannot tell a split from a
  rename.
- **`brokerfd.FDEnum` gains 3**; `code` and `name` are the same facts it already spells
  `Value` and `Label`, so both sets are carried.
- **`brokerfd.FDOrder` gains `client_order_id`**, which the place and replace pages
  declare and no Broker FD order type carried.
- **`data.FinancialAlert` gains 9**, `data.FundInfo` 6, `data.ScreenerStock` 12 across
  three pages, `data.DSNewsSummaryItem` 10 across four, `data.Logo` 2,
  `data.FundNav` 2, `data.FinancialIndicator` 2, and `data.StockInstrument` 2.

### Type decisions taken from the page, not the fixture

- **Free-form objects are carried untyped.** `FinancialIndicator.Values`,
  `DSNewsSummaryItem.Args`, `.Headers` and `.Rows`, and
  `FDCorporateActionTarget.Payouts` are objects the page declares without naming any
  property inside, so there is nothing to type them as. They are `map[string]any`, and
  each says why on the field. A typed struct here would be invented, not read.
- **`is_incumbent` is an `int64`, not a bool.** The page documents it as an integer, so
  a caller reads a number. `FundInfoManager` also carries `tenure_days` as an integer.
- **`id` on the news pages is an `int64`.** The values are 17 digits wide, so a
  string-backed field would have dropped them.
- **`is_adr` is a `string` on a page whose example is `"false"`.** The page sends a
  string, so a `bool` field would have rejected it. `financial-alerts` `fiscal_year`
  and `fiscal_period` are integers; the four money figures beside them are decimal
  strings, matching how Webull sends money.
- **The sector-detail counts are `string`.** `advanced`, `declined` and `flat` are counts,
  but the page types them as strings, so they are carried as strings rather than
  widened into an int the wire does not send.

### Noted

- **The `DeclaredNameCoverage` bite test has now outlived two of its subjects.** It was
  built on `brokerfd.GetFDCorporateActions`, which v2.1.33 fixed, then on
  `data.GetFundInfo`, which this release fixed. A bite that names a fixed defect asserts
  the defect is still there, so it moved to `data.GetMarketSectorDetail`, whose
  `data`/`pagination_key` rows still reproduce. That is the **fourth** bite test to need
  this treatment: `TopLevelShape` moved in v2.1.33 and `LeafType` already manufactures
  its condition. Three of the four moved because their subject was fixed, which is a
  sign the subjects were chosen from a class under active repair rather than from one
  that is stable.
- **The declared-name class is now homogeneous, and that is the finding.** 76 of its 104
  rows are fixed, and all 28 that remain are `data` or `pagination_key` on the 14 methods
  that return a bare slice where the page documents an envelope. The name check has
  stopped finding missing fields and started reporting the same missing wrapper 28 times.
  **That is the next piece of work and it is a functional gap, not a readability one:**
  `pagination_key` is unreachable and 12 of the 14 methods cannot send one either, so
  these endpoints cannot be paged at all. It is 14 public signatures, so it is a release
  of its own rather than a second half of this one.
- **10 of the 25 shape rows stay recorded as judgement calls, not schema readings.** A
  `.../list` path whose documented 200 is a single object (`GetDSLatestNews`, `GetLogos`,
  `GetStockProfilesV3`) is a documentation question, and a wrapper would be a decision
  rather than a reading. `brokerfd.ListAccountForms` stays recorded because its page
  publishes no properties at all, so there is nothing to compare.

## [2.1.33] - 2026-09-28

**Breaking release.** 62 recorded response-contract divergences are resolved: 58
required names across 18 Broker FD symbols, 4 on the two Display Solution bars
types, and one breaking envelope wrapper. The required-name class falls from 62 rows
over 20 symbols to **2 rows over 1 symbol**, which is the whole of what remains and
is now left recorded deliberately rather than fixed. Root baseline 219 to 157 rows
over 35 symbols. No fixture changed. Nothing is live-verified.

### Fixed - not breaking

- **58 required names across 18 Broker FD response types reached no field.** Each one
  decoded to the zero value with no error reported, so a caller following the
  documentation read an empty string, a zero or a false. `BankAccount` and
  `ACHAccount` gained their relationship, account-name, account-number, routing and
  code fields; `CreditInfo` gained all seven it was missing, including the `account`
  and `amount` that say which account a credit applied to and how much; `Transfer`
  gained `transfer_type` and `direction`; `TransferFee` gained `fee_id`, `account`,
  `contra_account` and `status`; `InstantFunding` gained `instant_funding_id`,
  `currency` and `type`; `Agreement` gained `id`, `name` and `content_type`;
  `FDCashJournal` gained `from_account` and `to_account`; `FDTradeCalendarEntry` gained
  `is_trading_day` and `is_settlement_day`; `FDAccount` gained `application_id`; and
  `FDOrderPreview` gained `estimated_cost` and `estimated_transaction_fee`, which
  splits a quoted total into base cost and brokerage.
- **`StockBars` and `BatchBars` gained `times` and `special_times`.** The Display
  Solution bars pages require both, and without them the display host's trading
  calendar was unreadable, so a caller could not tell why a session was missing from a
  bar series.
- **Dual spellings are carried side by side, as on `FDPosition`.** Several of these
  types have always carried an SDK spelling for a fact the page spells differently
  (`bank_id` against `bank_relationship_id`, `funding_id` against
  `instant_funding_id`, `type` against `transfer_type`, `agreement_id` against `id`).
  Which of the two a live server sends is unverified, so both are carried and a
  response populates whichever it carries. Nothing here asserts which a live server
  sends.

### Fixed - breaking

- **`GetFDAssetsDetail` now returns the documented envelope.** The page documents the
  200 body as an object with `account_currency_assets`, `total_asset_currency` and
  `total_cash_balance`, all required; the method returned the currency array alone, so
  the two totals reached no field and were dropped with no error reported. A caller
  holding a multi-currency account could see cash per currency but never the total,
  which is the figure an account is normally reconciled against. It returns
  `*FDAssetsDetail`; a caller reading the per-currency entries now reads
  `out.AccountCurrencyAssets[i]` rather than `out[i]`.

### Two distinct types where one would have to lie

- **`data.ExchangeTimes` and `data.SpecialExchangeTimes` are separate types.** The
  Display bars page requires `start`, `end` and `trading_session` on both hours arrays,
  but documents `start` and `end` as **strings** in `times` and as **integers** in
  `special_times`. One element type with a widened field would silently drop whichever
  form it could not hold, so there are two. A negative test pins the disagreement, so
  merging them later fails a test instead of quietly losing a form.

### Noted

- **The `TopLevelShape` bite test was repointed, not deleted.** It had been constructed
  on `GetFDAssetsDetail`, which this release fixes, so it was asserting a defect that
  no longer existed. The shape class still holds 25 live rows, so the bite moved to
  `GetFDTransferFees`, which inverts the same way. This is the second time a bite test
  has outlived the defect it bit on; the `LeafType` bite already manufactures its
  condition for the same reason.
- **The last 2 required-name rows are now recorded as deliberately unfixed.**
  `GetFDOrderDetail` reconciles to the *list* page, so its `combo_type` and `orders`
  are envelope fields belonging to a list response. Tagging them on `FDOrder` would
  put two envelope fields on a one-order type shared by seven other methods, and would
  make the gate green while the detail call still could not read either name. The
  recorded reason says so, because a required reason is not a required true reason and
  the previous text implied a retag was the fix.
- **The four Display rows were decidable all along.** The paid Display entitlement
  blocks calling the endpoint, not declaring a field the page requires; the two methods
  send the path the page documents, so the missing tags were a defect rather than a
  naming choice. They are the only `data` rows in the required-name class and they are
  now gone, leaving the `data` column of that class at zero.
- **Still 105 declared-name rows, all in `data`.** The required-name class is
  effectively closed; the declared-name class is not, and is where any further
  response-contract work belongs.

## [2.1.32] - 2026-09-28

**Breaking release.** Three event-contract endpoints decoded the wrong shape, and
fixing them changes public return types. Item 27 recorded them as decidable and
awaiting a decision; this is the decision, taken in favour of the documentation. The
evidence is the pages themselves, which are unambiguous where item 19's page is not.
Root baseline 239 to 219 rows over 52 symbols. No fixture changed. Nothing is
live-verified.

### Fixed - breaking

- **`data.GetEventBars` and `data.GetEventTick` decoded one nesting level too deep.**
  Both pages document the 200 body as an array of objects, each requiring
  `instrument_id`, `result` and `symbol`, where `result` is itself an array of bars or
  ticks. The SDK returned the elements of the inner array, so the grouping key was
  discarded and bars or ticks from several instruments arrived in one flat list with
  nothing to say which instrument each belonged to. They now return `[]EventBarsResult`
  and `[]EventTickResult`, one entry per instrument. **A caller reading `out.Result[i]`
  instead of `out[i]` adapts in one line; a caller that used `out[i]` as a bar must read
  `out[i].Result`.**
- **`data.GetEventDepth` returned one object where the page documents an array**, and
  carried no no-side book at all - so a caller could not read the no side of an event
  contract's depth, which is the side such an instrument is named for. It now returns
  `[]EventDepth` with `InstrumentID`, `QuoteTime`, `NoBids` and `NoAsks`. This inverts
  which shape fails to decode: previously a conforming response could not be read, and
  now a single-object response cannot be. Neither shape is verified against a live host.
- **The depth time and bar time were named `timestamp` where the pages require
  `quote_time` and `time`.** Both SDK spellings are kept alongside the documented ones,
  so a required name always reads whichever the server sends.
- **The conformance symbol table was updated to match.** It is maintained by hand and
  still pointed at the old decode targets, so the harness compared the bars and tick
  endpoints against the inner element and reported no change - the first 14 rows dropped
  and the bars and tick rows did not, which is what exposed it.

### Fixed - not breaking

- **`data.EventSnapshot` could not read 10 of the 15 names its page requires.** The type
  had `last_price` where the page says `price`, and no size companion for any side, so a
  caller could not read the depth a snapshot exists to report. The documented names are
  carried alongside the existing seven.

### Noted

- **Why this is a decision taken rather than a defect merely recorded.** The four pages
  are not self-contradictory - unlike the Display Solution page in item 19, whose
  `operationId` says GET while its method says POST - so on the available evidence the
  SDK is wrong in all three methods. Item 26 had deferred these as needing a live probe
  on the grounds that an envelope presents object-against-object; reading the schemas
  showed the documented shape is unambiguous and only the obstacle was the break.
- **The break is larger than the two field-level ones in v2.1.28.** Those changed which
  fields a type carried; this changes what a method returns. A caller following the
  CHANGELOG's `out.Result[i]` note adapts without reading the method.
- The `missing-required-name` class is now 62 rows over 20 symbols, down from 119 over
  28 when the sizing work began. The `trade` column remains empty.
- No fixture or manifest byte changed.

## [2.1.31] - 2026-09-28

Corrective release. 11 more required-name rows are closed by additive fields on three
response types, and the `trade` package's column in the class table is now empty. It
also corrects item 26's own conclusion: the envelope cases it deferred as needing a
live probe are decidable from the page schema, and four of them are recorded with
their evidence and the decision they now need. Root baseline 250 to 239 rows over 56
symbols. No fixture changed. Nothing is live-verified.

### Fixed

- **`trade.BatchPlaceOrder` could not report how much of a batch was placed.** The page
  documents `total`, `success`, `failed` and `batch_orders` and marks all four
  required, with the per-order element carrying `error_code` and `message`. The SDK's
  type had a single `Results` field that appears on **no page at all**. So a caller
  could not learn how much of a batch succeeded, and a rejected order inside a batch was
  indistinguishable from a success that returned no identifier. All four are now
  carried, plus the element's `error_code` and `message`; `Results` is kept so nothing
  that compiles today breaks, with the type comment naming `BatchOrders` as the
  documented field to read.
- **The event-contract category and event types could not read their required values.**
  `EventContractCategory` carried `category` and `name` and none of the three required
  names the page declares. `EventContractEvent` documented `symbol` and `series_id`
  where the SDK had `event_symbol` and `series_symbol`, and had no `short_name` or
  `mutually_exclusive`. Both sets are now carried, which is the additive route rather
  than a rename, so which naming a live server sends does not have to be known first.
- **The `trade` column of the class table is now empty.** Its four rows were all
  `missing-required-name` on `BatchPlaceOrder`, so closing them leaves no `trade` row in
  that class. Stated here because a zero in a table reads as "measured and clean",
  which is a stronger claim than "nothing is recorded here yet".

### Noted

- **Item 26 was too cautious, and this release says so.** Item 26 recorded that a page
  wrapping its payload in an envelope presents object-against-object, which the shape
  check cannot see, and deferred those rows as needing a live call. Reading the page
  schemas shows the documented shape is unambiguous, so the mismatch is determinable
  offline and the only obstacle is a public API break.
- **`data.GetEventBars` and `GetEventTick` decode one nesting level too deep.** Both
  pages document `array<{instrument_id, symbol, result: array<...>}>`, requiring
  `instrument_id`, `result` and `symbol`. The SDK returns the **elements of the inner
  `result` array**, so the documented grouping key is discarded and bars from different
  instruments arrive in one flat list. The missing `result` name is the symptom of that
  depth error, not of an absent field.
- **`data.GetEventSnapshot` and `GetEventDepth` are wrong shapes, not shallower ones.**
  The snapshot page requires 15 names, including `price` where the SDK has `last_price`
  and four `*_size` companions the SDK does not have at all. The depth page requires
  `no_bids` and `no_asks`, which the SDK's type does not carry, and returns a single
  object where the page documents an array.
- **What stopped them is a decision, not a block, and the evidence is now complete.**
  Fixing any of the four changes a public return type, so callers stop compiling - a
  larger break than the two field-level ones taken in v2.1.28. Unlike the Display
  Solution page in item 19, these four pages are not self-contradictory, so on the
  available evidence the SDK is wrong in all four. That makes it a maintainer call
  between accepting the break and recording the pages as wrong, and it is not made
  here.
- No fixture or manifest byte changed.

## [2.1.30] - 2026-09-28

Corrective release, and a small one with a large lesson in it. 28 recorded
divergences are fixed by 19 additive fields, and the largest remaining class is sized
and split into work that is safe and work that would encode the wrong contract. It also
adds the gate whose absence let three separate false premises accumulate in the status
document, and finds that a fourth recorded item was already closed. Root baseline 278 to
250 rows over 59 symbols. No fixture changed. Nothing is live-verified.

### Fixed

- **28 `missing-required-name` rows removed, by 19 fields, with no probe.** The three
  names are `client_request_id` on 11 symbols, `update_time` on 11 and `create_time` on
  6, across 8 Broker FD response types. A *required* name is one a conforming server
  always sends, so adding a field cannot break a call that works today - the same
  additive argument item 20 made, and the reason it needs no credential.
  `client_request_id` already appears in the sibling `broker/` module in exactly this
  role, so the spelling follows it rather than being invented. The exact-set gate
  confirmed all 28 stopped reproducing and produced **no new divergence**, which is the
  check that the chosen field types are compatible; a wrong type would have shown up as
  a new leaf-type or decode row.
- **Item 15's premise was false, and it had been closed since v2.1.4.** It recorded
  nested coverage and a strict documentation build as Makefile/release gates rather
  than CI gates. Both are CI jobs - `docs (strict)` runs `mkdocs build --strict`,
  `nested coverage` runs `go test -cover` across the module matrix with a floor on
  `broker/`, and `citations` is a third - all three verified passing on the v2.1.29
  tag run, and the document's own version history recorded adding them in v2.1.4. What
  survives is the half that was never a gap: a measured percentage is not a behaviour
  guarantee, and a coverage floor is a floor rather than proof.

### Added

- **A gate on the status document's own numbers.** The class table is parsed out of
  both status documents and compared against the committed baselines, so a count can
  only change by changing the baseline and the table together. It exists because the
  table had to be corrected by hand in v2.1.27 - a hand-counted figure said 29 symbols
  where the baseline and the gate both said 28 - and because a hand-recomputed figure
  is a figure that will be wrong again. The gate caught the same class of error
  immediately on its first run: the 28-row fix left the table at 278 while the
  baseline was 250. It is proven to bite on a single wrong cell, a wrong total, and the
  mirror drifting from the root.

### Noted

- **The remaining 91 `missing-required-name` rows are not one kind of work, and 53 were
  deliberately left alone.** A missing required name is not always a missing field: a
  page that wraps its payload in `{"result": [...]}` against an SDK that decodes the
  inner item presents object against object, which the container-kind check cannot
  see, so every top-level required name reads as missing when the real defect is a
  depth mismatch. `data.GetEventBars` and `GetEventTick` are missing `result` while
  `data.StockBars` has one; `trade.BatchPlaceOrder` is missing `batch_orders`, `total`,
  `success` and `failed`, which is a batch-result envelope rather than four attributes
  of an order. Adding fields for those would encode a contract the page does not
  describe.
- **A measurable warning sign, which is the part worth keeping.** 4 of the 28 symbols
  carrying required-name rows also record a container-kind mismatch -
  `GetAgreementDetail`, `GetFDAssetsDetail`, `GetFDTransferFees` and
  `data.GetEventDepth` - so their 17 rows are secondary symptoms and must not be fixed
  by adding fields. The other 24 record no mismatch, which is necessary but **not
  sufficient**: the envelope case is invisible to the shape check, so absence of a
  mismatch is not evidence that a name belongs on the object. The class should be
  worked one symbol at a time, not in bulk.
- **A second blind spot this exposed.** The shape check compares container *kind*, so
  it cannot see a depth mismatch when both sides are objects. That is a real limit of
  the instrument. Only the 4 correlated symbols show it in the committed baseline,
  which is too few to size the class, and a depth-scoped check would report names the
  SDK does decode - which is why `CheckNameDepth` exists as a check that does not run.
- **Three false premises in one document, in a row.** Item 12's counts, item 19's
  "differs from both official sources", and now item 15's CI gaps were each wrong, and
  each would have sent a reader down a path the evidence does not support. The table
  gate covers the numbers; the prose claims are outside its scope and remain a manual
  read, which is stated here so the gate is not mistaken for more than it is.
- No fixture or manifest byte changed.

## [2.1.29] - 2026-09-28

Documentation release. It changes no SDK request and no fixture. It completes the
investigation item 17 had been holding open, and it corrects a factual error that had
propagated into item 19, the next-steps list and `AGENTS.md`: the SDK's `/openapi/*`
paths were recorded as disagreeing with **both** official Webull sources, and they
disagree with only one.

### Fixed

- **The four "summary-only" rows are documentation drift, not SDK defects.** All four
  are the `/openapi/*` namespace, and in each the SDK path equals the **llms.txt
  summary** - which is one of the two official machine-readable sources, not a
  third-party index - while the OpenAPI JSON documents a reorganised path:
  `client.CreateToken` (`/openapi/auth/token/create` against `POST
  /auth/tokens/create`), `client.CheckToken`, `data.GetStockInstruments` and
  `data.GetDisplaySnapshot`. The two Webull sources disagree with each other and the
  SDK sides with one of them, which makes its path *supported* rather than *proved*:
  no live call in this session tested either spelling.
- **The next step that proposed aligning them was withdrawn, because taken literally
  it would have broken authentication.** It asked to "resolve the four summary-only
  matches ... through the doc generator and official OpenAPI", which means moving
  `client.CreateToken` off the path that every authenticated call in this SDK depends
  on. The OpenAPI JSON spelling has never been exercised by this code. A gate or a
  checklist that proposes this is worse than one that says nothing, so the step is
  withdrawn in both status documents rather than left to be executed.
- **`data.GetDisplaySnapshot`'s recorded premise was wrong.** Item 19 read that the
  method "differs from both official sources". It does not: the llms.txt summary for
  that endpoint is `/openapi/market-data/stock/snapshot`, exactly the path the SDK
  sends. That weakens the case for the large breaking change rather than strengthening
  it, because the path is now the part an official source actively supports, while the
  request-shape problems - a required JSON body with a `category_symbols` array against
  an SDK that sends flat query parameters - are untouched and still need the paid
  entitlement.
- **The `/openapi/*` sibling alignment, which the same record relied on, now looks
  like a coincidence.** The item cited the four sibling Display constants having been
  aligned to `/market-data/stocks/*` in commit `2b29c88`, leaving this one as a
  "holdout". If the OpenAPI JSON is the reorganised source, that commit aligned four
  endpoints *away* from a spelling Webull still documents, and this one was left alone
  correctly. That is a second, larger question this release records rather than
  answers: it needs a live call per endpoint, and no entitlement is available.

### Noted

- **All 14 `/broker-fd/*` path literals are now classified, and none is alignable.**
  The recorded judgement that "4 can be aligned" was optimistic; against the docgen
  cache the answer is 0. The cache holds 50 Broker FD pages and every one is under
  `/broker/...`, so none of the 14 appears in the documentation at all. They split into
  **6 contested** (two SDK symbols plausibly own the one documented page, so it cannot
  be assigned to either without a decision the docs do not make), **6 with no
  documented endpoint at all**, and **1 ambiguous** with two candidate pages. The
  contested pairs are `GetAccountsSummary`/`GetFDAssetsSummary`,
  `GetFDPositions`/`GetPositions`, `GetFDCorporateActions`/
  `GetFDCorporateActionDetail`, `GetFDCashJournalDetail`/`ListFDCashJournals` and
  `CreateFDAccount`/`SubmitAccountForm`.
- **Six of the 14 can never be resolved by a credential.** `ListDocuments`,
  `GetDocumentDetail`, `GetFDAchAccountDetail`, `GetFDBankAccountDetail`,
  `GetFDStockLocate` and `GetFDECInstrumentDetail` have no documented counterpart of
  any kind, because Webull publishes no page for them. No probe can return a path no
  page describes, so these are "unanswerable from the documentation" rather than
  "blocked on a credential" - a different category, and not one to wait on.
- **A near miss, aligned and then reverted, is recorded so it is not re-attempted.**
  `GetAccountFormDetail` looked like the one clear case: the cached Form Content page
  declares `GET /broker/forms/get` and is the only single-form fetch. It was aligned,
  the fixture regenerated, and the change reverted. That endpoint "retrieves the JSON
  schema for the specified form code and version" - it takes `form_code` and `version`
  as query parameters and its 200 body is a JSON Schema fragment whose properties are
  the schema keywords `required`, `type`, `format`, `min_items`, `max_items`,
  `max_length`, `enum_values`, `description` and `example`, while the method fetches a
  form *instance* by `form_id` and decodes `brokerfd.AccountForm`. The harness caught
  it: it reported `form_code` and `version` as required response names that no
  `AccountForm` field carries, which is a true statement about the page and a false one
  about the method. Matching endpoints on a shared word paired different *kinds* of
  endpoint, and the same trap applies to every remaining candidate. The reasoning now
  sits next to the three form path constants in `brokerfd/accounts.go` and in the
  docgen manifest's Form Content row.
- No fixture, manifest or SDK request changed. The committed conformance tree is
  byte-identical to `v2.1.28`, which the fixture self-check confirms. Nothing is
  live-verified.

## [2.1.28] - 2026-09-28

Corrective release. It closes three of the four remaining live-blocked defects, two
of them without the credential those items said was required, because both were
facts about the code rather than about the server. It adds the two harness checks
whose absence let those defects hide, and it fixes a response type that could not
decode half of Webull's own documentation. Root baseline 278 rows over 60 symbols,
broker 59 over 20; no fixture changed; nothing is live-verified.

### Fixed

- **`brokerfd` sent every request to the core host.** `brokerfd` called `c.core.Do`
  while its sibling `broker/` module, one directory over, called `c.core.DoBroker`
  for the same class of request. The transport, the `client.Endpoints.BrokerHTTP`
  field and the sibling's use of it all already existed - the two packages simply
  disagreed. The status recorded a US Broker FD credential as the unblock; that was
  needed to observe the 404, not to establish the cause. `DoBroker` is now used, and
  `TestEndpointsForAlwaysSuppliesBrokerHost` asserts that every region and environment
  supplies the field, because an empty one would make the package unusable rather
  than merely wrong.
- **The routing fix exposed a live-traffic hazard and it is closed.**
  `client.WithBaseURL` overrides one field and sets the internal override flag, so
  `client.New` leaves every other endpoint at its default - the *production* Broker
  host. 54 call sites in `brokerfd`'s tests used it, so correcting the routing turned
  the whole suite into calls against a production API. All 54 now use
  `client.WithEndpoints` with both hosts, as `broker/`'s tests already did.
  `client/client_test.go` asserts the `WithBaseURL` behaviour so it cannot change
  silently, because the failure mode is an offline suite reaching production rather
  than a failing assertion.
- **`broker.UpdateVirtualAccount` sent the wrong verb, the wrong body, and a field
  that does not exist.** The status recorded "the wrong verb and body shape". The
  official page says more: the verb is POST, `account_id` is a **body** field, the
  request takes no query parameters beyond the auth headers, and two body fields are
  required - `account_id` and `client_request_id`. The old request type carried one
  field, `AccountName`, which the page does not declare at all, so the request was
  wrong in every respect at once. It now issues POST with the documented body, and
  `UpdateVirtualAccountRequest` carries `ClientRequestID`, `TradingPermissions`,
  `OptionLevel`, `CommissionCode`, `W8BENInfo` and `ChinaConnectInvestorInfo`.
  **Source-level break, taken deliberately:** `AccountName` is gone rather than kept
  and ignored, because a field that was never in the contract and is silently dropped
  is the worst of the three options - the caller compiles, sets it, and learns
  nothing. The unexported `put` helper had this as its only call site and went with it.
- **`data.Quote.QuoteTime` could not decode half of Webull's own documentation, and
  this was a total failure rather than a silent zero.** Four documented pages carry a
  `quote_time` property and Webull's own published examples disagree about its type:
  the stock and Display Solution depth pages show a quoted string
  (`"1640688000000"`), the futures and event-contract depth pages show a bare number
  (`1761131409276`). The field was `int64`, so it decoded the two numeric pages and
  **failed the other two with a decode error that abandoned the whole response** - not
  one field, the entire body. `data.QuoteTime` reads a quoted decimal, a bare number,
  an exponent form, null and an omitted field, and reports anything else through one
  `ErrQuoteTimeFormat` sentinel, so it is correct whichever the server sends. A second
  field carrying the same name was not an alternative: two fields claiming one json
  tag is the defect the new duplicate-tag scan exists to catch. Baselines: 4 rows
  removed, root 282 to 278 and 62 to 60 symbols, and the `leaf-type-mismatch` class is
  now **empty**.
- **The additive route that closed item 20 does not apply to a type conflict, and the
  difference is worth recording.** In item 20 the documentation and the SDK used
  *different* wire names for one value, so a twin field was possible. Here they use the
  *same* name with different types, so a twin cannot populate from it - the fix has to
  be a type that accepts both.

### Added

- **A verb check, because all five pre-existing checks compare the response and none
  looked at the request.** `broker.UpdateVirtualAccount` reconciled as a clean `match`
  for the whole life of the harness while issuing PUT where the page documents POST -
  the sharpest single illustration of the limit the status states about a green row.
  `conformance.CompareVerb` reads the verb the method actually sends **from the Go
  source** and compares it with `documented.method`, which the generator already
  recorded for all 193 fixtures. Reading it from source rather than a hand-maintained
  table follows the precedent `conformance/envelopes.go` set for unexported types. The
  extractor follows a method, a sub-service helper, an `http.MethodX` passed as an
  argument, a verb written as a string literal, and same-package delegation - which it
  must, because `trade.GetOpenOrders` delegates to `GetOpenOrdersPage` and
  `display.Service.EnsureToken` to `fetchToken`. On the committed surface: **183 mapped
  symbols, all readable; 178 rows agree, 0 disagree, 15 not comparable on path.** Zero
  is honest only because the fix above landed first; reverting one line produces
  exactly one finding naming it with both sides. Wired into the `broker/` gate too.
- **A whole-module duplicate-`json`-tag scan, because the class is invisible to a
  name-lookup harness.** `brokerfd.FDPosition` tagged both `UnrealizedPL` and
  `RealizedPL` `json:"unrealized_pl"`, and `encoding/json` drops **both** fields when
  two fields of one struct claim the same name - with no error reported - so open P&L
  had been silently zero for every caller. It survived 285 rows of conformance work
  because the harness looks for names that are *missing* and nothing was missing:
  `tagsOf` keeps one of the two colliding names and reported the name as covered while
  the decoder filled neither field. The two defects look identical from outside and
  only the level at which you look tells them apart.
  `conformance.ScanDuplicateTags` walks **370 structs across 6 module directories**,
  matching on the name before the comma so `json:"a"` and `json:"a,omitempty"` are
  correctly a collision, and excluding a nested module from the parent walk by its own
  `go.mod` so no struct is counted twice.
- **Two refusals that keep the new checks honest.** A verb the extractor cannot read
  yields no divergence, because that is a defect in the extractor; it is required to be
  zero by a test rather than recorded in the baseline, which would put a tool defect
  in the one file whose purpose is to record SDK defects against the documentation.
  And `client.CreateToken` and `client.CheckToken` each send several verbs, so they
  have no single counterpart verb and are reported as not comparable rather than forced
  into an answer.
- **`leafVerdict` now asks a type instead of naming it.** It listed `money.Money` as the
  one exempt type, which is a fact about the SDK that goes stale the moment a second
  such type appears - and it had to be extended the moment `data.QuoteTime` did. It
  sends the documented kind through the type's own `UnmarshalJSON` instead,
  deliberately one-sided: a type that accepts the sample becomes positively confirmed,
  one that rejects it stays not judgeable rather than being reported as a mismatch. The
  change can only reduce the set of unexamined fields, never invent a finding.

### Noted

- `TestComparisonBites/LeafType` used `data.Quote.QuoteTime` as a live defect to bite on.
  That class is now empty, so the bite is manufactured instead: the subtest re-types
  `quote_time` back to a plain `int64` - exactly the field the SDK carried - and
  requires the check to name the disagreement again. A bite test that kept asserting a
  fixed defect would have gone on passing for the wrong reason, and one that deleted
  itself would have left the check unproven.
- The `data.QuoteTime` public type change means `Quote.QuoteTime` is no longer an
  `int64`; a caller assigning it to an `int64` variable needs a conversion. Comparisons
  against untyped constants are unaffected, and every use in this repository is one.
- Nothing here is live-verified. The two SDK fixes are derived from the official pages
  and from the code's own internal inconsistency, not from an observation of a server.

## [2.1.27] - 2026-09-28

Corrective release. It closes the silent part of live-blocked defect item 20 without
the credential that item says is required, lints all six modules in CI rather than
one, and fixes a defect the conformance harness was structurally unable to see. The
`brokerfd.FDPosition` and `broker.Position` changes are additive, so no call that
works today can break. No fixture changed; the two baselines lost 6 rows between
them, from defects that no longer reproduce. Nothing is live-verified.

### Fixed

- **Item 20's three silently zeroed values are no longer zeroed, and no credential
  was needed.** `GET /broker/assets/positions/list` documents `cost_price`,
  `last_price` and `unrealized_profit_loss`; both `brokerfd.FDPosition` and
  `broker.Position` carried only `average_cost`, `market_value` and `unrealized_pl`.
  The item's own minimal fix offered two routes and gated both on a probe. Only one
  of them actually needs one: **the documented names are now carried alongside the
  SDK's own** rather than replacing them, so whichever set a live server sends, the
  caller reads a populated value. A response populates whichever of the two names
  it carries, so a caller reading either name reads a value when the server sends
  it; which of the two a live server sends is precisely what stays unverified. The
  change is additive rather than a retag — a retag would have been
  breaking *and* still speculative. `brokerfd.GetFDPositions` now records **no
  divergence of any kind** for the first time, and the broker baseline falls by 3.
  Root: 285 to 282 rows, 63 to 62 symbols. Broker: 62 to 59. Combined 341 over 82.
- **A worse defect in the same struct, which the harness could not see.**
  `RealizedPL` was tagged `json:"unrealized_pl"`, duplicating `UnrealizedPL`. When
  two fields of one struct claim the same name, `encoding/json` drops **both** rather
  than picking one: verified directly, an untagged body fills neither and
  `encoding/json` reports no error. So **open P&L had been silently zero for every
  caller since the field was written**, and `realized_pl` was never decodable at all.
  Nothing here reads as *missing*, which is exactly why 285 rows of divergence
  hunting never found it: `tagsOf` collects wire names into a map and keeps one of
  the two, so the harness reported `unrealized_pl` as **covered** while the decoder
  filled neither field. A tag-lookup harness has a blind spot at the field level
  that a name-lookup harness does not, and this is it.
  `TestFDPositionHasNoDuplicateWireName` and `TestPositionHasNoDuplicateWireName`
  now pin the class by walking the struct rather than naming the pair. A scan of all
  **365** structs in the root module and the examples found this to be the only
  occurrence. Reintroducing the duplicate tag was confirmed to fail three tests and
  show both fields at `"0"`.
- **The existing test could not have caught either.** `TestGetFDPositions` marshals
  `[]FDPosition` and decodes it back, which is green on a type that reads nothing,
  and it never asserted `UnrealizedPL` or `RealizedPL`. The new tests decode literal
  bodies instead: documented names, SDK names, both at once, and each field carrying
  a shape the decoder must survive. `TestComparisonBites/RequiredNameCoverage` used
  to assert the three missing names still reproduced and had to be inverted to assert
  the fix holds, with its three mutation checks unchanged so the check itself is
  still proven to bite.
- **CI now lints all six modules, not one.** `golangci-lint-action` resolves
  `./...` through the module it is invoked in, so a single run at the repository root
  never saw `broker/` or any example module — while the `gofmt` step in the same job,
  being a plain directory walk, did. That asymmetry is why a dead three-line method
  sat in `broker/client.go` from 2026-09-21 until a local run found it. The lint job
  is now a six-entry matrix mirroring the Makefile's `MODULES`, and gofmt runs once
  from the root rather than six times. Reintroducing the dead method is confirmed to
  be caught by the `broker` entry while a root-level run still reports `0 issues`.

### Changed

- **The three fuzz targets gained real seeds.** Each had one trivial
  `{"symbol":"AAPL"}` and discarded the decode error, so `make fuzz` was a
  one-second exercise of almost nothing that could only catch a panic. 88 seed runs
  now replace 3, chosen from what this SDK actually gets wrong: money fields as both
  JSON string and number, `quote_time` as the documented string against the SDK's
  `int64` (item 21's open leaf-type defect), each type seeded with the container
  inversion its own page is ambiguous about, and the degenerate cases a careless
  server produces — null, empty, wrong-typed, truncated, non-ASCII. Go replays a
  fuzz target's seeds as ordinary tests, so every one is now a standing regression
  case in `go test ./data/`. All three targets fuzz clean at `-fuzztime=1s`, writing
  no crasher corpus.

### Noted

- **The citation gate caught the fallout of the doc comment, which is the second
  time this session it has earned its place.** Adding 14 lines above
  `FDPosition` moved the struct from line 79 to line 98, so the citation in item 20
  pointed at a range that no longer contained any code. The gate failed on exactly
  that and nothing else, and the range is corrected to `brokerfd/assets.go:98-121`.
  Nothing in a build would otherwise have noticed a `file:line` citation pointing
  into a comment.
- **`brokerfd.FDPosition` and `broker.Position` grew by 3 fields each.** Six fields
  in total, all `money.Money` in the first and `string` in the second, matching what
  those types already used. Any caller that constructs either type with a positional
  literal would break; a keyed literal is unaffected, and every use in this
  repository is keyed.

## [2.1.26] - 2026-09-28

Harness and test release. The Broker API HK surface stops being the largest
unexamined area in the repository: its 29 documented endpoints are now measured by
the same five checks as the root module, and **62 further divergences are recorded
for them**. One dead unexported method is removed. No fixture, no baseline entry in
the root module, and no SDK response type changed, and nothing is live-verified.

### Added

- **The `broker/` module is inside the wire-conformance harness.** Its 29 rows were
  out of scope because the module is separate and no symbol in the root module can
  name a broker type — the same unexported-type problem the envelope machinery
  exists to solve. `broker/conformance_test.go` resolves all 29 rows from inside the
  module and runs the identical checks: **28 compared, 2 of them against an
  unexported envelope read from source, 1 decoding no body, 0 not comparable, and 62
  divergences recorded**. Every check runs the root module's code; nothing in the new
  file reimplements one, so a change to a check reaches these rows as well as the
  root module's 285. Only the symbol table is written twice, which is unavoidable
  and is asserted in both directions so a hole fails the build rather than quietly
  reducing what is measured.
- **Two exports from the root module, both narrow.**
  `RegisterEnvelopeNamedType` lets a module register a type the root cannot name,
  and refuses a spelling already taken so it cannot redefine a root fact;
  `EnvelopeType` exposes the resolver for the same reason. The registration map is
  now read under the mutex that already guarded the envelope cache, and the
  root module's 285 recorded divergences and its race-clean run are unchanged.
- **The broker baseline carries a reason and a pointer per entry**, and the loader
  refuses a file whose entries lack either, rejects a duplicate key, and rejects an
  unsorted `entries` array — the same rules the root baseline holds itself to, for
  the same reasons.

### Fixed

- **62 recorded divergences on the Broker API HK surface, 52 of them silent, over
  20 symbols.** By kind: 29 `missing-required-name`, 23 `missing-declared-name`,
  5 `top-level-shape-mismatch` and 5 `decode-failure`; the two name classes are the
  silent ones, and the 5 decode failures restate a shape row on the same symbol.
  The largest single finding is `broker.GetPositions`, missing `cost_price`,
  `last_price` and `unrealized_profit_loss`: **the identical three names
  `brokerfd.GetFDPositions` is missing**, so a caller reading a Broker HK position
  gets the same three silently zeroed values it already got from Broker FD. The
  defect is duplicated across two surfaces, which no amount of measurement of one of
  them would have revealed; the row also misses `option_strategy` and `position_id`
  of its own. `broker.UpdateVirtualAccount` is missing
  `client_request_id`, the field item 18's verb-and-body defect is about, so the two
  findings describe one problem from two directions.
- **A dead unexported method removed from `broker/client.go`.** A three-line
  "convenience wrapper" around `Client.do` for DELETE requests, with no caller
  anywhere in the module. It survived because the root `lint` job cannot reach a
  separate Go module: CI lints the root module only, so nothing had ever examined
  this file. It is the third blind spot of one shape in this release — `broker/` was
  outside the root's tests, then outside its lint, and now outside neither.

### Noted

- **The two divergence sets are reported separately and not merged.** The root
  module's 285 and the broker module's 62 are two measurements of two modules, and
  adding them would lose the fact that one comes from a place the other cannot see.
  A limit that survives: the root module still cannot name a broker type, so a
  *change* to a broker response type is visible only to the broker module's own CI
  job, not to the root gate.
- **Proven to bite, in three directions, on the real tree.** Removing one baseline
  entry produces a stale-entry failure; blanking one reason is refused at load with
  the entry named; and dropping one row from the symbol table produces
  "manifest row ... names broker.GetFXRate, which no table here can resolve" rather
  than a smaller coverage number reading as green. The first attempt at the third
  proof passed spuriously because the removal string did not match after gofmt
  realigned the map, so the proof was redone against the aligned text.

## [2.1.25] - 2026-09-28

Harness release. The committed wire-conformance fixtures gained a per-fixture
checksum, closing the last hole in the integrity of the evidence the 285 recorded
divergences rest on. No SDK `.go` file, no fixture **content**, and no baseline
entry changed: the 193 fixture files are byte-identical and only the manifest
gained a field. Nothing is live-verified.

### Fixed

- **A length-preserving edit to a committed fixture passed every check in CI.**
  The manifest recorded each fixture's `bytes` and nothing else, so the assertion
  in `TestFixtureMatchesItsRecord` could see a reshaped or truncated artifact but
  not one whose length was preserved. A hand-edited token, a single substituted
  hex digit, or any same-length reformat would have left the gate green while the
  evidence base changed underneath the baseline that claims to describe it.
- **Each fixture now records a sha256 over its CRLF-collapsed bytes.** The
  normalisation is the same one the length tripwire already applied, so the digest
  is the same value on a Windows checkout as on a Linux one, and the manifest
  carries one number for both. The two together pin the bytes exactly: the digest
  fixes the content, the length reports the mismatch, and
  `checkNoStrayCarriageReturn` holds every CR to be half of a CRLF so the
  tolerated set cannot widen.
- **Why it was worth adding, since the generator already checks the whole tree.**
  `gen_fixtures.py --check` does compare the committed tree against the
  documentation, so a fixture edit is caught by anyone who runs
  `make conformance-fixtures`. But `--check` needs the docgen cache, and
  `AGENTS.md` records that both fixture targets **skip with exit 0 when the cache
  is absent** — which is the case in CI. Nothing in CI performed that comparison.
  The length tripwire was therefore the only fixture integrity check CI ran, and
  it could not see a length-preserving edit. The digest needs no cache, no
  network and no toolchain, so it runs everywhere the tests do.

### Noted

- **The gap was closed deliberately, and it was a deliberate absence rather than
  an oversight.** The comment above the length tripwire said a byte count "is not
  a checksum, and calling it one would be the overclaim." That was right about the
  byte count, and wrong as a reason to have no checksum. Both the claim and the
  reason are now corrected in place, and the reason it stayed absent for this long
  is recorded: not an oversight but an unexamined gap between a check that runs
  locally and a check that runs in CI.
- **Proven to bite, on the real tree, not by inspection.** One fixture's token hex
  was changed from `ccb071...` to `ccb072...`: 91 bytes before, 91 bytes after.
  The length assertion passed and the digest failed, naming both values, and only
  that row failed out of the authentication group. `TestCommittedDigestProperties`
  additionally pins the two properties the design rests on, so neither can be
  lost quietly: a CRLF input and an LF input must digest identically, and a
  length-preserving content change must not.

## [2.1.24] - 2026-09-28

Corrective release, lint only. `[2.1.23]` shipped with 26 of 27 jobs green and the
`lint` job red; this is the release that clears it. The `v2.1.23` tag stays as the
record. Both defects are in the test file added by `[2.1.22]`, so the three
consecutive corrective releases share one origin: a check written without running
the linter against it. No SDK behaviour, no fixture and no baseline entry changed,
and nothing is live-verified.

### Fixed

- **`G304: Potential file inclusion via variable` in `readRepoFile`.** The helper
  reads a repository document named by a variable, which is what gosec G304 guards
  against. The paths are constants at every call site and none is
  request-derived, which is the condition the rule exists for, so the read carries
  `//nolint:gosec` with that reason rather than a blanket suppression.
- **`S1039: unnecessary use of fmt.Sprintf`.** A single-argument `Sprintf` with no
  format verb, now a plain string.

Both were confirmed locally against `golangci-lint` v2.9.0, the version CI pins,
which reports 0 issues across the module. The first release in this sequence to have
its lint verified locally before tagging rather than by a red CI run afterwards.

## [2.1.23] - 2026-09-28

**The `lint` job is red; the other 26 are green**, including the `release records`
job added below, which ran the gate for real over 13 tags. Cause and correction are
in `[2.1.24]`. The tag stays as the record.

Corrective release. `[2.1.22]` shipped red — 18 of 26 jobs passed — and this
is the release that fixes it. The `v2.1.22` tag stays as the record of the red one,
as `v2.1.15` does. The cause was in the release-record gate added by `[2.1.22]`
itself, so the defect this repository has spent several releases chasing — a check
added without being exercised in the environment it runs in — repeated one release
later, and the tag trigger added in the same change is what caught it. Nothing here
changes SDK behaviour, no fixture and no baseline entry was touched, and nothing is
live-verified.

### Fixed

- **The release-record gate failed on a shallow clone, turning 8 CI jobs red.**
  `.github/workflows/ci.yml` checks out with the default `fetch-depth: 1`, which
  fetches the ref being built and no history. On the first tag-triggered run the
  clone held exactly one tag: the one being built. The gate listed annotated tags,
  found a non-empty list, and so its "no tags, nothing to check" skip did not fire;
  `TestReleaseRecordCutoffNamesATag` then concluded `v2.1.10` was not an annotated
  tag and failed, taking the six build/vet/test matrix jobs, the coverage job and
  the lint job with it. The failure was honest — a partial tag set genuinely cannot
  answer the question — but it was reported as a misconfigured constant rather than
  as an environment that cannot answer, which is the distinction that matters.
- **The gate now distinguishes "cannot answer" from "answered no."** It asks git
  whether the repository is a shallow clone and, if so, skips with that reason
  rather than treating a partial tag set as a complete one. `git rev-parse
  --is-shallow-repository` is used rather than looking for `.git/shallow`, because
  it is authoritative and works from a worktree. Asking and being unable to answer
  is treated as *not* shallow, so an unanswerable question never becomes a silent
  skip.
- **Verified in the exact failing condition, not by inspection.** A clone of
  `--depth 1 --branch v2.1.22` reproduces the CI condition precisely — shallow, one
  tag, non-empty — and on the old code it fails with the same message CI produced;
  on the fixed code both tag-dependent tests skip with the reason and the job goes
  green. A full clone of the same tree, which is what the new job checks out, has 54
  tags and runs the gate for real over 13 of them.

### Added

- **A dedicated `release-records` CI job, so the gate is enforced rather than
  skipped everywhere.** The other eight jobs check out shallow and will therefore
  skip the gate by design; a gate that only ever skips in CI is the failure mode the
  file was written against. This one job checks out with `fetch-depth: 0` and runs
  `go test . -run TestReleaseRecord`, so the check that every tagged release carries
  a `CHANGELOG.md` entry, a `docs/runs/index.md` row and a current status-document
  header actually executes on every push and every tag.

## [2.1.22] - 2026-09-28

**This release is red.** 18 of 26 CI jobs passed; the six build/vet/test matrix jobs,
the coverage job and the lint job failed. Cause and correction are in `[2.1.23]`: the
release-record gate added below could not answer its question in a shallow clone and
reported the environment as a misconfigured constant. The tag stays as the record,
as `v2.1.15` does.

Corrective release, and the first one whose purpose is to stop the same class of
defect recurring. It closes the loop opened by `[2.1.21]`: the release record now
has a gate, a tag push is verified, and a public API field that could never be
populated is gone. The one source-level break is named below. No fixture, no
baseline entry and no SDK behaviour outside the removed field has changed, and
nothing is live-verified.

### Added

- **A release-record gate, so a tag can no longer ship with no record.**
  `release_records_test.go` asserts that every annotated tag from `v2.1.10` onward
  carries a `CHANGELOG.md` entry, a `docs/runs/index.md` row, and is named by both
  status-document headers. This exists because the record fell behind three
  releases running and nothing in the build could see it: a path comparison cannot
  observe how a number in a document was arrived at, and a fixture byte check
  cannot observe that a release has no record. Each gap was found by a person
  checking after the fact, which is also how the wrong `trade` denominator in
  `[2.1.18]` was found.
- **The gate is bounded, and the bound is stated.** It applies from `v2.1.10`, the
  earliest tag from which every release carries both records, because the
  alternative is a gate that cannot pass: 24 earlier tags have no runs-index row
  and `v0.8.0` has no changelog entry, and writing those now would be inventing
  records for releases whose run artifacts were never kept. A gate that cannot
  pass is worth less than a narrower one that does.
- **It skips visibly rather than passing silently.** CI checks out with
  `fetch-depth: 1` and fetches no tags, so the gate reports a skip with the reason
  in the test log. A gate that treated "no tags" as "all recorded" would be green
  in CI while never having looked at anything, which is the one failure mode the
  file exists to remove. `TestReleaseRecordCutoffNamesATag` separately stops the
  cutoff being raised past the newest tag to make the gate pass vacuously.
- **It is known to bite.** `TestReleaseRecordGapsAreDetected` feeds the checker a
  repository missing exactly one record in each of five ways and requires it to
  notice, and each of the three repository-level arms — a missing index row, a
  missing changelog entry, a stale header — was additionally confirmed to fail
  end to end against the real tree before being restored. "The gate passes" and
  "the gate works" are otherwise the same claim.

### Fixed

- **A tag push now runs CI.** `.github/workflows/ci.yml` triggers on
  `push: branches: [main]` only, so no tag had ever been built, and `v2.1.20` was
  tagged on a commit nothing had run: its branch push was rejected while its tag
  push succeeded. The trigger now also carries `tags: ['v*']`. The `v` prefix keeps
  a stray ref from triggering the 26-job matrix.
- **Five public fields that could never hold a value are removed.** `data.FundNav`,
  `data.FundInfo`, `data.FundDividend`, `data.FundListItem` and `data.EventSnapshot`
  each carried an `Extra map[string]string` tagged `json:"-"`, which
  `encoding/json` never populates. Four were silent; `data.EventSnapshot`'s was
  documented with the comment "Extra holds additional fields not mapped to the
  struct", which is false, so a caller reading it expected a populated map and
  received nil. All five were unreferenced anywhere in the repository.

  **This is a source-level break** for any caller naming the field, which is
  deliberate: a field that is provably always nil is a silent runtime trap, and
  removing it turns that into a compile error a caller notices immediately. The
  wire contract is untouched — a tag-by-tag comparison of every `json` value in
  both files shows the only difference is `-` disappearing 4 times in
  `data/fund_data.go` and once in `data/eventcontracts_market.go`, with no other tag
  added, removed or changed.

### Noted

- **The `Current hardening` header field is clarified, not renumbered.** It read
  "tagged in repository `v2.1.4`; introduced in `v2.1.1`", which is easy to
  misread as the newest tag carrying SDK code. It is not: git history shows it
  tracked the release at which the hardening was last *re-declared*, moving from
  `v2.1.1` to `v2.1.4`, and `v2.1.4` changed no `.go` file at all. No release has
  re-declared the hardening since, so the value is current rather than stale and is
  left alone. The wording now says so, and records the underlying facts: `v2.1.1`
  and `v2.1.3` carry the hardening code, `v2.1.5` is the last tag to change any SDK
  `.go` file, and nothing from `v2.1.6` onward has changed SDK code.

## [2.1.21] - 2026-09-28

Release-records release. It exists to make this repository's own history legible,
because the last two tags shipped with no changelog entry, the runs index had
skipped four tags rather than two, and both status documents still named a
three-releases-old tag as current. It also records a release-process defect that
cost a release its verification. No `.go` file, no fixture, no baseline entry and
no workflow was changed; no SDK behaviour differs from `[2.1.20]`; nothing is
live-verified.

### Fixed

- **`[2.1.19]` and `[2.1.20]` were tagged with no changelog entry at all.** Both
  were written here from the tag and commit records rather than recalled, and both
  carry the negatives rather than only the results. `[2.1.19]` records that the
  divergence set was unchanged — 180 observed against 180 recorded,
  `conformance/known-divergences.json` byte-identical at sha256 `dec687e6` — and
  that this invariance is precisely why the check in the following release could be
  generated against fixed code. `[2.1.20]` records that 68 of its 104 new rows
  restate a container-kind row and were deliberately kept rather than suppressed,
  that the free-form-map guard covers 105 declared names and without it 68% of the
  finding set would have been false positives, and that its 104 rows need the same
  live evidence as the 122 before any retag.
- **`docs/runs/index.md` gained four rows, not two.** `v2.1.16` had never been
  recorded either, and it is the release that repaired the red `v2.1.15`, so its row
  is what makes the red/green pair readable from the index alone.
- **Both status documents moved to naming the current tag.** Each read
  `Latest repository tag: v2.1.17` while `v2.1.20` existed.

### Noted

- **A tag push produces no CI run, and that is why `v2.1.20` was never built.**
  `.github/workflows/ci.yml` triggers on `push: branches: [main]` only; `tags` is
  absent from its `on:` block. Every "tagged CI green" in this release history has
  in fact been the *branch* push at the tagged commit — content-equivalent, but
  never literally verified at the tag. It cost a real release: the branch push for
  `3c32230` was rejected because both remotes already held `510deae`, while the tag
  push succeeded, so `v2.1.20` landed on a commit that was on no branch and a merge
  was required to keep the published tag valid. The green 26-job run is that merge,
  `49e3eca`. Adding `tags: [v*]` to the `on:` block would close this; **it is not
  done here, because the workflow is not to be edited without an explicit request.**
- **The process lesson, recorded so it is not repeated.** A rejected branch push
  alongside a successful tag push is a state to detect *before* tagging, not after.
  This release pushes `main`, confirms both remotes hold the commit, and only then
  creates the tag.
- **A third status-document header field is stale in a way the repository cannot
  settle.** Both documents carry
  `Current hardening: tagged in repository v2.1.4; introduced in v2.1.1`, left
  unchanged by this release. Diff evidence: `v2.1.1` does introduce the hardening,
  but `v2.1.4` contains **no** `.go` changes at all — it is a docs-and-records
  release — while `v2.1.5` does change `brokerfd/events/events.go` and `stream/`,
  and no tag from `v2.1.6` onward changes SDK code at all. So if the field means
  the newest tag carrying hardening code it should read `v2.1.5`; if it means the
  tag at which hardening was declared closed, `v2.1.4` may be correct. The
  repository cannot distinguish the two readings, so no value was guessed into a
  status document. **This needs a maintainer decision.**
- **The repetition itself is the reason for this release existing.** Three
  consecutive releases shipped with records owed. The structural fix is a gate —
  asserting that every annotated tag has a changelog entry, a runs-index row, and a
  matching status-document header — which is about thirty lines in `conformance/`
  or a new test and needs no workflow change. It is not added here, having not been
  asked for.

## [2.1.20] - 2026-09-28

Harness release. The wire-conformance baseline moves from 180 recorded divergences
to 285, because a second name check now runs over the property names a page
declares without marking `required` — the route item 22 named as its third
direction and recorded as unapplied. No SDK `.go` file under `data/`, `trade/`,
`brokerfd/`, `broker/`, `client/`, `stream/`, `events/` or `display/` was changed,
no fixture was regenerated, and nothing is live-verified.

### Added

- **`missing-declared-name`, a second name check over 84 of the 90 rows the
  required-name check cannot reach.** 90 of the 154 compared rows come from pages
  publishing no `required` list, so the stronger check had nothing to look at and a
  clean row there was closer to unexamined than to correct. Those pages do declare
  property names, and the manifest already carried them as
  `checks.declaredTopLevelNames` and `checks.declaredElementNames`. The new check
  runs only where the required list is absent and records
  `missing-declared-name` (104 rows over 30 symbols, `brokerfd` 31 and `data` 73)
  plus `declared-inventory-empty` (1 row, `brokerfd.ListAccountForms`). It is a
  separate kind and a separate check rather than an extension of
  `missing-required-name` because the evidence differs and a baseline entry has to
  keep that readable: a required name is a promise the page makes, a declared name
  is a description of one response, and a name absent from a page that does not
  require it may be optional or conditionally sent. Each of the 104 reasons states
  that per entry. The 122 required-backed rows are untouched and the two figures
  are never summed.
- **It was not a formality.** 9 of the 30 rows that gained a finding recorded no
  divergence of any kind before it ran, so the 67 "fully clean" rows were not
  clean. The worst rows are `brokerfd.GetFDCorporateActions`, missing 12 of the 14
  names its page declares, and `data.GetFinancialAlert`, missing 9 of 14 and so
  carrying `eps_est`, `rev_est` and `fiscal_year` among them. The structurally
  interesting pair is `data.GetHighDividendRank` and `data.GetWeek52HighLow`:
  `data.ScreenerStock` is shared with `data.GetMostActive` and
  `data.GetTopGainersLosers`, which report no divergence at all, so those pages
  declare names a shared DTO cannot carry for any of its callers. The 90 now split
  32/58 rather than 23/67, and both status documents record the new totals.
- **`declared-inventory-empty` names a hole instead of filing it as a pass.**
  `brokerfd.ListAccountForms` on `broker-fd-us/GET-broker-forms-list` declares no
  property name at all, so neither name check had anything to look at. It is a gap
  in Webull's published documentation, not in the SDK, and it is its own kind
  precisely so the row does not read as one that examined clean when it was not
  examined. The row already carried a container-kind and a `decode-failure` row, so
  this adds no new divergent row.

### Noted

- **68 of the 104 are recorded rather than suppressed, and the dead end is worth
  keeping.** They sit on rows that also record a container-kind mismatch, 26 of
  them the `data`/`pagination_key` pair a bare-array SDK cannot carry. The obvious
  suppression is wrong, because that row set is not homogeneous:
  `brokerfd.ListFDAccounts` declares 2 names against a page of 6, so its object is
  an envelope and its names are the wrapper, while `data.GetDSLatestNews` declares
  6 against a page of 6, so its object is the payload and its 5 missing names are a
  second, independent defect; `data.GetMarketSectorDetail` is both at once,
  carrying 4 content names beside 2 wrapper keys. No per-name test exists — the
  manifest records no type per declared name, and `propertyNameCountInFixture` is 0
  on these rows because the committed minimal instances are empty — so suppressing
  per row would have discarded 20 real findings. The overlap is documented instead,
  which is the treatment the 29 `decode-failure` rows already get.
- **A guard that mattered more than the check.** `data.GetBalanceSheet`,
  `data.GetCashFlow` and `data.GetIncomeStatement` decode into a free-form
  `map[string]any`, which carries every documented name and declares none of them.
  Their pages declare 105 names between them, so without the `Carrier == nil`
  guard 68% of the whole finding set would have been false positives. The check is
  also tag-only: the generator builds each minimal instance from `required` names
  alone, so every declared name is absent from it and the required-check's
  "absent from the instance" half would have reported 100% false positives. A biting
  subtest asserts the three map rows report the check as skipped.
- **`3c32230`, the commit this tag names, was never built.** The workflow triggers
  on `push: branches: [main]` only — `tags` is absent from its `on:` block — so a
  tag push produces no run, and the branch push for `3c32230` was **rejected**
  because both remotes already held `510deae` while the tag push succeeded. The tag
  therefore landed on a commit that was on no branch, which forced a merge rather
  than a rebase to keep the published tag valid. The green 26-job run is the merge
  `49e3eca`, which contains all nine changed files, so the code is verified; the
  tagged commit itself is not. Every "tagged CI green" in this release history has
  in fact been the branch push at the tagged commit.

## [2.1.19] - 2026-09-28

Harness release, deliberately separate from the check that follows it. One latent
bug in `conformance/envelopes.go`, and the divergence set is **unchanged**: 180
observed against 180 recorded, with `conformance/known-divergences.json` byte-identical
at sha256 `dec687e6`. No SDK `.go` file was changed, no fixture was regenerated,
and nothing is live-verified.

### Fixed

- **A reconstructed envelope's struct tag could not be read at all.**
  `parseEnvelopeFields` stored `f.Tag.Value`, which `go/ast` reports as the source
  literal, so a raw string arrived wrapped in backticks.
  `reflect.StructTag.Lookup("json")` then failed on every tag and `wireNameOf`
  fell back to the Go field name. That fallback agrees for camelCase and silently
  loses a snake_case name, so a documented `pagination_key` compared as absent from
  an SDK type that carries it. The tag is now unquoted at capture, and a literal
  that will not unquote is refused rather than defaulted to the empty tag, because
  defaulting is the exact quiet failure the line exists to remove.
- **The fault was latent by construction, and that is the reusable lesson.**
  `assertEnvelopeMatchesSource` compares the rebuilt tag against the tag the parser
  read, so both sides carried the backticks and the round trip was faithful to its
  own wrong input. A test cannot catch a defect by agreeing with its own source of
  truth. `TestEnvelopeTagsAreUsableAsStructTags` asserts *usability* instead, and
  fails on the previous code with the predicted symptom; it also requires at least
  one tag and at least one snake_case json name, so it cannot pass by checking
  nothing.
- **Shipping it alone is what made it safe, and the invariant is checkable.** The
  only envelope row carrying a required name is `data.GetOptionBars`, whose sole
  required name `result` is exactly the camelCase case the fallback happened to
  rescue, so a correct fix *must not* move the 122. It did not. Had the new check
  landed on the unfixed code, it would have recorded 4 false positives against
  `data.GetCorporateActions`, `data.GetCorporateActionsByMarket`,
  `data.GetCryptoInstruments` and `data.GetOptionContracts`, each on
  `pagination_key`, and they would have entered the baseline as findings.

## [2.1.18] - 2026-09-27

Corrective documentation release. `2.1.17` documented a `trade` split whose
denominator and two counts derived from it were one too low, in three files. The
argument those numbers supported was correct and is unchanged; only the numbers
are. No `.go` file, no fixture, no generator and no baseline entry was touched, no
SDK behaviour changed, and nothing is live-verified.

### Fixed

- **The `trade` split's denominator was one low, and the two counts derived from
  it followed.** `IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md` and
  the `[2.1.17]` entry itself each restate the same split of the compared trading
  rows, and each carried the wrong figures in 6 places across the 3 files: 13
  trading endpoints are compared rather than 12, 8 publish no `required` list
  rather than 7, and 12 record no divergence rather than 11. The error was a row
  excluded from the denominator, not a count of divergences: the
  `trading/GET-trading-instruments-stocks-profiles-list` row is one the harness
  reports as not comparable, because `data.GetStockInstruments` sends a path other
  than the one the page documents, and a not-comparable row is still a compared
  row — it is the fifth of the five such rows item 22 already accounts for. Every
  figure is recounted from `conformance/testdata/manifest.json` joined to
  `conformance/known-divergences.json` on `entries[].fixture` → `fixtures[].id`,
  which matches all 180 entries; joining `fixtures[].fixture` instead matches none,
  because the entries carry the bare id and the manifest appends `.json`. The two
  halves that were already right are unchanged: 5 endpoints publish a `required`
  list, and 4 of those 5 are among the clean ones, all 4 `trade` divergence rows
  sitting on the fifth. The 7 attached to "of the clean, are not name-checkable"
  was a derived count of the same split and corrects to 8 with it.
- **No gate found this, and that is worth recording plainly.** A path comparison
  and a fixture byte check cannot observe how a number in a status document was
  arrived at, so the defect survived the release that shipped it. It was found by
  a post-release check of the approved content, not by the conformance gate, the
  citation checker, or CI. Nothing here is a claim that a green run would have
  caught it.

## [2.1.17] - 2026-09-27

Repository documentation release. `2.1.15` measured the wire-conformance
divergence set and deliberately left it unwritten; this release writes it up. Both
status documents now carry the recorded divergences as a class rather than
leaving a reader to infer from five live-blocked defects that the response
contracts are conformant. Documentation only: no `.go` file, no fixture, no
baseline entry and no generator was touched, no SDK behaviour changed, and nothing
is newly live-verified.

### Added

- **Live-blocked SDK defect item 21 and a preceding subsection now record the
  measured response-contract class in `IMPLEMENTATION_STATUS.md` and
  `docs/implementation-status.md`.** Item 21 states the class; the subsection
  before it carries the detail, so a reader is no longer invited to read the five
  items above it as the complete list of known divergences. Items 16 to 19 are
  request defects and item 20 is a single response DTO; the class is what item 20
  belongs to, taken across the surface.
- The class table, by check:

  | Class | Rows | Symbols | `brokerfd` | `data` | `trade` | How it fails |
  |---|---:|---:|---:|---:|---:|---|
  | `missing-required-name` | 122 | 29 | 89 | 29 | 4 | **Silently** — the documented value decodes to its zero value, no error reported |
  | Container kind (`top-level-shape-mismatch` 26, `element-type-mismatch` 1) | 27 | 27 | 12 | 15 | 0 | Loudly, at decode time |
  | `decode-failure` | 29 | 29 | 12 | 17 | 0 | Restates a container-kind or leaf-type row for the same symbol |
  | `leaf-type-mismatch` | 2 | 2 | 0 | 2 | 0 | Loudly — a JSON type error on the documented value |
  | **Total** | **180** | **54** | **113** | **63** | **4** | |

- **The silent class is the one that matters, and it is named as such.** The 122
  `missing-required-name` rows leave the documented value at its zero value with no
  error reported, so nothing in the return path distinguishes a zeroed field from a
  genuine zero. The container-kind and leaf-type classes fail loudly at decode
  time, and the two severities are not merged into one number.

### Noted

- **Two caveats travel with those counts, and they are what make them honest.**
  **122 is a floor, not a total**: the name check takes the union of json tags at
  every depth, because `encoding/json` flattens a response body into one name
  space, so 19 rows carry 41 of their 90 required names only through a nested
  object and count as covered. **19 of the 27 container-kind rows are
  indeterminate**: they sit on pages whose name ends `-list`, of which 12 document
  a `{data, pagination_key}` pagination envelope the SDK does not unwrap and 4
  document a single item where the path says list (`data.GetDSLatestNews`,
  `data.GetDSMarketNews`, `data.GetDSSymbolNews`, `data.GetLogos`) — the strongest
  documentation-error candidates in the set. The harness has no basis to prefer
  the page over the type, so it records all 27 as open; **do not read all 27 as
  SDK defects**, and one live body per pattern would settle them.
- **"Green is not correct" is given its denominator.** Of 154 compared rows, **90
  come from pages publishing no `required` list at all** — 101 of the 193 fixtures
  are in that state — so a clean row there is closer to unexamined than to
  correct. The name check could bite on only the other 64 rows, and it fired on 29
  of them.
- **180 is not 180 independent problems.** The 29 `decode-failure` rows are the
  weakest of the four checks and every one restates a row for the same symbol —
  27 duplicate a container-kind row and 2 a leaf-type row — so they are recorded
  so the count stays stable and no separate work is attached to them. The total
  therefore overstates the work, and the 151 rows that remain are not a defect
  count either, because 19 of them are the indeterminate container-kind rows.
- **`trade` carries 4 rows, all on `trade.BatchPlaceOrder`**, and 12 of the 13
  compared trading endpoints record none. That is a statement about what the
  instrument did not find, on the 5 trading pages that publish a `required` list
  and give it something to check — not a correctness verdict.
- **No fix is applied and nothing is live-verified.** The five live-blocked SDK
  defects are unchanged and still credential-gated; Broker FD, Display Solution,
  the US-only surfaces and footprint remain exactly as blocked. The 29 `broker/`
  endpoints are still out of scope, being a separate Go module.

### Changed

- `conformance/known-divergences.json`: the `notInStatus` note no longer claims
  the status documents under-report the divergence set by 175 rows, which stopped
  being true when item 21 and its subsection were written. It now records that the
  write-up is class-level while the 175 rows' own `recordedIn` strings are
  unchanged, and that a row in a documented class is not a row that has been
  individually triaged. No `entries` element was altered: the `recordedIn` values
  are this file's own data, each verified non-empty and meaningful, and rewriting
  175 of them to claim a row-level record that does not exist would be the less
  accurate option. The `entries` array hashes identically before and after.
- `IMPLEMENTATION_STATUS.md` and `docs/implementation-status.md`: the
  `Latest repository tag` header now reads `v2.1.17` (2026-09-27) instead of the
  stale `v2.1.8`, the `Last updated` date is 2026-09-27, and the version-history
  preamble no longer names `v2.1.6` as the current authorized repository tag. The
  historical rows of both version tables are untouched.
- **The `[2.1.16]` gap is filled, from the tag and the commit record.** That
  release's instruction listed the six files to commit and omitted this one, so
  the omission was a defect in the instruction rather than an unrecorded
  release, and the section below is written from the annotated tag and
  `cd9dd4e` — its file set, its diff and its own commit message — rather than
  invented; `[2.1.15]` and above stay byte-unchanged.

The files edited are `CHANGELOG.md`, `IMPLEMENTATION_STATUS.md`,
`docs/implementation-status.md`, `conformance/known-divergences.json` (the
`notInStatus` note only) and `docs/runs/index.md`. No SDK endpoint was
live-verified by this release, no Webull host was called, and none is claimed to
be.

## [2.1.16] - 2026-09-27

Corrective tooling release. `2.1.15` shipped red — 4 of its 26 CI jobs failed,
all four on `TestKnownDivergenceBaseline` and `TestFixtureMatchesItsRecord` in
the conformance package — and this release fixes the two portability defects
those failures named while changing no finding. Six files, every one of them
inside the conformance instrument or its generator: no production `.go` file was
edited, no SDK behaviour changed, and nothing is live-verified. The `v2.1.15`
tag is not amended and was never intended to be: it stays the record of the red
run, and rewriting pushed history to erase that record was not an acceptable
alternative.

### Fixed

- **The recorded baseline was keyed on `encoding/json`'s error text, which moves
  with the toolchain.** All 29 `decode-failure` details carried the standard
  library's message verbatim — `json: cannot unmarshal string into Go value of
  type brokerfd.AccountForm` — and a Go release that rephrases the same failure
  as `... into .0 of type brokerfd.AccountForm` makes every one of those entries
  read twice over: as a new divergence, and as a recorded entry that has stopped
  reproducing, which is the one thing the gate cannot tell apart from a fix.
  `Detail` is part of a baseline entry's identity, so a string the standard
  library owns cannot be part of it. Each detail is now composed from the
  check's own inputs — the recorded documented top-level kind, the compared
  type's shape by reflection, and the kind and name of each sharper finding — so
  it states what disagrees instead of how the decoder phrased it, and a reader
  meets one account of a defect in two entries rather than two accounts of it.
  `TestDecodeFailureDetailsSurviveTheDecoder` proves it in both directions: no
  recorded detail may contain a fragment of a decoder message, and each row is
  re-derived against a body that fails for a demonstrably different reason — a
  top-level JSON number, which the diagnosis never reads — and must come out
  byte-identical. 27 rows are re-proved that way and 2 leaf-only rows by
  inspection, because moving those would mean inventing a second real
  disagreement.
- **The fixture byte comparison was exact while the checkout it ran against was
  not.** The manifest records the length the generator wrote, and the generator
  opens every file with `newline="\n"`, so it records an LF count; but the tree
  is embedded and embedding reads the working tree, so a checkout applying a
  line-ending policy embeds a longer file than the manifest describes. This
  repository's `.gitattributes` is `* text=auto` with no `eol` rule for
  `conformance/testdata`, so nothing pins the fixtures to LF, and
  `windows-latest` checks out with `core.autocrlf=true` — CRLF is that runner's
  default, not a choice. All 193 fixtures hold at least one LF, so all 193 were
  affected, for a difference carrying no JSON meaning: CRLF is legal
  inter-token whitespace to `encoding/json`, which is why only the byte counts
  disagreed. The comparison is now exact after CRLF is collapsed to LF, and a
  second check refuses any carriage return that is not half of a CRLF, so the
  embedded bytes must equal the committed bytes plus one CR before each LF and
  line-ending policy is the only free parameter. The mismatch reproduces from a
  checkout: `authentication/POST-auth-tokens-create` is 91 bytes as committed
  and 92 as CRLF, which is the `fixture is 92 bytes, manifest says 91` the run
  named. What the tolerance gives up is exactly one thing — a fixture differing
  from the manifest only in line terminators.
- Two smaller defects in the same files, both the same shape of mistake: an
  instrument asserting less than it appeared to.
  `TestManifestTotalsAgreeWithRecords` summed the manifest's *recorded* byte
  counts rather than the tree's, so it compared the manifest against itself and
  could not fail on a real disagreement; it now reads every fixture and sums
  that. And `gen_fixtures.py` printed a hardcoded `12` while `self_test()` held
  13 checks, so the one line a reader looks at under-reported the tool's own
  coverage; it now counts what ran, and its failure line reports how many of how
  many.

### Noted

- **The same keying was in the two `unreadable` paths, and there it was an OS
  message rather than a Go one.** `conformance/shapes.go` built a
  `Divergence.Detail` by concatenating `err.Error()` — for an embed miss and for
  an unparseable instance — and an embed miss reads `no such file or directory`
  on one operating system and `The system cannot find the file specified` on
  another, so the identity carried the host as well as the toolchain.
  `unreadable` now takes the reason and the underlying cause apart: the detail
  is written from the fixture's own identity, and the OS message reaches
  `SkippedCheck.Reason`, which nothing compares. A sweep of all 180 recorded
  `detail` strings now finds no standard-library fragment and no OS message,
  where `v2.1.15` had 29 of the former. Its `reason` strings find none of
  either — but 3 of them do cite an SDK source file, on purpose, as the pointer
  to where the defect is written down, and those citations are byte-identical
  to `v2.1.15`. "No path fragment" is true of a host path and false of a source
  citation, so it is worth not resting on the phrase.
- **No finding changed.** Still 180 divergences across 54 symbols; every
  `reason` and every `recordedIn` is byte-identical to `v2.1.15`, and so are
  `coverage`, `policy` and `baselineCommit`. The per-kind distribution is the
  one `[2.1.17]` tabulates and is not restated here, so the two sections cannot
  drift apart. Only the 29 `decode-failure` detail strings were reworded. The
  class that most needs writing up — the rows that leave a documented value
  silently zeroed — is untouched, so this release makes the instrument reliable
  and fixes nothing it measures.
- **The job-level attribution is the CI run's, not the repository's.** Which 4
  of the 26 jobs failed, the Docs and citations jobs passing, and the pattern
  that localised the two defects — Go-stable matrix jobs failing where the
  pinned-toolchain macOS and Ubuntu jobs passed, both Windows jobs failing
  either way — are read from a run log that is not in this repository and cannot
  be re-derived from a checkout. The two mechanisms are each reproducible on
  their own: the toolchain wording against the 29 recorded details, and the CRLF
  mismatch against the manifest's byte counts.
- **A length is not a digest, and this release does not make it one.** The byte
  count is a tripwire against a reshaped or truncated artifact, and a content
  edit that preserves the length stays invisible to it, before and after. What
  catches a real change is the structural set around it: valid JSON, the
  documented top-level kind, every documented required name present at the level
  the page declares it, and the recorded distinct-property-name count.
- Nothing here is live-verified. No Webull host was called and no credential was
  used. The five live-blocked SDK defects are unchanged and still
  credential-gated, and the Broker FD, Display Solution, US-only and footprint
  surfaces remain exactly as blocked.

The files edited are `conformance/divergence_test.go`, `conformance/doc.go`,
`conformance/fixtures_test.go`, `conformance/shapes.go` and
`tools/conformance/gen_fixtures.py`, plus `conformance/known-divergences.json`
for the 29 reworded details and no other field.

## [2.1.15] - 2026-09-27

Tooling release. A wire-conformance harness is added: it compares every SDK
response type against the shape Webull documents for the call that decodes it,
and records 180 divergences across 54 symbols. It measures; it repairs nothing.
No production `.go` file was edited and no SDK behaviour changed.

### Added

- `tools/conformance/gen_fixtures.py` emits a committed fixture per documented
  endpoint — 193 fixtures totalling 19,192 bytes (mean 99 B, median 5 B, max
  1,293 B) plus a 285,507-byte `manifest.json` carrying per-fixture provenance
  — and passes 12 emitter self-checks before it writes anything. It is Python 3
  standard library only and imports `tools/webull-docgen/_common.py` for the
  cache format without modifying it. No fixture byte comes from Go: the SDK
  contributes path constants only, so a fixture cannot drift toward the code it
  exists to contradict.
- `conformance/` holds the comparison and its gate. Four checks run in priority
  order: required-name coverage (122 rows), top-level shape (26 plus one
  element-type row), leaf type (2), and a decode-without-error check (29) that
  is reported last and is called weak, because a decode which succeeds is
  consistent with a type that ignores every documented name.
  `TestObservedDivergenceReport` prints the whole observed set grouped by check,
  which is what a baseline edit has to be made against.
- **180 divergences are recorded across 54 symbols** in the embedded
  `conformance/known-divergences.json`: `brokerfd` 113, `data` 63, `trade` 4.
  The gate is green because the observed set is exactly the recorded set, and it
  fails in both directions because both need a person: a divergence absent from
  the file is a new finding, and a recorded entry that stops reproducing is
  either a fix or a check that changed meaning.
- Four Makefile targets: `conformance-fixtures` (`--check --self-test`),
  `conformance-fixtures-update` (`--write --self-test`), `conformance-gate` and
  `conformance-report`. `--check` regenerates in memory and fails on any
  difference, so drift from the documentation is a signal and a stale fixture
  cannot quietly become a fresh one; adopting new output is deliberately a
  separate verb. Both skip with exit 0 when the gitignored docgen cache is
  absent, since a fresh checkout cannot judge drift either way.
- **The motivating weakness is in `brokerfd`'s own tests.**
  `TestGetFDPositions` serves `json.NewEncoder(w).Encode([]FDPosition{...})` —
  the SDK's own struct — and then asserts the SDK decodes it, which proves the
  type is self-consistent and not that it matches the endpoint. The root
  module's tests serve 85 response bodies from a local server this way, 50 of
  them in `brokerfd`, and 49 of those 50 encode an SDK struct type.

### Fixed

- These are corrections to the instrument, and an independent review ran before
  any finding from it reached a status document. It found the harness
  misreporting, which is the reason to trust the numbers above rather than a
  reason to discount them.
- **Six symbol-table entries named the public projection a method returns rather
  than the envelope it actually decodes.** Those six envelopes are unexported
  and cannot be named in a `reflect.TypeOf` call at all. Comparing a tagless
  struct against a page reports a missing name for every documented property,
  which produced 12 rows the SDK does not have. The types are now read from the
  SDK sources with `go/parser` and rebuilt with `reflect.StructOf`, failing
  closed on any field it cannot represent, and each rebuilt type is compared
  field by field against the declaration on every run.
- **A reasoned waiver was suppressing a live defect, and was removed.**
  `data.Quote.QuoteTime` is `int64` where the page the method cites requires a
  `string`. The waiver argued that the sandbox sends a number, so retyping the
  field would be a regression; that reasoning fails twice over. The method
  cites the very page the harness read, so the SDK contradicts a contract it
  names, and one environment's observation does not stop a published type
  applying. On the Display Depth row there was no evidence at all, that host
  returning 403 here. Both rows are now ordinary open defects, and
  `pkg/domain/money` already accepts a JSON string or a number, so a
  two-way-tolerant type satisfies both sources.
- **`tagSite.Depth` was written and never read** while the comment claimed depth
  scoping, so the required-name count was a floor presented as a count. Depth is
  now read, and the report names what it costs: 19 rows on which 41 of 90
  required names reach a tag only from below the level the page declares them
  at. They stay uncounted, because `encoding/json` flattens a body into one
  name space and a top-level name reached one level down still decodes.
- **One false reason had been copied 88 times** — that a DTO is named after the
  undocumented `/broker-fd/*` payloads — and it was false for 86 of them,
  because those methods send the path the page documents. Every entry now
  carries a reason true of that entry. A required non-empty reason is not a
  required true one, and these two rows are the case that proved it: the guard
  was satisfied perfectly by a false sentence, identically on both entries.
- Also corrected: two latent `encoding/json` embedding bugs, now covered by a
  regression test; an over-permissive `money.Money` leaf exemption; a stray `%s`
  in two provenance strings; and a note claiming one waiver where there were two.
- The fixing pass then corrected the review itself on six points, among them
  that the review's 170/14/16 split does not decompose.

### Noted

- **The status documents materially under-report, and this release does not fix
  that.** Only 5 of the 180 recorded rows are already written up in
  `IMPLEMENTATION_STATUS.md`: three `GetFDPositions` rows and two
  `GetFDAssetsDetail` rows, all under live-blocked defect item 20. The other 175
  state in their own `recordedIn` field that no status-document entry exists for
  them. A reader who takes the five live-blocked SDK defects as the complete
  list of known problems is wrong by 175 rows, and this release is not evidence
  against that reading — it is the measurement that establishes it. Writing the
  ~170 genuine rows up is deliberate follow-up work, deferred until the
  instrument has run in CI, so that no claim is made with a tool CI has never
  executed.
- **180 is a row count, not 180 defects.** 16 of the 27 shape rows are pages
  that document a bare object on an endpoint whose page name ends in `-list`
  while the SDK decodes an array. The harness cannot tell whether the page
  omitted the array wrapper or the type is wrong, and records each as an open
  defect because it has no basis to prefer the page over the type. One live
  probe against any of them would settle it.
- **What the harness proves is narrow, and the narrowness is the point.** It
  shows that a documented response shape is or is not representable in the SDK's
  types. It does not show the SDK is correct where the harness is green, and it
  cannot: `data.GetDisplaySnapshot` reconciles as a clean path match and is
  still defective, so a green row in `docs/reconciliation.md` remains a path
  comparison rather than a correctness verdict.
- **The instrument's limits, which a green run does not lift.** A required name
  reachable only through a nested object is not reported, so 122 is a floor and
  not a total. The decode-only check is the weakest of the four and its 29 rows
  mostly duplicate a shape or name row for the same symbol. 5 compared rows are
  marked not comparable, because the SDK method sends a different path than the
  page documents and the page is therefore not that call's contract; those rows
  claim nothing either way. 29 of the 193 endpoints are out of scope entirely,
  because `broker/` is a separate Go module whose types are not importable from
  the root module without a `go.mod` change, so no `broker/` row is covered
  here. 154 of 193 endpoints are actually compared.
- **No divergence is fixed and none is claimed to be.** The 175 unwritten rows
  are open defects. The five live-blocked SDK defects recorded through `2.1.11`
  remain blocked and unchanged; this release describes them and fixes none.
- Nothing here is live-verified. No Webull host was called and no credential was
  used. The five live-blocked defects are still credential-gated, and the Broker
  FD, Display Solution, US-only and footprint surfaces remain exactly as blocked
  as `2.1.11` recorded them.

No SDK endpoint was live-verified by this release, no Webull host was called,
and none is claimed to be. No existing `.go` file was modified; no SDK code,
generator or generated documentation was changed. The files edited are
`CHANGELOG.md`, `AGENTS.md` and `docs/runs/index.md`.
`IMPLEMENTATION_STATUS.md` and `docs/implementation-status.md` are deliberately
untouched, because the under-report above is a finding to be written up, not one
to be papered over in the release that measured it.

## [2.1.14] - 2026-09-27

Repository documentation patch release. The `[2.1.13]` citation gate resolves a
fully qualified `file.go:25` directly but has to *infer* the file behind a bare
`assets.go:25`, and this release makes every citation in the three authority
documents name its own file, so those inference rules now resolve nothing. No
line number changed and no citation target moved. No SDK behaviour changed and
no `.go` file was edited.

### Changed

- 35 bare-filename citation spans in `IMPLEMENTATION_STATUS.md` and
  `docs/implementation-status.md` are now fully qualified, so all 86 full
  citation spans name their own file; `AGENTS.md` already named every file it
  cites and is unchanged. The gate's four bare-filename rules — recall an
  earlier citation of that name, a sibling of the most recent citation, a file
  at the repository root, then a guess — now match nothing. The run reports 0
  recalled, 0 sibling, 0 root and 0 guessed where the previous release shipped
  23 recalled and 12 sibling, and the gate moves from 0 errors with 23
  `ambiguous-basename` warnings to 0 errors and 0 warnings.
- **The line numbers are untouched; only the path was spelled out.** Resolving
  both versions of all three documents through the checker itself and comparing
  them row by row, all 104 resolved citations — the 86 full spans, the 8 bare
  `:NN` continuations, and the further items that comma lists and ranges expand
  to — resolve to the same file, the same line spec and the same adjacent label
  as before, and 0 moved. The line-number multiset is identical at 132
  occurrences per version with 0 disagreements. 45 of those 104 items were
  resolved by inference before and are resolved literally now; the 51 that
  already named their own file and the 8 continuations are untouched. A word
  diff over both documents confirms the stronger form: every prose token
  outside a code span is identical, 4837 and 3874 of them, so no prose word was
  removed or altered and the only changes are 35 code spans gaining a
  directory.
- **This is a spelling fix, not a change of meaning, and the file was settled
  independently of the tool.** An analysis pass established, for all 45
  inferred targets, that the file the prose intends is the file the tool chose:
  zero disagreements. The decisive evidence is a 1:1 correspondence rather than
  an argument — the 14 undocumented `/broker-fd/*` literals these citations
  point at all live in `brokerfd/`, and there are exactly 14, so the file
  follows from the enumeration the prose already gives rather than from the
  tool's guess. The only other `/broker-fd/` literals in `brokerfd/` are four in
  `brokerfd/brokerfd_test.go`, which nothing cites.
- Two working assumptions were wrong and are recorded here rather than left
  implicit. First, the two status documents do not share a wrap style: the two
  affected list items in `IMPLEMENTATION_STATUS.md` are single lines of 1375
  and 1654 characters and 83 of its 262 lines exceed 100 columns, so its
  citations were qualified in place, preserving that file's style; only
  `docs/implementation-status.md` is hand-wrapped, and its bullets have no
  single wrap width, since no width from 70 to 94 reproduces the affected
  bullet's own line breaks, so that one bullet was re-wrapped by hand to an
  80-column budget, which is its new maximum. Second, the unit was first
  miscounted as "45 citations" when 45 is the number of resolved *items*: there
  are 35 bare *spans*, and a span such as `instruments.go:27,29,31` expands to
  three items. Both numbers are correct for what they measure and this entry
  does not conflate them.

### Noted

- **The gate is not more correct now, only less dependent on a heuristic.** It
  verifies that a cited line exists, is in range, is not blank, and is not a
  pure comment. A citation pointing at a wrong but real line of the right file
  is still invisible to it, so the boundary `[2.1.13]` states is unchanged: a
  green run is evidence that the citations resolve to existing, non-blank,
  non-comment lines inside the cited file, and nothing more. This release is
  not a claim that any citation is correct, only that the file each one was read
  from is now written down instead of recalled.
- The five live-blocked SDK defects recorded through `2.1.11` remain blocked
  and unchanged. This release documents citations about them; it fixes,
  verifies and unblocks none of them, and a 0-warning run says nothing about
  whether any of them has been fixed.

No SDK endpoint was live-verified by this release, no Webull host was called,
and none is claimed to be. No SDK code, no generator and no generated
documentation was changed, and neither was the gate itself:
`tools/citations/check.py` is byte-unchanged, so the checker that produced the
measurements above is the one `[2.1.13]` shipped. The files edited are
`IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md`, `CHANGELOG.md`,
and `docs/runs/index.md`.

## [2.1.13] - 2026-09-27

Repository documentation-and-tooling patch release. A `file:line` citation
validator is added, after five consecutive releases shipped stale citations and
nothing mechanical checked any of them, and a third self-correction is recorded
from the `v2.1.12` release record. No SDK behaviour changed and no `.go` file was
edited.

### Added

- `tools/citations/check.py` validates the `file:line` citations in the status
  documents. It is Python 3 standard library only: no dependency, no
  `pip install`, and no module's `go.mod` touched. The trigger is a documented
  requirement rather than bad luck — `AGENTS.md` requires live-blocked defects to
  be recorded with `file:line` citations, and across five consecutive releases
  those citations went stale and shipped. A comment inserted above one const
  block invalidated nine references across three documents, one manifest citation
  was transposed at birth, another pointed at a blank line, and a human reviewer
  caught some of them.
- A `make citations` target runs it from the repository root, and a `citations`
  job runs it in CI. The target reports only; a non-zero exit fails the build.
  The CI job invokes the checker directly rather than through the target, which
  matches every other gate in that workflow. The edit to
  `.github/workflows/ci.yml`, a protected file, was explicitly approved on
  2026-09-27.
- Checked by default: `IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md`,
  and `AGENTS.md`. `CHANGELOG.md` is an explicit waived constant,
  `WAIVED_DOCUMENTS` in the tool, because it holds released history whose seven
  known-stale citations must stay byte-unchanged. The waiver is a named constant
  so that it is greppable rather than an omission, and the document list can be
  overridden on the command line for testing while the defaults stay those three.
- Five failure conditions: `missing-file`, `line-out-of-range`, `blank-line`,
  `comment-line`, a pure comment being `//` in Go, `#` in Python, and `<!--` in
  Markdown, and `empty-waiver`, a `# citation-waiver:` marker carrying no reason.
  Plus `unresolved-continuation`: the bare `:NN` form the status documents use
  heavily, when it has no earlier full citation to resolve against. A bare `:NN`
  is resolved against the most recent full citation in the same document rather
  than skipped, because a checker that ignored the form would check far less than
  it appears to. `missing-file` is reachable only by a path carrying its own
  directory: a mistyped bare filename names no file at all, so it is reported as
  `guessed-basename` or `unresolved-basename`, both warnings, and the gate stays
  green. A typo in a fully qualified path fails; a typo in a bare one does not.
- A bare filename is resolved against the document's own citation history, so the
  tool, not the document, decides which file it reads — and 41 Go basenames here
  name more than one file. Every such resolution is reported as
  `ambiguous-basename`, naming the file actually checked and every candidate, and
  `--verbose` prints the resolution of every citation whether it was inferred or
  named outright. The exit status stays zero: on the three default documents this
  reports 23 ambiguous resolutions over 86 full citations and still exits green,
  which is the honest state of those documents rather than a green light on their
  correctness.
- A range is held to a weaker rule on purpose, and the difference is a real
  limitation rather than a convenience: a range fails only when *no* line in it is
  code, so a range starting on a GoDoc comment, or spanning a struct whose
  interior is field comments, still passes. The count of the passing ranges that
  exercise the weakening is reported in the summary so it stays visible, and a
  range that failed is not counted among them. The weakening is not theoretical:
  `data/display_quotes.go:41-49` is a range that drifted onto a const, and this
  tool does not report it.
- A `# citation-waiver: <reason>` marker on the citing line skips that citation
  and counts it, and the marker must sit outside every code span, because one
  inside a span documents the syntax rather than waiving anything. The reason is
  enforced, not decorative: a marker with nothing after the colon is not a
  waiver, so the citation is checked and `empty-waiver` fails the gate.
  `AGENTS.md` records the escape hatch, because a maintainer whose build just
  failed on the new job needs to know it exists.

### Noted

- **The gate does not prevent citation errors.** It detects one specific
  mechanical subset: a cited file that is missing, a line that is out of range,
  blank, or a pure comment, a continuation that cannot be resolved, and a waiver
  that does not say why. It does not check that the cited line contains what the
  prose claims, so it cannot detect a label/number transposition, where the
  number is a valid line of the right file but the wrong line. Nor can it see a
  stale citation that drifted onto a line which still reads as code: of the seven
  historical misses the release record names, it reports the four now reading as
  a blank or comment line and cannot see the three now reading as code, which are
  `tools/webull-docgen/_common.py:298`, `data/display_quotes.go:52`, and the
  `data/display_quotes.go:41-49` range. Measured on 2026-09-27 against a copy of
  `CHANGELOG.md`, the four are reported as 5 `comment-line` and 2 `blank-line`
  findings over 4 distinct targets. The `duplicate-label` warning is a partial
  mitigation for the transposition class, not a fix, and fires far too rarely to
  be trusted. A green run is evidence that the citations resolve to existing,
  non-blank, non-comment lines inside the cited file, and nothing more.
  Overstating what the gate does would repeat the exact error this entry records
  below.
- The tool never rewrites, repairs, or auto-fixes a citation, and it writes to no
  file, including the documents it reads. There is no `--fix` and there will not
  be one. An auto-rewritten citation is a wrong citation that looks verified,
  which is worse than a red gate.
- Neither status document's recorded defect analysis is touched by this release.
  The five live-blocked defects recorded through `2.1.11` remain blocked and
  unchanged, and a green citation run says nothing about whether any of them has
  been fixed.

### Fixed

- Three defects in the citation checker itself, found by review of the unreleased
  tool and corrected before the tag rather than after it, because a gate that
  validates the wrong file is worse than no gate. Each shipped a claim in this
  entry that was the opposite of what the code did.
- **A bare basename could bind to the wrong package and the run stayed green.**
  The checker resolved a bare `accounts.go` plus a line number by recalling the
  most recent full citation of that name, and 41 Go basenames here name more than
  one file, so it could read `broker/accounts.go` where the prose meant
  `brokerfd/accounts.go` — two files that share a name and nothing else. On the
  three default documents 35 of 86 full citations have their file inferred rather
  than named, 23 of them against a name the repository holds more than once, and CI
  runs without `--verbose`, so nothing identified which file each one landed on.
  Every such resolution is now reported as `ambiguous-basename` with the file
  checked and every candidate named, and `--verbose` prints the resolution of
  every citation. The requirement that every citation be fully qualified was
  deliberately *not* imposed: that would mean rewriting 35 citations in the same
  change that checks them, and each rewrite is a chance to introduce a fresh
  error.
- **The "the reason is required" control did not exist.** The checker matched the
  `# citation-waiver:` marker and never read the text after it, so a waiver with
  an empty reason waived silently — while this entry, the tool's `--help`, its
  module docstring, and `AGENTS.md` all asserted that it could not. Enforced now:
  an unexplained marker is not a waiver, the citation is checked, and
  `empty-waiver` fails the gate.
- **"Catches two cleanly of seven real historical misses" was false.** Measured
  against `CHANGELOG.md`, the only document still holding the historical misses,
  the checker reports 4 distinct stale targets, not 2, as 5 `comment-line` and 2
  `blank-line` findings — 7 findings, which is where the number collision in the
  prose came from. The adjacent claim that both historical misses were
  single-line citations was false too, since one of the three it cannot see is
  the `data/display_quotes.go:41-49` range. The limits text in all three places
  now states the measurement and names what is being counted.
- Smaller corrections in the same pass: `Makefile` told the reader to run from
  the repository root so cited paths would resolve, when the root is derived from
  the tool's own location and the target works from any directory; the failure
  list read as if a mistyped bare filename failed the gate, when only a path
  carrying its own directory can; the footer's range counters were incremented
  before the range was checked, so a range that *failed* was counted under the
  heading "not failures"; and the comment claiming the tool never opens
  `CHANGELOG.md` was true only because the file is waived, not because the tool
  declines to.
- The `[2.1.12]` entry's first bullet claims that the `AGENTS.md` edit "is
  described correctly in the `[2.1.11]` body". That is false: the `[2.1.11]`
  section contains no occurrence of `AGENTS` at all, and the edit is described
  only in the `59f7859` commit message. Only the location claim was wrong. The
  rest of the clause is true and stands — the `AGENTS.md` edit is real, its
  content is accurate, and the omission that bullet reports is a file-list
  omission — so nothing else in it needs correcting. The error was caught by the
  release agent reviewing its own output after `v2.1.12` had been pushed, and was
  deliberately left uncorrected at the time because the release had already
  shipped.
- `[2.1.12]` is itself **left byte-unchanged**, on the same principle that left
  `[2.1.11]`, `[2.1.6]`, and `[2.1.8]` byte-unchanged, so this correction is
  recorded here instead. Amending a tagged entry would be the first exception in
  six corrections and would contradict the invariant this same release installs
  as a waived constant: the `[2.1.12]` entry states that rule about itself in the
  very sentence that carries the error, so an amendment would leave the entry
  citing a rule it violates. The error is also inert — it concerns where a
  description of a documentation edit lives, not SDK behaviour, an API, a defect,
  or any verification claim, so no reader acts differently because of it.
- This correction is folded into the new work rather than released on its own. A
  docs-only release per typo has negative expected value: each new entry is
  another chance to introduce a fresh error while the underlying problem goes
  unaddressed.

No SDK endpoint was live-verified by this release, no Webull host was called, and
none is claimed to be. No SDK code was changed. The files edited are
`tools/citations/check.py`, `Makefile`, `.github/workflows/ci.yml`, `AGENTS.md`,
`CHANGELOG.md`, and `docs/runs/index.md`.

## [2.1.12] - 2026-09-27

Repository documentation-only patch release recording two textual errors found in
the `v2.1.11` release record after it was pushed. Both are prose errors, not
code errors, and both are recorded here rather than silently amended, because
`59f7859` is already on two remotes and correcting either would mean rewriting
pushed history. No SDK behaviour changed and no `.go` file was edited.

### Fixed

- The `[2.1.11]` entry's closing file list omits `AGENTS.md`. That release edited
  five files — `AGENTS.md`, `CHANGELOG.md`, `IMPLEMENTATION_STATUS.md`,
  `docs/implementation-status.md`, and `docs/runs/index.md` — and `git show --stat
  59f7859` names all five. The omission is in the file list only: the `AGENTS.md`
  edit itself is real, is described correctly in the `[2.1.11]` body, and its
  content is accurate. The shipped entry is **left byte-unchanged**, on the same
  principle that kept the five stale citations in `[2.1.6]` and `[2.1.8]`
  unchanged, so this correction is recorded here instead.
- The `v2.1.11` commit message contains a wrong clause: it says `FDPosition`
  "declares tags for none of the first, second, and last of those" required
  properties. The correct three are the **first, fourth, and last**:
  `cost_price` (the struct has `average_cost`), `last_price` (the struct has
  `market_value`), and `unrealized_profit_loss` (the struct has
  `unrealized_pl`). The second, `currency`, **is** declared, at
  `brokerfd/assets.go:89`, so the clause both omits the real fourth and
  misidentifies `currency`. Only the enumeration is wrong: the immediately
  following sentence in the same commit body names the three properties
  correctly, so the substantive claim that three required properties cannot be
  received and decode as silent zeroes stands. The `[2.1.11]` entry here and the
  defect entry in both status documents also name the correct three, so the error
  is confined to the commit message. The commit is left as pushed rather than
  force-pushed.

No SDK endpoint was live-verified by this release, no Webull host was called, and
none is claimed to be. No SDK code was changed: the only files edited are
`CHANGELOG.md` and `docs/runs/index.md`. Neither correction affects the code or
the defect analysis; both are errors in the prose of a release record. The five
live-blocked defects recorded through `2.1.11` remain blocked and unchanged.

## [2.1.11] - 2026-09-27

Repository documentation-only patch release: two stale citations are corrected, a
newly-found silent data-loss defect is recorded, and a citation-sweep finds the
correction in `2.1.10` itself was incomplete. No SDK behaviour changed and no
`.go` file was edited.

### Fixed

- Two stale `file:line` citations in `IMPLEMENTATION_STATUS.md` and
  `docs/implementation-status.md` are corrected. Both pointed one line above the
  thing they name. The assets-summary manifest entry was cited as
  `tools/webull-docgen/_common.py:313`, where `:313` is the last line of a
  4-line comment inserted in `2.1.10` and the "Assets Summary" tuple is at
  `:314`; the List Accounts manifest entry was cited as
  `tools/webull-docgen/_common.py:298`, where `:298` is the `[` that opens the
  endpoint list and the "List Accounts" tuple is at `:299`. Both are now `:314`
  and `:299`.
  The `2.1.10` correction pass did not land: it moved the assets-summary
  citation from `:298` to `:313` while the tuple had already moved to `:314`, so
  it was still one line short when shipped, and it never re-examined the
  List Accounts citation. The latter is worse than a shift — at `5dcb008`, when
  it was added, `:298` was the *assets-summary* tuple, not the List Accounts one,
  which by then was at `:287`; the label and the number were transposed, and
  `2.1.10` renamed the entry to "List Accounts" without correcting the number.
  The citation-scope note in both status documents is widened accordingly: it
  now covers `tools/webull-docgen/_common.py` and
  `tools/webull-docgen/docgen.py` line numbers as well as
  `data/display_quotes.go`, states that a `file:line` is valid only for the
  commit that fixed it, and records that a second pass over one file missed
  another.
- Every `file:line` citation in the four documentation files was swept against
  the current source, not only the two named above: all 45 citations across
  `IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md`, `CHANGELOG.md`,
  and `docs/runs/index.md`, covering `brokerfd/client.go`,
  `internal/region/region.go`, `broker/client.go`, seven `brokerfd` path-literal
  sites, `broker/accounts.go`, `broker/accounts_test.go`,
  `data/display_quotes.go`, `data/snapshot.go`, and
  `tools/webull-docgen/_common.py`. The citations in both status documents all
  resolve to the declarations they name. Five further stale citations were found
  in `CHANGELOG.md` and are **left unchanged**, because they sit in the
  `[2.1.6]` and `[2.1.8]` sections, which are historical records and must stay
  byte-unchanged. The stale numbers are recorded here in words rather than in
  `file:line` form, so this entry adds no citation of its own that fails to
  resolve. In the `[2.1.6]` `data.GetDisplaySnapshot` bullet, the snapshot
  constant and the `Get` call site are cited at lines 28 and 52, now 33 and 63
  after the `2.1.10` comment insertion; in the `[2.1.8]` bullet, the query
  construction and `GetDisplaySnapshot` are cited at lines 41-49 and 40, now
  52-60 and 51; and `examples/options-multi-leg/main.go:193` points at a blank
  line, while the `client_order_id` change it describes is at line 199.

### Noted

- A fifth live-blocked SDK defect is recorded in `IMPLEMENTATION_STATUS.md` and
  `docs/implementation-status.md`: `brokerfd.GetFDPositions` silently zeroes
  three of the eight required response properties. The request is correct —
  `GetFDPositions` (`brokerfd/assets.go:92-100`) requests
  `pathFDAssetsPositions` (`brokerfd/assets.go:27`), exactly the documented
  `GET /broker/assets/positions/list` — so the reconciliation row is a clean
  `✅ match` and the report gives no signal, the same path-only limitation already
  recorded as a caveat. The defect is on the response side: the documented `200`
  is a `type: array` whose `items` require `cost_price`, `currency`,
  `instrument_type`, `last_price`, `position_id`, `quantity`, `symbol`, and
  `unrealized_profit_loss`, while `FDPosition` (`brokerfd/assets.go:79-90`)
  declares tags for `position_id`, `account_id`, `symbol`, `quantity`,
  `average_cost`, `market_value`, `unrealized_pl`, `realized_pl`,
  `instrument_type`, and `currency`. **`cost_price`, `last_price`, and
  `unrealized_profit_loss` have no matching tag.** All three corresponding fields
  are `money.Money`, so a missing key leaves them at zero and `encoding/json`
  returns no error, which makes this the most dangerous of the five: the path
  defects fail loudly with a `404`, while this one hands the caller a successful
  call and three wrong zeroes indistinguishable from a genuinely zero position.
  A bare retag is not recorded as the fix, because the two shapes diverge in both
  directions and the documentation cannot say which the server honours; the
  recommendation is conditional on a probe. Unblocked by the same US-sandbox
  credential as the two `brokerfd` path defects, plus one question added to the
  Webull enquiry that item already needs, so no new round trip. A cross-reference
  is added from the `/broker-fd/*` path defect, but the two stay separate
  entries, because a response-schema defect and a request-path defect warrant
  different urgency. The adjacent `GetFDAssetsDetail` container-type mismatch is
  recorded in the same entry and explicitly marked as failing loudly, not
  silently.
- The static tag-versus-required-name mismatch is certain: the struct tags and the
  cached reference page's `required` list were compared directly, with no
  interpretation involved. Which side the live server honours is unverified.

No SDK endpoint was live-verified by this release, no Webull host was called, and
none is claimed to be. No SDK code was changed: the only files edited are
`IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md`, `CHANGELOG.md`, and
`docs/runs/index.md`, and the retag in the new defect entry is a recommendation
recorded in the status documents, not an edit. The four defects recorded in
`2.1.7` remain blocked and unchanged.

## [2.1.10] - 2026-09-27

Repository patch release of three honesty fixes that need no credential, one
regression guard against a new failure mode, and a read-only diagnostic cache
refresh. No SDK behaviour changed: the only `.go` edit is comments, in
`data/display_quotes.go`, and no `.go` statement was altered.

### Fixed

- The unverified marker on `data.GetDisplaySnapshot` was missing. Commit
  `2b29c88` aligned four sibling Display constants from `/openapi/...` to
  `/market-data/stocks/...` in the same const block, left `pathDSSnapshot` as
  the sole `/openapi/...` holdout, and deleted the
  `TODO(ds): Confirm exact paths against US sandbox` comment that had flagged
  it, so nothing in the source signalled that one constant of five is
  unverified and it read as already confirmed. The const comment and the method
  GoDoc now state it, name the two things that make it unresolvable from the
  docs alone (the reference page contradicts itself, with `operationId`
  `snapshotUsingGET` against `method` `post`, and the Display host answers `403`
  without a paid entitlement), and warn that aligning it by analogy to the
  aligned siblings is the mistake the deleted marker existed to prevent. The
  constant's value, the method's signature, the query construction and the call
  site are byte-unchanged. The defect recorded in `2.1.7` is unchanged and still
  blocked; this entry restores the warning, it does not resolve the path.
- The docgen page cache could not be safely refreshed. `CACHE` was a hard-coded
  module constant, so a refresh could not be redirected away from the committed
  cache, and `fetch()` called `urllib.request.urlopen` unwrapped, so one timeout
  aborted a whole run partway. The cache-hit test is `getsize > 0`, so such a
  partially refreshed directory was then reused indefinitely and silently.
  `tools/webull-docgen/_common.py` now resolves the directory per call through
  `cache_dir()` with a `WEBULL_DOCGEN_CACHE` override, blank or unset still
  meaning `.cache/`; `fetch()` retries up to `FETCH_ATTEMPTS = 3` with `1.0` and
  `2.0` backoff and the existing 0.2s politeness pause now applied after every
  attempt, successful or not; and a page that fails every attempt is recorded and
  skipped rather than aborting the run. A skip is surfaced four ways: on stdout,
  in `changes.json` beside the cache directory it was rendered from, as the
  absence of a file, and as an in-page `*Unavailable*` marker. Nothing is written
  on failure, because a zero-byte file fails the cache-hit test and is re-fetched
  while a truncated one passes it, which would turn a transient outage into a
  permanent, invisible hole in the evidence base.
- The skip-don't-abort behaviour above introduced a regression, and it is
  guarded. A missing cached page could flip a row's status and still produce a
  full-length report that passed `mkdocs --strict`, and a missing `llms.txt`
  index could collapse the status partition from `summary=4 differs=1` to
  `summary=0 differs=5` with the row count, the totals and the strict build all
  unchanged, so nothing downstream would reveal the loss.
  `generate_reconciliation()` now refuses to write when a manifest-mapped page
  or an `llms.txt` index is missing, or when any page failed for a reason other
  than a permanent `404`/`410`, naming the URLs on stderr. The permanent-404
  exemption is required rather than cosmetic: 8 URLs that Webull's own `llms.txt`
  indexes advertise return a permanent `404` for their `.md` variant, so a broad
  "any skip is fatal" rule made the report unproducible. The other two
  generation targets keep marking the loss in the page itself, where a lost page
  thins documentation instead of changing a verdict.
- Nine stale `file:line` citations in `IMPLEMENTATION_STATUS.md`,
  `docs/implementation-status.md` and `AGENTS.md` were invalidated by the
  comment insertion above and are corrected. Both status documents also gained a
  citation-scope note recording that those line numbers are anchored to the
  2026-09-27 state of the files and must be re-checked after an edit.

### Changed

- The two `broker-fd-us` manifest entries that map two SDK symbols onto one
  documented page, "Assets Summary" and "Positions", carried an empty note, so a
  reader could not tell which symbol the reported path belonged to. The notes in
  `tools/webull-docgen/_common.py` now name it, because a row renders a single
  SDK path. "Positions" is the case that carried information:
  `brokerfd.GetFDPositions` sends the documented
  `GET /broker/assets/positions/list` exactly and so is a `✅ match`, and the
  uncovered gap was that the green row silently stood for a second symbol,
  `brokerfd.GetPositions`, which sends the undocumented `/broker-fd/positions`
  with no `account_id`. "Assets Summary" additionally records that the documented
  `GET /broker/assets/summaries/get` is sent by neither mapped symbol and that
  the response DTO is a flat struct rather than the documented `balance` and
  `positions` envelope. Both notes state that no fix was applied and point at
  `IMPLEMENTATION_STATUS.md`; the 14 `brokerfd` `/broker-fd/*` literals are
  unchanged.
- `docs/reconciliation.md` is re-rendered with the two notes above and a
  `2026-09-27` snapshot date. Its partition is unmoved at 184 exact OpenAPI JSON
  path matches, 4 summary-only, 1 differing, 0 unresolved SDK paths, 3 rows
  labelled `no OpenAPI schema on page` and 17 manifest entries deliberately
  mapped to no SDK symbol, summing to 209 implemented endpoints and 0
  documented-only gaps. It remains a path comparison, not a correctness verdict.

### Noted

- A diagnostic cache refresh was run read-only and deliberately not adopted. 317
  pages were fetched into a scratch directory through `WEBULL_DOCGEN_CACHE`; 306
  of the 309 pages comparable with the 2026-09-22 cache were byte-identical, and
  the 3 that differed were host-table and prose edits changing no `path`, method
  or required-body property. The four live-blocked defects recorded in `2.1.7`
  were re-checked against their full OpenAPI specs and are unchanged, so the
  refreshed cache was not adopted and the committed report remains a render of
  the 2026-09-22 evidence.

No SDK endpoint was live-verified by this release, and none is claimed to be.
The only network traffic was the read-only scratch refresh above, which observed
that 8 `llms.txt`-indexed URLs permanently `404` their `.md` variant and read
host-table text; it authenticated to no Webull host and called no SDK endpoint.
The four live-blocked defects from `2.1.7` remain blocked and unchanged, still
waiting on credentials or entitlement this environment does not have.

## [2.1.9] - 2026-09-26

### Fixed

- The `AGENTS.md` release-status paragraph hardcoded the current repository tag,
  so it went stale on every release and was four tags behind by this one. It now
  states only the durable facts and points at `CHANGELOG.md` for the current
  release, so a future release does not require editing it again, and it gains
  the reconciler's path-only caveat as a cross-reference: a `✅ match` row in
  `docs/reconciliation.md` is a path comparison, not a correctness verdict.

Nothing in this entry is live-verified; no network call was made.

## [2.1.8] - 2026-09-26

### Fixed

- Two of the four live-blocked SDK defects recorded in `2.1.7` understated the
  size of their own fixes. The `2.1.7` wording described
  `data.GetDisplaySnapshot` as an "align the path and switch to the documented
  POST body" change; it is larger, and correcting the request is a **breaking
  public API change**. The documented request is a required JSON body whose only
  required property is `category_symbols`, an array of `{category, symbols[]}`
  objects, while `data.SnapshotQuery` (`data/snapshot.go:35-47`) models a single
  category with a flat symbol list and `data/display_quotes.go:41-49` sends it as
  query parameters. Because `data.SnapshotQuery` is the public parameter type of
  both `GetDisplaySnapshot` (`data/display_quotes.go:40`) and `GetSnapshot`
  (`data/snapshot.go:142`), moving it to the documented array breaks consumers of
  either method, so the fix needs a maintainer decision on how to version a
  breaking change on a module that stays on the v1 import path. The
  `extend_hour_required` and `overnight_required` flags are already
  type-compatible. The same `2.1.7` entry implied a uniform one-line fix for the
  14 undocumented `brokerfd` `/broker-fd/*` literals; only 4 have an unambiguous
  documented counterpart, 1 is a probable duplicate, 6 are plausibly ambiguous,
  and 3 have no documented counterpart at all, which credentials cannot settle
  and which need a written answer from Webull. The unblock requirements recorded
  in `2.1.7` are unchanged.
- The reconciler's limits are now stated as a caveat in
  `IMPLEMENTATION_STATUS.md` and `docs/implementation-status.md`. It compares
  path strings only, so a `✅ match` row is not evidence that an endpoint is
  correct: it reports `broker.UpdateVirtualAccount` as a clean match although its
  request is defective, and it cannot observe HTTP verbs, request bodies, or
  transport-host routing. This is a limit of a path comparison, not a defect in
  the generator, which is doing its stated job.

Nothing in this entry is live-verified; no network call was made.

## [2.1.7] - 2026-09-26

### Fixed

- The doc generator mislabelled two categories of reconciliation row, which
  inflated the "unresolved SDK path" headline. A gRPC reference page that embeds
  no OpenAPI schema is not a REST endpoint, and a manifest entry deliberately
  mapped to `-` is not a missing symbol, but both were reported as unresolved
  SDK paths. The generator status labels were corrected in
  `tools/webull-docgen/`, not in the generated output, and four false-positive
  unresolved flags were removed as a result.
- The hand-written documentation still quoted the superseded "180 exact, 4
  summary-only, 0 differing, 25 unresolved" figures. `AGENTS.md`,
  `IMPLEMENTATION_STATUS.md`, `docs/implementation-status.md`, `docs/index.md`,
  `docs/webull-api.md`, `ARCHITECTURE.md`, and `PLAN.md` now carry the
  2026-09-26 figures, and the two status docs record the four live-blocked SDK
  defects below.

### Changed

- Regenerated the 2026-09-26 reconciliation snapshot. Its partition is 209
  implemented endpoints and 0 documented-only gaps, of which 184 exact OpenAPI
  JSON path matches, 4 summary-only matches, 1 path differing from both
  sources, 0 unresolved SDK paths, 3 rows labelled `no OpenAPI schema on page`,
  and 17 manifest entries deliberately mapped to no SDK symbol. The 189 rows
  with a verified path plus those 3 and 17 account for all 209, so the two
  non-defect categories are why no endpoint is missing rather than a gap. The
  snapshot is still not a zero-discrepancy report: 4 summary-only and 1
  differing remain. Of the 25 rows the previous snapshot reported as
  unresolved, 20 were generator artifacts; the remaining 5 were investigated
  individually, and 4 were correct SDK code the generator could not statically
  follow while 1 is the real path mismatch now reported as ⚠️.
- That `no OpenAPI schema on page` count is a label count, not a page count, and
  quoting it as pages is corrected here. 7 gRPC reference pages embed no OpenAPI
  schema, but only 3 rows carry the label: the generator evaluates the unmapped
  status first, so the 4 pages that the manifest also maps to no SDK symbol are
  counted as unmapped instead, which is what keeps the status table a partition
  of the 209 rows. A further 3 rows have a JSON block that yields no `path` and
  are also labelled unmapped, so 10 rows in total have no usable official path.
  This restores the reading the `2.1.5` entry below already gave ("Seven are
  gRPC pages carrying no OpenAPI schema").
- `CHANGELOG.md`'s `2.1.5` entry already recorded that 24 of the 29 non-exact
  states needed no SDK change; this entry is the status-label fix that
  `2.1.5` deferred.

### Noted

- Four live-blocked SDK defects found by static analysis on 2026-09-26 are now
  recorded with `file:line`, impact, minimal fix, and unblock requirement in
  `IMPLEMENTATION_STATUS.md` and `docs/implementation-status.md`. None is
  live-verified; nothing in that pass touched the network.
  - `brokerfd/client.go:43` routes the whole `brokerfd` package to the core host
    through `c.core.Do` instead of the documented Broker host
    (`internal/region/region.go:191`); the HK `broker` package routes correctly
    through `DoBroker` at `broker/client.go:63`. Unblocked by US sandbox
    credentials.
  - `brokerfd` retains 14 undocumented `/broker-fd/*` path literals while every
    other cached `broker-fd-api` page uses `/broker/...`;
    `brokerfd/assets.go:25` sends `/broker-fd/assets/summary` against a
    documented `GET /broker/assets/summaries/get` and is the single ⚠️ row in
    the snapshot. Unblocked by US sandbox credentials.
  - `broker.UpdateVirtualAccount` (`broker/accounts.go:66-68`) issues PUT with
    `account_id` as a query parameter and an undocumented `account_name` body
    field, while the documented endpoint is POST requiring `account_id` and
    `client_request_id` in the JSON body;
    `broker/accounts_test.go:151-160` currently certifies the wrong contract.
    Unblocked by a production or US-scoped Broker credential.
  - `data.GetDisplaySnapshot` (`data/display_quotes.go:28`, `:52`) uses GET
    `/openapi/market-data/stock/snapshot` where both official sources say POST
    `/market-data/stocks/snapshots/list`; the four sibling Display paths were
    aligned in `2b29c88`, the same commit that removed the unverified marker
    covering it. Unblocked by a paid Display Solution entitlement.

## [2.1.6] - 2026-09-26

Repository patch release fixing the nested coverage CI job. `v2.1.6` is a
repository Git tag, not a published Go-semver v2 module; the module stays on the
v1 import path by decision, so `v1.1.1` remains the newest installable version.

### Fixed

- The `nested-coverage` job introduced in `v2.1.4` failed for all four example
  modules because of two independent defects. GitHub artifact names may not
  contain `/`, so `coverage-${{ matrix.module }}` produced
  `coverage-examples/broker-probe` and the upload was rejected as an invalid
  name. Separately, `go test -coverprofile` writes no profile at all for a
  module with no test files, and the four example modules are single-file
  `package main` programs, so the upload then found no files. Each matrix entry
  now carries an artifact-safe label, the upload sets `if-no-files-found: ignore`,
  and the coverage report step is guarded on the profile existing so
  `go tool cover` is not run against a missing file. The `broker` coverage floor
  is unchanged and still gated at 75.0 percent, and it passed throughout.

Both `v2.1.4` and `v2.1.5` are tagged with a red CI run because of this defect.
Those tags are immutable and are left as an honest record; `v2.1.6` is the first
tag whose CI is fully green.

## [2.1.5] - 2026-09-26

Repository patch release closing three roadmap items that did not require live
access: Broker FD event metadata delivery, the first stream dispatch benchmarks,
and the classification of every non-exact endpoint reconciliation state.
`v2.1.5` is a repository Git tag, not a published Go-semver v2 module; the module
stays on the v1 import path by decision, so `v1.1.1` remains the newest
installable version. Not newly live-verified.

### Added

- `brokerfd/events.Client.OnDataEvent`, which delivers the whole `DataEvent`
  including the server-assigned `RequestId` and `Timestamp`. Those fields were
  already modelled and populated by `SubscribeResponse.ToDataEvent`, but
  `OnData`'s three-argument callback could not carry them, so they were discarded
  on delivery. The addition is purely additive: `OnData` keeps its exact
  signature and behavior, and the two registrations are independent.
- Stream dispatch benchmarks covering callback fan-out at 0 to 64 handlers,
  per-policy channel dispatch, subscriber-count scaling, subscription release, a
  mixed 8-handler and 32-subscriber case, and the head-of-line measurement. The
  repository previously contained no benchmarks at all.

### Changed

- `docs/broker-fd-us.md` no longer describes the missing request ID and timestamp
  as a current API limitation and instead documents both registrations, their
  independence, and the fact that `RequestId` is empty and `Timestamp` zero when
  the server does not supply them.
- The 29 non-exact endpoint reconciliation states are classified. Seven are gRPC
  pages carrying no OpenAPI schema, 13 are docgen manifest entries deliberately
  mapped to `-`, 5 name an SDK symbol whose path could not be resolved, and 4 are
  upstream disagreements between the OpenAPI JSON and the llms.txt summary path.
  24 of the 29 therefore need no SDK change, which the previous "25 unresolved"
  headline overstated. The generator status labels are identified as the real
  defect and are recorded for a separate change.

### Measured

Stream dispatch on 2026-09-26, windows/amd64, i7-13700, at `benchtime 1000x`.
These are measurements rather than thresholds; no latency objective has been
agreed, so the retain-or-change decision for synchronous dispatch is still open.

- With a subscriber registered first that is never read, a later ready subscriber
  is delayed until the blocked one is cancelled: 2396494 blocked-ns/op under
  `DropBlock`, against 0 blocked-ns/op for the identical shape under
  `DropOldest`.
- Channel fan-out is superlinear in subscriber count: 249 ns/op at one
  subscriber rising to 14860 ns/op at 64.

## [2.1.4] - 2026-09-26

Repository patch release covering nested CI gates and documentation accuracy.
`v2.1.4` is a repository Git tag, not a published Go-semver v2 module; the module
stays on the v1 import path by decision, so `v1.1.1` remains the newest
installable version. Not newly live-verified.

### Added

- A nested-module coverage job. Only `broker` is gated, at a 75.0% floor against
  its measured 80.8%, because it is the only nested module with library code to
  cover. The four example modules are single-file `package main` programs with no
  test files, so their coverage is structurally 0.0%; they are measured and
  uploaded as artifacts for visibility but not gated.
- `govulncheck` now runs as a matrix across all six modules instead of the root
  only, using the already-pinned `v1.8.0` so results stay reproducible. The
  vulnerability database is still fetched live from vuln.go.dev, so pinning the
  tool does not stale findings.
- A strict documentation gate that runs `mkdocs build --strict` on every push and
  pull request. The Docs workflow already built the site on push to `main`, but as
  a deploy job rather than a gate, and not on pull requests, so a broken link or
  nav entry previously could only be found after merge.

### Fixed

- `SECURITY.md` reported `v2.1.2` as the latest repository tag and its
  supported-versions table had no `v2.1.3` row, because the `v2.1.3` release
  updated the status files but not the security policy. A reader checking whether
  a version was supported would have been given a stale answer. Both the tag
  reference and the missing rows are corrected.
- The run index reported `v2.1.1` with post-release commits through `7d1489d`,
  which stopped being accurate once `v2.1.2` and `v2.1.3` were tagged.
- The install guidance described a pinned commit as "the current tree", which
  goes stale on every release. It now states that any commit at or after the
  desired tag can be pinned and marks the existing commit as a dated example.

### Changed

- The internal roadmap was corrected to the current release: the review range
  moved from `b3c647b..0ec3105` to `0e6f978..v2.1.4`, the pre-P1 uncommitted-tree
  risk is marked resolved, the coverage candidate is marked completed, and the
  root coverage figure moved from 71.6% to 73.6%.
- Two findings are recorded rather than acted on. The repository contains no
  `func Benchmark` at all, so the stream head-of-line candidate is greenfield.
  And the Broker FD `subscribeType` the SDK transmits contradicts Webull's
  published value: the documentation states that only `1` is supported, while the
  unexported config field is never assigned and so sends `0`, with no public
  option to correct it. That is documented as a risk and deliberately left
  unchanged, because the endpoint is US-only and has never been exercised live.
- Maintainer question 9 is resolved: the OpenTelemetry SDK modules stay direct
  `go.mod` requirements, and the decision is parked rather than reopened because
  no module publication is planned.

## [2.1.3] - 2026-09-26

Repository patch release covering CI signal honesty, dead-code removal, and
targeted coverage. `v2.1.3` is a repository Git tag, not a published Go-semver
v2 module; the module stays on the v1 import path by decision, so `v1.1.1`
remains the newest installable version. Not newly live-verified.

### Fixed

- The nightly live-sandbox workflow reported `success` while verifying nothing.
  When the required repository secrets were absent the guard set a skip flag,
  every live test step was skipped, and skipped steps still count as success, so
  "not verified" was indistinguishable from "verification passed". The guard now
  distinguishes three states: all secrets present runs the tests, none present
  skips and stays green so forks still succeed, and a partial configuration now
  fails. Each run publishes a job summary naming every secret and an explicit
  `RAN`, `SKIPPED`, or `MISCONFIGURED` verdict, and a new `require_secrets`
  input on `workflow_dispatch` turns an unconfigured clone into a hard failure so
  the pipeline can be proven to execute.

### Removed

- Deleted `internal/errs`, `internal/mqtt`, `internal/transport`, and the whole
  `internal/resilience` tree, together with `internal/resilience/breaker`,
  `ratelimit`, `retry`, and `clock`. All had zero importers, and the resilience
  packages were legacy duplicates of the better-covered `pkg/resilience`
  packages. About 1119 lines of duplicate production code are gone. This is safe
  because an `internal/` path cannot be imported outside this module.
- The `AGENTS.md` claim that `internal/errs` existed for backward compatibility
  is withdrawn; an `internal/` path is unreachable from outside the module, so it
  could never have served that purpose.
- The six client resilience tests that happened to live in
  `internal/resilience/integration_test.go` are relocated to
  `client/resilience_integration_test.go` rather than deleted. They are the only
  end-to-end coverage of the retry, circuit-breaker, and rate-limiter wiring
  around the request pipeline.

### Added

- Direct tests for the order reconciliation surface, which
  `trade.Client.ReconcileOrderStatus` and `ReconcileOrderState` call: the
  multi-step event paths, the deliberate no-op shapes, the authoritative terminal
  snapshot, invalid regressions, and concurrent reconciliation.
- Direct tests for the MQTT paho callback adapters and connect error paths,
  using in-package fakes so no broker is required.
- Direct tests for the `money.Money` comparison, rounding, shifting, and
  accessor surface.
- Direct tests for the `pkg/errors` rendering, identity-marker, and
  API-message-extraction behaviour, and for the `internal/region` accessors and
  domain table.

### Changed

- Aggregate root coverage measured 73.6% on 2026-09-26, up from 71.6%. The
  per-package work moved `pkg/domain/order` from 53.4% to 98.9%,
  `pkg/transport/mqtt` from 71.8% to 90.1%, `pkg/domain/money` from 71.0% to
  98.4%, `pkg/errors` from 79.2% to 97.4%, and `internal/region` from 75.0% to
  100%. Nested `broker` remains 80.8%. These are dated measurements, not
  behavior guarantees.

## [2.1.2] - 2026-09-26

Repository patch release covering dependency, CI, and documentation maintenance
after `v2.1.1`. `v2.1.2` is a repository Git tag, not a published Go-semver v2
module; the root module path remains `github.com/shing1211/webullapi4go` and
stays on the v1 import path by decision, so the module proxy serves only the
`v1.x` line and `v1.1.1` remains the newest installable version. Consumers reach
this work by pinning a commit. This release was not newly live-verified.

### Changed

- Bumped the OpenTelemetry Go modules to v1.46.0 across the root module and every
  nested module (`broker/` and the four example modules). All five modules —
  `go.opentelemetry.io/otel`, `otel/metric`, `otel/trace`, `otel/sdk`, and
  `otel/sdk/metric` — are kept on the same version because they are released in
  lockstep and a mixed set is untested. `github.com/go-logr/logr` moves to
  v1.4.4 as a transitive requirement.
- Bumped GitHub Actions: `upload-artifact` v4 to v7 and `golangci-lint-action`
  v8 to v9 in CI, and `deploy-pages` v4 to v5, `setup-python` v5 to v7, and
  `upload-pages-artifact` v4 to v5 in the docs workflow. `upload-artifact` v4
  is being retired by GitHub. Versions v6 and later of the affected actions run
  on Node.js 24 and require runner v2.327.1 or newer; every job uses a
  GitHub-hosted runner.
- Pinned `govulncheck` to v1.8.0 in the CI workflow and the `make vuln` target,
  which previously resolved `@latest` and made scan results non-reproducible.
  The vulnerability database is still fetched live from vuln.go.dev, so pinning
  the tool does not stale reported findings.
- Grouped OpenTelemetry updates in `.github/dependabot.yml` so the lockstep
  modules arrive in a single pull request. The previous configuration produced
  three byte-identical pull requests, each bumping the core modules while
  leaving the SDK modules behind.

### Documentation

- Sandbox test credentials are no longer inlined in hand-written documentation.
  `docs/sandbox.md` and `AGENTS.md` now link the published Webull test-accounts
  page instead, so a rotated or retired shared account cannot leave a stale copy
  behind. The generated pages under `docs/webull-api/**` are unchanged and still
  mirror Webull's published material verbatim.
- Recorded the module-path decision. Documentation previously described v2 module
  publication as deferred, which implied a pending choice; it is declined. The
  documentation now states the consequence, that `v2.x` tags are repository-only
  and `v1.1.1` is the newest installable version, and gives the working install
  recipe that pins a commit.
- `SECURITY.md` now lists `v1.1.1` as the supported installable line and
  reclassifies the `v2.x` tags as repository-only.

## [2.1.1] - 2026-09-25

Repository patch release of the current hardening. `v2.1.1` is a repository Git
tag, not a published Go-semver v2 module; the root module path remains
`github.com/shing1211/webullapi4go` and stays on the v1 import path by
decision, with no `/v2` migration planned. The hardening was not newly
live-verified.

### Added

- `client.Client.ObservabilityConfig` exposes the shared, read-only telemetry
  configuration inherited by service clients.
- `pkg/errors.NewSentinel` creates identity-specific semantic sentinels while
  preserving category checks through `errs.Is`.
- `pkg/observability.SafeErrorText` reduces errors to credential-free,
  stable telemetry text without copying API bodies, status messages, or causes.
- `events` and `brokerfd/events` now emit one client-kind OpenTelemetry span
  per gRPC stream attempt, `event_stream_attempts` and
  `event_stream_attempt_duration` metrics, structured start/done logs, and
  correlation/trace metadata. Credentials are not included in telemetry.
- `trade.Client` exposes account-scoped OMS inspection and reconciliation:
  `GetTrackedOrder`, `TrackedOrder`, `TrackedOrderState`, `TrackedOrders`,
  `ApplyOrderEvent`, `ReconcileOrderStatus`, and `ReconcileOrderState`.
- `pkg/domain/order.Machine` and `order.Order` now support concurrency-safe
  status reconciliation, stale non-terminal no-ops, terminal-state protection,
  and distinct failed cancel/modify/reject scene events.
- `examples/account-monitor`: a read-only, single-worker balance monitor with
  bounded requests, signal-aware cancellation, and a consecutive-failure stop.

### Changed

- Ordinary `pkg/errors` values continue to match by category, while semantic
  sentinels now match only themselves or wrappers preserving their identity.
  This applies to subscription expiry, event/MQTT connection limits, broker
  refusal, circuit-open, and explicit-token-required sentinels.
- HTTP 417 remains mapped to `INVALID_TOKEN` for compatibility, while tests and
  documentation now make clear that the same status can carry business
  validation messages such as invalid symbols or unsupported categories.
- `WithMeterProvider` and `WithResiliencePreset` are order-independent, and an
  explicitly supplied `WithBreaker` remains caller-owned.
- Stream handlers and channel subscribers dispatch synchronously in
  registration order. A full `DropBlock` subscriber therefore has documented
  head-of-line behavior until cancellation or terminal close.
- Trading `events.Close` cancels every active `Run`; each run returns
  `context.Canceled` without a shutdown `OnError` callback.
- `client.Do`, `client.DoBroker`, and `client.DoStream` now share one
  per-attempt request pipeline for rate limiting, circuit breaking, hooks,
  interceptors, signing, correlation, logging, tracing, and metrics. Hook
  attempt numbers are one-based and increase across retries.
- `client.Do` and `client.DoBroker` reuse one correlation ID across retries;
  `DoBroker` and `DoStream` also learn the bounded clock offset from successful
  response `Date` headers when clock correction is enabled.
- REST requests inject the configured W3C trace propagator. `json.RawMessage`
  request bodies remain byte-for-byte stable across retries.
- The default HTTP client now clones and tunes Go's default transport for idle
  connections, TLS handshakes, and expect-continue timing; explicit
  `WithHTTPTransport` still wins.
- `trade.WithAutoClientOrderID(true)` derives stable IDs for the same logical
  place/batch request, does not mutate the caller's request slice, and keeps
  repeated requests idempotent. Successful batch results are tracked locally.
- The trade order registry is keyed by `(account_id, client_order_id)`.
  Successful replace/cancel actions advance matching local state; failed
  actions leave it unchanged.
- Order guardrail failures retain the historical outer `invalid_config` code
  and also wrap `pkg/errors.ErrOrderGuardrail`.
- `stream.Close` is terminal, channel cancel functions are idempotent, blocked
  channel sends unblock on cancellation/close, and reconnect replay is
  serialized with subscription mutations.

### Fixed

- `client.Close` now forwards `CloseIdleConnections` through the token-injection
  transport wrapper, so caller-supplied transports release idle connections.
- MQTT `Connect` no longer starts a `Token.Wait` waiter goroutine, checks an
  already-cancelled context before connecting, and rejects use after `Close`.
  MQTT close is idempotent and suppresses late callbacks/messages.
- Stream state changes use compare-and-swap expected-state transitions, so
  stale health ticks and late data cannot overwrite reconnect or terminal-close
  state.
- REST, MQTT, and gRPC failure telemetry now records sanitized error text
  instead of raw response bodies, gRPC status messages, or wrapped causes.
- `DropOldest` now removes the oldest unread value. Cancelling or closing while
  dispatch is blocked no longer deadlocks the stream.
- Stream health only treats quote/snapshot/tick messages as fresh data; notices
  and echo heartbeats neither prevent nor falsely recover `StateDegraded`.
  Repeated stale checks no longer restore `StateConnected` without new data.
- OMS scene mapping no longer treats failed cancel, modify, or reject operations
  as successful transitions. Status reconciliation never regresses a terminal
  order and ignores stale non-terminal snapshots.
- Correlation context lookup is type-safe, response status/latency hooks receive
  real attempt values, and stream interceptors cannot accidentally return a
  nil successful response.

### Tests

- Added offline regression coverage for the shared request pipeline, raw-body
  retries, clock correction, cancellation cleanup, account-scoped OMS tracking,
  status reconciliation, concurrent order machines, stable auto IDs, category
  versus semantic-sentinel matching, HTTP 417 diagnostics, MQTT/channel
  shutdown, compare-and-swap stream state, health recovery, resubscription
  serialization, deterministic drop policies, `money.Money` wire forms, and
  Trading/Broker FD event telemetry.
- Added direct offline tests for the public resilience, transport, shared type,
  and `webull` alias packages. Data and both event packages now use strict
  goroutine-leak checks; event tests now cancel and wait for completed runs.
- The module-aware offline race/vet checks recorded for this run passed on
  2026-09-25. `make cover` measured 71.6% aggregate root coverage and 80.8% in
  the nested `broker/` module; these are measurements, not behavior
  guarantees.
- No live verification was added for the v2.1.1 repository-tagged hardening.

### Documentation

- Reconciled the architecture and status docs with the preserved root service
  packages and public `money.Money`; marked the full `pkg` relocation,
  root-service shim plan, raw `decimal.Decimal` DTO migration, and `sync.Pool`
  decisions as superseded or closed.
- Added OMS reconciliation, compare-and-swap streaming state/health/channel
  behavior, category versus semantic-sentinel matching, HTTP 417 caveats,
  MQTT/Connect cancellation semantics, exact OTel instrument contracts,
  `SafeErrorText`, event all-runs shutdown behavior, and explicit
  implemented/offline-tested/live-verified/blocked distinctions.
- Corrected generated reconciliation summaries so they no longer claim zero
  path discrepancies while the generated report contains summary-only and
  unresolved paths.
- The current error, request, OMS, streaming, event-telemetry, testing, and
  documentation hardening is tagged in repository `v2.1.1`; the module stays on
  the v1 import path by decision, so `v2.x` tags are repository-only.

## [2.1.0] - 2026-09-24

Phase 13 vet warnings fixed.

### Fixed

- `pkg/observability/otel.go`: `Tracer`, `Meter`, `InjectTraceContext`, and
  `ExtractTraceContext` now use `*Config` pointer receivers, eliminating the
  `passes lock by value` vet warnings.
- `client/config.go`: `Validate` now uses `*Config` pointer receiver.
- `client/config.go`: `client.Config.otel` field changed from embedded
  `observability.Config` to `*observability.Config` (pointer), so that
  returning `Config` by value does not copy the embedded `sync.RWMutex`.
- `client/client.go`: `Client.cfg` changed to `*Config` (pointer) to avoid
  copying the mutex when constructing a `Client`; `Config()` returns
  `*c.cfg` by value-copy of the pointer.

## [2.0.9] - 2026-09-24

Phase 9 structured errors.

### Added

- `pkg/errors/errors.go`: three new error codes (`CodeNotInitialized`,
  `CodeValidation`, `CodeOrderGuardrail`) and four new sentinels
  (`ErrNotInitialized`, `ErrValidation`, `ErrOrderGuardrail`,
  `ErrSubscriptionExpired`, `ErrConnectionLimitExceeded`).
- `pkg/errors/errors.go`: `Error.Is` now compares by pointer equality first,
  ensuring `ErrX.Is(ErrX)` returns true even when `ErrX` is a typed
  `*Error`.

### Changed

- `client/request.go`: `ErrCircuitOpen` is now `*pkgerrors.Error` with
  `CodeTransport`; existing `errors.Is(err, client.ErrCircuitOpen)` checks
  continue to work.
- `pkg/transport/mqtt/mqtt.go`: `ErrConnectionRefused` and `ErrConnectionLimit`
  are now `*pkgerrors.Error` with `CodeTransport`; `errors.Is(err,
  mqtt.ErrConnectionLimit)` and `errors.Is(err, mqtt.ErrConnectionRefused)`
  continue to work.
- `trade/orders.go`: `enforceGuardrails` now uses `errs.Wrap` instead of
  `errs.New` when prefixing batch-order index, preserving the full error chain.
- `events/client.go`, `brokerfd/events/events.go`: subscription-expired errors
  now wrap `ErrSubscriptionExpired` so callers can check
  `errors.Is(err, errs.ErrSubscriptionExpired)`.

## [2.0.8] - 2026-09-24

Phase 8 context hygiene.

### Fixed

- `pkg/transport/mqtt/mqtt.go`: `Connect()` now uses `sync.WaitGroup` so the
  `token.Wait()` goroutine exits promptly when the context is cancelled, instead
  of leaking until the MQTT stack processes the disconnect.
- `stream/client.go`: `resubscribeContext()` panics on a nil `resubCtx` instead of
  silently using an uncancellable `context.Background()`.

### Documentation

- `stream/client.go`, `stream/channels.go`: added comments explaining that OTel
  metric `Add` calls use `context.Background()` intentionally (fire-and-forget;
  the metric SDK records synchronously and cannot block).

## [2.0.7] - 2026-09-24

Phase 6.3 OTel metrics hooks.

### Added

- `pkg/observability/otel.go`: lazy OTel instrument accessors on `Config`:
  `ClientLatencyHistogram()`, `BreakerTransitionsCounter()`,
  `StreamReconnectsCounter()`, `StreamChannelDropsCounter()`.
- `pkg/resilience/breaker/breaker.go`: `WithMeter(m metric.Meter)` option and
  `transitionCounter` instrument. `transitionTo()` records `from`/`to` state
  attributes on every breaker transition.
- `client/request.go`: `request_latency` histogram recorded on every `attempt()`
  call with `http.route` attribute.
- `stream/client.go`: `streamMetrics` with reconnect and per-topic drop
  instruments in scope `webullapi4go/stream`. `handleReconnecting()` increments
  the `reconnects` counter. `WithMeter(m metric.Meter)` is available in
  `stream/option.go`.
- `stream/channels.go`: per-topic drop counters feed the `channel_drops`
  instrument with a `topic` attribute.

## [2.0.6] - 2026-09-24

Phase 5 streaming engine hardening + Phase 6 observability.

### Added

- `stream/state.go`: new `State` type with six values (`Disconnected`,
  `Connecting`, `Connected`, `Reconnecting`, `Degraded`, `Closed`). `Client.State()`
  returns the current state; `Client.OnStateChange(func(prev, next State))` fires
  on every transition.
- `stream/option.go`: `WithHealthWatchdog(interval)` — enables a background
  watchdog that transitions the connection to `StateDegraded` when no data
  message (quote/snapshot/tick) arrives within the configured interval, and
  recovers to `StateConnected` when a message arrives.
- `stream/client.go`: `Client.OnReconnecting(func())` handler invoked when the
  client begins attempting to reconnect after a connection loss.
- `stream/channels.go`: new channel-mux API with configurable backpressure.
  `Client.SubscribeQuoteChan`, `Client.SubscribeSnapshotChan`,
  `Client.SubscribeTickChan` each return a `<-chan T` with a per-subscription
  cancel function. `ChannelConfig.Policy` supports `DropBlock` (default,
  blocks sender), `DropOldest` (drops oldest unread), and `DropSample`
  (probabilistic). Drop counters track discards per channel.
- `client/request.go`: each `Client.Do` and `Client.DoBroker` call now establishes
  an OpenTelemetry span (`SpanKindClient`) scoped to the configured
  `TracerProvider`. Span attributes include `http.method`, `http.route`,
  `webull.attempt`, and `http.duration_ms`. Correlation ID (`X-Correlation-ID`
  header) is generated per-request using `crypto/rand` and threaded through the
  context.
- `client/request.go`: `WithCorrelationID(ctx, id)` and
  `CorrelationIDFromContext(ctx)` helpers for explicit correlation ID management.
- `client/option.go`: `WithLogger(*slog.Logger)` — configures per-request
  structured log output (request start and completion with method, path,
  correlation ID, latency, and error).
- `client/option.go`: `WithTracerProvider`, `WithMeterProvider`, and
  `WithPropagator` options supply real OpenTelemetry providers, replacing the
  no-op defaults.
- `pkg/observability/otel.go`: new package providing API-only OTel integration.
  `Config.Tracer()`, `Config.Meter()`, `Config.InjectTraceContext`, and
  `Config.ExtractTraceContext` use the configured providers; all types default
  to the global no-op implementations.
- `pkg/observability/otel.go`: `SpanAttributes(method, path, attempt)` and
  `SpanName(method, path)` helpers for consistent span naming and attribute
  sets.

## [2.0.5] - 2026-09-24

Phase 4 production hardening: interceptor pipeline, initial OMS integration,
and typed order builders.

### Added

- `client/option.go`: `Interceptor` func type and `WithInterceptor(Interceptor)`
  option — composable request pipeline interceptors that run after rate-limiting
  and circuit-breaking but before signing and send.
- `client/option.go`: `Hooks` struct with `OnRequest`, `OnResponse`, `OnError`,
  and `OnLatency` callbacks — observability integration point for Phase 6.
- `client/request.go`: `attempt()` refactored to run an ordered interceptor
  chain; hooks fire on every request attempt including retries.
- `pkg/domain/order/order.go`: new `Order` struct embeds `PlaceOrderResult`
  plus `AccountID` and a local `*Machine` state machine; `SceneTypeToEvent`
  maps Webull gRPC scene types to domain events.
- `trade/client.go`: added `orderRegistry` map and `registerOrder`/`getOrder`
  helpers for local order state tracking.
- `trade/orders.go`: `PlaceOrder` now returns `*order.Order` (not
  `*PlaceOrderResult`); the order is registered with `StatePending` on
  placement. `PlaceOrderResult` is embedded so `order.OrderID` and
  `order.ClientOrderID` remain accessible. `BatchPlaceOrder` unchanged.
- `trade/order_actions.go`: `CancelOrder` and `ReplaceOrder` check local order
  state before sending; terminal orders (filled, cancelled, failed, expired)
  return `errs.CodeInvalidTransition` without an API call.
- `trade/orders.go`: `NewPlaceOrderRequest(accountID, orders...)`,
  `NewEquityOrder(symbol, side, qty)`, and `EquityOrderBuilder` fluent API —
  typed request constructors with US equity defaults.
- `trade/order_actions.go`: `NewCancelOrderRequest(accountID, clientOrderID)` and
  `NewModifyOrderRequest(accountID, clientOrderID)` convenience constructors.
- `pkg/errors/errors.go`: added `CodeInvalidTransition` and `ErrInvalidTransition`
  for local state validation failures.

## [2.0.4] - 2026-09-24

Phase 3 production hardening: clock-drift correction, idempotency helpers,
transport tuning, and resilience presets.

### Added

- `client/option.go`: `WithClockDriftCorrection(bool)` — learns the clock offset
  between the client and Webull server from the `Date` response header and applies
  it to subsequent request signing timestamps, clamped to ±5 minutes.
- `client/option.go`: `WithHTTPTransport(*http.Transport)` — sets a custom HTTP
  transport on the client, applied after `WithHTTPClient` so callers can tune
  connection pooling and TLS without replacing the whole client.
- `client/option.go`: `WithResiliencePreset(ProductionPreset)` — applies a
  production-ready resilience composition: per-path rate limiter (10 req/s, burst
  20), circuit breaker (5-failure threshold, 30 s cooldown), and exponential
  backoff retry with full jitter (base 200 ms, cap 2 s, up to 3 attempts).
- `pkg/resilience/retry/`: `WithFullJitter(bool)` option — enables the full
  jitter formula `rand(0, min(cap, base·2ⁿ))` giving uniformly distributed
  delays in `[0, cap]` instead of ±20% jitter around the exponential.
- `trade/orders.go`: `NewClientOrderID()` — generates a fresh 20-character
  client-order identifier using `crypto/rand`, suitable for idempotent order
  placement.
- `trade/orders.go`: `ClientOrderIDFrom([]byte)` — derives a deterministic
  32-character client-order identifier from arbitrary content via SHA-256.
- `trade/orders.go`: `ValidClientOrderID(string) bool` — exported validation
  helper for the client-order identifier character set.
- `trade/option.go`: `WithAutoClientOrderID(bool)` — configures `PlaceOrder`
  and `BatchPlaceOrder` to auto-generate and assign a `NewClientOrderID` to
  each order whose `ClientOrderID` is empty.

### Documentation

- Reorganized the docs site nav into Guides / API Reference / Official Webull
  Docs / Coverage / Architecture Decisions, and renamed pages for consistency
  (`broker.md` → `broker-hk.md`, `brokerfd.md` → `broker-fd-us.md`,
  `webull-api/crypto.md` → `webull-api/market-data-crypto.md`); old URLs are
  preserved with `mkdocs-redirects`.
- Marked every generated file with a "do not edit" banner; refreshed stale
  v1.0-era content (feature matrices, hub path status, provisional warnings) for
  v1.1.0; consolidated the official Webull doc links into `AGENTS.md`; merged
  `docs/adr/README.md` into `docs/adr/index.md`; standardized READMEs across all
  17 example directories.
- Fixed 80+ doc-code discrepancies across all endpoint docs: wrong function names
  (`GetOptionChain`→`GetOptionContracts`, `GetStockProfilesList`→`GetStockProfilesV3`),
  parameter mismatches, outdated code examples, and missing type documentation.
- Replaced the v0.1 coverage table in `market-data.md` with a comprehensive
  current table of all ~105 functions organized by group.
- Added SDK compatibility notes to all Display Solution endpoints documenting
  differences between the SDK implementation and the official OpenAPI spec.
- Added shared patterns reference (`patterns.md`) covering client construction,
  options, pagination, numeric strings, error handling, display service, and
  streaming patterns.
- Enhanced `errors.md` with transient vs permanent error classification table
  and rate-limit retry pattern.
- Added architecture diagram, "What's new in v1.1" section, and cross-reference
  links to `docs/index.md`.
- Added prerequisites to streaming, trading, and events docs.
- Documented previously undocumented types and functions across authentication,
  streaming, events, trading, broker-fd, and connect-api docs.
- Added prerequisites admonitions to all 10 guide pages (authentication,
  market-data, streaming, trading, events, fundamentals, errors, patterns,
  broker-hk, broker-fd-us).
- Added error handling tips to 5 guide pages (market-data, streaming, trading,
  events, fundamentals).
- Created glossary (`glossary.md`) with 30+ SDK-specific terms.
- Added ASCII token lifecycle diagram to authentication guide.
- Restructured API reference index (`webull-api.md`) with service-level grouping.
- Added common first-call failures table to getting started guide.
- Added `connect` and `display` packages to Go Packages reference (`api.md`).
- Added Documentation Standards section to `CONTRIBUTING.md`.
- Fixed broken Options table in streaming guide.

## [2.0.3] - 2026-09-23

### Fixed

- `stream/` and `data/`: `go.uber.org/goleak` false positives from
  `net/http.(*http2clientConnReadLoop).run` goroutines left after paho-mqtt
  WebSocket disconnect. Added `goleak.IgnoreAnyFunction` filters to both
  packages' `TestMain`. The goroutines are cleaned up asynchronously by the Go
  runtime and are not an SDK leak.

## [2.0.2] - 2026-09-23

### Changed (breaking)

- All numeric string fields in `brokerfd/` converted to `*money.Money`
  (optional/request-side) or `money.Money` (required/response-side).
  Approximately 30 fields affected across `assets.go`, `orders.go`,
  `funding.go`, `activity.go`, `journals.go`, and `instruments.go`.
  JSON serialization is preserved (decimal strings).

## [2.0.1] - 2026-09-23

Phase 2 production hardening: public API restructuring and type-safe numeric fields.

### Added

- `pkg/errors`: typed errors (`Error` with `Code` and `Message`) extracted from
  `internal/errs` to the public `pkg/errors/` package. A deprecated shim at
  `internal/errs/errs.go` preserves backward compatibility with existing imports.
- `pkg/transport/http.go`: public `Doer` interface (`Do(ctx, method, path, query,
  req, resp) error`) extracted from `client.Client`. New `pkg/transport/mqtt/`
  sub-package for MQTT transport types. A deprecated shim at
  `internal/transport/http.go` preserves backward compatibility.
- `pkg/resilience/{breaker,clock,ratelimit,retry}/`: public resilience primitives
  extracted from `internal/resilience`. A deprecated shim at
  `internal/resilience/resilience.go` preserves backward compatibility.
- `pkg/domain/money`: type-safe `Money` wrapper around `shopspring/decimal`
  (`github.com/shopspring/decimal v1.4.0`). Provides JSON round-trip fidelity,
  decimal arithmetic, `MarshalJSON`/`UnmarshalJSON`, and `Rat()` for rational
  arithmetic. Package-level helpers: `NewFromString`, `Must`, `ParseMoney`,
  `ParseDecimal`, `Zero`.
- `pkg/domain/order`: OMS state machine (`StateMachine`, `ApplyEvent`,
  `FromWebullStatus`) for order lifecycle tracking.
- `webull/`: thin type-alias convenience package re-exporting selected
  `client` and `trade` constructors and types. The root service packages remain
  canonical; this is not an aggregate service facade.

### Changed (breaking)

- `pkg/errors` is now the canonical import path. Old `internal/errs` is
  deprecated and will be removed in a future release.
- `pkg/transport` is now the canonical import path for `Doer`. Old
  `internal/transport` is deprecated.
- `pkg/resilience` is now the canonical import path. Old `internal/resilience`
  is deprecated.
- All numeric string fields converted to `*money.Money` (optional/request-side)
  or `money.Money` (required/response-side) throughout `trade/` and `data/`
  packages. Approximately 250 fields affected. JSON serialization is preserved
  (decimal strings); the change is transparent for most callers.
- `shopspring/decimal` is now a direct dependency (v1.4.0). Previously it was
  an indirect dependency.
- `money.Rat()` now delegates to the shopspring library's built-in `Rat()`
  method, fixing a bug where fractional quantities were incorrectly computed.

### Deprecated

- `internal/errs/errs.go`: use `pkg/errors` instead.
- `internal/transport/http.go`: use `pkg/transport` instead.
- `internal/resilience/resilience.go`: use `pkg/resilience` instead.
- Direct use of `client.Client` remains canonical; `webull.New` and
  `webull.Client` are optional aliases for callers that prefer one import.

## [2.0.0] - 2026-09-23

Phase 2 architecture: introduce the public SDK foundations and the optional
core-client alias package. The later decision to retain the root service
packages supersedes any broader service-relocation reading of this historical
tag; the current tree keeps those services at the repository root.

### Added

- `pkg/errors`: public typed errors, codes, and sentinels, with a deprecated
  `internal/errs` compatibility shim.
- `pkg/transport`: public HTTP and MQTT transport foundations.
- `pkg/resilience`: public retry, rate-limit, circuit-breaker, and clock
  primitives.
- `pkg/domain/money`: exact decimal `Money` support for JSON financial values.
- `pkg/domain/order`: the public order lifecycle state machine.
- `webull`: optional aliases for selected `client` and `trade` constructors and
  types; it is not an aggregate service facade.

### Changed

- The public foundation packages became the canonical import paths while the
  root service packages remained available as the service layer.

## [1.1.1] - 2026-09-23

### Dev/Tooling

- Added `Makefile` targets: `build`, `test`, `test-race`, `cover`, `lint`,
  `fuzz`, `vuln`, `docs`.
- Enabled `gosec` in `.golangci.yml`; excluded `internal/auth`
  (protocol-mandated HMAC-SHA1/MD5 body digest) and
  `internal/resilience/retry` (intentional `math/rand` jitter).
- Expanded CI matrix to ubuntu/macos/windows; added coverage job
  (baseline measurement) and `govulncheck` job.
- Added `.github/dependabot.yml` (weekly gomod + github-actions updates).
- Added `go.uber.org/goleak` for goroutine-leak detection; `TestMain` with
  `goleak.VerifyTestMain` for `client`, `stream`, `data`, `internal/mqtt`.
- Added `data/fuzz_test.go` with `FuzzDecodeQuote`/`Snapshot`/`Tick`
  JSON deserialization fuzz tests.

## [1.1.0] - 2026-09-22

Full SDK parity with the official Webull OpenAPI.

### Added

- `connect/`: OAuth 2.0 authorization-code flow (`AuthorizationURL`,
  `CreateToken`) with `GrantTypeAuthorizationCode`/`GrantTypeRefreshToken`.
- `data/`: crypto (`GetCryptoBars`, `GetCryptoSnapshot`, `GetCryptoInstruments`),
  Display event contracts (`GetEventContractTags`, `GetEventContractEventsList`,
  `GetEventContractMilestones`, `GetEventContractSeriesList`,
  `GetEventContractSportsFilters`, `GetEventGameStats`, `GetEventLiveData`,
  `GetEventMarketBars`, `GetEventMarketBarsByEvent`, `GetEventMarketDepth`,
  `GetEventMarketSnapshot`), and fund extras (`GetFundPerformance`,
  `GetFundHoldings`, `GetFundRating`, `GetFundSplits`, `GetFundFiles`,
  `GetFundAllocation`).
- `display`: `Service.RefreshClientToken`.

### Changed

- Aligned 69 SDK paths to the official OpenAPI definition across `broker`,
  `brokerfd` and `data`. `broker.GetTradeCalendar` is now `GET` with query
  parameters; Broker FD account-update, order-replace and order-cancel are now
  `POST`.
- Removed all 41 provisional TODO markers. The SDK now covers every documented
  endpoint (0 documented-only gaps in `docs/reconciliation.md`). Path
  reconciliation remains tracked separately by the generated report.

## [1.0.3] - 2026-09-22

Documentation release: verbatim Webull master reference, an SDK↔API
reconciliation, a reproducible doc generator, and a path probe.

### Added

- `docs/webull-api.md` and `docs/webull-api/**`: an SDK-mapped reference for
  every documented endpoint, plus **verbatim** Webull master guides
  (`master-guides.md`) and per-area OpenAPI definitions (`reference/*.md`).
- `docs/reconciliation.md`: maps every documented endpoint to its SDK function
  and status (match / differs / not implemented).
- `tools/webull-docgen/`: generator CLI (`docgen.py
  reference|master|reconciliation|all`) that fetches Webull's machine-readable
  `.md` pages and renders the docs above.
- `examples/path-probe/`: env-gated probe comparing SDK paths against the
  official OpenAPI paths for the endpoints where they disagree.

### Changed

- `IMPLEMENTATION_STATUS.md`: coverage/gap section, Known Issues 12–13 (path
  drift), version history updated to v1.0.3.
- `mkdocs.yml`: added the Webull API Reference section.

## [1.0.2] - 2026-09-22

Reconciled the codebase against the official Webull API.

### Removed (breaking)

- Undocumented functions: crypto data (`data/crypto_data.go`), screener v2
  (`data/screener_v2.go`), option expirations/chains, and the redundant HK
  futures market-data variants.

### Changed

- `data.GetOptionContracts` now uses the Trading API path
  `/trading/instruments/options/contracts/list`.

## [1.0.1] - 2026-09-22

HK sandbox probe: confirmed futures product-codes path, fixed FuturesInstrument.Unit
flexible type to handle numeric API responses, fixed client_order_id length overflow
in options-multi-leg example, and updated TODO markers with HK sandbox findings.

### Fixed

- `data/futures.go`: `FuturesInstrument.Unit` now uses a `StringOrNumber` flexible
  type that handles both string (`"1-index points"`) and numeric (`1`) JSON values
  from the live API.
- `examples/options-multi-leg/main.go:193`: `client_order_id` now uses `Unix()` instead
  of `UnixNano()` to stay within the 32-character limit.

### Tests Added

- `data/futures_test.go`: `TestGetFuturesInstrumentsNumericUnit` verifies numeric unit
  deserialization.

### Changed

- `trade/types.go`: `TODO(t8)` comment updated to reflect HK sandbox finding that
  multi-leg strategies are rejected (only SINGLE accepted).
- `trade/options.go`: `TODO(t8)` comments updated to reflect HK sandbox findings.
- `data/futures_market.go`: header comment notes that product-codes path is confirmed
  while market data paths remain unconfirmed.

## [1.0.0] - 2026-09-22

Futures market-data bug fix and two new probe examples. The release also
recorded the remaining provisional audit items that required US sandbox access.

### Fixed

- `data/futures_market.go`: added the missing `Category` field to all five
  futures market-data query structs, fixing the hard-coded US-futures query
  behavior that prevented HK futures queries.

### Added

- `examples/futures-probe/`: futures discovery and market-data probe.
- `examples/options-multi-leg/`: multi-leg options preview probe.
- The 49 remaining provisional audit items were documented as requiring US
  sandbox credentials for verification.

## [0.9.2] - 2026-09-22

Critical path corrections for the Broker API HK package. All paths were wrong
— the SDK used `/openapi/v1/broker/...` but the official API uses `/broker/...`
directly. Sandbox returned 404 not because the API was missing but because every
path was wrong. Additionally, request/response shapes were updated to match the
official OpenAPI specification.

### Fixed

- `broker/accounts.go`: paths corrected to `/broker/accounts/virtual-accounts/{list,get,create,update}`;
  `VirtualAccount` struct updated to official fields; `ListVirtualAccounts` now handles
  wrapped `{data:[...]}` response.
- `broker/activities.go`: path corrected to `/broker/activities/list`.
- `broker/assets.go`: paths corrected to `/broker/assets/{balances/get,positions/list}`;
  `Balance` struct updated to official nested shape with `total_cash_balance`,
  `total_market_value`, `total_unrealized_profit_loss`; `GetPositions` handles
  wrapped `{data:[...]}` response.
- `broker/instruments.go`: paths corrected to `/broker/instruments/{stocks/list,stock-locate/get,corporate-actions/get}`.
- `broker/orders.go`: paths corrected to `/broker/orders/{preview,place,replace,cancel,get,history,open}`;
  `ReplaceOrder` and `CancelOrder` now use POST with JSON body instead of PUT/DELETE
  with query params; `PreviewOrderRequest` adapted to official nested wire format.
- `broker/funding.go`: paths corrected to `/broker/funding/fx-rates/get`,
  `/broker/funding/fx-exchanges/{create,get}`,
  `/broker/funding/instant-fx/{create,get}`, `/broker/funding/instant/{create,get}`;
  `FXRate` struct updated (`fx_rate`, `rate_effective_time`, `rate_expire_time`);
  detail queries now require `client_request_id` + `account_id`.
- `broker/journals.go`: paths corrected to `/broker/journals/{cash-journals/{create,get},position-journals/{create,get}}`;
  detail queries now require `client_request_id` + `account_id`.
- `broker/masterdata.go`: path corrected to `/broker/master-data/trade-calendar/query`;
  changed from GET with query params to POST with JSON body.
- `broker/eventcontracts.go`: paths corrected to `/broker/event-contracts/{categories/list,series/get,events/get,instruments/get}`.
- `examples/broker-probe/main.go`: updated to use new field names and paths.

### Noted

- Broker API HK (`/broker/...`) returns `401 ROUTE_NOT_PERMITTED` on the HK
  sandbox — the app does not have Broker API scope enabled in the sandbox.
  Paths are confirmed correct (no more 404); the 401 means auth/scope issue.

## [0.9.1] - 2026-09-21

Bug fixes found during remaining endpoint verification sweep.

### Fixed

- `data/watchlist.go`: `AddWatchlistInstruments`, `RemoveWatchlistInstruments`,
  `UpdateWatchlistInstruments`, `UpdateWatchlist`, and `DeleteWatchlist` now
  handle bare boolean (`true`/`false`) responses from the API using a
  `BoolOrSuccess` type that accepts both raw bool and `{"success":bool}` JSON.
- `broker/client.go`: `DoBroker` method added to route Broker API calls to the
  correct `broker-api.sandbox.webull.hk` base URL instead of the market data URL.

### Added

- `examples/watchlist-cmd/`: Watchlist CRUD example (create, add instruments,
  update, remove, delete) guarded by `WEBULL_WATCHLIST_TEST=1`.
- `examples/broker-probe/`: Broker HK read-only endpoint probe program.

### Noted

- Broker API HK (`/openapi/v1/broker/...`) returns `404 Route Not Found` on the
  HK sandbox — the endpoint group is not available in the sandbox environment.
  Broker HK remains unverified pending production or US sandbox access.

## [0.9.0] - 2026-09-22

GoDoc coverage, HK options stubs, HK futures market data, new examples, graceful
credential handling, and full sandbox verification.

### Added

- `brokerfd/brokerfd.go`: GoDoc on all 12 files and ~80 exported identifiers.
- `brokerfd/events/events.go`, `brokerfd/events/option.go`, `brokerfd/events/sign.go`:
  GoDoc on all event types and sign functions.
- `examples/brokerfd/`: Broker FD US read-only endpoint probe (accounts, orders,
  assets, instruments, funding, activity, journals, master data, agreements, documents).
- `examples/brokerfd-events/`: Broker FD gRPC event subscription probe (order, option,
  position streams).
- `examples/options/`: HK options discovery probe (expirations, option chain).
- `examples/`: graceful credential handling in all 9 existing examples — no more
  panic on missing env vars; instead a descriptive message and zero-value client.
- `data/futures.go`: `FuturesCategoryCN` constant for CN futures queries.
- `data/futures_market.go`: `GetHKFuturesTick`, `GetHKFuturesSnapshot`,
  `GetHKFuturesBars`, `GetHKFuturesDepth`, `GetHKFuturesFootprint` — 5 new HK
  futures market data functions.
- `data/options.go`: `OptionCategoryHK`, `OptionCategoryCN` constants.
- `data/options.go`: `GetHKOptionExpirations`, `GetHKOptionChain` — HK options
  discovery stubs (TODO t10, paths unconfirmed).
- `trade/derivatives_hk.go`: new file for HK derivatives-specific trading helpers.
- `examples/`: all examples verified against HK sandbox (auth, marketdata,
  account, watchlist, order preview, data-fundamentals, streaming, events).
- Documentation complete: `IMPLEMENTATION_STATUS.md` with full codebase audit
  (49 TODOs, 522 tests), README roadmap sync, AGENTS.md constraints updated.

### Fixed

- `data/fundamentals.go`: `FinancialsItem` changed from `map[string]string` to
  `map[string]any` to handle numeric values returned by the income statement,
  balance sheet, and cash flow endpoints.

### Tests Added

- `brokerfd/brokerfd_test.go`: unit tests for broker FD root package.
- `broker/options_test.go`: unit tests for broker options.
- `brokerfd/events/option_test.go`: unit tests for broker FD option events.

## [0.7.0] - 2026-09-21

Complete API coverage: market data extensions, event contracts, Broker API HK, Broker FD API US, and Broker FD gRPC events.

### Added

- `data/eventcontracts.go`: Event contract instrument discovery — `GetEventContractCategories`, `GetEventContractSeries`, `GetEventContractEvents`, `GetEventContractMarkets`
- `data/eventcontracts_market.go`: Event contract market data (provisional) — `GetEventSnapshot`, `GetEventDepth`, `GetEventBars`, `GetEventTick` (TODO)
- `trade/types.go`: `InstrumentTypeEvent` and `EventOutcome` (`yes`/`no`) for event contract trading
- `trade/rules.go`: `validateEventRules` — LIMIT-only, DAY-only, QTY-only, integer quantity ≤50,000 (TODO)
- `data/screener.go`: Non-display screener extensions — `GetMarketSectors`, `GetMarketSectorDetail`, `GetHighDividendRank`, `GetWeek52HighLow`
- `data/crypto_data.go`: US crypto dedicated paths — `GetCryptoSnapshotList`, `GetCryptoBarsList`
- `trade/orders.go`: `BatchPlaceOrder` for multi-order submission
- `trade/accounts.go`: Cash activity accessors — `GetCashActivities`, `GetCashActivitiesPage`, `GetAllCashActivities`
- `data/futures_market.go`: Futures market data (provisional) — `GetFuturesTick`, `GetFuturesSnapshot`, `GetFuturesBars`, `GetFuturesDepth`, `GetFuturesFootprint` (TODO)
- `data/display_screener.go`: Display Solution screeners (provisional) — `GetDisplayGainersLosers`, `GetDisplayTopActive` (TODO)
- `data/display_quotes.go`: Display Solution quotes (provisional) — `GetDisplaySnapshot`, `GetDisplayBars`, `GetDisplayBarsSingle`, `GetDisplayTick`, `GetDisplayDepth` (TODO)
- `data/display_instruments.go`, `data/display_news.go`, `data/display_streaming.go`: Display Solution instruments, news, and streaming (provisional) (TODO)
- `broker/` module: New Go module for Broker API HK — virtual accounts, instruments, assets, orders, cash activities, funding FX, instant funding, journals, master data, event contracts
- `brokerfd/` module: Complete rewrite of Broker FD API US — agreements, accounts, documents, assets, activity, funding, instruments, orders, journals, master data
- `brokerfd/events/` module: Broker FD gRPC events client using `grpc.event.EventService`
- `internal/region/region.go`: `BrokerHTTP` field added to `Endpoints` struct
- `gen/webull/brokerfd/events/v1/`: Generated protobuf for Broker FD events (P5.2)

### Changed

- `data/instrument_v3.go`: `GetStockProfilesV3` auth routed through `DisplayService()`
- `data/logos.go`: `GetLogos` auth routed through `DisplayService()`

## [0.6.0] - 2026-09-20

Multi-leg options orders, futures order validation, speculative option chain discovery, and documentation sync.

### Added

- `trade.OptionStrategy` constants for multi-leg options: `VERTICAL`, `STRADDLE`, `STRANGLE`, `IRON_CONDOR`, `IRON_BUTTERFLY`, `BUTTERFLY`, `COLLAR`, `CALENDAR`, `DIAGONAL`, `RATIO` — provisional (TODO(t8))
- `trade.validateFuturesRules`: dedicated futures order-type matrix (US/HK), QTY-only/whole-contract/DAY-GTC validation — provisional (TODO(t9))
- `trade/multileg_test.go`, `trade/futures_order_test.go`: offline tests for multi-leg and futures validation
- `data.GetOptionExpirations`, `data.GetOptionChain`: speculative option chain discovery (not in published Webull API) — marked TODO(t10)
- `data/options_test.go`: offline tests for option chain functions

### Changed

- `trade.orders.go`: instrument-dispatched validation (equity/option/futures branches)
- `trade.options.go`: multi-leg order validation (2+ leg enforcement, duplicate/degenerate detection, canonicalStrike normalization)
- `trade.rules.go`: futures order-type matrix and `isPositiveInteger` helper
- `trade.option.go`: `WithMaxOrderNotional` documentation noting multi-leg bypass
- `README.md`, `AGENTS.md`, `docs/trading.md`, `docs/api.md`, `docs/market-data.md` updated to reflect new features

## [0.5.0] - 2026-09-19

Display Solution integration, Corporate Actions, Fund Data, Crypto Data, and Screener v2 stubs.

### Added

- `display` package: Display Solution authentication service — HMAC-SHA1 signed
  client-token fetch, `Authorization: Bearer` header, lazy initialization, HK
  sandbox/prod hosts.
- `data.GetCorporateActionsList`, `data.GetCorporateActionsMarket`: Corporate
  actions by symbol or by market (Display Solution). Paths confirmed via live probe.
- `data.GetStockProfilesList`: Batch instrument profiles (POST, Display Solution).
- `data.GetLogosBatch`: Batch company logos (POST, Display Solution).
- `data.GetFundNav`, `data.GetFundInfo`, `data.GetFundDividends`, `data.GetFundList`:
  Fund NAV, info, dividends, and fund list. **Paths are best-effort; require US
  sandbox credentials to verify.**
- `data.GetCryptoBars`, `data.GetCryptoTick`, `data.GetCryptoDepth`,
  `data.GetCryptoSnapshot`: Crypto OHLCV, tick, depth, and snapshot. Paths corrected
  to `/market-data/` with `category=CRYPTO` query param. **HK sandbox returns 417
  (CRYPTO category unsupported); US credentials required to verify response schemas.**
- `data.GetScreenerV2`: Screener v2 POST query. **Path is best-effort; requires
  US sandbox credentials to verify.**
- `brokerfd` package: Broker FD HTTP stub with `GetAccountsSummary` and `GetPositions`.
  **All paths are best-effort; returns 404 in HK sandbox.**
- `examples/probe`: Live probe binary for systematic endpoint testing.

### Changed

- `data.GetCorporateActions`: signature changed from `(ctx, query) (data, error)`
  to `(ctx, query) (data, paginationKey, error)` to match the paginated API response.

### Fixed

- `data/crypto_data.go`: paths changed from `/market-data/crypto/{symbol}/bars` to
  `/market-data/bars` with `category=CRYPTO&symbol=` query params.

## [0.4.0] - 2026-09-18

Market data fundamentals: capital flows, industry comparisons, earnings and dividend
calendars, SEC filings, and financial statements.

### Added

- `data.GetCapitalFlow`: capital-flow data (large/medium/small in/out flows) for
  the trailing N trading days per symbol.
- `data.GetIndustryComparison`: industry-relative performance metrics (e.g. EPS_TTM)
  with per-company ranks and values.
- `data.GetEarningsCalendar`: upcoming and historical earnings-release dates, EPS
  actual vs. estimate, and revenue actual vs. estimate per fiscal period.
- `data.GetDividendCalendar`: dividend and split events with declare/ex-div/record/pay
  dates and per-share amounts.
- `data.GetFilings`: SEC filings list (8-K, 10-K, 10-Q, etc.) with titles, URLs,
  and publish dates.
- `data.GetIncomeStatement`: multi-period income statement data (revenue, net income,
  EPS, etc.).
- `data.GetBalanceSheet`: multi-period balance sheet data (total assets, liabilities,
  equity, etc.).
- `data.GetCashFlow`: multi-period cash flow statement data (operating, investing,
  financing cash flows).
- `data.GetFinancialIndicators`: key financial ratios and metrics (ROA, ROE, EPS,
  net margin, debt ratio).
- `data.GetFinancialAlert`: upcoming earnings-release alert with expected date and
  estimated vs. last-year EPS.
- `data.GetForecastEPS`: analyst consensus EPS forecasts for the next 5 quarters.
- `internal/auth`: added `DigestCase` (`DigestUpper`/`DigestLower`) to `SignParams`
  so callers can control the casing of the body digest hex output; the REST path
  remains unchanged.
- `examples/data-fundamentals`: Runnable program demonstrating all fundamentals
  endpoints for AAPL.

## [0.3.0] - 2026-09-18

Trading events over gRPC.

### Added

- `events` package: a server-streaming gRPC client for Webull's trade events,
  built from the core `client.Client` with `events.New`. It subscribes to order,
  event-contract position, and option streams (`SubscribeOrder`,
  `SubscribePosition`, `SubscribeOption`, and the `SubscribeAll` bitmask),
  optionally scoped to trading accounts with `events.WithAccounts`.
- Typed event payloads: `OrderEvent`, `PositionEvent`, and `OptionEvent` decode
  the JSON carried by data events, with every numeric value kept as a string and
  `Raw` preserving the exact wire bytes. Handlers are `OnOrder`, `OnPosition`,
  `OnOption`, and the raw `OnEvent`, alongside `OnConnect`, `OnPing`, and
  `OnError`.
- HMAC-SHA256 signing for the event service: no `host` participates in the
  canonical string and the body digest is lower-case SHA-256, distinct from the
  REST HMAC-SHA1 signer. The algorithm-parameterised signer in `internal/auth`
  keeps the REST path byte-for-byte unchanged (golden vectors) and adds the
  SHA-256 path used by `events`.
- Vendored `events.proto` from the Apache-2.0 Webull Python SDK, generated into
  `gen/webull/trade/events/v1`, with provenance and attribution recorded in
  `THIRD_PARTY_NOTICES.md` and the accepted ADR-0002.
- Reconnect and re-subscribe with exponential backoff plus jitter, bounded by
  `events.WithMaxReconnectAttempts`; terminal server events (auth, connection
  limit, expired subscription) end the run with a typed error.
- Runnable `examples/events` program that subscribes to order events, prints
  each decoded event, and shuts down on Ctrl+C.
- Trading events documentation page covering the endpoint, the HMAC-SHA256
  signing rules, the subscribe bitmask, the dispatch model, the order-event JSON
  schema, reconnect options, and sandbox caveats.
- A nightly, read-only live-sandbox CI job that connects to the event stream
  without mutating any account.

### Changed

- The README feature matrix and roadmap mark Trading events as supported in
  v0.3.0, and the documentation site gains a Trading Events page.

## [0.2.6] - 2026-09-18

Route the news Server-Sent Events stream through the core client pipeline.
No public API changes.

### Added

- `client.DoStream`, a signed streaming request that returns the raw, open
  response for endpoints that stream their result. It applies the same signing,
  per-path API version, token, rate limiter, and circuit breaker as
  `client.Do`, but never retries and leaves the successful response body for
  the caller to read and close.

### Changed

- The `data` package news summary (`GetNewsSummary`) now uses
  `client.DoStream` instead of its own private stream path, so the access token,
  API version, and resilience policies apply to it exactly as to buffered
  requests. Its removed private helpers are not part of the public API.

## [0.2.5] - 2026-09-18

US combo orders for the `trade` package.

### Added

- US combo-order groups: `PlaceOrderRequest` now validates the whole order set as
  a combo group when any order uses a non-NORMAL `combo_type`. Supported groups
  are take-profit/stop-loss (`MASTER` with optional `STOP_PROFIT` and
  `STOP_LOSS`), `OTO`, `OCO`, and `OTOCO`. Combo orders are US equity only.
- Composition and leg-count rules: one request must contain a single group kind,
  every order must be a US equity order, and `client_combo_order_id` is required.
  OTO and OTOCO require exactly one `MASTER` plus one to six legs, OCO requires
  two to six legs with no `MASTER`, and take-profit/stop-loss allows at most one
  `MASTER`, one `STOP_PROFIT`, and one `STOP_LOSS`.
- Per-role order-type rules: a take-profit/stop-loss `MASTER` is `MARKET` or
  `LIMIT`; an OTO/OTOCO `MASTER` and an OTO leg add `STOP_LOSS` and
  `STOP_LOSS_LIMIT`; take-profit/stop-loss, OCO, and OTOCO legs are `LIMIT`,
  `STOP_LOSS`, or `STOP_LOSS_LIMIT`.
- Sell-to-close take-profit/stop-loss groups: a group with no `MASTER` must use
  side `SELL` for every sub-order.
- Table-driven tests for the combo composition, leg-count, order-type, and
  sell-to-close rules, plus a combo preview test.
- A "Combo types" section in the trading documentation covering the group kinds,
  the composition and leg-count rules, the sell-to-close form, and a
  MASTER/STOP_PROFIT/STOP_LOSS example.

### Changed

- `PlaceOrderRequest.Validate` now runs the combo-group validation after the
  per-order checks.

## [0.2.4] - 2026-09-18

Single-leg options orders for the `trade` package.

### Added

- Single-leg options orders: `OrderRequest` accepts `instrument_type` `OPTION`
  with `option_strategy` `SINGLE` and exactly one `legs` entry, validated before
  any network call. Option orders support `LIMIT`, `STOP_LOSS`, and
  `STOP_LOSS_LIMIT`; `MARKET` and the other stock order types are rejected.
- Options side and time-in-force rules: only `BUY` and `SELL` are accepted
  (`SHORT` is rejected), sell-side orders must use `DAY`, and `GTD` is rejected
  for options.
- `OrderLeg.Validate` validates a single option leg: `instrument_type` `OPTION`,
  `market` `US`, a non-blank `symbol`, `BUY`/`SELL`, a positive decimal
  `strike_price`, a calendar `option_expire_date` in `YYYY-MM-DD` form,
  `option_type` `CALL` or `PUT`, and a positive decimal `quantity`.
- Table-driven tests for the option order and leg rules, plus an env-gated,
  live-tolerant preview test.
- An "Options orders" section in the trading documentation covering the
  supported order types, side and time-in-force rules, the `legs` schema, a code
  example, and sandbox option-contract caveats.

### Changed

- `OrderRequest.Validate` now enforces the option-specific rules after the
  market rules; a non-OPTION order that carries `option_strategy` or `legs` is
  rejected.

## [0.2.3] - 2026-09-18

Market-specific order validation and Hong Kong BCAN support for the `trade`
package.

### Added

- Market-specific equity order-type validation: US accepts limit, market, stop,
  stop-limit, market-on-open, market-on-close, touch, and trailing stop orders;
  HK accepts enhanced limit, at-auction, at-auction limit, stop, stop-limit,
  touch, and trailing stop orders; CN accepts only `LIMIT`.
- Hong Kong BCAN: `HK` equity orders now require at least one `no_party_ids`
  entry with a non-blank `party_id`, `party_id_source` `"D"`, and `party_role`
  `"3"`. A `no_party_ids` list on a non-HK-equity order is rejected.
- US `support_trading_session` validation (`CORE`, `ALL`, `NIGHT`, `ALL_DAY`).
  The deprecated `Y` and `N` aliases and use on a non-US order are rejected.
- At-auction price rules: `AT_AUCTION` rejects `limit_price` and
  `AT_AUCTION_LIMIT` requires it.
- A-share (`CN`) documentation noting that only `LIMIT` is accepted and that
  A-share trading is disabled by default until enabled by Webull support.
- Table-driven tests covering the market matrix, Hong Kong BCAN, US trading
  sessions, and the at-auction price rules.

### Changed

- The trading documentation adds a "Market rules" section, and the `OrderType`,
  `TradingSession`, and `PartyID` GoDoc describe the enforced rules.

## [0.2.2] - 2026-09-18

Stock-order lifecycle and queries for the `trade` package.

### Added

- `trade` order methods: `PreviewOrder` and `PlaceOrder` (with guardrail
  enforcement), `ReplaceOrder`, and `CancelOrder`. Requests are validated before
  any network call and the v3 modify endpoints identify an order by its client
  order ID.
- `trade` order queries: `GetOpenOrders`, `GetOpenOrdersPage`,
  `GetAllOpenOrders`, `GetOrderHistory`, `GetOrderHistoryPage`,
  `GetAllOrderHistory`, and `GetOrderDetail`. Pagination follows the cursor to
  exhaustion, bounded by `trade.MaxOrderQueryPages`.
- Order domain types and enums: `OrderRequest`, `PlaceOrderRequest`,
  `PlaceOrderResult`, `PreviewResult`, `ModifyOrderRequest`,
  `ReplaceOrderRequest`, `ReplaceOrderResult`, `CancelOrderRequest`,
  `CancelOrderResult`, `OrderGroup`, `OrderPage`, `Order`, `OrderLeg`,
  `OrderLegDetail`, `OrderCommission`, `OrderFee`, `OrderHistoryQuery`,
  `OrderSide`, `OrderType`, `TimeInForce`, `ComboType`, `EntrustType`,
  `TradingSession`, `TriggerPriceType`, `TrailingType`, `OrderStatus`, and
  `PartyID`.
- `OrderRequest.Validate`, `PlaceOrderRequest.Validate`,
  `ReplaceOrderRequest.Validate`, and `CancelOrderRequest.Validate`, returning
  typed `invalid_config` errors for the first problem found.
- Runnable `examples/order` program that previews a non-marketable AAPL limit
  buy and, only when `WEBULL_ORDER_PLACE=1` is set, places and cancels it.
- Expanded trading documentation covering the order lifecycle, order-type,
  time-in-force, and combo-type tables, validation rules, guardrails, queries,
  and the explicit warning that placing orders mutates a real account.

## [0.2.1] - 2026-09-18

Trading HTTP foundation: read-only accounts and assets.

### Added

- `trade` package: a typed Trading HTTP client built on the core `client.Client`
  with `trade.New`. Endpoints: `ListAccounts`, `GetBalance`, and `GetPositions`,
  covering account listing and per-account balances and positions. Asset calls
  reject an empty account ID before any network request.
- Trading domain types and enums: `Account`, `AccountType`, `AccountClass`,
  `AssetsBalance`, `AssetsCurrencyAssets`, `Position`, `PositionLeg`, `Market`,
  `InstrumentType`, `OptionType`, and `OptionStrategy`. Numeric fields stay
  strings to preserve precision.
- Order guardrail options `trade.WithMaxOrderNotional` and
  `trade.WithMaxOrderQuantity`, enforced by the order methods before an order is
  built (order placement lands in a later patch).
- Built-in per-path API version defaults: requests under `/trading/` now send
  `x-version: v3`, while Market Data paths remain on `v2`. Overridable with
  `client.WithAPIVersion` and `client.WithAPIVersionFor`.
- Runnable `examples/account` program that lists accounts and prints the balance
  and positions for an account.
- Trading documentation page covering authentication, accounts and assets, the
  v3 default, guardrails, and sandbox usage.

## [0.1.0] - 2026-09-18

Initial public release.

### Added

- Core `client` package: `client.New`, the single signed-request entry point
  `Client.Do`, and functional options for credentials, region, environment,
  endpoints, transport, retries, rate limiting, circuit breaking, API version,
  and environment configuration (`WithEnv`).
- Authentication: HMAC-SHA1 request signing over a percent-encoded canonical
  string, the token lifecycle (`CreateToken`, `CheckToken`, `EnsureToken`,
  `CurrentToken`, `AccessToken`, `SetToken`), automatic token injection, and
  automatic sandbox token acquisition with `WithAutoToken`.
- `data` package: Market Data HTTP endpoints for instruments, company profile,
  analyst data, futures static data, snapshot, tick, quotes/depth, single and
  batch bars, footprint, NOII, screener, watchlist CRUD, options, and news.
- `stream` package: Market Data streaming over MQTT and MQTT-over-WebSocket with
  typed `Quote`, `Snapshot`, and `Tick` handlers, automatic reconnect, and
  idempotent automatic re-subscription.
- Generated protobuf types in `gen/webull/marketdata/v1` for streamed messages.
- `pkg/types` shared domain types.
- Resilient transport primitives: retry with exponential backoff, a keyed token
  bucket rate limiter, and a circuit breaker.
- Documentation site (MkDocs Material) with getting-started, authentication,
  market-data, streaming, sandbox, errors, API reference, and troubleshooting
  pages, plus Architecture Decision Records.
- Runnable examples under `examples/` for auth, market data, streaming, and
  watchlists.

[Unreleased]: https://github.com/shing1211/webullapi4go/compare/v2.1.38...HEAD
[2.1.1]: https://github.com/shing1211/webullapi4go/releases/tag/v2.1.1
[2.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v2.1.0
[2.0.9]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.9
[2.0.8]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.8
[2.0.7]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.7
[2.0.6]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.6
[2.0.5]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.5
[2.0.4]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.4
[2.0.3]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.3
[2.0.2]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.2
[2.0.1]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.1
[2.0.0]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.0
[1.1.1]: https://github.com/shing1211/webullapi4go/releases/tag/v1.1.1
[1.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v1.1.0
[1.0.3]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.3
[1.0.2]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.2
[1.0.1]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.1
[1.0.0]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.0
[0.9.2]: https://github.com/shing1211/webullapi4go/releases/tag/v0.9.2
[0.9.1]: https://github.com/shing1211/webullapi4go/releases/tag/v0.9.1
[0.9.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.9.0
[0.7.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.7.0
[0.6.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.6.0
[0.5.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.5.0
[0.4.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.4.0
[0.3.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.3.0
[0.2.6]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.6
[0.2.5]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.5
[0.2.4]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.4
[0.2.3]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.3
[0.2.2]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.2
[0.2.1]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.1
[0.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.1.0

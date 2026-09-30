# fix-futures-tick-instrument-id

Accept both documented and live instrument identifier spellings for data.GetFuturesTick without breaking the public StockTicks API.

## Status: Completed

Implementation merged to `main` in `v2.1.37`. All 16 tasks closed.
The `data.StockTicks.UnmarshalJSON` decoder accepts both `instrument_id` (documented)
and `instrumentId` (observed from sandbox) spellings. The field tag, marshalling
behaviour, and method signatures are unchanged. The live-divergences row is
retained as a record about Webull's server, not about the SDK.

# Next Phase

## Completed This Run

Fixed all Broker API HK paths to match official Webull API specification. Discovered that the entire `broker/` package had wrong paths — SDK used `/openapi/v1/broker/...` but official API uses `/broker/...` directly. Release v0.9.2.

## Gaps and Tech Debt

- Broker API HK returns `401 ROUTE_NOT_PERMITTED` in sandbox — app lacks scope
- 49 provisional TODO items still need US sandbox or Display Solution entitlement
- `golangci-lint` CI already configured (no action needed)
- No GitHub Actions CI pipeline (golangci-lint already in workflow)

## Candidate Next-Phase Items

| # | Title | Objective | Why Now | Effort | Dependencies |
|---|-------|-----------|---------|--------|-------------|
| 1 | Obtain Broker API HK scope | Enable Broker API in Webull developer portal | Unblocks all 32 broker methods | S | User enables in portal |
| 2 | US sandbox registration | Get Webull US developer portal access | Unblocks 36+ TODO items | M | User registers |
| 3 | v1.0 API stability review | Freeze public API surface | Required before semver guarantees | L | All verification complete |
| 4 | Futures/event-contract live probe | Test futures bars, event-contract market data | Verify HK paths for these endpoints | M | Specific symbols |

## Recommended Next Phase

**v1.0 API stability review** — largest remaining effort, no credentials needed. Freezes the public API surface and confirms which provisional TODO items should be removed vs. kept.

## Open Questions

1. Do you have Broker API HK scope enabled in your developer portal?
2. Do you plan to register for US sandbox access?

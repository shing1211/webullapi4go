# Report — v02-trading-http (v0.2.x)

## Outcome
The Trading HTTP phase is complete, delivered as six patch releases `0.2.1` … `0.2.6`,
all published to GitHub and Gitee.

| Patch | Theme | Release commit | Tag |
|---|---|---|---|
| 0.2.1 | `trade` package + accounts/assets + `/trading` v3 + guardrails | `5958dcf` | `v0.2.1` |
| 0.2.2 | Stock order lifecycle (preview/place/replace/cancel/query) | `025588b` | `v0.2.2` |
| 0.2.3 | Market-specific rules + HK BCAN | `d9ba14a` | `v0.2.3` |
| 0.2.4 | Single-leg options orders | `b08b83c` | `v0.2.4` |
| 0.2.5 | US combo orders (TP/SL, OTO, OCO, OTOCO) | `fb76828` | `v0.2.5` |
| 0.2.6 | News SSE refactor onto `client.DoStream` | `df6040d` | `v0.2.6` |

## Shipped
- **`trade` package**: `New` + typed enums/structs; `ListAccounts`, `GetBalance`,
  `GetPositions`; `PreviewOrder`, `PlaceOrder`, `ReplaceOrder`, `CancelOrder`;
  `GetOpenOrders`/`GetAllOpenOrders`, `GetOrderHistory`/`GetAllOrderHistory`,
  `GetOrderDetail`; cursor pagination with a max-page guard.
- **Validation**: required fields, `client_order_id` format, conditional prices;
  US/HK/CN order-type matrices; HK BCAN `no_party_ids`; US trading-session rules;
  at-auction price rules; option `SINGLE` rules; combo composition/leg-count rules.
- **Guardrails**: `WithMaxOrderNotional`, `WithMaxOrderQuantity` (opt-in).
- **Client**: `/trading/**` defaults to `x-version: v3`; public `client.DoStream`.
- **Docs/examples**: `docs/trading.md` (accounts, orders, market rules, options, combos),
  `examples/account`, `examples/order`; README/CHANGELOG/docs site updated.

## Verification
- `go build ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run ./...` clean at each
  patch; `go test ./... -count=1` (offline) and `-race` green.
- Live sandbox (dedicated Account #2 `OGG4…`): accounts/balance/positions, preview,
  **place → cancel** (mutate-gated), open/history queries, HK enhanced-limit preview,
  option preview, combo preview all pass.
- `mkdocs build --strict` succeeds; `docs/runs/**` excluded from the site.

## Deviations / notes
1. Order paths differ from early guesses and were confirmed from the reference docs:
   history is `/trading/orders/historical-orders/list`; detail is `/trading/orders/get`
   (keyed by `client_order_id`, not `order_id`).
2. Guardrail option constructors panic on non-numeric input (matches the existing
   `client.NewRateLimiter` convention); invalid numeric values cannot be represented.
3. Sandbox cancel of a just-placed non-marketable order can return `417 The current status
   cannot be modified`; the test tolerates this and does not guarantee cleanup.
4. `data` news live call returns an upstream `504 GATEWAY_TIMEOUT` from the sandbox
   summarizer; the typed error proves auth/version were accepted. Test tolerates it.
5. `TestSandboxWatchlists` intermittently returns `401` on the shared account — reproduced
   on baseline `fb76828` (pre-existing shared-sandbox flake), outside this run's scope.
6. A released commit message for `v0.2.3` says "74 subtests"; the accurate count is 62.
   Immutable after push; noted here.

## Risks / follow-ups
- gRPC trade events remain unimplemented (ADR-0002 proposes extracting `.proto` from the
  Apache-2.0 Python SDK).
- Display Solution and Broker API not implemented.
- Sandbox test reliability on shared accounts (token churn / 401) — consider consolidating
  to a single cached token per package.

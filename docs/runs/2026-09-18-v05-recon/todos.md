# Todos — webullapi4go Phase 0: API Recon Spike

Single source of truth. Statuses: `todo`, `doing`, `blocked`, `review`, `done`.

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|-----------|------------|
| T0.1 | Probe Display Solution host + auth flow | backend | done | — | Hosts confirmed from official docs; live probing blocked — no credentials. findings.md updated. |
| T0.2 | Corporate Actions schemas | backend | todo | T0.1 | Paths + response schema documented in findings.md |
| T0.3 | Fund data endpoints | backend | todo | T0.1 | All fund data paths + schemas in findings.md |
| T0.4 | Crypto data endpoints | backend | todo | T0.1 | All crypto paths (bars, tick, depth, orderbook) + schemas in findings.md |
| T0.5 | Screener v2 endpoint | backend | todo | T0.1 | Screener v2 path + schema in findings.md |
| T0.6 | Broker FD HTTP schemas (12 categories) | backend | todo | T0.1 | All 12 category schemas in findings.md |
| T0.7 | Broker FD gRPC proto + event types | architect | todo | T0.1 | Proto files in gen/webull/brokerfd/v1/, event types in findings.md |

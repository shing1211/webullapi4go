# Broker FD Example

Demonstrates read-only Broker FD (US) API operations.

## Credentials

Set these environment variables before running:

```sh
export WEBULL_APP_KEY=your_app_key
export WEBULL_APP_SECRET=your_app_secret
export WEBULL_TRADE_ACCOUNT_ID=your_account_id   # optional; needed for open orders
export WEBULL_ENVIRONMENT=sandbox                # or omit for production
```

## Run

```sh
go run ./examples/brokerfd
```

The program will exit early with a graceful message if `WEBULL_APP_KEY` is not set.

## Operations

- **Accounts Summary**: Cash and buying power for all Broker FD accounts
- **Positions**: All positions held across Broker FD accounts
- **Open Orders**: Open (unfilled) fractional orders for a specific account (requires `WEBULL_TRADE_ACCOUNT_ID`)

This is a read-only example. It never places, replaces, or cancels orders.

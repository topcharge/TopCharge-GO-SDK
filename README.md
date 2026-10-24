# TopCharge Go SDK

Dependency-free Go SDK for the TopCharge merchant API. It uses `net/http`, works with standard HTTP services, Gin, Echo, Fiber, Chi, gRPC workers, and command-line applications.

## Release status

The SDK source is publicly available at [github.com/topcharge/TopCharge-GO-SDK](https://github.com/topcharge/TopCharge-GO-SDK).

## Install

```bash
go get github.com/topcharge/TopCharge-GO-SDK
```

## Install locally

Add this local module to your application during development:

```bash
go mod edit -require=github.com/topcharge/TopCharge-GO-SDK@v0.0.0
go mod edit -replace=github.com/topcharge/TopCharge-GO-SDK=/absolute/path/to/sdk/go
go mod tidy
```

## Quick start

```go
package main

import (
    "context"
    "os"

    topcharge "github.com/topcharge/TopCharge-GO-SDK"
)

func main() {
    client, err := topcharge.NewClient(
        os.Getenv("TOPCHARGE_CLIENT_ID"),
        os.Getenv("TOPCHARGE_SECRET_KEY"),
    )
    if err != nil { panic(err) }

    payment, err := client.Payments.Create(context.Background(), topcharge.PaymentCreateRequest{
        PaymentAmount: "25000",
        Description: "Invoice #1001",
        RedirectURL: "https://merchant.example/success",
        CallbackURL: "https://merchant.example/webhooks/payment",
    })
    if err != nil { panic(err) }

    println(payment.PaymentURL)
}
```

The default API host is `https://api.topcharge.app`. Use `topcharge.WithBaseURL(...)` or `topcharge.WithHTTPClient(...)` when needed. API failures return `*topcharge.APIError`, exposing `StatusCode` and `ResponseBody`.

The client sends `x-api-key` and `x-api-secret`. Store their values in server-side environment variables or secret storage only; never commit or expose them to a browser.

## Services

```go
ctx := context.Background()
if err := client.Account.VerifyCredentials(ctx); err != nil { panic(err) }
if _, err := client.Account.GetBalance(ctx); err != nil { panic(err) }

payment, err := client.Payments.Create(ctx, topcharge.PaymentCreateRequest{PaymentAmount: "25000"})
if err != nil { panic(err) }
if _, err := client.Payments.GetStatus(ctx, payment.PaymentID); err != nil { panic(err) }
if _, err := client.Payments.Cancel(ctx, payment.PaymentID); err != nil { panic(err) }

if _, err := client.Withdrawals.Create(ctx, topcharge.WithdrawalCreateRequest{
    PhoneNumberCountryCode: "964", PhoneNumber: "7701234567", Amount: "50000",
}); err != nil { panic(err) }

```

## Tests and individual integration examples

```bash
go test ./...
go run ./cmd/integration help
```

For a live payment-link smoke test (create → status → cancel):

```bash
export TOPCHARGE_RUN_INTEGRATION=1
export TOPCHARGE_CLIENT_ID='mch_...'
export TOPCHARGE_SECRET_KEY='sk_...'
go run ./cmd/integration payment:smoke
```

Each operation is separate: `payment:create`, `payment:status`, `payment:cancel`, `account:balance`, and `withdrawal:create`. Operations that move value require their matching `TOPCHARGE_ALLOW_*` confirmation variable. Run `go run ./cmd/integration help` for the complete list.

Keep merchant credentials in server environment variables and never commit them to Git.

## Reliability

The default HTTP client timeout is 30 seconds, and every request accepts a `context.Context` for cancellation. The SDK does not retry automatically. Reads can generally be retried; do not blindly retry payment creation or withdrawals after an unknown result.

## License

[MIT License](LICENSE)

## Source

[github.com/topcharge/TopCharge-GO-SDK](https://github.com/topcharge/TopCharge-GO-SDK)

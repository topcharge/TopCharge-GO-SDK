# TopCharge Go SDK

Dependency-free Go SDK for the TopCharge merchant API. It uses `net/http`, works with standard HTTP services, Gin, Echo, Fiber, Chi, gRPC workers, and command-line applications.

## Release status

This module is currently consumed from this repository for local development. `github.com/topcharge/topcharge-go` is its intended module path, not a claim that the module is publicly available yet.

## Install locally

Add this local module to your application during development:

```bash
go mod edit -require=github.com/topcharge/topcharge-go@v0.0.0
go mod edit -replace=github.com/topcharge/topcharge-go=/absolute/path/to/sdk/go
go mod tidy
```

## Quick start

```go
package main

import (
    "context"
    "os"

    topcharge "github.com/topcharge/topcharge-go"
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
batchID := "01..." // An ID returned by Issue or ListBatches.
cardID := "01..."  // An ID returned by a card-list operation.

if err := client.Account.VerifyCredentials(ctx); err != nil { panic(err) }
if _, err := client.Account.GetBalance(ctx); err != nil { panic(err) }

payment, err := client.Payments.Create(ctx, topcharge.PaymentCreateRequest{PaymentAmount: "25000"})
if err != nil { panic(err) }
if _, err := client.Payments.GetStatus(ctx, payment.PaymentID); err != nil { panic(err) }
if _, err := client.Payments.Cancel(ctx, payment.PaymentID); err != nil { panic(err) }

if _, err := client.Withdrawals.Create(ctx, topcharge.WithdrawalCreateRequest{
    PhoneNumberCountryCode: "964", PhoneNumber: "7701234567", Amount: "50000",
}); err != nil { panic(err) }

if _, err := client.GiftCards.Issue(ctx, topcharge.GiftCardIssueRequest{Amount: "25000", IdempotencyKey: "order-1001-card-01"}); err != nil { panic(err) }
if _, err := client.GiftCards.ListBatches(ctx, topcharge.PaginationOptions{Limit: 20}); err != nil { panic(err) }
if _, err := client.GiftCards.GetBatch(ctx, batchID); err != nil { panic(err) }
if _, err := client.GiftCards.ListBatchCards(ctx, batchID, topcharge.PaginationOptions{Limit: 20}, true); err != nil { panic(err) }
if _, err := client.GiftCards.RevokeBatch(ctx, batchID, topcharge.GiftCardRevokeRequest{Reason: "Campaign cancelled"}); err != nil { panic(err) }
if _, err := client.GiftCards.ListCards(ctx, topcharge.GiftCardListOptions{PaginationOptions: topcharge.PaginationOptions{Limit: 20}, Status: "ACTIVE"}); err != nil { panic(err) }
if _, err := client.GiftCards.GetCard(ctx, cardID, true); err != nil { panic(err) }
if _, err := client.GiftCards.RevokeCard(ctx, cardID, topcharge.GiftCardRevokeRequest{Reason: "Issued in error"}); err != nil { panic(err) }
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

Each operation is separate: `payment:create`, `payment:status`, `payment:cancel`, `account:balance`, `withdrawal:create`, and all `giftcard:*` commands. Operations that move or reserve value require their matching `TOPCHARGE_ALLOW_*` confirmation variable. Run `go run ./cmd/integration help` for the complete list.

Keep merchant credentials in server environment variables and never commit them to Git.

## Reliability

The default HTTP client timeout is 30 seconds, and every request accepts a `context.Context` for cancellation. The SDK does not retry automatically. Reads can generally be retried; do not blindly retry payment creation or withdrawals after an unknown result. Retry gift-card issuance only with the same persisted idempotency key.

## License

[MIT License](LICENSE)

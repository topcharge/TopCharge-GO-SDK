package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	topcharge "github.com/topcharge/topcharge-go"
)

func main() {
	action := "help"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	if action == "help" {
		printHelp()
		return
	}
	if os.Getenv("TOPCHARGE_RUN_INTEGRATION") != "1" {
		fail("Refusing to call the API. Set TOPCHARGE_RUN_INTEGRATION=1.")
	}

	options := []topcharge.Option{}
	if baseURL := os.Getenv("TOPCHARGE_BASE_URL"); baseURL != "" {
		options = append(options, topcharge.WithBaseURL(baseURL))
	}
	client, err := topcharge.NewClient(required("TOPCHARGE_CLIENT_ID"), required("TOPCHARGE_SECRET_KEY"), options...)
	if err != nil {
		fail(err.Error())
	}
	ctx := context.Background()

	switch action {
	case "account:verify":
		err = client.Account.VerifyCredentials(ctx)
		if err == nil {
			fmt.Println("Credentials are valid.")
		}
	case "account:balance":
		var result *topcharge.AccountBalance
		result, err = client.Account.GetBalance(ctx)
		output(result)
	case "payment:create":
		var result *topcharge.Payment
		result, err = client.Payments.Create(ctx, paymentRequest("manual test"))
		output(result)
	case "payment:status":
		var result *topcharge.Payment
		result, err = client.Payments.GetStatus(ctx, required("TOPCHARGE_PAYMENT_ID"))
		output(result)
	case "payment:cancel":
		var result *topcharge.PaymentCancelResponse
		result, err = client.Payments.Cancel(ctx, required("TOPCHARGE_PAYMENT_ID"))
		output(result)
	case "payment:smoke":
		err = paymentSmoke(ctx, client)
	case "withdrawal:create":
		confirmation("TOPCHARGE_ALLOW_MONEY_MOVEMENT")
		var result *topcharge.Withdrawal
		result, err = client.Withdrawals.Create(ctx, topcharge.WithdrawalCreateRequest{PhoneNumberCountryCode: required("TOPCHARGE_PHONE_COUNTRY_CODE"), PhoneNumber: required("TOPCHARGE_PHONE_NUMBER"), Amount: required("TOPCHARGE_WITHDRAWAL_AMOUNT")})
		output(result)
	case "giftcard:issue":
		confirmation("TOPCHARGE_ALLOW_VALUE_ISSUANCE")
		var result *topcharge.GiftCardIssueResponse
		result, err = client.GiftCards.Issue(ctx, topcharge.GiftCardIssueRequest{Amount: required("TOPCHARGE_GIFTCARD_AMOUNT"), IdempotencyKey: required("TOPCHARGE_IDEMPOTENCY_KEY"), Quantity: positiveInt("TOPCHARGE_GIFTCARD_QUANTITY"), Name: os.Getenv("TOPCHARGE_GIFTCARD_NAME")})
		output(result)
	case "giftcard:batches":
		var result *topcharge.CursorPage[topcharge.GiftCardBatch]
		result, err = client.GiftCards.ListBatches(ctx, pagination())
		output(result)
	case "giftcard:batch":
		var result *topcharge.GiftCardBatchDetail
		result, err = client.GiftCards.GetBatch(ctx, required("TOPCHARGE_BATCH_ID"))
		output(result)
	case "giftcard:batch-cards":
		var result *topcharge.CursorPage[topcharge.GiftCard]
		result, err = client.GiftCards.ListBatchCards(ctx, required("TOPCHARGE_BATCH_ID"), pagination(), boolean("TOPCHARGE_REVEAL"))
		output(result)
	case "giftcard:revoke-batch":
		confirmation("TOPCHARGE_ALLOW_VALUE_REVOCATION")
		var result *topcharge.GiftCardRevokeResponse
		result, err = client.GiftCards.RevokeBatch(ctx, required("TOPCHARGE_BATCH_ID"), topcharge.GiftCardRevokeRequest{Reason: os.Getenv("TOPCHARGE_REVOKE_REASON")})
		output(result)
	case "giftcard:cards":
		var result *topcharge.CursorPage[topcharge.GiftCard]
		result, err = client.GiftCards.ListCards(ctx, topcharge.GiftCardListOptions{PaginationOptions: pagination(), Status: os.Getenv("TOPCHARGE_GIFTCARD_STATUS"), BatchID: os.Getenv("TOPCHARGE_BATCH_ID"), Standalone: nullableBoolean("TOPCHARGE_STANDALONE"), Reveal: boolean("TOPCHARGE_REVEAL")})
		output(result)
	case "giftcard:card":
		var result *topcharge.GiftCard
		result, err = client.GiftCards.GetCard(ctx, required("TOPCHARGE_CARD_ID"), boolean("TOPCHARGE_REVEAL"))
		output(result)
	case "giftcard:revoke-card":
		confirmation("TOPCHARGE_ALLOW_VALUE_REVOCATION")
		var result *topcharge.GiftCardRevokeResponse
		result, err = client.GiftCards.RevokeCard(ctx, required("TOPCHARGE_CARD_ID"), topcharge.GiftCardRevokeRequest{Reason: os.Getenv("TOPCHARGE_REVOKE_REASON")})
		output(result)
	default:
		fail("Unknown action '" + action + "'. Run: go run ./cmd/integration help")
	}
	if err != nil {
		fail(err.Error())
	}
}

func paymentSmoke(ctx context.Context, client *topcharge.Client) error {
	payment, err := client.Payments.Create(ctx, paymentRequest("smoke "+randomID()))
	if err != nil {
		return err
	}
	if payment.PaymentID == "" || payment.PaymentURL == "" {
		return fmt.Errorf("payment creation response is missing paymentId or paymentUrl")
	}
	defer func() {
		result, cancelErr := client.Payments.Cancel(ctx, payment.PaymentID)
		if cancelErr != nil {
			fmt.Fprintln(os.Stderr, "payment cleanup failed:", cancelErr)
			return
		}
		output(result)
		fmt.Println("Test payment cancelled.")
	}()
	status, err := client.Payments.GetStatus(ctx, payment.PaymentID)
	if err != nil {
		return err
	}
	if status.PaymentID != payment.PaymentID {
		return fmt.Errorf("payment status response does not match the created payment")
	}
	fmt.Println("Payment link created and verified for payment ID:", payment.PaymentID)
	return nil
}

func paymentRequest(label string) topcharge.PaymentCreateRequest {
	return topcharge.PaymentCreateRequest{PaymentAmount: envOr("TOPCHARGE_TEST_PAYMENT_AMOUNT", "50000"), Description: envOr("TOPCHARGE_PAYMENT_DESCRIPTION", "TopCharge Go SDK "+label), RedirectURL: envOr("TOPCHARGE_REDIRECT_URL", "https://example.com/topcharge-sdk-test"), CallbackURL: envOr("TOPCHARGE_CALLBACK_URL", "https://example.com/topcharge-sdk-webhook")}
}
func pagination() topcharge.PaginationOptions {
	return topcharge.PaginationOptions{Cursor: os.Getenv("TOPCHARGE_CURSOR"), Limit: valueOrZero("TOPCHARGE_LIMIT"), Sort: sort()}
}
func required(name string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	fail(name + " is required.")
	return ""
}
func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
func confirmation(name string) {
	if os.Getenv(name) != "1" {
		fail("This action changes money or value. Set " + name + "=1 to continue.")
	}
}
func valueOrZero(name string) int {
	if os.Getenv(name) == "" {
		return 0
	}
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil || value < 1 {
		fail(name + " must be a positive integer.")
	}
	return value
}
func positiveInt(name string) *int {
	value := valueOrZero(name)
	if value == 0 {
		return nil
	}
	return &value
}
func sort() string {
	value := os.Getenv("TOPCHARGE_SORT")
	if value != "" && value != "asc" && value != "desc" {
		fail("TOPCHARGE_SORT must be asc or desc.")
	}
	return value
}
func boolean(name string) bool {
	value := os.Getenv(name)
	if value == "" {
		return false
	}
	switch value {
	case "1", "true", "yes":
		return true
	case "0", "false", "no":
		return false
	default:
		fail(name + " must be true or false.")
		return false
	}
}
func nullableBoolean(name string) *bool {
	if os.Getenv(name) == "" {
		return nil
	}
	value := boolean(name)
	return &value
}
func output(value any) {
	encoded, _ := json.MarshalIndent(value, "", "  ")
	fmt.Println(string(encoded))
}
func randomID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
func fail(message string) { fmt.Fprintln(os.Stderr, message); os.Exit(1) }
func printHelp() {
	fmt.Println(`TopCharge Go integration examples

Set TOPCHARGE_RUN_INTEGRATION=1, TOPCHARGE_CLIENT_ID, and TOPCHARGE_SECRET_KEY first.

Account: account:verify | account:balance
Payment: payment:create | payment:status (TOPCHARGE_PAYMENT_ID) | payment:cancel (TOPCHARGE_PAYMENT_ID) | payment:smoke
Withdrawal: withdrawal:create (TOPCHARGE_ALLOW_MONEY_MOVEMENT=1 plus recipient variables)
Gift cards: giftcard:issue | giftcard:batches | giftcard:batch | giftcard:batch-cards | giftcard:cards | giftcard:card | giftcard:revoke-batch | giftcard:revoke-card

Run payment:smoke to create, verify, then cancel a test payment.`)
}

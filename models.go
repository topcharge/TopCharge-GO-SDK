package topcharge

import "time"

type AccountBalance struct {
	TotalBalance         string `json:"totalBalance"`
	ReservedBalance      string `json:"reservedBalance"`
	LiabilityBalance     string `json:"liabilityBalance"`
	AvailableBalance     string `json:"availableBalance"`
	CashAvailableBalance string `json:"cashAvailableBalance"`
}

type PaymentCreateRequest struct {
	PaymentAmount string `json:"paymentAmount"`
	Description   string `json:"description,omitempty"`
	RedirectURL   string `json:"redirectUrl,omitempty"`
	CallbackURL   string `json:"callbackUrl,omitempty"`
}

type Payment struct {
	PaymentID       string     `json:"paymentId"`
	PaymentURL      string     `json:"paymentUrl"`
	RequestedAmount string     `json:"requestedAmount"`
	FeeAmount       string     `json:"feeAmount"`
	NetAmount       string     `json:"netAmount"`
	Status          string     `json:"status"`
	ExpiresAt       time.Time  `json:"expiresAt"`
	PaidAt          *time.Time `json:"paidAt"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type PaymentCancelResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	PaymentID string `json:"paymentId"`
	Status    string `json:"status"`
}

type WithdrawalCreateRequest struct {
	PhoneNumberCountryCode string `json:"phoneNumberCountryCode"`
	PhoneNumber            string `json:"phoneNumber"`
	Amount                 string `json:"amount"`
}

type WithdrawalRecipient struct {
	ID                     string `json:"id"`
	Username               string `json:"username"`
	PhoneNumberCountryCode string `json:"phoneNumberCountryCode"`
	PhoneNumber            string `json:"phoneNumber"`
}

type Withdrawal struct {
	WithdrawalID    string              `json:"withdrawalId"`
	RequestedAmount string              `json:"requestedAmount"`
	FeePercent      string              `json:"feePercent"`
	FeeAmount       string              `json:"feeAmount"`
	NetAmount       string              `json:"netAmount"`
	TransactionID   string              `json:"transactionId"`
	Status          string              `json:"status"`
	Recipient       WithdrawalRecipient `json:"recipient"`
	CompletedAt     time.Time           `json:"completedAt"`
	CreatedAt       time.Time           `json:"createdAt"`
}

type GiftCardIssueRequest struct {
	Amount         string `json:"amount"`
	IdempotencyKey string `json:"idempotencyKey"`
	Quantity       *int   `json:"quantity,omitempty"`
	Name           string `json:"name,omitempty"`
}

type GiftCardIssueResponse struct {
	Batch    *GiftCardBatch `json:"batch"`
	Cards    []GiftCard     `json:"cards"`
	Replayed bool           `json:"replayed"`
}

type GiftCard struct {
	ID           string     `json:"id"`
	BatchID      *string    `json:"batchId"`
	Code         *string    `json:"code"`
	MaskedCode   string     `json:"maskedCode"`
	Amount       string     `json:"amount"`
	Currency     string     `json:"currency"`
	Status       string     `json:"status"`
	RedeemedAt   *time.Time `json:"redeemedAt"`
	RevokedAt    *time.Time `json:"revokedAt"`
	RevokeReason *string    `json:"revokeReason"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type GiftCardBatch struct {
	ID          string    `json:"id"`
	Name        *string   `json:"name"`
	Quantity    int       `json:"quantity"`
	UnitAmount  string    `json:"unitAmount"`
	TotalAmount string    `json:"totalAmount"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GiftCardBatchDetail struct {
	GiftCardBatch
	ActiveCount   int `json:"activeCount"`
	RedeemedCount int `json:"redeemedCount"`
	RevokedCount  int `json:"revokedCount"`
}

type GiftCardRevokeRequest struct {
	Reason string `json:"reason,omitempty"`
}

type GiftCardRevokeResponse struct {
	RevokedCount   int    `json:"revokedCount"`
	SkippedCount   int    `json:"skippedCount"`
	RefundedAmount string `json:"refundedAmount"`
}

type PaginationOptions struct {
	Cursor string
	Limit  int
	Sort   string
}

type GiftCardListOptions struct {
	PaginationOptions
	Status     string
	BatchID    string
	Standalone *bool
	Reveal     bool
}

type CursorPage[T any] struct {
	Items  []T     `json:"items"`
	Cursor *string `json:"cursor"`
}

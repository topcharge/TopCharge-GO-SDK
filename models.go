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


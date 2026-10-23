package topcharge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.topcharge.app"

type APIError struct {
	StatusCode   int
	ResponseBody string
	Message      string
}

func (e *APIError) Error() string { return e.Message }

type Client struct {
	clientID   string
	secretKey  string
	baseURL    string
	httpClient *http.Client

	Account     *AccountService
	Payments    *PaymentService
	Withdrawals *WithdrawalService
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(client *Client) { client.baseURL = strings.TrimRight(baseURL, "/") }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) { client.httpClient = httpClient }
}

func NewClient(clientID, secretKey string, options ...Option) (*Client, error) {
	if strings.TrimSpace(clientID) == "" {
		return nil, errors.New("topcharge client ID must not be empty")
	}
	client := &Client{
		clientID: clientID, secretKey: secretKey, baseURL: DefaultBaseURL, httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, option := range options {
		option(client)
	}
	if strings.TrimSpace(client.baseURL) == "" {
		return nil, errors.New("topcharge base URL must not be empty")
	}
	if client.httpClient == nil {
		return nil, errors.New("topcharge HTTP client must not be nil")
	}
	client.Account = &AccountService{client: client}
	client.Payments = &PaymentService{client: client}
	client.Withdrawals = &WithdrawalService{client: client}
	return client, nil
}

func (c *Client) do(ctx context.Context, method, path string, requestBody any, query url.Values, result any) error {
	endpoint := fmt.Sprintf("%s/%s/", c.baseURL, strings.Trim(path, "/"))
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var body io.Reader
	if requestBody != nil {
		encoded, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("encode TopCharge request: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return fmt.Errorf("create TopCharge request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "topcharge-go-sdk/1.0")
	req.Header.Set("x-api-key", c.clientID)
	if c.secretKey != "" {
		req.Header.Set("x-api-secret", c.secretKey)
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send TopCharge request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read TopCharge response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return &APIError{StatusCode: response.StatusCode, ResponseBody: string(responseBody), Message: errorMessage(responseBody, response.StatusCode)}
	}
	if result == nil || len(responseBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(responseBody, result); err != nil {
		return fmt.Errorf("decode TopCharge response: %w", err)
	}
	return nil
}

type AccountService struct{ client *Client }

func (s *AccountService) VerifyCredentials(ctx context.Context) error {
	return s.client.do(ctx, http.MethodGet, "/api/v1/account/verify-credentials", nil, nil, nil)
}

func (s *AccountService) GetBalance(ctx context.Context) (*AccountBalance, error) {
	var result AccountBalance
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/account/balance", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type PaymentService struct{ client *Client }

func (s *PaymentService) Create(ctx context.Context, input PaymentCreateRequest) (*Payment, error) {
	var result Payment
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/payment/create", input, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *PaymentService) GetStatus(ctx context.Context, paymentID string) (*Payment, error) {
	var result Payment
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/payment/status/"+url.PathEscape(paymentID), nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *PaymentService) Cancel(ctx context.Context, paymentID string) (*PaymentCancelResponse, error) {
	var result PaymentCancelResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/payment/cancel/"+url.PathEscape(paymentID), nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type WithdrawalService struct{ client *Client }

func (s *WithdrawalService) Create(ctx context.Context, input WithdrawalCreateRequest) (*Withdrawal, error) {
	var result Withdrawal
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/withdrawal/create", input, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func errorMessage(body []byte, statusCode int) string {
	var response struct {
		Message string `json:"message"`
		Err     struct {
			Message string `json:"message"`
		} `json:"err"`
	}
	if json.Unmarshal(body, &response) == nil {
		if response.Message != "" {
			return response.Message
		}
		if response.Err.Message != "" {
			return response.Err.Message
		}
	}
	return fmt.Sprintf("TopCharge API request failed with HTTP status %d.", statusCode)
}

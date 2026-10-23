package topcharge

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientRequestsAndErrors(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		if request.Header.Get("x-api-key") != "mch_test_example" || request.Header.Get("x-api-secret") != "sk_test_example" {
			t.Fatal("credentials headers are missing")
		}
		switch requestCount {
		case 1:
			if request.URL.Path != "/api/v1/account/balance/" {
				t.Fatalf("unexpected balance path: %s", request.URL.Path)
			}
			_, _ = io.WriteString(writer, `{"availableBalance":"140000.0000"}`)
		case 2:
			if request.URL.Path != "/api/v1/payment/create/" {
				t.Fatalf("unexpected create path: %s", request.URL.Path)
			}
			body, _ := io.ReadAll(request.Body)
			if string(body) != `{"paymentAmount":"25000","description":"Invoice #1001"}` {
				t.Fatalf("unexpected payload: %s", body)
			}
			writer.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(writer, `{"paymentId":"pay_123"}`)
		default:
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(writer, `{"ok":false,"err":{"message":"invalid API credentials"}}`)
		}
	}))
	defer server.Close()

	client, err := NewClient("mch_test_example", "sk_test_example", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}

	balance, err := client.Account.GetBalance(context.Background())
	if err != nil || balance.AvailableBalance != "140000.0000" {
		t.Fatalf("unexpected balance: %#v, %v", balance, err)
	}
	payment, err := client.Payments.Create(context.Background(), PaymentCreateRequest{PaymentAmount: "25000", Description: "Invoice #1001"})
	if err != nil || payment.PaymentID != "pay_123" {
		t.Fatalf("unexpected payment: %#v, %v", payment, err)
	}
	err = client.Account.VerifyCredentials(context.Background())
	apiError, ok := err.(*APIError)
	if !ok || apiError.StatusCode != http.StatusUnauthorized || !strings.Contains(apiError.Error(), "invalid API credentials") {
		t.Fatalf("expected API error, got %v", fmt.Sprintf("%#v", err))
	}
}

package payment

import (
	"context"
	"testing"
)

func TestNewStripeProvider(t *testing.T) {
	apiKey := "sk_test_123"
	provider := NewStripeProvider(apiKey)

	if provider == nil {
		t.Fatal("Expected Stripe provider to be created")
	}
	if provider.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, provider.apiKey)
	}
}

func TestStripeProvider_Charge(t *testing.T) {
	provider := NewStripeProvider("sk_test_123")
	ctx := context.Background()

	t.Run("successfully charges valid amount", func(t *testing.T) {
		result, err := provider.Charge(ctx, 1000, "usd", "tok_visa")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected payment result to be returned")
		}
		if result.Amount != 1000 {
			t.Errorf("Expected amount 1000, got %d", result.Amount)
		}
		if result.Currency != "usd" {
			t.Errorf("Expected currency usd, got %s", result.Currency)
		}
		if result.Status != "succeeded" {
			t.Errorf("Expected status succeeded, got %s", result.Status)
		}
		if result.ID == "" {
			t.Error("Expected charge ID to be set")
		}
	})

	t.Run("returns error for zero amount", func(t *testing.T) {
		_, err := provider.Charge(ctx, 0, "usd", "tok_visa")

		if err != ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount, got %v", err)
		}
	})

	t.Run("returns error for negative amount", func(t *testing.T) {
		_, err := provider.Charge(ctx, -100, "usd", "tok_visa")

		if err != ErrInvalidAmount {
			t.Errorf("Expected ErrInvalidAmount, got %v", err)
		}
	})
}

func TestStripeProvider_Refund(t *testing.T) {
	provider := NewStripeProvider("sk_test_123")
	ctx := context.Background()

	t.Run("successfully refunds charge", func(t *testing.T) {
		err := provider.Refund(ctx, "ch_test123", 1000)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("refunds partial amount", func(t *testing.T) {
		err := provider.Refund(ctx, "ch_test123", 500)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestStripeProvider_CreateSubscription(t *testing.T) {
	provider := NewStripeProvider("sk_test_123")
	ctx := context.Background()

	t.Run("successfully creates subscription", func(t *testing.T) {
		subscription, err := provider.CreateSubscription(ctx, "cus_test123", "plan_test123")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if subscription == nil {
			t.Fatal("Expected subscription to be returned")
		}
		if subscription.CustomerID != "cus_test123" {
			t.Errorf("Expected customer ID cus_test123, got %s", subscription.CustomerID)
		}
		if subscription.PlanID != "plan_test123" {
			t.Errorf("Expected plan ID plan_test123, got %s", subscription.PlanID)
		}
		if subscription.Status != "active" {
			t.Errorf("Expected status active, got %s", subscription.Status)
		}
		if subscription.ID == "" {
			t.Error("Expected subscription ID to be set")
		}
	})
}

func TestNewPayPalProvider(t *testing.T) {
	clientID := "client_id_123"
	clientSecret := "client_secret_123"
	provider := NewPayPalProvider(clientID, clientSecret)

	if provider == nil {
		t.Fatal("Expected PayPal provider to be created")
	}
	if provider.clientID != clientID {
		t.Errorf("Expected client ID %s, got %s", clientID, provider.clientID)
	}
	if provider.clientSecret != clientSecret {
		t.Errorf("Expected client secret %s, got %s", clientSecret, provider.clientSecret)
	}
}

func TestPayPalProvider_Charge(t *testing.T) {
	provider := NewPayPalProvider("client_id", "client_secret")
	ctx := context.Background()

	t.Run("successfully charges amount", func(t *testing.T) {
		result, err := provider.Charge(ctx, 2000, "usd", "paypal_account")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected payment result to be returned")
		}
		if result.Amount != 2000 {
			t.Errorf("Expected amount 2000, got %d", result.Amount)
		}
		if result.Currency != "usd" {
			t.Errorf("Expected currency usd, got %s", result.Currency)
		}
		if result.Status != "COMPLETED" {
			t.Errorf("Expected status COMPLETED, got %s", result.Status)
		}
		if result.ID == "" {
			t.Error("Expected payment ID to be set")
		}
	})
}

func TestPayPalProvider_Refund(t *testing.T) {
	provider := NewPayPalProvider("client_id", "client_secret")
	ctx := context.Background()

	t.Run("successfully refunds payment", func(t *testing.T) {
		err := provider.Refund(ctx, "PAYID-test123", 2000)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

func TestPayPalProvider_CreateSubscription(t *testing.T) {
	provider := NewPayPalProvider("client_id", "client_secret")
	ctx := context.Background()

	t.Run("returns error for not implemented subscription", func(t *testing.T) {
		subscription, err := provider.CreateSubscription(ctx, "cus_test123", "plan_test123")

		if err == nil {
			t.Error("Expected error for not implemented subscription")
		}
		if subscription != nil {
			t.Error("Expected subscription to be nil")
		}
	})
}

func TestPaymentProviderInterface(t *testing.T) {
	ctx := context.Background()

	providers := []struct {
		name     string
		provider PaymentProvider
	}{
		{"Stripe", NewStripeProvider("sk_test_123")},
		{"PayPal", NewPayPalProvider("client_id", "client_secret")},
	}

	for _, p := range providers {
		t.Run(p.name+" implements PaymentProvider", func(t *testing.T) {
			// Test that the provider implements the interface
			var _ PaymentProvider = p.provider

			// Test Charge method exists
			_, err := p.provider.Charge(ctx, 100, "usd", "source")
			// We don't check error because PayPal doesn't validate amount
			_ = err

			// Test Refund method exists
			err = p.provider.Refund(ctx, "charge_id", 100)
			_ = err

			// Test CreateSubscription method exists
			_, err = p.provider.CreateSubscription(ctx, "customer", "plan")
			_ = err
		})
	}
}

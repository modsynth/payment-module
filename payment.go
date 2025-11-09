package payment

import (
	"context"
	"errors"
)

var (
	ErrInvalidAmount = errors.New("invalid payment amount")
	ErrPaymentFailed = errors.New("payment failed")
)

type PaymentProvider interface {
	Charge(ctx context.Context, amount int64, currency, source string) (*PaymentResult, error)
	Refund(ctx context.Context, chargeID string, amount int64) error
	CreateSubscription(ctx context.Context, customerID, planID string) (*Subscription, error)
}

type PaymentResult struct {
	ID       string
	Amount   int64
	Currency string
	Status   string
}

type Subscription struct {
	ID         string
	CustomerID string
	PlanID     string
	Status     string
}

type StripeProvider struct {
	apiKey string
}

func NewStripeProvider(apiKey string) *StripeProvider {
	return &StripeProvider{apiKey: apiKey}
}

func (s *StripeProvider) Charge(ctx context.Context, amount int64, currency, source string) (*PaymentResult, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	// Stripe charge implementation
	return &PaymentResult{
		ID:       "ch_test123",
		Amount:   amount,
		Currency: currency,
		Status:   "succeeded",
	}, nil
}

func (s *StripeProvider) Refund(ctx context.Context, chargeID string, amount int64) error {
	// Stripe refund implementation
	return nil
}

func (s *StripeProvider) CreateSubscription(ctx context.Context, customerID, planID string) (*Subscription, error) {
	// Stripe subscription implementation
	return &Subscription{
		ID:         "sub_test123",
		CustomerID: customerID,
		PlanID:     planID,
		Status:     "active",
	}, nil
}

type PayPalProvider struct {
	clientID     string
	clientSecret string
}

func NewPayPalProvider(clientID, clientSecret string) *PayPalProvider {
	return &PayPalProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (p *PayPalProvider) Charge(ctx context.Context, amount int64, currency, source string) (*PaymentResult, error) {
	// PayPal charge implementation
	return &PaymentResult{
		ID:       "PAYID-test123",
		Amount:   amount,
		Currency: currency,
		Status:   "COMPLETED",
	}, nil
}

func (p *PayPalProvider) Refund(ctx context.Context, chargeID string, amount int64) error {
	// PayPal refund implementation
	return nil
}

func (p *PayPalProvider) CreateSubscription(ctx context.Context, customerID, planID string) (*Subscription, error) {
	// PayPal subscription implementation
	return nil, errors.New("not implemented")
}

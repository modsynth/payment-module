# Payment Module

> Payment service integration (Stripe, PayPal)

Part of the [Modsynth](https://github.com/modsynth) ecosystem.

## Features

- Stripe payment processing
- PayPal payment processing
- Refund management
- Subscription management
- Webhook support

## Installation

```bash
go get github.com/modsynth/payment-module
```

## Quick Start

### Stripe

```go
package main

import (
    "context"
    "github.com/modsynth/payment-module"
)

func main() {
    stripe := payment.NewStripeProvider("sk_test_...")

    // Charge
    result, err := stripe.Charge(context.Background(), 5000, "usd", "tok_visa")
    if err != nil {
        panic(err)
    }
    println("Charge ID:", result.ID)

    // Refund
    err = stripe.Refund(context.Background(), result.ID, 5000)

    // Subscription
    sub, err := stripe.CreateSubscription(context.Background(), "cus_123", "plan_123")
}
```

### PayPal

```go
paypal := payment.NewPayPalProvider("client_id", "client_secret")

result, err := paypal.Charge(context.Background(), 5000, "usd", "")
```

## Supported Providers

- **Stripe** - Full support (charge, refund, subscriptions)
- **PayPal** - Basic support (charge, refund)

## Version

Current version: `v0.1.0`

## License

MIT

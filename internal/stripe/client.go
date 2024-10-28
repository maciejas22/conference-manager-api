package stripe

import (
	"log/slog"

	"github.com/stripe/stripe-go/v80"
)

func SetupStripe(l *slog.Logger) {
	stripe.DefaultLeveledLogger = NewStripeLogger(l)
}

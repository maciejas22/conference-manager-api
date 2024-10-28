package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
)

type PaymentIntentMetadata struct {
	ConferenceId int    `json:"conference_id"`
	UserId       int    `json:"user_id"`
	Destination  string `json:"destination"`
}

func CreateConferencePaymentIntent(amount int, metadata PaymentIntentMetadata) (*string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		TransferData: &stripe.PaymentIntentTransferDataParams{
			Destination: stripe.String(metadata.Destination),
		},
		Metadata: map[string]string{
			"conference_id": strconv.Itoa(metadata.ConferenceId),
			"user_id":       strconv.Itoa(metadata.UserId),
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, errors.New("error creating payment intent")
	}

	return stripe.String(pi.ClientSecret), nil
}

func HandlePaymentIntentConfirmation(ctx context.Context, dbClient *db.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request body", http.StatusServiceUnavailable)
			return
		}

		event := stripe.Event{}
		if err := json.Unmarshal(payload, &event); err != nil {
			http.Error(w, "failed to parse webhook body json", http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "payment_intent.succeeded":
			var paymentIntent stripe.PaymentIntent
			err := json.Unmarshal(event.Data.Raw, &paymentIntent)
			if err != nil {
				http.Error(w, "error parsing webhook json", http.StatusBadRequest)
				return
			}

			userId, err := strconv.Atoi(paymentIntent.Metadata["user_id"])
			if err != nil {
				http.Error(w, "invalid userId in metadata", http.StatusBadRequest)
				return
			}

			conferenceId, err := strconv.Atoi(paymentIntent.Metadata["conference_id"])
			if err != nil {
				http.Error(w, "invalid conferenceId in metadata", http.StatusBadRequest)
				return
			}

			if _, err := AddUserToConference(context.Background(), dbClient, userId, conferenceId); err != nil {
				http.Error(w, "error adding user to conference", http.StatusInternalServerError)
			}

		default:
			logger.Error("unhandled event type", "type", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}

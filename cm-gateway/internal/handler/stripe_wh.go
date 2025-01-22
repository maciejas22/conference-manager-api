package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	conferencePb "github.com/maciejas22/conference-manager-api/cm-proto/conference"
	stripe "github.com/stripe/stripe-go/v81"
)

func HandlePaymentIntentConfirmation(ctx context.Context, logger *slog.Logger, conferenceService conferencePb.ConferenceServiceClient) http.HandlerFunc {
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

			if _, err := conferenceService.AddUserToConference(ctx, &conferencePb.AddUserToConferenceRequest{
				UserId:       int32(userId),
				ConferenceId: int32(conferenceId),
			}); err != nil {
				http.Error(w, "error adding user to conference", http.StatusInternalServerError)
			}

		default:
			logger.Error("unhandled event type", "type", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}

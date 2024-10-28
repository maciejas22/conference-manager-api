package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
	s "github.com/maciejas22/conference-manager/api/internal/stripe"
	"github.com/stripe/stripe-go/v80/account"
)

func GetUserData(ctx context.Context, dbClient *db.DB, userId int) (*models.User, error) {
	var u repositories.User
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		u, err = repositories.GetUserByID(tx, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       u.Id,
		Name:     u.Name,
		Surname:  u.Surname,
		Username: u.Username,
		Email:    *u.Email,
		Role:     models.Role(u.Role),
	}, nil
}

func StripeOnboard(ctx context.Context, dbClient *db.DB, userId int, returnURL string, refreshURL string) (string, error) {
	var u repositories.User
	var url *string
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		u, err = repositories.GetUserByID(tx, userId)
		if err != nil {
			return err
		}

		url, err = s.CreateAccountLink(*u.StripeAccountId, returnURL, refreshURL)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", errors.New("Failed to create account link")
	}

	return *url, nil
}

func GetStripeAccount(ctx context.Context, dbClient *db.DB, userId int) (*models.StripeAccountDetails, error) {
	var stripeAccountId string
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		u, err := repositories.GetUserByID(tx, userId)
		if err != nil {
			return err
		}

		stripeAccountId = *u.StripeAccountId
		return nil
	})
	if err != nil {
		return nil, err
	}

	acc, err := account.GetByID(stripeAccountId, nil)
	if err != nil {
		return nil, errors.New("Failed to get account details")
	}

	return &models.StripeAccountDetails{
		ID:         acc.ID,
		IsVerified: acc.ChargesEnabled && acc.PayoutsEnabled,
	}, nil

}

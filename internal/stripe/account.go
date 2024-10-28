package stripe

import (
	"errors"

	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/account"
	"github.com/stripe/stripe-go/v80/accountlink"
)

func CreateStripeAccount(userEmail string) (*string, error) {
	account, err := account.New(&stripe.AccountParams{
		Email: stripe.String(userEmail),
	})

	if err != nil {
		return nil, errors.New("Error creating stripe account")
	}

	return &account.ID, nil
}

func CreateAccountLink(userAccountId, returnUrl, refreshUrl string) (*string, error) {
	accountLink, err := accountlink.New(&stripe.AccountLinkParams{
		Account:    stripe.String(userAccountId),
		ReturnURL:  stripe.String(returnUrl),
		RefreshURL: stripe.String(refreshUrl),
		Type:       stripe.String("account_onboarding"),
	})

	if err != nil {
		return nil, errors.New("Error creating account link")
	}
	return &accountLink.URL, nil
}

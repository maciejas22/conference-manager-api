package repository

type ToSRepoInterface interface {
	GetTermsOfService() (*TermsOfService, error)
}

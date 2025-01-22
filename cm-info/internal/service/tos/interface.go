package service

type ToSServiceInterface interface {
	GetTermsOfService() (*TermsOfService, error)
}

package service

type Role string

const (
	Organizer   Role = "Organizer"
	Participant Role = "Participant"
)

type User struct {
	Id              int
	SessionId       string
	Email           string
	PasswordHash    string
	Name            *string
	Surname         *string
	Username        *string
	Role            Role
	StripeAccountId *string
}

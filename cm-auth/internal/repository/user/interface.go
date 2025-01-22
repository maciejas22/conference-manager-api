package repository

type UserRepoInterface interface {
	GetUserByID(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserBySessionId(sessionId string) (User, error)
	CreateUser(email string, passwordHash string, role Role, stripeAccountId *string) (int, error)
	UpdateUser(id int, input UpdateUserInput) (int, error)
}

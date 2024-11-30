package service

type AuthServiceInterface interface {
	RegisterUser(userData RegisterUserInput) (*int, error)
	LoginUser(userData LoginUserInput) (*string, error)
	UpdateSession(sessionId string) (string, int, error)
	VerifySession(sessionId string) (bool, error)
	UpdateUser(id int, input UpdateUserInput) (int, error)
	LogoutUser(sessionId string) (*bool, error)
	GetUser(userId int) (*User, error)
	GetUserBySession(sessionId string) (*User, error)
}

package service

import (
	"errors"
	"time"

	sessionRepository "github.com/maciejas22/conference-manager-api/cm-auth/internal/repository/session"
	userRepository "github.com/maciejas22/conference-manager-api/cm-auth/internal/repository/user"
	"github.com/maciejas22/conference-manager-api/cm-auth/internal/utils/hash"
)

type UserService struct {
	userRepo    userRepository.UserRepoInterface
	sessionRepo sessionRepository.SessionRepoInterface
}

func NewAuthService(userRepo userRepository.UserRepoInterface, sessionRepo sessionRepository.SessionRepoInterface) AuthServiceInterface {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type RegisterUserInput struct {
	Email           string
	Password        string
	Role            Role
	StripeAccountId *string
}

func (s *UserService) RegisterUser(userData RegisterUserInput) (*int, error) {
	hashedPassword, err := hash.HashPassword(userData.Password)
	if err != nil {
		return nil, errors.New("Could not hash password")
	}

	uId, err := s.userRepo.CreateUser(userData.Email, hashedPassword, userRepository.Role(userData.Role), userData.StripeAccountId)
	if err != nil {
		return nil, err
	}

	return &uId, nil
}

type LoginUserInput struct {
	Email    string
	Password string
}

func (s *UserService) LoginUser(userData LoginUserInput) (*string, error) {
	u, err := s.userRepo.GetUserByEmail(userData.Email)
	if err != nil {
		return nil, err
	}

	if !hash.CheckPasswordHash(userData.Password, u.PasswordHash) {
		return nil, errors.New("Invalid email or password")
	}

	sessionId, err := s.sessionRepo.CreateSession(u.Id)
	if err != nil {
		return nil, err
	}

	return sessionId, nil
}

func (s *UserService) GetUser(userId int) (*User, error) {
	u, err := s.userRepo.GetUserByID(userId)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:              u.Id,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		Name:            u.Name,
		Surname:         u.Surname,
		Username:        u.Username,
		StripeAccountId: u.StripeAccountId,
	}, nil
}

func (s *UserService) GetUserBySession(sessionId string) (*User, error) {
	u, err := s.userRepo.GetUserBySessionId(sessionId)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:              u.Id,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		Name:            u.Name,
		Surname:         u.Surname,
		Username:        u.Username,
		Role:            Role(u.Role),
		StripeAccountId: u.StripeAccountId,
	}, nil
}

func (s *UserService) UpdateSession(sessionId string) (string, int, error) {
	session, err := s.sessionRepo.GetSessionBySessionId(sessionId)
	if err != nil {
		return "", 0, err
	}

	now := time.Now()
	if now.After(session.ExpiresAt) {
		return "", 0, errors.New("Session expired")
	}

	newSession := &sessionId
	if now.After(session.ExpiresAt.Add(-5 * time.Minute)) {
		newSession, err = s.sessionRepo.CreateSession(session.UserId)
		if err != nil {
			return "", 0, errors.New("Could not create new session")
		}
	}

	return *newSession, session.UserId, nil
}

func (s *UserService) LogoutUser(sessionId string) (*bool, error) {
	err := s.sessionRepo.DestroySession(sessionId)
	if err != nil {
		return nil, err
	}

	ok := true
	return &ok, nil
}

type UpdateUserInput struct {
	Email    *string
	Password *string
	Name     *string
	Surname  *string
	Username *string
}

func (s *UserService) UpdateUser(userId int, input UpdateUserInput) (int, error) {
	var updateUserInput userRepository.UpdateUserInput
	if input.Email != nil {
		updateUserInput.Email = input.Email
	}
	if input.Password != nil {
		hashedPassword, err := hash.HashPassword(*input.Password)
		if err != nil {
			return 0, errors.New("Could not hash password")
		}
		updateUserInput.Password = &hashedPassword
	}
	if input.Name != nil {
		updateUserInput.Name = input.Name
	}
	if input.Surname != nil {
		updateUserInput.Surname = input.Surname
	}
	if input.Username != nil {
		updateUserInput.Username = input.Username
	}

	id, err := s.userRepo.UpdateUser(userId, updateUserInput)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) VerifySession(sessionId string) (bool, error) {
	session, err := s.sessionRepo.GetSessionBySessionId(sessionId)
	if err != nil {
		return false, err
	}

	now := time.Now()
	if now.After(session.ExpiresAt) {
		return false, errors.New("Session expired")
	}

	return true, nil
}

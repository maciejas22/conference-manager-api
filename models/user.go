package models

import "github.com/maciejas22/conference-manager/api/db/repositories"

type User struct {
	ID       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Surname  *string `json:"surname,omitempty"`
	Username *string `json:"username"`
	Email    string  `json:"email"`
	Role     Role    `json:"role"`
}

func (u User) ToRepo() repositories.User {
	return repositories.User{
		PublicUser: repositories.PublicUser{
			Id:       u.ID,
			Name:     u.Name,
			Surname:  u.Surname,
			Username: u.Username,
			Role:     repositories.Role(u.Role),
		},
		Email: u.Email,
	}
}

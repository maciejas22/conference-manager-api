package models

type User struct {
	ID       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Surname  *string `json:"surname,omitempty"`
	Username *string `json:"username"`
	Email    string  `json:"email"`
	Role     Role    `json:"role"`
}

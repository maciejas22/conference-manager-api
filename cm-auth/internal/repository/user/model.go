package repository

import "time"

type Role string

const (
	Organizer   Role = "Organizer"
	Participant Role = "Participant"
)

type User struct {
	Id              int       `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"password" db:"password"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Name            *string   `json:"name" db:"name"`
	Surname         *string   `json:"surname" db:"surname"`
	Username        *string   `json:"username" db:"username"`
	Role            Role      `json:"role" db:"role"`
	StripeAccountId *string   `json:"stripe_account_id" db:"stripe_account_id"`
}

func (u *User) TableName() string {
	return "public.users"
}

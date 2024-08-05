package repositories

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type Role string

const (
	Organizer   Role = "Organizer"
	Participant Role = "Participant"
)

type PublicUser struct {
	Id        string  `json:"id" db:"id"`
	Name      *string `json:"name" db:"name"`
	Surname   *string `json:"surname" db:"surname"`
	Username  *string `json:"username" db:"username"`
	Role      Role    `json:"role" db:"role"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}

func (u *PublicUser) TableName() string {
	return "public.users"
}

type AuthUser struct {
	InstanceId               *string          `json:"instance_id,omitempty" db:"instance_id"`
	Id                       string           `json:"id" db:"id"`
	Aud                      *string          `json:"aud,omitempty" db:"aud"`
	Role                     *string          `json:"role,omitempty" db:"role"`
	Email                    *string          `json:"email,omitempty" db:"email"`
	EncryptedPassword        *string          `json:"encrypted_password,omitempty" db:"encrypted_password"`
	EmailConfirmedAt         *string          `json:"email_confirmed_at,omitempty" db:"email_confirmed_at"`
	InvitedAt                *string          `json:"invited_at,omitempty" db:"invited_at"`
	ConfirmationToken        *string          `json:"confirmation_token,omitempty" db:"confirmation_token"`
	ConfirmationSentAt       *string          `json:"confirmation_sent_at,omitempty" db:"confirmation_sent_at"`
	RecoveryToken            *string          `json:"recovery_token,omitempty" db:"recovery_token"`
	RecoverySentAt           *string          `json:"recovery_sent_at,omitempty" db:"recovery_sent_at"`
	EmailChangeTokenNew      *string          `json:"email_change_token_new,omitempty" db:"email_change_token_new"`
	EmailChange              *string          `json:"email_change,omitempty" db:"email_change"`
	EmailChangeSentAt        *string          `json:"email_change_sent_at,omitempty" db:"email_change_sent_at"`
	LastSignInAt             *string          `json:"last_sign_in_at,omitempty" db:"last_sign_in_at"`
	RawAppMetaData           *json.RawMessage `json:"raw_app_meta_data,omitempty" db:"raw_app_meta_data"`
	RawUserMetaData          *json.RawMessage `json:"raw_user_meta_data,omitempty" db:"raw_user_meta_data"`
	IsSuperAdmin             *bool            `json:"is_super_admin,omitempty" db:"is_super_admin"`
	CreatedAt                *string          `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt                *string          `json:"updated_at,omitempty" db:"updated_at"`
	Phone                    *string          `json:"phone,omitempty" db:"phone"`
	PhoneConfirmedAt         *string          `json:"phone_confirmed_at,omitempty" db:"phone_confirmed_at"`
	PhoneChange              *string          `json:"phone_change,omitempty" db:"phone_change"`
	PhoneChangeToken         *string          `json:"phone_change_token,omitempty" db:"phone_change_token"`
	PhoneChangeSentAt        *string          `json:"phone_change_sent_at,omitempty" db:"phone_change_sent_at"`
	ConfirmedAt              *string          `json:"confirmed_at,omitempty" db:"confirmed_at"`
	EmailChangeTokenCurrent  *string          `json:"email_change_token_current,omitempty" db:"email_change_token_current"`
	EmailChangeConfirmStatus *int16           `json:"email_change_confirm_status,omitempty" db:"email_change_confirm_status"`
	BannedUntil              *string          `json:"banned_until,omitempty" db:"banned_until"`
	ReauthenticationToken    *string          `json:"reauthentication_token,omitempty" db:"reauthentication_token"`
	ReauthenticationSentAt   *string          `json:"reauthentication_sent_at,omitempty" db:"reauthentication_sent_at"`
	IsSSOUser                bool             `json:"is_sso_user" db:"is_sso_user"`
	DeletedAt                *string          `json:"deleted_at,omitempty" db:"deleted_at"`
	IsAnonymous              bool             `json:"is_anonymous" db:"is_anonymous"`
}

func (u *AuthUser) TableName() string {
	return "auth.users"
}

type User struct {
	PublicUser
	Email string `json:"email"`
}

func GetUserByID(tx *sqlx.Tx, id string) (User, error) {
	var user User

	p := &PublicUser{}
	a := &AuthUser{}
	query := "SELECT u.id, u.name, u.surname, u.username, u.role, a.email, a.created_at, a.updated_at FROM " + p.TableName() + " u JOIN " + a.TableName() + " a ON u.id = a.id WHERE u.id = $1"
	err := tx.Get(
		&user,
		query,
		id,
	)
	if err != nil {
		log.Println(err.Error())
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func UpdateUser(
	tx *sqlx.Tx,
	id string,
	username string,
	email string,
	name string,
	surname string,
) (User, error) {
	var user User

	p := &PublicUser{}
	a := &AuthUser{}

	query := "UPDATE " + p.TableName() + " SET name = $1, surname = $2, username = $3 WHERE id = $4"
	_, err := tx.Exec(query, name, surname, username, id)
	if err != nil {
		return User{}, err
	}

	query = "UPDATE " + a.TableName() + " SET email = $1 WHERE id = $2"
	_, err = tx.Exec(query, email, id)
	if err != nil {
		return User{}, errors.New("could not update user")
	}

	user, err = GetUserByID(tx, id)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

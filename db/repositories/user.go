package repositories

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Role string

const (
	Organizer   Role = "Organizer"
	Participant Role = "Participant"
)

type User struct {
	Id        int     `json:"id" db:"id"`
	Name      *string `json:"name" db:"name"`
	Surname   *string `json:"surname" db:"surname"`
	Username  *string `json:"username" db:"username"`
	Email     *string `json:"email" db:"email"`
	Role      Role    `json:"role" db:"role"`
	Password  string  `json:"password" db:"password"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}

func (u *User) TableName() string {
	return "public.users"
}

func GetUserByID(tx *sqlx.Tx, id int) (User, error) {
	var user User

	query := `
    SELECT id, name, surname, username, role, email, created_at, updated_at, password 
    FROM ` + user.TableName() + ` 
    WHERE id = $1 
  `
	err := tx.Get(
		&user,
		query,
		id,
	)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

type UpdateUserInput struct {
	Username *string
	Email    *string
	Name     *string
	Surname  *string
	Password *string
}

func UpdateUser(
	tx *sqlx.Tx,
	id int,
	input UpdateUserInput,
) (User, error) {
	var user User

	fields := []string{}
	args := []interface{}{}

	if input.Username != nil {
		fields = append(fields, "username = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Username)
	}
	if input.Email != nil {
		fields = append(fields, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Email)
	}
	if input.Name != nil {
		fields = append(fields, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Name)
	}
	if input.Surname != nil {
		fields = append(fields, "surname = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Surname)
	}
	if input.Password != nil {
		fields = append(fields, "password = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Password)
	}

	if len(fields) == 0 {
		return User{}, errors.New("no fields to update")
	}

	query := `
    UPDATE ` + user.TableName() + ` 
    SET ` + strings.Join(fields, ", ") + ` 
    WHERE id = $` + strconv.Itoa(len(args)+1)
	args = append(args, id)

	_, err := tx.Exec(query, args...)
	if err != nil {
		return User{}, errors.New("could not update user")
	}

	user, err = GetUserByID(tx, id)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func CreateUser(tx *sqlx.Tx, email string, passwordHash string, role string) (int, error) {
	query := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"
	var userId int
	err := tx.QueryRow(query, email, passwordHash, role).Scan(&userId)
	if err != nil {
		return 0, errors.New("could not create user")
	}

	return userId, nil
}

func GetUserBySessionId(tx *sqlx.Tx, sessionId string) (User, error) {
	var user User

	query := `
    SELECT u.id, u.name, u.surname, u.username, u.role, u.email, u.created_at, u.updated_at, u.password
    FROM ` + user.TableName() + ` u 
    JOIN sessions s ON u.id = s.user_id 
    WHERE s.session_id = $1
  `
	err := tx.Get(
		&user,
		query,
		sessionId,
	)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByEmail(tx *sqlx.Tx, email string) (User, error) {
	var user User

	query := `
    SELECT u.id, u.name, u.surname, u.username, u.role, u.email, u.role, u.created_at, u.updated_at, u.password
    FROM ` + user.TableName() + ` u 
    WHERE u.email = $1
  `
	err := tx.Get(
		&user,
		query,
		email,
	)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

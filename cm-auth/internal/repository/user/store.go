package repository

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	Db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepoInterface {
	return &UserRepo{Db: db}
}

func (r *UserRepo) GetUserByID(id int) (User, error) {
	var user User

	query := `
    SELECT id, email, password, created_at, updated_at, name, surname, username, role, stripe_account_id
    FROM ` + user.TableName() + ` 
    WHERE id = $1 
  `
	err := r.Db.Get(
		&user,
		query,
		id,
	)
	if err != nil {
		return User{}, errors.New("User not found")
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (User, error) {
	var user User

	query := `
    SELECT id, email, password, created_at, updated_at, name, surname, username, role, stripe_account_id
    FROM ` + user.TableName() + `
    WHERE email = $1
  `
	err := r.Db.Get(
		&user,
		query,
		email,
	)
	if err != nil {
		return User{}, errors.New("User not found")
	}

	return user, nil
}

func (r *UserRepo) GetUserBySessionId(sessionId string) (User, error) {
	var user User

	query := `
    SELECT u.id, u.email, u.password, u.created_at, u.updated_at, u.name, u.surname, u.username, u.role, u.stripe_account_id
    FROM ` + user.TableName() + ` u 
    JOIN sessions s ON u.id = s.user_id 
    WHERE s.session_id = $1
  `
	err := r.Db.Get(
		&user,
		query,
		sessionId,
	)
	if err != nil {
		return User{}, errors.New("User not found")
	}

	return user, nil
}

func (r *UserRepo) CreateUser(email string, passwordHash string, role Role, stripeAccountId *string) (int, error) {
	query := "INSERT INTO users (email, password, role, stripe_account_id) VALUES ($1, $2, $3, $4) RETURNING id"
	var userId int
	err := r.Db.QueryRow(query, email, passwordHash, role, stripeAccountId).Scan(&userId)
	if err != nil {
		return 0, errors.New("Could not create user")
	}

	return userId, nil
}

type UpdateUserInput struct {
	Email    *string
	Password *string
	Name     *string
	Surname  *string
	Username *string
}

func (r *UserRepo) UpdateUser(
	id int,
	input UpdateUserInput,
) (int, error) {
	var user User

	fields := []string{}
	args := []interface{}{}

	if input.Email != nil {
		fields = append(fields, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Email)
	}
	if input.Password != nil {
		fields = append(fields, "password = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Password)
	}
	if input.Name != nil {
		fields = append(fields, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Name)
	}
	if input.Surname != nil {
		fields = append(fields, "surname = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Surname)
	}
	if input.Username != nil {
		fields = append(fields, "username = $"+strconv.Itoa(len(args)+1))
		args = append(args, *input.Username)
	}

	if len(fields) == 0 {
		return 0, errors.New("No fields to update")
	}
	fields = append(fields, "updated_at = NOW()")

	query := `
    UPDATE ` + user.TableName() + ` 
    SET ` + strings.Join(fields, ", ") + ` 
    WHERE id = $` + strconv.Itoa(len(args)+1)
	args = append(args, id)

	_, err := r.Db.Exec(query, args...)
	if err != nil {
		return 0, errors.New("Could not update user")
	}

	return id, nil
}

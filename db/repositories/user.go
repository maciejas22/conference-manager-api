package repositories

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
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
	return "users"
}

func GetUserByID(qe *db.QueryExecutor, id int) (User, error) {
	var user User

	query := `
    SELECT id, name, surname, username, role, email, created_at, updated_at, password 
    FROM ` + user.TableName() + ` 
    WHERE id = ?
  `
	err := sqlx.Get(
		qe,
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
	qe *db.QueryExecutor,
	id int,
	username string,
	email string,
	name string,
	surname string,
) (User, error) {
	var user User

	query := `
    UPDATE ` + user.TableName() + ` 
    SET name = ?, surname = ?, username = ?, email = ? 
    WHERE id = ?
  `
	_, err := qe.Exec(query, name, surname, username, email, id)
	if err != nil {
		return User{}, errors.New("could not update user")
	}

	user, err = GetUserByID(qe, id)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func CreateUser(qe *db.QueryExecutor, email string, passwordHash string, role string) (int, error) {
	query := "INSERT INTO users (email, password, role) VALUES (?, ?, ?)"
	res, err := qe.Exec(
		query,
		email,
		passwordHash,
		role,
	)
	if err != nil {
		return 0, errors.New("could not create user")
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("could not get user id")
	}

	return int(userId), nil
}

func GetUserBySessionId(qe *db.QueryExecutor, sessionId string) (User, error) {
	var user User

	query := `
    SELECT u.id, u.name, u.surname, u.username, u.role, u.email, u.created_at, u.updated_at, u.password
    FROM ` + user.TableName() + ` u 
    JOIN sessions s ON u.id = s.user_id 
    WHERE s.session_id = ?
  `
	err := sqlx.Get(
		qe,
		&user,
		query,
		sessionId,
	)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByEmail(qe *db.QueryExecutor, email string) (User, error) {
	var user User

	query := `
    SELECT u.id, u.name, u.surname, u.username, u.role, u.email, u.role, u.created_at, u.updated_at, u.password
    FROM ` + user.TableName() + ` u 
    WHERE u.email = ?
  `
	err := sqlx.Get(
		qe,
		&user,
		query,
		email,
	)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

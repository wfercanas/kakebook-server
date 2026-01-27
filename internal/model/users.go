package model

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	userId uuid.UUID
	name   string
	email  string
}

func (u *User) String() string {
	return fmt.Sprintf(
		`{
		user_id: %s,
		name: %s,
		email: %s
	}`, u.userId, u.name, u.email)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Get(userId uuid.UUID) (User, error) {
	stmt := `SELECT user_id, name, email
	FROM users
	WHERE user_id = $1`

	row := m.DB.QueryRow(stmt, userId)

	var u User
	err := row.Scan(&u.userId, &u.name, &u.email)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

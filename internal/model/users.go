package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	UserId uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Get(userId uuid.UUID) (User, error) {
	stmt := `SELECT user_id, name, email
	FROM users
	WHERE user_id = $1`

	row := m.DB.QueryRow(stmt, userId)

	var user User
	err := row.Scan(&user.UserId, &user.Name, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

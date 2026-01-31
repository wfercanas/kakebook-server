package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Account struct {
	AccountId       uuid.UUID
	Name            string
	AccountCategory string
	ProjectId       uuid.UUID
}

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(name string, accountCategory string, projectId uuid.UUID) error {
	stmt := `INSERT INTO accounts (account_id, name, account_category, project_id)
	VALUES (gen_random_uuid(), $1, $2, $3) RETURNING account_id`

	var accountId uuid.UUID
	err := m.DB.QueryRow(stmt, name, accountCategory, projectId).Scan(&accountId)
	if err != nil {
		return err
	}

	return nil
}

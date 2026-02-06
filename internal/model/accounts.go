package model

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Account struct {
	AccountId       uuid.UUID `json:"account_id"`
	AccountName     string    `json:"account_name"`
	AccountCategory string    `json:"account_category"`
	ProjectId       uuid.UUID `json:"project_id"`
	Balance         float64   `json:"balance"`
}

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(name string, accountCategory string, projectId uuid.UUID) error {
	stmt := `INSERT INTO accounts (account_id, account_name, account_category, project_id)
	VALUES (gen_random_uuid(), $1, $2, $3) RETURNING account_id`

	var accountId uuid.UUID
	err := m.DB.QueryRow(stmt, name, accountCategory, projectId).Scan(&accountId)
	if err != nil {
		return err
	}

	return nil
}

func (m *AccountModel) GetAccountById(id uuid.UUID) (Account, error) {
	accountStmt := `SELECT account_name, project_id, account_category
	FROM accounts 
	WHERE account_id = $1`

	var account Account
	account.AccountId = id

	result := m.DB.QueryRow(accountStmt, id)
	err := result.Scan(&account.AccountName, &account.ProjectId, &account.AccountCategory)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Account{}, ErrNoRecord
		} else {
			return Account{}, err
		}
	}

	movementsStmt := `SELECT movement_type, sum(value)
	FROM movements
	WHERE account_id = $1
	GROUP BY movement_type
	`

	balance := 0.0
	type subtotal struct {
		movement_type string
		amount        float64
	}

	rows, err := m.DB.Query(movementsStmt, account.AccountId)
	if err != nil {
		return Account{}, err
	}

	for rows.Next() {
		var s subtotal
		err = rows.Scan(&s.movement_type, &s.amount)
		if err != nil {
			return Account{}, err
		}

		if account.AccountCategory == "assets" || account.AccountCategory == "expenses" {
			if s.movement_type == "debit" {
				balance += s.amount
			} else {
				balance -= s.amount
			}
		}

		if account.AccountCategory != "assets" && account.AccountCategory != "expenses" {
			if s.movement_type == "credit" {
				balance += s.amount
			} else {
				balance -= s.amount
			}
		}
	}

	account.Balance = balance

	return account, nil
}

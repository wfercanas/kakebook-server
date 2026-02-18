package model

import (
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type LedgerAccount struct {
	AccountName     string           `json:"account_name"`
	AccountCategory string           `json:"account_category"`
	Balance         float64          `json:"balance"`
	AccountId       uuid.UUID        `json:"account_id"`
	ProjectId       uuid.UUID        `json:"project_id"`
	Movements       []LedgerMovement `json:"movements"`
}

type LedgerMovement struct {
	Description  string    `json:"description"`
	Date         time.Time `json:"date"`
	MovementType string    `json:"movement_type"`
	Value        float64   `json:"value"`
	Balance      float64   `json:"balance"`
	EntryId      uuid.UUID `json:"entry_id"`
}

type LedgerModel struct {
	DB *sql.DB
}

func (m *LedgerModel) GetLedgerAccountById(accountId uuid.UUID) (LedgerAccount, error) {
	accountStmt := `SELECT account_id, account_name, account_category, project_id
	FROM accounts
	WHERE account_id = $1`

	movementsStmt := `SELECT mov.movement_type, mov.value, mov.entry_id, ent.date, ent.description 
	FROM movements mov
	JOIN entries ent
	ON mov.entry_id = ent.entry_id
	WHERE mov.account_id = $1`

	var account LedgerAccount

	err := m.DB.QueryRow(accountStmt, accountId).Scan(&account.AccountId, &account.AccountName, &account.AccountCategory, &account.ProjectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LedgerAccount{}, ErrNoRecord
		} else {
			return LedgerAccount{}, err
		}
	}

	result, err := m.DB.Query(movementsStmt, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LedgerAccount{}, ErrNoRecord
		} else {
			return LedgerAccount{}, err
		}
	}

	balance := 0.0
	for result.Next() {
		var movement LedgerMovement
		err = result.Scan(&movement.MovementType, &movement.Value, &movement.EntryId, &movement.Date, &movement.Description)
		if err != nil {
			return LedgerAccount{}, err
		}

		var debitAccountCategories []string
		debitAccountCategories = append(debitAccountCategories, "assets", "expenses")

		if slices.Contains(debitAccountCategories, account.AccountCategory) {
			if movement.MovementType == "debit" {
				balance += movement.Value
			} else {
				balance -= movement.Value
			}
		} else {
			if movement.MovementType == "debit" {
				balance += movement.Value
			} else {
				balance -= movement.Value
			}
		}

		movement.Balance = balance
		account.Movements = append(account.Movements, movement)
	}

	account.Balance = balance
	return account, nil
}

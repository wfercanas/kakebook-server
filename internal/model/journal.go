package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	Amount      float64    `json:"amount"`
	ProjectId   uuid.UUID  `json:"project_id"`
	EntryId     uuid.UUID  `json:"entry_id"`
	Movements   []Movement `json:"movements"`
}

type Movement struct {
	AccountName     string    `json:"account_name"`
	AccountCategory string    `json:"account_category"`
	MovementType    string    `json:"movement_type"`
	Value           float32   `json:"value"`
	AccountId       uuid.UUID `json:"account_id"`
}

type NewEntry struct {
	ProjectId   uuid.UUID     `json:"project_id"`
	Date        string        `json:"date"`
	Description string        `json:"description"`
	Movements   []NewMovement `json:"movements"`
	Amount      float64       `json:"amount"`
}

type NewMovement struct {
	AccountId    uuid.UUID `json:"account_id"`
	MovementType string    `json:"movement_type"`
	Value        float64   `json:"value"`
}

type EntryModel struct {
	DB *sql.DB
}

func (m *EntryModel) Get(entryId uuid.UUID) (Entry, error) {
	entryStmt := `SELECT description, date, amount, project_id
	FROM entries
	WHERE entry_id = $1`

	movementsStmt := `SELECT ac.account_name, ac.account_category, mv.movement_type, mv.value, mv.account_id
	FROM movements mv
	JOIN accounts ac
	ON mv.account_id = ac.account_id
	WHERE mv.entry_id = $1`

	var entry Entry

	result := m.DB.QueryRow(entryStmt, entryId)
	err := result.Scan(&entry.Description, &entry.Date, &entry.Amount, &entry.ProjectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNoRecord
		} else {
			return Entry{}, err
		}
	}

	entry.EntryId = entryId

	results, err := m.DB.Query(movementsStmt, entryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNoRecord
		} else {
			return Entry{}, err
		}
	}

	for results.Next() {
		var movement Movement

		err := results.Scan(&movement.AccountName, &movement.AccountCategory, &movement.MovementType, &movement.Value, &movement.AccountId)
		if err != nil {
			return Entry{}, err
		}

		entry.Movements = append(entry.Movements, movement)
	}

	return entry, nil
}

func (m *EntryModel) Delete(entryId uuid.UUID) error {
	movementStmt := `DELETE FROM movements
	WHERE entry_id = $1`

	entryStmt := `DELETE FROM entries
	WHERE entry_id = $1`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(movementStmt, entryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}

	_, err = tx.Exec(entryStmt, entryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}

	tx.Commit()
	return nil

}

func (m *EntryModel) Insert(newEntry NewEntry) error {
	entryStmt := `INSERT INTO entries (entry_id, date, description, project_id, amount)
	VALUES (gen_random_uuid(), $1, $2, $3, $4) RETURNING entry_id`

	movementStmt := `INSERT INTO movements (value, movement_type, account_id, entry_id)
	VALUES ($1, $2, $3, $4)`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var entryId uuid.UUID

	err = tx.QueryRow(entryStmt, newEntry.Date, newEntry.Description, newEntry.ProjectId, newEntry.Amount).Scan(&entryId)
	if err != nil {
		return err
	}

	for _, movement := range newEntry.Movements {
		_, err := tx.Exec(movementStmt, movement.Value, movement.MovementType, movement.AccountId, entryId)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

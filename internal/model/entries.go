package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	entryId     uuid.UUID
	date        time.Time
	description string
	projectId   uuid.UUID
	amount      float64
	movements   []Movement
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
	return Entry{}, nil
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

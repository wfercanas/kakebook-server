package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	entryId     uuid.UUID
	date        time.Time
	description string
	projectId   uuid.UUID
	amount      float32
	movements   []Movement
}

type EntryModel struct {
	DB *sql.DB
}

func (m *EntryModel) Get(entryId uuid.UUID) (Entry, error) {
	return Entry{}, nil
}

func (m *EntryModel) Insert(date time.Time, description string, projectId uuid.UUID, amount float32) error {
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

	err = tx.QueryRow(entryStmt, date, description, projectId, amount).Scan(&entryId)
	if err != nil {
		return err
	}

	bancolombiaId, err := uuid.Parse("418d3211-fb86-44e4-b83f-309107e51473")
	if err != nil {
		return err
	}
	nuId, err := uuid.Parse("e1d7c3ac-4156-4b19-97e2-782097b14e60")
	if err != nil {
		return err
	}

	_, err = tx.Exec(movementStmt, amount, "debit", bancolombiaId, entryId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(movementStmt, amount, "credit", nuId, entryId)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

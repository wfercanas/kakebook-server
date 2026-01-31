package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Movement struct {
	Value        float32
	MovementType string
	AccountId    uuid.UUID
	EntryId      uuid.UUID
}

type MovementModel struct {
	DB *sql.DB
}

func (m *MovementModel) Insert(tx *sql.Tx, value float32, movementType string, accountId uuid.UUID, entryId uuid.UUID) error {
	stmt := `INSERT INTO movements (value, movement_type, account_id, entry_id)
	VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(stmt, value, movementType, accountId, entryId)
	if err != nil {
		return err
	}

	return nil
}

package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Movement struct {
	AccountName     string    `json:"account_name"`
	AccountCategory string    `json:"account_category"`
	MovementType    string    `json:"movement_type"`
	Value           float32   `json:"value"`
	AccountId       uuid.UUID `json:"account_id"`
	EntryId         uuid.UUID `json:"entry_id"`
}

type MovementModel struct {
	DB *sql.DB
}

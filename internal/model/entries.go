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
}

type EntryModel struct {
	DB *sql.DB
}

func (m *EntryModel) Get(entryId uuid.UUID) (Entry, error) {
	return Entry{}, nil
}

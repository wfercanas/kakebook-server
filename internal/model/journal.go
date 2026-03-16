package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Journal []Entry

type Entry struct {
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	Amount      float64    `json:"amount"`
	ProjectId   uuid.UUID  `json:"project_id"`
	EntryId     int        `json:"entry_id"`
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

type JournalModel struct {
	DB *sql.DB
}

func (m *JournalModel) GetJournalByProjectId(projectId uuid.UUID) (Journal, error) {
	entriesStmt := `SELECT description, date, amount, entry_id, project_id
	FROM entries
	WHERE project_id = $1
	ORDER BY entry_id DESC`

	movementsStmt := `SELECT ac.account_name, ac.account_category, mv.movement_type, mv.value, mv.account_id
	FROM movements mv
	JOIN accounts ac
	ON mv.account_id = ac.account_id
	WHERE mv.entry_id = $1
	ORDER BY mv.movement_type DESC`

	results, err := m.DB.Query(entriesStmt, projectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Journal{}, ErrNoRecord
		} else {
			return Journal{}, err
		}
	}
	defer results.Close()

	var journal Journal

	for results.Next() {
		var entry Entry

		err := results.Scan(&entry.Description, &entry.Date, &entry.Amount, &entry.EntryId, &entry.ProjectId)
		if err != nil {
			return Journal{}, err
		}

		journal = append(journal, entry)
	}

	for i := range journal {
		results, err := m.DB.Query(movementsStmt, journal[i].EntryId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Journal{}, ErrNoRecord
			} else {
				return Journal{}, err
			}
		}
		defer results.Close()

		for results.Next() {
			var movement Movement
			err = results.Scan(&movement.AccountName, &movement.AccountCategory, &movement.MovementType, &movement.Value, &movement.AccountId)
			if err != nil {
				return Journal{}, err
			}

			journal[i].Movements = append(journal[i].Movements, movement)
		}
	}

	return journal, nil
}

func (m *JournalModel) GetEntryById(entryId int) (Entry, error) {
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
	defer results.Close()

	for results.Next() {
		var movement Movement

		err := results.Scan(&movement.AccountName, &movement.AccountCategory, &movement.MovementType, &movement.Value, &movement.AccountId)
		if err != nil {
			return Entry{}, err
		}

		entry.Movements = append(entry.Movements, movement)
	}

	err = results.Err()
	if err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (m *JournalModel) DeleteEntryById(entryId int) error {
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

func (m *JournalModel) InsertEntry(newEntry NewEntry) error {
	entryStmt := `INSERT INTO entries (date, description, project_id, amount)
	VALUES ($1, $2, $3, $4) RETURNING entry_id`

	movementStmt := `INSERT INTO movements (value, movement_type, account_id, entry_id)
	VALUES ($1, $2, $3, $4)`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var entryId int

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

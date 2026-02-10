package model

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Project struct {
	ProjectId uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
}

type ProjectModel struct {
	DB *sql.DB
}

func (m *ProjectModel) GetProjectsByUserId(userId uuid.UUID) ([]Project, error) {
	stmt := `SELECT p.project_id, title
	FROM projects p
	JOIN projects_users pu
	ON p.project_id = pu.project_id
	WHERE user_id = $1`

	result, err := m.DB.Query(stmt, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	var projects []Project

	for result.Next() {
		var project Project
		err = result.Scan(&project.ProjectId, &project.Title)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (m *ProjectModel) GetAccountsByProjectId(projectId uuid.UUID) ([]Account, error) {
	stmt := `SELECT account_id, account_name, account_category, project_id
	FROM accounts
	WHERE project_id = $1
	ORDER BY account_category`

	results, err := m.DB.Query(stmt, projectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	var accounts []Account

	for results.Next() {
		var account Account

		err = results.Scan(&account.AccountId, &account.AccountName, &account.AccountCategory, &account.ProjectId)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

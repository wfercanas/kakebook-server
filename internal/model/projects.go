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

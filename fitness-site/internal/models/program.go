package models

import (
	"context"
	"database/sql"
	"errors"
)

// Program описывает фитнес-программу.
type Program struct {
	ID          int
	Name       string
	Description string
    Days        []string
}

var ErrProgramNotFound = errors.New("program not found")

// GetAllPrograms возвращает все программы.
func GetAllPrograms(ctx context.Context) ([]Program, error) {
	rows, err := DB.QueryContext(ctx,
		`SELECT id, title, description FROM programs ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progs []Program
	for rows.Next() {
		var p Program
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		progs = append(progs, p)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return progs, nil
}

// GetProgramsForUser просто возвращает GetAllPrograms (можно расширить логику).
func GetProgramsForUser(ctx context.Context, userID int) ([]Program, error) {
	return GetAllPrograms(ctx)
}

// GetProgramByID возвращает программу по ID или ErrProgramNotFound.
func GetProgramByID(ctx context.Context, programID int) (*Program, error) {
	row := DB.QueryRowContext(ctx,
		`SELECT id, title, description FROM programs WHERE id = $1`, programID)

	var p Program
	if err := row.Scan(&p.ID, &p.Name, &p.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProgramNotFound
		}
		return nil, err
	}
	return &p, nil
}

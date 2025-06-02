package models

import (
	"context"
	"database/sql"
	"time"
    "errors"
)

// Workout описывает одну тренировку (день) внутри программы.
type Workout struct {
	ID         int
	ProgramID  int
	DayNumber  int
	Date       time.Time
	Exercises  string
	Notes      string
    Description string
    Duration    int
}

var ErrWorkoutNotFound = errors.New("workout not found")

// GetWorkoutsByProgram возвращает все тренировки для указанной программы.
func GetWorkoutsByProgram(ctx context.Context, programID int) ([]Workout, error) {
	rows, err := DB.QueryContext(ctx,
		`SELECT id, program_id, day_number, date, exercises, notes
         FROM workouts
         WHERE program_id = $1
         ORDER BY day_number ASC`,
		programID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Workout
	for rows.Next() {
		var w Workout
		if err := rows.Scan(&w.ID, &w.ProgramID, &w.DayNumber, &w.Date, &w.Exercises, &w.Notes); err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return list, nil
}

// GetWorkoutByID возвращает одну тренировку по её ID.
func GetWorkoutByID(ctx context.Context, workoutID int) (*Workout, error) {
	row := DB.QueryRowContext(ctx,
		`SELECT id, program_id, day_number, date, exercises, notes FROM workouts WHERE id = $1`, workoutID)

	var w Workout
	if err := row.Scan(&w.ID, &w.ProgramID, &w.DayNumber, &w.Date, &w.Exercises, &w.Notes); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWorkoutNotFound
		}
		return nil, err
	}
	return &w, nil
}

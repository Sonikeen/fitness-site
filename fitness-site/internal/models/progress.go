package models

import (
	"context"
	"time"
    "errors"
)


// Progress хранит факт выполнения дня программ пользователем.
type Progress struct {
	ID          int
	UserID      int
	ProgramID   int
	DayNumber   int
	CompletedAt time.Time
}

var ErrProgressNotFound = errors.New("progress not found")

// MarkDayCompleted создаёт запись о том, что userID выполнил dayNumber в programID.
func MarkDayCompleted(ctx context.Context, userID, programID, dayNumber int) error {
	_, err := DB.ExecContext(ctx,
		`INSERT INTO progress (user_id, program_id, day_number, completed_at)
         VALUES ($1, $2, $3, $4)`,
		userID, programID, dayNumber, time.Now(),
	)
	return err
}

// ListProgress возвращает все записи прогресса для userID и programID.
func ListProgress(ctx context.Context, userID, programID int) ([]Progress, error) {
	rows, err := DB.QueryContext(ctx,
		`SELECT id, user_id, program_id, day_number, completed_at
         FROM progress
         WHERE user_id = $1 AND program_id = $2
         ORDER BY day_number ASC`,
		userID, programID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Progress
	for rows.Next() {
		var p Progress
		if err := rows.Scan(&p.ID, &p.UserID, &p.ProgramID, &p.DayNumber, &p.CompletedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return out, nil
}

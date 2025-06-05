package storage

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "fitness-site/internal/models"
)

type ProgressPGStorage struct {
    pool *pgxpool.Pool
}

type ProgressStorage interface {
    Create(ctx context.Context, userID, programID, day int) error
    Delete(ctx context.Context, userID, programID, day int) error
    List(ctx context.Context, userID, programID int) ([]models.Progress, error)
}

func NewProgressStorage(pool *pgxpool.Pool) *ProgressPGStorage {
    return &ProgressPGStorage{pool: pool}
}

func (s *ProgressPGStorage) Create(ctx context.Context, userID, programID, day int) error {
    _, err := s.pool.Exec(
        ctx,
        "INSERT INTO progress(user_id, program_id, day, completed_at) VALUES ($1, $2, $3, NOW()) ON CONFLICT DO NOTHING",
        userID, programID, day,
    )
    return err
}

func (s *ProgressPGStorage) Delete(ctx context.Context, userID, programID, day int) error {
    _, err := s.pool.Exec(
        ctx,
        "DELETE FROM progress WHERE user_id = $1 AND program_id = $2 AND day = $3",
        userID, programID, day,
    )
    return err
}

func (s *ProgressPGStorage) List(ctx context.Context, userID, programID int) ([]models.Progress, error) {
    rows, err := s.pool.Query(
        ctx,
        "SELECT id, user_id, program_id, day, completed_at FROM progress WHERE user_id=$1 AND program_id=$2",
        userID, programID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []models.Progress
    for rows.Next() {
        var p models.Progress
        if err := rows.Scan(&p.ID, &p.UserID, &p.ProgramID, &p.DayNumber, &p.CompletedAt); err != nil {
            return nil, err
        }
        res = append(res, p)
    }
    return res, nil
}

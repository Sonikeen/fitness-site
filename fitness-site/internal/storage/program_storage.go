package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"fitness-site/internal/models"
)

// ProgramPGStorage работает через pgxpool.Pool
type ProgramPGStorage struct {
	pool *pgxpool.Pool
}

func NewProgramStorage(pool *pgxpool.Pool) *ProgramPGStorage {
	return &ProgramPGStorage{pool: pool}
}

// GetAllPrograms возвращает все программы из БД (колонка name вместо title)
func (s *ProgramPGStorage) GetAllPrograms(ctx context.Context) ([]models.Program, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, description FROM programs ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Program
	for rows.Next() {
		var p models.Program
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// GetByID возвращает одну программу по её ID (с name вместо title)
func (s *ProgramPGStorage) GetByID(ctx context.Context, id int) (*models.Program, error) {
	row := s.pool.QueryRow(ctx, "SELECT id, name, description FROM programs WHERE id = $1", id)
	var p models.Program
	if err := row.Scan(&p.ID, &p.Name, &p.Description); err != nil {
		return nil, err
	}
	return &p, nil
}

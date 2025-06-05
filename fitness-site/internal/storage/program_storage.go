// internal/storage/program_storage.go
package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"fitness-site/internal/models"
)

// ProgramStorage описывает операции с таблицей programs и их днями (days).
type ProgramStorage interface {
	GetAllPrograms(ctx context.Context) ([]models.Program, error)
	GetByID(ctx context.Context, id int) (*models.Program, error)
}

// ProgramPGStorage — реализация ProgramStorage через pgxpool.Pool.
type ProgramPGStorage struct {
	pool *pgxpool.Pool
}

// NewProgramStorage создаёт новое хранилище программ.
func NewProgramStorage(pool *pgxpool.Pool) *ProgramPGStorage {
	return &ProgramPGStorage{pool: pool}
}

// GetAllPrograms возвращает все программы (без деталей по дням).
func (s *ProgramPGStorage) GetAllPrograms(ctx context.Context) ([]models.Program, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, description
         FROM programs
         ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("GetAllPrograms: %w", err)
	}
	defer rows.Close()

	var result []models.Program
	for rows.Next() {
		var p models.Program
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, fmt.Errorf("GetAllPrograms.Scan: %w", err)
		}
		// Поле Days оставляем пустым — подробности подтянем только при GetByID
		result = append(result, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("GetAllPrograms.Rows.Err: %w", rows.Err())
	}
	return result, nil
}

// GetByID возвращает одну программу по её ID, включая массив строк DayDescriptions.
// Использует таблицу `days` для получения описаний: (program_id, day_number, description).
func (s *ProgramPGStorage) GetByID(ctx context.Context, id int) (*models.Program, error) {
	// 1) Сначала читаем базовые поля программы: id, name, description.
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, description
         FROM programs
         WHERE id = $1`, id)

	var p models.Program
	if err := row.Scan(&p.ID, &p.Name, &p.Description); err != nil {
		return nil, fmt.Errorf("GetByID.Scan: %w", err)
	}

	// 2) Теперь читаем все строки из таблицы `days` для этой программы:
	dayRows, err := s.pool.Query(ctx,
		`SELECT day_number, description
         FROM days
         WHERE program_id = $1
         ORDER BY day_number`, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID.DaysQuery: %w", err)
	}
	defer dayRows.Close()

	var days []string
	for dayRows.Next() {
		var dayNum int
		var desc string
		if err := dayRows.Scan(&dayNum, &desc); err != nil {
			return nil, fmt.Errorf("GetByID.DaysScan: %w", err)
		}
		days = append(days, desc)
	}
	if dayRows.Err() != nil {
		return nil, fmt.Errorf("GetByID.DaysRows.Err: %w", dayRows.Err())
	}

	p.Days = days
	return &p, nil
}

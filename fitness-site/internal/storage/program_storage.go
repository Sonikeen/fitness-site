package storage
import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"fitness-site/internal/models"
)
type ProgramStorage interface {
	GetAllPrograms(ctx context.Context) ([]models.Program, error)
	GetByID(ctx context.Context, id int) (*models.Program, error)
	Create(ctx context.Context, p *models.Program) error
	Update(ctx context.Context, p *models.Program) error
	Delete(ctx context.Context, id int) error
}
type ProgramPGStorage struct {
	pool *pgxpool.Pool
}
func NewProgramStorage(pool *pgxpool.Pool) *ProgramPGStorage {
	return &ProgramPGStorage{pool: pool}
}
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
		result = append(result, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("GetAllPrograms.Rows.Err: %w", rows.Err())
	}
	return result, nil
}
func (s *ProgramPGStorage) GetByID(ctx context.Context, id int) (*models.Program, error) {	row := s.pool.QueryRow(ctx,
		`SELECT id, name, description
         FROM programs
         WHERE id = $1`, id)
	var p models.Program
	if err := row.Scan(&p.ID, &p.Name, &p.Description); err != nil {
		return nil, fmt.Errorf("GetByID.Scan: %w", err)
	}
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
func (s *ProgramPGStorage) Create(ctx context.Context, p *models.Program) error {
	var newID int
	err := s.pool.QueryRow(ctx,
		`INSERT INTO programs (name, description)
         VALUES ($1, $2)
         RETURNING id`, p.Name, p.Description).Scan(&newID)
	if err != nil {
		return fmt.Errorf("Create.ProgramInsert: %w", err)
	}
	for i, desc := range p.Days {
		dayNum := i + 1
		_, err := s.pool.Exec(ctx,
			`INSERT INTO days (program_id, day_number, description)
             VALUES ($1, $2, $3)`,
			newID, dayNum, desc,
		)
		if err != nil {
			return fmt.Errorf("Create.DaysInsert: %w", err)
		}
	}
	return nil
}
func (s *ProgramPGStorage) Update(ctx context.Context, p *models.Program) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE programs
         SET name = $1,
             description = $2
         WHERE id = $3`,
		p.Name, p.Description, p.ID,
	)
	if err != nil {
		return fmt.Errorf("Update.ProgramUpdate: %w", err)
	}
	_, err = s.pool.Exec(ctx,
		`DELETE FROM days
         WHERE program_id = $1`, p.ID)
	if err != nil {
		return fmt.Errorf("Update.DaysDelete: %w", err)
	}
	for i, desc := range p.Days {
		dayNum := i + 1
		_, err := s.pool.Exec(ctx,
			`INSERT INTO days (program_id, day_number, description)
             VALUES ($1, $2, $3)`,
			p.ID, dayNum, desc,
		)
		if err != nil {
			return fmt.Errorf("Update.DaysInsert: %w", err)
		}
	}
	return nil
}
func (s *ProgramPGStorage) Delete(ctx context.Context, id int) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM days
         WHERE program_id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete.Days: %w", err)
	}
	_, err = s.pool.Exec(ctx,
		`DELETE FROM programs
         WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete.Programs: %w", err)
	}
	return nil
}
// internal/service/program.go
package service

import (
	"context"
	"fitness-site/internal/models"
)

type ProgramStorage interface {
	GetAllPrograms(ctx context.Context) ([]models.Program, error)
	GetByID(ctx context.Context, id int) (*models.Program, error)
	Create(ctx context.Context, p *models.Program) error
	Update(ctx context.Context, p *models.Program) error
	Delete(ctx context.Context, id int) error
}

type ProgramService struct {
	store ProgramStorage
}

func NewProgramService(store ProgramStorage) *ProgramService {
	return &ProgramService{store: store}
}

func (s *ProgramService) GetAllPrograms(ctx context.Context) ([]models.Program, error) {
	return s.store.GetAllPrograms(ctx)
}

func (s *ProgramService) GetProgramByID(ctx context.Context, id int) (*models.Program, error) {
	return s.store.GetByID(ctx, id)
}

// Create создаёт новую программу
func (s *ProgramService) Create(ctx context.Context, p *models.Program) error {
	return s.store.Create(ctx, p)
}

// Update обновляет существующую программу
func (s *ProgramService) Update(ctx context.Context, p *models.Program) error {
	return s.store.Update(ctx, p)
}

// Delete удаляет программу по ID
func (s *ProgramService) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}

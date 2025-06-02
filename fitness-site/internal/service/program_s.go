package service

import (
    "context"
    "fitness-site/internal/models"
)

type ProgramStorage interface {
    GetAllPrograms(ctx context.Context) ([]models.Program, error)
    GetByID(ctx context.Context, id int) (*models.Program, error)
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

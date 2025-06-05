package service

import (
    "context"
    "fitness-site/internal/models"
)

// Интерфейс хранилища прогресса
type ProgressStorage interface {
    Create(ctx context.Context, userID, programID, day int) error
    Delete(ctx context.Context, userID, programID, day int) error
    List(ctx context.Context, userID, programID int) ([]models.Progress, error)
}

// ProgressService — бизнес-логика работы с прогрессом
type ProgressService struct {
    store ProgressStorage
}

// Конструктор
func NewProgressService(store ProgressStorage) *ProgressService {
    return &ProgressService{store: store}
}

func (s *ProgressService) MarkCompleted(ctx context.Context, userID, programID, day int) error {
    return s.store.Create(ctx, userID, programID, day)
}

func (s *ProgressService) MarkIncomplete(ctx context.Context, userID, programID, day int) error {
    return s.store.Delete(ctx, userID, programID, day)
}

func (s *ProgressService) ListProgress(ctx context.Context, userID, programID int) ([]models.Progress, error) {
    return s.store.List(ctx, userID, programID)
}

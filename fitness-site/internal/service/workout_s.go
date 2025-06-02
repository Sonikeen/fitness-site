package service

import (
	"context"

	"fitness-site/internal/models"
	"fitness-site/internal/storage"
)

// WorkoutService — бизнес-логика для работы с тренировками.
type WorkoutService struct {
	store storage.WorkoutStorage
}

// NewWorkoutService конструктор.
func NewWorkoutService(store storage.WorkoutStorage) *WorkoutService {
	return &WorkoutService{store: store}
}

// GetByProgramID возвращает список тренировок для указанной программы.
func (s *WorkoutService) GetByProgramID(ctx context.Context, programID int) ([]models.Workout, error) {
	return s.store.GetByProgramID(ctx, programID)
}

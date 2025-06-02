package storage

import (
	"context"
	"fitness-site/internal/models"
)

// WorkoutStorage описывает методы для работы с тренировками.
type WorkoutStorage interface {
	GetByProgramID(ctx context.Context, programID int) ([]models.Workout, error)
	GetByID(ctx context.Context, workoutID int) (*models.Workout, error)
}

// WorkoutPGStorage — реальная реализация WorkoutStorage через модели.
type WorkoutPGStorage struct{}

// NewWorkoutStorage возвращает экземпляр реализацией WorkoutStorage.
func NewWorkoutStorage() *WorkoutPGStorage {
	return &WorkoutPGStorage{}
}

// GetByProgramID вызывает models.GetWorkoutsByProgram и возвращает результат.
func (s *WorkoutPGStorage) GetByProgramID(ctx context.Context, programID int) ([]models.Workout, error) {
	return models.GetWorkoutsByProgram(ctx, programID)
}

// GetByID вызывает models.GetWorkoutByID и возвращает результат.
func (s *WorkoutPGStorage) GetByID(ctx context.Context, workoutID int) (*models.Workout, error) {
	return models.GetWorkoutByID(ctx, workoutID)
}

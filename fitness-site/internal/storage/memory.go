package storage

import (
	"sync"

	"fitness-site/internal/models"
)

var (
	// Вспомогательное in-memory хранилище тренировок (для тестов/демо).
	workoutsInMemory = []models.Workout{}
	nextWorkoutID    = 1
	workoutInMemMux  sync.Mutex
)

// GetAllWorkouts возвращает копию всех тренировок из памяти.
func GetAllWorkouts() []models.Workout {
	workoutInMemMux.Lock()
	defer workoutInMemMux.Unlock()

	copied := make([]models.Workout, len(workoutsInMemory))
	copy(copied, workoutsInMemory)
	return copied
}

// GetWorkoutByID ищет тренировку по ID и возвращает копию или nil.
func GetWorkoutByID(id int) *models.Workout {
	workoutInMemMux.Lock()
	defer workoutInMemMux.Unlock()

	for _, w := range workoutsInMemory {
		if w.ID == id {
			copyW := w
			return &copyW
		}
	}
	return nil
}

// AddWorkout присваивает уникальный ID и добавляет тренировку в память.
func AddWorkout(w models.Workout) {
	workoutInMemMux.Lock()
	defer workoutInMemMux.Unlock()

	w.ID = nextWorkoutID
	nextWorkoutID++
	workoutsInMemory = append(workoutsInMemory, w)
}

// DeleteWorkout удаляет из памяти тренировку с заданным ID.
func DeleteWorkout(id int) {
	workoutInMemMux.Lock()
	defer workoutInMemMux.Unlock()

	for i, w := range workoutsInMemory {
		if w.ID == id {
			workoutsInMemory = append(workoutsInMemory[:i], workoutsInMemory[i+1:]...)
			return
		}
	}
}

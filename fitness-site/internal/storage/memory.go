package storage

import (
	"fitness-site/internal/models"
	"sync"
)

var (
	workouts   = []model.Workout{}
	nextID     = 1
	workoutMux sync.Mutex
)

func GetAllWorkouts() []model.Workout {
	return workouts
}

func GetWorkoutByID(id int) *model.Workout {
	for _, w := range workouts {
		if w.ID == id {
			return &w
		}
	}
	return nil
}

func AddWorkout(w model.Workout) {
	workoutMux.Lock()
	defer workoutMux.Unlock()
	w.ID = nextID
	nextID++
	workouts = append(workouts, w)
}

func DeleteWorkout(id int) {
	workoutMux.Lock()
	defer workoutMux.Unlock()
	for i, w := range workouts {
		if w.ID == id {
			workouts = append(workouts[:i], workouts[i+1:]...)
			return
		}
	}
}

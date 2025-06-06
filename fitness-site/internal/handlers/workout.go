package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"fitness-site/internal/models"
	"fitness-site/internal/storage"
)

// WorkoutListHandler показывает HTML-список всех тренировок (in-memory).
func WorkoutListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Список тренировок</h1>")
	for _, workout := range storage.GetAllWorkouts() {
		fmt.Fprintf(w, "<p>ID: %d | Название: %s | Длительность: %d мин</p>",
			workout.ID, workout.Description, workout.Duration)
	}
	fmt.Fprint(w, `<br><a href="/workouts/new">Добавить тренировку</a>`)
}

// WorkoutCreateHandler обрабатывает GET/POST для добавления тренировки (in-memory).
func WorkoutCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `
			<h1>Добавить тренировку</h1>
			<form method="POST">
				<label>Название:</label><br>
				<input type="text" name="title" required><br><br>
				<label>Длительность (в минутах):</label><br>
				<input type="number" name="duration" required><br><br>
				<input type="submit" value="Сохранить">
			</form>
			<a href="/workouts">← Назад к списку</a>
		`)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}
		title := r.FormValue("title")
		durationStr := r.FormValue("duration")
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			http.Error(w, "Длительность должна быть числом", http.StatusBadRequest)
			return
		}
		newWorkout := models.Workout{
			Description: title,
			Duration:    duration,
		}
		storage.AddWorkout(newWorkout)
		http.Redirect(w, r, "/workouts", http.StatusSeeOther)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"fitness-site/internal/middleware"
	"fitness-site/internal/models"
)

// programSummaryJSON используется для JSON-ответа списка программ.
type programSummaryJSON struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// workoutInfoJSON описывает одну тренировку для JSON.
type workoutInfoJSON struct {
	DayNumber int    `json:"day_number"`
	Date      string `json:"date"`
	Exercises string `json:"exercises"`
	Notes     string `json:"notes"`
}

// programDetailJSON описывает подробную информацию о программе.
type programDetailJSON struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Workouts    []workoutInfoJSON `json:"workouts"`
}

// DashboardHandler возвращает JSON со списком программ для текущего пользователя.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем все программы для пользователя (или все, если нет разделения)
	programs, err := models.GetProgramsForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Ошибка при загрузке программ", http.StatusInternalServerError)
		return
	}

	var resp []programSummaryJSON
	for _, p := range programs {
		resp = append(resp, programSummaryJSON{
			ID:    p.ID,
			Title: p.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ProgramHandler возвращает JSON с деталями конкретной программы по ID.
func ProgramHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	programID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем программу
	prog, err := models.GetProgramByID(r.Context(), programID)
	if err != nil {
		http.Error(w, "Программа не найдена", http.StatusNotFound)
		return
	}

	// Получаем тренировки
	workouts, err := models.GetWorkoutsByProgram(r.Context(), programID)
	if err != nil {
		http.Error(w, "Ошибка при получении тренировок", http.StatusInternalServerError)
		return
	}

	var resp programDetailJSON
	resp.ID = prog.ID
	resp.Title = prog.Name
	resp.Description = prog.Description

	for _, wkt := range workouts {
		resp.Workouts = append(resp.Workouts, workoutInfoJSON{
			DayNumber: wkt.DayNumber,
			Date:      wkt.Date.Format("2006-01-02"),
			Exercises: wkt.Exercises,
			Notes:     wkt.Notes,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

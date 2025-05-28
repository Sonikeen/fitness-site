package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
	"fitness-site/internal/middleware"
    "github.com/go-chi/chi/v5"

    "fitness-site/internal/models"
)

// TrackProgress обрабатывает POST /programs/{id}/progress
func TrackProgress(w http.ResponseWriter, r *http.Request) {
    pid, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Неверный ID программы", http.StatusBadRequest)
        return
    }
    day, err := strconv.Atoi(r.FormValue("day"))
    if err != nil {
        http.Error(w, "Неверный день", http.StatusBadRequest)
        return
    }
    userID := r.Context().Value(middleware.UserIDKey).(int)


    p := models.Progress{
        UserID:    userID,
        ProgramID: pid,
        Day:       day,
    }
    if err := progressService.MarkCompleted(p); err != nil {
        http.Error(w, "Ошибка сохранения прогресса", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

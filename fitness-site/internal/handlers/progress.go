package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "fitness-site/internal/middleware"
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

    userIDVal := r.Context().Value(middleware.UserIDKey)
    if userIDVal == nil {
        http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
        return
    }
    userID := userIDVal.(int)

    checked := r.FormValue("checked") == "true"

    var progressErr error
    if checked {
        progressErr = ProgressService.MarkCompleted(r.Context(), userID, pid, day)
    } else {
        progressErr = ProgressService.MarkIncomplete(r.Context(), userID, pid, day)
    }

    if progressErr != nil {
        http.Error(w, "Не удалось обновить прогресс", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

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
    // 1) Получаем ID программы из URL
    pid, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Неверный ID программы", http.StatusBadRequest)
        return
    }

    // 2) Парсим номер дня из формы (например, ?day=1)
    day, err := strconv.Atoi(r.FormValue("day"))
    if err != nil {
        http.Error(w, "Неверный день", http.StatusBadRequest)
        return
    }

    // 3) Получаем userID из контекста
    userIDVal := r.Context().Value(middleware.UserIDKey)
    if userIDVal == nil {
        http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
        return
    }
    userID := userIDVal.(int)

    // 4) Сохраняем прогресс через сервис
    if err := ProgressService.MarkCompleted(r.Context(), userID, pid, day); err != nil {
        http.Error(w, "Не удалось сохранить прогресс", http.StatusInternalServerError)
        return
    }

    // 5) Отправляем JSON-ответ
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

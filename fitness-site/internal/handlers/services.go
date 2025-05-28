package handlers

import (
    "net/http"

    "fitness-site/internal/middleware"
    "fitness-site/internal/models"
)

// ServicesPageData — данные для шаблона services.html
type ServicesPageData struct {
    Programs []models.Program
    LoggedIn bool
}

// ServicesPage выводит список программ и кнопку «Начать» или «Войти»
func ServicesPage(w http.ResponseWriter, r *http.Request) {
    progs, err := programService.GetAllPrograms()
    if err != nil {
        http.Error(w, "Ошибка получения программ", http.StatusInternalServerError)
        return
    }

    // Вот здесь мы должны смотреть в middleware.UserIDKey
    _, loggedIn := r.Context().Value(middleware.UserIDKey).(int)

    data := ServicesPageData{
        Programs: progs,
        LoggedIn: loggedIn,
    }
    render(w, "services.html", data)
}

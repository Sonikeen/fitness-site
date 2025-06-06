package handlers

import (
    "log"
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
    progs, err := ProgramService.GetAllPrograms(r.Context())
    if err != nil {
        log.Printf("Ошибка получения программ: %v", err)
        http.Error(w, "Ошибка получения программ: "+err.Error(), http.StatusInternalServerError)
        return
    }

    _, loggedIn := middleware.UserIDFromContext(r.Context())

    data := ServicesPageData{
        Programs: progs,
        LoggedIn: loggedIn,
    }
    render(w, "services.html", data)
}

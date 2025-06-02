package handlers

import (
    "net/http"
    "fitness-site/internal/models"
    "fitness-site/internal/middleware"
    "log"
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
        log.Printf("Ошибка получения программ: %v", err) // <-- будет видно в консоли
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

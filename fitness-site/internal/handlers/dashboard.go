package handlers

import (
    "html/template"
    "log"
    "net/http"
    "path/filepath"

    "fitness-site/internal/middleware"
    "fitness-site/internal/models"
)

// DashboardItem описывает одну строку в личном кабинете.
type DashboardItem struct {
    Program   models.Program
    Completed int
    Total     int
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.UserIDFromContext(r.Context())
    if !ok {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Все программы
    progs, err := ProgramService.GetAllPrograms(r.Context())

    if err != nil {
        http.Error(w, "Ошибка получения программ", http.StatusInternalServerError)
        return
    }

    // Для каждой программы считаем сколько дней выполнено
    var items []DashboardItem
    for _, p := range progs {
        done, err := ProgressService.ListProgress(r.Context(), userID, p.ID)

        if err != nil {
            http.Error(w, "Ошибка загрузки прогресса", http.StatusInternalServerError)
            return
        }
        // p.Days должен быть []string!
        items = append(items, DashboardItem{
            Program:   p,
            Completed: len(done),
            Total:     len(p.Days),
        })
    }

    // Рендерим шаблон
    data := map[string]interface{}{"Items": items}
    t := template.Must(template.ParseFiles(
        filepath.Join("internal", "templates", "base.html"),
        filepath.Join("internal", "templates", "dashboard.html"),
    ))
    if err := t.ExecuteTemplate(w, "base", data); err != nil {
        log.Printf("Dashboard: ошибка при выполнении шаблона: %v", err)
        http.Error(w, "Ошибка серверного рендеринга", http.StatusInternalServerError)
    }
}

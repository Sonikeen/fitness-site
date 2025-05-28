package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
	"fitness-site/internal/middleware"
    "fitness-site/internal/models"
)

// DashboardItem описывает одну строку в личном кабинете
type DashboardItem struct {
    Program   models.Program
    Completed int
    Total     int
}

// Dashboard рендерит личный кабинет с прогрессом по всем программам
func Dashboard(w http.ResponseWriter, r *http.Request) {
    // 1) Получаем ID пользователя из контекста
	userID := r.Context().Value(middleware.UserIDKey).(int)


    // 2) Берём все программы
    progs, err := programService.GetAllPrograms()
    if err != nil {
        http.Error(w, "Ошибка получения программ", http.StatusInternalServerError)
        return
    }

    // 3) Для каждой программы считаем, сколько дней выполнено
    var items []DashboardItem
    for _, p := range progs {
        done, err := progressService.ListProgress(userID, p.ID)
        if err != nil {
            http.Error(w, "Ошибка загрузки прогресса", http.StatusInternalServerError)
            return
        }
        items = append(items, DashboardItem{
            Program:   p,
            Completed: len(done),
            Total:     len(p.Days),
        })
    }

    // 4) Рендерим шаблон
    data := map[string]interface{}{"Items": items}
    t := template.Must(template.ParseFiles(
        filepath.Join("internal", "templates", "base.html"),
        filepath.Join("internal", "templates", "dashboard.html"),
    ))
    t.ExecuteTemplate(w, "base", data)
}

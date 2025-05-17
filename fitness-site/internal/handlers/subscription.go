package handler

import (
    "fitness-site/internal/models"
    "fitness-site/internal/storage"
    "fmt"
    "net/http"
    "strconv"
)

// SubscriptionListHandler — GET /subscriptions
func SubscriptionListHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintln(w, "<h1>Список подписок</h1>")
    for _, sub := range storage.GetAllSubscriptions() {
        fmt.Fprintf(w,
            "<p>ID: %d | Пользователь: %d | План: %s | Активна: %t</p>",
            sub.ID, sub.UserID, sub.Plan, sub.IsActive,
        )
    }
    fmt.Fprint(w, `<br><a href="/subscriptions/new">Добавить подписку</a>`)
}

// SubscriptionCreateHandler — GET & POST /subscriptions/new
func SubscriptionCreateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprint(w, `
            <h1>Добавить подписку</h1>
            <form method="POST" action="/subscriptions/new">
                <label>ID пользователя:</label><br>
                <input type="number" name="user_id" required><br><br>
                <label>План:</label><br>
                <input type="text" name="plan" required><br><br>
                <label>Активна:</label>
                <input type="checkbox" name="active"><br><br>
                <button type="submit">Сохранить</button>
            </form>
            <a href="/subscriptions">← Назад к списку</a>
        `)
        return
    }

    // POST
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Ошибка формы", http.StatusBadRequest)
        return
    }

    userID, err := strconv.Atoi(r.FormValue("user_id"))
    if err != nil {
        http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
        return
    }

    isActive := r.FormValue("active") == "on"

    newSub := model.Subscription{
        UserID:   userID,
        Plan:     r.FormValue("plan"),
        IsActive: isActive,
    }

    storage.AddSubscription(newSub)
    http.Redirect(w, r, "/subscriptions", http.StatusSeeOther)
}

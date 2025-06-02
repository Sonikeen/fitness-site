package handlers

import (
    "net/http"
    "strings"

    "github.com/jackc/pgx/v5/pgconn"
    "fitness-site/internal/models"
    "fitness-site/internal/middleware"
)

// ShowRegister отображает форму регистрации.
func ShowRegister(w http.ResponseWriter, r *http.Request) {
    render(w, "register.html", map[string]interface{}{"Error": ""})
}

// HandleRegister обрабатывает POST /register.
func HandleRegister(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        render(w, "register.html", map[string]interface{}{"Error": "Неверные данные формы"})
        return
    }

    username := strings.TrimSpace(r.FormValue("username"))
    email := strings.TrimSpace(r.FormValue("email"))
    pass := strings.TrimSpace(r.FormValue("password"))

    // Простая валидация на пустые поля
    if username == "" || email == "" || pass == "" {
        render(w, "register.html", map[string]interface{}{
            "Error": "Пожалуйста, заполните имя пользователя, email и пароль",
        })
        return
    }

    u := &models.User{
        Name:     username,
        Email:    email,
        PasswordHash: pass,
    }

    // Пытаемся зарегистрировать через сервис
    if err := UserService.Register(r.Context(), u); err != nil {
        // Пытаемся распознать код PG 23505 — уникальное ограничение нарушено
        if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
            // constraint может быть "users_username_key" или "users_email_key"
            if strings.Contains(pgErr.ConstraintName, "users_username_key") {
                render(w, "register.html", map[string]interface{}{"Error": "Имя пользователя уже занято"})
                return
            }
            if strings.Contains(pgErr.ConstraintName, "users_email_key") {
                render(w, "register.html", map[string]interface{}{"Error": "Email уже зарегистрирован"})
                return
            }
        }
        // Любая другая ошибка
        render(w, "register.html", map[string]interface{}{
            "Error": "Ошибка регистрации: " + err.Error(),
        })
        return
    }

    // Успешно — перенаправляем на страницу входа
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ShowLogin отображает форму входа.
func ShowLogin(w http.ResponseWriter, r *http.Request) {
    render(w, "login.html", map[string]interface{}{"Error": ""})
}

// HandleLogin обрабатывает POST /login.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        render(w, "login.html", map[string]interface{}{"Error": "Неверные данные"})
        return
    }
    email := strings.TrimSpace(r.FormValue("email"))
    pass := strings.TrimSpace(r.FormValue("password"))

    user, err := UserService.Authenticate(r.Context(), email, pass)
    if err != nil {
        render(w, "login.html", map[string]interface{}{"Error": "Неверный email или пароль"})
        return
    }

    sid, err := middleware.CreateSession(user.ID)
    if err != nil {
        render(w, "login.html", map[string]interface{}{"Error": "Ошибка создания сессии"})
        return
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sid,
        Path:     "/",
        HttpOnly: true,
    })

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// LogoutHandler разрывает сессию и перенаправляет на /login.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        sid := cookie.Value
        middleware.DeleteSession(sid)
        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    "",
            Path:     "/",
            HttpOnly: true,
            MaxAge:   -1,
        })
    }
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

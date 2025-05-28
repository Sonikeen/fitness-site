package handlers

import (
    "fmt"
    "net/http"
    "strings"
    "github.com/jackc/pgx/v5/pgconn"
    "fitness-site/internal/middleware"
    "fitness-site/db"
    "fitness-site/internal/models"
    "fitness-site/internal/service"
    "fitness-site/internal/storage"
)

var userSvc *service.UserService

func init() {
    db.Connect()
    us := storage.NewUserStorage(db.Conn)
    userSvc = service.NewUserService(us)
}

func ShowRegister(w http.ResponseWriter, r *http.Request) {
    render(w, "register.html", map[string]interface{}{"Error": ""})
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        render(w, "register.html", map[string]interface{}{"Error": "Неверные данные формы"})
        return
    }

    // Теперь читаем из поля "username"
    username := strings.TrimSpace(r.FormValue("username"))
    email := strings.TrimSpace(r.FormValue("email"))
    pass := strings.TrimSpace(r.FormValue("password"))

    fmt.Printf("DEBUG Register: username=%q email=%q password=%q\n", username, email, pass)

    if username == "" || email == "" || pass == "" {
        render(w, "register.html", map[string]interface{}{
            "Error": "Пожалуйста, заполните имя пользователя, email и пароль",
        })
        return
    }

    u := &models.User{
        Name:     username,  // модельное поле Name сохраняется в колонку username
        Email:    email,
        Password: pass,
    }
    if err := userSvc.Register(u); err != nil {
        if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
            if strings.Contains(pgErr.ConstraintName, "username") {
                render(w, "register.html", map[string]interface{}{"Error": "Имя пользователя уже занято"})
                return
            }
            if strings.Contains(pgErr.ConstraintName, "email") {
                render(w, "register.html", map[string]interface{}{"Error": "Email уже зарегистрирован"})
                return
            }
        }
        render(w, "register.html", map[string]interface{}{
            "Error": "Ошибка регистрации: " + err.Error(),
        })
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
    render(w, "login.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Неверные данные", http.StatusBadRequest)
        return
    }
    email := strings.TrimSpace(r.FormValue("email"))
    pass := strings.TrimSpace(r.FormValue("password"))

    user, err := userSvc.Authenticate(email, pass)
    if err != nil {
        http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
        return
    }
    // вместо контекста — создаём сессию
    sid, err := middleware.CreateSession(user.ID)
    if err != nil {
        http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
        return
    }
    // ставим HTTP-cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sid,
        Path:     "/",
        HttpOnly: true,
        // Secure: true, // в проде по https
    })
    http.Redirect(w, r, "/services", http.StatusSeeOther)
}

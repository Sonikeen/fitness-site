package handlers
import (
    "context"
    "html/template"
    "net/http"
    "path/filepath"
    "fitness-site/db"
    "golang.org/x/crypto/bcrypt"
)
var registerTmpl = template.Must(template.ParseFiles(
    filepath.Join("internal", "templates", "base.html"),
    filepath.Join("internal", "templates", "contact.html"),
))
func RegisterHandlers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := registerTmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
            http.Error(w, "Ошибка рендеринга: "+err.Error(), http.StatusInternalServerError)
        }
    case http.MethodPost:
        username := r.FormValue("username")
        email    := r.FormValue("email")
        password := r.FormValue("password")
        confirm  := r.FormValue("confirm")
        if username == "" || email == "" || password == "" || confirm == "" {
            http.Error(w, "Все поля обязательны", http.StatusBadRequest)
            return
        }
        if password != confirm {
            http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
            return
        }
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Ошибка хеширования", http.StatusInternalServerError)
            return
        }
        if _, err := db.Conn.Exec(
            context.Background(),
            `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`,
            username, email, string(hashed),
        ); err != nil {
            http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

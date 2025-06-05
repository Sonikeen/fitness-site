package handlers
import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"fitness-site/internal/middleware"
	"fitness-site/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)
func loadAuthTemplates(child string) (*template.Template, error) {
	basePath := filepath.Join("internal", "templates", "base.html")
	childPath := filepath.Join("internal", "templates", child)
	return template.ParseFiles(basePath, childPath)
}
func renderAuthError(w http.ResponseWriter, activeTab, message string) {
	tmpl, err := loadAuthTemplates("dashboard.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Error":     message,
		"ActiveTab": activeTab,
		"Items":     nil,
	}
	tmpl.ExecuteTemplate(w, "base", data)
}
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderAuthError(w, "register", "Неверные данные формы")
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	pass := strings.TrimSpace(r.FormValue("password"))
	confirm := strings.TrimSpace(r.FormValue("confirm"))
	if username == "" || email == "" || pass == "" || confirm == "" || pass != confirm {
		renderAuthError(w, "register", "Пожалуйста, заполните все поля корректно и подтвердите пароль")
		return
	}
	u := &models.User{
		Name:         username,
		Email:        email,
		PasswordHash: pass,
	}
	if err := UserService.Register(r.Context(), u); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "users_username_key") {
				renderAuthError(w, "register", "Имя пользователя уже занято")
				return
			}
			if strings.Contains(pgErr.ConstraintName, "users_email_key") {
				renderAuthError(w, "register", "Email уже зарегистрирован")
				return
			}
		}
		renderAuthError(w, "register", "Ошибка регистрации: "+err.Error())
		return
	}
	user, err := UserService.Authenticate(r.Context(), email, pass)
	if err != nil {
		http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
		return
	}
	sid, err := middleware.CreateSession(user.ID)
	if err != nil {
		renderAuthError(w, "register", "Ошибка создания сессии")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/programs", http.StatusSeeOther)

}
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderAuthError(w, "login", "Неверные данные")
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	pass := strings.TrimSpace(r.FormValue("password"))
	fmt.Printf("Попытка логина: email=%q\n", email)
	user, err := UserService.Authenticate(r.Context(), email, pass)
	if err != nil {
		fmt.Printf("Authenticate вернул ошибку: %v\n", err)
		renderAuthError(w, "login", "Неверный email или пароль")
		return
	}
	sid, err := middleware.CreateSession(user.ID)
	if err != nil {
		renderAuthError(w, "login", "Ошибка создания сессии")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, 
	})
	http.Redirect(w, r, "/programs", http.StatusSeeOther)

}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		middleware.DeleteSession(cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}
	http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
}

package handlers
import (
	"fitness-site/internal/models"
	"html/template"
	"log"
	"net/http"
)
var users []models.User
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		user := models.User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"), // В реальном проекте нужно хешировать!
		}
		users = append(users, user)
		log.Printf("Зарегистрирован пользователь: %s (%s)\n", user.Name, user.Email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	renderAuthPage(w, "register")
}
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		for _, u := range users {
			if u.Email == email && u.Password == password {
				log.Printf("Успешный вход: %s", email)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
		log.Println("Ошибка входа: неверные данные")
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}
	renderAuthPage(w, "login")
}
func renderAuthPage(w http.ResponseWriter, page string) {
	tmpl, err := template.ParseFiles(
		"internal/templates/base.html",
		"internal/templates/" + page + ".html",
	)
	if err != nil {
		log.Println("Ошибка при парсинге шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

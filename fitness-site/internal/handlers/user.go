package handlers
import (
	"fitness-site/internal/models"
	"html/template"
	"log"
	"net/http"
)
var users []models.User


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

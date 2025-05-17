package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(
		"internal/templates/base.html",
		"internal/templates/"+tmpl+".html",
	)
	if err != nil {
		log.Println("Ошибка шаблона:", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about")
}

func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "services")
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		message := r.FormValue("message")
		log.Printf("Новое сообщение от %s: %s\n", name, message)

		fmt.Fprintf(w, "<h1>Спасибо, %s!</h1><p>Ваше сообщение получено.</p><a href='/'>На главную</a>", name)
		return
	}
	renderTemplate(w, "contact")
}

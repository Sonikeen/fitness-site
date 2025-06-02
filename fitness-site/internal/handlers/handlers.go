package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
    "log"
)

func render(w http.ResponseWriter, tmpl string, data interface{}) {
    t := template.Must(template.ParseFiles(
        filepath.Join("internal", "templates", "base.html"),
        filepath.Join("internal", "templates", tmpl),
    ))
    if err := t.ExecuteTemplate(w, "base", data); err != nil {
        log.Printf("render: ошибка при выполнении шаблона: %v", err)
        http.Error(w, "Ошибка серверного рендеринга", http.StatusInternalServerError)
    }
}

func HomePage(w http.ResponseWriter, r *http.Request) {
    render(w, "index.html", map[string]interface{}{"Title": "Главная"})
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
    render(w, "about.html", map[string]interface{}{"Title": "О нас"})
}

func ContactPage(w http.ResponseWriter, r *http.Request) {
    render(w, "contact.html", map[string]interface{}{"Title": "Контакты"})
}

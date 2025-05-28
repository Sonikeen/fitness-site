package handlers

import (
    "net/http"
    "html/template"
    "path/filepath"
	"fmt"
)

// render упрощает рендеринг: базовый шаблон + конкретный tmpl
func render(w http.ResponseWriter, tmpl string, data interface{}) {
    t := template.Must(template.ParseFiles(
        filepath.Join("internal", "templates", "base.html"),
        filepath.Join("internal", "templates", tmpl),
    ))
    t.ExecuteTemplate(w, "base", data)
}

// HomePage — обработчик корня "/"
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== HomePage reached ===") 
    render(w, "index.html", map[string]string{"Title": "Главная"})
}

// AboutPage — обработчик "/about"
func AboutPage(w http.ResponseWriter, r *http.Request) {
    render(w, "about.html", map[string]string{"Title": "О нас"})
}

// ContactPage — обработчик "/contact"
func ContactPage(w http.ResponseWriter, r *http.Request) {
    render(w, "contact.html", map[string]string{"Title": "Контакты"})
}

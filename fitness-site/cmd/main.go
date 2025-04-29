package main

import (
    "fmt"
    "log"
    "net/http"
    "html/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("internal/templates/index.html")
    if err != nil {
        http.Error(w, "Template parsing error", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func main() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    http.HandleFunc("/", homeHandler)

    fmt.Println("Сервер запущен на http://localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}

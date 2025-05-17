package main

import (
    "fmt"
    "log"
    "net/http"

    "fitness-site/db"
    "fitness-site/internal/router"
)

func main() {
    // 1) Подключаемся к БД
    db.Connect()

    // 2) Настраиваем маршруты через ваш router
    mux := router.SetupRoutes()

    fmt.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

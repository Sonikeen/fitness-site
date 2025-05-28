package router

import (
    "net/http"

    chi "github.com/go-chi/chi/v5"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"

    "fitness-site/internal/handlers"
    "fitness-site/internal/middleware"
)

// SetupRouter настраивает все маршруты и возвращает HTTP-роутер
func SetupRouter() http.Handler {
    r := chi.NewRouter()
    r.Use(middleware.SessionMiddleware)

    // Chi встроенные middleware
    r.Use(chiMiddleware.RealIP, chiMiddleware.Logger, chiMiddleware.Recoverer)

    // Статика
    fileServer := http.FileServer(http.Dir("static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

    // Публичные страницы
    r.Get("/", handlers.HomePage)
    r.Get("/services", handlers.ServicesPage)
    r.Get("/about", handlers.AboutPage)
    r.Get("/contact", handlers.ContactPage)

    // Регистрация и вход
    r.Get("/register", handlers.ShowRegister)
    r.Post("/register", handlers.HandleRegister)
    r.Get("/login", handlers.ShowLogin)
    r.Post("/login", handlers.HandleLogin)

    // Защищённые маршруты: список программ и прогресс
    r.Route("/programs", func(r chi.Router) {
        r.Use(middleware.AuthMiddleware)
        r.Get("/", handlers.ProgramList)         // ← теперь определён в program.go
        r.Get("/{id}", handlers.ProgramDetail)
        r.Post("/{id}/progress", handlers.TrackProgress)
    })
    // Личный кабинет
    r.With(middleware.AuthMiddleware).
        Get("/dashboard", handlers.Dashboard)
    
        return r
}

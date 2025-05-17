package router

import (
    "net/http"

    "fitness-site/internal/handlers"
)

func SetupRoutes() http.Handler {
    mux := http.NewServeMux()

    // 1) Раздача статики
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // 2) Основные страницы
    mux.HandleFunc("/",         handler.HomeHandler)
    mux.HandleFunc("/about",    handler.AboutHandler)
    mux.HandleFunc("/services", handler.ServicesHandler)
    mux.HandleFunc("/contact",  handler.ContactHandler)

    // 3) Регистрация (HTML и JSON)
    mux.HandleFunc("/register",     handler.RegisterHandlers)
    mux.HandleFunc("/api/register", handler.APIRegisterHandler)

    // 4) Тренировки
    mux.HandleFunc("/workouts",     handler.WorkoutListHandler)
    mux.HandleFunc("/workouts/new", handler.WorkoutCreateHandler)

    // 5) Подписки
    mux.HandleFunc("/subscriptions",     handler.SubscriptionListHandler)
    mux.HandleFunc("/subscriptions/new", handler.SubscriptionCreateHandler)

    // 6) Остальные ваши обработчики (subscription, user, etc.)
    // e.g. mux.HandleFunc("/user", handlers.UserHandler) ...

    return mux
}

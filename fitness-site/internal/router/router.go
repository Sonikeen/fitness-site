package router

import (
	"net/http"
	"fitness-site/internal/handlers"
	"fitness-site/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RealIP, chiMiddleware.Logger, chiMiddleware.Recoverer)

	// Статика
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Публичные страницы
	r.Get("/", handlers.HomePage)
	r.Get("/services", handlers.ServicesPage)
	r.Get("/about", handlers.AboutPage)

	// Точка входа: Dashboard (покажет либо форму логина/регистрации, либо список программ)
	r.Get("/dashboard", handlers.Dashboard)



	// Логин / регистрация
	r.Post("/login", handlers.HandleLogin)
	r.Post("/register", handlers.HandleRegister)
	r.Get("/logout", handlers.LogoutHandler)

	// Профиль (редактирование) — только для авторизованных
	r.With(middleware.AuthMiddleware).Get("/profile", handlers.ProfileEditHandler)
	r.With(middleware.AuthMiddleware).
Route("/profile/edit", func(r chi.Router) {
		r.Get("/", handlers.ProfileEditHandler)
		r.Post("/", handlers.ProfileEditHandler)
	})

	// программы и прогресс (только для авторизованных)
	r.Route("/programs", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", handlers.Dashboard)             // список программ
		r.Get("/{id}", handlers.ProgramPageHandler) // детали программы
		r.Post("/{id}/progress", handlers.TrackProgress)
	})

	// тренировки (публичные)
	r.Get("/workouts", handlers.WorkoutListHandler)
	r.Get("/workouts/new", handlers.WorkoutCreateHandler)
	r.Post("/workouts/new", handlers.WorkoutCreateHandler)

	return r
}
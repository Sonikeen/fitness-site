package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"fitness-site/internal/handlers"
	"fitness-site/internal/middleware"
)

// SetupRouter настраивает все маршруты и возвращает http.Handler (Chi).
func SetupRouter() http.Handler {
	r := chi.NewRouter()

	// Встроенные middleware Chi
	r.Use(
		chiMiddleware.RealIP,
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
		// Если у тебя есть SessionMiddleware — добавь здесь:
		// middleware.SessionMiddleware,
	)

	// Статика (/static/* → папка static/)
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

	// Выход — только для авторизованных
	r.With(middleware.AuthMiddleware).
		Get("/logout", handlers.LogoutHandler)

	// Личный кабинет (HTML)
	r.With(middleware.AuthMiddleware).
		Get("/dashboard", handlers.Dashboard)

	// Работа с программами (JSON API, только для авторизованных)
	r.Route("/programs", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", handlers.DashboardHandler)           // GET  /programs
		r.Get("/{id}", handlers.ProgramHandler)         // GET  /programs/{id}
		r.Post("/{id}/progress", handlers.TrackProgress) // POST /programs/{id}/progress
	})

	// Тренировки (in-memory)
	r.Get("/workouts", handlers.WorkoutListHandler)
	r.Get("/workouts/new", handlers.WorkoutCreateHandler)
	r.Post("/workouts/new", handlers.WorkoutCreateHandler)

	return r
}

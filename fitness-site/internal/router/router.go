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

	// Авторизация / Личный кабинет
	r.Get("/dashboard", handlers.Dashboard)
	r.Post("/login", handlers.HandleLogin)
	r.Post("/register", handlers.HandleRegister)
	r.Get("/logout", handlers.LogoutHandler)

	// Профиль (только авторизованные)
	r.With(middleware.AuthMiddleware).Get("/profile", handlers.ProfileEditHandler)
	r.With(middleware.AuthMiddleware).Route("/profile/edit", func(r chi.Router) {
		r.Get("/", handlers.ProfileEditHandler)
		r.Post("/", handlers.ProfileEditHandler)
	})

	// Программы и прогресс (только авторизованные)
	r.Route("/programs", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", handlers.Dashboard)
		r.Get("/{id}", handlers.ProgramPageHandler)
		r.Post("/{id}/progress", handlers.TrackProgress)
	})

	// Админка: управление программами
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/programs", handlers.AdminProgramsList)
		r.Get("/programs/new", handlers.AdminNewProgramForm)
		r.Post("/programs/new", handlers.AdminNewProgramSubmit)
		r.Get("/programs/{id}/edit", handlers.AdminEditProgramForm)
		r.Post("/programs/{id}/edit", handlers.AdminEditProgramSubmit)
		r.Get("/programs/{id}/delete", handlers.AdminDeleteProgram)
	})

	return r
}

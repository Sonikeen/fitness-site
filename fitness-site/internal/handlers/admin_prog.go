package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"fitness-site/internal/middleware"
	"fitness-site/internal/models"
)

func AdminProgramsList(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || !user.IsAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	programs, err := ProgramService.GetAllPrograms(r.Context())
	if err != nil {
		http.Error(w, "Ошибка загрузки программ", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "admin_prog.html"),
	)
	if err != nil {
		log.Printf("AdminProgramsList: ошибка ParseFiles: %v", err)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Programs":   programs,
		"IsLoggedIn": true,
		"IsAdmin":    true,
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("AdminProgramsList: ExecuteTemplate вернул ошибку: %v", err)
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}

func AdminNewProgramForm(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || !user.IsAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "admin_program_form.html"),
	)
	if err != nil {
		log.Printf("AdminNewProgramForm: ошибка ParseFiles: %v", err)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"IsLoggedIn": true,
		"IsAdmin":    true,
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("AdminNewProgramForm: ExecuteTemplate вернул ошибку: %v", err)
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}

func AdminNewProgramSubmit(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || !user.IsAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	desc := strings.TrimSpace(r.FormValue("description"))
	daysText := strings.TrimSpace(r.FormValue("days"))
	if name == "" || desc == "" || daysText == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}
	days := strings.Split(daysText, "\n")
	prog := &models.Program{
		Name:        name,
		Description: desc,
		Days:        days,
	}
	if err := ProgramService.Create(r.Context(), prog); err != nil {
		log.Printf("AdminNewProgramSubmit: ошибка сохранения: %v", err)
		http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/programs", http.StatusSeeOther)
}

func AdminEditProgramForm(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || !user.IsAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	prog, err := ProgramService.GetProgramByID(r.Context(), id)
	if err != nil || prog == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "admin_program_edit.html"),
	)
	if err != nil {
		log.Printf("AdminEditProgramForm: ошибка ParseFiles: %v", err)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Program":    prog,
		"IsLoggedIn": true,
		"IsAdmin":    true,
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("AdminEditProgramForm: ExecuteTemplate вернул ошибку: %v", err)
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}

func AdminEditProgramSubmit(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || !user.IsAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	desc := strings.TrimSpace(r.FormValue("description"))
	daysText := strings.TrimSpace(r.FormValue("days"))
	if name == "" || desc == "" || daysText == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}
	days := strings.Split(daysText, "\n")

	prog := &models.Program{
		ID:          id,
		Name:        name,
		Description: desc,
		Days:        days,
	}
	if err := ProgramService.Update(r.Context(), prog); err != nil {
		log.Printf("AdminEditProgramSubmit: ошибка обновления: %v", err)
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/programs", http.StatusSeeOther)
}

func AdminDeleteProgram(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || !user.IsAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := ProgramService.Delete(r.Context(), id); err != nil {
		log.Printf("AdminDeleteProgram: ошибка удаления: %v", err)
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/programs", http.StatusSeeOther)
}

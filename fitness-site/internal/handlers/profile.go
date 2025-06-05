package handlers

import (
    "net/http"
    "html/template"
    "path/filepath"
    "strings"
    "fitness-site/internal/middleware"
)

// ProfileEditHandler — страница редактирования профиля.
func ProfileEditHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.UserIDFromContext(r.Context())
    if !ok {
        http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
        return
    }

    // Получаем пользователя через UserService (добавь метод GetByID если надо)
    user, err := UserService.GetByID(r.Context(), userID)
    if err != nil {
        http.Error(w, "Ошибка загрузки профиля", http.StatusInternalServerError)
        return
    }

    if r.Method == http.MethodPost {
        name := strings.TrimSpace(r.FormValue("name"))
        email := strings.TrimSpace(r.FormValue("email"))
        if name != "" && email != "" {
            user.Name = name
            user.Email = email
            if err := UserService.Update(r.Context(), user); err != nil {
                http.Error(w, "Ошибка обновления профиля", http.StatusInternalServerError)
                return
            }
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
            return
        }
    }

    tmpl, err := template.ParseFiles(
        filepath.Join("internal", "templates", "base.html"),
        filepath.Join("internal", "templates", "profile_edit.html"),
    )
    if err != nil {
        http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
        return
    }
    tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"User": user,
	})
	
}

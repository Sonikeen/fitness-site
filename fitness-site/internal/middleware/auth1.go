// internal/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"net/http"

	"fitness-site/db"
	"fitness-site/internal/models"
)

type ctxKey string

const (
	UserIDKey ctxKey = "userID"
	UserKey   ctxKey = "user"
)

// UserIDFromContext возвращает ID пользователя из контекста и true, если он там есть.
func UserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(UserIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}

// GetUser возвращает (*models.User, true), если объект user есть в контексте.
func GetUser(r *http.Request) (*models.User, bool) {
	v := r.Context().Value(UserKey)
	if u, ok := v.(*models.User); ok {
		return u, true
	}
	return nil, false
}

// AuthMiddleware проверяет куку, получает userID из сессии,
// затем загружает пользователя из БД и кладёт в контекст {userID, *models.User}.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1) Проверяем наличие сессионной куки:
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
			return
		}

		// 2) Получаем userID по session_id:
		userID, err := GetUserIDBySession(cookie.Value)
		if err != nil {
			// если сессия недействительна, удаляем куку и редиректим:
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   -1,
			})
			http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
			return
		}

		// 3) Загружаем пользователя из БД (первый аргумент – контекст):
		row := db.Pool.QueryRow(r.Context(), `
SELECT id, username, email, password_hash,
       is_admin, age, height_cm, weight_kg, goals, avatar_url
  FROM users
 WHERE id = $1
`, userID)

		var u models.User
		if err := row.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.PasswordHash,
			&u.IsAdmin,
			&u.Age,
			&u.HeightCM,
			&u.WeightKG,
			&u.Goals,
			&u.AvatarURL,
		); err != nil {
			fmt.Println("AuthMiddleware: пользователь не найден или ошибка при Scan:", err)
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}

		// 4) Кладём в контекст и userID, и сам объект user:
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserKey, &u)

		// 5) Переходим к следующему обработчику:
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

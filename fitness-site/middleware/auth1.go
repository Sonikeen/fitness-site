package middleware

import (
	"context"
	"net/http"
)

type ctxKey string
const UserIDKey ctxKey = "userID"

func UserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(UserIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}


// AuthMiddleware — проверяет куку и сессию
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
			return
		}
		userID, err := GetUserIDBySession(cookie.Value)
		if err != nil {
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
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

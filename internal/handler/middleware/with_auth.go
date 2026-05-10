package middleware

import (
	"context"
	"net/http"

	"github.com/Vadich007/shortener/internal/model"
)

// AuthMiddleware проверяет наличие валидной JWT куки
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(model.COOKIE_NAME)

		var userID int

		if err == nil {
			userID, err = model.GetUserID(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		} else {
			userID = model.GenerateUserID()

			tokenString, tokenErr := model.BuildJWTString(userID)
			if tokenErr != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     model.COOKIE_NAME,
				Value:    tokenString,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   int(model.TOKEN_EXP.Seconds()),
			})
		}

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext извлекает UserID из контекста запроса.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("userID").(int)
	return userID, ok
}

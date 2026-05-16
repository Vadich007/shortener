package middleware

import (
	"context"
	"net/http"

	"github.com/Vadich007/shortener/internal/model"
)

type contextKey string

const userIDKey contextKey = "userID"

type AuthMiddleware struct {
	SecretKey string
}

func (m AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(model.CookieName)

		var userID int

		if err == nil {
			userID, err = model.GetUserID(cookie.Value, m.SecretKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		} else {
			userID = model.GenerateUserID()

			tokenString, tokenErr := model.BuildJWTString(userID, m.SecretKey)
			if tokenErr != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     model.CookieName,
				Value:    tokenString,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   int(model.TokenExp.Seconds()),
			})
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

func ContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

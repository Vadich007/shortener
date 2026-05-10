package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vadich007/shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewareNoCookie(t *testing.T) {
	var gotUserID int
	var gotOK bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, gotOK = GetUserIDFromContext(r.Context())
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	AuthMiddleware(next).ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, gotOK, "userID must be set in context")
	assert.NotZero(t, gotUserID)

	var authCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == model.COOKIE_NAME {
			authCookie = c
		}
	}
	require.NotNil(t, authCookie, "Set-Cookie header must be present")
	assert.NotEmpty(t, authCookie.Value)
}

func TestAuthMiddlewareValidCookie(t *testing.T) {
	const expectedUserID = 42

	token, err := model.BuildJWTString(expectedUserID)
	require.NoError(t, err)

	var gotUserID int
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = GetUserIDFromContext(r.Context())
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: model.COOKIE_NAME, Value: token})
	w := httptest.NewRecorder()
	AuthMiddleware(next).ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedUserID, gotUserID)
}

func TestAuthMiddlewareInvalidCookie(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: model.COOKIE_NAME, Value: "not.a.valid.jwt"})
	w := httptest.NewRecorder()
	AuthMiddleware(next).ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.False(t, called, "next handler must not be called on invalid cookie")
}

func TestGetUserIDFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, ok := GetUserIDFromContext(req.Context())
	assert.False(t, ok, "must return false when userID is not in context")
}

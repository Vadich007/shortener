package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/handler/middleware"
	"github.com/Vadich007/shortener/internal/model"
	"github.com/Vadich007/shortener/internal/repository/memory"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/Vadich007/shortener/pkg/shorter"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const storagePath = "test_storage.json"

func Fixture(t *testing.T) {
	file, err := os.Create(storagePath)
	if err != nil {
		panic(err)
	}
	file.Close()
	t.Cleanup(func() {
		os.Remove(storagePath)
	})
}

func TestHandleGetMethodNotAllowed(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	w := httptest.NewRecorder()

	hand.HandleGet(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestHandleGetNotFound(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodGet, "/asdsad", nil)
	w := httptest.NewRecorder()

	hand.HandleGet(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestHandleGetSuccess(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Get("/{shortedLink}", hand.HandleGet)

	originalLink := "https://practicum.yandex.ru/"
	shortedLink, _ := serv.AddLink(originalLink, 0)

	req := httptest.NewRequest(http.MethodGet, "/"+strings.Split(shortedLink, "/")[3], nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, resp.Header.Get("Location"), originalLink)
}

func TestHandlePostSuccess(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"
	shortedLink := "http://localhost:8080/" + shorter.Shorten(originalLink)
	body := bytes.NewBufferString(originalLink)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	hand.HandlePost(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	actual, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, string(actual), shortedLink)
}

func TestHandlePostEmptyBody(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	hand.HandlePost(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestHandlePostEmptyStringBody(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	hand.HandlePost(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestHandlePostJsonNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"
	shortedLink := "http://localhost:8080/" + shorter.Shorten(originalLink)

	rawResp := model.Response{
		Result: shortedLink,
	}

	rawReq := model.Request{
		URL: originalLink,
	}

	jsonData, _ := json.Marshal(rawReq)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	hand.HandlePostJSON(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	actual, _ := io.ReadAll(resp.Body)

	jsonDataResp, _ := json.Marshal(rawResp)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, string(actual), string(jsonDataResp)+"\n")
}

func TestHandlePostJsonExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"
	shortedLink := "http://localhost:8080/" + shorter.Shorten(originalLink)

	serv.AddLink(originalLink, 0)

	rawResp := model.Response{
		Result: shortedLink,
	}

	rawReq := model.Request{
		URL: originalLink,
	}

	jsonData, _ := json.Marshal(rawReq)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	hand.HandlePostJSON(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	actual, _ := io.ReadAll(resp.Body)

	jsonDataResp, _ := json.Marshal(rawResp)
	assert.Equal(t, resp.StatusCode, http.StatusConflict)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, string(actual), string(jsonDataResp)+"\n")
}

func TestGetUserUrlsNoContent(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	w := httptest.NewRecorder()

	hand.GetUserUrls(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGetUserUrlsSuccess(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	const userID = 5
	originalURL := "https://example.com"
	serv.AddLink(originalURL, userID)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", userID))
	w := httptest.NewRecorder()

	hand.GetUserUrls(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var urls []model.UserURLResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&urls))
	assert.Len(t, urls, 1)
	assert.Equal(t, originalURL, urls[0].OriginalURL)
	assert.Equal(t, "http://localhost:8080/"+shorter.Shorten(originalURL), urls[0].ShortURL)
}

func TestGetUserUrlsOnlyOwnUrls(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	serv.AddLink("https://user10.com", 10)
	serv.AddLink("https://user20.com", 20)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", 10))
	w := httptest.NewRecorder()

	hand.GetUserUrls(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var urls []model.UserURLResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&urls))
	assert.Len(t, urls, 1)
	assert.Equal(t, "https://user10.com", urls[0].OriginalURL)
}

func TestGetUserUrlsInvalidCookie(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)
	r.Get("/api/user/urls", hand.GetUserUrls)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "this.is.invalid"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestHandleGetGone(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Get("/{shortedLink}", hand.HandleGet)

	repo.AddLink("abc123", "https://example.com", 1)
	repo.DeleteURLsBatch(1, []string{"abc123"})

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusGone, resp.StatusCode)
}

func TestDeleteUserUrlsAccepted(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	body, _ := json.Marshal([]string{"abc123", "def456"})
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	w := httptest.NewRecorder()

	hand.DeleteUserUrls(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestDeleteUserUrlsWrongContentType(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString(`["abc"]`))
	req.Header.Set("Content-Type", "text/plain")
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	w := httptest.NewRecorder()

	hand.DeleteUserUrls(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestDeleteUserUrlsInvalidJSON(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	w := httptest.NewRecorder()

	hand.DeleteUserUrls(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteUserUrlsThenGetGone(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)
	r.Get("/{shortedLink}", hand.HandleGet)
	r.Delete("/api/user/urls", hand.DeleteUserUrls)

	// создаём JWT-куку для userID=42
	token, _ := model.BuildJWTString(42)

	// сохраняем ссылку напрямую в репозиторий
	shortKey := "myshort"
	repo.AddLink(shortKey, "https://example.com", 42)

	// DELETE запрос
	body, _ := json.Marshal([]string{shortKey})
	delReq := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(body))
	delReq.Header.Set("Content-Type", "application/json")
	delReq.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	delW := httptest.NewRecorder()
	r.ServeHTTP(delW, delReq)
	assert.Equal(t, http.StatusAccepted, delW.Result().StatusCode)

	// ждём асинхронного удаления
	assert.Eventually(t, func() bool {
		getReq := httptest.NewRequest(http.MethodGet, "/"+shortKey, nil)
		getReq.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
		getW := httptest.NewRecorder()
		r.ServeHTTP(getW, getReq)
		return getW.Result().StatusCode == http.StatusGone
	}, 2*time.Second, 100*time.Millisecond)
}

func TestGetUserUrlsNoCookieSetsNewCookie(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)
	r.Get("/api/user/urls", hand.GetUserUrls)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	var authCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "auth_token" {
			authCookie = c
		}
	}
	require.NotNil(t, authCookie, "must set auth_token cookie for new user")
}

func TestHandlePostJsonWrongHeader(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"

	serv.AddLink(originalLink, 0)

	rawReq := model.Request{
		URL: originalLink,
	}

	jsonData, _ := json.Marshal(rawReq)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	hand.HandlePostJSON(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	actual, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(t, string(actual), "Unprocessable entity\n")
}

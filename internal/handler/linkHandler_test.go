package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/Vadich007/shortener/pkg/shorter"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTPMethodNotAllowed(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
}

func TestServeHTTPNotFound(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodGet, "/asdsad", nil)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestServeHTTPGetSuccess(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	originalLink := "https://practicum.yandex.ru/"
	shortedLink, _ := serv.AddLink(originalLink)

	req := httptest.NewRequest(http.MethodGet, "/"+strings.Split(shortedLink, "/")[3], nil)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, resp.Header.Get("Location"), originalLink)
}

func TestServeHTTPPostSuccess(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"
	shortedLink := "http://localhost:8080/" + shorter.Shorten(originalLink)
	body := bytes.NewBufferString(originalLink)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	actual, _ := io.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, string(actual), shortedLink)
}

func TestServeHTTPPostEmptyBody(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestServeHTTPPostEmptyStringBody(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

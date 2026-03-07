package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/Vadich007/shortener/pkg/shorter"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTPMethodNotAllowed(t *testing.T) {
	repo, _ := repository.NewInMemoryLinkRepository()
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
	repo, _ := repository.NewInMemoryLinkRepository()
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
	repo, _ := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	originalLink := "asdsad"
	shortedLink, _ := serv.AddLink(originalLink)

	req := httptest.NewRequest(http.MethodGet, "/"+shortedLink, nil)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, resp.Header["Location"], originalLink)
}

func TestServeHTTPPostSuccess(t *testing.T) {
	repo, _ := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"
	shortedLink := shorter.Shorten(originalLink)
	body := bytes.NewBufferString(originalLink)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, resp.Body, shortedLink)
}

func TestServeHTTPPostEmptyBody(t *testing.T) {
	repo, _ := repository.NewInMemoryLinkRepository()
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
	repo, _ := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := NewLinkHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	hand.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

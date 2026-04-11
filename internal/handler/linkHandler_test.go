package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/Vadich007/shortener/pkg/shorter"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Get("/{shortedLink}", hand.HandleGet)

	originalLink := "https://practicum.yandex.ru/"
	shortedLink, _ := serv.AddLink(originalLink)

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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"
	shortedLink := "http://localhost:8080/" + shorter.Shorten(originalLink)

	serv.AddLink(originalLink)

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

func TestHandlePostJsonWrongHeader(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := service.NewLinkService(repo, conf)
	hand := NewLinkHandler(serv)

	originalLink := "example.com"

	serv.AddLink(originalLink)

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

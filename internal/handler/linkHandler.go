package handler

import (
	"io"
	"net/http"

	"github.com/Vadich007/shortener/internal/service"
)

type LinkHandler struct {
	service service.Service
}

func NewLinkHandler(service service.Service) *LinkHandler {
	return &LinkHandler{service: service}
}

func (h *LinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LinkHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	originalLink, err := h.service.GetLink(r.URL.Path[1:])
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *LinkHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	bodyString := string(body)
	if bodyString == "" {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	shortedLink, err := h.service.AddLink(bodyString)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	data := []byte(shortedLink)
	w.Write(data)
}

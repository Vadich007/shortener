package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Vadich007/shortener/internal/model"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	service service.Service
}

func NewLinkHandler(service service.Service) *LinkHandler {
	return &LinkHandler{service: service}
}

func (h *LinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortedLink := chi.URLParam(r, "shortedLink")
	originalLink, err := h.service.GetLink(shortedLink)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *LinkHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
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

func (h *LinkHandler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	var req model.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	shortedLink, err := h.service.AddLink(req.URL)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	resp := model.Response{
		Result: shortedLink,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Vadich007/shortener/internal/handler/middleware"
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
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	shortedLink, err := h.service.AddLink(bodyString, userID)
	if err != nil {
		var linkAlreadyExistError *model.LinkAlreadyExistError
		if errors.As(err, &linkAlreadyExistError) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortedLink))
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortedLink))
}

func (h *LinkHandler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	var req model.Request

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	shortedLink, err := h.service.AddLink(req.URL, userID)

	resp := model.Response{
		Result: shortedLink,
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		var linkAlreadyExistError *model.LinkAlreadyExistError
		if errors.As(err, &linkAlreadyExistError) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(resp)
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (h *LinkHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.service.PingDB(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *LinkHandler) Batch(w http.ResponseWriter, r *http.Request) {
	var req []model.BatchRecordRequest

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	resp, err := h.service.AddLinksBatch(req, userID)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (h *LinkHandler) GetUserUrls(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	urls, err := h.service.GetUserUrls(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

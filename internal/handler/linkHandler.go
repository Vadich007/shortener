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
		h.service.GetLink(r.URL.Path[1:])
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		bodyString := string(body)
		h.service.AddLink(bodyString)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

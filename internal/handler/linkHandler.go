package handler

import (
	"net/http"

	"github.com/Vadich007/shortener/internal/service"
)

type LinkHandler struct {
	service *service.LinkService
}

func (h *LinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.service.GetLink(w)
	case http.MethodPost:
		h.service.AddLink(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

package main

import (
	"net/http"

	"github.com/Vadich007/shortener/internal/handler"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo)
	hand := handler.NewLinkHandler(serv)

	r := chi.NewRouter()
	r.Get("/{shortedLink}", hand.HandleGet)
	r.Post("/", hand.HandlePost)
	http.ListenAndServe(":8080", r)
}

package main

import (
	"log"
	"net/http"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/handler"
	"github.com/Vadich007/shortener/internal/handler/middleware"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	conf := config.GetConfig()
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := handler.NewLinkHandler(serv)

	r := chi.NewRouter()

	r.Use(middleware.WithLogging)

	r.Get("/{shortedLink}", hand.HandleGet)
	r.Post("/", hand.HandlePost)

	if err := http.ListenAndServe(conf.ServerAddress, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/Vadich007/shortener/internal/config/flags"
	"github.com/Vadich007/shortener/internal/handler"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	f := flags.ProcessingFlags()
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo, f)
	hand := handler.NewLinkHandler(serv)

	r := chi.NewRouter()

	r.Get(f.B+"/{shortedLink}", hand.HandleGet)
	r.Post(f.B+"/", hand.HandlePost)

	if err := http.ListenAndServe(f.A, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

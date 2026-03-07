package main

import (
	"net/http"

	"github.com/Vadich007/shortener/internal/handler"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
)

func main() {
	repo, err := repository.NewInMemoryLinkRepository()
	if err != nil {
		return
	}
	serv := service.NewLinkService(repo)
	hand := handler.NewLinkHandler(serv)

	mux := http.NewServeMux()
	mux.Handle("/", hand)
	http.ListenAndServe(":8080", mux)
}

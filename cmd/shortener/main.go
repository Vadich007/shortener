package main

import (
	"net/http"

	"github.com/Vadich007/shortener/internal/handler"
)

func main() {
	handler = 
	mux := http.NewServeMux()
	mux.Handle("/", handler.LinkHandler)
}

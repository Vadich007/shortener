package main

import (
	"net/http"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/handler"
	"github.com/Vadich007/shortener/internal/handler/middleware"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar = *logger.Sugar()
	conf := config.GetConfig()
	repo := repository.NewInMemoryLinkRepository()
	serv := service.NewLinkService(repo, conf)
	hand := handler.NewLinkHandler(serv)
	loggingMiddleware := middleware.LoggingMiddleware{Sugar: sugar}

	sugar.Infow(
		"Starting server",
		"addr", conf.ServerAddress,
	)

	r := chi.NewRouter()

	r.Use(loggingMiddleware.WithLogging)
	r.Use(middleware.WithCompress)

	r.Get("/{shortedLink}", hand.HandleGet)
	r.Post("/", hand.HandlePost)
	r.Post("/api/shorten", hand.HandlePostJson)

	if err := http.ListenAndServe(conf.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}

package main

import (
	"net/http"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/handler"
	"github.com/Vadich007/shortener/internal/handler/middleware"
	"github.com/Vadich007/shortener/internal/repository/factory"
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
	repo, err := factory.GetRepository(conf)
	if err != nil {
		panic(err)
	}
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
	r.Use(middleware.AuthMiddleware)

	r.Get("/{shortedLink}", hand.HandleGet)
	r.Get("/ping", hand.PingDB)
	r.Post("/", hand.HandlePost)
	r.Post("/api/shorten", hand.HandlePostJSON)
	r.Post("/api/shorten/batch", hand.Batch)
	r.Get("/api/user/urls", hand.GetUserUrls)

	if err := http.ListenAndServe(conf.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}

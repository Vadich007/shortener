package main

import (
	"net/http"
	"os/signal"
	"syscall"

	"context"

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	serv := service.NewLinkService(ctx, repo, conf)
	hand := handler.NewLinkHandler(serv)
	loggingMiddleware := middleware.LoggingMiddleware{Sugar: sugar}
	authMiddleware := middleware.AuthMiddleware{SecretKey: conf.SecretKey}

	sugar.Infow(
		"Starting server",
		"addr", conf.ServerAddress,
	)

	r := chi.NewRouter()

	r.Use(loggingMiddleware.WithLogging)
	r.Use(middleware.WithCompress)
	r.Use(authMiddleware.Handle)

	r.Get("/{shortedLink}", hand.HandleGet)
	r.Get("/ping", hand.PingDB)
	r.Post("/", hand.HandlePost)
	r.Post("/api/shorten", hand.HandlePostJSON)
	r.Post("/api/shorten/batch", hand.Batch)
	r.Get("/api/user/urls", hand.GetUserUrls)
	r.Delete("/api/user/urls", hand.DeleteUserUrls)

	if err := http.ListenAndServe(conf.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}

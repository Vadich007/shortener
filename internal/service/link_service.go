package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type LinkService struct {
	repository repository.LinkRepository
	conf       config.Config
	db         *sql.DB
}

func NewLinkService(r repository.LinkRepository, conf config.Config) *LinkService {
	db, err := sql.Open("pgx", conf.DatabaseDsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return &LinkService{repository: r, conf: conf, db: db}
}

func (s *LinkService) GetLink(shortedLink string) (string, error) {
	return s.repository.GetLink(shortedLink)
}

func (s *LinkService) AddLink(originalLink string) (string, error) {
	shortedLink := shorter.Shorten(originalLink)
	return s.conf.BaseURL + "/" + shortedLink, s.repository.AddLink(shortedLink, originalLink)
}

func (s *LinkService) PingDb() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return sql.ErrConnDone
	}
	return nil
}

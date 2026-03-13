package service

import (
	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
)

type LinkService struct {
	repository repository.LinkRepository
	conf       config.Config
}

func NewLinkService(r repository.LinkRepository, conf config.Config) *LinkService {
	return &LinkService{repository: r, conf: conf}
}

func (s *LinkService) GetLink(shortedLink string) (string, error) {
	return s.repository.GetLink(shortedLink)
}

func (s *LinkService) AddLink(originalLink string) (string, error) {
	shortedLink := shorter.Shorten(originalLink)
	return s.conf.BaseURL + "/" + shortedLink, s.repository.AddLink(shortedLink, originalLink)
}

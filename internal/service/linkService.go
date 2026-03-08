package service

import (
	"github.com/Vadich007/shortener/internal/config/flags"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
)

type LinkService struct {
	repository repository.LinkRepository
	flags      flags.Flags
}

func NewLinkService(r repository.LinkRepository, f flags.Flags) *LinkService {
	return &LinkService{repository: r, flags: f}
}

func (s *LinkService) GetLink(shortedLink string) (string, error) {
	return s.repository.GetLink(shortedLink)
}

func (s *LinkService) AddLink(originalLink string) (string, error) {
	shortedLink := shorter.Shorten(originalLink)
	return "http://" + s.flags.A + s.flags.B + "/" + shortedLink, s.repository.AddLink(shortedLink, originalLink)
}

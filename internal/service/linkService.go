package service

import "github.com/Vadich007/shortener/internal/repository"

type LinkService struct {
	repository repository.LinkRepository
}

func NewLinkService(r repository.LinkRepository) *LinkService {
	return &LinkService{repository: r}
}

func (s *LinkService) GetLink(shortedLink string) (string, error) {
	return s.repository.GetLink(shortedLink)
}

func (s *LinkService) AddLink(originalLink string) (string, error) {
	return s.repository.AddLink(originalLink)
}

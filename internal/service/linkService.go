package service

import "github.com/Vadich007/shortener/internal/repository"

type LinkService struct {
	repository *repository.LinkRepository
}

func NewLinkService(r repository.LinkRepository) *LinkService {
	return &LinkService{repository: &r}
}

func (s *LinkService) GetLink(id string) (string, error) {
	return "", nil
}

func (s *LinkService) AddLink(originalLink string) (string, error) {
	return "", nil
}

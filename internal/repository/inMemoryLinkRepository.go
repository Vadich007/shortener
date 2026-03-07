package repository

import (
	"errors"
)

type InMemoryLinkRepository struct {
	m map[string]string
}

func NewInMemoryLinkRepository() (*InMemoryLinkRepository, error) {
	m := make(map[string]string)
	return &InMemoryLinkRepository{m: m}, nil
}

func (r *InMemoryLinkRepository) GetLink(shortedLink string) (string, error) {
	originalLink, exist := r.m[shortedLink]
	if exist {
		return originalLink, nil
	}
	return "", errors.New("link doesn't exist")
}

func (r *InMemoryLinkRepository) AddLink(shortedLink string, originalLink string) error {
	_, err := r.GetLink(shortedLink)
	if err == nil {
		return nil
	}

	r.m[shortedLink] = originalLink
	return nil
}

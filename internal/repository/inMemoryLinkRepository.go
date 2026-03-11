package repository

import (
	"errors"
	"sync"
)

type InMemoryLinkRepository struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewInMemoryLinkRepository() *InMemoryLinkRepository {
	m := make(map[string]string)
	return &InMemoryLinkRepository{m: m}
}

func (r *InMemoryLinkRepository) GetLink(shortedLink string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	originalLink, exist := r.m[shortedLink]
	if exist {
		return originalLink, nil
	}
	return "", errors.New("link doesn't exist")
}

func (r *InMemoryLinkRepository) AddLink(shortedLink string, originalLink string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.m[shortedLink]; exist {
		return nil
	}

	r.m[shortedLink] = originalLink
	return nil
}

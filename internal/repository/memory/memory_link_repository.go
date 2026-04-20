package memory

import (
	"errors"
	"sync"
)

type MemoryLinkRepository struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewMemoryLinkRepository() (*MemoryLinkRepository, error) {
	m := make(map[string]string)

	return &MemoryLinkRepository{m: m}, nil
}

func (r *MemoryLinkRepository) GetLink(shortedLink string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	originalLink, exist := r.m[shortedLink]
	if exist {
		return originalLink, nil
	}
	return "", errors.New("link doesn't exist")
}

func (r *MemoryLinkRepository) AddLink(shortedLink string, originalLink string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.m[shortedLink]; exist {
		return nil
	}

	r.m[shortedLink] = originalLink
	return nil
}

func (r *MemoryLinkRepository) PingDB() error {
	return nil
}

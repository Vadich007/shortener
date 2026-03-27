package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
)

type InMemoryLinkRepository struct {
	mu   sync.RWMutex
	m    map[string]string
	path string
}

func NewInMemoryLinkRepository(conf config.Config) (*InMemoryLinkRepository, error) {
	m := make(map[string]string)
	data, err := os.ReadFile(conf.FileStoragePath)
	if err != nil {
		return nil, err
	}
	var records []model.StorageRecord

	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		m[record.ShortedURL] = record.OriginalURL
	}

	return &InMemoryLinkRepository{m: m, path: conf.FileStoragePath}, nil
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
	return r.addToFile(shortedLink, originalLink)
}

func (r *InMemoryLinkRepository) addToFile(shortedLink string, originalLink string) error {
	file, err := os.OpenFile(r.path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	record := model.StorageRecord{ShortedURL: shortedLink,
		OriginalURL: originalLink}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = file.Write(data)

	return err
}

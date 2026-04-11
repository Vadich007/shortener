package repository

import (
	"encoding/json"
	"errors"
	"io"
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
	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &InMemoryLinkRepository{m: m, path: conf.FileStoragePath}, nil
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
	return r.saveFile()
}

func (r *InMemoryLinkRepository) saveFile() error {
	file, err := os.OpenFile(r.path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var records []model.StorageRecord
	for shortedLink, originalLink := range r.m {
		records = append(records, model.StorageRecord{ShortedURL: shortedLink,
			OriginalURL: originalLink})
	}
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}
	_, err = file.Write(data)

	return err
}

package memory

import (
	"errors"
	"sync"

	"github.com/Vadich007/shortener/internal/model"
)

type MemoryLinkRepository struct {
	mu sync.RWMutex
	m  map[string]model.StorageRecord
}

func NewMemoryLinkRepository() (*MemoryLinkRepository, error) {
	m := make(map[string]model.StorageRecord)
	return &MemoryLinkRepository{m: m}, nil
}

func (r *MemoryLinkRepository) GetLink(shortedLink string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exist := r.m[shortedLink]
	if exist {
		return record.OriginalURL, nil
	}
	return "", errors.New("link doesn't exist")
}

func (r *MemoryLinkRepository) AddLink(shortedLink, originalLink string, userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.m[shortedLink]; exist {
		return model.NewLinkAlreadyExistError(shortedLink)
	}

	r.m[shortedLink] = model.StorageRecord{
		ShortedURL:  shortedLink,
		OriginalURL: originalLink,
		UserID:      userID,
	}
	return nil
}

func (r *MemoryLinkRepository) PingDB() error {
	return nil
}

func (r *MemoryLinkRepository) AddLinksBatch(request []model.BatchRecordRequest, shortedMap map[string]string, userID int) error {
	for _, record := range request {
		r.AddLink(shortedMap[record.CorrelationID], record.OriginalURL, userID)
	}
	return nil
}

func (r *MemoryLinkRepository) GetUserUrls(userID int) ([]model.UserURLResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.UserURLResponse
	for _, record := range r.m {
		if record.UserID == userID {
			result = append(result, model.UserURLResponse{
				ShortURL:    record.ShortedURL,
				OriginalURL: record.OriginalURL,
			})
		}
	}
	return result, nil
}

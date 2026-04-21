package file

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
)

type FileLinkRepository struct {
	mu   sync.RWMutex
	m    map[string]string
	path string
}

func NewFileLinkRepository(conf config.Config) (*FileLinkRepository, error) {
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
		return &FileLinkRepository{m: m, path: conf.FileStoragePath}, nil
	}

	var records []model.StorageRecord

	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		m[record.ShortedURL] = record.OriginalURL
	}

	return &FileLinkRepository{m: m, path: conf.FileStoragePath}, nil
}

func (r *FileLinkRepository) GetLink(shortedLink string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	originalLink, exist := r.m[shortedLink]
	if exist {
		return originalLink, nil
	}
	return "", errors.New("link doesn't exist")
}

func (r *FileLinkRepository) AddLink(shortedLink string, originalLink string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.m[shortedLink]; exist {
		return model.NewLinkAlreadyExistError(shortedLink)
	}

	r.m[shortedLink] = originalLink
	return r.saveFile()
}

func (r *FileLinkRepository) saveFile() error {
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

func (r *FileLinkRepository) AddLinksBatch(request *model.BatchRequest, shortedMap map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range request.Records {
		shortedLink := shortedMap[record.CorrelationID]
		if _, exist := r.m[shortedLink]; exist {
			continue
		}

		r.m[shortedLink] = record.OriginalURL
	}

	return r.saveFile()
}

func (r *FileLinkRepository) PingDB() error {
	return nil
}

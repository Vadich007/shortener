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
	m    map[string]model.StorageRecord
	path string
}

func NewFileLinkRepository(conf config.Config) (*FileLinkRepository, error) {
	m := make(map[string]model.StorageRecord)
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
	if err = json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	for _, record := range records {
		m[record.ShortedURL] = record
	}

	return &FileLinkRepository{m: m, path: conf.FileStoragePath}, nil
}

func (r *FileLinkRepository) GetLink(shortedLink string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exist := r.m[shortedLink]
	if exist {
		return record.OriginalURL, nil
	}
	return "", errors.New("link doesn't exist")
}

func (r *FileLinkRepository) AddLink(shortedLink, originalLink string, userID int) error {
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
	return r.saveFile()
}

func (r *FileLinkRepository) saveFile() error {
	file, err := os.OpenFile(r.path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var records []model.StorageRecord
	for _, record := range r.m {
		records = append(records, record)
	}
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func (r *FileLinkRepository) AddLinksBatch(request []model.BatchRecordRequest, shortedMap map[string]string, userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range request {
		shortedLink := shortedMap[record.CorrelationID]
		if _, exist := r.m[shortedLink]; exist {
			continue
		}
		r.m[shortedLink] = model.StorageRecord{
			ShortedURL:  shortedLink,
			OriginalURL: record.OriginalURL,
			UserID:      userID,
		}
	}
	return r.saveFile()
}

func (r *FileLinkRepository) PingDB() error {
	return nil
}

func (r *FileLinkRepository) GetUserUrls(userID int) ([]model.UserURLResponse, error) {
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

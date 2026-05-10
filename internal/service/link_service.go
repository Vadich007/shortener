package service

import (
	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
)

type LinkService struct {
	repository repository.LinkRepository
	conf       config.Config
}

func NewLinkService(r repository.LinkRepository, conf config.Config) *LinkService {
	return &LinkService{repository: r, conf: conf}
}

func (s *LinkService) GetLink(shortedLink string) (string, error) {
	return s.repository.GetLink(shortedLink)
}

func (s *LinkService) AddLink(originalLink string, userID int) (string, error) {
	shortedLink := shorter.Shorten(originalLink)
	return s.conf.BaseURL + "/" + shortedLink, s.repository.AddLink(shortedLink, originalLink, userID)
}

func (s *LinkService) PingDB() error {
	return s.repository.PingDB()
}

func (s *LinkService) AddLinksBatch(request []model.BatchRecordRequest, userID int) ([]model.BatchRecordResponse, error) {
	m := make(map[string]string)
	var response []model.BatchRecordResponse

	for _, originalRecord := range request {
		shortedLink := shorter.Shorten(originalRecord.OriginalURL)
		m[originalRecord.CorrelationID] = shortedLink
		response = append(response, model.BatchRecordResponse{
			CorrelationID: originalRecord.CorrelationID,
			ShortedURL:    s.conf.BaseURL + "/" + shortedLink,
		})
	}

	return response, s.repository.AddLinksBatch(request, m, userID)
}

func (s *LinkService) GetUserUrls(userID int) ([]model.UserURLResponse, error) {
	records, err := s.repository.GetUserUrls(userID)
	if err != nil {
		return nil, err
	}
	for i := range records {
		records[i].ShortURL = s.conf.BaseURL + "/" + records[i].ShortURL
	}
	return records, nil
}

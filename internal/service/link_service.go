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

func (s *LinkService) AddLink(originalLink string) (string, error) {
	shortedLink := shorter.Shorten(originalLink)
	return s.conf.BaseURL + "/" + shortedLink, s.repository.AddLink(shortedLink, originalLink)
}

func (s *LinkService) PingDB() error {
	return s.repository.PingDB()
}

func (s *LinkService) AddLinksBatch(request *model.BatchRequest) (*model.BatchResponse, error) {
	m := make(map[string]string)
	var response model.BatchResponse

	for _, originalRecord := range request.Records {
		shortedLink := shorter.Shorten(originalRecord.OriginalURL)
		m[originalRecord.CorrelationID] = shortedLink
		response.Records = append(response.Records, model.BatchRecordResponse{
			CorrelationID: originalRecord.CorrelationID,
			ShortedURL:    shortedLink,
		})
	}

	return &response, s.repository.AddLinksBatch(request, m)
}

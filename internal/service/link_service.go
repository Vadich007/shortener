package service

import (
	"time"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
)

type deleteTask struct {
	userID   int
	shortURL string
}

type LinkService struct {
	repository repository.LinkRepository
	conf       config.Config
	deleteCh   chan deleteTask
}

func NewLinkService(r repository.LinkRepository, conf config.Config) *LinkService {
	s := &LinkService{
		repository: r,
		conf:       conf,
		deleteCh:   make(chan deleteTask, 1024),
	}
	go s.flushDeletes()
	return s
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

// DeleteURLs принимает список коротких идентификаторов и асинхронно помечает их удалёнными.
// Каждый запрос порождает горутину, которая проталкивает задачи в общий deleteCh (паттерн fanIn).
func (s *LinkService) DeleteURLs(userID int, shortURLs []string) {
	go func() {
		for _, u := range shortURLs {
			s.deleteCh <- deleteTask{userID: userID, shortURL: u}
		}
	}()
}

// flushDeletes — единственный потребитель deleteCh. Накапливает задачи и батчами отправляет в репозиторий.
func (s *LinkService) flushDeletes() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var batch []deleteTask
	for {
		select {
		case task := <-s.deleteCh:
			batch = append(batch, task)
		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}
			byUser := make(map[int][]string)
			for _, t := range batch {
				byUser[t.userID] = append(byUser[t.userID], t.shortURL)
			}
			for uid, urls := range byUser {
				s.repository.DeleteURLsBatch(uid, urls)
			}
			batch = batch[:0]
		}
	}
}

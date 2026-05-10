package repository

import "github.com/Vadich007/shortener/internal/model"

type LinkRepository interface {
	GetLink(string) (string, error)
	AddLink(shortedLink, originalLink string, userID int) error
	PingDB() error
	AddLinksBatch(request []model.BatchRecordRequest, shortedMap map[string]string, userID int) error
	GetUserUrls(userID int) ([]model.UserURLResponse, error)
}

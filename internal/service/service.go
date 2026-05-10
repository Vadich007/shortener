package service

import "github.com/Vadich007/shortener/internal/model"

type Service interface {
	GetLink(shortedLink string) (string, error)
	AddLink(originalLink string, userID int) (string, error)
	PingDB() error
	AddLinksBatch(request []model.BatchRecordRequest, userID int) ([]model.BatchRecordResponse, error)
	GetUserUrls(userID int) ([]model.UserURLResponse, error)
}

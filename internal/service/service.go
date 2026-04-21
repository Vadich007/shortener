package service

import "github.com/Vadich007/shortener/internal/model"

type Service interface {
	GetLink(shortedLink string) (string, error)
	AddLink(originalLink string) (string, error)
	PingDB() error
	AddLinksBatch(*model.BatchRequest) (*model.BatchResponse, error)
}

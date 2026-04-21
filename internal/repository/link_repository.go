package repository

import "github.com/Vadich007/shortener/internal/model"

type LinkRepository interface {
	GetLink(string) (string, error)
	AddLink(string, string) error
	PingDB() error
	AddLinksBatch(*model.BatchRequest, map[string]string) error
}

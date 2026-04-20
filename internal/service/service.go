package service

type Service interface {
	GetLink(shortedLink string) (string, error)
	AddLink(originalLink string) (string, error)
	PingDB() error
}

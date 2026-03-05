package service

type Service interface {
	GetLink(id string) (string, error)
	AddLink(originalLink string) (string, error)
}

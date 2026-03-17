package service

type Service interface {
	GetLink(shortedLink string) (string, error)
	GetLinkByOriginal(originalLink string) (string, error)
	AddLink(originalLink string) (string, error)
}

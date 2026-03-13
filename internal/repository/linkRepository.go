package repository

type LinkRepository interface {
	GetLink(string) (string, error)
	AddLink(string, string) error
}

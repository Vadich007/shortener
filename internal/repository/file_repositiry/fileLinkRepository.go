package repository

type FileLinkRepository struct {
	path string
}

func NewLinkRepository(p string) *FileLinkRepository {
	return &FileLinkRepository{
		path: p,
	}
}

func GetLink(shortedLink string) (string, error) {
	return "", nil
}

func AddLink(originalLink string) (string, error) {
	return "", nil
}

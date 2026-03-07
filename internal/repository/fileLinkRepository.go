package repository

import (
	"errors"
)

type FileLinkRepository struct {
	m map[string]string
}

func NewFileLinkRepository(path string) (*FileLinkRepository, error) {
	m := make(map[string]string)
	return &FileLinkRepository{m: m}, nil
}

func (r *FileLinkRepository) GetLink(shortedLink string) (string, error) {
	originalLink, exist := r.m[shortedLink]
	if exist {
		return originalLink, nil
	}
	return "", errors.New("Link doesn't exist")
}

func (r *FileLinkRepository) AddLink(shortedLink string, originalLink string) error {
	_, err := r.GetLink(shortedLink)
	if err == nil {
		return nil
	}

	r.m[shortedLink] = originalLink
	return nil
}

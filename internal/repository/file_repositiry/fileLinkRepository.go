package repository

import (
	"encoding/csv"
	"errors"
	"os"
)

type FileLinkRepository struct {
	m map[string]string
}

func NewLinkRepository(path string) (*FileLinkRepository, error) {
	m := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		m[record[0]] = record[1]
	}

	return &FileLinkRepository{m: m}, nil
}

func (r *FileLinkRepository) GetLink(shortedLink string) (string, error) {
	originalLink, exist := r.m[shortedLink]
	if exist {
		return originalLink, nil
	}
	return "", errors.New("Link doesn't exist")
}

func (r *FileLinkRepository) AddLink(originalLink string) (string, error) {
	shortedLink, err := r.GetLink(originalLink)
	if err == nil {
		return shortedLink, nil
	}

	// TODO: логика генерации строки
	return "", nil
}

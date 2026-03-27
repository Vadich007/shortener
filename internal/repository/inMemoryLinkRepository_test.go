package repository

import (
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetLinkNotExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "./test_storage.json"}
	repo, _ := NewInMemoryLinkRepository(conf)
	link, err := repo.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "./test_storage.json"}
	repo, _ := NewInMemoryLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
	link, err := repo.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "./test_storage.json"}
	repo, _ := NewInMemoryLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkNotExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "./test_storage.json"}
	repo, _ := NewInMemoryLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

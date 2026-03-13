package service

import (
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
	"github.com/stretchr/testify/assert"
)

func TestGetLinkNotExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	link, err := serv.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo, conf)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)

	link, err := serv.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo, conf)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)
	repo.AddLink(shorter.Shorten(originalName), originalName)

	actualShortedLink, err := serv.AddLink(originalName)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

func TestAddLinkNotExist(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080"}
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo, conf)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)

	actualShortedLink, err := serv.AddLink(originalName)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

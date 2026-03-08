package service

import (
	"testing"

	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
	"github.com/stretchr/testify/assert"
)

func TestGetLinkNotExist(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo)

	link, err := serv.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)

	link, err := serv.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)
	repo.AddLink(shorter.Shorten(originalName), originalName)

	actualShortedLink, err := serv.AddLink(originalName)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

func TestAddLinkNotExist(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	serv := NewLinkService(repo)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)

	actualShortedLink, err := serv.AddLink(originalName)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

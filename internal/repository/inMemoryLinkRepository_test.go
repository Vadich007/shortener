package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLinkNotExist(t *testing.T) {
	repo := NewInMemoryLinkRepository()
	link, err := repo.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	repo := NewInMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
	link, err := repo.GetLink(originalName)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	repo := NewInMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkNotExist(t *testing.T) {
	repo := NewInMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

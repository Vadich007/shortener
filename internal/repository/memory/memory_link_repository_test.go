package memory

import (
	"testing"

	"github.com/Vadich007/shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetLinkNotExist(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	link, err := repo.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
	link, err := repo.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, model.NewLinkAlreadyExistError(shortedLink))
}

func TestAddLinkNotExist(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

func TestPingDB(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	assert.Equal(t, repo.PingDB(), nil)
}

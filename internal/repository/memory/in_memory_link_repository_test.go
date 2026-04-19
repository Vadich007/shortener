package memory

import (
	"os"
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/stretchr/testify/assert"
)

const storagePath = "test_storage.json"

func Fixture(t *testing.T) {
	file, err := os.Create(storagePath)
	if err != nil {
		panic(err)
	}
	file.Close()
	t.Cleanup(func() {
		os.Remove(storagePath)
	})
}

func TestGetLinkNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewInMemoryLinkRepository(conf)
	link, err := repo.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
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
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewInMemoryLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewInMemoryLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName)
	assert.Equal(t, err, nil)
}

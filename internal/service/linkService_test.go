package service

import (
	"os"
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/repository"
	"github.com/Vadich007/shortener/pkg/shorter"
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
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := NewLinkService(repo, conf)

	link, err := serv.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := NewLinkService(repo, conf)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName)

	link, err := serv.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := NewLinkService(repo, conf)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)
	repo.AddLink(shorter.Shorten(originalName), originalName)

	actualShortedLink, err := serv.AddLink(originalName)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

func TestAddLinkNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := NewLinkService(repo, conf)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)

	actualShortedLink, err := serv.AddLink(originalName)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

func TestGetLinkByOriginNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := NewLinkService(repo, conf)

	link, err := serv.GetLinkByOriginal("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkByOriginExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := repository.NewInMemoryLinkRepository(conf)
	serv := NewLinkService(repo, conf)
	originalLink := "link"
	shortedLink := shorter.Shorten(originalLink)
	repo.AddLink(shortedLink, originalLink)

	link, err := serv.GetLinkByOriginal(originalLink)

	assert.Equal(t, link, "http://localhost:8080/"+shortedLink)
	assert.Equal(t, err, nil)
}

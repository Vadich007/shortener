package service

import (
	"os"
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
	"github.com/Vadich007/shortener/internal/repository/memory"
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
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	link, err := serv.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName, 0)

	link, err := serv.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)
	repo.AddLink(shorter.Shorten(originalName), originalName, 0)

	actualShortedLink, err := serv.AddLink(originalName, 0)
	assert.Equal(t, err, model.NewLinkAlreadyExistError(shorter.Shorten(originalName)))
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

func TestAddLinkNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)
	originalName := "link"
	expectedShortedLink := "http://localhost:8080/" + shorter.Shorten(originalName)

	actualShortedLink, err := serv.AddLink(originalName, 0)
	assert.Equal(t, err, nil)
	assert.Equal(t, actualShortedLink, expectedShortedLink)
}

func TestGetUserUrlsEmptyService(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	urls, err := serv.GetUserUrls(1)
	assert.NoError(t, err)
	assert.Empty(t, urls)
}

func TestGetUserUrlsService(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	originalURL := "https://example.com"
	serv.AddLink(originalURL, 7)

	urls, err := serv.GetUserUrls(7)
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, originalURL, urls[0].OriginalURL)
	assert.Equal(t, "http://localhost:8080/"+shorter.Shorten(originalURL), urls[0].ShortURL)
}

func TestGetUserUrlsServiceIsolation(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	serv.AddLink("https://user1.com", 1)
	serv.AddLink("https://user2.com", 2)

	urls1, _ := serv.GetUserUrls(1)
	assert.Len(t, urls1, 1)
	assert.Equal(t, "https://user1.com", urls1[0].OriginalURL)

	urls2, _ := serv.GetUserUrls(2)
	assert.Len(t, urls2, 1)
	assert.Equal(t, "https://user2.com", urls2[0].OriginalURL)
}

package service

import (
	"errors"
	"os"
	"testing"
	"time"

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

func TestGetLinkDeletedService(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	repo.AddLink("short", "https://example.com", 1)
	repo.DeleteURLsBatch(1, []string{"short"})

	_, err := serv.GetLink("short")
	var deletedErr *model.LinkDeletedError
	assert.True(t, errors.As(err, &deletedErr))
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

func TestDeleteURLsAsyncEventuallyDeletes(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	originalURL := "https://example.com/async"
	serv.AddLink(originalURL, 1)
	shortKey := shorter.Shorten(originalURL)

	// возвращает управление немедленно
	serv.DeleteURLs(1, []string{shortKey})

	// воркер тикает каждые 500ms; ждём до 2s
	assert.Eventually(t, func() bool {
		_, err := serv.GetLink(shortKey)
		var deletedErr *model.LinkDeletedError
		return errors.As(err, &deletedErr)
	}, 2*time.Second, 100*time.Millisecond)
}

func TestDeleteURLsAsyncExcludedFromUserUrls(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	serv.AddLink("https://keep.com", 1)
	serv.AddLink("https://delete.com", 1)
	shortKey := shorter.Shorten("https://delete.com")

	serv.DeleteURLs(1, []string{shortKey})

	assert.Eventually(t, func() bool {
		urls, _ := serv.GetUserUrls(1)
		return len(urls) == 1 && urls[0].OriginalURL == "https://keep.com"
	}, 2*time.Second, 100*time.Millisecond)
}

func TestDeleteURLsWrongOwnerNoEffect(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := memory.NewMemoryLinkRepository()
	serv := NewLinkService(repo, conf)

	originalURL := "https://owned-by-user1.com"
	serv.AddLink(originalURL, 1)
	shortKey := shorter.Shorten(originalURL)

	// user 2 пытается удалить URL пользователя 1
	serv.DeleteURLs(2, []string{shortKey})

	time.Sleep(1 * time.Second)

	link, err := serv.GetLink(shortKey)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, link)
}

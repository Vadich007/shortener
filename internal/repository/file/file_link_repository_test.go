package file

import (
	"os"
	"testing"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
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
	repo, _ := NewFileLinkRepository(conf)
	link, err := repo.GetLink("notExist")

	assert.Equal(t, link, "")
	assert.Equal(t, err.Error(), "link doesn't exist")
}

func TestGetLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName, 0)
	assert.Equal(t, err, nil)
	link, err := repo.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName, 0)
	err := repo.AddLink(shortedLink, originalName, 0)
	assert.Equal(t, err, model.NewLinkAlreadyExistError(shortedLink))
}

func TestAddLinkNotExist(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName, 0)
	assert.Equal(t, err, nil)
}

func TestPingDB(t *testing.T) {
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	assert.Equal(t, repo.PingDB(), nil)
}

func TestGetLinkDeletedFile(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	repo.AddLink("short", "original", 1)
	repo.DeleteURLsBatch(1, []string{"short"})

	_, err := repo.GetLink("short")
	var deletedErr *model.LinkDeletedError
	assert.ErrorAs(t, err, &deletedErr)
}

func TestDeleteURLsBatchFile(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	repo.AddLink("s1", "original1", 1)
	repo.AddLink("s2", "original2", 1)

	err := repo.DeleteURLsBatch(1, []string{"s1"})
	assert.NoError(t, err)

	_, err = repo.GetLink("s1")
	var deletedErr *model.LinkDeletedError
	assert.ErrorAs(t, err, &deletedErr)

	link, err := repo.GetLink("s2")
	assert.NoError(t, err)
	assert.Equal(t, "original2", link)
}

func TestDeleteURLsBatchWrongOwnerFile(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	repo.AddLink("s1", "original1", 1)

	repo.DeleteURLsBatch(2, []string{"s1"})

	link, err := repo.GetLink("s1")
	assert.NoError(t, err)
	assert.Equal(t, "original1", link)
}

func TestGetUserUrlsExcludesDeletedFile(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}
	repo, _ := NewFileLinkRepository(conf)
	repo.AddLink("s1", "original1", 1)
	repo.AddLink("s2", "original2", 1)
	repo.DeleteURLsBatch(1, []string{"s1"})

	urls, err := repo.GetUserUrls(1)
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, "original2", urls[0].OriginalURL)
}

func TestDeletedFlagPersistedToFile(t *testing.T) {
	Fixture(t)
	conf := config.Config{ServerAddress: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: storagePath}

	repo, _ := NewFileLinkRepository(conf)
	repo.AddLink("s1", "original1", 1)
	repo.DeleteURLsBatch(1, []string{"s1"})

	// reload from file
	repo2, _ := NewFileLinkRepository(conf)
	_, err := repo2.GetLink("s1")
	var deletedErr *model.LinkDeletedError
	assert.ErrorAs(t, err, &deletedErr)
}

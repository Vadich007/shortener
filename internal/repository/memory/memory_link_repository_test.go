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
	err := repo.AddLink(shortedLink, originalName, 0)
	assert.Equal(t, err, nil)
	link, err := repo.GetLink(shortedLink)

	assert.Equal(t, link, originalName)
	assert.Equal(t, err, nil)
}

func TestAddLinkExist(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	repo.AddLink(shortedLink, originalName, 0)
	err := repo.AddLink(shortedLink, originalName, 0)
	assert.Equal(t, err, model.NewLinkAlreadyExistError(shortedLink))
}

func TestAddLinkNotExist(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	originalName := "link"
	shortedLink := "short"
	err := repo.AddLink(shortedLink, originalName, 0)
	assert.Equal(t, err, nil)
}

func TestPingDB(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	assert.Equal(t, repo.PingDB(), nil)
}

func TestGetUserUrlsEmpty(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	urls, err := repo.GetUserUrls(99)
	assert.NoError(t, err)
	assert.Empty(t, urls)
}

func TestGetUserUrlsReturnsOnlyOwnUrls(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	repo.AddLink("s1", "original1", 1)
	repo.AddLink("s2", "original2", 2)
	repo.AddLink("s3", "original3", 1)

	urls, err := repo.GetUserUrls(1)
	assert.NoError(t, err)
	assert.Len(t, urls, 2)

	originals := make([]string, 0, len(urls))
	for _, u := range urls {
		originals = append(originals, u.OriginalURL)
	}
	assert.ElementsMatch(t, []string{"original1", "original3"}, originals)
}

func TestGetUserUrlsOtherUserGetsSeparateResult(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	repo.AddLink("s1", "original1", 1)
	repo.AddLink("s2", "original2", 2)

	urls, err := repo.GetUserUrls(2)
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, "original2", urls[0].OriginalURL)
	assert.Equal(t, "s2", urls[0].ShortURL)
}

func TestGetLinkDeleted(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	repo.AddLink("short", "original", 1)
	repo.DeleteURLsBatch(1, []string{"short"})

	_, err := repo.GetLink("short")
	var deletedErr *model.LinkDeletedError
	assert.ErrorAs(t, err, &deletedErr)
}

func TestDeleteURLsBatchMarksDeleted(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
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

func TestDeleteURLsBatchWrongOwnerNoEffect(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	repo.AddLink("s1", "original1", 1)

	// user 2 tries to delete user 1's URL
	repo.DeleteURLsBatch(2, []string{"s1"})

	link, err := repo.GetLink("s1")
	assert.NoError(t, err)
	assert.Equal(t, "original1", link)
}

func TestDeleteURLsBatchEmpty(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	repo.AddLink("s1", "original1", 1)

	err := repo.DeleteURLsBatch(1, []string{})
	assert.NoError(t, err)

	link, err := repo.GetLink("s1")
	assert.NoError(t, err)
	assert.Equal(t, "original1", link)
}

func TestGetUserUrlsExcludesDeleted(t *testing.T) {
	repo, _ := NewMemoryLinkRepository()
	repo.AddLink("s1", "original1", 1)
	repo.AddLink("s2", "original2", 1)
	repo.DeleteURLsBatch(1, []string{"s1"})

	urls, err := repo.GetUserUrls(1)
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, "original2", urls[0].OriginalURL)
}

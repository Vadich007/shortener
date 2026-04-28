package factory

import (
	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/repository"
	file "github.com/Vadich007/shortener/internal/repository/file"
	"github.com/Vadich007/shortener/internal/repository/memory"
	"github.com/Vadich007/shortener/internal/repository/postgres"
)

func GetRepository(conf config.Config) (repository.LinkRepository, error) {
	if conf.DatabaseDsn != "" {
		return postgres.NewPostrgesLinkRepository(conf)
	} else if conf.FileStoragePath != "" {
		return file.NewFileLinkRepository(conf)
	}

	return memory.NewMemoryLinkRepository()
}

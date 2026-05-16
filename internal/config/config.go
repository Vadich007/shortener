package config

import (
	"github.com/Vadich007/shortener/internal/config/env"
	"github.com/Vadich007/shortener/internal/config/flags"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDsn     string
	SecretKey       string
}

func GetConfig() Config {
	env := env.GetEnv()
	serverAddress := env.ServerAddress
	baseURL := env.BaseURL
	fileStoragePath := env.FileStoragePath
	databaseDsn := env.DatabaseDsn
	secretKey := env.SecretKey

	flag := flags.ProcessingFlags()

	if serverAddress == "" {
		serverAddress = flag.A
	}

	if baseURL == "" {
		baseURL = flag.B
	}

	if fileStoragePath == "" {
		fileStoragePath = flag.F
	}

	if databaseDsn == "" {
		databaseDsn = flag.D
	}

	if secretKey == "" {
		secretKey = flag.S
	}

	return Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		DatabaseDsn:     databaseDsn,
		SecretKey:       secretKey,
	}
}

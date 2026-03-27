package config

import (
	"github.com/Vadich007/shortener/internal/config/env"
	"github.com/Vadich007/shortener/internal/config/flags"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func GetConfig() Config {
	var serverAddress, baseURL, fileStoragePath string
	env := env.GetEnv()
	serverAddress = env.ServerAddress
	baseURL = env.BaseURL
	fileStoragePath = env.FileStoragePath

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

	return Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
	}
}

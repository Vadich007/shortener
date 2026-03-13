package config

import (
	"github.com/Vadich007/shortener/internal/config/env"
	"github.com/Vadich007/shortener/internal/config/flags"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func GetConfig() Config {
	var serverAddress, baseURL string
	env := env.GetEnv()
	serverAddress = env.ServerAddress
	baseURL = env.BaseURL

	flag := flags.ProcessingFlags()

	if serverAddress == "" {
		serverAddress = flag.A
	}

	if baseURL == "" {
		baseURL = flag.B
	}

	return Config{
		ServerAddress: serverAddress,
		BaseURL:       baseURL,
	}
}

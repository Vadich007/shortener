package config

import (
	"github.com/Vadich007/shortener/internal/config/env"
	"github.com/Vadich007/shortener/internal/config/flags"
)

type Config struct {
	ServerAddress string
	BaseUrl       string
}

func GetConfig() Config {
	var serverAddress, baseUrl string
	env := env.GetEnv()
	serverAddress = env.ServerAddress
	baseUrl = env.BaseUrl

	flag := flags.ProcessingFlags()

	if serverAddress == "" {
		serverAddress = flag.A
	}

	if baseUrl == "" {
		baseUrl = flag.B
	}

	return Config{
		ServerAddress: serverAddress,
		BaseUrl:       baseUrl,
	}
}

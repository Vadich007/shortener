package env

import (
	"os"
)

type Env struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDsn     string
	SecretKey       string
}

func GetEnv() Env {
	return Env{
		ServerAddress:   os.Getenv("SERVER_ADDRESS"),
		BaseURL:         os.Getenv("BASE_URL"),
		FileStoragePath: os.Getenv("FILE_STORAGE_PATH"),
		DatabaseDsn:     os.Getenv("DATABASE_DSN"),
		SecretKey:       os.Getenv("SECRET_KEY"),
	}
}

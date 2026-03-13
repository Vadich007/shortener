package env

import (
	"os"
)

type Env struct {
	ServerAddress string
	BaseURL       string
}

func GetEnv() Env {
	return Env{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		BaseURL:       os.Getenv("BASE_URL"),
	}
}

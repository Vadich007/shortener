package env

import (
	"os"
)

type Env struct {
	ServerAddress string
	BaseUrl       string
}

func GetEnv() Env {
	return Env{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		BaseUrl:       os.Getenv("BASE_URL"),
	}
}

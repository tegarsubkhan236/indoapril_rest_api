package pkg

import (
	"github.com/joho/godotenv"
	"os"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		panic("error loading .env file !")
	}
	return os.Getenv(key)
}

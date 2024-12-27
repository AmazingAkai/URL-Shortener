package utils

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	godotenv.Load(".env.local")
}

package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	CHARSET          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	SHORT_URL_LENGTH = 8
)

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateShortUrl() string {
	var builder strings.Builder
	builder.Grow(SHORT_URL_LENGTH)

	for i := 0; i < SHORT_URL_LENGTH; i++ {
		randomIndex := randomGenerator.Intn(len(CHARSET))
		builder.WriteByte(CHARSET[randomIndex])
	}

	return builder.String()
}

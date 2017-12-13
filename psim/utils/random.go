package utils

import (
	"math/rand"
	"time"
)

var (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randomString(n int, source string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func GenerateToken() string {
	return randomString(20, letters)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

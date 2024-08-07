package router

import (
	"crypto/rand"
	"time"
)

func generateshortkey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixMicro())
	shortkey := make([]byte, keyLength)
	for i := range shortkey {
		shortkey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortkey)
	
}
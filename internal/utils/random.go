package utils

import (
	"crypto/rand"
	"io"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() (string, error) {
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}

	for i := 0; i < 12; i++ {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}

	return string(b), nil
}

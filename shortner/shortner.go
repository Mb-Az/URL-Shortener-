package shortner

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

const urlLen = 8

func RandomString(length int) (string, error) {

	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("Error in randomstring(): %w", err)
		}

		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}

func CreateShortUrl(url string, length int) (string, error) {
	if len(url) == 0 {
		return "", fmt.Errorf("url is empty")
	}
	if length == 0 {
		length = urlLen
	}
	id, err := RandomString(length)
	if err != nil {
		return "", err
	}
	// fmt.Println("Generated UUID:")
	// fmt.Printf("%s", id)
	return id, nil
}

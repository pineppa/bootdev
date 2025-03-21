package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// Generates a random 256-bit hex-encoded string
func MakeRefreshToken() (string, error) {
	var b = make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error:  %s\n", err)
		return "", err
	}
	return hex.EncodeToString(b), nil
}

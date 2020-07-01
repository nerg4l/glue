package glue

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomKey() (string, error) {
	ra := make([]byte, 32)
	if _, err := rand.Read(ra); err != nil {
		return "", err
	}
	return "base64:" + base64.StdEncoding.EncodeToString(ra), nil
}

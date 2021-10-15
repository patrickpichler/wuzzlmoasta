package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(len int) string {
	bytes := make([]byte, len)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

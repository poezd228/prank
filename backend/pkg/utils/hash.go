package utils

import (
	"crypto/sha256"
	"fmt"
)

func HashSha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

package logic

import (
	"crypto/sha256"
	"encoding/hex"
)

func hashing(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

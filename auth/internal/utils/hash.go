package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashPassword(s, salt string) string {
	return hex.EncodeToString(Sha1(s + salt))
}

func Sha1(s string) []byte {
	h := sha1.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

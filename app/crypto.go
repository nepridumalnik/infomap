package app

import (
	"crypto/sha512"
	"encoding/hex"
)

func Sha512(text string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(text + salt))
	return hex.EncodeToString(algorithm.Sum(nil))
}

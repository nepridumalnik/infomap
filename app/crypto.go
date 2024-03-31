package app

import (
	"crypto/sha512"
	"encoding/hex"
)

const (
	salt = "qweASD123zxcASDqwe=-124-980r7102370129edqwu98"
)

func Sha512(text string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(text + salt))
	return hex.EncodeToString(algorithm.Sum(nil))
}

package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetToken(email string) string {
	hashString := email + GetEnv("DECODING_SECRET", "")
	return GetMd5(hashString)
}

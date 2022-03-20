package helpers

import (
	"math/rand"
	"strings"
	"time"
)

func GetStaticUrl(hash string) string {
	return "https://voiceline.me/" + hash
}

func FormatURL(url string) string {
	s3 := strings.Split(url, "/")
	if len(s3) > 1 {
		return GetStaticUrl(s3[len(s3)-1])
	}
	return s3[len(s3)-1]
}

func GetStaticProfileUrl(hash string) string {
	return "https://voiceline.me/p/" + hash
}

func GetUserUrl(email, workspaceId string) string {
	hashString := email + GetEnv("DECODING_SECRET", "")
	hashedString := GetMd5(hashString)
	domain := GetEnv("APP_ENVIRONMENT", "staging")
	return "https://" + domain + ".getvoiceline.com/profile/" + workspaceId + "/" + hashedString
}

func GetVLStaticUrl(hash string) string {
	env := GetEnv("APP_ENVIRONMENT", "")
	if env == "app" {
		return "https://voiceline.me/vl/" + hash
	}
	return "https://voiceline.me/staging/vl/" + hash
}

func GetUniqueHash(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

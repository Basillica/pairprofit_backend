package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func ArrayToString(a []int, delim string) (res *string) {
	val := strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	res = &val
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
	return
}

func FloatToString(a []float64, delim string) (res *string) {
	val := strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	res = &val
	return
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func FormatExpTime(expirationData time.Time) string {
	t := strings.ReplaceAll(expirationData.Format(time.RFC3339), ":", "-")
	time := strings.Split(t, "+")[0]
	return time
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetHashedString(clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return secretHash
}

package helpers

import (
	"fmt"
	"strings"
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

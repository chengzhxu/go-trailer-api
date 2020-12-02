package util

import "time"

func ExistIntElement(element int64, array []int64) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}

	return false
}

func GetNowTimeStamp() int {
	return int(time.Now().Unix())
}
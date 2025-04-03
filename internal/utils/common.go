package utils

import "time"

func GetUnixTimestamp() int64 {
	return time.Now().Unix()
}

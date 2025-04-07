package utils

import (
	"fmt"
	"os"
	"time"
)

func GetUnixTimestamp() int64 {
	return time.Now().Unix()
}

func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error file read %w", err)
	}
	return string(data), nil
}

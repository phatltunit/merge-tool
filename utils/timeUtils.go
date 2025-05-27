package utils

import (
	"time"
)

func GetCurrentTime(format string) string {
	return time.Now().Format(format)
}

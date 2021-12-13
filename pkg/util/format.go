package util

import "time"

func TimeHHMMSS() string {
	return time.Now().Format("15:04:05")
}

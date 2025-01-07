package utils

import "time"

func GetDateForFilename() string {
	now := time.Now()
	dateString := now.Format("2006-01-02_15-04-05")
	return dateString
}

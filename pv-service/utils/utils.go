package utils

import (
	"time"
)

const decade = 315532800
const Hour = 3600

func ConvertUnixToTimeStamp(unixtime uint32) time.Time {
	unixtime += decade
	timestamp := time.Unix(int64(unixtime), 0)
	return timestamp
}

func ConvertTimestampToUnix(timestamp *time.Time) uint32 {
	unixtime := timestamp.Unix()
	unixtime -= decade
	return uint32(unixtime)
}

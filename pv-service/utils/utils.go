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
	unixTime := timestamp.Unix()
	unixTime -= decade
	return uint32(unixTime)
}

func GetAverage(arr []uint32) uint32 {
	sum := uint32(0)
	for _, n := range arr {
		sum += n
	}
	return sum / uint32(len(arr))
}

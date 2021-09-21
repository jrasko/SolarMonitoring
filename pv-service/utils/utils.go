package utils

import (
	"time"
)

const decade = 315532800
const Hour = 3600

var UTC1 = time.FixedZone("UTC+01", 1)

func ConvertUnixToTimeStamp(unixtime uint32) time.Time {
	unixtime += decade
	return time.Unix(int64(unixtime), 0).In(UTC1)
}

func ConvertTimestampToUnix(timestamp *time.Time) uint32 {
	unixTime := timestamp.Unix()
	unixTime -= decade
	return uint32(unixTime)
}

func GetAverageEnergy(arr []uint32) uint32 {
	sum := uint32(0)
	for _, n := range arr {
		sum += n
	}
	return sum / uint32(len(arr))
}
func GetAverageTime(arr []int64) int64 {
	sum := int64(0)
	for _, n := range arr {
		sum += n
	}
	return sum / int64(len(arr))
}

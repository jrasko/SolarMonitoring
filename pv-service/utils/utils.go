package utils

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

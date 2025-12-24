package model

type Number interface {
	int64 | uint32 | int
}

func GetAverage[T Number](arr []T) T {
	var sum T
	for _, n := range arr {
		sum += n
	}
	return sum / T(len(arr))
}

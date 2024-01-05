package model

import (
	"fmt"
)

type MinuteDataSet []MinuteData

func (s *MinuteDataSet) CurrentAverage(intervalTime int) error {
	set := *s
	if intervalTime == 1 {
		return nil
	}
	if intervalTime > len(set) {
		return fmt.Errorf("intervalTime is larger than timespan")
	}

	arrays := len(set[0].DcI)
	averagedDCI := make([][]int, arrays)
	for i := range averagedDCI {
		averagedDCI[i] = make([]int, 0, intervalTime)
	}
	for i, minuteData := range set {
		for j, dc := range minuteData.DcI {
			if i >= intervalTime {
				averagedDCI[j][i%intervalTime] = dc
			} else {
				averagedDCI[j] = append(averagedDCI[j], dc)
			}
			set[i].DcI[j] = GetAverage(averagedDCI[j])
		}
	}
	*s = set
	return nil
}

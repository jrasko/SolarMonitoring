package model

import (
	"fmt"
)

type DailyDataSet []DailyData

func (s *DailyDataSet) Average(timeInterval int, energyInterval int) error {
	err := s.timeAverage(timeInterval)
	if err != nil {
		return err
	}
	return s.energyAverage(energyInterval)
}

func (s *DailyDataSet) energyAverage(intervalTime int) error {
	if intervalTime <= 1 {
		return nil
	}

	dataset := *s
	if intervalTime > len(dataset) {
		return fmt.Errorf("intervalTime is larger than timespan")
	}

	var averagedEnergy []int
	for i, dailyData := range dataset {
		if i >= intervalTime {
			averagedEnergy[i%intervalTime] = dailyData.ProducedEnergy
		} else {
			averagedEnergy = append(averagedEnergy, dailyData.ProducedEnergy)
		}
		dataset[i].ProducedEnergy = GetAverage(averagedEnergy)
	}
	*s = dataset
	return nil
}

func (s *DailyDataSet) timeAverage(intervalTime int) error {
	if intervalTime <= 1 {
		return nil
	}

	dataset := *s
	if intervalTime > len(dataset) {
		return fmt.Errorf("intervalTime is larger than timespan")
	}
	var averagedTime []int64
	for i, dailyData := range dataset {
		currentTimeStamp := dailyData.StartupTime.ToSeconds()
		if i >= intervalTime {
			averagedTime[i%intervalTime] = currentTimeStamp
		} else {
			averagedTime = append(averagedTime, currentTimeStamp)
		}
		dailyData.StartupTime = ClockFromSeconds(GetAverage(averagedTime))
	}
	*s = dataset
	return nil
}

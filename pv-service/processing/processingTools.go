package processing

import (
	"fmt"
	"pv-service/entities/dao"
	"pv-service/entities/dto"
	"pv-service/graph/model"
	"pv-service/utils"
	"time"
)

type processingData struct {
	date   dto.Date
	time   dto.TimeOfDay
	totalE uint16
}

func getDate(t time.Time) dto.Date {
	return dto.Date{
		Day:   uint8(t.Day()),
		Month: uint8(t.Month()),
		Year:  uint16(t.Year()),
	}
}

func averagizeEnergy(intervalTime uint32, dailyDataArray []*model.DailyData) ([]*model.DailyData, error) {
	if intervalTime == 1 {
		return dailyDataArray, nil
	}
	if intervalTime > uint32(len(dailyDataArray)) {
		return nil, fmt.Errorf("intervalTime is larger than timespan")
	}
	var averagedEnergy []uint32
	for i, dailyData := range dailyDataArray {
		if uint32(i) >= intervalTime {
			averagedEnergy[uint32(i)%intervalTime] = dailyData.ProducedEnergy
		} else {
			averagedEnergy = append(averagedEnergy, dailyData.ProducedEnergy)
		}
		dailyData.ProducedEnergy = utils.GetAverageEnergy(averagedEnergy)
	}
	return dailyDataArray, nil
}

func averagizeTime(intervalTime uint32, dailyDataArray []*model.DailyData) ([]*model.DailyData, error) {
	if intervalTime == 1 {
		return dailyDataArray, nil
	}
	if intervalTime > uint32(len(dailyDataArray)) {
		return nil, fmt.Errorf("intervalTime is larger than timespan")
	}
	var averagedTime []int64
	for i, dailyData := range dailyDataArray {
		currentTimeStamp := time.Date(
			1970, 1, 1,
			int(dailyData.StartupTime.Hours), int(dailyData.StartupTime.Minutes), int(dailyData.StartupTime.Seconds),
			0, utils.UTC1).Unix()
		if uint32(i) >= intervalTime {
			averagedTime[uint32(i)%intervalTime] = currentTimeStamp
		} else {
			averagedTime = append(averagedTime, currentTimeStamp)
		}
		aT := time.Unix(utils.GetAverageTime(averagedTime), 0)
		dailyData.StartupTime = dto.TimeOfDay{
			Hours:   uint8(aT.Hour()),
			Minutes: uint8(aT.Minute()),
			Seconds: uint8(aT.Second()),
		}
	}
	return dailyDataArray, nil
}

func mapDataAndRemoveDuplicates(data *[]dao.PVData) *[]processingData {
	mappedData := make([]*processingData, 0, len(*data))
	for _, pvData := range *data {
		timeOfDatapoint := utils.ConvertUnixToTimeStamp(pvData.Time)
		mappedData = append(mappedData, &processingData{
			date: getDate(timeOfDatapoint),
			time: dto.TimeOfDay{
				Seconds: uint8(timeOfDatapoint.Second()),
				Minutes: uint8(timeOfDatapoint.Minute()),
				Hours:   uint8(timeOfDatapoint.Hour()),
			},
			totalE: pvData.TotalE,
		})
	}
	currentIndex := 0
	var purgedDataArray = []processingData{*mappedData[0]}
	for _, dataPoint := range mappedData {
		if dataPoint.date != purgedDataArray[currentIndex].date {
			purgedDataArray = append(purgedDataArray, *dataPoint)
			currentIndex++
		}
	}

	return &purgedDataArray
}

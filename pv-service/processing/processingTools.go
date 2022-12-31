package processing

import (
	"fmt"
	"pv-service/entities/dao"
	"pv-service/entities/dto"
	"pv-service/graph/model"
	"pv-service/utils"
	"time"
)

var (
	inverterSwitch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

const (
	inverterSwitchTotalE = 107758
)

type processingData struct {
	date   dto.Date
	time   dto.TimeOfDay
	totalE uint32
}

func getDate(t dto.PVTime) dto.Date {
	return dto.Date{
		Day:   uint8(t.ToTime().Day()),
		Month: uint8(t.ToTime().Month()),
		Year:  uint16(t.ToTime().Year()),
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

func averagizeCurrent(intervalTime uint32, minuteDataArray []*model.MinuteDataOfDay) ([]*model.MinuteDataOfDay, error) {
	if intervalTime == 1 {
		return minuteDataArray, nil
	}
	if intervalTime > uint32(len(minuteDataArray)) {
		return nil, fmt.Errorf("intervalTime is larger than timespan")
	}
	var averagedDC1I []uint32
	var averagedDC2I []uint32
	var averagedDC3I []uint32
	for i, minuteData := range minuteDataArray {
		if uint32(i) >= intervalTime {
			averagedDC1I[uint32(i)%intervalTime] = minuteData.Dc1i
			averagedDC2I[uint32(i)%intervalTime] = minuteData.Dc2i
			averagedDC3I[uint32(i)%intervalTime] = minuteData.Dc3i
		} else {
			averagedDC1I = append(averagedDC1I, minuteData.Dc1i)
			averagedDC2I = append(averagedDC2I, minuteData.Dc2i)
			averagedDC3I = append(averagedDC3I, minuteData.Dc3i)
		}
		minuteData.Dc1i = utils.GetAverageEnergy(averagedDC1I)
		minuteData.Dc2i = utils.GetAverageEnergy(averagedDC2I)
		minuteData.Dc3i = utils.GetAverageEnergy(averagedDC3I)
	}
	return minuteDataArray, nil
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
			0, dto.UTC1).Unix()
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
	// Map DB Datasets to Model
	for _, pvData := range *data {
		timeOfDatapoint := dto.PVTimeFromUnix(pvData.Time).ToTime()
		mappedData = append(mappedData, &processingData{
			time: dto.TimeOfDay{
				Seconds: uint8(timeOfDatapoint.Second()),
				Minutes: uint8(timeOfDatapoint.Minute()),
				Hours:   uint8(timeOfDatapoint.Hour()),
			},
			date:   getDate(dto.PVTimeFromTime(timeOfDatapoint)),
			totalE: pvData.TotalE,
		})
	}
	// Purge Datasets on the same Day
	currentIndex := 0
	purgedDataArray := []processingData{*mappedData[0]}
	for _, dataPoint := range mappedData {
		if dataPoint.date != purgedDataArray[currentIndex].date {
			purgedDataArray = append(purgedDataArray, *dataPoint)
			currentIndex++
		}
	}
	// Transform Total Energy after inverter replacement
	if purgedDataArray[len(purgedDataArray)-1].date.Year >= uint16(inverterSwitch.Year()) {
		for i, p := range purgedDataArray {
			if p.date.Year >= uint16(inverterSwitch.Year()) {
				purgedDataArray[i].totalE += inverterSwitchTotalE
			}
		}
	}
	// Remove startup timestamp before inverter replacement
	if purgedDataArray[0].date.Year < 2020 {
		for i := range purgedDataArray {
			purgedDataArray[i].time = dto.TimeOfDay{}
		}
	}
	return &purgedDataArray
}

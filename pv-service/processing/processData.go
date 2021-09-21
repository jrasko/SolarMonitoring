package processing

import (
	"fmt"
	"pv-service/database"
	"pv-service/entities/dao"
	"pv-service/entities/dto"
	"pv-service/graph/model"
	"pv-service/utils"
	"time"
)

type Processor interface {
	GetDailyDataBetweenDates(start *time.Time, end *time.Time, energyInterval uint32, startupInterval uint32) ([]*model.DailyData, error)
	GetRawDataBetweenDates(begin *time.Time, end *time.Time) ([]*model.RawData, error)
}

type processor struct {
	db database.DBConnection
}

type processingData struct {
	date   dto.Date
	time   dto.TimeOfDay
	totalE uint16
}

func GetProcessor() Processor {
	return &processor{
		db: database.GetDBConnection(),
	}
}

func (p *processor) GetDailyDataBetweenDates(
	start *time.Time, end *time.Time, energyInterval uint32, startupInterval uint32) ([]*model.DailyData, error) {
	data, err := p.db.GetDailyDataBetweenStartAndEndTime(
		utils.ConvertTimestampToUnix(start),
		utils.ConvertTimestampToUnix(end)+12*utils.Hour,
	)
	if err != nil {
		return nil, err
	}

	mappedData := mapDataAndRemoveDuplicates(data)

	lastE := uint16(0)
	var (
		lastTime dto.TimeOfDay
	)
	dailyDataArray := make([]*model.DailyData, 0, len(*mappedData)-1)
	for i, pvData := range *mappedData {
		if i == 0 {
			lastE = pvData.totalE
			lastTime = pvData.time
			continue
		}
		date := pvData.date
		date.Yesterday()
		dailyDataArray = append(dailyDataArray, &model.DailyData{
			Date:           date,
			StartupTime:    lastTime,
			ProducedEnergy: uint32(pvData.totalE - lastE),
		})
		lastE = pvData.totalE
		lastTime = pvData.time
	}

	dailyDataArray, err = averagizeTime(energyInterval, dailyDataArray)
	if err != nil {
		return nil, err
	}
	dailyDataArray, err = averagizeEnergy(startupInterval, dailyDataArray)
	if err != nil {
		return nil, err
	}
	return dailyDataArray, nil
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
	var mappedData []*processingData
	for _, pvData := range *data {
		timeOfDatapoint := utils.ConvertUnixToTimeStamp(pvData.Time)
		mappedData = append(mappedData, &processingData{
			date: dto.Date{
				Day:   uint8(timeOfDatapoint.Day()),
				Month: uint8(timeOfDatapoint.Month()),
				Year:  uint16(timeOfDatapoint.Year()),
			},
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

func (p *processor) GetRawDataBetweenDates(start *time.Time, end *time.Time) ([]*model.RawData, error) {
	startTime := utils.ConvertTimestampToUnix(start)
	endTime := utils.ConvertTimestampToUnix(end)
	data, err := p.db.GetNonDailyDataBetweenStartAndEndTime(startTime, endTime)
	if err != nil {
		return nil, err
	}
	var rawDataArray []*model.RawData
	for _, pvData := range *data {
		rawDataArray = append(rawDataArray, &model.RawData{
			Time:   pvData.Time,
			Dc1U:   uint32(pvData.Dc1U),
			Dc1I:   uint32(pvData.Dc1I),
			Dc1P:   uint32(pvData.Dc1P),
			Dc1T:   pvData.Dc1T,
			Dc1S:   pvData.Dc1S,
			Dc2U:   uint32(pvData.Dc2U),
			Dc2I:   uint32(pvData.Dc2I),
			Dc2P:   uint32(pvData.Dc2P),
			Dc2T:   pvData.Dc2T,
			Dc2S:   pvData.Dc2S,
			Dc3U:   uint32(pvData.Dc3U),
			Dc3I:   uint32(pvData.Dc3I),
			Dc3P:   uint32(pvData.Dc3P),
			Dc3T:   pvData.Dc3T,
			Dc3S:   pvData.Dc3S,
			Ac1U:   uint32(pvData.Ac1U),
			Ac1I:   uint32(pvData.Ac1I),
			Ac1P:   int32(pvData.Ac1P),
			Ac1T:   pvData.Ac1T,
			Ac2U:   uint32(pvData.Ac1U),
			Ac2I:   uint32(pvData.Ac2I),
			Ac2P:   int32(pvData.Ac2P),
			Ac2T:   pvData.Ac2T,
			Ac3U:   uint32(pvData.Ac3U),
			Ac3I:   uint32(pvData.Ac3I),
			Ac3P:   int32(pvData.Ac3P),
			Ac3T:   pvData.Ac3T,
			AcF:    float64(pvData.AcF),
			FcI:    int32(pvData.FcI),
			Ain1:   int32(pvData.Ain1),
			Ain2:   int32(pvData.Ain2),
			Ain3:   int32(pvData.Ain3),
			AcS:    uint32(pvData.AcS),
			Err:    int32(pvData.Err),
			EnsErr: uint32(pvData.EnsErr),
			Event:  pvData.Event,
		})
	}
	return rawDataArray, nil
}

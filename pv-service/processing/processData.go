package processing

import (
	"pv-service/database"
	"pv-service/entities/dto"
	"pv-service/graph/model"
	"pv-service/utils"
	"time"
)

type Processor interface {
	GetDailyData(start *time.Time, end *time.Time, energyInterval uint32, startupInterval uint32) ([]*model.DailyData, error)
	GetMinuteDataOfDay(start *time.Time, end *time.Time, currentInterval uint32) ([]*model.MinuteDataOfDay, error)
	GetRawDataBetweenDates(start *time.Time, end *time.Time) ([]*model.RawData, error)
}

const averageDataPerDay = 60

type processor struct {
	db database.DBConnection
}

func GetProcessor() Processor {
	return &processor{
		db: database.GetDBConnection(),
	}
}
func (p *processor) GetMinuteDataOfDay(start *time.Time, end *time.Time, currentInterval uint32) ([]*model.MinuteDataOfDay, error) {
	rawData, err := p.db.GetNonDailyDataBetweenStartAndEndTime(
		utils.ConvertTimestampToUnix(start),
		utils.ConvertTimestampToUnix(end))
	if err != nil {
		return nil, err
	}
	processedData := make([]*model.MinuteDataOfDay, 0, len(*rawData)/averageDataPerDay)
	lastDate := getDate(utils.ConvertUnixToTimeStamp((*rawData)[0].Time))
	dcI := [3]uint32{0, 0, 0}
	adding := 0
	for _, data := range *rawData {
		thisDate := getDate(utils.ConvertUnixToTimeStamp(data.Time))
		if thisDate != lastDate && adding > 0 {
			processedData = append(processedData, &model.MinuteDataOfDay{
				Date: lastDate,
				Dc1i: dcI[0] / uint32(adding),
				Dc2i: dcI[1] / uint32(adding),
				Dc3i: dcI[2] / uint32(adding),
			})
			adding = 0
			dcI = [3]uint32{0, 0, 0}
		}
		lastDate = thisDate
		dcI[0] += uint32(data.Dc1I)
		dcI[1] += uint32(data.Dc2I)
		dcI[2] += uint32(data.Dc3I)
		adding++
	}
	processedData = append(processedData, &model.MinuteDataOfDay{
		Date: lastDate,
		Dc1i: dcI[0] / uint32(adding),
		Dc2i: dcI[1] / uint32(adding),
		Dc3i: dcI[2] / uint32(adding),
	})
	averagizedData, err := averagizeCurrent(currentInterval, processedData)
	if err != nil {
		return nil, err
	}
	return averagizedData, nil
}

func (p *processor) GetDailyData(
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
		date := pvData.date.Yesterday()
		dailyDataArray = append(dailyDataArray, &model.DailyData{
			Date:           date,
			StartupTime:    lastTime,
			ProducedEnergy: uint32(pvData.totalE - lastE),
		})
		lastE = pvData.totalE
		lastTime = pvData.time
	}

	dailyDataArray, err = averagizeTime(startupInterval, dailyDataArray)
	if err != nil {
		return nil, err
	}
	dailyDataArray, err = averagizeEnergy(energyInterval, dailyDataArray)
	if err != nil {
		return nil, err
	}
	return dailyDataArray, nil
}

func (p *processor) GetRawDataBetweenDates(start *time.Time, end *time.Time) ([]*model.RawData, error) {
	startTime := utils.ConvertTimestampToUnix(start)
	endTime := utils.ConvertTimestampToUnix(end)
	data, err := p.db.GetNonDailyDataBetweenStartAndEndTime(startTime, endTime)
	if err != nil {
		return nil, err
	}
	var rawDataArray = make([]*model.RawData, 0, len(*data))
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

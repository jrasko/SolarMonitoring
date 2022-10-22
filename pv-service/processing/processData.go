package processing

import (
	"context"
	"pv-service/database"
	"pv-service/entities/dto"
	"pv-service/graph/model"
	"pv-service/utils"
	"time"
)

const averageDataPerDay = 60

type Processor struct {
	db *database.DBConnection
}

func GetProcessor() *Processor {
	return &Processor{
		db: database.GetDBConnection(),
	}
}
func (p *Processor) GetMinuteDataOfDay(ctx context.Context, start *time.Time, end *time.Time, currentInterval uint32) ([]*model.MinuteDataOfDay, error) {
	rawData, err := p.db.GetNonDailyDataBetweenStartAndEndTime(
		ctx,
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

func (p *Processor) GetDailyData(ctx context.Context, start *time.Time, end *time.Time, energyInterval uint32, startupInterval uint32) ([]*model.DailyData, error) {
	data, err := p.db.GetDailyDataBetweenStartAndEndTime(
		ctx,
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
			Date:             date,
			StartupTime:      lastTime,
			ProducedEnergy:   uint32(pvData.totalE - lastE),
			CumulativeEnergy: uint32(pvData.totalE),
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

func (p *Processor) GetRawDataBetweenDates(ctx context.Context, start *time.Time, end *time.Time) ([]*model.RawData, error) {
	startTime := utils.ConvertTimestampToUnix(start)
	endTime := utils.ConvertTimestampToUnix(end)
	data, err := p.db.GetNonDailyDataBetweenStartAndEndTime(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}
	var rawDataArray = make([]*model.RawData, 0, len(*data))
	for _, pvData := range *data {
		pvModel := pvData.ToModel()
		rawDataArray = append(rawDataArray, &pvModel)
	}
	return rawDataArray, nil
}

func (p *Processor) GetZappiDataBetweenDates(ctx context.Context, begin *time.Time, end *time.Time) ([]*model.ZappiData, error) {
	data, err := p.db.GetZappiDataBetweenStartAndEnddate(ctx, begin, end)
	if err != nil {
		return nil, err
	}
	var zappiDataArray []*model.ZappiData
	for _, zappiData := range *data {
		zModel := zappiData.ToModel()
		zappiDataArray = append(zappiDataArray, &zModel)
	}
	return zappiDataArray, nil
}

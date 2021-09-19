package processing

import (
	"pv-service/database"
	"pv-service/graph/model"
	"pv-service/utils"
	"time"
)

type Processor interface {
	GetDailyDataBetweenDates(start *time.Time, end *time.Time) ([]*model.DailyData, error)
	GetRawDataBetweenDates(begin *time.Time, end *time.Time) ([]*model.RawData, error)
}

type processor struct {
	db database.DBConnection
}

func GetProcessor() Processor {
	return &processor{
		db: database.GetDBConnection(),
	}
}

func (p *processor) GetDailyDataBetweenDates(start *time.Time, end *time.Time) ([]*model.DailyData, error) {
	startTime := utils.ConvertTimestampToUnix(start)
	endTime := utils.ConvertTimestampToUnix(end)
	data, err := p.db.GetDailyDataBetweenStartAndEndTime(startTime, endTime+12*utils.Hour)
	if err != nil {
		return nil, err
	}

	lastE := uint16(0)
	lastT := uint32(0)
	var dailyDataArray []*model.DailyData
	for i, pvData := range *data {
		if i == 0 {
			lastE = pvData.TotalE
			lastT = pvData.Time
			continue
		}
		d := utils.ConvertUnixToTimeStamp(pvData.Time - 12*utils.Hour)
		dailyDataArray = append(dailyDataArray, &model.DailyData{
			Date:           time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC),
			StartUpTime:    utils.ConvertUnixToTimeStamp(lastT),
			ProducedEnergy: uint32(pvData.TotalE - lastE),
		})
		lastE = pvData.TotalE
		lastT = pvData.Time
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

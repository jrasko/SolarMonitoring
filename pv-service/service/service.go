package processing

import (
	"context"
	"pv-service/database"
	"pv-service/graph/model"
	"pv-service/processing"
)

type Service struct {
	db database.DBConnection
}

func GetService() (Service, error) {
	connection, err := database.GetDBConnection()
	if err != nil {
		return Service{}, err
	}
	return Service{db: connection}, nil
}

func (p Service) GetMinuteDataOfDay(ctx context.Context, start *model.Time, end *model.Time, currentInterval int) (model.MinuteDataSet, error) {
	unixStart := start.AsStartUnix()
	unixEnd := end.AsEndUnix()

	rawData, err := p.db.GetNonDailyData(ctx, unixStart, unixEnd)
	if err != nil || len(rawData) == 0 {
		return nil, err
	}
	minuteDataArray := processing.MinuteData(rawData)
	err = minuteDataArray.CurrentAverage(currentInterval)
	if err != nil {
		return nil, err
	}
	return minuteDataArray, nil
}

func (p Service) GetDailyData(ctx context.Context, start *model.Time, end *model.Time, energyInterval int, startupInterval int) (model.DailyDataSet, error) {
	unixStart := start.AsStartUnix()
	unixEnd := end.AsEndUnix()

	rawData, err := p.db.GetDailyData(ctx, unixStart, unixEnd)
	if err != nil || len(rawData) == 0 {
		return nil, err
	}

	dailyDataArray := processing.DailyData(rawData)
	err = dailyDataArray.Average(startupInterval, energyInterval)
	if err != nil {
		return nil, err
	}

	return dailyDataArray, err
}

func (p Service) GetZappiDataBetweenDates(ctx context.Context, begin *model.Time, end *model.Time) ([]model.ZappiData, error) {
	data, err := p.db.GetZappiData(ctx, begin.Time(), end.Time())
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var zappiDataArray []model.ZappiData
	for _, zappiData := range data {
		zModel := zappiData.ToModel()
		zappiDataArray = append(zappiDataArray, zModel)
	}
	return zappiDataArray, nil
}

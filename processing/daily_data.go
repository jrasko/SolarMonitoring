package processing

import (
	"pv-service/database"
	"pv-service/graph/model"
	"time"
)

var (
	inverterSwitch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

const (
	inverterSwitchTotalE = 107758
)

func DailyData(rawData []database.DailyData) model.DailyDataSet {
	internal := mapToInternal(rawData)
	// Purge Datasets on the same Day
	internal.purgeDuplicates()

	// apply inverter
	internal.applyInverter()

	return internal.calculateDailyData()
}

type dailyData struct {
	Time   model.Time
	Date   model.Date
	Clock  model.Clock
	TotalE int
}

type dailyDataSet []dailyData

func (d *dailyDataSet) purgeDuplicates() {
	data := *d
	currentIndex := 0
	purgedDataArray := dailyDataSet{data[0]}
	for _, dataPoint := range data {
		if dataPoint.Date != purgedDataArray[currentIndex].Date {
			purgedDataArray = append(purgedDataArray, dataPoint)
			currentIndex++
		}
	}
	*d = purgedDataArray
}

func (d *dailyDataSet) applyInverter() {
	dataSet := *d
	// Transform Total Energy after inverter replacement
	firstDay := dataSet[0]
	lastDay := dataSet[len(dataSet)-1]
	if lastDay.Time.Time().After(inverterSwitch) {
		for i, p := range dataSet {
			if p.Time.Time().After(inverterSwitch) {
				dataSet[i].TotalE += inverterSwitchTotalE
			}
		}
	}

	// Remove startup timestamp before inverter replacement
	if firstDay.Time.Time().Before(inverterSwitch) {
		for i, p := range dataSet {
			if p.Time.Time().Before(inverterSwitch) {
				dataSet[i].Clock = model.DefaultClock()
			} else {
				break
			}
		}
	}
	*d = dataSet
}

func (d dailyDataSet) calculateDailyData() model.DailyDataSet {
	// set total energy and startup time
	// the TotalE of a point always refers to the day before
	dailyDataArray := make(model.DailyDataSet, len(d)-1)
	lastTime := d[0].Clock
	lastE := d[0].TotalE
	for i, pvData := range d[1:] {
		date := pvData.Date.Yesterday()
		dailyDataArray[i] = model.DailyData{
			Date:             date,
			StartupTime:      lastTime,
			ProducedEnergy:   pvData.TotalE - lastE,
			CumulativeEnergy: pvData.TotalE,
		}
		lastE = pvData.TotalE
		lastTime = pvData.Clock
	}
	return dailyDataArray
}

func mapToInternal(data []database.DailyData) dailyDataSet {
	internal := make([]dailyData, 0, len(data))
	// Map DB Datasets to Model
	for _, pvData := range data {
		timeOfDatapoint := model.PVTimeFromUnix(pvData.Time)
		internal = append(internal, dailyData{
			Time:   timeOfDatapoint,
			Clock:  timeOfDatapoint.Clock(),
			Date:   timeOfDatapoint.GetDate(),
			TotalE: pvData.TotalE,
		})
	}
	return internal
}

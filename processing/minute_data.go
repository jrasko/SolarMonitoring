package processing

import (
	"pv-service/database"
	"pv-service/graph/model"
	"slices"
)

func MinuteData(rawData []database.MinuteData) model.MinuteDataSet {
	dataPerDay := sumPerDay(rawData)

	return dataPerDay.toAveragedModel()
}

type summedCurrents struct {
	SolarArraySum []int
	Entries       int
}

func (c summedCurrents) averages() []int {
	for i := range c.SolarArraySum {
		c.SolarArraySum[i] /= c.Entries
	}
	return c.SolarArraySum
}

type dailyCurrents map[model.Date]summedCurrents

func (c dailyCurrents) toAveragedModel() model.MinuteDataSet {
	external := make(model.MinuteDataSet, 0, len(c))

	// map to model, retreive averages
	for date, currents := range c {
		external = append(external, model.MinuteData{
			Date: date,
			DcI:  currents.averages(),
		})
	}

	slices.SortFunc(external, func(a, b model.MinuteData) int {
		return a.Date.Compare(b.Date)
	})
	return external
}

func sumPerDay(rawData []database.MinuteData) dailyCurrents {
	dataPerDay := dailyCurrents{}

	// sum all currents per solar array up
	for _, raw := range rawData {
		date := model.PVTimeFromUnix(raw.Time).GetDate()
		if _, ok := dataPerDay[date]; !ok {
			dataPerDay[date] = summedCurrents{
				SolarArraySum: make([]int, len(raw.DcI)),
			}
		}
		data := dataPerDay[date]
		for j, dc := range raw.DcI {
			data.SolarArraySum[j] += dc
		}
		data.Entries++
		dataPerDay[date] = data
	}
	return dataPerDay
}

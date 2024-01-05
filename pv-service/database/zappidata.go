package database

import (
	"pv-service/graph/model"
	"time"
)

type ZappiData struct {
	ZappiSN        int
	PluggedIn      time.Time
	Unplugged      time.Time
	ChargeDuration int
	Electricity    float64
}

func (ZappiData) TableName() string {
	return "MyEnergi"
}

func (z ZappiData) ToModel() model.ZappiData {
	return model.ZappiData{
		ZappiSn:        z.ZappiSN,
		PluggedIn:      model.Time(z.PluggedIn),
		Unplugged:      model.Time(z.Unplugged),
		ChargeDuration: z.ChargeDuration,
		Electricity:    z.Electricity,
	}
}

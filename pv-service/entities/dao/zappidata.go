package dao

import (
	"pv-service/entities/dto"
	"pv-service/graph/model"
	"time"
)

type ZappiData struct {
	ZappiSN        int32
	PluggedIn      time.Time
	Unplugged      time.Time
	ChargeDuration int32
	Electricity    float64
}

func (ZappiData) TableName() string {
	return "MyEnergi"
}

func (z ZappiData) ToModel() model.ZappiData {
	return model.ZappiData{
		ZappiSn:        z.ZappiSN,
		PluggedIn:      dto.PVTimeFromTime(z.PluggedIn),
		Unplugged:      dto.PVTimeFromTime(z.Unplugged),
		ChargeDuration: z.ChargeDuration,
		Electricity:    z.Electricity,
	}
}

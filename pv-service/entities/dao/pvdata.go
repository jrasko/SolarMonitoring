package dao

import "pv-service/graph/model"

type PVData struct {
	Time uint32

	Dc1U uint16
	Dc1I uint16
	Dc1P uint16
	Dc1T uint32
	Dc1S uint32

	Dc2U uint16
	Dc2I uint16
	Dc2P uint16
	Dc2T uint32
	Dc2S uint32

	Dc3U uint16
	Dc3I uint16
	Dc3P uint16
	Dc3T uint32
	Dc3S uint32

	Ac1U uint16
	Ac1I uint16
	Ac1P int16
	Ac1T uint32

	Ac2U uint16
	Ac2I uint16
	Ac2P int16
	Ac2T uint32

	Ac3U uint16
	Ac3I uint16
	Ac3P int16
	Ac3T uint32

	AcF    float32
	FcI    int16
	Ain1   int16
	Ain2   int16
	Ain3   int16
	Ain4   int16
	AcS    uint16
	Err    uint16
	EnsS   uint16
	EnsErr uint16
	KbS    string
	TotalE uint32
	IsoR   uint32
	Event  string
}

func (PVData) TableName() string {
	return "public.pvdata"
}

func (p PVData) ToModel() model.RawData {
	return model.RawData{
		Time:   p.Time,
		Dc1U:   uint32(p.Dc1U),
		Dc1I:   uint32(p.Dc1I),
		Dc1P:   uint32(p.Dc1P),
		Dc1T:   p.Dc1T,
		Dc1S:   p.Dc1S,
		Dc2U:   uint32(p.Dc2U),
		Dc2I:   uint32(p.Dc2I),
		Dc2P:   uint32(p.Dc2P),
		Dc2T:   p.Dc2T,
		Dc2S:   p.Dc2S,
		Dc3U:   uint32(p.Dc3U),
		Dc3I:   uint32(p.Dc3I),
		Dc3P:   uint32(p.Dc3P),
		Dc3T:   p.Dc3T,
		Dc3S:   p.Dc3S,
		Ac1U:   uint32(p.Ac1U),
		Ac1I:   uint32(p.Ac1I),
		Ac1P:   int32(p.Ac1P),
		Ac1T:   p.Ac1T,
		Ac2U:   uint32(p.Ac1U),
		Ac2I:   uint32(p.Ac2I),
		Ac2P:   int32(p.Ac2P),
		Ac2T:   p.Ac2T,
		Ac3U:   uint32(p.Ac3U),
		Ac3I:   uint32(p.Ac3I),
		Ac3P:   int32(p.Ac3P),
		Ac3T:   p.Ac3T,
		AcF:    float64(p.AcF),
		FcI:    int32(p.FcI),
		Ain1:   int32(p.Ain1),
		Ain2:   int32(p.Ain2),
		Ain3:   int32(p.Ain3),
		AcS:    uint32(p.AcS),
		Err:    int32(p.Err),
		EnsErr: uint32(p.EnsErr),
		Event:  p.Event,
	}
}

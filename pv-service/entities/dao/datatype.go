package dao

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
	TotalE uint16
	IsoR   uint32
	Event  string
}

func (PVData) TableName() string {
	return "public.pvdata"
}

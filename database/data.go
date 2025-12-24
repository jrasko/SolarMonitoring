package database

/*
type pvdata struct {
	Time int64

	Dc1U int
	Dc1I int
	Dc1P int
	Dc1T int
	Dc1S int

	Dc2U int
	Dc2I int
	Dc2P int
	Dc2T int
	Dc2S int

	Dc3U int
	Dc3I int
	Dc3P int
	Dc3T int
	Dc3S int

	Ac1U int
	Ac1I int
	Ac1P int
	Ac1T int

	Ac2U int
	Ac2I int
	Ac2P int
	Ac2T int

	Ac3U int
	Ac3I int
	Ac3P int
	Ac3T int

	AcF    float64
	FcI    int
	Ain1   int
	Ain2   int
	Ain3   int
	Ain4   int
	AcS    int
	Err    int
	EnsS   int
	EnsErr int
	KbS    string
	TotalE int
	IsoR   int
	Event  string
}
*/

type DailyData struct {
	Time   int64
	TotalE int
}

func (d DailyData) TableName() string {
	return "pvdata"
}

type MinuteData struct {
	Time int64
	DcI  []int
}

type dbMinuteData struct {
	Time int64
	Dc1I int
	Dc2I int
	Dc3I int
}

type dbMinuteDataSet []dbMinuteData

func (d dbMinuteDataSet) toExternal() []MinuteData {
	external := make([]MinuteData, len(d))
	for i, data := range d {
		external[i] = MinuteData{
			Time: data.Time,
			DcI:  []int{data.Dc1I, data.Dc2I, data.Dc3I},
		}
	}
	return external
}

func (d dbMinuteData) TableName() string {
	return "pvdata"
}

package model

import (
	"fmt"
	"io"
	"time"
)

type Date struct {
	Day   int
	Month int
	Year  int
}

func (d Date) toTime() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

func (d Date) Compare(b Date) int {
	return d.toTime().Compare(b.toTime())
}

func (d *Date) UnmarshalGQL(interface{}) error {
	return nil
}

func (d Date) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprintf(w, `"%02d.%02d.%04d"`, d.Day, d.Month, d.Year)
	if err != nil {
		fmt.Println(err)
	}
}

func (d Date) Yesterday() Date {
	yesterday := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	year, month, day := yesterday.Date()
	return Date{
		Day:   day,
		Month: int(month),
		Year:  year,
	}
}

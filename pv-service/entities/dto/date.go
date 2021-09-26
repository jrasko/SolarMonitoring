package dto

import (
	"fmt"
	"io"
)

type Date struct {
	Day   uint8
	Month uint8
	Year  uint16
}

var monthLengths = [12]uint8{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func (t *Date) UnmarshalGQL(interface{}) error {
	return nil
}

func (t Date) MarshalGQL(w io.Writer) {
	dateString := fmt.Sprintf("\"%02d.%02d.%04d\"", t.Day, t.Month, t.Year)
	_, err := w.Write([]byte(dateString))
	if err != nil {
		fmt.Println(err)
	}
}

func (t Date) isLeapYear() bool {
	return t.Year%400 == 0 || (t.Year%4 == 0 && t.Year%100 != 0)
}
func (t Date) Yesterday() Date {
	t.Day--
	if t.Day >= 1 {
		return t
	}
	t.Month--
	switch {
	case t.isLeapYear() && t.Month == 2:
		t.Day = 29
	case t.Month == 0:
		t.Year--
		t.Month = 12
		t.Day = 31
	default:
		t.Day = monthLengths[t.Month-1]
	}
	return t
}

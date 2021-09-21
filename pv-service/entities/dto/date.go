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

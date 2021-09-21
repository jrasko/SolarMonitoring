package dto

import (
	"fmt"
	"io"
	"strconv"
)

type TimeOfDay struct {
	Hours   uint8
	Minutes uint8
	Seconds uint8
}

func (t *TimeOfDay) UnmarshalGQL(interface{}) error {
	return nil
}

func (t TimeOfDay) MarshalGQL(w io.Writer) {
	dateInt := int(t.Hours)*60*60 + int(t.Minutes)*60 + int(t.Seconds)
	// dateString := fmt.Sprintf("%d%02d%02d", t.Hours, t.Minutes, t.Seconds)
	_, err := w.Write([]byte(strconv.Itoa(dateInt)))
	if err != nil {
		fmt.Println(err)
	}
}

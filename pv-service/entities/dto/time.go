package dto

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	decade = 315532800
	hour   = 3600
)

var (
	UTC1 = time.FixedZone("UTC+01", 1)
)

type PVTime time.Time

func PVTimeFromUnix(unixTime uint32) PVTime {
	unixTime += decade
	timestamp := time.Unix(int64(unixTime), 0).In(UTC1)
	return PVTime(timestamp)
}

func PVTimeFromTime(timestamp time.Time) PVTime {
	return PVTime(timestamp)
}

func (t PVTime) ToTime() time.Time {
	return time.Time(t)
}

func (t PVTime) ToUnix() uint32 {
	unixTime := t.ToTime().Unix()
	unixTime -= decade
	return uint32(unixTime)
}

func (t PVTime) ToUnixPlus12Hours() uint32 {
	unixTime := t.ToTime().Unix()
	unixTime -= decade
	return uint32(unixTime) + 12*hour
}

func (t *PVTime) UnmarshalGQL(v interface{}) error {
	if tmpStr, ok := v.(string); ok {
		ts, err := time.Parse(time.RFC3339Nano, tmpStr)
		if err != nil {
			return err
		}
		*t = PVTimeFromTime(ts)
		return nil
	}
	return errors.New("time should be RFC3339Nano formatted string")
}

func (t PVTime) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(time.Time(t).Format(time.RFC3339Nano)))
}

type TimeOfDay struct {
	Hours   uint8
	Minutes uint8
	Seconds uint8
}

func (t *TimeOfDay) UnmarshalGQL(interface{}) error {
	return nil
}

func (t TimeOfDay) MarshalGQL(w io.Writer) {
	// as seconds
	dateInt := int(t.Hours)*60*60 + int(t.Minutes)*60 + int(t.Seconds)
	// dateString := fmt.Sprintf("%d%02d%02d", t.Hours, t.Minutes, t.Seconds)
	_, err := w.Write([]byte(strconv.Itoa(dateInt)))
	if err != nil {
		fmt.Println(err)
	}
}

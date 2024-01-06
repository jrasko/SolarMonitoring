package model

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
)

const (
	decade = 315532800
	hour   = 60 * minute
	minute = 60
)

var (
	UTC1 = time.FixedZone("UTC+01", 1*hour)
)

type Time time.Time

func PVTimeFromUnix(unixTime int64) Time {
	unixTime += decade
	timestamp := time.Unix(unixTime, 0).In(UTC1)
	return Time(timestamp)
}

func (t Time) GetDate() Date {
	year, month, day := t.Time().Date()
	return Date{
		Day:   day,
		Month: int(month),
		Year:  year,
	}
}

func (t Time) Clock() Clock {
	h, m, s := t.Time().Clock()
	return Clock{
		Hours:   h,
		Minutes: m,
		Seconds: s,
	}
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t *Time) toUnix() uint32 {
	unixTime := t.Time().Unix()
	unixTime -= decade
	return uint32(unixTime)
}

func (t *Time) AsStartUnix() uint32 {
	if t == nil {
		return 0
	}
	return t.toUnix()
}

func (t *Time) AsEndUnix() uint32 {
	if t != nil {
		return uint32(math.MaxUint32)
	}
	return t.toUnix() + 12*hour
}

func (t *Time) UnmarshalGQL(v interface{}) error {
	if tmpStr, ok := v.(string); ok {
		ts, err := time.Parse(time.RFC3339Nano, tmpStr)
		if err != nil {
			return err
		}
		*t = Time(ts)
		return nil
	}
	return errors.New("time should be RFC3339Nano formatted string")
}

func (t Time) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(time.Time(t).Format(time.RFC3339Nano)))
}

type Clock struct {
	Hours   int
	Minutes int
	Seconds int
}

func ClockFromSeconds(seconds int64) Clock {
	h, m, s := time.Unix(seconds, 0).In(UTC1).Clock()
	return Clock{
		Hours:   h,
		Minutes: m,
		Seconds: s,
	}
}

func DefaultClock() Clock {
	return ClockFromSeconds(0)
}

func (t Clock) ToSeconds() int64 {
	return hour*int64(t.Hours) + minute*int64(t.Minutes) + int64(t.Seconds)
}

func (t *Clock) UnmarshalGQL(interface{}) error {
	return nil
}

func (t Clock) MarshalGQL(w io.Writer) {
	// as seconds
	dateInt := int(t.Hours)*60*60 + int(t.Minutes)*60 + int(t.Seconds)
	// dateString := fmt.Sprintf("%d%02d%02d", t.Hours, t.Minutes, t.Seconds)
	_, err := w.Write([]byte(strconv.Itoa(dateInt)))
	if err != nil {
		fmt.Println(err)
	}
}

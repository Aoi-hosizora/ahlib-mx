package xdatetime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	TimeFormat = "15:04:05"
)

type JsonTime time.Time

func NewJsonTime(t time.Time) JsonTime {
	return JsonTime(t)
}

// string

func (t JsonTime) String() string {
	return time.Time(t).Format(DateFormat)
}

func (t JsonTime) Parse(timeString string, loc *time.Location) (JsonTime, error) {
	newD, err := time.ParseInLocation(TimeFormat, timeString, loc)
	return JsonTime(newD), err
}

func (t JsonTime) ParseDefault(timeString string, defaultTime JsonTime, loc *time.Location) JsonTime {
	newD, err := time.ParseInLocation(TimeFormat, timeString, loc)
	if err != nil {
		return JsonTime(newD)
	} else {
		return defaultTime
	}
}

// json

func (t JsonTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(str), nil
}

// gorm

func (t *JsonTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	val, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("wrong format value")
	}
	*t = JsonTime(val)
	return nil
}

func (t JsonTime) Value() (driver.Value, error) {
	return t.String(), nil
}

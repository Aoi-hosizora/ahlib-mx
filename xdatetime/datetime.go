package xdatetime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	ISO8601DateTimeFormat = "2006-01-02T15:04:05Z07:00"
	LocalDateTimeFormat   = "2006-01-02 15:04:05"
)

type JsonDateTime time.Time

func NewJsonDateTime(t time.Time) JsonDateTime {
	return JsonDateTime(t)
}

func (dt JsonDateTime) Time() time.Time {
	return time.Time(dt)
}

// string

func (dt JsonDateTime) String() string {
	return dt.Time().Format(ISO8601DateTimeFormat)
}

func (dt JsonDateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", dt.Time().Format(ISO8601DateTimeFormat))
	return []byte(str), nil
}

// parse

func ParseISO8601DateTime(dateTimeString string) (JsonDateTime, error) {
	newDt, err := time.Parse(ISO8601DateTimeFormat, dateTimeString)
	return JsonDateTime(newDt), err
}

func ParseISO8601DateTimeDefault(dateTimeString string, defaultDateTime JsonDateTime) JsonDateTime {
	newDt, err := time.Parse(ISO8601DateTimeFormat, dateTimeString)
	if err != nil {
		return JsonDateTime(newDt)
	} else {
		return defaultDateTime
	}
}

func ParseDateTime(dateTimeString string, layout string, loc *time.Location) (JsonDateTime, error) {
	newDt, err := time.ParseInLocation(layout, dateTimeString, loc)
	return JsonDateTime(newDt), err
}

func ParseDateTimeDefault(dateTimeString string, defaultDateTime JsonDateTime, layout string, loc *time.Location) JsonDateTime {
	newDt, err := time.ParseInLocation(layout, dateTimeString, loc)
	if err != nil {
		return JsonDateTime(newDt)
	} else {
		return defaultDateTime
	}
}

// gorm

func (dt *JsonDateTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	val, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("wrong format value")
	}
	*dt = JsonDateTime(val)
	return nil
}

func (dt JsonDateTime) Value() (driver.Value, error) {
	return dt.Time(), nil
}

package xdatetime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	DateTimeFormat = "2006-01-02 15:04:05"
)

type JsonDateTime time.Time

func NewJsonDateTime(t time.Time) JsonDateTime {
	return JsonDateTime(t)
}

// string

func (dt JsonDateTime) String() string {
	return time.Time(dt).Format(DateTimeFormat)
}

func (dt JsonDateTime) Parse(dateTimeString string, loc *time.Location) (JsonDateTime, error) {
	newDt, err := time.ParseInLocation(DateTimeFormat, dateTimeString, loc)
	return JsonDateTime(newDt), err
}

func (dt JsonDateTime) ParseDefault(dateTimeString string, defaultDateTime JsonDateTime, loc *time.Location) JsonDateTime {
	newDt, err := time.ParseInLocation(DateTimeFormat, dateTimeString, loc)
	if err != nil {
		return JsonDateTime(newDt)
	} else {
		return defaultDateTime
	}
}

// json

func (dt JsonDateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(dt).Format(DateTimeFormat))
	return []byte(str), nil
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
	return dt.String(), nil
}

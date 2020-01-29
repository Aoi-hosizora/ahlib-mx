package xdatetime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	DateFormat = "2006-01-02"
)

type JsonDate time.Time

func NewJsonDate(t time.Time) JsonDate {
	return JsonDate(t)
}

// string

func (d JsonDate) String() string {
	return time.Time(d).Format(DateFormat)
}

func (d JsonDate) Parse(dateString string, loc *time.Location) (JsonDate, error) {
	newD, err := time.ParseInLocation(DateFormat, dateString, loc)
	return JsonDate(newD), err
}

func (d JsonDate) ParseDefault(dateString string, defaultDate JsonDate, loc *time.Location) JsonDate {
	newD, err := time.ParseInLocation(DateFormat, dateString, loc)
	if err != nil {
		return JsonDate(newD)
	} else {
		return defaultDate
	}
}

// json

func (d JsonDate) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(d).Format(DateFormat))
	return []byte(str), nil
}

// gorm

func (d *JsonDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	val, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("wrong format value")
	}
	*d = JsonDate(val)
	return nil
}

func (d JsonDate) Value() (driver.Value, error) {
	return d.String(), nil
}

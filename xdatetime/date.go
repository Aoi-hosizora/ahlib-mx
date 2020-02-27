package xdatetime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	ISO8601DateFormat = "2006-01-02"
	LocalDateFormat = "2006-01-02"
)

type JsonDate time.Time

func NewJsonDate(t time.Time) JsonDate {
	return JsonDate(t)
}

func (d JsonDate) Time() time.Time {
	return time.Time(d)
}

// string

func (d JsonDate) String() string {
	return d.Time().Format(ISO8601DateFormat)
}

func (d JsonDate) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", d.Time().Format(ISO8601DateFormat))
	return []byte(str), nil
}

// parse

func ParseISO8601Date(dateString string) (JsonDate, error) {
	newD, err := time.Parse(ISO8601DateFormat, dateString)
	return JsonDate(newD), err
}

func ParseISO8601DateDefault(dateString string, defaultDate JsonDate) JsonDate {
	newD, err := time.Parse(ISO8601DateFormat, dateString)
	if err != nil {
		return JsonDate(newD)
	} else {
		return defaultDate
	}
}

func ParseDate(dateString string, layout string, loc *time.Location) (JsonDate, error) {
	newD, err := time.ParseInLocation(layout, dateString, loc)
	return JsonDate(newD), err
}

func ParseDateDefault(dateString string, defaultDate JsonDate, layout string, loc *time.Location) JsonDate {
	newD, err := time.ParseInLocation(layout, dateString, loc)
	if err != nil {
		return JsonDate(newD)
	} else {
		return defaultDate
	}
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
	return d.Time(), nil
}

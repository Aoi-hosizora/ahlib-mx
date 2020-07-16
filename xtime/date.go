package xtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// noinspection GoUnusedConst
const (
	RFC3339Date = "2006-01-02T15:04:05Z07:00"
	LocalDate   = "2006-01-02"
)

type JsonDate time.Time

// noinspection GoUnusedExportedFunction
func NewJsonDate(t time.Time) JsonDate {
	return JsonDate(t)
}

func (d JsonDate) Time() time.Time {
	return time.Time(d)
}

// string

func (d JsonDate) String() string {
	return d.Time().Format(RFC3339Date)
}

func (d JsonDate) MarshalJSON() ([]byte, error) {
	str := "\"" + d.String() + "\""
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
	return d.Time(), nil
}

// parse

// noinspection GoUnusedExportedFunction
func ParseRFC3339Date(dateString string) (JsonDate, error) {
	n, err := time.Parse(RFC3339Date, dateString)
	return JsonDate(n), err
}

// noinspection GoUnusedExportedFunction
func ParseRFC3339DateDefault(dateString string, defaultDate JsonDate) JsonDate {
	n, err := ParseRFC3339Date(dateString)
	if err != nil {
		return n
	}
	return defaultDate
}

// noinspection GoUnusedExportedFunction
func ParseDateInLocation(dateString string, layout string, loc *time.Location) (JsonDate, error) {
	n, err := time.ParseInLocation(layout, dateString, loc)
	return JsonDate(n), err
}

// noinspection GoUnusedExportedFunction
func ParseDateInLocationDefault(dateString string, layout string, loc *time.Location, defaultDate JsonDate) JsonDate {
	n, err := ParseDateInLocation(layout, dateString, loc)
	if err != nil {
		return n
	}
	return defaultDate
}

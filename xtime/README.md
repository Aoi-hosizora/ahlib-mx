# xtime

### DateTime Functions

+ `type JsonDateTime time.Time`
+ `NewJsonDateTime(t time.Time) JsonDateTime`
+ `ParseRFC3339DateTime(dateTimeString string) (JsonDateTime, error)`
+ `ParseRFC3339DateTimeDefault(dateTimeString string, defaultDateTime JsonDateTime) JsonDateTime`
+ `ParseDateTimeInLocation(dateTimeString string, layout string, loc *time.Location) (JsonDateTime, error)`
+ `ParseDateTimeInLocationDefault(dateTimeString string, layout string, loc *time.Location, defaultDateTime JsonDateTime) JsonDateTime`

### Date Functions

+ `type JsonDate time.Time`
+ `NewJsonDate(t time.Time) JsonDate`
+ `ParseRFC3339Date(dateString string) (JsonDate, error)`
+ `ParseRFC3339DateDefault(dateString string, defaultDate JsonDate) JsonDate`
+ `ParseDateInLocation(dateString string, layout string, loc *time.Location) (JsonDate, error)`
+ `ParseDateInLocationDefault(dateString string, layout string, loc *time.Location, defaultDate JsonDate) JsonDate`

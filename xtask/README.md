# xtask

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/robfig/cron/v3

## Documents

### Types

+ `type CronWrapper struct`
+ `type FuncJob struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func NewCronWrapper(c *cron.Cron) *CronWrapper`

### Methods

+ `func (c *CronWrapper) Cron() *cron.Cron`
+ `func (c *CronWrapper) Jobs() []*FuncJob`
+ `func (c *CronWrapper) ScheduleParser() cron.ScheduleParser`
+ `func (c *CronWrapper) AddJobByCronSpec(title string, spec string, f func() error) (cron.EntryID, error)`
+ `func (c *CronWrapper) AddJobBySchedule(title string, schedule cron.Schedule, f func() error) cron.EntryID`
+ `func (c *CronWrapper) RemoveJob(id cron.EntryID)`
+ `func (c *CronWrapper) SetJobAddedCallback(cb func(job *FuncJob))`
+ `func (c *CronWrapper) SetJobRemovedCallback(cb func(job *FuncJob))`
+ `func (c *CronWrapper) SetPanicHandler(handler func(v interface{}))`
+ `func (c *CronWrapper) SetErrorHandler(handler func(err error))`
+ `func (f *FuncJob) Title() string`
+ `func (f *FuncJob) CronSpec() string`
+ `func (f *FuncJob) Schedule() cron.Schedule`
+ `func (f *FuncJob) ScheduleExpr() string`
+ `func (f *FuncJob) Funcname() string`
+ `func (f *FuncJob) Entry() *cron.Entry`
+ `func (f *FuncJob) EntryID() cron.EntryID`
+ `func (f *FuncJob) Run()`

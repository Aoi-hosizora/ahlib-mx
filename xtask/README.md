# xtask

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/robfig/cron/v3

## Documents

### Types

+ `type CronTask struct`
+ `type FuncJob struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func NewCronTask(c *cron.Cron) *CronTask`
+ `func DefaultAddedCallback(j *FuncJob)`
+ `func DefaultColorizedAddedCallback(j *FuncJob)`

### Methods

+ `func (c *CronTask) Cron() *cron.Cron`
+ `func (c *CronTask) ScheduleParser() cron.ScheduleParser`
+ `func (c *CronTask) Jobs() []*FuncJob`
+ `func (c *CronTask) AddJobByCronSpec(title string, spec string, f func()) (cron.EntryID, error)`
+ `func (c *CronTask) AddJobBySchedule(title string, schedule cron.Schedule, f func()) cron.EntryID`
+ `func (c *CronTask) RemoveJob(id cron.EntryID)`
+ `func (c *CronTask) SetAddedCallback(cb func(job *FuncJob))`
+ `func (c *CronTask) SetRemovedCallback(cb func(job *FuncJob))`
+ `func (c *CronTask) SetScheduledCallback(cb func(job *FuncJob))`
+ `func (c *CronTask) SetPanicHandler(handler func(job *FuncJob, v interface{}))`
+ `func (f *FuncJob) Title() string`
+ `func (f *FuncJob) CronSpec() string`
+ `func (f *FuncJob) Schedule() cron.Schedule`
+ `func (f *FuncJob) ScheduleExpr() string`
+ `func (f *FuncJob) Funcname() string`
+ `func (f *FuncJob) Entry() *cron.Entry`
+ `func (f *FuncJob) EntryID() cron.EntryID`
+ `func (f *FuncJob) Run()`

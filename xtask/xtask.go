package xtask

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/robfig/cron/v3"
	"log"
	"reflect"
	"runtime"
)

// CronTask represents a task, or a job collection, which is implemented by wrapping cron.Cron.
type CronTask struct {
	cron *cron.Cron
	jobs []*FuncJob

	jobAddedCallback     func(job *FuncJob)
	jobRemovedCallback   func(job *FuncJob)
	jobScheduledCallback func(job *FuncJob)

	panicHandler func(job *FuncJob, v interface{})
	errorHandler func(job *FuncJob, err error)
}

// FuncJob represents a cron.Job with some information such as title, cron.Schedule and cron.Entry, stored in CronTask.
type FuncJob struct {
	title    string
	cronSpec string
	schedule cron.Schedule
	function func() error
	entry    *cron.Entry
	entryID  cron.EntryID

	parent *CronTask
}

var _ cron.Job = (*FuncJob)(nil)

// ========
// CronTask
// ========

// NewCronTask creates a default CronTask with given cron.Cron and default callbacks and handlers.
func NewCronTask(c *cron.Cron) *CronTask {
	return &CronTask{
		cron: c,
		jobs: make([]*FuncJob, 0),

		jobAddedCallback:     defaultJobAddedCallback,
		jobRemovedCallback:   defaultJobRemovedCallback,
		jobScheduledCallback: defaultJobScheduledCallback,
		panicHandler:         func(job *FuncJob, v interface{}) { log.Printf("Warning: Job %s panics with `%v`", job.title, v) },
		errorHandler:         nil, // skip
	}
}

// Cron returns cron.Cron from CronTask.
func (c *CronTask) Cron() *cron.Cron {
	return c.cron
}

// Jobs returns FuncJob slice from CronTask.
func (c *CronTask) Jobs() []*FuncJob {
	return c.jobs
}

// ScheduleParser returns cron.ScheduleParser from cron.Cron in CronTask.
func (c *CronTask) ScheduleParser() cron.ScheduleParser {
	return xreflect.GetUnexportedField(xreflect.FieldValueOf(c.cron, "parser")).Interface().(cron.ScheduleParser)
}

// newFuncJob creates a FuncJob with given parameters with CronTask parent.
func (c *CronTask) newFuncJob(title string, spec string, schedule cron.Schedule, f func() error) *FuncJob {
	return &FuncJob{parent: c, title: title, cronSpec: spec, schedule: schedule, function: f}
}

const (
	panicNilFunction = "xtask: nil function"
	panicNilSchedule = "xtask: nil schedule"
)

// AddJobByCronSpec adds a FuncJob to cron.Cron and CronTask by given title, cron spec and function.
func (c *CronTask) AddJobByCronSpec(title string, spec string, f func() error) (cron.EntryID, error) {
	if f == nil {
		panic(panicNilFunction)
	}
	job := c.newFuncJob(title, spec, nil, f)
	id, err := c.cron.AddJob(spec, job) // <<<
	if err != nil {
		return 0, err
	}

	entry := c.cron.Entry(id)
	job.entry = &entry
	job.entryID = id
	c.jobs = append(c.jobs, job)
	if c.jobAddedCallback != nil {
		c.jobAddedCallback(job)
	}
	return id, nil
}

// AddJobBySchedule adds a FuncJob to cron.Cron and CronTask by given title, cron.Schedule and function.
func (c *CronTask) AddJobBySchedule(title string, schedule cron.Schedule, f func() error) cron.EntryID {
	if schedule == nil {
		panic(panicNilSchedule)
	}
	if f == nil {
		panic(panicNilFunction)
	}
	job := c.newFuncJob(title, "", schedule, f)
	id := c.cron.Schedule(schedule, job) // <<<

	entry := c.cron.Entry(id)
	job.entry = &entry
	job.entryID = id
	c.jobs = append(c.jobs, job)
	if c.jobAddedCallback != nil {
		c.jobAddedCallback(job)
	}
	return id
}

// RemoveJob removes a cron.Entry by given cron.EntryID from cron.Cron and CronTask.
func (c *CronTask) RemoveJob(id cron.EntryID) {
	c.cron.Remove(id)
	c.jobs = xslice.DeleteAllWithG(c.jobs, &FuncJob{entryID: id}, func(i, j interface{}) bool {
		if i.(*FuncJob).entryID == j.(*FuncJob).entryID {
			if c.jobRemovedCallback != nil {
				c.jobRemovedCallback(i.(*FuncJob))
			}
			return true
		}
		return false
	}).([]*FuncJob)
}

// defaultJobAddedCallback is the default jobAddedCallback, can be modified by CronTask.SetJobAddedCallback.
//
// The default callback logs like:
// 	[Task-debug] job1, 0/1 * * * * *             --> ... (EntryID: 1)
// 	[Task-debug] job3, every 3s                  --> ... (EntryID: 3)
// 	[Task-debug] job4, <parsed SpecSchedule>     --> ... (EntryID: 4)
// 	            |-------------------------------|   |----------------|
// 	                           31                          ...
func defaultJobAddedCallback(j *FuncJob) {
	fmt.Printf("[Task-debug] %-31s --> %s (EntryID: %d)\n", fmt.Sprintf("%s, %s", j.Title(), j.ScheduleExpr()), j.Funcname(), j.EntryID())
}

// defaultJobRemovedCallback is the default jobRemovedCallback, can be modified by CronTask.SetJobRemovedCallback
//
// The default callback logs like:
// 	[Task-debug] Remove job: job1, EntryID: 1
// 	[Task-debug] Remove job: job2, EntryID: 2
func defaultJobRemovedCallback(j *FuncJob) {
	fmt.Printf("[Task-debug] Remove job: %s, EntryID: %d\n", j.Title(), j.EntryID())
}

// defaultJobScheduledCallback is the default jobRemovedCallback, can be modified by CronTask.SetJobScheduledCallback
//
// The default callback does nothing.
func defaultJobScheduledCallback(*FuncJob) {
	// skip
}

// SetJobAddedCallback sets job added callback, this will be invoked after FuncJob added.
func (c *CronTask) SetJobAddedCallback(cb func(job *FuncJob)) {
	c.jobAddedCallback = cb
}

// SetJobRemovedCallback sets job removed callback, this will be invoked after FuncJob removed.
func (c *CronTask) SetJobRemovedCallback(cb func(job *FuncJob)) {
	c.jobRemovedCallback = cb
}

// SetJobScheduledCallback sets job scheduled callback, this will be invoked when after FuncJob scheduled.
func (c *CronTask) SetJobScheduledCallback(cb func(job *FuncJob)) {
	c.jobScheduledCallback = cb
}

// SetPanicHandler sets panic handler for jobs executing, defaults to print warning message.
func (c *CronTask) SetPanicHandler(handler func(job *FuncJob, v interface{})) {
	c.panicHandler = handler
}

// SetErrorHandler sets error handler for jobs executing, defaults to do nothing.
func (c *CronTask) SetErrorHandler(handler func(job *FuncJob, err error)) {
	c.errorHandler = handler
}

// =======
// FuncJob
// =======

// Title returns title from FuncJob.
func (f *FuncJob) Title() string {
	return f.title
}

// CronSpec returns cron spec string from FuncJob.
func (f *FuncJob) CronSpec() string {
	return f.cronSpec
}

// Schedule returns cron.Schedule from FuncJob.
func (f *FuncJob) Schedule() cron.Schedule {
	return f.schedule
}

// ScheduleExpr returns schedule expr from FuncJob, and this is generated by cron spec string and cron.Schedule.
func (f *FuncJob) ScheduleExpr() string {
	if f.cronSpec != "" {
		return f.cronSpec
	}
	if _, ok := f.schedule.(*cron.SpecSchedule); ok {
		return "<parsed SpecSchedule>"
	}
	if s, ok := f.schedule.(cron.ConstantDelaySchedule); ok {
		return fmt.Sprintf("every %s", s.Delay.String())
	}
	return "<unknown Schedule>"
}

// Funcname returns job function name from FuncJob.
func (f *FuncJob) Funcname() string {
	return runtime.FuncForPC(reflect.ValueOf(f.function).Pointer()).Name()
}

// Entry returns cron.Entry from FuncJob.
func (f *FuncJob) Entry() *cron.Entry {
	return f.entry
}

// EntryID returns cron.EntryID from FuncJob.
func (f *FuncJob) EntryID() cron.EntryID {
	return f.entryID
}

// Run runs the FuncJob with panic handler and error handler, and implements cron.Job interface.
func (f *FuncJob) Run() {
	defer func() {
		v := recover()
		if v != nil && f.parent.panicHandler != nil {
			f.parent.panicHandler(f, v) // defaults to log warning
		}
	}()

	if f.parent.jobScheduledCallback != nil {
		f.parent.jobScheduledCallback(f) // defaults to ignore
	}
	err := f.function()
	if err != nil && f.parent.errorHandler != nil {
		f.parent.errorHandler(f, err) // defaults to ignore
	}
}

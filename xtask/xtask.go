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

// CronWrapper represents a cron.Cron wrapper type with some helper functions.
type CronWrapper struct {
	cron *cron.Cron
	jobs []*FuncJob

	jobAddedCallback   func(job *FuncJob)
	jobRemovedCallback func(job *FuncJob)
	panicHandler       func(v interface{})
	errorHandler       func(err error)
}

// FuncJob represents a cron.Job with some information such as title, cron.Schedule and cron.Entry, stored in CronWrapper.
type FuncJob struct {
	title    string
	cronSpec string
	schedule cron.Schedule
	function func() error
	entry    *cron.Entry
	entryID  cron.EntryID

	parent *CronWrapper
}

var _ cron.Job = (*FuncJob)(nil)

// ===========
// CronWrapper
// ===========

// NewCronWrapper creates a default CronWrapper with given cron.Cron and default callbacks and handlers.
func NewCronWrapper(c *cron.Cron) *CronWrapper {
	return &CronWrapper{
		cron: c,
		jobs: make([]*FuncJob, 0),

		jobAddedCallback: func(j *FuncJob) {
			fmt.Printf("[Task-debug] %-29s --> %s (EntryID: %d)\n", fmt.Sprintf("%s, %s", j.Title(), j.ScheduleExpr()), j.Funcname(), j.EntryID())
		},
		jobRemovedCallback: func(j *FuncJob) {
			fmt.Printf("[Task-debug] Remove job: %s, EntryID: %d\n", j.Title(), j.EntryID())
		},
		panicHandler: func(v interface{}) {
			log.Printf("Warning: Panic with `%v`", v)
		},
	}
}

// Cron returns cron.Cron from CronWrapper.
func (c *CronWrapper) Cron() *cron.Cron {
	return c.cron
}

// Jobs returns FuncJob slice from CronWrapper.
func (c *CronWrapper) Jobs() []*FuncJob {
	return c.jobs
}

// ScheduleParser returns cron.ScheduleParser from cron.Cron in CronWrapper.
func (c *CronWrapper) ScheduleParser() cron.ScheduleParser {
	return xreflect.GetUnexportedField(xreflect.FieldValueOf(c.cron, "parser")).Interface().(cron.ScheduleParser)
}

// newFuncJob creates a FuncJob with given parameters with CronWrapper parent.
func (c *CronWrapper) newFuncJob(title string, spec string, schedule cron.Schedule, f func() error) *FuncJob {
	return &FuncJob{parent: c, title: title, cronSpec: spec, schedule: schedule, function: f}
}

const (
	panicNilFunction = "xtask: nil function"
	panicNilSchedule = "xtask: nil schedule"
)

// AddJobByCronSpec adds a FuncJob to cron.Cron and CronWrapper by given title, cron spec and function.
func (c *CronWrapper) AddJobByCronSpec(title string, spec string, f func() error) (cron.EntryID, error) {
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

// AddJobBySchedule adds a FuncJob to cron.Cron and CronWrapper by given title, cron.Schedule and function.
func (c *CronWrapper) AddJobBySchedule(title string, schedule cron.Schedule, f func() error) cron.EntryID {
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

// RemoveJob removes a cron.Entry by given cron.EntryID from cron.Cron and CronWrapper.
func (c *CronWrapper) RemoveJob(id cron.EntryID) {
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

// SetJobAddedCallback sets job added callback, this will be invoked after FuncJob added.
//
// The default callback logs like:
// 	[Task-debug] job1, 0/1 * * * * *           --> ... (EntryID: 1)
// 	[Task-debug] job3, every 3s                --> ... (EntryID: 3)
// 	[Task-debug] job4, <parsed SpecSchedule>   --> ... (EntryID: 4)
// 	            |-----------------------------|   |----------------|
// 	                          29                          ...
func (c *CronWrapper) SetJobAddedCallback(cb func(job *FuncJob)) {
	c.jobAddedCallback = cb
}

// SetJobRemovedCallback sets job removed callback, this will be invoked after FuncJob removed.
//
// The default callback logs like:
// 	[Task-debug] Remove job: job3, EntryID: 3
func (c *CronWrapper) SetJobRemovedCallback(cb func(job *FuncJob)) {
	c.jobRemovedCallback = cb
}

// SetPanicHandler sets panic handler for jobs executing.
func (c *CronWrapper) SetPanicHandler(handler func(v interface{})) {
	c.panicHandler = handler
}

// SetErrorHandler sets error handler for jobs executing.
func (c *CronWrapper) SetErrorHandler(handler func(err error)) {
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
			f.parent.panicHandler(v) // defaults to log warning
		}
	}()

	err := f.function()
	if err != nil && f.parent.errorHandler != nil {
		f.parent.errorHandler(err) // defaults to ignore
	}
}

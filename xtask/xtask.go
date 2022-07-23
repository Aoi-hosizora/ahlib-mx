package xtask

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/robfig/cron/v3"
	"log"
	"reflect"
	"runtime"
	"sync"
)

// CronTask represents a job collection, or called a task, which is implemented by wrapping cron.Cron.
type CronTask struct {
	cron  *cron.Cron
	jobs  []*FuncJob
	muJob sync.RWMutex

	addedCallback     func(job *FuncJob)
	removedCallback   func(job *FuncJob)
	scheduledCallback func(job *FuncJob)
	panicHandler      func(job *FuncJob, v interface{})
}

// FuncJob represents a cron.Job with some information such as title, cron.Schedule and cron.Entry, stored in CronTask.
type FuncJob struct {
	title    string
	cronSpec string        // can be empty
	schedule cron.Schedule // can be nil
	function func()
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

		addedCallback:     DefaultAddedCallback,
		removedCallback:   defaultJobRemovedCallback,
		scheduledCallback: nil, // defaults to do nothing
		panicHandler: func(job *FuncJob, v interface{}) {
			log.Printf("xtask warning: Job \"%s\" panicked with `%v`", job.title, v)
		},
	}
}

// Cron returns cron.Cron from CronTask.
func (c *CronTask) Cron() *cron.Cron {
	return c.cron
}

// ScheduleParser returns cron.ScheduleParser from cron.Cron in CronTask.
func (c *CronTask) ScheduleParser() cron.ScheduleParser {
	return xreflect.GetUnexportedField(xreflect.FieldValueOf(c.cron, "parser")).Interface().(cron.ScheduleParser)
}

// Jobs returns FuncJob slice from CronTask.
func (c *CronTask) Jobs() []*FuncJob {
	c.muJob.RLock()
	out := c.jobs
	c.muJob.RUnlock()
	return out
}

const (
	panicNilFunction = "xtask: nil job function"
	panicNilSchedule = "xtask: nil cron schedule"
)

// AddJobByCronSpec adds a FuncJob to cron.Cron and CronTask by given repeatable title, cron spec and function.
func (c *CronTask) AddJobByCronSpec(title string, spec string, f func()) (cron.EntryID, error) {
	if f == nil {
		panic(panicNilFunction)
	}
	schedule, err := c.ScheduleParser().Parse(spec) // use cron.Schedule() rather than cron.AddJob()
	if err != nil {
		return 0, err
	}
	return c.addJob(title, spec, schedule, f), nil // by spec when spec is not empty
}

// AddJobBySchedule adds a FuncJob to cron.Cron and CronTask by given repeatable title, cron.Schedule and function.
func (c *CronTask) AddJobBySchedule(title string, schedule cron.Schedule, f func()) cron.EntryID {
	if schedule == nil {
		panic(panicNilSchedule)
	}
	if f == nil {
		panic(panicNilFunction)
	}
	return c.addJob(title, "", schedule, f) // by schedule when spec is empty
}

// addJob adds given FuncJob's fields to cron.Cron using given parsed cron.Schedule and returns the added cron.EntryID.
func (c *CronTask) addJob(title string, spec string, schedule cron.Schedule, f func()) cron.EntryID {
	job := &FuncJob{parent: c, title: title, cronSpec: spec /* maybe empty */, schedule: schedule, function: f}

	c.muJob.Lock()
	id := c.cron.Schedule(schedule, job) // always use schedule
	entry := c.cron.Entry(id)
	job.entry = &entry
	job.entryID = entry.ID
	c.jobs = append(c.jobs, job)
	c.muJob.Unlock()

	if c.addedCallback != nil {
		c.addedCallback(job)
	}
	return id
}

// RemoveJob removes a cron.Entry by given cron.EntryID from cron.Cron and CronTask.
func (c *CronTask) RemoveJob(id cron.EntryID) {
	c.muJob.Lock()
	c.cron.Remove(id)
	c.jobs = xslice.DeleteAllWithG(c.jobs, &FuncJob{entryID: id}, func(i, j interface{}) bool {
		if i.(*FuncJob).entryID == j.(*FuncJob).entryID {
			if c.removedCallback != nil {
				c.removedCallback(i.(*FuncJob))
			}
			return true
		}
		return false
	}).([]*FuncJob)
	c.muJob.Unlock()
}

// DefaultAddedCallback is the default CronTask's addedCallback, can be modified by CronTask.SetAddedCallback.
//
// The default callback logs like (just like gin.DebugPrintRouteFunc):
// 	[Task] job1, 0/1 * * * * *             --> ... (EntryID: 1)
// 	[Task] job3, every 3s                  --> ... (EntryID: 3)
// 	[Task] job4, <parsed SpecSchedule>     --> ... (EntryID: 4)
// 	      |-------------------------------|   |----------------|
// 	                     31                          ...
func DefaultAddedCallback(j *FuncJob) {
	fmt.Printf("[Task] %-31s --> %s (EntryID: %d)\n", fmt.Sprintf("%s, %s", j.Title(), j.ScheduleExpr()), j.Funcname(), j.EntryID())
}

// DefaultColorizedAddedCallback is the DefaultAddedCallback (CronTask's addedCallback) in color.
//
// The default callback logs like (just like gin.DebugPrintRouteFunc):
// 	[Task] job1, 0/1 * * * * *             --> ... (EntryID: 1)
// 	[Task] job3, every 3s                  --> ... (EntryID: 3)
// 	[Task] job4, <parsed SpecSchedule>     --> ... (EntryID: 4)
// 	      |-------------------------------|   |----------------|
// 	                  31 (blue)                      ...
func DefaultColorizedAddedCallback(j *FuncJob) {
	fmt.Printf("[Task] %s --> %s (EntryID: %d)\n", xcolor.Blue.ASprintf(-31, "%s, %s", j.Title(), j.ScheduleExpr()), j.Funcname(), j.EntryID())
}

// defaultJobRemovedCallback is the default removedCallback, can be modified by CronTask.SetRemovedCallback
//
// The default callback logs like:
// 	[Task] Remove job: job1, EntryID: 1
func defaultJobRemovedCallback(j *FuncJob) {
	fmt.Printf("[Task] Remove job: %s, EntryID: %d\n", j.Title(), j.EntryID())
}

// SetAddedCallback sets job added callback, this will be invoked after FuncJob added, defaults to DefaultAddedCallback.
func (c *CronTask) SetAddedCallback(cb func(job *FuncJob)) {
	c.addedCallback = cb
}

// SetRemovedCallback sets job removed callback, this will be invoked after FuncJob removed, defaults to defaultJobRemovedCallback.
func (c *CronTask) SetRemovedCallback(cb func(job *FuncJob)) {
	c.removedCallback = cb
}

// SetScheduledCallback sets job scheduled callback, this will be invoked when after FuncJob scheduled, defaults to do nothing.
func (c *CronTask) SetScheduledCallback(cb func(job *FuncJob)) {
	c.scheduledCallback = cb
}

// SetPanicHandler sets panic handler for jobs executing, defaults to print warning message.
func (c *CronTask) SetPanicHandler(handler func(job *FuncJob, v interface{})) {
	c.panicHandler = handler
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

// ScheduleExpr returns schedule expr from FuncJob, is generated by cron spec string and cron.Schedule.
func (f *FuncJob) ScheduleExpr() string {
	if f.cronSpec != "" { // by spec
		return f.cronSpec
	}
	if _, ok := f.schedule.(*cron.SpecSchedule); ok {
		return "<parsed SpecSchedule>" // hide origin spec
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

// Run runs the FuncJob with panic handler, this implements cron.Job interface.
func (f *FuncJob) Run() {
	defer func() {
		v := recover()
		if v != nil && f.parent.panicHandler != nil {
			f.parent.panicHandler(f, v) // defaults to print warning message
		}
	}()

	if f.parent.scheduledCallback != nil {
		f.parent.scheduledCallback(f) // defaults to do nothing
	}
	f.function()
}

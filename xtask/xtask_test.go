package xtask

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestCronTask(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		// c := cron.New(cron.WithSeconds(), cron.WithLogger(cron.DiscardLogger))
		task := NewCronTask(c)
		xtesting.Equal(t, task.Cron(), c)
		xtesting.Equal(t, len(task.Jobs()), 0)
		schedule, err := task.ScheduleParser().Parse("30 12 * * * *")
		xtesting.Nil(t, err)
		parser2 := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		schedule2, err := parser2.Parse("30 12 * * * *")
		xtesting.Nil(t, err)
		xtesting.Equal(t, schedule, schedule2)
	})

	t.Run("add and remove", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		task := NewCronTask(c)

		xtesting.Panic(t, func() {
			task.AddJobByCronSpec("", "", nil)
		})
		xtesting.Panic(t, func() {
			task.AddJobBySchedule("", nil, func() { return })
		})
		xtesting.Panic(t, func() {
			task.AddJobBySchedule("", cron.Every(time.Second), nil)
		})

		// 1
		id1, err := task.AddJobByCronSpec("job1", "0/1 * * * * *", func() { return })
		xtesting.Nil(t, err)
		xtesting.Equal(t, id1, cron.EntryID(1))
		// 2
		_, err = task.AddJobByCronSpec("job2", "@", func() { return })
		xtesting.NotNil(t, err)
		id2, err := task.AddJobByCronSpec("job2", "0/2 * * * * *", func() { return })
		xtesting.Equal(t, id2, cron.EntryID(2))
		sch2, _ := task.ScheduleParser().Parse("0/2 * * * * *")
		// 3
		sch3 := cron.Every(time.Second * 3)
		id3 := task.AddJobBySchedule("job3", sch3, func() { return })
		xtesting.Equal(t, id3, cron.EntryID(3))
		// remove 3
		xtesting.Equal(t, task.Jobs()[2].ScheduleExpr(), "every 3s")
		task.RemoveJob(id3)
		// 4
		sch4, err := task.ScheduleParser().Parse("0/4 * * * * *")
		xtesting.Nil(t, err)
		id4 := task.AddJobBySchedule("job4", sch4, func() { return })
		xtesting.Equal(t, id4, cron.EntryID(4))

		xtesting.Equal(t, task.Jobs()[0].Title(), "job1")
		xtesting.Equal(t, task.Jobs()[0].EntryID(), cron.EntryID(1))
		xtesting.Equal(t, task.Jobs()[1].CronSpec(), "0/2 * * * * *")
		xtesting.Equal(t, task.Jobs()[1].Schedule(), sch2)
		xtesting.Equal(t, task.Jobs()[2].Schedule(), sch4)
		xtesting.Equal(t, task.Jobs()[2].CronSpec(), "")
		xtesting.Equal(t, task.Jobs()[2].Entry().ID, cron.EntryID(4))
		xtesting.Equal(t, task.Jobs()[0].ScheduleExpr(), "0/1 * * * * *")
		xtesting.Equal(t, task.Jobs()[1].ScheduleExpr(), "0/2 * * * * *")
		xtesting.Equal(t, task.Jobs()[2].ScheduleExpr(), "<parsed SpecSchedule>")
		f := &FuncJob{schedule: nil}
		xtesting.Equal(t, f.ScheduleExpr(), "<unknown Schedule>") // fake
	})
}

func TestFuncJob(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		task := NewCronTask(c)

		task.AddJobByCronSpec("every1s", "0/1 * * * * *", func() {
			log.Printf("every1s_1_%s", time.Now().Format(time.RFC3339Nano))
		})
		task.SetAddedCallback(DefaultColorizedAddedCallback)
		task.AddJobByCronSpec("every2s", "0/2 * * * * *", func() {
			log.Printf("every2s_2_%s", time.Now().Format(time.RFC3339Nano))
		})
		task.SetRemovedCallback(func(j *FuncJob) {
			log.Printf("[Task] Remove job: %s | EntryID: %d", j.Title(), j.EntryID())
		})
		task.RemoveJob(2)
		task.SetAddedCallback(func(j *FuncJob) {
			log.Printf("[Task] %-29s | %s (EntryID: %d)", fmt.Sprintf("%s, %s", j.Title(), j.ScheduleExpr()), j.Funcname(), j.EntryID())
		})
		task.AddJobByCronSpec("every2s", "0/2 * * * * *", func() {
			log.Printf("every2s_3_%s", time.Now().Format(time.RFC3339Nano))
		})
		task.AddJobByCronSpec("every1s", "0/1 * * * * *", func() {
			log.Printf("every1s_4_%s", time.Now().Format(time.RFC3339Nano))
		})
		task.SetScheduledCallback(func(j *FuncJob) {
			log.Printf("[Task] Executing job: %s", j.Title())
		})
		xtesting.Equal(t, len(task.Jobs()), 3)
		xtesting.Equal(t, len(task.Cron().Entries()), 3)

		task.Cron().Start()
		time.Sleep(time.Second * 3)
		<-task.Cron().Stop().Done()
	})

	t.Run("panic", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		task := NewCronTask(c)
		task.AddJobBySchedule("panic", cron.Every(time.Second), func() {
			panic("test")
		})

		task.Cron().Start()
		time.Sleep(time.Second + time.Millisecond*200)
		panicV := (interface{})(nil)
		task.SetPanicHandler(func(job *FuncJob, v interface{}) { log.Printf("panic: %v | %s", v, job.Title()); panicV = v })
		time.Sleep(time.Second + time.Millisecond*200)
		xtesting.Equal(t, panicV, "test")
		<-task.Cron().Stop().Done()
	})
}

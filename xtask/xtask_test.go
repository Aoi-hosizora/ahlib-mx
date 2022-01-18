package xtask

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestCronWrapper(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		c := cron.New(cron.WithSeconds() /*, cron.WithLogger(cron.DiscardLogger)*/)
		cw := NewCronWrapper(c)
		xtesting.Equal(t, cw.Cron(), c)
		xtesting.Equal(t, len(cw.Jobs()), 0)
		schedule, err := cw.ScheduleParser().Parse("30 12 * * * *")
		xtesting.Nil(t, err)
		parser2 := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		schedule2, err := parser2.Parse("30 12 * * * *")
		xtesting.Nil(t, err)
		xtesting.Equal(t, schedule, schedule2)
	})

	t.Run("add and remove", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		cw := NewCronWrapper(c)

		xtesting.Panic(t, func() {
			cw.AddJobByCronSpec("", "", nil)
		})
		xtesting.Panic(t, func() {
			cw.AddJobBySchedule("", nil, func() error { return nil })
		})
		xtesting.Panic(t, func() {
			cw.AddJobBySchedule("", cron.Every(time.Second), nil)
		})

		// 1
		id1, err := cw.AddJobByCronSpec("job1", "0/1 * * * * *", func() error { return nil })
		xtesting.Nil(t, err)
		xtesting.Equal(t, id1, cron.EntryID(1))
		// 2
		_, err = cw.AddJobByCronSpec("job2", "@", func() error { return nil })
		xtesting.NotNil(t, err)
		id2, err := cw.AddJobByCronSpec("job2", "0/2 * * * * *", func() error { return nil })
		xtesting.Equal(t, id2, cron.EntryID(2))
		// 3
		sch3 := cron.Every(time.Second * 3)
		id3 := cw.AddJobBySchedule("job3", sch3, func() error { return nil })
		xtesting.Equal(t, id3, cron.EntryID(3))
		// remove 3
		xtesting.Equal(t, cw.Jobs()[2].ScheduleExpr(), "every 3s")
		cw.RemoveJob(id3)
		// 4
		sch4, err := cw.ScheduleParser().Parse("0/4 * * * * *")
		xtesting.Nil(t, err)
		id4 := cw.AddJobBySchedule("job4", sch4, func() error { return nil })
		xtesting.Equal(t, id4, cron.EntryID(4))

		xtesting.Equal(t, cw.Jobs()[0].Title(), "job1")
		xtesting.Equal(t, cw.Jobs()[0].EntryID(), cron.EntryID(1))
		xtesting.Equal(t, cw.Jobs()[1].CronSpec(), "0/2 * * * * *")
		xtesting.Equal(t, cw.Jobs()[1].Schedule(), nil)
		xtesting.Equal(t, cw.Jobs()[2].Schedule(), sch4)
		xtesting.Equal(t, cw.Jobs()[2].CronSpec(), "")
		xtesting.Equal(t, cw.Jobs()[2].Entry().ID, cron.EntryID(4))
		xtesting.Equal(t, cw.Jobs()[0].ScheduleExpr(), "0/1 * * * * *")
		xtesting.Equal(t, cw.Jobs()[1].ScheduleExpr(), "0/2 * * * * *")
		xtesting.Equal(t, cw.Jobs()[2].ScheduleExpr(), "<parsed SpecSchedule>")
		xtesting.Equal(t, cw.newFuncJob("", "", nil, nil).ScheduleExpr(), "<unknown Schedule>") // fake
	})
}

func TestFuncJob(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		cw := NewCronWrapper(c)

		cw.SetJobAddedCallback(func(j *FuncJob) {
			log.Printf("[Task] %-29s | %s (EntryID: %d)", fmt.Sprintf("%s, %s", j.Title(), j.ScheduleExpr()), j.Funcname(), j.EntryID())
		})
		cw.SetJobRemovedCallback(func(j *FuncJob) {
			log.Printf("[Task] Remove job: %s | EntryID: %d", j.Title(), j.EntryID())
		})
		cw.AddJobByCronSpec("every1s", "0/1 * * * * *", func() error {
			log.Printf("every1s_1_%s", time.Now().Format(time.RFC3339Nano))
			return nil
		})
		cw.AddJobByCronSpec("every2s", "0/2 * * * * *", func() error {
			log.Printf("every2s_2_%s", time.Now().Format(time.RFC3339Nano))
			return nil
		})
		cw.RemoveJob(2)
		cw.AddJobByCronSpec("every2s", "0/2 * * * * *", func() error {
			log.Printf("every2s_3_%s", time.Now().Format(time.RFC3339Nano))
			return nil
		})
		cw.AddJobByCronSpec("every1s", "0/1 * * * * *", func() error {
			log.Printf("every1s_4_%s", time.Now().Format(time.RFC3339Nano))
			return nil
		})
		xtesting.Equal(t, len(cw.Jobs()), 3)
		xtesting.Equal(t, len(cw.Cron().Entries()), 3)

		cw.Cron().Start()
		time.Sleep(time.Second * 3)
		<-cw.Cron().Stop().Done()
	})

	t.Run("panic and error", func(t *testing.T) {
		c := cron.New(cron.WithSeconds())
		cw := NewCronWrapper(c)
		cw.AddJobBySchedule("panic", cron.Every(time.Second), func() error {
			panic("test")
		})
		cw.AddJobBySchedule("error", cron.Every(time.Second), func() error {
			return errors.New("test")
		})

		cw.Cron().Start()
		time.Sleep(time.Second + time.Millisecond*200)
		panicV := (interface{})(nil)
		errorV := error(nil)
		cw.SetPanicHandler(func(v interface{}) { log.Printf("panic: %v", v); panicV = v })
		cw.SetErrorHandler(func(err error) { log.Printf("error: %v", err); errorV = err })
		time.Sleep(time.Second + time.Millisecond*200)
		xtesting.Equal(t, panicV, "test")
		xtesting.Equal(t, errorV.Error(), "test")
		<-cw.Cron().Stop().Done()
	})
}

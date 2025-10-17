package scheduler

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithLogger(cron.VerbosePrintfLogger(log.Default())),
		),
	}
}

func (s *Scheduler) At(hour, minute, second int) *TimeTaskBuilder {
	return &TimeTaskBuilder{
		scheduler: s,
		hour:      hour,
		minute:    minute,
		second:    second,
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

type TimeTaskBuilder struct {
	scheduler *Scheduler
	hour      int
	minute    int
	second    int
	id        string
}

func (ttb *TimeTaskBuilder) WithID(id string) *TimeTaskBuilder {
	ttb.id = id
	return ttb
}

func (ttb *TimeTaskBuilder) Do(fn func()) (cron.EntryID, error) {
	spec := cronSpec(ttb.hour, ttb.minute, ttb.second)
	return ttb.scheduler.cron.AddFunc(spec, fn)
}

func cronSpec(hour, minute, second int) string {
	return fmt.Sprintf("%d %d %d * * *", second, minute, hour)
}

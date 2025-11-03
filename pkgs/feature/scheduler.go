package scheduler

import (
	"fmt"
	"sync"

	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/robfig/cron/v3"
)

type TaskInfo struct {
	ID       string
	Name     string
	Schedule string
	EntryID  cron.EntryID
}

type Scheduler struct {
	cron  *cron.Cron
	tasks map[string]*TaskInfo
	mu    sync.RWMutex
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
		),
		tasks: make(map[string]*TaskInfo),
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

func (s *Scheduler) RegisterTask(id, name, schedule string, entryID cron.EntryID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[id] = &TaskInfo{
		ID:       id,
		Name:     name,
		Schedule: schedule,
		EntryID:  entryID,
	}
}

func (s *Scheduler) GetTasks() map[string]*TaskInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]*TaskInfo)
	for k, v := range s.tasks {
		result[k] = v
	}
	return result
}

func (s *Scheduler) GetCronEntries() map[cron.EntryID]cron.Entry {
	entries := make(map[cron.EntryID]cron.Entry)
	for _, entry := range s.cron.Entries() {
		entries[entry.ID] = entry
	}
	return entries
}

func (s *Scheduler) TriggerTask(id string) bool {
	s.mu.RLock()
	task, exists := s.tasks[id]
	s.mu.RUnlock()
	if !exists {
		logger.Warn(fmt.Sprintf("TriggerTask: task %s not found", id))
		return false
	}
	entry := s.cron.Entry(task.EntryID)
	if entry.ID == 0 {
		logger.Warn(fmt.Sprintf("TriggerTask: entry not found for task %s (EntryID: %d)", id, task.EntryID))
		return false
	}
	logger.Info(fmt.Sprintf("TriggerTask: triggering task %s (%s)", id, task.Name))
	go entry.Job.Run()
	return true
}

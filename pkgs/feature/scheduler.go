package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type TaskFunc func()

type Task struct {
	ID       string
	Interval time.Duration
	Fn       TaskFunc
	cancel   context.CancelFunc
	running  bool
	mu       sync.Mutex
}

type Scheduler struct {
	tasks map[string]*Task
	mu    sync.RWMutex
	ctx   context.Context
}

func New() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
		ctx:   context.Background(),
	}
}

func (s *Scheduler) AddTask(id string, interval time.Duration, fn TaskFunc) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[id]; exists {
		task.Stop()
	}

	task := &Task{
		ID:       id,
		Interval: interval,
		Fn:       fn,
	}

	s.tasks[id] = task
	return task
}

func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[id]; exists {
		task.Stop()
		delete(s.tasks, id)
	}
}

func (s *Scheduler) GetTask(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	return task, exists
}

func (s *Scheduler) Start(id string) error {
	s.mu.RLock()
	task, exists := s.tasks[id]
	s.mu.RUnlock()

	if !exists {
		return nil
	}

	task.Start()
	return nil
}

func (s *Scheduler) Stop(id string) {
	s.mu.RLock()
	task, exists := s.tasks[id]
	s.mu.RUnlock()

	if exists {
		task.Stop()
	}
}

func (s *Scheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		task.Stop()
	}
}

func (s *Scheduler) StartAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		task.Start()
	}
}

func (t *Task) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	t.running = true

	go func() {
		ticker := time.NewTicker(t.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t.Fn()
			}
		}
	}()
}

func (t *Task) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return
	}

	if t.cancel != nil {
		t.cancel()
	}
	t.running = false
}

func (t *Task) IsRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

func (t *Task) UpdateInterval(interval time.Duration) {
	t.mu.Lock()
	wasRunning := t.running
	t.mu.Unlock()

	if wasRunning {
		t.Stop()
		t.Interval = interval
		t.Start()
	} else {
		t.Interval = interval
	}
}

func EverySecond(fn TaskFunc) *Task {
	return &Task{
		Interval: time.Second,
		Fn:       fn,
	}
}

func EveryMinute(fn TaskFunc) *Task {
	return &Task{
		Interval: time.Minute,
		Fn:       fn,
	}
}

func EveryHour(fn TaskFunc) *Task {
	return &Task{
		Interval: time.Hour,
		Fn:       fn,
	}
}

func EveryDay(fn TaskFunc) *Task {
	return &Task{
		Interval: 24 * time.Hour,
		Fn:       fn,
	}
}

func Every(duration time.Duration, fn TaskFunc) *Task {
	return &Task{
		Interval: duration,
		Fn:       fn,
	}
}

func (s *Scheduler) Every(duration time.Duration) *TaskBuilder {
	return &TaskBuilder{
		scheduler: s,
		interval:  duration,
	}
}

type TaskBuilder struct {
	scheduler *Scheduler
	interval  time.Duration
	id        string
}

func (tb *TaskBuilder) WithID(id string) *TaskBuilder {
	tb.id = id
	return tb
}

func (tb *TaskBuilder) Do(fn TaskFunc) *Task {
	if tb.id == "" {
		tb.id = generateTaskID()
	}
	task := tb.scheduler.AddTask(tb.id, tb.interval, fn)
	task.Start()
	return task
}

func generateTaskID() string {
	return time.Now().Format("20060102150405.000000000")
}

func At(hour, minute, second int, fn TaskFunc) *Task {
	return &Task{
		Interval: 24 * time.Hour,
		Fn: func() {
			now := time.Now()
			target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

			if now.After(target) {
				target = target.Add(24 * time.Hour)
			}

			duration := target.Sub(now)
			time.Sleep(duration)
			fn()
		},
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

func (ttb *TimeTaskBuilder) Do(fn TaskFunc) *Task {
	if ttb.id == "" {
		ttb.id = generateTaskID()
	}

	task := ttb.scheduler.AddTask(ttb.id, 24*time.Hour, func() {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(), ttb.hour, ttb.minute, ttb.second, 0, now.Location())

		if now.After(target) {
			target = target.Add(24 * time.Hour)
		}

		duration := target.Sub(now)
		log.Printf("[Scheduler] Task %s will execute in %v at %s", ttb.id, duration, target.Format("2006-01-02 15:04:05"))

		time.Sleep(duration)
		fn()
	})

	task.Start()
	return task
}

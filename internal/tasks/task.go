package tasks

import (
	"log"

	"github.com/crayon/bot_golang/internal/config"
	scheduler "github.com/crayon/bot_golang/pkgs/feature"
)

type Task interface {
	Name() string
	Schedule(sched *scheduler.Scheduler, cfg *config.Config) error
}

func RegisterAll(sched *scheduler.Scheduler, cfg *config.Config) {
	tasks := []Task{
		&TechPushTask{},
	}

	for _, task := range tasks {
		if err := task.Schedule(sched, cfg); err != nil {
			log.Printf("Failed to schedule task %s: %v", task.Name(), err)
		} else {
			log.Printf("Scheduled task: %s", task.Name())
		}
	}
}

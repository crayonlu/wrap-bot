package tasks

import (
	"github.com/crayon/wrap-bot/internal/config"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Task interface {
	Name() string
	Schedule(sched *scheduler.Scheduler, cfg *config.Config) error
}

func RegisterAll(sched *scheduler.Scheduler, cfg *config.Config) {
	tasks := []Task{
		NewTechPushTask(cfg),
		NewRssPushTask(cfg),
	}

	for _, task := range tasks {
		if err := task.Schedule(sched, cfg); err != nil {
			logger.Error("Failed to schedule task " + task.Name() + ": " + err.Error())
		} else {
			logger.Info("Scheduled task: " + task.Name())
		}
	}
}

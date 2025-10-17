package tasks

import (
	"log"

	"github.com/crayon/bot_golang/internal/config"
	scheduler "github.com/crayon/bot_golang/pkgs/feature"
	"github.com/crayon/bot_golang/pkgs/feature/rss"
)

type RssPushTask struct {
	service *rss.RssPush
}

func NewRssPushTask(cfg *config.Config) *RssPushTask {
	return &RssPushTask{
		service: rss.NewRssPush(cfg),
	}
}

func (t *RssPushTask) Name() string {
	return "rss-push-daily"
}

func (t *RssPushTask) Schedule(sched *scheduler.Scheduler, cfg *config.Config) error {
	if cfg.RSSApiHost == "" {
		log.Println("RssPushTask: RSS_API_HOST not set, skipping")
		return nil
	}

	if len(cfg.RssPushGroups) == 0 && len(cfg.RssPushUsers) == 0 {
		log.Println("RssPushTask: no target groups or users configured, skipping")
		return nil
	}

	_, err := sched.At(13, 0, 0).WithID(t.Name()).Do(func() {
		if err := t.service.SendRssPush(); err != nil {
			log.Printf("[Rss Push]RssPushTask execution failed: %v", err)
		} else {
			log.Println("[Rss Push]RssPushTask executed successfully")
		}
	})

	return err
}

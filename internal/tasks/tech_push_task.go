package tasks

import (
	"log"

	"github.com/crayon/bot_golang/internal/config"
	scheduler "github.com/crayon/bot_golang/pkgs/feature"
	"github.com/crayon/bot_golang/pkgs/feature/tech_push"
)

type TechPushTask struct {
	cache map[string][]byte
}

func (t *TechPushTask) Name() string {
	return "tech-push-daily"
}

func (t *TechPushTask) Schedule(sched *scheduler.Scheduler, cfg *config.Config) error {
	if cfg.HotApiHost == "" || cfg.HotApiKey == "" {
		log.Println("TechPushTask: HOT_API_HOST or HOT_API_KEY not set, skipping")
		return nil
	}

	if len(cfg.TechPushGroups) == 0 && len(cfg.TechPushUsers) == 0 {
		log.Println("TechPushTask: no target groups or users configured, skipping")
		return nil
	}

	t.cache = make(map[string][]byte)

	sched.At(8, 0, 0).WithID(t.Name()).Do(func() {
		if err := tech_push.SendTechPush(cfg, t.cache); err != nil {
			log.Printf("TechPushTask execution failed: %v", err)
		} else {
			log.Println("TechPushTask executed successfully")
		}
	})

	return nil
}

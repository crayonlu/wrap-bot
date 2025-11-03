package tasks

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/feature/tech_push"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type TechPushTask struct {
	cache   map[string][]byte
	service *tech_push.TechPush
}

func NewTechPushTask(cfg *config.Config) *TechPushTask {
	var aiAnalyzer tech_push.AIAnalyzer

	if cfg.AIEnabled {
		aiService := ai.NewService(ai.Config{
			APIURL:           cfg.AIURL,
			APIKey:           cfg.AIKey,
			Model:            cfg.AIModel,
			SystemPromptPath: cfg.SystemPromptPath,
			MaxHistory:       5,
			Temperature:      0.7,
			TopP:             0.9,
			MaxTokens:        2000,
		})

		aiAnalyzer = ai.NewAnalyzer(ai.AnalyzerConfig{
			Service:            aiService,
			AnalyzerPromptPath: cfg.AnalyzerPromptPath,
		})
	}

	return &TechPushTask{
		cache:   make(map[string][]byte),
		service: tech_push.NewTechPush(cfg, aiAnalyzer),
	}
}

func (t *TechPushTask) Name() string {
	return "tech-push-daily"
}

func (t *TechPushTask) Schedule(sched *scheduler.Scheduler, cfg *config.Config) error {
	if cfg.HotApiHost == "" || cfg.HotApiKey == "" {
		logger.Warn("TechPushTask: HOT_API_HOST or HOT_API_KEY not set, skipping")
		return nil
	}

	if len(cfg.TechPushGroups) == 0 && len(cfg.TechPushUsers) == 0 {
		logger.Warn("TechPushTask: no target groups or users configured, skipping")
		return nil
	}

	entryID, err := sched.At(12, 0, 0).WithID(t.Name()).Do(func() {
		if err := t.service.SendTechPush(t.cache); err != nil {
			logger.Error(fmt.Sprintf("TechPushTask execution failed: %v", err))
		} else {
			logger.Info("TechPushTask executed successfully")
		}
	})

	if err == nil {
		sched.RegisterTask(t.Name(), "Tech Push Daily", "0 0 12 * * *", entryID)
	}

	return err
}

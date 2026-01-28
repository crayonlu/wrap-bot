package tasks

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
	aiconfig "github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/crayon/wrap-bot/pkgs/feature/analyzer"
	"github.com/crayon/wrap-bot/pkgs/feature/tech_push"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type TechPushTask struct {
	cache   map[string][]byte
	service *tech_push.TechPush
}

func NewTechPushTask(cfg *config.Config) *TechPushTask {
	var aiAnalyzer *analyzer.Analyzer

	if cfg.AIEnabled {
		aiCfg := &aiconfig.Config{
			APIURL: cfg.AIURL,
			APIKey: cfg.AIKey,

			Model: cfg.AIModel,

			Temperature: 0.7,
			TopP:        0.9,
			MaxTokens:   2000,
			MaxHistory:  5,

			ToolsEnabled: cfg.AIToolsEnabled,

			SystemPromptPath: cfg.SystemPromptPath,
			SerpAPIKey:       cfg.SerpAPIKey,
			WeatherAPIKey:    cfg.WeatherAPIKey,
		}

		factory := factory.NewFactory(aiCfg)
		chatAgent := factory.CreateAgent()

		aiAnalyzer = analyzer.NewAnalyzer(analyzer.AnalyzerConfig{
			Agent:              chatAgent,
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
		sched.RegisterTask(t.Name(), "Tech Push Daily", "每日推送技术新闻", "0 0 12 * * *", entryID)
	}

	return err
}

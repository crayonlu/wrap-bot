package tasks

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
	aiconfig "github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/crayon/wrap-bot/pkgs/feature/analyzer"
	"github.com/crayon/wrap-bot/pkgs/feature/rss"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type RssPushTask struct {
	service *rss.RssPush
}

func NewRssPushTask(cfg *config.Config) *RssPushTask {
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

	return &RssPushTask{
		service: rss.NewRssPush(cfg, aiAnalyzer),
	}
}

func (t *RssPushTask) Name() string {
	return "rss-push-daily"
}

func (t *RssPushTask) Schedule(sched *scheduler.Scheduler, cfg *config.Config) error {
	if cfg.RSSApiHost == "" {
		logger.Warn("RssPushTask: RSS_API_HOST not set, skipping")
		return nil
	}

	if len(cfg.RssPushGroups) == 0 && len(cfg.RssPushUsers) == 0 {
		logger.Warn("RssPushTask: no target groups or users configured, skipping")
		return nil
	}

	entryID, err := sched.At(13, 0, 0).WithID(t.Name()).Do(func() {
		if err := t.service.SendRssPush(); err != nil {
			logger.Error(fmt.Sprintf("[Rss Push]RssPushTask execution failed: %v", err))
		}
	})

	if err == nil {
		sched.RegisterTask(t.Name(), "RSS Daily Push", "每日推送 RSS 订阅内容", "0 0 13 *", entryID)
	}

	return err
}

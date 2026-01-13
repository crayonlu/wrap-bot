package plugins

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	aiconfig "github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/crayon/wrap-bot/pkgs/feature/analyzer"
	"github.com/crayon/wrap-bot/pkgs/feature/rss"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

var rssPushService *rss.RssPush

func RssPushPlugin(cfg *config.Config) bot.HandlerFunc {
	var aiAnalyzer *analyzer.Analyzer

	if cfg.AIEnabled {
		aiCfg := &aiconfig.Config{
			APIURL:           cfg.AIURL,
			APIKey:           cfg.AIKey,
			TextModel:        cfg.AITextModel,
			VisionModel:      cfg.AIVisionModel,
			Temperature:      0.7,
			TopP:             0.9,
			MaxTokens:        2000,
			MaxHistory:       5,
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

	rssPushService = rss.NewRssPush(cfg, aiAnalyzer)

	return func(ctx *bot.Context) {
		if ctx.Event.RawMessage == "/rss" {
			go func() {
				if err := rssPushService.SendRssPush(); err != nil {
					logger.Error(fmt.Sprintf("RSS push failed: %v", err))
				} else {
					logger.Info("RSS push succeeded")
				}
			}()
			return
		}
		ctx.Next()
	}
}

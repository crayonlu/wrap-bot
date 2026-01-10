package plugins

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/feature/rss"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

var rssPushService *rss.RssPush

func RssPushPlugin(cfg *config.Config) bot.HandlerFunc {
	var aiAnalyzer rss.AIAnalyzer

	if cfg.AIEnabled {
		aiService := ai.NewService(ai.Config{
			APIURL:           cfg.AIURL,
			APIKey:           cfg.AIKey,
			TextModel:        cfg.AITextModel,
			VisionModel:      cfg.AIVisionModel,
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

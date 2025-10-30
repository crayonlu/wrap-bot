package plugins

import (
	"log"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/feature/rss"
)

var rssPushService *rss.RssPush

func RssPushPlugin(cfg *config.Config) bot.HandlerFunc {
	var aiAnalyzer rss.AIAnalyzer

	if cfg.AIEnabled {
		aiService := ai.NewService(ai.Config{
			APIURL:           cfg.AIURL,
			APIKey:           cfg.AIKey,
			Model:            cfg.AIModel,
			SystemPromptPath: cfg.AnalyzerPromptPath,
			MaxHistory:       5,
			Temperature:      0.7,
			TopP:             0.9,
			MaxTokens:        2000,
		})
		aiAnalyzer = ai.NewAnalyzer(aiService)
	}

	rssPushService = rss.NewRssPush(cfg, aiAnalyzer)

	return func(ctx *bot.Context) {
		if ctx.Event.RawMessage == "/rss" {
			go func() {
				if err := rssPushService.SendRssPush(); err != nil {
					log.Printf("RSS push failed: %v", err)
				} else {
					log.Printf("RSS push succeeded")
				}
			}()
			return

		}
		ctx.Next()
	}
}

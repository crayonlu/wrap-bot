package plugins

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/feature/tech_push"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

var techPushCache = make(map[string][]byte)
var techPushService *tech_push.TechPush

func TechPushPlugin(cfg *config.Config) bot.HandlerFunc {
	var aiAnalyzer tech_push.AIAnalyzer

	if cfg.AIEnabled {
		aiService := ai.NewService(ai.Config{
			APIURL:           cfg.AIURL,
			APIKey:           cfg.AIKey,
			Model:            cfg.AIModel,
			SystemPromptPath: cfg.SystemPromptPath,
			MaxHistory:       20,
			Temperature:      0.7,
			TopP:             0.9,
			MaxTokens:        2000,
		})
		aiAnalyzer = ai.NewAnalyzer(aiService)
	}

	techPushService = tech_push.NewTechPush(cfg, aiAnalyzer)

	return func(ctx *bot.Context) {
		if ctx.Event.RawMessage == "/tech" {
			go func() {
				if err := techPushService.SendTechPush(techPushCache); err != nil {
					logger.Error(fmt.Sprintf("Tech push failed: %v", err))
				}
			}()
			return
		}
		ctx.Next()
	}
}

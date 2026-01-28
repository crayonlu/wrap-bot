package plugins

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	aiconfig "github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/crayon/wrap-bot/pkgs/feature/analyzer"
	"github.com/crayon/wrap-bot/pkgs/feature/tech_push"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

var techPushCache = make(map[string][]byte)
var techPushService *tech_push.TechPush

func TechPushPlugin(cfg *config.Config) bot.HandlerFunc {
	var aiAnalyzer *analyzer.Analyzer

	if cfg.AIEnabled {
		aiCfg := &aiconfig.Config{
			APIURL:           cfg.AIURL,
			APIKey:           cfg.AIKey,
			Model:            cfg.AIModel,
			Temperature:      0.7,
			TopP:             0.9,
			MaxTokens:        2000,
			MaxHistory:       20,
			ToolsEnabled:     cfg.AIToolsEnabled,
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

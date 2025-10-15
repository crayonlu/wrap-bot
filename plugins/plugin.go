package plugins

import (
	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
	tech_push "github.com/crayon/bot_golang/plugins/tech_push"
)

func Register(engine *bot.Engine, cfg *config.Config) {
	engine.Use(PingPlugin(cfg))
	engine.Use(EchoPlugin(cfg))
	if cfg.AIEnabled {
		engine.Use(AIChatPlugin(cfg))
	}
	if cfg.HotApiHost != "" && cfg.HotApiKey != "" {
		engine.Use(tech_push.TechPushPlugin(cfg))
	}
}

package plugins

import (
	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
)

func Register(engine *bot.Engine, cfg *config.Config) {
	engine.Use(PingPlugin(cfg))
	engine.Use(EchoPlugin(cfg))
	if cfg.AIEnabled {
		engine.Use(AIChatPlugin(cfg))
	}
	if cfg.HotApiHost != "" && cfg.HotApiKey != "" {
		engine.Use(TechPushPlugin(cfg))
	}
	if cfg.RSSApiHost != "" {
		engine.Use(RssPushPlugin(cfg))
	}
}

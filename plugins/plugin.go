package plugins

import (
	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
)

func Register(engine *bot.Engine, cfg *config.Config) {
	engine.Use(PingPlugin(cfg))
	engine.Use(EchoPlugin(cfg))
	engine.Use(HelpPlugin(cfg))
	engine.Use(AIChatPlugin(cfg))
}

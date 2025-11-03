package plugins

import (
	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
)

func Register(engine *bot.Engine, cfg *config.Config) {
	engine.RegisterPlugin("ping", "Simple ping-pong command", PingPlugin(cfg))
	engine.RegisterPlugin("echo", "Echo back user messages", EchoPlugin(cfg))
	engine.RegisterPlugin("help", "Show available commands", HelpPlugin(cfg))

	if cfg.AIEnabled {
		engine.RegisterPlugin("ai_chat", "AI conversation plugin", AIChatPlugin(cfg))
	}
	if cfg.HotApiHost != "" && cfg.HotApiKey != "" {
		engine.RegisterPlugin("tech_push", "Tech news push service", TechPushPlugin(cfg))
	}
	if cfg.RSSApiHost != "" {
		engine.RegisterPlugin("rss_push", "RSS feed push service", RssPushPlugin(cfg))
	}
}

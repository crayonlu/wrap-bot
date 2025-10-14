package plugins

import (
	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
	scheduler "github.com/crayon/bot_golang/pkgs/feature"
	tech_push "github.com/crayon/bot_golang/plugins/tech_push"
)

func Register(engine *bot.Engine, cfg *config.Config, sched *scheduler.Scheduler) {
	engine.Use(PingPlugin(cfg))
	engine.Use(EchoPlugin(cfg))
	if cfg.AIEnabled {
		engine.Use(AIChatPlugin(cfg))
	}
	if cfg.HotApiHost != "" && cfg.HotApiKey != "" {
		engine.Use(tech_push.TechPushPlugin(cfg, sched))
	}
}

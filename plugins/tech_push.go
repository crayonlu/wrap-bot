package plugins

import (
	"log"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
	"github.com/crayon/bot_golang/pkgs/feature/tech_push"
)

var techPushCache = make(map[string][]byte)

func TechPushPlugin(cfg *config.Config) bot.HandlerFunc {
	return func(ctx *bot.Context) {
		if ctx.Event.RawMessage == "/tech" {
			go func() {
				if err := tech_push.SendTechPush(cfg, techPushCache); err != nil {
					log.Printf("Tech push failed: %v", err)
				}
			}()
			return
		}
		ctx.Next()
	}
}

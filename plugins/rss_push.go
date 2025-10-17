package plugins

import (
	"log"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
	"github.com/crayon/bot_golang/pkgs/feature/rss"
)

var rssPushService *rss.RssPush

func RssPushPlugin(cfg *config.Config) bot.HandlerFunc {
	rssPushService = rss.NewRssPush(cfg)

	return func(ctx *bot.Context) {
		if ctx.Event.RawMessage == "/rss" {
			go func() {
				if err := rssPushService.SendRssPush(); err != nil {
					log.Printf("[Rss Push]Error details: " + err.Error())
				} else {
					log.Printf("[Rss Push]RSS推送成功！")
				}
			}()
			return
		}
		ctx.Next()
	}
}

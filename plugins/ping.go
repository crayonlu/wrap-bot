package plugins

import (
	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
)

func PingPlugin(cfg *config.Config) bot.HandlerFunc {
	return bot.OnCommand(cfg.CommandPrefix, "ping", func(ctx *bot.Context) {
		ctx.ReplyText("pong!")
	})
}

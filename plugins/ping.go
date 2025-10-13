package plugins

import (
	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
)

func PingPlugin(cfg *config.Config) bot.HandlerFunc {
	return bot.OnCommand(cfg.CommandPrefix, "ping", func(ctx *bot.Context) {
		ctx.ReplyText("pong!")
	})
}

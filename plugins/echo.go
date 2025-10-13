package plugins

import (
	"strings"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
)

func EchoPlugin(cfg *config.Config) bot.HandlerFunc {
	return bot.OnCommand(cfg.CommandPrefix, "echo", func(ctx *bot.Context) {
		text := ctx.Event.GetText()
		prefix := cfg.CommandPrefix + "echo"
		
		if len(text) <= len(prefix) {
			ctx.ReplyText("Usage: /echo <message>")
			return
		}
		
		message := strings.TrimSpace(text[len(prefix):])
		ctx.ReplyText(message)
	})
}

package plugins

import (
	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
)

func HelpPlugin(cfg *config.Config) bot.HandlerFunc {
	return bot.OnCommand(cfg.CommandPrefix, "help", func(ctx *bot.Context) {
		help := `Available commands:
/ping - Check if bot is alive
/echo <message> - Echo your message
/help - Show this help message`

		ctx.ReplyText(help)
	})
}

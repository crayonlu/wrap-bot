package plugins

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

func AIChatPlugin(cfg *config.Config) bot.HandlerFunc {
	if !cfg.AIEnabled {
		return func(ctx *bot.Context) {}
	}

	aiService := ai.NewService(ai.Config{
		APIURL:           cfg.AIURL,
		APIKey:           cfg.AIKey,
		Model:            cfg.AIModel,
		SystemPromptPath: cfg.SystemPromptPath,
		MaxHistory:       20,
		Temperature:      0.7,
		TopP:             0.9,
		MaxTokens:        2000,
	})

	return func(ctx *bot.Context) {
		if !ctx.Event.IsGroupMessage() && !ctx.Event.IsPrivateMessage() {
			return
		}

		text := ctx.Event.GetText()
		if text == "" {
			return
		}

		if len(text) > 0 && text[0] == '/' {
			return
		}

		conversationID := fmt.Sprintf("%d_%d", ctx.Event.GroupID, ctx.Event.UserID)
		if ctx.Event.IsPrivateMessage() {
			conversationID = fmt.Sprintf("private_%d", ctx.Event.UserID)
		}

		if text == "清除历史" || text == "reset" {
			aiService.ClearHistory(conversationID)
			ctx.ReplyText("空空如也了")
			return
		}

		response, err := aiService.Chat(conversationID, text, true)
		if err != nil {
			logger.Error(fmt.Sprintf("AI chat error: %v", err))
			ctx.ReplyText("坠机了嘻嘻...")
			return
		}

		if ctx.Event.IsGroupMessage() {
			ctx.ReplyAt(response)
		} else {
			ctx.ReplyText(response)
		}
	}
}

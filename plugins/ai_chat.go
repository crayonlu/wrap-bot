package plugins

import (
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/crayon/wrap-bot/pkgs/napcat"
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
		ToolsEnabled:     cfg.AIToolsEnabled,
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

		var response *ai.ChatResult
		var err error

		if cfg.AIVisionEnabled {
			imageURLs := ctx.Event.GetImages()
			if len(imageURLs) > 0 {
				response, err = aiService.ChatWithImages(conversationID, text, imageURLs, cfg.AIImageDetail, true)
			} else {
				response, err = aiService.Chat(conversationID, text, true)
			}
		} else {
			response, err = aiService.Chat(conversationID, text, true)
		}

		if err != nil {
			logger.Error(fmt.Sprintf("AI chat error: %v", err))
			ctx.ReplyText("坠机了嘻嘻...")
			return
		}

		if response.Thinking != "" {
			apiClient := ctx.GetAPIClient()
			if napcatClient, ok := apiClient.(*napcat.Client); ok {
				thinkingMsg := fmt.Sprintf("thinking: \n---\n%s\n---", response.Thinking)
				node := napcat.NewMixedForwardNode(
					"AI Thinking",
					ctx.Event.SelfID,
					napcat.NewTextSegment(thinkingMsg),
				)

				if ctx.Event.IsGroupMessage() {
					_, err := napcatClient.SendGroupForwardMsg(ctx.Event.GroupID, []napcat.ForwardNode{node})
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to send forward message: %v", err))
					}
				} else {
					_, err := napcatClient.SendPrivateForwardMsg(ctx.Event.UserID, []napcat.ForwardNode{node})
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to send forward message: %v", err))
					}
				}
			}
		}

		if ctx.Event.IsGroupMessage() {
			ctx.ReplyAt(response.Content)
		} else {
			ctx.ReplyText(response.Content)
		}
	}
}

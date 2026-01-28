package plugins

import (
	"context"
	"fmt"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
	aiconfig "github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/crayon/wrap-bot/pkgs/napcat"
)

func AIChatPlugin(cfg *config.Config) bot.HandlerFunc {
	if !cfg.AIEnabled {
		return func(ctx *bot.Context) {}
	}

	aiCfg := &aiconfig.Config{
		APIURL: cfg.AIURL,
		APIKey: cfg.AIKey,

		Model: cfg.AIModel,

		Temperature:      0.7,
		TopP:             0.9,
		MaxTokens:        2000,
		MaxHistory:       20,
		SystemPromptPath: cfg.SystemPromptPath,
		SerpAPIKey:       cfg.SerpAPIKey,
		WeatherAPIKey:    cfg.WeatherAPIKey,
		ToolsEnabled:     cfg.AIToolsEnabled,
	}

	factory := factory.NewFactory(aiCfg)
	chatAgent := factory.CreateAgent()

	logger.Info(fmt.Sprintf("[AIChatPlugin] Initialized with %d tools",
		len(aiCfg.ToolsEnabled)))

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
			chatAgent.ClearHistory(conversationID)
			ctx.ReplyText("空空如也了")
			return
		}

		var response *agent.ChatResult
		var err error

		imageURLs := ctx.Event.GetImages()
		if len(imageURLs) > 0 {
			response, err = chatAgent.ChatWithImages(context.Background(), conversationID, text, imageURLs)
		} else {
			response, err = chatAgent.Chat(context.Background(), conversationID, text)
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

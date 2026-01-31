package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	aiconfig "github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/crayon/wrap-bot/pkgs/feature/chat_explainer"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/crayon/wrap-bot/pkgs/napcat"
)

func ChatExplainerPlugin(cfg *config.Config) bot.HandlerFunc {
	if !cfg.AIEnabled {
		return func(ctx *bot.Context) {}
	}

	aiCfg := &aiconfig.Config{
		APIURL:           cfg.AIURL,
		APIKey:           cfg.AIKey,
		Model:            cfg.AIModel,
		Temperature:      0.7,
		TopP:             0.9,
		MaxTokens:        2000,
		MaxHistory:       0,
		SystemPromptPath: cfg.SystemPromptPath,
		ToolsEnabled:     []string{},
	}

	factory := factory.NewFactory(aiCfg)
	chatAgent := factory.CreateAgent()

	systemPrompt := ""
	if cfg.SystemPromptPath != "" {
		data, err := os.ReadFile(cfg.SystemPromptPath)
		if err == nil {
			systemPrompt = string(data)
		}
	}

	analyzer := chat_explainer.NewAnalyzer(chatAgent, systemPrompt, 20)
	parser := chat_explainer.NewParser()

	return func(ctx *bot.Context) {
		if !ctx.Event.IsGroupMessage() && !ctx.Event.IsPrivateMessage() {
			return
		}

		eventData, err := json.Marshal(ctx.Event)
		if err != nil {
			return
		}

		var eventMap map[string]interface{}
		if err := json.Unmarshal(eventData, &eventMap); err != nil {
			return
		}

		if !parser.IsForwardMessage(eventMap) {
			return
		}

		forwardID := parser.GetForwardID(eventMap)
		if forwardID == "" {
			ctx.ReplyText("无法读取合并转发消息")
			return
		}

		apiClient := ctx.GetAPIClient()
		napcatClient, ok := apiClient.(*napcat.Client)
		if !ok {
			ctx.ReplyText("API客户端不可用")
			return
		}

		forwardData, err := napcatClient.GetForwardMsg(forwardID)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get forward message: %v", err))
			ctx.ReplyText("无法读取合并转发消息内容")
			return
		}

		forwardedChat, err := parser.ParseForwardMessage(forwardData)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse forward message: %v", err))
			ctx.ReplyText("解析消息失败")
			return
		}

		if len(forwardedChat.Messages) == 0 {
			ctx.ReplyText("消息内容为空")
			return
		}

		ctx.ReplyText("正在分析对话，请稍候...")

		analysis, err := analyzer.AnalyzeAll(context.Background(), forwardedChat.Messages)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to analyze chat: %v", err))
			ctx.ReplyText("分析失败，请稍后重试")
			return
		}

		nodes := chat_explainer.BuildForwardNodes(
			forwardedChat.Messages,
			analysis.MessageAnalyses,
			analysis.Summary,
			ctx.Event.SelfID,
		)

		if ctx.Event.IsGroupMessage() {
			_, err = napcatClient.SendGroupForwardMsg(ctx.Event.GroupID, nodes)
		} else {
			_, err = napcatClient.SendPrivateForwardMsg(ctx.Event.UserID, nodes)
		}

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to send forward message: %v", err))
			ctx.ReplyText("发送解读结果失败")
		}
	}
}

package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/feature/ai"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/memory"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/provider"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/tool"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type AgentConfig struct {
	Provider     provider.LLMProvider
	History      *memory.HistoryManager
	ToolRegistry tool.ToolRegistry
	SystemPrompt string
	TextModel    string
	VisionModel  string
	Temperature  float64
	TopP         float64
	MaxTokens    int
	TextTools    []string
	VisionTools  []string
}

type ChatAgent struct {
	config AgentConfig
}

type ChatResult struct {
	Thinking string
	Content  string
}

func NewChatAgent(cfg AgentConfig) *ChatAgent {
	return &ChatAgent{
		config: cfg,
	}
}

func (a *ChatAgent) Chat(ctx context.Context, conversationID, message string) (*ChatResult, error) {
	logger.Info(fmt.Sprintf("[Chat] ConversationID: %s, Message: %s", conversationID, message))

	userMsg := memory.Message{Role: "user", Content: message}
	a.config.History.AddMessage(conversationID, userMsg)

	messages := []memory.Message{{Role: "system", Content: a.config.SystemPrompt}}
	history, _ := a.config.History.GetHistory(conversationID)
	messages = append(messages, history...)

	logger.Info(fmt.Sprintf("[Chat] History size: %d messages", len(history)))

	req := ai.ChatRequest{
		Model:       a.config.TextModel,
		Messages:    convertMessagesToChatRequest(messages),
		Stream:      false,
		Temperature: a.config.Temperature,
		TopP:        a.config.TopP,
		MaxTokens:   a.config.MaxTokens,
	}

	tools := a.getToolsForModel("text")
	if len(tools) > 0 {
		req.Tools = convertToolsToChatRequest(tools)
		logger.Info(fmt.Sprintf("[Chat] Added %d tool(s) to request", len(tools)))
	}

	resp, err := a.config.Provider.Complete(ctx, req)
	if err != nil {
		logger.Error(fmt.Sprintf("[Chat] API request failed: %v", err))
		return nil, err
	}

	if len(resp.Choices) == 0 {
		logger.Error("[Chat] No response from AI")
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]

	logger.Info(fmt.Sprintf("[Chat] Choice: %+v", choice))

	if len(choice.Message.ToolCalls) > 0 {
		logger.Info(fmt.Sprintf("[Chat] Model requested %d tool call(s)", len(choice.Message.ToolCalls)))
		return a.handleToolCalls(ctx, conversationID, choice.Message)
	}

	assistantMsg := memory.Message{Role: "assistant", Content: choice.Message.Content}
	a.config.History.AddMessage(conversationID, assistantMsg)

	if contentStr, ok := choice.Message.Content.(string); ok {
		thinking := choice.Message.ReasoningContent
		content := contentStr

		if thinkContent, cleanContent := parseThinkTags(content); thinkContent != "" {
			if thinking != "" {
				thinking = thinking + "\n\n" + thinkContent
			} else {
				thinking = thinkContent
			}
			content = cleanContent
		}

		return &ChatResult{
			Thinking: thinking,
			Content:  content,
		}, nil
	}
	return nil, fmt.Errorf("unexpected content type in response")
}

func (a *ChatAgent) ChatWithImages(ctx context.Context, conversationID, message string, imageURLs []string) (*ChatResult, error) {
	logger.Info(fmt.Sprintf("[ChatWithImages] ConversationID: %s, Message: %s, Images: %d", conversationID, message, len(imageURLs)))

	var content interface{}
	if len(imageURLs) > 0 {
		contentItems := []ai.ContentItem{}
		for _, url := range imageURLs {
			contentItems = append(contentItems, ai.ContentItem{
				Type: "image_url",
				ImageURL: &ai.ImageURL{
					URL:    url,
					Detail: "auto",
				},
			})
		}
		if message != "" {
			contentItems = append(contentItems, ai.ContentItem{
				Type: "text",
				Text: message,
			})
		}
		content = contentItems
	} else {
		content = message
	}

	userMsg := memory.Message{Role: "user", Content: content}
	a.config.History.AddMessage(conversationID, userMsg)

	messages := []memory.Message{{Role: "system", Content: a.config.SystemPrompt}}
	history, _ := a.config.History.GetHistory(conversationID)
	filteredHistory := filterToolCallMessages(history)
	logger.Info(fmt.Sprintf("[ChatWithImages] History size: %d -> %d (after filtering tool calls)", len(history), len(filteredHistory)))
	messages = append(messages, filteredHistory...)

	req := ai.ChatRequest{
		Model:       a.config.VisionModel,
		Messages:    convertMessagesToChatRequest(messages),
		Stream:      false,
		Temperature: a.config.Temperature,
		TopP:        a.config.TopP,
		MaxTokens:   a.config.MaxTokens,
	}

	tools := a.getToolsForModel("vision")
	if len(tools) > 0 {
		req.Tools = convertToolsToChatRequest(tools)
		logger.Info(fmt.Sprintf("[ChatWithImages] Added %d tool(s) to request", len(tools)))
	}

	resp, err := a.config.Provider.Complete(ctx, req)
	if err != nil {
		logger.Error(fmt.Sprintf("[ChatWithImages] API request failed: %v", err))
		return nil, err
	}

	if len(resp.Choices) == 0 {
		logger.Error("[ChatWithImages] No response from AI")
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]
	logger.Info(fmt.Sprintf("[Chat] Choice: %+v", choice))

	if len(choice.Message.ToolCalls) > 0 {
		logger.Info(fmt.Sprintf("[ChatWithImages] Model requested %d tool call(s)", len(choice.Message.ToolCalls)))
		return a.handleToolCalls(ctx, conversationID, choice.Message)
	}

	assistantMsg := memory.Message{Role: "assistant", Content: choice.Message.Content}
	a.config.History.AddMessage(conversationID, assistantMsg)

	if contentStr, ok := choice.Message.Content.(string); ok {
		thinking := choice.Message.ReasoningContent
		content := contentStr

		if thinkContent, cleanContent := parseThinkTags(content); thinkContent != "" {
			if thinking != "" {
				thinking = thinking + "\n\n" + thinkContent
			} else {
				thinking = thinkContent
			}
			content = cleanContent
		}

		return &ChatResult{
			Thinking: thinking,
			Content:  content,
		}, nil
	}
	return nil, fmt.Errorf("unexpected content type in response")
}

func (a *ChatAgent) ClearHistory(conversationID string) error {
	return a.config.History.ClearHistory(conversationID)
}

func (a *ChatAgent) getToolsForModel(modelType string) []tool.Tool {
	var enabledToolNames []string

	switch modelType {
	case "text":
		enabledToolNames = a.config.TextTools
	case "vision":
		enabledToolNames = a.config.VisionTools
	default:
		return []tool.Tool{}
	}

	tools := a.config.ToolRegistry.GetAll()
	var result []tool.Tool
	for _, t := range tools {
		if t.Enabled && contains(enabledToolNames, t.Name) {
			result = append(result, t)
		}
	}

	return result
}

func (a *ChatAgent) handleToolCalls(ctx context.Context, conversationID string, assistantMsg ai.Message) (*ChatResult, error) {
	logger.Info(fmt.Sprintf("[ToolCall] Received %d tool call(s)", len(assistantMsg.ToolCalls)))

	assistantMemoryMsg := memory.Message{
		Role:      "assistant",
		Content:   assistantMsg.Content,
		ToolCalls: convertToolCallsToMemory(assistantMsg.ToolCalls),
	}
	a.config.History.AddMessage(conversationID, assistantMemoryMsg)

	for _, toolCall := range assistantMsg.ToolCalls {
		logger.Info(fmt.Sprintf("[ToolCall] Executing tool: %s with args: %s", toolCall.Function.Name, toolCall.Function.Arguments))

		result, err := a.config.ToolRegistry.Execute(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			logger.Error(fmt.Sprintf("[ToolCall] Tool %s execution failed: %v", toolCall.Function.Name, err))
			result = fmt.Sprintf("Error: %v", err)
		} else {
			logger.Info(fmt.Sprintf("[ToolCall] Tool %s executed successfully, result length: %d", toolCall.Function.Name, len(result)))
		}

		toolMsg := memory.Message{
			Role:       "tool",
			Content:    result,
			ToolCallID: toolCall.ID,
		}
		a.config.History.AddMessage(conversationID, toolMsg)
	}

	logger.Info("[ToolCall] Sending final request to model with tool results")

	messages := []ai.Message{{Role: "system", Content: a.config.SystemPrompt}}
	history, _ := a.config.History.GetHistory(conversationID)
	messages = append(messages, convertMessagesToChatRequest(history)...)

	req := ai.ChatRequest{
		Model:       a.config.TextModel,
		Messages:    messages,
		Stream:      false,
		Temperature: a.config.Temperature,
	}

	tools := a.getToolsForModel("text")
	if len(tools) > 0 {
		req.Tools = convertToolsToChatRequest(tools)
	}

	resp, err := a.config.Provider.Complete(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		if contentStr, ok := content.(string); ok {
			assistantMsg := memory.Message{Role: "assistant", Content: contentStr}
			a.config.History.AddMessage(conversationID, assistantMsg)

			thinking := resp.Choices[0].Message.ReasoningContent
			content := contentStr

			if thinkContent, cleanContent := parseThinkTags(content); thinkContent != "" {
				if thinking != "" {
					thinking = thinking + "\n\n" + thinkContent
				} else {
					thinking = thinkContent
				}
				content = cleanContent
			}

			return &ChatResult{
				Thinking: thinking,
				Content:  content,
			}, nil
		}
		return nil, fmt.Errorf("unexpected content type in response")
	}

	return nil, fmt.Errorf("no response after tool call")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func filterToolCallMessages(messages []memory.Message) []memory.Message {
	result := make([]memory.Message, 0, len(messages))
	for _, msg := range messages {
		if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
			cleanMsg := memory.Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
			result = append(result, cleanMsg)
			continue
		}
		if msg.Role == "tool" {
			continue
		}
		result = append(result, msg)
	}
	return result
}

func convertMessagesToChatRequest(messages []memory.Message) []ai.Message {
	result := make([]ai.Message, 0, len(messages))
	for _, msg := range messages {
		aiMsg := ai.Message{
			Role:       msg.Role,
			Content:    msg.Content,
			ToolCalls:  convertToolCallsToAI(msg.ToolCalls),
			ToolCallID: msg.ToolCallID,
		}
		result = append(result, aiMsg)
	}
	return result
}

func convertToolCallsToMemory(toolCalls []ai.ToolCall) []memory.ToolCall {
	result := make([]memory.ToolCall, 0, len(toolCalls))
	for _, tc := range toolCalls {
		result = append(result, memory.ToolCall{
			ID:   tc.ID,
			Type: tc.Type,
			Function: memory.FunctionCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		})
	}
	return result
}

func convertToolCallsToAI(toolCalls []memory.ToolCall) []ai.ToolCall {
	result := make([]ai.ToolCall, 0, len(toolCalls))
	for _, tc := range toolCalls {
		result = append(result, ai.ToolCall{
			ID:   tc.ID,
			Type: tc.Type,
			Function: ai.FunctionCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		})
	}
	return result
}

func convertToolsToChatRequest(tools []tool.Tool) []ai.Tool {
	result := make([]ai.Tool, 0, len(tools))
	for _, t := range tools {
		result = append(result, ai.Tool{
			Type: "function",
			Function: ai.FunctionDef{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}
	return result
}

func parseThinkTags(content string) (string, string) {
	re := regexp.MustCompile(`(?s)<think>(.*?)</think>`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		thinkContent := matches[1]
		cleanContent := re.ReplaceAllString(content, "")
		cleanContent = strings.TrimSpace(cleanContent)
		return thinkContent, cleanContent
	}

	return "", content
}

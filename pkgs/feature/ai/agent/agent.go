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
	userMsg := memory.Message{Role: "user", Content: message}
	a.config.History.AddMessage(conversationID, userMsg)

	messages := []memory.Message{{Role: "system", Content: a.config.SystemPrompt}}
	history, _ := a.config.History.GetHistory(conversationID)
	messages = append(messages, history...)

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
	}

	resp, err := a.config.Provider.Complete(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) > 0 {
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
	messages = append(messages, history...)

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
	}

	resp, err := a.config.Provider.Complete(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) > 0 {
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
	assistantMemoryMsg := memory.Message{
		Role:    "assistant",
		Content: assistantMsg.Content,
	}
	a.config.History.AddMessage(conversationID, assistantMemoryMsg)

	for _, toolCall := range assistantMsg.ToolCalls {
		result, err := a.config.ToolRegistry.Execute(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			result = fmt.Sprintf("Error: %v", err)
		}

		toolMsg := memory.Message{
			Role:    "tool",
			Content: result,
		}
		a.config.History.AddMessage(conversationID, toolMsg)
	}

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

func convertMessagesToChatRequest(messages []memory.Message) []ai.Message {
	result := make([]ai.Message, 0, len(messages))
	for _, msg := range messages {
		result = append(result, ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
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

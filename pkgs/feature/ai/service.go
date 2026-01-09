package ai

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Service interface {
	Chat(conversationID, userMessage string, addToHistory bool) (string, error)
	ChatWithImages(conversationID, userMessage string, imageURLs []string, imageDetail string, addToHistory bool) (string, error)
	ClearHistory(conversationID string)
}

type AIService struct {
	provider         Provider
	history          History
	toolRegistry     ToolRegistry
	model            string
	systemPrompt     string
	systemPromptPath string
	promptModTime    time.Time
	promptMu         sync.RWMutex
	temperature      float64
	topP             float64
	maxTokens        int
}

type Config struct {
	APIURL           string
	APIKey           string
	Model            string
	SystemPromptPath string
	MaxHistory       int
	Temperature      float64
	TopP             float64
	MaxTokens        int
}

func NewService(cfg Config) Service {
	systemPrompt, modTime := loadSystemPromptWithTime(cfg.SystemPromptPath)

	return &AIService{
		provider:         NewHTTPProvider(cfg.APIURL, cfg.APIKey),
		history:          NewMemoryHistory(cfg.MaxHistory),
		toolRegistry:     NewDefaultToolRegistry(),
		model:            cfg.Model,
		systemPrompt:     systemPrompt,
		systemPromptPath: cfg.SystemPromptPath,
		promptModTime:    modTime,
		temperature:      cfg.Temperature,
		topP:             cfg.TopP,
		maxTokens:        cfg.MaxTokens,
	}
}

func loadSystemPromptWithTime(path string) (string, time.Time) {
	stat, err := os.Stat(path)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to stat system prompt %s: %v, using default", path, err))
		return "", time.Time{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to load system prompt from %s: %v, using default", path, err))
		return "", time.Time{}
	}
	return string(data), stat.ModTime()
}

func (s *AIService) getSystemPrompt() string {
	stat, err := os.Stat(s.systemPromptPath)
	if err == nil && stat.ModTime().After(s.promptModTime) {
		s.promptMu.Lock()
		if stat.ModTime().After(s.promptModTime) {
			newPrompt, modTime := loadSystemPromptWithTime(s.systemPromptPath)
			s.systemPrompt = newPrompt
			s.promptModTime = modTime
			logger.Info(fmt.Sprintf("System prompt auto-reloaded from %s", s.systemPromptPath))
		}
		s.promptMu.Unlock()
	}

	s.promptMu.RLock()
	defer s.promptMu.RUnlock()
	return s.systemPrompt
}

func (s *AIService) Chat(conversationID, userMessage string, addToHistory bool) (string, error) {
	userMsg := Message{Role: "user", Content: userMessage}

	if addToHistory {
		s.history.Add(conversationID, userMsg)
	}

	systemPrompt := s.getSystemPrompt()
	messages := []Message{{Role: "system", Content: systemPrompt}}
	messages = append(messages, s.history.Get(conversationID)...)

	req := ChatRequest{
		Model:       s.model,
		Messages:    messages,
		Stream:      false,
		Temperature: s.temperature,
		TopP:        s.topP,
		MaxTokens:   s.maxTokens,
		Tools:       s.toolRegistry.GetTools(),
	}

	resp, err := s.provider.Complete(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) > 0 {
		return s.handleToolCalls(conversationID, choice.Message, addToHistory)
	}

	if addToHistory {
		assistantMsg := Message{Role: "assistant", Content: choice.Message.Content}
		s.history.Add(conversationID, assistantMsg)
	}

	if contentStr, ok := choice.Message.Content.(string); ok {
		return contentStr, nil
	}
	return "", fmt.Errorf("unexpected content type in response")
}

func (s *AIService) ChatWithImages(conversationID, userMessage string, imageURLs []string, imageDetail string, addToHistory bool) (string, error) {
	var content interface{}
	if len(imageURLs) > 0 {
		contentItems := []ContentItem{}
		for _, url := range imageURLs {
			contentItems = append(contentItems, ContentItem{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL:    url,
					Detail: imageDetail,
				},
			})
		}
		if userMessage != "" {
			contentItems = append(contentItems, ContentItem{
				Type: "text",
				Text: userMessage,
			})
		}
		content = contentItems
	} else {
		content = userMessage
	}

	userMsg := Message{Role: "user", Content: content}

	if addToHistory {
		s.history.Add(conversationID, userMsg)
	}

	systemPrompt := s.getSystemPrompt()
	messages := []Message{{Role: "system", Content: systemPrompt}}
	messages = append(messages, s.history.Get(conversationID)...)

	req := ChatRequest{
		Model:       s.model,
		Messages:    messages,
		Stream:      false,
		Temperature: s.temperature,
		TopP:        s.topP,
		MaxTokens:   s.maxTokens,
	}

	resp, err := s.provider.Complete(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) > 0 {
		return s.handleToolCalls(conversationID, choice.Message, addToHistory)
	}

	if addToHistory {
		assistantMsg := Message{Role: "assistant", Content: choice.Message.Content}
		s.history.Add(conversationID, assistantMsg)
	}

	if contentStr, ok := choice.Message.Content.(string); ok {
		return contentStr, nil
	}
	return "", fmt.Errorf("unexpected content type in response")
}

func (s *AIService) handleToolCalls(conversationID string, assistantMsg Message, addToHistory bool) (string, error) {
	if addToHistory {
		s.history.Add(conversationID, assistantMsg)
	}

	for _, toolCall := range assistantMsg.ToolCalls {
		result := s.toolRegistry.Execute(toolCall.Function.Name, toolCall.Function.Arguments)

		if addToHistory {
			toolMsg := Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			}
			s.history.Add(conversationID, toolMsg)
		}
	}

	systemPrompt := s.getSystemPrompt()
	messages := []Message{{Role: "system", Content: systemPrompt}}
	messages = append(messages, s.history.Get(conversationID)...)

	req := ChatRequest{
		Model:       s.model,
		Messages:    messages,
		Stream:      false,
		Temperature: s.temperature,
		Tools:       s.toolRegistry.GetTools(),
	}

	resp, err := s.provider.Complete(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		if contentStr, ok := content.(string); ok {
			if addToHistory {
				s.history.Add(conversationID, Message{Role: "assistant", Content: contentStr})
			}
			return contentStr, nil
		}
		return "", fmt.Errorf("unexpected content type in response")
	}

	return "", fmt.Errorf("no response after tool call")
}

func (s *AIService) ClearHistory(conversationID string) {
	s.history.Clear(conversationID)
}

package ai

import (
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Service interface {
	Chat(conversationID, userMessage string, addToHistory bool) (*ChatResult, error)
	ChatWithImages(conversationID, userMessage string, imageURLs []string, imageDetail string, addToHistory bool) (*ChatResult, error)
	ClearHistory(conversationID string)
}

type ChatResult struct {
	Thinking string
	Content  string
}

type AIService struct {
	provider         Provider
	history          History
	toolRegistry     ToolRegistry
	textModel        string
	visionModel      string
	systemPrompt     string
	systemPromptPath string
	promptModTime    time.Time
	promptMu         sync.RWMutex
	temperature      float64
	topP             float64
	maxTokens        int
	toolsEnabled     bool
}

type Config struct {
	APIURL           string
	APIKey           string
	TextModel        string
	VisionModel      string
	SystemPromptPath string
	MaxHistory       int
	Temperature      float64
	TopP             float64
	MaxTokens        int
	ToolsEnabled     bool
}

func NewService(cfg Config) Service {
	systemPrompt, modTime := loadSystemPromptWithTime(cfg.SystemPromptPath)

	return &AIService{
		provider:         NewHTTPProvider(cfg.APIURL, cfg.APIKey),
		history:          NewMemoryHistory(cfg.MaxHistory),
		toolRegistry:     NewDefaultToolRegistry(),
		textModel:        cfg.TextModel,
		visionModel:      cfg.VisionModel,
		systemPrompt:     systemPrompt,
		systemPromptPath: cfg.SystemPromptPath,
		promptModTime:    modTime,
		temperature:      cfg.Temperature,
		topP:             cfg.TopP,
		maxTokens:        cfg.MaxTokens,
		toolsEnabled:     cfg.ToolsEnabled,
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

func parseThinkTags(content string) (string, string) {
	re := regexp.MustCompile(`(?s)<think>(.*?)</think>`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		thinkContent := matches[1]
		cleanContent := re.ReplaceAllString(content, "")
		cleanContent = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(cleanContent, "")
		return thinkContent, cleanContent
	}

	return "", content
}

func convertMessagesToText(messages []Message) []Message {
	result := make([]Message, len(messages))
	for i, msg := range messages {
		result[i] = Message{
			Role:             msg.Role,
			ReasoningContent: msg.ReasoningContent,
			ToolCalls:        msg.ToolCalls,
			ToolCallID:       msg.ToolCallID,
		}

		if items, ok := msg.Content.([]ContentItem); ok {
			var textParts []string
			for _, item := range items {
				if item.Type == "text" {
					textParts = append(textParts, item.Text)
				}
			}
			if len(textParts) > 0 {
				result[i].Content = textParts[0]
			} else {
				result[i].Content = ""
			}
		} else {
			result[i].Content = msg.Content
		}
	}
	return result
}

func (s *AIService) Chat(conversationID, userMessage string, addToHistory bool) (*ChatResult, error) {
	userMsg := Message{Role: "user", Content: userMessage}

	if addToHistory {
		s.history.Add(conversationID, userMsg)
	}

	systemPrompt := s.getSystemPrompt()
	messages := []Message{{Role: "system", Content: systemPrompt}}
	messages = append(messages, convertMessagesToText(s.history.Get(conversationID))...)

	req := ChatRequest{
		Model:       s.textModel,
		Messages:    messages,
		Stream:      false,
		Temperature: s.temperature,
		TopP:        s.topP,
		MaxTokens:   s.maxTokens,
	}

	if s.toolsEnabled {
		req.Tools = s.toolRegistry.GetTools()
	}

	resp, err := s.provider.Complete(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
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

func (s *AIService) ChatWithImages(conversationID, userMessage string, imageURLs []string, imageDetail string, addToHistory bool) (*ChatResult, error) {
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
		Model:       s.visionModel,
		Messages:    messages,
		Stream:      false,
		Temperature: s.temperature,
		TopP:        s.topP,
		MaxTokens:   s.maxTokens,
	}

	resp, err := s.provider.Complete(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
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

func (s *AIService) handleToolCalls(conversationID string, assistantMsg Message, addToHistory bool) (*ChatResult, error) {
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
	messages = append(messages, convertMessagesToText(s.history.Get(conversationID))...)

	req := ChatRequest{
		Model:       s.textModel,
		Messages:    messages,
		Stream:      false,
		Temperature: s.temperature,
	}

	if s.toolsEnabled {
		req.Tools = s.toolRegistry.GetTools()
	}

	resp, err := s.provider.Complete(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		if contentStr, ok := content.(string); ok {
			if addToHistory {
				s.history.Add(conversationID, Message{Role: "assistant", Content: contentStr})
			}

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

func (s *AIService) ClearHistory(conversationID string) {
	s.history.Clear(conversationID)
}

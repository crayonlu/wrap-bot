package ai

import (
	"fmt"
	"log"
	"os"
)

type Service interface {
	Chat(conversationID, userMessage string, addToHistory bool) (string, error)
	ClearHistory(conversationID string)
}

type AIService struct {
	provider     Provider
	history      History
	toolRegistry ToolRegistry
	model        string
	systemPrompt string
	temperature  float64
	topP         float64
	maxTokens    int
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
	systemPrompt := loadSystemPrompt(cfg.SystemPromptPath)

	return &AIService{
		provider:     NewHTTPProvider(cfg.APIURL, cfg.APIKey),
		history:      NewMemoryHistory(cfg.MaxHistory),
		toolRegistry: NewDefaultToolRegistry(),
		model:        cfg.Model,
		systemPrompt: systemPrompt,
		temperature:  cfg.Temperature,
		topP:         cfg.TopP,
		maxTokens:    cfg.MaxTokens,
	}
}

func loadSystemPrompt(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load system prompt from %s: %v, using default", path, err)
		return "你是一个友好、乐于助人的猫娘小管家"
	}
	return string(data)
}

func (s *AIService) Chat(conversationID, userMessage string, addToHistory bool) (string, error) {
	userMsg := Message{Role: "user", Content: userMessage}

	if addToHistory {
		s.history.Add(conversationID, userMsg)
	}

	messages := []Message{{Role: "system", Content: s.systemPrompt}}
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

	return choice.Message.Content, nil
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

	messages := []Message{{Role: "system", Content: s.systemPrompt}}
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
		if addToHistory {
			s.history.Add(conversationID, Message{Role: "assistant", Content: content})
		}
		return content, nil
	}

	return "", fmt.Errorf("no response after tool call")
}

func (s *AIService) ClearHistory(conversationID string) {
	s.history.Clear(conversationID)
}

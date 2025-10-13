package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
)

type AIClient struct {
	apiURL       string
	apiKey       string
	model        string
	client       *http.Client
	systemPrompt string
	history      *ConversationHistory
}

type ConversationHistory struct {
	mu           sync.RWMutex
	conversations map[string][]Message
	maxMessages  int
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function FunctionDef  `json:"function"`
}

type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewConversationHistory(maxMessages int) *ConversationHistory {
	return &ConversationHistory{
		conversations: make(map[string][]Message),
		maxMessages:   maxMessages,
	}
}

func (h *ConversationHistory) Add(conversationID string, msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.conversations[conversationID]; !exists {
		h.conversations[conversationID] = make([]Message, 0)
	}

	h.conversations[conversationID] = append(h.conversations[conversationID], msg)

	if len(h.conversations[conversationID]) > h.maxMessages {
		h.conversations[conversationID] = h.conversations[conversationID][len(h.conversations[conversationID])-h.maxMessages:]
	}
}

func (h *ConversationHistory) Get(conversationID string) []Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if msgs, exists := h.conversations[conversationID]; exists {
		result := make([]Message, len(msgs))
		copy(result, msgs)
		return result
	}
	return []Message{}
}

func (h *ConversationHistory) Clear(conversationID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conversations, conversationID)
}

func NewAIClient(apiURL, apiKey, model, systemPromptPath string) *AIClient {
	systemPrompt := loadSystemPrompt(systemPromptPath)
	
	return &AIClient{
		apiURL:       apiURL,
		apiKey:       apiKey,
		model:        model,
		client:       &http.Client{Timeout: 60 * time.Second},
		systemPrompt: systemPrompt,
		history:      NewConversationHistory(20),
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

func (ai *AIClient) Chat(conversationID, userMessage string) (string, error) {
	userMsg := Message{Role: "user", Content: userMessage}
	ai.history.Add(conversationID, userMsg)

	messages := []Message{{Role: "system", Content: ai.systemPrompt}}
	messages = append(messages, ai.history.Get(conversationID)...)

	tools := ai.getTools()

	reqBody := ChatRequest{
		Model:       ai.model,
		Messages:    messages,
		Stream:      false,
		Temperature: 0.7,
		TopP:        0.9,
		MaxTokens:   2000,
		Tools:       tools,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ai.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.apiKey)

	resp, err := ai.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	choice := chatResp.Choices[0]

	if len(choice.Message.ToolCalls) > 0 {
		return ai.handleToolCalls(conversationID, messages, choice.Message)
	}

	assistantMsg := Message{Role: "assistant", Content: choice.Message.Content}
	ai.history.Add(conversationID, assistantMsg)

	return choice.Message.Content, nil
}

func (ai *AIClient) handleToolCalls(conversationID string, messages []Message, assistantMsg Message) (string, error) {
	ai.history.Add(conversationID, assistantMsg)

	for _, toolCall := range assistantMsg.ToolCalls {
		result := ai.executeFunction(toolCall.Function.Name, toolCall.Function.Arguments)
		
		toolMsg := Message{
			Role:       "tool",
			Content:    result,
			ToolCallID: toolCall.ID,
		}
		ai.history.Add(conversationID, toolMsg)
	}

	messages = []Message{{Role: "system", Content: ai.systemPrompt}}
	messages = append(messages, ai.history.Get(conversationID)...)

	reqBody := ChatRequest{
		Model:       ai.model,
		Messages:    messages,
		Stream:      false,
		Temperature: 0.7,
		Tools:       ai.getTools(),
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", ai.apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.apiKey)

	resp, err := ai.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chatResp ChatResponse
	json.Unmarshal(body, &chatResp)

	if len(chatResp.Choices) > 0 {
		content := chatResp.Choices[0].Message.Content
		ai.history.Add(conversationID, Message{Role: "assistant", Content: content})
		return content, nil
	}

	return "", fmt.Errorf("no response after tool call")
}

func (ai *AIClient) executeFunction(name, argsJSON string) string {
	switch name {
	case "get_current_time":
		return time.Now().Format("2006-01-02 15:04:05")
	case "calculate":
		var args struct {
			Expression string `json:"expression"`
		}
		json.Unmarshal([]byte(argsJSON), &args)
		return fmt.Sprintf("计算结果: %s", args.Expression)
	default:
		return fmt.Sprintf("Unknown function: %s", name)
	}
}

func (ai *AIClient) getTools() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "get_current_time",
				Description: "获取当前时间",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "calculate",
				Description: "进行数学计算",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"expression": map[string]interface{}{
							"type":        "string",
							"description": "要计算的数学表达式",
						},
					},
					"required": []string{"expression"},
				},
			},
		},
	}
}

func (ai *AIClient) ClearHistory(conversationID string) {
	ai.history.Clear(conversationID)
}

func AIChatPlugin(cfg *config.Config) bot.HandlerFunc {
	if !cfg.AIEnabled {
		return func(ctx *bot.Context) {}
	}

	aiClient := NewAIClient(
		cfg.AIURL,
		cfg.AIKey,
		cfg.AIModel,
		cfg.SystemPromptPath,
	)

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
			aiClient.ClearHistory(conversationID)
			ctx.ReplyText("空空如也了")
			return
		}

		response, err := aiClient.Chat(conversationID, text)
		if err != nil {
			log.Printf("AI chat error: %v", err)
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

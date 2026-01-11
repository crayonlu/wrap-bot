package service

import (
	"context"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
)

type ChatService interface {
	Chat(conversationID, message string) (*agent.ChatResult, error)
	ChatWithImages(conversationID, message string, imageURLs []string) (*agent.ChatResult, error)
	ClearHistory(conversationID string) error
}

type ChatServiceImpl struct {
	agent *agent.ChatAgent
}

func NewChatService(agent *agent.ChatAgent) *ChatServiceImpl {
	return &ChatServiceImpl{
		agent: agent,
	}
}

func (s *ChatServiceImpl) Chat(conversationID, message string) (*agent.ChatResult, error) {
	return s.agent.Chat(context.Background(), conversationID, message)
}

func (s *ChatServiceImpl) ChatWithImages(conversationID, message string, imageURLs []string) (*agent.ChatResult, error) {
	return s.agent.ChatWithImages(context.Background(), conversationID, message, imageURLs)
}

func (s *ChatServiceImpl) ClearHistory(conversationID string) error {
	return s.agent.ClearHistory(conversationID)
}

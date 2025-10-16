package ai

import "sync"

type History interface {
	Add(conversationID string, msg Message)
	Get(conversationID string) []Message
	Clear(conversationID string)
}

type MemoryHistory struct {
	mu            sync.RWMutex
	conversations map[string][]Message
	maxMessages   int
}

func NewMemoryHistory(maxMessages int) *MemoryHistory {
	return &MemoryHistory{
		conversations: make(map[string][]Message),
		maxMessages:   maxMessages,
	}
}

func (h *MemoryHistory) Add(conversationID string, msg Message) {
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

func (h *MemoryHistory) Get(conversationID string) []Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if msgs, exists := h.conversations[conversationID]; exists {
		result := make([]Message, len(msgs))
		copy(result, msgs)
		return result
	}
	return []Message{}
}

func (h *MemoryHistory) Clear(conversationID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conversations, conversationID)
}

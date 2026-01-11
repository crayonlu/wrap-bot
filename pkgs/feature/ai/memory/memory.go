package memory

import "sync"

type Message struct {
	Role    string
	Content interface{}
}

type ConversationStore interface {
	Add(conversationID string, msg Message) error
	Get(conversationID string) ([]Message, error)
	Clear(conversationID string) error
}

type MemoryStore struct {
	mu            sync.RWMutex
	conversations map[string][]Message
	maxMessages   int
}

func NewMemoryStore(maxMessages int) *MemoryStore {
	return &MemoryStore{
		conversations: make(map[string][]Message),
		maxMessages:   maxMessages,
	}
}

func (s *MemoryStore) Add(conversationID string, msg Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.conversations[conversationID]; !exists {
		s.conversations[conversationID] = make([]Message, 0)
	}

	s.conversations[conversationID] = append(s.conversations[conversationID], msg)

	if len(s.conversations[conversationID]) > s.maxMessages {
		s.conversations[conversationID] = s.conversations[conversationID][len(s.conversations[conversationID])-s.maxMessages:]
	}

	return nil
}

func (s *MemoryStore) Get(conversationID string) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if msgs, exists := s.conversations[conversationID]; exists {
		result := make([]Message, len(msgs))
		copy(result, msgs)
		return result, nil
	}
	return []Message{}, nil
}

func (s *MemoryStore) Clear(conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.conversations, conversationID)
	return nil
}

type HistoryManager struct {
	store ConversationStore
}

func NewHistoryManager(store ConversationStore) *HistoryManager {
	return &HistoryManager{
		store: store,
	}
}

func (m *HistoryManager) AddMessage(conversationID string, msg Message) error {
	return m.store.Add(conversationID, msg)
}

func (m *HistoryManager) GetHistory(conversationID string) ([]Message, error) {
	return m.store.Get(conversationID)
}

func (m *HistoryManager) ClearHistory(conversationID string) error {
	return m.store.Clear(conversationID)
}

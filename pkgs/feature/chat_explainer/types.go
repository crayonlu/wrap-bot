package chat_explainer

import (
	"time"
)

type ChatMessage struct {
	MessageID   int64
	SenderName  string
	SenderID    int64
	Content     string
	Images      []string
	ReplyTo     *ChatMessage
	Timestamp   int64
	MessageType string
}

type ForwardedChat struct {
	SourceGroup string
	Messages    []ChatMessage
}

type MessageAnalysis struct {
	Content string
}

type ChatAnalysis struct {
	MessageAnalyses []MessageAnalysis
	Summary         string
	OriginalCount   int
	MergedCount     int
}

type MessageProcessingStatus struct {
	MessageIndex int
	MessageID    int64
	Processed    bool
	Analysis     string
	Error        error
	Timestamp    time.Time
}

type ChainRequestConfig struct {
	ConversationID    string
	Messages          []ChatMessage
	TimeoutPerRequest time.Duration
}

type ChainRequestResult struct {
	MessageAnalyses []MessageAnalysis
	Summary         string
	TotalTokens     int
	SuccessCount    int
	FailedCount     int
	Errors          []error
}

type MergedMessage struct {
	SenderName  string
	SenderID    int64
	Contents    []string
	Images      []string
	MessageType string
	MessageIDs  []int64
	Timestamps  []int64
}

type MessageGroup struct {
	MergedMessages []MergedMessage
	OriginalCount  int
	MergedCount    int
}

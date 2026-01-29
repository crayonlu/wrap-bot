package chat_explainer

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
}

package types

type PluginStatus struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type TaskStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Schedule    string `json:"schedule"`
	NextRun     string `json:"next_run"`
	LastRun     string `json:"last_run"`
	Status      string `json:"status"`
	CanTrigger  bool   `json:"can_trigger"`
	Description string `json:"description,omitempty"`
}

type BotStatus struct {
	Running   bool   `json:"running"`
	Uptime    int64  `json:"uptime"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
}

type ConfigUpdate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ConfigItem struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type AITool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type AIChatRequest struct {
	Message        string   `json:"message"`
	Images         []string `json:"images,omitempty"`
	Model          string   `json:"model"`
	ConversationID string   `json:"conversation_id,omitempty"`
}

type AIChatResponse struct {
	Response       string     `json:"response"`
	ToolCalls      []ToolCall `json:"tool_calls,omitempty"`
	ConversationID string     `json:"conversation_id"`
}

type ToolCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

package ai

import "fmt"

type Analyzer struct {
	service Service
}

func NewAnalyzer(service Service) *Analyzer {
	return &Analyzer{service: service}
}

func (a *Analyzer) Analyze(content string) (string, error) {
	prompt := fmt.Sprintf(`请分析以下今日技术热点：

%s`, content)

	conversationID := "tech_analysis"
	a.service.ClearHistory(conversationID)

	return a.service.Chat(conversationID, prompt, true)
}

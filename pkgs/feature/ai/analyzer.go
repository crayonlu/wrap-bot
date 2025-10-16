package ai

import "fmt"

type Analyzer struct {
	service Service
}

func NewAnalyzer(service Service) *Analyzer {
	return &Analyzer{service: service}
}

func (a *Analyzer) Analyze(content string) (string, error) {
	prompt := fmt.Sprintf(`总结今日热点趋势：

%s

请分析主要技术趋势和值得关注的热点。`, content)
	return a.service.Chat("tech_analysis", prompt, false)
}

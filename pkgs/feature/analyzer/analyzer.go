package analyzer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Analyzer struct {
	agent                 *agent.ChatAgent
	analyzerPrompt        string
	analyzerPromptPath    string
	analyzerPromptModTime time.Time
	promptMu              sync.RWMutex
}

type AnalyzerConfig struct {
	Agent              *agent.ChatAgent
	AnalyzerPromptPath string
}

func NewAnalyzer(cfg AnalyzerConfig) *Analyzer {
	analyzerPrompt, modTime := loadAnalyzerPromptWithTime(cfg.AnalyzerPromptPath)
	return &Analyzer{
		agent:                 cfg.Agent,
		analyzerPrompt:        analyzerPrompt,
		analyzerPromptPath:    cfg.AnalyzerPromptPath,
		analyzerPromptModTime: modTime,
	}
}

func loadAnalyzerPromptWithTime(path string) (string, time.Time) {
	stat, err := os.Stat(path)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to stat analyzer prompt %s: %v, using default", path, err))
		return "请分析以下今日技术热点：", time.Time{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to load analyzer prompt from %s: %v, using default", path, err))
		return "请分析以下今日技术热点：", time.Time{}
	}
	return string(data), stat.ModTime()
}

func (a *Analyzer) getAnalyzerPrompt() string {
	stat, err := os.Stat(a.analyzerPromptPath)
	if err == nil && stat.ModTime().After(a.analyzerPromptModTime) {
		a.promptMu.Lock()
		if stat.ModTime().After(a.analyzerPromptModTime) {
			newPrompt, modTime := loadAnalyzerPromptWithTime(a.analyzerPromptPath)
			a.analyzerPrompt = newPrompt
			a.analyzerPromptModTime = modTime
			logger.Info(fmt.Sprintf("Analyzer prompt auto-reloaded from %s", a.analyzerPromptPath))
		}
		a.promptMu.Unlock()
	}

	a.promptMu.RLock()
	defer a.promptMu.RUnlock()
	return a.analyzerPrompt
}

func (a *Analyzer) Analyze(content string) (string, error) {
	analyzerPrompt := a.getAnalyzerPrompt()
	prompt := fmt.Sprintf(`%s

%s`, analyzerPrompt, content)

	conversationID := "tech_analysis"
	a.agent.ClearHistory(conversationID)

	result, err := a.agent.Chat(context.Background(), conversationID, prompt)
	if err != nil {
		return "", err
	}
	return result.Content, nil
}

package chat_explainer

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/memory"
	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Analyzer struct {
	agent        *agent.ChatAgent
	memoryStore  *memory.MemoryStore
	systemPrompt string
	logger       *logger.Logger
}

func NewAnalyzer(chatAgent *agent.ChatAgent, systemPrompt string, maxHistory int) *Analyzer {
	return &Analyzer{
		agent:        chatAgent,
		memoryStore:  memory.NewMemoryStore(maxHistory),
		systemPrompt: systemPrompt,
		logger:       logger.NewLogger(1000),
	}
}

func (a *Analyzer) AnalyzeAll(ctx context.Context, messages []ChatMessage) (*ChatAnalysis, error) {
	if len(messages) == 0 {
		return &ChatAnalysis{}, nil
	}

	chainResult, err := a.AnalyzeWithChain(ctx, messages)
	if err != nil {
		return nil, err
	}

	return &ChatAnalysis{
		MessageAnalyses: chainResult.MessageAnalyses,
		Summary:         chainResult.Summary,
	}, nil
}

func (a *Analyzer) AnalyzeWithChain(ctx context.Context, messages []ChatMessage) (*ChainRequestResult, error) {
	conversationID := fmt.Sprintf("forward_%d_%d", messages[0].MessageID, time.Now().Unix())

	config := &ChainRequestConfig{
		ConversationID:    conversationID,
		Messages:          messages,
		TimeoutPerRequest: 30 * time.Second,
	}

	return a.chainAnalyze(ctx, config)
}

func (a *Analyzer) chainAnalyze(ctx context.Context, config *ChainRequestConfig) (*ChainRequestResult, error) {
	result := &ChainRequestResult{
		MessageAnalyses: make([]MessageAnalysis, 0, len(config.Messages)),
	}

	a.logger.Info(fmt.Sprintf("开始链式分析消息，共 %d 条", len(config.Messages)))

	for i, msg := range config.Messages {
		a.logger.Info(fmt.Sprintf("处理第 %d/%d 条消息，发送者: %s", i+1, len(config.Messages), msg.SenderName))

		analysis, err := a.processSingleMessage(ctx, config.ConversationID, msg, i, result.MessageAnalyses)
		if err != nil {
			a.logger.Error(fmt.Sprintf("消息 %d 处理失败: %v", i+1, err))
			result.Errors = append(result.Errors, fmt.Errorf("消息%d处理失败: %w", i+1, err))
			result.FailedCount++

			result.MessageAnalyses = append(result.MessageAnalyses, MessageAnalysis{
				Content: fmt.Sprintf("[分析失败] %v", err),
			})
			continue
		}

		result.MessageAnalyses = append(result.MessageAnalyses, MessageAnalysis{
			Content: analysis,
		})
		result.SuccessCount++

		a.logger.Info(fmt.Sprintf("消息 %d 处理完成", i+1))
	}

	a.logger.Info("开始生成整体总结")
	summary, err := a.generateSummary(ctx, config.ConversationID, result.MessageAnalyses)
	if err != nil {
		a.logger.Error(fmt.Sprintf("总结生成失败: %v", err))
		result.Errors = append(result.Errors, fmt.Errorf("总结生成失败: %w", err))
		result.Summary = "总结生成失败"
	} else {
		result.Summary = summary
		a.logger.Info("总结生成成功")
	}

	a.logger.Info(fmt.Sprintf("链式分析完成，成功 %d/%d", result.SuccessCount, len(config.Messages)))

	return result, nil
}

func (a *Analyzer) processSingleMessage(ctx context.Context, conversationID string, msg ChatMessage, index int, previousAnalyses []MessageAnalysis) (string, error) {
	prompt := a.buildSingleMessagePrompt(msg, previousAnalyses)

	a.logger.Info(fmt.Sprintf("[Prompt] 消息%d，发送者: %s, 内容长度: %d", index+1, msg.SenderName, len(prompt)))
	a.logger.Debug(fmt.Sprintf("[Prompt] 完整内容: %s", prompt))

	var result *agent.ChatResult
	var err error

	if len(msg.Images) > 0 {
		result, err = a.agent.ChatWithImages(ctx, conversationID, prompt, msg.Images)
	} else {
		result, err = a.agent.Chat(ctx, conversationID, prompt)
	}

	if err != nil {
		return "", err
	}

	analysis := parseAnalysisResult(result.Content)
	return analysis, nil
}

func (a *Analyzer) buildSingleMessagePrompt(msg ChatMessage, previousAnalyses []MessageAnalysis) string {
	var prompt strings.Builder

	prompt.WriteString("你是群聊对话解读助手。\n\n")

	if len(previousAnalyses) > 0 {
		prompt.WriteString("【对话上下文】以下是之前消息的分析结果，帮助你理解对话脉络：\n")
		for i, analysis := range previousAnalyses {
			prompt.WriteString(fmt.Sprintf("消息%d：%s\n", i+1, analysis.Content))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("【当前消息】请分析以下这条群聊消息：\n")
	prompt.WriteString(fmt.Sprintf("发送者：%s (ID: %d)\n", msg.SenderName, msg.SenderID))

	if msg.Content != "" {
		prompt.WriteString(fmt.Sprintf("消息内容：%s\n", msg.Content))
	}

	if len(msg.Images) > 0 {
		prompt.WriteString(fmt.Sprintf("图片数量：%d\n", len(msg.Images)))
	}

	if msg.ReplyTo != nil {
		prompt.WriteString(fmt.Sprintf("回复给：%s 的消息\n", msg.ReplyTo.SenderName))
	}

	prompt.WriteString("\n请结合上面的对话上下文，用2-3句话解释这条消息在讨论什么，在对话中起什么作用。")

	return prompt.String()
}

func (a *Analyzer) generateSummary(ctx context.Context, conversationID string, analyses []MessageAnalysis) (string, error) {
	if len(analyses) == 0 {
		return "", nil
	}

	var prompt strings.Builder

	prompt.WriteString("基于以下群聊消息的逐条分析，请总结整段对话：\n\n")

	for i, analysis := range analyses {
		prompt.WriteString(fmt.Sprintf("消息%d：%s\n", i+1, analysis.Content))
	}

	prompt.WriteString("\n请总结：\n")
	prompt.WriteString("1. 这段对话的整体脉络（发生了什么、讨论了什么话题）\n")
	prompt.WriteString("2. 出现的专业术语或关键概念\n")

	result, err := a.agent.Chat(ctx, conversationID, prompt.String())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Content), nil
}

func parseAnalysisResult(content string) string {
	re := regexp.MustCompile(`(?s)<think>\n(.*?)\n</think>\n(.*)`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 2 {
		return strings.TrimSpace(matches[2])
	}

	re2 := regexp.MustCompile(`(?s)(.*?)\n</think>\n(.*)`)
	matches2 := re2.FindStringSubmatch(content)
	if len(matches2) > 2 {
		return strings.TrimSpace(matches2[2])
	}

	return strings.TrimSpace(content)
}

func (a *Analyzer) AnalyzeMessage(ctx context.Context, msg ChatMessage, prevMessages []ChatMessage) (*MessageAnalysis, error) {
	var contextBuilder strings.Builder
	if len(prevMessages) > 0 {
		contextBuilder.WriteString("前文：\n")
		start := 0
		if len(prevMessages) > 3 {
			start = len(prevMessages) - 3
		}
		for _, prev := range prevMessages[start:] {
			content := prev.Content
			if content == "" {
				if len(prev.Images) > 0 {
					content = "[图片]"
				} else {
					content = "[表情]"
				}
			}
			contextBuilder.WriteString(fmt.Sprintf("%s: %s\n", prev.SenderName, content))
		}
		contextBuilder.WriteString("\n")
	}

	displayContent := msg.Content
	if displayContent == "" {
		if len(msg.Images) > 0 {
			displayContent = "[图片消息，请直接分析图片内容]"
		} else {
			displayContent = "[表情或无文字内容]"
		}
	}

	prompt := fmt.Sprintf(`%s需要解读的群聊消息：

发送者：%s
消息内容：%s

结合上下文理解这条消息在讨论什么，如果涉及专业术语请顺便解释一下。保持简洁，2-3句话即可。`,
		contextBuilder.String(), msg.SenderName, displayContent)

	result, err := a.agent.ChatWithImages(ctx, fmt.Sprintf("explainer_%d", msg.MessageID), prompt, msg.Images)
	if err != nil {
		return nil, err
	}

	return &MessageAnalysis{
		Content: result.Content,
	}, nil
}

func (a *Analyzer) AnalyzeBatch(ctx context.Context, messages []ChatMessage, individualAnalyses []MessageAnalysis) (*ChatAnalysis, error) {
	var conversation strings.Builder
	for i, msg := range messages {
		content := msg.Content
		if content == "" {
			if len(msg.Images) > 0 {
				content = "[图片]"
			} else {
				content = "[表情]"
			}
		}
		conversation.WriteString(fmt.Sprintf("%s: %s\n", msg.SenderName, content))
		if i < len(individualAnalyses) {
			conversation.WriteString(fmt.Sprintf("[解读] %s\n", individualAnalyses[i].Content))
		}
	}

	prompt := fmt.Sprintf(`需要总结的群聊对话：

%s

总结要求：
1. 对话的整体脉络（发生了什么、讨论了什么话题、结论是什么）
2. 涉及的专业术语解释（按合适的方式组织）

**重要规则：**
- 必须基于上述实际提供的对话内容进行总结
- 严禁编造或假设不存在的内容
- 如果对话内容较少或简单，直接说明即可，不要过度解读

用自然语言自由发挥，让总结清晰易懂。`, conversation.String())

	result, err := a.agent.Chat(ctx, "explainer_summary", prompt)
	if err != nil {
		return nil, err
	}

	return &ChatAnalysis{
		MessageAnalyses: individualAnalyses,
		Summary:         result.Content,
	}, nil
}

package chat_explainer

import (
	"context"
	"fmt"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
)

type Analyzer struct {
	agent *agent.ChatAgent
}

func NewAnalyzer(chatAgent *agent.ChatAgent) *Analyzer {
	return &Analyzer{
		agent: chatAgent,
	}
}

func (a *Analyzer) AnalyzeAll(ctx context.Context, messages []ChatMessage) (*ChatAnalysis, error) {
	analyses := make([]MessageAnalysis, len(messages))

	for i, msg := range messages {
		prevMessages := messages[:i]

		analysis, err := a.AnalyzeMessage(ctx, msg, prevMessages)
		if err != nil {
			return nil, err
		}
		analyses[i] = *analysis
	}

	return a.AnalyzeBatch(ctx, messages, analyses)
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
			contextBuilder.WriteString(fmt.Sprintf("%s: %s\n", prev.SenderName, prev.Content))
		}
		contextBuilder.WriteString("\n")
	}

	prompt := fmt.Sprintf(`%s请解读这条群聊消息：

发送者：%s
消息内容：%s

请结合上下文理解这条消息在讨论什么，如果涉及专业术语请顺便解释一下。保持简洁，2-3句话即可。`,
		contextBuilder.String(), msg.SenderName, msg.Content)

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
		conversation.WriteString(fmt.Sprintf("%s: %s\n", msg.SenderName, msg.Content))
		if i < len(individualAnalyses) {
			conversation.WriteString(fmt.Sprintf("[解读] %s\n", individualAnalyses[i].Content))
		}
	}

	prompt := fmt.Sprintf(`请总结这段群聊对话：

%s

请提供：
1. 对话的整体脉络（发生了什么、讨论了什么话题、结论是什么）
2. 涉及的专业术语解释（按你觉得合适的方式组织）

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

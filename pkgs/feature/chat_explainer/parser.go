package chat_explainer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/napcat"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseForwardMessage(data map[string]interface{}) (*ForwardedChat, error) {
	messages, ok := data["messages"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid messages format")
	}

	var chat ForwardedChat
	if groupName, ok := data["group_name"].(string); ok {
		chat.SourceGroup = groupName
	}

	for _, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			continue
		}

		chatMsg := p.parseSingleMessage(msgMap)
		chat.Messages = append(chat.Messages, *chatMsg)
	}

	return &chat, nil
}

func (p *Parser) parseSingleMessage(data map[string]interface{}) *ChatMessage {
	msg := &ChatMessage{}

	if sender, ok := data["sender"].(map[string]interface{}); ok {
		if nickname, ok := sender["nickname"].(string); ok {
			msg.SenderName = nickname
		}
		if userID, ok := sender["user_id"].(float64); ok {
			msg.SenderID = int64(userID)
		}
	}

	if time, ok := data["time"].(float64); ok {
		msg.Timestamp = int64(time)
	}

	if messageID, ok := data["message_id"].(float64); ok {
		msg.MessageID = int64(messageID)
	}

	content := p.extractContent(data)
	msg.Content = content.Text
	msg.Images = content.Images

	if len(msg.Images) > 0 && msg.Content != "" {
		msg.MessageType = "mixed"
	} else if len(msg.Images) > 0 {
		msg.MessageType = "image"
	} else {
		msg.MessageType = "text"
	}

	return msg
}

type extractedContent struct {
	Text   string
	Images []string
}

func (p *Parser) extractContent(data map[string]interface{}) extractedContent {
	var result extractedContent

	if rawMsg, ok := data["raw_message"].(string); ok && rawMsg != "" {
		result.Text = p.cleanCQCode(rawMsg)
		result.Images = p.extractImageURLs(rawMsg)
		return result
	}

	if message, ok := data["message"].([]interface{}); ok {
		for _, seg := range message {
			segMap, ok := seg.(map[string]interface{})
			if !ok {
				continue
			}

			msgType, _ := segMap["type"].(string)
			data, _ := segMap["data"].(map[string]interface{})

			switch msgType {
			case "text":
				if text, ok := data["text"].(string); ok {
					result.Text += text
				}
			case "image":
				if url, ok := data["url"].(string); ok {
					result.Images = append(result.Images, url)
				} else if file, ok := data["file"].(string); ok {
					result.Images = append(result.Images, file)
				}
			}
		}
	}

	return result
}

func (p *Parser) cleanCQCode(text string) string {
	re := regexp.MustCompile(`\[CQ:[^\]]+\]`)
	return re.ReplaceAllString(text, "")
}

func (p *Parser) extractImageURLs(text string) []string {
	var urls []string
	re := regexp.MustCompile(`\[CQ:image,[^\]]*url=([^,\]]+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.ReplaceAll(match[1], "&amp;", "&")
			urls = append(urls, url)
		}
	}
	return urls
}

func (p *Parser) IsForwardMessage(event map[string]interface{}) bool {
	if message, ok := event["message"].([]interface{}); ok && len(message) > 0 {
		if firstSeg, ok := message[0].(map[string]interface{}); ok {
			if msgType, ok := firstSeg["type"].(string); ok {
				return msgType == "forward"
			}
		}
	}
	return false
}

func (p *Parser) GetForwardID(event map[string]interface{}) string {
	if message, ok := event["message"].([]interface{}); ok && len(message) > 0 {
		if firstSeg, ok := message[0].(map[string]interface{}); ok {
			if data, ok := firstSeg["data"].(map[string]interface{}); ok {
				if id, ok := data["id"].(string); ok {
					return id
				}
			}
		}
	}
	return ""
}

func BuildOriginalContent(msg ChatMessage) string {
	var content strings.Builder

	if msg.Content != "" {
		content.WriteString(msg.Content)
	}

	if len(msg.Images) > 0 {
		if content.Len() > 0 {
			content.WriteString("\n")
		}
		content.WriteString(fmt.Sprintf("[图片 x%d]", len(msg.Images)))
	}

	return content.String()
}

func BuildForwardNodes(messages []ChatMessage, analyses []MessageAnalysis, summary string, selfID int64) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

	for i, msg := range messages {
		var content strings.Builder

		original := BuildOriginalContent(msg)
		content.WriteString(original)

		if i < len(analyses) {
			content.WriteString("\n-------\n")
			content.WriteString(analyses[i].Content)
		}

		node := napcat.NewMixedForwardNode(
			msg.SenderName,
			msg.SenderID,
			napcat.NewTextSegment(content.String()),
		)
		nodes = append(nodes, node)
	}

	if summary != "" {
		summaryNode := napcat.NewMixedForwardNode(
			"整体总结",
			selfID,
			napcat.NewTextSegment(summary),
		)
		nodes = append(nodes, summaryNode)
	}

	return nodes
}

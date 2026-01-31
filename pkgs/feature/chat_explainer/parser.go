package chat_explainer

import (
	"fmt"
	"regexp"
	"strconv"
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

	parsedMessages := make([]ChatMessage, 0, len(messages))

	for _, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			continue
		}

		chatMsg := p.parseSingleMessage(msgMap)
		parsedMessages = append(parsedMessages, *chatMsg)
	}

	for i := range parsedMessages {
		if replyID, ok := p.getReplyIDFromRaw(i, messages); ok {
			parsedMessages[i].ReplyTo = p.FindMessageByID(parsedMessages, replyID)
		}
	}

	chat.Messages = parsedMessages

	return &chat, nil
}

func (p *Parser) getReplyIDFromRaw(index int, messages []interface{}) (int64, bool) {
	if index >= len(messages) {
		return 0, false
	}

	msgMap, ok := messages[index].(map[string]interface{})
	if !ok {
		return 0, false
	}

	return p.GetReplyID(msgMap)
}

func (p *Parser) parseSingleMessage(data map[string]interface{}) *ChatMessage {
	msg := &ChatMessage{}

	if sender, ok := data["sender"].(map[string]interface{}); ok {
		if nickname, ok := sender["nickname"].(string); ok {
			msg.SenderName = nickname
		}
		if userID, ok := sender["user_id"].(float64); ok {
			msg.SenderID = int64(userID)
		} else if userIDStr, ok := sender["user_id"].(string); ok {
			if parsedID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
				msg.SenderID = parsedID
			}
		}
	} else {
		if nickname, ok := data["sender_nickname"].(string); ok {
			msg.SenderName = nickname
		}
		if userID, ok := data["sender_user_id"].(float64); ok {
			msg.SenderID = int64(userID)
		} else if userIDStr, ok := data["sender_user_id"].(string); ok {
			if parsedID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
				msg.SenderID = parsedID
			}
		}
	}

	if msg.SenderName == "" {
		if nickname, ok := data["nickname"].(string); ok {
			msg.SenderName = nickname
		}
	}

	if msg.SenderID == 0 {
		if userID, ok := data["user_id"].(float64); ok {
			msg.SenderID = int64(userID)
		} else if userIDStr, ok := data["user_id"].(string); ok {
			if parsedID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
				msg.SenderID = parsedID
			}
		}
	}

	if time, ok := data["time"].(float64); ok {
		msg.Timestamp = int64(time)
	}

	if messageID, ok := data["message_id"].(float64); ok {
		msg.MessageID = int64(messageID)
	} else if msgIDStr, ok := data["message_id"].(string); ok {
		if parsedID, err := strconv.ParseInt(msgIDStr, 10, 64); err == nil {
			msg.MessageID = parsedID
		}
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
		var segments []napcat.MessageSegment

		if msg.Content != "" {
			segments = append(segments, napcat.NewTextSegment(msg.Content))
		}

		for _, imageURL := range msg.Images {
			segments = append(segments, napcat.NewImageSegment(imageURL))
		}

		if i < len(analyses) && analyses[i].Content != "" {
			segments = append(segments, napcat.NewTextSegment("\n-------\n"+analyses[i].Content))
		}

		if len(segments) == 0 {
			segments = append(segments, napcat.NewTextSegment("[无内容]"))
		}

		node := napcat.NewMixedForwardNode(
			msg.SenderName,
			msg.SenderID,
			segments...,
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

func (p *Parser) GetReplyID(data map[string]interface{}) (int64, bool) {
	if message, ok := data["message"].([]interface{}); ok {
		for _, seg := range message {
			segMap, ok := seg.(map[string]interface{})
			if !ok {
				continue
			}
			if msgType, _ := segMap["type"].(string); msgType == "reply" {
				if data, ok := segMap["data"].(map[string]interface{}); ok {
					if id, ok := data["id"].(string); ok {
						replyID, err := strconv.ParseInt(id, 10, 64)
						if err == nil {
							return replyID, true
						}
					}
				}
			}
		}
	}
	return 0, false
}

func (p *Parser) FindMessageByID(messages []ChatMessage, id int64) *ChatMessage {
	for i := range messages {
		if messages[i].MessageID == id {
			return &messages[i]
		}
	}
	return nil
}

func (p *Parser) IsMergeable(msg ChatMessage) bool {
	return msg.MessageType == "text" && msg.ReplyTo == nil && len(msg.Images) == 0
}

func (p *Parser) MergeConsecutiveMessages(messages []ChatMessage) MessageGroup {
	if len(messages) == 0 {
		return MessageGroup{}
	}

	var merged []MergedMessage
	var current *MergedMessage

	for _, msg := range messages {
		if p.IsMergeable(msg) {
			if current == nil {
				current = &MergedMessage{
					SenderName:  msg.SenderName,
					SenderID:    msg.SenderID,
					MessageType: msg.MessageType,
					MessageIDs:  []int64{msg.MessageID},
					Timestamps:  []int64{msg.Timestamp},
					Contents:    []string{msg.Content},
				}
			} else if current.SenderID == msg.SenderID {
				current.Contents = append(current.Contents, msg.Content)
				current.MessageIDs = append(current.MessageIDs, msg.MessageID)
				current.Timestamps = append(current.Timestamps, msg.Timestamp)
			} else {
				if current != nil {
					merged = append(merged, *current)
				}
				current = &MergedMessage{
					SenderName:  msg.SenderName,
					SenderID:    msg.SenderID,
					MessageType: msg.MessageType,
					MessageIDs:  []int64{msg.MessageID},
					Timestamps:  []int64{msg.Timestamp},
					Contents:    []string{msg.Content},
				}
			}
		} else {
			if current != nil {
				merged = append(merged, *current)
				current = nil
			}
			merged = append(merged, MergedMessage{
				SenderName:  msg.SenderName,
				SenderID:    msg.SenderID,
				Contents:    []string{msg.Content},
				Images:      msg.Images,
				MessageType: msg.MessageType,
				MessageIDs:  []int64{msg.MessageID},
				Timestamps:  []int64{msg.Timestamp},
			})
		}
	}

	if current != nil {
		merged = append(merged, *current)
	}

	return MessageGroup{
		MergedMessages: merged,
		OriginalCount:  len(messages),
		MergedCount:    len(merged),
	}
}

func BuildForwardNodesWithMerge(messages []ChatMessage, analyses []MessageAnalysis, summary string, selfID int64, mergedInfo *MessageGroup) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

	if mergedInfo != nil {
		for _, merged := range mergedInfo.MergedMessages {
			var segments []napcat.MessageSegment

			content := strings.Join(merged.Contents, "\n")
			if content != "" {
				segments = append(segments, napcat.NewTextSegment(content))
			}

			for _, imageURL := range merged.Images {
				segments = append(segments, napcat.NewImageSegment(imageURL))
			}

			analysisContent := merged.AnalysisContent
			if len(merged.Contents) > 1 {
				analysisContent = fmt.Sprintf("【连续消息 x%d】%s", len(merged.Contents), analysisContent)
			}

			if analysisContent != "" {
				segments = append(segments, napcat.NewTextSegment("\n-------\n"+analysisContent))
			}

			if len(segments) == 0 {
				segments = append(segments, napcat.NewTextSegment("[无内容]"))
			}

			node := napcat.NewMixedForwardNode(
				merged.SenderName,
				merged.SenderID,
				segments...,
			)
			nodes = append(nodes, node)
		}
	} else {
		for i, msg := range messages {
			var segments []napcat.MessageSegment

			if msg.Content != "" {
				segments = append(segments, napcat.NewTextSegment(msg.Content))
			}

			for _, imageURL := range msg.Images {
				segments = append(segments, napcat.NewImageSegment(imageURL))
			}

			if i < len(analyses) && analyses[i].Content != "" {
				segments = append(segments, napcat.NewTextSegment("\n-------\n"+analyses[i].Content))
			}

			if len(segments) == 0 {
				segments = append(segments, napcat.NewTextSegment("[无内容]"))
			}

			node := napcat.NewMixedForwardNode(
				msg.SenderName,
				msg.SenderID,
				segments...,
			)
			nodes = append(nodes, node)
		}
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

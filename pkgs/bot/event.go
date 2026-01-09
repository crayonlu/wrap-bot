package bot

import (
	"regexp"
	"strings"
)

type EventType string

const (
	EventTypeMessage       EventType = "message"
	EventTypeNotice        EventType = "notice"
	EventTypeRequest       EventType = "request"
	EventTypeMetaEvent     EventType = "meta_event"
	EventTypeMessageSent   EventType = "message_sent"
)

type MessageType string

const (
	MessageTypePrivate MessageType = "private"
	MessageTypeGroup   MessageType = "group"
)

type NoticeType string

const (
	NoticeTypeGroupUpload     NoticeType = "group_upload"
	NoticeTypeGroupAdmin      NoticeType = "group_admin"
	NoticeTypeGroupDecrease   NoticeType = "group_decrease"
	NoticeTypeGroupIncrease   NoticeType = "group_increase"
	NoticeTypeGroupBan        NoticeType = "group_ban"
	NoticeTypeFriendAdd       NoticeType = "friend_add"
	NoticeTypeGroupRecall     NoticeType = "group_recall"
	NoticeTypeFriendRecall    NoticeType = "friend_recall"
	NoticeTypeNotify          NoticeType = "notify"
)

type RequestType string

const (
	RequestTypeFriend RequestType = "friend"
	RequestTypeGroup  RequestType = "group"
)

type Event struct {
	Time        int64                  `json:"time"`
	SelfID      int64                  `json:"self_id"`
	PostType    EventType              `json:"post_type"`
	MessageType MessageType            `json:"message_type,omitempty"`
	SubType     string                 `json:"sub_type,omitempty"`
	MessageID   int32                  `json:"message_id,omitempty"`
	UserID      int64                  `json:"user_id,omitempty"`
	Message     interface{}            `json:"message,omitempty"`
	RawMessage  string                 `json:"raw_message,omitempty"`
	Font        int32                  `json:"font,omitempty"`
	Sender      *Sender                `json:"sender,omitempty"`
	GroupID     int64                  `json:"group_id,omitempty"`
	NoticeType  NoticeType             `json:"notice_type,omitempty"`
	RequestType RequestType            `json:"request_type,omitempty"`
	Comment     string                 `json:"comment,omitempty"`
	Flag        string                 `json:"flag,omitempty"`
	MetaType    string                 `json:"meta_event_type,omitempty"`
	Extra       map[string]interface{} `json:"-"`
}

type Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Age      int32  `json:"age,omitempty"`
	Area     string `json:"area,omitempty"`
	Level    string `json:"level,omitempty"`
	Role     string `json:"role,omitempty"`
	Title    string `json:"title,omitempty"`
}

type MessageSegment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (e *Event) IsGroupMessage() bool {
	return e.PostType == EventTypeMessage && e.MessageType == MessageTypeGroup
}

func (e *Event) IsPrivateMessage() bool {
	return e.PostType == EventTypeMessage && e.MessageType == MessageTypePrivate
}

func (e *Event) IsNotice() bool {
	return e.PostType == EventTypeNotice
}

func (e *Event) IsRequest() bool {
	return e.PostType == EventTypeRequest
}

func (e *Event) GetText() string {
	if e.RawMessage != "" {
		return e.RawMessage
	}
	return ""
}

func (e *Event) GetImages() []string {
	if e.RawMessage == "" {
		return nil
	}

	re := regexp.MustCompile(`\[CQ:image,[^\]]*url=([^,\]]+)`)
	matches := re.FindAllStringSubmatch(e.RawMessage, -1)

	var urls []string
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.ReplaceAll(match[1], "&amp;", "&")
			urls = append(urls, url)
		}
	}

	return urls
}

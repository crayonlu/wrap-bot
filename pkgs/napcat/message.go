package napcat

import (
	"encoding/json"
	"github.com/crayon/bot_golang/pkgs/bot"
)

type ForwardNode struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type MessageSegment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (c *Client) SendGroupMessage(groupID int64, message interface{}) (int32, error) {
	payload := map[string]interface{}{
		"group_id": groupID,
		"message":  message,
	}
	
	resp, err := c.post("/send_group_msg", payload)
	if err != nil {
		return 0, err
	}
	
	var result struct {
		MessageID int32 `json:"message_id"`
	}
	
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return 0, err
	}
	
	return result.MessageID, nil
}

func (c *Client) SendPrivateMessage(userID int64, message interface{}) (int32, error) {
	payload := map[string]interface{}{
		"user_id": userID,
		"message": message,
	}
	
	resp, err := c.post("/send_private_msg", payload)
	if err != nil {
		return 0, err
	}
	
	var result struct {
		MessageID int32 `json:"message_id"`
	}
	
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return 0, err
	}
	
	return result.MessageID, nil
}

func (c *Client) DeleteMessage(messageID int32) error {
	payload := map[string]interface{}{
		"message_id": messageID,
	}
	
	_, err := c.post("/delete_msg", payload)
	return err
}

func (c *Client) GetGroupList() ([]bot.Group, error) {
	resp, err := c.get("/get_group_list")
	if err != nil {
		return nil, err
	}
	
	var groups []bot.Group
	if err := json.Unmarshal(resp.Data, &groups); err != nil {
		return nil, err
	}
	
	return groups, nil
}

func (c *Client) GetGroupInfo(groupID int64) (*bot.GroupInfo, error) {
	payload := map[string]interface{}{
		"group_id": groupID,
	}
	
	resp, err := c.post("/get_group_info", payload)
	if err != nil {
		return nil, err
	}
	
	var info bot.GroupInfo
	if err := json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	
	return &info, nil
}

func (c *Client) GetGroupMemberList(groupID int64) ([]bot.GroupMember, error) {
	payload := map[string]interface{}{
		"group_id": groupID,
	}
	
	resp, err := c.post("/get_group_member_list", payload)
	if err != nil {
		return nil, err
	}
	
	var members []bot.GroupMember
	if err := json.Unmarshal(resp.Data, &members); err != nil {
		return nil, err
	}
	
	return members, nil
}

func (c *Client) GetFriendList() ([]bot.Friend, error) {
	resp, err := c.get("/get_friend_list")
	if err != nil {
		return nil, err
	}
	
	var friends []bot.Friend
	if err := json.Unmarshal(resp.Data, &friends); err != nil {
		return nil, err
	}
	
	return friends, nil
}

func NewTextSegment(text string) MessageSegment {
	return MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	}
}

func NewImageSegment(file string) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

func NewAtSegment(qq int64) MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": qq,
		},
	}
}

func NewAtAllSegment() MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": "all",
		},
	}
}

func NewFaceSegment(id int) MessageSegment {
	return MessageSegment{
		Type: "face",
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

func NewVideoSegment(file string) MessageSegment {
	return MessageSegment{
		Type: "video",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

func NewRecordSegment(file string) MessageSegment {
	return MessageSegment{
		Type: "record",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

func NewCustomForwardNode(name string, uin int64, content interface{}) ForwardNode {
	return ForwardNode{
		Type: "node",
		Data: map[string]interface{}{
			"name":    name,
			"uin":     uin,
			"content": content,
		},
	}
}

func NewMixedForwardNode(name string, uin int64, segments ...MessageSegment) ForwardNode {
	content := make([]interface{}, len(segments))
	for i, seg := range segments {
		content[i] = seg
	}
	
	return ForwardNode{
		Type: "node",
		Data: map[string]interface{}{
			"name":    name,
			"uin":     uin,
			"content": content,
		},
	}
}

func NewMessageForwardNode(messageID int32) ForwardNode {
	return ForwardNode{
		Type: "node",
		Data: map[string]interface{}{
			"id": messageID,
		},
	}
}

func (c *Client) SendGroupForwardMsg(groupID int64, nodes []ForwardNode) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"group_id": groupID,
		"messages": nodes,
	}
	
	resp, err := c.post("/send_group_forward_msg", payload)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (c *Client) SendPrivateForwardMsg(userID int64, nodes []ForwardNode) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"user_id":  userID,
		"messages": nodes,
	}
	
	resp, err := c.post("/send_private_forward_msg", payload)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (c *Client) GetForwardMsg(messageID string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"message_id": messageID,
	}
	
	resp, err := c.post("/get_forward_msg", payload)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (c *Client) ForwardMsgToGroup(groupID int64, messageID int32) error {
	payload := map[string]interface{}{
		"group_id":   groupID,
		"message_id": messageID,
	}
	
	_, err := c.post("/forward_group_single_msg", payload)
	return err
}

func (c *Client) ForwardMsgToPrivate(userID int64, messageID int32) error {
	payload := map[string]interface{}{
		"user_id":    userID,
		"message_id": messageID,
	}
	
	_, err := c.post("/forward_friend_single_msg", payload)
	return err
}

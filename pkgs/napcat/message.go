package napcat

import (
	"encoding/json"
	"github.com/crayon/bot_golang/pkgs/bot"
)

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

package handlers

import (
	"encoding/json"
	"fmt"
)

type BilibiliVideo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Cover     string `json:"cover"`
	Author    string `json:"author"`
	Timestamp int64  `json:"timestamp"`
	Hot       int    `json:"hot"`
	URL       string `json:"url"`
	MobileUrl string `json:"mobile_url"`
}

type BilibiliRes struct {
	Name     string          `json:"name"`
	Videos  []BilibiliVideo  `json:"videos"`
}

func BilibiliHandler(data []byte) (*BilibiliRes, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse bilibili data: %w", err)
	}
	name := result["title"].(string)
	videos := result["data"].([]BilibiliVideo)

	return &BilibiliRes{
		Name:    name,
		Videos: videos,
	}, nil
}

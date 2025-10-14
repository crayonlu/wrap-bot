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
	Title  string          `json:"title"`
	Videos []BilibiliVideo `json:"data"`
}

func BilibiliHandler(data []byte) (*BilibiliRes, error) {
	var res BilibiliRes
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to parse bilibili data: %w", err)
	}

	return &res, nil
}

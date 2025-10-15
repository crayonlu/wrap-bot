package handlers

import (
	"encoding/json"
	"fmt"
)

type JuejinArticle struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Hot    int    `json:"hot"`
	URL    string `json:"url"`
}

type JuejinRes struct {
	Title    string          `json:"title"`
	Articles []JuejinArticle `json:"data"`
}

func JuejinHandler(data []byte) (*JuejinRes, error) {
	var res JuejinRes
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to parse juejin data: %w", err)
	}

	return &res, nil
}

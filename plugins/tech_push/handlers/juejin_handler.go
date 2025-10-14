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
	Name     string          `json:"name"`
	Articles []JuejinArticle `json:"articles"`
}

func JuejinHandler(data []byte) (*JuejinRes, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse juejin data: %w", err)
	}
	name := result["title"].(string)
	articles := result["data"].([]JuejinArticle)

	return &JuejinRes{
		Name:     name,
		Articles: articles,
	}, nil
}

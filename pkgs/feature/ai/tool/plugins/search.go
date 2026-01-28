package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/tool"
)

func RegisterSearchTools(registry tool.ToolRegistry, client *SerpAPIClient) error {
	tools := []tool.Tool{
		{
			Name:        "web_search",
			Description: "网络搜索",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "搜索关键词",
					},
					"num_results": map[string]interface{}{
						"type":        "integer",
						"description": "结果数量",
					},
				},
				"required": []string{"query"},
			},
			Handler: WebSearch(client),
			Enabled: true,
		},
	}

	for _, t := range tools {
		if err := registry.Register(t); err != nil {
			return err
		}
	}
	return nil
}

func WebSearch(client *SerpAPIClient) tool.ToolHandler {
	return func(ctx context.Context, args string) (string, error) {
		var params struct {
			Query      string `json:"query"`
			NumResults int    `json:"num_results"`
		}

		json.Unmarshal([]byte(args), &params)

		if params.NumResults == 0 {
			params.NumResults = 5
		}

		result, err := client.Search(ctx, params.Query, params.NumResults)
		if err != nil {
			return "", err
		}

		var builder strings.Builder
		for i, item := range result.Results {
			if i > 0 {
				builder.WriteString("\n\n")
			}
			builder.WriteString(fmt.Sprintf("%d. %s\n%s\n%s", i+1, item.Title, item.URL, item.Snippet))
		}

		return builder.String(), nil
	}
}

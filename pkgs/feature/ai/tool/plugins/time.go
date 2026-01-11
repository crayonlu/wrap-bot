package plugins

import (
	"context"
	"encoding/json"
	"time"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/tool"
)

func RegisterTimeTools(registry tool.ToolRegistry) error {
	tools := []tool.Tool{
		{
			Name:        "get_current_time",
			Description: "获取当前时间",
			Category:   tool.CategoryBoth,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"timezone": map[string]interface{}{
						"type":        "string",
						"description": "时区，如 Asia/Shanghai",
					},
				},
			},
			Handler: GetCurrentTime,
			Enabled: true,
		},
		{
			Name:        "parse_relative_time",
			Description: "解析相对时间，如'3天后'",
			Category:   tool.CategoryBoth,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "相对时间表达式",
					},
				},
				"required": []string{"expression"},
			},
			Handler: ParseRelativeTime,
			Enabled: true,
		},
		{
			Name:        "format_time",
			Description: "格式化时间显示",
			Category:   tool.CategoryBoth,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"type":        "integer",
						"description": "Unix 时间戳",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"description": "时间格式，默认为 2006-01-02 15:04:05",
					},
				},
			},
			Handler: FormatTime,
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

func GetCurrentTime(ctx context.Context, args string) (string, error) {
	var params struct {
		Timezone string `json:"timezone"`
	}

	json.Unmarshal([]byte(args), &params)

	now := time.Now()

	if params.Timezone != "" {
		loc, err := time.LoadLocation(params.Timezone)
		if err != nil {
			now = now.UTC()
		} else {
			now = now.In(loc)
		}
	}

	return now.Format("2006-01-02 15:04:05"), nil
}

func ParseRelativeTime(ctx context.Context, args string) (string, error) {
	var params struct {
		Expression string `json:"expression"`
	}

	json.Unmarshal([]byte(args), &params)

	parser := NewTimeParser()
	result, err := parser.Parse(params.Expression)
	if err != nil {
		return "", err
	}

	return result.Format("2006-01-02 15:04:05"), nil
}

func FormatTime(ctx context.Context, args string) (string, error) {
	var params struct {
		Timestamp int64  `json:"timestamp"`
		Format    string `json:"format"`
	}

	json.Unmarshal([]byte(args), &params)

	if params.Timestamp == 0 {
		return time.Now().Format("2006-01-02 15:04:05"), nil
	}

	t := time.Unix(params.Timestamp, 0)

	format := params.Format
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	return t.Format(format), nil
}

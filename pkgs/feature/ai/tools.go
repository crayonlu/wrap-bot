package ai

import (
	"encoding/json"
	"fmt"
	"time"
)

type ToolRegistry interface {
	GetTools() []Tool
	Execute(name, argsJSON string) string
}

type DefaultToolRegistry struct {
	tools []Tool
}

func NewDefaultToolRegistry() *DefaultToolRegistry {
	return &DefaultToolRegistry{
		tools: []Tool{
			{
				Type: "function",
				Function: FunctionDef{
					Name:        "get_current_time",
					Description: "获取当前时间",
					Parameters: map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
			},
			{
				Type: "function",
				Function: FunctionDef{
					Name:        "calculate",
					Description: "进行数学计算",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"expression": map[string]interface{}{
								"type":        "string",
								"description": "要计算的数学表达式",
							},
						},
						"required": []string{"expression"},
					},
				},
			},
		},
	}
}

func (r *DefaultToolRegistry) GetTools() []Tool {
	return r.tools
}

func (r *DefaultToolRegistry) Execute(name, argsJSON string) string {
	switch name {
	case "get_current_time":
		return time.Now().Format("2006-01-02 15:04:05")
	case "calculate":
		var args struct {
			Expression string `json:"expression"`
		}
		json.Unmarshal([]byte(argsJSON), &args)
		return fmt.Sprintf("计算结果: %s", args.Expression)
	default:
		return fmt.Sprintf("Unknown function: %s", name)
	}
}

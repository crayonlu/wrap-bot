package tool

import (
	"context"
	"fmt"
	"sync"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Handler     ToolHandler
	Enabled     bool
}

type ToolHandler func(ctx context.Context, args string) (string, error)

type ToolRegistry interface {
	Register(tool Tool) error
	Unregister(name string) error
	Get(name string) (Tool, bool)
	GetAll() []Tool
	GetEnabled() []Tool
	Execute(ctx context.Context, name string, args string) (string, error)
}

type DefaultToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewToolRegistry() *DefaultToolRegistry {
	return &DefaultToolRegistry{
		tools: make(map[string]Tool),
	}
}

func (r *DefaultToolRegistry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name] = tool
	return nil
}

func (r *DefaultToolRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tools, name)
	return nil
}

func (r *DefaultToolRegistry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, exists := r.tools[name]
	return tool, exists
}

func (r *DefaultToolRegistry) GetAll() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool)
	}
	return result
}

func (r *DefaultToolRegistry) GetEnabled() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Tool
	for _, tool := range r.tools {
		if tool.Enabled {
			result = append(result, tool)
		}
	}
	return result
}

func (r *DefaultToolRegistry) Execute(ctx context.Context, name string, args string) (string, error) {
	logger.Info(fmt.Sprintf("[ToolRegistry] Executing tool: %s", name))

	tool, exists := r.Get(name)
	if !exists {
		logger.Error(fmt.Sprintf("[ToolRegistry] Tool not found: %s", name))
		return "", fmt.Errorf("tool not found")
	}

	if !tool.Enabled {
		logger.Warn(fmt.Sprintf("[ToolRegistry] Tool is disabled: %s", name))
		return "", fmt.Errorf("tool is disabled")
	}

	logger.Info(fmt.Sprintf("[ToolRegistry] Tool %s is enabled, executing with args: %s", name, args))
	result, err := tool.Handler(ctx, args)
	if err != nil {
		logger.Error(fmt.Sprintf("[ToolRegistry] Tool %s execution error: %v", name, err))
	} else {
		logger.Info(fmt.Sprintf("[ToolRegistry] Tool %s executed successfully", name))
	}

	return result, err
}

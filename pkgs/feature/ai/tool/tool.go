package tool

import (
	"context"
	"fmt"
	"sync"
)

type ToolCategory string

const (
	CategoryText   ToolCategory = "text"
	CategoryVision ToolCategory = "vision"
	CategoryBoth   ToolCategory = "both"
)

type Tool struct {
	Name        string
	Description string
	Category    ToolCategory
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
	GetByCategory(category ToolCategory) []Tool
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

func (r *DefaultToolRegistry) GetByCategory(category ToolCategory) []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Tool
	for _, tool := range r.tools {
		if tool.Category == category || tool.Category == CategoryBoth {
			result = append(result, tool)
		}
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
	tool, exists := r.Get(name)
	if !exists {
		return "", fmt.Errorf("tool not found")
	}

	if !tool.Enabled {
		return "", fmt.Errorf("tool is disabled")
	}

	return tool.Handler(ctx, args)
}

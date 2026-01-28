package factory

import (
	"os"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/memory"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/provider"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/service"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/tool"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/tool/plugins"
)

type Factory struct {
	config *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		config: cfg,
	}
}

func (f *Factory) CreateProvider() provider.LLMProvider {
	return provider.NewHTTPProvider(f.config.APIURL, f.config.APIKey)
}

func (f *Factory) CreateMemoryStore() *memory.MemoryStore {
	return memory.NewMemoryStore(f.config.MaxHistory)
}

func (f *Factory) CreateToolRegistry() tool.ToolRegistry {
	registry := tool.NewToolRegistry()

	plugins.RegisterTimeTools(registry)

	if f.config.SerpAPIKey != "" {
		searchClient := plugins.NewSerpAPIClient(f.config.SerpAPIKey)
		plugins.RegisterSearchTools(registry, searchClient)
	}

	if f.config.WeatherAPIKey != "" {
		weatherClient := plugins.NewWeatherAPIClient(f.config.WeatherAPIKey)
		plugins.RegisterWeatherTools(registry, weatherClient)
	}

	return registry
}

func (f *Factory) CreateAgent() *agent.ChatAgent {
	return agent.NewChatAgent(agent.AgentConfig{
		Provider:     f.CreateProvider(),
		History:      memory.NewHistoryManager(f.CreateMemoryStore()),
		ToolRegistry: f.CreateToolRegistry(),
		SystemPrompt: f.loadSystemPrompt(),

		Model: f.config.Model,

		Temperature: f.config.Temperature,
		TopP:        f.config.TopP,
		MaxTokens:   f.config.MaxTokens,

		ToolsEnabled: f.config.ToolsEnabled,
	})
}

func (f *Factory) CreateChatService() service.ChatService {
	return service.NewChatService(f.CreateAgent())
}

func (f *Factory) loadSystemPrompt() string {
	data, err := os.ReadFile(f.config.SystemPromptPath)
	if err != nil {
		return ""
	}
	return string(data)
}

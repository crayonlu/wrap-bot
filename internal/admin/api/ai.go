package api

import (
	"net/http"
	"time"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/agent"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/config"
	"github.com/crayon/wrap-bot/pkgs/feature/ai/factory"
	"github.com/labstack/echo/v4"
)

var aiFactory *factory.Factory
var aiAgent *agent.ChatAgent

func initAI() {
	if aiFactory == nil {
		cfg := config.Load()
		aiFactory = factory.NewFactory(cfg)
		aiAgent = aiFactory.CreateAgent()
	}
}

func GetAITools(c echo.Context) error {
	initAI()

	cfg := config.Load()
	enabledTools := cfg.ToolsEnabled

	tools := []types.AITool{
		{
			Name:        "get_current_time",
			Description: "获取当前时间",
			Enabled:     contains(enabledTools, "get_current_time"),
		},
		{
			Name:        "parse_relative_time",
			Description: "解析相对时间表达式（如'3天后'）",
			Enabled:     contains(enabledTools, "parse_relative_time"),
		},
		{
			Name:        "web_search",
			Description: "网络搜索",
			Enabled:     contains(enabledTools, "web_search"),
		},
		{
			Name:        "get_weather",
			Description: "获取当前天气",
			Enabled:     contains(enabledTools, "get_weather"),
		},
		{
			Name:        "get_weather_forecast",
			Description: "获取天气预报",
			Enabled:     contains(enabledTools, "get_weather_forecast"),
		},
	}

	return c.JSON(http.StatusOK, tools)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TestAIChat(c echo.Context) error {
	initAI()

	req := new(types.AIChatRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	conversationID := req.ConversationID
	if conversationID == "" {
		conversationID = "test-" + time.Now().Format("20060102-150405")
	}

	result, err := aiAgent.Chat(c.Request().Context(), conversationID, req.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := types.AIChatResponse{
		Response:       result.Content,
		ConversationID: conversationID,
	}

	return c.JSON(http.StatusOK, response)
}

func TestAIImageChat(c echo.Context) error {
	initAI()

	req := new(types.AIChatRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if len(req.Images) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no images provided"})
	}

	conversationID := req.ConversationID
	if conversationID == "" {
		conversationID = "test-" + time.Now().Format("20060102-150405")
	}

	result, err := aiAgent.ChatWithImages(c.Request().Context(), conversationID, req.Message, req.Images)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := types.AIChatResponse{
		Response:       result.Content,
		ConversationID: conversationID,
	}

	return c.JSON(http.StatusOK, response)
}

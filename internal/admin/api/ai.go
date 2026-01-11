package api

import (
	"net/http"
	"os"
	"strings"
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

	tools := []types.AITool{
		{
			Name:        "get_current_time",
			Description: "获取当前时间",
			Category:    "both",
		},
		{
			Name:        "parse_relative_time",
			Description: "解析相对时间表达式（如'3天后'）",
			Category:    "both",
		},
		{
			Name:        "web_search",
			Description: "网络搜索",
			Category:    "both",
		},
		{
			Name:        "get_weather",
			Description: "获取当前天气",
			Category:    "both",
		},
		{
			Name:        "get_weather_forecast",
			Description: "获取天气预报",
			Category:    "both",
		},
	}

	textTools := strings.Split(os.Getenv("AI_TEXT_MODEL_TOOLS"), ",")

	for i := range tools {
		tools[i].Enabled = false
		for _, t := range textTools {
			if strings.TrimSpace(t) == tools[i].Name {
				tools[i].Enabled = true
				break
			}
		}
	}

	return c.JSON(http.StatusOK, tools)
}

func GetAIStats(c echo.Context) error {
	stats := types.AIStats{
		TotalCalls:  0,
		ToolUsage:   make(map[string]int),
		SuccessRate: 100.0,
		RecentCalls: []types.AICallRecord{},
	}

	return c.JSON(http.StatusOK, stats)
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

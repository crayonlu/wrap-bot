package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

var configDescriptions = map[string]string{
	"NAPCAT_HTTP_URL":       "NapCat HTTP API address",
	"NAPCAT_WS_URL":         "NapCat WebSocket address",
	"NAPCAT_HTTP_TOKEN":     "NapCat HTTP authentication token",
	"NAPCAT_WS_TOKEN":       "NapCat WebSocket authentication token",
	"SERVER_PORT":           "Admin backend port",
	"SERVER_ENABLED":        "Whether admin backend is enabled",
	"DEBUG":                 "DEBUG mode",
	"COMMAND_PREFIX":        "Command prefix",
	"AI_ENABLED":            "Whether AI features are enabled",
	"AI_URL":                "AI API address",
	"AI_KEY":                "AI API key",
	"AI_USE_UNIFIED":        "Whether to use unified model mode",
	"AI_UNIFIED_MODEL":      "Unified model name",
	"AI_TEXT_MODEL":         "Text model name",
	"AI_VISION_MODEL":       "Vision model name",
	"AI_TEMPERATURE":        "AI temperature parameter",
	"AI_TOP_P":              "AI Top-P parameter",
	"AI_MAX_TOKENS":         "AI max tokens",
	"AI_MAX_HISTORY":        "AI max history records",
	"AI_IMAGE_DETAIL":       "Image processing detail (high/low/auto)",
	"AI_TEXT_MODEL_TOOLS":   "Enabled tools for text model (comma-separated)",
	"AI_VISION_MODEL_TOOLS": "Enabled tools for vision model (comma-separated)",
	"SYSTEM_PROMPT_PATH":    "System prompt path",
	"ANALYZER_PROMPT_PATH":  "Analyzer prompt path",
	"HOT_API_HOST":          "Hot API URL",
	"HOT_API_KEY":           "Hot API key",
	"RSS_API_HOST":          "RSS API URL",
	"TECH_PUSH_GROUPS":      "Tech push group IDs (comma-separated)",
	"TECH_PUSH_USERS":       "Tech push user IDs (comma-separated)",
	"RSS_PUSH_GROUPS":       "RSS push group IDs (comma-separated)",
	"RSS_PUSH_USERS":        "RSS push user IDs (comma-separated)",
	"ALLOWED_USERS":         "Allowed user IDs (comma-separated)",
	"ALLOWED_GROUPS":        "Allowed group IDs (comma-separated)",
	"ADMIN_IDS":             "Admin user IDs (comma-separated)",
	"SERP_API_KEY":          "SerpAPI key (web search)",
	"WEATHER_API_KEY":       "WeatherAPI key (weather query)",
}

func GetConfig(c echo.Context) error {
	configKeys := []string{
		"NAPCAT_HTTP_URL",
		"NAPCAT_WS_URL",
		"NAPCAT_HTTP_TOKEN",
		"NAPCAT_WS_TOKEN",
		"SERVER_PORT",
		"SERVER_ENABLED",
		"DEBUG",
		"COMMAND_PREFIX",
		"AI_ENABLED",
		"AI_URL",
		"AI_KEY",
		"AI_USE_UNIFIED",
		"AI_UNIFIED_MODEL",
		"AI_TEXT_MODEL",
		"AI_VISION_MODEL",
		"AI_TEMPERATURE",
		"AI_TOP_P",
		"AI_MAX_TOKENS",
		"AI_MAX_HISTORY",
		"AI_IMAGE_DETAIL",
		"AI_TEXT_MODEL_TOOLS",
		"AI_VISION_MODEL_TOOLS",
		"SYSTEM_PROMPT_PATH",
		"ANALYZER_PROMPT_PATH",
		"HOT_API_HOST",
		"HOT_API_KEY",
		"RSS_API_HOST",
		"TECH_PUSH_GROUPS",
		"TECH_PUSH_USERS",
		"RSS_PUSH_GROUPS",
		"RSS_PUSH_USERS",
		"ALLOWED_USERS",
		"ALLOWED_GROUPS",
		"ADMIN_IDS",
		"SERP_API_KEY",
		"WEATHER_API_KEY",
	}

	config := make([]types.ConfigItem, 0, len(configKeys))
	for _, key := range configKeys {
		config = append(config, types.ConfigItem{
			Key:         key,
			Value:       os.Getenv(key),
			Description: configDescriptions[key],
		})
	}

	return c.JSON(http.StatusOK, config)
}

func UpdateConfig(c echo.Context) error {
	req := new([]types.ConfigItem)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	envFile := os.Getenv("APP_ENV_FILE")

	if envFile == "" {
		envFile = ".env"
	}

	content, err := os.ReadFile(envFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read .env"})
	}

	lines := strings.Split(string(content), "\n")
	updatedKeys := make(map[string]bool)

	for _, item := range *req {
		found := false
		for i, line := range lines {
			if strings.HasPrefix(line, item.Key+"=") {
				lines[i] = item.Key + "=" + item.Value
				found = true
				break
			}
		}

		if !found {
			lines = append(lines, item.Key+"="+item.Value)
		}

		os.Setenv(item.Key, item.Value)
		updatedKeys[item.Key] = true
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(envFile, []byte(newContent), 0644); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to write .env"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":        "updated",
		"updated_count": len(*req),
		"updated_keys":  updatedKeys,
	})
}

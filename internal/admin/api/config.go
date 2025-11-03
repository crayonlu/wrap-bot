package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

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
		"AI_MODEL",
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
	}

	config := make([]types.ConfigItem, 0, len(configKeys))
	for _, key := range configKeys {
		config = append(config, types.ConfigItem{
			Key:   key,
			Value: os.Getenv(key),
		})
	}

	return c.JSON(http.StatusOK, config)
}

func UpdateConfig(c echo.Context) error {
	req := new(types.ConfigUpdate)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	envFile := ".env"
	content, err := os.ReadFile(envFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read .env"})
	}

	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, req.Key+"=") {
			lines[i] = req.Key + "=" + req.Value
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, req.Key+"="+req.Value)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(envFile, []byte(newContent), 0644); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to write .env"})
	}

	os.Setenv(req.Key, req.Value)
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

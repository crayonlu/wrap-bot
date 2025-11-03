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
		"AI_ENABLED",
		"AI_URL",
		"AI_MODEL",
		"HOT_API_HOST",
		"RSS_API_HOST",
		"NAPCAT_HTTP_URL",
		"NAPCAT_WS_URL",
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

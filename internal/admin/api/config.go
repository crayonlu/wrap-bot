package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

func GetConfig(c echo.Context) error {
	config := map[string]string{
		"AI_ENABLED":      os.Getenv("AI_ENABLED"),
		"AI_URL":          os.Getenv("AI_URL"),
		"AI_MODEL":        os.Getenv("AI_MODEL"),
		"HOT_API_HOST":    os.Getenv("HOT_API_HOST"),
		"RSS_API_HOST":    os.Getenv("RSS_API_HOST"),
		"NAPCAT_HTTP_URL": os.Getenv("NAPCAT_HTTP_URL"),
		"NAPCAT_WS_URL":   os.Getenv("NAPCAT_WS_URL"),
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

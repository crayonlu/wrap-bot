package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

var configDescriptions = map[string]string{
	"NAPCAT_HTTP_URL":      "NapCat HTTP API 地址",
	"NAPCAT_WS_URL":        "NapCat WebSocket 地址",
	"NAPCAT_HTTP_TOKEN":    "NapCat HTTP 认证令牌",
	"NAPCAT_WS_TOKEN":      "NapCat WebSocket 认证令牌",
	"SERVER_PORT":          "管理后台端口",
	"SERVER_ENABLED":       "是否启用管理后台",
	"DEBUG":                "DEBUG模式",
	"COMMAND_PREFIX":       "命令前缀",
	"AI_ENABLED":           "是否启用 AI 功能",
	"AI_URL":               "AI API 地址",
	"AI_KEY":               "AI API 密钥",
	"AI_MODEL":             "AI 模型名称",
	"SYSTEM_PROMPT_PATH":   "系统提示词路径",
	"ANALYZER_PROMPT_PATH": "Analyzer 提示词路径",
	"HOT_API_HOST":         "热点 API URL",
	"HOT_API_KEY":          "热点 API 密钥",
	"RSS_API_HOST":         "RSS API URL",
	"TECH_PUSH_GROUPS":     "技术推送群号列表（逗号分隔）",
	"TECH_PUSH_USERS":      "技术推送 QQ 号列表（逗号分隔）",
	"RSS_PUSH_GROUPS":      "RSS 推送群号列表（逗号分隔）",
	"RSS_PUSH_USERS":       "RSS 推送 QQ 号列表（逗号分隔）",
	"ALLOWED_USERS":        "允许的 QQ 号列表（逗号分隔）",
	"ALLOWED_GROUPS":       "允许的群号列表（逗号分隔）",
	"ADMIN_IDS":            "管理员 QQ 号列表（逗号分隔）",
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
			Key:         key,
			Value:       os.Getenv(key),
			Description: configDescriptions[key],
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

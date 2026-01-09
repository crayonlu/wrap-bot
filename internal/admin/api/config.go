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
	"AI_VISION_ENABLED":    "是否启用 AI 图片分析功能",
	"AI_IMAGE_DETAIL":      "图片处理精度 (high/low/auto)",
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
		"AI_VISION_ENABLED",
		"AI_IMAGE_DETAIL",
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
	req := new([]types.ConfigItem)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	envFile := ".env"
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

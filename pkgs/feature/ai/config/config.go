package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Config struct {
	APIURL string
	APIKey string

	AIUnifiedModel string
	AIUseUnified   bool

	TextModel   string
	VisionModel string

	Temperature float64
	TopP        float64
	MaxTokens   int

	TextToolsEnabled    []string
	VisionToolsEnabled  []string

	MaxHistory       int
	SystemPromptPath string
	SerpAPIKey       string
	WeatherAPIKey    string
}

func Load() *Config {
	return &Config{
		APIURL: getEnv("AI_URL", "https://api.siliconflow.cn/v1/chat/completions"),
		APIKey: getEnv("AI_KEY", "YOUR_API_KEY_HERE"),

		AIUnifiedModel: getEnv("AI_UNIFIED_MODEL", ""),
		AIUseUnified:   getEnvBool("AI_USE_UNIFIED", false),

		TextModel:   getEnv("AI_TEXT_MODEL", "deepseek/deepseek-r1-turbo"),
		VisionModel: getEnv("AI_VISION_MODEL", "qwen/qwen3-vl-235b-a22b-thinking"),

		Temperature:      getEnvFloat64("AI_TEMPERATURE", 0.7),
		TopP:             getEnvFloat64("AI_TOP_P", 0.9),
		MaxTokens:        getEnvInt("AI_MAX_TOKENS", 2000),
		TextToolsEnabled:    getEnvStringSlice("AI_TEXT_MODEL_TOOLS", []string{}),
		VisionToolsEnabled:  getEnvStringSlice("AI_VISION_MODEL_TOOLS", []string{}),
		MaxHistory:       getEnvInt("AI_MAX_HISTORY", 20),
		SystemPromptPath: getEnv("SYSTEM_PROMPT_PATH", "configs/system_prompt.md"),
		SerpAPIKey:       getEnv("SERP_API_KEY", ""),
		WeatherAPIKey:    getEnv("WEATHER_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
		logger.Warn("Invalid int value for " + key + ": " + value + ", using default")
	}
	return defaultValue
}

func getEnvFloat64(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		var result float64
		_, err := fmt.Sscanf(value, "%f", &result)
		if err == nil {
			return result
		}
		logger.Warn("Invalid float64 value for " + key + ": " + value + ", using default")
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err == nil {
			return b
		}
		logger.Warn("Invalid boolean value for " + key + ": " + value + ", using default")
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

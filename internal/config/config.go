package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NapCatHTTPURL    string
	NapCatWSURL      string
	NapCatHTTPToken  string
	NapCatWSToken    string
	ServerPort       string
	ServerEnabled    bool
	Debug            bool
	AdminIDs         []int64
	CommandPrefix    string
	AIEnabled        bool
	AIURL            string
	AIKey            string
	AIModel          string
	SystemPromptPath string
	HotApiHost       string
	HotApiKey        string
	TechPushGroups   []int64
	TechPushUsers    []int64
}

func Load() *Config {
	return &Config{
		NapCatHTTPURL:    getEnv("NAPCAT_HTTP_URL", "http://localhost:3000"),
		NapCatWSURL:      getEnv("NAPCAT_WS_URL", "ws://localhost:3001"),
		NapCatHTTPToken:  getEnv("NAPCAT_HTTP_TOKEN", ""),
		NapCatWSToken:    getEnv("NAPCAT_WS_TOKEN", ""),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		ServerEnabled:    getEnvBool("SERVER_ENABLED", true),
		Debug:            getEnvBool("DEBUG", false),
		AdminIDs:         getEnvInt64Slice("ADMIN_IDS", []int64{}),
		CommandPrefix:    getEnv("COMMAND_PREFIX", "/"),
		AIEnabled:        getEnvBool("AI_ENABLED", false),
		AIURL:            getEnv("AI_URL", "https://api.siliconflow.cn/v1/chat/completions"),
		AIKey:            getEnv("AI_KEY", "YOUR_API_KEY_HERE"),
		AIModel:          getEnv("AI_MODEL", "deepseek-ai/DeepSeek-V3.1"),
		SystemPromptPath: getEnv("SYSTEM_PROMPT_PATH", "configs/system_prompt.md"),
		HotApiHost:       getEnv("HOT_API_HOST", "https://hot-api.crayoncreator.top"),
		HotApiKey:        getEnv("HOT_API_KEY", "keykeykey"),
		TechPushGroups:   getEnvInt64Slice("TECH_PUSH_GROUPS", []int64{}),
		TechPushUsers:    getEnvInt64Slice("TECH_PUSH_USERS", []int64{}),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid boolean value for %s: %s, using default: %v", key, value, defaultValue)
			return defaultValue
		}
		return b
	}
	return defaultValue
}

func getEnvInt64Slice(key string, defaultValue []int64) []int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			log.Printf("Invalid int64 value in %s: %s, skipping", key, part)
			continue
		}

		result = append(result, num)
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Config struct {
	NapCatHTTPURL      string
	NapCatWSURL        string
	NapCatHTTPToken    string
	NapCatWSToken      string
	ServerPort         string
	ServerEnabled      bool
	Debug              bool
	AdminIDs           []int64
	CommandPrefix      string
	AIEnabled          bool
	AIURL              string
	AIKey              string
	AIModel            string
	AIToolsEnabled     bool
	AIVisionEnabled    bool
	AIImageDetail      string
	SystemPromptPath   string
	AnalyzerPromptPath string
	HotApiHost         string
	HotApiKey          string
	TechPushGroups     []int64
	TechPushUsers      []int64
	RssPushGroups      []int64
	RssPushUsers       []int64
	AllowedUsers       []int64
	AllowedGroups      []int64
	RSSApiHost         string
	AdminUsername      string
	AdminPassword      string
	JWTSecret          string
}

func Load() *Config {
	return &Config{
		NapCatHTTPURL:      getEnv("NAPCAT_HTTP_URL", "http://localhost:3000"),
		NapCatWSURL:        getEnv("NAPCAT_WS_URL", "ws://localhost:3001"),
		NapCatHTTPToken:    getEnv("NAPCAT_HTTP_TOKEN", ""),
		NapCatWSToken:      getEnv("NAPCAT_WS_TOKEN", ""),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		ServerEnabled:      getEnvBool("SERVER_ENABLED", true),
		Debug:              getEnvBool("DEBUG", false),
		AdminIDs:           getEnvInt64Slice("ADMIN_IDS", []int64{}),
		CommandPrefix:      getEnv("COMMAND_PREFIX", "/"),
		AIEnabled:          getEnvBool("AI_ENABLED", false),
		AIURL:              getEnv("AI_URL", "https://api.siliconflow.cn/v1/chat/completions"),
		AIKey:              getEnv("AI_KEY", "YOUR_API_KEY_HERE"),
		AIModel:            getEnv("AI_MODEL", "deepseek-ai/DeepSeek-V3.1"),
		AIToolsEnabled:     getEnvBool("AI_TOOLS_ENABLED", true),
		AIVisionEnabled:    getEnvBool("AI_VISION_ENABLED", false),
		AIImageDetail:      getEnv("AI_IMAGE_DETAIL", "auto"),
		SystemPromptPath:   getEnv("SYSTEM_PROMPT_PATH", "configs/system_prompt.md"),
		AnalyzerPromptPath: getEnv("ANALYZER_PROMPT_PATH", "configs/analyzer_prompt.md"),
		HotApiHost:         getEnv("HOT_API_HOST", "https://hot-api.crayoncreator.top"),
		HotApiKey:          getEnv("HOT_API_KEY", "keykeykey"),
		TechPushGroups:     getEnvInt64Slice("TECH_PUSH_GROUPS", []int64{}),
		TechPushUsers:      getEnvInt64Slice("TECH_PUSH_USERS", []int64{}),
		RssPushGroups:      getEnvInt64Slice("RSS_PUSH_GROUPS", []int64{}),
		RssPushUsers:       getEnvInt64Slice("RSS_PUSH_USERS", []int64{}),
		AllowedUsers:       getEnvInt64Slice("ALLOWED_USERS", []int64{}),
		AllowedGroups:      getEnvInt64Slice("ALLOWED_GROUPS", []int64{}),
		RSSApiHost:         getEnv("RSS_API_HOST", "https://rsshub.rssforever.com"),
		AdminUsername:      getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:      getEnv("ADMIN_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
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
			logger.Warn("Invalid boolean value for " + key + ": " + value + ", using default")
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
			logger.Warn("Invalid int64 value in " + key + ": " + part + ", skipping")
			continue
		}

		result = append(result, num)
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type Config struct {
	NapCatHTTPURL   string
	NapCatWSURL     string
	NapCatHTTPToken string
	NapCatWSToken   string
	ServerPort      string
	ServerEnabled   bool
	Debug           bool
	AdminIDs        []int64
	CommandPrefix   string

	AIEnabled bool
	AIURL     string
	AIKey     string

	AIModel string

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
	SerpAPIKey         string
	WeatherAPIKey      string
	AIToolsEnabled     []string
}

func Load() *Config {
	cfg := &Config{
		NapCatHTTPURL:   getEnv("NAPCAT_HTTP_URL", "http://localhost:3000"),
		NapCatWSURL:     getEnv("NAPCAT_WS_URL", "ws://localhost:3001"),
		NapCatHTTPToken: getEnv("NAPCAT_HTTP_TOKEN", ""),
		NapCatWSToken:   getEnv("NAPCAT_WS_TOKEN", ""),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		ServerEnabled:   getEnvBool("SERVER_ENABLED", true),
		Debug:           getEnvBool("DEBUG", false),
		AdminIDs:        getEnvInt64Slice("ADMIN_IDS", []int64{}),
		CommandPrefix:   getEnv("COMMAND_PREFIX", "/"),
		AIEnabled:       getEnvBool("AI_ENABLED", false),
		AIURL:           getEnv("AI_URL", "https://api.siliconflow.cn/v1/chat/completions"),
		AIKey:           getEnv("AI_KEY", "YOUR_API_KEY_HERE"),

		AIModel: getEnv("AI_MODEL", "deepseek/deepseek-r1-turbo"),

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
		SerpAPIKey:         getEnv("SERP_API_KEY", ""),
		WeatherAPIKey:      getEnv("WEATHER_API_KEY", ""),
		AIToolsEnabled:     getEnvStringSlice("AI_TOOLS", []string{}),
	}

	logger.Info("================================================")
	logger.Info("Loaded configuration:")
	logger.Info("================================================")
	logger.Info("  NapCatHTTPURL: " + cfg.NapCatHTTPURL)
	logger.Info("  NapCatWSURL: " + cfg.NapCatWSURL)
	logger.Info("  ServerPort: " + cfg.ServerPort)
	logger.Info("  ServerEnabled: " + strconv.FormatBool(cfg.ServerEnabled))
	logger.Info("  Debug: " + strconv.FormatBool(cfg.Debug))
	logger.Info("  AdminIDs: " + strings.Join(int64SliceToString(cfg.AdminIDs), ","))
	logger.Info("  CommandPrefix: " + cfg.CommandPrefix)
	logger.Info("  AIEnabled: " + strconv.FormatBool(cfg.AIEnabled))
	logger.Info("  AIURL: " + cfg.AIURL)
	logger.Info("  AIModel: " + cfg.AIModel)
	logger.Info("  AIImageDetail: " + cfg.AIImageDetail)
	logger.Info("  SystemPromptPath: " + cfg.SystemPromptPath)
	logger.Info("  AnalyzerPromptPath: " + cfg.AnalyzerPromptPath)
	logger.Info("  HotApiHost: " + cfg.HotApiHost)
	logger.Info("  TechPushGroups: " + strings.Join(int64SliceToString(cfg.TechPushGroups), ","))
	logger.Info("  TechPushUsers: " + strings.Join(int64SliceToString(cfg.TechPushUsers), ","))
	logger.Info("  RssPushGroups: " + strings.Join(int64SliceToString(cfg.RssPushGroups), ","))
	logger.Info("  RssPushUsers: " + strings.Join(int64SliceToString(cfg.RssPushUsers), ","))
	logger.Info("  AllowedUsers: " + strings.Join(int64SliceToString(cfg.AllowedUsers), ","))
	logger.Info("  AllowedGroups: " + strings.Join(int64SliceToString(cfg.AllowedGroups), ","))
	logger.Info("  RSSApiHost: " + cfg.RSSApiHost)
	logger.Info("  AdminUsername: " + cfg.AdminUsername)
	logger.Info("  AIToolsEnabled: " + strings.Join(cfg.AIToolsEnabled, ","))
	logger.Info("  SerpAPIKey: " + maskKey(cfg.SerpAPIKey))
	logger.Info("  WeatherAPIKey: " + maskKey(cfg.WeatherAPIKey))
	logger.Info("================================================")
	return cfg
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

func int64SliceToString(slice []int64) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = strconv.FormatInt(v, 10)
	}
	return result
}

func maskKey(key string) string {
	if key == "" {
		return "(empty)"
	}
	if len(key) <= 4 {
		return "***"
	}
	return key[:2] + "***" + key[len(key)-2:]
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/crayon/wrap-bot/internal/admin"
	adminws "github.com/crayon/wrap-bot/internal/admin/websocket"
	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/internal/shared"
	"github.com/crayon/wrap-bot/internal/tasks"
	"github.com/crayon/wrap-bot/pkgs/bot"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/crayon/wrap-bot/pkgs/napcat"
	"github.com/crayon/wrap-bot/plugins"
	"github.com/joho/godotenv"
)

func main() {
	logger.SetupStdLogger()

	envFile := os.Getenv("APP_ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		logger.Warn("Failed to load .env file: " + err.Error())
	}

	cfg := config.Load()

	engine := bot.New()
	sched := scheduler.New()
	wsHub := adminws.NewHub()

	apiClient := napcat.NewClient(cfg.NapCatHTTPURL, cfg.NapCatHTTPToken)
	wsClient := napcat.NewWSClient(cfg.NapCatWSURL, cfg.NapCatWSToken)

	engine.SetAPIClient(apiClient)
	engine.SetWebSocketClient(wsClient)

	engine.Use(bot.Recovery())
	engine.Use(bot.Logger())
	engine.Use(bot.Authentication(cfg.AllowedUsers, cfg.AllowedGroups))
	engine.Use(bot.InjectAPIClient(apiClient))

	plugins.Register(engine, cfg)
	tasks.RegisterAll(sched, cfg)

	sched.Start()

	shared.SetAdminContext(&shared.AdminContext{
		Engine:    engine,
		Scheduler: sched,
		Config:    cfg,
		WSHub:     wsHub,
	})

	go wsHub.Run()

	logger.GetLogger().SetBroadcastFunc(func(entry logger.LogEntry) {
		wsHub.BroadcastLog(entry)
	})

	go adminws.StartStatusBroadcaster(wsHub, engine, 3*time.Second)

	if cfg.ServerEnabled {
		logger.Info(fmt.Sprintf("Starting admin server on port %s", cfg.ServerPort))
		admin.StartServer(cfg.ServerPort)
	}

	logger.Info(fmt.Sprintf("Starting bot with NapCat WebSocket: %s", cfg.NapCatWSURL))
	if err := engine.Run(); err != nil {
		logger.Error(fmt.Sprintf("Bot stopped with error: %v", err))
		os.Exit(1)
	}
}

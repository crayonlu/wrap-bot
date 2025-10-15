package main

import (
	"log"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/internal/tasks"
	"github.com/crayon/bot_golang/pkgs/bot"
	scheduler "github.com/crayon/bot_golang/pkgs/feature"
	"github.com/crayon/bot_golang/pkgs/napcat"
	"github.com/crayon/bot_golang/plugins"
)

func main() {
	cfg := config.Load()

	engine := bot.New()
	sched := scheduler.New()

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

	log.Printf("Starting bot with NapCat WebSocket: %s", cfg.NapCatWSURL)
	if err := engine.Run(); err != nil {
		log.Fatalf("Bot stopped with error: %v", err)
	}
}

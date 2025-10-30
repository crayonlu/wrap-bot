package main

import (
	"log"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/internal/tasks"
	"github.com/crayon/wrap-bot/pkgs/bot"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
	"github.com/crayon/wrap-bot/pkgs/napcat"
	"github.com/crayon/wrap-bot/plugins"
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

	sched.Start()

	log.Printf("Starting bot with NapCat WebSocket: %s", cfg.NapCatWSURL)
	if err := engine.Run(); err != nil {
		log.Fatalf("Bot stopped with error: %v", err)
	}
}

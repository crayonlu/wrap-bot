package shared

import (
	"github.com/crayon/wrap-bot/internal/admin/websocket"
	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
)

type AdminContext struct {
	Engine    *bot.Engine
	Scheduler *scheduler.Scheduler
	Config    *config.Config
	WSHub     *websocket.Hub
}

var globalContext *AdminContext

func SetAdminContext(ctx *AdminContext) {
	globalContext = ctx
}

func GetAdminContext() *AdminContext {
	return globalContext
}

package admin

import (
	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/bot"
	scheduler "github.com/crayon/wrap-bot/pkgs/feature"
)

type Context struct {
	Engine    *bot.Engine
	Scheduler *scheduler.Scheduler
	Config    *config.Config
}

var globalContext *Context

func SetContext(ctx *Context) {
	globalContext = ctx
}

func GetContext() *Context {
	return globalContext
}

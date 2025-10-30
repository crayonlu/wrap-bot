package tech_push

import (
	"github.com/crayon/wrap-bot/pkgs/feature/tech_push/handlers"
)

type DataSource struct {
	Endpoint string
	Handler  interface{}
}

var dataSources = map[string]DataSource{
	"juejin": {
		Endpoint: "/juejin",
		Handler:  handlers.JuejinHandler,
	},
	"bilibili": {
		Endpoint: "/bilibili",
		Handler:  handlers.BilibiliHandler,
	},
}

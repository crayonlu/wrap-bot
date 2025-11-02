package admin

import (
	"github.com/crayon/wrap-bot/internal/admin/api"
	"github.com/crayon/wrap-bot/internal/admin/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func StartServer(port string) *echo.Echo {
	e := echo.New()

	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	e.POST("/api/auth/login", api.Login)

	admin := e.Group("/api")
	admin.Use(middleware.JWTAuth())

	admin.GET("/status", api.GetStatus)
	admin.GET("/plugins", api.GetPlugins)
	admin.POST("/plugins/:name/toggle", api.TogglePlugin)
	admin.GET("/tasks", api.GetTasks)
	admin.POST("/tasks/:id/trigger", api.TriggerTask)
	admin.GET("/config", api.GetConfig)
	admin.POST("/config", api.UpdateConfig)
	admin.GET("/logs", api.GetLogs)

	go e.Start(":" + port)
	return e
}

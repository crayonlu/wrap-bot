package admin

import (
	"os"

	"github.com/crayon/wrap-bot/internal/admin/api"
	"github.com/crayon/wrap-bot/internal/admin/middleware"
	adminws "github.com/crayon/wrap-bot/internal/admin/websocket"
	"github.com/crayon/wrap-bot/internal/shared"
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
	admin.GET("/presets", api.GetPresets)
	admin.GET("/presets/:filename", api.GetPreset)
	admin.PUT("/presets/:filename", api.UpdatePreset)
	admin.GET("/ai/tools", api.GetAITools)
	admin.GET("/ai/stats", api.GetAIStats)
	admin.POST("/ai/chat", api.TestAIChat)
	admin.POST("/ai/chat/image", api.TestAIImageChat)

	ctx := shared.GetAdminContext()
	if ctx != nil && ctx.WSHub != nil {
		admin.GET("/ws", adminws.HandleWebSocket(ctx.WSHub))
	}

	serveStaticFiles(e)

	go e.Start(":" + port)
	return e
}

func serveStaticFiles(e *echo.Echo) {
	distPath := "web/dist"
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		e.Logger.Warn("web/dist not found, skipping static file serving")
		return
	}

	e.Static("/assets", distPath+"/assets")

	e.GET("/*", func(c echo.Context) error {
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			return echo.ErrNotFound
		}

		indexPath := distPath + "/index.html"
		if _, err := os.Stat(indexPath); err == nil {
			return c.File(indexPath)
		}
		return echo.ErrNotFound
	})
}

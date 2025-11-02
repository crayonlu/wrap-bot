package api

import (
	"net/http"
	"runtime"
	"time"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

var startTime = time.Now()

func GetStatus(c echo.Context) error {
	uptime := time.Since(startTime).Seconds()
	status := types.BotStatus{
		Running:   true,
		Uptime:    int64(uptime),
		Version:   "1.0.0",
		GoVersion: runtime.Version(),
	}
	return c.JSON(http.StatusOK, status)
}

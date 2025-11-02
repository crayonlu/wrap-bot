package api

import (
	"net/http"
	"time"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

func GetLogs(c echo.Context) error {
	logs := []types.LogEntry{
		{
			Timestamp: time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
			Level:     "info",
			Message:   "Bot started successfully",
		},
		{
			Timestamp: time.Now().Add(-30 * time.Second).Format(time.RFC3339),
			Level:     "info",
			Message:   "Connected to NapCat WebSocket",
		},
	}
	return c.JSON(http.StatusOK, logs)
}

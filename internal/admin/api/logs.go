package api

import (
	"net/http"
	"strconv"

	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/labstack/echo/v4"
)

func GetLogs(c echo.Context) error {
	level := c.QueryParam("level")
	limitStr := c.QueryParam("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	log := logger.GetLogger()
	logs := log.GetLogs(level, limit)

	return c.JSON(http.StatusOK, logs)
}

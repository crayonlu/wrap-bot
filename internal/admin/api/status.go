package api

import (
	"net/http"

	"github.com/crayon/wrap-bot/internal/shared"
	"github.com/labstack/echo/v4"
)

func GetStatus(c echo.Context) error {
	ctx := shared.GetAdminContext()
	if ctx == nil || ctx.Engine == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Engine not available",
		})
	}

	status := ctx.Engine.GetStatus()
	return c.JSON(http.StatusOK, status)
}

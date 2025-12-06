package api

import (
	"net/http"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/crayon/wrap-bot/internal/shared"
	"github.com/labstack/echo/v4"
)

func GetPlugins(c echo.Context) error {
	ctx := shared.GetAdminContext()
	if ctx == nil || ctx.Engine == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "engine not available"})
	}

	pluginsMap := ctx.Engine.GetPlugins()
	plugins := []types.PluginStatus{}
	for name, info := range pluginsMap {
		plugins = append(plugins, types.PluginStatus{
			Name:    name,
			Enabled: info.Enabled,
			Discription: info.Description,
		})
	}
	return c.JSON(http.StatusOK, plugins)
}

func TogglePlugin(c echo.Context) error {
	ctx := shared.GetAdminContext()
	if ctx == nil || ctx.Engine == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "engine not available"})
	}

	name := c.Param("name")
	if !ctx.Engine.TogglePlugin(name) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "plugin not found"})
	}

	if ctx.WSHub != nil {
		pluginsMap := ctx.Engine.GetPlugins()
		plugins := []types.PluginStatus{}
		for name, info := range pluginsMap {
			plugins = append(plugins, types.PluginStatus{
				Name:    name,
				Enabled: info.Enabled,
			})
		}
		ctx.WSHub.BroadcastPlugins(plugins)
	}

	enabled := ctx.Engine.IsPluginEnabled(name)
	return c.JSON(http.StatusOK, types.PluginStatus{
		Name:    name,
		Enabled: enabled,
	})
}

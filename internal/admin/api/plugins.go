package api

import (
	"net/http"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

var pluginsState = map[string]bool{
	"ai_chat":   true,
	"tech_push": true,
	"rss_push":  true,
	"ping":      true,
	"echo":      true,
}

func GetPlugins(c echo.Context) error {
	plugins := []types.PluginStatus{}
	for name, enabled := range pluginsState {
		plugins = append(plugins, types.PluginStatus{
			Name:    name,
			Enabled: enabled,
		})
	}
	return c.JSON(http.StatusOK, plugins)
}

func TogglePlugin(c echo.Context) error {
	name := c.Param("name")
	if _, exists := pluginsState[name]; !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "plugin not found"})
	}

	pluginsState[name] = !pluginsState[name]
	return c.JSON(http.StatusOK, types.PluginStatus{
		Name:    name,
		Enabled: pluginsState[name],
	})
}

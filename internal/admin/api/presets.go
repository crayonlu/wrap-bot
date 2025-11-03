package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/labstack/echo/v4"
)

type PresetFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type UpdatePresetRequest struct {
	Content string `json:"content"`
}

func getConfiguredPath(filename string) string {
	cfg := config.Load()

	switch filename {
	case "system_prompt.md":
		return cfg.SystemPromptPath
	case "analyzer_prompt.md":
		return cfg.AnalyzerPromptPath
	default:
		return ""
	}
}

func GetPresets(c echo.Context) error {
	cfg := config.Load()

	presets := []PresetFile{}
	configFiles := map[string]string{
		"system_prompt.md":   cfg.SystemPromptPath,
		"analyzer_prompt.md": cfg.AnalyzerPromptPath,
	}

	for filename, path := range configFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to read preset %s from %s: %v", filename, path, err))
			continue
		}

		presets = append(presets, PresetFile{
			Name:    filename,
			Path:    path,
			Content: string(content),
		})
	}

	return c.JSON(http.StatusOK, presets)
}

func GetPreset(c echo.Context) error {
	filename := c.Param("filename")

	if filepath.Ext(filename) != ".md" || filepath.Base(filename) != filename {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid filename",
		})
	}

	filePath := getConfiguredPath(filename)
	if filePath == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Preset not found",
		})
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Preset not found",
		})
	}

	return c.JSON(http.StatusOK, PresetFile{
		Name:    filename,
		Path:    filePath,
		Content: string(content),
	})
}

func UpdatePreset(c echo.Context) error {
	filename := c.Param("filename")

	if filepath.Ext(filename) != ".md" || filepath.Base(filename) != filename {
		logger.Warn(fmt.Sprintf("Invalid filename attempt: %s", filename))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid filename",
		})
	}

	var req UpdatePresetRequest
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("Failed to bind request: %v", err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	writePath := getConfiguredPath(filename)
	if writePath == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Preset not found",
		})
	}

	dir := filepath.Dir(writePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create directory %s: %v", dir, err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create directory",
		})
	}

	if err := os.WriteFile(writePath, []byte(req.Content), 0644); err != nil {
		logger.Error(fmt.Sprintf("Failed to write file %s: %v", writePath, err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to update preset: %v", err),
		})
	}

	logger.Info(fmt.Sprintf("Successfully updated preset: %s at %s", filename, writePath))

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Preset updated successfully",
	})
}

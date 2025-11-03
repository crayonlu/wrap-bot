package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

func GetPresets(c echo.Context) error {
	configsDir := "configs"

	files, err := os.ReadDir(configsDir)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read configs directory",
		})
	}

	var presets []PresetFile
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(configsDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		presets = append(presets, PresetFile{
			Name:    file.Name(),
			Path:    filePath,
			Content: string(content),
		})
	}

	return c.JSON(http.StatusOK, presets)
}

func GetPreset(c echo.Context) error {
	filename := c.Param("filename")
	filePath := filepath.Join("configs", filename)

	if filepath.Ext(filename) != ".md" || filepath.Base(filename) != filename {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid filename",
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
	filePath := filepath.Join("configs", filename)

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

	configsDir := "configs"
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create configs directory: %v", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create configs directory",
		})
	}

	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		logger.Error(fmt.Sprintf("Failed to write file %s: %v", filePath, err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to update preset: %v", err),
		})
	}

	logger.Info(fmt.Sprintf("Successfully updated preset: %s", filename))
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Preset updated successfully",
	})
}

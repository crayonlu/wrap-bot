package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/labstack/echo/v4"
)

const (
	defaultReadDir  = "configs"
	defaultWriteDir = "/data/configs"
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
		return filepath.Join(defaultReadDir, filename)
	}
}

func getWritablePath(filename string) string {
	configuredPath := getConfiguredPath(filename)

	dir := filepath.Dir(configuredPath)

	testFile := filepath.Join(dir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err == nil {
		os.Remove(testFile)
		return configuredPath
	}

	return filepath.Join(defaultWriteDir, filename)
}

func GetPresets(c echo.Context) error {
	var presets []PresetFile
	seen := make(map[string]bool)

	dirs := []string{defaultWriteDir, defaultReadDir}
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
				continue
			}

			if seen[file.Name()] {
				continue
			}
			seen[file.Name()] = true

			filePath := filepath.Join(dir, file.Name())
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
	}

	cfg := config.Load()
	configFiles := map[string]string{
		"system_prompt.md":   cfg.SystemPromptPath,
		"analyzer_prompt.md": cfg.AnalyzerPromptPath,
	}

	for filename, path := range configFiles {
		if seen[filename] {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
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

	dirs := []string{defaultWriteDir, defaultReadDir}
	var content []byte
	var err error
	var foundPath string

	for _, dir := range dirs {
		filePath := filepath.Join(dir, filename)
		content, err = os.ReadFile(filePath)
		if err == nil {
			foundPath = filePath
			break
		}
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Preset not found",
		})
	}

	return c.JSON(http.StatusOK, PresetFile{
		Name:    filename,
		Path:    foundPath,
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

	writePath := getWritablePath(filename)

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

	if strings.Contains(writePath, defaultWriteDir) {
		configuredPath := getConfiguredPath(filename)
		logger.Info(fmt.Sprintf("Note: Configured path %s is read-only, saved to %s instead", configuredPath, writePath))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Preset updated successfully",
	})
}

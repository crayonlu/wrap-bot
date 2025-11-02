package api

import (
	"net/http"
	"time"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/labstack/echo/v4"
)

func GetTasks(c echo.Context) error {
	tasks := []types.TaskStatus{
		{
			ID:         "tech-push-daily",
			Name:       "Tech Push Daily",
			Schedule:   "0 12 * * *",
			NextRun:    time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			LastRun:    time.Now().Add(-22 * time.Hour).Format(time.RFC3339),
			Status:     "active",
			CanTrigger: true,
		},
		{
			ID:         "rss-push-daily",
			Name:       "RSS Push Daily",
			Schedule:   "0 9 * * *",
			NextRun:    time.Now().Add(5 * time.Hour).Format(time.RFC3339),
			LastRun:    time.Now().Add(-19 * time.Hour).Format(time.RFC3339),
			Status:     "active",
			CanTrigger: true,
		},
	}
	return c.JSON(http.StatusOK, tasks)
}

func TriggerTask(c echo.Context) error {
	taskID := c.Param("id")
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "triggered",
		"task_id": taskID,
	})
}

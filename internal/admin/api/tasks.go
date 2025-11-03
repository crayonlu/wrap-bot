package api

import (
	"net/http"

	"github.com/crayon/wrap-bot/internal/admin/types"
	"github.com/crayon/wrap-bot/internal/shared"
	"github.com/labstack/echo/v4"
)

func GetTasks(c echo.Context) error {
	ctx := shared.GetAdminContext()
	if ctx == nil || ctx.Scheduler == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "scheduler not available",
		})
	}

	schedulerTasks := ctx.Scheduler.GetTasks()
	entries := ctx.Scheduler.GetCronEntries()
	tasks := make([]types.TaskStatus, 0, len(schedulerTasks))

	for _, task := range schedulerTasks {
		var nextRun, lastRun string
		if entry, ok := entries[task.EntryID]; ok {
			nextRun = entry.Next.Format("2006-01-02T15:04:05Z07:00")
			if !entry.Prev.IsZero() {
				lastRun = entry.Prev.Format("2006-01-02T15:04:05Z07:00")
			}
		}

		tasks = append(tasks, types.TaskStatus{
			ID:         task.ID,
			Name:       task.Name,
			Schedule:   task.Schedule,
			NextRun:    nextRun,
			LastRun:    lastRun,
			Status:     "active",
			CanTrigger: true,
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

func TriggerTask(c echo.Context) error {
	taskID := c.Param("id")

	ctx := shared.GetAdminContext()
	if ctx == nil || ctx.Scheduler == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "scheduler not available",
		})
	}

	success := ctx.Scheduler.TriggerTask(taskID)
	if !success {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "task not found or cannot be triggered",
		})
	}

	if ctx.WSHub != nil {
		schedulerTasks := ctx.Scheduler.GetTasks()
		entries := ctx.Scheduler.GetCronEntries()
		tasks := make([]types.TaskStatus, 0, len(schedulerTasks))

		for _, task := range schedulerTasks {
			var nextRun, lastRun string
			if entry, ok := entries[task.EntryID]; ok {
				nextRun = entry.Next.Format("2006-01-02T15:04:05Z07:00")
				if !entry.Prev.IsZero() {
					lastRun = entry.Prev.Format("2006-01-02T15:04:05Z07:00")
				}
			}

			tasks = append(tasks, types.TaskStatus{
				ID:         task.ID,
				Name:       task.Name,
				Schedule:   task.Schedule,
				NextRun:    nextRun,
				LastRun:    lastRun,
				Status:     "active",
				CanTrigger: true,
			})
		}
		ctx.WSHub.BroadcastTasks(tasks)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "triggered",
		"task_id": taskID,
	})
}

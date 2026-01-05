package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	db "tasksphere-backend/db/sqlc"
)

type TaskHandler struct {
	DB      *pgxpool.Pool
	Queries *db.Queries
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var task struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  string `json:"status"`
	}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if task.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	status := task.Status
	if status == "" {
		status = "todo"
	}
	if !isValidTaskStatus(status) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = h.Queries.CreateTask(context.Background(), h.DB, db.CreateTaskParams{
		Title:   task.Title,
		UserID:  userID,
		Content: task.Content,
		Status:  status,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save task"})
	}
	return c.JSON(fiber.Map{"message": "Task created successfully!"})
}

func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	results, err := h.Queries.ListTasksByUser(context.Background(), h.DB, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}

	var tasks []map[string]interface{}
	for _, task := range results {
		tasks = append(tasks, map[string]interface{}{
			"id":      task.ID,
			"title":   task.Title,
			"content": task.Content,
			"status":  task.Status,
		})
	}
	return c.JSON(tasks)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	idParam := c.Params("id")
	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task id"})
	}

	var task struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  string `json:"status"`
	}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if task.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	status := task.Status
	if status == "" {
		currentStatus, err := h.Queries.GetTaskStatusByID(context.Background(), h.DB, db.GetTaskStatusByIDParams{
			ID:     taskID,
			UserID: userID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
		}
		status = currentStatus
	} else if !isValidTaskStatus(status) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	updated, err := h.Queries.UpdateTask(context.Background(), h.DB, db.UpdateTaskParams{
		Title:   task.Title,
		Content: task.Content,
		Status:  status,
		ID:      taskID,
		UserID:  userID,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}
	if updated == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}
	return c.JSON(fiber.Map{"message": "Task updated successfully"})
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")
	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task id"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	deleted, err := h.Queries.DeleteTask(context.Background(), h.DB, db.DeleteTaskParams{
		ID:     taskID,
		UserID: userID,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}
	if deleted == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted successfully!"})
}

func getUserIDFromContext(c *fiber.Ctx) (int64, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	val, ok := claims["user_id"]
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	floatVal, ok := val.(float64)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return int64(floatVal), nil
}

func isValidTaskStatus(status string) bool {
	switch status {
	case "todo", "in_progress", "done":
		return true
	default:
		return false
	}
}

package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskHandler struct {
	DB *pgxpool.Pool
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var task struct {
		Title string `json:"title"`
	}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if task.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	_, err := h.DB.Exec(context.Background(),
		"INSERT INTO tasks (title, user_id) VALUES ($1, $2)", task.Title, userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save task"})
	}
	return c.JSON(fiber.Map{"message": "Task created successfully!"})
}

func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	rows, err := h.DB.Query(context.Background(),
		"SELECT id, title FROM tasks WHERE user_id=$1", userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse tasks"})
		}
		tasks = append(tasks, map[string]interface{}{
			"id":    id,
			"title": title,
		})
	}
	return c.JSON(tasks)

}

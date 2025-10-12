package main

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/jwt/v3"
	// "github.com/jackc/pgx/v5"
)

func main() {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/tasks", func(c *fiber.Ctx) error {
		tasks := []map[string]string{
			{"id": "1", "title": "Task 1"},
			{"id": "2", "title": "Task 2"},
		}
		return c.JSON(tasks)
	})
	app.Post("/tasks", func(c *fiber.Ctx) error {
		type Task struct {
			Title string `json:"title"`
		}

		var newTask Task

		if err := c.BodyParser(&newTask); err != nil {
			return c.Status(400).SendString("Geçersiz istek gövdesi")
		}

		return c.Status(201).JSON(fiber.Map{
			"id":    "3",
			"title": newTask.Title,
		})
	})

	app.Listen(":8080")
}

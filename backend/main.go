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
	app.Listen(":8080")
}

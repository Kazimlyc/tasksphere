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
	app.Listen(":8080")
}

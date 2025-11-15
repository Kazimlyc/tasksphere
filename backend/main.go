package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"tasksphere-backend/config"
	"tasksphere-backend/handlers"
	"tasksphere-backend/middleware"
)

func main() {
	cfg := config.Load()

	var pool *pgxpool.Pool

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	app := fiber.New()

	authHandler := handlers.AuthHandler{
		DB:        pool,
		JWTSecret: cfg.JWTSecret,
	}
	taskHandler := handlers.TaskHandler{
		DB: pool,
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Use("/tasks", middleware.JWTProtected((cfg.JWTSecret)))
	app.Post("/tasks", taskHandler.CreateTask)
	app.Get("/tasks", taskHandler.GetTasks)
	app.Put("/tasks/:id", taskHandler.UpdateTask)
	app.Delete("/tasks/:id", taskHandler.DeleteTask)

	log.Fatal(app.Listen(":" + cfg.Port))
}

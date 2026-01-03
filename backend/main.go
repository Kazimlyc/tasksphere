package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"tasksphere-backend/config"
	db "tasksphere-backend/db/sqlc"
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
	app.Use(cors.New())
	queries := db.New()

	authHandler := handlers.AuthHandler{
		DB:        pool,
		Queries:   queries,
		JWTSecret: cfg.JWTSecret,
	}
	taskHandler := handlers.TaskHandler{
		DB:      pool,
		Queries: queries,
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

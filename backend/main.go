package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"

	"tasksphere-backend/handlers"
	"tasksphere-backend/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var pool *pgxpool.Pool
	// Connect to PostgreSQL
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL ✅")

	app := fiber.New()

	authHandler := handlers.AuthHandler{
		DB:        pool,
		JWTSecret: jwtSecret,
	}
	taskHandler := handlers.TaskHandler{
		DB: pool,
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Use("/tasks", middleware.JWTProtected((jwtSecret)))
	app.Post("/tasks", taskHandler.CreateTask)
	app.Get("/tasks", taskHandler.GetTasks)
	app.Put("/tasks/:id", taskHandler.UpdateTask)
	app.Delete("tasks/:id", taskHandler.DeleteTask)

	log.Fatal(app.Listen(":8080"))
}

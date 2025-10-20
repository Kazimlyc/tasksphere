package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

func main() {
	var conn *pgx.Conn
	// Connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), "postgres://admin:secret@localhost:5432/tasksphere")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	log.Println("Connected to PostgreSQL ✅")

	_, err = conn.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS  tasks(
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
	log.Println("Table 'tasks' is ready ✅")

	_, err = conn.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v\n", err)
	}
	log.Println("Table 'users' is ready ✅")

	app := fiber.New()

	export JWT_SECRET="supergizlianahtar123"
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
	jwtSecret = "dev-secret-change-me"
}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Post("/login", func(c *fiber.Ctx)error{
		var payload struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&payload); er != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		if payload.Email == "" || payload.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error":"email and password required"})
		}
		
	})
	app.Post("/tasks", func(c *fiber.Ctx) error {
		var task struct {
			Title string `json:"title"`
		}

		if err := c.BodyParser(&task); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		if task.Title == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
		}
		_, err := conn.Exec(context.Background(),
			"INSERT INTO tasks (title) VALUES ($1)", task.Title)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save task"})
		}
		return c.JSON(fiber.Map{"message": "Task created successfully!"})
	})
	app.Post("/register", func(c *fiber.Ctx) error {
		var user struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		if user.Email == "" || user.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Email and Password are required"})
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		_, err = conn.Exec(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)`, user.Email, string(hashedPassword))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
		}
		return c.JSON(fiber.Map{"message": "User registered successfully!"})
	})
	app.Get("/tasks", func(c *fiber.Ctx) error {
		rows, err := conn.Query(context.Background(), "SELECT id, title FROM tasks")
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
	})
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var task struct {
			Title string `json:"title"`
		}
		if err := c.BodyParser(&task); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		if task.Title == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
		}
		_, err := conn.Exec(context.Background(),
			"UPDATE tasks SET title=$1 WHERE id=$2", task.Title, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
		}
		return c.JSON(fiber.Map{"message": "Task updated successfully!"})
	})
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := conn.Exec(context.Background(),
			"DELETE FROM tasks WHERE id=$1", id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
		}
		return c.JSON(fiber.Map{"message": "Task deleted successfully!"})
	})

	app.Listen(":8080")
}

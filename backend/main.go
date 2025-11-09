package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		if payload.Email == "" || payload.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "email and password required"})
		}

		var id int
		var hashedPassword string
		err := pool.QueryRow(context.Background(),
			"SELECT id, password FROM users WHERE email=$1", payload.Email).Scan(&id, &hashedPassword)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
		}

		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password)) != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": id,
			"email":   payload.Email,
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}
		return c.JSON(fiber.Map{
			"token": tokenString,
		})
	})
	app.Use("/tasks", jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
	}))
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
		_, err := pool.Exec(context.Background(),
			"INSERT INTO tasks (title) VALUES ($1)", task.Title)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save task"})
		}

		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))
		_, err = pool.Exec(context.Background(),
			"INSERT INTO tasks (title, user_id) VALUES ($1, $2)", task.Title, userID)

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
		_, err = pool.Exec(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)`, user.Email, string(hashedPassword))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
		}
		return c.JSON(fiber.Map{"message": "User registered successfully!"})
	})
	app.Get("/tasks", func(c *fiber.Ctx) error {
		rows, err := pool.Query(context.Background(), "SELECT id, title FROM tasks")
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
		_, err := pool.Exec(context.Background(),
			"UPDATE tasks SET title=$1 WHERE id=$2", task.Title, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
		}
		return c.JSON(fiber.Map{"message": "Task updated successfully!"})
	})
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := pool.Exec(context.Background(),
			"DELETE FROM tasks WHERE id=$1", id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
		}
		return c.JSON(fiber.Map{"message": "Task deleted successfully!"})
	})

	log.Fatal(app.Listen(":8080"))
}

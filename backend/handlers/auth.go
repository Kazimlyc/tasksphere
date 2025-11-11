package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB        *pgxpool.Pool
	JWTSecret string
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
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
	_, err = h.DB.Exec(context.Background(), `
		INSERT INTO users (email,password) VALUES ($1, $2)`, user.Email, string(hashedPassword))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
	}
	return c.JSON(fiber.Map{"message": "User registered successfully!"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

	err := h.DB.QueryRow(context.Background(),
		"SELECT id, password FROM users WHERE email=$1", payload.Email).Scan(&id, &hashedPassword)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   payload.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		log.Println("JWT Error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create token"})
	}
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}

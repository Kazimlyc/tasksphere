package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func getUserIDFromContext(c *fiber.Ctx) (int64, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	val, ok := claims["user_id"]
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	floatVal, ok := val.(float64)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return int64(floatVal), nil
}

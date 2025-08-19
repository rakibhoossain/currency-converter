package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware that validates Bearer token
func AuthMiddleware(authToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Authorization header is required",
			})
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Authorization header must start with 'Bearer '",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate token
		if token != authToken {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid authorization token",
			})
		}

		// Token is valid, continue to next handler
		return c.Next()
	}
}

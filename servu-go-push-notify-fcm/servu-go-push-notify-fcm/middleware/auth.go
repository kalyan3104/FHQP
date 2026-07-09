package middleware

import (
	"strings"

	"github.com/baryogenesis2025/servu-go/auth"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware validates JWT tokens and attaches user info to context
func JWTMiddleware(authService *auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check Bearer format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
		}

		token := parts[1]

		// Verify token
		claims, err := authService.VerifyAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Attach user info to context
		c.Locals("user_id", claims.UserID)
		c.Locals("phone_number", claims.PhoneNumber)
		c.Locals("role", claims.Role)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User role not found",
			})
		}

		roleStr := userRole.(string)

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if roleStr == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// OptionalJWT - similar to JWTMiddleware but doesn't fail if token is missing
func OptionalJWT(authService *auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		claims, err := authService.VerifyAccessToken(token)
		if err == nil {
			// Valid token - attach to context
			c.Locals("user_id", claims.UserID)
			c.Locals("phone_number", claims.PhoneNumber)
			c.Locals("role", claims.Role)
			c.Locals("claims", claims)
		}

		return c.Next()
	}
}

// GetUserID - helper to extract user ID from context
func GetUserID(c *fiber.Ctx) (uint, bool) {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0, false
	}
	return userID.(uint), true
}

// GetUserRole - helper to extract user role from context
func GetUserRole(c *fiber.Ctx) (string, bool) {
	role := c.Locals("role")
	if role == nil {
		return "", false
	}
	return role.(string), true
}

// MustGetUserID - helper that returns error if user ID not found
func MustGetUserID(c *fiber.Ctx) (uint, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return userID.(uint), nil
}

// GetClaims - helper to extract full c

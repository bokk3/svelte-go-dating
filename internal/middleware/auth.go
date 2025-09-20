package middleware

import (
    "strings"
    
    "github.com/gofiber/fiber/v2"
    "dating-svelte/internal/auth"
)

// AuthRequired middleware validates JWT tokens
func AuthRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Authorization header required",
            })
        }
        
        // Extract token from "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid authorization header format",
            })
        }
        
        tokenString := parts[1]
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }
        
        // Store user information in context for use in handlers
        c.Locals("user_id", claims.UserID)
        c.Locals("user_email", claims.Email)
        c.Locals("is_premium", claims.IsPremium)
        
        return c.Next()
    }
}

// PremiumRequired middleware checks if user has premium subscription
func PremiumRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        isPremium, ok := c.Locals("is_premium").(bool)
        if !ok || !isPremium {
            return c.Status(403).JSON(fiber.Map{
                "error": "Premium subscription required",
            })
        }
        
        return c.Next()
    }
}
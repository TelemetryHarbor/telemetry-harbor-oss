package middleware

import (
	"go-ingest-service/internal/config"
	"go-ingest-service/internal/models"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth extracts and verifies the X-API-Key header.
func APIKeyAuth(c *fiber.Ctx) error {
	log.Printf("[Auth] APIKeyAuth: Starting for %s %s", c.Method(), c.OriginalURL())
	apiKey := c.Get("X-API-Key")

	// Check if the API key is missing
	if apiKey == "" {
		log.Printf("[Auth] APIKeyAuth: API key missing for %s %s. Returning 403 Forbidden.", c.Method(), c.OriginalURL())
		return c.Status(http.StatusForbidden).JSON(models.APIError{
			Message: "API key is missing",
		})
	}

	// Verify the API key against the one from the configuration
	if apiKey != config.AppConfig.APIKey {
		log.Printf("[Auth] APIKeyAuth: Invalid API key provided for %s %s. Returning 401 Unauthorized.", c.Method(), c.OriginalURL())
		return c.Status(http.StatusUnauthorized).JSON(models.APIError{
			Message: "Invalid API key",
		})
	}

	// If verification passes, store the key and proceed
	c.Locals("apiKey", apiKey)
	log.Printf("[Auth] APIKeyAuth: API key verified successfully. Proceeding to next middleware/handler.")
	return c.Next()
}
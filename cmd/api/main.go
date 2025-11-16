package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"go-ingest-service/internal/cache"
	"go-ingest-service/internal/config"
	"go-ingest-service/internal/db"
	general_handler "go-ingest-service/internal/ingest/general"
	mw "go-ingest-service/internal/middleware"
	"go-ingest-service/internal/models"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("[API] Failed to load configuration: %v", err)
	}
 

	app := fiber.New(fiber.Config{
		ErrorHandler: customHTTPErrorHandler,
	})

	// --- Standard Middleware ---
	app.Use(logger.New(logger.Config{
		Format: "[API] ${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.APP_URL,
		AllowMethods:     "POST,GET,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,X-API-Key",
		AllowCredentials: true,
	}))
	app.Use(compress.New())

	// --- Connections ---
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("[API] Failed to connect to database: %v", err)
	}
	defer db.CloseDB()
	if err := cache.ConnectRedis(); err != nil {
		log.Fatalf("[API] CRITICAL: Failed to connect to Redis: %v", err)
	}
	defer cache.CloseRedis()

	// --- Health Check Route ---
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	// --- ROUTING ---
	apiv1 := app.Group("/api/v2")

	// Define the shared middleware chain for ingest routes.
	ingestChain := []fiber.Handler{
		mw.APIKeyAuth,
	}
	apiv1.Post("/ingest", append(ingestChain, general_handler.IngestData)...)
	apiv1.Post("/ingest/batch", append(ingestChain, general_handler.IngestBatchData)...)


	// --- Graceful Shutdown ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	listenAddr := ":" + config.AppConfig.Port
	go func() {
		log.Printf("[API] Starting Go Ingest Service on %s", listenAddr)
		if err := app.Listen(listenAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[API] Server listener failed: %v", err)
		}
	}()

	sig := <-sigChan
	log.Printf("[API] Received signal %v, initiating graceful shutdown...", sig)
	if err := app.Shutdown(); err != nil {
		log.Printf("[API] Server shutdown failed: %v", err)
	}
	log.Println("[API] Server gracefully shut down.")
}

// customHTTPErrorHandler handles all HTTP errors.
func customHTTPErrorHandler(c *fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	message := "An unexpected internal server error occurred."

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("[API] ErrorHandler: URL=%s, Error=%v, Status=%d", c.OriginalURL(), err, code)

	return c.Status(code).JSON(models.APIError{
		Message: "Internal Server Error",
		Details: message, // The specific error message goes into 'details'
	})
}
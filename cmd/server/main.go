package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"

	"dating-svelte/internal/database"
	"dating-svelte/internal/handlers"
	"dating-svelte/internal/middleware"
	wshandler "dating-svelte/internal/websocket"
)

var wsHub *wshandler.Hub

func main() {
	// Initialize database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:pass@localhost:5432/dating_db?sslmode=disable"
	}

	db, err := database.New(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize WebSocket hub
	wsHub = wshandler.NewHub(db)
	go wsHub.Run()

	// Initialize handlers with database
	handlers.InitializeHandlers(db)

	app := fiber.New(fiber.Config{
		Prefork:     false, // Disable for development
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:      100,
		Duration: time.Minute,
	}))

	// Serve static files from the dist directory
	app.Static("/", "./dist", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		Index:         "index.html",
		CacheDuration: 1 * time.Hour,
		MaxAge:        3600,
	})

	// SPA fallback - redirect all non-API routes to index.html
	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if path == "/" ||
			strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/ws") ||
			strings.Contains(path, ".") {
			return c.Next()
		}
		return c.SendFile("./dist/index.html")
	})

	// Routes
	setupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Auth routes
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	api.Post("/refresh", handlers.RefreshToken)

	// Protected routes
	protected := api.Use(middleware.AuthRequired())
	protected.Get("/profile", handlers.GetProfile)
	protected.Put("/profile", handlers.UpdateProfile)
	protected.Get("/matches", handlers.GetMatches)
	protected.Post("/swipe", handlers.Swipe)

	// WebSocket for real-time messaging (with auth)
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		return websocket.New(wshandler.HandleWebSocket(wsHub, userID))(c)
	})

	// Payment routes
	protected.Post("/subscribe", handlers.CreateSubscription)
	protected.Post("/crypto-payment", handlers.ProcessCryptoPayment)

	// GDPR routes
	protected.Post("/gdpr/export", handlers.ExportData)
	protected.Delete("/gdpr/delete", handlers.DeleteAccount)
}

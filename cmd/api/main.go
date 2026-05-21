package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rwndy/bookmark-api/internal/config"
	"github.com/rwndy/bookmark-api/internal/handler"
	"github.com/rwndy/bookmark-api/internal/middleware"
	"github.com/rwndy/bookmark-api/internal/repository"
	"github.com/rwndy/bookmark-api/internal/service"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	log.Printf("DSN: %s", cfg.DBDsn)

	// Database
	db, err := gorm.Open(postgres.Open(cfg.DBDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	log.Println("database connected")

	// Dependency injection
	userRepo := repository.NewUserRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	authSvc := service.NewAuthService(userRepo)
	bookmarkSvc := service.NewBookmarkService(bookmarkRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkSvc)

	// Router
	app := fiber.New()
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Public
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	// Protected
	api := app.Group("/api", middleware.AuthRequired())
	api.Post("/bookmarks", bookmarkHandler.Create)
	api.Get("/bookmarks", bookmarkHandler.List)
	api.Put("/bookmarks/:id", bookmarkHandler.Update)
	api.Delete("/bookmarks/:id", bookmarkHandler.Delete)

	log.Printf("server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
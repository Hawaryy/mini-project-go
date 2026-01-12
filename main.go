package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourusername/book-management-api/config"
	"github.com/yourusername/book-management-api/handlers"
	"github.com/yourusername/book-management-api/middleware"
	"github.com/yourusername/book-management-api/migration"
	"github.com/yourusername/book-management-api/repository"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := migration.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(db)
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	bookHandler := handlers.NewBookHandler(bookRepo, categoryRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup Gin router
	router := gin.Default()

	// Public routes (no auth required)
	router.POST("/api/users/login", userHandler.Login)

	// Protected routes with JWT or Basic Auth
	authMiddleware := middleware.AuthMiddleware(userRepo)

	// Category routes
	router.GET("/api/categories", authMiddleware, categoryHandler.GetAllCategories)
	router.POST("/api/categories", authMiddleware, categoryHandler.CreateCategory)
	router.GET("/api/categories/:id", authMiddleware, categoryHandler.GetCategoryByID)
	router.DELETE("/api/categories/:id", authMiddleware, categoryHandler.DeleteCategory)
	router.GET("/api/categories/:id/books", authMiddleware, bookHandler.GetBooksByCategory)

	// Book routes
	router.GET("/api/books", authMiddleware, bookHandler.GetAllBooks)
	router.POST("/api/books", authMiddleware, bookHandler.CreateBook)
	router.GET("/api/books/:id", authMiddleware, bookHandler.GetBookByID)
	router.DELETE("/api/books/:id", authMiddleware, bookHandler.DeleteBook)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

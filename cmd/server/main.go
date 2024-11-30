package main

import (
	"log"
	"os"

	"todo/docs"
	"todo/internal/db"
	"todo/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Todo API
// @version 1.0
// @description A simple todo list API with PostgreSQL backend
// @host localhost:8081
// @BasePath /
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize database
	db.InitDB()

	// Create a new gin router with default middleware
	router := gin.Default()

	// Configure Swagger documentation
	setupSwagger()

	// Setup routes
	setupRoutes(router)

	// Get port from environment or use default
	port := getPort()

	// Start server
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupSwagger() {
	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "A simple todo list API with PostgreSQL backend"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func setupRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "OK")
	})

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Todo endpoints
	router.POST("/todos", handlers.CreateTodo)
	router.GET("/todos", handlers.GetTodos)
	router.GET("/todos/:id", handlers.GetTodo)
	router.PUT("/todos/:id", handlers.UpdateTodo)
	router.DELETE("/todos/:id", handlers.DeleteTodo)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	return port
}

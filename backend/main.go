package main

import (
	"log"
	"os"

	"appdrop/database"
	"appdrop/handlers"
	"appdrop/middleware"
	"appdrop/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	dbConfig := database.NewConfigFromEnv()
	if err := database.Connect(dbConfig); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize repositories
	db := database.GetDB()
	pageRepo := repository.NewPageRepository(db)
	widgetRepo := repository.NewWidgetRepository(db)

	// Initialize handlers
	pageHandler := handlers.NewPageHandler(pageRepo, widgetRepo)
	widgetHandler := handlers.NewWidgetHandler(widgetRepo, pageRepo)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "Mini App Config API",
			"version": "1.0.0",
		})
	})

	// API routes
	api := router.Group("/")
	{
		// Page routes
		api.GET("/pages", pageHandler.ListPages)
		api.GET("/pages/:id", pageHandler.GetPage)
		api.POST("/pages", pageHandler.CreatePage)
		api.PUT("/pages/:id", pageHandler.UpdatePage)
		api.DELETE("/pages/:id", pageHandler.DeletePage)

		// Widget routes
		api.GET("/pages/:id/widgets", widgetHandler.GetWidgets)
		api.POST("/pages/:id/widgets", widgetHandler.CreateWidget)
		api.POST("/pages/:id/widgets/reorder", widgetHandler.ReorderWidgets)
		api.PUT("/widgets/:id", widgetHandler.UpdateWidget)
		api.DELETE("/widgets/:id", widgetHandler.DeleteWidget)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Mini App Config API running on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

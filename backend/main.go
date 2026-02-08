// Package main serves as the entry point for the Mini App Configuration API.
// It initializes the database connection, repositories, handlers, and the web server.
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

// main is the primary execution function that sets up the application infrastructure,
// configures the routing engine, and starts the HTTP server.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbConfig := database.NewConfigFromEnv()
	if err := database.Connect(dbConfig); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	db := database.GetDB()
	pageRepo := repository.NewPageRepository(db)
	widgetRepo := repository.NewWidgetRepository(db)

	pageHandler := handlers.NewPageHandler(pageRepo, widgetRepo)
	widgetHandler := handlers.NewWidgetHandler(widgetRepo, pageRepo)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "Mini App Config API",
			"version": "1.0.0",
		})
	})

	api := router.Group("/")
	{
		api.GET("/pages", pageHandler.ListPages)
		api.GET("/pages/:id", pageHandler.GetPage)
		api.POST("/pages", pageHandler.CreatePage)
		api.PUT("/pages/:id", pageHandler.UpdatePage)
		api.DELETE("/pages/:id", pageHandler.DeletePage)

		api.GET("/pages/:id/widgets", widgetHandler.GetWidgets)
		api.POST("/pages/:id/widgets", widgetHandler.CreateWidget)
		api.POST("/pages/:id/widgets/reorder", widgetHandler.ReorderWidgets)
		api.PUT("/widgets/:id", widgetHandler.UpdateWidget)
		api.DELETE("/widgets/:id", widgetHandler.DeleteWidget)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Mini App Config API running on port %s", port)
	log.Printf("API Documentation: http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

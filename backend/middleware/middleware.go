// Package middleware provides HTTP interceptors for logging, security, and error recovery.
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger initializes and returns a middleware handler for structured request/response logging.
// It captures execution latency, status codes, and client information for each incoming request.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%s] %d | %v | %s | %s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %s", e.Error())
			}
		}
	}
}

// CORS initializes and returns a middleware handler for Cross-Origin Resource Sharing (CORS) configuration.
// It sets permissive headers to facilitate communication with front-end applications hosted on different origins.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Recovery initializes and returns a panic-recovery middleware provided by the Gin framework.
// It ensures that the server gracefully recovers from unexpected runtime panics and returns a 500 error instead of crashing.
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

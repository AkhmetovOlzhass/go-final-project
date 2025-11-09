package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		println(
			"Method:", c.Request.Method,
			"Path:", c.Request.URL.Path,
			"Status:", c.Writer.Status(),
			"Duration:", duration.String(),
		)
	}
}
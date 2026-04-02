package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start).Milliseconds()
		log.Printf("method=%s path=%s status=%d latency_ms=%d ip=%s",
			c.Request.Method, c.FullPath(), c.Writer.Status(), latency, c.ClientIP())
	}
}

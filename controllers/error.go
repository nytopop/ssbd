package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/logs"
)

// Logger is a logging middleware, tying into the ssbd/logs package for raw log management.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// start timer
		start := time.Now()
		path := c.Request.URL.Path

		// process request
		c.Next()

		// collect data for log
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		errors := c.Errors.String()

		// write log
		switch {
		case status >= 200 && status < 400: // All good signals
			logs.Access(status, method, latency,
				clientIP, path, errors)
		case status >= 400 && status < 600: // Errors
			logs.Error(status, method, latency,
				clientIP, path, errors)
		}
	}
}

package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Logger middleware provides structured request logging using slog
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Use structured logging with slog
		slog.Info("HTTP Request",
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.Int("status", param.StatusCode),
			slog.Duration("latency", param.Latency),
			slog.String("client_ip", param.ClientIP),
			slog.Time("timestamp", param.TimeStamp),
			slog.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}

// LoggerWithLevel creates a logger middleware with configurable log level
func LoggerWithLevel(level slog.Level) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Only log if the current level allows it
		if !slog.Default().Enabled(nil, level) {
			return ""
		}

		// Log at the specified level
		slog.Log(nil, level, "HTTP Request",
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.Int("status", param.StatusCode),
			slog.Duration("latency", param.Latency),
			slog.String("client_ip", param.ClientIP),
			slog.Time("timestamp", param.TimeStamp),
			slog.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}

// ErrorLogger logs errors with structured logging
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred during request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error("Request processing error",
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("error", err.Error()),
					slog.Any("type", err.Type),
					slog.String("client_ip", c.ClientIP()),
				)
			}
		}
	}
}
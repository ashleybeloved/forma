package main

import (
	"forma/internal/config"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	slog.Info("Server running on port " + cfg.ServerPort)
	r.Run(cfg.ServerPort)
}

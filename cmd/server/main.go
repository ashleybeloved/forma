package main

import (
	"forma/internal/config"
	"forma/internal/handler"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Handlers
	pingHandler := handler.NewPingHandler(cfg)

	// Routes
	r.GET("/ping", pingHandler.Handle)

	slog.Info("Forma Server running on port " + cfg.ServerPort)
	r.Run(cfg.ServerPort)
}

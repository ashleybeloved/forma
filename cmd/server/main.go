package main

import (
	"forma/internal/config"
	"forma/internal/handler"
	"forma/internal/middleware"
	"forma/internal/repository"
	"forma/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Connect Database
	db := repository.Connect(cfg.DatabasePath)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Dependency Injection
	pingHandler := handler.NewPingHandler(cfg)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, cfg)

	// - Routes -
	r.GET("/ping", pingHandler.Handle)

	// -- No Auth --
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/logout", userHandler.Logout)

	// -- Need Auth --
	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(cfg))

	slog.Info("Forma Server running on port " + cfg.ServerPort)
	r.Run(cfg.ServerPort)
}

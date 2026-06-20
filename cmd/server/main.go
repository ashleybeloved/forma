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

	pollRepo := repository.NewPollRepository(db)
	pollService := service.NewPollService(pollRepo, cfg)
	pollHandler := handler.NewPollHandler(pollService, cfg)

	// - Routes -
	r.GET("/ping", pingHandler.Handle)

	// -- No Auth --
	r.GET("/poll/:short_id", pollHandler.GetPollByShortID)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/logout", userHandler.Logout)

	// -- Need Auth --
	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(cfg))
	{
		auth.POST("/poll", pollHandler.CreatePoll)    // Create Poll
		auth.PATCH("/poll", pollHandler.UpdatePoll)   // Edit Poll
		auth.DELETE("/poll", pollHandler.DeletePoll)  // Delete Poll
		auth.GET("/polls", pollHandler.GetAllMyPolls) // Get All Profile Polls | Queries LIMIT and OFFSET
	}

	slog.Info("Forma Server running on port " + cfg.ServerPort)
	r.Run(cfg.ServerPort)
}

package main

import (
	"context"
	"forma/internal/config"
	"forma/internal/handler"
	"forma/internal/middleware"
	"forma/internal/repository"
	"forma/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Connect Database
	db := repository.Connect(cfg.DatabasePath)
	defer func() {
		slog.Info("Close connection with database...")
		err := db.Close()
		if err != nil {
			slog.Error("failed to close connection with database", "error", err)
		}
	}()

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Dependency Injection
	pingHandler := handler.NewPingHandler(cfg)
	geoIPService := service.NewGeoIPService(cfg.GeoIPDatabasePath)
	defer func() {
		slog.Info("Close connection with GeoIP Service...")
		geoIPService.Close()
	}()
	validatorService := service.NewValidatorService(cfg)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, validatorService)
	userHandler := handler.NewUserHandler(userService, cfg)

	pollRepo := repository.NewPollRepository(db)
	pollService := service.NewPollService(pollRepo, cfg, geoIPService, validatorService)
	pollHandler := handler.NewPollHandler(pollService, cfg)

	// - Routes -
	// -- No Auth --
	r.GET("/ping", pingHandler.Handle)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/logout", userHandler.Logout)

	// -- Guest Token Required --
	guest := r.Group("")
	guest.Use(middleware.GuestMiddleware(cfg))
	{
		guest.GET("/poll/:short_id", pollHandler.GetPollByShortID) // Get Poll
		guest.POST("/poll/:short_id/vote", pollHandler.Vote)       // Vote in Poll
		guest.POST("/poll/:short_id/check", pollHandler.CheckVote) // Check vote user
	}

	// -- Need Auth --
	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(cfg))
	{
		auth.GET("/me", userHandler.Me)
		auth.POST("/poll", pollHandler.CreatePoll)   // Create Poll
		auth.PATCH("/poll", pollHandler.UpdatePoll)  // Edit Poll
		auth.DELETE("/poll", pollHandler.DeletePoll) // Delete Poll
		auth.GET("/poll", pollHandler.GetAllMyPolls) // Get All Profile Polls | Queries LIMIT & OFFSET

		auth.GET("/poll/:short_id/stats", pollHandler.GetPollStats)
	}

	// -!- Start server & Graceful Shutdown -!-
	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
	}

	go func() {
		slog.Info("Forma Server running on port " + cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Server get signal for Graceful Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}

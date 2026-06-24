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

	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // For Dev
			"https://" + cfg.Domain, // For Prod
			"http://" + cfg.Domain,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	api := r.Group("/api")
	{
		// -- No Auth --
		api.GET("/ping", pingHandler.Handle)
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/logout", userHandler.Logout)

		// -- Guest Token Required --
		guest := api.Group("")
		guest.Use(middleware.GuestMiddleware(cfg))
		{
			guest.GET("/poll/:short_id", pollHandler.GetPollByShortID) // Get Poll
			guest.POST("/poll/:short_id/vote", pollHandler.Vote)       // Vote in Poll
			guest.POST("/poll/:short_id/check", pollHandler.CheckVote) // Check vote user
		}

		// -- Need Auth --
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(cfg))
		{
			auth.GET("/me", userHandler.Me)
			auth.POST("/poll", pollHandler.CreatePoll)            // Create Poll
			auth.PATCH("/poll/:short_id", pollHandler.UpdatePoll) // Edit Poll
			auth.DELETE("/poll/short_id", pollHandler.DeletePoll) // Delete Poll
			auth.GET("/poll", pollHandler.GetAllMyPolls)          // Get All Profile Polls | Queries LIMIT & OFFSET

			auth.GET("/poll/:short_id/stats", pollHandler.GetPollStats) // Statistics
		}
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

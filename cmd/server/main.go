package main

import (
	"context"
	"forma/internal/config"
	"forma/internal/handler"
	"forma/internal/repository"
	"forma/internal/service"
	"forma/router"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// ROUTER
	r := router.Setup(cfg, pingHandler, userHandler, pollHandler)

	// Start server & Graceful Shutdown
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

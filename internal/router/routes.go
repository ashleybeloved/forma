package router

import (
	"forma/internal/config"
	"forma/internal/handler"
	"forma/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config, pingHandler *handler.PingHandler, userHandler *handler.AuthHandler, pollHandler *handler.PollHandler) *gin.Engine {
	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
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
			auth.POST("/poll", pollHandler.CreatePoll)             // Create Poll
			auth.PATCH("/poll/:short_id", pollHandler.UpdatePoll)  // Edit Poll
			auth.DELETE("/poll/:short_id", pollHandler.DeletePoll) // Delete Poll
			auth.GET("/poll", pollHandler.GetAllMyPolls)           // Get All Profile Polls | Queries LIMIT & OFFSET

			auth.GET("/poll/:short_id/stats", pollHandler.GetPollStats) // Statistics
		}
	}

	return r
}

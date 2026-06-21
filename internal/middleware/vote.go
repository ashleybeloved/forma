package middleware

import (
	"forma/internal/config"

	"github.com/gin-gonic/gin"
)

// Middleware for check user vote by IP and Guest Token
func VoteMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

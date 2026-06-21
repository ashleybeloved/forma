package middleware

import (
	"forma/internal/config"

	"github.com/gin-gonic/gin"
)

// Middleware to check and give a any User (with non-auth) a guest jwt-token :-)
func GuestMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

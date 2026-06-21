package middleware

import (
	"forma/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Middleware to check and give to any User (including non-auth) a guest jwt-token :-)
func GuestMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("guest_token")

		var token string

		if err != nil || cookie == "" {
			token := uuid.New().String()

			c.SetCookie("guest_token", token, cfg.JWTTimeToLive*3600, "/", cfg.Domain, false, true)
		} else {
			token = cookie
		}

		c.Set("guest_token", token)
		c.Next()
	}
}

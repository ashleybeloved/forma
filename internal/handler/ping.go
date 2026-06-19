package handler

import (
	"forma/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	Config *config.Config
}

func NewPingHandler(cfg *config.Config) *PingHandler {
	return &PingHandler{
		Config: cfg,
	}
}

func (h *PingHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": h.Config.AppVersion,
	})
}

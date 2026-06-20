package handler

import (
	"forma/internal/config"
	"forma/internal/model"
	"forma/internal/pkg"
	"forma/internal/repository"
	"forma/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
	Config  *config.Config
}

func NewUserHandler(service *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Config:  cfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	req := model.RegisterRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	user, err := h.Service.Register(req.Username, req.Password)
	if err != nil {
		switch err {
		case repository.ErrUsernameAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": repository.ErrUsernameAlreadyExists.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to register user",
			})
		}
		return
	}

	token, err := pkg.GenerateToken(user.ID, h.Config.JWTTimeToLive, h.Config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.SetCookie("forma_token", token, h.Config.JWTTimeToLive*3600, "/", h.Config.Domain, false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	req := model.LoginRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	user, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case repository.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrUserNotFound.Error()})
		case service.ErrInvalidPassword:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrInvalidPassword.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to register user",
			})
		}
		return
	}

	token, err := pkg.GenerateToken(user.ID, h.Config.JWTTimeToLive, h.Config.JWTSecretKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.SetCookie("forma_token", token, h.Config.JWTTimeToLive*3600, "/", h.Config.Domain, false, true)
	c.JSON(http.StatusCreated, gin.H{
		"message": "login into " + user.Username + " successful",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("forma_token", "", -1, "/", h.Config.Domain, false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "logout successful",
	})
}

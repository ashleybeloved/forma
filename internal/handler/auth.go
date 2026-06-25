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

func (h *AuthHandler) Me(c *gin.Context) {
	id, _ := c.Get("user_id")
	if id == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user_id in context",
		})
		return
	}

	user, err := h.Service.Me(id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user information",
		})
		return
	}

	response := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
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
		case service.ErrMustContainLetter:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrMustContainLetter.Error()})
		case service.ErrDoubleUnderscores:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrDoubleUnderscores.Error()})
		case service.ErrUsernameTooBig:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrUsernameTooBig.Error()})
		case service.ErrUsernameTooSmall:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrUsernameTooSmall.Error()})
		case service.ErrPasswordTooBig:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrPasswordTooBig.Error()})
		case service.ErrPasswordTooSmall:
			c.JSON(http.StatusConflict, gin.H{"error": service.ErrPasswordTooSmall.Error()})
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

	c.SetCookie("forma_token", token, h.Config.JWTTimeToLive*3600, "/", h.Config.Domain, h.Config.HTTPS, true)

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
				"error": "failed to login user",
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

	c.SetCookie("forma_token", token, h.Config.JWTTimeToLive*3600, "/", h.Config.Domain, h.Config.HTTPS, true)
	c.JSON(http.StatusCreated, gin.H{
		"message": "login into " + user.Username + " successful",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("forma_token", "", -1, "/", h.Config.Domain, h.Config.HTTPS, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "logout successful",
	})
}

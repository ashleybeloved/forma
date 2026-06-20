package handler

import (
	"forma/internal/config"
	"forma/internal/model"
	"forma/internal/repository"
	"forma/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PollHandler struct {
	Service *service.PollService
	Config  *config.Config
}

func NewPollHandler(service *service.PollService, cfg *config.Config) *PollHandler {
	return &PollHandler{
		Service: service,
		Config:  cfg,
	}
}

func (h *PollHandler) GetPollByShortID(c *gin.Context) {
	shortID := c.Param("short_id")

	poll, err := h.Service.Repo.GetPollByShortID(shortID)
	if err != nil {
		switch err {
		case repository.ErrPollNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": repository.ErrPollNotFound.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get poll"})
		}
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (h *PollHandler) CreatePoll(c *gin.Context) {
	req := model.NewPollRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input, title and config are required fields",
		})
		return
	}

	creatorID, _ := c.Get("user_id")
	if creatorID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user_id in context",
		})
		return
	}

	poll, err := h.Service.CreatePoll(req.Title, req.Description, req.Config, creatorID.(int))
	if err != nil {
		switch err {
		case service.ErrMarshalJSON:
			c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrMarshalJSON.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create poll"})
		}
		return
	}

	c.JSON(http.StatusCreated, poll)
}

func (h *PollHandler) GetAllMyPolls(c *gin.Context) {
	creatorID, _ := c.Get("user_id")
	if creatorID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user_id in context",
		})
		return
	}

	polls, err := h.Service.GetAllMyPolls(creatorID.(int))
	if err != nil {
		switch err {
		case repository.ErrPollsNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrPollsNotFound.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to found polls"})
		}
		return
	}

	c.JSON(http.StatusOK, polls)
}

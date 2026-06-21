package handler

import (
	"forma/internal/config"
	"forma/internal/model"
	"forma/internal/repository"
	"forma/internal/service"
	"net/http"
	"strconv"

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

func (h *PollHandler) UpdatePoll(c *gin.Context) {
	req := &model.UpdatePollRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input, id, title and config are required fields",
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

	err = h.Service.UpdatePoll(req.ID, req.Title, req.Description, req.Config, creatorID.(int))
	if err != nil {
		switch err {
		case repository.ErrPollNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": repository.ErrPollNotFound.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create poll"})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *PollHandler) DeletePoll(c *gin.Context) {
	req := &model.DeletePollRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input, id are required field",
		})
		return
	}

	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user_id in context",
		})
		return
	}

	err = h.Service.DeletePoll(req.ID, userID.(int))
	if err != nil {
		switch err {
		case repository.ErrPollNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": repository.ErrPollNotFound.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete poll"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PollHandler) GetAllMyPolls(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user_id in context",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	polls, err := h.Service.GetAllMyPolls(userID.(int), limit, offset)
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

func (h *PollHandler) Vote(c *gin.Context) {
	req := &model.NewVoteRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input, answers are required field",
		})
		return
	}

	tokenStr, err := c.Cookie("forma_token")
	if err != nil {
		tokenStr = ""
	}

	guestToken, _ := c.Get("guest_token")
	if guestToken == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "missed guest token",
		})
		return
	}

	err = h.Service.Vote(tokenStr, c.Param("short_id"), guestToken.(string), c.ClientIP(), &req.Answers)
	if err != nil {
		switch err {
		case service.ErrMarshalJSON:
			c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrMarshalJSON.Error()})
		case service.ErrAlreadyVoted:
			c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrAlreadyVoted.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to vote"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully voted"})
}

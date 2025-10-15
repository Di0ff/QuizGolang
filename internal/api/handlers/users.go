package handlers

import (
	"net/http"
	"quiz/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateOrGet(c *gin.Context) {
	var req models.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, err, http.StatusBadRequest)
		return
	}

	user, err := h.S.Find(req.TelegramID)
	if err == nil {
		c.JSON(http.StatusOK, user)
		return
	}

	newUser, err := h.S.CreateUser(req.TelegramID, req.Username, req.Name)
	if err != nil {
		h.Error(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (h *Handler) UpdateStreak(c *gin.Context) {
	var req models.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, err, http.StatusBadRequest)
		return
	}

	user, err := h.S.Find(req.TelegramID)
	if err != nil {
		h.Error(c, err, http.StatusNotFound)
		return
	}

	if err := h.S.UpdateStreak(user); err != nil {
		h.Error(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

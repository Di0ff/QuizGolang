package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLeaderboard(c *gin.Context) {
	leaders, err := h.S.Top(10)
	if err != nil {
		h.Error(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, leaders)
}

func (h *Handler) GetUser(c *gin.Context) {
	temp := c.Query("user_id")
	userID, err := strconv.Atoi(temp)
	if err != nil {
		h.Error(c, err, http.StatusBadRequest)
		return
	}

	user, err := h.S.GetUserStats(userID)
	if err != nil {
		h.Error(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetLeaderboardWithUser(c *gin.Context) {
	temp := c.Query("user_id")
	userID, err := strconv.Atoi(temp)
	if err != nil {
		h.Error(c, err, http.StatusBadRequest)
		return
	}

	data, err := h.S.GetLeaderboardWithUser(10, userID)
	if err != nil {
		h.Error(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

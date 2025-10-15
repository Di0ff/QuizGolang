package handlers

import (
	"net/http"
	"quiz/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetQuestions(c *gin.Context) {
	temp := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(temp)
	if err != nil || limit <= 0 {
		limit = 3
	}

	questions, err := h.S.GetQuestions(limit)
	if err != nil {
		h.Error(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *Handler) Save(c *gin.Context) {
	var req models.RequestQuiz
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, err, http.StatusBadRequest)
		return
	}

	if err := h.S.Save(req.UserID, req.Score, req.TotalQuestions); err != nil {
		h.Error(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Результат сохранён"})
}

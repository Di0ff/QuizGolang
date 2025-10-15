package handlers

import (
	"quiz/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	S *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{S: s}
}

func (h *Handler) Error(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

package repository

import (
	"quiz/internal/models"

	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &Database{DB: db}
}

type Repository interface {
	CreateUser(user *models.Users) error
	GetUser(tgID int64) (*models.Users, error)
	UpdateUser(user *models.Users) error

	GetRandomQuestions(limit int) ([]models.Questions, error)
	CreateQuestion(q *models.Questions) error

	AddResult(l *models.Leaderboards) error
	GetLeaderboards(limit int) ([]models.Leaderboards, error)
	GetUserLeaderboard(userID int) ([]models.Leaderboards, error)
}

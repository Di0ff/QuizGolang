package service

import (
	"quiz/internal/database/repository"
	"quiz/internal/logger"
)

type Service struct {
	repo repository.Repository
	log  *logger.Zap
}

func New(repo repository.Repository, log *logger.Zap) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

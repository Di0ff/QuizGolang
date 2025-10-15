package service

import (
	"quiz/internal/models"

	"go.uber.org/zap"
)

func (s *Service) Top(limit int) (leaderboards []models.Leaderboards, err error) {
	leaderboards, err = s.repo.GetLeaderboards(limit)
	if err != nil {
		s.log.Error("Ошибка получения списка лидеров", zap.Error(err))
		return nil, err
	}

	return leaderboards, nil
}

func (s *Service) GetUserStats(userID int) (stats []models.Leaderboards, err error) {
	stats, err = s.repo.GetUserLeaderboard(userID)
	if err != nil {
		s.log.Error("Ошибка получения статистики", zap.Int("user_id", userID), zap.Error(err))
		return nil, err
	}

	return stats, nil
}

func (s *Service) GetLeaderboardWithUser(limit int, userID int) (*models.LeaderboardWithUser, error) {
	leaderboards, err := s.Top(limit)
	if err != nil {
		return nil, err
	}

	stats, err := s.GetUserStats(userID)
	if err != nil {
		return nil, err
	}

	return &models.LeaderboardWithUser{
		Leaderboard: leaderboards,
		UserStats:   stats,
	}, nil
}

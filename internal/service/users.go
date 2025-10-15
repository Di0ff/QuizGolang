package service

import (
	"time"

	"quiz/internal/models"

	"go.uber.org/zap"
)

func (s *Service) Find(tgID int64) (user *models.Users, err error) {
	user, err = s.repo.GetUser(tgID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(tgID int64, username, name string) (new *models.Users, err error) {
	new = &models.Users{
		TelegramID: tgID,
		Username:   username,
		Name:       name,
	}

	if err = s.repo.CreateUser(new); err != nil {
		s.log.Error("Ошибка! Пользователь не создан", zap.Int64("telegram_id", tgID), zap.Error(err))
		return nil, err
	}

	return new, nil
}

func (s *Service) CreateOrGetUser(tgID int64, username, name string) (*models.QuizUser, error) {
	user, err := s.Find(tgID)
	if err == nil {
		return &models.QuizUser{
			ID:         user.ID,
			TelegramID: user.TelegramID,
			Username:   user.Username,
			Name:       user.Name,
		}, nil
	}

	newUser, err := s.CreateUser(tgID, username, name)
	if err != nil {
		return nil, err
	}

	return &models.QuizUser{
		ID:         newUser.ID,
		TelegramID: newUser.TelegramID,
		Username:   newUser.Username,
		Name:       newUser.Name,
	}, nil
}

func (s *Service) UpdateStreak(user *models.Users) (err error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	if user.LastActivity != nil {
		last := user.LastActivity.Truncate(24 * time.Hour)

		switch {
		case last.Equal(today):
			return nil
		case last.Equal(today.AddDate(0, 0, -1)):
			user.Streak++
		default:
			user.Streak = 1
		}
	} else {
		user.Streak = 1
	}

	user.LastActivity = &today

	if err = s.repo.UpdateUser(user); err != nil {
		s.log.Error("Ошибка обновления стрика", zap.Int("user_id", user.ID), zap.Error(err))
		return err
	}

	return nil
}

package service

import (
	"go.uber.org/zap"

	"quiz/internal/models"
)

func (s *Service) GetQuestions(limit int) (questions []models.Questions, err error) {
	questions, err = s.repo.GetRandomQuestions(limit)
	if err != nil {
		s.log.Error("Ошибка получения вопросов", zap.Error(err))
		return nil, err
	}

	return questions, nil
}

func (s *Service) GetRandomQuestions(limit int) ([]models.QuizQuestion, error) {
	questions, err := s.repo.GetRandomQuestions(limit)
	if err != nil {
		s.log.Error("Ошибка получения вопросов", zap.Error(err))
		return nil, err
	}

	result := make([]models.QuizQuestion, len(questions))
	for i, q := range questions {
		result[i] = models.QuizQuestion{
			ID:            q.ID,
			Question:      q.Question,
			Options:       q.Options,
			CorrectOption: q.CorrectOption,
		}
	}

	return result, nil
}

func (s *Service) SaveResult(userID int, score, total int) error {
	result := &models.Leaderboards{
		UserID:         userID,
		Score:          score,
		TotalQuestions: total,
	}

	if err := s.repo.AddResult(result); err != nil {
		s.log.Error("Ошибка сохранения результата", zap.Int("user_id", userID), zap.Error(err))
		return err
	}

	return nil
}

func (s *Service) Save(userID int, score, total int) (err error) {
	result := &models.Leaderboards{
		UserID:         userID,
		Score:          score,
		TotalQuestions: total,
	}

	if err = s.repo.AddResult(result); err != nil {
		s.log.Error("Ошибка сохранения результата", zap.Int("user_id", userID), zap.Error(err))
		return err
	}

	return nil
}

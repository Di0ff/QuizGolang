package repository

import "quiz/internal/models"

func (d *Database) CreateQuestion(m *models.Questions) error {
	return d.DB.Create(m).Error
}

func (d *Database) GetRandomQuestions(limit int) (questions []models.Questions, err error) {
	return questions, d.DB.Order("RANDOM()").Limit(limit).Find(&questions).Error
}

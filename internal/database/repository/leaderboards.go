package repository

import "quiz/internal/models"

func (d *Database) AddResult(m *models.Leaderboards) error {
	var existing models.Leaderboards
	err := d.DB.Where("user_id = ?", m.UserID).First(&existing).Error

	if err != nil {
		return d.DB.Create(m).Error
	}

	existing.Score += m.Score
	existing.TotalQuestions += m.TotalQuestions
	return d.DB.Save(&existing).Error
}

func (d *Database) GetLeaderboards(limit int) (leaderboards []models.Leaderboards, err error) {
	return leaderboards, d.DB.
		Preload("User").
		Order("score DESC, created_at ASC").
		Limit(limit).
		Find(&leaderboards).Error
}

func (d *Database) GetUserLeaderboard(userID int) (leaderboards []models.Leaderboards, err error) {
	return leaderboards, d.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&leaderboards).Error
}

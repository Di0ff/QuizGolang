package repository

import (
	"quiz/internal/models"
)

func (d *Database) CreateUser(m *models.Users) error {
	return d.DB.Create(m).Error
}

func (d *Database) GetUser(tgID int64) (m *models.Users, err error) {
	return m, d.DB.Where("telegram_id = ?", tgID).First(&m).Error
}

func (d *Database) UpdateUser(m *models.Users) error {
	return d.DB.Save(m).Error
}

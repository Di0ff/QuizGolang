package models

import "time"

type Users struct {
	ID           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	TelegramID   int64      `json:"telegram_id" gorm:"uniqueIndex;not null" binding:"required"`
	Username     string     `json:"username" gorm:"type:text"`
	Name         string     `json:"name" gorm:"type:text"`
	LastActivity *time.Time `json:"last_activity,omitempty" gorm:"type:date"`
	Streak       int        `json:"streak" gorm:"default:0"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

type Questions struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Question      string `json:"question" gorm:"type:text;not null" binding:"required"`
	Options       string `json:"options" gorm:"type:jsonb;not null" binding:"required"`
	CorrectOption int    `json:"correct_option" gorm:"not null" binding:"required"`
	Difficulty    int16  `json:"difficulty,omitempty" gorm:"type:smallint"`
	Topic         string `json:"topic,omitempty" gorm:"type:text"`
}

type Leaderboards struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"not null;index" binding:"required"`
	User           Users     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Score          int       `json:"score" gorm:"not null" binding:"required,gte=0"`
	TotalQuestions int       `json:"total_questions" gorm:"not null" binding:"required,gte=1"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// -------------------------------------------------------------------------------------------------

type LeaderboardWithUser struct {
	Leaderboard []Leaderboards `json:"leaderboard"`
	UserStats   []Leaderboards `json:"user_stats"`
}

type Request struct {
	TelegramID int64  `json:"telegram_id" binding:"required"`
	Username   string `json:"username"`
	Name       string `json:"name"`
}

type RequestQuiz struct {
	UserID         int `json:"user_id" binding:"required"`
	Score          int `json:"score" binding:"required,gte=0"`
	TotalQuestions int `json:"total_questions" binding:"required,gte=1"`
}

// -------------------------------------------------------------------------------------------------

type QuizQuestion struct {
	ID            int    `json:"id"`
	Question      string `json:"question"`
	Options       string `json:"options"`
	CorrectOption int    `json:"correct_option"`
}

type QuizUser struct {
	ID         int    `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
}

type QuizSession struct {
	Questions    []QuizQuestion
	CurrentIndex int
	CorrectCount int
	UserID       int
}

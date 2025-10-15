package routes

import (
	"quiz/internal/api/handlers"
	"quiz/internal/config"
	"quiz/internal/database/repository"
	"quiz/internal/logger"
	"quiz/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(db *gorm.DB, log *logger.Zap, cfg *config.Cfg) *gin.Engine {
	r := gin.Default()

	repo := repository.New(db)
	ser := service.New(repo, log)
	hand := handlers.New(ser)

	miniapp := r.Group("/miniapp")
	{
		miniapp.GET("/leaderboard", hand.GetLeaderboard)
		miniapp.GET("/leaderboard/with-user", hand.GetLeaderboardWithUser)
	}

	tg := r.Group("/tg")
	{
		tg.POST("/users", hand.CreateOrGet)
		tg.POST("/users/streak", hand.UpdateStreak)
		tg.GET("/questions", hand.GetQuestions)
		tg.POST("/results", hand.Save)
	}

	r.Static("/static", cfg.App.StaticDir)
	r.StaticFile("/logo.png", cfg.App.FrontendDir+"/logo.png")
	r.StaticFile("/leaderboard", cfg.App.FrontendDir+"/leaderboard.html")
	r.StaticFile("/", cfg.App.FrontendDir+"/leaderboard.html")

	return r
}

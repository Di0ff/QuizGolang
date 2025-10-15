package handlers

import (
	"quiz/internal/config"
	"quiz/internal/database/redis"
	"quiz/internal/logger"
	"quiz/internal/service"

	"gopkg.in/telebot.v4"
)

func Start(bot *telebot.Bot, menu *telebot.ReplyMarkup, cfg *config.Cfg, log *logger.Zap, redis *redis.Redis, srv *service.Service) {
	bot.Handle("/start", func(c telebot.Context) error {
		text := "👋 Привет! Я помогу тебе подготовиться к собесам на Golang\n\nВыбери действие:"
		return c.Send(text, menu)
	})

	bot.Handle("🎯 Пройти квест", StartQuiz(srv, redis, cfg))
	bot.Handle(telebot.OnCallback, HandleAnswer(redis, srv))

	bot.Handle("🤖 Совет от ИИ", func(c telebot.Context) error {
		return HandleAI(c, cfg, log)
	})

	bot.Handle("ℹ️ О боте", func(c telebot.Context) error {
		return c.Send("Этот бот создан для подготовки к Go собесам\nby https://github.com/Di0ff")
	})
}

package bot

import (
	"quiz/internal/bot/handlers"
	"quiz/internal/bot/menu"
	"quiz/internal/config"
	"quiz/internal/database/redis"
	"quiz/internal/logger"
	"quiz/internal/service"
	"time"

	"gopkg.in/telebot.v4"
)

func Start(token, url string, log *logger.Zap, cfg *config.Cfg, redis *redis.Redis, srv *service.Service) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Error("Ошибка запуска тг бота")
	}

	menu := menu.Init(url)
	handlers.Start(bot, menu, cfg, log, redis, srv)

	log.Info("Telegram бот запущен")
	bot.Start()
}

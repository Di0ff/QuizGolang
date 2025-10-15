package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"quiz/internal/api/routes"
	"quiz/internal/bot"
	"quiz/internal/config"
	"quiz/internal/database/migrations"
	"quiz/internal/database/postgres"
	"quiz/internal/database/redis"
	"quiz/internal/database/repository"
	"quiz/internal/graceful"
	"quiz/internal/logger"
	"quiz/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	logger, err := logger.New(cfg.Logger.Env, cfg.Logger.Level)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer logger.Sync()

	db, err := postgres.New(cfg, logger)
	if err != nil {
		logger.Fatal("Ошибка подключения к БД")
	}
	defer db.Close(logger)

	if err := migrations.Run(cfg, logger); err != nil {
		logger.Error("Ошибка выполнения миграций")
		os.Exit(1)
	}

	redisStore, err := redis.New(&cfg.Redis, logger)
	if err != nil {
		logger.Fatal("Ошибка подключения к Redis: " + err.Error())
	}
	defer redisStore.Close()

	repo := repository.New(db.DB)
	service := service.New(repo, logger)

	go func() {
		logger.Info("Запуск Telegram бота")
		bot.Start(cfg.Telegram.Token, cfg.Telegram.URL, logger, cfg, redisStore, service)
	}()

	logger.Info("Запуск сервера")
	router := routes.Init(db.DB, logger, cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServeTLS(
			"/etc/letsencrypt/live/quizgolang.ru/fullchain.pem",
			"/etc/letsencrypt/live/quizgolang.ru/privkey.pem",
		); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка запуска сервера")
		}
	}()

	graceful.Shutdown(srv, logger, 30*time.Second)
}

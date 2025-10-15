package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quiz/internal/logger"

	"go.uber.org/zap"
)

func Shutdown(srv *http.Server, log *logger.Zap, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Получен сигнал, останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Ошибка при остановке сервера", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Сервер корректно остановлен")
}

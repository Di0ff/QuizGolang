package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"quiz/internal/config"
	"quiz/internal/logger"
	"quiz/internal/models"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb        *redis.Client
	sessionTTL time.Duration
	log        *logger.Zap
}

func New(cfg *config.ConfigRedis, log *logger.Zap) (*Redis, error) {
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		log.Error("Неверный формат DB: " + cfg.DB)
		return nil, err
	}

	ttl, err := strconv.Atoi(cfg.SessionTTL)
	if err != nil {
		log.Error("Неверный формат TTL: " + cfg.SessionTTL)
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Ошибка подключения к Redis: " + err.Error())
		return nil, err
	}

	sessionTTL := time.Duration(ttl) * time.Second
	if sessionTTL == 0 {
		sessionTTL = time.Hour
	}

	return &Redis{
		rdb:        rdb,
		sessionTTL: sessionTTL,
		log:        log,
	}, nil
}

func (s *Redis) Set(userID int64, session *models.QuizSession) error {
	ctx := context.Background()

	data, err := json.Marshal(session)
	if err != nil {
		s.log.Error("Ошибка маршалинга сессии " + strconv.FormatInt(userID, 10))
		return err
	}

	key := "quiz_session:" + strconv.FormatInt(userID, 10)
	if err := s.rdb.Set(ctx, key, data, s.sessionTTL).Err(); err != nil {
		s.log.Error("Ошибка сохранения сессии " + strconv.FormatInt(userID, 10))
		return err
	}

	return nil
}

func (s *Redis) Get(userID int64) (*models.QuizSession, error) {
	ctx := context.Background()

	key := "quiz_session:" + strconv.FormatInt(userID, 10)
	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		s.log.Error("Ошибка получения сессии " + strconv.FormatInt(userID, 10))
		return nil, err
	}

	var session models.QuizSession
	if err = json.Unmarshal([]byte(val), &session); err != nil {
		s.log.Error("Ошибка анмаршалинга сессии " + strconv.FormatInt(userID, 10))
		return nil, err
	}

	return &session, nil
}

func (s *Redis) Delete(userID int64) error {
	ctx := context.Background()

	key := "quiz_session:" + strconv.FormatInt(userID, 10)
	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		s.log.Error("Ошибка удаления сессии " + strconv.FormatInt(userID, 10))
		return err
	}
	return nil
}

func (s *Redis) Close() error {
	return s.rdb.Close()
}

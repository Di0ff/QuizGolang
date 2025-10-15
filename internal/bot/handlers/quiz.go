package handlers

import (
	"encoding/json"
	"fmt"
	"quiz/internal/config"
	"quiz/internal/database/redis"
	"quiz/internal/models"
	"quiz/internal/service"
	"strconv"
	"strings"

	"gopkg.in/telebot.v4"
)

func StartQuiz(srv *service.Service, redisStore *redis.Redis, cfg *config.Cfg) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		name := strings.TrimSpace(c.Sender().FirstName + " " + c.Sender().LastName)
		user, err := srv.CreateOrGetUser(c.Sender().ID, c.Sender().Username, name)
		if err != nil {
			return c.Send("Ошибка! пользователь не найден")
		}

		limit, _ := strconv.Atoi(cfg.Quiz.QuestionsLimit)
		if limit == 0 {
			limit = 3
		}

		questions, err := srv.GetRandomQuestions(limit)
		if err != nil {
			return c.Send("Ошибка! не удалось загрузить вопросы")
		}

		session := &models.QuizSession{
			Questions:    questions,
			CurrentIndex: 0,
			CorrectCount: 0,
			UserID:       user.ID,
		}

		if err := redisStore.Set(c.Sender().ID, session); err != nil {
			return c.Send("Ошибка! создания сессии")
		}

		return showQuestion(c, session, redisStore, srv)
	}
}

func showQuestion(c telebot.Context, session *models.QuizSession, redisStore *redis.Redis, srv *service.Service) error {
	if session.CurrentIndex >= len(session.Questions) {
		return finishQuiz(c, session, redisStore, srv)
	}

	q := session.Questions[session.CurrentIndex]

	var options []string
	if err := json.Unmarshal([]byte(q.Options), &options); err != nil {
		return c.Send("Ошибка! парсинга вопроса")
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	for i, opt := range options {
		btn := markup.Data(opt, fmt.Sprintf("ans:%d:%d", session.CurrentIndex, i))
		rows = append(rows, markup.Row(btn))
	}
	markup.Inline(rows...)

	text := fmt.Sprintf("Вопрос %d из %d:\n\n%s",
		session.CurrentIndex+1,
		len(session.Questions),
		q.Question)

	return c.Send(text, markup)
}

func HandleAnswer(redisStore *redis.Redis, srv *service.Service) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		data := c.Callback().Data
		userID := c.Sender().ID

		session, err := redisStore.Get(userID)
		if err != nil || session == nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Сессия истекла или не найдена"})
		}

		parts := strings.Split(data, ":")
		if len(parts) != 3 {
			return c.Respond(&telebot.CallbackResponse{Text: "Ошибка данных"})
		}

		questionIndex, _ := strconv.Atoi(parts[1])
		answerIndex, _ := strconv.Atoi(parts[2])

		q := session.Questions[questionIndex]
		if answerIndex == q.CorrectOption {
			session.CorrectCount++
			c.Respond(&telebot.CallbackResponse{Text: "Правильно!"})
		} else {
			c.Respond(&telebot.CallbackResponse{Text: "Неправильно"})
		}

		session.CurrentIndex++

		if err := redisStore.Set(userID, session); err != nil {
			return c.Send("Ошибка! не удалось обновить сессию")
		}

		return showQuestion(c, session, redisStore, srv)
	}
}

func finishQuiz(c telebot.Context, session *models.QuizSession, redisStore *redis.Redis, srv *service.Service) error {
	err := srv.SaveResult(session.UserID, session.CorrectCount, len(session.Questions))
	if err != nil {
		return c.Send("Ошибка! сохранения результата")
	}

	redisStore.Delete(c.Sender().ID)

	percentage := (session.CorrectCount * 100) / len(session.Questions)
	emoji := "😢"
	if percentage >= 80 {
		emoji = "🎉"
	} else if percentage >= 50 {
		emoji = "😊"
	}

	text := fmt.Sprintf("%s Квест завершён!\n\nПравильных ответов: %d из %d",
		emoji,
		session.CorrectCount,
		len(session.Questions))

	return c.Send(text)
}

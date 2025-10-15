package handlers

import (
	"context"
	"quiz/internal/config"
	"quiz/internal/logger"
	"strconv"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

func HandleAI(c telebot.Context, cfg *config.Cfg, log *logger.Zap) error {
	question := c.Text()

	c.Notify(telebot.Typing)

	apiKey := cfg.OpenAI.KeyAI
	if apiKey == "" {
		log.Error("Не найден OpenAI ключ в конфиге")
		return c.Send("Ошибка! Сервис не отвечает")
	}

	maxTokens, _ := strconv.Atoi(cfg.OpenAI.MaxTokens)
	if maxTokens == 0 {
		maxTokens = 500
	}

	answer, err := askChatGPT(apiKey, cfg.OpenAI.Model, question, maxTokens, log)
	if err != nil {
		log.Error("Ошибка при запросе к ChatGPT", zap.Error(err))
		return c.Send("Ошибка! Попробуй позже")
	}

	return c.Send(answer)
}

func askChatGPT(apiKey, model, question string, maxTokens int, log *logger.Zap) (string, error) {
	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Ты опытный Go-разработчик. Отвечай кратко и по делу."},
				{Role: openai.ChatMessageRoleUser, Content: question},
			},
			MaxTokens: maxTokens,
		},
	)
	if err != nil {
		log.Error("Ошибка OpenAI запроса", zap.Error(err))
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

package menu

import "gopkg.in/telebot.v4"

func Init(url string) *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}

	quiz := menu.Text("🎯 Пройти квест")
	ai := menu.Text("🤖 Совет от ИИ")
	leaderboard := menu.WebApp("📊 Таблица лидеров", &telebot.WebApp{URL: url})
	about := menu.Text("ℹ️ О боте")

	menu.Reply(
		menu.Row(quiz, ai),
		menu.Row(leaderboard, about),
	)

	return menu
}

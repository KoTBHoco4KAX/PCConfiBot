package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Критическая ошибка: переменная окружения TELEGRAM_BOT_TOKEN не задана")
	}

	initDatabase()

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}

	// Сессии живут 2 часа после последнего действия
	sessionMgr := NewSessionManager(2 * time.Hour)
	handlers := NewBotHandlers(b, sessionMgr)
	handlers.RegisterHandlers()

	log.Println("Бот-конфигуратор успешно запущен...")
	b.Start()
}

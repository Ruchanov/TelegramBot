package main

import (
	"github.com/Ruchanov/TelegramBot/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"sync"
)

var (
	bot     *tgbotapi.BotAPI
	wg      sync.WaitGroup
	mu      sync.Mutex
	counter int
)

func incrementCounter() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		counter++
	}()
}

func main() {
	var err error
	const telegramToken = "6110459357:AAG9b7-9REqy0k_-lCkB_GTR1uzlVJpb_is"
	bot, err = tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {
			if update.Message.IsCommand() || strings.ToLower(update.Message.Text) == "image" {
				incrementCounter()
				image, err := service.FetchRandomImage()
				if err != nil {
					log.Println(err)
					continue
				}

				photo := tgbotapi.NewPhotoShare(update.Message.Chat.ID, image.URLs.Regular)
				photo.Caption = image.Description
				_, err = bot.Send(photo)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

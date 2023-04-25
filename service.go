package main

import (
	"context"
	"fmt"
	"github.com/Ruchanov/TelegramBot/unsplash"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"sync"
)

type Service struct {
	unsplash unsplash.Service
}

func NewService(unsplash unsplash.Service) *Service {
	return &Service{
		unsplash: unsplash,
	}
}

func (s *Service) GetUpdates(ctx context.Context, wg *sync.WaitGroup, token string) {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)

	go func() {
		<-ctx.Done()
		fmt.Println("stopping getUpdates")
		wg.Done()
		return
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	unsplashService := unsplash.NewService(os.Getenv("UNSPLASH_ACCESS_KEY"))

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "image" || update.Message.Command() == "image" {
			wg.Add(1)
			go func() {
				defer wg.Done()

				photo, _ := unsplashService.GetRandomPhoto()

				file := tgbotapi.FileURL(photo.Urls.Regular)
				file.NeedsUpload()

				photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, nil)
				photoMsg.File = file
				photoMsg.Caption = photo.Description

				bot.Send(photoMsg)
			}()
		}
	}

	wg.Wait()
}

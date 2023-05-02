package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	unsplashAPIURL    = "https://api.unsplash.com/photos/random"
	unsplashAccessKey = "un_Gg3hwwtJ69cIC7ntFo4tPIwOt_KR8vcWzz2sPWQQ"
	telegramToken     = "6110459357:AAG9b7-9REqy0k_-lCkB_GTR1uzlVJpb_is"
)

var (
	bot     *tgbotapi.BotAPI
	wg      sync.WaitGroup
	mu      sync.Mutex
	counter int
)

func main() {
	var err error

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
				image, err := fetchRandomImage()
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

func incrementCounter() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		counter++
	}()
}

type unsplashImage struct {
	Description string `json:"description"`
	URLs        struct {
		Regular string `json:"regular"`
	} `json:"urls"`
}

func fetchRandomImage() (*unsplashImage, error) {
	req, err := http.NewRequest("GET", unsplashAPIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Version", "v1")
	req.Header.Set("Authorization", "Client-ID "+unsplashAccessKey)

	client := http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var image unsplashImage
	err = json.NewDecoder(res.Body).Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

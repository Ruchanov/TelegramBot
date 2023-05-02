package service

import (
	"encoding/json"
	"github.com/Ruchanov/TelegramBot/models"
	"net/http"
	"time"
)

const (
	unsplashAPIURL    = "https://api.unsplash.com/photos/random"
	unsplashAccessKey = "un_Gg3hwwtJ69cIC7ntFo4tPIwOt_KR8vcWzz2sPWQQ"
)

func FetchRandomImage() (*models.UnsplashImage, error) {
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

	var image models.UnsplashImage
	err = json.NewDecoder(res.Body).Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

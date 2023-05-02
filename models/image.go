package models

type UnsplashImage struct {
	Description string `json:"description"`
	URLs        struct {
		Regular string `json:"regular"`
	} `json:"urls"`
}

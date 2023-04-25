package models

type Photo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Urls        struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
	} `json:"urls"`
}

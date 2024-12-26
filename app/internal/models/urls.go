package models

type URL struct {
	ShortURL string `json:"short_url" validate:"required,min=8,max=8"`
	LongURL  string `json:"long_url" validate:"required,url"`
}

package movie

import (
	"net/http"
	"time"
)

type MovieStorage struct {
	apiKey string
	apiURL string
	client *http.Client
}

func NewMovieStorage(apiKey, apiURL string) *MovieStorage {
	return &MovieStorage{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

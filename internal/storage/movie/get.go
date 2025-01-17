package movie

import (
	"context"
	"encoding/json"
	"fmt"
	models "match/internal/models/movie"
	"net/http"
)

func (m *MovieStorage) GetMoviesByGenre(ctx context.Context, genre string, page int) ([]models.Movie, error) {
	reqURL := fmt.Sprintf("%s/?apikey=%s&s=%s&page=%d", m.apiURL, m.apiKey, genre, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from Movie API: %d", resp.StatusCode)
	}

	var result struct {
		Search   []models.Movie `json:"Search"`
		Response string         `json:"Response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}
	if result.Response == "False" {
		return []models.Movie{}, nil
	}
	return result.Search, nil
}

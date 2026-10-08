package tmdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrMissingAPIKey = errors.New("TMDB API key is not configured")

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type SearchResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type Movie struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Overview      string `json:"overview"`
	PosterPath    string `json:"poster_path"`
	ReleaseDate   string `json:"release_date"`
	OriginalTitle string `json:"original_title"`
}

func (c *Client) SearchMovies(ctx context.Context, query string) (SearchResponse, error) {
	if c.apiKey == "" {
		return SearchResponse{}, ErrMissingAPIKey
	}

	endpoint, err := url.Parse("https://api.themoviedb.org/3/search/movie")
	if err != nil {
		return SearchResponse{}, fmt.Errorf("parse TMDB endpoint: %w", err)
	}
	params := endpoint.Query()
	params.Set("api_key", c.apiKey)
	params.Set("query", query)
	params.Set("include_adult", "false")
	endpoint.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("create TMDB request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("send TMDB request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return SearchResponse{}, fmt.Errorf("TMDB returned status %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 5<<20)).Decode(&result); err != nil {
		return SearchResponse{}, fmt.Errorf("decode TMDB response: %w", err)
	}
	return result, nil
}
